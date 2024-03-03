package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/stretchr/testify/assert"
)

var (
	newLinesRegexp = regexp.MustCompile(`\n(\s|\t)*\n`)
	tabRegexp      = regexp.MustCompile(`\t`)
)

func formatTextContent(input string) string {
	return tabRegexp.ReplaceAllString(newLinesRegexp.ReplaceAllString(input, "\n"), "  ")
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
			ModuleName: "hasura.dev/connector",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			schemaBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/schema.json"))
			assert.NoError(t, err)
			connectorContentBytes, err := os.ReadFile(path.Join(tc.BasePath, "expected/connector.go.tmpl"))
			assert.NoError(t, err)

			fset := token.NewFileSet()
			rawSchema := NewRawConnectorSchema()
			schemaParser := &SchemaParser{
				fset:  fset,
				files: make(map[string]*ast.File),
			}

			err = filepath.Walk(path.Join(tc.BasePath, "source"), func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}

				contentBytes, err := os.ReadFile(filePath)
				if err != nil {
					return err
				}

				f, err := parser.ParseFile(fset, filePath, contentBytes, parser.ParseComments)
				if err != nil {
					t.Errorf("failed to parse src: %s", err)
					t.FailNow()
				}
				schemaParser.files[filePath] = f

				return nil
			})
			assert.NoError(t, err)
			assert.NoError(t, schemaParser.checkAndParseRawSchemaFromAstFiles(rawSchema))

			schemaOutput := rawSchema.Schema()

			var schema schema.SchemaResponse
			assert.NoError(t, json.Unmarshal(schemaBytes, &schema))

			assert.Equal(t, schema.Collections, schemaOutput.Collections)
			assert.Equal(t, schema.Functions, schemaOutput.Functions)
			assert.Equal(t, schema.Procedures, schemaOutput.Procedures)
			assert.Equal(t, schema.ScalarTypes, schemaOutput.ScalarTypes)
			assert.Equal(t, schema.ObjectTypes, schemaOutput.ObjectTypes)

			cg := NewConnectorGenerator(".", "hasura.dev/connector", rawSchema)
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			assert.NoError(t, cg.genConnectorCodeFromTemplate(w))
			w.Flush()
			outputText := formatTextContent(string(buf.String()))
			assert.Equal(t, formatTextContent(string(connectorContentBytes)), outputText)

			//
			assert.NoError(t, cg.genFunctionArgumentConstructors())
			assert.NoError(t, cg.genObjectMethods())
			assert.NoError(t, cg.genCustomScalarMethods())
			for name, builder := range cg.typeBuilders {
				fnContent, err := os.ReadFile(path.Join(tc.BasePath, "expected", fmt.Sprintf("%s.go.tmpl", name)))
				assert.NoError(t, err)
				assert.Equal(t, formatTextContent(string(fnContent)), formatTextContent(builder.String()), name)
			}
		})
	}
}
