package internal

import (
	"fmt"
	"strings"

	"github.com/fatih/structtag"
)

// NDCTagInfo holds information of an object field that is parsed from the struct tag.
type NDCTagInfo struct {
	Name                string
	Ignored             bool
	OmitEmpty           bool
	PredicateObjectName string
}

// parses connector and json tag information from the raw struct tag.
func parseNDCTagInfo(input string) (NDCTagInfo, error) {
	result := NDCTagInfo{}
	if input == "" {
		return result, nil
	}

	tags, err := structtag.Parse(input)
	if err != nil || tags == nil {
		return result, err
	}

	jsonTag, err := tags.Get("json")
	if err == nil {
		if jsonTag.Value() == "-" {
			result.Ignored = true
		}

		result.Name = strings.TrimSpace(jsonTag.Name)
		result.OmitEmpty = jsonTag.HasOption("omitempty")
	}

	rawTag, err := tags.Get("ndc")
	if err != nil {
		return result, nil //nolint:nilerr
	}

	for _, tag := range append(rawTag.Options, rawTag.Name) {
		keyValue := strings.Split(strings.TrimSpace(tag), "=")
		if len(keyValue) != 2 {
			return result, fmt.Errorf("invalid tag: %s", tag)
		}

		if keyValue[1] == "" {
			return result, fmt.Errorf("value of %s tag is empty", keyValue[0])
		}

		switch keyValue[0] {
		case "predicate":
			result.PredicateObjectName = keyValue[1]
		default:
			return result, fmt.Errorf("invalid tag key: %s", keyValue[0])
		}
	}

	return result, nil
}
