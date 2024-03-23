package schema

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
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
			assertNoError(t, json.Unmarshal([]byte(tc.raw), &procedure))
			assertDeepEqual(t, tc.expected, procedure)

			var jsonMap map[string]json.RawMessage
			assertNoError(t, json.Unmarshal([]byte(tc.raw), &jsonMap))

			var result ProcedureInfo
			assertNoError(t, result.UnmarshalJSONMap(jsonMap))
			assertDeepEqual(t, procedure, result)
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

func assertNoError(t *testing.T, err error, msgs ...string) {
	if err != nil {
		t.Errorf("%s expected no error, got: %s", strings.Join(msgs, " "), err)
		t.FailNow()
	}
}

func assertDeepEqual(t *testing.T, expected any, reality any, msgs ...string) {
	if !internal.DeepEqual(expected, reality) {
		t.Errorf("%s not equal, expected: %+v got: %+v", strings.Join(msgs, " "), expected, reality)
		t.FailNow()
	}
}
