package utils

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hasura/ndc-sdk-go/schema"
)

// EncodeRow encodes an object row to a map[string]any, using json tag to convert object keys
func EncodeRow(row any) (map[string]any, error) {
	value, ok := row.(map[string]any)
	if ok {
		return value, nil
	}
	return encodeRows[map[string]any](row)
}

// EncodeRows encodes an object rows to a slice of map[string]any, using json tag to convert object keys
func EncodeRows(rows any) ([]map[string]any, error) {
	return encodeRows[[]map[string]any](rows)
}

func encodeRows[R any](rows any) (R, error) {
	var result R
	if rows == nil {
		return result, errors.New("expected object fields, got nil")
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &result,
		TagName:    "json",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(timeDecodeHookFunc()),
	})
	if err != nil {
		return result, err
	}
	err = decoder.Decode(rows)

	return result, err
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// EvalNestedColumnFields evaluate and prune nested fields without relationship
func EvalNestedColumnFields(fields schema.NestedField, value any) (any, error) {
	if isNil(value) {
		return nil, nil
	}

	iNestedField, err := fields.InterfaceT()
	switch nf := iNestedField.(type) {
	case *schema.NestedObject:
		row, err := EncodeRow(value)
		if err != nil {
			return nil, fmt.Errorf("expected object, got %s", reflect.ValueOf(value).Kind())
		}

		return EvalColumnFields(nf.Fields, row)
	case *schema.NestedArray:
		array, err := EncodeRows(value)
		if err != nil {
			return nil, err
		}

		result := []any{}
		for _, item := range array {
			val, err := EvalNestedColumnFields(nf.Fields, item)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		}
		return result, nil
	default:
		return nil, err
	}
}

// EvalColumnFields evaluate and prune column fields without relationship
func EvalColumnFields(fields map[string]schema.Field, result any) (map[string]any, error) {
	outputMap, err := EncodeRow(result)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return outputMap, nil
	}

	output := make(map[string]any)
	for key, field := range fields {
		switch fi := field.Interface().(type) {
		case *schema.ColumnField:
			if col, ok := outputMap[fi.Column]; ok {
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

// ResolveArguments resolve variables in arguments and map them to struct
func ResolveArguments[R any](arguments map[string]schema.Argument, variables map[string]any) (*R, error) {
	resolvedArgs, err := ResolveArgumentVariables(arguments, variables)
	if err != nil {
		return nil, err
	}

	var result R

	if err = mapstructure.Decode(resolvedArgs, &result); err != nil {
		return nil, err
	}

	return &result, nil
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

// timeDecodeHookFunc creates a mapstructure decode hook function to decode time.Time
func timeDecodeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if ty, ok := data.(*time.Time); ok {
			return *ty, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}
