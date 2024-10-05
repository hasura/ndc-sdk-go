package utils

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// GetDefault returns the value or default one if value is empty
func GetDefault[T comparable](value T, defaultValue T) T {
	var empty T
	if value == empty {
		return defaultValue
	}
	return value
}

// GetDefaultPtr returns the first pointer or default one if GetDefaultPtr is nil
func GetDefaultPtr[T any](value *T, defaultValue *T) *T {
	if value == nil {
		return defaultValue
	}
	return value
}

// GetDefaultValuePtr return the value of pointer or default one if the value of pointer is null or empty
func GetDefaultValuePtr[T comparable](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	var empty T
	if *value == empty {
		return defaultValue
	}
	return *value
}

// GetKeys gets keys of a map
func GetKeys[K cmp.Ordered, V any](input map[K]V) []K {
	var results []K
	for key := range input {
		results = append(results, key)
	}
	return results
}

// GetSortedKeys gets keys of a map and sorts them
func GetSortedKeys[K cmp.Ordered, V any](input map[K]V) []K {
	results := GetKeys(input)
	slices.Sort(results)
	return results
}

// GetSortedValuesByKey gets values of a map and sorts by keys
func GetSortedValuesByKey[K cmp.Ordered, V any](input map[K]V) []V {
	if len(input) == 0 {
		return []V{}
	}
	keys := GetSortedKeys(input)
	results := make([]V, len(input))
	for i, k := range keys {
		results[i] = input[k]
	}
	return results
}

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
}

// ToPtrs converts the value slice to pointer slice
func ToPtrs[T any](input []T) []*T {
	results := make([]*T, len(input))
	for i := range input {
		v := input[i]
		results[i] = &v
	}
	return results
}

// PointersToValues converts the pointer slice to value slice
func PointersToValues[T any](input []*T) ([]T, error) {
	results := make([]T, len(input))
	for i, v := range input {
		if IsNil(v) {
			return nil, fmt.Errorf("element at %d must not be nil", i)
		}
		results[i] = *v
	}
	return results, nil
}

// UnwrapPointerFromReflectValue unwraps pointers from the reflect value
func UnwrapPointerFromReflectValue(reflectValue reflect.Value) (reflect.Value, bool) {
	for reflectValue.Kind() == reflect.Pointer {
		if reflectValue.IsNil() {
			return reflectValue, false
		}
		reflectValue = reflectValue.Elem()
	}
	kind := reflectValue.Kind()
	if (kind == reflect.Slice || kind == reflect.Interface || kind == reflect.Map) && reflectValue.IsNil() {
		return reflectValue, false
	}
	return reflectValue, true
}

// UnwrapPointerFromAny unwraps pointers from the input any type
func UnwrapPointerFromAny(value any) (any, bool) {
	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, false
	}
	return reflectValue.Interface(), true
}

// IsDebug checks if the log level is debug
func IsDebug(logger *slog.Logger) bool {
	return logger.Enabled(context.TODO(), slog.LevelDebug)
}

// MergeMap merges two value maps into one
func MergeMap[K comparable, V any](dest map[K]V, src map[K]V) map[K]V {
	result := dest
	if result == nil {
		result = map[K]V{}
	}
	for k, v := range src {
		result[k] = v
	}
	return result
}

// ParseIntMapFromString parses an integer map from a string with format:
//
//	<key1>=<value1>;<key2>=<value2>
func ParseIntMapFromString(input string) (map[string]int, error) {
	result := make(map[string]int)
	if input == "" {
		return result, nil
	}
	rawItems := strings.Split(input, ";")
	for _, rawItem := range rawItems {
		keyValue := strings.Split(rawItem, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid int map string %s, expected <key1>=<value1>;<key2>=<value2>", input)
		}
		value, err := strconv.ParseInt(keyValue[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %s in item %s", keyValue[1], rawItem)
		}
		result[keyValue[0]] = int(value)
	}

	return result, nil
}
