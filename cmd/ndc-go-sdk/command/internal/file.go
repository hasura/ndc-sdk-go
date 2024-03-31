package internal

import (
	"encoding/json"
	"os"
)

// WritePrettyFileJSON writes JSON data with indent
func WritePrettyFileJSON(fileName string, data any) error {
	if _, err := os.Stat(fileName); err == nil {
		return nil
	}
	rawBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, rawBytes, 0644)
}
