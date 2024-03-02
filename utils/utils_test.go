package utils

import (
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
)

func TestEvalNestedFields(t *testing.T) {
	testCases := []struct {
		Name      string
		Input     any
		Selection schema.NestedField
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
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{}).Encode(),
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
			Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{
				"id": schema.NewColumnField("id", nil),
			}).Encode(),
			Expected: map[string]any{
				"id": "1",
			},
		},
		// {
		// 	Name: "nested_fields",
		// 	Input: struct {
		// 		ID       string `json:"id"`
		// 		Name     string `json:"name"`
		// 		Articles []struct {
		// 			ID        int       `json:"id"`
		// 			Name      string    `json:"name"`
		// 			CreatedAt time.Time `json:"created_at"`
		// 		} `json:"articles"`
		// 	}{
		// 		ID:   "1",
		// 		Name: "John",
		// 		Articles: []struct {
		// 			ID        int       `json:"id"`
		// 			Name      string    `json:"name"`
		// 			CreatedAt time.Time `json:"created_at"`
		// 		}{
		// 			{
		// 				ID:        1,
		// 				Name:      "Article 1",
		// 				CreatedAt: time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
		// 			},
		// 		},
		// 	},
		// 	Selection: schema.NewNestedObject(map[string]schema.FieldEncoder{
		// 		"id": schema.NewColumnField("id", nil),
		// 		"articles": schema.NewColumnField("articles", schema.NewNestedArray(schema.NewNestedObject(map[string]schema.FieldEncoder{
		// 			"id":         schema.NewColumnField("id", nil),
		// 			"created_at": schema.NewColumnField("created_at", nil),
		// 		}))),
		// 	}).Encode(),
		// 	Expected: map[string]any{
		// 		"id": "1",
		// 		"articles": []map[string]any{
		// 			{
		// 				"id":         1,
		// 				"created_at": time.Date(2020, 01, 01, 00, 0, 0, 0, time.UTC),
		// 			},
		// 		},
		// 	},
		// },
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
