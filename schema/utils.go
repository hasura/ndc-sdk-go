package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// isNil a safe function to check null value.
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

func getStringValueByKey(collection map[string]any, key string) (string, error) {
	if len(collection) == 0 {
		return "", nil
	}

	anyValue, ok := collection[key]
	if !ok || anyValue == nil {
		return "", nil
	}

	strValue, err := decodeStringValue(anyValue)
	if err != nil {
		return "", err
	}

	if strValue != nil {
		return *strValue, nil
	}

	return "", nil
}

func decodeStringValue(anyValue any) (*string, error) {
	if arg, ok := anyValue.(string); ok {
		return &arg, nil
	}

	if arg, ok := anyValue.(*string); ok {
		return arg, nil
	}

	return nil, fmt.Errorf("expected string, got %v", anyValue)
}

func getStringSliceByKey(collection map[string]any, key string) ([]string, error) { //nolint:unparam
	if len(collection) == 0 {
		return nil, nil
	}

	anyValue, ok := collection[key]
	if !ok || anyValue == nil {
		return nil, nil
	}

	if args, ok := anyValue.([]string); ok {
		return args, nil
	}

	args, ok := anyValue.([]any)
	if !ok {
		return nil, fmt.Errorf("expected []string, got %v", anyValue)
	}

	results := make([]string, len(args))

	for i, item := range args {
		str, err := decodeStringValue(item)
		if err != nil {
			return nil, fmt.Errorf("failed to parse element at %d: %w", i, err)
		}

		if str == nil {
			return nil, fmt.Errorf("failed to parse element at %d: string value must not be null", i)
		}

		results[i] = *str
	}

	return results, nil
}

func getArgumentMapByKey(collection map[string]any, key string) (map[string]Argument, error) { //nolint:unparam
	rawArguments, ok := collection[key]
	if !ok || rawArguments == nil {
		return nil, nil
	}

	args, ok := rawArguments.(map[string]Argument)
	if ok {
		return args, nil
	}

	rawArgumentsMap, ok := rawArguments.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected object, got %v", rawArguments)
	}

	if rawArgumentsMap == nil {
		return nil, nil
	}

	arguments := map[string]Argument{}

	for key, rawArg := range rawArgumentsMap {
		argMap, ok := rawArg.(map[string]any)
		if !ok || argMap == nil {
			return nil, fmt.Errorf("field %s in map[string]Argument: expected object, got %v", key, argMap)
		}

		argument := Argument{}

		if err := argument.FromValue(argMap); err != nil {
			return nil, fmt.Errorf("field %s in map[string]Argument: %w", key, err)
		}

		arguments[key] = argument
	}

	return arguments, nil
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
