package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
)

// ToMaps encodes objects to a slice of map[string]any, using json tag to convert object keys
func ToMaps[T MapEncoder](inputs []T) []map[string]any {
	var results []map[string]any
	for _, item := range inputs {
		results = append(results, item.ToMap())
	}
	return results
}

// EncodeObject encodes an unknown type to a map[string]any, using json tag to convert object keys
func EncodeObject(input any) (map[string]any, error) {
	if IsNil(input) {
		return nil, nil
	}
	return encodeObject(input)
}

// EncodeObjectSlice encodes array of unknown type to map[string]any slice, using json tag to convert object keys
func EncodeObjectSlice[T any](input []T) ([]map[string]any, error) {
	results := make([]map[string]any, len(input))
	for i, item := range input {
		result, err := EncodeObject(item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

func encodeObject(input any) (map[string]any, error) {
	switch value := input.(type) {
	case map[string]any:
		return value, nil
	case MapEncoder:
		return value.ToMap(), nil
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, time.Time, time.Duration, time.Ticker, *bool, *string, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, *complex64, *complex128, *time.Time, *time.Duration, *time.Ticker, []bool, []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128, []time.Time, []time.Duration, []time.Ticker:
		return nil, schema.BadRequestError("failed to encode object", map[string]any{
			"value": input,
		})
	default:
		inputValue := reflect.ValueOf(input)
		kind := inputValue.Kind()
		switch kind {
		case reflect.Pointer:
			return encodeObject(inputValue.Elem().Interface())
		case reflect.Struct:
			return encodeStruct(inputValue), nil
		default:
			return nil, schema.BadRequestError(fmt.Sprintf("failed to encode object, got %s", kind.String()), map[string]any{
				"value": input,
			})
		}
	}
}

func encodeField(input reflect.Value) (any, bool) {
	switch input.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return input.Interface(), true
	case reflect.Struct:
		inputType := input.Type()
		switch inputType.PkgPath() {
		case "time":
			switch inputType.Name() {
			case "Time":
				return input.Interface(), true
			}
		default:
			toMap := input.MethodByName("ToMap")
			if toMap.IsValid() {
				results := toMap.Call([]reflect.Value{})
				if len(results) == 0 {
					return nil, true
				}
				return results[0].Interface(), true
			}
			return encodeStruct(input), true
		}
	case reflect.Pointer:
		return encodeField(input.Elem())
	case reflect.Array, reflect.Slice:
		len := input.Len()
		var result []any
		for i := 0; i < len; i++ {
			item, ok := encodeField(input.Index(i))
			if ok {
				result = append(result, item)
			}
		}
		return result, true
	}
	return nil, false
}

func encodeStruct(input reflect.Value) map[string]any {
	result := make(map[string]any)
	for i := 0; i < input.NumField(); i++ {
		fieldValue := input.Field(i)
		fieldType := input.Type().Field(i)
		fieldJSONTag := fieldType.Tag.Get("json")
		fieldName := fieldType.Name
		if fieldJSONTag == "-" {
			continue
		}
		if fieldJSONTag != "" {
			fieldName = strings.Split(fieldJSONTag, ",")[0]
		}

		value, ok := encodeField(fieldValue)
		if ok {
			result[fieldName] = value
		}
	}
	return result
}

// EncodeObjects encodes an object rows to a slice of map[string]any, using json tag to convert object keys
func EncodeObjects(input any) ([]map[string]any, error) {
	if IsNil(input) {
		return nil, nil
	}
	inputValue := reflect.ValueOf(input)
	if inputValue.Kind() != reflect.Array && inputValue.Kind() != reflect.Slice {
		return nil, schema.BadRequestError("failed to encode array objects", map[string]any{
			"value": input,
		})
	}
	len := inputValue.Len()
	results := make([]map[string]any, len)

	for i := 0; i < len; i++ {
		item, err := EncodeObject(inputValue.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		results[i] = item
	}
	return results, nil
}

// IsNil a safe function to check null value
func IsNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Ptr && v.IsNil()
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
