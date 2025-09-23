package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

// WriteFileStrategy decides the strategy to do when the written file exists.
type WriteFileStrategy string

const (
	WriteFileStrategyNone           WriteFileStrategy = "none"
	WriteFileStrategyOverride       WriteFileStrategy = "override"
	WriteFileStrategyUpdateResponse WriteFileStrategy = "update_response"
)

var enumValues_WriteFileStrategy = []WriteFileStrategy{
	WriteFileStrategyNone,
	WriteFileStrategyOverride,
	WriteFileStrategyUpdateResponse,
}

// ParseWriteFileStrategy parses a WriteFileStrategy enum from string.
func ParseWriteFileStrategy(input string) (WriteFileStrategy, error) {
	result := WriteFileStrategy(input)
	if !slices.Contains(enumValues_WriteFileStrategy, result) {
		return WriteFileStrategy(
				"",
			), fmt.Errorf(
				"failed to parse WriteFileStrategy, expect one of %v, got: %s",
				enumValues_WriteFileStrategy,
				input,
			)
	}

	return result, nil
}

// WritePrettyFileJSON writes JSON data with indent.
func WritePrettyFileJSON(fileName string, data any) error {
	rawBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, rawBytes, 0o644)
}
