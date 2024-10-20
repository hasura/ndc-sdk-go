package internal

import (
	"fmt"
	"strings"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

// WriteGoSchema writes the schema as Go codes
func (rcs RawConnectorSchema) WriteGoSchema(packageName string) (string, error) {
	builder := strings.Builder{}
	writeFileHeader(&builder, packageName)
	builder.WriteString(`
import (
  "github.com/hasura/ndc-sdk-go/schema"
)


func toPtr[V any](value V) *V {
	return &value
}

// GetConnectorSchema gets the generated connector schema
func GetConnectorSchema() *schema.SchemaResponse {
	return &schema.SchemaResponse{
		Collections: []schema.CollectionInfo{},
		ObjectTypes: schema.SchemaResponseObjectTypes{`)

	objectKeys := utils.GetSortedKeys(rcs.ObjectSchemas)
	for _, key := range objectKeys {
		objectType := rcs.ObjectSchemas[key]
		if err := rcs.writeObjectType(&builder, key, objectType); err != nil {
			return "", err
		}
	}
	builder.WriteString(`
		},
		Functions: []schema.FunctionInfo{`)
	for _, fn := range rcs.Functions {
		fnSchema := fn.Schema()
		if err := rcs.writeOperationInfo(&builder, fnSchema.Name, fnSchema.Description, fnSchema.Arguments, fnSchema.ResultType); err != nil {
			return "", err
		}
	}

	builder.WriteString(`
		},
		Procedures: []schema.ProcedureInfo{`)
	for _, proc := range rcs.Procedures {
		procSchema := proc.Schema()
		if err := rcs.writeOperationInfo(&builder, procSchema.Name, procSchema.Description, procSchema.Arguments, procSchema.ResultType); err != nil {
			return "", err
		}
	}

	builder.WriteString(`
		},
		ScalarTypes: schema.SchemaResponseScalarTypes{`)
	scalarKeys := utils.GetSortedKeys(rcs.ScalarSchemas)
	for _, key := range scalarKeys {
		scalarType := rcs.ScalarSchemas[key]
		if err := rcs.writeScalarType(&builder, key, scalarType); err != nil {
			return "", err
		}
	}

	builder.WriteString("\n  	},\n	}\n}")
	return builder.String(), nil
}

func (rcs RawConnectorSchema) writeOperationInfo(builder *strings.Builder, name string, desc *string, arguments map[string]schema.ArgumentInfo, resultType schema.Type) error {
	baseIndent := 6
	builder.WriteString(`
			{
				Name: "`)
	builder.WriteString(name)
	builder.WriteString("\",\n")
	rcs.writeDescription(builder, desc)
	writeIndent(builder, baseIndent+2)

	builder.WriteString("ResultType: ")
	retType, err := rcs.writeType(resultType, 0)
	if err != nil {
		return fmt.Errorf("failed to render function %s: %w", name, err)
	}
	builder.WriteString(retType)
	builder.WriteString(",\n")
	writeIndent(builder, baseIndent+2)
	builder.WriteString("Arguments: map[string]schema.ArgumentInfo{")
	argumentKeys := utils.GetSortedKeys(arguments)
	for _, argKey := range argumentKeys {
		argument := arguments[argKey]
		builder.WriteRune('\n')
		writeIndent(builder, baseIndent+4)
		builder.WriteRune('"')
		builder.WriteString(argKey)
		builder.WriteString("\": {\n")
		rcs.writeDescription(builder, argument.Description)
		writeIndent(builder, baseIndent+6)
		builder.WriteString("Type: ")

		argType, err := rcs.writeType(argument.Type, 0)
		if err != nil {
			return fmt.Errorf("failed to render argument %s of function %s: %w", argKey, name, err)
		}
		builder.WriteString(argType)
		builder.WriteString(",\n")
		writeIndent(builder, baseIndent+4)
		builder.WriteString("},")
	}
	builder.WriteRune('\n')
	writeIndent(builder, baseIndent+2)
	builder.WriteString("},\n")
	writeIndent(builder, baseIndent)
	builder.WriteString("},")

	return nil
}

func (rcs RawConnectorSchema) writeScalarType(builder *strings.Builder, key string, scalarType schema.ScalarType) error {
	baseIndent := 6
	builder.WriteRune('\n')
	writeIndent(builder, baseIndent)
	builder.WriteRune('"')

	builder.WriteString(key)
	builder.WriteString(`": schema.ScalarType{
		  	AggregateFunctions:  schema.ScalarTypeAggregateFunctions{},
		  	ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},`)

	if scalarType.Representation != nil {
		builder.WriteRune('\n')
		writeIndent(builder, baseIndent+2)
		builder.WriteString("Representation:      schema.NewTypeRepresentation")
		rep, err := scalarType.Representation.InterfaceT()
		switch t := rep.(type) {
		case *schema.TypeRepresentationBoolean:
			builder.WriteString("Boolean()")
		case *schema.TypeRepresentationBigDecimal:
			builder.WriteString("BigDecimal()")
		case *schema.TypeRepresentationInt8:
			builder.WriteString("Int8()")
		case *schema.TypeRepresentationInt16:
			builder.WriteString("Int16()")
		case *schema.TypeRepresentationInt32:
			builder.WriteString("Int32()")
		case *schema.TypeRepresentationInt64:
			builder.WriteString("Int64()")
		case *schema.TypeRepresentationBigInteger:
			builder.WriteString("BigInteger()")
		case *schema.TypeRepresentationBytes:
			builder.WriteString("Bytes()")
		case *schema.TypeRepresentationDate:
			builder.WriteString("Date()")
		case *schema.TypeRepresentationFloat32:
			builder.WriteString("Float32()")
		case *schema.TypeRepresentationFloat64:
			builder.WriteString("Float64()")
		case *schema.TypeRepresentationJSON:
			builder.WriteString("JSON()")
		case *schema.TypeRepresentationString:
			builder.WriteString("String()")
		case *schema.TypeRepresentationTimestamp:
			builder.WriteString("Timestamp()")
		case *schema.TypeRepresentationTimestampTZ:
			builder.WriteString("TimestampTZ()")
		case *schema.TypeRepresentationUUID:
			builder.WriteString("UUID()")
		case *schema.TypeRepresentationGeography:
			builder.WriteString("Geography()")
		case *schema.TypeRepresentationGeometry:
			builder.WriteString("Geometry()")
		case *schema.TypeRepresentationEnum:
			builder.WriteString("Enum([]string{")
			for i, enum := range t.OneOf {
				if i > 0 {
					builder.WriteString(", ")
				}
				builder.WriteRune('"')
				builder.WriteString(enum)
				builder.WriteRune('"')
			}
			builder.WriteString("})")
		default:
			return err
		}
	}
	builder.WriteString(".Encode(),")
	builder.WriteString("\n    	},")
	return nil
}

func (rcs RawConnectorSchema) writeDescription(builder *strings.Builder, description *string) {
	if description != nil {
		builder.WriteString(`      	Description: toPtr("`)
		builder.WriteString(*description)
		builder.WriteString("\"),\n")
	}
}

func (rcs RawConnectorSchema) writeObjectType(builder *strings.Builder, key string, objectType schema.ObjectType) error {
	baseIndent := 6
	builder.WriteRune('\n')
	writeIndent(builder, baseIndent)
	builder.WriteRune('"')
	builder.WriteString(key)
	builder.WriteString("\": schema.ObjectType{\n")
	rcs.writeDescription(builder, objectType.Description)
	builder.WriteString(strings.Repeat(" ", baseIndent))
	builder.WriteString("  Fields: schema.ObjectTypeFields{\n")

	fieldKeys := utils.GetSortedKeys(objectType.Fields)
	for _, fieldKey := range fieldKeys {
		field := objectType.Fields[fieldKey]
		writeIndent(builder, baseIndent+4)
		builder.WriteRune('"')
		builder.WriteString(fieldKey)
		builder.WriteString("\": schema.ObjectField{\n")
		rcs.writeDescription(builder, field.Description)

		ft, err := rcs.writeType(field.Type, 0)
		if err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
		writeIndent(builder, baseIndent+6)
		builder.WriteString("Type: ")
		builder.WriteString(ft)
		builder.WriteString(",\n")
		writeIndent(builder, baseIndent+4)
		builder.WriteString("},\n")
	}
	writeIndent(builder, baseIndent+2)
	builder.WriteString("},\n")
	writeIndent(builder, baseIndent)
	builder.WriteString("},")

	return nil
}

func (rcs RawConnectorSchema) writeType(schemaType schema.Type, depth uint) (string, error) {
	ty, err := schemaType.InterfaceT()
	switch t := ty.(type) {
	case *schema.ArrayType:
		nested, err := rcs.writeType(t.ElementType, depth+1)
		if err != nil {
			return "", err
		}
		if depth == 0 {
			return fmt.Sprintf("schema.NewArrayType(%s).Encode()", nested), nil
		}
		return fmt.Sprintf("schema.NewArrayType(%s)", nested), nil
	case *schema.NullableType:
		nested, err := rcs.writeType(t.UnderlyingType, depth+1)
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
		return "", fmt.Errorf("invalid schema type: %w", err)
	}
}
