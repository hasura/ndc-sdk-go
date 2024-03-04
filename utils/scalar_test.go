package utils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
		t.FailNow()
	}
}

func assertError(t *testing.T, err error, message string) {
	if err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	} else if !strings.Contains(err.Error(), message) {
		t.Errorf("expected error with content: %s, got: %s", err.Error(), message)
		t.FailNow()
	}
}

func assertEqual(t *testing.T, expected any, reality any) {
	if expected != reality {
		t.Errorf("not equal, expected: %+v got: %+v", expected, reality)
		t.FailNow()
	}
}

// func assertDeepEqual(t *testing.T, expected any, reality any) {
// 	if !internal.DeepEqual(expected, reality) {
// 		t.Errorf("not equal, expected: %+v got: %+v", expected, reality)
// 		t.FailNow()
// 	}
// }

func TestDecodeBool(t *testing.T) {
	value, err := DecodeBoolean(true)
	assertNoError(t, err)
	assertEqual(t, true, value)

	ptr, err := DecodeBooleanPtr(true)
	assertNoError(t, err)
	assertEqual(t, true, *ptr)

	ptr2, err := DecodeBooleanPtr(ptr)
	assertNoError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeBoolean(nil)
	assertError(t, err, "the Boolean value must not be null")

	_, err = DecodeBoolean("failure")
	assertError(t, err, "failed to convert Boolean, got: failure")
}

func TestDecodeString(t *testing.T) {
	value, err := DecodeString("success")
	assertNoError(t, err)
	assertEqual(t, "success", value)

	ptr, err := DecodeStringPtr("pointer")
	assertNoError(t, err)
	assertEqual(t, "pointer", *ptr)

	ptr2, err := DecodeStringPtr(ptr)
	assertNoError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeString(nil)
	assertError(t, err, "the String value must not be null")

	_, err = DecodeString(0)
	assertError(t, err, "failed to convert String, got: 0")
}

func TestDecodeDateTime(t *testing.T) {
	now := time.Now()
	t.Run("decode_value", func(t *testing.T) {
		value, err := DecodeDateTime(now)
		assertNoError(t, err)
		assertEqual(t, now, value)
	})

	t.Run("unix_int", func(t *testing.T) {
		iNow := now.UnixMilli()
		value, err := DecodeDateTime(iNow)
		assertNoError(t, err)
		assertEqual(t, now.UnixMilli(), value.UnixMilli())

		value, err = DecodeDateTime(&iNow)
		assertNoError(t, err)
		assertEqual(t, now.UnixMilli(), value.UnixMilli())

		var nilI64 *int64 = nil
		ptr, err := DecodeDateTimePtr(nilI64)
		assertNoError(t, err)
		assertEqual(t, true, IsNil(ptr))
	})

	t.Run("from_string", func(t *testing.T) {
		nowStr := now.Format(time.RFC3339)
		value, err := DecodeDateTime(nowStr)
		assertNoError(t, err)
		assertEqual(t, now.Unix(), value.Unix())

		value, err = DecodeDateTime(&nowStr)
		assertNoError(t, err)
		assertEqual(t, now.Unix(), value.Unix())

		invalidStr := "test"
		_, err = DecodeDateTime(invalidStr)
		assertError(t, err, "failed to parse time from string: test")
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeDateTimePtr(&now)
		assertNoError(t, err)
		assertEqual(t, now, *ptr)

		ptr2, err := DecodeDateTimePtr(ptr)
		assertNoError(t, err)
		assertEqual(t, *ptr, *ptr2)
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeDateTime(nil)
		assertError(t, err, "the DateTime value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeDateTime(false)
		assertError(t, err, "failed to convert DateTime, got: false")
	})
}

func TestDecodeDuration(t *testing.T) {
	duration := 10 * time.Second
	t.Run("decode_value", func(t *testing.T) {
		value, err := DecodeDuration(duration)
		assertNoError(t, err)
		assertEqual(t, duration, value)
	})

	t.Run("unix_int", func(t *testing.T) {
		iDuration := int64(duration)
		value, err := DecodeDuration(iDuration)
		assertNoError(t, err)
		assertEqual(t, duration, value)

		value, err = DecodeDuration(&iDuration)
		assertNoError(t, err)
		assertEqual(t, duration, value)
	})

	t.Run("from_string", func(t *testing.T) {
		durationStr := "10s"
		value, err := DecodeDuration(durationStr)
		assertNoError(t, err)
		assertEqual(t, duration, value)

		value, err = DecodeDuration(&durationStr)
		assertNoError(t, err)
		assertEqual(t, duration, value)

		invalidStr := "test"
		_, err = DecodeDuration(invalidStr)
		assertError(t, err, "time: invalid duration \"test\"")

		_, err = DecodeDuration(&invalidStr)
		assertError(t, err, "time: invalid duration \"test\"")

		var nilStr *string
		nilValue, err := DecodeDurationPtr(nilStr)
		assertNoError(t, err)
		assertEqual(t, true, IsNil(nilValue))
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeDurationPtr(&duration)
		assertNoError(t, err)
		assertEqual(t, duration, *ptr)

		ptr2, err := DecodeDurationPtr(ptr)
		assertNoError(t, err)
		assertEqual(t, *ptr, *ptr2)
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeDuration(nil)
		assertError(t, err, "the Duration value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeDuration(false)
		assertError(t, err, "failed to convert Duration, got: false")
	})
}

func TestDecodeInt(t *testing.T) {

	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeInt[int64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeIntPtr[int64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeIntPtr[int64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vInt := int(1)
		ptr, err := DecodeIntPtr[int64](&vInt)
		assertNoError(t, err)
		assertEqual(t, int64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeIntPtr[int64](&ptrInt)
		assertNoError(t, err)
		assertEqual(t, int64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeIntPtr[int64](&vInt8)
		assertNoError(t, err)
		assertEqual(t, int64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeIntPtr[int64](&vInt16)
		assertNoError(t, err)
		assertEqual(t, int64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeIntPtr[int64](&vInt32)
		assertNoError(t, err)
		assertEqual(t, int64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeIntPtr[int64](&vInt64)
		assertNoError(t, err)
		assertEqual(t, int64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeIntPtr[int64](&vUint)
		assertNoError(t, err)
		assertEqual(t, int64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeIntPtr[int64](&vUint8)
		assertNoError(t, err)
		assertEqual(t, int64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeIntPtr[int64](&vUint16)
		assertNoError(t, err)
		assertEqual(t, int64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeIntPtr[int64](&vUint32)
		assertNoError(t, err)
		assertEqual(t, int64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeIntPtr[int64](&vUint64)
		assertNoError(t, err)
		assertEqual(t, int64(vUint64), *ptr)

		var vAny any = float64(1.1)
		_, err = DecodeIntPtr[int64](&vAny)
		assertError(t, err, "failed to convert integer, got: *interface {} (1.1)")

		var vFn any = func() {}
		_, err = DecodeIntPtr[int64](&vFn)
		assertError(t, err, "failed to convert integer")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeInt[rune](nil)
		assertError(t, err, "the Int value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeInt[rune]("failure")
		assertError(t, err, "failed to convert integer, got: failure")
	})
}

func TestDecodeUint(t *testing.T) {

	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeUint[uint64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeUintPtr[uint64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeUintPtr[uint64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vUint := uint(1)
		ptr, err := DecodeUintPtr[uint64](&vUint)
		assertNoError(t, err)
		assertEqual(t, uint64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeUintPtr[uint64](&ptrUint)
		assertNoError(t, err)
		assertEqual(t, uint64(vUint), *ptr)

		var vAny any = float64(1.1)
		_, err = DecodeUintPtr[uint64](&vAny)
		assertError(t, err, "failed to convert integer, got: *interface {} (1.1)")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeUint[uint](nil)
		assertError(t, err, "the Uint value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeUint[uint]("failure")
		assertError(t, err, "failed to convert integer, got: failure")
	})
}

func TestDecodeFloat(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", value))

			ptr, err := DecodeFloatPtr[float64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", *ptr))

			ptr2, err := DecodeFloatPtr[float64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	for _, expected := range []any{float32(1.1), float64(2.2)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", value))

			ptr, err := DecodeFloatPtr[float64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", *ptr))

			ptr2, err := DecodeFloatPtr[float64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {

		vInt := int(1)
		ptr, err := DecodeFloatPtr[float64](&vInt)
		assertNoError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeFloatPtr[float64](&ptrInt)
		assertNoError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeFloatPtr[float64](&vInt8)
		assertNoError(t, err)
		assertEqual(t, float64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeFloatPtr[float64](&vInt16)
		assertNoError(t, err)
		assertEqual(t, float64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeFloatPtr[float64](&vInt32)
		assertNoError(t, err)
		assertEqual(t, float64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeFloatPtr[float64](&vInt64)
		assertNoError(t, err)
		assertEqual(t, float64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeFloatPtr[float64](&vUint)
		assertNoError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeFloatPtr[float64](&ptrUint)
		assertNoError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeFloatPtr[float64](&vUint8)
		assertNoError(t, err)
		assertEqual(t, float64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeFloatPtr[float64](&vUint16)
		assertNoError(t, err)
		assertEqual(t, float64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeFloatPtr[float64](&vUint32)
		assertNoError(t, err)
		assertEqual(t, float64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeFloatPtr[float64](&vUint64)
		assertNoError(t, err)
		assertEqual(t, float64(vUint64), *ptr)

		vFloat32 := float32(1)
		ptr, err = DecodeFloatPtr[float64](&vFloat32)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		ptrFloat32 := &vFloat32
		ptr, err = DecodeFloatPtr[float64](&ptrFloat32)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		vFloat64 := float64(2.2)
		ptr, err = DecodeFloatPtr[float64](&vFloat64)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat64), *ptr)

		var vAny any = "test"
		_, err = DecodeFloatPtr[float64](&vAny)
		assertError(t, err, "failed to convert Float, got: *interface {} (test)")

		var vFn any = func() {}
		_, err = DecodeFloatPtr[float64](&vFn)
		assertError(t, err, "failed to convert Float")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeFloat[float64](nil)
		assertError(t, err, "the Float value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeFloat[float64]("failure")
		assertError(t, err, "failed to convert Float, got: failure")
	})
}

func TestDecodeComplex(t *testing.T) {
	for _, expected := range []any{complex(1, 1), complex64(2.2)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeComplex[complex128](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", value))

			ptr, err := DecodeComplexPtr[complex128](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", *ptr))

			ptr2, err := DecodeComplexPtr[complex128](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {

		vc64 := complex64(1)
		ptr, err := DecodeComplexPtr[complex64](&vc64)
		assertNoError(t, err)
		assertEqual(t, complex64(vc64), *ptr)

		ptrInt := &vc64
		ptr, err = DecodeComplexPtr[complex64](&ptrInt)
		assertNoError(t, err)
		assertEqual(t, complex64(vc64), *ptr)

		vc128 := complex128(1)
		ptr, err = DecodeComplexPtr[complex64](&vc64)
		assertNoError(t, err)
		assertEqual(t, complex64(vc128), *ptr)

		var vAny any = "test"
		_, err = DecodeComplexPtr[complex128](&vAny)
		assertError(t, err, "failed to convert Complex, got: *interface {} (test)")

		var vFn any = func() {}
		_, err = DecodeComplexPtr[complex128](&vFn)
		assertError(t, err, "failed to convert Complex")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeComplex[complex64](nil)
		assertError(t, err, "the Complex value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeComplex[complex64]("failure")
		assertError(t, err, "failed to convert Complex, got: failure")
	})
}
