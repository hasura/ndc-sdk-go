package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go/token"
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

// ConnectorGenerationArguments represent input arguments of the ConnectorGenerator
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

// SetImport sets an import package into the import list
func (ctb *connectorTypeBuilder) SetImport(value string, alias string) {
	ctb.imports[value] = alias
}

// String renders generated Go types and methods
func (ctb connectorTypeBuilder) String() string {
	var bs strings.Builder
	writeFileHeader(&bs, ctb.packageName)
	if len(ctb.imports) > 0 {
		bs.WriteString("import (\n")
		sortedImports := utils.GetSortedKeys(ctb.imports)
		for _, pkg := range sortedImports {
			alias := ctb.imports[pkg]
			if alias != "" {
				alias = alias + " "
			}
			bs.WriteString(fmt.Sprintf("  %s\"%s\"\n", alias, pkg))
		}
		bs.WriteString(")\n")
	}

	bs.WriteString("var connector_Decoder = utils.NewDecoder()\n")
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

// ParseAndGenerateConnector parses and generate connector codes
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
	if err := os.WriteFile(schemaPath, schemaBytes, 0644); err != nil {
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

// generate encoding and decoding methods for schema types
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

		if err := os.WriteFile(schemaPath, []byte(builder.String()), 0644); err != nil {
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
		if object.IsAnonymous || !strings.HasPrefix(object.Type.PackagePath, cg.moduleName) {
			continue
		}
		sb := cg.getOrCreateTypeBuilder(object.Type.PackagePath)
		sb.builder.WriteString(fmt.Sprintf(`
// ToMap encodes the struct to a value map
func (j %s) ToMap() map[string]any {
  r := make(map[string]any)
`, object.Type.GetArgumentName(object.Type.PackagePath)))
		cg.genObjectToMap(sb, &object, "j", "r")
		sb.builder.WriteString(`
	return r
}`)
	}

	return nil
}

func (cg *connectorGenerator) genObjectToMap(sb *connectorTypeBuilder, object *ObjectInfo, selector string, name string) {
	fieldKeys := utils.GetSortedKeys(object.Fields)
	for _, fieldKey := range fieldKeys {
		field := object.Fields[fieldKey]
		fieldSelector := fmt.Sprintf("%s.%s", selector, field.Name)
		fieldAssigner := fmt.Sprintf("%s[\"%s\"]", name, fieldKey)
		cg.genToMapProperty(sb, &field, fieldSelector, fieldAssigner, field.Type, field.Type.TypeFragments)
	}
}

func (cg *connectorGenerator) genToMapProperty(sb *connectorTypeBuilder, field *ObjectField, selector string, assigner string, ty *TypeInfo, fragments []string) string {
	if ty.IsScalar {
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, selector))
		return selector
	}

	if isNullableFragments(fragments) {
		childFragments := fragments[1:]
		sb.builder.WriteString(fmt.Sprintf("  if %s != nil {\n", selector))
		propName := cg.genToMapProperty(sb, field, fmt.Sprintf("(*%s)", selector), assigner, ty, childFragments)
		sb.builder.WriteString("  }\n")
		return propName
	}

	if isArrayFragments(fragments) {
		varName := formatLocalFieldName(selector)
		valueName := varName + "_v"
		sb.builder.WriteString(fmt.Sprintf("  %s := make([]any, len(%s))\n", varName, selector))
		sb.builder.WriteString(fmt.Sprintf("  for i, %s := range %s {\n", valueName, selector))
		cg.genToMapProperty(sb, field, valueName, varName+"[i]", ty, fragments[1:])
		sb.builder.WriteString("  }\n")
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))
		return varName
	}

	isAnonymous := strings.HasPrefix(strings.Join(fragments, ""), "struct{")
	if isAnonymous {
		innerObject, ok := cg.rawSchema.Objects[ty.String()]
		if !ok {
			tyName := buildTypeNameFromFragments(ty.TypeFragments, ty, sb.packagePath)
			innerObject, ok = cg.rawSchema.Objects[tyName]
			if !ok {
				return selector
			}
		}

		// anonymous struct
		varName := formatLocalFieldName(selector, "obj")
		sb.builder.WriteString(fmt.Sprintf("  %s := make(map[string]any)\n", varName))
		cg.genObjectToMap(sb, &innerObject, selector, varName)
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))
		return varName
	}
	if !field.Type.Embedded {
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, selector))
		return selector
	}

	sb.builder.WriteString(fmt.Sprintf("  r = utils.MergeMap(r, %s.ToMap())\n", selector))
	return selector
}

// generate Scalar implementation for custom scalar types
func (cg *connectorGenerator) genCustomScalarMethods() error {
	if len(cg.rawSchema.CustomScalars) == 0 {
		return nil
	}

	scalarKeys := utils.GetSortedKeys(cg.rawSchema.CustomScalars)

	for _, scalarKey := range scalarKeys {
		scalar := cg.rawSchema.CustomScalars[scalarKey]
		sb := cg.getOrCreateTypeBuilder(scalar.PackagePath)
		sb.builder.WriteString(fmt.Sprintf(`
// ScalarName get the schema name of the scalar
func (j %s) ScalarName() string {
  return "%s"
}
`, scalar.Name, scalar.SchemaName))
		// generate enum and parsers if exist
		if scalar.ScalarRepresentation != nil {
			switch scalarRep := scalar.ScalarRepresentation.Interface().(type) {
			case *schema.TypeRepresentationEnum:
				sb.imports["errors"] = ""
				sb.imports["encoding/json"] = ""
				sb.imports["slices"] = ""

				sb.builder.WriteString("const (\n")
				pascalName := strcase.ToCamel(scalar.Name)
				enumConstants := make([]string, len(scalarRep.OneOf))
				for i, enum := range scalarRep.OneOf {
					enumConst := fmt.Sprintf("%s%s", pascalName, strcase.ToCamel(enum))
					enumConstants[i] = enumConst
					sb.builder.WriteString(fmt.Sprintf("  %s %s = \"%s\"\n", enumConst, scalar.Name, enum))
				}
				sb.builder.WriteString(fmt.Sprintf(`)

var enumValues_%s = []%s{%s}
`, pascalName, scalar.Name, strings.Join(enumConstants, ", ")))

				sb.builder.WriteString(fmt.Sprintf(`
// Parse%s parses a %s enum from string
func Parse%s(input string) (%s, error) {
	result := %s(input)
	if !slices.Contains(enumValues_%s, result) {
		return %s(""), errors.New("failed to parse %s, expect one of %s")
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
`, pascalName, scalar.Name, pascalName, scalar.Name, scalar.Name, pascalName, scalar.Name, scalar.Name, strings.Join(enumConstants, ", "),
					scalar.Name, pascalName, scalar.Name, pascalName, scalar.Name, pascalName,
				))
			}
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
	if len(info.Fields) == 0 || info.Type == nil || !info.Type.CanMethod() || info.IsAnonymous {
		return
	}

	sb := cg.getOrCreateTypeBuilder(info.Type.PackagePath)
	sb.builder.WriteString(`
// FromValue decodes values from map
func (j *`)
	sb.builder.WriteString(info.Type.GetArgumentName(info.Type.PackagePath))
	sb.builder.WriteString(`) FromValue(input map[string]any) error {
  var err error
`)
	argumentKeys := utils.GetSortedKeys(info.Fields)
	for _, key := range argumentKeys {
		arg := info.Fields[key]
		cg.genGetTypeValueDecoder(sb, arg.Type, key, arg.Name)
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
			imports: map[string]string{
				"github.com/hasura/ndc-sdk-go/utils": "",
			},
			builder: &strings.Builder{},
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

func (cg *connectorGenerator) genGetTypeValueDecoder(sb *connectorTypeBuilder, ty *TypeInfo, key string, fieldName string) {
	typeName := ty.TypeAST.String()
	if strings.Contains(typeName, "complex64") || strings.Contains(typeName, "complex128") || strings.Contains(typeName, "time.Duration") {
		panic(fmt.Errorf("unsupported type: %s", typeName))
	}

	fullTypeName := strings.Join(ty.TypeFragments, "")
	switch typeName {
	case "bool":
		funcName := "GetBoolean"
		if fullTypeName == "[]bool" {
			funcName = "GetBooleanSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*bool":
		funcName := "GetNullableBoolean"
		if fullTypeName == "[]*bool" {
			funcName = "GetBooleanPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*[]*bool":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableBooleanPtrSlice(input, "%s")`, fieldName, key))
	case "string":
		funcName := "GetString"
		if fullTypeName == "[]string" {
			funcName = "GetStringSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*string":
		funcName := "GetNullableString"
		if fullTypeName == "[]*string" {
			funcName = "GetStringPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*[]*string":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableStringPtrSlice(input, "%s")`, fieldName, key))
	case "int", "int8", "int16", "int32", "int64", "rune", "byte":
		funcName := "GetInt"
		if slices.Contains([]string{"[]int", "[]int8", "[]int16", "[]int32", "[]int64", "[]rune", "[]byte"}, fullTypeName) {
			funcName = "GetIntSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, typeName, key))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		funcName := "GetUint"
		if slices.Contains([]string{"[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64"}, fullTypeName) {
			funcName = "GetUintSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, typeName, key))
	case "*int", "*int8", "*int16", "*int32", "*int64", "*rune", "*byte":
		funcName := "GetNullableInt"
		if slices.Contains([]string{"[]*int", "[]*int8", "[]*int16", "[]*int32", "[]*int64", "[]*rune", "[]*byte"}, fullTypeName) {
			funcName = "GetIntPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, strings.TrimPrefix(typeName, "*"), key))
	case "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		funcName := "GetNullableUint"
		if slices.Contains([]string{"[]*uint", "[]*uint8", "[]*uint16", "[]*uint32", "[]*uint64"}, fullTypeName) {
			funcName = "GetUintPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, strings.TrimPrefix(typeName, "*"), key))
	case "*[]*int", "*[]*int8", "*[]*int16", "*[]*int32", "*[]*int64", "*[]*rune", "*[]*byte":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableIntPtrSlice[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*[]*"), key))
	case "*[]*uint", "*[]*uint8", "*[]*uint16", "*[]*uint32", "*[]*uint64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableUintPtrSlice[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*[]*"), key))
	case "float32", "float64":
		funcName := "GetFloat"
		if slices.Contains([]string{"[]float32", "[]float64"}, fullTypeName) {
			funcName = "GetFloatSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, typeName, key))
	case "*float32", "*float64":
		funcName := "GetNullableFloat"
		if slices.Contains([]string{"[]*float32", "[]*float64"}, fullTypeName) {
			funcName = "GetFloatPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s[%s](input, "%s")`, fieldName, funcName, strings.TrimPrefix(typeName, "*"), key))
	case "*[]*float32", "*[]*float64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableFloatPtrSlice[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*[]*"), key))
	case "time.Time":
		funcName := "GetDateTime"
		if fullTypeName == "[]Time" {
			funcName = "GetDateTimeSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*time.Time":
		funcName := "GetNullableDateTime"
		if fullTypeName == "[]*Time" {
			funcName = "GetDateTimePtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*[]*time.Time":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableDateTimePtrSlice(input, "%s")`, fieldName, key))
	case "github.com/google/uuid.UUID":
		funcName := "GetUUID"
		if fullTypeName == "[]UUID" {
			funcName = "GetUUIDSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*github.com/google/uuid.UUID":
		funcName := "GetNullableUUID"
		if fullTypeName == "[]*UUID" {
			funcName = "GetUUIDPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*[]*github.com/google/uuid.UUID":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableUUIDPtrSlice(input, "%s")`, fieldName, key))
	case "encoding/json.RawMessage":
		funcName := "GetRawJSON"
		if fullTypeName == "[]RawMessage" {
			funcName = "GetRawJSONSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))

	case "*encoding/json.RawMessage":
		funcName := "GetNullableRawJSON"
		if fullTypeName == "[]*RawMessage" {
			funcName = "GetRawJSONPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))
	case "*[]*encoding/json.RawMessage":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableRawJSONPtrSlice(input, "%s")`, fieldName, key))
	case "any", "interface{}":
		funcName := "GetArbitraryJSON"
		if slices.Contains([]string{"[]any", "[]interface{}"}, fullTypeName) {
			funcName = "GetArbitraryJSONSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))

	case "*any", "*interface{}":
		funcName := "GetNullableArbitraryJSON"
		if slices.Contains([]string{"[]*any", "[]*interface{}"}, fullTypeName) {
			funcName = "GetArbitraryJSONPtrSlice"
		}
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.%s(input, "%s")`, fieldName, funcName, key))

	case "*[]*any", "*[]*interface{}":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableArbitraryJSONPtrSlice(input, "%s")`, fieldName, key))

	default:
		if ty.IsNullable() {
			// typeName := strings.TrimLeft(typeName, "*")
			// pkgName, tyName, ok := findAndReplaceNativeScalarPackage(typeName)
			pkgName := ty.PackagePath
			tyName := buildTypeNameFromFragments(ty.TypeFragments[1:], ty, sb.packagePath)

			if pkgName != "" && pkgName != sb.packagePath {
				sb.imports[pkgName] = ""
			}
			sb.builder.WriteString("  j.")
			sb.builder.WriteString(fieldName)
			sb.builder.WriteString(" = new(")
			sb.builder.WriteString(tyName)
			sb.builder.WriteString(")\n  err = connector_Decoder.")
			if ty.Embedded {
				sb.builder.WriteString("DecodeObject(j.")
				sb.builder.WriteString(fieldName)
				sb.builder.WriteString(", input)")
			} else {
				sb.builder.WriteString("DecodeNullableObjectValue(j.")
				sb.builder.WriteString(fieldName)
				sb.builder.WriteString(`, input, "`)
				sb.builder.WriteString(key)
				sb.builder.WriteString(`")`)
			}
		} else if ty.Embedded {
			sb.builder.WriteString("  err = connector_Decoder.DecodeObject(&j.")
			sb.builder.WriteString(fieldName)
			sb.builder.WriteString(", input)")
		} else {
			sb.builder.WriteString("  err = connector_Decoder.DecodeObjectValue(&j.")
			sb.builder.WriteString(fieldName)
			sb.builder.WriteString(`, input, "`)
			sb.builder.WriteString(key)
			sb.builder.WriteString(`")`)
		}
	}
	sb.builder.WriteString(textBlockErrorCheck)
}

func buildTypeNameFromFragments(items []string, ti *TypeInfo, currentPackagePath string) string {
	var results string
	for _, item := range items {
		if item == "*" || item == "[]" {
			results += item
			continue
		}
		results += buildTypeWithAlias(item, ti, currentPackagePath)
	}

	return results
}

func buildTypeWithAlias(name string, ti *TypeInfo, currentPackagePath string) string {
	if ti.Name == "" ||
		// do not add alias to anonymous struct
		strings.HasPrefix(name, "struct{") {
		return name
	}
	return ti.GetArgumentName(currentPackagePath)
}

func formatLocalFieldName(input string, others ...string) string {
	name := fieldNameRegex.ReplaceAllString(input, "_")
	return strings.Trim(strings.Join(append([]string{name}, others...), "_"), "_")
}
