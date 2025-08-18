package command

import (
	"path/filepath"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestGenerateNewProject(t *testing.T) {
	tempDir := t.TempDir()
	assert.NilError(t, GenerateNewProject(&NewArguments{
		Name:    "test",
		Module:  "hasura.dev/connector",
		Output:  tempDir,
		Version: "v2.0.0",
	}, true))

	UpdateConnectorSchema(UpdateArguments{
		Path: tempDir,
	}, time.Now())

	assert.NilError(t, GenTestSnapshots(&GenTestSnapshotArguments{
		Schema: filepath.Join(tempDir, "schema.generated.json"),
		Dir:    filepath.Join(tempDir, "testdata"),
	}))
}
