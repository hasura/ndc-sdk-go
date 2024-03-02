package utils

import (
	"errors"
	"fmt"
	"time"
)

// ToIntPtr tries to convert an interface to an integer pointer
func ToIntPtr[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value any) (*T, error) {
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

// ToInt tries to convert an interface to an integer value
func ToInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value any) (T, error) {
	result, err := ToIntPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Int value must not be null")
	}
	return *result, nil
}

// ToStringPtr tries to convert an interface to a string pointer
func ToStringPtr(value any) (*string, error) {
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

// ToString tries to convert an interface to a string value
func ToString(value any) (string, error) {
	result, err := ToStringPtr(value)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("the String value must not be null")
	}
	return *result, nil
}

// ToFloatPtr tries to convert an interface to a float pointer
func ToFloatPtr[T float32 | float64](value any) (*T, error) {
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

// ToFloat tries to convert an interface to a float value
func ToFloat[T float32 | float64](value any) (T, error) {
	result, err := ToFloatPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Float value must not be null")
	}
	return *result, nil
}

// ToComplexPtr tries to convert an interface to a complex pointer
func ToComplexPtr[T complex64 | complex128](value any) (*T, error) {
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

// ToComplex tries to convert an interface to a complex value
func ToComplex[T complex64 | complex128](value any) (T, error) {
	result, err := ToComplexPtr[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Complex value must not be null")
	}
	return *result, nil
}

// ToBooleanPtr tries to convert an interface to a bool pointer
func ToBooleanPtr(value any) (*bool, error) {
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

// ToBoolean tries to convert an interface to a bool value
func ToBoolean(value any) (bool, error) {
	result, err := ToBooleanPtr(value)
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, errors.New("the Boolean value must not be null")
	}
	return *result, nil
}

// ToDateTimePtr tries to convert an interface to a time.Time pointer
func ToDateTimePtr(value any) (*time.Time, error) {
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
		i64, err := ToIntPtr[int64](v)
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

// ToDateTime tries to convert an interface to a time.Time value
func ToDateTime(value any) (time.Time, error) {
	result, err := ToDateTimePtr(value)
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

// ToDurationPtr tries to convert an interface to a duration pointer
func ToDurationPtr(value any) (*time.Duration, error) {
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
		i64, err := ToIntPtr[int64](v)
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

// ToDuration tries to convert an interface to a duration value
func ToDuration(value any) (time.Duration, error) {
	result, err := ToDurationPtr(value)
	if err != nil {
		return time.Duration(0), err
	}
	if result == nil {
		return time.Duration(0), errors.New("the Duration value must not be null")
	}
	return *result, nil
}
