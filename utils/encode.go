package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/v2/schema"
)

// MapEncoder abstracts a type with the ToMap method to encode type to map.
type MapEncoder interface {
	ToMap() map[string]any
}

// EncodeObject encodes an unknown type to a map[string]any, using json tag to convert object keys.
func EncodeObject(input any) (map[string]any, error) {
	if input == nil {
		return nil, nil
	}

	return encodeObject(input, "")
}

// EncodeObjectSlice encodes an array of unknown type to map[string]any slice, using json tag to convert object keys.
func EncodeObjectSlice[T any](input []T) ([]map[string]any, error) {
	return encodeObjectSlice(input, "")
}

// EncodeNullableObjectSlice encodes the pointer array of unknown type to map[string]any slice, using json tag to convert object keys.
func EncodeNullableObjectSlice[T any](inputs *[]T) ([]map[string]any, error) {
	if inputs == nil {
		return nil, nil
	}

	return encodeObjectSlice(*inputs, "")
}

func encodeObjectSlice[T any](input []T, fieldPath string) ([]map[string]any, error) {
	results := make([]map[string]any, len(input))

	for i, item := range input {
		result, err := encodeObject(item, fmt.Sprintf("%s[%d]", fieldPath, i))
		if err != nil {
			return nil, err
		}

		results[i] = result
	}

	return results, nil
}

func encodeObject(input any, fieldPath string) (map[string]any, error) {
	switch value := input.(type) {
	case map[string]any:
		return value, nil
	case MapEncoder:
		if value == nil {
			return nil, nil
		}

		return value.ToMap(), nil
	case Scalar:
		return nil, &schema.ErrorResponse{
			Message: "cannot encode scalar to object",
			Details: map[string]any{
				"reason": fmt.Sprintf("expected object, got %s", reflect.TypeOf(input).Kind()),
				"path":   fieldPath,
			},
		}
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, time.Time, time.Duration, time.Ticker, *bool, *string, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, *complex64, *complex128, *time.Time, *time.Duration, *time.Ticker, []bool, []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128, []time.Time, []time.Duration, []time.Ticker:
		return nil, &schema.ErrorResponse{
			Message: "cannot encode scalar to object",
			Details: map[string]any{
				"reason": fmt.Sprintf("expected object, got %s", reflect.TypeOf(input).Kind()),
				"path":   fieldPath,
			},
		}
	default:
		return encodeObjectReflection(reflect.ValueOf(input), fieldPath)
	}
}

func encodeObjectReflection(inputValue reflect.Value, fieldPath string) (map[string]any, error) {
	kind := inputValue.Kind()

	switch kind {
	case reflect.Pointer:
		v, ok := UnwrapPointerFromReflectValue(inputValue)
		if !ok {
			return nil, nil
		}

		return encodeObjectReflection(v, fieldPath)
	case reflect.Struct:
		return encodeStruct(inputValue), nil
	case reflect.Map:
		result := make(map[string]any)
		iter := inputValue.MapRange()

		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			key, err := DecodeString(k)
			if err != nil {
				return nil, &schema.ErrorResponse{
					Message: fmt.Sprintf(
						"cannot encode map; the object key must be a string, got: %v",
						k.Interface(),
					),
					Details: map[string]any{
						"path": fieldPath,
					},
				}
			}

			value, ok := encodeField(v)
			if ok {
				result[key] = value
			}
		}

		return result, nil
	default:
		return nil, &schema.ErrorResponse{
			Message: "cannot encode object",
			Details: map[string]any{
				"reason": fmt.Sprintf("expected object, got %s", kind),
				"path":   fieldPath,
			},
		}
	}
}

// encode fields with [type representation spec]
//
// [type representation spec]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#new-representations
func encodeField(input reflect.Value) (any, bool) {
	switch input.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return nil, false
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
		reflect.Map,
		reflect.Int64:
		return input.Interface(), true
	case reflect.Uint64:
		return strconv.FormatUint(input.Uint(), 10), true
	case reflect.Struct:
		inputType := input.Type()

		switch inputType.PkgPath() {
		case "time":
			switch inputType.Name() {
			case "Time":
				return input.Interface(), true
			default:
				return nil, false
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

			// determine if the type implements the Scalar interface
			if input.MethodByName("ScalarName").IsValid() {
				return input.Interface(), true
			}

			return encodeStruct(input), true
		}
	case reflect.Pointer:
		if input.IsNil() {
			return nil, true
		}

		return encodeField(input.Elem())
	case reflect.Array, reflect.Slice:
		if input.IsNil() {
			return nil, true
		}

		valueLength := input.Len()
		result := []any{}

		for i := range valueLength {
			item, ok := encodeField(input.Index(i))
			if ok {
				result = append(result, item)
			}
		}

		return result, true
	default:
	}

	return nil, false
}

func encodeStruct(input reflect.Value) map[string]any {
	result := make(map[string]any)

	for i := range input.NumField() {
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

// EncodeObjects encodes an object rows to a slice of map[string]any, using json tag to convert object keys.
func EncodeObjects(input any) ([]map[string]any, error) {
	return encodeObjects(input, "")
}

func encodeObjects(input any, fieldPath string) ([]map[string]any, error) {
	inputValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(input))
	if !ok {
		return nil, nil
	}

	inputKind := inputValue.Kind()
	if inputKind != reflect.Array && inputKind != reflect.Slice {
		return nil, &schema.ErrorResponse{
			Message: "failed to encode array objects",
			Details: map[string]any{
				"reason": fmt.Sprintf("expected array objects, got %s", inputKind),
				"path":   fieldPath,
			},
		}
	}

	valueLength := inputValue.Len()
	results := make([]map[string]any, valueLength)

	for i := range valueLength {
		item, err := encodeObject(inputValue.Index(i).Interface(), fieldPath)
		if err != nil {
			return nil, err
		}

		results[i] = item
	}

	return results, nil
}
