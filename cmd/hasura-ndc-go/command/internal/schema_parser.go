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
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"golang.org/x/tools/go/packages"
)

type SchemaParser struct {
	context      context.Context
	moduleName   string
	rawSchema    *RawConnectorSchema
	packages     []*packages.Package
	packageIndex int
	namingStyle  OperationNamingStyle
}

// GetCurrentPackage gets the current evaluating package.
func (sp SchemaParser) GetCurrentPackage() *packages.Package {
	return sp.packages[sp.packageIndex]
}

// FindPackageByPath finds the package by package path.
func (sp SchemaParser) FindPackageByPath(input string) *packages.Package {
	for _, p := range sp.packages {
		if p.ID == input {
			return p
		}
	}
	return nil
}

func parseRawConnectorSchemaFromGoCode(ctx context.Context, moduleName string, filePath string, args *ConnectorGenerationArguments) (*RawConnectorSchema, error) {
	var err error
	namingStyle := StyleCamelCase
	if args.Style != "" {
		namingStyle, err = ParseOperationNamingStyle(args.Style)
		if err != nil {
			return nil, err
		}
	}
	rawSchema := NewRawConnectorSchema()

	tempDirs := args.Directories
	if len(args.Directories) == 0 {
		// recursively walk directories if the user don't explicitly specify target folders
		entries, err := os.ReadDir(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read subdirectories of %s: %w", filePath, err)
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
				return nil, fmt.Errorf("failed to read subdirectories of %s/%s: %w", filePath, dir, err)
			}
			// cleanup types.generated.go files
			fileCount := 0
			for _, fp := range goFiles {
				if !strings.HasSuffix(fp, typeMethodsOutputFile) {
					fileCount++
					continue
				}
				if err := os.Remove(fp); err != nil {
					return nil, fmt.Errorf("failed to delete %s: %w", fp, err)
				}
			}
			if fileCount > 0 {
				directories = append(directories, dir)
				break
			}
		}
	}

	if len(directories) > 0 {
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
			parseSchemaCtx, parseSchemaTask := trace.NewTask(ctx, "parse_schema_"+packageList[i].ID)
			sp := &SchemaParser{
				context:      parseSchemaCtx,
				moduleName:   moduleName,
				packages:     packageList,
				packageIndex: i,
				rawSchema:    rawSchema,
				namingStyle:  namingStyle,
			}

			err = sp.parseRawConnectorSchema(packageList[i].Types)
			parseSchemaTask.End()
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Info().Msgf("no subdirectory in %s", filePath)
	}

	if rawSchema.StateType != nil {
		rawSchema.Imports[rawSchema.StateType.PackagePath] = true
	} else {
		pkgPathTypes, err := evalPackageTypesLocation(moduleName, filePath, args.ConnectorDir)
		if err != nil {
			return nil, err
		}
		rawSchema.StateType = &TypeInfo{
			Name:        "State",
			PackagePath: pkgPathTypes,
			PackageName: "types",
		}
		rawSchema.Imports[rawSchema.StateType.PackagePath] = true
	}
	return rawSchema, nil
}

// parse raw connector schema from Go code.
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
			return fmt.Errorf("%s: expect 2 or 3 parameters only (ctx context.Context, state *types.State, arguments *[ArgumentType]), got %s", opInfo.OriginName, params)
		}

		if resultTuple == nil || resultTuple.Len() != 2 {
			return fmt.Errorf("%s: expect result tuple ([type], error), got %s", opInfo.OriginName, resultTuple)
		}

		if sp.rawSchema.StateType == nil {
			ty := sp.getNamedType(params.At(1).Type())
			if ty != nil {
				so := ty.Obj()
				if so != nil {
					objPkg := so.Pkg()
					if objPkg != nil {
						sp.rawSchema.StateType = &TypeInfo{
							Name:        so.Name(),
							PackageName: objPkg.Name(),
							PackagePath: objPkg.Path(),
						}
					}
				}
			}
		}

		// parse arguments in the function if exists
		// ignore 2 first parameters (context and state)
		if params.Len() == 3 {
			arg := params.At(2)
			argumentInfo, err := sp.parseArgumentTypes(arg.Type(), &opInfo.Kind, []string{})
			if err != nil {
				return err
			}
			opInfo.ArgumentsType = argumentInfo.Type
			if opInfo.Kind == OperationFunction {
				sp.rawSchema.setFunctionArgument(*argumentInfo)
			}
			// convert argument schema
			for k, a := range argumentInfo.Fields {
				if !a.Type.Embedded {
					opInfo.Arguments[k] = schema.ArgumentInfo{
						Description: a.Description,
						Type:        a.Type.Schema.Encode(),
					}
					continue
				}

				embeddedObject, ok := sp.rawSchema.Objects[a.Type.String()]
				if ok {
					// flatten embedded object fields to the parent object
					for k, of := range embeddedObject.Fields {
						opInfo.Arguments[k] = schema.ArgumentInfo{
							Type: of.Type.Schema.Encode(),
						}
					}
				}
			}
		}

		resultType, err := sp.parseType(nil, resultTuple.At(0).Type(), []string{}, nil)
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

func (sp *SchemaParser) getNamedType(ty types.Type) *types.Named {
	switch t := ty.(type) {
	case *types.Pointer:
		return sp.getNamedType(t.Elem())
	case *types.Named:
		return t
	case *types.Slice:
		return sp.getNamedType(t.Elem())
	case *types.Array:
		return sp.getNamedType(t.Elem())
	default:
		return nil
	}
}

func (sp *SchemaParser) parseArgumentTypes(ty types.Type, argumentFor *OperationKind, fieldPaths []string) (*ObjectInfo, error) {
	switch inferredType := ty.(type) {
	case *types.Pointer:
		return sp.parseArgumentTypes(inferredType.Elem(), argumentFor, fieldPaths)
	case *types.Struct:
		result := &ObjectInfo{
			Fields: map[string]ObjectField{},
		}
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
			typeInfo.Embedded = fieldVar.Embedded()
			fieldType, err := sp.parseType(typeInfo, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), argumentFor)
			if err != nil {
				return nil, err
			}
			fieldName := getFieldNameOrTag(fieldVar.Name(), fieldTag)
			if fieldType.TypeAST == nil {
				fieldType.TypeAST = fieldVar.Type()
			}
			if !fieldType.IsScalar && argumentFor != nil && *argumentFor == OperationFunction {
				object, ok := sp.rawSchema.Objects[fieldType.GetArgumentName(fieldType.PackagePath)]
				if ok {
					sp.rawSchema.setFunctionArgument(object)
				} else {
					sp.rawSchema.setFunctionArgument(ObjectInfo{
						Type: fieldType,
					})
				}
			}

			result.Fields[fieldName] = ObjectField{
				Name: fieldVar.Name(),
				Type: fieldType,
			}
		}
		return result, nil
	case *types.Named:

		typeObj := inferredType.Obj()
		typeInfo := &TypeInfo{
			Name:       typeObj.Name(),
			SchemaName: typeObj.Name(),
			TypeAST:    typeObj.Type().Underlying(),
		}
		pkg := typeObj.Pkg()
		if pkg != nil {
			typeInfo.PackagePath = pkg.Path()
			typeInfo.PackageName = pkg.Name()
		}

		typeParams := inferredType.TypeParams()
		if typeParams != nil && typeParams.Len() > 0 {
			// unwrap the generic type parameters such as Foo[T]
			parseTypeParameters(typeInfo, inferredType.String())
			typeInfo.TypeAST = inferredType.Underlying()
		}

		arguments, err := sp.parseArgumentTypes(typeInfo.TypeAST, argumentFor, append(fieldPaths, typeObj.Name()))
		if err != nil {
			return nil, err
		}

		arguments.Type = typeInfo
		return arguments, nil
	default:
		return nil, fmt.Errorf("expected struct type, got %s", ty.String())
	}
}

func (sp *SchemaParser) parseType(rootType *TypeInfo, ty types.Type, fieldPaths []string, argumentFor *OperationKind) (*TypeInfo, error) {
	switch inferredType := ty.(type) {
	case *types.Pointer:
		innerType, err := sp.parseType(rootType, inferredType.Elem(), fieldPaths, argumentFor)
		if err != nil {
			return nil, err
		}
		innerType.TypeAST = ty
		innerType.TypeFragments = append([]string{"*"}, innerType.TypeFragments...)
		innerType.Schema = schema.NewNullableType(innerType.Schema)
		return innerType, nil
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
			rootType.Schema = schema.NewNamedType(rootType.SchemaName)
		}
		objType := schema.ObjectType{
			Description: rootType.Description,
			Fields:      make(schema.ObjectTypeFields),
		}
		objFields := ObjectInfo{
			IsAnonymous: isAnonymous,
			Type:        rootType,
			Fields:      map[string]ObjectField{},
		}
		// temporarily add the object type to raw schema to avoid infinite loop
		sp.rawSchema.ObjectSchemas[rootType.SchemaName] = objType
		sp.rawSchema.Objects[rootType.String()] = objFields

		for i := 0; i < inferredType.NumFields(); i++ {
			fieldVar := inferredType.Field(i)
			fieldTag := inferredType.Tag(i)

			fieldType, err := sp.parseType(nil, fieldVar.Type(), append(fieldPaths, fieldVar.Name()), argumentFor)
			if err != nil {
				return nil, err
			}
			fieldType.Embedded = fieldVar.Embedded()
			fieldType.TypeAST = fieldVar.Type()
			fieldKey := getFieldNameOrTag(fieldVar.Name(), fieldTag)
			if fieldType.Embedded {
				embeddedObject, ok := sp.rawSchema.ObjectSchemas[fieldType.Name]
				if ok {
					// flatten embedded object fields to the parent object
					for k, of := range embeddedObject.Fields {
						objType.Fields[k] = of
					}
				}
			} else {
				objType.Fields[fieldKey] = schema.ObjectField{
					Type: fieldType.Schema.Encode(),
				}
			}
			objFields.Fields[fieldKey] = ObjectField{
				Name: fieldVar.Name(),
				Type: fieldType,
			}
		}
		sp.rawSchema.ObjectSchemas[rootType.SchemaName] = objType
		sp.rawSchema.Objects[rootType.String()] = objFields

		return rootType, nil
	case *types.Named:

		innerType := inferredType.Obj()
		if innerType == nil {
			return nil, fmt.Errorf("failed to parse named type: %s", inferredType.String())
		}

		innerPkg := innerType.Pkg()
		typeInfo := &TypeInfo{
			Name:       innerType.Name(),
			SchemaName: innerType.Name(),
			TypeAST:    innerType.Type(),
			Schema:     schema.NewNamedType(innerType.Name()),
		}

		if rootType != nil {
			typeInfo.Embedded = rootType.Embedded
		}
		if innerPkg != nil {
			typeInfo.PackageName = innerPkg.Name()
			typeInfo.PackagePath = innerPkg.Path()
			typeInfo.TypeFragments = []string{innerType.Name()}
			typeParams := inferredType.TypeParams()
			if typeParams != nil && typeParams.Len() > 0 {
				// unwrap the generic type parameters such as Foo[T]
				parseTypeParameters(typeInfo, inferredType.String())
				typeInfo.TypeAST = inferredType.Underlying()
			}

			if _, ok := sp.rawSchema.Objects[typeInfo.String()]; ok {
				return typeInfo, nil
			}

			if err := sp.parseTypeInfoFromComments(typeInfo, innerType.Parent()); err != nil {
				return nil, err
			}

			var scalarName ScalarName
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
			if argumentFor != nil {
				return nil, fmt.Errorf("%s: native `error` interface isn't allowed in input arguments", strings.Join(fieldPaths, "."))
			}
			typeInfo.IsScalar = true
			typeInfo.SchemaName = string(ScalarJSON)
			typeInfo.Schema = schema.NewNullableType(schema.NewNamedType(string(ScalarJSON)))
			typeInfo.ScalarRepresentation = schema.NewTypeRepresentationJSON().Encode()

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

		if _, ok := sp.rawSchema.ObjectSchemas[typeInfo.SchemaName]; ok {
			// the object schema exists, rename to format <name>_<package_name>
			packagePath := strings.TrimPrefix(typeInfo.PackagePath, sp.moduleName)
			typeInfo.SchemaName = fieldNameRegex.ReplaceAllString(strings.Join([]string{typeInfo.Name, packagePath}, ""), "_")
			typeInfo.Schema = schema.NewNamedType(typeInfo.SchemaName)
		}

		return sp.parseType(typeInfo, typeInfo.TypeAST.Underlying(), append(fieldPaths, innerType.Name()), argumentFor)
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
		innerType, err := sp.parseType(nil, inferredType.Elem(), fieldPaths, argumentFor)
		if err != nil {
			return nil, err
		}
		innerType.TypeFragments = append([]string{"[]"}, innerType.TypeFragments...)
		innerType.Schema = schema.NewArrayType(innerType.Schema)
		return innerType, nil
	case *types.Slice:
		innerType, err := sp.parseType(nil, inferredType.Elem(), fieldPaths, argumentFor)
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

func (sp *SchemaParser) parseTypeInfoFromComments(typeInfo *TypeInfo, scope *types.Scope) error {
	comments := make([]string, 0)
	commentGroup := findCommentsFromPos(sp.FindPackageByPath(typeInfo.PackagePath), scope, typeInfo.Name)
	if commentGroup != nil {
		for i, line := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(line.Text, "/"))
			if text == "" {
				continue
			}
			if i == 0 {
				text = strings.TrimPrefix(text, typeInfo.Name+" ")
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
					return fmt.Errorf("require enum values in the comment of %s", typeInfo.Name)
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
						return fmt.Errorf("failed to parse type representation of scalar %s: %w", typeInfo.Name, err)
					}
					if typeRep == schema.TypeRepresentationTypeEnum {
						return errors.New("use @enum tag with values instead")
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
							return errors.New("use @enum tag with values instead")
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
		matches := ndcScalarNameRegex.FindStringSubmatch(typeInfo.Name)
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

	return nil
}

// format operation name with style.
func (sp SchemaParser) formatOperationName(name string) string {
	switch sp.namingStyle {
	case StyleSnakeCase:
		return strcase.ToSnake(name)
	default:
		return strcase.ToLowerCamel(name)
	}
}

func (sp *SchemaParser) parseOperationInfo(fn *types.Func) *OperationInfo {
	functionName := fn.Name()
	result := OperationInfo{
		OriginName: functionName,
		Arguments:  make(map[string]schema.ArgumentInfo),
	}

	var descriptions []string
	commentGroup := findCommentsFromPos(sp.GetCurrentPackage(), fn.Scope(), functionName)
	if commentGroup != nil {
		for i, comment := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(comment.Text, "/"))

			// trim the function name in the first line if exists
			if i == 0 {
				text = strings.TrimPrefix(text, functionName+" ")
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
					result.Name = sp.formatOperationName(functionName)
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
		result.Name = sp.formatOperationName(operationNameResults[2])
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
			exp := regexp.MustCompile(fmt.Sprintf(`^//\s+%s(\s|$)`, name))
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
// return the struct field name if not exist.
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

func evalPackageTypesLocation(moduleName string, filePath string, connectorDir string) (string, error) {
	matches, err := filepath.Glob(path.Join(filePath, "types", "*.go"))
	if err == nil && len(matches) > 0 {
		return moduleName + "/types", nil
	}

	if connectorDir != "" && !strings.HasPrefix(".", connectorDir) {
		matches, err = filepath.Glob(path.Join(filePath, connectorDir, "types", "*.go"))
		if err == nil && len(matches) > 0 {
			return fmt.Sprintf("%s/%s/types", moduleName, connectorDir), nil
		}
	}
	return "", fmt.Errorf("the `types` package where the State struct is in must be placed in root or connector directory, %w", err)
}

func parseTypeParameters(rootType *TypeInfo, input string) {
	paramsString := strings.TrimRight(strings.TrimLeft(strings.TrimPrefix(input, rootType.PackagePath+"."+rootType.Name), "["), "]")
	rawParams := strings.Split(paramsString, ",")

	for _, param := range rawParams {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}
		typeInfo := TypeInfo{}

		parts := strings.Split(param, ".")
		partsLen := len(parts)
		if partsLen == 1 {
			typeInfo.Name = parts[0]
		} else {
			typeInfo.PackagePath = strings.Join(parts[0:partsLen-1], ".")
			packageParts := strings.Split(typeInfo.PackagePath, "/")
			typeInfo.PackageName = packageParts[len(packageParts)-1]
			typeInfo.Name = parts[partsLen-1]
		}
		rootType.SchemaName += "_" + typeInfo.Name
		rootType.TypeParameters = append(rootType.TypeParameters, typeInfo)
	}
	rootType.Schema = schema.NewNamedType(rootType.SchemaName)
}
