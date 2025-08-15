package schema

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestUnmarshalFunctionInfo(t *testing.T) {
	limitDesc := "How many items to return at one time (max 100)"
	functionDesc := "List all pets"
	testCases := []struct {
		name     string
		raw      string
		expected FunctionInfo
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
			expected: FunctionInfo{
				Arguments: FunctionInfoArguments{
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
			var procedure FunctionInfo
			assert.NilError(t, json.Unmarshal([]byte(tc.raw), &procedure))
			assert.DeepEqual(t, tc.expected, procedure)

			var jsonMap map[string]json.RawMessage
			assert.NilError(t, json.Unmarshal([]byte(tc.raw), &jsonMap))

			var result FunctionInfo
			assert.NilError(t, result.UnmarshalJSONMap(jsonMap))
			assert.DeepEqual(t, procedure, result)
		})
	}
}

func TestUnmarshalQueryRequest(t *testing.T) {
	testCases := []struct {
		Name     string
		Raw      string
		Expected QueryRequest
	}{
		{
			Name: "nested_collection_with_grouping",
			Raw: `{
  "collection": "institutions",
  "arguments": {},
  "query": {
    "fields": {
      "id": {
        "type": "column",
        "column": "id"
      },
      "staff_groups": {
        "type": "column",
        "column": "staff",
        "arguments": {
          "limit": {
            "type": "literal",
            "value": null
          }
        },
        "field_path": [],
        "fields": {
          "type": "collection",
          "query": {
            "groups": {
              "dimensions": [
                {
                  "type": "column",
                  "column_name": "last_name",
                  "path": []
                }
              ],
              "aggregates": {
                "count": {
                  "type": "star_count"
                }
              }
            }
          }
        }
      },
      "staff": {
        "type": "column",
        "column": "staff",
        "arguments": {
          "limit": {
            "type": "literal",
            "value": null
          }
        },
        "fields": {
          "type": "array",
          "fields": {
            "type": "object",
            "fields": {
              "last_name": {
                "type": "column",
                "column": "last_name"
              },
              "first_name": {
                "type": "column",
                "column": "first_name"
              }
            }
          }
        }
      }
    }
  },
  "collection_relationships": {}
}`,
			Expected: QueryRequest{
				CollectionRelationships: QueryRequestCollectionRelationships{},
				Collection:              "institutions",
				Arguments:               make(QueryRequestArguments),
				Query: Query{
					Fields: QueryFields{
						"id": NewColumnField("id").Encode(),
						"staff": NewColumnField("staff").
							WithArgument("limit", NewArgumentLiteral(nil)).
							WithNestedField(NewNestedArray(NewNestedObject(map[string]FieldEncoder{
								"first_name": NewColumnField("first_name"),
								"last_name":  NewColumnField("last_name"),
							}))).
							Encode(),
						"staff_groups": NewColumnField("staff").
							WithArgument("limit", NewArgumentLiteral(nil)).
							WithNestedField(NewNestedCollection(Query{
								Groups: &Grouping{
									Aggregates: GroupingAggregates{
										"count": NewAggregateStarCount().Encode(),
									},
									Dimensions: []Dimension{
										NewDimensionColumn("last_name", nil).Wrap(),
									},
								},
							})).
							Encode(),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var query QueryRequest
			assert.NilError(t, json.Unmarshal([]byte(tc.Raw), &query))
			assert.DeepEqual(t, tc.Expected, query)
		})
	}
}
