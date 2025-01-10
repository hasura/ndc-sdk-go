package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"io"
	"os"
	"path"
	"runtime/trace"
	"slices"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"golang.org/x/tools/go/packages"
)

// ConnectorGenerationArguments represent input arguments of the ConnectorGenerator.
type ConnectorGenerationArguments struct {
	Path         string   `help:"The path of the root directory where the go.mod file is present" short:"p" env:"HASURA_PLUGIN_CONNECTOR_CONTEXT_PATH" default:"."`
	ConnectorDir string   `help:"The directory where the connector.go file is placed" default:"."`
	Directories  []string `help:"Folders contain NDC operation functions" short:"d"`
	Trace        string   `help:"Enable tracing and write to target file path."`
	SchemaFormat string   `help:"The output format for the connector schema. Accept: json, go" enum:"json,go" default:"json"`
	Style        string   `help:"The naming style for functions and procedures. Accept: camel-case, snake-case" enum:"camel-case,snake-case" default:"camel-case"`
	TypeOnly     bool     `help:"Generate type only" default:"false"`
}

type connectorTypeBuilder struct {
	packageName string
	packagePath string
	imports     map[string]string
	builder     *strings.Builder
	functions   []FunctionInfo
	procedures  []ProcedureInfo
}

// SetImport sets an import package into the import list.
func (ctb *connectorTypeBuilder) SetImport(value string, alias string) {
	ctb.imports[value] = alias
}

// String renders generated Go types and methods.
func (ctb connectorTypeBuilder) String() string {
	var bs strings.Builder
	writeFileHeader(&bs, ctb.packageName)
	if len(ctb.imports) > 0 {
		bs.WriteString("import (\n")
		sortedImports := utils.GetSortedKeys(ctb.imports)
		for _, pkg := range sortedImports {
			alias := ctb.imports[pkg]
			if alias != "" {
				alias += " "
			}
			bs.WriteString(fmt.Sprintf("  %s\"%s\"\n", alias, pkg))
		}
		bs.WriteString(")\n")
	}

	bs.WriteString(ctb.builder.String())
	return bs.String()
}

type connectorGenerator struct {
	basePath          string
	connectorDir      string
	moduleName        string
	schemaFormat      string
	rawSchema         *RawConnectorSchema
	typeBuilders      map[string]*connectorTypeBuilder
	typeOnly          bool
	functionHandlers  []string
	procedureHandlers []string
}

// ParseAndGenerateConnector parses and generate connector codes.
func ParseAndGenerateConnector(args ConnectorGenerationArguments, moduleName string) error {
	if args.Trace != "" {
		w, err := os.Create(args.Trace)
		if err != nil {
			return fmt.Errorf("failed to create trace file at %s", args.Trace)
		}

		defer func() {
			_ = w.Close()
		}()

		if err = trace.Start(w); err != nil {
			return fmt.Errorf("failed to start trace: %w", err)
		}
		defer trace.Stop()
	}

	parseCtx, parseTask := trace.NewTask(context.TODO(), "parse")
	sm, err := parseRawConnectorSchemaFromGoCode(parseCtx, moduleName, ".", &args)
	if err != nil {
		parseTask.End()
		return err
	}
	parseTask.End()

	_, genTask := trace.NewTask(context.TODO(), "generate_code")
	defer genTask.End()

	connectorGen := &connectorGenerator{
		basePath:     ".",
		connectorDir: args.ConnectorDir,
		moduleName:   moduleName,
		schemaFormat: args.SchemaFormat,
		rawSchema:    sm,
		typeBuilders: make(map[string]*connectorTypeBuilder),
		typeOnly:     args.TypeOnly,
	}

	connectorPkgName, err := connectorGen.loadConnectorPackage()
	if err != nil {
		return err
	}

	return connectorGen.generateConnector(connectorPkgName)
}

func (cg *connectorGenerator) loadConnectorPackage() (string, error) {
	connectorPath := path.Join(cg.basePath, cg.connectorDir)
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Dir:  connectorPath,
		Fset: fset,
	}
	pkgList, err := packages.Load(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to load the package in connector directory: %w", err)
	}

	if len(pkgList) == 0 || pkgList[0].Name == "" {
		return "main", nil
	}
	return pkgList[0].Name, nil
}

func (cg *connectorGenerator) generateConnector(name string) error {
	var schemaBytes []byte
	var err error

	schemaOutputFile := schemaOutputJSONFile
	if cg.schemaFormat == "go" {
		schemaOutputFile = schemaOutputGoFile
		cg.rawSchema.Imports["encoding/json"] = true
		output, err := cg.rawSchema.WriteGoSchema(name)
		if err != nil {
			return err
		}
		schemaBytes = []byte(output)
	} else {
		// generate schema.generated.json
		schemaBytes, err = json.MarshalIndent(cg.rawSchema.Schema(), "", "  ")
		if err != nil {
			return err
		}
	}

	schemaPath := path.Join(cg.basePath, cg.connectorDir, schemaOutputFile)
	if err := os.WriteFile(schemaPath, schemaBytes, 0o644); err != nil {
		return err
	}

	if err := cg.genTypeMethods(); err != nil {
		return err
	}

	if !cg.typeOnly {
		targetPath := path.Join(cg.basePath, cg.connectorDir, connectorOutputFile)
		f, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		w := bufio.NewWriter(f)
		defer func() {
			_ = w.Flush()
		}()

		if err := cg.genConnectorCodeFromTemplate(w, name); err != nil {
			return err
		}
		_ = w.Flush()
	}

	if len(cg.rawSchema.Functions) == 0 && len(cg.rawSchema.Procedures) == 0 {
		log.Warn().Msg("neither function nor procedure is generated. If your project uses Go Workspace please add the root path to the go.work file and update again")
	}

	return nil
}

func (cg *connectorGenerator) genConnectorCodeFromTemplate(w io.Writer, packageName string) error {
	sortedImports := utils.GetSortedKeys(cg.rawSchema.Imports)

	importLines := []string{}
	for _, importPath := range sortedImports {
		importLines = append(importLines, fmt.Sprintf(`"%s"`, importPath))
	}

	return connectorTemplate.Execute(w, map[string]any{
		"Imports":          strings.Join(importLines, "\n"),
		"PackageName":      packageName,
		"StateArgument":    cg.rawSchema.StateType.GetArgumentName(""),
		"SchemaFormat":     cg.schemaFormat,
		"QueryHandlers":    cg.renderOperationHandlers(cg.functionHandlers),
		"MutationHandlers": cg.renderOperationHandlers(cg.procedureHandlers),
	})
}

func (cg *connectorGenerator) renderOperationHandlers(values []string) string {
	if len(values) == 0 {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteRune('{')
	slices.Sort(values)

	for i, v := range values {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(v)
		sb.WriteString(".DataConnectorHandler{}")
	}

	sb.WriteRune('}')
	return sb.String()
}

// generate encoding and decoding methods for schema types.
func (cg *connectorGenerator) genTypeMethods() error {
	cg.genFunctionArgumentConstructors()
	if err := cg.genObjectMethods(); err != nil {
		return err
	}
	if err := cg.genCustomScalarMethods(); err != nil {
		return err
	}
	cg.genConnectorHandlers()

	log.Debug().Msg("Generating types...")
	for packagePath, builder := range cg.typeBuilders {
		if !strings.HasPrefix(packagePath, cg.moduleName) {
			continue
		}
		relativePath := strings.TrimPrefix(packagePath, cg.moduleName)
		schemaPath := path.Join(cg.basePath, relativePath, typeMethodsOutputFile)
		log.Debug().
			Str("package_name", builder.packageName).
			Str("package_path", packagePath).
			Str("module_name", cg.moduleName).
			Msg(schemaPath)

		if err := os.WriteFile(schemaPath, []byte(builder.String()), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (cg *connectorGenerator) genObjectMethods() error {
	if len(cg.rawSchema.Objects) == 0 {
		return nil
	}

	objectKeys := utils.GetSortedKeys(cg.rawSchema.Objects)

	for _, objectName := range objectKeys {
		object := cg.rawSchema.Objects[objectName]
		if object.Type.IsAnonymous || !object.Type.CanMethod() || !strings.HasPrefix(object.Type.PackagePath, cg.moduleName) {
			continue
		}
		sb := cg.getOrCreateTypeBuilder(object.Type.PackagePath)
		sb.builder.WriteString(fmt.Sprintf(`
// ToMap encodes the struct to a value map
func (j %s) ToMap() map[string]any {
  r := make(map[string]any)
`, object.Type.String()))
		cg.writeObjectToMap(sb, &object, "j", "r")
		sb.builder.WriteString(`
  return r
}`)
	}

	return nil
}

func (cg *connectorGenerator) writeObjectToMap(sb *connectorTypeBuilder, object *ObjectInfo, selector string, name string) {
	fieldKeys := utils.GetSortedKeys(object.Fields)
	for _, fieldKey := range fieldKeys {
		field := object.Fields[fieldKey]
		fieldSelector := fmt.Sprintf("%s.%s", selector, field.Name)
		fieldAssigner := fmt.Sprintf("%s[\"%s\"]", name, fieldKey)
		if cg.rawSchema.GetScalarFromType(field.Type) != nil {
			sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", fieldAssigner, fieldSelector))
			continue
		}

		cg.writeToMapProperty(sb, &field, fieldSelector, fieldAssigner, field.Type)
	}
}

func (cg *connectorGenerator) writeToMapProperty(sb *connectorTypeBuilder, field *Field, selector string, assigner string, ty Type) string {
	switch t := ty.(type) {
	case *NullableType:
		sb.builder.WriteString(fmt.Sprintf("  if %s != nil {\n", selector))
		newSelector := selector

		if t.UnderlyingType.Kind() != schema.TypePredicate {
			newSelector = fmt.Sprintf("(*%s)", selector)
		}

		propName := cg.writeToMapProperty(sb, field, newSelector, assigner, t.UnderlyingType)
		sb.builder.WriteString("  }\n")

		return propName
	case *ArrayType:
		varName := formatLocalFieldName(selector)
		valueName := varName + "_v"
		sb.builder.WriteString(fmt.Sprintf("  %s := make([]any, len(%s))\n", varName, selector))
		sb.builder.WriteString(fmt.Sprintf("  for i, %s := range %s {\n", valueName, selector))
		cg.writeToMapProperty(sb, field, valueName, varName+"[i]", t.ElementType)
		sb.builder.WriteString("  }\n")
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))

		return varName
	case *NamedType:
		innerObject, ok := cg.rawSchema.Objects[t.Name]
		if t.NativeType.IsAnonymous {
			if !ok {
				return selector
			}

			// anonymous struct
			varName := formatLocalFieldName(selector, "obj")
			sb.builder.WriteString(fmt.Sprintf("  %s := make(map[string]any)\n", varName))
			cg.writeObjectToMap(sb, &innerObject, selector, varName)
			sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))

			return varName
		}

		if !field.Embedded {
			sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, selector))

			return selector
		}

		sb.imports[packageSDKUtils] = ""
		sb.builder.WriteString(fmt.Sprintf("  r = utils.MergeMap(r, %s.ToMap())\n", selector))

		return selector
	case *PredicateType:
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, selector))

		return selector
	default:
		panic(fmt.Errorf("failed to write the ToMap method; invalid type: %s", ty))
	}
}

// generate Scalar implementation for custom scalar types.
func (cg *connectorGenerator) genCustomScalarMethods() error {
	if len(cg.rawSchema.Scalars) == 0 {
		return nil
	}

	scalarKeys := utils.GetSortedKeys(cg.rawSchema.Scalars)

	for _, scalarKey := range scalarKeys {
		scalar := cg.rawSchema.Scalars[scalarKey]
		if scalar.NativeType == nil || scalar.NativeType.TypeAST == nil {
			continue
		}
		sb := cg.getOrCreateTypeBuilder(scalar.NativeType.PackagePath)
		sb.builder.WriteString(`
// ScalarName get the schema name of the scalar
func (j `)
		sb.builder.WriteString(scalar.NativeType.Name)
		sb.builder.WriteString(`) ScalarName() string {
  return "`)
		sb.builder.WriteString(scalar.NativeType.SchemaName)
		sb.builder.WriteString(`"
}
`)
		// generate enum and parsers if exist
		switch scalarRep := scalar.Schema.Representation.Interface().(type) {
		case *schema.TypeRepresentationEnum:
			sb.imports["errors"] = ""
			sb.imports["encoding/json"] = ""
			sb.imports["slices"] = ""
			sb.imports[packageSDKUtils] = ""

			sb.builder.WriteString("const (\n")
			pascalName := strcase.ToCamel(scalar.NativeType.Name)
			enumConstants := make([]string, len(scalarRep.OneOf))
			for i, enum := range scalarRep.OneOf {
				enumConst := fmt.Sprintf("%s%s", pascalName, strcase.ToCamel(enum))
				enumConstants[i] = enumConst
				sb.builder.WriteString(fmt.Sprintf("  %s %s = \"%s\"\n", enumConst, scalarKey, enum))
			}
			sb.builder.WriteString(fmt.Sprintf(`)

var enumValues_%s = []%s{%s}
`, pascalName, scalarKey, strings.Join(enumConstants, ", ")))

			sb.builder.WriteString(fmt.Sprintf(`
// Parse%s parses a %s enum from string
func Parse%s(input string) (%s, error) {
  result := %s(input)
  if !slices.Contains(enumValues_%s, result) {
    return %s(""), errors.New("failed to parse %s, expect one of [%s]")
  }

  return result, nil
}

// IsValid checks if the value is invalid
func (j %s) IsValid() bool {
  return slices.Contains(enumValues_%s, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *%s) UnmarshalJSON(b []byte) error {
  var rawValue string
  if err := json.Unmarshal(b, &rawValue); err != nil {
    return err
  }

  value, err := Parse%s(rawValue)
  if err != nil {
    return err
  }

  *j = value
  return nil
}

// FromValue decodes the scalar from an unknown value
func (s *%s) FromValue(value any) error {
  valueStr, err := utils.DecodeNullableString(value)
  if err != nil {
    return err
  }
  if valueStr == nil {
    return nil
  }
  result, err := Parse%s(*valueStr)
  if err != nil {
    return err
  }

  *s = result
  return nil
}
`, pascalName, scalarKey, pascalName, scalarKey, scalarKey, pascalName, scalarKey, scalarKey, strings.Join(scalarRep.OneOf, ", "),
				scalarKey, pascalName, scalarKey, pascalName, scalarKey, pascalName,
			))
		default:
		}
	}

	return nil
}

func (cg *connectorGenerator) genFunctionArgumentConstructors() {
	if len(cg.rawSchema.Functions) == 0 {
		return
	}

	writtenObjects := make(map[string]bool)
	fnKeys := utils.GetSortedKeys(cg.rawSchema.FunctionArguments)
	for _, k := range fnKeys {
		argInfo := cg.rawSchema.FunctionArguments[k]
		writtenObjects[argInfo.Type.String()] = true
		cg.writeObjectFromValue(&argInfo)
	}
}

func (cg *connectorGenerator) writeObjectFromValue(info *ObjectInfo) {
	if len(info.Fields) == 0 || info.Type == nil || !info.Type.CanMethod() || info.Type.IsAnonymous {
		return
	}

	sb := cg.getOrCreateTypeBuilder(info.Type.PackagePath)
	sb.builder.WriteString(`
// FromValue decodes values from map
func (j *`)
	sb.builder.WriteString(info.Type.String())
	sb.builder.WriteString(`) FromValue(input map[string]any) error {
  var err error
`)
	argumentKeys := utils.GetSortedKeys(info.Fields)
	for _, key := range argumentKeys {
		arg := info.Fields[key]
		schemaField := info.SchemaFields[key]

		cg.writeGetTypeValueDecoder(sb, &arg, schemaField, key)
	}
	sb.builder.WriteString("  return nil\n}")
}

func (cg *connectorGenerator) getOrCreateTypeBuilder(packagePath string) *connectorTypeBuilder {
	pkgParts := strings.Split(packagePath, "/")
	pkgName := pkgParts[len(pkgParts)-1]
	bs, ok := cg.typeBuilders[packagePath]
	if !ok {
		bs = &connectorTypeBuilder{
			packageName: pkgName,
			packagePath: packagePath,
			imports:     map[string]string{},
			builder:     &strings.Builder{},
		}
		cg.typeBuilders[packagePath] = bs
	}
	return bs
}

func (cg *connectorGenerator) genConnectorHandlers() {
	for _, fn := range cg.rawSchema.Functions {
		builder := cg.getOrCreateTypeBuilder(fn.PackagePath)
		builder.functions = append(builder.functions, fn)
	}
	for _, fn := range cg.rawSchema.Procedures {
		builder := cg.getOrCreateTypeBuilder(fn.PackagePath)
		builder.procedures = append(builder.procedures, fn)
	}

	for _, bs := range cg.typeBuilders {
		fnLen := len(bs.functions)
		procLen := len(bs.procedures)
		if fnLen == 0 && procLen == 0 {
			continue
		}
		cg.rawSchema.Imports[bs.packagePath] = true
		if fnLen > 0 {
			cg.functionHandlers = append(cg.functionHandlers, bs.packageName)
		}
		if procLen > 0 {
			cg.procedureHandlers = append(cg.procedureHandlers, bs.packageName)
		}
		(&connectorHandlerBuilder{
			Builder:    bs,
			RawSchema:  cg.rawSchema,
			Functions:  bs.functions,
			Procedures: bs.procedures,
		}).Render()
	}
}

func (cg *connectorGenerator) writeGetTypeValueDecoder(sb *connectorTypeBuilder, field *Field, objectField schema.ObjectField, key string) {
	ty := field.Type
	fieldName := field.Name
	typeName := ty.String()
	fullTypeName := ty.FullName()

	if strings.Contains(typeName, "complex64") || strings.Contains(typeName, "complex128") || strings.Contains(typeName, "time.Duration") {
		panic(fmt.Errorf("unsupported type: %s", typeName))
	}

	switch fullTypeName {
	case "bool":
		cg.writeScalarDecodeValue(sb, fieldName, "GetBoolean", "", key, objectField, false)
	case "[]bool":
		cg.writeScalarDecodeValue(sb, fieldName, "GetBooleanSlice", "", key, objectField, false)
	case "*bool":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableBoolean", "", key, objectField, true)
	case "[]*bool":
		cg.writeScalarDecodeValue(sb, fieldName, "GetBooleanPtrSlice", "", key, objectField, false)
	case "*[]*bool":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableBooleanPtrSlice", "", key, objectField, true)
	case "string":
		cg.writeScalarDecodeValue(sb, fieldName, "GetString", "", key, objectField, false)
	case "[]string":
		cg.writeScalarDecodeValue(sb, fieldName, "GetStringSlice", "", key, objectField, false)
	case "*string":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableString", "", key, objectField, true)
	case "[]*string":
		cg.writeScalarDecodeValue(sb, fieldName, "GetStringPtrSlice", "", key, objectField, false)
	case "*[]*string":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableStringPtrSlice", "", key, objectField, true)
	case "int", "int8", "int16", "int32", "int64", "rune", "byte":
		cg.writeScalarDecodeValue(sb, fieldName, "GetInt", typeName, key, objectField, false)
	case "[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]rune", "[]byte":
		cg.writeScalarDecodeValue(sb, fieldName, "GetIntSlice", typeName[2:], key, objectField, false)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUint", typeName, key, objectField, false)
	case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUintSlice", typeName[2:], key, objectField, false)
	case "*int", "*int8", "*int16", "*int32", "*int64", "*rune", "*byte":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableInt", typeName[1:], key, objectField, true)
	case "[]*int", "[]*int8", "[]*int16", "[]*int32", "[]*int64", "[]*rune", "[]*byte":
		cg.writeScalarDecodeValue(sb, fieldName, "GetIntPtrSlice", typeName[3:], key, objectField, false)
	case "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableUint", typeName[1:], key, objectField, true)
	case "[]*uint", "[]*uint8", "[]*uint16", "[]*uint32", "[]*uint64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUintPtrSlice", typeName[3:], key, objectField, false)
	case "*[]*int", "*[]*int8", "*[]*int16", "*[]*int32", "*[]*int64", "*[]*rune", "*[]*byte":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableIntPtrSlice", typeName[4:], key, objectField, true)
	case "*[]*uint", "*[]*uint8", "*[]*uint16", "*[]*uint32", "*[]*uint64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableUintPtrSlice", typeName[4:], key, objectField, true)
	case "float32", "float64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetFloat", typeName, key, objectField, false)
	case "[]float32", "[]float64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetFloatSlice", typeName[2:], key, objectField, false)
	case "*float32", "*float64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableFloat", typeName[1:], key, objectField, true)
	case "[]*float32", "[]*float64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetFloatPtrSlice", typeName[3:], key, objectField, false)
	case "*[]*float32", "*[]*float64":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableFloatPtrSlice", typeName[4:], key, objectField, true)
	case "time.Time":
		cg.writeScalarDecodeValue(sb, fieldName, "GetDateTime", "", key, objectField, false)
	case "[]time.Time":
		cg.writeScalarDecodeValue(sb, fieldName, "GetDateTimeSlice", "", key, objectField, false)
	case "*time.Time":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableDateTime", "", key, objectField, true)
	case "[]*time.Time":
		cg.writeScalarDecodeValue(sb, fieldName, "GetDateTimePtrSlice", "", key, objectField, false)
	case "*[]*time.Time":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableDateTimePtrSlice", "", key, objectField, true)
	case "github.com/google/uuid.UUID":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUUID", "", key, objectField, false)
	case "[]github.com/google/uuid.UUID":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUUIDSlice", "", key, objectField, false)
	case "*github.com/google/uuid.UUID":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableUUID", "", key, objectField, true)
	case "[]*github.com/google/uuid.UUID":
		cg.writeScalarDecodeValue(sb, fieldName, "GetUUIDPtrSlice", "", key, objectField, false)
	case "*[]*github.com/google/uuid.UUID":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableUUIDPtrSlice", "", key, objectField, true)
	case "encoding/json.RawMessage":
		cg.writeScalarDecodeValue(sb, fieldName, "GetRawJSON", "", key, objectField, false)
	case "[]encoding/json.RawMessage":
		cg.writeScalarDecodeValue(sb, fieldName, "GetRawJSONSlice", "", key, objectField, false)
	case "*encoding/json.RawMessage":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableRawJSON", "", key, objectField, true)
	case "[]*encoding/json.RawMessage":
		cg.writeScalarDecodeValue(sb, fieldName, "GetRawJSONPtrSlice", "", key, objectField, false)
	case "*[]*encoding/json.RawMessage":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableRawJSONPtrSlice", "", key, objectField, true)
	case "any", "interface{}":
		cg.writeScalarDecodeValue(sb, fieldName, "GetArbitraryJSON", "", key, objectField, false)
	case "[]any", "[]interface{}":
		cg.writeScalarDecodeValue(sb, fieldName, "GetArbitraryJSONSlice", "", key, objectField, false)
	case "*any", "*interface{}":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableArbitraryJSON", "", key, objectField, true)
	case "[]*any", "[]*interface{}":
		cg.writeScalarDecodeValue(sb, fieldName, "GetArbitraryJSONPtrSlice", "", key, objectField, false)
	case "*[]*any", "*[]*interface{}":
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableArbitraryJSONPtrSlice", "", key, objectField, true)
	case fmt.Sprintf("*%s.Expression", packageSDKSchema):
		cg.writeScalarDecodeValue(sb, fieldName, "GetNullableUUID", "", key, objectField, true)
	default:
		sb.imports[packageSDKUtils] = ""
		sb.builder.WriteString("  j.")
		sb.builder.WriteString(fieldName)
		sb.builder.WriteString(", err = utils.")

		t := ty
		var tyName string
		var packagePaths []string

		if nt, ok := ty.(*NullableType); ok {
			if nt.UnderlyingType.Kind() != schema.TypePredicate {
				if nt.IsAnonymous() {
					tyName, packagePaths = cg.getAnonymousObjectTypeName(sb, field.TypeAST, true)
				} else {
					packagePaths = getTypePackagePaths(nt.UnderlyingType, sb.packagePath)
					tyName = getTypeArgumentName(nt.UnderlyingType, sb.packagePath, false)
				}

				for _, pkgPath := range packagePaths {
					sb.imports[pkgPath] = ""
				}

				if field.Embedded {
					sb.builder.WriteString("DecodeNullableObject[")
					sb.builder.WriteString(tyName)
					sb.builder.WriteString("](input)")
				} else {
					sb.builder.WriteString("DecodeNullableObjectValue[")
					sb.builder.WriteString(tyName)
					sb.builder.WriteString(`](input, "`)
					sb.builder.WriteString(key)
					sb.builder.WriteString(`")`)
				}

				break
			}

			t = nt.UnderlyingType
		}

		if t.IsAnonymous() {
			tyName, packagePaths = cg.getAnonymousObjectTypeName(sb, field.TypeAST, true)
		} else {
			packagePaths = getTypePackagePaths(t, sb.packagePath)
			tyName = getTypeArgumentName(t, sb.packagePath, false)
		}

		for _, pkgPath := range packagePaths {
			sb.imports[pkgPath] = ""
		}

		if field.Embedded {
			sb.builder.WriteString("DecodeObject")
			sb.builder.WriteRune('[')
			sb.builder.WriteString(tyName)
			sb.builder.WriteString("](input)")
		} else {
			sb.builder.WriteString("DecodeObjectValue")

			if len(objectField.Type) > 0 {
				if typeEnum, err := objectField.Type.Type(); err == nil && typeEnum == schema.TypeNullable {
					sb.builder.WriteString("Default")
				}
			}

			sb.builder.WriteRune('[')
			sb.builder.WriteString(tyName)
			sb.builder.WriteString(`](input, "`)
			sb.builder.WriteString(key)
			sb.builder.WriteString(`")`)
		}
	}

	writeErrorCheck(sb.builder, 1, 2)
}

func (cg *connectorGenerator) writeScalarDecodeValue(sb *connectorTypeBuilder, fieldName, functionName, typeParam, key string, objectField schema.ObjectField, isNullable bool) {
	sb.imports[packageSDKUtils] = ""
	sb.builder.WriteString("  j.")
	sb.builder.WriteString(fieldName)
	sb.builder.WriteString(", err = utils.")
	sb.builder.WriteString(functionName)
	if !isNullable && len(objectField.Type) > 0 {
		if typeEnum, err := objectField.Type.Type(); err == nil && typeEnum == schema.TypeNullable {
			sb.builder.WriteString("Default")
		}
	}
	if typeParam != "" {
		sb.builder.WriteRune('[')
		sb.builder.WriteString(typeParam)
		sb.builder.WriteRune(']')
	}
	sb.builder.WriteString(`(input, "`)
	sb.builder.WriteString(key)
	sb.builder.WriteString(`")`)
}

// generate anonymous object type name with absolute package paths removed
func (cg *connectorGenerator) getAnonymousObjectTypeName(sb *connectorTypeBuilder, goType types.Type, skipNullable bool) (string, []string) {
	switch inferredType := goType.(type) {
	case *types.Pointer:
		var result string
		if !skipNullable {
			result += "*"
		}
		underlyingName, packagePaths := cg.getAnonymousObjectTypeName(sb, inferredType.Elem(), false)
		return result + underlyingName, packagePaths
	case *types.Struct:
		packagePaths := []string{}
		result := "struct{"
		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)
			if i > 0 {
				result += "; "
			}
			result += fieldVar.Name() + " "
			underlyingName, pkgPaths := cg.getAnonymousObjectTypeName(sb, fieldVar.Type(), false)
			result += underlyingName
			packagePaths = append(packagePaths, pkgPaths...)
			if fieldTag != "" {
				result += " `" + fieldTag + "`"
			}
		}
		result += "}"
		return result, packagePaths
	case *types.Named:
		packagePaths := []string{}
		innerType := inferredType.Obj()
		if innerType == nil {
			return "", packagePaths
		}

		var result string
		typeInfo := &TypeInfo{
			Name: innerType.Name(),
		}

		innerPkg := innerType.Pkg()
		if innerPkg != nil && innerPkg.Name() != "" && innerPkg.Path() != sb.packagePath {
			packagePaths = append(packagePaths, innerPkg.Path())
			result += innerPkg.Name() + "."
			typeInfo.PackageName = innerPkg.Name()
			typeInfo.PackagePath = innerPkg.Path()
		}

		result += innerType.Name()
		typeParams := inferredType.TypeParams()
		if typeParams != nil && typeParams.Len() > 0 {
			// unwrap the generic type parameters such as Foo[T]
			if err := parseTypeParameters(typeInfo, inferredType.String()); err == nil {
				result += "["
				for i, typeParam := range typeInfo.TypeParameters {
					if i > 0 {
						result += ", "
					}
					packagePaths = append(packagePaths, getTypePackagePaths(typeParam, sb.packagePath)...)
					result += getTypeArgumentName(typeParam, sb.packagePath, false)
				}
				result += "]"
			}
		}

		return result, packagePaths
	case *types.Basic:
		return inferredType.Name(), []string{}
	case *types.Array:
		result, packagePaths := cg.getAnonymousObjectTypeName(sb, inferredType.Elem(), false)
		return "[]" + result, packagePaths
	case *types.Slice:
		result, packagePaths := cg.getAnonymousObjectTypeName(sb, inferredType.Elem(), false)
		return "[]" + result, packagePaths
	default:
		return inferredType.String(), []string{}
	}
}

func formatLocalFieldName(input string, others ...string) string {
	name := fieldNameRegex.ReplaceAllString(input, "_")
	return strings.Trim(strings.Join(append([]string{name}, others...), "_"), "_")
}
