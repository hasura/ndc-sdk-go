package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"path"
)

//go:embed templates/connector/connector.go.tmpl
var connectorTemplate []byte

func generateConnector(rawSchema *RawConnectorSchema, srcPath string) error {

	// generate schema.generated.json
	schemaBytes, err := json.MarshalIndent(rawSchema.Schema(), "", "  ")
	if err != nil {
		return err
	}

	schemaPath := path.Join(srcPath, "schema.generated.json")
	if err := os.WriteFile(schemaPath, schemaBytes, 0644); err != nil {
		return err
	}

	return nil
}
