package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
 * Types track the valid representations of values as JSON
 */

type TypeEnum string

const (
	TypeNamed     TypeEnum = "named"
	TypeNullable  TypeEnum = "nullable"
	TypeArray     TypeEnum = "array"
	TypePredicate TypeEnum = "predicate"
)

var enumValues_Type = []TypeEnum{
	TypeNamed,
	TypeNullable,
	TypeArray,
	TypePredicate,
}

// ParseTypeEnum parses a type enum from string
func ParseTypeEnum(input string) (TypeEnum, error) {
	result := TypeEnum(input)

	if !result.IsValid() {
		return TypeEnum(""), fmt.Errorf("failed to parse TypeEnum, expect one of %v, got %s", enumValues_Type, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j TypeEnum) IsValid() bool {
	return Contains(enumValues_Type, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TypeEnum) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseTypeEnum(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// Types track the valid representations of values as JSON
type Type map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Type) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in Type: required")
	}

	var ty TypeEnum
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in Type: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case TypeNamed:
		rawName, ok := raw["name"]
		if !ok {
			return errors.New("field name in Type is required for named type")
		}
		var name string
		if err := json.Unmarshal(rawName, &name); err != nil {
			return fmt.Errorf("field name in Type: %s", err)
		}
		if name == "" {
			return fmt.Errorf("field name in Type: required")
		}
		result["name"] = name
	case TypeNullable:
		rawUnderlyingType, ok := raw["underlying_type"]
		if !ok {
			return errors.New("field underlying_type in Type is required for nullable type")
		}
		var underlyingType Type
		if err := json.Unmarshal(rawUnderlyingType, &underlyingType); err != nil {
			return fmt.Errorf("field underlying_type in Type: %s", err)
		}
		result["underlying_type"] = underlyingType
	case TypeArray:
		rawElementType, ok := raw["element_type"]
		if !ok {
			return errors.New("field element_type in Type is required for array type")
		}
		var elementType Type
		if err := json.Unmarshal(rawElementType, &elementType); err != nil {
			return fmt.Errorf("field element_type in Type: %s", err)
		}
		result["element_type"] = elementType
	case TypePredicate:
		rawName, ok := raw["object_type_name"]
		if !ok {
			return errors.New("field object_type_name in Type is required for predicate type")
		}
		var objectTypeName string
		if err := json.Unmarshal(rawName, &objectTypeName); err != nil {
			return fmt.Errorf("field object_type_name in Type: %s", err)
		}
		if objectTypeName == "" {
			return fmt.Errorf("field object_type_name in Type: required")
		}
		result["object_type_name"] = objectTypeName
	}
	*j = result
	return nil
}

// IsZero checks whether is the instance is empty
func (ty Type) IsZero() bool {
	if len(ty) == 0 {
		return true
	}
	_, ok := ty["type"]
	return !ok
}

// Type gets the type enum of the current type
func (ty Type) Type() (TypeEnum, error) {
	t, ok := ty["type"]
	if !ok {
		return TypeEnum(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseTypeEnum(raw)
		if err != nil {
			return TypeEnum(""), err
		}
		return v, nil
	case TypeEnum:
		return raw, nil
	default:
		return TypeEnum(""), fmt.Errorf("invalid Type type: %+v", t)
	}
}

// AsNamed tries to convert the current type to NamedType
func (ty Type) AsNamed() (*NamedType, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeNamed {
		return nil, fmt.Errorf("invalid Type type; expected %s, got %s", TypeNamed, t)
	}
	return &NamedType{
		Type: t,
		Name: getStringValueByKey(ty, "name"),
	}, nil
}

// AsNullable tries to convert the current type to NullableType
func (ty Type) AsNullable() (*NullableType, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeNullable {
		return nil, fmt.Errorf("invalid Type type; expected %s, got %s", TypeNullable, t)
	}

	rawUnderlyingType, ok := ty["underlying_type"]
	if !ok {
		return nil, errors.New("underlying_type is required")
	}
	underlyingType, ok := rawUnderlyingType.(Type)
	if !ok {
		return nil, errors.New("underlying_type is not Type type")
	}
	return &NullableType{
		Type:           t,
		UnderlyingType: underlyingType,
	}, nil
}

// AsArray tries to convert the current type to ArrayType
func (ty Type) AsArray() (*ArrayType, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeArray {
		return nil, fmt.Errorf("invalid Type type; expected %s, got %s", TypeArray, t)
	}

	rawElementType, ok := ty["element_type"]
	if !ok {
		return nil, errors.New("element_type is required in Type")
	}
	elementType, ok := rawElementType.(Type)
	if !ok {
		return nil, errors.New("element_type is not Type type")
	}
	return &ArrayType{
		Type:        t,
		ElementType: elementType,
	}, nil
}

// AsPredicate tries to convert the current type to PredicateType
func (ty Type) AsPredicate() (*PredicateType, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypePredicate {
		return nil, fmt.Errorf("invalid Type type; expected %s, got %s", TypePredicate, t)
	}

	return &PredicateType{
		Type:           t,
		ObjectTypeName: getStringValueByKey(ty, "object_type_name"),
	}, nil
}

// Interface converts the instance to the TypeEncoder interface
func (ty Type) Interface() TypeEncoder {
	result, _ := ty.InterfaceT()
	return result
}

// InterfaceT converts the instance to the TypeEncoder interface safely with explicit error
func (ty Type) InterfaceT() (TypeEncoder, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case TypeNamed:
		return ty.AsNamed()
	case TypeNullable:
		return ty.AsNullable()
	case TypeArray:
		return ty.AsArray()
	case TypePredicate:
		return ty.AsPredicate()
	default:
		return nil, fmt.Errorf("invalid Type type: %s", t)
	}
}

// String implements the fmt.Stringer interface
func (ty Type) String() string {
	if ty.IsZero() {
		return "zero_type"
	}
	t, err := ty.InterfaceT()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint(t)
}

// TypeEncoder abstracts the Type interface
type TypeEncoder interface {
	Encode() Type
}

// NamedType represents a named type
type NamedType struct {
	Type TypeEnum `json:"type" yaml:"type" mapstructure:"type"`
	// The name can refer to a primitive type or a scalar type
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// NewNamedType creates a new NamedType instance
func NewNamedType(name string) *NamedType {
	return &NamedType{
		Type: TypeNamed,
		Name: name,
	}
}

// Encode returns the raw Type instance
func (ty NamedType) Encode() Type {
	return map[string]any{
		"type": ty.Type,
		"name": ty.Name,
	}
}

// String implements the fmt.Stringer interface
func (ty NamedType) String() string {
	return ty.Name
}

// NullableType represents a nullable type
type NullableType struct {
	Type TypeEnum `json:"type" yaml:"type" mapstructure:"type"`
	// The type of the non-null inhabitants of this type
	UnderlyingType Type `json:"underlying_type" yaml:"underlying_type" mapstructure:"underlying_type"`
}

// NewNullableType creates a new NullableType instance with underlying type
func NewNullableType(underlyingType TypeEncoder) *NullableType {
	return &NullableType{
		Type:           TypeNullable,
		UnderlyingType: underlyingType.Encode(),
	}
}

// NewNullableNamedType creates a new NullableType instance with underlying named type
func NewNullableNamedType(name string) *NullableType {
	return &NullableType{
		Type:           TypeNullable,
		UnderlyingType: NewNamedType(name).Encode(),
	}
}

// Encode returns the raw Type instance
func (ty NullableType) Encode() Type {
	return map[string]any{
		"type":            ty.Type,
		"underlying_type": ty.UnderlyingType,
	}
}

// String implements the fmt.Stringer interface
func (ty NullableType) String() string {
	return fmt.Sprintf("Nullable<%s>", ty.UnderlyingType)
}

// ArrayType represents an array type
type ArrayType struct {
	Type TypeEnum `json:"type" yaml:"type" mapstructure:"type"`
	// The type of the elements of the array
	ElementType Type `json:"element_type" yaml:"element_type" mapstructure:"element_type"`
}

// Encode returns the raw Type instance
func (ty ArrayType) Encode() Type {
	return map[string]any{
		"type":         ty.Type,
		"element_type": ty.ElementType,
	}
}

// String implements the fmt.Stringer interface
func (ty ArrayType) String() string {
	return fmt.Sprintf("Array<%s>", ty.ElementType)
}

// NewArrayType creates a new ArrayType instance
func NewArrayType(elementType TypeEncoder) *ArrayType {
	return &ArrayType{
		Type:        TypeArray,
		ElementType: elementType.Encode(),
	}
}

// PredicateType represents a predicate type for a given object type
type PredicateType struct {
	Type TypeEnum `json:"type" yaml:"type" mapstructure:"type"`
	// The name can refer to a primitive type or a scalar type
	ObjectTypeName string `json:"object_type_name" yaml:"object_type_name" mapstructure:"object_type_name"`
}

// NewPredicateType creates a new PredicateType instance
func NewPredicateType(objectTypeName string) *PredicateType {
	return &PredicateType{
		Type:           TypePredicate,
		ObjectTypeName: objectTypeName,
	}
}

// Encode returns the raw Type instance
func (ty PredicateType) Encode() Type {
	return map[string]any{
		"type":             ty.Type,
		"object_type_name": ty.ObjectTypeName,
	}
}

// String implements the fmt.Stringer interface
func (ty PredicateType) String() string {
	return fmt.Sprintf("Predicate<%s>", ty.ObjectTypeName)
}

// GetUnderlyingNamedType gets the underlying named type of the input type recursively if exists
func GetUnderlyingNamedType(input Type) *NamedType {
	if len(input) == 0 {
		return nil
	}
	switch ty := input.Interface().(type) {
	case *NullableType:
		return GetUnderlyingNamedType(ty.UnderlyingType)
	case *ArrayType:
		return GetUnderlyingNamedType(ty.ElementType)
	case *NamedType:
		return ty
	default:
		return nil
	}
}

// UnwrapNullableType unwraps nullable types from the input type recursively
func UnwrapNullableType(input Type) Type {
	if input == nil || input.IsZero() {
		return nil
	}
	switch ty := input.Interface().(type) {
	case *NullableType:
		return UnwrapNullableType(ty.UnderlyingType)
	default:
		return input
	}
}
