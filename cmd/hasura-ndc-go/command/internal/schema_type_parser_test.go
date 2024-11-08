package internal

import (
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseTypeFromString(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected Type
	}{
		{
			Input: "[]*github.com/hasura/ndc-codegen-example/types.Author",
			Expected: NewArrayType(NewNullableType(NewNamedType("Author", &TypeInfo{
				Name:        "Author",
				PackageName: "types",
				PackagePath: "github.com/hasura/ndc-codegen-example/types",
			}))),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.FormatInt(int64(i), 32), func(t *testing.T) {
			result, err := parseTypeFromString(tc.Input)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Expected, result)
		})
	}
}

func TestParseTypeParameters(t *testing.T) {
	testCases := []struct {
		Input    string
		Info     TypeInfo
		Expected TypeInfo
	}{
		{
			Input: "github.com/hasura/ndc-codegen-example/types.CustomHeadersResult[github.com/hasura/ndc-codegen-example/types.Author]",
			Info: TypeInfo{
				PackagePath: "github.com/hasura/ndc-codegen-example/types",
				Name:        "CustomHeadersResult",
			},
			Expected: TypeInfo{
				Name:        "CustomHeadersResult",
				PackagePath: "github.com/hasura/ndc-codegen-example/types",
				SchemaName:  "_Author",
				TypeParameters: []Type{
					NewNamedType("Author", &TypeInfo{
						Name:        "Author",
						PackageName: "types",
						PackagePath: "github.com/hasura/ndc-codegen-example/types",
					}),
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.FormatInt(int64(i), 32), func(t *testing.T) {
			result := tc.Info
			err := parseTypeParameters(&result, tc.Input)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.Expected, result)
		})
	}
}
