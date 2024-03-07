package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hasura/ndc-sdk-go/schema"
)

const (
	errFunctionValueFieldRequired = "__value field is required in query function type"
)

// Scalar abstracts a scalar interface to determine when evaluating
type Scalar interface {
	ScalarName() string
}

// EvalNestedColumnObject evaluate and prune nested fields from an object without relationship
func EvalNestedColumnObject(fields *schema.NestedObject, value any) (any, error) {
	row, err := EncodeObject(value)
	if err != nil {
		return nil, fmt.Errorf("expected object, got %s", reflect.ValueOf(value).Kind())
	}

	return EvalObjectWithColumnSelection(fields.Fields, row)
}

// EvalNestedColumnArrayIntoSlice evaluate and prune nested fields from array without relationship
func EvalNestedColumnArrayIntoSlice[T any](fields *schema.NestedArray, value []T) (any, error) {
	array, err := EncodeObjectSlice(value)
	if err != nil {
		return nil, err
	}

	result := []any{}
	for _, item := range array {
		val, err := EvalNestedColumnFields(fields.Fields, item)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

// EvalNestedColumnArray evaluate and prune nested fields from array without relationship
func EvalNestedColumnArray(fields *schema.NestedArray, value any) (any, error) {
	array, err := EncodeObjects(value)
	if err != nil {
		return nil, err
	}

	result := []any{}
	for _, item := range array {
		val, err := EvalNestedColumnFields(fields.Fields, item)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

// EvalNestedColumnFields evaluate and prune nested fields without relationship
func EvalNestedColumnFields(fields schema.NestedField, value any) (any, error) {
	if IsNil(value) {
		return nil, nil
	}

	iNestedField, err := fields.InterfaceT()
	switch nf := iNestedField.(type) {
	case *schema.NestedObject:
		return EvalNestedColumnObject(nf, value)
	case *schema.NestedArray:
		return EvalNestedColumnArray(nf, value)
	default:
		return nil, err
	}
}

// EncodeObjectsWithColumnSelection encodes objects with column fields selection without relationship
func EncodeObjectsWithColumnSelection[T any](fields map[string]schema.Field, data []T) ([]map[string]any, error) {
	objects, err := EncodeObjectSlice[T](data)
	if err != nil {
		return nil, err
	}
	return EvalObjectsWithColumnSelection(fields, objects)
}

// EncodeObjectWithColumnSelection encodes an object with column fields selection without relationship
func EncodeObjectWithColumnSelection[T any](fields map[string]schema.Field, data T) (map[string]any, error) {
	objects, err := EncodeObject(data)
	if err != nil {
		return nil, err
	}
	return EvalObjectWithColumnSelection(fields, objects)
}

// EvalObjectsWithColumnSelection evaluate and prune column fields of array objects without relationship
func EvalObjectsWithColumnSelection(fields map[string]schema.Field, data []map[string]any) ([]map[string]any, error) {
	results := make([]map[string]any, len(data))
	for i, item := range data {
		result, err := EvalObjectWithColumnSelection(fields, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	return results, nil
}

// EvalObjectWithColumnSelection evaluate and prune column fields without relationship
func EvalObjectWithColumnSelection(fields map[string]schema.Field, data map[string]any) (map[string]any, error) {
	if len(fields) == 0 {
		return data, nil
	}

	output := make(map[string]any)
	for key, field := range fields {
		switch fi := field.Interface().(type) {
		case *schema.ColumnField:
			if col, ok := data[fi.Column]; ok {
				if fi.Fields != nil {
					nestedValue, err := EvalNestedColumnFields(fi.Fields, col)
					if err != nil {
						return nil, err
					}
					output[fi.Column] = nestedValue
				} else {
					output[fi.Column] = col
				}
			} else {
				output[fi.Column] = nil
			}
		case *schema.RelationshipField:
			return nil, fmt.Errorf("unsupported relationship field,  %s", key)
		default:
			return nil, fmt.Errorf("invalid column field, %s", key)
		}
	}

	return output, nil
}

// ResolveArgumentVariables resolve variables in arguments if exist
func ResolveArgumentVariables(arguments map[string]schema.Argument, variables map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	for key, arg := range arguments {
		switch arg.Type {
		case schema.ArgumentTypeLiteral:
			results[key] = arg.Value
		case schema.ArgumentTypeVariable:
			value, ok := variables[arg.Name]
			if !ok {
				return nil, fmt.Errorf("variable %s not found", arg.Name)
			}
			results[key] = value
		default:
			return nil, fmt.Errorf("unsupported argument type: %s", arg.Type)
		}
	}

	return results, nil
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
		return nil, schema.BadRequestError(fmt.Sprintf("__value: %s", err), nil)
	}
	if valueColumn.Column != "__value" {
		return nil, errors.New(errFunctionValueFieldRequired)
	}
	return valueColumn.Fields, nil
}
