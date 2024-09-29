package ndctest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func assertResponseJSON[B any](t *testing.T, res *http.Response, statusCode int, expectedBody B) {
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	assert.NilError(t, err, "failed to read response body")
	assert.Equal(t, res.StatusCode, statusCode, "unexpected response status code")

	var body B
	assert.NilError(t, json.Unmarshal(bodyBytes, &body), "failed to decode json body")
	assert.DeepEqual(t, body, expectedBody)
}

func httpPostJSON(url string, body any) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
}

func readSnapshot(t *testing.T, dir string, skipResponseValidation bool) (map[string]any, any) {
	reqPath := filepath.Join(dir, "request.json")
	reqBytes, err := os.ReadFile(reqPath)
	if err != nil {
		// skip non-exist request.json file in the folder
		if os.IsNotExist(err) {
			return nil, nil
		}
		t.Errorf("failed to read the request snapshot at %s: %s", reqPath, err)
		t.FailNow()
	}

	var req map[string]any
	var expected any
	assert.NilError(t, json.Unmarshal(reqBytes, &req))

	if skipResponseValidation {
		return req, nil
	}

	expectedPath := filepath.Join(dir, "expected.json")
	expectedBytes, err := os.ReadFile(expectedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return req, nil
		}
		t.Errorf("failed to read the expected snapshot at %s: %s", expectedPath, err)
		t.FailNow()
	}

	assert.NilError(t, json.Unmarshal(expectedBytes, &expected))

	return req, expected
}
