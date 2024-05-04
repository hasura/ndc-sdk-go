package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"runtime/trace"
	"sort"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
)

var fieldNameRegex = regexp.MustCompile(`[^\w]`)

type connectorTypeBuilder struct {
	packageName string
	packagePath string
	imports     map[string]string
	builder     *strings.Builder
}

// SetImport sets an import package into the import list
func (ctb *connectorTypeBuilder) SetImport(value string, alias string) {
	ctb.imports[value] = alias
}

// GetDecoderName gets the global decoder name
func (ctb connectorTypeBuilder) GetDecoderName() string {
	return fmt.Sprintf("%s_Decoder", ctb.packageName)
}

// String renders generated Go types and methods
func (ctb connectorTypeBuilder) String() string {
	var bs strings.Builder
	bs.WriteString(genFileHeader(ctb.packageName))
	if len(ctb.imports) > 0 {
		bs.WriteString("import (\n")
		sortedImports := getSortedKeys(ctb.imports)
		for _, pkg := range sortedImports {
			alias := ctb.imports[pkg]
			if alias != "" {
				alias = alias + " "
			}
			bs.WriteString(fmt.Sprintf("  %s\"%s\"\n", alias, pkg))
		}
		bs.WriteString(")\n")
	}

	decoderName := ctb.GetDecoderName()
	bs.WriteString(fmt.Sprintf("var %s = utils.NewDecoder()\n", decoderName))
	bs.WriteString(ctb.builder.String())
	return bs.String()
}

type connectorGenerator struct {
	basePath     string
	moduleName   string
	rawSchema    *RawConnectorSchema
	typeBuilders map[string]*connectorTypeBuilder
}

func NewConnectorGenerator(basePath string, moduleName string, rawSchema *RawConnectorSchema) *connectorGenerator {
	return &connectorGenerator{
		basePath:     basePath,
		moduleName:   moduleName,
		rawSchema:    rawSchema,
		typeBuilders: make(map[string]*connectorTypeBuilder),
	}
}

func parseAndGenerateConnector(args *GenerateArguments, moduleName string) error {
	if cli.Generate.Trace != "" {
		w, err := os.Create(cli.Generate.Trace)
		if err != nil {
			return fmt.Errorf("failed to create trace file at %s", cli.Generate.Trace)
		}
		defer func() {
			_ = w.Close()
		}()
		if err = trace.Start(w); err != nil {
			return fmt.Errorf("failed to start trace: %v", err)
		}
		defer trace.Stop()
	}

	parseCtx, parseTask := trace.NewTask(context.TODO(), "parse")
	sm, err := parseRawConnectorSchemaFromGoCode(parseCtx, moduleName, ".", args.Directories)
	if err != nil {
		parseTask.End()
		return err
	}
	parseTask.End()

	_, genTask := trace.NewTask(context.TODO(), "generate_code")
	defer genTask.End()
	connectorGen := NewConnectorGenerator(".", moduleName, sm)
	return connectorGen.generateConnector()
}

func (cg *connectorGenerator) generateConnector() error {
	// generate schema.generated.json
	schemaBytes, err := json.MarshalIndent(cg.rawSchema.Schema(), "", "  ")
	if err != nil {
		return err
	}

	schemaPath := path.Join(cg.basePath, schemaOutputFile)
	if err := os.WriteFile(schemaPath, schemaBytes, 0644); err != nil {
		return err
	}

	targetPath := path.Join(cg.basePath, connectorOutputFile)
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

	if err := cg.genConnectorCodeFromTemplate(w); err != nil {
		return err
	}

	if err := cg.genTypeMethods(); err != nil {
		return err
	}

	return nil
}

func (cg *connectorGenerator) genConnectorCodeFromTemplate(w io.Writer) error {
	queries := cg.genConnectorFunctions(cg.rawSchema)
	procedures := cg.genConnectorProcedures(cg.rawSchema)

	sortedImports := getSortedKeys(cg.rawSchema.Imports)

	importLines := []string{}
	for _, importPath := range sortedImports {
		importLines = append(importLines, fmt.Sprintf(`"%s"`, importPath))
	}

	return connectorTemplate.Execute(w, map[string]any{
		"Imports":    strings.Join(importLines, "\n"),
		"Module":     cg.moduleName,
		"Queries":    queries,
		"Procedures": procedures,
	})
}

func (cg *connectorGenerator) genConnectorFunctions(rawSchema *RawConnectorSchema) string {
	if len(rawSchema.Functions) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, fn := range rawSchema.Functions {
		if _, ok := cg.rawSchema.Imports[fn.PackagePath]; !ok {
			cg.rawSchema.Imports[fn.PackagePath] = true
		}
		var argumentParamStr string
		sb.WriteString(fmt.Sprintf("  case \"%s\":", fn.Name))
		if fn.ResultType.IsScalar {
			sb.WriteString(`
    if len(queryFields) > 0 {
      return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
    }`)
		} else if fn.ResultType.IsArray() {
			sb.WriteString(`
    selection, err := queryFields.AsArray()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
        "cause": err.Error(),
      })
    }`)
		} else {
			sb.WriteString(`
    selection, err := queryFields.AsObject()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
        "cause": err.Error(),
      })
    }`)
		}

		if fn.ArgumentsType != "" {
			argumentStr := fmt.Sprintf(`
    rawArgs, err := utils.ResolveArgumentVariables(request.Arguments, variables)
    if err != nil {
      return nil, schema.UnprocessableContentError("failed to resolve argument variables", map[string]any{
        "cause": err.Error(),
      })
    }
    
		connector_addSpanEvent(span, logger, "resolve_arguments", map[string]any{
			"raw_arguments": rawArgs,
		})
		
    var args %s.%s
    if err = args.FromValue(rawArgs); err != nil {
      return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
        "cause": err.Error(),
      })
    }
		
		connector_addSpanEvent(span, logger, "execute_function", map[string]any{
			"arguments": args,
		})`, fn.PackageName, fn.ArgumentsType)
			sb.WriteString(argumentStr)
			argumentParamStr = ", &args"
		}

		if fn.ResultType.IsScalar {
			sb.WriteString(fmt.Sprintf("\n    return %s.%s(ctx, state%s)\n", fn.PackageName, fn.OriginName, argumentParamStr))
			continue
		}

		sb.WriteString(fmt.Sprintf("\n    rawResult, err := %s.%s(ctx, state%s)", fn.PackageName, fn.OriginName, argumentParamStr))
		genGeneralOperationResult(&sb, fn.ResultType)

		sb.WriteString(`
		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})`)
		if fn.ResultType.IsArray() {
			sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)")
		} else {
			sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)")
		}
		sb.WriteString(textBlockErrorCheck2)
		sb.WriteString("    return result, nil\n")
	}

	return sb.String()
}

func genGeneralOperationResult(sb *strings.Builder, resultType *TypeInfo) {
	sb.WriteString(textBlockErrorCheck2)
	if resultType.IsNullable() {
		sb.WriteString(`
    if rawResult == nil {
      return nil, nil
    }
`)
	} else {
		sb.WriteString(`
    if rawResult == nil {
      return nil, schema.UnprocessableContentError("expected not null result", nil)
    }
`)
	}
}

func (cg *connectorGenerator) genConnectorProcedures(rawSchema *RawConnectorSchema) string {
	if len(rawSchema.Procedures) == 0 {
		return ""
	}

	cg.rawSchema.Imports["encoding/json"] = true

	var sb strings.Builder
	for _, fn := range rawSchema.Procedures {
		if _, ok := cg.rawSchema.Imports[fn.PackagePath]; !ok {
			cg.rawSchema.Imports[fn.PackagePath] = true
		}
		var argumentParamStr string
		sb.WriteString(fmt.Sprintf("  case \"%s\":", fn.Name))
		if fn.ResultType.IsScalar {
			sb.WriteString(`
    if len(operation.Fields) > 0 {
      return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
    }`)
		} else if fn.ResultType.IsArray() {
			sb.WriteString(`
    selection, err := operation.Fields.AsArray()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
        "cause": err.Error(),
      })
    }`)
		} else {
			sb.WriteString(`
    selection, err := operation.Fields.AsObject()
    if err != nil {
      return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
        "cause": err.Error(),
      })
    }`)
		}
		if fn.ArgumentsType != "" {
			argumentStr := fmt.Sprintf(`
    var args %s.%s
    if err := json.Unmarshal(operation.Arguments, &args); err != nil {
      return nil, schema.UnprocessableContentError("failed to decode arguments", map[string]any{
        "cause": err.Error(),
      })
    }`, fn.PackageName, fn.ArgumentsType)
			sb.WriteString(argumentStr)
			argumentParamStr = ", &args"
		}

		sb.WriteString("\n    span.AddEvent(\"execute_procedure\")")
		if fn.ResultType.IsScalar {
			sb.WriteString(fmt.Sprintf(`
    var err error
    result, err := %s.%s(ctx, state%s)`, fn.PackageName, fn.OriginName, argumentParamStr))
		} else {
			sb.WriteString(fmt.Sprintf("\n    rawResult, err := %s.%s(ctx, state%s)\n", fn.PackageName, fn.OriginName, argumentParamStr))
			genGeneralOperationResult(&sb, fn.ResultType)

			sb.WriteString(`    connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})`)
			if fn.ResultType.IsArray() {
				sb.WriteString("\n    result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)\n")
			} else {
				sb.WriteString("\n    result, err := utils.EvalNestedColumnObject(selection, rawResult)\n")
			}
		}

		sb.WriteString(textBlockErrorCheck2)
		sb.WriteString("    return schema.NewProcedureResult(result).Encode(), nil\n")
	}

	return sb.String()
}

func (cg *connectorGenerator) genTypeMethods() error {
	if err := cg.genFunctionArgumentConstructors(); err != nil {
		return err
	}
	if err := cg.genObjectMethods(); err != nil {
		return err
	}
	if err := cg.genCustomScalarMethods(); err != nil {
		return err
	}
	for folderPath, builder := range cg.typeBuilders {
		schemaPath := path.Join(cg.basePath, folderPath, typeMethodsOutputFile)
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

	objectKeys := getSortedKeys(cg.rawSchema.Objects)

	for _, objectName := range objectKeys {
		object := cg.rawSchema.Objects[objectName]
		if object.IsAnonymous {
			continue
		}
		sb := cg.getOrCreateTypeBuilder(object.PackageName, object.PackageName, object.PackagePath)
		sb.builder.WriteString(fmt.Sprintf(`
// ToMap encodes the struct to a value map
func (j %s) ToMap() map[string]any {
  r := make(map[string]any)
`, objectName))
		cg.genObjectToMap(sb, object, "j", "r")
		sb.builder.WriteString(`
	return r
}`)
	}

	return nil
}

func (cg *connectorGenerator) genObjectToMap(sb *connectorTypeBuilder, object *ObjectInfo, selector string, name string) {

	fieldKeys := getSortedKeys(object.Fields)
	for _, fieldKey := range fieldKeys {
		field := object.Fields[fieldKey]
		fieldSelector := fmt.Sprintf("%s.%s", selector, field.Name)
		fieldAssigner := fmt.Sprintf("%s[\"%s\"]", name, field.Key)
		cg.genToMapProperty(sb, field, fieldSelector, fieldAssigner, field.Type, field.Type.TypeFragments)
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
		valueName := fmt.Sprintf("%s_v", varName)
		sb.builder.WriteString(fmt.Sprintf("  %s := make([]map[string]any, len(%s))\n", varName, selector))
		sb.builder.WriteString(fmt.Sprintf("  for i, %s := range %s {\n", valueName, selector))
		cg.genToMapProperty(sb, field, valueName, fmt.Sprintf("%s[i]", varName), ty, fragments[1:])
		sb.builder.WriteString("  }\n")
		sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))
		return varName
	}

	isAnonymous := strings.HasPrefix(strings.Join(fragments, ""), "struct{")
	if !isAnonymous {
		sb.builder.WriteString(fmt.Sprintf("  %s = utils.EncodeMap(%s)\n", assigner, selector))
		return selector
	}
	innerObject, ok := cg.rawSchema.Objects[ty.Name]
	if !ok {
		_, tyName := buildTypeNameFromFragments(ty.TypeFragments, sb.packagePath)
		innerObject, ok = cg.rawSchema.Objects[tyName]
		if !ok {
			return selector
		}
	}

	// anonymous struct
	varName := formatLocalFieldName(selector, "obj")
	sb.builder.WriteString(fmt.Sprintf("  %s := make(map[string]any)\n", varName))
	cg.genObjectToMap(sb, innerObject, selector, varName)
	sb.builder.WriteString(fmt.Sprintf("  %s = %s\n", assigner, varName))
	return varName
}

// generate Scalar implementation for custom scalar types
func (cg *connectorGenerator) genCustomScalarMethods() error {
	if len(cg.rawSchema.CustomScalars) == 0 {
		return nil
	}

	scalarKeys := getSortedKeys(cg.rawSchema.CustomScalars)

	for _, scalarKey := range scalarKeys {
		scalar := cg.rawSchema.CustomScalars[scalarKey]
		sb := cg.getOrCreateTypeBuilder(scalar.PackageName, scalar.PackageName, scalar.PackagePath)
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
				sb.imports["github.com/hasura/ndc-sdk-go/schema"] = ""

				sb.builder.WriteString("const (\n")
				pascalName := ToPascalCase(scalar.Name)
				enumConstants := make([]string, len(scalarRep.OneOf))
				for i, enum := range scalarRep.OneOf {
					enumConst := fmt.Sprintf("%s%s", pascalName, ToPascalCase(enum))
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
	if !schema.Contains(enumValues_%s, result) {
		return %s(""), errors.New("failed to parse %s, expect one of %s")
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j %s) IsValid() bool {
	return schema.Contains(enumValues_%s, j)
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

func (cg *connectorGenerator) genFunctionArgumentConstructors() error {
	if len(cg.rawSchema.Functions) == 0 {
		return nil
	}

	for _, fn := range cg.rawSchema.Functions {
		if len(fn.Arguments) == 0 {
			continue
		}
		sb := cg.getOrCreateTypeBuilder(fn.PackageName, fn.PackageName, fn.PackagePath)
		sb.builder.WriteString(fmt.Sprintf(`
// FromValue decodes values from map
func (j *%s) FromValue(input map[string]any) error {
  var err error
`, fn.ArgumentsType))

		argumentKeys := getSortedKeys(fn.Arguments)
		for _, key := range argumentKeys {
			arg := fn.Arguments[key]
			cg.genGetTypeValueDecoder(sb, arg.Type, key, arg.FieldName)
		}
		sb.builder.WriteString(`  return nil
}`)
	}

	return nil
}

func (cg *connectorGenerator) getOrCreateTypeBuilder(fileName string, packageName string, packagePath string) *connectorTypeBuilder {
	bs, ok := cg.typeBuilders[fileName]
	if !ok {
		bs = &connectorTypeBuilder{
			packageName: packageName,
			packagePath: packagePath,
			imports: map[string]string{
				"github.com/hasura/ndc-sdk-go/utils": "",
			},
			builder: &strings.Builder{},
		}
		cg.typeBuilders[fileName] = bs
	}
	return bs
}

func genFileHeader(packageName string) string {
	return fmt.Sprintf(`// Code generated by github.com/hasura/ndc-sdk-go/codegen, DO NOT EDIT.
package %s
`, packageName)
}

func (cg *connectorGenerator) genGetTypeValueDecoder(sb *connectorTypeBuilder, ty *TypeInfo, key string, fieldName string) {
	typeName := ty.TypeAST.String()
	if strings.Contains(typeName, "complex64") || strings.Contains(typeName, "complex128") || strings.Contains(typeName, "time.Duration") {
		panic(fmt.Errorf("unsupported type: %s", typeName))
	}

	switch typeName {
	case "bool":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetBool(input, "%s")`, fieldName, key))
	case "*bool":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableBool(input, "%s")`, fieldName, key))
	case "string":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetString(input, "%s")`, fieldName, key))
	case "*string":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableString(input, "%s")`, fieldName, key))
	case "int", "int8", "int16", "int32", "int64", "rune", "byte":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetInt[%s](input, "%s")`, fieldName, typeName, key))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetUint[%s](input, "%s")`, fieldName, typeName, key))
	case "*int", "*int8", "*int16", "*int32", "*int64", "*rune", "*byte":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableInt[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*"), key))
	case "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableUint[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*"), key))
	case "float32", "float64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetFloat[%s](input, "%s")`, fieldName, typeName, key))
	case "*float32", "*float64":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableFloat[%s](input, "%s")`, fieldName, strings.TrimPrefix(typeName, "*"), key))
	case "time.Time":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetDateTime(input, "%s")`, fieldName, key))
	case "*time.Time":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableDateTime(input, "%s")`, fieldName, key))
	case "github.com/google/uuid.UUID":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetObjectUUID(input, "%s")`, fieldName, key))
	case "*github.com/google/uuid.UUID":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableObjectUUID(input, "%s")`, fieldName, key))
	case "encoding/json.RawMessage":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetObjectRawJSON(input, "%s")`, fieldName, key))
	case "*encoding/json.RawMessage":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableObjectRawJSON(input, "%s")`, fieldName, key))
	case "any", "interface{}":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetArbitraryJSON(input, "%s")`, fieldName, key))
	case "*any", "*interface{}":
		sb.builder.WriteString(fmt.Sprintf(`  j.%s, err = utils.GetNullableArbitraryJSON(input, "%s")`, fieldName, key))
	default:
		if ty.IsNullable() {
			var pkgName, tyName string
			typeName := strings.TrimLeft(typeName, "*")
			scalarPackagePattern := fmt.Sprintf("%s.", sdkScalarPackageName)
			if strings.HasPrefix(typeName, scalarPackagePattern) {
				pkgName = sdkScalarPackageName
				tyName = fmt.Sprintf("scalar.%s", strings.TrimLeft(typeName, scalarPackagePattern))
			} else {
				pkgName, tyName = buildTypeNameFromFragments(ty.TypeFragments[1:], sb.packagePath)
			}
			if pkgName != "" {
				sb.imports[pkgName] = ""
			}
			sb.builder.WriteString(fmt.Sprintf(`  j.%s = new(%s)
  err = %s.DecodeNullableObjectValue(j.%s, input, "%s")`, fieldName, tyName, sb.GetDecoderName(), fieldName, key))
		} else {
			sb.builder.WriteString(fmt.Sprintf(`  err = %s.DecodeObjectValue(&j.%s, input, "%s")`, sb.GetDecoderName(), fieldName, key))
		}
	}
	sb.builder.WriteString(textBlockErrorCheck)
}

func buildTypeNameFromFragments(items []string, packagePath string) (string, string) {
	results := make([]string, len(items))
	importModule := ""
	for i, item := range items {
		if item == "*" || item == "[]" {
			results[i] = item
			continue
		}
		importModule, results[i] = extractPackageAndTypeName(item, packagePath)
	}

	return importModule, strings.Join(results, "")
}

func extractPackageAndTypeName(name string, packagePath string) (string, string) {
	if packagePath != "" {
		packagePath = fmt.Sprintf("%s.", packagePath)
	}

	if packagePath != "" && strings.HasPrefix(name, packagePath) {
		return "", strings.ReplaceAll(name, packagePath, "")
	}
	parts := strings.Split(name, "/")
	typeName := parts[len(parts)-1]
	typeNameParts := strings.Split(typeName, ".")
	if len(typeNameParts) < 2 {
		return "", typeName
	}
	if len(parts) == 1 {
		return typeNameParts[0], typeName
	}

	return strings.Join(append(parts[:len(parts)-1], typeNameParts[0]), "/"), typeName
}

func getSortedKeys[V any](input map[string]V) []string {
	var results []string
	for key := range input {
		results = append(results, key)
	}
	sort.Strings(results)
	return results
}

func formatLocalFieldName(input string, others ...string) string {
	name := fieldNameRegex.ReplaceAllString(input, "_")
	return strings.Trim(strings.Join(append([]string{name}, others...), "_"), "_")
}
