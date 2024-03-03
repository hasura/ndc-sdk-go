package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNewProject(t *testing.T) {
	tempDir := t.TempDir()
	assert.NoError(t, generateNewProject("test", "hasura.dev/connector", tempDir))
}
