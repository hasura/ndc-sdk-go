package internal

import (
	"regexp"
	"strings"
)

var nonAlphaDigitRegex = regexp.MustCompile(`[^\w]+`)

// ToCamelCase convert a string to camelCase
func ToCamelCase(input string) string {
	pascalCase := ToPascalCase(input)
	if pascalCase == "" {
		return pascalCase
	}
	return strings.ToLower(pascalCase[:1]) + pascalCase[1:]
}

// ToPascalCase convert a string to PascalCase
func ToPascalCase(input string) string {
	if input == "" {
		return input
	}
	input = nonAlphaDigitRegex.ReplaceAllString(input, "_")
	parts := strings.Split(input, "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}
