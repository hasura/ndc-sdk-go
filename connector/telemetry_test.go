package connector

import (
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewTelemetryHeaders(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          http.Header
		AllowedHeaders []string
		Expected       http.Header
	}{
		{
			Name: "basic",
			Input: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer abcdefghijkxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
				"Api-Key":       []string{"abcxyz"},
				"Secret-Key":    []string{"secret-key"},
				"X-Empty":       []string{},
			},
			Expected: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bea*******(65)"},
				"Api-Key":       []string{"******"},
				"Secret-Key":    []string{"s*********"},
			},
		},
		{
			Name: "allowed_list",
			Input: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer abcdefghijkxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
				"Api-Key":       []string{"abcxyz"},
				"Secret-Key":    []string{"secret-key"},
			},
			AllowedHeaders: []string{"Content-Type", "Api-Key"},
			Expected: http.Header{
				"Content-Type": []string{"application/json"},
				"Api-Key":      []string{"******"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.DeepEqual(t, tc.Expected, NewTelemetryHeaders(tc.Input, tc.AllowedHeaders...))
		})
	}
}
