package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/google/uuid"
)

type convertFunc[T any] func(value any) (*T, error)
type convertFuncReflection[T any] func(value reflect.Value) (*T, error)

// ValueDecoder abstracts a type with the FromValue method to decode any value
type ValueDecoder interface {
	FromValue(value any) error
}

// ObjectDecoder abstracts a type with the FromValue method to decode an object value
type ObjectDecoder interface {
	FromValue(value map[string]any) error
}

// IsNil a safe function to check null value
func IsNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Pointer && v.IsNil()
}

// Decoder is a wrapper of mapstructure decoder
type Decoder struct {
	decodeHook mapstructure.DecodeHookFunc
}

// NewDecoder creates a Decoder instance
func NewDecoder(decodeHooks ...mapstructure.DecodeHookFunc) *Decoder {
	return &Decoder{
		decodeHook: mapstructure.ComposeDecodeHookFunc(append(defaultDecodeFuncs, decodeHooks...)...),
	}
}

// DecodeObjectValue get and decode a value from object by key
func (d Decoder) DecodeObjectValue(target any, object map[string]any, key string) error {
	value, ok := GetAny(object, key)
	if !ok {
		return fmt.Errorf("%s: field is required", key)
	}
	err := d.DecodeValue(target, value)
	if err != nil {
		return fmt.Errorf("%s: %s", key, err)
	}
	return nil
}

// DecodeNullableObjectValue get and decode a nullable value from object by key
func (d Decoder) DecodeNullableObjectValue(target any, object map[string]any, key string) error {
	value, ok := GetAny(object, key)
	if !ok {
		return nil
	}
	err := d.DecodeNullableValue(target, value)
	if err != nil {
		return fmt.Errorf("%s: %s", key, err)
	}
	return nil
}

// DecodeValue tries to convert and set an unknown value into the target, the value must not be null
// fallback to mapstructure decoder
func (d Decoder) DecodeValue(target any, value any) error {
	if IsNil(value) {
		return errors.New("the value must not be null")
	}
	return d.decodeValue(target, value)
}

// DecodeNullableValue tries to convert and set an unknown value into the target,
// fallback to mapstructure decoder
func (d Decoder) DecodeNullableValue(target any, value any) error {
	if IsNil(value) {
		return nil
	}
	return d.decodeValue(target, value)
}

func (d Decoder) decodeValue(target any, value any) error {
	if IsNil(target) {
		return errors.New("the decoded target must be not null")
	}
	if IsNil(value) {
		return nil
	}

	switch t := target.(type) {
	case ObjectDecoder:
		switch v := value.(type) {
		case map[string]any:
			return t.FromValue(v)
		case MapEncoder:
			object, err := v.ToMap()
			if err != nil {
				return err
			}
			return t.FromValue(object)
		default:
			return errors.New("the value must be an object-liked")
		}
	case ValueDecoder:
		return t.FromValue(value)
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, time.Time, time.Duration, time.Ticker:
		return errors.New("the decoded target must be a pointer")
	case *complex64, *complex128:
		return errors.New("unsupported complex types")
	case *bool:
		v, err := DecodeNullableBoolean(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *string:
		v, err := DecodeNullableString(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *int:
		v, err := DecodeNullableInt[int](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *int8:
		v, err := DecodeNullableInt[int8](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *int16:
		v, err := DecodeNullableInt[int16](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *int32:
		v, err := DecodeNullableInt[int32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *int64:
		v, err := DecodeNullableInt[int64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uint:
		v, err := DecodeNullableUint[uint](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uint8:
		v, err := DecodeNullableUint[uint8](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uint16:
		v, err := DecodeNullableUint[uint16](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uint32:
		v, err := DecodeNullableUint[uint32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uint64:
		v, err := DecodeNullableUint[uint64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *float32:
		v, err := DecodeNullableFloat[float32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *float64:
		v, err := DecodeNullableFloat[float64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]bool:
		v, err := DecodeBooleanSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]string:
		v, err := DecodeStringSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]int:
		v, err := DecodeIntSlice[int](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]int8:
		v, err := DecodeIntSlice[int8](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]int16:
		v, err := DecodeIntSlice[int16](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]int32:
		v, err := DecodeIntSlice[int32](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]int64:
		v, err := DecodeIntSlice[int64](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]uint:
		v, err := DecodeUintSlice[uint](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]uint8:
		v, err := DecodeUintSlice[uint8](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]uint16:
		v, err := DecodeUintSlice[uint16](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]uint32:
		v, err := DecodeUintSlice[uint32](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]uint64:
		v, err := DecodeUintSlice[uint64](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]float32:
		v, err := DecodeFloatSlice[float32](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]float64:
		v, err := DecodeFloatSlice[float64](value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]json.RawMessage:
		v, err := DecodeRawJSONSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = v
	case *[]*bool:
		v, err := DecodeNullableBooleanPtrSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*string:
		v, err := DecodeNullableStringPtrSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*int:
		v, err := DecodeNullableIntPtrSlice[int](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*int8:
		v, err := DecodeNullableIntPtrSlice[int8](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*int16:
		v, err := DecodeNullableIntPtrSlice[int16](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*int32:
		v, err := DecodeNullableIntPtrSlice[int32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*int64:
		v, err := DecodeNullableIntPtrSlice[int64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uint:
		v, err := DecodeNullableUintPtrSlice[uint](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uint8:
		v, err := DecodeNullableUintPtrSlice[uint8](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uint16:
		v, err := DecodeNullableUintPtrSlice[uint16](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uint32:
		v, err := DecodeNullableUintPtrSlice[uint32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uint64:
		v, err := DecodeNullableUintPtrSlice[uint64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*float32:
		v, err := DecodeNullableFloatPtrSlice[float32](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*float64:
		v, err := DecodeNullableFloatPtrSlice[float64](value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*json.RawMessage:
		v, err := DecodeNullableRawJSONPtrSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *time.Time:
		v, err := DecodeNullableDateTime(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]time.Time:
		v, err := DecodeNullableDateTimeSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*time.Time:
		v, err := DecodeNullableDateTimePtrSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *time.Duration:
		v, err := DecodeNullableDuration(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *uuid.UUID:
		v, err := DecodeNullableUUID(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]uuid.UUID:
		v, err := DecodeNullableUUIDSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	case *[]*uuid.UUID:
		v, err := DecodeNullableUUIDPtrSlice(value)
		if err != nil || v == nil {
			return err
		}
		*t = *v
	default:
		return decodeAnyValue(target, value, d.decodeHook)
	}

	return nil
}

// DecodeNullableInt tries to convert an unknown value to a nullable integer
func DecodeNullableInt[T int | int8 | int16 | int32 | int64](value any) (*T, error) {
	return decodeNullableInt(value, convertNullableInt[T])
}

// DecodeInt tries to convert an unknown value to a not-null integer value
func DecodeInt[T int | int8 | int16 | int32 | int64](value any) (T, error) {
	result, err := DecodeNullableInt[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Int value must not be null")
	}
	return *result, nil
}

// DecodeNullableUint tries to convert an unknown value to a nullable unsigned integer pointer
func DecodeNullableUint[T uint | uint8 | uint16 | uint32 | uint64](value any) (*T, error) {
	return decodeNullableInt(value, convertNullableUint[T])
}

// DecodeUint tries to convert an unknown value to an unsigned integer value
func DecodeUint[T uint | uint8 | uint16 | uint32 | uint64](value any) (T, error) {
	result, err := DecodeNullableUint[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Uint value must not be null")
	}
	return *result, nil
}

func decodeNullableInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value any, convertFn convertFunc[T]) (*T, error) {
	if IsNil(value) {
		return nil, nil
	}
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
	case uint64:
		result = T(v)
	case float32:
		result = T(v)
	case float64:
		result = T(v)
	case string:
		newVal, err := convertFn(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert integer, got: %s", v)
		}
		return newVal, err
	case *int:
		result = T(*v)
	case *int8:
		result = T(*v)
	case *int16:
		result = T(*v)
	case *int32:
		result = T(*v)
	case *int64:
		result = T(*v)
	case *uint:
		result = T(*v)
	case *uint8:
		result = T(*v)
	case *uint16:
		result = T(*v)
	case *uint32:
		result = T(*v)
	case *uint64:
		result = T(*v)
	case *float32:
		result = T(*v)
	case *float64:
		result = T(*v)
	case *string:
		newVal, err := convertFn(*v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert integer, got: %s", *v)
		}
		return newVal, err
	case bool, complex64, complex128, time.Time, time.Duration, time.Ticker, *bool, *complex64, *complex128, *time.Time, *time.Duration, *time.Ticker, []bool, []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128, []time.Time, []time.Duration, []time.Ticker:
		return nil, fmt.Errorf("failed to convert integer, got: %+v", value)
	default:
		return decodeNullableIntRefection(convertFn)(reflect.ValueOf(value))
	}

	return &result, nil
}

func decodeNullableIntRefection[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](convertFn convertFunc[T]) convertFuncReflection[T] {
	return func(value reflect.Value) (*T, error) {
		inferredValue, ok := UnwrapPointerFromReflectValue(value)
		if !ok {
			return nil, nil
		}
		var result T
		switch inferredValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = T(inferredValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result = T(inferredValue.Uint())
		case reflect.Float32, reflect.Float64:
			result = T(inferredValue.Float())
		case reflect.Interface, reflect.String:
			newVal, parseErr := convertFn(inferredValue.Interface())
			if parseErr != nil {
				return nil, fmt.Errorf("failed to convert integer, got: %+v", inferredValue.Interface())
			}
			result = T(*newVal)
		default:
			return nil, fmt.Errorf("failed to convert integer, got: %+v", inferredValue.Interface())
		}

		return &result, nil
	}
}

func convertNullableInt[T int | int8 | int16 | int32 | int64](v any) (*T, error) {
	rawResult, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	if err != nil {
		return nil, err
	}
	result := T(rawResult)
	return &result, nil
}

func convertNullableUint[T uint | uint8 | uint16 | uint32 | uint64](v any) (*T, error) {
	rawResult, err := strconv.ParseUint(fmt.Sprint(v), 10, 64)
	if err != nil {
		return nil, err
	}
	result := T(rawResult)
	return &result, nil
}

// DecodeNullableString tries to convert an unknown value to a string pointer
func DecodeNullableString(value any) (*string, error) {
	if IsNil(value) {
		return nil, nil
	}
	var result string
	switch v := value.(type) {
	case string:
		result = v
	case *string:
		result = *v
	default:
		return decodeNullableStringReflection(reflect.ValueOf(value))
	}

	return &result, nil
}

func decodeNullableStringReflection(value reflect.Value) (*string, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	switch inferredValue.Kind() {
	case reflect.String:
		result := inferredValue.String()
		return &result, nil
	case reflect.Interface:
		result := fmt.Sprint(inferredValue.Interface())
		return &result, nil
	default:
		return nil, fmt.Errorf("failed to convert String, got: %v", value)
	}
}

// DecodeString tries to convert an unknown value to a string value
func DecodeString(value any) (string, error) {
	result, err := DecodeNullableString(value)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", errors.New("the String value must not be null")
	}
	return *result, nil
}

// DecodeNullableFloat tries to convert an unknown value to a float pointer
func DecodeNullableFloat[T float32 | float64](value any) (*T, error) {
	if IsNil(value) {
		return nil, nil
	}
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
	case uint64:
		result = T(v)
	case *int:
		result = T(*v)
	case *int8:
		result = T(*v)
	case *int16:
		result = T(*v)
	case *int32:
		result = T(*v)
	case *int64:
		result = T(*v)
	case *uint:
		result = T(*v)
	case *uint8:
		result = T(*v)
	case *uint16:
		result = T(*v)
	case *uint32:
		result = T(*v)
	case *uint64:
		result = T(*v)
	case float32:
		result = T(v)
	case float64:
		result = T(v)
	case *float32:
		result = T(*v)
	case *float64:
		result = T(*v)
	case bool, string, complex64, complex128, time.Time, time.Duration, time.Ticker, *bool, *string, *complex64, *complex128, *time.Time, *time.Duration, *time.Ticker, []bool, []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128, []time.Time, []time.Duration, []time.Ticker:
		return nil, fmt.Errorf("failed to convert Float, got: %+v", value)
	default:
		return decodeNullableFloatReflection[T](reflect.ValueOf(value))
	}

	return &result, nil
}

func decodeNullableFloatReflection[T float32 | float64](value reflect.Value) (*T, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}
	var result T
	switch inferredValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = T(inferredValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = T(inferredValue.Uint())
	case reflect.Float32, reflect.Float64:
		result = T(inferredValue.Float())
	case reflect.Interface:
		v := fmt.Sprint(inferredValue.Interface())
		newVal, parseErr := strconv.ParseFloat(v, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to convert Float, got: %s", v)
		}
		result = T(newVal)
	default:
		return nil, fmt.Errorf("failed to convert Float, got: %+v", inferredValue.Interface())
	}
	return &result, nil
}

// DecodeFloat tries to convert an unknown value to a float value
func DecodeFloat[T float32 | float64](value any) (T, error) {
	result, err := DecodeNullableFloat[T](value)
	if err != nil {
		return T(0), err
	}
	if result == nil {
		return T(0), errors.New("the Float value must not be null")
	}
	return *result, nil
}

// DecodeNullableBoolean tries to convert an unknown value to a bool pointer
func DecodeNullableBoolean(value any) (*bool, error) {
	if IsNil(value) {
		return nil, nil
	}

	var result bool
	switch v := value.(type) {
	case bool:
		result = v
	case *bool:
		result = *v
	default:
		return decodeNullableBooleanReflection(reflect.ValueOf(value))
	}

	return &result, nil
}

// DecodeBoolean tries to convert an unknown value to a bool value
func DecodeBoolean(value any) (bool, error) {
	result, err := DecodeNullableBoolean(value)
	if err != nil {
		return false, err
	}
	if result == nil {
		return false, errors.New("the Boolean value must not be null")
	}
	return *result, nil
}

func decodeNullableBooleanReflection(value reflect.Value) (*bool, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	kind := inferredValue.Kind()
	switch kind {
	case reflect.Bool:
		result := inferredValue.Bool()
		return &result, nil
	case reflect.Interface:
		v := fmt.Sprint(inferredValue.Interface())
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Boolean, got: %s", v)
		}
		return &b, nil
	default:
		return nil, fmt.Errorf("failed to convert Boolean, got: %v", kind)
	}
}

// DecodeNullableDateTime tries to convert an unknown value to a time.Time pointer
func DecodeNullableDateTime(value any) (*time.Time, error) {
	var result time.Time
	switch v := value.(type) {
	case time.Time:
		result = v
	case *time.Time:
		result = *v
	case string:
		return parseDateTime(v)
	case *string:
		if IsNil(v) {
			return nil, nil
		}
		return parseDateTime(*v)
	default:
		i64, err := DecodeNullableInt[int64](v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert DateTime, got: %v", value)
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
	result, err := DecodeNullableDateTime(value)
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
		result, err := time.Parse(format, value)
		if err != nil {
			continue
		}
		return &result, nil
	}

	return nil, fmt.Errorf("failed to parse time from string: %s", value)
}

// DecodeNullableDuration tries to convert an unknown value to a duration pointer
func DecodeNullableDuration(value any) (*time.Duration, error) {
	var result time.Duration
	switch v := value.(type) {
	case time.Duration:
		result = v
	case *time.Duration:
		result = *v
	case string:
		dur, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		result = dur
	case *string:
		if IsNil(v) {
			return nil, nil
		}
		dur, err := time.ParseDuration(*v)
		if err != nil {
			return nil, err
		}
		result = dur
	default:
		i64, err := DecodeNullableInt[int64](v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Duration, got: %v", value)
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
	result, err := DecodeNullableDuration(value)
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

// GetNullableArbitraryJSON get an arbitrary json pointer from object by key
func GetNullableArbitraryJSON(object map[string]any, key string) (*any, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	value, ok = UnwrapPointerFromAny(value)
	if !ok {
		return nil, nil
	}
	return &value, nil
}

// GetArbitraryJSON get an arbitrary json value from object by key
func GetArbitraryJSON(object map[string]any, key string) (any, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, fmt.Errorf("field `%s` is required", key)
	}

	value, ok = UnwrapPointerFromAny(value)
	if !ok {
		return nil, fmt.Errorf("field `%s` must not be null", key)
	}
	return value, nil
}

// GetNullableInt get an integer pointer from object by key
func GetNullableInt[T int | int8 | int16 | int32 | int64](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableInt[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetInt get an integer value from object by key
func GetInt[T int | int8 | int16 | int32 | int64](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeInt[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableUint get an unsigned integer pointer from object by key
func GetNullableUint[T uint | uint8 | uint16 | uint32 | uint64](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableUint[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetUint get an unsigned integer value from object by key
func GetUint[T uint | uint8 | uint16 | uint32 | uint64](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeUint[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableFloat get a float pointer from object by key
func GetNullableFloat[T float32 | float64](object map[string]any, key string) (*T, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableFloat[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetFloat get a float value from object by key
func GetFloat[T float32 | float64](object map[string]any, key string) (T, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeFloat[T](value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableString get a string pointer from object by key
func GetNullableString(object map[string]any, key string) (*string, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableString(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetString get a string value from object by key
func GetString(object map[string]any, key string) (string, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return "", fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeString(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableBoolean get a bool pointer from object by key
func GetNullableBoolean(object map[string]any, key string) (*bool, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableBoolean(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetBoolean get a bool value from object by key
func GetBoolean(object map[string]any, key string) (bool, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return false, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeBoolean(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableDateTime get a time.Time pointer from object by key
func GetNullableDateTime(object map[string]any, key string) (*time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableDateTime(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDateTime get a time.Time value from object by key
func GetDateTime(object map[string]any, key string) (time.Time, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return time.Time{}, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeDateTime(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableDuration get a time.Duration pointer from object by key
func GetNullableDuration(object map[string]any, key string) (*time.Duration, error) {
	value, ok := GetAny(object, key)
	if !ok || value == nil {
		return nil, nil
	}
	result, err := DecodeNullableDuration(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetDuration get a time.Duration value from object by key
func GetDuration(object map[string]any, key string) (time.Duration, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return 0, fmt.Errorf("field `%s` is required", key)
	}
	result, err := DecodeDuration(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// DecodeUUID decodes UUID from string
func DecodeUUID(value any) (uuid.UUID, error) {
	result, err := DecodeNullableUUID(value)
	if err != nil {
		return uuid.UUID{}, err
	}
	if result == nil {
		return uuid.UUID{}, errors.New("the uuid value must not be null")
	}
	return *result, nil
}

// DecodeNullableUUID decodes UUID pointer from string or bytes
func DecodeNullableUUID(value any) (*uuid.UUID, error) {
	if IsNil(value) {
		return nil, nil
	}
	switch v := value.(type) {
	case string:
		result, err := uuid.Parse(v)
		if err != nil {
			return nil, err
		}
		return &result, nil
	case *string:
		if v == nil {
			return nil, nil
		}
		result, err := uuid.Parse(*v)
		if err != nil {
			return nil, err
		}
		return &result, nil
	case [16]byte:
		result := uuid.UUID(v)
		return &result, nil
	case *[16]byte:
		if v == nil {
			return nil, nil
		}
		result := uuid.UUID(*v)
		return &result, nil
	case uuid.UUID:
		return &v, nil
	case *uuid.UUID:
		return v, nil
	default:
		return nil, fmt.Errorf("failed to parse uuid, got: %+v", value)
	}
}

// GetUUID get an UUID value from object by key
func GetUUID(object map[string]any, key string) (uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("field %s is required", key)
	}
	result, err := DecodeUUID(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetNullableUUID get an UUID pointer from object by key
func GetNullableUUID(object map[string]any, key string) (*uuid.UUID, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}
	result, err := DecodeNullableUUID(value)
	if err != nil {
		return result, fmt.Errorf("%s: %s", key, err)
	}
	return result, nil
}

// GetRawJSON get a raw json.RawMessage value from object by key
func GetRawJSON(object map[string]any, key string) (json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok || IsNil(value) {
		return nil, fmt.Errorf("field %s is required", key)
	}
	result, err := DecodeNullableRawJSON(value)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// DecodeNullableRawJSON decodes a raw json.RawMessage pointer from object by key
func DecodeNullableRawJSON(value any) (*json.RawMessage, error) {
	if IsNil(value) {
		return nil, nil
	}
	switch v := value.(type) {
	case []byte:
		vp := json.RawMessage(v)
		return &vp, nil
	case *[]byte:
		result := json.RawMessage(*v)
		return &result, nil
	default:
		result, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		rawResult := json.RawMessage(result)
		return &rawResult, nil
	}
}

// GetNullableRawJSON gets a raw json.RawMessage pointer from object by key
func GetNullableRawJSON(object map[string]any, key string) (*json.RawMessage, error) {
	value, ok := GetAny(object, key)
	if !ok {
		return nil, nil
	}
	return DecodeNullableRawJSON(value)
}

func decodeAnyValue(target any, value any, decodeHook mapstructure.DecodeHookFunc) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     target,
		TagName:    "json",
		DecodeHook: decodeHook,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(value)
}

var defaultDecodeFuncs = []mapstructure.DecodeHookFunc{
	decodeValueHookFunc(),
	decodeTimeHookFunc(),
	decodeUUIDHookFunc(),
}

func decodeValueHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Value, to reflect.Value) (any, error) {
		if from.Kind() == reflect.Pointer && from.IsNil() {
			return nil, nil
		}
		if to.Kind() == reflect.Pointer {
			if to.IsNil() {
				to = reflect.New(to.Type().Elem())
			}
		}
		toValue := to.Interface()
		fromValue := from.Interface()
		decoder, ok := toValue.(ValueDecoder)
		if ok {
			err := decoder.FromValue(fromValue)
			return decoder, err
		}

		if to.CanAddr() {
			decoder, ok = to.Addr().Interface().(ValueDecoder)
			if ok {
				if err := decoder.FromValue(fromValue); err != nil {
					return nil, err
				}
				return reflect.ValueOf(decoder).Elem().Interface(), nil
			}
		}

		objDecoder, ok := toValue.(ObjectDecoder)
		isObjectPtr := false
		if !ok && to.CanAddr() {
			objDecoder, ok = to.Addr().Interface().(ObjectDecoder)
			isObjectPtr = true
		}
		if ok {
			var err error
			switch v := fromValue.(type) {
			case map[string]any:
				err = objDecoder.FromValue(v)
			case MapEncoder:
				object, toMapErr := v.ToMap()
				if toMapErr != nil {
					return nil, toMapErr
				}
				mapValue := object
				err = objDecoder.FromValue(mapValue)
			}
			if err != nil {
				return nil, err
			}
			if isObjectPtr {
				return reflect.ValueOf(objDecoder).Elem().Interface(), nil
			}
			return objDecoder, nil
		}

		return fromValue, nil
	}
}

func decodeTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Value, to reflect.Value) (any, error) {
		toType := to.Type()
		if toType.PkgPath() != "time" || toType.Name() != "Time" {
			return from.Interface(), nil
		}

		kind := from.Type().Kind()
		switch kind {
		case reflect.String:
			return parseDateTime(from.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return time.UnixMilli(from.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return time.UnixMilli(int64(from.Uint())), nil
		case reflect.Float32, reflect.Float64:
			return time.UnixMilli(int64(from.Float())), nil
		default:
			return nil, fmt.Errorf("failed to decode time.Time, got: %s", kind.String())
		}
	}
}

func decodeUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if to.PkgPath() != "github.com/google/uuid" || to.Name() != "UUID" {
			return data, nil
		}
		result, err := DecodeNullableUUID(data)
		if err != nil || result == nil {
			return uuid.UUID{}, err
		}

		return *result, nil
	}
}
