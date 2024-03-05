package schema

import (
	"encoding/json"
	"errors"
)

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
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
