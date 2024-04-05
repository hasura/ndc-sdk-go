package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

// Decide the strategy to do when the written file exists
type WriteFileStrategy string

const (
	WriteFileStrategyNone     WriteFileStrategy = "none"
	WriteFileStrategyOverride WriteFileStrategy = "override"
)

var enumValues_WriteFileStrategy = []WriteFileStrategy{
	WriteFileStrategyNone,
	WriteFileStrategyOverride,
}

// ParseWriteFileStrategy parses a WriteFileStrategy enum from string
func ParseWriteFileStrategy(input string) (WriteFileStrategy, error) {
	result := WriteFileStrategy(input)
	if !slices.Contains(enumValues_WriteFileStrategy, result) {
		return WriteFileStrategy(""), fmt.Errorf("failed to parse WriteFileStrategy, expect one of %v, got: %s", enumValues_WriteFileStrategy, input)
	}

	return result, nil
}

// WritePrettyFileJSON writes JSON data with indent
func WritePrettyFileJSON(fileName string, data any, strategy WriteFileStrategy) error {
	if _, err := os.Stat(fileName); err == nil {
		switch strategy {
		case WriteFileStrategyNone, "":
			return nil
		}
	}
	rawBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, rawBytes, 0644)
}
