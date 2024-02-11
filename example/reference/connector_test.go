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

func createTestServer(t *testing.T) *connector.Server[RawConfiguration, Configuration, State] {
	server, err := connector.NewServer[RawConfiguration, Configuration, State](&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, connector.WithLogger(zerolog.Nop()))

	if err != nil {
		t.Errorf("NewServer: expected no error, got %s", err)
		t.FailNow()
	}

	return server
}

func fetchTestSample[R any](t *testing.T, uri string) *http.Response {
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

	if !internal.DeepEqual(body, expectedBody) {
		expectedBytes, _ := json.Marshal(expectedBody)
		t.Errorf("\nexpect: %+v\ngot		: %+v", string(expectedBytes), string(bodyBytes))
		t.FailNow()
	}
}

func TestQuery(t *testing.T) {
	server := createTestServer(t).BuildTestServer()

	testCases := []struct {
		name        string
		requestURL  string
		responseURL string
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
			name:        "get_all_articles",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_all_articles/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_all_articles/expected.json",
		},
		{
			name:        "get_max_article_id",
			requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article_id/request.json",
			responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/get_max_article_id/expected.json",
		},
		// {
		// 	name:        "nested_array_select",
		// 	requestURL:  "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_array_select/request.json",
		// 	responseURL: "https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/query/nested_array_select/expected.json",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := fetchTestSample[schema.QueryRequest](t, tc.requestURL)
			expectedRes := fetchTestSample[schema.QueryResponse](t, tc.responseURL)

			res, err := http.Post(fmt.Sprintf("%s/query", server.URL), "application/json", req.Body)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				t.FailNow()
			}

			var expected schema.QueryResponse
			if err := json.NewDecoder(expectedRes.Body).Decode(&expected); err != nil {
				t.Errorf("failed to decode expected response: %s", err)
				t.FailNow()
			}

			assertHTTPResponse[schema.QueryResponse](t, res, http.StatusOK, expected)
		})
	}
}
