package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// NewScalarType creates an empty ScalarType instance
func NewScalarType() *ScalarType {
	return &ScalarType{
		AggregateFunctions:  ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]ComparisonOperatorDefinition{},
	}
}

/*
 * Representations of scalar types
 */

// TypeRepresentationType represents the type enum of TypeRepresentation
type TypeRepresentationType string

const (
	TypeRepresentationTypeBoolean TypeRepresentationType = "boolean"
	TypeRepresentationTypeString  TypeRepresentationType = "string"
	TypeRepresentationTypeNumber  TypeRepresentationType = "number"
	TypeRepresentationTypeInteger TypeRepresentationType = "integer"
	TypeRepresentationTypeEnum    TypeRepresentationType = "enum"
)

var enumValues_TypeRepresentationType = []TypeRepresentationType{
	TypeRepresentationTypeBoolean,
	TypeRepresentationTypeString,
	TypeRepresentationTypeNumber,
	TypeRepresentationTypeInteger,
	TypeRepresentationTypeEnum,
}

// ParseTypeRepresentationType parses a TypeRepresentationType enum from string
func ParseTypeRepresentationType(input string) (TypeRepresentationType, error) {
	result := TypeRepresentationType(input)
	if !Contains(enumValues_TypeRepresentationType, result) {
		return TypeRepresentationType(""), fmt.Errorf("failed to parse TypeRepresentationType, expect one of %v, got: %s", enumValues_TypeRepresentationType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j TypeRepresentationType) IsValid() bool {
	return Contains(enumValues_TypeRepresentationType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TypeRepresentationType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseTypeRepresentationType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// Representations of scalar types
type TypeRepresentation map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *TypeRepresentation) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in TypeRepresentation: required")
	}

	var ty TypeRepresentationType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in TypeRepresentation: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case TypeRepresentationTypeEnum:
		rawOneOf, ok := raw["one_of"]
		if !ok {
			return errors.New("field one_of in TypeRepresentation is required for enum type")
		}
		var oneOf []string
		if err := json.Unmarshal(rawOneOf, &oneOf); err != nil {
			return fmt.Errorf("field one_of in TypeRepresentation: %s", err)
		}
		if len(oneOf) == 0 {
			return errors.New("TypeRepresentation requires at least 1 item in one_of field for enum type")
		}
		result["one_of"] = oneOf
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (ty TypeRepresentation) Type() (TypeRepresentationType, error) {
	t, ok := ty["type"]
	if !ok {
		return TypeRepresentationType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseTypeRepresentationType(raw)
		if err != nil {
			return TypeRepresentationType(""), err
		}
		return v, nil
	case TypeRepresentationType:
		return raw, nil
	default:
		return TypeRepresentationType(""), fmt.Errorf("invalid type: %+v", t)
	}
}

// AsBoolean tries to convert the current type to TypeRepresentationBoolean
func (ty TypeRepresentation) AsBoolean() (*TypeRepresentationBoolean, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeBoolean {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeBoolean, t)
	}

	return &TypeRepresentationBoolean{
		Type: t,
	}, nil
}

// AsString tries to convert the current type to TypeRepresentationString
func (ty TypeRepresentation) AsString() (*TypeRepresentationString, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeString {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeString, t)
	}

	return &TypeRepresentationString{
		Type: t,
	}, nil
}

// AsNumber tries to convert the current type to TypeRepresentationNumber
func (ty TypeRepresentation) AsNumber() (*TypeRepresentationNumber, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeNumber {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeNumber, t)
	}

	return &TypeRepresentationNumber{
		Type: t,
	}, nil
}

// AsInteger tries to convert the current type to TypeRepresentationNumber
func (ty TypeRepresentation) AsInteger() (*TypeRepresentationInteger, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeInteger {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeInteger, t)
	}

	return &TypeRepresentationInteger{
		Type: t,
	}, nil
}

var errTypeRepresentationOneOfRequired = errors.New("TypeRepresentationEnum must have at least 1 item in one_of array")

// AsEnum tries to convert the current type to TypeRepresentationEnum
func (ty TypeRepresentation) AsEnum() (*TypeRepresentationEnum, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeEnum {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeEnum, t)
	}

	rawOneOf, ok := ty["one_of"]
	if !ok {
		return nil, errTypeRepresentationOneOfRequired
	}

	oneOf, ok := rawOneOf.([]string)
	if !ok {
		return nil, errTypeRepresentationOneOfRequired
	}

	if len(oneOf) == 0 {
		return nil, errTypeRepresentationOneOfRequired
	}

	return &TypeRepresentationEnum{
		Type:  t,
		OneOf: oneOf,
	}, nil
}

// Interface converts the instance to the TypeRepresentationEncoder interface
func (ty TypeRepresentation) Interface() TypeRepresentationEncoder {
	result, _ := ty.InterfaceT()
	return result
}

// InterfaceT converts the instance to the TypeRepresentationEncoder interface safely with explicit error
func (ty TypeRepresentation) InterfaceT() (TypeRepresentationEncoder, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case TypeRepresentationTypeBoolean:
		return ty.AsBoolean()
	case TypeRepresentationTypeString:
		return ty.AsString()
	case TypeRepresentationTypeNumber:
		return ty.AsNumber()
	case TypeRepresentationTypeInteger:
		return ty.AsInteger()
	case TypeRepresentationTypeEnum:
		return ty.AsEnum()
	default:
		return nil, fmt.Errorf("invalid TypeRepresentation type: %s", t)
	}
}

// TypeRepresentationEncoder abstracts the TypeRepresentation interface
type TypeRepresentationEncoder interface {
	Encode() TypeRepresentation
}

// TypeRepresentationBoolean represents a JSON boolean type representation
type TypeRepresentationBoolean struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBoolean creates a new TypeRepresentationBoolean instance
func NewTypeRepresentationBoolean() *TypeRepresentationBoolean {
	return &TypeRepresentationBoolean{
		Type: TypeRepresentationTypeBoolean,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationBoolean) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationString represents a JSON string type representation
type TypeRepresentationString struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationString creates a new TypeRepresentationString instance
func NewTypeRepresentationString() *TypeRepresentationString {
	return &TypeRepresentationString{
		Type: TypeRepresentationTypeString,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationString) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationNumber represents a JSON number type representation
type TypeRepresentationNumber struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationNumber creates a new TypeRepresentationNumber instance
func NewTypeRepresentationNumber() *TypeRepresentationNumber {
	return &TypeRepresentationNumber{
		Type: TypeRepresentationTypeNumber,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationNumber) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInteger represents a JSON integer type representation
type TypeRepresentationInteger struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInteger creates a new TypeRepresentationInteger instance
func NewTypeRepresentationInteger() *TypeRepresentationInteger {
	return &TypeRepresentationInteger{
		Type: TypeRepresentationTypeInteger,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationInteger) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationEnum represents an enum type representation
type TypeRepresentationEnum struct {
	Type  TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
	OneOf []string               `json:"one_of" yaml:"one_of" mapstructure:"one_of"`
}

// NewTypeRepresentationEnum creates a new TypeRepresentationEnum instance
func NewTypeRepresentationEnum(oneOf []string) *TypeRepresentationEnum {
	return &TypeRepresentationEnum{
		Type:  TypeRepresentationTypeEnum,
		OneOf: oneOf,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationEnum) Encode() TypeRepresentation {
	return map[string]any{
		"type":   ty.Type,
		"one_of": ty.OneOf,
	}
}
