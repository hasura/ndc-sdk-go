package connector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
	"github.com/swaggest/jsonschema-go"
)

// buildTestConfigurationServer builds the http test configuration server for testing purpose
func buildTestConfigurationServer() *httptest.Server {
	s := NewConfigurationServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, WithLogger(zerolog.Nop()))
	return httptest.NewServer(s.buildHandler())
}

func TestConfigurationServer(t *testing.T) {
	httpServer := buildTestConfigurationServer()
	defer httpServer.Close()

	t.Run("GET /healthz", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/healthz", httpServer.URL))
		if err != nil {
			t.Errorf("GET /healthz: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponseStatus(t, "GET /healthz", res, http.StatusNoContent)
	})

	t.Run("GET /", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s", httpServer.URL))
		if err != nil {
			t.Errorf("GET /: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "GET /", res, http.StatusOK, mockRawConfiguration{
			Version: "1",
		})
	})

	t.Run("POST /", func(t *testing.T) {
		res, err := httpPostJSON(httpServer.URL, mockRawConfiguration{
			Version: "1",
		})
		if err != nil {
			t.Errorf("POST /: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "POST /", res, http.StatusOK, mockRawConfiguration{
			Version: "1",
		})
	})

	t.Run("POST / - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s", httpServer.URL), "")
		if err != nil {
			t.Errorf("POST /: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "POST /", res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]string{
				"cause": "json: cannot unmarshal string into Go value of type connector.mockRawConfiguration",
			},
		})
	})

	t.Run("GET /schema", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/schema", httpServer.URL))
		if err != nil {
			t.Errorf("GET /schema: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "GET /schema", res, http.StatusOK, jsonschema.Schema{
			ID: schema.ToPtr("test"),
		})
	})

	t.Run("POST /validate", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/validate", httpServer.URL), mockRawConfiguration{
			Version: "1",
		})
		if err != nil {
			t.Errorf("POST /validate: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "POST /validate", res, http.StatusOK, schema.ValidateResponse{
			Schema:                mockSchema,
			Capabilities:          mockCapabilities,
			ResolvedConfiguration: `{"version":1}`,
		})
	})

	t.Run("POST /validate - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/validate", httpServer.URL), "")
		if err != nil {
			t.Errorf("POST /validate: expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, "POST /validate", res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]string{
				"cause": "json: cannot unmarshal string into Go value of type connector.mockRawConfiguration",
			},
		})
	})

}
