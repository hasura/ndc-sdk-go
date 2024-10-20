package schema

import (
	"encoding/json"
	"errors"
	"reflect"
)

// isNil a safe function to check null value
func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func isNullJSON(value []byte) bool {
	return len(value) == 0 || string(value) == "null"
}

func getStringValueByKey(collection map[string]any, key string) string {
	if len(collection) == 0 {
		return ""
	}

	anyValue, ok := collection[key]
	if !ok || isNil(anyValue) {
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

	if len(collection) == 0 {
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
