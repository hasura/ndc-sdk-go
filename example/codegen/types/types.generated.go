// Code generated by github.com/hasura/ndc-sdk-go/codegen, DO NOT EDIT.
package types

import (
	"encoding/json"
	"errors"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

var types_Decoder = utils.NewDecoder()

// ToMap encodes the struct to a value map
func (j Author) ToMap() (map[string]any, error) {
	r := make(map[string]any)
	r["created_at"] = j.CreatedAt
	r["id"] = j.ID
	r["tags"] = j.Tags

	return r, nil
}

// ScalarName get the schema name of the scalar
func (j CommentText) ScalarName() string {
	return "CommentString"
}

// ScalarName get the schema name of the scalar
func (j SomeEnum) ScalarName() string {
	return "SomeEnum"
}

const (
	SomeEnumFoo SomeEnum = "foo"
	SomeEnumBar SomeEnum = "bar"
)

var enumValues_SomeEnum = []SomeEnum{SomeEnumFoo, SomeEnumBar}

// ParseSomeEnum parses a SomeEnum enum from string
func ParseSomeEnum(input string) (SomeEnum, error) {
	result := SomeEnum(input)
	if !schema.Contains(enumValues_SomeEnum, result) {
		return SomeEnum(""), errors.New("failed to parse SomeEnum, expect one of SomeEnumFoo, SomeEnumBar")
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j SomeEnum) IsValid() bool {
	return schema.Contains(enumValues_SomeEnum, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SomeEnum) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseSomeEnum(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// FromValue decodes the scalar from an unknown value
func (s *SomeEnum) FromValue(value any) error {
	valueStr, err := utils.DecodeNullableString(value)
	if err != nil {
		return err
	}
	if valueStr == nil {
		return nil
	}
	result, err := ParseSomeEnum(*valueStr)
	if err != nil {
		return err
	}

	*s = result
	return nil
}
