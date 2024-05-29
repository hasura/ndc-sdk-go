// Code generated by github.com/hasura/ndc-sdk-go/codegen, DO NOT EDIT.
package functions

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hasura/ndc-sdk-go/scalar"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"time"
)

var functions_Decoder = utils.NewDecoder()

// FromValue decodes values from map
func (j *GetTypesArguments) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObjectValue(&j.ArrayBigInt, input, "ArrayBigInt")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.ArrayBigIntPtr, input, "ArrayBigIntPtr")
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
	err = functions_Decoder.DecodeObjectValue(&j.ArrayMap, input, "ArrayMap")
	if err != nil {
		return err
	}
	j.ArrayMapPtr = new([]map[string]any)
	err = functions_Decoder.DecodeNullableObjectValue(j.ArrayMapPtr, input, "ArrayMapPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.ArrayObject, input, "ArrayObject")
	if err != nil {
		return err
	}
	j.ArrayObjectPtr = new([]struct {
		Content string "json:\"content\""
	})
	err = functions_Decoder.DecodeNullableObjectValue(j.ArrayObjectPtr, input, "ArrayObjectPtr")
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
	err = functions_Decoder.DecodeObjectValue(&j.BigInt, input, "BigInt")
	if err != nil {
		return err
	}
	j.BigIntPtr = new(scalar.BigInt)
	err = functions_Decoder.DecodeNullableObjectValue(j.BigIntPtr, input, "BigIntPtr")
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
	err = functions_Decoder.DecodeObjectValue(&j.Bytes, input, "Bytes")
	if err != nil {
		return err
	}
	j.BytesPtr = new(scalar.Bytes)
	err = functions_Decoder.DecodeNullableObjectValue(j.BytesPtr, input, "BytesPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.CustomScalar, input, "CustomScalar")
	if err != nil {
		return err
	}
	j.CustomScalarPtr = new(CommentText)
	err = functions_Decoder.DecodeNullableObjectValue(j.CustomScalarPtr, input, "CustomScalarPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.Date, input, "Date")
	if err != nil {
		return err
	}
	j.DatePtr = new(scalar.Date)
	err = functions_Decoder.DecodeNullableObjectValue(j.DatePtr, input, "DatePtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.Enum, input, "Enum")
	if err != nil {
		return err
	}
	j.EnumPtr = new(SomeEnum)
	err = functions_Decoder.DecodeNullableObjectValue(j.EnumPtr, input, "EnumPtr")
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
	err = functions_Decoder.DecodeObjectValue(&j.Map, input, "Map")
	if err != nil {
		return err
	}
	j.MapPtr = new(map[string]any)
	err = functions_Decoder.DecodeNullableObjectValue(j.MapPtr, input, "MapPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.NamedArray, input, "NamedArray")
	if err != nil {
		return err
	}
	j.NamedArrayPtr = new([]Author)
	err = functions_Decoder.DecodeNullableObjectValue(j.NamedArrayPtr, input, "NamedArrayPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.NamedObject, input, "NamedObject")
	if err != nil {
		return err
	}
	j.NamedObjectPtr = new(Author)
	err = functions_Decoder.DecodeNullableObjectValue(j.NamedObjectPtr, input, "NamedObjectPtr")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.Object, input, "Object")
	if err != nil {
		return err
	}
	j.ObjectPtr = new(struct {
		Long int
		Lat  int
	})
	err = functions_Decoder.DecodeNullableObjectValue(j.ObjectPtr, input, "ObjectPtr")
	if err != nil {
		return err
	}
	j.PtrArrayBigInt = new([]scalar.BigInt)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayBigInt, input, "PtrArrayBigInt")
	if err != nil {
		return err
	}
	j.PtrArrayBigIntPtr = new([]*scalar.BigInt)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayBigIntPtr, input, "PtrArrayBigIntPtr")
	if err != nil {
		return err
	}
	j.PtrArrayBool = new([]bool)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayBool, input, "PtrArrayBool")
	if err != nil {
		return err
	}
	j.PtrArrayBoolPtr, err = utils.GetNullableBooleanPtrSlice(input, "PtrArrayBoolPtr")
	if err != nil {
		return err
	}
	j.PtrArrayFloat32 = new([]float32)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayFloat32, input, "PtrArrayFloat32")
	if err != nil {
		return err
	}
	j.PtrArrayFloat32Ptr, err = utils.GetNullableFloatPtrSlice[float32](input, "PtrArrayFloat32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayFloat64 = new([]float64)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayFloat64, input, "PtrArrayFloat64")
	if err != nil {
		return err
	}
	j.PtrArrayFloat64Ptr, err = utils.GetNullableFloatPtrSlice[float64](input, "PtrArrayFloat64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt = new([]int)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayInt, input, "PtrArrayInt")
	if err != nil {
		return err
	}
	j.PtrArrayInt16 = new([]int16)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayInt16, input, "PtrArrayInt16")
	if err != nil {
		return err
	}
	j.PtrArrayInt16Ptr, err = utils.GetNullableIntPtrSlice[int16](input, "PtrArrayInt16Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt32 = new([]int32)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayInt32, input, "PtrArrayInt32")
	if err != nil {
		return err
	}
	j.PtrArrayInt32Ptr, err = utils.GetNullableIntPtrSlice[int32](input, "PtrArrayInt32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt64 = new([]int64)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayInt64, input, "PtrArrayInt64")
	if err != nil {
		return err
	}
	j.PtrArrayInt64Ptr, err = utils.GetNullableIntPtrSlice[int64](input, "PtrArrayInt64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayInt8 = new([]int8)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayInt8, input, "PtrArrayInt8")
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
	j.PtrArrayJSON = new([]any)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayJSON, input, "PtrArrayJSON")
	if err != nil {
		return err
	}
	j.PtrArrayJSONPtr, err = utils.GetNullableArbitraryJSONPtrSlice(input, "PtrArrayJSONPtr")
	if err != nil {
		return err
	}
	j.PtrArrayRawJSON = new([]json.RawMessage)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayRawJSON, input, "PtrArrayRawJSON")
	if err != nil {
		return err
	}
	j.PtrArrayRawJSONPtr, err = utils.GetNullableRawJSONPtrSlice(input, "PtrArrayRawJSONPtr")
	if err != nil {
		return err
	}
	j.PtrArrayString = new([]string)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayString, input, "PtrArrayString")
	if err != nil {
		return err
	}
	j.PtrArrayStringPtr, err = utils.GetNullableStringPtrSlice(input, "PtrArrayStringPtr")
	if err != nil {
		return err
	}
	j.PtrArrayTime = new([]time.Time)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayTime, input, "PtrArrayTime")
	if err != nil {
		return err
	}
	j.PtrArrayTimePtr, err = utils.GetNullableDateTimePtrSlice(input, "PtrArrayTimePtr")
	if err != nil {
		return err
	}
	j.PtrArrayUUID = new([]uuid.UUID)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUUID, input, "PtrArrayUUID")
	if err != nil {
		return err
	}
	j.PtrArrayUUIDPtr, err = utils.GetNullableUUIDPtrSlice(input, "PtrArrayUUIDPtr")
	if err != nil {
		return err
	}
	j.PtrArrayUint = new([]uint)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUint, input, "PtrArrayUint")
	if err != nil {
		return err
	}
	j.PtrArrayUint16 = new([]uint16)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUint16, input, "PtrArrayUint16")
	if err != nil {
		return err
	}
	j.PtrArrayUint16Ptr, err = utils.GetNullableUintPtrSlice[uint16](input, "PtrArrayUint16Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint32 = new([]uint32)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUint32, input, "PtrArrayUint32")
	if err != nil {
		return err
	}
	j.PtrArrayUint32Ptr, err = utils.GetNullableUintPtrSlice[uint32](input, "PtrArrayUint32Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint64 = new([]uint64)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUint64, input, "PtrArrayUint64")
	if err != nil {
		return err
	}
	j.PtrArrayUint64Ptr, err = utils.GetNullableUintPtrSlice[uint64](input, "PtrArrayUint64Ptr")
	if err != nil {
		return err
	}
	j.PtrArrayUint8 = new([]uint8)
	err = functions_Decoder.DecodeNullableObjectValue(j.PtrArrayUint8, input, "PtrArrayUint8")
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
	err = functions_Decoder.DecodeObjectValue(&j.Text, input, "Text")
	if err != nil {
		return err
	}
	j.TextPtr = new(Text)
	err = functions_Decoder.DecodeNullableObjectValue(j.TextPtr, input, "TextPtr")
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
	return nil
}

// FromValue decodes values from map
func (j *GetArticlesArguments) FromValue(input map[string]any) error {
	var err error
	j.Limit, err = utils.GetFloat[float64](input, "Limit")
	if err != nil {
		return err
	}
	return nil
}

// ToMap encodes the struct to a value map
func (j Author) ToMap() map[string]any {
	r := make(map[string]any)
	r["created_at"] = j.CreatedAt
	r["id"] = j.ID
	r["tags"] = j.Tags

	return r
}

// ToMap encodes the struct to a value map
func (j CreateArticleResult) ToMap() map[string]any {
	r := make(map[string]any)
	j_Authors := make([]map[string]any, len(j.Authors))
	for i, j_Authors_v := range j.Authors {
		j_Authors[i] = utils.EncodeMap(j_Authors_v)
	}
	r["authors"] = j_Authors
	r["id"] = j.ID

	return r
}

// ToMap encodes the struct to a value map
func (j CreateAuthorResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["created_at"] = j.CreatedAt
	r["id"] = j.ID
	r["name"] = j.Name

	return r
}

// ToMap encodes the struct to a value map
func (j GetArticlesResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["id"] = j.ID
	r["Name"] = j.Name

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
	j_ArrayObject := make([]map[string]any, len(j.ArrayObject))
	for i, j_ArrayObject_v := range j.ArrayObject {
		j_ArrayObject_v_obj := make(map[string]any)
		j_ArrayObject_v_obj["content"] = j_ArrayObject_v.Content
		j_ArrayObject[i] = j_ArrayObject_v_obj
	}
	r["ArrayObject"] = j_ArrayObject
	if j.ArrayObjectPtr != nil {
		j_ArrayObjectPtr := make([]map[string]any, len((*j.ArrayObjectPtr)))
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
	j_NamedArray := make([]map[string]any, len(j.NamedArray))
	for i, j_NamedArray_v := range j.NamedArray {
		j_NamedArray[i] = utils.EncodeMap(j_NamedArray_v)
	}
	r["NamedArray"] = j_NamedArray
	if j.NamedArrayPtr != nil {
		j_NamedArrayPtr := make([]map[string]any, len((*j.NamedArrayPtr)))
		for i, j_NamedArrayPtr_v := range *j.NamedArrayPtr {
			j_NamedArrayPtr[i] = utils.EncodeMap(j_NamedArrayPtr_v)
		}
		r["NamedArrayPtr"] = j_NamedArrayPtr
	}
	r["NamedObject"] = utils.EncodeMap(j.NamedObject)
	if j.NamedObjectPtr != nil {
		r["NamedObjectPtr"] = utils.EncodeMap((*j.NamedObjectPtr))
	}
	j_Object_obj := make(map[string]any)
	j_Object_obj["created_at"] = j.Object.CreatedAt
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

	return r
}

// ToMap encodes the struct to a value map
func (j HelloResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["foo"] = j.Foo
	r["id"] = j.ID
	r["num"] = j.Num
	r["text"] = j.Text

	return r
}

// ScalarName get the schema name of the scalar
func (j CommentText) ScalarName() string {
	return "CommentString"
}

// ScalarName get the schema name of the scalar
func (j ScalarFoo) ScalarName() string {
	return "Foo"
}

// ScalarName get the schema name of the scalar
func (j SomeEnum) ScalarName() string {
	return "SomeEnum"
}

const (
	SomeEnumFoo SomeEnum = "foo"
	SomeEnumBar SomeEnum = "bar"
)

var enumValues_SomeEnum = []SomeEnum{SomeEnumFoo, SomeEnumBar}

// ParseSomeEnum parses a SomeEnum enum from string
func ParseSomeEnum(input string) (SomeEnum, error) {
	result := SomeEnum(input)
	if !schema.Contains(enumValues_SomeEnum, result) {
		return SomeEnum(""), errors.New("failed to parse SomeEnum, expect one of SomeEnumFoo, SomeEnumBar")
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j SomeEnum) IsValid() bool {
	return schema.Contains(enumValues_SomeEnum, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SomeEnum) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseSomeEnum(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// FromValue decodes the scalar from an unknown value
func (s *SomeEnum) FromValue(value any) error {
	valueStr, err := utils.DecodeNullableString(value)
	if err != nil {
		return err
	}
	if valueStr == nil {
		return nil
	}
	result, err := ParseSomeEnum(*valueStr)
	if err != nil {
		return err
	}

	*s = result
	return nil
}
