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
