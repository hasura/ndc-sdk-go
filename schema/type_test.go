package schema

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestTypeStringer(t *testing.T) {
	testCases := []struct {
		Input    Type
		Expected string
	}{
		{
			Input:    Type{},
			Expected: "zero_type",
		},
		{
			Input:    Type{"type": "invalid"},
			Expected: "failed to parse TypeEnum, expect one of [named nullable array predicate], got invalid",
		},
		{
			Input:    NewNullableNamedType("Foo").Encode(),
			Expected: "Nullable<Foo>",
		},
		{
			Input:    NewNullableType(NewArrayType(NewNamedType("Baz"))).Encode(),
			Expected: "Nullable<Array<Baz>>",
		},
		{
			Input:    NewArrayType(NewNamedType("Bar")).Encode(),
			Expected: "Array<Bar>",
		},
		{
			Input:    NewPredicateType("Object").Encode(),
			Expected: "Predicate<Object>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Expected, func(t *testing.T) {
			assert.DeepEqual(t, tc.Expected, tc.Input.String())
		})
	}
}

func TestGetUnderlyingNamedType(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    Type
		Expected *NamedType
	}{
		{
			Name: "null",
		},
		{
			Name:     "nullable",
			Input:    NewNullableNamedType("Foo").Encode(),
			Expected: NewNamedType("Foo"),
		},
		{
			Name:     "nested_nullable",
			Input:    NewNullableType(NewNullableNamedType("Baz")).Encode(),
			Expected: NewNamedType("Baz"),
		},
		{
			Name:     "nested_nullable_array",
			Input:    NewArrayType(NewNamedType("Bar")).Encode(),
			Expected: NewNamedType("Bar"),
		},
		{
			Name:     "predicate",
			Input:    NewPredicateType("Object").Encode(),
			Expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.DeepEqual(t, tc.Expected, GetUnderlyingNamedType(tc.Input))
		})
	}
}

func TestUnwrapNullableType(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    Type
		Expected Type
	}{
		{
			Name: "null",
		},
		{
			Name:     "named",
			Input:    NewNullableNamedType("Foo").Encode(),
			Expected: NewNamedType("Foo").Encode(),
		},
		{
			Name:     "nested_nullable",
			Input:    NewNullableType(NewNullableNamedType("Baz")).Encode(),
			Expected: NewNamedType("Baz").Encode(),
		},
		{
			Name:     "nested_array_nullable",
			Input:    NewArrayType(NewNullableNamedType("Bar")).Encode(),
			Expected: NewArrayType(NewNullableNamedType("Bar")).Encode(),
		},
		{
			Name:     "nested_nullable_array",
			Input:    NewNullableType(NewArrayType(NewNamedType("Bar"))).Encode(),
			Expected: NewArrayType(NewNamedType("Bar")).Encode(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.DeepEqual(t, tc.Expected, UnwrapNullableType(tc.Input))
		})
	}
}
