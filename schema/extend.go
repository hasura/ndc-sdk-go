package schema

import (
	"encoding/json"
	"fmt"
)

/*
 * Types track the valid representations of values as JSON
 */

// NamedType represents a named type
type NamedType struct {
	Type string `json:"type"`
	// The name can refer to a primitive type or a scalar type
	Name string `json:"name"`
}

// NewNamedType creates a new NamedType instance
func NewNamedType(name string) *NamedType {
	return &NamedType{
		Type: "named",
		Name: name,
	}
}

// NullableType represents a nullable type
type NullableType struct {
	Type string `json:"type"`
	// The type of the non-null inhabitants of this type
	UnderlyingType any `json:"underlying_type"`
}

// NewNullableNamedType creates a new NullableType instance with underlying named type
func NewNullableNamedType(name string) *NullableType {
	return &NullableType{
		Type:           "nullable",
		UnderlyingType: NewNamedType(name),
	}
}

// NewNullableArrayType creates a new NullableType instance with underlying array type
func NewNullableArrayType(elementType any) *NullableType {
	return &NullableType{
		Type:           "nullable",
		UnderlyingType: NewArrayType(elementType),
	}
}

// ArrayType represents an array type
type ArrayType struct {
	Type string `json:"type"`
	// The type of the elements of the array
	ElementType any `json:"element_type"`
}

// NewArrayType creates a new ArrayType instance
func NewArrayType(elementType any) *ArrayType {
	return &ArrayType{
		Type:        "array",
		ElementType: elementType,
	}
}

type ArgumentType string

const (
	ArgumentLiteral  ArgumentType = "literal"
	ArgumentVariable ArgumentType = "variable"
)

// Argument is provided by reference to a variable or as a literal value
type Argument struct {
	Type  ArgumentType `json:"type" mapstructure:"type"`
	Name  string       `json:"name" mapstructure:"name"`
	Value any          `json:"value" mapstructure:"value"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Argument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	rawArgumentType, ok := raw["type"]
	if !ok || rawArgumentType == nil {
		return fmt.Errorf("field type in QueryRequestArguments: required")
	}

	argumentType, ok := rawArgumentType.(string)
	if !ok || (argumentType != string(ArgumentLiteral) && argumentType != string(ArgumentVariable)) {
		return fmt.Errorf("field type in QueryRequestArguments: expect %s or %s", ArgumentLiteral, ArgumentVariable)
	}

	arg := Argument{
		Type: ArgumentType(argumentType),
	}

	switch arg.Type {
	case ArgumentLiteral:
		if value, ok := raw["value"]; !ok {
			return fmt.Errorf("field value in QueryRequestArguments is required for literal type")
		} else {
			arg.Value = value
		}
		break
	case ArgumentVariable:
		rawName, ok := raw["name"]
		if !ok || rawName == nil {
			return fmt.Errorf("field name in QueryRequestArguments is required for variable type")
		}
		name, ok := rawName.(string)
		if !ok {
			return fmt.Errorf("field name in QueryRequestArguments is required for variable type")
		}
		arg.Name = name
		break
	}

	*j = arg
	return nil
}
