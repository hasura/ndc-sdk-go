package internal

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
	"gotest.tools/v3/assert"
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
		Name         string
		BasePath     string
		ConnectorDir string
		Directories  []string
		ModuleName   string
		PackageTypes string
		NamingStyle  string
		errorMsg     string
	}{
		{
			Name:        "basic",
			BasePath:    "./testdata/basic",
			ModuleName:  "github.com/hasura/ndc-codegen-test",
			Directories: []string{"functions"},
		},
		{
			Name:        "empty",
			BasePath:    "./testdata/empty",
			ModuleName:  "github.com/hasura/ndc-codegen-empty-test",
			Directories: []string{"functions"},
		},
		{
			Name:         "subdir",
			BasePath:     "./testdata/subdir",
			ConnectorDir: "connector",
			ModuleName:   "github.com/hasura/ndc-codegen-subdir-test",
			Directories:  []string{"connector/functions"},
		},
		{
			Name:        "snake_case",
			BasePath:    "./testdata/snake_case",
			ModuleName:  "github.com/hasura/ndc-codegen-test-snake-case",
			Directories: []string{"functions"},
			NamingStyle: string(StyleSnakeCase),
		},
	}

	rootDir, err := os.Getwd()
	assert.NilError(t, err)
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.NilError(t, os.Chdir(rootDir))
			expectedSchemaBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/schema.json"))
			assert.NilError(t, err)
			connectorContentBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/connector.go.tmpl"))
			assert.NilError(t, err)

			srcDir := path.Join(tc.BasePath, "source")
			assert.NilError(t, os.Chdir(srcDir))
			err = ParseAndGenerateConnector(ConnectorGenerationArguments{
				ConnectorDir: tc.ConnectorDir,
				Directories:  tc.Directories,
				Style:        tc.NamingStyle,
			}, tc.ModuleName)
			if tc.errorMsg != "" {
				assert.ErrorContains(t, err, tc.errorMsg)
				return
			}
			if err != nil {
				panic(err)
			}

			var expectedSchema schema.SchemaResponse
			assert.NilError(t, json.Unmarshal(expectedSchemaBytes, &expectedSchema))

			schemaBytes, err := os.ReadFile(path.Join(tc.ConnectorDir, "schema.generated.json"))
			assert.NilError(t, err)
			var schemaOutput schema.SchemaResponse
			assert.NilError(t, json.Unmarshal(schemaBytes, &schemaOutput))

			assert.DeepEqual(t, expectedSchema.Collections, schemaOutput.Collections)
			assert.DeepEqual(t, expectedSchema.Functions, schemaOutput.Functions)
			assert.DeepEqual(t, expectedSchema.Procedures, schemaOutput.Procedures)
			assert.DeepEqual(t, expectedSchema.ScalarTypes, schemaOutput.ScalarTypes)
			assert.DeepEqual(t, expectedSchema.ObjectTypes, schemaOutput.ObjectTypes)

			connectorBytes, err := os.ReadFile(path.Join(tc.ConnectorDir, "connector.generated.go"))
			assert.NilError(t, err)
			assert.Equal(t, formatTextContent(string(connectorContentBytes)), formatTextContent(string(connectorBytes)))

			// go to the base test directory
			assert.NilError(t, os.Chdir(".."))

			for _, td := range tc.Directories {
				expectedFunctionTypesBytes, err := os.ReadFile(path.Join("expected", "functions.go.tmpl"))
				if err == nil {
					functionTypesBytes, err := os.ReadFile(path.Join("source", td, "types.generated.go"))
					assert.NilError(t, err)
					expected := formatTextContent(string(expectedFunctionTypesBytes))
					reality := formatTextContent(string(functionTypesBytes))
					assert.Equal(t, expected, reality)
				} else if !os.IsNotExist(err) {
					assert.NilError(t, err)
				}
			}
		})
	}

	// go template
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.NilError(t, os.Chdir(rootDir))

			expectedSchemaBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/schema.go.tmpl"))
			if err != nil {
				if os.IsNotExist(err) {
					return
				}
				assert.NilError(t, err)
			}
			connectorContentBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/connector-go.go.tmpl"))
			if err != nil {
				if os.IsNotExist(err) {
					return
				}
				assert.NilError(t, err)
			}

			srcDir := path.Join(tc.BasePath, "source")
			assert.NilError(t, os.Chdir(srcDir))
			err = ParseAndGenerateConnector(ConnectorGenerationArguments{
				ConnectorDir: tc.ConnectorDir,
				Directories:  tc.Directories,
				Style:        tc.NamingStyle,
				SchemaFormat: "go",
			}, tc.ModuleName)
			if tc.errorMsg != "" {
				assert.ErrorContains(t, err, tc.errorMsg)
				return
			}
			assert.NilError(t, err)

			schemaBytes, err := os.ReadFile(path.Join(tc.ConnectorDir, "schema.generated.go"))
			assert.NilError(t, err)
			assert.Equal(t, formatTextContent(string(expectedSchemaBytes)), formatTextContent(string(schemaBytes)))

			connectorBytes, err := os.ReadFile(path.Join(tc.ConnectorDir, "connector.generated.go"))
			assert.NilError(t, err)
			assert.Equal(t, formatTextContent(string(connectorContentBytes)), formatTextContent(string(connectorBytes)))

			// go to the base test directory
			assert.NilError(t, os.Chdir(".."))

			for _, td := range tc.Directories {
				expectedFunctionTypesBytes, err := os.ReadFile(path.Join("expected", "functions.go.tmpl"))
				if err == nil {
					functionTypesBytes, err := os.ReadFile(path.Join("source", td, "types.generated.go"))
					assert.NilError(t, err)
					assert.Equal(t, formatTextContent(string(expectedFunctionTypesBytes)), formatTextContent(string(functionTypesBytes)))
				} else if !os.IsNotExist(err) {
					assert.NilError(t, err)
				}
			}
		})
	}
}
