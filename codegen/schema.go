package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/structtag"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog/log"
)

var defaultScalarTypes = schema.SchemaResponseScalarTypes{
	"String": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Int": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Float": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
	"Boolean": schema.ScalarType{
		AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
	},
}

var ndcOperationNameRegex = regexp.MustCompile(`^(Function|Procedure)([A-Z][A-Za-z0-9]*)$`)

type OperationKind string

var (
	OperationFunction  OperationKind = "Function"
	OperationProcedure OperationKind = "Procedure"
)

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function or procedure schema
type OperationInfo struct {
	Kind          OperationKind
	Name          string
	OriginName    string
	PackageName   string
	Description   string
	ArgumentsType string
	Arguments     schema.FunctionInfoArguments
	ResultType    schema.TypeEncoder
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function schema
type FunctionInfo OperationInfo

// Schema returns a NDC function schema
func (op FunctionInfo) Schema() schema.FunctionInfo {
	result := schema.FunctionInfo{
		Name:       op.Name,
		ResultType: op.ResultType.Encode(),
		Arguments:  op.Arguments,
	}
	if op.Description != "" {
		result.Description = &op.Description
	}
	return result
}

// ProcedureInfo represents a readable Go function info
// which can convert to a NDC procedure schema
type ProcedureInfo FunctionInfo

// Schema returns a NDC procedure schema
func (op ProcedureInfo) Schema() schema.ProcedureInfo {
	result := schema.ProcedureInfo{
		Name:       op.Name,
		ResultType: op.ResultType.Encode(),
		Arguments:  schema.ProcedureInfoArguments(op.Arguments),
	}
	if op.Description != "" {
		result.Description = &op.Description
	}
	return result
}

// RawConnectorSchema represents a readable Go schema object
// which can encode to NDC schema
type RawConnectorSchema struct {
	PackageNames []string
	Scalars      schema.SchemaResponseScalarTypes
	Objects      schema.SchemaResponseObjectTypes
	Functions    []FunctionInfo
	Procedures   []ProcedureInfo
}

// Schema converts to a NDC schema
func (rcs RawConnectorSchema) Schema() *schema.SchemaResponse {
	result := &schema.SchemaResponse{
		ScalarTypes: rcs.Scalars,
		ObjectTypes: rcs.Objects,
		Collections: []schema.CollectionInfo{},
	}
	for _, function := range rcs.Functions {
		result.Functions = append(result.Functions, function.Schema())
	}
	for _, procedure := range rcs.Procedures {
		result.Procedures = append(result.Procedures, procedure.Schema())
	}

	return result
}

func parseRawConnectorSchemaFromGoCode(filePath string, folders []string) (*RawConnectorSchema, error) {
	fset := token.NewFileSet()
	rawSchema := &RawConnectorSchema{
		PackageNames: make([]string, 0),
		Scalars:      make(schema.SchemaResponseScalarTypes),
		Objects:      make(schema.SchemaResponseObjectTypes),
		Functions:    []FunctionInfo{},
		Procedures:   []ProcedureInfo{},
	}

	for _, folder := range folders {
		packages, err := parser.ParseDir(fset, path.Join(filePath, folder), func(fi fs.FileInfo) bool {
			return !fi.IsDir() && !strings.Contains(fi.Name(), "generated")
		}, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		var files []*ast.File
		for _, pkg := range packages {
			for _, file := range pkg.Files {
				files = append(files, file)
			}
		}
		conf := types.Config{
			Importer:                 importer.ForCompiler(fset, "source", nil),
			IgnoreFuncBodies:         true,
			DisableUnusedImportCheck: true,
		}
		info := &types.Info{}
		pkg, err := conf.Check("", fset, files, info)
		if err != nil {
			return nil, err
		}

		rawSchema.PackageNames = append(rawSchema.PackageNames, pkg.Name())
		err = parseRawConnectorSchema(rawSchema, pkg, files)
		if err != nil {
			return nil, err
		}
	}

	return rawSchema, nil
}

// parse raw connector schema from Go code
func parseRawConnectorSchema(rawSchema *RawConnectorSchema, pkg *types.Package, files []*ast.File) error {

	for _, name := range pkg.Scope().Names() {
		switch obj := pkg.Scope().Lookup(name).(type) {
		case *types.Func:
			// parse function with following prefixes:
			// - FunctionXxx as a query function
			// - ProcedureXxx as a mutation procedure
			operationNameResults := ndcOperationNameRegex.FindStringSubmatch(obj.Name())
			if len(operationNameResults) < 3 {
				continue
			}
			opInfo := OperationInfo{
				Kind:        OperationKind(operationNameResults[1]),
				Name:        camelCase(operationNameResults[2]),
				OriginName:  operationNameResults[0],
				PackageName: pkg.Name(),
				Arguments:   make(schema.FunctionInfoArguments),
			}

			var resultTuple *types.Tuple
			var params *types.Tuple
			switch sig := obj.Type().(type) {
			case *types.Signature:
				params = sig.Params()
				resultTuple = sig.Results()
			default:
				return fmt.Errorf("expected function signature, got: %s", sig.String())
			}

			if params == nil || (params.Len() < 2 || params.Len() > 3) {
				return fmt.Errorf("%s: expect 2 or 3 parameters only (ctx context.Context, state types.State, arguments *[ArgumentType]), got %s", opInfo.OriginName, params)
			}

			if resultTuple == nil || resultTuple.Len() != 2 {
				return fmt.Errorf("%s: expect result tuple ([type], error), got %s", opInfo.OriginName, resultTuple)
			}

			// parse arguments in the function if exists
			// ignore 2 first parameters (context and state)
			if params.Len() == 3 {
				arg := params.At(2)
				arguments, argumentTypeName, err := parseArgumentTypes(rawSchema, arg.Type())
				if err != nil {
					return err
				}
				opInfo.ArgumentsType = argumentTypeName
				opInfo.Arguments = arguments
			}

			resultType, err := parseType(rawSchema, resultTuple.At(0).Type(), false)
			if err != nil {
				return err
			}
			opInfo.ResultType = resultType

			switch opInfo.Kind {
			case OperationProcedure:
				rawSchema.Procedures = append(rawSchema.Procedures, ProcedureInfo(opInfo))
			case OperationFunction:
				rawSchema.Functions = append(rawSchema.Functions, FunctionInfo(opInfo))
			}
		}
	}

	return nil
}

func camelCase(input string) string {
	if len(input) == 0 {
		return ""
	}
	return strings.ToLower(input[0:1]) + input[1:]
}

// format field name by json tag
// return the struct field name if not exist
func formatFieldName(name string, tag string) string {
	if tag == "" {
		return name
	}
	tags, err := structtag.Parse(tag)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to parse tag of struct field: %s", name)
		return name
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		log.Warn().Err(err).Msgf("json tag does not exist in struct field: %s", name)
		return name
	}

	return jsonTag.Name
}

func parseArgumentTypes(rawSchema *RawConnectorSchema, ty types.Type) (map[string]schema.ArgumentInfo, string, error) {
	switch inferredType := ty.(type) {
	case *types.Pointer:
		return parseArgumentTypes(rawSchema, inferredType.Elem())
	case *types.Struct:
		result := make(map[string]schema.ArgumentInfo)
		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)
			fieldType, err := parseType(rawSchema, fieldVar.Type(), false)
			if err != nil {
				return nil, "", err
			}
			fieldName := formatFieldName(fieldVar.Name(), fieldTag)
			result[fieldName] = schema.ArgumentInfo{
				Type: fieldType.Encode(),
			}
		}
		return result, "", nil
	case *types.Named:
		arguments, _, err := parseArgumentTypes(rawSchema, inferredType.Obj().Type().Underlying())
		if err != nil {
			return nil, "", err
		}
		return arguments, inferredType.Obj().Name(), nil
	default:
		return nil, "", fmt.Errorf("expected struct type, got %s", ty.String())
	}
}

func parseStructType(rawSchema *RawConnectorSchema, name string, ty types.Type) error {
	inferredType, ok := ty.(*types.Struct)
	if !ok {
		return fmt.Errorf("expected struct type, got %s", ty.String())
	}
	objType := schema.ObjectType{
		Fields: make(schema.ObjectTypeFields),
	}
	for i := 0; i < inferredType.NumFields(); i++ {
		fieldVar := inferredType.Field(i)
		fieldTag := inferredType.Tag(i)
		fieldType, err := parseType(rawSchema, fieldVar.Type(), false)
		if err != nil {
			return err
		}
		fieldName := formatFieldName(fieldVar.Name(), fieldTag)
		objType.Fields[fieldName] = schema.ObjectField{
			Type: fieldType.Encode(),
		}
	}
	rawSchema.Objects[name] = objType
	return nil
}

func parseType(rawSchema *RawConnectorSchema, ty types.Type, skipNullable bool) (schema.TypeEncoder, error) {

	switch inferredType := ty.(type) {
	case *types.Pointer:
		if skipNullable {
			return parseType(rawSchema, inferredType.Elem(), false)
		}
		innerType, err := parseType(rawSchema, inferredType.Elem(), false)
		if err != nil {
			return nil, err
		}
		return schema.NewNullableType(innerType), nil
	case *types.Struct:
		return schema.NewNamedType(inferredType.String()), nil
	case *types.Named:
		innerType := inferredType.Obj()
		if innerType != nil {
			// recursively parse object types
			err := parseStructType(rawSchema, innerType.Name(), innerType.Type().Underlying())
			if err != nil {
				return nil, err
			}
		}
		return schema.NewNamedType(innerType.Name()), nil
	case *types.Basic:
		var scalarName string
		switch inferredType.Info() {
		case types.IsBoolean:
			scalarName = "Boolean"
			rawSchema.Scalars[scalarName] = defaultScalarTypes[scalarName]
		case types.IsInteger, types.IsUnsigned:
			scalarName = "Int"
			rawSchema.Scalars[scalarName] = defaultScalarTypes[scalarName]
		case types.IsFloat, types.IsComplex:
			scalarName = "Float"
			rawSchema.Scalars[scalarName] = defaultScalarTypes[scalarName]
		case types.IsString:
			scalarName = "String"
			rawSchema.Scalars[scalarName] = defaultScalarTypes[scalarName]
		default:
			return nil, fmt.Errorf("unsupported scalar type: %s", inferredType.String())
		}
		return schema.NewNamedType(scalarName), nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", ty.String())
	}
}

// func parseOperationInfoFromComment(functionName string, comments []*ast.Comment) *OperationInfo {
// 	var result OperationInfo
// 	var descriptions []string
// 	for _, comment := range comments {
// 		text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))
// 		matches := operationNameCommentRegex.FindStringSubmatch(text)
// 		matchesLen := len(matches)
// 		if matchesLen > 1 {
// 			switch matches[1] {
// 			case string(OperationFunction):
// 				result.Kind = OperationFunction
// 			case string(OperationProcedure):
// 				result.Kind = OperationProcedure
// 			default:
// 				log.Println("unsupported operation kind:", matches[0])
// 			}

// 			if matchesLen > 3 && strings.TrimSpace(matches[3]) != "" {
// 				result.Name = strings.TrimSpace(matches[3])
// 			} else {
// 				result.Name = strings.ToLower(functionName[:1]) + functionName[1:]
// 			}
// 		} else {
// 			descriptions = append(descriptions, text)
// 		}
// 	}

// 	if result.Kind == "" {
// 		return nil
// 	}

// 	result.Description = strings.TrimSpace(strings.Join(descriptions, " "))
// 	return &result
// }
