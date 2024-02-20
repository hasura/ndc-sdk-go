package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
)

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

func TestSchema(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	expectedResp, err := http.Get("https://raw.githubusercontent.com/hasura/ndc-spec/ed9254cf16efeabaaa7ad92967fd6734e342d9c4/ndc-reference/tests/schema/expected.json")
	if err != nil {
		t.Errorf("failed to fetch expected schema: %s", err.Error())
		t.FailNow()
	}

	var expectedSchema schema.SchemaResponse
	err = json.NewDecoder(expectedResp.Body).Decode(&expectedSchema)
	if err != nil {
		t.Errorf("failed to read expected body: %s", err.Error())
		t.FailNow()
	}

	httpResp, err := http.Get(fmt.Sprintf("%s/schema", server.URL))
	if err != nil {
		t.Errorf("failed to fetch schema: %s", err.Error())
		t.FailNow()
	}

	assertHTTPResponse[schema.SchemaResponse](t, httpResp, http.StatusOK, expectedSchema)
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
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/aggregate_function/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/aggregate_function/expected.json",
		},
		{
			name:        "authors_with_article_aggregate",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/authors_with_article_aggregate/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/authors_with_article_aggregate/expected.json",
		},
		{
			name:        "authors_with_articles",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/authors_with_articles/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/authors_with_articles/expected.json",
		},
		{
			name:        "column_count",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/column_count/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/column_count/expected.json",
		},
		{
			name:        "get_max_article",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article/expected.json",
		},
		{
			name:        "get_all_articles",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_all_articles/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_all_articles/expected.json",
		},
		{
			name:        "get_max_article_id",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article_id/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article_id/expected.json",
		},
		{
			name:        "nested_array_select",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_array_select/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_array_select/expected.json",
		},
		{
			name:        "nested_object_select",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_object_select/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_object_select/expected.json",
		},
		{
			name:        "order_by_aggregate",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate/expected.json",
		},
		{
			name:        "order_by_aggregate_function",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate_function/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate_function/expected.json",
		},
		{
			name:        "order_by_aggregate_with_predicate",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate_with_predicate/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_aggregate_with_predicate/expected.json",
		},
		{
			name:        "order_by_column",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_column/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_column/expected.json",
		},
		{
			name:        "order_by_relationship",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_relationship/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/order_by_relationship/expected.json",
		},
		{
			name:        "pagination",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/pagination/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/pagination/expected.json",
		},
		{
			name:        "predicate_with_array_relationship",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_array_relationship/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_array_relationship/expected.json",
		},
		{
			name:        "predicate_with_eq",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_eq/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_eq/expected.json",
		},
		{
			name:        "predicate_with_exists",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_exists/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_exists/expected.json",
		},
		{
			name:        "predicate_with_in",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_in/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_in/expected.json",
		},
		{
			name:        "predicate_with_like",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_like/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_like/expected.json",
		},
		{
			name:        "predicate_with_nondet_in_1",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_nondet_in_1/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_nondet_in_1/expected.json",
		},
		{
			name:        "predicate_with_unrelated_exists",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_unrelated_exists/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_unrelated_exists/expected.json",
		},
		{
			name:       "predicate_with_unrelated_exists_and_relationship",
			requestURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/predicate_with_unrelated_exists_and_relationship/request.json",
			response:   []byte(`[{"rows":[{"author_if_has_functional_articles":{},"title":"The Next 700 Programming Languages"},{"author_if_has_functional_articles":{"rows":[{"articles":{"rows":[{"title":"Why Functional Programming Matters"},{"title":"The Design And Implementation Of Programming Languages"}]},"first_name":"John","last_name":"Hughes"}]},"title":"Why Functional Programming Matters"},{"author_if_has_functional_articles":{"rows":[{"articles":{"rows":[{"title":"Why Functional Programming Matters"},{"title":"The Design And Implementation Of Programming Languages"}]},"first_name":"John","last_name":"Hughes"}]},"title":"The Design And Implementation Of Programming Languages"}]}]`),
		},
		{
			name:        "star_count",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/star_count/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/star_count/expected.json",
		},
		{
			name:        "table_argument",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument/expected.json",
		},
		{
			name:        "table_argument_aggregate",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_aggregate/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_aggregate/expected.json",
		},
		{
			name:        "table_argument_exists",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_exists/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_exists/expected.json",
		},
		{
			name:        "table_argument_order_by",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_order_by/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_order_by/expected.json",
		},
		{
			name:        "table_argument_predicate",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_predicate/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_predicate/expected.json",
		},
		{
			name:        "table_argument_relationship_1",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_relationship_1/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_relationship_1/expected.json",
		},
		{
			name:        "table_argument_relationship_2",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_relationship_2/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_relationship_2/expected.json",
		},
		{
			name:        "table_argument_unrelated_exists",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_unrelated_exists/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/table_argument_unrelated_exists/expected.json",
		},
		{
			name:        "variables",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/variables/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/variables/expected.json",
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
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/upsert_article/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/upsert_article/expected.json",
		},
		{
			name:        "upsert_article_with_relationship",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/upsert_article_with_relationship/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/upsert_article_with_relationship/expected.json",
		},
		{
			name:        "delete_articles",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/delete_articles/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/mutation/delete_articles/expected.json",
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

			assertHTTPResponse[schema.MutationResponse](t, res, http.StatusOK, expected)
		})
	}
}
