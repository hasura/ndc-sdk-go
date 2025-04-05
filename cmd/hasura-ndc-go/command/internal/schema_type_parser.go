package internal

import (
	"errors"
	"fmt"
	"go/types"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
)

type TypeParser struct {
	schemaParser *SchemaParser
	field        *Field
	rootType     types.Type
	argumentFor  *OperationKind
	tagInfo      NDCTagInfo
	// cached parent named type if the underlying type is an object
	typeInfo *TypeInfo
}

func NewTypeParser(schemaParser *SchemaParser, field *Field, ty types.Type, tagInfo NDCTagInfo, argumentFor *OperationKind) *TypeParser {
	return &TypeParser{
		schemaParser: schemaParser,
		field:        field,
		rootType:     ty,
		tagInfo:      tagInfo,
		argumentFor:  argumentFor,
	}
}

func (tp *TypeParser) Parse(fieldPaths []string) (*Field, error) {
	ty, err := tp.parseType(tp.rootType, fieldPaths)
	if err != nil {
		return nil, err
	}

	if ty == nil {
		return nil, nil
	}

	tp.field.Type = ty

	return tp.field, nil
}

func (tp *TypeParser) ParseArgumentTypes(fieldPaths []string) (*ObjectInfo, error) {
	return tp.parseArgumentTypes(tp.rootType, fieldPaths)
}

func (tp *TypeParser) parseArgumentTypes(ty types.Type, fieldPaths []string) (*ObjectInfo, error) {
	switch inferredType := ty.(type) {
	case *types.Pointer:
		return tp.parseArgumentTypes(inferredType.Elem(), fieldPaths)
	case *types.Struct:
		result := &ObjectInfo{
			Fields:       map[string]Field{},
			SchemaFields: schema.ObjectTypeFields{},
		}

		if err := tp.parseStructType(result, inferredType, fieldPaths); err != nil {
			return nil, err
		}

		return result, nil
	case *types.Named:
		typeObj := inferredType.Obj()
		if typeObj == nil {
			return nil, fmt.Errorf("named type %s does not exist", inferredType.String())
		}

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
			if err := parseTypeParameters(typeInfo, inferredType.String()); err != nil {
				return nil, err
			}

			typeInfo.TypeAST = inferredType.Underlying()
		}

		arguments, err := tp.parseArgumentTypes(typeInfo.TypeAST, append(fieldPaths, typeObj.Name()))
		if err != nil {
			return nil, err
		}

		arguments.Type = typeInfo

		return arguments, nil
	default:
		return nil, fmt.Errorf("expected struct type, got %s", ty.String())
	}
}

func (tp *TypeParser) parseType(ty types.Type, fieldPaths []string) (Type, error) {
	switch inferredType := ty.(type) {
	case *types.Pointer:
		innerType, err := tp.parseType(inferredType.Elem(), fieldPaths)
		if err != nil {
			return nil, err
		}

		return NewNullableType(innerType), nil
	case *types.Struct:
		typeInfo := tp.typeInfo
		if typeInfo == nil {
			typeInfo = &TypeInfo{}
		}

		if typeInfo.Name == "" {
			typeInfo.Name = ty.String()
			typeInfo.SchemaName = strings.Join(fieldPaths, "")
			typeInfo.TypeAST = ty
			typeInfo.IsAnonymous = true
		}

		objFields := ObjectInfo{
			Description:  typeInfo.Description,
			Type:         typeInfo,
			Fields:       map[string]Field{},
			SchemaFields: schema.ObjectTypeFields{},
		}

		// temporarily add the object type to raw schema to avoid infinite loop
		tp.schemaParser.rawSchema.Objects[typeInfo.SchemaName] = objFields

		if err := tp.parseStructType(&objFields, inferredType, fieldPaths); err != nil {
			return nil, err
		}

		tp.schemaParser.rawSchema.Objects[typeInfo.SchemaName] = objFields

		if tp.argumentFor != nil && *tp.argumentFor == OperationFunction {
			tp.schemaParser.rawSchema.setFunctionArgument(objFields)
		}

		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	case *types.Named:
		return tp.parseNamedType(inferredType, fieldPaths)
	case *types.Basic:
		typeInfo := tp.typeInfo
		if typeInfo == nil {
			typeInfo = &TypeInfo{
				TypeAST: ty,
			}
		}

		if typeInfo.Name == "" {
			typeInfo.Name = inferredType.Name()
			typeInfo.PackagePath = ""
		}

		switch inferredType.Kind() {
		case types.Bool:
			typeInfo.SchemaName = string(ScalarBoolean)
		case types.Int8, types.Uint8:
			typeInfo.SchemaName = string(ScalarInt8)
		case types.Int16, types.Uint16:
			typeInfo.SchemaName = string(ScalarInt16)
		case types.Int, types.Int32, types.Uint, types.Uint32:
			typeInfo.SchemaName = string(ScalarInt32)
		case types.Int64, types.Uint64:
			typeInfo.SchemaName = string(ScalarInt64)
		case types.Float32:
			typeInfo.SchemaName = string(ScalarFloat32)
		case types.Float64:
			typeInfo.SchemaName = string(ScalarFloat64)
		case types.String:
			typeInfo.SchemaName = string(ScalarString)
		default:
			return nil, fmt.Errorf("%s: unsupported scalar type <%s>", strings.Join(fieldPaths, "."), inferredType.String())
		}

		tp.schemaParser.rawSchema.SetScalar(typeInfo.SchemaName, Scalar{
			Schema: defaultScalarTypes[ScalarName(typeInfo.SchemaName)],
		})

		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	case *types.Array:
		return tp.parseSliceType(inferredType.Elem(), fieldPaths)
	case *types.Slice:
		return tp.parseSliceType(inferredType.Elem(), fieldPaths)
	case *types.Map, *types.Interface:
		scalarName := ScalarJSON

		typeInfo := tp.typeInfo
		if typeInfo == nil {
			typeInfo = &TypeInfo{
				TypeAST: ty,
			}
		}

		if typeInfo.Name == "" {
			typeInfo.Name = inferredType.String()
		}

		typeInfo.PackagePath = ""
		typeInfo.SchemaName = string(scalarName)
		tp.schemaParser.rawSchema.SetScalar(string(scalarName), Scalar{
			Schema:     defaultScalarTypes[ScalarJSON],
			NativeType: typeInfo,
		})

		return NewNamedType(string(ScalarJSON), typeInfo), nil
	case *types.Alias:
		return tp.parseType(inferredType.Underlying(), fieldPaths)
	case *types.Chan, *types.Signature, *types.Tuple, *types.Union:
		return nil, nil
	default:
		return nil, fmt.Errorf("%s: unsupported type: %s", strings.Join(fieldPaths, "."), ty.String())
	}
}

func (tp *TypeParser) parseNamedType(inferredType *types.Named, fieldPaths []string) (Type, error) { //nolint:cyclop,funlen
	innerType := inferredType.Obj()
	if innerType == nil {
		return nil, fmt.Errorf("failed to parse named type: %s", inferredType.String())
	}

	typeInfo := &TypeInfo{
		Name:       innerType.Name(),
		SchemaName: innerType.Name(),
		TypeAST:    innerType.Type(),
	}

	if innerType.Name() == "error" {
		if tp.argumentFor != nil {
			return nil, fmt.Errorf("%s: native `error` interface isn't allowed in input arguments", strings.Join(fieldPaths, "."))
		}

		scalarName := string(ScalarJSON)
		errorScalar := defaultScalarTypes[ScalarJSON]
		typeInfo.SchemaName = scalarName
		tp.schemaParser.rawSchema.SetScalar(scalarName, Scalar{
			Schema: errorScalar,
		})

		return NewNullableType(NewNamedType(scalarName, typeInfo)), nil
	}

	innerPkg := innerType.Pkg()
	if innerPkg == nil {
		return nil, fmt.Errorf("%s: unsupported type <%s>", strings.Join(fieldPaths, "."), innerType.Name())
	}

	typeInfo.PackageName = innerPkg.Name()
	typeInfo.PackagePath = innerPkg.Path()
	typeParams := inferredType.TypeParams()

	if typeParams != nil && typeParams.Len() > 0 {
		// unwrap the generic type parameters such as Foo[T]
		if err := parseTypeParameters(typeInfo, inferredType.String()); err != nil {
			return nil, err
		}

		typeInfo.TypeAST = inferredType.Underlying()
	}

	if object, ok := tp.schemaParser.rawSchema.Objects[typeInfo.SchemaName]; ok {
		if tp.argumentFor != nil && *tp.argumentFor == OperationFunction {
			tp.schemaParser.rawSchema.setFunctionArgument(object)
		}

		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	}

	if _, ok := tp.schemaParser.rawSchema.Scalars[typeInfo.SchemaName]; ok {
		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	}

	scalarType, err := tp.parseTypeInfoFromComments(typeInfo, innerType.Parent())
	if err != nil {
		return nil, err
	}

	if scalarType != nil {
		if len(scalarType.Schema.Representation) == 0 {
			// requires representation since NDC spec v0.1.2
			scalarType.Schema = defaultScalarTypes[ScalarJSON]
		}

		tp.schemaParser.rawSchema.SetScalar(typeInfo.SchemaName, *scalarType)

		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	}

	switch innerPkg.Path() {
	case "time":
		switch innerType.Name() {
		case "Time":
			typeInfo.SchemaName = string(ScalarTimestampTZ)
			scalarType = &Scalar{
				Schema: defaultScalarTypes[ScalarTimestampTZ],
			}
		case "Duration":
			return nil, errUnsupportedTypeDuration
		default:
			return nil, fmt.Errorf("unsupported type %s.%s", innerPkg.Path(), innerType.Name())
		}
	case "encoding/json":
		switch innerType.Name() {
		case "RawMessage":
			typeInfo.SchemaName = string(ScalarRawJSON)
			scalarType = &Scalar{
				Schema: defaultScalarTypes[ScalarRawJSON],
			}
		default:
			return nil, fmt.Errorf("unsupported type %s.%s", innerPkg.Path(), innerType.Name())
		}
	case "github.com/google/uuid":
		switch innerType.Name() {
		case "UUID":
			typeInfo.SchemaName = string(ScalarUUID)
			scalarType = &Scalar{
				Schema: defaultScalarTypes[ScalarUUID],
			}
		default:
			return nil, fmt.Errorf("unsupported type %s.%s", innerPkg.Path(), innerType.Name())
		}
	case packageSDKSchema:
		if innerType.Name() == "Expression" && tp.argumentFor != nil {
			if tp.tagInfo.PredicateObjectName == "" {
				return nil, fmt.Errorf("%s: predicate field tag must be set `ndc:\"predicate=<object-name>\"`", strings.Join(fieldPaths, "."))
			}

			return NewNullableType(NewPredicateType(tp.tagInfo.PredicateObjectName)), nil
		}
	case "github.com/hasura/ndc-sdk-go/scalar":
		scalarName := ScalarName(innerType.Name())
		switch scalarName {
		case ScalarDate, ScalarBigInt, ScalarBytes, ScalarURL, ScalarDuration, ScalarDurationString, ScalarDurationInt64:
			typeInfo.SchemaName = innerType.Name()
			scalarType = &Scalar{
				Schema: defaultScalarTypes[scalarName],
			}
		default:
			return nil, fmt.Errorf("unsupported scalar type %s.%s", innerPkg.Path(), innerType.Name())
		}
	}

	if scalarType != nil {
		tp.schemaParser.rawSchema.SetScalar(typeInfo.SchemaName, *scalarType)

		return NewNamedType(typeInfo.SchemaName, typeInfo), nil
	}

	if _, ok := tp.schemaParser.rawSchema.Objects[typeInfo.SchemaName]; ok {
		// the object schema exists, rename to format <name>_<package_name>
		packagePath := strings.TrimPrefix(typeInfo.PackagePath, tp.schemaParser.moduleName)
		typeInfo.SchemaName = fieldNameRegex.ReplaceAllString(strings.Join([]string{typeInfo.Name, packagePath}, ""), "_")
	}

	tp.typeInfo = typeInfo

	return tp.parseType(typeInfo.TypeAST.Underlying(), append(fieldPaths, innerType.Name()))
}

func (tp *TypeParser) parseStructType(objectInfo *ObjectInfo, inferredType *types.Struct, fieldPaths []string) error {
	for i := range inferredType.NumFields() {
		fieldVar := inferredType.Field(i)
		if !fieldVar.Exported() {
			continue
		}

		fieldTag := inferredType.Tag(i)

		tagInfo, err := parseNDCTagInfo(fieldTag)
		if err != nil {
			return fmt.Errorf("%s: %w", strings.Join(fieldPaths, "."), err)
		}

		if tagInfo.Ignored {
			continue
		}

		fieldKey := tagInfo.Name
		if fieldKey == "" {
			fieldKey = fieldVar.Name()
		}

		fieldParser := NewTypeParser(tp.schemaParser, &Field{
			Name:     fieldVar.Name(),
			Embedded: fieldVar.Embedded(),
			TypeAST:  fieldVar.Type(),
		}, fieldVar.Type(), tagInfo, tp.argumentFor)

		field, err := fieldParser.Parse(append(fieldPaths, fieldVar.Name()))
		if err != nil {
			return err
		}

		if field == nil {
			continue
		}

		embeddedObject, ok := tp.schemaParser.rawSchema.Objects[field.Type.SchemaName(false)]
		if field.Embedded && ok {
			// flatten embedded object fields to the parent object
			for k, of := range embeddedObject.SchemaFields {
				objectInfo.SchemaFields[k] = of
			}
		} else {
			fieldSchema := field.Type.Schema()
			if tagInfo.OmitEmpty && field.Type.Kind() != schema.TypeNullable {
				fieldSchema = schema.NewNullableType(fieldSchema)
			}

			objectInfo.SchemaFields[fieldKey] = schema.ObjectField{
				Type: fieldSchema.Encode(),
			}
		}

		objectInfo.Fields[fieldKey] = *field
	}

	return nil
}

func (tp *TypeParser) parseSliceType(ty types.Type, fieldPaths []string) (Type, error) {
	innerType, err := tp.parseType(ty, fieldPaths)
	if err != nil {
		return nil, err
	}

	return NewArrayType(innerType), nil
}

func (tp *TypeParser) parseTypeInfoFromComments(typeInfo *TypeInfo, scope *types.Scope) (*Scalar, error) {
	var scalarType *Scalar

	comments := make([]string, 0)
	commentGroup := findCommentsFromPos(tp.schemaParser.FindPackageByPath(typeInfo.PackagePath), scope, typeInfo.Name)

	if commentGroup != nil { //nolint:nestif
		for i, line := range commentGroup.List {
			text := strings.TrimSpace(strings.TrimLeft(line.Text, "/"))
			if text == "" {
				continue
			}

			if i == 0 {
				text = strings.TrimPrefix(text, typeInfo.Name+" ")
			}

			enumMatches := ndcEnumCommentRegex.FindStringSubmatch(strings.TrimRight(text, "."))
			if len(enumMatches) == 2 {
				rawEnumItems := strings.Split(enumMatches[1], ",")

				var enums []string

				for _, item := range rawEnumItems {
					trimmed := strings.TrimSpace(item)
					if trimmed != "" {
						enums = append(enums, trimmed)
					}
				}

				if len(enums) == 0 {
					return nil, errors.New("require enum values in the comment of " + typeInfo.Name)
				}

				typeInfo.SchemaName = typeInfo.Name
				scalarType = &Scalar{
					Schema:     *schema.NewScalarType(),
					NativeType: typeInfo,
				}
				scalarType.Schema.Representation = schema.NewTypeRepresentationEnum(enums).Encode()

				continue
			}

			matches := ndcScalarCommentRegex.FindStringSubmatch(text)
			matchesLen := len(matches)

			if matchesLen > 1 {
				if matchesLen > 3 && matches[3] != "" {
					typeInfo.SchemaName = matches[2]
					typeRep, err := schema.ParseTypeRepresentationType(strings.TrimSpace(matches[3]))

					if err != nil {
						return nil, fmt.Errorf("failed to parse type representation of scalar %s: %w", typeInfo.Name, err)
					}

					if typeRep == schema.TypeRepresentationTypeEnum {
						return nil, errMustUseEnumTag
					}

					scalarType = &Scalar{
						Schema:     *schema.NewScalarType(),
						NativeType: typeInfo,
					}
					scalarType.Schema.Representation = schema.TypeRepresentation{
						"type": typeRep,
					}
				} else if matchesLen > 2 && matches[2] != "" {
					typeInfo.SchemaName = matches[2]
					scalarType = &Scalar{
						Schema:     *schema.NewScalarType(),
						NativeType: typeInfo,
					}
					// if the second string is a type representation, use it as a TypeRepresentation instead
					// e.g @scalar string
					typeRep, err := schema.ParseTypeRepresentationType(matches[2])
					if err == nil {
						if typeRep == schema.TypeRepresentationTypeEnum {
							return nil, errMustUseEnumTag
						}

						scalarType.Schema.Representation = schema.TypeRepresentation{
							"type": typeRep,
						}
					}
				}

				continue
			}

			comments = append(comments, text)
		}
	}

	if scalarType == nil {
		// fallback to parse scalar from type name with Scalar prefix
		matches := ndcScalarNameRegex.FindStringSubmatch(typeInfo.Name)
		if len(matches) > 1 {
			typeInfo.SchemaName = matches[1]
			scalarType = &Scalar{
				Schema:     *schema.NewScalarType(),
				NativeType: typeInfo,
			}
		}
	}

	desc := strings.Join(comments, " ")
	if desc != "" {
		typeInfo.Description = &desc
	}

	return scalarType, nil
}

func parseTypeParameters(rootType *TypeInfo, input string) error {
	paramsString := strings.TrimPrefix(input, rootType.PackagePath+"."+rootType.Name)
	if paramsString[0] == '[' {
		paramsString = paramsString[1:]
	}

	if paramsString[len(paramsString)-1] == ']' {
		paramsString = paramsString[:len(paramsString)-1]
	}

	rawParams := strings.Split(paramsString, ",")

	for _, param := range rawParams {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		ty, err := parseTypeFromString(param)
		if err != nil {
			return err
		}

		rootType.SchemaName += "_" + ty.SchemaName(true)
		rootType.TypeParameters = append(rootType.TypeParameters, ty)
	}

	return nil
}

func parseTypeFromString(input string) (Type, error) {
	if len(input) == 0 {
		return nil, errors.New("failed to parse type from string, the input value is empty")
	}

	if input[0] == '*' {
		underlyingType, err := parseTypeFromString(input[1:])
		if err != nil {
			return nil, err
		}

		return NewNullableType(underlyingType), nil
	}

	if len(input) >= 2 && input[0:2] == "[]" {
		elementType, err := parseTypeFromString(input[2:])
		if err != nil {
			return nil, err
		}

		return NewArrayType(elementType), nil
	}

	parts := strings.Split(input, ".")

	partsLen := len(parts)
	if partsLen == 1 {
		return NewNamedType(parts[0], &TypeInfo{
			Name: parts[0],
		}), nil
	}

	typeInfo := &TypeInfo{}
	typeInfo.PackagePath = strings.Join(parts[0:partsLen-1], ".")
	packageParts := strings.Split(typeInfo.PackagePath, "/")
	typeInfo.PackageName = packageParts[len(packageParts)-1]
	typeInfo.Name = parts[partsLen-1]

	return NewNamedType(typeInfo.Name, typeInfo), nil
}
