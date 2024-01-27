package schema

import (
	"encoding/json"
	"errors"
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
