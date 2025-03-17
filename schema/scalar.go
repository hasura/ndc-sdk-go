package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// NewScalarType creates an empty ScalarType instance.
func NewScalarType() *ScalarType {
	return &ScalarType{
		AggregateFunctions:  ScalarTypeAggregateFunctions{},
		ComparisonOperators: map[string]ComparisonOperatorDefinition{},
	}
}

/*
 * Representations of scalar types
 */

// TypeRepresentationType represents the type enum of TypeRepresentation.
type TypeRepresentationType string

const (
	// JSON booleans.
	TypeRepresentationTypeBoolean TypeRepresentationType = "boolean"
	// JSON booleans.
	TypeRepresentationTypeString TypeRepresentationType = "string"
	// One of the specified string values.
	TypeRepresentationTypeEnum TypeRepresentationType = "enum"
	// A 8-bit signed integer with a minimum value of -2^7 and a maximum value of 2^7 - 1.
	TypeRepresentationTypeInt8 TypeRepresentationType = "int8"
	// A 16-bit signed integer with a minimum value of -2^15 and a maximum value of 2^15 - 1.
	TypeRepresentationTypeInt16 TypeRepresentationType = "int16"
	// A 32-bit signed integer with a minimum value of -2^31 and a maximum value of 2^31 - 1.
	TypeRepresentationTypeInt32 TypeRepresentationType = "int32"
	// A 64-bit signed integer with a minimum value of -2^63 and a maximum value of 2^63 - 1.
	TypeRepresentationTypeInt64 TypeRepresentationType = "int64"
	// An IEEE-754 single-precision floating-point number.
	TypeRepresentationTypeFloat32 TypeRepresentationType = "float32"
	// An IEEE-754 double-precision floating-point number.
	TypeRepresentationTypeFloat64 TypeRepresentationType = "float64"
	// Arbitrary-precision integer string.
	TypeRepresentationTypeBigInteger TypeRepresentationType = "biginteger"
	// Arbitrary-precision decimal string.
	TypeRepresentationTypeBigDecimal TypeRepresentationType = "bigdecimal"
	// UUID string (8-4-4-4-12).
	TypeRepresentationTypeUUID TypeRepresentationType = "uuid"
	// ISO 8601 date.
	TypeRepresentationTypeDate TypeRepresentationType = "date"
	// ISO 8601 timestamp.
	TypeRepresentationTypeTimestamp TypeRepresentationType = "timestamp"
	// ISO 8601 timestamp-with-timezone.
	TypeRepresentationTypeTimestampTZ TypeRepresentationType = "timestamptz"
	// GeoJSON, per RFC 7946.
	TypeRepresentationTypeGeography TypeRepresentationType = "geography"
	// GeoJSON Geometry object, per RFC 7946.
	TypeRepresentationTypeGeometry TypeRepresentationType = "geometry"
	// Base64-encoded bytes.
	TypeRepresentationTypeBytes TypeRepresentationType = "bytes"
	// Arbitrary JSON.
	TypeRepresentationTypeJSON TypeRepresentationType = "json"
)

var enumValues_TypeRepresentationType = []TypeRepresentationType{
	TypeRepresentationTypeBoolean,
	TypeRepresentationTypeString,
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

// ParseTypeRepresentationType parses a TypeRepresentationType enum from string.
func ParseTypeRepresentationType(input string) (TypeRepresentationType, error) {
	result := TypeRepresentationType(input)
	if !result.IsValid() {
		return TypeRepresentationType(""), fmt.Errorf("failed to parse TypeRepresentationType, expect one of %v, got: %s", enumValues_TypeRepresentationType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j TypeRepresentationType) IsValid() bool {
	return slices.Contains(enumValues_TypeRepresentationType, j)
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

// Representations of scalar types.
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
		return fmt.Errorf("field type in TypeRepresentation: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	if ty == TypeRepresentationTypeEnum {
		rawOneOf, ok := raw["one_of"]
		if !ok {
			return errors.New("field one_of in TypeRepresentation is required for enum type")
		}
		var oneOf []string
		if err := json.Unmarshal(rawOneOf, &oneOf); err != nil {
			return fmt.Errorf("field one_of in TypeRepresentation: %w", err)
		}
		if len(oneOf) == 0 {
			return errors.New("TypeRepresentation requires at least 1 item in one_of field for enum type")
		}
		result["one_of"] = oneOf
	}

	*j = result
	return nil
}

// Type gets the type enum of the current type.
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

// AsBoolean tries to convert the current type to TypeRepresentationBoolean.
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

// AsString tries to convert the current type to TypeRepresentationString.
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

// AsInt8 tries to convert the current type to TypeRepresentationInt8.
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

// AsInt16 tries to convert the current type to TypeRepresentationInt16.
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

// AsInt32 tries to convert the current type to TypeRepresentationInt32.
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

// AsInt64 tries to convert the current type to TypeRepresentationInt64.
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

// AsFloat32 tries to convert the current type to TypeRepresentationFloat32.
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

// AsFloat64 tries to convert the current type to TypeRepresentationFloat64.
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

// AsBigInteger tries to convert the current type to TypeRepresentationBigInteger.
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

// AsBigDecimal tries to convert the current type to TypeRepresentationBigDecimal.
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

// AsUUID tries to convert the current type to TypeRepresentationUUID.
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

// AsDate tries to convert the current type to TypeRepresentationDate.
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

// AsTimestamp tries to convert the current type to TypeRepresentationTimestamp.
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

// AsTimestampTZ tries to convert the current type to TypeRepresentationTimestampTZ.
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

// AsGeography tries to convert the current type to TypeRepresentationGeography.
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

// AsGeometry tries to convert the current type to TypeRepresentationGeometry.
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

// AsBytes tries to convert the current type to TypeRepresentationBytes.
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

// AsJSON tries to convert the current type to TypeRepresentationJSON.
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

// AsEnum tries to convert the current type to TypeRepresentationEnum.
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

// Interface converts the instance to the TypeRepresentationEncoder interface.
func (ty TypeRepresentation) Interface() TypeRepresentationEncoder {
	result, _ := ty.InterfaceT()
	return result
}

// InterfaceT converts the instance to the TypeRepresentationEncoder interface safely with explicit error.
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

// TypeRepresentationEncoder abstracts the TypeRepresentation interface.
type TypeRepresentationEncoder interface {
	Encode() TypeRepresentation
}

// TypeRepresentationBoolean represents a JSON boolean type representation.
type TypeRepresentationBoolean struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBoolean creates a new TypeRepresentationBoolean instance.
func NewTypeRepresentationBoolean() *TypeRepresentationBoolean {
	return &TypeRepresentationBoolean{
		Type: TypeRepresentationTypeBoolean,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationBoolean) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationString represents a JSON string type representation.
type TypeRepresentationString struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationString creates a new TypeRepresentationString instance.
func NewTypeRepresentationString() *TypeRepresentationString {
	return &TypeRepresentationString{
		Type: TypeRepresentationTypeString,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationString) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt8 represents a 8-bit signed integer with a minimum value of -2^7 and a maximum value of 2^7 - 1.
type TypeRepresentationInt8 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt8 creates a new TypeRepresentationInt8 instance.
func NewTypeRepresentationInt8() *TypeRepresentationInt8 {
	return &TypeRepresentationInt8{
		Type: TypeRepresentationTypeInt8,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationInt8) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt16 represents a 16-bit signed integer with a minimum value of -2^15 and a maximum value of 2^15 - 1.
type TypeRepresentationInt16 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt16 creates a new TypeRepresentationInt16 instance.
func NewTypeRepresentationInt16() *TypeRepresentationInt16 {
	return &TypeRepresentationInt16{
		Type: TypeRepresentationTypeInt16,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationInt16) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt32 represents a 32-bit signed integer with a minimum value of -2^31 and a maximum value of 2^31 - 1.
type TypeRepresentationInt32 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt32 creates a new TypeRepresentationInt32 instance.
func NewTypeRepresentationInt32() *TypeRepresentationInt32 {
	return &TypeRepresentationInt32{
		Type: TypeRepresentationTypeInt32,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationInt32) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationInt64 represents a 64-bit signed integer with a minimum value of -2^63 and a maximum value of 2^63 - 1.
type TypeRepresentationInt64 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationInt64 creates a new TypeRepresentationInt64 instance.
func NewTypeRepresentationInt64() *TypeRepresentationInt64 {
	return &TypeRepresentationInt64{
		Type: TypeRepresentationTypeInt64,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationInt64) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationFloat32 represents an IEEE-754 single-precision floating-point number.
type TypeRepresentationFloat32 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationFloat32 creates a new TypeRepresentationFloat32 instance.
func NewTypeRepresentationFloat32() *TypeRepresentationFloat32 {
	return &TypeRepresentationFloat32{
		Type: TypeRepresentationTypeFloat32,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationFloat32) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationFloat64 represents an IEEE-754 double-precision floating-point number.
type TypeRepresentationFloat64 struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationFloat64 creates a new TypeRepresentationFloat64 instance.
func NewTypeRepresentationFloat64() *TypeRepresentationFloat64 {
	return &TypeRepresentationFloat64{
		Type: TypeRepresentationTypeFloat64,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationFloat64) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBigInteger represents an arbitrary-precision integer string.
type TypeRepresentationBigInteger struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBigInteger creates a new TypeRepresentationBigInteger instance.
func NewTypeRepresentationBigInteger() *TypeRepresentationBigInteger {
	return &TypeRepresentationBigInteger{
		Type: TypeRepresentationTypeBigInteger,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationBigInteger) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBigDecimal represents an arbitrary-precision decimal string.
type TypeRepresentationBigDecimal struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBigDecimal creates a new TypeRepresentationBigDecimal instance.
func NewTypeRepresentationBigDecimal() *TypeRepresentationBigDecimal {
	return &TypeRepresentationBigDecimal{
		Type: TypeRepresentationTypeBigDecimal,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationBigDecimal) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationUUID represents an UUID string (8-4-4-4-12).
type TypeRepresentationUUID struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationUUID creates a new TypeRepresentationUUID instance.
func NewTypeRepresentationUUID() *TypeRepresentationUUID {
	return &TypeRepresentationUUID{
		Type: TypeRepresentationTypeUUID,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationUUID) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationDate represents an ISO 8601 date.
type TypeRepresentationDate struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationDate creates a new TypeRepresentationDate instance.
func NewTypeRepresentationDate() *TypeRepresentationDate {
	return &TypeRepresentationDate{
		Type: TypeRepresentationTypeDate,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationDate) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationTimestamp represents an ISO 8601 timestamp.
type TypeRepresentationTimestamp struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationTimestamp creates a new TypeRepresentationTimestamp instance.
func NewTypeRepresentationTimestamp() *TypeRepresentationTimestamp {
	return &TypeRepresentationTimestamp{
		Type: TypeRepresentationTypeTimestamp,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationTimestamp) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationTimestampTZ represents an ISO 8601 timestamp-with-timezone.
type TypeRepresentationTimestampTZ struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationTimestampTZ creates a new TypeRepresentationTimestampTZ instance.
func NewTypeRepresentationTimestampTZ() *TypeRepresentationTimestampTZ {
	return &TypeRepresentationTimestampTZ{
		Type: TypeRepresentationTypeTimestampTZ,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationTimestampTZ) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationGeography represents a geography JSON object.
type TypeRepresentationGeography struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationGeography creates a new TypeRepresentationGeography instance.
func NewTypeRepresentationGeography() *TypeRepresentationGeography {
	return &TypeRepresentationGeography{
		Type: TypeRepresentationTypeGeography,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationGeography) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationGeometry represents a geography JSON object.
type TypeRepresentationGeometry struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationGeometry creates a new TypeRepresentationGeometry instance.
func NewTypeRepresentationGeometry() *TypeRepresentationGeometry {
	return &TypeRepresentationGeometry{
		Type: TypeRepresentationTypeGeometry,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationGeometry) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationBytes represent a base64-encoded bytes.
type TypeRepresentationBytes struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationBytes creates a new TypeRepresentationBytes instance.
func NewTypeRepresentationBytes() *TypeRepresentationBytes {
	return &TypeRepresentationBytes{
		Type: TypeRepresentationTypeBytes,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationBytes) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationJSON represents an arbitrary JSON.
type TypeRepresentationJSON struct {
	Type TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewTypeRepresentationJSON creates a new TypeRepresentationBytes instance.
func NewTypeRepresentationJSON() *TypeRepresentationJSON {
	return &TypeRepresentationJSON{
		Type: TypeRepresentationTypeJSON,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationJSON) Encode() TypeRepresentation {
	return map[string]any{
		"type": ty.Type,
	}
}

// TypeRepresentationEnum represents an enum type representation.
type TypeRepresentationEnum struct {
	Type  TypeRepresentationType `json:"type" yaml:"type" mapstructure:"type"`
	OneOf []string               `json:"one_of" yaml:"one_of" mapstructure:"one_of"`
}

// NewTypeRepresentationEnum creates a new TypeRepresentationEnum instance.
func NewTypeRepresentationEnum(oneOf []string) *TypeRepresentationEnum {
	return &TypeRepresentationEnum{
		Type:  TypeRepresentationTypeEnum,
		OneOf: oneOf,
	}
}

// Encode returns the raw TypeRepresentation instance.
func (ty TypeRepresentationEnum) Encode() TypeRepresentation {
	return map[string]any{
		"type":   ty.Type,
		"one_of": ty.OneOf,
	}
}

// ExtractionFunctionDefinitionType represents an extraction function definition type.
type ExtractionFunctionDefinitionType string

const (
	ExtractionFunctionDefinitionTypeNanosecond  ExtractionFunctionDefinitionType = "nanosecond"
	ExtractionFunctionDefinitionTypeMicrosecond ExtractionFunctionDefinitionType = "microsecond"
	ExtractionFunctionDefinitionTypeSecond      ExtractionFunctionDefinitionType = "second"
	ExtractionFunctionDefinitionTypeMinute      ExtractionFunctionDefinitionType = "minute"
	ExtractionFunctionDefinitionTypeHour        ExtractionFunctionDefinitionType = "hour"
	ExtractionFunctionDefinitionTypeDay         ExtractionFunctionDefinitionType = "day"
	ExtractionFunctionDefinitionTypeWeek        ExtractionFunctionDefinitionType = "week"
	ExtractionFunctionDefinitionTypeMonth       ExtractionFunctionDefinitionType = "month"
	ExtractionFunctionDefinitionTypeQuarter     ExtractionFunctionDefinitionType = "quarter"
	ExtractionFunctionDefinitionTypeYear        ExtractionFunctionDefinitionType = "year"
	ExtractionFunctionDefinitionTypeDayOfWeek   ExtractionFunctionDefinitionType = "day_of_week"
	ExtractionFunctionDefinitionTypeDayOfYear   ExtractionFunctionDefinitionType = "day_of_year"
	ExtractionFunctionDefinitionTypeCustom      ExtractionFunctionDefinitionType = "custom"
)

var enumValues_ExtractionFunctionDefinitionType = []ExtractionFunctionDefinitionType{
	ExtractionFunctionDefinitionTypeNanosecond,
	ExtractionFunctionDefinitionTypeMicrosecond,
	ExtractionFunctionDefinitionTypeSecond,
	ExtractionFunctionDefinitionTypeMinute,
	ExtractionFunctionDefinitionTypeHour,
	ExtractionFunctionDefinitionTypeDay,
	ExtractionFunctionDefinitionTypeWeek,
	ExtractionFunctionDefinitionTypeMonth,
	ExtractionFunctionDefinitionTypeYear,
	ExtractionFunctionDefinitionTypeDayOfWeek,
	ExtractionFunctionDefinitionTypeDayOfYear,
	ExtractionFunctionDefinitionTypeCustom,
}

// ParseExtractionFunctionDefinitionType parses a ordering target type argument type from string.
func ParseExtractionFunctionDefinitionType(input string) (ExtractionFunctionDefinitionType, error) {
	result := ExtractionFunctionDefinitionType(input)
	if !result.IsValid() {
		return ExtractionFunctionDefinitionType(""), fmt.Errorf("failed to parse ExtractionFunctionDefinitionType, expect one of %v, got %s", enumValues_ExtractionFunctionDefinitionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j ExtractionFunctionDefinitionType) IsValid() bool {
	return slices.Contains(enumValues_ExtractionFunctionDefinitionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExtractionFunctionDefinitionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseExtractionFunctionDefinitionType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// ExtractionFunctionDefinition represents the definition of an aggregation function on a scalar type.
type ExtractionFunctionDefinition map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExtractionFunctionDefinition) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in ExtractionFunctionDefinition: required")
	}

	var ty ExtractionFunctionDefinitionType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in ExtractionFunctionDefinition: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	switch ty {
	case ExtractionFunctionDefinitionTypeNanosecond, ExtractionFunctionDefinitionTypeMicrosecond, ExtractionFunctionDefinitionTypeSecond, ExtractionFunctionDefinitionTypeMinute, ExtractionFunctionDefinitionTypeHour, ExtractionFunctionDefinitionTypeDay, ExtractionFunctionDefinitionTypeWeek, ExtractionFunctionDefinitionTypeMonth, ExtractionFunctionDefinitionTypeYear, ExtractionFunctionDefinitionTypeDayOfWeek, ExtractionFunctionDefinitionTypeDayOfYear:
		resultType, err := unmarshalStringFromJsonMap(raw, "result_type", true)
		if err != nil {
			return fmt.Errorf("field result_type in ExtractionFunctionDefinition: %w", err)
		}

		result["result_type"] = resultType
	case ExtractionFunctionDefinitionTypeCustom:
		rawResultType, ok := raw["result_type"]
		if !ok {
			return errors.New("field result_type in ExtractionFunctionDefinition is required for custom type")
		}

		var resultType Type

		if err := json.Unmarshal(rawResultType, &resultType); err != nil {
			return fmt.Errorf("field result_type in ExtractionFunctionDefinition: %w", err)
		}

		result["result_type"] = resultType
	}

	*j = result

	return nil
}

// Type gets the type enum of the current type.
func (j ExtractionFunctionDefinition) Type() (ExtractionFunctionDefinitionType, error) {
	t, ok := j["type"]
	if !ok {
		return ExtractionFunctionDefinitionType(""), errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseExtractionFunctionDefinitionType(raw)
		if err != nil {
			return ExtractionFunctionDefinitionType(""), err
		}

		return v, nil
	case ExtractionFunctionDefinitionType:
		return raw, nil
	default:
		return ExtractionFunctionDefinitionType(""), fmt.Errorf("invalid ExtractionFunctionDefinition type: %+v", t)
	}
}

// AsNanosecond tries to convert the instance to ExtractionFunctionDefinitionNanosecond type.
func (j ExtractionFunctionDefinition) AsNanosecond() (*ExtractionFunctionDefinitionNanosecond, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeNanosecond {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeNanosecond, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionNanosecond{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsMicrosecond tries to convert the instance to ExtractionFunctionDefinitionMicrosecond type.
func (j ExtractionFunctionDefinition) AsMicrosecond() (*ExtractionFunctionDefinitionMicrosecond, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeMicrosecond {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeMicrosecond, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionMicrosecond{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsSecond tries to convert the instance to ExtractionFunctionDefinitionSecond type.
func (j ExtractionFunctionDefinition) AsSecond() (*ExtractionFunctionDefinitionSecond, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeSecond {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeSecond, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionSecond{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsMinute tries to convert the instance to ExtractionFunctionDefinitionMinute type.
func (j ExtractionFunctionDefinition) AsMinute() (*ExtractionFunctionDefinitionMinute, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeMinute {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeMinute, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionMinute{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsHour tries to convert the instance to ExtractionFunctionDefinitionHour type.
func (j ExtractionFunctionDefinition) AsHour() (*ExtractionFunctionDefinitionHour, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeHour {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeHour, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionHour{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsDay tries to convert the instance to ExtractionFunctionDefinitionDay type.
func (j ExtractionFunctionDefinition) AsDay() (*ExtractionFunctionDefinitionDay, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeDay {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeDay, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionDay{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsWeek tries to convert the instance to ExtractionFunctionDefinitionWeek type.
func (j ExtractionFunctionDefinition) AsWeek() (*ExtractionFunctionDefinitionWeek, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeWeek {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeWeek, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionWeek{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsMonth tries to convert the instance to ExtractionFunctionDefinitionMonth type.
func (j ExtractionFunctionDefinition) AsMonth() (*ExtractionFunctionDefinitionMonth, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeMonth {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeMonth, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionMonth{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsQuarter tries to convert the instance to ExtractionFunctionDefinitionQuarter type.
func (j ExtractionFunctionDefinition) AsQuarter() (*ExtractionFunctionDefinitionQuarter, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeQuarter {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeQuarter, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionQuarter{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsYear tries to convert the instance to ExtractionFunctionDefinitionYear type.
func (j ExtractionFunctionDefinition) AsYear() (*ExtractionFunctionDefinitionYear, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeYear {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeYear, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionYear{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsDayOfWeek tries to convert the instance to ExtractionFunctionDefinitionDayOfWeek type.
func (j ExtractionFunctionDefinition) AsDayOfWeek() (*ExtractionFunctionDefinitionDayOfWeek, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeDayOfWeek {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeDayOfWeek, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionDayOfWeek{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsDayOfYear tries to convert the instance to ExtractionFunctionDefinitionDayOfYear type.
func (j ExtractionFunctionDefinition) AsDayOfYear() (*ExtractionFunctionDefinitionDayOfYear, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExtractionFunctionDefinitionTypeDayOfYear {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeDayOfYear, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition result_type: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	return &ExtractionFunctionDefinitionDayOfYear{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// AsCustom tries to convert the instance to ComparisonOperatorIn type.
func (j ExtractionFunctionDefinition) AsCustom() (*ExtractionFunctionDefinitionCustom, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExtractionFunctionDefinitionTypeCustom {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type; expected: %s, got: %s", ExtractionFunctionDefinitionTypeCustom, t)
	}

	rawResultType, ok := j["result_type"]
	if !ok {
		return nil, errors.New("ExtractionFunctionDefinition.result_type is required")
	}

	resultType, ok := rawResultType.(Type)
	if !ok {
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition.result_type type; expected: Type, got: %+v", rawResultType)
	}

	return &ExtractionFunctionDefinitionCustom{
		Type:       t,
		ResultType: resultType,
	}, nil
}

// Interface tries to convert the instance to ComparisonOperatorDefinitionEncoder interface.
func (j ExtractionFunctionDefinition) Interface() ExtractionFunctionDefinitionEncoder {
	result, _ := j.InterfaceT()

	return result
}

// InterfaceT tries to convert the instance to ComparisonOperatorDefinitionEncoder interface safely with explicit error.
func (j ExtractionFunctionDefinition) InterfaceT() (ExtractionFunctionDefinitionEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case ExtractionFunctionDefinitionTypeNanosecond:
		return j.AsNanosecond()
	case ExtractionFunctionDefinitionTypeMicrosecond:
		return j.AsMicrosecond()
	case ExtractionFunctionDefinitionTypeSecond:
		return j.AsSecond()
	case ExtractionFunctionDefinitionTypeMinute:
		return j.AsMinute()
	case ExtractionFunctionDefinitionTypeHour:
		return j.AsHour()
	case ExtractionFunctionDefinitionTypeDay:
		return j.AsDay()
	case ExtractionFunctionDefinitionTypeWeek:
		return j.AsWeek()
	case ExtractionFunctionDefinitionTypeMonth:
		return j.AsMonth()
	case ExtractionFunctionDefinitionTypeQuarter:
		return j.AsQuarter()
	case ExtractionFunctionDefinitionTypeYear:
		return j.AsYear()
	case ExtractionFunctionDefinitionTypeDayOfWeek:
		return j.AsDayOfWeek()
	case ExtractionFunctionDefinitionTypeDayOfYear:
		return j.AsDayOfYear()
	case ExtractionFunctionDefinitionTypeCustom:
		return j.AsCustom()
	default:
		return nil, fmt.Errorf("invalid ExtractionFunctionDefinition type: %s", t)
	}
}

// ExtractionFunctionDefinitionEncoder abstracts the serialization interface for ExtractionFunctionDefinition.
type ExtractionFunctionDefinitionEncoder interface {
	Encode() ExtractionFunctionDefinition
}

// ExtractionFunctionDefinitionNanosecond presents a nanosecond extraction function definition.
type ExtractionFunctionDefinitionNanosecond struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionNanosecond create a new ExtractionFunctionDefinitionNanosecond instance.
func NewExtractionFunctionDefinitionNanosecond(resultType string) *ExtractionFunctionDefinitionNanosecond {
	return &ExtractionFunctionDefinitionNanosecond{
		Type:       ExtractionFunctionDefinitionTypeNanosecond,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionNanosecond) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionMicrosecond presents a microsecond extraction function definition.
type ExtractionFunctionDefinitionMicrosecond struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionMicrosecond create a new ExtractionFunctionDefinitionMicrosecond instance.
func NewExtractionFunctionDefinitionMicrosecond(resultType string) *ExtractionFunctionDefinitionMicrosecond {
	return &ExtractionFunctionDefinitionMicrosecond{
		Type:       ExtractionFunctionDefinitionTypeMicrosecond,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionMicrosecond) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionSecond presents a second extraction function definition.
type ExtractionFunctionDefinitionSecond struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionSecond create a new ExtractionFunctionDefinitionMicrosecond instance.
func NewExtractionFunctionDefinitionSecond(resultType string) *ExtractionFunctionDefinitionSecond {
	return &ExtractionFunctionDefinitionSecond{
		Type:       ExtractionFunctionDefinitionTypeSecond,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionSecond) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionMinute presents a minute extraction function definition.
type ExtractionFunctionDefinitionMinute struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionMinute create a new ExtractionFunctionDefinitionMinute instance.
func NewExtractionFunctionDefinitionMinute(resultType string) *ExtractionFunctionDefinitionMinute {
	return &ExtractionFunctionDefinitionMinute{
		Type:       ExtractionFunctionDefinitionTypeMinute,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionMinute) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionHour presents an hour extraction function definition.
type ExtractionFunctionDefinitionHour struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionHour create a new ExtractionFunctionDefinitionHour instance.
func NewExtractionFunctionDefinitionHour(resultType string) *ExtractionFunctionDefinitionHour {
	return &ExtractionFunctionDefinitionHour{
		Type:       ExtractionFunctionDefinitionTypeHour,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionHour) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionDay presents a day extraction function definition.
type ExtractionFunctionDefinitionDay struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionDay create a new ExtractionFunctionDefinitionDay instance.
func NewExtractionFunctionDefinitionDay(resultType string) *ExtractionFunctionDefinitionDay {
	return &ExtractionFunctionDefinitionDay{
		Type:       ExtractionFunctionDefinitionTypeDay,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionDay) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionWeek presents a week extraction function definition.
type ExtractionFunctionDefinitionWeek struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionWeek create a new ExtractionFunctionDefinitionWeek instance.
func NewExtractionFunctionDefinitionWeek(resultType string) *ExtractionFunctionDefinitionWeek {
	return &ExtractionFunctionDefinitionWeek{
		Type:       ExtractionFunctionDefinitionTypeWeek,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionWeek) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionMonth presents a month extraction function definition.
type ExtractionFunctionDefinitionMonth struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionMonth create a new ExtractionFunctionDefinitionMonth instance.
func NewExtractionFunctionDefinitionMonth(resultType string) *ExtractionFunctionDefinitionMonth {
	return &ExtractionFunctionDefinitionMonth{
		Type:       ExtractionFunctionDefinitionTypeMonth,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionMonth) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionQuarter presents a quarter extraction function definition.
type ExtractionFunctionDefinitionQuarter struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionQuarter create a new ExtractionFunctionDefinitionQuarter instance.
func NewExtractionFunctionDefinitionQuarter(resultType string) *ExtractionFunctionDefinitionQuarter {
	return &ExtractionFunctionDefinitionQuarter{
		Type:       ExtractionFunctionDefinitionTypeQuarter,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionQuarter) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionYear presents a year extraction function definition.
type ExtractionFunctionDefinitionYear struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionYear create a new ExtractionFunctionDefinitionYear instance.
func NewExtractionFunctionDefinitionYear(resultType string) *ExtractionFunctionDefinitionYear {
	return &ExtractionFunctionDefinitionYear{
		Type:       ExtractionFunctionDefinitionTypeYear,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionYear) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionDayOfWeek presents a day-of-week extraction function definition.
type ExtractionFunctionDefinitionDayOfWeek struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionDayOfWeek create a new ExtractionFunctionDefinitionDayOfWeek instance.
func NewExtractionFunctionDefinitionDayOfWeek(resultType string) *ExtractionFunctionDefinitionDayOfWeek {
	return &ExtractionFunctionDefinitionDayOfWeek{
		Type:       ExtractionFunctionDefinitionTypeDayOfWeek,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionDayOfWeek) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionDayOfYear presents a day-of-year extraction function definition.
type ExtractionFunctionDefinitionDayOfYear struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType string                           `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionDayOfYear create a new ExtractionFunctionDefinitionDayOfYear instance.
func NewExtractionFunctionDefinitionDayOfYear(resultType string) *ExtractionFunctionDefinitionDayOfYear {
	return &ExtractionFunctionDefinitionDayOfYear{
		Type:       ExtractionFunctionDefinitionTypeDayOfYear,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionDayOfYear) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// ExtractionFunctionDefinitionCustom presents a custom extraction function definition.
type ExtractionFunctionDefinitionCustom struct {
	Type       ExtractionFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	ResultType Type                             `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewExtractionFunctionDefinitionCustom create a new ExtractionFunctionDefinitionCustom instance.
func NewExtractionFunctionDefinitionCustom(resultType TypeEncoder) *ExtractionFunctionDefinitionCustom {
	return &ExtractionFunctionDefinitionCustom{
		Type:       ExtractionFunctionDefinitionTypeCustom,
		ResultType: resultType.Encode(),
	}
}

// Encode converts the instance to raw ExtractionFunctionDefinition.
func (efd ExtractionFunctionDefinitionCustom) Encode() ExtractionFunctionDefinition {
	return ExtractionFunctionDefinition{
		"type":        efd.Type,
		"result_type": efd.ResultType,
	}
}

// AggregateFunctionDefinitionType represents a type of AggregateFunctionDefinition
type AggregateFunctionDefinitionType string

const (
	AggregateFunctionDefinitionTypeMin     AggregateFunctionDefinitionType = "min"
	AggregateFunctionDefinitionTypeMax     AggregateFunctionDefinitionType = "max"
	AggregateFunctionDefinitionTypeSum     AggregateFunctionDefinitionType = "sum"
	AggregateFunctionDefinitionTypeAverage AggregateFunctionDefinitionType = "average"
	AggregateFunctionDefinitionTypeCustom  AggregateFunctionDefinitionType = "custom"
)

var enumValues_AggregateFunctionDefinitionType = []AggregateFunctionDefinitionType{
	AggregateFunctionDefinitionTypeMin,
	AggregateFunctionDefinitionTypeMax,
	AggregateFunctionDefinitionTypeSum,
	AggregateFunctionDefinitionTypeAverage,
	AggregateFunctionDefinitionTypeCustom,
}

// ParseAggregateFunctionDefinitionType parses a AggregateFunctionDefinitionType from string.
func ParseAggregateFunctionDefinitionType(input string) (AggregateFunctionDefinitionType, error) {
	result := AggregateFunctionDefinitionType(input)
	if !result.IsValid() {
		return AggregateFunctionDefinitionType(""), fmt.Errorf("failed to parse AggregateFunctionDefinitionType, expect one of %v, got %s", enumValues_AggregateFunctionDefinitionType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
func (j AggregateFunctionDefinitionType) IsValid() bool {
	return slices.Contains(enumValues_AggregateFunctionDefinitionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateFunctionDefinitionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseAggregateFunctionDefinitionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// AggregateFunctionDefinitionEncoder abstracts a generic interface of AggregateFunctionDefinition
type AggregateFunctionDefinitionEncoder interface {
	Encode() AggregateFunctionDefinition
}

// AggregateFunctionDefinition represents the definition of an aggregation function on a scalar type
type AggregateFunctionDefinition map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateFunctionDefinition) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty AggregateFunctionDefinitionType

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in AggregateFunctionDefinition: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in AggregateFunctionDefinition: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case AggregateFunctionDefinitionTypeSum, AggregateFunctionDefinitionTypeAverage:
		rawResultType, ok := raw["result_type"]
		if !ok || len(rawResultType) == 0 {
			return errors.New("field result_type in AggregateFunctionDefinition: required")
		}

		var resultType string
		if err := json.Unmarshal(rawResultType, &resultType); err != nil {
			return fmt.Errorf("field result_type in AggregateFunctionDefinition: %w", err)
		}

		results["result_type"] = resultType
	case AggregateFunctionDefinitionTypeCustom:
		rawResultType, ok := raw["result_type"]
		if !ok || len(rawResultType) == 0 {
			return errors.New("field result_type in AggregateFunctionDefinition: required")
		}

		var resultType Type
		if err := resultType.UnmarshalJSON(rawResultType); err != nil {
			return fmt.Errorf("field result_type in AggregateFunctionDefinition: %w", err)
		}

		results["result_type"] = resultType
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
func (j AggregateFunctionDefinition) Type() (AggregateFunctionDefinitionType, error) {
	t, ok := j["type"]
	if !ok {
		return AggregateFunctionDefinitionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseAggregateFunctionDefinitionType(raw)
		if err != nil {
			return AggregateFunctionDefinitionType(""), err
		}
		return v, nil
	case AggregateFunctionDefinitionType:
		return raw, nil
	default:
		return AggregateFunctionDefinitionType(""), fmt.Errorf("invalid AggregateFunctionDefinition type: %+v", t)
	}
}

// AsMin tries to convert the current type to AggregateFunctionDefinitionMin.
func (j AggregateFunctionDefinition) AsMin() (*AggregateFunctionDefinitionMin, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeMin {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeMin, t)
	}

	result := &AggregateFunctionDefinitionMin{
		Type: t,
	}

	return result, nil
}

// AsMax tries to convert the current type to AggregateFunctionDefinitionMax.
func (j AggregateFunctionDefinition) AsMax() (*AggregateFunctionDefinitionMax, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeMax {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeMax, t)
	}

	result := &AggregateFunctionDefinitionMax{
		Type: t,
	}

	return result, nil
}

// AsSum tries to convert the current type to AggregateFunctionDefinitionSum.
func (j AggregateFunctionDefinition) AsSum() (*AggregateFunctionDefinitionSum, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeSum {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeSum, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("field result_type in AggregateFunctionDefinitionSum: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionSum: required")
	}
	result := &AggregateFunctionDefinitionSum{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// AsAverage tries to convert the current type to AggregateFunctionDefinitionAverage.
func (j AggregateFunctionDefinition) AsAverage() (*AggregateFunctionDefinitionAverage, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeAverage {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeAverage, t)
	}

	resultType, err := getStringValueByKey(j, "result_type")
	if err != nil {
		return nil, fmt.Errorf("field result_type in AggregateFunctionDefinitionAverage: %w", err)
	}

	if resultType == "" {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionAverage: required")
	}

	result := &AggregateFunctionDefinitionAverage{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// AsCustom tries to convert the current type to AggregateFunctionDefinitionCustom.
func (j AggregateFunctionDefinition) AsCustom() (*AggregateFunctionDefinitionCustom, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeCustom {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeCustom, t)
	}

	rawResultType, ok := j["result_type"]
	if !ok {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionCustom: required")
	}

	resultType, ok := rawResultType.(Type)
	if !ok {
		return nil, fmt.Errorf("invalid result_type in AggregateFunctionDefinitionCustom, expected Type, got %v", rawResultType)
	}
	result := &AggregateFunctionDefinitionCustom{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j AggregateFunctionDefinition) Interface() AggregateFunctionDefinitionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j AggregateFunctionDefinition) InterfaceT() (AggregateFunctionDefinitionEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case AggregateFunctionDefinitionTypeMin:
		return j.AsMin()
	case AggregateFunctionDefinitionTypeMax:
		return j.AsMax()
	case AggregateFunctionDefinitionTypeSum:
		return j.AsSum()
	case AggregateFunctionDefinitionTypeAverage:
		return j.AsAverage()
	case AggregateFunctionDefinitionTypeCustom:
		return j.AsCustom()
	default:
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type: %s", ty)
	}
}

// AggregateFunctionDefinitionMin represents a min aggregate function definition
type AggregateFunctionDefinitionMin struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewAggregateFunctionDefinitionMin creates an AggregateFunctionDefinitionMin instance.
func NewAggregateFunctionDefinitionMin() *AggregateFunctionDefinitionMin {
	return &AggregateFunctionDefinitionMin{
		Type: AggregateFunctionDefinitionTypeMin,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionMin) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type": j.Type,
	}
	return result
}

// AggregateFunctionDefinitionMax represents a max aggregate function definition
type AggregateFunctionDefinitionMax struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewAggregateFunctionDefinitionMax creates an AggregateFunctionDefinitionMax instance.
func NewAggregateFunctionDefinitionMax() *AggregateFunctionDefinitionMax {
	return &AggregateFunctionDefinitionMax{
		Type: AggregateFunctionDefinitionTypeMax,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionMax) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type": j.Type,
	}
	return result
}

// AggregateFunctionDefinitionAverage represents an average aggregate function definition
type AggregateFunctionDefinitionAverage struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar type of the result of this function, which should have the type representation Float64
	ResultType string `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionAverage creates an AggregateFunctionDefinitionAverage instance.
func NewAggregateFunctionDefinitionAverage(resultType string) *AggregateFunctionDefinitionAverage {
	return &AggregateFunctionDefinitionAverage{
		Type:       AggregateFunctionDefinitionTypeAverage,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionAverage) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
}

// AggregateFunctionDefinitionSum represents a sum aggregate function definition
type AggregateFunctionDefinitionSum struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar type of the result of this function, which should have one of the type representations Int64 or Float64, depending on whether this function is defined on a scalar type with an integer or floating-point representation, respectively.
	ResultType string `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionSum creates an AggregateFunctionDefinitionSum instance.
func NewAggregateFunctionDefinitionSum(resultType string) *AggregateFunctionDefinitionSum {
	return &AggregateFunctionDefinitionSum{
		Type:       AggregateFunctionDefinitionTypeSum,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionSum) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
}

// AggregateFunctionDefinitionCustom represents a sum aggregate function definition
type AggregateFunctionDefinitionCustom struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar or object type of the result of this function.
	ResultType Type `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionCustom creates an AggregateFunctionDefinitionCustom instance.
func NewAggregateFunctionDefinitionCustom(resultType Type) *AggregateFunctionDefinitionCustom {
	return &AggregateFunctionDefinitionCustom{
		Type:       AggregateFunctionDefinitionTypeCustom,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionCustom) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
}

// ComparisonOperatorDefinitionType represents a binary comparison operator type enum.
type ComparisonOperatorDefinitionType string

const (
	ComparisonOperatorDefinitionTypeEqual  ComparisonOperatorDefinitionType = "equal"
	ComparisonOperatorDefinitionTypeIn     ComparisonOperatorDefinitionType = "in"
	ComparisonOperatorDefinitionTypeCustom ComparisonOperatorDefinitionType = "custom"
)

var enumValues_ComparisonOperatorDefinitionType = []ComparisonOperatorDefinitionType{
	ComparisonOperatorDefinitionTypeEqual,
	ComparisonOperatorDefinitionTypeIn,
	ComparisonOperatorDefinitionTypeCustom,
}

// ParseComparisonOperatorDefinitionType parses a type of a comparison operator definition.
func ParseComparisonOperatorDefinitionType(input string) (ComparisonOperatorDefinitionType, error) {
	result := ComparisonOperatorDefinitionType(input)
	if !result.IsValid() {
		return ComparisonOperatorDefinitionType(""), fmt.Errorf("failed to parse ComparisonOperatorDefinitionType, expect one of %v, got %s", enumValues_ComparisonOperatorDefinitionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j ComparisonOperatorDefinitionType) IsValid() bool {
	return slices.Contains(enumValues_ComparisonOperatorDefinitionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonOperatorDefinitionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseComparisonOperatorDefinitionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// ComparisonOperatorDefinition the definition of a comparison operator on a scalar type.
type ComparisonOperatorDefinition map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonOperatorDefinition) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in ComparisonOperatorDefinition: required")
	}

	var ty ComparisonOperatorDefinitionType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in ComparisonOperatorDefinition: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case ComparisonOperatorDefinitionTypeEqual:
	case ComparisonOperatorDefinitionTypeIn:
	case ComparisonOperatorDefinitionTypeCustom:
		rawArgumentType, ok := raw["argument_type"]
		if !ok {
			return errors.New("field argument_type in ComparisonOperatorDefinition is required for custom type")
		}
		var argumentType Type
		if err := json.Unmarshal(rawArgumentType, &argumentType); err != nil {
			return fmt.Errorf("field argument_type in ComparisonOperatorDefinition: %w", err)
		}
		result["argument_type"] = argumentType
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type.
func (j ComparisonOperatorDefinition) Type() (ComparisonOperatorDefinitionType, error) {
	t, ok := j["type"]
	if !ok {
		return ComparisonOperatorDefinitionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseComparisonOperatorDefinitionType(raw)
		if err != nil {
			return ComparisonOperatorDefinitionType(""), err
		}
		return v, nil
	case ComparisonOperatorDefinitionType:
		return raw, nil
	default:
		return ComparisonOperatorDefinitionType(""), fmt.Errorf("invalid ComparisonOperatorDefinition type: %+v", t)
	}
}

// AsEqual tries to convert the instance to ComparisonOperatorEqual type.
func (j ComparisonOperatorDefinition) AsEqual() (*ComparisonOperatorEqual, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ComparisonOperatorDefinitionTypeEqual {
		return nil, fmt.Errorf("invalid ComparisonOperatorDefinition type; expected: %s, got: %s", ComparisonOperatorDefinitionTypeEqual, t)
	}

	return &ComparisonOperatorEqual{
		Type: t,
	}, nil
}

// AsIn tries to convert the instance to ComparisonOperatorIn type.
func (j ComparisonOperatorDefinition) AsIn() (*ComparisonOperatorIn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ComparisonOperatorDefinitionTypeIn {
		return nil, fmt.Errorf("invalid ComparisonOperatorDefinition type; expected: %s, got: %s", ComparisonOperatorDefinitionTypeIn, t)
	}

	return &ComparisonOperatorIn{
		Type: t,
	}, nil
}

// AsCustom tries to convert the instance to ComparisonOperatorIn type.
func (j ComparisonOperatorDefinition) AsCustom() (*ComparisonOperatorCustom, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ComparisonOperatorDefinitionTypeCustom {
		return nil, fmt.Errorf("invalid ComparisonOperatorDefinition type; expected: %s, got: %s", ComparisonOperatorDefinitionTypeCustom, t)
	}

	rawArg, ok := j["argument_type"]
	if !ok {
		return nil, errors.New("ComparisonOperatorCustom.argument_type is required")
	}

	arg, ok := rawArg.(Type)
	if !ok {
		return nil, fmt.Errorf("invalid ComparisonOperatorCustom.argument_type type; expected: Type, got: %+v", rawArg)
	}

	return &ComparisonOperatorCustom{
		Type:         t,
		ArgumentType: arg,
	}, nil
}

// Interface tries to convert the instance to ComparisonOperatorDefinitionEncoder interface.
func (j ComparisonOperatorDefinition) Interface() ComparisonOperatorDefinitionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to ComparisonOperatorDefinitionEncoder interface safely with explicit error.
func (j ComparisonOperatorDefinition) InterfaceT() (ComparisonOperatorDefinitionEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case ComparisonOperatorDefinitionTypeEqual:
		return j.AsEqual()
	case ComparisonOperatorDefinitionTypeIn:
		return j.AsIn()
	case ComparisonOperatorDefinitionTypeCustom:
		return j.AsCustom()
	default:
		return nil, fmt.Errorf("invalid ComparisonOperatorDefinition type: %s", t)
	}
}

// ComparisonOperatorDefinitionEncoder abstracts the serialization interface for ComparisonOperatorDefinition.
type ComparisonOperatorDefinitionEncoder interface {
	Encode() ComparisonOperatorDefinition
}

// ComparisonOperatorEqual presents an equal comparison operator.
type ComparisonOperatorEqual struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewComparisonOperatorEqual create a new ComparisonOperatorEqual instance.
func NewComparisonOperatorEqual() *ComparisonOperatorEqual {
	return &ComparisonOperatorEqual{
		Type: ComparisonOperatorDefinitionTypeEqual,
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition.
func (ob ComparisonOperatorEqual) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type": ob.Type,
	}
}

// ComparisonOperatorIn presents an in comparison operator.
type ComparisonOperatorIn struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewComparisonOperatorIn create a new ComparisonOperatorIn instance.
func NewComparisonOperatorIn() *ComparisonOperatorIn {
	return &ComparisonOperatorIn{
		Type: ComparisonOperatorDefinitionTypeIn,
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition.
func (ob ComparisonOperatorIn) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type": ob.Type,
	}
}

// ComparisonOperatorCustom presents a custom comparison operator.
type ComparisonOperatorCustom struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The type of the argument to this operator
	ArgumentType Type `json:"argument_type" yaml:"argument_type" mapstructure:"argument_type"`
}

// NewComparisonOperatorCustom create a new ComparisonOperatorCustom instance.
func NewComparisonOperatorCustom(argumentType TypeEncoder) *ComparisonOperatorCustom {
	return &ComparisonOperatorCustom{
		Type:         ComparisonOperatorDefinitionTypeCustom,
		ArgumentType: argumentType.Encode(),
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition.
func (ob ComparisonOperatorCustom) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type":          ob.Type,
		"argument_type": ob.ArgumentType,
	}
}
