package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

var (
	errTypeRequired = errors.New("type field is required")
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

	if !slices.Contains(enumValues_Type, result) {
		return TypeEnum(""), fmt.Errorf("failed to parse TypeEnum, expect one of %v, got %s", enumValues_Type, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j TypeEnum) IsValid() bool {
	return slices.Contains(enumValues_Type, j)
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

// NewNullableArrayType creates a new NullableType instance with underlying array type
func NewNullableArrayType(elementType TypeEncoder) *NullableType {
	return &NullableType{
		Type:           TypeNullable,
		UnderlyingType: elementType.Encode(),
	}
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

// ArgumentType represents an argument type enum
type ArgumentType string

const (
	ArgumentTypeLiteral  ArgumentType = "literal"
	ArgumentTypeVariable ArgumentType = "variable"
)

var enumValues_ArgumentType = []ArgumentType{
	ArgumentTypeLiteral,
	ArgumentTypeVariable,
}

// ParseArgumentType parses an argument type from string
func ParseArgumentType(input string) (ArgumentType, error) {
	result := ArgumentType(input)
	if !slices.Contains(enumValues_ArgumentType, result) {
		return ArgumentType(""), fmt.Errorf("failed to parse ArgumentType, expect one of %v, got %s", enumValues_ArgumentType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid
func (j ArgumentType) IsValid() bool {
	return slices.Contains(enumValues_ArgumentType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArgumentType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseArgumentType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// Argument is provided by reference to a variable or as a literal value
type Argument struct {
	Type  ArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Name  string       `json:"name" yaml:"name" mapstructure:"name"`
	Value any          `json:"value" yaml:"value" mapstructure:"value"`
}

// MarshalJSON implements json.Marshaler.
func (j Argument) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"type": j.Type,
	}
	switch j.Type {
	case ArgumentTypeLiteral:
		result["value"] = j.Value
	case ArgumentTypeVariable:
		result["name"] = j.Name
	}
	return json.Marshal(result)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Argument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawArgumentType := getStringValueByKey(raw, "type")
	if rawArgumentType == "" {
		return errors.New("field type in Argument: required")
	}

	argumentType, err := ParseArgumentType(rawArgumentType)
	if err != nil {
		return fmt.Errorf("field type in Argument: %s", err)
	}

	arg := Argument{
		Type: argumentType,
	}

	switch arg.Type {
	case ArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			arg.Value = value
		}
	case ArgumentTypeVariable:
		name := getStringValueByKey(raw, "name")
		if name == "" {
			return errors.New("field name in Argument is required for variable type")
		}
		arg.Name = name
	}

	*j = arg
	return nil
}

// RelationshipArgumentType represents a relationship argument type enum
type RelationshipArgumentType string

const (
	RelationshipArgumentTypeLiteral  RelationshipArgumentType = "literal"
	RelationshipArgumentTypeVariable RelationshipArgumentType = "variable"
	RelationshipArgumentTypeColumn   RelationshipArgumentType = "column"
)

var enumValues_RelationshipArgumentType = []RelationshipArgumentType{
	RelationshipArgumentTypeLiteral,
	RelationshipArgumentTypeVariable,
	RelationshipArgumentTypeColumn,
}

// ParseRelationshipArgumentType parses a relationship argument type from string
func ParseRelationshipArgumentType(input string) (RelationshipArgumentType, error) {
	result := RelationshipArgumentType(input)
	if !slices.Contains(enumValues_RelationshipArgumentType, result) {
		return RelationshipArgumentType(""), fmt.Errorf("failed to parse RelationshipArgumentType, expect one of %v, got %s", enumValues_RelationshipArgumentType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid
func (j RelationshipArgumentType) IsValid() bool {
	return slices.Contains(enumValues_RelationshipArgumentType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationshipArgumentType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseRelationshipArgumentType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// RelationshipArgument is provided by reference to a variable or as a literal value
type RelationshipArgument struct {
	Type  RelationshipArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Name  string                   `json:"name" yaml:"name" mapstructure:"name"`
	Value any                      `json:"value" yaml:"value" mapstructure:"value"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationshipArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawArgumentType := getStringValueByKey(raw, "type")
	if rawArgumentType == "" {
		return errors.New("field type in Argument: required")
	}

	argumentType, err := ParseRelationshipArgumentType(rawArgumentType)
	if err != nil {
		return fmt.Errorf("field type in Argument: %s", err)
	}

	arg := RelationshipArgument{
		Type: argumentType,
	}

	switch arg.Type {
	case RelationshipArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			arg.Value = value
		}
	default:
		name := getStringValueByKey(raw, "name")
		if name == "" {
			return fmt.Errorf("field name in Argument is required for %s type", rawArgumentType)
		}
		arg.Name = name
	}

	*j = arg
	return nil
}

// FieldType represents a field type
type FieldType string

const (
	FieldTypeColumn       FieldType = "column"
	FieldTypeRelationship FieldType = "relationship"
)

var enumValues_FieldType = []FieldType{
	FieldTypeColumn,
	FieldTypeRelationship,
}

// ParseFieldType parses a field type from string
func ParseFieldType(input string) (FieldType, error) {
	result := FieldType(input)
	if !slices.Contains(enumValues_FieldType, result) {
		return FieldType(""), fmt.Errorf("failed to parse FieldType, expect one of %v, got %s", enumValues_FieldType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid
func (j FieldType) IsValid() bool {
	return slices.Contains(enumValues_FieldType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *FieldType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseFieldType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// Field represents a field
type Field map[string]any

type FieldEncoder interface {
	Encode() Field
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Field) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var fieldType FieldType

	rawFieldType, ok := raw["type"]
	if !ok {
		return errors.New("field type in Field: required")
	}
	err := json.Unmarshal(rawFieldType, &fieldType)
	if err != nil {
		return fmt.Errorf("field type in Field: %s", err)
	}

	results := map[string]any{
		"type": fieldType,
	}

	switch fieldType {
	case FieldTypeColumn:
		column, err := unmarshalStringFromJsonMap(raw, "column", true)

		if err != nil {
			return fmt.Errorf("field column in Field: %s", err)
		}

		results["column"] = column

		// decode fields
		var fields NestedField
		rawFields, ok := raw["fields"]
		if ok && !isNullJSON(rawFields) {
			if err = json.Unmarshal(rawFields, &fields); err != nil {
				return fmt.Errorf("field fields in Field: %s", err)
			}
			results["fields"] = fields
		}
	case FieldTypeRelationship:
		relationship, err := unmarshalStringFromJsonMap(raw, "relationship", true)
		if err != nil {
			return fmt.Errorf("field relationship in Field: %s", err)
		}
		results["relationship"] = relationship

		rawQuery, ok := raw["query"]
		if !ok {
			return errors.New("field query in Field: required")
		}
		var query Query
		if err = json.Unmarshal(rawQuery, &query); err != nil {
			return fmt.Errorf("field query in Field: %s", err)
		}
		results["query"] = query

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in Field: required")
		}

		var arguments map[string]RelationshipArgument
		if err = json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in Field: %s", err)
		}
		results["arguments"] = arguments
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type
func (j Field) Type() (FieldType, error) {
	t, ok := j["type"]
	if !ok {
		return FieldType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseFieldType(raw)
		if err != nil {
			return FieldType(""), err
		}
		return v, nil
	case FieldType:
		return raw, nil
	default:
		return FieldType(""), fmt.Errorf("invalid Field type: %+v", t)
	}
}

// AsColumn tries to convert the current type to ColumnField
func (j Field) AsColumn() (*ColumnField, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != FieldTypeColumn {
		return nil, fmt.Errorf("invalid Field type; expected %s, got %s", FieldTypeColumn, t)
	}
	column := getStringValueByKey(j, "column")
	if column == "" {
		return nil, errors.New("ColumnField.column is required")
	}

	result := &ColumnField{
		Type:   t,
		Column: column,
	}
	rawFields, ok := j["fields"]
	if ok && !isNil(rawFields) {
		fields, ok := rawFields.(NestedField)
		if !ok {
			return nil, fmt.Errorf("invalid ColumnField.fields type; expected NestedField, got %+v", rawFields)
		}
		result.Fields = fields
	}

	return result, nil
}

// AsRelationship tries to convert the current type to RelationshipField
func (j Field) AsRelationship() (*RelationshipField, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != FieldTypeRelationship {
		return nil, fmt.Errorf("invalid Field type; expected %s, got %s", FieldTypeRelationship, t)
	}
	relationship := getStringValueByKey(j, "relationship")
	if relationship == "" {
		return nil, errors.New("RelationshipField.relationship is required")
	}

	rawQuery, ok := j["query"]
	if !ok {
		return nil, errors.New("RelationshipField.query is required")
	}
	query, ok := rawQuery.(Query)
	if !ok {
		return nil, fmt.Errorf("invalid RelationshipField.query type; expected Query, got %+v", rawQuery)
	}

	rawArguments, ok := j["arguments"]
	if !ok {
		return nil, errors.New("RelationshipField.arguments is required")
	}
	arguments, ok := rawArguments.(map[string]RelationshipArgument)
	if !ok {
		return nil, fmt.Errorf("invalid RelationshipField.arguments type; expected map[string]RelationshipArgument, got %+v", rawArguments)
	}

	return &RelationshipField{
		Type:         t,
		Query:        query,
		Relationship: relationship,
		Arguments:    arguments,
	}, nil
}

// Interface converts the comparison value to its generic interface
func (j Field) Interface() FieldEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error
func (j Field) InterfaceT() (FieldEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case FieldTypeColumn:
		return j.AsColumn()
	case FieldTypeRelationship:
		return j.AsRelationship()
	default:
		return nil, fmt.Errorf("invalid Field type: %s", ty)
	}
}

// ColumnField represents a column field
type ColumnField struct {
	Type FieldType `json:"type" yaml:"type" mapstructure:"type"`
	// Column name
	Column string `json:"column" yaml:"column" mapstructure:"column"`
	// When the type of the column is a (possibly-nullable) array or object,
	// the caller can request a subset of the complete column data, by specifying fields to fetch here.
	// If omitted, the column data will be fetched in full.
	Fields NestedField `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields"`
}

// Encode converts the instance to raw Field
func (f ColumnField) Encode() Field {
	r := Field{
		"type":   f.Type,
		"column": f.Column,
	}

	if len(f.Fields) > 0 {
		r["fields"] = f.Fields
	}
	return r
}

// NewColumnField creates a new ColumnField instance
func NewColumnField(column string, fields NestedFieldEncoder) *ColumnField {
	var field NestedField
	if !isNil(fields) {
		field = fields.Encode()
	}
	return &ColumnField{
		Type:   FieldTypeColumn,
		Column: column,
		Fields: field,
	}
}

// RelationshipField represents a relationship field
type RelationshipField struct {
	Type FieldType `json:"type" yaml:"type" mapstructure:"type"`
	// The relationship query
	Query Query `json:"query" yaml:"query" mapstructure:"query"`
	// The name of the relationship to follow for the subquery
	Relationship string `json:"relationship" yaml:"relationship" mapstructure:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
}

// Encode converts the instance to raw Field
func (f RelationshipField) Encode() Field {
	return Field{
		"type":         f.Type,
		"query":        f.Query,
		"relationship": f.Relationship,
		"arguments":    f.Arguments,
	}
}

// NewRelationshipField creates a new RelationshipField instance
func NewRelationshipField(query Query, relationship string, arguments map[string]RelationshipArgument) *RelationshipField {
	return &RelationshipField{
		Type:         FieldTypeRelationship,
		Query:        query,
		Relationship: relationship,
		Arguments:    arguments,
	}
}

// ComparisonTargetType represents comparison target enums
type ComparisonTargetType string

const (
	ComparisonTargetTypeColumn               ComparisonTargetType = "column"
	ComparisonTargetTypeRootCollectionColumn ComparisonTargetType = "root_collection_column"
)

var enumValues_ComparisonTargetType = []ComparisonTargetType{
	ComparisonTargetTypeColumn,
	ComparisonTargetTypeRootCollectionColumn,
}

// ParseComparisonTargetType parses a comparison target type argument type from string
func ParseComparisonTargetType(input string) (ComparisonTargetType, error) {
	result := ComparisonTargetType(input)
	if !slices.Contains(enumValues_ComparisonTargetType, result) {
		return ComparisonTargetType(""), fmt.Errorf("failed to parse ComparisonTargetType, expect one of %v, got: %s", enumValues_ComparisonTargetType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j ComparisonTargetType) IsValid() bool {
	return slices.Contains(enumValues_ComparisonTargetType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonTargetType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseComparisonTargetType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// ComparisonTarget represents a comparison target object
type ComparisonTarget struct {
	Type ComparisonTargetType `json:"type" yaml:"type" mapstructure:"type"`
	Name string               `json:"name" yaml:"name" mapstructure:"name"`
	Path []PathElement        `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path"`
}

// ExpressionType represents the filtering expression enums
type ExpressionType string

const (
	ExpressionTypeAnd                      ExpressionType = "and"
	ExpressionTypeOr                       ExpressionType = "or"
	ExpressionTypeNot                      ExpressionType = "not"
	ExpressionTypeUnaryComparisonOperator  ExpressionType = "unary_comparison_operator"
	ExpressionTypeBinaryComparisonOperator ExpressionType = "binary_comparison_operator"
	ExpressionTypeExists                   ExpressionType = "exists"
)

var enumValues_ExpressionType = []ExpressionType{
	ExpressionTypeAnd,
	ExpressionTypeOr,
	ExpressionTypeNot,
	ExpressionTypeUnaryComparisonOperator,
	ExpressionTypeBinaryComparisonOperator,
	ExpressionTypeExists,
}

// ParseExpressionType parses a comparison target type argument type from string
func ParseExpressionType(input string) (ExpressionType, error) {
	result := ExpressionType(input)
	if !slices.Contains(enumValues_ExpressionType, ExpressionType(input)) {
		return ExpressionType(""), fmt.Errorf("failed to parse ExpressionType, expect one of %v, got %s", enumValues_ExpressionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j ExpressionType) IsValid() bool {
	return slices.Contains(enumValues_ExpressionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExpressionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseExpressionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// ComparisonValueType represents a comparison value type enum
type ComparisonValueType string

const (
	ComparisonValueTypeColumn   ComparisonValueType = "column"
	ComparisonValueTypeScalar   ComparisonValueType = "scalar"
	ComparisonValueTypeVariable ComparisonValueType = "variable"
)

var enumValues_ComparisonValueType = []ComparisonValueType{
	ComparisonValueTypeColumn,
	ComparisonValueTypeScalar,
	ComparisonValueTypeVariable,
}

// ParseComparisonValueType parses a comparison value type from string
func ParseComparisonValueType(input string) (ComparisonValueType, error) {
	result := ComparisonValueType(input)
	if !slices.Contains(enumValues_ComparisonValueType, ComparisonValueType(input)) {
		return ComparisonValueType(""), fmt.Errorf("failed to parse ComparisonValueType, expect one of %v, got %s", enumValues_ComparisonValueType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j ComparisonValueType) IsValid() bool {
	return slices.Contains(enumValues_ComparisonValueType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonValueType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseComparisonValueType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// ComparisonValue represents a raw comparison value object with validation
type ComparisonValue map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonValue) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in ComparisonValue: required")
	}

	var ty ComparisonValueType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in ComparisonValue: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case ComparisonValueTypeVariable:
		rawName, ok := raw["name"]
		if !ok {
			return errors.New("field name in ComparisonValue is required for variable type")
		}
		var name string
		if err := json.Unmarshal(rawName, &name); err != nil {
			return fmt.Errorf("field name in ComparisonValue: %s", err)
		}
		result["name"] = name
	case ComparisonValueTypeColumn:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in ComparisonValue is required for column type")
		}
		var column ComparisonTarget
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in ComparisonValue: %s", err)
		}
		result["column"] = column
	case ComparisonValueTypeScalar:
		rawValue, ok := raw["value"]
		if !ok {
			return errors.New("field value in ComparisonValue is required for scalar type")
		}
		var value any
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in ComparisonValue: %s", err)
		}
		result["value"] = value
	}
	*j = result
	return nil
}

// GetType gets the type of comparison value
func (cv ComparisonValue) Type() (ComparisonValueType, error) {
	t, ok := cv["type"]
	if !ok {
		return ComparisonValueType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseComparisonValueType(raw)
		if err != nil {
			return ComparisonValueType(""), err
		}
		return v, nil
	case ComparisonValueType:
		return raw, nil
	default:
		return ComparisonValueType(""), fmt.Errorf("invalid ComparisonValue type: %+v", t)
	}
}

// AsScalar tries to convert the comparison value to scalar
func (cv ComparisonValue) AsScalar() (*ComparisonValueScalar, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != ComparisonValueTypeScalar {
		return nil, fmt.Errorf("invalid ComparisonValue type; expected %s, got %s", ComparisonValueTypeScalar, ty)
	}

	value, ok := cv["value"]
	if !ok {
		return nil, errors.New("ComparisonValueScalar.value is required")
	}

	return &ComparisonValueScalar{
		Type:  ty,
		Value: value,
	}, nil
}

// AsColumn tries to convert the comparison value to column
func (cv ComparisonValue) AsColumn() (*ComparisonValueColumn, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != ComparisonValueTypeColumn {
		return nil, fmt.Errorf("invalid ComparisonValue type; expected %s, got %s", ComparisonValueTypeColumn, ty)
	}

	rawColumn, ok := cv["column"]
	if !ok {
		return nil, errors.New("ComparisonValueColumn.column is required")
	}

	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		return nil, fmt.Errorf("invalid ComparisonValueColumn.column; expected ComparisonTarget, got %+v", rawColumn)
	}
	return &ComparisonValueColumn{
		Type:   ty,
		Column: column,
	}, nil
}

// AsVariable tries to convert the comparison value to column
func (cv ComparisonValue) AsVariable() (*ComparisonValueVariable, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != ComparisonValueTypeVariable {
		return nil, fmt.Errorf("invalid ComparisonValue type; expected %s, got %s", ComparisonValueTypeVariable, ty)
	}

	name := getStringValueByKey(cv, "name")
	if name == "" {
		return nil, errors.New("ComparisonValueVariable.name is required")
	}
	return &ComparisonValueVariable{
		Type: ty,
		Name: name,
	}, nil
}

// Interface converts the comparison value to its generic interface
func (cv ComparisonValue) Interface() ComparisonValueEncoder {
	result, _ := cv.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error
func (cv ComparisonValue) InterfaceT() (ComparisonValueEncoder, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case ComparisonValueTypeColumn:
		return cv.AsColumn()
	case ComparisonValueTypeVariable:
		return cv.AsVariable()
	case ComparisonValueTypeScalar:
		return cv.AsScalar()
	default:
		return nil, fmt.Errorf("invalid ComparisonValue type: %s", ty)
	}
}

// ComparisonValueEncoder represents a comparison value encoder interface
type ComparisonValueEncoder interface {
	Encode() ComparisonValue
}

// ComparisonValueColumn represents a comparison value with column type
type ComparisonValueColumn struct {
	Type   ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Column ComparisonTarget    `json:"column" yaml:"column" mapstructure:"column"`
}

// Encode converts to the raw comparison value
func (cv ComparisonValueColumn) Encode() ComparisonValue {
	return map[string]any{
		"type":   cv.Type,
		"column": cv.Column,
	}
}

// ComparisonValueScalar represents a comparison value with scalar type
type ComparisonValueScalar struct {
	Type  ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Value any                 `json:"value" yaml:"value" mapstructure:"value"`
}

// Encode converts to the raw comparison value
func (cv ComparisonValueScalar) Encode() ComparisonValue {
	return map[string]any{
		"type":  cv.Type,
		"value": cv.Value,
	}
}

// ComparisonValueVariable represents a comparison value with variable type
type ComparisonValueVariable struct {
	Type ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Name string              `json:"name" yaml:"name" mapstructure:"name"`
}

// Encode converts to the raw comparison value
func (cv ComparisonValueVariable) Encode() ComparisonValue {
	return map[string]any{
		"type": cv.Type,
		"name": cv.Name,
	}
}

// ExistsInCollectionType represents an exists in collection type enum
type ExistsInCollectionType string

const (
	ExistsInCollectionTypeRelated   ExistsInCollectionType = "related"
	ExistsInCollectionTypeUnrelated ExistsInCollectionType = "unrelated"
)

var enumValues_ExistsInCollectionType = []ExistsInCollectionType{
	ExistsInCollectionTypeRelated,
	ExistsInCollectionTypeUnrelated,
}

// ParseExistsInCollectionType parses a comparison value type from string
func ParseExistsInCollectionType(input string) (ExistsInCollectionType, error) {
	result := ExistsInCollectionType(input)
	if !slices.Contains(enumValues_ExistsInCollectionType, result) {
		return result, fmt.Errorf("failed to parse ExistsInCollectionType, expect one of %v, got %s", enumValues_ExistsInCollectionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j ExistsInCollectionType) IsValid() bool {
	return slices.Contains(enumValues_ExistsInCollectionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExistsInCollectionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseExistsInCollectionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// ExistsInCollection represents an Exists In Collection object
type ExistsInCollection map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExistsInCollection) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in ExistsInCollection: required")
	}

	var ty ExistsInCollectionType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in ExistsInCollection: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case ExistsInCollectionTypeRelated:
		rawRelationship, ok := raw["relationship"]
		if !ok {
			return errors.New("field relationship in ExistsInCollection is required for related type")
		}
		var relationship string
		if err := json.Unmarshal(rawRelationship, &relationship); err != nil {
			return fmt.Errorf("field name in ExistsInCollection: %s", err)
		}
		result["relationship"] = relationship

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in ExistsInCollection is required for related type")
		}
		var arguments map[string]RelationshipArgument
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in ExistsInCollection: %s", err)
		}
		result["arguments"] = arguments
	case ExistsInCollectionTypeUnrelated:
		rawCollection, ok := raw["collection"]
		if !ok {
			return errors.New("field collection in ExistsInCollection is required for unrelated type")
		}
		var collection string
		if err := json.Unmarshal(rawCollection, &collection); err != nil {
			return fmt.Errorf("field collection in ExistsInCollection: %s", err)
		}
		result["collection"] = collection

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in ExistsInCollection is required for unrelated type")
		}
		var arguments map[string]RelationshipArgument
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in ExistsInCollection: %s", err)
		}
		result["arguments"] = arguments
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j ExistsInCollection) Type() (ExistsInCollectionType, error) {
	t, ok := j["type"]
	if !ok {
		return ExistsInCollectionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseExistsInCollectionType(raw)
		if err != nil {
			return ExistsInCollectionType(""), err
		}
		return v, nil
	case ExistsInCollectionType:
		return raw, nil
	default:
		return ExistsInCollectionType(""), fmt.Errorf("invalid ExistsInCollection type: %+v", t)
	}
}

// AsRelated tries to convert the instance to related type
func (j ExistsInCollection) AsRelated() (*ExistsInCollectionRelated, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExistsInCollectionTypeRelated {
		return nil, fmt.Errorf("invalid ExistsInCollection type; expected: %s, got: %s", ExistsInCollectionTypeRelated, t)
	}

	relationship := getStringValueByKey(j, "relationship")
	if relationship == "" {
		return nil, errors.New("ExistsInCollectionRelated.relationship is required")
	}
	rawArgs, ok := j["arguments"]
	if !ok {
		return nil, errors.New("ExistsInCollectionRelated.arguments is required")
	}
	args, ok := rawArgs.(map[string]RelationshipArgument)
	if !ok {
		return nil, fmt.Errorf("invalid ExistsInCollectionRelated.arguments type; expected: map[string]RelationshipArgument, got: %+v", rawArgs)
	}

	return &ExistsInCollectionRelated{
		Type:         t,
		Relationship: relationship,
		Arguments:    args,
	}, nil
}

// AsRelated tries to convert the instance to unrelated type
func (j ExistsInCollection) AsUnrelated() (*ExistsInCollectionUnrelated, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExistsInCollectionTypeUnrelated {
		return nil, fmt.Errorf("invalid ExistsInCollection type; expected: %s, got: %s", ExistsInCollectionTypeUnrelated, t)
	}

	collection := getStringValueByKey(j, "collection")
	if collection == "" {
		return nil, errors.New("ExistsInCollectionUnrelated.collection is required")
	}
	rawArgs, ok := j["arguments"]
	if !ok {
		return nil, errors.New("ExistsInCollectionUnrelated.arguments is required")
	}
	args, ok := rawArgs.(map[string]RelationshipArgument)
	if !ok {
		return nil, fmt.Errorf("invalid ExistsInCollectionUnrelated.arguments type; expected: map[string]RelationshipArgument, got: %+v", rawArgs)
	}

	return &ExistsInCollectionUnrelated{
		Type:       t,
		Collection: collection,
		Arguments:  args,
	}, nil
}

// Interface tries to convert the instance to the ExistsInCollectionEncoder interface
func (j ExistsInCollection) Interface() ExistsInCollectionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to the ExistsInCollectionEncoder interface safely with explicit error
func (j ExistsInCollection) InterfaceT() (ExistsInCollectionEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case ExistsInCollectionTypeRelated:
		return j.AsRelated()
	case ExistsInCollectionTypeUnrelated:
		return j.AsUnrelated()
	default:
		return nil, fmt.Errorf("invalid ExistsInCollection type: %s", t)
	}
}

// ExistsInCollectionEncoder abstracts the ExistsInCollection serialization interface
type ExistsInCollectionEncoder interface {
	Encode() ExistsInCollection
}

// ExistsInCollectionRelated represents [Related collections] that are related to the original collection by a relationship in the collection_relationships field of the top-level QueryRequest.
//
// [Related collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#related-collections
type ExistsInCollectionRelated struct {
	Type         ExistsInCollectionType `json:"type" yaml:"type" mapstructure:"type"`
	Relationship string                 `json:"relationship" yaml:"relationship" mapstructure:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
}

// Encode converts the instance to its raw type
func (ei ExistsInCollectionRelated) Encode() ExistsInCollection {
	return ExistsInCollection{
		"type":         ei.Type,
		"relationship": ei.Relationship,
		"arguments":    ei.Arguments,
	}
}

// ExistsInCollectionUnrelated represents [unrelated collections].
//
// [unrelated collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#unrelated-collections
type ExistsInCollectionUnrelated struct {
	Type ExistsInCollectionType `json:"type" yaml:"type" mapstructure:"type"`
	// The name of a collection
	Collection string `json:"collection" yaml:"collection" mapstructure:"collection"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
}

// Encode converts the instance to its raw type
func (ei ExistsInCollectionUnrelated) Encode() ExistsInCollection {
	return ExistsInCollection{
		"type":       ei.Type,
		"collection": ei.Collection,
		"arguments":  ei.Arguments,
	}
}

// Expression represents the query expression object
type Expression map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Expression) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in Expression: required")
	}

	var ty ExpressionType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in Expression: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case ExpressionTypeAnd, ExpressionTypeOr:
		rawExpressions, ok := raw["expressions"]
		if !ok {
			return fmt.Errorf("field expressions in Expression is required for '%s' type", ty)
		}
		var expressions []Expression
		if err := json.Unmarshal(rawExpressions, &expressions); err != nil {
			return fmt.Errorf("field expressions in Expression: %s", err)
		}
		result["expressions"] = expressions
	case ExpressionTypeNot:
		rawExpression, ok := raw["expression"]
		if !ok {
			return fmt.Errorf("field expressions in Expression is required for '%s' type", ty)
		}
		var expression Expression
		if err := json.Unmarshal(rawExpression, &expression); err != nil {
			return fmt.Errorf("field expression in Expression: %s", err)
		}
		result["expression"] = expression
	case ExpressionTypeUnaryComparisonOperator:
		rawOperator, ok := raw["operator"]
		if !ok {
			return fmt.Errorf("field operator in Expression is required for '%s' type", ty)
		}
		var operator UnaryComparisonOperator
		if err := json.Unmarshal(rawOperator, &operator); err != nil {
			return fmt.Errorf("field operator in Expression: %s", err)
		}
		result["operator"] = operator

		rawColumn, ok := raw["column"]
		if !ok {
			return fmt.Errorf("field column in Expression is required for '%s' type", ty)
		}
		var column ComparisonTarget
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Expression: %s", err)
		}
		result["column"] = column
	case ExpressionTypeBinaryComparisonOperator:
		rawOperator, ok := raw["operator"]
		if !ok {
			return fmt.Errorf("field operator in Expression is required for '%s' type", ty)
		}
		var operator string
		if err := json.Unmarshal(rawOperator, &operator); err != nil {
			return fmt.Errorf("field operator in Expression: %s", err)
		}

		if operator == "" {
			return fmt.Errorf("field operator in Expression is required for '%s' type", ty)
		}
		result["operator"] = operator

		rawColumn, ok := raw["column"]
		if !ok {
			return fmt.Errorf("field column in Expression is required for '%s' type", ty)
		}
		var column ComparisonTarget
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Expression: %s", err)
		}
		result["column"] = column

		rawValue, ok := raw["value"]
		if !ok {
			return fmt.Errorf("field value in Expression is required for '%s' type", ty)
		}
		var value ComparisonValue
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in Expression: %s", err)
		}
		result["value"] = value
	case ExpressionTypeExists:
		rawPredicate, ok := raw["predicate"]
		if ok {
			var predicate Expression
			if err := json.Unmarshal(rawPredicate, &predicate); err != nil {
				return fmt.Errorf("field predicate in Expression: %s", err)
			}
			result["predicate"] = predicate
		}

		rawInCollection, ok := raw["in_collection"]
		if !ok {
			return fmt.Errorf("field in_collection in Expression is required for '%s' type", ty)
		}
		var inCollection ExistsInCollection
		if err := json.Unmarshal(rawInCollection, &inCollection); err != nil {
			return fmt.Errorf("field in_collection in Expression: %s", err)
		}
		result["in_collection"] = inCollection
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j Expression) Type() (ExpressionType, error) {
	t, ok := j["type"]
	if !ok {
		return ExpressionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseExpressionType(raw)
		if err != nil {
			return ExpressionType(""), err
		}
		return v, nil
	case ExpressionType:
		return raw, nil
	default:
		return ExpressionType(""), fmt.Errorf("invalid Expression type: %+v", t)
	}
}

// AsAnd tries to convert the instance to and type
func (j Expression) AsAnd() (*ExpressionAnd, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeAnd {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeAnd, t)
	}

	rawExpressions, ok := j["expressions"]
	if !ok {
		return nil, errors.New("ExpressionAnd.expression is required")
	}
	expressions, ok := rawExpressions.([]Expression)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionAnd.expression type; expected: []Expression, got: %+v", rawExpressions)
	}

	return &ExpressionAnd{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsOr tries to convert the instance to ExpressionOr instance
func (j Expression) AsOr() (*ExpressionOr, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeOr {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeOr, t)
	}

	rawExpressions, ok := j["expressions"]
	if !ok {
		return nil, errors.New("ExpressionOr.expression is required")
	}
	expressions, ok := rawExpressions.([]Expression)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionOr.expression type; expected: []Expression, got: %+v", rawExpressions)
	}

	return &ExpressionOr{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsNot tries to convert the instance to ExpressionNot instance
func (j Expression) AsNot() (*ExpressionNot, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeNot {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeNot, t)
	}

	rawExpression, ok := j["expression"]
	if !ok {
		return nil, errors.New("ExpressionNot.expression is required")
	}
	expression, ok := rawExpression.(Expression)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionNot.expression type; expected: Expression, got: %+v", rawExpression)
	}

	return &ExpressionNot{
		Type:       t,
		Expression: expression,
	}, nil
}

// AsUnaryComparisonOperator tries to convert the instance to ExpressionUnaryComparisonOperator instance
func (j Expression) AsUnaryComparisonOperator() (*ExpressionUnaryComparisonOperator, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeUnaryComparisonOperator {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeUnaryComparisonOperator, t)
	}

	rawOperator, ok := j["operator"]
	if !ok {
		return nil, errors.New("ExpressionUnaryComparisonOperator.operator is required")
	}
	operator, ok := rawOperator.(UnaryComparisonOperator)
	if !ok {
		operatorStr, ok := rawOperator.(string)
		if !ok {
			return nil, fmt.Errorf("invalid ExpressionUnaryComparisonOperator.operator type; expected: UnaryComparisonOperator, got: %v", rawOperator)
		}

		operator = UnaryComparisonOperator(operatorStr)
	}

	rawColumn, ok := j["column"]
	if !ok {
		return nil, errors.New("ExpressionUnaryComparisonOperator.column is required")
	}
	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionUnaryComparisonOperator.column type; expected: ComparisonTarget, got: %v", rawColumn)
	}

	return &ExpressionUnaryComparisonOperator{
		Type:     t,
		Operator: operator,
		Column:   column,
	}, nil
}

// AsBinaryComparisonOperator tries to convert the instance to ExpressionBinaryComparisonOperator instance
func (j Expression) AsBinaryComparisonOperator() (*ExpressionBinaryComparisonOperator, error) {

	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeBinaryComparisonOperator {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeBinaryComparisonOperator, t)
	}

	rawColumn, ok := j["column"]
	if !ok {
		return nil, errors.New("ExpressionBinaryComparisonOperator.column is required")
	}
	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionBinaryComparisonOperator.column type; expected: ComparisonTarget, got: %+v", rawColumn)
	}

	rawValue, ok := j["value"]
	if !ok {
		return nil, errors.New("ExpressionBinaryComparisonOperator.value is required")
	}
	value, ok := rawValue.(ComparisonValue)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionBinaryComparisonOperator.value type; expected: ComparisonValue, got: %+v", rawValue)
	}

	return &ExpressionBinaryComparisonOperator{
		Type:     t,
		Operator: getStringValueByKey(j, "operator"),
		Column:   column,
		Value:    value,
	}, nil
}

// AsExists tries to convert the instance to ExpressionExists instance
func (j Expression) AsExists() (*ExpressionExists, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExpressionTypeExists {
		return nil, fmt.Errorf("invalid Expression type; expected: %s, got: %s", ExpressionTypeExists, t)
	}

	rawInCollection, ok := j["in_collection"]
	if !ok {
		return nil, errors.New("ExpressionExists.in_collection is required")
	}
	inCollection, ok := rawInCollection.(ExistsInCollection)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionExists.in_collection type; expected: ExistsInCollection, got: %+v", rawInCollection)
	}

	result := &ExpressionExists{
		Type:         t,
		InCollection: inCollection,
	}
	rawPredicate, ok := j["predicate"]
	if ok && !isNil(rawPredicate) {
		predicate, ok := rawPredicate.(Expression)
		if !ok {
			return nil, fmt.Errorf("invalid ExpressionExists.predicate type; expected: Expression, got: %+v", rawPredicate)
		}
		result.Predicate = predicate
	}
	return result, nil
}

// Interface tries to convert the instance to the ExpressionEncoder interface
func (j Expression) Interface() ExpressionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to the ExpressionEncoder interface safely with explicit error
func (j Expression) InterfaceT() (ExpressionEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	switch t {
	case ExpressionTypeAnd:
		return j.AsAnd()
	case ExpressionTypeOr:
		return j.AsOr()
	case ExpressionTypeNot:
		return j.AsNot()
	case ExpressionTypeUnaryComparisonOperator:
		return j.AsUnaryComparisonOperator()
	case ExpressionTypeBinaryComparisonOperator:
		return j.AsBinaryComparisonOperator()
	case ExpressionTypeExists:
		return j.AsExists()
	default:
		return nil, fmt.Errorf("invalid Expression type: %s", t)
	}
}

// ExpressionEncoder abstracts the expression encoder interface
type ExpressionEncoder interface {
	Encode() Expression
}

// ExpressionAnd is an object which represents the [conjunction of expressions]
//
// [conjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#conjunction-of-expressions
type ExpressionAnd struct {
	Type        ExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expressions []Expression   `json:"expressions" yaml:"expressions" mapstructure:"expressions"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionAnd) Encode() Expression {
	return Expression{
		"type":        exp.Type,
		"expressions": exp.Expressions,
	}
}

// ExpressionOr is an object which represents the [disjunction of expressions]
//
// [disjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#disjunction-of-expressions
type ExpressionOr struct {
	Type        ExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expressions []Expression   `json:"expressions" yaml:"expressions" mapstructure:"expressions"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionOr) Encode() Expression {
	return Expression{
		"type":        exp.Type,
		"expressions": exp.Expressions,
	}
}

// ExpressionNot is an object which represents the [negation of an expression]
//
// [negation of an expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#negation
type ExpressionNot struct {
	Type       ExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expression Expression     `json:"expression" yaml:"expression" mapstructure:"expression"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionNot) Encode() Expression {
	return Expression{
		"type":       exp.Type,
		"expression": exp.Expression,
	}
}

// ExpressionUnaryComparisonOperator is an object which represents a [unary operator expression]
//
// [unary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionUnaryComparisonOperator struct {
	Type     ExpressionType          `json:"type" yaml:"type" mapstructure:"type"`
	Operator UnaryComparisonOperator `json:"operator" yaml:"operator" mapstructure:"operator"`
	Column   ComparisonTarget        `json:"column" yaml:"column" mapstructure:"column"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionUnaryComparisonOperator) Encode() Expression {
	return Expression{
		"type":     exp.Type,
		"operator": exp.Operator,
		"column":   exp.Column,
	}
}

// ExpressionBinaryComparisonOperator is an object which represents an [binary operator expression]
//
// [binary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionBinaryComparisonOperator struct {
	Type     ExpressionType   `json:"type" yaml:"type" mapstructure:"type"`
	Operator string           `json:"operator" yaml:"operator" mapstructure:"operator"`
	Column   ComparisonTarget `json:"column" yaml:"column" mapstructure:"column"`
	Value    ComparisonValue  `json:"value" yaml:"value" mapstructure:"value"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionBinaryComparisonOperator) Encode() Expression {
	return Expression{
		"type":     exp.Type,
		"operator": exp.Operator,
		"column":   exp.Column,
		"value":    exp.Value,
	}
}

// ExpressionExists is an object which represents an [EXISTS expression]
//
// [EXISTS expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#exists-expressions
type ExpressionExists struct {
	Type         ExpressionType     `json:"type" yaml:"type" mapstructure:"type"`
	Predicate    Expression         `json:"predicate" yaml:"predicate" mapstructure:"predicate"`
	InCollection ExistsInCollection `json:"in_collection" yaml:"in_collection" mapstructure:"in_collection"`
}

// Encode converts the instance to a raw Expression
func (exp ExpressionExists) Encode() Expression {
	return Expression{
		"type":          exp.Type,
		"predicate":     exp.Predicate,
		"in_collection": exp.InCollection,
	}
}

// AggregateType represents an aggregate type
type AggregateType string

const (
	// AggregateTypeColumnCount aggregates count the number of rows with non-null values in the specified columns.
	// If the distinct flag is set, then the count should only count unique non-null values of those columns,
	AggregateTypeColumnCount AggregateType = "column_count"
	// AggregateTypeSingleColumn aggregates apply an aggregation function (as defined by the column's scalar type in the schema response) to a column
	AggregateTypeSingleColumn AggregateType = "single_column"
	// AggregateTypeStarCount aggregates count all matched rows
	AggregateTypeStarCount AggregateType = "star_count"
)

var enumValues_AggregateType = []AggregateType{
	AggregateTypeColumnCount,
	AggregateTypeSingleColumn,
	AggregateTypeStarCount,
}

// ParseAggregateType parses an aggregate type argument type from string
func ParseAggregateType(input string) (AggregateType, error) {
	result := AggregateType(input)
	if !slices.Contains(enumValues_AggregateType, result) {
		return AggregateType(""), fmt.Errorf("failed to parse AggregateType, expect one of %v, got %s", enumValues_AggregateType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j AggregateType) IsValid() bool {
	return slices.Contains(enumValues_AggregateType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseAggregateType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// Aggregate represents an [aggregated query] object
//
// [aggregated query]: https://hasura.github.io/ndc-spec/specification/queries/aggregates.html
type Aggregate map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Aggregate) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in Aggregate: required")
	}

	var ty AggregateType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in Aggregate: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case AggregateTypeStarCount:
	case AggregateTypeSingleColumn:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in Aggregate is required for single_column type")
		}
		var column string
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Aggregate: %s", err)
		}
		result["column"] = column

		rawFunction, ok := raw["function"]
		if !ok {
			return errors.New("field function in Aggregate is required for single_column type")
		}
		var function string
		if err := json.Unmarshal(rawFunction, &function); err != nil {
			return fmt.Errorf("field function in Aggregate: %s", err)
		}
		result["function"] = function

	case AggregateTypeColumnCount:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in Aggregate is required for column_count type")
		}
		var column string
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Aggregate: %s", err)
		}
		result["column"] = column

		rawDistinct, ok := raw["distinct"]
		if !ok {
			return errors.New("field distinct in Aggregate is required for column_count type")
		}
		var distinct bool
		if err := json.Unmarshal(rawDistinct, &distinct); err != nil {
			return fmt.Errorf("field distinct in Aggregate: %s", err)
		}
		result["distinct"] = distinct
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j Aggregate) Type() (AggregateType, error) {
	t, ok := j["type"]
	if !ok {
		return AggregateType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseAggregateType(raw)
		if err != nil {
			return AggregateType(""), err
		}
		return v, nil
	case AggregateType:
		return raw, nil
	default:
		return AggregateType(""), fmt.Errorf("invalid Aggregate type: %+v", t)
	}
}

// AsStarCount tries to convert the instance to AggregateStarCount type
func (j Aggregate) AsStarCount() (*AggregateStarCount, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateTypeStarCount {
		return nil, fmt.Errorf("invalid Aggregate type; expected: %s, got: %s", AggregateTypeStarCount, t)
	}

	return &AggregateStarCount{
		Type: t,
	}, nil
}

// AsSingleColumn tries to convert the instance to AggregateSingleColumn type
func (j Aggregate) AsSingleColumn() (*AggregateSingleColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateTypeSingleColumn {
		return nil, fmt.Errorf("invalid Aggregate type; expected: %s, got: %s", AggregateTypeSingleColumn, t)
	}

	column := getStringValueByKey(j, "column")
	if column == "" {
		return nil, errors.New("AggregateSingleColumn.column is required")
	}

	function := getStringValueByKey(j, "function")
	if function == "" {
		return nil, errors.New("AggregateSingleColumn.function is required")
	}
	return &AggregateSingleColumn{
		Type:     t,
		Column:   column,
		Function: function,
	}, nil
}

// AsColumnCount tries to convert the instance to AggregateColumnCount type
func (j Aggregate) AsColumnCount() (*AggregateColumnCount, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateTypeColumnCount {
		return nil, fmt.Errorf("invalid Aggregate type; expected: %s, got: %s", AggregateTypeColumnCount, t)
	}

	column := getStringValueByKey(j, "column")
	if column == "" {
		return nil, errors.New("AggregateColumnCount.column is required")
	}

	rawDistinct, ok := j["distinct"]
	if !ok {
		return nil, errors.New("AggregateColumnCount.distinct is required")
	}
	distinct, ok := rawDistinct.(bool)
	if !ok {
		return nil, fmt.Errorf("invalid AggregateColumnCount.distinct type; expected bool, got %+v", rawDistinct)
	}
	return &AggregateColumnCount{
		Type:     t,
		Column:   column,
		Distinct: distinct,
	}, nil
}

// Interface tries to convert the instance to AggregateEncoder interface
func (j Aggregate) Interface() AggregateEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to AggregateEncoder interface safely with explicit error
func (j Aggregate) InterfaceT() (AggregateEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case AggregateTypeStarCount:
		return j.AsStarCount()
	case AggregateTypeColumnCount:
		return j.AsColumnCount()
	case AggregateTypeSingleColumn:
		return j.AsSingleColumn()
	default:
		return nil, fmt.Errorf("invalid Aggregate type: %s", t)
	}
}

// AggregateEncoder abstracts the serialization interface for Aggregate
type AggregateEncoder interface {
	Encode() Aggregate
}

// AggregateStarCount represents an aggregate object which counts all matched rows
type AggregateStarCount struct {
	Type AggregateType `json:"type" mapstructure:"type"`
}

// Encode converts the instance to raw Aggregate
func (ag AggregateStarCount) Encode() Aggregate {
	return Aggregate{
		"type": ag.Type,
	}
}

// NewAggregateStarCount creates a new AggregateStarCount instance
func NewAggregateStarCount() *AggregateStarCount {
	return &AggregateStarCount{
		Type: AggregateTypeStarCount,
	}
}

// AggregateSingleColumn represents an aggregate object which applies an aggregation function (as defined by the column's scalar type in the schema response) to a column.
type AggregateSingleColumn struct {
	Type AggregateType `json:"type" yaml:"type" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" yaml:"column" mapstructure:"column"`
	// Single column aggregate function name.
	Function string `json:"function" yaml:"function" mapstructure:"function"`
}

// Encode converts the instance to raw Aggregate
func (ag AggregateSingleColumn) Encode() Aggregate {
	return Aggregate{
		"type":     ag.Type,
		"column":   ag.Column,
		"function": ag.Function,
	}
}

// NewAggregateSingleColumn creates a new AggregateSingleColumn instance
func NewAggregateSingleColumn(column string, function string) *AggregateSingleColumn {
	return &AggregateSingleColumn{
		Type:     AggregateTypeSingleColumn,
		Column:   column,
		Function: function,
	}
}

// AggregateColumnCount represents an aggregate object which count the number of rows with non-null values in the specified columns.
// If the distinct flag is set, then the count should only count unique non-null values of those columns.
type AggregateColumnCount struct {
	Type AggregateType `json:"type" yaml:"type" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" yaml:"column" mapstructure:"column"`
	// Whether or not only distinct items should be counted.
	Distinct bool `json:"distinct" yaml:"distinct" mapstructure:"distinct"`
}

// Encode converts the instance to raw Aggregate
func (ag AggregateColumnCount) Encode() Aggregate {
	return Aggregate{
		"type":     ag.Type,
		"column":   ag.Column,
		"distinct": ag.Distinct,
	}
}

// NewAggregateColumnCount creates a new AggregateColumnCount instance
func NewAggregateColumnCount(column string, distinct bool) *AggregateColumnCount {
	return &AggregateColumnCount{
		Type:     AggregateTypeColumnCount,
		Column:   column,
		Distinct: distinct,
	}
}

// OrderByTargetType represents a ordering target type
type OrderByTargetType string

const (
	OrderByTargetTypeColumn                OrderByTargetType = "column"
	OrderByTargetTypeSingleColumnAggregate OrderByTargetType = "single_column_aggregate"
	OrderByTargetTypeStarCountAggregate    OrderByTargetType = "star_count_aggregate"
)

var enumValues_OrderByTargetType = []OrderByTargetType{
	OrderByTargetTypeColumn,
	OrderByTargetTypeSingleColumnAggregate,
	OrderByTargetTypeStarCountAggregate,
}

// ParseOrderByTargetType parses a ordering target type argument type from string
func ParseOrderByTargetType(input string) (OrderByTargetType, error) {
	result := OrderByTargetType(input)
	if !slices.Contains(enumValues_OrderByTargetType, result) {
		return OrderByTargetType(""), fmt.Errorf("failed to parse OrderByTargetType, expect one of %v, got %s", enumValues_OrderByTargetType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j OrderByTargetType) IsValid() bool {
	return slices.Contains(enumValues_OrderByTargetType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *OrderByTargetType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseOrderByTargetType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// OrderByTarget represents an [order_by field] of the Query object
//
// [order_by field]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html
type OrderByTarget map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *OrderByTarget) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in OrderByTarget: required")
	}

	var ty OrderByTargetType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in OrderByTarget: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case OrderByTargetTypeColumn:
		rawName, ok := raw["name"]
		if !ok {
			return errors.New("field name in OrderByTarget is required for column type")
		}
		var name string
		if err := json.Unmarshal(rawName, &name); err != nil {
			return fmt.Errorf("field name in OrderByTarget: %s", err)
		}
		result["name"] = name

		rawPath, ok := raw["path"]
		if !ok {
			return errors.New("field path in OrderByTarget is required for column type")
		}
		var pathElem []PathElement
		if err := json.Unmarshal(rawPath, &pathElem); err != nil {
			return fmt.Errorf("field path in OrderByTarget: %s", err)
		}
		result["path"] = pathElem
	case OrderByTargetTypeSingleColumnAggregate:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in OrderByTarget is required for single_column_aggregate type")
		}
		var column string
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in OrderByTarget: %s", err)
		}
		result["column"] = column

		rawFunction, ok := raw["function"]
		if !ok {
			return errors.New("field function in OrderByTarget is required for single_column_aggregate type")
		}
		var function string
		if err := json.Unmarshal(rawFunction, &function); err != nil {
			return fmt.Errorf("field function in OrderByTarget: %s", err)
		}
		result["function"] = function

		rawPath, ok := raw["path"]
		if !ok {
			return errors.New("field path in OrderByTarget is required for single_column_aggregate type")
		}
		var pathElem []PathElement
		if err := json.Unmarshal(rawPath, &pathElem); err != nil {
			return fmt.Errorf("field path in OrderByTarget: %s", err)
		}
		result["path"] = pathElem
	case OrderByTargetTypeStarCountAggregate:
		rawPath, ok := raw["path"]
		if !ok {
			return errors.New("field path in OrderByTarget is required for star_count_aggregate type")
		}
		var pathElem []PathElement
		if err := json.Unmarshal(rawPath, &pathElem); err != nil {
			return fmt.Errorf("field path in OrderByTarget: %s", err)
		}
		result["path"] = pathElem
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j OrderByTarget) Type() (OrderByTargetType, error) {
	t, ok := j["type"]
	if !ok {
		return OrderByTargetType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseOrderByTargetType(raw)
		if err != nil {
			return OrderByTargetType(""), err
		}
		return v, nil
	case OrderByTargetType:
		return raw, nil
	default:
		return OrderByTargetType(""), fmt.Errorf("invalid OrderByTarget type: %+v", t)
	}
}

// AsColumn tries to convert the instance to OrderByColumn type
func (j OrderByTarget) AsColumn() (*OrderByColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != OrderByTargetTypeColumn {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeColumn, t)
	}

	name := getStringValueByKey(j, "name")
	if name == "" {
		return nil, errors.New("OrderByColumn.name is required")
	}
	rawPath, ok := j["path"]
	if !ok {
		return nil, errors.New("OrderByColumn.path is required")
	}
	p, ok := rawPath.([]PathElement)
	if !ok {
		return nil, fmt.Errorf("invalid OrderByColumn.path type; expected: []PathElement, got: %+v", rawPath)
	}
	return &OrderByColumn{
		Type: t,
		Name: name,
		Path: p,
	}, nil
}

// AsSingleColumnAggregate tries to convert the instance to OrderBySingleColumnAggregate type
func (j OrderByTarget) AsSingleColumnAggregate() (*OrderBySingleColumnAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != OrderByTargetTypeSingleColumnAggregate {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeSingleColumnAggregate, t)
	}

	column := getStringValueByKey(j, "column")
	if column == "" {
		return nil, errors.New("OrderBySingleColumnAggregate.column is required")
	}

	function := getStringValueByKey(j, "function")
	if function == "" {
		return nil, errors.New("OrderBySingleColumnAggregate.function is required")
	}
	rawPath, ok := j["path"]
	if !ok {
		return nil, errors.New("OrderBySingleColumnAggregate.path is required")
	}
	p, ok := rawPath.([]PathElement)
	if !ok {
		return nil, fmt.Errorf("invalid OrderBySingleColumnAggregate.path type; expected: []PathElement, got: %+v", rawPath)
	}
	return &OrderBySingleColumnAggregate{
		Type:     t,
		Column:   column,
		Function: function,
		Path:     p,
	}, nil
}

// AsStarCountAggregate tries to convert the instance to OrderByStarCountAggregate type
func (j OrderByTarget) AsStarCountAggregate() (*OrderByStarCountAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != OrderByTargetTypeStarCountAggregate {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeStarCountAggregate, t)
	}

	rawPath, ok := j["path"]
	if !ok {
		return nil, errors.New("OrderByStarCountAggregate.path is required")
	}
	p, ok := rawPath.([]PathElement)
	if !ok {
		return nil, fmt.Errorf("invalid OrderByStarCountAggregate.path type; expected: []PathElement, got: %+v", rawPath)
	}
	return &OrderByStarCountAggregate{
		Type: t,
		Path: p,
	}, nil
}

// Interface tries to convert the instance to OrderByTargetEncoder interface
func (j OrderByTarget) Interface() OrderByTargetEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to OrderByTargetEncoder interface safely with explicit error
func (j OrderByTarget) InterfaceT() (OrderByTargetEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case OrderByTargetTypeColumn:
		return j.AsColumn()
	case OrderByTargetTypeSingleColumnAggregate:
		return j.AsSingleColumnAggregate()
	case OrderByTargetTypeStarCountAggregate:
		return j.AsStarCountAggregate()
	default:
		return nil, fmt.Errorf("invalid OrderByTarget type: %s", t)
	}
}

// OrderByTargetEncoder abstracts the serialization interface for OrderByTarget
type OrderByTargetEncoder interface {
	Encode() OrderByTarget
}

// OrderByColumn represents an ordering object which compares the value in the selected column
type OrderByColumn struct {
	Type OrderByTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// The name of the column
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Any relationships to traverse to reach this column
	Path []PathElement `json:"path" yaml:"path" mapstructure:"path"`
}

// Encode converts the instance to raw OrderByTarget
func (ob OrderByColumn) Encode() OrderByTarget {
	return OrderByTarget{
		"type": ob.Type,
		"name": ob.Name,
		"path": ob.Path,
	}
}

// OrderBySingleColumnAggregate An ordering of type [single_column_aggregate] orders rows by an aggregate computed over rows in some related collection.
// If the respective aggregates are incomparable, the ordering should continue to the next OrderByElement.
//
// [single_column_aggregate]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html#type-single_column_aggregate
type OrderBySingleColumnAggregate struct {
	Type OrderByTargetType `json:"type" yaml:"name" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" yaml:"column" mapstructure:"column"`
	// Single column aggregate function name.
	Function string `json:"function" yaml:"function" mapstructure:"function"`
	// Non-empty collection of relationships to traverse
	Path []PathElement `json:"path" yaml:"path" mapstructure:"path"`
}

// Encode converts the instance to raw OrderByTarget
func (ob OrderBySingleColumnAggregate) Encode() OrderByTarget {
	return OrderByTarget{
		"type":     ob.Type,
		"column":   ob.Column,
		"function": ob.Function,
		"path":     ob.Path,
	}
}

// OrderByStarCountAggregate An ordering of type [star_count_aggregate] orders rows by a count of rows in some related collection.
// If the respective aggregates are incomparable, the ordering should continue to the next OrderByElement.
//
// [star_count_aggregate]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html#type-star_count_aggregate
type OrderByStarCountAggregate struct {
	Type OrderByTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// Non-empty collection of relationships to traverse
	Path []PathElement `json:"path" yaml:"path" mapstructure:"path"`
}

// Encode converts the instance to raw OrderByTarget
func (ob OrderByStarCountAggregate) Encode() OrderByTarget {
	return OrderByTarget{
		"type": ob.Type,
		"path": ob.Path,
	}
}

// ComparisonOperatorDefinitionType represents a binary comparison operator type enum
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

// ParseComparisonOperatorDefinitionType parses a type of a comparison operator definition
func ParseComparisonOperatorDefinitionType(input string) (ComparisonOperatorDefinitionType, error) {
	result := ComparisonOperatorDefinitionType(input)
	if !slices.Contains(enumValues_ComparisonOperatorDefinitionType, result) {
		return ComparisonOperatorDefinitionType(""), fmt.Errorf("failed to parse ComparisonOperatorDefinitionType, expect one of %v, got %s", enumValues_ComparisonOperatorDefinitionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
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

// ComparisonOperatorDefinition the definition of a comparison operator on a scalar type
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
		return fmt.Errorf("field type in ComparisonOperatorDefinition: %s", err)
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
			return fmt.Errorf("field argument_type in ComparisonOperatorDefinition: %s", err)
		}
		result["argument_type"] = argumentType
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
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

// AsEqual tries to convert the instance to ComparisonOperatorEqual type
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

// AsIn tries to convert the instance to ComparisonOperatorIn type
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

// AsCustom tries to convert the instance to ComparisonOperatorIn type
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

// Interface tries to convert the instance to ComparisonOperatorDefinitionEncoder interface
func (j ComparisonOperatorDefinition) Interface() ComparisonOperatorDefinitionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to ComparisonOperatorDefinitionEncoder interface safely with explicit error
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

// ComparisonOperatorDefinitionEncoder abstracts the serialization interface for ComparisonOperatorDefinition
type ComparisonOperatorDefinitionEncoder interface {
	Encode() ComparisonOperatorDefinition
}

// ComparisonOperatorEqual presents an equal comparison operator
type ComparisonOperatorEqual struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewComparisonOperatorEqual create a new ComparisonOperatorEqual instance
func NewComparisonOperatorEqual() *ComparisonOperatorEqual {
	return &ComparisonOperatorEqual{
		Type: ComparisonOperatorDefinitionTypeEqual,
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition
func (ob ComparisonOperatorEqual) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type": ob.Type,
	}
}

// ComparisonOperatorIn presents an in comparison operator
type ComparisonOperatorIn struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewComparisonOperatorIn create a new ComparisonOperatorIn instance
func NewComparisonOperatorIn() *ComparisonOperatorIn {
	return &ComparisonOperatorIn{
		Type: ComparisonOperatorDefinitionTypeIn,
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition
func (ob ComparisonOperatorIn) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type": ob.Type,
	}
}

// ComparisonOperatorCustom presents a custom comparison operator
type ComparisonOperatorCustom struct {
	Type ComparisonOperatorDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The type of the argument to this operator
	ArgumentType Type `json:"argument_type" yaml:"argument_type" mapstructure:"argument_type"`
}

// NewComparisonOperatorCustom create a new ComparisonOperatorCustom instance
func NewComparisonOperatorCustom(argumentType TypeEncoder) *ComparisonOperatorCustom {
	return &ComparisonOperatorCustom{
		Type:         ComparisonOperatorDefinitionTypeCustom,
		ArgumentType: argumentType.Encode(),
	}
}

// Encode converts the instance to raw ComparisonOperatorDefinition
func (ob ComparisonOperatorCustom) Encode() ComparisonOperatorDefinition {
	return ComparisonOperatorDefinition{
		"type":          ob.Type,
		"argument_type": ob.ArgumentType,
	}
}

// NestedFieldType represents a nested field type enum
type NestedFieldType string

const (
	NestedFieldTypeObject NestedFieldType = "object"
	NestedFieldTypeArray  NestedFieldType = "array"
)

var enumValues_NestedFieldType = []NestedFieldType{
	NestedFieldTypeObject,
	NestedFieldTypeArray,
}

// ParseNestedFieldType parses the type of nested field
func ParseNestedFieldType(input string) (NestedFieldType, error) {
	result := NestedFieldType(input)
	if !slices.Contains(enumValues_NestedFieldType, result) {
		return NestedFieldType(""), fmt.Errorf("failed to parse NestedFieldType, expect one of %v, got %s", enumValues_NestedFieldType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid
func (j NestedFieldType) IsValid() bool {
	return slices.Contains(enumValues_NestedFieldType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NestedFieldType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseNestedFieldType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// NestedField represents a nested field
type NestedField map[string]any

// IsNil checks if the field is null or empty
func (j NestedField) IsNil() bool {
	return len(j) == 0
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NestedField) UnmarshalJSON(b []byte) error {
	if isNullJSON(b) {
		return nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in NestedField: required")
	}

	var ty NestedFieldType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in NestedField: %s", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case NestedFieldTypeObject:
		rawFields, ok := raw["fields"]
		if !ok {
			return errors.New("field fields in NestedField is required for object type")
		}
		var fields map[string]Field
		if err := json.Unmarshal(rawFields, &fields); err != nil {
			return fmt.Errorf("field fields in NestedField object: %s", err)
		}
		result["fields"] = fields
	case NestedFieldTypeArray:
		rawFields, ok := raw["fields"]
		if !ok {
			return errors.New("field fields in NestedField is required for array type")
		}
		var fields NestedField
		if err := json.Unmarshal(rawFields, &fields); err != nil {
			return fmt.Errorf("field fields in NestedField array: %s", err)
		}
		result["fields"] = fields
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j NestedField) Type() (NestedFieldType, error) {
	t, ok := j["type"]
	if !ok {
		return NestedFieldType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseNestedFieldType(raw)
		if err != nil {
			return NestedFieldType(""), err
		}
		return v, nil
	case NestedFieldType:
		return raw, nil
	default:
		return NestedFieldType(""), fmt.Errorf("invalid NestedField type: %+v", t)
	}
}

// AsObject tries to convert the instance to NestedObject type
func (j NestedField) AsObject() (*NestedObject, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != NestedFieldTypeObject {
		return nil, fmt.Errorf("invalid NestedField type; expected: %s, got: %s", NestedFieldTypeObject, t)
	}

	rawFields, ok := j["fields"]
	if !ok {
		return nil, errors.New("NestedObject.fields is required")
	}

	fields, ok := rawFields.(map[string]Field)
	if !ok {
		return nil, fmt.Errorf("invalid NestedObject.fields type; expected: map[string]Field, got: %+v", rawFields)
	}

	return &NestedObject{
		Type:   t,
		Fields: fields,
	}, nil
}

// AsArray tries to convert the instance to NestedArray type
func (j NestedField) AsArray() (*NestedArray, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != NestedFieldTypeArray {
		return nil, fmt.Errorf("invalid NestedField type; expected: %s, got: %s", NestedFieldTypeArray, t)
	}

	rawFields, ok := j["fields"]
	if !ok {
		return nil, errors.New("NestedArray.fields is required")
	}

	fields, ok := rawFields.(NestedField)
	if !ok {
		return nil, fmt.Errorf("invalid NestedArray.fields type; expected: NestedField, got: %+v", rawFields)
	}

	return &NestedArray{
		Type:   t,
		Fields: fields,
	}, nil
}

// Interface tries to convert the instance to NestedFieldEncoder interface
func (j NestedField) Interface() NestedFieldEncoder {
	result, _ := j.InterfaceT()
	return result
}

// Interface tries to convert the instance to NestedFieldEncoder interface safely with explicit error
func (j NestedField) InterfaceT() (NestedFieldEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case NestedFieldTypeObject:
		return j.AsObject()
	case NestedFieldTypeArray:
		return j.AsArray()
	default:
		return nil, fmt.Errorf("invalid NestedField type: %s", t)
	}
}

// NestedFieldEncoder abstracts the serialization interface for NestedField
type NestedFieldEncoder interface {
	Encode() NestedField
}

// NestedObject presents a nested object field
type NestedObject struct {
	Type   NestedFieldType  `json:"type" yaml:"type" mapstructure:"type"`
	Fields map[string]Field `json:"fields" yaml:"fields" mapstructure:"fields"`
}

// NewNestedObject create a new NestedObject instance
func NewNestedObject(fields map[string]FieldEncoder) *NestedObject {
	fieldMap := make(map[string]Field)
	for k, v := range fields {
		fieldMap[k] = v.Encode()
	}
	return &NestedObject{
		Type:   NestedFieldTypeObject,
		Fields: fieldMap,
	}
}

// Encode converts the instance to raw NestedField
func (ob NestedObject) Encode() NestedField {
	return NestedField{
		"type":   ob.Type,
		"fields": ob.Fields,
	}
}

// NestedArray presents a nested array field
type NestedArray struct {
	Type   NestedFieldType `json:"type" yaml:"type" mapstructure:"type"`
	Fields NestedField     `json:"fields" yaml:"fields" mapstructure:"fields"`
}

// NewNestedArray create a new NestedArray instance
func NewNestedArray(fields NestedFieldEncoder) *NestedArray {
	var nf NestedField
	if fields != nil {
		nf = fields.Encode()
	}
	return &NestedArray{
		Type:   NestedFieldTypeArray,
		Fields: nf,
	}
}

// Encode converts the instance to raw NestedField
func (ob NestedArray) Encode() NestedField {
	return NestedField{
		"type":   ob.Type,
		"fields": ob.Fields,
	}
}
