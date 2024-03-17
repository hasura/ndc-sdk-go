package main

// var (
// 	newLinesRegexp = regexp.MustCompile(`\n(\s|\t)*\n`)
// 	tabRegexp      = regexp.MustCompile(`\t`)
// 	spacesRegexp   = regexp.MustCompile(`\n\s+`)
// )

// func formatTextContent(input string) string {
// 	return spacesRegexp.ReplaceAllString(tabRegexp.ReplaceAllString(newLinesRegexp.ReplaceAllString(input, "\n"), "  "), "\n")
// }

// func TestConnectorGeneration(t *testing.T) {

// 	testCases := []struct {
// 		Name       string
// 		BasePath   string
// 		ModuleName string
// 	}{
// 		{
// 			Name:       "basic",
// 			BasePath:   "testdata/basic",
// 			ModuleName: "github.com/hasura/ndc-codegen-test",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			_, err := os.ReadFile(path.Join(tc.BasePath, "expected/schema.json"))
// 			assert.NoError(t, err)
// 			_, err = os.ReadFile(path.Join(tc.BasePath, "expected/connector.go.tmpl"))
// 			assert.NoError(t, err)

// 			assert.NoError(t, parseAndGenerateConnector(&GenerateArguments{
// 				Path:        path.Join(tc.BasePath, "source"),
// 				Directories: []string{"functions"},
// 			}, tc.ModuleName))
// 			// schemaOutput := rawSchema.Schema()

// 			// var schema schema.SchemaResponse
// 			// assert.NoError(t, json.Unmarshal(schemaBytes, &schema))

// 			// assert.Equal(t, schema.Collections, schemaOutput.Collections)
// 			// assert.Equal(t, schema.Functions, schemaOutput.Functions)
// 			// assert.Equal(t, schema.Procedures, schemaOutput.Procedures)
// 			// assert.Equal(t, schema.ScalarTypes, schemaOutput.ScalarTypes)
// 			// assert.Equal(t, schema.ObjectTypes, schemaOutput.ObjectTypes)

// 			// assert.Equal(t, formatTextContent(string(connectorContentBytes)), outputText)

// 			//
// 			// assert.NoError(t, cg.genFunctionArgumentConstructors())
// 			// assert.NoError(t, cg.genObjectMethods())
// 			// assert.NoError(t, cg.genCustomScalarMethods())
// 			// for name, builder := range cg.typeBuilders {
// 			// 	fnContent, err := os.ReadFile(path.Join(tc.BasePath, "expected", fmt.Sprintf("%s.go.tmpl", name)))
// 			// 	assert.NoError(t, err)
// 			// 	assert.Equal(t, formatTextContent(string(fnContent)), formatTextContent(builder.String()), name)
// 			// }
// 		})
// 	}
// }
