package utils

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
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
func GetKeys[V any](input map[string]V) []string {
	var results []string
	for key := range input {
		results = append(results, key)
	}
	return results
}

// GetSortedKeys gets keys of a map and sorts them
func GetSortedKeys[V any](input map[string]V) []string {
	results := GetKeys(input)
	sort.Strings(results)
	return results
}

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
}

// ToPtrs converts the value slice to pointer slice
func ToPtrs[T any](input []T) []*T {
	results := make([]*T, len(input))
	for i, v := range input {
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
