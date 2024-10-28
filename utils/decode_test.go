package utils

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestDecodeBool(t *testing.T) {
	value, err := DecodeBoolean(true)
	assert.NilError(t, err)
	assert.Equal(t, true, value)

	ptr, err := DecodeNullableBoolean(true)
	assert.NilError(t, err)
	assert.Equal(t, true, *ptr)

	ptr2, err := DecodeNullableBoolean(ptr)
	assert.NilError(t, err)
	assert.Equal(t, *ptr, *ptr2)

	_, err = DecodeBoolean(nil)
	assert.ErrorContains(t, err, "the Boolean value must not be null")

	_, err = DecodeBoolean("failure")
	assert.ErrorContains(t, err, "failed to convert Boolean, got: string")
}

func TestDecodeBooleanSlice(t *testing.T) {
	value, err := DecodeBooleanSlice([]bool{true, false})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true, false}, value)

	value, err = DecodeBooleanSlice([]any{false, true})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{false, true}, value)

	value, err = DecodeBooleanSlice(&[]any{true})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true}, value)

	_, err = DecodeBooleanSlice(nil)
	assert.ErrorContains(t, err, "boolean slice must not be null")

	_, err = DecodeBooleanSlice("failure")
	assert.ErrorContains(t, err, "expected a boolean slice, got: string")

	_, err = DecodeBooleanSlice(time.Now())
	assert.ErrorContains(t, err, "expected a boolean slice, got: struct")

	_, err = DecodeBooleanSlice([]any{nil})
	assert.ErrorContains(t, err, "boolean element at 0 must not be null")

	_, err = DecodeBooleanSlice(&[]any{nil})
	assert.ErrorContains(t, err, "boolean element at 0 must not be null")
}

func TestDecodeNullableBooleanSlice(t *testing.T) {
	value, err := DecodeNullableBooleanSlice([]bool{true, false})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]*bool{ToPtr(true), ToPtr(false)})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice(&[]*bool{ToPtr(true), ToPtr(false)})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]any{false, true})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{false, true}, *value)

	value, err = DecodeNullableBooleanSlice(&[]any{true})
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true}, *value)

	_, err = DecodeNullableBooleanSlice([]any{nil})
	assert.ErrorContains(t, err, "boolean element at 0 must not be null")

	value, err = DecodeNullableBooleanSlice(nil)
	assert.NilError(t, err)
	assert.Check(t, value == nil)

	value, err = DecodeNullableBooleanSlice((*string)(nil))
	assert.NilError(t, err)
	assert.Check(t, value == nil)

	_, err = DecodeNullableBooleanSlice([]any{"true"})
	assert.ErrorContains(t, err, "failed to decode boolean element at 0: failed to convert Boolean, got: interface")

	_, err = DecodeNullableBooleanSlice("failure")
	assert.ErrorContains(t, err, "expected a boolean slice, got: string")

	_, err = DecodeNullableBooleanSlice(time.Now())
	assert.ErrorContains(t, err, "expected a boolean slice, got: struct")

}

func TestDecodeString(t *testing.T) {
	value, err := DecodeString("success")
	assert.NilError(t, err)
	assert.Equal(t, "success", value)

	ptr, err := DecodeNullableString("pointer")
	assert.NilError(t, err)
	assert.Equal(t, "pointer", *ptr)

	ptr2, err := DecodeNullableString(ptr)
	assert.NilError(t, err)
	assert.Equal(t, *ptr, *ptr2)

	_, err = DecodeString(nil)
	assert.ErrorContains(t, err, "the String value must not be null")

	_, err = DecodeString(0)
	assert.ErrorContains(t, err, "failed to convert String, got: 0")
}

func TestDecodeStringSlice(t *testing.T) {
	value, err := DecodeStringSlice([]string{"foo", "bar"})
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice([]*string{ToPtr("foo"), ToPtr("bar")})
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice(&[]*string{ToPtr("foo")})
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice([]any{"bar", "foo"})
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"bar", "foo"}, value)

	value, err = DecodeStringSlice(&[]any{"foo"})
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"foo"}, value)

	_, err = DecodeStringSlice(nil)
	assert.ErrorContains(t, err, "string slice must not be null")

	_, err = DecodeStringSlice("failure")
	assert.ErrorContains(t, err, "expected a string slice, got: string")

	_, err = DecodeStringSlice(time.Now())
	assert.ErrorContains(t, err, "expected a string slice, got: struct")

	_, err = DecodeStringSlice([]any{nil})
	assert.ErrorContains(t, err, "string element at 0 must not be null")

	_, err = DecodeStringSlice(&[]any{nil})
	assert.ErrorContains(t, err, "string element at 0 must not be null")
}

func TestDecodeNullableStringSlice(t *testing.T) {
	_, err := DecodeNullableStringSlice([]string{"foo", "bar"})
	assert.NilError(t, err)

	_, err = DecodeNullableStringSlice([]*string{ToPtr("foo"), ToPtr("bar")})
	assert.NilError(t, err)

	_, err = DecodeNullableStringSlice(&[]*string{ToPtr("foo")})
	assert.NilError(t, err)

	_, err = DecodeNullableStringSlice([]any{"bar", "foo"})
	assert.NilError(t, err)

	_, err = DecodeNullableStringSlice(&[]any{"foo"})
	assert.NilError(t, err)

	value, err := DecodeNullableStringSlice(nil)
	assert.NilError(t, err)
	assert.Check(t, value == nil)

	_, err = DecodeNullableStringSlice("failure")
	assert.ErrorContains(t, err, "expected a string slice, got: string")

	_, err = DecodeNullableStringSlice(time.Now())
	assert.ErrorContains(t, err, "expected a string slice, got: struct")

	_, err = DecodeNullableStringSlice([]any{nil})
	assert.ErrorContains(t, err, "string element at 0 must not be null")
}

func TestDecodeDateTime(t *testing.T) {
	now := time.Now()
	t.Run("decode_value", func(t *testing.T) {
		value, err := DecodeDateTime(now)
		assert.NilError(t, err)
		assert.Equal(t, now, value)
	})

	t.Run("unix_int", func(t *testing.T) {
		iNow := now.UnixMilli()
		value, err := DecodeDateTime(iNow)
		assert.NilError(t, err)
		assert.Equal(t, now.UnixMilli(), value.UnixMilli())

		value, err = DecodeDateTime(&iNow)
		assert.NilError(t, err)
		assert.Equal(t, now.UnixMilli(), value.UnixMilli())

		var nilI64 *int64 = nil
		ptr, err := DecodeNullableDateTime(nilI64)
		assert.NilError(t, err)
		assert.Equal(t, true, IsNil(ptr))
	})

	t.Run("from_string", func(t *testing.T) {
		nowStr := now.Format(time.RFC3339)
		value, err := DecodeDateTime(nowStr)
		assert.NilError(t, err)
		assert.Equal(t, now.Unix(), value.Unix())

		value, err = DecodeDateTime(&nowStr)
		assert.NilError(t, err)
		assert.Equal(t, now.Unix(), value.Unix())

		invalidStr := "test"
		_, err = DecodeDateTime(invalidStr)
		assert.ErrorContains(t, err, "failed to parse time from string: test")
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeNullableDateTime(&now)
		assert.NilError(t, err)
		assert.Equal(t, now, *ptr)

		ptr2, err := DecodeNullableDateTime(ptr)
		assert.NilError(t, err)
		assert.Equal(t, *ptr, *ptr2)
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeDateTime(nil)
		assert.ErrorContains(t, err, "the DateTime value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeDateTime(false)
		assert.ErrorContains(t, err, "failed to convert DateTime, got: false")
	})
}

func TestDecodeDate(t *testing.T) {
	now := time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC)
	t.Run("decode_value", func(t *testing.T) {
		value, err := DecodeDate(now)
		assert.NilError(t, err)
		assert.Equal(t, now, value)
	})

	t.Run("from_string", func(t *testing.T) {
		nowStr := now.Format(time.DateOnly)
		value, err := DecodeDate(nowStr)
		assert.NilError(t, err)
		assert.Equal(t, now.Unix(), value.Unix())

		value, err = DecodeDate(&nowStr)
		assert.NilError(t, err)
		assert.Equal(t, now.Unix(), value.Unix())

		invalidStr := "test"
		_, err = DecodeDate(invalidStr)
		assert.ErrorContains(t, err, "parsing time \"test\" as \"2006-01-02\": cannot parse \"test\" as \"2006\"")
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeNullableDate(&now)
		assert.NilError(t, err)
		assert.Equal(t, now, *ptr)

		ptr2, err := DecodeNullableDate(ptr)
		assert.NilError(t, err)
		assert.Equal(t, *ptr, *ptr2)
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeDate(nil)
		assert.ErrorContains(t, err, "the Date value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeDate(false)
		assert.ErrorContains(t, err, "failed to convert Date, got: false")
	})
}

func TestDecodeDuration(t *testing.T) {
	duration := 10 * time.Second
	t.Run("decode_value", func(t *testing.T) {
		value, err := DecodeDuration(duration)
		assert.NilError(t, err)
		assert.Equal(t, duration, value)
	})

	t.Run("unix_int", func(t *testing.T) {
		iDuration := int64(duration)
		value, err := DecodeDuration(iDuration)
		assert.NilError(t, err)
		assert.Equal(t, duration, value)

		value, err = DecodeDuration(&iDuration)
		assert.NilError(t, err)
		assert.Equal(t, duration, value)
	})

	t.Run("from_string", func(t *testing.T) {
		durationStr := "10s"
		value, err := DecodeDuration(durationStr)
		assert.NilError(t, err)
		assert.Equal(t, duration, value)

		value, err = DecodeDuration(&durationStr)
		assert.NilError(t, err)
		assert.Equal(t, duration, value)

		invalidStr := "test"
		_, err = DecodeDuration(invalidStr)
		assert.ErrorContains(t, err, "not a valid duration string: \"test\"")

		_, err = DecodeDuration(&invalidStr)
		assert.ErrorContains(t, err, "not a valid duration string: \"test\"")

		var nilStr *string
		nilValue, err := DecodeNullableDuration(nilStr)
		assert.NilError(t, err)
		assert.Equal(t, true, IsNil(nilValue))
	})

	t.Run("decode_pointer", func(t *testing.T) {
		ptr, err := DecodeNullableDuration(&duration)
		assert.NilError(t, err)
		assert.Equal(t, duration, *ptr)

		ptr2, err := DecodeNullableDuration(ptr)
		assert.NilError(t, err)
		assert.Equal(t, *ptr, *ptr2)
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeDuration(nil)
		assert.ErrorContains(t, err, "the Duration value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeDuration(false)
		assert.ErrorContains(t, err, "failed to convert Duration, got: false")
	})
}

func TestDecodeInt(t *testing.T) {

	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10), "11"} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeInt[int64](expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeNullableInt[int64](&expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeNullableInt[int64](ptr)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vInt := int(1)
		ptr, err := DecodeNullableInt[int64](&vInt)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeNullableInt[int64](&ptrInt)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeNullableInt[int64](&vInt8)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeNullableInt[int64](&vInt16)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeNullableInt[int64](&vInt32)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeNullableInt[int64](&vInt64)
		assert.NilError(t, err)
		assert.Equal(t, int64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeNullableInt[int64](&vUint)
		assert.NilError(t, err)
		assert.Equal(t, int64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeNullableInt[int64](&vUint8)
		assert.NilError(t, err)
		assert.Equal(t, int64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeNullableInt[int64](&vUint16)
		assert.NilError(t, err)
		assert.Equal(t, int64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeNullableInt[int64](&vUint32)
		assert.NilError(t, err)
		assert.Equal(t, int64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeNullableInt[int64](&vUint64)
		assert.NilError(t, err)
		assert.Equal(t, int64(vUint64), *ptr)

		var vAny any = "failure"
		_, err = DecodeNullableInt[int64](&vAny)
		assert.ErrorContains(t, err, "failed to convert integer, got failure")

		var vFn any = func() {}
		_, err = DecodeNullableInt[int64](&vFn)
		assert.ErrorContains(t, err, "failed to convert integer")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeInt[rune](nil)
		assert.ErrorContains(t, err, "the Int value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeInt[rune]("failure")
		assert.ErrorContains(t, err, "failed to convert integer, got failure")
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
		assert.NilError(t, err)
		assert.DeepEqual(t, tc.expected, value)
	}

	value, err := DecodeIntSlice[int](&[]any{"1", 2})
	assert.NilError(t, err)
	assert.DeepEqual(t, []int{1, 2}, value)

	_, err = DecodeIntSlice[int](nil)
	assert.ErrorContains(t, err, "the int slice must not be null")

	_, err = DecodeIntSlice[int]("failure")
	assert.ErrorContains(t, err, "expected a number slice, got: string")

	_, err = DecodeIntSlice[int](time.Now())
	assert.ErrorContains(t, err, "expected a number slice, got: struct")

	_, err = DecodeIntSlice[int]([]any{nil})
	assert.ErrorContains(t, err, "number element at 0 must not be null")

	_, err = DecodeIntSlice[int](&[]any{nil})
	assert.ErrorContains(t, err, "number element at 0 must not be null")
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
		assert.NilError(t, err)
		assert.DeepEqual(t, tc.expected, *value)
	}

	value, err := DecodeNullableIntSlice[int](&[]any{"1", 2})
	assert.NilError(t, err)
	assert.DeepEqual(t, []int{1, 2}, *value)

	value, err = DecodeNullableIntSlice[int](nil)
	assert.NilError(t, err)
	assert.Check(t, value == nil)

	_, err = DecodeNullableIntSlice[int]([]any{nil})
	assert.ErrorContains(t, err, "number element at 0 must not be null")

	_, err = DecodeNullableIntSlice[int]("failure")
	assert.ErrorContains(t, err, "expected a number slice, got: string")

	_, err = DecodeNullableIntSlice[int](time.Now())
	assert.ErrorContains(t, err, "expected a number slice, got: struct")

}

func TestDecodeUint(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10), "11"} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeUint[uint64](expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(value))

			value, err = DecodeUint[uint64](expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(value))

			ptr, err := DecodeNullableUint[uint64](&expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(*ptr))

			ptr2, err := DecodeNullableUint[uint64](ptr)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(*ptr), fmt.Sprint(*ptr2))
		})
	}

	t.Run("decode_pointers", func(t *testing.T) {
		vUint := uint(1)
		ptr, err := DecodeNullableUint[uint64](&vUint)
		assert.NilError(t, err)
		assert.Equal(t, uint64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeNullableUint[uint64](&ptrUint)
		assert.NilError(t, err)
		assert.Equal(t, uint64(vUint), *ptr)

		var vAny any = string("1,1")
		_, err = DecodeNullableUint[uint64](&vAny)
		assert.ErrorContains(t, err, "failed to convert integer, got 1,1")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeUint[uint](nil)
		assert.ErrorContains(t, err, "the Uint value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeUint[uint]("failure")
		assert.ErrorContains(t, err, "failed to convert integer, got failure")
	})
}

func TestDecodeUintSlice(t *testing.T) {
	for _, expected := range []any{
		[]int{1}, []int8{2}, []int16{3}, []int32{4}, []int64{5}, []uint{6}, []uint8{7}, []uint16{8}, []uint32{9}, []uint64{10}, []string{"11"}, []float32{12}, []float64{13},
		&[]int{1}, &[]int8{2}, &[]int16{3}, &[]int32{4}, &[]int64{5}, &[]uint{6}, &[]uint8{7}, &[]uint16{8}, &[]uint32{9}, &[]uint64{10}, &[]string{"11"}, &[]float32{12}, &[]float64{13},
	} {
		value, err := DecodeUintSlice[uint64](expected)
		assert.NilError(t, err)
		rexpected := reflect.ValueOf(expected)
		if rexpected.Kind() == reflect.Pointer {
			assert.DeepEqual(t, fmt.Sprint(rexpected.Elem().Interface()), fmt.Sprint(value))
		} else {
			assert.DeepEqual(t, fmt.Sprint(expected), fmt.Sprint(value))
		}
	}

	value, err := DecodeUintSlice[uint64](&[]any{"1", 2})
	assert.NilError(t, err)
	assert.DeepEqual(t, []uint64{1, 2}, value)

	_, err = DecodeUintSlice[uint64](nil)
	assert.ErrorContains(t, err, "the uint slice must not be null")

	_, err = DecodeUintSlice[uint64]("failure")
	assert.ErrorContains(t, err, "expected a number slice, got: string")

	_, err = DecodeUintSlice[uint64](time.Now())
	assert.ErrorContains(t, err, "expected a number slice, got: struct")

	_, err = DecodeUintSlice[uint64]([]any{nil})
	assert.ErrorContains(t, err, "number element at 0 must not be null")

	_, err = DecodeUintSlice[uint64](&[]any{nil})
	assert.ErrorContains(t, err, "number element at 0 must not be null")
}

func TestDecodeFloat(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", value))

			ptr, err := DecodeNullableFloat[float64](&expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", *ptr))

			ptr2, err := DecodeNullableFloat[float64](ptr)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	for _, expected := range []any{float32(1.1), float64(2.2)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {

			value, err := DecodeFloat[float64](expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", value))

			ptr, err := DecodeNullableFloat[float64](&expected)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", *ptr))

			ptr2, err := DecodeNullableFloat[float64](ptr)
			assert.NilError(t, err)
			assert.Equal(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	t.Run("decode_string", func(t *testing.T) {

		expected := "0"
		value, err := DecodeFloat[float64](expected)
		assert.NilError(t, err)
		assert.Equal(t, "0", fmt.Sprintf("%.0f", value))

		ptr, err := DecodeNullableFloat[float64](&expected)
		assert.NilError(t, err)
		assert.Equal(t, expected, fmt.Sprintf("%.0f", *ptr))

		ptr2, err := DecodeNullableFloat[float64](ptr)
		assert.NilError(t, err)
		assert.Equal(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
	})

	t.Run("decode_pointers", func(t *testing.T) {

		vInt := int(1)
		ptr, err := DecodeNullableFloat[float64](&vInt)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeNullableFloat[float64](&ptrInt)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeNullableFloat[float64](&vInt8)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeNullableFloat[float64](&vInt16)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeNullableFloat[float64](&vInt32)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeNullableFloat[float64](&vInt64)
		assert.NilError(t, err)
		assert.Equal(t, float64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeNullableFloat[float64](&vUint)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeNullableFloat[float64](&ptrUint)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeNullableFloat[float64](&vUint8)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeNullableFloat[float64](&vUint16)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeNullableFloat[float64](&vUint32)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeNullableFloat[float64](&vUint64)
		assert.NilError(t, err)
		assert.Equal(t, float64(vUint64), *ptr)

		vFloat32 := float32(1)
		ptr, err = DecodeNullableFloat[float64](&vFloat32)
		assert.NilError(t, err)
		assert.Equal(t, float64(vFloat32), *ptr)

		ptrFloat32 := &vFloat32
		ptr, err = DecodeNullableFloat[float64](&ptrFloat32)
		assert.NilError(t, err)
		assert.Equal(t, float64(vFloat32), *ptr)

		vFloat64 := float64(2.2)
		ptr, err = DecodeNullableFloat[float64](&vFloat64)
		assert.NilError(t, err)
		assert.Equal(t, float64(vFloat64), *ptr)

		var vAny any = "test"
		_, err = DecodeNullableFloat[float64](&vAny)
		assert.ErrorContains(t, err, "failed to convert Float, got: test")

		var vFn any = func() {}
		_, err = DecodeNullableFloat[float64](&vFn)
		assert.ErrorContains(t, err, "failed to convert Float")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeFloat[float64](nil)
		assert.ErrorContains(t, err, "the Float value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeFloat[float64]("failure")
		assert.ErrorContains(t, err, "failed to convert Float, got: failure")
	})
}

func TestDecodeUUID(t *testing.T) {
	t.Run("decode_value", func(t *testing.T) {
		expected := uuid.NewString()
		value, err := DecodeUUID(expected)
		assert.NilError(t, err)
		assert.Equal(t, expected, value.String())
	})

	t.Run("decode_pointer", func(t *testing.T) {
		expected := uuid.NewString()
		ptr, err := DecodeNullableUUID(&expected)
		assert.NilError(t, err)
		assert.Equal(t, expected, (*ptr).String())

		ptr2, err := DecodeNullableUUID(ptr)
		assert.NilError(t, err)
		assert.Equal(t, ptr.String(), ptr2.String())
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeUUID(nil)
		assert.ErrorContains(t, err, "the uuid value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeUUID(false)
		assert.ErrorContains(t, err, "failed to parse uuid, got: false")
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
	assert.NilError(t, err)
	assert.Equal(t, true, b)

	s, err := GetString(fixture, "string")
	assert.NilError(t, err)
	assert.Equal(t, "test", s)

	i, err := GetInt[int](fixture, "int")
	assert.NilError(t, err)
	assert.Equal(t, 1, i)

	f, err := GetFloat[float32](fixture, "float")
	assert.NilError(t, err)
	assert.Equal(t, fmt.Sprint(1.1), fmt.Sprint(f))

	b, err = GetBoolean(fixture, "boolPtr")
	assert.NilError(t, err)
	assert.Equal(t, true, b)

	s, err = GetString(fixture, "stringPtr")
	assert.NilError(t, err)
	assert.Equal(t, "test", s)

	i, err = GetInt[int](fixture, "intPtr")
	assert.NilError(t, err)
	assert.Equal(t, 1, i)

	f, err = GetFloat[float32](fixture, "floatPtr")
	assert.NilError(t, err)
	assert.Equal(t, fmt.Sprint(1.1), fmt.Sprint(f))

	arrItemAny := fixture["array"].([]any)[0]
	arrItem := arrItemAny.(map[string]any)

	b, err = GetBoolean(arrItem, "bool")
	assert.NilError(t, err)
	assert.Equal(t, true, b)

	s, err = GetString(arrItem, "string")
	assert.NilError(t, err)
	assert.Equal(t, "test", s)

	i, err = GetInt[int](arrItem, "int")
	assert.NilError(t, err)
	assert.Equal(t, 1, i)

	f, err = GetFloat[float32](arrItem, "float")
	assert.NilError(t, err)
	assert.Equal(t, fmt.Sprint(1.1), fmt.Sprint(f))

	b, err = GetBoolean(arrItem, "boolPtr")
	assert.NilError(t, err)
	assert.Equal(t, true, b)

	s, err = GetString(arrItem, "stringPtr")
	assert.NilError(t, err)
	assert.Equal(t, "test", s)

	i, err = GetInt[int](arrItem, "intPtr")
	assert.NilError(t, err)
	assert.Equal(t, 1, i)

	f, err = GetFloat[float32](arrItem, "floatPtr")
	assert.NilError(t, err)
	assert.Equal(t, fmt.Sprint(1.1), fmt.Sprint(f))

	j, err := GetArbitraryJSON(arrItem, "int")
	assert.NilError(t, err)
	assert.Equal(t, 1, j)

	jp, err := GetNullableArbitraryJSON(arrItem, "int")
	assert.NilError(t, err)
	assert.Equal(t, 1, *jp)

	jp, err = GetNullableArbitraryJSON(arrItem, "floatPtr")
	assert.NilError(t, err)
	assert.Equal(t, fmt.Sprint(1.1), fmt.Sprint(*jp))

	jp, err = GetNullableArbitraryJSON(arrItem, "not_exist")
	assert.NilError(t, err)
	assert.Equal(t, true, jp == nil)

	bs, err := GetBooleanSlice(fixture, "boolSlicePtr")
	assert.NilError(t, err)
	assert.DeepEqual(t, []bool{true, false}, bs)

	_, err = GetNullableBooleanSlice(fixture, "boolSlicePtr")
	assert.NilError(t, err)

	ss, err := GetStringSlice(fixture, "stringSlicePtr")
	assert.NilError(t, err)
	assert.DeepEqual(t, []string{"foo", "bar"}, ss)

	_, err = GetNullableStringSlice(fixture, "stringSlicePtr")
	assert.NilError(t, err)

	is, err := GetIntSlice[int64](fixture, "intSlicePtr")
	assert.NilError(t, err)
	assert.DeepEqual(t, []int64{1, 2, 3}, is)

	_, err = GetNullableIntSlice[int64](fixture, "intSlicePtr")
	assert.NilError(t, err)

	uis, err := GetUintSlice[uint64](fixture, "uintSlicePtr")
	assert.NilError(t, err)
	assert.DeepEqual(t, []uint64{1, 2, 3}, uis)

	_, err = GetNullableUintSlice[uint64](fixture, "uintSlicePtr")
	assert.NilError(t, err)

	fis, err := GetFloatSlice[float32](fixture, "floatSlicePtr")
	assert.NilError(t, err)
	assert.DeepEqual(t, []float32{1.1, 2.2, 3.3}, fis)

	_, err = GetNullableFloatSlice[float32](fixture, "floatSlicePtr")
	assert.NilError(t, err)
}
