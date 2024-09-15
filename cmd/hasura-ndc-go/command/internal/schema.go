package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/trace"
	"strings"

	"github.com/fatih/structtag"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog/log"
	"golang.org/x/tools/go/packages"
)

type OperationKind string

const (
	OperationFunction  OperationKind = "Function"
	OperationProcedure OperationKind = "Procedure"
)

type TypeKind string

// TypeInfo represents the serialization information of a type
type TypeInfo struct {
	Name                 string
	SchemaName           string
	Description          *string
	PackagePath          string
	PackageName          string
	IsScalar             bool
	ScalarRepresentation schema.TypeRepresentation
	TypeFragments        []string
	TypeAST              types.Type
	Schema               schema.TypeEncoder
}

// IsNullable checks if the current type is nullable
func (ti *TypeInfo) IsNullable() bool {
	return isNullableFragments(ti.TypeFragments)
}

func isNullableFragment(fragment string) bool {
	return fragment == "*"
}

func isNullableFragments(fragments []string) bool {
	return len(fragments) > 0 && isNullableFragment(fragments[0])
}

// IsArray checks if the current type is an array
func (ti *TypeInfo) IsArray() bool {
	return isArrayFragments(ti.TypeFragments)
}

func isArrayFragment(fragment string) bool {
	return fragment == "[]"
}

func isArrayFragments(fragments []string) bool {
	return len(fragments) > 0 && isArrayFragment(fragments[0])
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
	PackagePath   string
	Description   *string
	ArgumentsType *TypeInfo
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
		Functions:   []schema.FunctionInfo{},
		Procedures:  []schema.ProcedureInfo{},
	}
	for _, function := range rcs.Functions {
		result.Functions = append(result.Functions, function.Schema())
	}
	for _, procedure := range rcs.Procedures {
		result.Procedures = append(result.Procedures, procedure.Schema())
	}

	return result
}

// IsCustomType checks if the type name is a custom scalar or an exported object
func (rcs RawConnectorSchema) IsCustomType(name string) bool {
	if _, ok := rcs.CustomScalars[name]; ok {
		return true
	}
	if obj, ok := rcs.Objects[name]; ok {
		return !obj.IsAnonymous
	}
	return false
}

type SchemaParser struct {
	context      context.Context
	moduleName   string
	rawSchema    *RawConnectorSchema
	packages     []*packages.Package
	packageIndex int
}

// GetCurrentPackage gets the current evaluating package
func (sp SchemaParser) GetCurrentPackage() *packages.Package {
	return sp.packages[sp.packageIndex]
}

// FindPackageByPath finds the package by package path
func (sp SchemaParser) FindPackageByPath(input string) *packages.Package {
	for _, p := range sp.packages {
		if p.ID == input {
			return p
		}
	}
	return nil
}

func parseRawConnectorSchemaFromGoCode(ctx context.Context, moduleName string, filePath string, args *ConnectorGenerationArguments) (*RawConnectorSchema, error) {
	rawSchema := NewRawConnectorSchema()

	pkgTypes, err := evalPackageTypesLocation(args.PackageTypes, moduleName, filePath, args.ConnectorDir)
	if err != nil {
		return nil, err
	}
	rawSchema.Imports[pkgTypes] = true

	tempDirs := args.Directories
	if len(args.Directories) == 0 {
		// recursively walk directories if the user don't explicitly specify target folders
		entries, err := os.ReadDir(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read subdirectories of %s: %s", filePath, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			tempDirs = append(tempDirs, entry.Name())
		}
	}
	var directories []string
	for _, dir := range tempDirs {
		for _, globPath := range []string{path.Join(filePath, dir, "*.go"), path.Join(filePath, dir, "**", "*.go")} {
			goFiles, err := filepath.Glob(globPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read subdirectories of %s/%s: %s", filePath, dir, err)
			}
			// cleanup types.generated.go files
			fileCount := 0
			for _, fp := range goFiles {
				if !strings.HasSuffix(fp, typeMethodsOutputFile) {
					fileCount++
					continue
				}
				if err := os.Remove(fp); err != nil {
					return nil, fmt.Errorf("failed to delete %s: %s", fp, err)
				}
			}
			if fileCount > 0 {
				directories = append(directories, dir)
				break
			}
		}
	}

	if len(directories) == 0 {
		log.Info().Msgf("no subdirectory in %s", filePath)
		return rawSchema, nil
	}

	log.Info().Interface("directories", directories).Msgf("parsing connector schema...")

	var packageList []*packages.Package
	fset := token.NewFileSet()
	for _, folder := range directories {
		_, parseCodeTask := trace.NewTask(ctx, fmt.Sprintf("parse_%s_code", folder))
		folderPath := path.Join(filePath, folder)

		cfg := &packages.Config{
			Mode: packages.NeedSyntax | packages.NeedTypes,
			Dir:  folderPath,
			Fset: fset,
		}
		pkgList, err := packages.Load(cfg, flag.Args()...)
		parseCodeTask.End()
		if err != nil {
			return nil, err
		}
		packageList = append(packageList, pkgList...)
	}

	for i := range packageList {
		parseSchemaCtx, parseSchemaTask := trace.NewTask(ctx, fmt.Sprintf("parse_schema_%s", packageList[i].ID))
		sp := &SchemaParser{
			context:      parseSchemaCtx,
			moduleName:   moduleName,
			packages:     packageList,
			packageIndex: i,
			rawSchema:    rawSchema,
		}

		err = sp.parseRawConnectorSchema(packageList[i].Types)
		parseSchemaTask.End()
		if err != nil {
			return nil, err
		}
	}

	return rawSchema, nil
}

func evalPackageTypesLocation(name string, moduleName string, filePath string, connectorDir string) (string, error) {
	if name != "" {
		// assume that the absolute package name should have domain, e.g. github.com/...
		if strings.Contains(name, ".") {
			return name, nil
		}
		return fmt.Sprintf("%s/%s", moduleName, name), nil
	}

	matches, err := filepath.Glob(path.Join(filePath, "types", "*.go"))
	if err == nil && len(matches) > 0 {
		return fmt.Sprintf("%s/types", moduleName), nil
	}

	if connectorDir != "" && !strings.HasPrefix(".", connectorDir) {
		matches, err = filepath.Glob(path.Join(filePath, connectorDir, "types", "*.go"))
		if err == nil && len(matches) > 0 {
			return fmt.Sprintf("%s/%s/types", moduleName, connectorDir), nil
		}
	}
	return "", fmt.Errorf("the `types` package where the State struct is in must be placed in root or connector directory, %s", err)
}

// parse raw connector schema from Go code
func (sp *SchemaParser) parseRawConnectorSchema(pkg *types.Package) error {

	for _, name := range pkg.Scope().Names() {
		_, task := trace.NewTask(sp.context, fmt.Sprintf("parse_%s_schema_%s", sp.GetCurrentPackage().Name, name))
		err := sp.parsePackageScope(pkg, name)
		task.End()
		if err != nil {
			return err
		}
	}

	return nil
}

func (sp *SchemaParser) parsePackageScope(pkg *types.Package, name string) error {
	switch obj := pkg.Scope().Lookup(name).(type) {
	case *types.Func:
		// only parse public functions
		if !obj.Exported() {
			return nil
		}
		opInfo := sp.parseOperationInfo(obj)
		if opInfo == nil {
			return nil
		}
		opInfo.PackageName = pkg.Name()
		opInfo.PackagePath = pkg.Path()
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
			arguments, argumentType, err := sp.parseArgumentTypes(arg.Type(), []string{})
			if err != nil {
				return err
			}
			opInfo.ArgumentsType = argumentType
			opInfo.Arguments = arguments
		}

		resultType, err := sp.parseType(nil, resultTuple.At(0).Type(), []string{}, false, false)
		if err != nil {
			return err
		}
		opInfo.ResultType = resultType

		switch opInfo.Kind {
		case OperationProcedure:
			sp.rawSchema.Procedures = append(sp.rawSchema.Procedures, ProcedureInfo(*opInfo))
		case OperationFunction:
			sp.rawSchema.Functions = append(sp.rawSchema.Functions, FunctionInfo(*opInfo))
		}
	}
	return nil
}

func (sp *SchemaParser) parseArgumentTypes(ty types.Type, fieldPaths []string) (map[string]ArgumentInfo, *TypeInfo, error) {

	switch inferredType := ty.(type) {
	case *types.Pointer:
		return sp.parseArgumentTypes(inferredType.Elem(), fieldPaths)
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
			fieldType, err := sp.parseType(typeInfo, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), false, true)
			if err != nil {
				return nil, nil, err
			}
			fieldName := getFieldNameOrTag(fieldVar.Name(), fieldTag)
			if fieldType.TypeAST == nil {
				fieldType.TypeAST = fieldVar.Type()
			}
			result[fieldName] = ArgumentInfo{
				FieldName: fieldVar.Name(),
				Type:      fieldType,
			}
		}
		return result, nil, nil
	case *types.Named:
		arguments, _, err := sp.parseArgumentTypes(inferredType.Obj().Type().Underlying(), append(fieldPaths, inferredType.Obj().Name()))
		if err != nil {
			return nil, nil, err
		}

		typeObj := inferredType.Obj()
		typeInfo := &TypeInfo{
			Name:       typeObj.Name(),
			SchemaName: typeObj.Name(),
		}
		pkg := typeObj.Pkg()
		if pkg != nil {
			typeInfo.PackagePath = pkg.Path()
			typeInfo.PackageName = pkg.Name()
		}
		return arguments, typeInfo, nil
	default:
		return nil, nil, fmt.Errorf("expected struct type, got %s", ty.String())
	}
}

func (sp *SchemaParser) parseType(rootType *TypeInfo, ty types.Type, fieldPaths []string, skipNullable bool, isArgument bool) (*TypeInfo, error) {

	switch inferredType := ty.(type) {
	case *types.Pointer:
		if skipNullable {
			return sp.parseType(rootType, inferredType.Elem(), fieldPaths, false, isArgument)
		}
		innerType, err := sp.parseType(rootType, inferredType.Elem(), fieldPaths, false, isArgument)
		if err != nil {
			return nil, err
		}
		return &TypeInfo{
			Name:          innerType.Name,
			SchemaName:    innerType.Name,
			Description:   innerType.Description,
			PackagePath:   innerType.PackagePath,
			PackageName:   innerType.PackageName,
			TypeAST:       ty,
			TypeFragments: append([]string{"*"}, innerType.TypeFragments...),
			IsScalar:      innerType.IsScalar,
			Schema:        schema.NewNullableType(innerType.Schema),
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
			rootType.TypeFragments = append(rootType.TypeFragments, ty.String())
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
		// temporarily add the object type to raw schema to avoid infinite loop
		sp.rawSchema.ObjectSchemas[rootType.Name] = objType
		sp.rawSchema.Objects[rootType.Name] = objFields

		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)
			fieldType, err := sp.parseType(nil, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), false, isArgument)
			if err != nil {
				return nil, err
			}
			fieldKey := getFieldNameOrTag(fieldVar.Name(), fieldTag)
			objType.Fields[fieldKey] = schema.ObjectField{
				Type: fieldType.Schema.Encode(),
			}
			objFields.Fields[fieldVar.Name()] = &ObjectField{
				Name: fieldVar.Name(),
				Key:  fieldKey,
				Type: fieldType,
			}
		}
		sp.rawSchema.ObjectSchemas[rootType.Name] = objType
		sp.rawSchema.Objects[rootType.Name] = objFields

		return rootType, nil
	case *types.Named:

		innerType := inferredType.Obj()
		if innerType == nil {
			return nil, fmt.Errorf("failed to parse named type: %s", inferredType.String())
		}

		innerPkg := innerType.Pkg()
		var typeInfo *TypeInfo
		if innerPkg != nil {
			if _, ok := sp.rawSchema.Objects[innerType.Name()]; ok {
				ty := &TypeInfo{
					Name:          innerType.Name(),
					SchemaName:    innerType.Name(),
					PackageName:   innerPkg.Name(),
					PackagePath:   innerPkg.Path(),
					TypeAST:       innerType.Type(),
					Schema:        schema.NewNamedType(innerType.Name()),
					TypeFragments: []string{innerType.Name()},
				}
				return ty, nil
			}

			var err error
			typeInfo, err = sp.parseTypeInfoFromComments(innerType.Name(), innerPkg.Path(), innerType.Parent())
			if err != nil {
				return nil, err
			}
			var scalarName ScalarName
			typeInfo.PackageName = innerPkg.Name()
			typeInfo.PackagePath = innerPkg.Path()
			scalarSchema := schema.NewScalarType()

			switch innerPkg.Path() {
			case "time":
				switch innerType.Name() {
				case "Time":
					scalarName = ScalarTimestampTZ
					scalarSchema.Representation = schema.NewTypeRepresentationTimestampTZ().Encode()
				case "Duration":
					return nil, errors.New("unsupported type time.Duration. Create a scalar type wrapper with FromValue method to decode the any value")
				}
			case "encoding/json":
				switch innerType.Name() {
				case "RawMessage":
					scalarName = ScalarRawJSON
					scalarSchema.Representation = schema.NewTypeRepresentationJSON().Encode()
				}
			case "github.com/google/uuid":
				switch innerType.Name() {
				case "UUID":
					scalarName = ScalarUUID
					scalarSchema.Representation = schema.NewTypeRepresentationUUID().Encode()
				}
			case "github.com/hasura/ndc-sdk-go/scalar":
				scalarName = ScalarName(innerType.Name())
				switch innerType.Name() {
				case "Date":
					scalarSchema.Representation = schema.NewTypeRepresentationDate().Encode()
				case "BigInt":
					scalarSchema.Representation = schema.NewTypeRepresentationBigInteger().Encode()
				case "Bytes":
					scalarSchema.Representation = schema.NewTypeRepresentationBytes().Encode()
				case "URL":
					scalarSchema.Representation = schema.NewTypeRepresentationString().Encode()
				}
			}

			if scalarName != "" {
				typeInfo.IsScalar = true
				typeInfo.Schema = schema.NewNamedType(string(scalarName))
				typeInfo.TypeAST = ty
				sp.rawSchema.ScalarSchemas[string(scalarName)] = *scalarSchema
				return typeInfo, nil
			}
		} else if innerType.Name() == "error" {
			if isArgument {
				return nil, fmt.Errorf("%s: native `error` interface isn't allowed in input arguments", strings.Join(fieldPaths, "."))
			}
			typeInfo = &TypeInfo{
				Name:                 innerType.Name(),
				SchemaName:           string(ScalarJSON),
				TypeAST:              innerType.Type(),
				Schema:               schema.NewNullableType(schema.NewNamedType(string(ScalarJSON))),
				IsScalar:             true,
				ScalarRepresentation: schema.NewTypeRepresentationJSON().Encode(),
			}
			if _, ok := sp.rawSchema.ScalarSchemas[typeInfo.SchemaName]; !ok {
				sp.rawSchema.ScalarSchemas[typeInfo.SchemaName] = defaultScalarTypes[ScalarJSON]
			}
			return typeInfo, nil
		} else {
			return nil, fmt.Errorf("%s: unsupported type <%s>", strings.Join(fieldPaths, "."), innerType.Name())
		}

		if typeInfo.IsScalar {
			sp.rawSchema.CustomScalars[typeInfo.Name] = typeInfo
			scalarSchema := schema.NewScalarType()
			if typeInfo.ScalarRepresentation != nil {
				scalarSchema.Representation = typeInfo.ScalarRepresentation
			} else {
				// requires representation since NDC spec v0.1.2
				scalarSchema.Representation = schema.NewTypeRepresentationJSON().Encode()
			}
			sp.rawSchema.ScalarSchemas[typeInfo.SchemaName] = *scalarSchema
			return typeInfo, nil
		}

		return sp.parseType(typeInfo, innerType.Type().Underlying(), append(fieldPaths, innerType.Name()), false, isArgument)
	case *types.Basic:
		var scalarName ScalarName
		switch inferredType.Kind() {
		case types.Bool:
			scalarName = ScalarBoolean
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Int8, types.Uint8:
			scalarName = ScalarInt8
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Int16, types.Uint16:
			scalarName = ScalarInt16
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Int, types.Int32, types.Uint, types.Uint32:
			scalarName = ScalarInt32
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Int64, types.Uint64:
			scalarName = ScalarInt64
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Float32:
			scalarName = ScalarFloat32
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.Float64:
			scalarName = ScalarFloat64
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		case types.String:
			scalarName = ScalarString
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		default:
			return nil, fmt.Errorf("%s: unsupported scalar type <%s>", strings.Join(fieldPaths, "."), inferredType.String())
		}
		if rootType == nil {
			rootType = &TypeInfo{
				Name:          inferredType.Name(),
				SchemaName:    inferredType.Name(),
				TypeFragments: []string{inferredType.Name()},
				TypeAST:       ty,
			}
		}

		rootType.Schema = schema.NewNamedType(string(scalarName))
		rootType.IsScalar = true

		return rootType, nil
	case *types.Array:
		innerType, err := sp.parseType(nil, inferredType.Elem(), fieldPaths, false, isArgument)
		if err != nil {
			return nil, err
		}
		innerType.TypeFragments = append([]string{"[]"}, innerType.TypeFragments...)
		innerType.Schema = schema.NewArrayType(innerType.Schema)
		return innerType, nil
	case *types.Slice:
		innerType, err := sp.parseType(nil, inferredType.Elem(), fieldPaths, false, isArgument)
		if err != nil {
			return nil, err
		}

		innerType.TypeFragments = append([]string{"[]"}, innerType.TypeFragments...)
		innerType.Schema = schema.NewArrayType(innerType.Schema)
		return innerType, nil
	case *types.Map, *types.Interface:
		scalarName := ScalarJSON
		if rootType == nil {
			rootType = &TypeInfo{
				Name:       inferredType.String(),
				SchemaName: string(scalarName),
				TypeAST:    ty,
			}
		} else {
			rootType.PackagePath = ""
		}

		if _, ok := sp.rawSchema.ScalarSchemas[string(scalarName)]; !ok {
			sp.rawSchema.ScalarSchemas[string(scalarName)] = defaultScalarTypes[scalarName]
		}
		rootType.TypeFragments = append(rootType.TypeFragments, inferredType.String())
		rootType.Schema = schema.NewNamedType(string(scalarName))
		rootType.IsScalar = true

		return rootType, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", ty.String())
	}
}

func (sp *SchemaParser) parseTypeInfoFromComments(typeName string, packagePath string, scope *types.Scope) (*TypeInfo, error) {
	typeInfo := &TypeInfo{
		Name:          typeName,
		SchemaName:    typeName,
		IsScalar:      false,
		TypeFragments: []string{typeName},
		Schema:        schema.NewNamedType(typeName),
	}
	comments := make([]string, 0)
	commentGroup := findCommentsFromPos(sp.FindPackageByPath(packagePath), scope, typeName)
	if commentGroup != nil {
		for i, line := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(line.Text, "/"))
			if text == "" {
				continue
			}
			if i == 0 {
				text = strings.TrimPrefix(text, fmt.Sprintf("%s ", typeName))
			}

			enumMatches := ndcEnumCommentRegex.FindStringSubmatch(text)

			if len(enumMatches) == 2 {
				typeInfo.IsScalar = true
				rawEnumItems := strings.Split(enumMatches[1], ",")
				var enums []string
				for _, item := range rawEnumItems {
					trimmed := strings.TrimSpace(item)
					if trimmed != "" {
						enums = append(enums, trimmed)
					}
				}
				if len(enums) == 0 {
					return nil, fmt.Errorf("require enum values in the comment of %s", typeName)
				}
				typeInfo.ScalarRepresentation = schema.NewTypeRepresentationEnum(enums).Encode()
				continue
			}

			matches := ndcScalarCommentRegex.FindStringSubmatch(text)
			matchesLen := len(matches)
			if matchesLen > 1 {
				typeInfo.IsScalar = true
				if matchesLen > 3 && matches[3] != "" {
					typeInfo.SchemaName = matches[2]
					typeInfo.Schema = schema.NewNamedType(matches[2])
					typeRep, err := schema.ParseTypeRepresentationType(strings.TrimSpace(matches[3]))
					if err != nil {
						return nil, fmt.Errorf("failed to parse type representation of scalar %s: %s", typeName, err)
					}
					if typeRep == schema.TypeRepresentationTypeEnum {
						return nil, errors.New("use @enum tag with values instead")
					}
					typeInfo.ScalarRepresentation = schema.TypeRepresentation{
						"type": typeRep,
					}
				} else if matchesLen > 2 && matches[2] != "" {
					// if the second string is a type representation, use it as a TypeRepresentation instead
					// e.g @scalar string
					typeRep, err := schema.ParseTypeRepresentationType(matches[2])
					if err == nil {
						if typeRep == schema.TypeRepresentationTypeEnum {
							return nil, errors.New("use @enum tag with values instead")
						}
						typeInfo.ScalarRepresentation = schema.TypeRepresentation{
							"type": typeRep,
						}
						continue
					}

					typeInfo.SchemaName = matches[2]
					typeInfo.Schema = schema.NewNamedType(matches[2])
				}
				continue
			}

			comments = append(comments, text)
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

	return typeInfo, nil
}

func (sp *SchemaParser) parseOperationInfo(fn *types.Func) *OperationInfo {
	functionName := fn.Name()
	result := OperationInfo{
		OriginName: functionName,
		Arguments:  make(map[string]ArgumentInfo),
	}

	var descriptions []string
	commentGroup := findCommentsFromPos(sp.GetCurrentPackage(), fn.Scope(), functionName)
	if commentGroup != nil {
		for i, comment := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))

			// trim the function name in the first line if exists
			if i == 0 {
				text = strings.TrimPrefix(text, fmt.Sprintf("%s ", functionName))
			}
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
					result.Name = ToCamelCase(functionName)
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
		result.Name = ToCamelCase(operationNameResults[2])
	}

	desc := strings.TrimSpace(strings.Join(descriptions, " "))
	if desc != "" {
		result.Description = &desc
	}

	return &result
}

func findCommentsFromPos(pkg *packages.Package, scope *types.Scope, name string) *ast.CommentGroup {
	if pkg == nil {
		return nil
	}

	for _, f := range pkg.Syntax {
		for _, cg := range f.Comments {
			if len(cg.List) == 0 {
				continue
			}
			exp := regexp.MustCompile(fmt.Sprintf(`^//\s+%s`, name))
			if !exp.MatchString(cg.List[0].Text) {
				continue
			}
			if _, obj := scope.LookupParent(name, cg.Pos()); obj != nil {
				return cg
			}
		}
	}
	return nil
}

// get field name by json tag
// return the struct field name if not exist
func getFieldNameOrTag(name string, tag string) string {
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

func findAndReplaceNativeScalarPackage(input string) (string, string, bool) {
	for alias, pkg := range nativeScalarPackages {
		if pkg.Pattern.MatchString(input) {
			return pkg.PackageName, strings.ReplaceAll(input, pkg.PackageName, alias), true
		}
	}
	return "", "", false
}
