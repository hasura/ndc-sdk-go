package internal

import (
	"fmt"
	"slices"
)

// OperationNamingStyle the enum for operation naming style
type OperationNamingStyle string

const (
	StyleCamelCase OperationNamingStyle = "camel-case"
	StyleSnakeCase OperationNamingStyle = "snake-case"
)

var enumValues_OperationNamingStyle = []OperationNamingStyle{
	StyleCamelCase, StyleSnakeCase,
}

// IsValid checks if the value is invalid
func (j OperationNamingStyle) IsValid() bool {
	return slices.Contains(enumValues_OperationNamingStyle, j)
}

// ParseOperationNamingStyle parses a OperationNamingStyle enum from string
func ParseOperationNamingStyle(input string) (OperationNamingStyle, error) {
	result := OperationNamingStyle(input)
	if !result.IsValid() {
		return OperationNamingStyle(""), fmt.Errorf("failed to parse OperationNamingStyle, expect one of %v, got: %s", enumValues_OperationNamingStyle, input)
	}

	return result, nil
}
