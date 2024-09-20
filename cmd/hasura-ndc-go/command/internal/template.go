package internal

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

const (
	connectorOutputFile   = "connector.generated.go"
	schemaOutputJSONFile  = "schema.generated.json"
	schemaOutputGoFile    = "schema.generated.go"
	typeMethodsOutputFile = "types.generated.go"
)

//go:embed templates/connector/connector.go.tmpl
var connectorTemplateStr string
var connectorTemplate *template.Template

func init() {
	var err error
	connectorTemplate, err = template.New(connectorOutputFile).Parse(connectorTemplateStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse connector template: %s", err))
	}

	strcase.ConfigureAcronym("API", "Api")
	strcase.ConfigureAcronym("REST", "Rest")
	strcase.ConfigureAcronym("HTTP", "Http")
	strcase.ConfigureAcronym("SQL", "sql")
}

func renderFileHeader(builder *strings.Builder, packageName string) {
	_, _ = builder.WriteString(`// Code generated by github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go, DO NOT EDIT.
package `)
	_, _ = builder.WriteString(packageName)
	_, _ = builder.WriteRune('\n')
}

func writeIndent(builder *strings.Builder, num int) {
	for i := 0; i < num; i++ {
		_, _ = builder.WriteRune(' ')
	}
}