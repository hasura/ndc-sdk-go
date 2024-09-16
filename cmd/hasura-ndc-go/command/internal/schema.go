package internal

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
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

// IsArray checks if the current type is an array
func (ti *TypeInfo) IsArray() bool {
	return isArrayFragments(ti.TypeFragments)
}

func isNullableFragment(fragment string) bool {
	return fragment == "*"
}

func isNullableFragments(fragments []string) bool {
	return len(fragments) > 0 && isNullableFragment(fragments[0])
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
	TypesPackagePath string
	Imports          map[string]bool
	CustomScalars    map[string]*TypeInfo
	ScalarSchemas    schema.SchemaResponseScalarTypes
	Objects          map[string]*ObjectInfo
	ObjectSchemas    schema.SchemaResponseObjectTypes
	Functions        []FunctionInfo
	Procedures       []ProcedureInfo
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

// Render renders the schema to Go codes
func (rcs RawConnectorSchema) Render(packageName string) (string, error) {
	builder := strings.Builder{}
	_, _ = builder.WriteString(`// Code generated by github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go, DO NOT EDIT.
package `)
	_, _ = builder.WriteString(packageName)
	_, _ = builder.WriteString(`
import (
  "github.com/hasura/ndc-sdk-go/schema"
  "github.com/hasura/ndc-sdk-go/utils"
)

// GetConnectorSchema gets the generated connector schema
func GetConnectorSchema() schema.SchemaResponse {
	return schema.SchemaResponse{
		Collections: []schema.CollectionInfo{},
		ObjectTypes: schema.SchemaResponseObjectTypes{`)

	objectKeys := utils.GetSortedKeys(rcs.ObjectSchemas)
	for _, key := range objectKeys {
		objectType := rcs.ObjectSchemas[key]
		if err := rcs.renderObjectType(&builder, key, objectType); err != nil {
			return "", err
		}
	}
	_, _ = builder.WriteString(`
		},
		Functions: []schema.FunctionInfo{`)
	for _, fn := range rcs.Functions {
		fnSchema := fn.Schema()
		if err := rcs.renderOperationInfo(&builder, fnSchema.Name, fnSchema.Description, fnSchema.Arguments, fnSchema.ResultType); err != nil {
			return "", err
		}
	}

	_, _ = builder.WriteString(`
		},
		Procedures: []schema.ProcedureInfo{`)
	for _, proc := range rcs.Procedures {
		procSchema := proc.Schema()
		if err := rcs.renderOperationInfo(&builder, procSchema.Name, procSchema.Description, procSchema.Arguments, procSchema.ResultType); err != nil {
			return "", err
		}
	}

	_, _ = builder.WriteString(`
		},
		ScalarTypes: schema.SchemaResponseScalarTypes{`)
	scalarKeys := utils.GetSortedKeys(rcs.ScalarSchemas)
	for _, key := range scalarKeys {
		scalarType := rcs.ScalarSchemas[key]
		if err := rcs.renderScalarType(&builder, key, scalarType); err != nil {
			return "", err
		}
	}

	_, _ = builder.WriteString("\n  	},\n	}\n}")
	return builder.String(), nil
}

func (rcs RawConnectorSchema) renderOperationInfo(builder *strings.Builder, name string, desc *string, arguments map[string]schema.ArgumentInfo, resultType schema.Type) error {
	_, _ = builder.WriteString(`
		{
			Name: "`)
	_, _ = builder.WriteString(name)
	_, _ = builder.WriteString("\",\n")
	rcs.renderDescription(builder, desc)
	_, _ = builder.WriteString(`      	ResultType: `)
	retType, err := rcs.renderType(resultType, 0)
	if err != nil {
		return fmt.Errorf("failed to render function %s: %s", name, err)
	}
	_, _ = builder.WriteString(retType)
	_, _ = builder.WriteString(",\n      	Arguments: map[string]schema.ArgumentInfo{")
	for argKey, argument := range arguments {
		_, _ = builder.WriteString("\n        \"")
		_, _ = builder.WriteString(argKey)
		_, _ = builder.WriteString("\": {\n")
		rcs.renderDescription(builder, argument.Description)
		_, _ = builder.WriteString("          	Type: ")

		argType, err := rcs.renderType(argument.Type, 0)
		if err != nil {
			return fmt.Errorf("failed to render argument %s of function %s: %s", argKey, name, err)
		}
		_, _ = builder.WriteString(argType)
		_, _ = builder.WriteString(",\n	        },")
	}
	_, _ = builder.WriteString("\n	      },\n	    },")

	return nil
}

func (rcs RawConnectorSchema) renderScalarType(builder *strings.Builder, key string, scalarType schema.ScalarType) error {
	_, _ = builder.WriteString("\n    \"")
	_, _ = builder.WriteString(key)
	_, _ = builder.WriteString(`": schema.ScalarType{
		  	AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		  	ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},`)

	if scalarType.Representation != nil {
		_, _ = builder.WriteString("\n      Representation:      schema.NewTypeRepresentation")
		rep, err := scalarType.Representation.InterfaceT()
		switch t := rep.(type) {
		case *schema.TypeRepresentationBoolean:
			_, _ = builder.WriteString("Boolean()")
		case *schema.TypeRepresentationBigDecimal:
			_, _ = builder.WriteString("BigDecimal()")
		case *schema.TypeRepresentationInt8:
			_, _ = builder.WriteString("Int8()")
		case *schema.TypeRepresentationInt16:
			_, _ = builder.WriteString("Int16()")
		case *schema.TypeRepresentationInt32:
			_, _ = builder.WriteString("Int32()")
		case *schema.TypeRepresentationInt64:
			_, _ = builder.WriteString("Int64()")
		case *schema.TypeRepresentationBigInteger:
			_, _ = builder.WriteString("BigInteger()")
		case *schema.TypeRepresentationBytes:
			_, _ = builder.WriteString("Bytes()")
		case *schema.TypeRepresentationDate:
			_, _ = builder.WriteString("Date()")
		case *schema.TypeRepresentationFloat32:
			_, _ = builder.WriteString("Float32()")
		case *schema.TypeRepresentationFloat64:
			_, _ = builder.WriteString("Float64()")
		case *schema.TypeRepresentationJSON:
			_, _ = builder.WriteString("JSON()")
		case *schema.TypeRepresentationString:
			_, _ = builder.WriteString("String()")
		case *schema.TypeRepresentationTimestamp:
			_, _ = builder.WriteString("Timestamp()")
		case *schema.TypeRepresentationTimestampTZ:
			_, _ = builder.WriteString("TimestampTZ()")
		case *schema.TypeRepresentationUUID:
			_, _ = builder.WriteString("UUID()")
		case *schema.TypeRepresentationGeography:
			_, _ = builder.WriteString("Geography()")
		case *schema.TypeRepresentationGeometry:
			_, _ = builder.WriteString("Geometry()")
		case *schema.TypeRepresentationEnum:
			_, _ = builder.WriteString("Enum([]string{")
			for i, enum := range t.OneOf {
				if i > 0 {
					_, _ = builder.WriteString(", ")
				}
				_, _ = builder.WriteRune('"')
				_, _ = builder.WriteString(enum)
				_, _ = builder.WriteRune('"')
			}
			_, _ = builder.WriteString("})")
		default:
			return err
		}
	}
	_, _ = builder.WriteString(".Encode(),")
	_, _ = builder.WriteString("\n    	},")
	return nil
}

func (rcs RawConnectorSchema) renderDescription(builder *strings.Builder, description *string) {
	if description != nil {
		_, _ = builder.WriteString(`      	Description: utils.ToPtr("`)
		_, _ = builder.WriteString(*description)
		_, _ = builder.WriteString("\"),\n")
	}
}

func (rcs RawConnectorSchema) renderObjectType(builder *strings.Builder, key string, objectType schema.ObjectType) error {
	_, _ = builder.WriteString("\n    \"")
	_, _ = builder.WriteString(key)
	_, _ = builder.WriteString("\": schema.ObjectType{\n")

	rcs.renderDescription(builder, objectType.Description)

	_, _ = builder.WriteString("      Fields: schema.ObjectTypeFields{\n")
	for fieldKey, field := range objectType.Fields {
		_, _ = builder.WriteString(`        "`)
		_, _ = builder.WriteString(fieldKey)
		_, _ = builder.WriteString("\": schema.ObjectField{\n")
		rcs.renderDescription(builder, field.Description)

		ft, err := rcs.renderType(field.Type, 0)
		if err != nil {
			return fmt.Errorf("%s: %s", key, err)
		}
		_, _ = builder.WriteString("          Type: ")
		_, _ = builder.WriteString(ft)
		_, _ = builder.WriteString(",\n        },\n")
	}
	_, _ = builder.WriteString("      },\n    },")
	return nil
}

func (rcs RawConnectorSchema) renderType(schemaType schema.Type, depth uint) (string, error) {
	ty, err := schemaType.InterfaceT()
	switch t := ty.(type) {
	case *schema.ArrayType:
		nested, err := rcs.renderType(t.ElementType, depth+1)
		if err != nil {
			return "", err
		}
		if depth == 0 {
			return fmt.Sprintf("schema.NewArrayType(%s).Encode()", nested), nil
		}
		return fmt.Sprintf("schema.NewArrayType(%s)", nested), nil
	case *schema.NullableType:
		nested, err := rcs.renderType(t.UnderlyingType, depth+1)
		if err != nil {
			return "", err
		}
		if depth == 0 {
			return fmt.Sprintf("schema.NewNullableType(%s).Encode()", nested), nil
		}
		return fmt.Sprintf("schema.NewNullableType(%s)", nested), nil
	case *schema.NamedType:
		if depth == 0 {
			return fmt.Sprintf(`schema.NewNamedType("%s").Encode()`, t.Name), nil
		}
		return fmt.Sprintf(`schema.NewNamedType("%s")`, t.Name), nil
	default:
		return "", fmt.Errorf("invalid schema type: %s", err)
	}
}
