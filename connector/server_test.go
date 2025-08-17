package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/hasura/ndc-sdk-go/utils/compression"
	"gotest.tools/v3/assert"
)

type mockConfiguration struct {
	Version int `json:"version"`
}
type (
	mockState     struct{}
	mockConnector struct{}
)

var mockCapabilities = schema.CapabilitiesResponse{
	Version: schema.NDCVersion,
	Capabilities: schema.Capabilities{
		Query: schema.QueryCapabilities{
			Aggregates: &schema.AggregateCapabilities{},
			Variables:  &schema.LeafCapability{},
		},
		Relationships: &schema.RelationshipCapabilities{
			OrderByAggregate:    &schema.LeafCapability{},
			RelationComparisons: &schema.LeafCapability{},
		},
	},
}

var mockSchema = schema.SchemaResponse{
	ScalarTypes: schema.SchemaResponseScalarTypes{
		"String": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"like": schema.NewComparisonOperatorCustom(schema.NewNamedType("String")).Encode(),
			},
			Representation:      schema.NewTypeRepresentationString().Encode(),
			ExtractionFunctions: schema.ScalarTypeExtractionFunctions{},
		},
		"Int": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"max": schema.NewAggregateFunctionDefinitionMax().Encode(),
				"min": schema.NewAggregateFunctionDefinitionMin().Encode(),
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
			Representation:      schema.NewTypeRepresentationInt32().Encode(),
			ExtractionFunctions: schema.ScalarTypeExtractionFunctions{},
		},
	},
	ObjectTypes: schema.SchemaResponseObjectTypes{
		"article": schema.ObjectType{
			Description: utils.ToPtr("An article"),
			Fields: schema.ObjectTypeFields{
				"id": schema.ObjectField{
					Description: utils.ToPtr("The article's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"title": schema.ObjectField{
					Description: utils.ToPtr("The article's title"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"author_id": schema.ObjectField{
					Description: utils.ToPtr("The article's author ID"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
			},
			ForeignKeys: schema.ObjectTypeForeignKeys{},
		},
	},
	Collections: []schema.CollectionInfo{
		{
			Name:        "articles",
			Type:        "articles",
			Description: utils.ToPtr("A collection of articles"),
			Arguments:   schema.CollectionInfoArguments{},
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"ArticleByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
	},
	Functions: []schema.FunctionInfo{
		{
			Name:        "latest_article_id",
			Description: utils.ToPtr("Get the ID of the most recent article"),
			ResultType:  schema.NewNullableNamedType("Int").Encode(),
			Arguments:   schema.FunctionInfoArguments{},
		},
	},
	Procedures: []schema.ProcedureInfo{
		{
			Name:        "upsert_article",
			Description: utils.ToPtr("Insert or update an article"),
			Arguments: schema.ProcedureInfoArguments{
				"article": schema.ArgumentInfo{
					Description: utils.ToPtr("The article to insert or update"),
					Type:        schema.NewNamedType("article").Encode(),
				},
			},
			ResultType: schema.NewNullableNamedType("article").Encode(),
		},
	},
}

func (mc *mockConnector) ParseConfiguration(
	ctx context.Context,
	configurationDir string,
) (*mockConfiguration, error) {
	return &mockConfiguration{
		Version: 1,
	}, nil
}

func (mc *mockConnector) TryInitState(
	ctx context.Context,
	configuration *mockConfiguration,
	metrics *TelemetryState,
) (*mockState, error) {
	return &mockState{}, nil
}

func (mc *mockConnector) HealthCheck(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
) error {
	return nil
}

func (mc *mockConnector) GetCapabilities(
	configuration *mockConfiguration,
) schema.CapabilitiesResponseMarshaler {
	return mockCapabilities
}

func (mc *mockConnector) GetSchema(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
) (schema.SchemaResponseMarshaler, error) {
	return mockSchema, nil
}

func (mc *mockConnector) QueryExplain(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
	request *schema.QueryRequest,
) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *mockConnector) MutationExplain(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
	request *schema.MutationRequest,
) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *mockConnector) Mutation(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
	request *schema.MutationRequest,
) (*schema.MutationResponse, error) {
	results := []schema.MutationOperationResults{}
	for _, operation := range request.Operations {
		if operation.Name != "upsert_article" {
			return nil, schema.BadRequestError(
				fmt.Sprintf("operation not found: %s", operation.Name),
				nil,
			)
		}

		results = append(results, schema.NewProcedureResult(nil).Encode())
	}

	return &schema.MutationResponse{
		OperationResults: results,
	}, nil
}

func (mc *mockConnector) Query(
	ctx context.Context,
	configuration *mockConfiguration,
	state *mockState,
	request *schema.QueryRequest,
) (schema.QueryResponse, error) {
	if request.Collection != "articles" {
		return nil, schema.BadRequestError(
			fmt.Sprintf("collection not found: %s", request.Collection),
			nil,
		)
	}

	return schema.QueryResponse{
		{
			Aggregates: schema.RowSetAggregates{},
			Rows: []map[string]any{
				{
					"id":        1,
					"title":     "Hello world",
					"author_id": 1,
				},
			},
		},
	}, nil
}

func httpPostJSON(url string, headers map[string]string, body any) (*http.Response, error) {
	return httpPostJSONWithNDCVersion(url, schema.NDCVersion, headers, body)
}

func httpPostJSONWithNDCVersion(url string, ndcVersion string, headers map[string]string, body any) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(schema.XHasuraNDCVersion, ndcVersion)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func assertHTTPResponse[B any](t *testing.T, res *http.Response, statusCode int, expectedBody B) {
	t.Helper()

	defer res.Body.Close()

	var bodyBytes []byte
	var err error
	contentEnc := res.Header.Get(contentEncodingHeader)

	if contentEnc != "" {
		db, err := compression.DefaultCompressor.Decompress(res.Body, contentEnc)
		assert.NilError(t, err)

		bodyBytes, err = io.ReadAll(db)
	} else {
		bodyBytes, err = io.ReadAll(res.Body)
	}

	if err != nil {
		t.Errorf("failed to read response body: %s", err)
		t.FailNow()
	}

	if res.StatusCode != statusCode {
		t.Errorf(
			"expected status %d, got %d. Body: %s",
			statusCode,
			res.StatusCode,
			string(bodyBytes),
		)
		t.FailNow()
	}

	var body B
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		t.Errorf("failed to decode json body, got error: %s; body: %s", err, string(bodyBytes))
		t.FailNow()
	}

	assert.DeepEqual(t, body, expectedBody)
}

func TestNewServer(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		s, err := NewServer(&mockConnector{}, &ServerOptions{
			HTTPServerConfig: HTTPServerConfig{
				ServerMaxBodyMegabytes: 1,
			},
		})
		if err != nil {
			t.Errorf("NewServerWithoutConfig: expected no error, got %s", err)
			t.FailNow()
		}

		// test max size validation
		collectionName := strings.Repeat("xxxxxxxxxxxxxxxx", 1024*1024/16)
		mockBody, err := json.Marshal(schema.QueryRequest{
			Collection: collectionName,
		})
		assert.NilError(t, err)

		w := httptest.NewRecorder()
		r := &http.Request{
			Method: http.MethodPost,
			Body:   io.NopCloser(bytes.NewBuffer(mockBody)),
		}
		s.Query(w, r)

		assertHTTPResponse(t, w.Result(), http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "request body size exceeded 1 MB(s)",
			Details: map[string]any{},
		})

		go func() {
			if err := s.ListenAndServe(18080); err != nil {
				t.Errorf("error happened when running http server: %s", err)
			}
		}()

		time.Sleep(2 * time.Second)
		s.stop()
	})
}

func TestServerAuth(t *testing.T) {
	server, err := NewServer(&mockConnector{}, &ServerOptions{
		Configuration:      "{}",
		InlineConfig:       true,
		ServiceTokenSecret: "random-secret",
		HTTPServerConfig: HTTPServerConfig{
			ServerMaxBodyMegabytes: 1,
		},
	})
	if err != nil {
		t.Errorf("NewServerAuth: expected no error, got %s", err)
		t.FailNow()
	}

	httpServer := server.BuildTestServer()
	defer httpServer.Close()

	t.Run("Unauthorized GET /schema", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/schema", httpServer.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnauthorized, schema.ErrorResponse{
			Message: "Unauthorized",
			Details: map[string]any{
				"cause": "Bearer token does not match.",
			},
		})
	})

	t.Run("Authorized GET /schema", func(t *testing.T) {
		authRequest, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("%s/schema", httpServer.URL),
			nil,
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		authRequest.Header.Add("Authorization", "Bearer random-secret")

		res, err := http.DefaultClient.Do(authRequest)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, mockSchema)
	})
}

func TestServerConnector(t *testing.T) {
	server, err := NewServer(&mockConnector{}, &ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
		HTTPServerConfig: HTTPServerConfig{
			ServerMaxBodyMegabytes: 1,
		},
	})
	if err != nil {
		t.Errorf("NewServer: expected no error, got %s", err)
		t.FailNow()
	}

	httpServer := server.BuildTestServer()
	defer httpServer.Close()

	t.Run("GET /capabilities", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/capabilities", httpServer.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, mockCapabilities)
	})

	t.Run("GET /health", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/health", httpServer.URL))
		assert.NilError(t, err)
		defer res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("GET /metrics", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/metrics", httpServer.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("\n%s: expected 404 got status %d", "/metrics", res.StatusCode)
			t.FailNow()
		}
	})

	t.Run("POST /query", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), nil, schema.QueryRequest{
			Collection:              "articles",
			Arguments:               schema.QueryRequestArguments{},
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
			Query:                   schema.Query{},
			Variables:               []schema.QueryRequestVariablesElem{},
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, schema.QueryResponse{
			{
				Rows: []map[string]any{
					{
						"id":        float64(1),
						"title":     "Hello world",
						"author_id": float64(1),
					},
				},
			},
		})
	})

	t.Run("POST /query - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), nil, "")
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": "json: cannot unmarshal string into Go value of type map[string]interface {}",
			},
		})
	})

	t.Run("POST /query - collection not found", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), nil, schema.QueryRequest{
			Collection:              "test",
			Arguments:               schema.QueryRequestArguments{},
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
			Query:                   schema.Query{},
			Variables:               []schema.QueryRequestVariablesElem{},
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "collection not found: test",
			Details: map[string]any{},
		})
	})

	t.Run("POST_query_invalid_ndc_version", func(t *testing.T) {
		res, err := httpPostJSONWithNDCVersion(
			fmt.Sprintf("%s/query", httpServer.URL),
			"unknown",
			nil,
			schema.QueryRequest{
				Collection:              "test",
				Arguments:               schema.QueryRequestArguments{},
				CollectionRelationships: schema.QueryRequestCollectionRelationships{},
				Query:                   schema.Query{},
				Variables:               []schema.QueryRequestVariablesElem{},
			},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "Invalid X-Hasura-NDC-Version header, expected a semver version string, got: unknown", Details: map[string]any{},
		})
	})

	t.Run("POST_query_max_body_size", func(t *testing.T) {
		collectionName := strings.Repeat("xxxxxxxxxxxxxxxx", 1024*1024/16)

		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), nil, schema.QueryRequest{
			Collection:              collectionName,
			Arguments:               schema.QueryRequestArguments{},
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
			Query:                   schema.Query{},
			Variables:               []schema.QueryRequestVariablesElem{},
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "request body size exceeded 1 MB(s)",
			Details: map[string]any{},
		})
	})

	t.Run("POST /mutation", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), nil, schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Type: "procedure",
					Name: "upsert_article",
				},
			},
			CollectionRelationships: schema.MutationRequestCollectionRelationships{},
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, schema.MutationResponse{
			OperationResults: []schema.MutationOperationResults{
				schema.NewProcedureResult(nil).Encode(),
			},
		})
	})

	t.Run("POST /mutation - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), nil, "")
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": "json: cannot unmarshal string into Go value of type map[string]interface {}",
			},
		})
	})

	t.Run("POST /mutation - operation not found", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), nil, schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Type: "procedure",
					Name: "test",
				},
			},
			CollectionRelationships: schema.MutationRequestCollectionRelationships{},
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "operation not found: test",
			Details: map[string]any{},
		})
	})

	t.Run("POST_mutation_invalid_ndc_version", func(t *testing.T) {
		res, err := httpPostJSONWithNDCVersion(
			fmt.Sprintf("%s/mutation", httpServer.URL),
			"v0.1.6",
			nil,
			schema.MutationRequest{
				Operations: []schema.MutationOperation{
					{
						Type: "procedure",
						Name: "test",
					},
				},
				CollectionRelationships: schema.MutationRequestCollectionRelationships{},
			},
		)

		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: fmt.Sprintf(
				"NDC version range ^%s does not match implemented version v0.1.6",
				schema.NDCVersion,
			),
			Details: map[string]any{},
		})
	})

	t.Run("POST_mutation_max_body_size", func(t *testing.T) {
		collectionName := strings.Repeat("xxxxxxxxxxxxxxxx", 1024*1024/16)
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), nil, schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Type: "procedure",
					Name: collectionName,
				},
			},
			CollectionRelationships: schema.MutationRequestCollectionRelationships{},
		})

		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "request body size exceeded 1 MB(s)",
			Details: map[string]any{},
		})
	})

	t.Run("POST /query/explain", func(t *testing.T) {
		res, err := httpPostJSON(
			fmt.Sprintf("%s/query/explain", httpServer.URL),
			nil,
			schema.QueryRequest{
				Collection:              "articles",
				Arguments:               schema.QueryRequestArguments{},
				CollectionRelationships: schema.QueryRequestCollectionRelationships{},
				Query:                   schema.Query{},
				Variables:               []schema.QueryRequestVariablesElem{},
			},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	for _, encoding := range []compression.CompressionFormat{
		compression.EncodingDeflate,
		compression.EncodingGzip,
		compression.EncodingZstd,
	} {
		t.Run("POST_query_explain_"+string(encoding), func(t *testing.T) {
			res, err := httpPostJSON(
				fmt.Sprintf("%s/query/explain", httpServer.URL),
				map[string]string{
					acceptEncodingHeader: string(encoding),
				},
				schema.QueryRequest{
					Collection:              "articles",
					Arguments:               schema.QueryRequestArguments{},
					CollectionRelationships: schema.QueryRequestCollectionRelationships{},
					Query:                   schema.Query{},
					Variables:               []schema.QueryRequestVariablesElem{},
				},
			)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				t.FailNow()
			}

			assert.Equal(t, string(encoding), res.Header.Get(contentEncodingHeader))

			assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
				Details: schema.ExplainResponseDetails{},
			})
		})
	}

	t.Run("POST /query/explain - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", httpServer.URL), nil, map[string]any{})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": "field arguments in QueryRequest: required",
			},
		})
	})

	t.Run("POST_query_explain_invalid_ndc_version", func(t *testing.T) {
		res, err := httpPostJSONWithNDCVersion(
			fmt.Sprintf("%s/query/explain", httpServer.URL),
			"unknown",
			nil,
			map[string]any{},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "Invalid X-Hasura-NDC-Version header, expected a semver version string, got: unknown",
			Details: map[string]any{},
		})
	})

	t.Run("POST /mutation/explain", func(t *testing.T) {
		res, err := httpPostJSON(
			fmt.Sprintf("%s/mutation/explain", httpServer.URL),
			nil,
			schema.MutationRequest{
				Operations:              []schema.MutationOperation{},
				CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
			},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("POST /mutation/explain - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(
			fmt.Sprintf("%s/mutation/explain", httpServer.URL),
			nil,
			map[string]any{},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": "field collection_relationships in MutationRequest: required",
			},
		})
	})

	t.Run("POST_mutation_explain_invalid_ndc_version", func(t *testing.T) {
		res, err := httpPostJSONWithNDCVersion(
			fmt.Sprintf("%s/mutation/explain", httpServer.URL),
			"unknown",
			nil,
			map[string]any{},
		)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		assertHTTPResponse(t, res, http.StatusBadRequest, schema.ErrorResponse{
			Message: "Invalid X-Hasura-NDC-Version header, expected a semver version string, got: unknown",
			Details: map[string]any{},
		})
	})
}

func TestConnectorWithPrometheusEnabled(t *testing.T) {
	server, err := NewServer(&mockConnector{}, &ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
		OTLPConfig: OTLPConfig{
			MetricsExporter: string(otelMetricsExporterPrometheus),
		},
	})
	if err != nil {
		t.Errorf("NewServer: expected no error, got %s", err)
		t.FailNow()
	}

	httpServer := server.BuildTestServer()
	defer httpServer.Close()

	t.Run("GET /metrics", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/metrics", httpServer.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}

		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusInternalServerError {
			t.Errorf("\n%s: got status %d", "/metrics", res.StatusCode)
			t.FailNow()
		}
	})
}
