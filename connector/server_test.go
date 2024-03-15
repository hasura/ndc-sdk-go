package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/rs/zerolog"
)

type mockConfiguration struct {
	Version int `json:"version"`
}
type mockState struct{}
type mockConnector struct{}

var mockCapabilities = schema.CapabilitiesResponse{
	Version: "^0.1.0",
	Capabilities: schema.Capabilities{
		Query: schema.QueryCapabilities{
			Aggregates: schema.LeafCapability{},
			Variables:  schema.LeafCapability{},
		},
		Relationships: schema.RelationshipCapabilities{
			OrderByAggregate:    schema.LeafCapability{},
			RelationComparisons: schema.LeafCapability{},
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
		},
		"Int": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"max": schema.AggregateFunctionDefinition{
					ResultType: schema.NewNullableNamedType("Int").Encode(),
				},
				"min": schema.AggregateFunctionDefinition{
					ResultType: schema.NewNullableNamedType("Int").Encode(),
				},
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{},
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
		},
	},
	Collections: []schema.CollectionInfo{
		{
			Name:        "articles",
			Description: utils.ToPtr("A collection of articles"),
			ForeignKeys: schema.CollectionInfoForeignKeys{},
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

func (mc *mockConnector) ParseConfiguration(ctx context.Context, configurationDir string) (*mockConfiguration, error) {
	return &mockConfiguration{
		Version: 1,
	}, nil
}
func (mc *mockConnector) TryInitState(ctx context.Context, configuration *mockConfiguration, metrics *TelemetryState) (*mockState, error) {
	return &mockState{}, nil
}

func (mc *mockConnector) HealthCheck(ctx context.Context, configuration *mockConfiguration, state *mockState) error {
	return nil
}

func (mc *mockConnector) GetCapabilities(configuration *mockConfiguration) schema.CapabilitiesResponseMarshaler {
	return mockCapabilities
}

func (mc *mockConnector) GetSchema(ctx context.Context, configuration *mockConfiguration, state *mockState) (schema.SchemaResponseMarshaler, error) {
	return mockSchema, nil
}
func (mc *mockConnector) QueryExplain(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *mockConnector) MutationExplain(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.MutationRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}
func (mc *mockConnector) Mutation(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	results := []schema.MutationOperationResults{}
	for _, operation := range request.Operations {
		if operation.Name != "upsert_article" {
			return nil, schema.BadRequestError(fmt.Sprintf("operation not found: %s", operation.Name), nil)
		}

		results = append(results, schema.NewProcedureResult(nil).Encode())
	}

	return &schema.MutationResponse{
		OperationResults: results,
	}, nil
}
func (mc *mockConnector) Query(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (schema.QueryResponse, error) {
	if request.Collection != "articles" {
		return nil, schema.BadRequestError(fmt.Sprintf("collection not found: %s", request.Collection), nil)
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

func httpPostJSON(url string, body any) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
}

func assertHTTPResponseStatus(t *testing.T, name string, res *http.Response, statusCode int) {
	if res.StatusCode != statusCode {
		t.Errorf("\n%s: expected status %d, got %d", name, statusCode, res.StatusCode)
		t.FailNow()
	}
}

func assertHTTPResponse[B any](t *testing.T, res *http.Response, statusCode int, expectedBody B) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error("failed to read response body")
		t.FailNow()
	}

	if res.StatusCode != statusCode {
		t.Errorf("expected status %d, got %d. Body: %s", statusCode, res.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	var body B
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		t.Errorf("failed to decode json body, got error: %s; body: %s", err, string(bodyBytes))
		t.FailNow()
	}

	if !internal.DeepEqual(body, expectedBody) {
		expectedBytes, _ := json.Marshal(expectedBody)
		t.Errorf("\nexpect: %+v\ngot: %+v", string(expectedBytes), string(bodyBytes))
		t.FailNow()
	}
}

func TestNewServer(t *testing.T) {

	t.Run("start server", func(t *testing.T) {
		s, err := NewServer[mockConfiguration, mockState](&mockConnector{}, &ServerOptions{}, WithLogger(zerolog.Nop()))
		if err != nil {
			t.Errorf("NewServerWithoutConfig: expected no error, got %s", err)
			t.FailNow()
		}

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
	server, err := NewServer[mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration:      "{}",
		InlineConfig:       true,
		ServiceTokenSecret: "random-secret",
	}, WithLogger(zerolog.Nop()))

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
		authRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/schema", httpServer.URL), nil)
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
	server, err := NewServer[mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, WithLogger(zerolog.Nop()))

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
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponseStatus(t, "GET /health", res, http.StatusOK)
	})

	t.Run("GET /metrics", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/metrics", httpServer.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("\n%s: expected 404 got status %d", "/metrics", res.StatusCode)
			t.FailNow()
		}
	})

	t.Run("POST /query", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), schema.QueryRequest{
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
				Aggregates: schema.RowSetAggregates{},
				Rows: []map[string]any{
					{
						"id":        1,
						"title":     "Hello world",
						"author_id": 1,
					},
				},
			},
		})
	})

	t.Run("POST /query - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), "")
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
		res, err := httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), schema.QueryRequest{
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

	t.Run("POST /mutation", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
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
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), "")
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
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
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

	t.Run("POST /query/explain", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", httpServer.URL), schema.QueryRequest{
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
		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("POST /query/explain - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", httpServer.URL), schema.QueryRequest{})
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

	t.Run("POST /mutation/explain", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", httpServer.URL), schema.MutationRequest{
			Operations:              []schema.MutationOperation{},
			CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
		})
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("POST /mutation/explain - json decode failure", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", httpServer.URL), schema.MutationRequest{})
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
}

func TestConnectorWithPrometheusEnabled(t *testing.T) {

	server, err := NewServer[mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
		OTLPConfig: OTLPConfig{
			MetricsExporter: string(otelMetricsExporterPrometheus),
		},
	}, WithLogger(zerolog.Nop()))

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
