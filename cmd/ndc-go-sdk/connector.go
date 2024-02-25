package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const (
	connectorOutputFile = "connector.generated.go"
)

func parseAndGenerateConnector(basePath string, directories []string, moduleName string) error {
	if err := os.Chdir(basePath); err != nil {
		return err
	}

	sm, err := parseRawConnectorSchemaFromGoCode(moduleName, ".", directories)
	if err != nil {
		return err
	}

	return generateConnector(sm, ".", moduleName)
}

func generateConnector(rawSchema *RawConnectorSchema, srcPath string, moduleName string) error {
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
	defer func() {
		_ = w.Flush()
	}()

	return genConnectorCodeFromTemplate(w, moduleName, rawSchema)
}

func genConnectorCodeFromTemplate(w io.Writer, moduleName string, rawSchema *RawConnectorSchema) error {
	importLines := []string{}
	for importPath := range rawSchema.Imports {
		importLines = append(importLines, fmt.Sprintf(`"%s"`, importPath))
	}

	return connectorTemplate.Execute(w, map[string]any{
		"Imports":    strings.Join(importLines, "\n"),
		"Module":     moduleName,
		"Queries":    genConnectorFunctions(rawSchema),
		"Procedures": genConnectorProcedures(rawSchema),
	})
}

func genConnectorFunctions(rawSchema *RawConnectorSchema) string {
	if len(rawSchema.Functions) == 0 {
		return ""
	}

	var functionCases []string
	for _, fn := range rawSchema.Functions {
		var argumentStr string
		var argumentParamStr string
		if fn.ArgumentsType != "" {
			argumentStr = fmt.Sprintf(`args, err := schema.ResolveArguments[%s.%s](request.Arguments, variables)
		if err != nil {
			return nil, schema.BadRequestError("failed to resolve arguments", map[string]any{
				"cause": err.Error(),
			})
		}`, fn.PackageName, fn.ArgumentsType)
			argumentParamStr = ", args"
		}
		fnCase := fmt.Sprintf(`	case "%s":
		%s
		return %s.%s(ctx, state%s)`, fn.Name, argumentStr, fn.PackageName, fn.OriginName, argumentParamStr)
		functionCases = append(functionCases, fnCase)
	}

	return strings.Join(functionCases, "\n")
}

func genConnectorProcedures(rawSchema *RawConnectorSchema) string {
	if len(rawSchema.Procedures) == 0 {
		return ""
	}

	var cases []string
	for _, fn := range rawSchema.Procedures {
		var argumentStr string
		var argumentParamStr string
		if fn.ArgumentsType != "" {
			argumentStr = fmt.Sprintf(`var args %s.%s
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}`, fn.PackageName, fn.ArgumentsType)
			argumentParamStr = ", &args"
		}
		fnCase := fmt.Sprintf(`	case "%s":
		%s
		rawResult, err = %s.%s(ctx, state%s)`, fn.Name, argumentStr, fn.PackageName, fn.OriginName, argumentParamStr)
		cases = append(cases, fnCase)
	}

	return strings.Join(cases, "\n")
}
