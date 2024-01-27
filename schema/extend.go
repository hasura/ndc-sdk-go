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
	TypeNamed    TypeEnum = "named"
	TypeNullable TypeEnum = "nullable"
	TypeArray    TypeEnum = "array"
)

// Types track the valid representations of values as JSON
type Type interface {
	GetType() TypeEnum
}

// NamedType represents a named type
type NamedType struct {
	Type TypeEnum `json:"type" mapstructure:"type"`
	// The name can refer to a primitive type or a scalar type
	Name string `json:"name" mapstructure:"name"`
}

// GetType get the type of type
func (ty NamedType) GetType() TypeEnum {
	return ty.Type
}

// NewNamedType creates a new NamedType instance
func NewNamedType(name string) *NamedType {
	return &NamedType{
		Type: TypeNamed,
		Name: name,
	}
}

// NullableType represents a nullable type
type NullableType struct {
	Type TypeEnum `json:"type" mapstructure:"type"`
	// The type of the non-null inhabitants of this type
	UnderlyingType Type `json:"underlying_type" mapstructure:"underlying_type"`
}

// NewNullableNamedType creates a new NullableType instance with underlying named type
func NewNullableNamedType(name string) *NullableType {
	return &NullableType{
		Type:           TypeNullable,
		UnderlyingType: NewNamedType(name),
	}
}

// GetType get the type of type
func (ty NullableType) GetType() TypeEnum {
	return ty.Type
}

// NewNullableArrayType creates a new NullableType instance with underlying array type
func NewNullableArrayType(elementType *ArrayType) *NullableType {
	return &NullableType{
		Type:           TypeNullable,
		UnderlyingType: elementType,
	}
}

// ArrayType represents an array type
type ArrayType struct {
	Type TypeEnum `json:"type" mapstructure:"type"`
	// The type of the elements of the array
	ElementType Type `json:"element_type" mapstructure:"element_type"`
}

// GetType get the type of type
func (ty ArrayType) GetType() TypeEnum {
	return ty.Type
}

// NewArrayType creates a new ArrayType instance
func NewArrayType(elementType Type) *ArrayType {
	return &ArrayType{
		Type:        TypeArray,
		ElementType: elementType,
	}
}

// ArgumentType represents an argument type enum
type ArgumentType string

const (
	ArgumentTypeLiteral  ArgumentType = "literal"
	ArgumentTypeVariable ArgumentType = "variable"
)

// ParseArgumentType parses an argument type from string
func ParseArgumentType(input string) (*ArgumentType, error) {
	if input != string(ArgumentTypeLiteral) && input != string(ArgumentTypeVariable) {
		return nil, fmt.Errorf("failed to parse ArgumentType, expect one of %v", []ArgumentType{ArgumentTypeLiteral, ArgumentTypeVariable})
	}
	result := ArgumentType(input)
	return &result, nil
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

	*j = *value
	return nil
}

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

	rawArgumentType := getStringValueByKey(raw, "type")
	if rawArgumentType == "" {
		return errors.New("field type in Argument: required")
	}

	argumentType, err := ParseArgumentType(rawArgumentType)
	if err != nil {
		return fmt.Errorf("field type in Argument: %s", err)
	}

	arg := Argument{
		Type: *argumentType,
	}

	switch arg.Type {
	case ArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			arg.Value = value
		}
		break
	case ArgumentTypeVariable:
		name := getStringValueByKey(raw, "name")
		if name == "" {
			return errors.New("field name in Argument is required for variable type")
		}
		arg.Name = name
		break
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

// ParseRelationshipArgumentType parses a relationship argument type from string
func ParseRelationshipArgumentType(input string) (*RelationshipArgumentType, error) {
	if input != string(RelationshipArgumentTypeLiteral) && input != string(RelationshipArgumentTypeVariable) && input != string(RelationshipArgumentTypeColumn) {
		return nil, fmt.Errorf("failed to parse ArgumentType, expect one of %v", []RelationshipArgumentType{RelationshipArgumentTypeLiteral, RelationshipArgumentTypeVariable, RelationshipArgumentTypeColumn})
	}
	result := RelationshipArgumentType(input)
	return &result, nil
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

	*j = *value
	return nil
}

// RelationshipArgument is provided by reference to a variable or as a literal value
type RelationshipArgument struct {
	Type  RelationshipArgumentType `json:"type" mapstructure:"type"`
	Name  string                   `json:"name" mapstructure:"name"`
	Value any                      `json:"value" mapstructure:"value"`
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
		Type: *argumentType,
	}

	switch arg.Type {
	case RelationshipArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			arg.Value = value
		}
		break
	default:
		name := getStringValueByKey(raw, "name")
		if name == "" {
			return fmt.Errorf("field name in Argument is required for %s type", rawArgumentType)
		}
		arg.Name = name
		break
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

// ParseFieldType parses a field type from string
func ParseFieldType(input string) (*FieldType, error) {
	if input != string(FieldTypeColumn) && input != string(FieldTypeRelationship) {
		return nil, fmt.Errorf("failed to parse FieldType, expect one of %v", []FieldType{FieldTypeColumn, FieldTypeRelationship})
	}
	result := FieldType(input)
	return &result, nil
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

	*j = *value
	return nil
}

// Field represents a fielded
type Field struct {
	Type FieldType `json:"type" mapstructure:"type"`
	// Column name, only available for column type
	Column string `json:"column" mapstructure:"column"`

	// The relationship query, only available for relationship type
	Query *Query `json:"query" mapstructure:"query"`
	// The name of the relationship to follow for the subquery
	Relationship string `json:"relationship" mapstructure:"relationship"`
	// Values to be provided to any collection arguments, only available for relationship type
	Arguments map[string]RelationshipArgument `json:"arguments" mapstructure:"arguments"`
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

	value := Field{
		Type: fieldType,
	}

	switch value.Type {
	case FieldTypeColumn:
		column, err := unmarshalStringFromJsonMap(raw, "column", true)

		if err != nil {
			return fmt.Errorf("field column in Field: %s", err)
		}

		value.Column = column
		break
	case FieldTypeRelationship:
		relationship, err := unmarshalStringFromJsonMap(raw, "relationship", true)
		if err != nil {
			return fmt.Errorf("field relationship in Field: %s", err)
		}
		value.Relationship = relationship

		rawQuery, ok := raw["query"]
		if !ok {
			return errors.New("field query in Field: required")
		}
		var query Query
		if err = json.Unmarshal(rawQuery, &query); err != nil {
			return fmt.Errorf("field query in Field: %s", err)
		}
		value.Query = &query

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in Field: required")
		}

		var arguments map[string]RelationshipArgument
		if err = json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in Field: %s", err)
		}
		value.Arguments = arguments
		break
	}

	*j = value
	return nil
}

// MutationOperationType represents the mutation operation type enum
type MutationOperationType string

const (
	MutationOperationProcedure MutationOperationType = "procedure"
)

// ParseMutationOperationType parses a mutation operation type argument type from string
func ParseMutationOperationType(input string) (*MutationOperationType, error) {
	if input != string(MutationOperationProcedure) {
		return nil, fmt.Errorf("failed to parse MutationOperationType, expect one of %v", []MutationOperationType{MutationOperationProcedure})
	}
	result := MutationOperationType(input)

	return &result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MutationOperationType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseMutationOperationType(rawValue)
	if err != nil {
		return err
	}

	*j = *value
	return nil
}

// MutationOperation represents a mutation operation
type MutationOperation struct {
	Type MutationOperationType `json:"type" mapstructure:"type"`
	// The name of the operation
	Name string `json:"name" mapstructure:"name"`
	// Any named procedure arguments
	Arguments json.RawMessage `json:"arguments" mapstructure:"arguments"`
	// The fields to return
	Fields map[string]Field
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MutationOperation) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var operationType MutationOperationType

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in MutationOperation: required")
	}
	err := json.Unmarshal(rawType, &operationType)
	if err != nil {
		return fmt.Errorf("field type in MutationOperation: %s", err)
	}

	value := MutationOperation{
		Type: operationType,
	}

	switch value.Type {
	case MutationOperationProcedure:
		name, err := unmarshalStringFromJsonMap(raw, "name", true)
		if err != nil {
			return fmt.Errorf("field name in MutationOperation: %s", err)
		}
		value.Name = name

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in MutationOperation: required")
		}

		value.Arguments = rawArguments

		rawFields, ok := raw["fields"]
		if ok && rawFields != nil {
			var fields map[string]Field
			if err = json.Unmarshal(rawFields, &fields); err != nil {
				return fmt.Errorf("field fields in MutationOperation: %s", err)
			}
			value.Fields = fields
		}

		break
	}

	*j = value
	return nil
}

// ComparisonTarget represents comparison target enums
type ComparisonTargetType string

const (
	ComparisonTargetTypeColumn               ComparisonTargetType = "column"
	ComparisonTargetTypeRootCollectionColumn ComparisonTargetType = "root_collection_column"
)

// ParseComparisonTargetType parses a comparison target type argument type from string
func ParseComparisonTargetType(input string) (*ComparisonTargetType, error) {
	if input != string(ComparisonTargetTypeColumn) && input != string(ComparisonTargetTypeRootCollectionColumn) {
		return nil, fmt.Errorf("failed to parse ComparisonTargetType, expect one of %v", []ComparisonTargetType{ComparisonTargetTypeColumn, ComparisonTargetTypeRootCollectionColumn})
	}
	result := ComparisonTargetType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// ComparisonTarget represents a comparison target object
type ComparisonTarget struct {
	Type ComparisonTargetType `json:"type" mapstructure:"type"`
	Name string               `json:"name" mapstructure:"name"`
	Path []PathElement        `json:"path,omitempty" mapstructure:"path"`
}

// ExpressionType represents the filtering expression enums
type ExpressionType string

const (
	ExpressionTypeAnd                           ExpressionType = "and"
	ExpressionTypeOr                            ExpressionType = "or"
	ExpressionTypeNot                           ExpressionType = "not"
	ExpressionTypeUnaryComparisonOperator       ExpressionType = "unary_comparison_operator"
	ExpressionTypeBinaryComparisonOperator      ExpressionType = "binary_comparison_operator"
	ExpressionTypeBinaryArrayComparisonOperator ExpressionType = "binary_array_comparison_operator"
	ExpressionTypeExists                        ExpressionType = "exists"
)

var enumValues_ExpressionType = []ExpressionType{
	ExpressionTypeAnd,
	ExpressionTypeOr,
	ExpressionTypeNot,
	ExpressionTypeUnaryComparisonOperator,
	ExpressionTypeBinaryComparisonOperator,
	ExpressionTypeBinaryArrayComparisonOperator,
	ExpressionTypeExists,
}

// ParseExpressionType parses a comparison target type argument type from string
func ParseExpressionType(input string) (*ExpressionType, error) {
	if !Contains(enumValues_ExpressionType, ExpressionType(input)) {
		return nil, fmt.Errorf("failed to parse ExpressionType, expect one of %v", enumValues_ExpressionType)
	}
	result := ExpressionType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// BinaryComparisonOperatorType represents a binary comparison operator type enum
type BinaryComparisonOperatorType string

const (
	BinaryComparisonOperatorTypeEqual BinaryComparisonOperatorType = "equal"
	BinaryComparisonOperatorTypeOther BinaryComparisonOperatorType = "other"
)

var enumValues_BinaryComparisonOperatorType = []BinaryComparisonOperatorType{
	BinaryComparisonOperatorTypeEqual,
	BinaryComparisonOperatorTypeOther,
}

// ParseBinaryComparisonOperatorType parses a comparison target type argument type from string
func ParseBinaryComparisonOperatorType(input string) (*BinaryComparisonOperatorType, error) {
	if !Contains(enumValues_BinaryComparisonOperatorType, BinaryComparisonOperatorType(input)) {
		return nil, fmt.Errorf("failed to parse BinaryComparisonOperatorType, expect one of %v", enumValues_BinaryComparisonOperatorType)
	}
	result := BinaryComparisonOperatorType(input)

	return &result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BinaryComparisonOperatorType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseBinaryComparisonOperatorType(rawValue)
	if err != nil {
		return err
	}

	*j = *value
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
func ParseComparisonValueType(input string) (*ComparisonValueType, error) {
	if !Contains(enumValues_ComparisonValueType, ComparisonValueType(input)) {
		return nil, fmt.Errorf("failed to parse ComparisonValueType, expect one of %v", enumValues_ComparisonValueType)
	}
	result := ComparisonValueType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// ComparisonValue represents a comparison value
type ComparisonValue interface {
	GetType() ComparisonValueType
}

// ComparisonValueColumn represents a comparison value with column type
type ComparisonValueColumn struct {
	Type   ComparisonValueType `json:"type" mapstructure:"type"`
	Column ComparisonTarget    `json:"column" mapstructure:"column"`
}

// ComparisonValueColumn gets the type of comparison value
func (cv ComparisonValueColumn) GetType() ComparisonValueType {
	return cv.Type
}

// ComparisonValueScalar represents a comparison value with scalar type
type ComparisonValueScalar struct {
	Type  ComparisonValueType `json:"type" mapstructure:"type"`
	Value any                 `json:"value" mapstructure:"value"`
}

// ComparisonValueScalar gets the type of comparison value
func (cv ComparisonValueScalar) GetType() ComparisonValueType {
	return cv.Type
}

// ComparisonValueVariable represents a comparison value with variable type
type ComparisonValueVariable struct {
	Type ComparisonValueType `json:"type" mapstructure:"type"`
	Name string              `json:"name" mapstructure:"name"`
}

// ComparisonValueVariable gets the type of comparison value
func (cv ComparisonValueVariable) GetType() ComparisonValueType {
	return cv.Type
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
func ParseExistsInCollectionType(input string) (*ExistsInCollectionType, error) {
	if !Contains(enumValues_ExistsInCollectionType, ExistsInCollectionType(input)) {
		return nil, fmt.Errorf("failed to parse ExistsInCollectionType, expect one of %v", enumValues_ExistsInCollectionType)
	}
	result := ExistsInCollectionType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// ExistsInCollection represents an Exists In Collection object
type ExistsInCollection interface {
	GetType() ExistsInCollectionType
}

type ExistsInCollectionRelated struct {
	Type         ExistsInCollectionType `json:"type" mapstructure:"type"`
	Relationship string                 `json:"relationship" mapstructure:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" mapstructure:"arguments"`
}

// GetType get the type of ExistsInCollection
func (ei ExistsInCollectionRelated) GetType() ExistsInCollectionType {
	return ei.Type
}

type ExistsInCollectionUnrelated struct {
	Type ExistsInCollectionType `json:"type" mapstructure:"type"`
	// The name of a collection
	Collection string `json:"collection" mapstructure:"collection"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" mapstructure:"arguments"`
}

// GetType get the type of ExistsInCollection
func (ei ExistsInCollectionUnrelated) GetType() ExistsInCollectionType {
	return ei.Type
}

// BinaryComparisonOperator represents a binary comparison operator object
type BinaryComparisonOperator struct {
	Type BinaryComparisonOperatorType `json:"type" mapstructure:"type"`
	Name string                       `json:"name,omitempty" mapstructure:"name"`
}

// Expression represents the query expression object
type Expression interface {
	GetType() ExpressionType
}

// ExpressionAnd is an object which represents the [conjunction of expressions]
//
// [conjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#conjunction-of-expressions
type ExpressionAnd struct {
	Type        ExpressionType `json:"type" mapstructure:"type"`
	Expressions []Expression   `json:"expressions" mapstructure:"expressions"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionAnd) GetType() ExpressionType {
	return exp.Type
}

// ExpressionOr is an object which represents the [disjunction of expressions]
//
// [disjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#disjunction-of-expressions
type ExpressionOr struct {
	Type        ExpressionType `json:"type" mapstructure:"type"`
	Expressions []Expression   `json:"expressions" mapstructure:"expressions"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionOr) GetType() ExpressionType {
	return exp.Type
}

// ExpressionNot is an object which represents the [negation of an expression]
//
// [negation of an expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#negation
type ExpressionNot struct {
	Type       ExpressionType `json:"type" mapstructure:"type"`
	Expression Expression     `json:"expression" mapstructure:"expression"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionNot) GetType() ExpressionType {
	return exp.Type
}

// ExpressionUnaryComparisonOperator is an object which represents a [unary operator expression]
//
// [unary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionUnaryComparisonOperator struct {
	Type     ExpressionType          `json:"type" mapstructure:"type"`
	Operator UnaryComparisonOperator `json:"operator" mapstructure:"operator"`
	Column   ComparisonTarget        `json:"column" mapstructure:"column"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionUnaryComparisonOperator) GetType() ExpressionType {
	return exp.Type
}

// ExpressionBinaryComparisonOperator is an object which represents an [binary operator expression]
//
// [binary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionBinaryComparisonOperator struct {
	Type     ExpressionType           `json:"type" mapstructure:"type"`
	Operator BinaryComparisonOperator `json:"operator" mapstructure:"operator"`
	Column   ComparisonTarget         `json:"column" mapstructure:"column"`
	Value    ComparisonValue          `json:"value" mapstructure:"value"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionBinaryComparisonOperator) GetType() ExpressionType {
	return exp.Type
}

// ExpressionBinaryArrayComparisonOperator is an object which represents an [binary array-valued comparison operators expression]
//
// [binary array-valued comparison operators expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#binary-array-valued-comparison-operators
type ExpressionBinaryArrayComparisonOperator struct {
	Type     ExpressionType                `json:"type" mapstructure:"type"`
	Operator BinaryArrayComparisonOperator `json:"operator" mapstructure:"operator"`
	Column   ComparisonTarget              `json:"column" mapstructure:"column"`
	Values   []ComparisonValue             `json:"values" mapstructure:"values"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionBinaryArrayComparisonOperator) GetType() ExpressionType {
	return exp.Type
}

// ExpressionExists is an object which represents an [EXISTS expression]
//
// [EXISTS expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#exists-expressions
type ExpressionExists struct {
	Type         ExpressionType     `json:"type" mapstructure:"type"`
	Where        Expression         `json:"where" mapstructure:"where"`
	InCollection ExistsInCollection `json:"in_collection" mapstructure:"in_collection"`
}

// GetType returns the expression type. Implement the Expression interface
func (exp ExpressionExists) GetType() ExpressionType {
	return exp.Type
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
func ParseAggregateType(input string) (*AggregateType, error) {
	if !Contains(enumValues_AggregateType, AggregateType(input)) {
		return nil, fmt.Errorf("failed to parse AggregateType, expect one of %v", enumValues_AggregateType)
	}
	result := AggregateType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// Aggregate represents an [aggregated query] object
//
// [aggregated query]: https://hasura.github.io/ndc-spec/specification/queries/aggregates.html
type Aggregate interface {
	GetType() AggregateType
}

// AggregateStarCount represents an aggregate object which counts all matched rows
type AggregateStarCount struct {
	Type AggregateType `json:"type" mapstructure:"type"`
}

// GetType gets the type of aggregate object
func (ag AggregateStarCount) GetType() AggregateType {
	return ag.Type
}

// AggregateSingleColumn represents an aggregate object which applies an aggregation function (as defined by the column's scalar type in the schema response) to a column.
type AggregateSingleColumn struct {
	Type AggregateType `json:"type" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" mapstructure:"column"`
	// Single column aggregate function name.
	Function string `json:"function" mapstructure:"function"`
}

// GetType gets the type of aggregate object
func (ag AggregateSingleColumn) GetType() AggregateType {
	return ag.Type
}

// AggregateColumnCount represents an aggregate object which count the number of rows with non-null values in the specified columns.
// If the distinct flag is set, then the count should only count unique non-null values of those columns.
type AggregateColumnCount struct {
	Type AggregateType `json:"type" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" mapstructure:"column"`
	// Whether or not only distinct items should be counted.
	Distinct bool `json:"distinct" mapstructure:"distinct"`
}

// GetType gets the type of aggregate object
func (ag AggregateColumnCount) GetType() AggregateType {
	return ag.Type
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
func ParseOrderByTargetType(input string) (*OrderByTargetType, error) {
	if !Contains(enumValues_OrderByTargetType, OrderByTargetType(input)) {
		return nil, fmt.Errorf("failed to parse OrderByTargetType, expect one of %v", enumValues_OrderByTargetType)
	}
	result := OrderByTargetType(input)

	return &result, nil
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

	*j = *value
	return nil
}

// OrderByTarget represents an [order_by field] of the Query object
//
// [order_by field]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html
type OrderByTarget interface {
	GetType() OrderByTargetType
}

// OrderByColumn represents an ordering object which compares the value in the selected column
type OrderByColumn struct {
	Type OrderByTargetType `json:"type" mapstructure:"type"`
	// The name of the column
	Column string `json:"column" mapstructure:"column"`
	// Any relationships to traverse to reach this column
	Path []PathElement `json:"path" mapstructure:"path"`
}

// GetType gets the type of order by
func (ob OrderByColumn) GetType() OrderByTargetType {
	return ob.Type
}

// OrderBySingleColumnAggregate An ordering of type [single_column_aggregate] orders rows by an aggregate computed over rows in some related collection.
// If the respective aggregates are incomparable, the ordering should continue to the next OrderByElement.
//
// [single_column_aggregate]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html#type-single_column_aggregate
type OrderBySingleColumnAggregate struct {
	Type OrderByTargetType `json:"type" mapstructure:"type"`
	// The column to apply the aggregation function to
	Column string `json:"column" mapstructure:"column"`
	// Single column aggregate function name.
	Function string `json:"function" mapstructure:"function"`
	// Non-empty collection of relationships to traverse
	Path []PathElement `json:"path" mapstructure:"path"`
}

// GetType gets the type of order by
func (ob OrderBySingleColumnAggregate) GetType() OrderByTargetType {
	return ob.Type
}

// OrderByStarCountAggregate An ordering of type [star_count_aggregate] orders rows by a count of rows in some related collection.
// If the respective aggregates are incomparable, the ordering should continue to the next OrderByElement.
//
// [star_count_aggregate]: https://hasura.github.io/ndc-spec/specification/queries/sorting.html#type-star_count_aggregate
type OrderByStarCountAggregate struct {
	Type OrderByTargetType `json:"type" mapstructure:"type"`
	// Non-empty collection of relationships to traverse
	Path []PathElement `json:"path" mapstructure:"path"`
}

// GetType gets the type of order by
func (ob OrderByStarCountAggregate) GetType() OrderByTargetType {
	return ob.Type
}
