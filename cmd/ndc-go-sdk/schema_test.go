package main

import (
	"embed"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"testing"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/basic/source
var basicSource embed.FS

//go:embed testdata/basic/schema.json
var basicSchemaBytes []byte

func TestParseCodesToNdcSchema(t *testing.T) {
	testCases := []struct {
		Name   string
		Src    embed.FS
		Schema []byte
	}{
		{
			Name:   "basic",
			Src:    basicSource,
			Schema: basicSchemaBytes,
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

			if err != nil {
				t.Errorf("failed to read source code: %s", err)
				t.FailNow()
			}

			if err := schemaParser.checkAndParseRawSchemaFromAstFiles(rawSchema); err != nil {
				t.Errorf("failed to parse raw schema: %s", err)
				t.FailNow()
			}

			schemaOutput := rawSchema.Schema()
			var schema schema.SchemaResponse
			if err := json.Unmarshal(basicSchemaBytes, &schema); err != nil {
				t.Errorf("failed to decode expected schema: %s", err)
				t.FailNow()
			}

			assert.Equal(t, schema.Collections, schemaOutput.Collections)
			assert.Equal(t, schema.Functions, schemaOutput.Functions)
			assert.Equal(t, schema.Procedures, schemaOutput.Procedures)
			assert.Equal(t, schema.ScalarTypes, schemaOutput.ScalarTypes)
			assert.Equal(t, schema.ObjectTypes, schemaOutput.ObjectTypes)
		})
	}
}
