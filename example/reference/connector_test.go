package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/v2/connector"
	"github.com/hasura/ndc-sdk-go/v2/schema"
	"gotest.tools/v3/assert"
)

var baseTestSnapshotURL = fmt.Sprintf(
	"https://raw.githubusercontent.com/hasura/ndc-spec/refs/tags/v%s/ndc-reference/tests",
	schema.NDCVersion,
)

func TestCapabilitySchema(t *testing.T) {
	server := createTestServer(t).BuildTestServer()
	defer server.Close()

	t.Run("capabilities", func(t *testing.T) {
		res, err := http.DefaultClient.Get(server.URL + "/capabilities")
		assert.NilError(t, err)

		expected := fetchSnapshotSample[schema.CapabilitiesResponse](
			t,
			baseTestSnapshotURL+"/capabilities/expected.snap",
		)

		assertHTTPResponse(t, res, http.StatusOK, *expected)
	})

	t.Run("schema", func(t *testing.T) {
		res, err := http.DefaultClient.Get(server.URL + "/schema")
		assert.NilError(t, err)

		expected := fetchSnapshotSample[schema.SchemaResponse](
			t,
			baseTestSnapshotURL+"/schema/expected.snap",
		)

		assertHTTPResponse(t, res, http.StatusOK, *expected)
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
		res, err := httpPostJSON(
			fmt.Sprintf("%s/mutation/explain", server.URL),
			schema.MutationRequest{
				Operations: []schema.MutationOperation{
					{
						Name: "upsert_article",
						Type: schema.MutationOperationProcedure,
					},
				},
				CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
			},
		)

		assert.NilError(t, err)
		assertHTTPResponse(t, res, http.StatusOK, schema.ExplainResponse{
			Details: schema.ExplainResponseDetails{},
		})
	})

	t.Run("mutation_explain_invalid", func(t *testing.T) {
		res, err := httpPostJSON(
			fmt.Sprintf("%s/mutation/explain", server.URL),
			schema.MutationRequest{
				Operations: []schema.MutationOperation{
					{
						Name: "test",
						Type: schema.MutationOperationProcedure,
					},
				},
				CollectionRelationships: make(schema.MutationRequestCollectionRelationships),
			},
		)

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
		name     string
		response []byte
	}{
		{name: "aggregate_function"},
		{name: "aggregate_function_with_nested_field"},
		{name: "authors_with_article_aggregate"},
		{name: "authors_with_articles"},
		{name: "column_count"},
		{name: "get_max_article"},
		{name: "get_all_articles"},
		{name: "get_max_article_id"},
		{name: "group_by"},
		{name: "group_by_with_extraction"},
		{name: "group_by_with_having"},
		{name: "group_by_with_nested_relationship_in_dimension"},
		{name: "group_by_with_order_by"},
		{name: "group_by_with_order_by_dimension"},
		{name: "group_by_with_path_in_dimension"},
		{name: "group_by_with_relationship"},
		{name: "group_by_with_star_count"},
		{name: "group_by_with_where"},
		{name: "named_scopes"},
		{name: "nested_array_select"},
		{name: "nested_array_select_with_limit"},
		{name: "nested_collection_with_aggregates"},
		{name: "nested_collection_with_grouping"},
		{name: "nested_object_select"},
		{name: "nested_object_select_with_relationship"},
		{name: "order_by_aggregate"},
		{name: "order_by_aggregate_function"},
		{name: "order_by_aggregate_nested_relationship"},
		{name: "order_by_aggregate_with_predicate"},
		{name: "order_by_column"},
		{name: "order_by_nested_relationship"},
		{name: "order_by_relationship"},
		{name: "pagination"},
		{name: "predicate_with_array_contains"},
		{name: "predicate_with_array_is_empty"},
		{name: "predicate_with_column_aggregate"},
		{name: "predicate_with_column_aggregate"},
		{name: "predicate_with_eq"},
		{name: "predicate_with_exists"},
		{name: "predicate_with_exists_and_in"},
		{name: "predicate_with_exists_and_relationship"},
		{name: "predicate_with_exists_from_nested_field"},
		{name: "predicate_with_exists_in_nested_collection"},
		{name: "predicate_with_exists_in_nested_scalar_collection"},
		{name: "predicate_with_in"},
		{name: "predicate_with_like"},
		{name: "predicate_with_nested_field_eq"},
		{name: "predicate_with_star_count"},
		{name: "predicate_with_starts_with"},
		{name: "predicate_with_starts_with"},
		{name: "star_count"},
		{name: "table_argument"},
		{name: "table_argument_aggregate"},
		{name: "table_argument_exists"},
		{name: "table_argument_order_by"},
		{name: "table_argument_relationship_1"},
		{name: "table_argument_relationship_2"},
		{name: "table_argument_unrelated_exists"},
		{name: "variables"},
	}

	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestURL := fmt.Sprintf("%s/query/%s/request.json", baseTestSnapshotURL, tc.name)
			req := fetchTestSample(t, requestURL)
			defer req.Body.Close()

			var expected *schema.QueryResponse
			var err error

			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, &expected)
			} else {
				responseURL := fmt.Sprintf("%s/query/%s/expected.snap", baseTestSnapshotURL, tc.name)
				expected = fetchSnapshotSample[schema.QueryResponse](t, responseURL)
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
		name     string
		response []byte
	}{
		{name: "upsert_article"},
		{name: "upsert_article_with_relationship"},
		{name: "delete_articles"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestURL := fmt.Sprintf("%s/mutation/%s/request.json", baseTestSnapshotURL, tc.name)

			req := fetchTestSample(t, requestURL)
			var expected *schema.MutationResponse
			var err error

			if len(tc.response) > 0 {
				err = json.Unmarshal(tc.response, expected)
			} else {
				responseURL := fmt.Sprintf("%s/mutation/%s/expected.snap", baseTestSnapshotURL, tc.name)
				expected = fetchSnapshotSample[schema.MutationResponse](t, responseURL)
			}

			if err != nil {
				t.Errorf("failed to decode expected response: %s", err)
				t.FailNow()
			}

			res, err := http.Post(
				fmt.Sprintf("%s/mutation", server.URL),
				"application/json",
				req.Body,
			)
			if err != nil {
				t.Errorf("expected no error, got %s", err)
				t.FailNow()
			}

			assertHTTPResponse(t, res, http.StatusOK, expected)
		})
	}
}

func createTestServer(t *testing.T) *connector.Server[Configuration, State] {
	t.Helper()

	server, err := connector.NewServer(&Connector{}, &connector.ServerOptions{
		Configuration: "{}",
		InlineConfig:  true,
	}, connector.WithoutRecovery(), connector.WithLogger(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
	}))))
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

func fetchSnapshotSample[R any](t *testing.T, uri string) *R {
	t.Helper()

	res, err := http.Get(uri)
	if err != nil {
		t.Errorf("failed to fetch test sample at %s: %s", uri, err)
		t.FailNow()
	}

	defer res.Body.Close()

	rawBody, err := io.ReadAll(res.Body)
	assert.NilError(t, err)

	linesToSkip := 0
	lines := strings.Split(string(rawBody), "\n")

	for _, line := range lines {
		if line == "{" || line == "[" {
			break
		}

		linesToSkip++
	}

	bodyString := strings.Join(lines[linesToSkip:], "\n")

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

	assert.DeepEqual(t, expectedBody, body)
}
