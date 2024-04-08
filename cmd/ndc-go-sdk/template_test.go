package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNewProject(t *testing.T) {
	tempDir := t.TempDir()
	assert.NoError(t, generateNewProject(&NewArguments{
		Name:    "test",
		Module:  "hasura.dev/connector",
		Output:  tempDir,
		Version: "v1.1.0",
	}, true))
}
