package schema

import (
	"testing"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
)

func TestEvalNestedFields(t *testing.T) {
	testCases := []struct {
		Name      string
		Input     any
		Selection NestedField
		Expected  any
	}{
		{
			Name: "no_field",
			Input: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "1",
				Name: "John",
			},
			Selection: NewNestedObject(map[string]FieldEncoder{}).Encode(),
			Expected: map[string]any{
				"id":   "1",
				"name": "John",
			},
		},
		{
			Name: "column_field",
			Input: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "1",
				Name: "John",
			},
			Selection: NewNestedObject(map[string]FieldEncoder{
				"id": NewColumnField("id", nil),
			}).Encode(),
			Expected: map[string]any{
				"id": "1",
			},
		},
		{
			Name: "nested_fields",
			Input: struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Articles []struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					CreatedAt time.Time `json:"created_at"`
				} `json:"articles"`
			}{
				ID:   "1",
				Name: "John",
				Articles: []struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					CreatedAt time.Time `json:"created_at"`
				}{
					{
						ID:        1,
						Name:      "Article 1",
						CreatedAt: time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
					},
				},
			},
			Selection: NewNestedObject(map[string]FieldEncoder{
				"id": NewColumnField("id", nil),
				"articles": NewColumnField("articles", NewNestedArray(NewNestedObject(map[string]FieldEncoder{
					"id":         NewColumnField("id", nil),
					"created_at": NewColumnField("created_at", nil),
				}))),
			}).Encode(),
			Expected: map[string]any{
				"id": "1",
				"articles": []map[string]any{
					{
						"id":         1,
						"created_at": time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := EvalNestedColumnFields(tc.Selection, tc.Input)
			if err != nil {
				t.Errorf("failed to evaluate nested column fields: %s", err)
				t.FailNow()
			}
			if !internal.DeepEqual(tc.Expected, result) {
				t.Errorf("evaluated result does not equal\nexpected: %+v, got: %+v", tc.Expected, result)
				t.FailNow()
			}
		})
	}
}
