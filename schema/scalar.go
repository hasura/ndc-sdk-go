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
	// JSON booleans
	TypeRepresentationTypeBoolean TypeRepresentationType = "boolean"
	// JSON booleans
	TypeRepresentationTypeString TypeRepresentationType = "string"
	// Any JSON number
	//
	// Deprecated: [Deprecate Int and Number representations]
	//
	// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
	TypeRepresentationTypeNumber TypeRepresentationType = "number"
	// Any JSON number, with no decimal part
	//
	// Deprecated: [Deprecate Int and Number representations]
	//
	// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
	TypeRepresentationTypeInteger TypeRepresentationType = "integer"
	// One of the specified string values
	TypeRepresentationTypeEnum TypeRepresentationType = "enum"
	// A 8-bit signed integer with a minimum value of -2^7 and a maximum value of 2^7 - 1
	TypeRepresentationTypeInt8 TypeRepresentationType = "int8"
	// A 16-bit signed integer with a minimum value of -2^15 and a maximum value of 2^15 - 1
	TypeRepresentationTypeInt16 TypeRepresentationType = "int16"
	// A 32-bit signed integer with a minimum value of -2^31 and a maximum value of 2^31 - 1
	TypeRepresentationTypeInt32 TypeRepresentationType = "int32"
	// A 64-bit signed integer with a minimum value of -2^63 and a maximum value of 2^63 - 1
	TypeRepresentationTypeInt64 TypeRepresentationType = "int64"
	// An IEEE-754 single-precision floating-point number
	TypeRepresentationTypeFloat32 TypeRepresentationType = "float32"
	// An IEEE-754 double-precision floating-point number
	TypeRepresentationTypeFloat64 TypeRepresentationType = "float64"
	// Arbitrary-precision integer string
	TypeRepresentationTypeBigInteger TypeRepresentationType = "biginteger"
	// Arbitrary-precision decimal string
	TypeRepresentationTypeBigDecimal TypeRepresentationType = "bigdecimal"
	// UUID string (8-4-4-4-12)
	TypeRepresentationTypeUUID TypeRepresentationType = "uuid"
	// ISO 8601 date
	TypeRepresentationTypeDate TypeRepresentationType = "date"
	// ISO 8601 timestamp
	TypeRepresentationTypeTimestamp TypeRepresentationType = "timestamp"
	// ISO 8601 timestamp-with-timezone
	TypeRepresentationTypeTimestampTZ TypeRepresentationType = "timestamptz"
	// GeoJSON, per RFC 7946
	TypeRepresentationTypeGeography TypeRepresentationType = "geography"
	// GeoJSON Geometry object, per RFC 7946
	TypeRepresentationTypeGeometry TypeRepresentationType = "geometry"
	// Base64-encoded bytes
	TypeRepresentationTypeBytes TypeRepresentationType = "bytes"
	// Arbitrary JSON
	TypeRepresentationTypeJSON TypeRepresentationType = "json"
)

var enumValues_TypeRepresentationType = []TypeRepresentationType{
	TypeRepresentationTypeBoolean,
	TypeRepresentationTypeString,
	TypeRepresentationTypeNumber,
	TypeRepresentationTypeInteger,
	TypeRepresentationTypeEnum,
	TypeRepresentationTypeInt8,
	TypeRepresentationTypeInt16,
	TypeRepresentationTypeInt32,
	TypeRepresentationTypeInt64,
	TypeRepresentationTypeFloat32,
	TypeRepresentationTypeFloat64,
	TypeRepresentationTypeBigInteger,
	TypeRepresentationTypeBigDecimal,
	TypeRepresentationTypeUUID,
	TypeRepresentationTypeDate,
	TypeRepresentationTypeTimestamp,
	TypeRepresentationTypeTimestampTZ,
	TypeRepresentationTypeGeography,
	TypeRepresentationTypeGeometry,
	TypeRepresentationTypeBytes,
	TypeRepresentationTypeJSON,
}

// ParseTypeRepresentationType parses a TypeRepresentationType enum from string
func ParseTypeRepresentationType(input string) (TypeRepresentationType, error) {
	result := TypeRepresentationType(input)
	if !result.IsValid() {
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
//
// Deprecated: [Deprecate Int and Number representations]
//
// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
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

// AsInteger tries to convert the current type to TypeRepresentationInteger
//
// Deprecated: [Deprecate Int and Number representations]
//
// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
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

// AsInt8 tries to convert the current type to TypeRepresentationInt8
func (ty TypeRepresentation) AsInt8() (*TypeRepresentationInt8, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeInt8 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeInt8, t)
	}

	return &TypeRepresentationInt8{
		Type: t,
	}, nil
}

// AsInt16 tries to convert the current type to TypeRepresentationInt16
func (ty TypeRepresentation) AsInt16() (*TypeRepresentationInt16, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeInt16 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeInt16, t)
	}

	return &TypeRepresentationInt16{
		Type: t,
	}, nil
}

// AsInt32 tries to convert the current type to TypeRepresentationInt32
func (ty TypeRepresentation) AsInt32() (*TypeRepresentationInt32, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeInt32 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeInt32, t)
	}

	return &TypeRepresentationInt32{
		Type: t,
	}, nil
}

// AsInt64 tries to convert the current type to TypeRepresentationInt64
func (ty TypeRepresentation) AsInt64() (*TypeRepresentationInt64, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeInt64 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeInt64, t)
	}

	return &TypeRepresentationInt64{
		Type: t,
	}, nil
}

// AsFloat32 tries to convert the current type to TypeRepresentationFloat32
func (ty TypeRepresentation) AsFloat32() (*TypeRepresentationFloat32, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeFloat32 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeFloat32, t)
	}

	return &TypeRepresentationFloat32{
		Type: t,
	}, nil
}

// AsFloat64 tries to convert the current type to TypeRepresentationFloat64
func (ty TypeRepresentation) AsFloat64() (*TypeRepresentationFloat64, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeFloat64 {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeFloat64, t)
	}

	return &TypeRepresentationFloat64{
		Type: t,
	}, nil
}

// AsBigInteger tries to convert the current type to TypeRepresentationBigInteger
func (ty TypeRepresentation) AsBigInteger() (*TypeRepresentationBigInteger, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeBigInteger {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeBigInteger, t)
	}

	return &TypeRepresentationBigInteger{
		Type: t,
	}, nil
}

// AsBigDecimal tries to convert the current type to TypeRepresentationBigDecimal
func (ty TypeRepresentation) AsBigDecimal() (*TypeRepresentationBigDecimal, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeBigDecimal {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeBigDecimal, t)
	}

	return &TypeRepresentationBigDecimal{
		Type: t,
	}, nil
}

// AsUUID tries to convert the current type to TypeRepresentationUUID
func (ty TypeRepresentation) AsUUID() (*TypeRepresentationUUID, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeUUID {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeUUID, t)
	}

	return &TypeRepresentationUUID{
		Type: t,
	}, nil
}

// AsDate tries to convert the current type to TypeRepresentationDate
func (ty TypeRepresentation) AsDate() (*TypeRepresentationDate, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeDate {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeDate, t)
	}

	return &TypeRepresentationDate{
		Type: t,
	}, nil
}

// AsTimestamp tries to convert the current type to TypeRepresentationTimestamp
func (ty TypeRepresentation) AsTimestamp() (*TypeRepresentationTimestamp, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeTimestamp {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeTimestamp, t)
	}

	return &TypeRepresentationTimestamp{
		Type: t,
	}, nil
}

// AsTimestampTZ tries to convert the current type to TypeRepresentationTimestampTZ
func (ty TypeRepresentation) AsTimestampTZ() (*TypeRepresentationTimestampTZ, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeTimestampTZ {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeTimestampTZ, t)
	}

	return &TypeRepresentationTimestampTZ{
		Type: t,
	}, nil
}

// AsGeography tries to convert the current type to TypeRepresentationGeography
func (ty TypeRepresentation) AsGeography() (*TypeRepresentationGeography, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeGeography {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeGeography, t)
	}

	return &TypeRepresentationGeography{
		Type: t,
	}, nil
}

// AsGeometry tries to convert the current type to TypeRepresentationGeometry
func (ty TypeRepresentation) AsGeometry() (*TypeRepresentationGeometry, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeGeometry {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeGeometry, t)
	}

	return &TypeRepresentationGeometry{
		Type: t,
	}, nil
}

// AsBytes tries to convert the current type to TypeRepresentationBytes
func (ty TypeRepresentation) AsBytes() (*TypeRepresentationBytes, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeBytes {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeBytes, t)
	}

	return &TypeRepresentationBytes{
		Type: t,
	}, nil
}

// AsJSON tries to convert the current type to TypeRepresentationJSON
func (ty TypeRepresentation) AsJSON() (*TypeRepresentationJSON, error) {
	t, err := ty.Type()
	if err != nil {
		return nil, err
	}
	if t != TypeRepresentationTypeJSON {
		return nil, fmt.Errorf("invalid TypeRepresentation type; expected %s, got %s", TypeRepresentationTypeJSON, t)
	}

	return &TypeRepresentationJSON{
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
	case TypeRepresentationTypeInt8:
		return ty.AsInt8()
	case TypeRepresentationTypeInt16:
		return ty.AsInt16()
	case TypeRepresentationTypeInt32:
		return ty.AsInt32()
	case TypeRepresentationTypeInt64:
		return ty.AsInt64()
	case TypeRepresentationTypeFloat32:
		return ty.AsFloat32()
	case TypeRepresentationTypeFloat64:
		return ty.AsFloat64()
	case TypeRepresentationTypeBigInteger:
		return ty.AsBigInteger()
	case TypeRepresentationTypeBigDecimal:
		return ty.AsBigDecimal()
	case TypeRepresentationTypeUUID:
		return ty.AsUUID()
	case TypeRepresentationTypeDate:
		return ty.AsDate()
	case TypeRepresentationTypeTimestamp:
		return ty.AsTimestamp()
	case TypeRepresentationTypeTimestampTZ:
		return ty.AsTimestampTZ()
	case TypeRepresentationTypeGeography:
		return ty.AsGeography()
	case TypeRepresentationTypeGeometry:
		return ty.AsGeometry()
	case TypeRepresentationTypeBytes:
		return ty.AsBytes()
	case TypeRepresentationTypeJSON:
		return ty.AsJSON()
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
//
// Deprecated: [Deprecate Int and Number representations]
//
// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
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
//
// Deprecated: [Deprecate Int and Number representations]
//
// [Deprecate Int and Number representations]: https://github.com/hasura/ndc-spec/blob/main/rfcs/0007-additional-type-representations.md#deprecate-int-and-number-representations
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

// TypeRepresentationInt8 represents a 8-bit signed integer with a minimum value of -2^7 and a maximum value of 2^7 - 1
type TypeRepresentationInt8 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt8 creates a new TypeRepresentationInt8 instance
func NewTypeRepresentationInt8() *TypeRepresentationInt8 {
	return &TypeRepresentationInt8{
		Type: TypeRepresentationTypeInt8,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationInt8) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt16 represents a 16-bit signed integer with a minimum value of -2^15 and a maximum value of 2^15 - 1
type TypeRepresentationInt16 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt16 creates a new TypeRepresentationInt16 instance
func NewTypeRepresentationInt16() *TypeRepresentationInt16 {
	return &TypeRepresentationInt16{
		Type: TypeRepresentationTypeInt16,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationInt16) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt32 represents a 32-bit signed integer with a minimum value of -2^31 and a maximum value of 2^31 - 1
type TypeRepresentationInt32 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt32 creates a new TypeRepresentationInt32 instance
func NewTypeRepresentationInt32() *TypeRepresentationInt32 {
	return &TypeRepresentationInt32{
		Type: TypeRepresentationTypeInt32,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationInt32) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt64 represents a 64-bit signed integer with a minimum value of -2^63 and a maximum value of 2^63 - 1
type TypeRepresentationInt64 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt64 creates a new TypeRepresentationInt64 instance
func NewTypeRepresentationInt64() *TypeRepresentationInt64 {
	return &TypeRepresentationInt64{
		Type: TypeRepresentationTypeInt64,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationInt64) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationFloat32 represents an IEEE-754 single-precision floating-point number
type TypeRepresentationFloat32 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationFloat32 creates a new TypeRepresentationFloat32 instance
func NewTypeRepresentationFloat32() *TypeRepresentationFloat32 {
	return &TypeRepresentationFloat32{
		Type: TypeRepresentationTypeFloat32,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationFloat32) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationFloat64 represents an IEEE-754 double-precision floating-point number
type TypeRepresentationFloat64 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationFloat64 creates a new TypeRepresentationFloat64 instance
func NewTypeRepresentationFloat64() *TypeRepresentationFloat64 {
	return &TypeRepresentationFloat64{
		Type: TypeRepresentationTypeFloat64,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationFloat64) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBigInteger represents an arbitrary-precision integer string
type TypeRepresentationBigInteger struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBigInteger creates a new TypeRepresentationBigInteger instance
func NewTypeRepresentationBigInteger() *TypeRepresentationBigInteger {
	return &TypeRepresentationBigInteger{
		Type: TypeRepresentationTypeBigInteger,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationBigInteger) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBigDecimal represents an arbitrary-precision decimal string
type TypeRepresentationBigDecimal struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBigDecimal creates a new TypeRepresentationBigDecimal instance
func NewTypeRepresentationBigDecimal() *TypeRepresentationBigDecimal {
	return &TypeRepresentationBigDecimal{
		Type: TypeRepresentationTypeBigDecimal,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationBigDecimal) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationUUID represents an UUID string (8-4-4-4-12)
type TypeRepresentationUUID struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationUUID creates a new TypeRepresentationUUID instance
func NewTypeRepresentationUUID() *TypeRepresentationUUID {
	return &TypeRepresentationUUID{
		Type: TypeRepresentationTypeUUID,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationUUID) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationDate represents an ISO 8601 date
type TypeRepresentationDate struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationDate creates a new TypeRepresentationDate instance
func NewTypeRepresentationDate() *TypeRepresentationDate {
	return &TypeRepresentationDate{
		Type: TypeRepresentationTypeDate,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationDate) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationTimestamp represents an ISO 8601 timestamp
type TypeRepresentationTimestamp struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationTimestamp creates a new TypeRepresentationTimestamp instance
func NewTypeRepresentationTimestamp() *TypeRepresentationTimestamp {
	return &TypeRepresentationTimestamp{
		Type: TypeRepresentationTypeTimestamp,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationTimestamp) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationTimestampTZ represents an ISO 8601 timestamp-with-timezone
type TypeRepresentationTimestampTZ struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationTimestampTZ creates a new TypeRepresentationTimestampTZ instance
func NewTypeRepresentationTimestampTZ() *TypeRepresentationTimestampTZ {
	return &TypeRepresentationTimestampTZ{
		Type: TypeRepresentationTypeTimestampTZ,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationTimestampTZ) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationGeography represents a geography JSON object
type TypeRepresentationGeography struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationGeography creates a new TypeRepresentationGeography instance
func NewTypeRepresentationGeography() *TypeRepresentationGeography {
	return &TypeRepresentationGeography{
		Type: TypeRepresentationTypeGeography,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationGeography) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationGeometry represents a geography JSON object
type TypeRepresentationGeometry struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationGeometry creates a new TypeRepresentationGeometry instance
func NewTypeRepresentationGeometry() *TypeRepresentationGeometry {
	return &TypeRepresentationGeometry{
		Type: TypeRepresentationTypeGeometry,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationGeometry) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBytes represent a base64-encoded bytes
type TypeRepresentationBytes struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBytes creates a new TypeRepresentationBytes instance
func NewTypeRepresentationBytes() *TypeRepresentationBytes {
	return &TypeRepresentationBytes{
		Type: TypeRepresentationTypeBytes,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationBytes) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationJSON represents an arbitrary JSON
type TypeRepresentationJSON struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationJSON creates a new TypeRepresentationBytes instance
func NewTypeRepresentationJSON() *TypeRepresentationJSON {
	return &TypeRepresentationJSON{
		Type: TypeRepresentationTypeJSON,
	}
}

// Encode returns the raw TypeRepresentation instance
func (ty TypeRepresentationJSON) Encode() TypeRepresentation {
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
