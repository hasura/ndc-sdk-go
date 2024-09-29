package schema

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTypeRepresentation(t *testing.T) {
	t.Run("boolean", func(t *testing.T) {
		typeRep := NewTypeRepresentationBoolean()
		rawType := typeRep.Encode()

		_, err := rawType.AsInteger()
		assertError(t, err, "invalid TypeRepresentation type; expected integer, got boolean")

		anyType, ok := rawType.Interface().(*TypeRepresentationBoolean)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("string", func(t *testing.T) {
		typeRep := NewTypeRepresentationString()
		rawType := typeRep.Encode()

		_, err := rawType.AsInteger()
		assertError(t, err, "invalid TypeRepresentation type; expected integer, got string")

		anyType, ok := rawType.Interface().(*TypeRepresentationString)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("integer", func(t *testing.T) {
		typeRep := NewTypeRepresentationInteger()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got integer")

		anyType, ok := rawType.Interface().(*TypeRepresentationInteger)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("number", func(t *testing.T) {
		typeRep := NewTypeRepresentationNumber()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got number")

		anyType, ok := rawType.Interface().(*TypeRepresentationNumber)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("enum", func(t *testing.T) {
		rawBytes := []byte(`{
			"type": "enum",
			"one_of": ["foo"]
		}`)

		rawEmptyBytes := []byte(`{
			"type": "enum",
			"one_of": []
		}`)

		var enumType TypeRepresentation
		assertError(t, json.Unmarshal(rawEmptyBytes, &enumType), "TypeRepresentation requires at least 1 item in one_of field for enum type")

		assert.NilError(t, json.Unmarshal(rawBytes, &enumType))

		typeRep := NewTypeRepresentationEnum([]string{"foo"})
		rawType := typeRep.Encode()
		assert.DeepEqual(t, enumType, rawType)

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got enum")

		anyType, ok := rawType.Interface().(*TypeRepresentationEnum)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)

		invalidTypeRep := NewTypeRepresentationEnum([]string{})
		invalidRawType := invalidTypeRep.Encode()

		_, err = invalidRawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got enum")

		_, ok = invalidRawType.Interface().(*TypeRepresentationEnum)
		assert.DeepEqual(t, true, ok)

		_, err = (TypeRepresentation{
			"type": TypeRepresentationTypeEnum,
		}).AsEnum()
		assertError(t, err, "TypeRepresentationEnum must have at least 1 item in one_of array")
	})

	t.Run("int8", func(t *testing.T) {
		typeRep := NewTypeRepresentationInt8()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got int8")

		anyType, ok := rawType.Interface().(*TypeRepresentationInt8)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("int16", func(t *testing.T) {
		typeRep := NewTypeRepresentationInt16()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got int16")

		anyType, ok := rawType.Interface().(*TypeRepresentationInt16)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("int32", func(t *testing.T) {
		typeRep := NewTypeRepresentationInt32()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got int32")

		anyType, ok := rawType.Interface().(*TypeRepresentationInt32)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("int64", func(t *testing.T) {
		typeRep := NewTypeRepresentationInt64()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got int64")

		anyType, ok := rawType.Interface().(*TypeRepresentationInt64)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})
	t.Run("float32", func(t *testing.T) {
		typeRep := NewTypeRepresentationFloat32()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got float32")

		anyType, ok := rawType.Interface().(*TypeRepresentationFloat32)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("float64", func(t *testing.T) {
		typeRep := NewTypeRepresentationFloat64()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got float64")

		anyType, ok := rawType.Interface().(*TypeRepresentationFloat64)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("big integer", func(t *testing.T) {
		typeRep := NewTypeRepresentationBigInteger()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got biginteger")

		anyType, ok := rawType.Interface().(*TypeRepresentationBigInteger)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("big decimal", func(t *testing.T) {
		typeRep := NewTypeRepresentationBigDecimal()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got bigdecimal")

		anyType, ok := rawType.Interface().(*TypeRepresentationBigDecimal)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("uuid", func(t *testing.T) {
		typeRep := NewTypeRepresentationUUID()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got uuid")

		anyType, ok := rawType.Interface().(*TypeRepresentationUUID)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("date", func(t *testing.T) {
		typeRep := NewTypeRepresentationDate()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got date")

		anyType, ok := rawType.Interface().(*TypeRepresentationDate)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("timestamp", func(t *testing.T) {
		typeRep := NewTypeRepresentationTimestamp()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got timestamp")

		anyType, ok := rawType.Interface().(*TypeRepresentationTimestamp)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("timestamptz", func(t *testing.T) {
		typeRep := NewTypeRepresentationTimestampTZ()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got timestamptz")

		anyType, ok := rawType.Interface().(*TypeRepresentationTimestampTZ)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("geography", func(t *testing.T) {
		typeRep := NewTypeRepresentationGeography()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got geography")

		anyType, ok := rawType.Interface().(*TypeRepresentationGeography)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("geometry", func(t *testing.T) {
		typeRep := NewTypeRepresentationGeometry()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got geometry")

		anyType, ok := rawType.Interface().(*TypeRepresentationGeometry)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("bytes", func(t *testing.T) {
		typeRep := NewTypeRepresentationBytes()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got bytes")

		anyType, ok := rawType.Interface().(*TypeRepresentationBytes)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("json", func(t *testing.T) {
		typeRep := NewTypeRepresentationJSON()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got json")

		anyType, ok := rawType.Interface().(*TypeRepresentationJSON)
		assert.DeepEqual(t, true, ok)
		assert.DeepEqual(t, anyType, typeRep)
	})

	t.Run("invalid", func(t *testing.T) {
		rawType := TypeRepresentation{}

		_, err := rawType.AsString()
		assertError(t, err, "type field is required")

		_, ok := rawType.Interface().(*TypeRepresentationEnum)
		assert.DeepEqual(t, false, ok)

		assertError(t, json.Unmarshal([]byte(`{"type": "enum"}`), &rawType), "required")
	})
}
