package schema

import (
	"encoding/json"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshalProcedureInfo(t *testing.T) {
	limitDesc := "How many items to return at one time (max 100)"
	functionDesc := "List all pets"
	testCases := []struct {
		name     string
		raw      string
		expected ProcedureInfo
	}{
		{
			name: "success",
			raw: ` {
				"arguments": {
					"limit": {
						"description": "How many items to return at one time (max 100)",
						"type": {
							"type": "nullable",
							"underlying_type": { "name": "Int", "type": "named" }
						}
					}
				},
				"description": "List all pets",
				"name": "listPets",
				"result_type": {
					"element_type": { "name": "Pet", "type": "named" },
					"type": "array"
				}
			}`,
			expected: ProcedureInfo{
				Arguments: ProcedureInfoArguments{
					"limit": ArgumentInfo{
						Description: &limitDesc,
						Type:        NewNullableNamedType("Int").Encode(),
					},
				},
				Description: &functionDesc,
				Name:        "listPets",
				ResultType:  NewArrayType(NewNamedType("Pet")).Encode(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var procedure ProcedureInfo
			assert.NilError(t, json.Unmarshal([]byte(tc.raw), &procedure))
			assert.DeepEqual(t, tc.expected, procedure)

			var jsonMap map[string]json.RawMessage
			assert.NilError(t, json.Unmarshal([]byte(tc.raw), &jsonMap))

			var result ProcedureInfo
			assert.NilError(t, result.UnmarshalJSONMap(jsonMap))
			assert.DeepEqual(t, procedure, result)
		})
	}
}

func TestUnmarshalMutationRequest(t *testing.T) {
	testCases := []struct {
		Name              string
		Raw               string
		Expected          MutationRequest
		ExpectedArguments any
	}{
		{
			Name: "delete_articles",
			Raw: `{
  "operations": [
    {
      "type": "procedure",
      "name": "delete_articles",
      "arguments": {
        "where": {
          "type": "binary_comparison_operator",
          "column": { "type": "column", "name": "author_id" },
          "operator": "eq",
          "value": { "type": "scalar", "value": 1 }
        }
      },
      "fields": {
        "type": "array",
        "fields": {
          "type": "object",
          "fields": {
            "id": { "type": "column", "column": "id" },
            "title": { "type": "column", "column": "title" }
          }
        }
      }
    }
  ],
  "collection_relationships": {}
}`,
			Expected: MutationRequest{
				CollectionRelationships: MutationRequestCollectionRelationships{},
				Operations: []MutationOperation{
					{
						Type: "procedure",
						Name: "delete_articles",
						Fields: NewNestedArray(NewNestedObject(map[string]FieldEncoder{
							"id":    NewColumnField("id"),
							"title": NewColumnField("title"),
						})).Encode(),
					},
				},
			},
			ExpectedArguments: map[string]Expression{
				"where": NewExpressionBinaryComparisonOperator(
					NewComparisonTargetColumn("author_id"),
					"eq",
					NewComparisonValueScalar(float64(1)),
				).Encode(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var mutation MutationRequest
			assert.NilError(t, json.Unmarshal([]byte(tc.Raw), &mutation))
			assert.DeepEqual(t, tc.Expected.CollectionRelationships, mutation.CollectionRelationships)
			assert.DeepEqual(t, tc.Expected.Operations[0].Type, mutation.Operations[0].Type)
			assert.DeepEqual(t, tc.Expected.Operations[0].Name, mutation.Operations[0].Name)
			assert.DeepEqual(t, tc.Expected.Operations[0].Fields, mutation.Operations[0].Fields)

			var arguments map[string]Expression
			assert.NilError(t, json.Unmarshal(mutation.Operations[0].Arguments, &arguments))
			assert.DeepEqual(t, tc.ExpectedArguments, arguments)
		})
	}
}

func assertError(t *testing.T, err error, msg string) {
	if err == nil {
		t.Error("expected error, got: nil")
		t.FailNow()
	}
	if !strings.Contains(err.Error(), msg) {
		t.Errorf("expected error: %s, got: %s", msg, err.Error())
		t.FailNow()
	}
}
