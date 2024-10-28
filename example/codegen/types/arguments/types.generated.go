// Code generated by github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go, DO NOT EDIT.
package arguments

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/scalar"
	"github.com/hasura/ndc-sdk-go/utils"
	"time"
)

// FromValue decodes values from map
func (j *GetCustomHeadersInput) FromValue(input map[string]any) error {
	var err error
	j.ID, err = utils.GetUUID(input, "id")
	if err != nil {
		return err
	}
	j.Num, err = utils.GetInt[int](input, "num")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *GetTypesArguments) FromValue(input map[string]any) error {
	var err error
	j.ArrayBigInt, err = utils.DecodeObjectValue[[]scalar.BigInt](input, "ArrayBigInt")
	if err != nil {
		return err
	}
	j.ArrayBigIntPtr, err = utils.DecodeObjectValue[[]*scalar.BigInt](input, "ArrayBigIntPtr")
	if err != nil {
		return err
	}
	j.ArrayBool, err = utils.GetBooleanSlice(input, "ArrayBool")
	if err != nil {
		return err
	}
	j.ArrayBoolPtr, err = utils.GetBooleanPtrSlice(input, "ArrayBoolPtr")
	if err != nil {
		return err
	}
	j.ArrayFloat32, err = utils.GetFloatSlice[float32](input, "ArrayFloat32")
	if err != nil {
		return err
	}
	j.ArrayFloat32Ptr, err = utils.GetFloatPtrSlice[float32](input, "ArrayFloat32Ptr")
	if err != nil {
		return err
	}
	j.ArrayFloat64, err = utils.GetFloatSlice[float64](input, "ArrayFloat64")
	if err != nil {
		return err
	}
	j.ArrayFloat64Ptr, err = utils.GetFloatPtrSlice[float64](input, "ArrayFloat64Ptr")
	if err != nil {
		return err
	}
	j.ArrayInt, err = utils.GetIntSlice[int](input, "ArrayInt")
	if err != nil {
		return err
	}
	j.ArrayInt16, err = utils.GetIntSlice[int16](input, "ArrayInt16")
	if err != nil {
		return err
	}
	j.ArrayInt16Ptr, err = utils.GetIntPtrSlice[int16](input, "ArrayInt16Ptr")
	if err != nil {
		return err
	}
	j.ArrayInt32, err = utils.GetIntSlice[int32](input, "ArrayInt32")
	if err != nil {
		return err
	}
	j.ArrayInt32Ptr, err = utils.GetIntPtrSlice[int32](input, "ArrayInt32Ptr")
	if err != nil {
		return err
	}
	j.ArrayInt64, err = utils.GetIntSlice[int64](input, "ArrayInt64")
	if err != nil {
		return err
	}
	j.ArrayInt64Ptr, err = utils.GetIntPtrSlice[int64](input, "ArrayInt64Ptr")
	if err != nil {
		return err
	}
	j.ArrayInt8, err = utils.GetIntSlice[int8](input, "ArrayInt8")
	if err != nil {
		return err
	}
	j.ArrayInt8Ptr, err = utils.GetIntPtrSlice[int8](input, "ArrayInt8Ptr")
	if err != nil {
		return err
	}
	j.ArrayIntPtr, err = utils.GetIntPtrSlice[int](input, "ArrayIntPtr")
	if err != nil {
		return err
	}
	j.ArrayJSON, err = utils.GetArbitraryJSONSlice(input, "ArrayJSON")
	if err != nil {
		return err
	}
	j.ArrayJSONPtr, err = utils.GetArbitraryJSONPtrSlice(input, "ArrayJSONPtr")
	if err != nil {
		return err
	}
	j.ArrayMap, err = utils.DecodeObjectValue[[]map[string]any](input, "ArrayMap")
	if err != nil {
		return err
	}
	j.ArrayMapPtr, err = utils.DecodeNullableObjectValue[[]map[string]any](input, "ArrayMapPtr")
	if err != nil {
		return err
	}
	j.ArrayObject, err = utils.DecodeObjectValue[[]struct {
		Content string `json:"content"`
	}](input, "ArrayObject")
	if err != nil {
		return err
	}
	j.ArrayObjectPtr, err = utils.DecodeNullableObjectValue[[]struct {
		Content string `json:"content"`
	}](input, "ArrayObjectPtr")
	if err != nil {
		return err
	}
	j.ArrayRawJSON, err = utils.GetRawJSONSlice(input, "ArrayRawJSON")
	if err != nil {
		return err
	}
	j.ArrayRawJSONPtr, err = utils.GetRawJSONPtrSlice(input, "ArrayRawJSONPtr")
	if err != nil {
		return err
	}
	j.ArrayString, err = utils.GetStringSlice(input, "ArrayString")
	if err != nil {
		return err
	}
	j.ArrayStringPtr, err = utils.GetStringPtrSlice(input, "ArrayStringPtr")
	if err != nil {
		return err
	}
	j.ArrayTime, err = utils.GetDateTimeSlice(input, "ArrayTime")
	if err != nil {
		return err
	}
	j.ArrayTimePtr, err = utils.GetDateTimePtrSlice(input, "ArrayTimePtr")
	if err != nil {
		return err
	}
	j.ArrayUUID, err = utils.GetUUIDSlice(input, "ArrayUUID")
	if err != nil {
		return err
	}
	j.ArrayUUIDPtr, err = utils.GetUUIDPtrSlice(input, "ArrayUUIDPtr")
	if err != nil {
		return err
	}
	j.ArrayUint, err = utils.GetUintSlice[uint](input, "ArrayUint")
	if err != nil {
		return err
	}
	j.ArrayUint16, err = utils.GetUintSlice[uint16](input, "ArrayUint16")
	if err != nil {
		return err
	}
	j.ArrayUint16Ptr, err = utils.GetUintPtrSlice[uint16](input, "ArrayUint16Ptr")
	if err != nil {
		return err
	}
	j.ArrayUint32, err = utils.GetUintSlice[uint32](input, "ArrayUint32")
	if err != nil {
		return err
	}
	j.ArrayUint32Ptr, err = utils.GetUintPtrSlice[uint32](input, "ArrayUint32Ptr")
	if err != nil {
		return err
	}
	j.ArrayUint64, err = utils.GetUintSlice[uint64](input, "ArrayUint64")
	if err != nil {
		return err
	}
	j.ArrayUint64Ptr, err = utils.GetUintPtrSlice[uint64](input, "ArrayUint64Ptr")
	if err != nil {
		return err
	}
	j.ArrayUint8, err = utils.GetUintSlice[uint8](input, "ArrayUint8")
	if err != nil {
		return err
	}
	j.ArrayUint8Ptr, err = utils.GetUintPtrSlice[uint8](input, "ArrayUint8Ptr")
	if err != nil {
		return err
	}
	j.ArrayUintPtr, err = utils.GetUintPtrSlice[uint](input, "ArrayUintPtr")
	if err != nil {
		return err
	}
	j.BigInt, err = utils.DecodeObjectValue[scalar.BigInt](input, "BigInt")
	if err != nil {
		return err
	}
	j.BigIntPtr, err = utils.DecodeNullableObjectValue[scalar.BigInt](input, "BigIntPtr")
	if err != nil {
		return err
	}
	j.Bool, err = utils.GetBoolean(input, "Bool")
	if err != nil {
		return err
	}
	j.BoolPtr, err = utils.GetNullableBoolean(input, "BoolPtr")
	if err != nil {
		return err
	}
	j.Bytes, err = utils.DecodeObjectValue[scalar.Bytes](input, "Bytes")
	if err != nil {
		return err
	}
	j.BytesPtr, err = utils.DecodeNullableObjectValue[scalar.Bytes](input, "BytesPtr")
	if err != nil {
		return err
	}
	j.CustomScalar, err = utils.DecodeObjectValue[types.CommentText](input, "CustomScalar")
	if err != nil {
		return err
	}
	j.CustomScalarPtr, err = utils.DecodeNullableObjectValue[types.CommentText](input, "CustomScalarPtr")
	if err != nil {
		return err
	}
	j.Date, err = utils.DecodeObjectValue[scalar.Date](input, "Date")
	if err != nil {
		return err
	}
	j.DatePtr, err = utils.DecodeNullableObjectValue[scalar.Date](input, "DatePtr")
	if err != nil {
		return err
	}
	j.Duration, err = utils.DecodeObjectValue[scalar.Duration](input, "Duration")
	if err != nil {
		return err
	}
	j.Enum, err = utils.DecodeObjectValue[types.SomeEnum](input, "Enum")
	if err != nil {
		return err
	}
	j.EnumPtr, err = utils.DecodeNullableObjectValue[types.SomeEnum](input, "EnumPtr")
	if err != nil {
		return err
	}
	j.Float32, err = utils.GetFloat[float32](input, "Float32")
	if err != nil {
		return err
	}
	j.Float32Ptr, err = utils.GetNullableFloat[float32](input, "Float32Ptr")
	if err != nil {
		return err
	}
	j.Float64, err = utils.GetFloat[float64](input, "Float64")
	if err != nil {
		return err
	}
	j.Float64Ptr, err = utils.GetNullableFloat[float64](input, "Float64Ptr")
	if err != nil {
		return err
	}
	j.Int, err = utils.GetInt[int](input, "Int")
	if err != nil {
		return err
	}
	j.Int16, err = utils.GetInt[int16](input, "Int16")
	if err != nil {
		return err
	}
	j.Int16Ptr, err = utils.GetNullableInt[int16](input, "Int16Ptr")
	if err != nil {
		return err
	}
	j.Int32, err = utils.GetInt[int32](input, "Int32")
	if err != nil {
		return err
	}
	j.Int32Ptr, err = utils.GetNullableInt[int32](input, "Int32Ptr")
	if err != nil {
		return err
	}
	j.Int64, err = utils.GetInt[int64](input, "Int64")
	if err != nil {
		return err
	}
	j.Int64Ptr, err = utils.GetNullableInt[int64](input, "Int64Ptr")
	if err != nil {
		return err
	}
	j.Int8, err = utils.GetInt[int8](input, "Int8")
	if err != nil {
		return err
	}
	j.Int8Ptr, err = utils.GetNullableInt[int8](input, "Int8Ptr")
	if err != nil {
		return err
	}
	j.IntPtr, err = utils.GetNullableInt[int](input, "IntPtr")
	if err != nil {
		return err
	}
	j.JSON, err = utils.GetArbitraryJSON(input, "JSON")
	if err != nil {
		return err
	}
	j.JSONPtr, err = utils.GetNullableArbitraryJSON(input, "JSONPtr")
	if err != nil {
		return err
	}
	j.Map, err = utils.DecodeObjectValue[map[string]any](input, "Map")
	if err != nil {
		return err
	}
	j.MapPtr, err = utils.DecodeNullableObjectValue[map[string]any](input, "MapPtr")
	if err != nil {
		return err
	}
	j.NamedArray, err = utils.DecodeObjectValue[[]types.Author](input, "NamedArray")
	if err != nil {
		return err
	}
	j.NamedArrayPtr, err = utils.DecodeNullableObjectValue[[]types.Author](input, "NamedArrayPtr")
	if err != nil {
		return err
	}
	j.NamedObject, err = utils.DecodeObjectValue[types.Author](input, "NamedObject")
	if err != nil {
		return err
	}
	j.NamedObjectPtr, err = utils.DecodeNullableObjectValue[types.Author](input, "NamedObjectPtr")
	if err != nil {
		return err
	}
	j.Object, err = utils.DecodeObjectValue[struct {
		ID           uuid.UUID                               `json:"id"`
		CreatedAt    time.Time                               `json:"created_at"`
		GenericField types.CustomHeadersResult[types.Author] `json:"generic_field,omitempty"`
	}](input, "Object")
	if err != nil {
		return err
	}
	j.ObjectPtr, err = utils.DecodeNullableObjectValue[struct {
		Long int
		Lat  int
	}](input, "ObjectPtr")
	if err != nil {
		return err
	}
	j.PtrArrayBigInt, err = utils.DecodeNullableObjectValue[[]scalar.BigInt](input, "PtrArrayBigInt")
	if err != nil {
		return err
	}
	j.PtrArrayBigIntPtr, err = utils.DecodeNullableObjectValue[[]*scalar.BigInt](input, "PtrArrayBigIntPtr")
	if err != nil {
		return err
	}
	j.PtrArrayBool, err = utils.DecodeNullableObjectValue[[]bool](input, "PtrArrayBool")
	if err != nil {
		return err
	}
	j.PtrArrayBoolPtr, err = utils.GetNullableBooleanPtrSlice(input, "PtrArrayBoolPtr")
	if err != nil {
		return err
	}
	j.PtrArrayFloat32, err = utils.DecodeNullableObjectValue[[]float32](input, "PtrArrayFloat32")
	if err != nil {
		return err
	}
	j.PtrArrayFloat32Ptr, err = utils.GetNullableFloatPtrSlice[float32](input, "PtrArrayFloat32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayFloat64, err = utils.DecodeNullableObjectValue[[]float64](input, "PtrArrayFloat64")
	if err != nil {
		return err
	}
	j.PtrArrayFloat64Ptr, err = utils.GetNullableFloatPtrSlice[float64](input, "PtrArrayFloat64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt, err = utils.DecodeNullableObjectValue[[]int](input, "PtrArrayInt")
	if err != nil {
		return err
	}
	j.PtrArrayInt16, err = utils.DecodeNullableObjectValue[[]int16](input, "PtrArrayInt16")
	if err != nil {
		return err
	}
	j.PtrArrayInt16Ptr, err = utils.GetNullableIntPtrSlice[int16](input, "PtrArrayInt16Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt32, err = utils.DecodeNullableObjectValue[[]int32](input, "PtrArrayInt32")
	if err != nil {
		return err
	}
	j.PtrArrayInt32Ptr, err = utils.GetNullableIntPtrSlice[int32](input, "PtrArrayInt32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt64, err = utils.DecodeNullableObjectValue[[]int64](input, "PtrArrayInt64")
	if err != nil {
		return err
	}
	j.PtrArrayInt64Ptr, err = utils.GetNullableIntPtrSlice[int64](input, "PtrArrayInt64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt8, err = utils.DecodeNullableObjectValue[[]int8](input, "PtrArrayInt8")
	if err != nil {
		return err
	}
	j.PtrArrayInt8Ptr, err = utils.GetNullableIntPtrSlice[int8](input, "PtrArrayInt8Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayIntPtr, err = utils.GetNullableIntPtrSlice[int](input, "PtrArrayIntPtr")
	if err != nil {
		return err
	}
	j.PtrArrayJSON, err = utils.DecodeNullableObjectValue[[]any](input, "PtrArrayJSON")
	if err != nil {
		return err
	}
	j.PtrArrayJSONPtr, err = utils.GetNullableArbitraryJSONPtrSlice(input, "PtrArrayJSONPtr")
	if err != nil {
		return err
	}
	j.PtrArrayRawJSON, err = utils.DecodeNullableObjectValue[[]json.RawMessage](input, "PtrArrayRawJSON")
	if err != nil {
		return err
	}
	j.PtrArrayRawJSONPtr, err = utils.GetNullableRawJSONPtrSlice(input, "PtrArrayRawJSONPtr")
	if err != nil {
		return err
	}
	j.PtrArrayString, err = utils.DecodeNullableObjectValue[[]string](input, "PtrArrayString")
	if err != nil {
		return err
	}
	j.PtrArrayStringPtr, err = utils.GetNullableStringPtrSlice(input, "PtrArrayStringPtr")
	if err != nil {
		return err
	}
	j.PtrArrayTime, err = utils.DecodeNullableObjectValue[[]time.Time](input, "PtrArrayTime")
	if err != nil {
		return err
	}
	j.PtrArrayTimePtr, err = utils.GetNullableDateTimePtrSlice(input, "PtrArrayTimePtr")
	if err != nil {
		return err
	}
	j.PtrArrayUUID, err = utils.DecodeNullableObjectValue[[]uuid.UUID](input, "PtrArrayUUID")
	if err != nil {
		return err
	}
	j.PtrArrayUUIDPtr, err = utils.GetNullableUUIDPtrSlice(input, "PtrArrayUUIDPtr")
	if err != nil {
		return err
	}
	j.PtrArrayUint, err = utils.DecodeNullableObjectValue[[]uint](input, "PtrArrayUint")
	if err != nil {
		return err
	}
	j.PtrArrayUint16, err = utils.DecodeNullableObjectValue[[]uint16](input, "PtrArrayUint16")
	if err != nil {
		return err
	}
	j.PtrArrayUint16Ptr, err = utils.GetNullableUintPtrSlice[uint16](input, "PtrArrayUint16Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint32, err = utils.DecodeNullableObjectValue[[]uint32](input, "PtrArrayUint32")
	if err != nil {
		return err
	}
	j.PtrArrayUint32Ptr, err = utils.GetNullableUintPtrSlice[uint32](input, "PtrArrayUint32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint64, err = utils.DecodeNullableObjectValue[[]uint64](input, "PtrArrayUint64")
	if err != nil {
		return err
	}
	j.PtrArrayUint64Ptr, err = utils.GetNullableUintPtrSlice[uint64](input, "PtrArrayUint64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint8, err = utils.DecodeNullableObjectValue[[]uint8](input, "PtrArrayUint8")
	if err != nil {
		return err
	}
	j.PtrArrayUint8Ptr, err = utils.GetNullableUintPtrSlice[uint8](input, "PtrArrayUint8Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUintPtr, err = utils.GetNullableUintPtrSlice[uint](input, "PtrArrayUintPtr")
	if err != nil {
		return err
	}
	j.RawJSON, err = utils.GetRawJSON(input, "RawJSON")
	if err != nil {
		return err
	}
	j.RawJSONPtr, err = utils.GetNullableRawJSON(input, "RawJSONPtr")
	if err != nil {
		return err
	}
	j.String, err = utils.GetString(input, "String")
	if err != nil {
		return err
	}
	j.StringPtr, err = utils.GetNullableString(input, "StringPtr")
	if err != nil {
		return err
	}
	j.Text, err = utils.DecodeObjectValue[types.Text](input, "Text")
	if err != nil {
		return err
	}
	j.TextPtr, err = utils.DecodeNullableObjectValue[types.Text](input, "TextPtr")
	if err != nil {
		return err
	}
	j.Time, err = utils.GetDateTime(input, "Time")
	if err != nil {
		return err
	}
	j.TimePtr, err = utils.GetNullableDateTime(input, "TimePtr")
	if err != nil {
		return err
	}
	j.URL, err = utils.DecodeObjectValue[scalar.URL](input, "URL")
	if err != nil {
		return err
	}
	j.UUID, err = utils.GetUUID(input, "UUID")
	if err != nil {
		return err
	}
	j.UUIDPtr, err = utils.GetNullableUUID(input, "UUIDPtr")
	if err != nil {
		return err
	}
	j.Uint, err = utils.GetUint[uint](input, "Uint")
	if err != nil {
		return err
	}
	j.Uint16, err = utils.GetUint[uint16](input, "Uint16")
	if err != nil {
		return err
	}
	j.Uint16Ptr, err = utils.GetNullableUint[uint16](input, "Uint16Ptr")
	if err != nil {
		return err
	}
	j.Uint32, err = utils.GetUint[uint32](input, "Uint32")
	if err != nil {
		return err
	}
	j.Uint32Ptr, err = utils.GetNullableUint[uint32](input, "Uint32Ptr")
	if err != nil {
		return err
	}
	j.Uint64, err = utils.GetUint[uint64](input, "Uint64")
	if err != nil {
		return err
	}
	j.Uint64Ptr, err = utils.GetNullableUint[uint64](input, "Uint64Ptr")
	if err != nil {
		return err
	}
	j.Uint8, err = utils.GetUint[uint8](input, "Uint8")
	if err != nil {
		return err
	}
	j.Uint8Ptr, err = utils.GetNullableUint[uint8](input, "Uint8Ptr")
	if err != nil {
		return err
	}
	j.UintPtr, err = utils.GetNullableUint[uint](input, "UintPtr")
	if err != nil {
		return err
	}
	j.JSONEmpty, err = utils.GetArbitraryJSONDefault(input, "any_empty")
	if err != nil {
		return err
	}
	j.ArrayBigIntEmpty, err = utils.DecodeObjectValueDefault[[]scalar.BigInt](input, "array_bigint_empty")
	if err != nil {
		return err
	}
	j.ArrayBigIntPtrEmpty, err = utils.DecodeObjectValueDefault[[]*scalar.BigInt](input, "array_bigint_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayBoolEmpty, err = utils.GetBooleanSliceDefault(input, "array_bool_empty")
	if err != nil {
		return err
	}
	j.ArrayBoolPtrEmpty, err = utils.GetBooleanPtrSliceDefault(input, "array_bool_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayFloat32Empty, err = utils.GetFloatSliceDefault[float32](input, "array_float32_empty")
	if err != nil {
		return err
	}
	j.ArrayFloat32PtrEmpty, err = utils.GetFloatPtrSliceDefault[float32](input, "array_float32_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayFloat64Empty, err = utils.GetFloatSliceDefault[float64](input, "array_float64_empty")
	if err != nil {
		return err
	}
	j.ArrayFloat64PtrEmpty, err = utils.GetFloatPtrSliceDefault[float64](input, "array_float64_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayInt16Empty, err = utils.GetIntSliceDefault[int16](input, "array_int16_empty")
	if err != nil {
		return err
	}
	j.ArrayInt16PtrEmpty, err = utils.GetIntPtrSliceDefault[int16](input, "array_int16_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayInt32Empty, err = utils.GetIntSliceDefault[int32](input, "array_int32_empty")
	if err != nil {
		return err
	}
	j.ArrayInt32PtrEmpty, err = utils.GetIntPtrSliceDefault[int32](input, "array_int32_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayInt64Empty, err = utils.GetIntSliceDefault[int64](input, "array_int64_empty")
	if err != nil {
		return err
	}
	j.ArrayInt64PtrEmpty, err = utils.GetIntPtrSliceDefault[int64](input, "array_int64_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayInt8Empty, err = utils.GetIntSliceDefault[int8](input, "array_int8_empty")
	if err != nil {
		return err
	}
	j.ArrayInt8PtrEmpty, err = utils.GetIntPtrSliceDefault[int8](input, "array_int8_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayIntEmpty, err = utils.GetIntSliceDefault[int](input, "array_int_empty")
	if err != nil {
		return err
	}
	j.ArrayIntPtrEmpty, err = utils.GetIntPtrSliceDefault[int](input, "array_int_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayJSONEmpty, err = utils.GetArbitraryJSONSliceDefault(input, "array_json_empty")
	if err != nil {
		return err
	}
	j.ArrayJSONPtrEmpty, err = utils.GetArbitraryJSONPtrSliceDefault(input, "array_json_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayMapEmpty, err = utils.DecodeObjectValueDefault[[]map[string]any](input, "array_map_empty")
	if err != nil {
		return err
	}
	j.ArrayRawJSONEmpty, err = utils.GetRawJSONSliceDefault(input, "array_raw_json_empty")
	if err != nil {
		return err
	}
	j.ArrayRawJSONPtrEmpty, err = utils.GetRawJSONPtrSliceDefault(input, "array_raw_json_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayStringEmpty, err = utils.GetStringSliceDefault(input, "array_string_empty")
	if err != nil {
		return err
	}
	j.ArrayStringPtrEmpty, err = utils.GetStringPtrSliceDefault(input, "array_string_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayTimeEmpty, err = utils.GetDateTimeSliceDefault(input, "array_time_empty")
	if err != nil {
		return err
	}
	j.ArrayTimePtrEmpty, err = utils.GetDateTimePtrSliceDefault(input, "array_time_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUint16Empty, err = utils.GetUintSliceDefault[uint16](input, "array_uint16_empty")
	if err != nil {
		return err
	}
	j.ArrayUint16PtrEmpty, err = utils.GetUintPtrSliceDefault[uint16](input, "array_uint16_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUint32Empty, err = utils.GetUintSliceDefault[uint32](input, "array_uint32_empty")
	if err != nil {
		return err
	}
	j.ArrayUint32PtrEmpty, err = utils.GetUintPtrSliceDefault[uint32](input, "array_uint32_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUint64Empty, err = utils.GetUintSliceDefault[uint64](input, "array_uint64_empty")
	if err != nil {
		return err
	}
	j.ArrayUint64PtrEmpty, err = utils.GetUintPtrSliceDefault[uint64](input, "array_uint64_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUint8Empty, err = utils.GetUintSliceDefault[uint8](input, "array_uint8_empty")
	if err != nil {
		return err
	}
	j.ArrayUint8PtrEmpty, err = utils.GetUintPtrSliceDefault[uint8](input, "array_uint8_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUintEmpty, err = utils.GetUintSliceDefault[uint](input, "array_uint_empty")
	if err != nil {
		return err
	}
	j.ArrayUintPtrEmpty, err = utils.GetUintPtrSliceDefault[uint](input, "array_uint_ptr_empty")
	if err != nil {
		return err
	}
	j.ArrayUUIDEmpty, err = utils.GetUUIDSliceDefault(input, "array_uuid_empty")
	if err != nil {
		return err
	}
	j.ArrayUUIDPtrEmpty, err = utils.GetUUIDPtrSliceDefault(input, "array_uuid_ptr_empty")
	if err != nil {
		return err
	}
	j.BigIntEmpty, err = utils.DecodeObjectValueDefault[scalar.BigInt](input, "bigint_empty")
	if err != nil {
		return err
	}
	j.BoolEmpty, err = utils.GetBooleanDefault(input, "bool_empty")
	if err != nil {
		return err
	}
	j.CustomScalarEmpty, err = utils.DecodeObjectValueDefault[types.CommentText](input, "custom_scalar_empty")
	if err != nil {
		return err
	}
	j.DateEmpty, err = utils.DecodeObjectValueDefault[scalar.Date](input, "date_empty")
	if err != nil {
		return err
	}
	j.EnumEmpty, err = utils.DecodeObjectValueDefault[types.SomeEnum](input, "enum_empty")
	if err != nil {
		return err
	}
	j.Float32Empty, err = utils.GetFloatDefault[float32](input, "float32_empty")
	if err != nil {
		return err
	}
	j.Float64Empty, err = utils.GetFloatDefault[float64](input, "float64_empty")
	if err != nil {
		return err
	}
	j.Int16Empty, err = utils.GetIntDefault[int16](input, "int16_empty")
	if err != nil {
		return err
	}
	j.Int32Empty, err = utils.GetIntDefault[int32](input, "int32_empty")
	if err != nil {
		return err
	}
	j.Int64Empty, err = utils.GetIntDefault[int64](input, "int64_empty")
	if err != nil {
		return err
	}
	j.Int8Empty, err = utils.GetIntDefault[int8](input, "int8_empty")
	if err != nil {
		return err
	}
	j.IntEmpty, err = utils.GetIntDefault[int](input, "int_empty")
	if err != nil {
		return err
	}
	j.MapEmpty, err = utils.DecodeObjectValueDefault[map[string]any](input, "map_empty")
	if err != nil {
		return err
	}
	j.RawJSONEmpty, err = utils.GetRawJSONDefault(input, "raw_json_empty")
	if err != nil {
		return err
	}
	j.StringEmpty, err = utils.GetStringDefault(input, "string_empty")
	if err != nil {
		return err
	}
	j.TextEmpty, err = utils.DecodeObjectValueDefault[types.Text](input, "text_empty")
	if err != nil {
		return err
	}
	j.TimeEmpty, err = utils.GetDateTimeDefault(input, "time_empty")
	if err != nil {
		return err
	}
	j.Uint16Empty, err = utils.GetUintDefault[uint16](input, "uint16_empty")
	if err != nil {
		return err
	}
	j.Uint32Empty, err = utils.GetUintDefault[uint32](input, "uint32_empty")
	if err != nil {
		return err
	}
	j.Uint64Empty, err = utils.GetUintDefault[uint64](input, "uint64_empty")
	if err != nil {
		return err
	}
	j.Uint8Empty, err = utils.GetUintDefault[uint8](input, "uint8_empty")
	if err != nil {
		return err
	}
	j.UintEmpty, err = utils.GetUintDefault[uint](input, "uint_empty")
	if err != nil {
		return err
	}
	j.URLEmpty, err = utils.DecodeObjectValueDefault[scalar.URL](input, "url_empty")
	if err != nil {
		return err
	}
	j.UUIDEmpty, err = utils.GetUUIDDefault(input, "uuid_empty")
	if err != nil {
		return err
	}
	return nil
}

// ToMap encodes the struct to a value map
func (j GetCustomHeadersInput) ToMap() map[string]any {
	r := make(map[string]any)
	r["id"] = j.ID
	r["num"] = j.Num

	return r
}

// ToMap encodes the struct to a value map
func (j GetTypesArguments) ToMap() map[string]any {
	r := make(map[string]any)
	r["ArrayBigInt"] = j.ArrayBigInt
	r["ArrayBigIntPtr"] = j.ArrayBigIntPtr
	r["ArrayBool"] = j.ArrayBool
	r["ArrayBoolPtr"] = j.ArrayBoolPtr
	r["ArrayFloat32"] = j.ArrayFloat32
	r["ArrayFloat32Ptr"] = j.ArrayFloat32Ptr
	r["ArrayFloat64"] = j.ArrayFloat64
	r["ArrayFloat64Ptr"] = j.ArrayFloat64Ptr
	r["ArrayInt"] = j.ArrayInt
	r["ArrayInt16"] = j.ArrayInt16
	r["ArrayInt16Ptr"] = j.ArrayInt16Ptr
	r["ArrayInt32"] = j.ArrayInt32
	r["ArrayInt32Ptr"] = j.ArrayInt32Ptr
	r["ArrayInt64"] = j.ArrayInt64
	r["ArrayInt64Ptr"] = j.ArrayInt64Ptr
	r["ArrayInt8"] = j.ArrayInt8
	r["ArrayInt8Ptr"] = j.ArrayInt8Ptr
	r["ArrayIntPtr"] = j.ArrayIntPtr
	r["ArrayJSON"] = j.ArrayJSON
	r["ArrayJSONPtr"] = j.ArrayJSONPtr
	r["ArrayMap"] = j.ArrayMap
	r["ArrayMapPtr"] = j.ArrayMapPtr
	j_ArrayObject := make([]any, len(j.ArrayObject))
	for i, j_ArrayObject_v := range j.ArrayObject {
		j_ArrayObject_v_obj := make(map[string]any)
		j_ArrayObject_v_obj["content"] = j_ArrayObject_v.Content
		j_ArrayObject[i] = j_ArrayObject_v_obj
	}
	r["ArrayObject"] = j_ArrayObject
	if j.ArrayObjectPtr != nil {
		j_ArrayObjectPtr := make([]any, len((*j.ArrayObjectPtr)))
		for i, j_ArrayObjectPtr_v := range *j.ArrayObjectPtr {
			j_ArrayObjectPtr_v_obj := make(map[string]any)
			j_ArrayObjectPtr_v_obj["content"] = j_ArrayObjectPtr_v.Content
			j_ArrayObjectPtr[i] = j_ArrayObjectPtr_v_obj
		}
		r["ArrayObjectPtr"] = j_ArrayObjectPtr
	}
	r["ArrayRawJSON"] = j.ArrayRawJSON
	r["ArrayRawJSONPtr"] = j.ArrayRawJSONPtr
	r["ArrayString"] = j.ArrayString
	r["ArrayStringPtr"] = j.ArrayStringPtr
	r["ArrayTime"] = j.ArrayTime
	r["ArrayTimePtr"] = j.ArrayTimePtr
	r["ArrayUUID"] = j.ArrayUUID
	r["ArrayUUIDPtr"] = j.ArrayUUIDPtr
	r["ArrayUint"] = j.ArrayUint
	r["ArrayUint16"] = j.ArrayUint16
	r["ArrayUint16Ptr"] = j.ArrayUint16Ptr
	r["ArrayUint32"] = j.ArrayUint32
	r["ArrayUint32Ptr"] = j.ArrayUint32Ptr
	r["ArrayUint64"] = j.ArrayUint64
	r["ArrayUint64Ptr"] = j.ArrayUint64Ptr
	r["ArrayUint8"] = j.ArrayUint8
	r["ArrayUint8Ptr"] = j.ArrayUint8Ptr
	r["ArrayUintPtr"] = j.ArrayUintPtr
	r["BigInt"] = j.BigInt
	r["BigIntPtr"] = j.BigIntPtr
	r["Bool"] = j.Bool
	r["BoolPtr"] = j.BoolPtr
	r["Bytes"] = j.Bytes
	r["BytesPtr"] = j.BytesPtr
	r["CustomScalar"] = j.CustomScalar
	r["CustomScalarPtr"] = j.CustomScalarPtr
	r["Date"] = j.Date
	r["DatePtr"] = j.DatePtr
	r["Duration"] = j.Duration
	r["Enum"] = j.Enum
	r["EnumPtr"] = j.EnumPtr
	r["Float32"] = j.Float32
	r["Float32Ptr"] = j.Float32Ptr
	r["Float64"] = j.Float64
	r["Float64Ptr"] = j.Float64Ptr
	r["Int"] = j.Int
	r["Int16"] = j.Int16
	r["Int16Ptr"] = j.Int16Ptr
	r["Int32"] = j.Int32
	r["Int32Ptr"] = j.Int32Ptr
	r["Int64"] = j.Int64
	r["Int64Ptr"] = j.Int64Ptr
	r["Int8"] = j.Int8
	r["Int8Ptr"] = j.Int8Ptr
	r["IntPtr"] = j.IntPtr
	r["JSON"] = j.JSON
	r["JSONPtr"] = j.JSONPtr
	r["Map"] = j.Map
	r["MapPtr"] = j.MapPtr
	j_NamedArray := make([]any, len(j.NamedArray))
	for i, j_NamedArray_v := range j.NamedArray {
		j_NamedArray[i] = j_NamedArray_v
	}
	r["NamedArray"] = j_NamedArray
	if j.NamedArrayPtr != nil {
		j_NamedArrayPtr := make([]any, len((*j.NamedArrayPtr)))
		for i, j_NamedArrayPtr_v := range *j.NamedArrayPtr {
			j_NamedArrayPtr[i] = j_NamedArrayPtr_v
		}
		r["NamedArrayPtr"] = j_NamedArrayPtr
	}
	r["NamedObject"] = j.NamedObject
	if j.NamedObjectPtr != nil {
		r["NamedObjectPtr"] = (*j.NamedObjectPtr)
	}
	j_Object_obj := make(map[string]any)
	j_Object_obj["created_at"] = j.Object.CreatedAt
	j_Object_obj["generic_field"] = j.Object.GenericField
	j_Object_obj["id"] = j.Object.ID
	r["Object"] = j_Object_obj
	if j.ObjectPtr != nil {
		j_ObjectPtr__obj := make(map[string]any)
		j_ObjectPtr__obj["Lat"] = (*j.ObjectPtr).Lat
		j_ObjectPtr__obj["Long"] = (*j.ObjectPtr).Long
		r["ObjectPtr"] = j_ObjectPtr__obj
	}
	r["PtrArrayBigInt"] = j.PtrArrayBigInt
	r["PtrArrayBigIntPtr"] = j.PtrArrayBigIntPtr
	r["PtrArrayBool"] = j.PtrArrayBool
	r["PtrArrayBoolPtr"] = j.PtrArrayBoolPtr
	r["PtrArrayFloat32"] = j.PtrArrayFloat32
	r["PtrArrayFloat32Ptr"] = j.PtrArrayFloat32Ptr
	r["PtrArrayFloat64"] = j.PtrArrayFloat64
	r["PtrArrayFloat64Ptr"] = j.PtrArrayFloat64Ptr
	r["PtrArrayInt"] = j.PtrArrayInt
	r["PtrArrayInt16"] = j.PtrArrayInt16
	r["PtrArrayInt16Ptr"] = j.PtrArrayInt16Ptr
	r["PtrArrayInt32"] = j.PtrArrayInt32
	r["PtrArrayInt32Ptr"] = j.PtrArrayInt32Ptr
	r["PtrArrayInt64"] = j.PtrArrayInt64
	r["PtrArrayInt64Ptr"] = j.PtrArrayInt64Ptr
	r["PtrArrayInt8"] = j.PtrArrayInt8
	r["PtrArrayInt8Ptr"] = j.PtrArrayInt8Ptr
	r["PtrArrayIntPtr"] = j.PtrArrayIntPtr
	r["PtrArrayJSON"] = j.PtrArrayJSON
	r["PtrArrayJSONPtr"] = j.PtrArrayJSONPtr
	r["PtrArrayRawJSON"] = j.PtrArrayRawJSON
	r["PtrArrayRawJSONPtr"] = j.PtrArrayRawJSONPtr
	r["PtrArrayString"] = j.PtrArrayString
	r["PtrArrayStringPtr"] = j.PtrArrayStringPtr
	r["PtrArrayTime"] = j.PtrArrayTime
	r["PtrArrayTimePtr"] = j.PtrArrayTimePtr
	r["PtrArrayUUID"] = j.PtrArrayUUID
	r["PtrArrayUUIDPtr"] = j.PtrArrayUUIDPtr
	r["PtrArrayUint"] = j.PtrArrayUint
	r["PtrArrayUint16"] = j.PtrArrayUint16
	r["PtrArrayUint16Ptr"] = j.PtrArrayUint16Ptr
	r["PtrArrayUint32"] = j.PtrArrayUint32
	r["PtrArrayUint32Ptr"] = j.PtrArrayUint32Ptr
	r["PtrArrayUint64"] = j.PtrArrayUint64
	r["PtrArrayUint64Ptr"] = j.PtrArrayUint64Ptr
	r["PtrArrayUint8"] = j.PtrArrayUint8
	r["PtrArrayUint8Ptr"] = j.PtrArrayUint8Ptr
	r["PtrArrayUintPtr"] = j.PtrArrayUintPtr
	r["RawJSON"] = j.RawJSON
	r["RawJSONPtr"] = j.RawJSONPtr
	r["String"] = j.String
	r["StringPtr"] = j.StringPtr
	r["Text"] = j.Text
	r["TextPtr"] = j.TextPtr
	r["Time"] = j.Time
	r["TimePtr"] = j.TimePtr
	r["URL"] = j.URL
	r["UUID"] = j.UUID
	r["UUIDPtr"] = j.UUIDPtr
	r["Uint"] = j.Uint
	r["Uint16"] = j.Uint16
	r["Uint16Ptr"] = j.Uint16Ptr
	r["Uint32"] = j.Uint32
	r["Uint32Ptr"] = j.Uint32Ptr
	r["Uint64"] = j.Uint64
	r["Uint64Ptr"] = j.Uint64Ptr
	r["Uint8"] = j.Uint8
	r["Uint8Ptr"] = j.Uint8Ptr
	r["UintPtr"] = j.UintPtr
	r["any_empty"] = j.JSONEmpty
	r["array_bigint_empty"] = j.ArrayBigIntEmpty
	r["array_bigint_ptr_empty"] = j.ArrayBigIntPtrEmpty
	r["array_bool_empty"] = j.ArrayBoolEmpty
	r["array_bool_ptr_empty"] = j.ArrayBoolPtrEmpty
	r["array_float32_empty"] = j.ArrayFloat32Empty
	r["array_float32_ptr_empty"] = j.ArrayFloat32PtrEmpty
	r["array_float64_empty"] = j.ArrayFloat64Empty
	r["array_float64_ptr_empty"] = j.ArrayFloat64PtrEmpty
	r["array_int16_empty"] = j.ArrayInt16Empty
	r["array_int16_ptr_empty"] = j.ArrayInt16PtrEmpty
	r["array_int32_empty"] = j.ArrayInt32Empty
	r["array_int32_ptr_empty"] = j.ArrayInt32PtrEmpty
	r["array_int64_empty"] = j.ArrayInt64Empty
	r["array_int64_ptr_empty"] = j.ArrayInt64PtrEmpty
	r["array_int8_empty"] = j.ArrayInt8Empty
	r["array_int8_ptr_empty"] = j.ArrayInt8PtrEmpty
	r["array_int_empty"] = j.ArrayIntEmpty
	r["array_int_ptr_empty"] = j.ArrayIntPtrEmpty
	r["array_json_empty"] = j.ArrayJSONEmpty
	r["array_json_ptr_empty"] = j.ArrayJSONPtrEmpty
	r["array_map_empty"] = j.ArrayMapEmpty
	r["array_raw_json_empty"] = j.ArrayRawJSONEmpty
	r["array_raw_json_ptr_empty"] = j.ArrayRawJSONPtrEmpty
	r["array_string_empty"] = j.ArrayStringEmpty
	r["array_string_ptr_empty"] = j.ArrayStringPtrEmpty
	r["array_time_empty"] = j.ArrayTimeEmpty
	r["array_time_ptr_empty"] = j.ArrayTimePtrEmpty
	r["array_uint16_empty"] = j.ArrayUint16Empty
	r["array_uint16_ptr_empty"] = j.ArrayUint16PtrEmpty
	r["array_uint32_empty"] = j.ArrayUint32Empty
	r["array_uint32_ptr_empty"] = j.ArrayUint32PtrEmpty
	r["array_uint64_empty"] = j.ArrayUint64Empty
	r["array_uint64_ptr_empty"] = j.ArrayUint64PtrEmpty
	r["array_uint8_empty"] = j.ArrayUint8Empty
	r["array_uint8_ptr_empty"] = j.ArrayUint8PtrEmpty
	r["array_uint_empty"] = j.ArrayUintEmpty
	r["array_uint_ptr_empty"] = j.ArrayUintPtrEmpty
	r["array_uuid_empty"] = j.ArrayUUIDEmpty
	r["array_uuid_ptr_empty"] = j.ArrayUUIDPtrEmpty
	r["bigint_empty"] = j.BigIntEmpty
	r["bool_empty"] = j.BoolEmpty
	r["custom_scalar_empty"] = j.CustomScalarEmpty
	r["date_empty"] = j.DateEmpty
	r["enum_empty"] = j.EnumEmpty
	r["float32_empty"] = j.Float32Empty
	r["float64_empty"] = j.Float64Empty
	r["int16_empty"] = j.Int16Empty
	r["int32_empty"] = j.Int32Empty
	r["int64_empty"] = j.Int64Empty
	r["int8_empty"] = j.Int8Empty
	r["int_empty"] = j.IntEmpty
	r["map_empty"] = j.MapEmpty
	r["raw_json_empty"] = j.RawJSONEmpty
	r["string_empty"] = j.StringEmpty
	r["text_empty"] = j.TextEmpty
	r["time_empty"] = j.TimeEmpty
	r["uint16_empty"] = j.Uint16Empty
	r["uint32_empty"] = j.Uint32Empty
	r["uint64_empty"] = j.Uint64Empty
	r["uint8_empty"] = j.Uint8Empty
	r["uint_empty"] = j.UintEmpty
	r["url_empty"] = j.URLEmpty
	r["uuid_empty"] = j.UUIDEmpty

	return r
}
