package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
)

const test_SpecVersion = "v0.1.4"

func createTestServer(t *testing.T) *connector.Server[Configuration, State] {
	server, err := connector.NewServer[Configuration, State](&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, connector.WithLogger(zerolog.Nop()), connector.WithoutRecovery())

	if err != nil {
		t.Errorf("NewServer: expected no error, got %s", err)
		t.FailNow()
	}

	return server
}

func fetchTestSample(t *testing.T, uri string) *http.Response {
	res, err := http.Get(uri)
	if err != nil {
		t.Errorf("failed to fetch test sample at %s: %s", uri, err)
		t.FailNow()
	}

	return res
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

	if !internal.DeepEqual(expectedBody, body) {
		expectedBytes, _ := json.Marshal(expectedBody)
		t.Errorf("\nexpect: %+v\ngot: %+v", string(expectedBytes), string(bodyBytes))
		t.FailNow()
	}
}

func TestGeneralMethods(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	t.Run("capabilities", func(t *testing.T) {
		expectedBytes, err := os.ReadFile("./testdata/capabilities")
		if err != nil {
			t.Errorf("failed to get expected capabilities: %s", err.Error())
			t.FailNow()
		}

		var expectedResult schema.CapabilitiesResponse
		err = json.Unmarshal(expectedBytes, &expectedResult)
		if err != nil {
			t.Errorf("failed to read expected body: %s", err.Error())
			t.FailNow()
		}

		httpResp, err := http.Get(fmt.Sprintf("%s/capabilities", server.URL))
		if err != nil {
			t.Errorf("failed to fetch capabilities: %s", err.Error())
			t.FailNow()
		}

		assertHTTPResponse(t, httpResp, http.StatusOK, expectedResult)
	})

	t.Run("schema", func(t *testing.T) {
		expectedBytes, err := os.ReadFile("./testdata/schema")
		if err != nil {
			t.Errorf("failed to fetch expected schema: %s", err.Error())
			t.FailNow()
		}

		var expectedSchema schema.SchemaResponse
		err = json.Unmarshal(expectedBytes, &expectedSchema)
		if err != nil {
			t.Errorf("failed to read expected body: %s", err.Error())
			t.FailNow()
		}

		httpResp, err := http.Get(fmt.Sprintf("%s/schema", server.URL))
		if err != nil {
			t.Errorf("failed to fetch schema: %s", err.Error())
			t.FailNow()
		}

		assertHTTPResponse(t, httpResp, http.StatusOK, expectedSchema)
	})

	t.Run("GET /health", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/health", server.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		assertHTTPResponseStatus(t, "GET /health", res, http.StatusOK)
	})

	t.Run("GET /metrics", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/metrics", server.URL))
		if err != nil {
			t.Errorf("expected no error, got %s", err)
			t.FailNow()
		}
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("\n%s: expected 404 got status %d", "/metrics", res.StatusCode)
			t.FailNow()
		}
	})

	t.Run("POST /query/explain", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/query/explain", server.URL), schema.QueryRequest{
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

	t.Run("POST /mutation/explain", func(t *testing.T) {
		res, err := httpPostJSON(fmt.Sprintf("%s/mutation/explain", server.URL), schema.MutationRequest{
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

}

func TestQuery(t *testing.T) {
	server := createTestServer(t).BuildTestServer()

	testCases := []struct {
		name        string
		requestURL  string
		responseURL string
		response    []byte
	}{
		{
			name:        "aggregate_function",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/aggregate_function/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/aggregate_function/expected.json", test_SpecVersion),
		},
		{
			name:        "authors_with_article_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_article_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_article_aggregate/expected.json", test_SpecVersion),
		},
		{
			name:        "authors_with_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/authors_with_articles/expected.json", test_SpecVersion),
		},
		{
			name:        "column_count",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/column_count/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/column_count/expected.json", test_SpecVersion),
		},
		{
			name:        "get_max_article",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article/expected.json", test_SpecVersion),
		},
		{
			name:        "get_all_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_all_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_all_articles/expected.json", test_SpecVersion),
		},
		{
			name:        "get_max_article_id",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article_id/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/get_max_article_id/expected.json", test_SpecVersion),
		},
		{
			name:        "nested_array_select",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_array_select/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_array_select/expected.json", test_SpecVersion),
		},
		{
			name:        "nested_object_select",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_object_select/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/nested_object_select/expected.json", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate/expected.json", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate_function",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_function/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_function/expected.json", test_SpecVersion),
		},
		{
			name:        "order_by_aggregate_with_predicate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_with_predicate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_aggregate_with_predicate/expected.json", test_SpecVersion),
		},
		{
			name:        "order_by_column",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_column/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_column/expected.json", test_SpecVersion),
		},
		{
			name:        "order_by_relationship",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_relationship/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/order_by_relationship/expected.json", test_SpecVersion),
		},
		{
			name:        "pagination",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/pagination/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/pagination/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_array_relationship",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_array_relationship/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_array_relationship/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_eq",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_eq/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_eq/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_exists/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_in",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_in/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_in/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_like",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_like/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_like/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_nondet_in_1",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_nondet_in_1/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_nondet_in_1/expected.json", test_SpecVersion),
		},
		{
			name:        "predicate_with_unrelated_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_unrelated_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_unrelated_exists/expected.json", test_SpecVersion),
		},
		{
			name:       "predicate_with_unrelated_exists_and_relationship",
			requestURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/predicate_with_unrelated_exists_and_relationship/request.json", test_SpecVersion),
			response:   []byte(`[{"rows":[{"author_if_has_functional_articles":{},"title":"The Next 700 Programming Languages"},{"author_if_has_functional_articles":{"rows":[{"articles":{"rows":[{"title":"Why Functional Programming Matters"},{"title":"The Design And Implementation Of Programming Languages"}]},"first_name":"John","last_name":"Hughes"}]},"title":"Why Functional Programming Matters"},{"author_if_has_functional_articles":{"rows":[{"articles":{"rows":[{"title":"Why Functional Programming Matters"},{"title":"The Design And Implementation Of Programming Languages"}]},"first_name":"John","last_name":"Hughes"}]},"title":"The Design And Implementation Of Programming Languages"}]}]`),
		},
		{
			name:        "star_count",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/star_count/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/star_count/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_aggregate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_aggregate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_aggregate/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_exists/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_order_by",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_order_by/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_order_by/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_predicate",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_predicate/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_predicate/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_relationship_1",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_1/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_1/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_relationship_2",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_2/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_relationship_2/expected.json", test_SpecVersion),
		},
		{
			name:        "table_argument_unrelated_exists",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_unrelated_exists/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/table_argument_unrelated_exists/expected.json", test_SpecVersion),
		},
		{
			name:        "variables",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/variables/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/query/variables/expected.json", test_SpecVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := fetchTestSample(t, tc.requestURL)
			var expected schema.QueryResponse
			var err error
			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, &expected)
			} else {
				expectedRes := fetchTestSample(t, tc.responseURL)
				err = json.NewDecoder(expectedRes.Body).Decode(&expected)
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

			assertHTTPResponse[schema.QueryResponse](t, res, http.StatusOK, expected)
		})
	}
}

func TestMutation(t *testing.T) {
	server := createTestServer(t).BuildTestServer()

	testCases := []struct {
		name        string
		requestURL  string
		responseURL string
		response    []byte
	}{
		{
			name:        "upsert_article",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article/expected.json", test_SpecVersion),
		},
		{
			name:        "upsert_article_with_relationship",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article_with_relationship/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/upsert_article_with_relationship/expected.json", test_SpecVersion),
		},
		{
			name:        "delete_articles",
			requestURL:  fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/delete_articles/request.json", test_SpecVersion),
			responseURL: fmt.Sprintf("https://raw.githubusercontent.com/hasura/ndc-spec/%s/ndc-reference/tests/mutation/delete_articles/expected.json", test_SpecVersion),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := fetchTestSample(t, tc.requestURL)
			var expected schema.MutationResponse
			var err error
			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, &expected)
			} else {
				expectedRes := fetchTestSample(t, tc.responseURL)
				err = json.NewDecoder(expectedRes.Body).Decode(&expected)
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
