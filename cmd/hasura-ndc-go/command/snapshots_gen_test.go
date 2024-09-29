package command_test

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command"
	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command/internal"
	"gotest.tools/v3/assert"
)

//go:embed testdata/snapshots/schema.json
var testSchema string

func TestGenTestSnapshots(t *testing.T) {
	tmpDir := t.TempDir()
	mux := http.NewServeMux()
	mux.HandleFunc("/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testSchema))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	assert.NilError(t, command.GenTestSnapshots(&command.GenTestSnapshotArguments{
		Endpoint: server.URL,
		Dir:      tmpDir,
		Depth:    10,
		Strategy: internal.WriteFileStrategyOverride,
	}))
}
