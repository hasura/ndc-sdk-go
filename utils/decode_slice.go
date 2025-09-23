package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// DecodeNullableIntSlice decodes an integer pointer slice from an unknown value.
func DecodeNullableIntSlice[T int | int8 | int16 | int32 | int64](value any) (*[]T, error) {
	return decodeNullableNumberSlice(value, decodeNullableIntRefection(convertInt[T]))
}

// DecodeNullableUintSlice decodes an unsigned integer slice from an unknown value.
func DecodeNullableUintSlice[T uint | uint8 | uint16 | uint32 | uint64](value any) (*[]T, error) {
	return decodeNullableNumberSlice(value, decodeNullableIntRefection(convertUint[T]))
}

// DecodeNullableFloatSlice decodes a float slice from an unknown value.
func DecodeNullableFloatSlice[T float32 | float64](value any) (*[]T, error) {
	return decodeNullableNumberSlice(value, DecodeNullableFloatReflection[T])
}

// DecodeIntSlice decodes an integer slice from an unknown value.
func DecodeIntSlice[T int | int8 | int16 | int32 | int64](value any) ([]T, error) {
	results, err := decodeNullableNumberSlice(value, decodeNullableIntRefection(convertInt[T]))
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the int slice must not be null")
	}

	return *results, nil
}

// DecodeUintSlice decodes an unsigned integer slice from an unknown value.
func DecodeUintSlice[T uint | uint8 | uint16 | uint32 | uint64](value any) ([]T, error) {
	results, err := decodeNullableNumberSlice(value, decodeNullableIntRefection(convertUint[T]))
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the uint slice must not be null")
	}

	return *results, nil
}

// DecodeFloatSlice decodes a float slice from an unknown value.
func DecodeFloatSlice[T float32 | float64](value any) ([]T, error) {
	results, err := decodeNullableNumberSlice(value, DecodeNullableFloatReflection[T])
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the float slice must not be null")
	}

	return *results, nil
}

// DecodeNullableIntPtrSlice decodes an integer pointer slice from an unknown value.
func DecodeNullableIntPtrSlice[T int | int8 | int16 | int32 | int64](value any) (*[]*T, error) {
	return decodeNullableNumberPtrSlice(value, DecodeNullableInt[T])
}

// DecodeNullableUintPtrSlice decodes an unsigned integer pointer slice from an unknown value.
func DecodeNullableUintPtrSlice[T uint | uint8 | uint16 | uint32 | uint64](
	value any,
) (*[]*T, error) {
	return decodeNullableNumberPtrSlice(value, DecodeNullableUint[T])
}

// DecodeNullableFloatPtrSlice decodes a float pointer slice from an unknown value.
func DecodeNullableFloatPtrSlice[T float32 | float64](value any) (*[]*T, error) {
	return decodeNullableNumberPtrSlice(value, DecodeNullableFloat[T])
}

// DecodeIntPtrSlice decodes an integer slice from an unknown value.
func DecodeIntPtrSlice[T int | int8 | int16 | int32 | int64](value any) ([]*T, error) {
	results, err := decodeNullableNumberPtrSlice(value, DecodeNullableInt[T])
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the int slice must not be null")
	}

	return *results, nil
}

// DecodeUintPtrSlice decodes an unsigned integer pointer slice from an unknown value.
func DecodeUintPtrSlice[T uint | uint8 | uint16 | uint32 | uint64](value any) ([]*T, error) {
	results, err := decodeNullableNumberPtrSlice(value, DecodeNullableUint[T])
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the uint slice must not be null")
	}

	return *results, nil
}

// DecodeFloatPtrSlice decodes a float pointer slice from an unknown value.
func DecodeFloatPtrSlice[T float32 | float64](value any) ([]*T, error) {
	results, err := decodeNullableNumberPtrSlice(value, DecodeNullableFloat[T])
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the float slice must not be null")
	}

	return *results, nil
}

// GetNullableIntSlice get an integer pointer slice from object by key.
func GetNullableIntSlice[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) (*[]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableIntSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetNullableUintSlice get an unsigned integer pointer slice from object by key.
func GetNullableUintSlice[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) (*[]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableUintSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetNullableFloatSlice get a float pointer slice from object by key.
func GetNullableFloatSlice[T float32 | float64](object map[string]any, key string) (*[]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableFloatSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetIntSlice get an integer slice from object by key.
func GetIntSlice[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeIntSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetIntSliceDefault get an integer slice from object by key.
func GetIntSliceDefault[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableIntSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetUintSlice get an unsigned integer slice from object by key.
func GetUintSlice[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeUintSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetUintSliceDefault get an unsigned integer slice from object by key.
func GetUintSliceDefault[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableUintSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetFloatSlice get a float slice from object by key.
func GetFloatSlice[T float32 | float64](object map[string]any, key string) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeFloatSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetFloatSliceDefault get a float slice from object by key.
func GetFloatSliceDefault[T float32 | float64](object map[string]any, key string) ([]T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableFloatSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableIntPtrSlice get an integer pointer slice from object by key.
func GetNullableIntPtrSlice[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) (*[]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableIntPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetNullableUintPtrSlice get an unsigned integer pointer slice from object by key.
func GetNullableUintPtrSlice[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) (*[]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableUintPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetNullableFloatPtrSlice get a float pointer slice from object by key.
func GetNullableFloatPtrSlice[T float32 | float64](
	object map[string]any,
	key string,
) (*[]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableFloatPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetIntPtrSlice get an integer slice from object by key.
func GetIntPtrSlice[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeIntPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetIntPtrSliceDefault get an integer slice from object by key.
func GetIntPtrSliceDefault[T int | int8 | int16 | int32 | int64](
	object map[string]any,
	key string,
) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableIntPtrSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetUintPtrSlice get an unsigned integer slice from object by key.
func GetUintPtrSlice[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeUintPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetUintPtrSliceDefault get an unsigned integer slice from object by key.
func GetUintPtrSliceDefault[T uint | uint8 | uint16 | uint32 | uint64](
	object map[string]any,
	key string,
) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableUintPtrSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetFloatPtrSlice get a float slice from object by key.
func GetFloatPtrSlice[T float32 | float64](object map[string]any, key string) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeFloatPtrSlice[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetFloatPtrSliceDefault get a float slice from object by key.
func GetFloatPtrSliceDefault[T float32 | float64](object map[string]any, key string) ([]*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableFloatPtrSlice[T](value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// decodeNullableNumberSlice tries to convert an unknown value to a number slice.
func decodeNullableNumberSlice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	value any,
	convertFn convertNullableFuncReflection[T],
) (*[]T, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a number slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]T, valueLen)

	for i := range valueLen {
		item := reflectValue.Index(i)

		result, err := convertFn(item)
		if err != nil {
			return nil, fmt.Errorf("failed to convert number element at %d: %w", i, err)
		}

		if result == nil {
			return nil, fmt.Errorf("number element at %d must not be null", i)
		}

		results[i] = *result
	}

	return &results, nil
}

// decodeNullableNumberPtrSlice tries to convert an unknown value to a number slice.
func decodeNullableNumberPtrSlice[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	value any,
	convertFn convertNullableFunc[T],
) (*[]*T, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a number slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*T, valueLen)

	for i := range valueLen {
		item := reflectValue.Index(i)

		result, err := convertFn(item.Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to convert number element at %d: %w", i, err)
		}

		results[i] = result
	}

	return &results, nil
}

// DecodeNullableStringSlice decodes a nullable string slice from an unknown value.
func DecodeNullableStringSlice(value any) (*[]string, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a string slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]string, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableStringReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode string element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("string element at %d must not be null", i)
		}

		results[i] = *elem
	}

	return &results, nil
}

// DecodeStringSlice decodes a string slice from an unknown value.
func DecodeStringSlice(value any) ([]string, error) {
	results, err := DecodeNullableStringSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("string slice must not be null")
	}

	return *results, nil
}

// DecodeNullableStringPtrSlice decodes a nullable string slice from an unknown value.
func DecodeNullableStringPtrSlice(value any) (*[]*string, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a string slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*string, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableStringReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode string element at %d: %w", i, err)
		}

		results[i] = elem
	}

	return &results, nil
}

// DecodeStringPtrSlice decodes a string pointer slice from an unknown value.
func DecodeStringPtrSlice(value any) ([]*string, error) {
	results, err := DecodeNullableStringPtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the string pointer slice must not be null")
	}

	return *results, nil
}

// GetStringSlice get a string slice value from object by key.
func GetStringSlice(object map[string]any, key string) ([]string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeStringSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetStringSliceDefault get a string slice value from object by key.
func GetStringSliceDefault(object map[string]any, key string) ([]string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableStringSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableStringSlice get a string pointer slice from object by key.
func GetNullableStringSlice(object map[string]any, key string) (*[]string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableStringSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetStringPtrSlice get a string pointer slice value from object by key.
func GetStringPtrSlice(object map[string]any, key string) ([]*string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeStringPtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetStringPtrSliceDefault get a string pointer slice value from object by key.
func GetStringPtrSliceDefault(object map[string]any, key string) ([]*string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableStringPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableStringPtrSlice get a nullable string pointer slice from object by key.
func GetNullableStringPtrSlice(object map[string]any, key string) (*[]*string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableStringPtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// DecodeNullableBooleanSlice decodes a nullable boolean slice from an unknown value.
func DecodeNullableBooleanSlice(value any) (*[]bool, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a boolean slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]bool, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableBooleanReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("boolean element at %d must not be null", i)
		}

		results[i] = *elem
	}

	return &results, nil
}

// DecodeBooleanSlice decodes a boolean slice from an unknown value.
func DecodeBooleanSlice(value any) ([]bool, error) {
	results, err := DecodeNullableBooleanSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("boolean slice must not be null")
	}

	return *results, nil
}

// DecodeNullableBooleanPtrSlice decodes a boolean pointer slice from an unknown value.
func DecodeNullableBooleanPtrSlice(value any) (*[]*bool, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a boolean slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*bool, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableBooleanReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
		}

		results[i] = elem
	}

	return &results, nil
}

// DecodeBooleanPtrSlice decodes a boolean pointer slice from an unknown value.
func DecodeBooleanPtrSlice(value any) ([]*bool, error) {
	results, err := DecodeNullableBooleanPtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the boolean pointer slice must not be null")
	}

	return *results, nil
}

// GetBooleanSlice get a boolean slice value from object by key.
func GetBooleanSlice(object map[string]any, key string) ([]bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeBooleanSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetBooleanSliceDefault get a boolean slice value from object by key.
func GetBooleanSliceDefault(object map[string]any, key string) ([]bool, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableBooleanSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableBooleanSlice get a nullable boolean slice from object by key.
func GetNullableBooleanSlice(object map[string]any, key string) (*[]bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableBooleanSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetBooleanPtrSlice get a boolean pointer slice value from object by key.
func GetBooleanPtrSlice(object map[string]any, key string) ([]*bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeBooleanPtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetBooleanPtrSliceDefault get a boolean pointer slice value from object by key.
func GetBooleanPtrSliceDefault(object map[string]any, key string) ([]*bool, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableBooleanPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableBooleanPtrSlice get a nullable boolean slice from object by key.
func GetNullableBooleanPtrSlice(object map[string]any, key string) (*[]*bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableBooleanPtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// DecodeUUIDSlice decodes UUID slice from array string.
func DecodeUUIDSlice(value any) ([]uuid.UUID, error) {
	results, err := DecodeNullableUUIDSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the uuid slice must not be null")
	}

	return *results, nil
}

// DecodeNullableUUIDSlice decodes a nullable UUID slice from array string.
func DecodeNullableUUIDSlice(value any) (*[]uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	strSlice, err := DecodeNullableStringPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid slice, got: %+v", value)
	}

	if strSlice == nil {
		return nil, err
	}

	results := make([]uuid.UUID, len(*strSlice))

	for i, str := range *strSlice {
		if str == nil {
			return nil, fmt.Errorf("uuid element at %d must not be null", i)
		}

		uid, err := uuid.Parse(*str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid element at %d: %w", i, err)
		}

		results[i] = uid
	}

	return &results, nil
}

// DecodeUUIDPtrSlice decodes UUID slice from array string.
func DecodeUUIDPtrSlice(value any) ([]*uuid.UUID, error) {
	results, err := DecodeNullableUUIDPtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the uuid pointer slice must not be null")
	}

	return *results, nil
}

// DecodeNullableUUIDPtrSlice decodes UUID pointer slice from array string.
func DecodeNullableUUIDPtrSlice(value any) (*[]*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	strSlice, err := DecodeNullableStringPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid slice, got: %+v", value)
	}

	if strSlice == nil {
		return nil, err
	}

	results := make([]*uuid.UUID, len(*strSlice))

	for i, str := range *strSlice {
		if str == nil {
			continue
		}

		uid, err := uuid.Parse(*str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse uuid element at %d: %w", i, err)
		}

		results[i] = &uid
	}

	return &results, nil
}

// GetUUIDSlice get an UUID slice from object by key.
func GetUUIDSlice(object map[string]any, key string) ([]uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field %s is required", key)
	}

	result, err := DecodeUUIDSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetUUIDSliceDefault get an UUID slice from object by key.
func GetUUIDSliceDefault(object map[string]any, key string) ([]uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableUUIDSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableUUIDSlice get an UUID pointer slice from object by key.
func GetNullableUUIDSlice(object map[string]any, key string) (*[]uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableUUIDSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetUUIDPtrSlice get an UUID slice from object by key.
func GetUUIDPtrSlice(object map[string]any, key string) ([]*uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field %s is required", key)
	}

	result, err := DecodeUUIDPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetUUIDPtrSliceDefault get an UUID slice from object by key.
func GetUUIDPtrSliceDefault(object map[string]any, key string) ([]*uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableUUIDPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableUUIDPtrSlice get an UUID pointer slice from object by key.
func GetNullableUUIDPtrSlice(object map[string]any, key string) (*[]*uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableUUIDPtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// DecodeNullableDateTimePtrSlice tries to convert an unknown value to a nullable time.Time pointer slice.
func DecodeNullableDateTimePtrSlice(value any) (*[]*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a time.Time slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*time.Time, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableDateTime(reflectValue.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to decode time.Time element at %d: %w", i, err)
		}

		if elem == nil {
			continue
		}

		results[i] = elem
	}

	return &results, nil
}

// DecodeDateTimePtrSlice tries to convert an unknown value to a time.Time pointer slice.
func DecodeDateTimePtrSlice(value any) ([]*time.Time, error) {
	results, err := DecodeNullableDateTimePtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("the time.Time slice must not be null")
	}

	return *results, nil
}

// DecodeNullableDateTimeSlice tries to convert an unknown value to a nullable time.Time slice.
func DecodeNullableDateTimeSlice(value any) (*[]time.Time, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a time.Time slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]time.Time, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableDateTime(reflectValue.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to decode time.Time element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("time.Time element at %d must not be null", i)
		}

		results[i] = *elem
	}

	return &results, nil
}

// DecodeDateTimeSlice tries to convert an unknown value to a time.Time slice.
func DecodeDateTimeSlice(value any) ([]time.Time, error) {
	results, err := DecodeNullableDateTimeSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("time.Time slice must not be null")
	}

	return *results, nil
}

// GetNullableDateTimeSlice get a nullable time.Time slice from object by key.
func GetNullableDateTimeSlice(object map[string]any, key string) (*[]time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableDateTimeSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetDateTimeSlice get a time.Time slice from object by key.
func GetDateTimeSlice(object map[string]any, key string) ([]time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeDateTimeSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetDateTimeSliceDefault get a time.Time slice from object by key.
func GetDateTimeSliceDefault(object map[string]any, key string) ([]time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableDateTimeSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableDateTimePtrSlice get a nullable time.Time pointer slice from object by key.
func GetNullableDateTimePtrSlice(object map[string]any, key string) (*[]*time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	result, err := DecodeNullableDateTimePtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetDateTimePtrSlice get a time.Time pointer slice from object by key.
func GetDateTimePtrSlice(object map[string]any, key string) ([]*time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	result, err := DecodeDateTimePtrSlice(value)
	if err != nil {
		return result, fmt.Errorf("%s: %w", key, err)
	}

	return result, nil
}

// GetDateTimePtrSliceDefault get a time.Time pointer slice from object by key.
func GetDateTimePtrSliceDefault(object map[string]any, key string) ([]*time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableDateTimePtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// DecodeNullableArbitraryJSONPtrSlice decodes a nullable arbitrary json pointer slice from object by key.
func DecodeNullableArbitraryJSONPtrSlice(value any) (*[]*any, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected an arbitrary json slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*any, valueLen)

	for i := range valueLen {
		elem, ok := UnwrapPointerFromReflectValue(reflectValue.Index(i))
		if !ok {
			continue
		}

		elemValue := elem.Interface()
		results[i] = &elemValue
	}

	return &results, nil
}

// DecodeArbitraryJSONPtrSlice decodes an arbitrary json pointer slice from object by key.
func DecodeArbitraryJSONPtrSlice(value any) ([]*any, error) {
	results, err := DecodeNullableArbitraryJSONPtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("arbitrary json pointer slice must not be null")
	}

	return *results, nil
}

// DecodeNullableArbitraryJSONSlice decodes a nullable arbitrary json slice from object by key.
func DecodeNullableArbitraryJSONSlice(value any) (*[]any, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected an arbitrary json slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]any, valueLen)

	for i := range valueLen {
		elem, ok := UnwrapPointerFromReflectValue(reflectValue.Index(i))
		if !ok {
			return nil, fmt.Errorf("arbitrary json element at %d must not be null", i)
		}

		results[i] = elem.Interface()
	}

	return &results, nil
}

// DecodeArbitraryJSONSlice decodes an arbitrary json slice from object by key.
func DecodeArbitraryJSONSlice(value any) ([]any, error) {
	results, err := DecodeNullableArbitraryJSONSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("arbitrary json slice must not be null")
	}

	return *results, nil
}

// GetNullableArbitraryJSONPtrSlice get a nullable arbitrary json pointer slice from object by key.
func GetNullableArbitraryJSONPtrSlice(object map[string]any, key string) (*[]*any, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	return DecodeNullableArbitraryJSONPtrSlice(value)
}

// GetArbitraryJSONPtrSlice get an arbitrary json pointer slice from object by key.
func GetArbitraryJSONPtrSlice(object map[string]any, key string) ([]*any, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	return DecodeArbitraryJSONPtrSlice(value)
}

// GetArbitraryJSONPtrSliceDefault get an arbitrary json pointer slice from object by key.
func GetArbitraryJSONPtrSliceDefault(object map[string]any, key string) ([]*any, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableArbitraryJSONPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableArbitraryJSONSlice get a nullable arbitrary json slice from object by key.
func GetNullableArbitraryJSONSlice(object map[string]any, key string) (*[]any, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	return DecodeNullableArbitraryJSONSlice(value)
}

// GetArbitraryJSONSlice get an arbitrary json slice from object by key.
func GetArbitraryJSONSlice(object map[string]any, key string) ([]any, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	return DecodeArbitraryJSONSlice(value)
}

// GetArbitraryJSONSliceDefault get an arbitrary json slice from object by key.
func GetArbitraryJSONSliceDefault(object map[string]any, key string) ([]any, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableArbitraryJSONSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// DecodeNullableRawJSONSlice decodes a nullable json.RawMessage slice from object by key.
func DecodeNullableRawJSONSlice(value any) (*[]json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a raw string or byte slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]json.RawMessage, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableRawJSON(reflectValue.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to decode json.RawMessage element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("json.RawMessage element at %d must not be null", i)
		}

		results[i] = *elem
	}

	return &results, nil
}

// DecodeRawJSONSlice decodes a json.RawMessage slice from object by key.
func DecodeRawJSONSlice(value any) ([]json.RawMessage, error) {
	results, err := DecodeNullableRawJSONSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("raw json slice must not be null")
	}

	return *results, nil
}

// DecodeNullableRawJSONPtrSlice decodes a nullable json.RawMessage pointer slice from object by key.
func DecodeNullableRawJSONPtrSlice(value any) (*[]*json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("expected a raw string or byte slice, got: %s", valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]*json.RawMessage, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableRawJSON(reflectValue.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to decode json.RawMessage element at %d: %w", i, err)
		}

		results[i] = elem
	}

	return &results, nil
}

// DecodeRawJSONPtrSlice decodes a json.RawMessage pointer slice from object by key.
func DecodeRawJSONPtrSlice(value any) ([]*json.RawMessage, error) {
	results, err := DecodeNullableRawJSONPtrSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("raw json pointer slice must not be null")
	}

	return *results, nil
}

// GetNullableRawJSONSlice get a nullable json.RawMessage slice from object by key.
func GetNullableRawJSONSlice(object map[string]any, key string) (*[]json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	return DecodeNullableRawJSONSlice(value)
}

// GetRawJSONSlice get a raw json.RawMessage slice from object by key.
func GetRawJSONSlice(object map[string]any, key string) ([]json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field %s is required", key)
	}

	return DecodeRawJSONSlice(value)
}

// GetRawJSONSliceDefault get a raw json.RawMessage slice from object by key.
func GetRawJSONSliceDefault(object map[string]any, key string) ([]json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableRawJSONSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}

// GetNullableRawJSONPtrSlice get a nullable json.RawMessage pointer slice from object by key.
func GetNullableRawJSONPtrSlice(object map[string]any, key string) (*[]*json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}

	return DecodeNullableRawJSONPtrSlice(value)
}

// GetRawJSONPtrSlice get a json.RawMessage pointer slice from object by key.
func GetRawJSONPtrSlice(object map[string]any, key string) ([]*json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, fmt.Errorf("field %s is required", key)
	}

	return DecodeRawJSONPtrSlice(value)
}

// GetRawJSONPtrSliceDefault get a json.RawMessage pointer slice from object by key.
func GetRawJSONPtrSliceDefault(object map[string]any, key string) ([]*json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNullableRawJSONPtrSlice(value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}

	if result == nil {
		return nil, nil
	}

	return *result, nil
}
