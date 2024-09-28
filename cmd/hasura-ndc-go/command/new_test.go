package command

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateNewProject(t *testing.T) {
	tempDir := t.TempDir()
	assert.NilError(t, GenerateNewProject(&NewArguments{
		Name:    "test",
		Module:  "hasura.dev/connector",
		Output:  tempDir,
		Version: "v1.1.1",
	}, true))
}
