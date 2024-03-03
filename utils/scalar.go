package utils

import (
	"errors"
	"fmt"
	"time"
)

// Scalar abstracts a scalar interface to determine when evaluating
type Scalar interface {
	ScalarName() string
}

// MapEncoder abstracts a type with the ToMap method to encode type to map
type MapEncoder interface {
	ToMap() map[string]any
}

// ValueDecoder abstracts a type with the FromValue method to decode any value
type ValueDecoder interface {
	FromValue(value any) error
}

// DecodeValue tries to convert an unknown value to a value
func DecodeValue(decoder ValueDecoder, value any) error {
	if IsNil(value) {
		return nil
	}
	return decoder.FromValue(value)
}

// DecodeObjectValue get and decode a value from object by key
func DecodeObjectValue(decoder ValueDecoder, object map[string]any, key string) error {
	value, ok := GetAny(object, key)
	if !ok || IsNil(value) {
		return nil
	}
	err := decoder.FromValue(value)
	if err != nil {
		return fmt.Errorf("%s: %s", key, err)
	}
	return nil
}

// DecodeIntPtr tries to convert an unknown value to an integer pointer
func DecodeIntPtr[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value any) (*T, error) {
	var result T
	switch v := value.(type) {
	case int:
		result = T(v)
	case int8:
		result = T(v)
	case int16:
		result = T(v)
	case int32:
		result = T(v)
	case int64:
		result = T(v)
	case uint:
		result = T(v)
	case uint8:
		result = T(v)
	case uint16:
		result = T(v)
	case uint32:
		result = T(v)
	case *int:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *int8:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *int16:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *int32:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *int64:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *uint:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *uint8:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *uint16:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *uint32:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	default:
		return nil, fmt.Errorf("failed to convert Int, got %v", value)
	}

	return &result, nil
}

// DecodeInt tries to convert an unknown value to an integer value
func DecodeInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value any) (T, error) {
	result, err := DecodeIntPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Int value must not be null")
	}
	return *result, nil
}

// DecodeStringPtr tries to convert an unknown value to a string pointer
func DecodeStringPtr(value any) (*string, error) {
	var result string
	switch v := value.(type) {
	case string:
		result = v
	case *string:
		if v == nil {
			return nil, nil
		}
		result = *v
	default:
		return nil, fmt.Errorf("failed to convert String, got %v", value)
	}

	return &result, nil
}

// DecodeString tries to convert an unknown value to a string value
func DecodeString(value any) (string, error) {
	result, err := DecodeStringPtr(value)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("the String value must not be null")
	}
	return *result, nil
}

// DecodeFloatPtr tries to convert an unknown value to a float pointer
func DecodeFloatPtr[T float32 | float64](value any) (*T, error) {
	var result T
	switch v := value.(type) {
	case float32:
		result = T(v)
	case float64:
		result = T(v)
	case *float32:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *float64:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	default:
		return nil, fmt.Errorf("failed to convert Float, got %v", value)
	}

	return &result, nil
}

// DecodeFloat tries to convert an unknown value to a float value
func DecodeFloat[T float32 | float64](value any) (T, error) {
	result, err := DecodeFloatPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Float value must not be null")
	}
	return *result, nil
}

// DecodeComplexPtr tries to convert an unknown value to a complex pointer
func DecodeComplexPtr[T complex64 | complex128](value any) (*T, error) {
	var result T
	switch v := value.(type) {
	case complex64:
		result = T(v)
	case complex128:
		result = T(v)
	case *complex64:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	case *complex128:
		if v == nil {
			return nil, nil
		}
		result = T(*v)
	default:
		return nil, fmt.Errorf("failed to convert Complex, got %v", value)
	}

	return &result, nil
}

// DecodeComplex tries to convert an unknown value to a complex value
func DecodeComplex[T complex64 | complex128](value any) (T, error) {
	result, err := DecodeComplexPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Complex value must not be null")
	}
	return *result, nil
}

// DecodeBooleanPtr tries to convert an unknown value to a bool pointer
func DecodeBooleanPtr(value any) (*bool, error) {
	var result bool
	switch v := value.(type) {
	case bool:
		result = v
	case *bool:
		if v == nil {
			return nil, nil
		}
		result = *v
	default:
		return nil, fmt.Errorf("failed to convert Boolean, got %v", value)
	}

	return &result, nil
}

// DecodeBoolean tries to convert an unknown value to a bool value
func DecodeBoolean(value any) (bool, error) {
	result, err := DecodeBooleanPtr(value)
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, errors.New("the Boolean value must not be null")
	}
	return *result, nil
}

// DecodeDateTimePtr tries to convert an unknown value to a time.Time pointer
func DecodeDateTimePtr(value any) (*time.Time, error) {
	var result time.Time
	switch v := value.(type) {
	case time.Time:
		result = v
	case *time.Time:
		if v == nil {
			return nil, nil
		}
		result = *v
	case string:
		return parseDateTime(v)
	case *string:
		if v == nil {
			return nil, nil
		}
		return parseDateTime(*v)
	default:
		i64, err := DecodeIntPtr[int64](v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert DateTime, got %v", value)
		}
		if i64 == nil {
			return nil, nil
		}
		result = time.UnixMilli(*i64)
	}

	return &result, nil
}

// DecodeDateTime tries to convert an unknown value to a time.Time value
func DecodeDateTime(value any) (time.Time, error) {
	result, err := DecodeDateTimePtr(value)
	if err != nil {
		return time.Time{}, err
	}
	if result == nil {
		return time.Time{}, errors.New("the DateTime value must not be null")
	}
	return *result, nil
}

// parse date time with fallback ISO8601 formats
func parseDateTime(value string) (*time.Time, error) {
	for _, format := range []string{time.RFC3339, "2006-01-02T15:04:05Z0700", "2006-01-02T15:04:05-0700", time.RFC3339Nano} {
		result, err := time.Parse(value, format)
		if err != nil {
			continue
		}
		return &result, nil
	}

	return nil, fmt.Errorf("failed to parse time from string: %s", value)
}

// DecodeDurationPtr tries to convert an unknown value to a duration pointer
func DecodeDurationPtr(value any) (*time.Duration, error) {
	var result time.Duration
	switch v := value.(type) {
	case time.Duration:
		result = v
	case *time.Duration:
		if v == nil {
			return nil, nil
		}
		result = *v
	case string:
		dur, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		result = dur
	case *string:
		if v == nil {
			return nil, nil
		}
		dur, err := time.ParseDuration(*v)
		if err != nil {
			return nil, err
		}
		result = dur
	default:
		i64, err := DecodeIntPtr[int64](v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert DateTime, got %v", value)
		}
		if i64 == nil {
			return nil, nil
		}
		result = time.Duration(*i64)
	}

	return &result, nil
}

// DecodeDuration tries to convert an unknown value to a duration value
func DecodeDuration(value any) (time.Duration, error) {
	result, err := DecodeDurationPtr(value)
	if err != nil {
		return time.Duration(0), err
	}
	if result == nil {
		return time.Duration(0), errors.New("the Duration value must not be null")
	}
	return *result, nil
}

// GetAny get an unknown value from object by key
func GetAny(object map[string]any, key string) (any, bool) {
	if object == nil {
		return nil, false
	}
	value, ok := object[key]
	return value, ok
}

// GetIntPtr get an integer pointer from object by key
func GetIntPtr[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeIntPtr[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetInt get an integer value from object by key
func GetInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeInt[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetFloatPtr get a float pointer from object by key
func GetFloatPtr[T float32 | float64](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeFloatPtr[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetFloat get a float value from object by key
func GetFloat[T float32 | float64](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeFloat[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetComplexPtr get a complex pointer from object by key
func GetComplexPtr[T complex64 | complex128](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeComplexPtr[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetComplex get a complex value from object by key
func GetComplex[T complex64 | complex128](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeComplex[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetStringPtr get a string pointer from object by key
func GetStringPtr(object map[string]any, key string) (*string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeStringPtr(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetString get a bool value from object by key
func GetString(object map[string]any, key string) (string, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return "", fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeString(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetBoolPtr get a bool pointer from object by key
func GetBoolPtr(object map[string]any, key string) (*bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeBooleanPtr(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetBool get a bool value from object by key
func GetBool(object map[string]any, key string) (bool, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return false, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeBoolean(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDateTimePtr get a time.Time pointer from object by key
func GetDateTimePtr(object map[string]any, key string) (*time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeDateTimePtr(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDateTime get a time.Time value from object by key
func GetDateTime(object map[string]any, key string) (time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return time.Time{}, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeDateTime(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDurationPtr get a time.Duration pointer from object by key
func GetDurationPtr(object map[string]any, key string) (*time.Duration, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeDurationPtr(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDuration get a time.Duration value from object by key
func GetDuration(object map[string]any, key string) (time.Duration, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` does not exist", key)
	}
	result, err := DecodeDuration(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}
