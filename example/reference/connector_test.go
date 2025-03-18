package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/ndctest"
	"github.com/hasura/ndc-sdk-go/schema"
	"gotest.tools/v3/assert"
)

const test_SpecVersion = "v0.2.0"

func createTestServer(t *testing.T) *connector.Server[Configuration, State] {
	t.Helper()

	server, err := connector.NewServer(&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, connector.WithoutRecovery())

	if err != nil {
		t.Errorf("NewServer: expected no error, got %s", err)
		t.FailNow()
	}

	return server
}

func fetchTestSample(t *testing.T, uri string) *http.Response {
	t.Helper()

	res, err := http.Get(uri)
	if err != nil {
		t.Errorf("failed to fetch test sample at %s: %s", uri, err)
		t.FailNow()
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("failed to fetch test sample at %s: %s", uri, res.Status)
		t.FailNow()
	}

	return res
}

func fetchResponseSample[R any](t *testing.T, uri string) *R {
	t.Helper()

	res, err := http.Get(uri)
	if err != nil {
		t.Errorf("failed to fetch test sample at %s: %s", uri, err)
		t.FailNow()
	}

	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	assert.NilError(t, err)

	lines := strings.Split(string(rawBody), "\n")
	if len(lines) < 5 {
		return nil
	}

	bodyString := strings.Join(lines[5:], "\n")

	var result R
	assert.NilError(t, json.Unmarshal([]byte(bodyString), &result))

	return &result
}

func httpPostJSON(url string, body any) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
}

func assertHTTPResponse[B any](t *testing.T, res *http.Response, statusCode int, expectedBody B) {
	t.Helper()

	defer res.Body.Close()
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

	assert.DeepEqual(t, expectedBody, body)
}

func TestConnector(t *testing.T) {
	ndctest.TestConnector(t, &Connector{}, ndctest.TestConnectorOptions{
		Configuration: "{}",
		InlineConfig:  true,
		TestDataDir:   "testdata",
	})
}

func TestExplain(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	defer server.Close()

	t.Run("query_explain_success", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", server.URL), schema.QueryRequest{
			Collection:              "articles",
			Arguments:               schema.QueryRequestArguments{},
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
			Query:                   schema.Query{},
			Variables:               []schema.QueryRequestVariablesElem{},
		})
		assert.NilError(t, err)
		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("query_explain_invalid", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", server.URL), schema.QueryRequest{
			Collection:              "test",
			Arguments:               schema.QueryRequestArguments{},
			CollectionRelationships: schema.QueryRequestCollectionRelationships{},
			Query:                   schema.Query{},
			Variables:               []schema.QueryRequestVariablesElem{},
		})
		assert.NilError(t, err)
		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "invalid query name: test",
			Details: map[string]any{},
		})
	})

	t.Run("mutation_explain_success", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", server.URL), schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Name: "upsert_article",
					Type: schema.MutationOperationProcedure,
				},
			},
			CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
		})
		assert.NilError(t, err)
		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("mutation_explain_invalid", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", server.URL), schema.MutationRequest{
			Operations: []schema.MutationOperation{
				{
					Name: "test",
					Type: schema.MutationOperationProcedure,
				},
			},
			CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
		})
		assert.NilError(t, err)
		assertHTTPResponse(t, res, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "invalid mutation name: test",
			Details: map[string]any{},
		})
	})
}

func TestQuery(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	defer server.Close()

	testCases := []struct {
		name        string
		requestURL  string
		responseURL string
		response    []byte
	}{
		{
			name:        "aggregate_function",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/aggregate_function/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/aggregate_function/expected.snap", test_SpecVersion),
		},
		{
			name:        "authors_with_article_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_article_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_article_aggregate/expected.snap", test_SpecVersion),
		},
		{
			name:        "authors_with_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_articles/expected.snap", test_SpecVersion),
		},
		{
			name:        "column_count",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/column_count/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/column_count/expected.snap", test_SpecVersion),
		},
		{
			name:        "get_max_article",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article/expected.snap", test_SpecVersion),
		},
		{
			name:        "get_all_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_all_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_all_articles/expected.snap", test_SpecVersion),
		},
		{
			name:        "get_max_article_id",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article_id/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article_id/expected.snap", test_SpecVersion),
		},
		{
			name:        "nested_array_select",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_array_select/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_array_select/expected.snap", test_SpecVersion),
		},
		{
			name:        "nested_object_select",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_object_select/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_object_select/expected.snap", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate/expected.snap", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate_function",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_function/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_function/expected.snap", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate_with_predicate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_with_predicate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_with_predicate/expected.snap", test_SpecVersion),
		},
		{
			name:        "order_by_column",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_column/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_column/expected.snap", test_SpecVersion),
		},
		{
			name:        "order_by_relationship",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_relationship/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_relationship/expected.snap", test_SpecVersion),
		},
		{
			name:        "pagination",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/pagination/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/pagination/expected.snap", test_SpecVersion),
		},
		{
			name:        "predicate_with_eq",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_eq/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_eq/expected.snap", test_SpecVersion),
		},
		{
			name:        "predicate_with_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_exists/expected.snap", test_SpecVersion),
		},
		{
			name:        "predicate_with_in",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_in/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_in/expected.snap", test_SpecVersion),
		},
		{
			name:        "predicate_with_like",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_like/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_like/expected.snap", test_SpecVersion),
		},
		{
			name:        "star_count",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/star_count/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/star_count/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_aggregate/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_exists/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument_order_by",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_order_by/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_order_by/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument_relationship_2",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_2/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_2/expected.snap", test_SpecVersion),
		},
		{
			name:        "table_argument_unrelated_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_unrelated_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_unrelated_exists/expected.snap", test_SpecVersion),
		},
		{
			name:        "variables",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/variables/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/variables/expected.snap", test_SpecVersion),
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := fetchTestSample(t, tc.requestURL)
			defer req.Body.Close()

			var expected *schema.QueryResponse
			var err error

			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, &expected)
			} else {
				expected = fetchResponseSample[schema.QueryResponse](t, tc.responseURL)
			}

			if err != nil {
				t.Errorf("failed to decode expected response: %s", err)
				t.FailNow()
			}

			res, err := http.Post(fmt.Sprintf("%s/query", server.URL), "application/json", req.Body)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				t.FailNow()
			}

			assertHTTPResponse(t, res, http.StatusOK, expected)
		})
	}
}

func TestMutation(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	defer server.Close()

	testCases := []struct {
		name        string
		requestURL  string
		responseURL string
		response    []byte
	}{
		{
			name:        "upsert_article",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article/expected.snap", test_SpecVersion),
		},
		{
			name:        "upsert_article_with_relationship",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article_with_relationship/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article_with_relationship/expected.snap", test_SpecVersion),
		},
		{
			name:        "delete_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/delete_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/delete_articles/expected.snap", test_SpecVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := fetchTestSample(t, tc.requestURL)
			var expected *schema.MutationResponse
			var err error

			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, expected)
			} else {
				expected = fetchResponseSample[schema.MutationResponse](t, tc.responseURL)
			}

			if err != nil {
				t.Errorf("failed to decode expected response: %s", err)
				t.FailNow()
			}

			res, err := http.Post(fmt.Sprintf("%s/mutation", server.URL), "application/json", req.Body)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				t.FailNow()
			}

			assertHTTPResponse(t, res, http.StatusOK, expected)
		})
	}
}
