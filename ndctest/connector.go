package ndctest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"gotest.tools/v3/assert"
)

// TestConnectorOptions options for the test connector runner
type TestConnectorOptions struct {
	Configuration string
	InlineConfig  bool
	TestDataDir   string
}

// TestConnector the native test runner for the data connector.
// This is a port of [ndc-test]. ndc-test is awesome. However, it doesn't help increase the go coverage ratio.
//
// [ndc-test]: https://github.com/hasura/ndc-spec/tree/main/ndc-test
func TestConnector[Configuration any, State any](t *testing.T, ndc connector.Connector[Configuration, State], options TestConnectorOptions) {
	server, err := connector.NewServer(ndc, &connector.ServerOptions{
		OTLPConfig: connector.OTLPConfig{
			MetricsExporter: "prometheus",
		},
		Configuration: options.Configuration,
		InlineConfig:  options.InlineConfig,
	}, connector.WithoutRecovery())
	assert.NilError(t, err)

	httpServer := server.BuildTestServer()
	defer httpServer.Close()

	if options.TestDataDir == "" {
		options.TestDataDir = "testdata"
	}

	// evaluate capabilities
	res, err := http.Get(fmt.Sprintf("%s/capabilities", httpServer.URL))
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		t.FailNow()
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var capabilities schema.CapabilitiesResponse
	capabilitiesPath := filepath.Join(options.TestDataDir, "capabilities")
	expectedBytes, err := os.ReadFile(capabilitiesPath)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("failed to read file %s: %s", capabilitiesPath, err)
			t.FailNow()
		}
		assert.Equal(t, res.StatusCode, http.StatusOK)
		assert.NilError(t, json.NewDecoder(res.Body).Decode(&capabilities))
		res.Body.Close()
	} else {
		assert.NilError(t, json.Unmarshal(expectedBytes, &capabilities))
		assertResponseJSON(t, res, http.StatusOK, capabilities)
	}

	t.Run("get_health", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/health", httpServer.URL))
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("get_schema", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/schema", httpServer.URL))
		assert.NilError(t, err)
		defer res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var schemaResp schema.SchemaResponse
		schemaPath := filepath.Join(options.TestDataDir, "schema")
		expectedBytes, err := os.ReadFile(schemaPath)
		if err != nil {
			if !os.IsNotExist(err) {
				t.Errorf("failed to read file %s: %s", schemaPath, err)
				t.FailNow()
			}
			assert.Equal(t, res.StatusCode, http.StatusOK)
			assert.NilError(t, json.NewDecoder(res.Body).Decode(&schemaResp))
			return
		}

		assert.NilError(t, json.Unmarshal(expectedBytes, &schemaResp))
		assertResponseJSON(t, res, http.StatusOK, schemaResp)
	})

	t.Run("get_metrics", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/metrics", httpServer.URL))
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("explain_query_failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", httpServer.URL), schema.QueryRequest{})
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Check(t, res.StatusCode >= 400)
	})

	t.Run("explain_mutation_failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", httpServer.URL), schema.MutationRequest{})
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Check(t, res.StatusCode >= 400)
	})

	t.Run("query_failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), schema.QueryRequest{
			Collection: uuid.NewString(),
		})
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Check(t, res.StatusCode >= 400)
	})

	t.Run("mutation_failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Type: schema.MutationOperationProcedure,
					Name: uuid.NewString(),
				},
			},
		})
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Check(t, res.StatusCode >= 400)
	})

	t.Run("mutation_invalid_operation_type", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Type: "invalid_type",
				},
			},
		})
		assert.NilError(t, err)
		_ = res.Body.Close()
		assert.Check(t, res.StatusCode >= 400)
	})

	// replay query tests
	queryDirs, err := os.ReadDir(filepath.Join(options.TestDataDir, "query"))
	if err != nil {
		if os.IsNotExist(err) {
			t.Logf("skip running query snapshots. The query folder does not exist")
		} else {
			t.Errorf("failed to read query snapshots at %s: %s", options.TestDataDir, err)
			t.FailNow()
		}
	}

	queryURL := fmt.Sprintf("%s/query", httpServer.URL)
	for _, dir := range queryDirs {
		t.Run(fmt.Sprintf("query/%s", dir.Name()), func(t *testing.T) {
			snapshotDir := filepath.Join(options.TestDataDir, "query", dir.Name())
			req, expected := readSnapshot(t, snapshotDir)
			if len(req) == 0 {
				return
			}
			res, err := httpPostJSON(queryURL, req)
			assert.NilError(t, err)

			defer func() {
				_ = res.Body.Close()
			}()

			if utils.IsNil(expected) {
				assert.Equal(t, res.StatusCode, http.StatusOK)
				var r schema.QueryResponse
				assert.NilError(t, json.NewDecoder(res.Body).Decode(&r))
				return
			}

			assertResponseJSON(t, res, http.StatusOK, expected)
		})
	}

	// replay mutation tests
	mutationDirs, err := os.ReadDir(filepath.Join(options.TestDataDir, "mutation"))
	if err != nil {
		if os.IsNotExist(err) {
			t.Logf("skip running mutation snapshots. The mutation folder does not exist")
		} else {
			t.Errorf("failed to read mutation snapshots at %s: %s", options.TestDataDir, err)
			t.FailNow()
		}
	}

	mutationURL := fmt.Sprintf("%s/mutation", httpServer.URL)
	for _, dir := range mutationDirs {
		t.Run(fmt.Sprintf("mutation/%s", dir.Name()), func(t *testing.T) {
			snapshotDir := filepath.Join(options.TestDataDir, "mutation", dir.Name())
			req, expected := readSnapshot(t, snapshotDir)
			if len(req) == 0 {
				return
			}
			res, err := httpPostJSON(mutationURL, req)
			assert.NilError(t, err)
			defer func() {
				_ = res.Body.Close()
			}()

			if utils.IsNil(expected) {
				assert.Equal(t, res.StatusCode, http.StatusOK)
				var r schema.MutationResponse
				assert.NilError(t, json.NewDecoder(res.Body).Decode(&r))
				return
			}

			assertResponseJSON(t, res, http.StatusOK, expected)
		})
	}
}
