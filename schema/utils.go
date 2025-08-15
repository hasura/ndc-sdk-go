package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// GetUnderlyingNamedType gets the underlying named type of the input type recursively if exists.
func GetUnderlyingNamedType(input Type) *NamedType {
	if len(input) == 0 {
		return nil
	}

	switch ty := input.Interface().(type) {
	case *NullableType:
		return GetUnderlyingNamedType(ty.UnderlyingType)
	case *ArrayType:
		return GetUnderlyingNamedType(ty.ElementType)
	case *NamedType:
		return ty
	default:
		return nil
	}
}

// UnwrapNullableType unwraps nullable types from the input type recursively.
func UnwrapNullableType(input Type) Type {
	if input == nil || input.IsZero() {
		return nil
	}

	switch ty := input.Interface().(type) {
	case *NullableType:
		return UnwrapNullableType(ty.UnderlyingType)
	default:
		return input
	}
}

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
			return nil, fmt.Errorf(
				"failed to parse element at %d: string value must not be null",
				i,
			)
		}

		results[i] = *str
	}

	return results, nil
}

func getPathElementByKey(
	collection map[string]any,
	key string,
) ([]PathElement, error) {
	if len(collection) == 0 {
		return nil, nil
	}

	anyValue, ok := collection[key]
	if !ok || anyValue == nil {
		return nil, nil
	}

	if args, ok := anyValue.([]PathElement); ok {
		return args, nil
	}

	args, ok := anyValue.([]any)
	if !ok {
		return nil, fmt.Errorf("expected []PathElement, got %v", anyValue)
	}

	results := make([]PathElement, len(args))

	for i, item := range args {
		if elem, ok := item.(PathElement); ok {
			results[i] = elem

			continue
		}

		valueMap, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"failed to parse path element at %d: expected object, got: %v",
				i,
				item,
			)
		}

		if valueMap == nil {
			return nil, fmt.Errorf(
				"failed to parse path element at %d: value must not be null",
				i,
			)
		}

		result := PathElement{}

		if err := result.FromValue(valueMap); err != nil {
			return nil, fmt.Errorf("failed to parse path element at %d: %w", i, err)
		}

		results[i] = result
	}

	return results, nil
}

func getArgumentMapByKey(
	collection map[string]any,
	key string, //nolint:unparam
) (map[string]Argument, error) {
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
			return nil, fmt.Errorf(
				"field %s in map[string]Argument: expected object, got %v",
				key,
				argMap,
			)
		}

		argument := Argument{}

		if err := argument.FromValue(argMap); err != nil {
			return nil, fmt.Errorf("field %s in map[string]Argument: %w", key, err)
		}

		arguments[key] = argument
	}

	return arguments, nil
}

func getRelationshipArgumentMapByKey(
	collection map[string]any,
	key string, //nolint:unparam
) (map[string]RelationshipArgument, error) {
	rawArguments, ok := collection[key]
	if !ok || rawArguments == nil {
		return nil, nil
	}

	args, ok := rawArguments.(map[string]RelationshipArgument)
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

	arguments := map[string]RelationshipArgument{}

	for key, rawArg := range rawArgumentsMap {
		argMap, ok := rawArg.(map[string]any)
		if !ok || argMap == nil {
			return nil, fmt.Errorf(
				"field %s in map[string]Argument: expected object, got %v",
				key,
				argMap,
			)
		}

		argument := RelationshipArgument{}

		if err := argument.FromValue(argMap); err != nil {
			return nil, fmt.Errorf("field %s in map[string]Argument: %w", key, err)
		}

		arguments[key] = argument
	}

	return arguments, nil
}

func unmarshalStringFromJsonMap(
	collection map[string]json.RawMessage,
	key string,
	required bool, //nolint:unparam
) (string, error) {
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
