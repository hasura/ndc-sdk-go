package utils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-sdk-go/internal"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
}

func assertError(t *testing.T, err error, message string) {
	if err == nil {
		t.Fatal("expected error, got nil")
	} else if !strings.Contains(err.Error(), message) {
		t.Fatalf("expected error with content: %s, got: %s", message, err.Error())
	}
}

func assertEqual(t *testing.T, expected any, reality any) {
	if expected != reality {
		t.Fatalf("not equal, expected: %+v got: %+v", expected, reality)
	}
}

func assertDeepEqual(t *testing.T, expected any, reality any) {
	if !internal.DeepEqual(expected, reality) {
		t.Fatalf("not equal, expected: %+v got: %+v", expected, reality)
	}
}

func TestDecodeBool(t *testing.T) {
	value, err := DecodeBoolean(true)
	assertNoError(t, err)
	assertEqual(t, true, value)

	ptr, err := DecodeNullableBoolean(true)
	assertNoError(t, err)
	assertEqual(t, true, *ptr)

	ptr2, err := DecodeNullableBoolean(ptr)
	assertNoError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeBoolean(nil)
	assertError(t, err, "the Boolean value must not be null")

	_, err = DecodeBoolean("failure")
	assertError(t, err, "failed to convert Boolean, got: string")
}

func TestDecodeBooleanSlice(t *testing.T) {
	value, err := DecodeBooleanSlice([]bool{true, false})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	value, err = DecodeBooleanSlice([]any{false, true})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{false, true}, value)

	value, err = DecodeBooleanSlice(&[]any{true})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true}, value)

	_, err = DecodeBooleanSlice(nil)
	assertError(t, err, "boolean slice must not be null")

	_, err = DecodeBooleanSlice("failure")
	assertError(t, err, "expected a boolean slice, got: string")

	_, err = DecodeBooleanSlice(time.Now())
	assertError(t, err, "expected a boolean slice, got: struct")

	_, err = DecodeBooleanSlice([]any{nil})
	assertError(t, err, "boolean element at 0 must not be null")

	_, err = DecodeBooleanSlice(&[]any{nil})
	assertError(t, err, "boolean element at 0 must not be null")
}

func TestDecodeNullableBooleanSlice(t *testing.T) {
	value, err := DecodeNullableBooleanSlice([]bool{true, false})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	value, err = DecodeNullableBooleanSlice([]*bool{ToPtr(true), ToPtr(false)})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	value, err = DecodeNullableBooleanSlice(&[]*bool{ToPtr(true), ToPtr(false)})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	value, err = DecodeNullableBooleanSlice([]any{false, true})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{false, true}, value)

	value, err = DecodeNullableBooleanSlice(&[]any{true})
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true}, value)

	_, err = DecodeNullableBooleanSlice([]any{nil})
	assertError(t, err, "boolean element at 0 must not be null")

	value, err = DecodeNullableBooleanSlice(nil)
	assertNoError(t, err)
	assertDeepEqual(t, []bool{}, value)

	value, err = DecodeNullableBooleanSlice((*string)(nil))
	assertNoError(t, err)
	assertDeepEqual(t, []bool{}, value)

	_, err = DecodeNullableBooleanSlice([]any{"true"})
	assertError(t, err, "failed to decode boolean element at 0: failed to convert Boolean, got: interface")

	_, err = DecodeNullableBooleanSlice("failure")
	assertError(t, err, "expected a boolean slice, got: string")

	_, err = DecodeNullableBooleanSlice(time.Now())
	assertError(t, err, "expected a boolean slice, got: struct")

}

func TestDecodeString(t *testing.T) {
	value, err := DecodeString("success")
	assertNoError(t, err)
	assertEqual(t, "success", value)

	ptr, err := DecodeNullableString("pointer")
	assertNoError(t, err)
	assertEqual(t, "pointer", *ptr)

	ptr2, err := DecodeNullableString(ptr)
	assertNoError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeString(nil)
	assertError(t, err, "the String value must not be null")

	_, err = DecodeString(0)
	assertError(t, err, "failed to convert String, got: 0")
}

func TestDecodeStringSlice(t *testing.T) {
	value, err := DecodeStringSlice([]string{"foo", "bar"})
	assertNoError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice([]*string{ToPtr("foo"), ToPtr("bar")})
	assertNoError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice(&[]*string{ToPtr("foo")})
	assertNoError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice([]any{"bar", "foo"})
	assertNoError(t, err)
	assertDeepEqual(t, []string{"bar", "foo"}, value)

	value, err = DecodeStringSlice(&[]any{"foo"})
	assertNoError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	_, err = DecodeStringSlice(nil)
	assertError(t, err, "string slice must not be null")

	_, err = DecodeStringSlice("failure")
	assertError(t, err, "expected a string slice, got: string")

	_, err = DecodeStringSlice(time.Now())
	assertError(t, err, "expected a string slice, got: struct")

	_, err = DecodeStringSlice([]any{nil})
	assertError(t, err, "string element at 0 must not be null")

	_, err = DecodeStringSlice(&[]any{nil})
	assertError(t, err, "string element at 0 must not be null")
}

func TestDecodeNullableStringSlice(t *testing.T) {
	_, err := DecodeNullableStringSlice([]string{"foo", "bar"})
	assertNoError(t, err)

	_, err = DecodeNullableStringSlice([]*string{ToPtr("foo"), ToPtr("bar")})
	assertNoError(t, err)

	_, err = DecodeNullableStringSlice(&[]*string{ToPtr("foo")})
	assertNoError(t, err)

	_, err = DecodeNullableStringSlice([]any{"bar", "foo"})
	assertNoError(t, err)

	_, err = DecodeNullableStringSlice(&[]any{"foo"})
	assertNoError(t, err)

	value, err := DecodeNullableStringSlice(nil)
	assertNoError(t, err)
	assertDeepEqual(t, []string{}, value)

	_, err = DecodeNullableStringSlice("failure")
	assertError(t, err, "expected a string slice, got: string")

	_, err = DecodeNullableStringSlice(time.Now())
	assertError(t, err, "expected a string slice, got: struct")

	_, err = DecodeNullableStringSlice([]any{nil})
	assertError(t, err, "string element at 0 must not be null")
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
		ptr, err := DecodeNullableDateTime(nilI64)
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
		ptr, err := DecodeNullableDateTime(&now)
		assertNoError(t, err)
		assertEqual(t, now, *ptr)

		ptr2, err := DecodeNullableDateTime(ptr)
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
		nilValue, err := DecodeNullableDuration(nilStr)
		assertNoError(t, err)
		assertEqual(t, true, IsNil(nilValue))
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeNullableDuration(&duration)
		assertNoError(t, err)
		assertEqual(t, duration, *ptr)

		ptr2, err := DecodeNullableDuration(ptr)
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

	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10), "11"} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeInt[int64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeNullableInt[int64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeNullableInt[int64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vInt := int(1)
		ptr, err := DecodeNullableInt[int64](&vInt)
		assertNoError(t, err)
		assertEqual(t, int64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeNullableInt[int64](&ptrInt)
		assertNoError(t, err)
		assertEqual(t, int64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeNullableInt[int64](&vInt8)
		assertNoError(t, err)
		assertEqual(t, int64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeNullableInt[int64](&vInt16)
		assertNoError(t, err)
		assertEqual(t, int64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeNullableInt[int64](&vInt32)
		assertNoError(t, err)
		assertEqual(t, int64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeNullableInt[int64](&vInt64)
		assertNoError(t, err)
		assertEqual(t, int64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeNullableInt[int64](&vUint)
		assertNoError(t, err)
		assertEqual(t, int64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeNullableInt[int64](&vUint8)
		assertNoError(t, err)
		assertEqual(t, int64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeNullableInt[int64](&vUint16)
		assertNoError(t, err)
		assertEqual(t, int64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeNullableInt[int64](&vUint32)
		assertNoError(t, err)
		assertEqual(t, int64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeNullableInt[int64](&vUint64)
		assertNoError(t, err)
		assertEqual(t, int64(vUint64), *ptr)

		var vAny any = float64(1.1)
		_, err = DecodeNullableInt[int64](&vAny)
		assertError(t, err, "failed to convert integer, got: 1.1")

		var vFn any = func() {}
		_, err = DecodeNullableInt[int64](&vFn)
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

func TestDecodeIntSlice(t *testing.T) {
	for _, tc := range []struct {
		input    any
		expected []int
	}{
		{[]int{1}, []int{1}}, {[]int8{2}, []int{2}}, {[]int16{3}, []int{3}}, {[]int32{4}, []int{4}}, {[]int64{5}, []int{5}},
		{[]uint{6}, []int{6}}, {[]uint8{7}, []int{7}}, {[]uint16{8}, []int{8}}, {[]uint32{9}, []int{9}}, {[]uint64{10}, []int{10}},
		{[]string{"11"}, []int{11}}, {[]float32{12}, []int{12}}, {[]float64{13}, []int{13}},

		{[]*int{ToPtr(1)}, []int{1}}, {[]*int8{ToPtr(int8(2))}, []int{2}}, {[]*int16{ToPtr[int16](3)}, []int{3}}, {[]*int32{ToPtr[int32](4)}, []int{4}}, {[]*int64{ToPtr(int64(5))}, []int{5}},
		{[]*uint{ToPtr[uint](6)}, []int{6}}, {[]*uint8{ToPtr(uint8(7))}, []int{7}}, {[]*uint16{ToPtr(uint16(8))}, []int{8}}, {[]*uint32{ToPtr(uint32(9))}, []int{9}}, {[]*uint64{ToPtr(uint64(10))}, []int{10}},
		{[]*string{ToPtr("11")}, []int{11}}, {[]*float32{ToPtr[float32](12)}, []int{12}}, {[]*float64{ToPtr[float64](13)}, []int{13}},

		{&[]int{1}, []int{1}}, {&[]int8{2}, []int{2}}, {&[]int16{3}, []int{3}}, {&[]int32{4}, []int{4}}, {&[]int64{5}, []int{5}},
		{&[]uint{6}, []int{6}}, {&[]uint8{7}, []int{7}}, {&[]uint16{8}, []int{8}}, {&[]uint32{9}, []int{9}}, {&[]uint64{10}, []int{10}},
		{&[]string{"11"}, []int{11}}, {&[]float32{12}, []int{12}}, {&[]float64{13}, []int{13}},

		{&[]*int{ToPtr(1)}, []int{1}}, {&[]*int8{ToPtr(int8(2))}, []int{2}}, {&[]*int16{ToPtr[int16](3)}, []int{3}}, {&[]*int32{ToPtr[int32](4)}, []int{4}}, {&[]*int64{ToPtr(int64(5))}, []int{5}},
		{&[]*uint{ToPtr[uint](6)}, []int{6}}, {&[]*uint8{ToPtr(uint8(7))}, []int{7}}, {&[]*uint16{ToPtr(uint16(8))}, []int{8}}, {&[]*uint32{ToPtr(uint32(9))}, []int{9}}, {&[]*uint64{ToPtr(uint64(10))}, []int{10}},
		{&[]*string{ToPtr("11")}, []int{11}}, {&[]*float32{ToPtr[float32](12)}, []int{12}}, {&[]*float64{ToPtr[float64](13)}, []int{13}},
	} {
		value, err := DecodeIntSlice[int](tc.input)
		assertNoError(t, err)
		assertDeepEqual(t, tc.expected, value)
	}

	value, err := DecodeIntSlice[int](&[]any{"1", 2})
	assertNoError(t, err)
	assertDeepEqual(t, []int{1, 2}, value)

	_, err = DecodeIntSlice[int](nil)
	assertError(t, err, "the int slice must not be null")

	_, err = DecodeIntSlice[int]("failure")
	assertError(t, err, "expected a number slice, got: string")

	_, err = DecodeIntSlice[int](time.Now())
	assertError(t, err, "expected a number slice, got: struct")

	_, err = DecodeIntSlice[int]([]any{nil})
	assertError(t, err, "number element at 0 must not be null")

	_, err = DecodeIntSlice[int](&[]any{nil})
	assertError(t, err, "number element at 0 must not be null")
}

func TestDecodeNullableIntSlice(t *testing.T) {
	for _, tc := range []struct {
		input    any
		expected []int
	}{
		{[]int{1}, []int{1}}, {[]int8{2}, []int{2}}, {[]int16{3}, []int{3}}, {[]int32{4}, []int{4}}, {[]int64{5}, []int{5}},
		{[]uint{6}, []int{6}}, {[]uint8{7}, []int{7}}, {[]uint16{8}, []int{8}}, {[]uint32{9}, []int{9}}, {[]uint64{10}, []int{10}},
		{[]string{"11"}, []int{11}}, {[]float32{12}, []int{12}}, {[]float64{13}, []int{13}},

		{[]*int{ToPtr(1)}, []int{1}}, {[]*int8{ToPtr(int8(2))}, []int{2}}, {[]*int16{ToPtr[int16](3)}, []int{3}}, {[]*int32{ToPtr[int32](4)}, []int{4}}, {[]*int64{ToPtr(int64(5))}, []int{5}},
		{[]*uint{ToPtr[uint](6)}, []int{6}}, {[]*uint8{ToPtr(uint8(7))}, []int{7}}, {[]*uint16{ToPtr(uint16(8))}, []int{8}}, {[]*uint32{ToPtr(uint32(9))}, []int{9}}, {[]*uint64{ToPtr(uint64(10))}, []int{10}},
		{[]*string{ToPtr("11")}, []int{11}}, {[]*float32{ToPtr[float32](12)}, []int{12}}, {[]*float64{ToPtr[float64](13)}, []int{13}},

		{&[]int{1}, []int{1}}, {&[]int8{2}, []int{2}}, {&[]int16{3}, []int{3}}, {&[]int32{4}, []int{4}}, {&[]int64{5}, []int{5}},
		{&[]uint{6}, []int{6}}, {&[]uint8{7}, []int{7}}, {&[]uint16{8}, []int{8}}, {&[]uint32{9}, []int{9}}, {&[]uint64{10}, []int{10}},
		{&[]string{"11"}, []int{11}}, {&[]float32{12}, []int{12}}, {&[]float64{13}, []int{13}},

		{&[]*int{ToPtr(1)}, []int{1}}, {&[]*int8{ToPtr(int8(2))}, []int{2}}, {&[]*int16{ToPtr[int16](3)}, []int{3}}, {&[]*int32{ToPtr[int32](4)}, []int{4}}, {&[]*int64{ToPtr(int64(5))}, []int{5}},
		{&[]*uint{ToPtr[uint](6)}, []int{6}}, {&[]*uint8{ToPtr(uint8(7))}, []int{7}}, {&[]*uint16{ToPtr(uint16(8))}, []int{8}}, {&[]*uint32{ToPtr(uint32(9))}, []int{9}}, {&[]*uint64{ToPtr(uint64(10))}, []int{10}},
		{&[]*string{ToPtr("11")}, []int{11}}, {&[]*float32{ToPtr[float32](12)}, []int{12}}, {&[]*float64{ToPtr[float64](13)}, []int{13}},
	} {
		value, err := DecodeNullableIntSlice[int](tc.input)
		assertNoError(t, err)
		assertDeepEqual(t, tc.expected, value)
	}

	value, err := DecodeNullableIntSlice[int](&[]any{"1", 2})
	assertNoError(t, err)
	assertDeepEqual(t, []int{1, 2}, value)

	value, err = DecodeNullableIntSlice[int](nil)
	assertNoError(t, err)
	assertDeepEqual(t, []int{}, value)

	_, err = DecodeNullableIntSlice[int]([]any{nil})
	assertError(t, err, "number element at 0 must not be null")

	_, err = DecodeNullableIntSlice[int]("failure")
	assertError(t, err, "expected a number slice, got: string")

	_, err = DecodeNullableIntSlice[int](time.Now())
	assertError(t, err, "expected a number slice, got: struct")

}

func TestDecodeUint(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10), "11"} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeUint[uint64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(value))

			value, err = DecodeUint[uint64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeNullableUint[uint64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeNullableUint[uint64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vUint := uint(1)
		ptr, err := DecodeNullableUint[uint64](&vUint)
		assertNoError(t, err)
		assertEqual(t, uint64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeNullableUint[uint64](&ptrUint)
		assertNoError(t, err)
		assertEqual(t, uint64(vUint), *ptr)

		var vAny any = float64(1.1)
		_, err = DecodeNullableUint[uint64](&vAny)
		assertError(t, err, "failed to convert integer, got: 1.1")
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

func TestDecodeUintSlice(t *testing.T) {
	for _, expected := range []any{
		[]int{1}, []int8{2}, []int16{3}, []int32{4}, []int64{5}, []uint{6}, []uint8{7}, []uint16{8}, []uint32{9}, []uint64{10}, []string{"11"}, []float32{12}, []float64{13},
		&[]int{1}, &[]int8{2}, &[]int16{3}, &[]int32{4}, &[]int64{5}, &[]uint{6}, &[]uint8{7}, &[]uint16{8}, &[]uint32{9}, &[]uint64{10}, &[]string{"11"}, &[]float32{12}, &[]float64{13},
	} {
		value, err := DecodeUintSlice[uint64](expected)
		assertNoError(t, err)
		rexpected := reflect.ValueOf(expected)
		if rexpected.Kind() == reflect.Pointer {
			assertDeepEqual(t, fmt.Sprint(rexpected.Elem().Interface()), fmt.Sprint(value))
		} else {
			assertDeepEqual(t, fmt.Sprint(expected), fmt.Sprint(value))
		}
	}

	value, err := DecodeUintSlice[uint64](&[]any{"1", 2})
	assertNoError(t, err)
	assertDeepEqual(t, []int{1, 2}, value)

	_, err = DecodeUintSlice[uint64](nil)
	assertError(t, err, "the uint slice must not be null")

	_, err = DecodeUintSlice[uint64]("failure")
	assertError(t, err, "expected a number slice, got: string")

	_, err = DecodeUintSlice[uint64](time.Now())
	assertError(t, err, "expected a number slice, got: struct")

	_, err = DecodeUintSlice[uint64]([]any{nil})
	assertError(t, err, "number element at 0 must not be null")

	_, err = DecodeUintSlice[uint64](&[]any{nil})
	assertError(t, err, "number element at 0 must not be null")
}

func TestDecodeFloat(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", value))

			ptr, err := DecodeNullableFloat[float64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", *ptr))

			ptr2, err := DecodeNullableFloat[float64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	for _, expected := range []any{float32(1.1), float64(2.2)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", value))

			ptr, err := DecodeNullableFloat[float64](&expected)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", *ptr))

			ptr2, err := DecodeNullableFloat[float64](ptr)
			assertNoError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {

		vInt := int(1)
		ptr, err := DecodeNullableFloat[float64](&vInt)
		assertNoError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeNullableFloat[float64](&ptrInt)
		assertNoError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeNullableFloat[float64](&vInt8)
		assertNoError(t, err)
		assertEqual(t, float64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeNullableFloat[float64](&vInt16)
		assertNoError(t, err)
		assertEqual(t, float64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeNullableFloat[float64](&vInt32)
		assertNoError(t, err)
		assertEqual(t, float64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeNullableFloat[float64](&vInt64)
		assertNoError(t, err)
		assertEqual(t, float64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeNullableFloat[float64](&vUint)
		assertNoError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeNullableFloat[float64](&ptrUint)
		assertNoError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeNullableFloat[float64](&vUint8)
		assertNoError(t, err)
		assertEqual(t, float64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeNullableFloat[float64](&vUint16)
		assertNoError(t, err)
		assertEqual(t, float64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeNullableFloat[float64](&vUint32)
		assertNoError(t, err)
		assertEqual(t, float64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeNullableFloat[float64](&vUint64)
		assertNoError(t, err)
		assertEqual(t, float64(vUint64), *ptr)

		vFloat32 := float32(1)
		ptr, err = DecodeNullableFloat[float64](&vFloat32)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		ptrFloat32 := &vFloat32
		ptr, err = DecodeNullableFloat[float64](&ptrFloat32)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		vFloat64 := float64(2.2)
		ptr, err = DecodeNullableFloat[float64](&vFloat64)
		assertNoError(t, err)
		assertEqual(t, float64(vFloat64), *ptr)

		var vAny any = "test"
		_, err = DecodeNullableFloat[float64](&vAny)
		assertError(t, err, "failed to convert Float, got: test")

		var vFn any = func() {}
		_, err = DecodeNullableFloat[float64](&vFn)
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

func TestDecodeUUID(t *testing.T) {
	t.Run("decode_value", func(t *testing.T) {
		expected := uuid.NewString()
		value, err := DecodeUUID(expected)
		assertNoError(t, err)
		assertEqual(t, expected, value.String())
	})

	t.Run("decode_pointer", func(t *testing.T) {
		expected := uuid.NewString()
		ptr, err := DecodeNullableUUID(&expected)
		assertNoError(t, err)
		assertEqual(t, expected, (*ptr).String())

		ptr2, err := DecodeNullableUUID(ptr)
		assertNoError(t, err)
		assertEqual(t, ptr.String(), ptr2.String())
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeUUID(nil)
		assertError(t, err, "the uuid value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeUUID(false)
		assertError(t, err, "failed to parse uuid, got: false")
	})
}

func TestDecodeNestedInterface(t *testing.T) {
	boolValue := true
	stringValue := "test"
	intValue := 1
	floatValue := 1.1
	boolSliceValue := []bool{true, false}
	stringSliceValue := []string{"foo", "bar"}
	intSliceValue := []int{1, 2, 3}
	uintSliceValue := []uint{1, 2, 3}
	floatSliceValue := []float64{1.1, 2.2, 3.3}

	fixture := map[string]any{
		"bool":           any(any(boolValue)),
		"string":         any(any(stringValue)),
		"int":            any(any(intValue)),
		"float":          any(any(floatValue)),
		"boolPtr":        any(any(&boolValue)),
		"stringPtr":      any(any(&stringValue)),
		"intPtr":         any(any(&intValue)),
		"floatPtr":       any(any(&floatValue)),
		"boolSlicePtr":   any(any(&boolSliceValue)),
		"stringSlicePtr": any(any(&stringSliceValue)),
		"intSlicePtr":    any(any(&intSliceValue)),
		"uintSlicePtr":   any(any(&uintSliceValue)),
		"floatSlicePtr":  any(any(&floatSliceValue)),
		"array": any([]any{
			any(
				map[string]any{
					"bool":      any(any(boolValue)),
					"string":    any(any(stringValue)),
					"int":       any(any(intValue)),
					"float":     any(any(floatValue)),
					"boolPtr":   any(any(&boolValue)),
					"stringPtr": any(any(&stringValue)),
					"intPtr":    any(any(&intValue)),
					"floatPtr":  any(any(&floatValue)),
				},
			),
		}),
	}

	b, err := GetBoolean(fixture, "bool")
	assertNoError(t, err)
	assertEqual(t, true, b)

	s, err := GetString(fixture, "string")
	assertNoError(t, err)
	assertEqual(t, "test", s)

	i, err := GetInt[int](fixture, "int")
	assertNoError(t, err)
	assertEqual(t, 1, i)

	f, err := GetFloat[float32](fixture, "float")
	assertNoError(t, err)
	assertEqual(t, fmt.Sprint(1.1), fmt.Sprint(f))

	b, err = GetBoolean(fixture, "boolPtr")
	assertNoError(t, err)
	assertEqual(t, true, b)

	s, err = GetString(fixture, "stringPtr")
	assertNoError(t, err)
	assertEqual(t, "test", s)

	i, err = GetInt[int](fixture, "intPtr")
	assertNoError(t, err)
	assertEqual(t, 1, i)

	f, err = GetFloat[float32](fixture, "floatPtr")
	assertNoError(t, err)
	assertEqual(t, fmt.Sprint(1.1), fmt.Sprint(f))

	arrItemAny := fixture["array"].([]any)[0]
	arrItem := arrItemAny.(map[string]any)

	b, err = GetBoolean(arrItem, "bool")
	assertNoError(t, err)
	assertEqual(t, true, b)

	s, err = GetString(arrItem, "string")
	assertNoError(t, err)
	assertEqual(t, "test", s)

	i, err = GetInt[int](arrItem, "int")
	assertNoError(t, err)
	assertEqual(t, 1, i)

	f, err = GetFloat[float32](arrItem, "float")
	assertNoError(t, err)
	assertEqual(t, fmt.Sprint(1.1), fmt.Sprint(f))

	b, err = GetBoolean(arrItem, "boolPtr")
	assertNoError(t, err)
	assertEqual(t, true, b)

	s, err = GetString(arrItem, "stringPtr")
	assertNoError(t, err)
	assertEqual(t, "test", s)

	i, err = GetInt[int](arrItem, "intPtr")
	assertNoError(t, err)
	assertEqual(t, 1, i)

	f, err = GetFloat[float32](arrItem, "floatPtr")
	assertNoError(t, err)
	assertEqual(t, fmt.Sprint(1.1), fmt.Sprint(f))

	j, err := GetArbitraryJSON(arrItem, "int")
	assertNoError(t, err)
	assertEqual(t, 1, j)

	jp, err := GetNullableArbitraryJSON(arrItem, "int")
	assertNoError(t, err)
	assertEqual(t, 1, *jp)

	jp, err = GetNullableArbitraryJSON(arrItem, "floatPtr")
	assertNoError(t, err)
	assertEqual(t, fmt.Sprint(1.1), fmt.Sprint(*jp))

	jp, err = GetNullableArbitraryJSON(arrItem, "not_exist")
	assertNoError(t, err)
	assertEqual(t, true, jp == nil)

	bs, err := GetBooleanSlice(fixture, "boolSlicePtr")
	assertNoError(t, err)
	assertDeepEqual(t, []bool{true, false}, bs)

	_, err = GetNullableBooleanSlice(fixture, "boolSlicePtr")
	assertNoError(t, err)

	ss, err := GetStringSlice(fixture, "stringSlicePtr")
	assertNoError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, ss)

	_, err = GetNullableStringSlice(fixture, "stringSlicePtr")
	assertNoError(t, err)

	is, err := GetIntSlice[int64](fixture, "intSlicePtr")
	assertNoError(t, err)
	assertDeepEqual(t, []int64{1, 2, 3}, is)

	_, err = GetNullableIntSlice[int64](fixture, "intSlicePtr")
	assertNoError(t, err)

	uis, err := GetUintSlice[uint64](fixture, "uintSlicePtr")
	assertNoError(t, err)
	assertDeepEqual(t, []uint64{1, 2, 3}, uis)

	_, err = GetNullableUintSlice[uint64](fixture, "uintSlicePtr")
	assertNoError(t, err)

	fis, err := GetFloatSlice[float32](fixture, "floatSlicePtr")
	assertNoError(t, err)
	assertDeepEqual(t, []float32{1.1, 2.2, 3.3}, fis)

	_, err = GetNullableFloatSlice[float32](fixture, "floatSlicePtr")
	assertNoError(t, err)
}
