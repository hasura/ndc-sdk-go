package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/basic/source
var basicSource embed.FS

//go:embed testdata/basic/schema.json
var basicSchemaBytes []byte

//go:embed testdata/basic/expected/connector.go.tmpl
var basicExpectedContent string

func TestConnectorGeneration(t *testing.T) {
	trimNewLinesRegexp := regexp.MustCompile(`\n\t*\n`)

	testCases := []struct {
		Name      string
		Src       embed.FS
		Schema    []byte
		Generated string
	}{
		{
			Name:      "basic",
			Src:       basicSource,
			Schema:    basicSchemaBytes,
			Generated: basicExpectedContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fset := token.NewFileSet()
			rawSchema := NewRawConnectorSchema()
			schemaParser := &SchemaParser{
				fset:  fset,
				files: make(map[string]*ast.File),
			}

			err := fs.WalkDir(tc.Src, ".", func(filePath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}

				contentBytes, err := tc.Src.ReadFile(filePath)
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
			assert.NoError(t, json.Unmarshal(basicSchemaBytes, &schema))

			assert.Equal(t, schema.Collections, schemaOutput.Collections)
			assert.Equal(t, schema.Functions, schemaOutput.Functions)
			assert.Equal(t, schema.Procedures, schemaOutput.Procedures)
			assert.Equal(t, schema.ScalarTypes, schemaOutput.ScalarTypes)
			assert.Equal(t, schema.ObjectTypes, schemaOutput.ObjectTypes)

			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			assert.NoError(t, genConnectorCodeFromTemplate(w, "hasura.dev/connector", rawSchema))
			w.Flush()
			outputText := trimNewLinesRegexp.ReplaceAllString(string(buf.String()), "\n")
			assert.Equal(t, basicExpectedContent, outputText)
		})
	}
}
