package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
)

// MapEncoder abstracts a type with the ToMap method to encode type to map
type MapEncoder interface {
	ToMap() map[string]any
}

// EncodeMap encodes an object to a map[string]any, using json tag to convert object keys
func EncodeMap[T MapEncoder](input T) map[string]any {
	if IsNil(input) {
		return nil
	}
	return input.ToMap()
}

// EncodeMaps encode objects to a slice of map[string]any, using json tag to convert object keys
func EncodeMaps[T MapEncoder](inputs []T) []map[string]any {
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

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
}

func encodeObject(input any) (map[string]any, error) {
	switch value := input.(type) {
	case map[string]any:
		return value, nil
	case MapEncoder:
		return value.ToMap(), nil
	case Scalar:
		return nil, schema.BadRequestError("cannot encode scalar to object", map[string]any{
			"value": input,
		})
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
	case reflect.Complex64, reflect.Complex128:
		return nil, false
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
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
			// determine if the type implements the Scalar interface
			if input.MethodByName("ScalarName").IsValid() {
				return input.Interface(), true
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
