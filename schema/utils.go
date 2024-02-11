package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
}

// ToAnySlice converts a typed slice to any slice
func ToAnySlice[V any](slice []V) []any {
	results := make([]any, len(slice))
	for i, v := range slice {
		results[i] = v
	}
	return results
}

// Index returns the index of the first occurrence of item in slice,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// Contains checks whether the value is present in slice.
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

func getStringValueByKey(collection map[string]any, key string) string {
	if collection == nil {
		return ""
	}

	anyValue, ok := collection[key]
	if !ok || anyValue == nil {
		return ""
	}

	if arg, ok := anyValue.(string); ok {
		return arg
	}

	return ""
}

func unmarshalStringFromJsonMap(collection map[string]json.RawMessage, key string, required bool) (string, error) {

	emptyFn := func() (string, error) {
		if !required {
			return "", nil
		}

		return "", errors.New("required")
	}

	if collection == nil {
		return emptyFn()
	}

	rawValue, ok := collection[key]
	if !ok || len(rawValue) == 0 {
		return emptyFn()
	}

	var result string
	if err := json.Unmarshal(rawValue, &result); err != nil {
		return "", err
	}

	if result == "" {
		return emptyFn()
	}

	return result, nil
}

// EncodeRow encodes an object row to a map[string]any, using json tag to convert object keys
func EncodeRow(row any) (map[string]any, error) {
	if row == nil {
		return nil, errors.New("expected object fields, got nil")
	}

	var outputMap map[string]any
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &outputMap,
		TagName: "json",
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(row); err != nil {
		return nil, err
	}

	return outputMap, nil
}

// PruneFields prune unnecessary fields from selection
func PruneFields(fields map[string]Field, result any) (map[string]any, error) {
	outputMap, err := EncodeRow(result)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return outputMap, nil
	}

	output := make(map[string]any)
	for key, field := range fields {
		f, err := field.Interface()
		switch fi := f.(type) {
		case *ColumnField:
			if col, ok := outputMap[fi.Column]; ok {
				output[fi.Column] = col
			} else {
				output[fi.Column] = nil
			}
		case *RelationshipField:
			return nil, fmt.Errorf("unsupported relationship field,  %s", key)
		default:
			return nil, err
		}
	}

	return output, nil
}

// ResolveArguments resolve variables in arguments and map them to struct
func ResolveArguments[R any](arguments map[string]Argument, variables map[string]any) (*R, error) {
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
func ResolveArgumentVariables(arguments map[string]Argument, variables map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	for key, arg := range arguments {
		switch arg.Type {
		case ArgumentTypeLiteral:
			results[key] = arg.Value
		case ArgumentTypeVariable:
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
