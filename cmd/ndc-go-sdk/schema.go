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
	"String":   *schema.NewScalarType(),
	"Int":      *schema.NewScalarType(),
	"BigInt":   *schema.NewScalarType(),
	"Float":    *schema.NewScalarType(),
	"Complex":  *schema.NewScalarType(),
	"Boolean":  *schema.NewScalarType(),
	"DateTime": *schema.NewScalarType(),
	"Duration": *schema.NewScalarType(),
	"UUID":     *schema.NewScalarType(),
}

var ndcOperationNameRegex = regexp.MustCompile(`^(Function|Procedure)([A-Z][A-Za-z0-9]*)$`)
var ndcOperationCommentRegex = regexp.MustCompile(`^@(function|procedure)(\s+([A-Za-z]\w*))?`)
var ndcScalarNameRegex = regexp.MustCompile(`^Scalar([A-Z]\w*)$`)
var ndcScalarCommentRegex = regexp.MustCompile(`^@scalar(\s+([A-Z]\w*))?`)

type OperationKind string

var (
	OperationFunction  OperationKind = "Function"
	OperationProcedure OperationKind = "Procedure"
)

// TypeInfo represents the serialization information of a type
type TypeInfo struct {
	Name        string
	SchemaName  string
	Description *string
	PackagePath string
	PackageName string
	IsScalar    bool
	IsNullable  bool
	IsArray     bool
	TypeAST     types.Type
	Schema      schema.TypeEncoder
}

// ObjectField represents the serialization information of an object field
type ObjectField struct {
	Name string
	Key  string
	Type *TypeInfo
}

// ObjectInfo represents the serialization information of an object type
type ObjectInfo struct {
	PackagePath string
	PackageName string
	IsAnonymous bool
	Fields      map[string]*ObjectField
}

// ArgumentInfo represents the serialization information of an argument type
type ArgumentInfo struct {
	FieldName   string
	Description *string
	Type        *TypeInfo
}

// Schema converts to ArgumentInfo schema
func (ai ArgumentInfo) Schema() schema.ArgumentInfo {
	return schema.ArgumentInfo{
		Description: ai.Description,
		Type:        ai.Type.Schema.Encode(),
	}
}

func buildArgumentInfosSchema(input map[string]ArgumentInfo) map[string]schema.ArgumentInfo {
	result := make(map[string]schema.ArgumentInfo)
	for k, arg := range input {
		result[k] = arg.Schema()
	}
	return result
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function or procedure schema
type OperationInfo struct {
	Kind          OperationKind
	Name          string
	OriginName    string
	PackageName   string
	Description   *string
	ArgumentsType string
	Arguments     map[string]ArgumentInfo
	ResultType    *TypeInfo
}

// FunctionInfo represents a readable Go function info
// which can convert to a NDC function schema
type FunctionInfo OperationInfo

// Schema returns a NDC function schema
func (op FunctionInfo) Schema() schema.FunctionInfo {
	result := schema.FunctionInfo{
		Name:        op.Name,
		Description: op.Description,
		ResultType:  op.ResultType.Schema.Encode(),
		Arguments:   buildArgumentInfosSchema(op.Arguments),
	}
	return result
}

// ProcedureInfo represents a readable Go function info
// which can convert to a NDC procedure schema
type ProcedureInfo FunctionInfo

// Schema returns a NDC procedure schema
func (op ProcedureInfo) Schema() schema.ProcedureInfo {
	result := schema.ProcedureInfo{
		Name:        op.Name,
		Description: op.Description,
		ResultType:  op.ResultType.Schema.Encode(),
		Arguments:   schema.ProcedureInfoArguments(buildArgumentInfosSchema(op.Arguments)),
	}
	return result
}

// RawConnectorSchema represents a readable Go schema object
// which can encode to NDC schema
type RawConnectorSchema struct {
	Imports       map[string]bool
	CustomScalars map[string]*TypeInfo
	ScalarSchemas schema.SchemaResponseScalarTypes
	Objects       map[string]*ObjectInfo
	ObjectSchemas schema.SchemaResponseObjectTypes
	Functions     []FunctionInfo
	Procedures    []ProcedureInfo
}

// NewRawConnectorSchema creates an empty RawConnectorSchema instance
func NewRawConnectorSchema() *RawConnectorSchema {
	return &RawConnectorSchema{
		Imports:       make(map[string]bool),
		CustomScalars: make(map[string]*TypeInfo),
		ScalarSchemas: make(schema.SchemaResponseScalarTypes),
		Objects:       make(map[string]*ObjectInfo),
		ObjectSchemas: make(schema.SchemaResponseObjectTypes),
		Functions:     []FunctionInfo{},
		Procedures:    []ProcedureInfo{},
	}
}

// Schema converts to a NDC schema
func (rcs RawConnectorSchema) Schema() *schema.SchemaResponse {
	result := &schema.SchemaResponse{
		ScalarTypes: rcs.ScalarSchemas,
		ObjectTypes: rcs.ObjectSchemas,
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

type SchemaParser struct {
	moduleName string
	files      map[string]*ast.File
	fset       *token.FileSet
}

func parseRawConnectorSchemaFromGoCode(moduleName string, filePath string, folders []string) (*RawConnectorSchema, error) {
	fset := token.NewFileSet()
	rawSchema := NewRawConnectorSchema()

	for _, folder := range folders {
		packages, err := parser.ParseDir(fset, path.Join(filePath, folder), func(fi fs.FileInfo) bool {
			return !fi.IsDir() && !strings.Contains(fi.Name(), "generated")
		}, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		files := make(map[string]*ast.File)
		for _, pkg := range packages {
			for n, pFiles := range pkg.Files {
				files[n] = pFiles
			}
		}
		sp := &SchemaParser{
			moduleName: moduleName,
			fset:       fset,
			files:      files,
		}

		if err = sp.checkAndParseRawSchemaFromAstFiles(rawSchema); err != nil {
			return nil, err
		}
	}

	return rawSchema, nil
}

func (sp *SchemaParser) checkAndParseRawSchemaFromAstFiles(rawSchema *RawConnectorSchema) error {

	conf := types.Config{
		Importer:                 importer.ForCompiler(sp.fset, "source", nil),
		IgnoreFuncBodies:         true,
		DisableUnusedImportCheck: true,
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	fileList := make([]*ast.File, 0, len(sp.files))
	for _, f := range sp.files {
		fileList = append(fileList, f)
	}
	pkg, err := conf.Check("", sp.fset, fileList, info)
	if err != nil {
		return err
	}
	rawSchema.Imports[fmt.Sprintf("%s/%s", sp.moduleName, pkg.Name())] = true
	return sp.parseRawConnectorSchema(rawSchema, pkg)
}

// parse raw connector schema from Go code
func (sp *SchemaParser) parseRawConnectorSchema(rawSchema *RawConnectorSchema, pkg *types.Package) error {

	for _, name := range pkg.Scope().Names() {
		switch obj := pkg.Scope().Lookup(name).(type) {
		case *types.Func:
			// only parse public functions
			if !obj.Exported() {
				continue
			}
			opInfo := sp.parseOperationInfo(obj.Name(), obj.Pos())
			if opInfo == nil {
				continue
			}
			opInfo.PackageName = pkg.Name()

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
				arguments, argumentTypeName, err := sp.parseArgumentTypes(rawSchema, arg.Type(), []string{})
				if err != nil {
					return err
				}
				opInfo.ArgumentsType = argumentTypeName
				opInfo.Arguments = arguments
			}

			resultType, err := sp.parseType(rawSchema, nil, resultTuple.At(0).Type(), []string{}, false)
			if err != nil {
				return err
			}
			opInfo.ResultType = resultType

			switch opInfo.Kind {
			case OperationProcedure:
				rawSchema.Procedures = append(rawSchema.Procedures, ProcedureInfo(*opInfo))
			case OperationFunction:
				rawSchema.Functions = append(rawSchema.Functions, FunctionInfo(*opInfo))
			}
		}
	}

	return nil
}

func (sp *SchemaParser) parseArgumentTypes(rawSchema *RawConnectorSchema, ty types.Type, fieldPaths []string) (map[string]ArgumentInfo, string, error) {

	switch inferredType := ty.(type) {
	case *types.Pointer:
		return sp.parseArgumentTypes(rawSchema, inferredType.Elem(), fieldPaths)
	case *types.Struct:
		result := make(map[string]ArgumentInfo)
		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)
			fieldPackage := fieldVar.Pkg()
			var typeInfo *TypeInfo
			if fieldPackage != nil {
				typeInfo = &TypeInfo{
					PackageName: fieldPackage.Name(),
					PackagePath: fieldPackage.Path(),
				}
			}
			fieldType, err := sp.parseType(rawSchema, typeInfo, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), false)
			if err != nil {
				return nil, "", err
			}
			fieldName := formatFieldName(fieldVar.Name(), fieldTag)
			if fieldType.TypeAST == nil {
				fieldType.TypeAST = fieldVar.Type()
			}
			result[fieldName] = ArgumentInfo{
				FieldName: fieldVar.Name(),
				Type:      fieldType,
			}
		}
		return result, "", nil
	case *types.Named:
		arguments, _, err := sp.parseArgumentTypes(rawSchema, inferredType.Obj().Type().Underlying(), append(fieldPaths, inferredType.Obj().Name()))
		if err != nil {
			return nil, "", err
		}
		return arguments, inferredType.Obj().Name(), nil
	default:
		return nil, "", fmt.Errorf("expected struct type, got %s", ty.String())
	}
}

func (sp *SchemaParser) parseType(rawSchema *RawConnectorSchema, rootType *TypeInfo, ty types.Type, fieldPaths []string, skipNullable bool) (*TypeInfo, error) {

	switch inferredType := ty.(type) {
	case *types.Pointer:
		if skipNullable {
			return sp.parseType(rawSchema, rootType, inferredType.Elem(), fieldPaths, false)
		}
		innerType, err := sp.parseType(rawSchema, rootType, inferredType.Elem(), fieldPaths, false)
		if err != nil {
			return nil, err
		}
		return &TypeInfo{
			Name:        innerType.Name,
			SchemaName:  innerType.Name,
			Description: innerType.Description,
			PackagePath: innerType.PackagePath,
			PackageName: innerType.PackageName,
			TypeAST:     ty,
			IsNullable:  true,
			Schema:      schema.NewNullableType(innerType.Schema),
		}, nil
	case *types.Struct:
		isAnonymous := false
		if rootType == nil {
			rootType = &TypeInfo{}
		}

		name := strings.Join(fieldPaths, "")
		if rootType.Name == "" {
			rootType.Name = name
			isAnonymous = true
		}
		if rootType.SchemaName == "" {
			rootType.SchemaName = name
		}
		if rootType.TypeAST == nil {
			rootType.TypeAST = ty
		}

		if rootType.Schema == nil {
			rootType.Schema = schema.NewNamedType(name)
		}
		objType := schema.ObjectType{
			Description: rootType.Description,
			Fields:      make(schema.ObjectTypeFields),
		}
		objFields := &ObjectInfo{
			PackagePath: rootType.PackagePath,
			PackageName: rootType.PackageName,
			IsAnonymous: isAnonymous,
			Fields:      map[string]*ObjectField{},
		}
		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)
			fieldType, err := sp.parseType(rawSchema, nil, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), false)
			if err != nil {
				return nil, err
			}
			fieldKey := formatFieldName(fieldVar.Name(), fieldTag)
			objType.Fields[fieldKey] = schema.ObjectField{
				Type: fieldType.Schema.Encode(),
			}
			objFields.Fields[fieldVar.Name()] = &ObjectField{
				Name: fieldVar.Name(),
				Key:  fieldKey,
				Type: fieldType,
			}
		}
		rawSchema.ObjectSchemas[rootType.Name] = objType
		rawSchema.Objects[rootType.Name] = objFields

		return rootType, nil
	case *types.Named:
		innerType := inferredType.Obj()
		if innerType == nil {
			return nil, fmt.Errorf("failed to parse named type: %s", inferredType.String())
		}
		typeInfo := sp.parseTypeInfoFromComments(innerType.Name(), innerType.Pos())

		innerPkg := innerType.Pkg()
		if innerPkg != nil {
			typeInfo.PackageName = innerPkg.Name()
			typeInfo.PackagePath = innerPkg.Path()
			var scalarName string
			switch innerPkg.Path() {
			case "time":
				switch innerType.Name() {
				case "Time":
					scalarName = "DateTime"
				case "Duration":
					scalarName = "Duration"
				}
			case "github.com/google/uuid":
				switch innerType.Name() {
				case "UUID":
					scalarName = "UUID"
				}
			}
			if scalarName != "" {
				rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
				typeInfo.Schema = schema.NewNamedType(scalarName)
				return typeInfo, nil
			}
		}

		if typeInfo.IsScalar {
			rawSchema.CustomScalars[typeInfo.SchemaName] = typeInfo
			rawSchema.ScalarSchemas[typeInfo.SchemaName] = *schema.NewScalarType()
			return typeInfo, nil
		}

		return sp.parseType(rawSchema, typeInfo, innerType.Type().Underlying(), append(fieldPaths, innerType.Name()), false)
	case *types.Basic:
		var scalarName string
		switch inferredType.Kind() {
		case types.Bool:
			scalarName = "Boolean"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		case types.Int8, types.Int, types.Int16, types.Int32, types.Uint, types.Uint16, types.Uint32:
			scalarName = "Int"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		case types.Int64, types.Uint64:
			scalarName = "BigInt"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		case types.Float32, types.Float64:
			scalarName = "Float"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		case types.Complex64, types.Complex128:
			scalarName = "Complex"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		case types.String:
			scalarName = "String"
			rawSchema.ScalarSchemas[scalarName] = defaultScalarTypes[scalarName]
		default:
			return nil, fmt.Errorf("unsupported scalar type: %s", inferredType.String())
		}
		if rootType == nil {
			rootType = &TypeInfo{
				Name:       inferredType.Name(),
				SchemaName: inferredType.Name(),
				TypeAST:    ty,
			}
		}

		rootType.Schema = schema.NewNamedType(scalarName)
		rootType.IsScalar = true

		return rootType, nil
	case *types.Array:
		innerType, err := sp.parseType(rawSchema, nil, inferredType.Elem(), fieldPaths, false)
		if err != nil {
			return nil, err
		}

		innerType.Schema = schema.NewArrayType(innerType.Schema)
		return innerType, nil
	case *types.Slice:
		innerType, err := sp.parseType(rawSchema, nil, inferredType.Elem(), fieldPaths, false)
		if err != nil {
			return nil, err
		}

		innerType.IsArray = true
		innerType.Schema = schema.NewArrayType(innerType.Schema)
		return innerType, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", ty.String())
	}
}

func (sp *SchemaParser) parseTypeInfoFromComments(typeName string, pos token.Pos) *TypeInfo {
	typeInfo := &TypeInfo{
		Name:       typeName,
		SchemaName: typeName,
		IsScalar:   false,
		Schema:     schema.NewNamedType(typeName),
	}
	comments := make([]string, 0)
	commentGroup := findCommentsFromPos(sp.fset, sp.files, pos)
	if commentGroup != nil {
		for _, line := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(line.Text, "/"))
			if text == "" {
				continue
			}
			matches := ndcScalarCommentRegex.FindStringSubmatch(text)
			matchesLen := len(matches)
			if matchesLen < 1 {
				comments = append(comments, text)
				continue
			}

			typeInfo.IsScalar = true
			if matchesLen > 2 && matches[2] != "" {
				typeInfo.SchemaName = matches[2]
				typeInfo.Schema = schema.NewNamedType(matches[2])
			}
		}
	}

	if !typeInfo.IsScalar {
		// fallback to parse scalar from type name with Scalar prefix
		matches := ndcScalarNameRegex.FindStringSubmatch(typeName)
		if len(matches) > 1 {
			typeInfo.IsScalar = true
			typeInfo.SchemaName = matches[1]
			typeInfo.Schema = schema.NewNamedType(matches[1])
		}
	}

	desc := strings.Join(comments, " ")
	if desc != "" {
		typeInfo.Description = &desc
	}

	return typeInfo
}

func (sp *SchemaParser) parseOperationInfo(functionName string, pos token.Pos) *OperationInfo {

	result := OperationInfo{
		OriginName: functionName,
		Arguments:  make(map[string]ArgumentInfo),
	}

	var descriptions []string
	commentGroup := findCommentsFromPos(sp.fset, sp.files, pos)
	if commentGroup != nil {
		for _, comment := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))
			matches := ndcOperationCommentRegex.FindStringSubmatch(text)
			matchesLen := len(matches)
			if matchesLen > 1 {
				switch matches[1] {
				case strings.ToLower(string(OperationFunction)):
					result.Kind = OperationFunction
				case strings.ToLower(string(OperationProcedure)):
					result.Kind = OperationProcedure
				default:
					log.Debug().Msgf("unsupported operation kind: %s", matches)
				}

				if matchesLen > 3 && strings.TrimSpace(matches[3]) != "" {
					result.Name = strings.TrimSpace(matches[3])
				} else {
					result.Name = camelCase(functionName)
				}
			} else {
				descriptions = append(descriptions, text)
			}
		}
	}

	// try to parse function with following prefixes:
	// - FunctionXxx as a query function
	// - ProcedureXxx as a mutation procedure
	if result.Kind == "" {
		operationNameResults := ndcOperationNameRegex.FindStringSubmatch(functionName)
		if len(operationNameResults) < 3 {
			return nil
		}
		result.Kind = OperationKind(operationNameResults[1])
		result.Name = camelCase(operationNameResults[2])
	}

	desc := strings.TrimSpace(strings.Join(descriptions, " "))
	if desc != "" {
		result.Description = &desc
	}

	return &result
}

func findCommentsFromPos(fset *token.FileSet, files map[string]*ast.File, pos token.Pos) *ast.CommentGroup {
	for fName, f := range files {
		position := fset.Position(pos)
		if position.Filename != fName {
			continue
		}
		offset := pos - 10
		for _, cg := range f.Comments {
			if len(cg.List) > 0 && (cg.Pos() <= token.Pos(offset) && cg.End() >= token.Pos(offset)) {
				return cg
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
