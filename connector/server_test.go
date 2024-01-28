package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
	"github.com/swaggest/jsonschema-go"
)

type mockRawConfiguration struct {
	Version string `json:"version"`
}

type mockConfiguration struct {
	Version int `json:"version"`
}
type mockState struct{}
type mockConnector struct{}

var mockCapabilities = schema.CapabilitiesResponse{
	Versions: "^0.1.0",
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
			ComparisonOperators: schema.ScalarTypeComparisonOperators{
				"like": schema.ComparisonOperatorDefinition{
					ArgumentType: schema.NewNamedType("String").Serialize(),
				},
			},
		},
		"Int": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"max": schema.AggregateFunctionDefinition{
					ResultType: schema.NewNullableNamedType("Int").Serialize(),
				},
				"min": schema.AggregateFunctionDefinition{
					ResultType: schema.NewNullableNamedType("Int").Serialize(),
				},
			},
			ComparisonOperators: schema.ScalarTypeComparisonOperators{},
		},
	},
	ObjectTypes: schema.SchemaResponseObjectTypes{
		"article": schema.ObjectType{
			Description: schema.ToPtr("An article"),
			Fields: schema.ObjectTypeFields{
				"id": schema.ObjectField{
					Description: schema.ToPtr("The article's primary key"),
					Type:        schema.NewNamedType("Int").Serialize(),
				},
				"title": schema.ObjectField{
					Description: schema.ToPtr("The article's title"),
					Type:        schema.NewNamedType("String").Serialize(),
				},
				"author_id": schema.ObjectField{
					Description: schema.ToPtr("The article's author ID"),
					Type:        schema.NewNamedType("Int").Serialize(),
				},
			},
		},
	},
	Collections: []schema.CollectionInfo{
		{
			Name:        "articles",
			Description: schema.ToPtr("A collection of articles"),
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
			Description: schema.ToPtr("Get the ID of the most recent article"),
			ResultType:  schema.NewNullableNamedType("Int").Serialize(),
			Arguments:   schema.FunctionInfoArguments{},
		},
	},
	Procedures: []schema.ProcedureInfo{
		{
			Name:        "upsert_article",
			Description: schema.ToPtr("Insert or update an article"),
			Arguments: schema.ProcedureInfoArguments{
				"article": schema.ArgumentInfo{
					Description: schema.ToPtr("The article to insert or update"),
					Type:        schema.NewNamedType("article").Serialize(),
				},
			},
			ResultType: schema.NewNullableNamedType("article").Serialize(),
		},
	},
}

func (mc *mockConnector) GetRawConfigurationSchema() *jsonschema.Schema {
	return nil
}
func (mc *mockConnector) MakeEmptyConfiguration() *mockRawConfiguration {
	return &mockRawConfiguration{
		Version: "1",
	}
}

func (mc *mockConnector) UpdateConfiguration(ctx context.Context, rawConfiguration *mockRawConfiguration) (*mockRawConfiguration, error) {
	return &mockRawConfiguration{
		Version: "1",
	}, nil
}
func (mc *mockConnector) ValidateRawConfiguration(rawConfiguration *mockRawConfiguration) (*mockConfiguration, error) {
	return &mockConfiguration{
		Version: 1,
	}, nil
}
func (mc *mockConnector) TryInitState(configuration *mockConfiguration, metrics *TelemetryState) (*mockState, error) {
	return &mockState{}, nil
}

func (mc *mockConnector) HealthCheck(ctx context.Context, configuration *mockConfiguration, state *mockState) error {
	return nil
}

func (mc *mockConnector) GetCapabilities(configuration *mockConfiguration) *schema.CapabilitiesResponse {
	return &mockCapabilities
}

func (mc *mockConnector) GetSchema(configuration *mockConfiguration) (*schema.SchemaResponse, error) {
	return &mockSchema, nil
}
func (mc *mockConnector) Explain(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{}, nil
}
func (mc *mockConnector) Mutation(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	results := []schema.MutationOperationResults{}
	for _, operation := range request.Operations {
		if operation.Name != "upsert_article" {
			return nil, schema.BadRequestError(fmt.Sprintf("operation not found: %s", operation.Name), nil)
		}

		results = append(results, schema.MutationOperationResults{
			AffectedRows: 1,
		})
	}

	return &schema.MutationResponse{
		OperationResults: results,
	}, nil
}
func (mc *mockConnector) Query(ctx context.Context, configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (*schema.QueryResponse, error) {
	if request.Collection != "articles" {
		return nil, schema.BadRequestError(fmt.Sprintf("collection not found: %s", request.Collection), nil)
	}
	return &schema.QueryResponse{
		{
			Aggregates: schema.RowSetAggregates{},
			Rows: []schema.Row{
				map[string]any{
					"id":        1,
					"title":     "Hello world",
					"author_id": 1,
				},
			},
		},
	}, nil
}

// buildTestServer builds the http test server for testing purpose
func buildTestServer(s *Server[mockRawConfiguration, mockConfiguration, mockState]) *httptest.Server {
	s.telemetry.Shutdown(context.Background())
	return httptest.NewServer(s.buildHandler())
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
		t.Errorf("%s: expected status %d, got %d", name, statusCode, res.StatusCode)
		t.FailNow()
	}
}

func assertHTTPResponse[B any](t *testing.T, name string, res *http.Response, statusCode int, expectedBody B) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s: failed to read response body", name)
		t.FailNow()
	}

	if res.StatusCode != statusCode {
		t.Errorf("%s: expected status %d, got %d. Body: %s", name, statusCode, res.StatusCode, string(bodyBytes))
		t.FailNow()
	}

	var body B
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		t.Errorf("%s: failed to decode json body, got error: %s; body: %s", name, err, string(bodyBytes))
		t.FailNow()
	}

	if !internal.DeepEqual(body, expectedBody) {
		t.Errorf("%s: expect body: %+v, got: %+v", name, body, expectedBody)
		t.FailNow()
	}
}

func TestNewServer(t *testing.T) {
	_, err := NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{}, WithLogger(zerolog.Nop()))
	if err == nil {
		t.Error("NewServerEmptyConfig: expect error, got nil")
		t.FailNow()
	}
	if errConfigurationRequired != err {
		t.Errorf("NewServerEmptyConfig: expected error %s, got %s", errConfigurationRequired, err)
		t.FailNow()
	}

	_, err = NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: "/tmp/any-file",
	}, WithLogger(zerolog.Nop()))
	if err == nil {
		t.Errorf("NewServerWithConfigFile: expected error, got nil")
		t.FailNow()
	}
	if !strings.Contains(err.Error(), "Invalid configuration provided: open /tmp/any-file: no such file or directory") {
		t.Errorf("NewServerWithConfigFile: expected file not found error, got %s", err)
		t.FailNow()
	}

	randomFilePath := fmt.Sprintf("%s/test-%d", os.TempDir(), rand.Int())
	if err := os.WriteFile(randomFilePath, []byte{}, 0666); err != nil {
		t.Errorf("NewServerWithEmptyConfigFile: expected no error, got %s", err)
		t.FailNow()
	}

	_, err = NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: randomFilePath,
	}, WithLogger(zerolog.Nop()))
	if err == nil {
		t.Errorf("NewServerWithEmptyConfigFile: expected error, got nil")
		t.FailNow()
	}
	if err != errConfigurationRequired {
		t.Errorf("NewServerWithEmptyConfigFile: expected required file error, got %s", err)
		t.FailNow()
	}

	randomFilePath = fmt.Sprintf("%s/test-%d", os.TempDir(), rand.Int())
	if err := os.WriteFile(randomFilePath, []byte("{"), 0666); err != nil {
		t.Errorf("NewServerWithInvalidConfigFile: expected no error, got %s", err)
		t.FailNow()
	}

	_, err = NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: randomFilePath,
	}, WithLogger(zerolog.Nop()))
	if err == nil {
		t.Errorf("NewServerWithInvalidConfigFile: expected error, got nil")
		t.FailNow()
	}
	if !strings.Contains(err.Error(), "Invalid configuration provided: unexpected end of JSON input") {
		t.Errorf("NewServerWithInvalidConfigFile: expected invalid json error, got %s", err)
		t.FailNow()
	}

	_, err = NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{}, WithoutConfig(), WithLogger(zerolog.Nop()))
	if err != nil {
		t.Errorf("NewServerWithoutConfig: expected no error, got %s", err)
		t.FailNow()
	}
}

func TestServerAuth(t *testing.T) {
	server, err := NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration:      "{}",
		InlineConfig:       true,
		ServiceTokenSecret: "random-secret",
	}, WithLogger(zerolog.Nop()))

	if err != nil {
		t.Errorf("NewServerAuth: expected no error, got %s", err)
		t.FailNow()
	}

	httpServer := buildTestServer(server)
	defer httpServer.Close()

	res, err := http.Get(fmt.Sprintf("%s/schema", httpServer.URL))
	if err != nil {
		t.Errorf("Unauthorized GET /schema: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "Unauthorized GET /schema", res, http.StatusUnauthorized, schema.ErrorResponse{
		Message: "Unauthorized",
		Details: map[string]any{
			"cause": "Bearer token does not match.",
		},
	})

	authRequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/schema", httpServer.URL), nil)
	if err != nil {
		t.Errorf("Authorized GET /schema: expected no error, got %s", err)
		t.FailNow()
	}
	authRequest.Header.Add("Authorization", "Bearer random-secret")

	res, err = http.DefaultClient.Do(authRequest)
	if err != nil {
		t.Errorf("Authorized GET /schema: expected no error, got %s", err)
		t.FailNow()
	}

	assertHTTPResponse(t, "Authorized GET /schema", res, http.StatusOK, mockSchema)
}

func TestServerConnector(t *testing.T) {
	server, err := NewServer[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}, &ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, WithLogger(zerolog.Nop()))

	if err != nil {
		t.Errorf("NewServerAuth: expected no error, got %s", err)
		t.FailNow()
	}

	httpServer := buildTestServer(server)
	defer httpServer.Close()

	res, err := http.Get(fmt.Sprintf("%s/capabilities", httpServer.URL))
	if err != nil {
		t.Errorf("GET /capabilities: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "GET /capabilities", res, http.StatusOK, mockCapabilities)

	res, err = http.Get(fmt.Sprintf("%s/healthz", httpServer.URL))
	if err != nil {
		t.Errorf("GET /healthz: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponseStatus(t, "GET /healthz", res, http.StatusNoContent)

	res, err = http.Get(fmt.Sprintf("%s/metrics", httpServer.URL))
	if err != nil {
		t.Errorf("GET /metrics: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponseStatus(t, "GET /metrics", res, http.StatusOK)

	res, err = httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), schema.QueryRequest{
		Collection:              "articles",
		Arguments:               schema.QueryRequestArguments{},
		CollectionRelationships: schema.QueryRequestCollectionRelationships{},
		Query:                   schema.Query{},
		Variables:               []schema.QueryRequestVariablesElem{},
	})
	if err != nil {
		t.Errorf("POST /query: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "POST /query", res, http.StatusOK, schema.QueryResponse{
		{
			Aggregates: schema.RowSetAggregates{},
			Rows: []schema.Row{
				map[string]any{
					"id":        1,
					"title":     "Hello world",
					"author_id": 1,
				},
			},
		},
	})

	res, err = httpPostJSON(fmt.Sprintf("%s/query", httpServer.URL), schema.QueryRequest{
		Collection:              "test",
		Arguments:               schema.QueryRequestArguments{},
		CollectionRelationships: schema.QueryRequestCollectionRelationships{},
		Query:                   schema.Query{},
		Variables:               []schema.QueryRequestVariablesElem{},
	})
	if err != nil {
		t.Errorf("POST /query: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "POST /query", res, http.StatusBadRequest, schema.ErrorResponse{
		Message: "collection not found: test",
		Details: map[string]any{},
	})

	res, err = httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
		Operations: []schema.MutationOperation{
			{
				Type: "procedure",
				Name: "upsert_article",
			},
		},
		CollectionRelationships: schema.MutationRequestCollectionRelationships{},
	})
	if err != nil {
		t.Errorf("POST /mutation: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "POST /mutation", res, http.StatusOK, schema.MutationResponse{
		OperationResults: []schema.MutationOperationResults{
			{
				AffectedRows: 1,
			},
		},
	})

	res, err = httpPostJSON(fmt.Sprintf("%s/mutation", httpServer.URL), schema.MutationRequest{
		Operations: []schema.MutationOperation{
			{
				Type: "procedure",
				Name: "test",
			},
		},
		CollectionRelationships: schema.MutationRequestCollectionRelationships{},
	})
	if err != nil {
		t.Errorf("POST /mutation: expected no error, got %s", err)
		t.FailNow()
	}
	assertHTTPResponse(t, "POST /mutation", res, http.StatusBadRequest, schema.ErrorResponse{
		Message: "operation not found: test",
		Details: map[string]any{},
	})
}
