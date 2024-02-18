package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"os"
	"path"
	"text/template"
)

//go:embed templates/connector/connector.go.tmpl
var connectorTemplate string

const (
	connectorOutputFile = "connector.generated.go"
)

func parseAndGenerateConnector(basePath string, directories []string, moduleName string) error {
	if err := os.Chdir(basePath); err != nil {
		return err
	}
	sm, err := parseRawConnectorSchemaFromGoCode(".", directories)
	if err != nil {
		return err
	}

	return generateConnector(sm, ".", moduleName)
}

func generateConnector(rawSchema *RawConnectorSchema, srcPath string, moduleName string) error {
	fileTemplate, err := template.New(connectorOutputFile).Parse(connectorTemplate)
	if err != nil {
		return err
	}

	// generate schema.generated.json
	schemaBytes, err := json.MarshalIndent(rawSchema.Schema(), "", "  ")
	if err != nil {
		return err
	}

	schemaPath := path.Join(srcPath, "schema.generated.json")
	if err := os.WriteFile(schemaPath, schemaBytes, 0644); err != nil {
		return err
	}

	targetPath := path.Join(srcPath, connectorOutputFile)
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	w := bufio.NewWriter(f)
	err = fileTemplate.Execute(w, map[string]any{
		"Module": moduleName,
	})
	if err != nil {
		return err
	}
	return w.Flush()
}
