package schema

import (
	"encoding/json"
	"testing"
)

func TestTypeRepresentation(t *testing.T) {
	t.Run("boolean", func(t *testing.T) {
		typeRep := NewTypeRepresentationBoolean()
		rawType := typeRep.Encode()

		_, err := rawType.AsInteger()
		assertError(t, err, "invalid TypeRepresentation type; expected integer, got boolean")

		anyType, ok := rawType.Interface().(*TypeRepresentationBoolean)
		assertDeepEqual(t, true, ok)
		assertDeepEqual(t, anyType, typeRep)
	})

	t.Run("string", func(t *testing.T) {
		typeRep := NewTypeRepresentationString()
		rawType := typeRep.Encode()

		_, err := rawType.AsInteger()
		assertError(t, err, "invalid TypeRepresentation type; expected integer, got string")

		anyType, ok := rawType.Interface().(*TypeRepresentationString)
		assertDeepEqual(t, true, ok)
		assertDeepEqual(t, anyType, typeRep)
	})

	t.Run("integer", func(t *testing.T) {
		typeRep := NewTypeRepresentationInteger()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got integer")

		anyType, ok := rawType.Interface().(*TypeRepresentationInteger)
		assertDeepEqual(t, true, ok)
		assertDeepEqual(t, anyType, typeRep)
	})

	t.Run("number", func(t *testing.T) {
		typeRep := NewTypeRepresentationNumber()
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got number")

		anyType, ok := rawType.Interface().(*TypeRepresentationNumber)
		assertDeepEqual(t, true, ok)
		assertDeepEqual(t, anyType, typeRep)
	})

	t.Run("enum", func(t *testing.T) {
		typeRep := NewTypeRepresentationEnum([]string{"foo"})
		rawType := typeRep.Encode()

		_, err := rawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got enum")

		anyType, ok := rawType.Interface().(*TypeRepresentationEnum)
		assertDeepEqual(t, true, ok)
		assertDeepEqual(t, anyType, typeRep)

		invalidTypeRep := NewTypeRepresentationEnum([]string{})
		invalidRawType := invalidTypeRep.Encode()

		_, err = invalidRawType.AsString()
		assertError(t, err, "invalid TypeRepresentation type; expected string, got enum")

		_, ok = invalidRawType.Interface().(*TypeRepresentationEnum)
		assertDeepEqual(t, true, ok)

		_, err = (TypeRepresentation{
			"type": TypeRepresentationTypeEnum,
		}).AsEnum()
		assertError(t, err, "TypeRepresentationEnum must have at least 1 item in one_of array")
	})

	t.Run("invalid", func(t *testing.T) {
		rawType := TypeRepresentation{}

		_, err := rawType.AsString()
		assertError(t, err, "type field is required")

		_, ok := rawType.Interface().(*TypeRepresentationEnum)
		assertDeepEqual(t, false, ok)

		assertError(t, json.Unmarshal([]byte(`{"type": "enum"}`), &rawType), "required")
	})
}
