package schema

import (
	"encoding/json"
	"fmt"
	"slices"
)

// CastTypeEnum represents a cast enum.
type CastTypeEnum string

const (
	CastTypeEnumBoolean    CastTypeEnum = "boolean"
	CastTypeEnumUTF8       CastTypeEnum = "utf8"
	CastTypeEnumInt8       CastTypeEnum = "int8"
	CastTypeEnumInt16      CastTypeEnum = "int16"
	CastTypeEnumInt32      CastTypeEnum = "int32"
	CastTypeEnumInt64      CastTypeEnum = "int64"
	CastTypeEnumUint8      CastTypeEnum = "uint8"
	CastTypeEnumUint16     CastTypeEnum = "uint16"
	CastTypeEnumUint32     CastTypeEnum = "uint32"
	CastTypeEnumUint64     CastTypeEnum = "uint64"
	CastTypeEnumFloat32    CastTypeEnum = "float32"
	CastTypeEnumFloat64    CastTypeEnum = "float64"
	CastTypeEnumDecimal128 CastTypeEnum = "decimal128"
	CastTypeEnumDecimal256 CastTypeEnum = "decimal256"
	CastTypeEnumDate       CastTypeEnum = "date"
	CastTypeEnumTime       CastTypeEnum = "time"
	CastTypeEnumTimestamp  CastTypeEnum = "timestamp"
	CastTypeEnumDuration   CastTypeEnum = "duration"
	CastTypeEnumInterval   CastTypeEnum = "interval"
)

var enumValues_CastTypeEnum = []CastTypeEnum{
	CastTypeEnumBoolean,
	CastTypeEnumUTF8,
	CastTypeEnumInt8,
	CastTypeEnumInt16,
	CastTypeEnumInt32,
	CastTypeEnumInt64,
	CastTypeEnumUint8,
	CastTypeEnumUint16,
	CastTypeEnumUint32,
	CastTypeEnumUint64,
	CastTypeEnumFloat32,
	CastTypeEnumFloat64,
	CastTypeEnumDecimal128,
	CastTypeEnumDecimal256,
	CastTypeEnumDate,
	CastTypeEnumTime,
	CastTypeEnumTimestamp,
	CastTypeEnumDuration,
	CastTypeEnumInterval,
}

// ParseCastTypeEnum parses a cast type enum from string.
func ParseCastTypeEnum(input string) (CastTypeEnum, error) {
	result := CastTypeEnum(input)
	if !result.IsValid() {
		return CastTypeEnum(
				"",
			), fmt.Errorf(
				"failed to parse CastTypeEnum, expect one of %v, got %s",
				enumValues_CastTypeEnum,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j CastTypeEnum) IsValid() bool {
	return slices.Contains(enumValues_CastTypeEnum, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CastTypeEnum) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseCastTypeEnum(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// CastType is provided by reference to a relational expression.
type CastType struct {
	inner CastTypeInner
}

// NewRelationalExpression creates a RelationalExpression instance.
func NewCastType[T CastTypeInner](inner T) CastType {
	return CastType{
		inner: inner,
	}
}

// CastTypeInner abstracts the interface for CastType.
type CastTypeInner interface {
	Type() CastTypeEnum
	ToMap() map[string]any
	Wrap() CastType
}

// IsEmpty checks if the inner type is empty.
func (j CastType) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j CastType) Type() CastTypeEnum {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// Interface tries to convert the instance to AggregateInner interface.
func (j CastType) Interface() CastTypeInner {
	return j.inner
}

// MarshalJSON implements json.Marshaler interface.
func (j CastType) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CastType) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseCastTypeEnum(rawType.Type)
	if err != nil {
		return fmt.Errorf("failed to decode CastType: %w", err)
	}

	switch ty {
	case CastTypeEnumBoolean:
		j.inner = NewCastTypeBoolean()
	case CastTypeEnumDate:
		j.inner = NewCastTypeDate()
	case CastTypeEnumDecimal128:
		var result CastTypeDecimal128

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode CastTypeDecimal128: %w", err)
		}

		j.inner = &result
	case CastTypeEnumDecimal256:
		var result CastTypeDecimal256

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode CastTypeDecimal256: %w", err)
		}

		j.inner = &result
	case CastTypeEnumDuration:
		j.inner = NewCastTypeDuration()
	case CastTypeEnumFloat32:
		j.inner = NewCastTypeFloat32()
	case CastTypeEnumFloat64:
		j.inner = NewCastTypeFloat64()
	case CastTypeEnumInt8:
		j.inner = NewCastTypeInt8()
	case CastTypeEnumInt16:
		j.inner = NewCastTypeInt16()
	case CastTypeEnumInt32:
		j.inner = NewCastTypeInt32()
	case CastTypeEnumInt64:
		j.inner = NewCastTypeInt64()
	case CastTypeEnumInterval:
		j.inner = NewCastTypeInterval()
	case CastTypeEnumTime:
		j.inner = NewCastTypeTime()
	case CastTypeEnumTimestamp:
		j.inner = NewCastTypeTimestamp()
	case CastTypeEnumUTF8:
		j.inner = NewCastTypeUTF8()
	case CastTypeEnumUint8:
		j.inner = NewCastTypeUint8()
	case CastTypeEnumUint16:
		j.inner = NewCastTypeUint16()
	case CastTypeEnumUint32:
		j.inner = NewCastTypeUint32()
	case CastTypeEnumUint64:
		j.inner = NewCastTypeUint64()
	default:
		return fmt.Errorf("unsupported cast type: %s", ty)
	}

	return nil
}

// CastTypeBoolean represents a CastType with the boolean type.
type CastTypeBoolean struct{}

// NewCastTypeBoolean creates a CastTypeBoolean instance.
func NewCastTypeBoolean() *CastTypeBoolean {
	return &CastTypeBoolean{}
}

// Type return the type name of the instance.
func (j CastTypeBoolean) Type() CastTypeEnum {
	return CastTypeEnumBoolean
}

// ToMap converts the instance to raw map.
func (j CastTypeBoolean) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeBoolean) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeUTF8 represents a CastType with the utf8 type.
type CastTypeUTF8 struct{}

// NewCastTypeUTF8 creates a CastTypeUTF8 instance.
func NewCastTypeUTF8() *CastTypeUTF8 {
	return &CastTypeUTF8{}
}

// Type return the type name of the instance.
func (j CastTypeUTF8) Type() CastTypeEnum {
	return CastTypeEnumUTF8
}

// ToMap converts the instance to raw map.
func (j CastTypeUTF8) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeUTF8) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeInt8 represents a CastType with the int8 type.
type CastTypeInt8 struct{}

// NewCastTypeInt8 creates a CastTypeInt8 instance.
func NewCastTypeInt8() *CastTypeInt8 {
	return &CastTypeInt8{}
}

// Type return the type name of the instance.
func (j CastTypeInt8) Type() CastTypeEnum {
	return CastTypeEnumInt8
}

// ToMap converts the instance to raw map.
func (j CastTypeInt8) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeInt8) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeInt16 represents a CastType with the int16 type.
type CastTypeInt16 struct{}

// NewCastTypeInt16 creates a CastTypeInt16 instance.
func NewCastTypeInt16() *CastTypeInt16 {
	return &CastTypeInt16{}
}

// Type return the type name of the instance.
func (j CastTypeInt16) Type() CastTypeEnum {
	return CastTypeEnumInt16
}

// ToMap converts the instance to raw map.
func (j CastTypeInt16) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeInt16) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeInt32 represents a CastType with the int32 type.
type CastTypeInt32 struct{}

// NewCastTypeInt32 creates a CastTypeInt32 instance.
func NewCastTypeInt32() *CastTypeInt32 {
	return &CastTypeInt32{}
}

// Type return the type name of the instance.
func (j CastTypeInt32) Type() CastTypeEnum {
	return CastTypeEnumInt32
}

// ToMap converts the instance to raw map.
func (j CastTypeInt32) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeInt32) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeInt64 represents a CastType with the int64 type.
type CastTypeInt64 struct{}

// NewCastTypeInt64 creates a CastTypeInt64 instance.
func NewCastTypeInt64() *CastTypeInt64 {
	return &CastTypeInt64{}
}

// Type return the type name of the instance.
func (j CastTypeInt64) Type() CastTypeEnum {
	return CastTypeEnumInt64
}

// ToMap converts the instance to raw map.
func (j CastTypeInt64) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeInt64) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeUint8 represents a CastType with the uint8 type.
type CastTypeUint8 struct{}

// NewCastTypeUint8 creates a CastTypeUint8 instance.
func NewCastTypeUint8() *CastTypeUint8 {
	return &CastTypeUint8{}
}

// Type return the type name of the instance.
func (j CastTypeUint8) Type() CastTypeEnum {
	return CastTypeEnumUint8
}

// ToMap converts the instance to raw map.
func (j CastTypeUint8) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeUint8) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeUint16 represents a CastType with the uint16 type.
type CastTypeUint16 struct{}

// NewCastTypeUint16 creates a CastTypeUint16 instance.
func NewCastTypeUint16() *CastTypeUint16 {
	return &CastTypeUint16{}
}

// Type return the type name of the instance.
func (j CastTypeUint16) Type() CastTypeEnum {
	return CastTypeEnumUint16
}

// ToMap converts the instance to raw map.
func (j CastTypeUint16) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeUint16) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeUint32 represents a CastType with the uint32 type.
type CastTypeUint32 struct{}

// NewCastTypeUint32 creates a CastTypeUint32 instance.
func NewCastTypeUint32() *CastTypeUint32 {
	return &CastTypeUint32{}
}

// Type return the type name of the instance.
func (j CastTypeUint32) Type() CastTypeEnum {
	return CastTypeEnumUint32
}

// ToMap converts the instance to raw map.
func (j CastTypeUint32) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeUint32) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeUint64 represents a CastType with the uint64 type.
type CastTypeUint64 struct{}

// NewCastTypeUint64 creates a CastTypeUint64 instance.
func NewCastTypeUint64() *CastTypeUint64 {
	return &CastTypeUint64{}
}

// Type return the type name of the instance.
func (j CastTypeUint64) Type() CastTypeEnum {
	return CastTypeEnumUint64
}

// ToMap converts the instance to raw map.
func (j CastTypeUint64) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeUint64) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeFloat32 represents a CastType with the float32 type.
type CastTypeFloat32 struct{}

// NewCastTypeFloat32 creates a CastTypeFloat32 instance.
func NewCastTypeFloat32() *CastTypeFloat32 {
	return &CastTypeFloat32{}
}

// Type return the type name of the instance.
func (j CastTypeFloat32) Type() CastTypeEnum {
	return CastTypeEnumFloat32
}

// ToMap converts the instance to raw map.
func (j CastTypeFloat32) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeFloat32) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeFloat64 represents a CastType with the float64 type.
type CastTypeFloat64 struct{}

// NewCastTypeFloat64 creates a CastTypeFloat64 instance.
func NewCastTypeFloat64() *CastTypeFloat64 {
	return &CastTypeFloat64{}
}

// Type return the type name of the instance.
func (j CastTypeFloat64) Type() CastTypeEnum {
	return CastTypeEnumFloat64
}

// ToMap converts the instance to raw map.
func (j CastTypeFloat64) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeFloat64) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeDecimal128 represents a CastType with unsigned 128-bit decimal.
type CastTypeDecimal128 struct {
	Scale int8  `json:"scale" mapstructure:"scale" yaml:"scale"`
	Spec  uint8 `json:"prec"  mapstructure:"prec"  yaml:"prec"`
}

// NewCastTypeDecimal128 creates a CastTypeDecimal128 instance.
func NewCastTypeDecimal128(scale int8, spec uint8) *CastTypeDecimal128 {
	return &CastTypeDecimal128{
		Scale: scale,
		Spec:  spec,
	}
}

// Type return the type name of the instance.
func (j CastTypeDecimal128) Type() CastTypeEnum {
	return CastTypeEnumDecimal128
}

// ToMap converts the instance to raw map.
func (j CastTypeDecimal128) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"scale": j.Scale,
		"prec":  j.Spec,
	}
}

// Encode returns the relation wrapper.
func (j CastTypeDecimal128) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeDecimal256 represents a CastType with unsigned 256-bit decimal.
type CastTypeDecimal256 struct {
	Scale int8  `json:"scale" mapstructure:"scale" yaml:"scale"`
	Spec  uint8 `json:"prec"  mapstructure:"prec"  yaml:"prec"`
}

// NewCastTypeDecimal256 creates a CastTypeDecimal256 instance.
func NewCastTypeDecimal256(scale int8, spec uint8) *CastTypeDecimal256 {
	return &CastTypeDecimal256{
		Scale: scale,
		Spec:  spec,
	}
}

// Type return the type name of the instance.
func (j CastTypeDecimal256) Type() CastTypeEnum {
	return CastTypeEnumDecimal256
}

// ToMap converts the instance to raw map.
func (j CastTypeDecimal256) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"scale": j.Scale,
		"prec":  j.Spec,
	}
}

// Encode returns the relation wrapper.
func (j CastTypeDecimal256) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeDate represents a CastType with the date type.
type CastTypeDate struct{}

// NewCastTypeDate creates a CastTypeDate instance.
func NewCastTypeDate() *CastTypeDate {
	return &CastTypeDate{}
}

// Type return the type name of the instance.
func (j CastTypeDate) Type() CastTypeEnum {
	return CastTypeEnumDate
}

// ToMap converts the instance to raw map.
func (j CastTypeDate) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeDate) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeTime represents a CastType with the time type.
type CastTypeTime struct{}

// NewCastTypeTime creates a CastTypeTime instance.
func NewCastTypeTime() *CastTypeTime {
	return &CastTypeTime{}
}

// Type return the type name of the instance.
func (j CastTypeTime) Type() CastTypeEnum {
	return CastTypeEnumTime
}

// ToMap converts the instance to raw map.
func (j CastTypeTime) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeTime) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeTimestamp represents a CastType with the timestamp type.
type CastTypeTimestamp struct{}

// NewCastTypeTimestamp creates a CastTypeTimestamp instance.
func NewCastTypeTimestamp() *CastTypeTimestamp {
	return &CastTypeTimestamp{}
}

// Type return the type name of the instance.
func (j CastTypeTimestamp) Type() CastTypeEnum {
	return CastTypeEnumTimestamp
}

// ToMap converts the instance to raw map.
func (j CastTypeTimestamp) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeTimestamp) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeDuration represents a CastType with the duration type.
type CastTypeDuration struct{}

// NewCastTypeDuration creates a CastTypeDuration instance.
func NewCastTypeDuration() *CastTypeDuration {
	return &CastTypeDuration{}
}

// Type return the type name of the instance.
func (j CastTypeDuration) Type() CastTypeEnum {
	return CastTypeEnumDuration
}

// ToMap converts the instance to raw map.
func (j CastTypeDuration) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeDuration) Wrap() CastType {
	return NewCastType(&j)
}

// CastTypeInterval represents a CastType with the interval type.
type CastTypeInterval struct{}

// NewCastTypeInterval creates a CastTypeInterval instance.
func NewCastTypeInterval() *CastTypeInterval {
	return &CastTypeInterval{}
}

// Type return the type name of the instance.
func (j CastTypeInterval) Type() CastTypeEnum {
	return CastTypeEnumInterval
}

// ToMap converts the instance to raw map.
func (j CastTypeInterval) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j CastTypeInterval) Wrap() CastType {
	return NewCastType(&j)
}
