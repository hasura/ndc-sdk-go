package schema

import (
	"encoding/json"
	"errors"
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

// ToRows converts a typed slice to Row slice
func ToRows[V any](slice []V) []Row {
	results := make([]Row, len(slice))
	for i, v := range slice {
		results[i] = v
	}
	return results
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
