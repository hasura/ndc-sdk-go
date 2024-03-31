package main

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/command"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/stretchr/testify/assert"
)

var (
	newLinesRegexp = regexp.MustCompile(`\n(\s|\t)*\n`)
	tabRegexp      = regexp.MustCompile(`\t`)
	spacesRegexp   = regexp.MustCompile(`\n\s+`)
)

func formatTextContent(input string) string {
	return strings.Trim(spacesRegexp.ReplaceAllString(tabRegexp.ReplaceAllString(newLinesRegexp.ReplaceAllString(input, "\n"), "  "), "\n"), "\n")
}

func TestConnectorGeneration(t *testing.T) {

	testCases := []struct {
		Name       string
		BasePath   string
		ModuleName string
	}{
		{
			Name:       "basic",
			BasePath:   "testdata/basic",
			ModuleName: "github.com/hasura/ndc-codegen-test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			expectedSchemaBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/schema.json"))
			assert.NoError(t, err)
			connectorContentBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/connector.go.tmpl"))
			assert.NoError(t, err)
			expectedFunctionTypesBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/functions.go.tmpl"))
			assert.NoError(t, err)

			srcDir := path.Join(tc.BasePath, "source")
			assert.NoError(t, os.Chdir(srcDir))

			assert.NoError(t, parseAndGenerateConnector(&GenerateArguments{
				Path:        srcDir,
				Directories: []string{"functions"},
			}, tc.ModuleName))

			var expectedSchema schema.SchemaResponse
			assert.NoError(t, json.Unmarshal(expectedSchemaBytes, &expectedSchema))

			schemaBytes, err := os.ReadFile("schema.generated.json")
			assert.NoError(t, err)
			var schemaOutput schema.SchemaResponse
			assert.NoError(t, json.Unmarshal(schemaBytes, &schemaOutput))

			assert.Equal(t, expectedSchema.Collections, schemaOutput.Collections)
			assert.Equal(t, expectedSchema.Functions, schemaOutput.Functions)
			assert.Equal(t, expectedSchema.Procedures, schemaOutput.Procedures)
			assert.Equal(t, expectedSchema.ScalarTypes, schemaOutput.ScalarTypes)
			assert.Equal(t, expectedSchema.ObjectTypes, schemaOutput.ObjectTypes)

			connectorBytes, err := os.ReadFile("connector.generated.go")
			assert.NoError(t, err)
			assert.Equal(t, formatTextContent(string(connectorContentBytes)), formatTextContent(string(connectorBytes)))

			functionTypesBytes, err := os.ReadFile("functions/types.generated.go")
			assert.NoError(t, err)
			assert.Equal(t, formatTextContent(string(expectedFunctionTypesBytes)), formatTextContent(string(functionTypesBytes)))

			// generate test cases
			assert.NoError(t, command.GenTestSnapshots(&command.GenTestSnapshotArguments{
				Dir:    "testdata",
				Schema: "schema.generated.json",
			}))
		})
	}
}
