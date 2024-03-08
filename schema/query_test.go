package schema

import (
	"encoding/json"
	"testing"
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
			assertNoError(t, json.Unmarshal([]byte(tc.raw), &procedure))
			assertDeepEqual(t, tc.expected, procedure)

			var jsonMap map[string]json.RawMessage
			assertNoError(t, json.Unmarshal([]byte(tc.raw), &jsonMap))

			var result FunctionInfo
			assertNoError(t, result.UnmarshalJSONMap(jsonMap))
			assertDeepEqual(t, procedure, result)
		})
	}
}
