package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hasura/ndc-sdk-go/schema"
)

const (
	errFunctionValueFieldRequired = "__value field is required in query function type"

	// CommandSelectionFieldKey the context key for the nested selection field in command.
	CommandSelectionFieldKey string = "ndc-command-selection-field"
)

var ErrHandlerNotfound = errors.New("connector handler not found")

// Scalar abstracts a scalar interface to determine when evaluating.
type Scalar interface {
	ScalarName() string
}

// EvalNestedColumnObject evaluate and prune nested fields from an object without relationship.
func EvalNestedColumnObject(fields *schema.NestedObject, value any) (any, error) {
	return evalNestedColumnObject(fields, value, "")
}

func evalNestedColumnObject(fields *schema.NestedObject, value any, fieldPath string) (any, error) {
	row, err := encodeObject(value, fieldPath)
	if err != nil {
		return nil, &schema.ErrorResponse{
			Message: "failed to evaluate nested field",
			Details: map[string]any{
				"reason": fmt.Sprintf("expected object, got %s", reflect.ValueOf(value).Kind()),
				"path":   fieldPath,
			},
		}
	}

	return evalObjectWithColumnSelection(fields.Fields, row, fieldPath)
}

// EvalNestedColumnArrayIntoSlice evaluate and prune nested fields from array without relationship.
func EvalNestedColumnArrayIntoSlice[T any](fields *schema.NestedArray, value []T) (any, error) {
	return evalNestedColumnArrayIntoSlice(fields, value, "")
}

func evalNestedColumnArrayIntoSlice[T any](
	fields *schema.NestedArray,
	value []T,
	fieldPath string,
) (any, error) {
	array, err := encodeObjectSlice(value, fieldPath)
	if err != nil {
		return nil, err
	}

	result := []any{}

	for i, item := range array {
		val, err := evalNestedColumnFields(fields.Fields, item, fmt.Sprintf("%s[%d]", fieldPath, i))
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

// EvalNestedColumnArray evaluate and prune nested fields from array without relationship.
func EvalNestedColumnArray(fields *schema.NestedArray, value any) (any, error) {
	return evalNestedColumnArray(fields, value, "")
}

func evalNestedColumnArray(fields *schema.NestedArray, value any, fieldPath string) (any, error) {
	array, err := encodeObjects(value, fieldPath)
	if err != nil {
		return nil, err
	}

	result := []any{}

	for i, item := range array {
		val, err := evalNestedColumnFields(fields.Fields, item, fmt.Sprintf("%s[%d]", fieldPath, i))
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

// EvalNestedColumnFields evaluate and prune nested fields without relationship.
func EvalNestedColumnFields(fields schema.NestedField, value any) (any, error) {
	return evalNestedColumnFields(fields, value, "")
}

func evalNestedColumnFields(fields schema.NestedField, value any, fieldPath string) (any, error) {
	if IsNil(value) {
		return nil, nil
	}

	iNestedField, err := fields.InterfaceT()
	switch nf := iNestedField.(type) {
	case *schema.NestedObject:
		return evalNestedColumnObject(nf, value, fieldPath)
	case *schema.NestedArray:
		return evalNestedColumnArray(nf, value, fieldPath)
	default:
		return nil, err
	}
}

// EncodeObjectsWithColumnSelection encodes objects with column fields selection without relationship.
func EncodeObjectsWithColumnSelection[T any](
	fields map[string]schema.Field,
	data []T,
) ([]map[string]any, error) {
	return encodeObjectsWithColumnSelection(fields, data, "")
}

func encodeObjectsWithColumnSelection[T any](
	fields map[string]schema.Field,
	data []T,
	fieldPath string,
) ([]map[string]any, error) {
	objects, err := encodeObjectSlice[T](data, fieldPath)
	if err != nil {
		return nil, err
	}

	return evalObjectsWithColumnSelection(fields, objects, fieldPath)
}

// EncodeObjectWithColumnSelection encodes an object with column fields selection without relationship.
func EncodeObjectWithColumnSelection[T any](
	fields map[string]schema.Field,
	data T,
) (map[string]any, error) {
	return encodeObjectWithColumnSelection(fields, data, "")
}

func encodeObjectWithColumnSelection[T any](
	fields map[string]schema.Field,
	data T,
	fieldPath string,
) (map[string]any, error) {
	objects, err := encodeObject(data, fieldPath)
	if err != nil {
		return nil, err
	}

	return evalObjectWithColumnSelection(fields, objects, fieldPath)
}

// EvalObjectsWithColumnSelection evaluate and prune column fields of array objects without relationship.
func EvalObjectsWithColumnSelection(
	fields map[string]schema.Field,
	data []map[string]any,
) ([]map[string]any, error) {
	return evalObjectsWithColumnSelection(fields, data, "")
}

func evalObjectsWithColumnSelection(
	fields map[string]schema.Field,
	data []map[string]any,
	fieldPath string,
) ([]map[string]any, error) {
	results := make([]map[string]any, len(data))

	for i, item := range data {
		result, err := evalObjectWithColumnSelection(
			fields,
			item,
			fmt.Sprintf("%s[%d]", fieldPath, i),
		)
		if err != nil {
			return nil, err
		}

		results[i] = result
	}

	return results, nil
}

// EvalObjectWithColumnSelection evaluate and prune column fields without relationship.
func EvalObjectWithColumnSelection(
	fields map[string]schema.Field,
	data map[string]any,
) (map[string]any, error) {
	return evalObjectWithColumnSelection(fields, data, "")
}

func evalObjectWithColumnSelection(
	fields map[string]schema.Field,
	data map[string]any,
	fieldPath string,
) (map[string]any, error) {
	if len(fields) == 0 {
		return data, nil
	}

	output := make(map[string]any)

	for key, field := range fields {
		switch fi := field.Interface().(type) {
		case *schema.ColumnField:
			if col, ok := data[fi.Column]; ok {
				if fi.Fields != nil {
					nestedValue, err := evalNestedColumnFields(fi.Fields, col, fmt.Sprintf("%s.%s", fieldPath, key))
					if err != nil {
						return nil, err
					}

					output[key] = nestedValue
				} else {
					output[key] = col
				}
			} else {
				output[key] = nil
			}
		case *schema.RelationshipField:
			return nil, &schema.ErrorResponse{
				Message: "failed to evaluate object field",
				Details: map[string]any{
					"reason": "unsupported relationship field",
					"path":   fmt.Sprintf("%s.%s", fieldPath, key),
				},
			}
		default:
			return nil, &schema.ErrorResponse{
				Message: "failed to evaluate object field",
				Details: map[string]any{
					"reason": "invalid column field",
					"path":   fmt.Sprintf("%s.%s", fieldPath, key),
				},
			}
		}
	}

	return output, nil
}

// ResolveArgumentVariables resolve variables in arguments if exist.
// Deprecated: use ResolveArguments instead.
func ResolveArgumentVariables(
	arguments map[string]schema.Argument,
	variables map[string]any,
) (map[string]any, error) {
	return ResolveArguments(arguments, variables)
}

// ResolveArguments resolve variables into request arguments if exist.
func ResolveArguments(
	arguments map[string]schema.Argument,
	variables map[string]any,
) (map[string]any, error) {
	results := make(map[string]any)

	for key, argument := range arguments {
		value, err := ResolveArgument(argument, variables)
		if err != nil {
			return nil, schema.UnprocessableContentError(
				fmt.Sprintf("failed to resolve argument %s: %s", key, err),
				map[string]any{
					"path": "." + key,
				},
			)
		}

		results[key] = value
	}

	return results, nil
}

// ResolveArgument resolve variables into the argument if exist.
func ResolveArgument(argument schema.Argument, variables map[string]any) (any, error) {
	argT, err := argument.InterfaceT()
	if err != nil {
		return nil, err
	}

	switch arg := argT.(type) {
	case *schema.ArgumentLiteral:
		return arg.Value, nil
	case *schema.ArgumentVariable:
		value, ok := variables[arg.Name]
		if !ok {
			return nil, fmt.Errorf("variable `%s` does not exist", arg.Name)
		}

		return value, nil
	default:
		return nil, fmt.Errorf("unsupported argument type: %+v", arg)
	}
}

// EvalFunctionSelectionFieldValue evaluates the __value field in a function query
// According to the [NDC spec], selection fields of the function type must follow this structure:
//
//	{
//		"fields": {
//			"__value": {
//				"type": "column",
//				"column": "__value",
//				"fields": {
//					"type": "object",
//					"fields": {
//						"fieldA": { "type": "column", "column": "fieldA", "fields": null }
//					} // or null
//				}
//			}
//		}
//	}
//
// [NDC spec]: https://hasura.github.io/ndc-spec/specification/queries/functions.html
func EvalFunctionSelectionFieldValue(request *schema.QueryRequest) (schema.NestedField, error) {
	if len(request.Query.Fields) == 0 {
		return nil, errors.New(errFunctionValueFieldRequired)
	}

	valueField, ok := request.Query.Fields["__value"]
	if !ok {
		return nil, errors.New(errFunctionValueFieldRequired)
	}

	valueColumn, err := valueField.AsColumn()
	if err != nil {
		return nil, schema.UnprocessableContentError(fmt.Sprintf("__value: %s", err), nil)
	}

	if valueColumn.Column != "__value" {
		return nil, errors.New(errFunctionValueFieldRequired)
	}

	return valueColumn.Fields, nil
}

// MergeSchemas merge multiple connector schemas into one schema.
func MergeSchemas(schemas ...*schema.SchemaResponse) (*schema.SchemaResponse, []error) {
	var errs []error

	result := schema.SchemaResponse{
		ObjectTypes: schema.SchemaResponseObjectTypes{},
		ScalarTypes: schema.SchemaResponseScalarTypes{},
	}
	collectionMap := map[string]schema.CollectionInfo{}
	functionMap := map[string]schema.FunctionInfo{}
	procedureMap := map[string]schema.ProcedureInfo{}

	for _, s := range schemas {
		if s == nil {
			continue
		}

		for _, col := range s.Collections {
			if _, ok := collectionMap[col.Name]; ok {
				errs = append(errs, fmt.Errorf("collection %s exists", col.Name))
			}

			collectionMap[col.Name] = col
		}

		for _, fn := range s.Functions {
			if _, ok := functionMap[fn.Name]; ok {
				errs = append(errs, fmt.Errorf("function %s exists", fn.Name))
			}

			functionMap[fn.Name] = fn
		}

		for _, fn := range s.Procedures {
			if _, ok := procedureMap[fn.Name]; ok {
				errs = append(errs, fmt.Errorf("procedure %s exists", fn.Name))
			}

			procedureMap[fn.Name] = fn
		}

		result.Collections = GetSortedValuesByKey(collectionMap)
		result.Functions = GetSortedValuesByKey(functionMap)
		result.Procedures = GetSortedValuesByKey(procedureMap)

		for k, obj := range s.ObjectTypes {
			if _, ok := result.ObjectTypes[k]; ok {
				errs = append(errs, fmt.Errorf("object %s exists", k))
			}

			result.ObjectTypes[k] = obj
		}

		for k, sl := range s.ScalarTypes {
			if _, ok := result.ScalarTypes[k]; ok {
				errs = append(errs, fmt.Errorf("scalar %s exists", k))
			}

			result.ScalarTypes[k] = sl
		}
	}

	return &result, errs
}

// CommandSelectionFieldFromContext gets the command's nested selection field from context.
func CommandSelectionFieldFromContext(ctx context.Context) schema.NestedField {
	value := ctx.Value(CommandSelectionFieldKey)
	if value != nil {
		if selection, ok := value.(schema.NestedField); ok {
			return selection
		}
	}

	return nil
}
