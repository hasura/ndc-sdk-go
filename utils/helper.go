package utils

import "fmt"

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
