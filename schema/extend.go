package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/go-viper/mapstructure/v2"
)

var errTypeRequired = errors.New("type field is required")

// ArgumentType represents an argument type enum.
type ArgumentType string

const (
	ArgumentTypeLiteral  ArgumentType = "literal"
	ArgumentTypeVariable ArgumentType = "variable"
)

var enumValues_ArgumentType = []ArgumentType{
	ArgumentTypeLiteral,
	ArgumentTypeVariable,
}

// ParseArgumentType parses an argument type from string.
func ParseArgumentType(input string) (ArgumentType, error) {
	result := ArgumentType(input)
	if !result.IsValid() {
		return ArgumentType(""), fmt.Errorf("failed to parse ArgumentType, expect one of %v, got %s", enumValues_ArgumentType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
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
type Argument map[string]any

// ArgumentEncoder abstracts the interface for Argument.
type ArgumentEncoder interface {
	Encode() Argument
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Argument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawArgumentType, err := getStringValueByKey(raw, "type")
	if err != nil {
		return fmt.Errorf("type in Argument: %w", err)
	}

	if rawArgumentType == "" {
		return errors.New("field type in Argument: required")
	}

	argumentType, err := ParseArgumentType(rawArgumentType)
	if err != nil {
		return fmt.Errorf("field type in Argument: %w", err)
	}

	arg := map[string]any{
		"type": argumentType,
	}

	switch argumentType {
	case ArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			arg["value"] = value
		}
	case ArgumentTypeVariable:
		name, err := getStringValueByKey(raw, "name")
		if err != nil {
			return fmt.Errorf("field name in Argument: %w", err)
		}

		if name == "" {
			return errors.New("field name in Argument is required for variable type")
		}

		arg["name"] = name
	}

	*j = arg
	return nil
}

// Type gets the type enum of the current type.
func (j Argument) Type() (ArgumentType, error) {
	t, ok := j["type"]
	if !ok {
		return ArgumentType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseArgumentType(raw)
		if err != nil {
			return ArgumentType(""), err
		}
		return v, nil
	case ArgumentType:
		return raw, nil
	default:
		return ArgumentType(""), fmt.Errorf("invalid Field type: %+v", t)
	}
}

// AsLiteral converts the instance to ArgumentLiteral.
func (j Argument) AsLiteral() (*ArgumentLiteral, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ArgumentTypeLiteral {
		return nil, fmt.Errorf("invalid ArgumentLiteral type; expected: %s, got: %s", ArgumentTypeLiteral, t)
	}
	value := j["value"]

	return &ArgumentLiteral{
		Type:  t,
		Value: value,
	}, nil
}

// AsVariable converts the instance to ArgumentVariable.
func (j Argument) AsVariable() (*ArgumentVariable, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ArgumentTypeVariable {
		return nil, fmt.Errorf("invalid ArgumentVariable type; expected: %s, got: %s", ArgumentTypeVariable, t)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("ArgumentVariable.name: %w", err)
	}

	if name == "" {
		return nil, errors.New("ArgumentVariable.name is required")
	}

	return &ArgumentVariable{
		Type: t,
		Name: name,
	}, nil
}

// Interface converts the comparison value to its generic interface.
func (j Argument) Interface() ArgumentEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j Argument) InterfaceT() (ArgumentEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case ArgumentTypeLiteral:
		return j.AsLiteral()
	case ArgumentTypeVariable:
		return j.AsVariable()
	default:
		return nil, err
	}
}

// ArgumentLiteral represents the literal argument.
type ArgumentLiteral struct {
	Type  ArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Value any          `json:"value" yaml:"value" mapstructure:"value"`
}

// NewArgumentLiteral creates an argument with a literal value.
func NewArgumentLiteral(value any) *ArgumentLiteral {
	return &ArgumentLiteral{
		Type:  ArgumentTypeLiteral,
		Value: value,
	}
}

// Encode converts the instance to raw Field.
func (j ArgumentLiteral) Encode() Argument {
	return Argument{
		"type":  j.Type,
		"value": j.Value,
	}
}

// ArgumentVariable represents the variable argument.
type ArgumentVariable struct {
	Type ArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Name string       `json:"name" yaml:"name" mapstructure:"name"`
}

// NewArgumentVariable creates an argument with a variable name.
func NewArgumentVariable(name string) *ArgumentVariable {
	return &ArgumentVariable{
		Type: ArgumentTypeVariable,
		Name: name,
	}
}

// Encode converts the instance to raw Field.
func (j ArgumentVariable) Encode() Argument {
	return Argument{
		"type": j.Type,
		"name": j.Name,
	}
}

// RelationshipArgumentType represents a relationship argument type enum.
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

// ParseRelationshipArgumentType parses a relationship argument type from string.
func ParseRelationshipArgumentType(input string) (RelationshipArgumentType, error) {
	result := RelationshipArgumentType(input)
	if !result.IsValid() {
		return RelationshipArgumentType(""), fmt.Errorf("failed to parse RelationshipArgumentType, expect one of %v, got %s", enumValues_RelationshipArgumentType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
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

// RelationshipArgument is provided by reference to a variable or as a literal value.
type RelationshipArgument map[string]any

// RelationshipArgumentEncoder abstracts the interface for RelationshipArgument.
type RelationshipArgumentEncoder interface {
	Encode() RelationshipArgument
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationshipArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawArgumentType, err := getStringValueByKey(raw, "type")
	if err != nil {
		return fmt.Errorf("field type in Argument: %w", err)
	}

	if rawArgumentType == "" {
		return errors.New("field type in Argument: required")
	}

	argumentType, err := ParseRelationshipArgumentType(rawArgumentType)
	if err != nil {
		return fmt.Errorf("field type in Argument: %w", err)
	}

	result := map[string]any{
		"type": argumentType,
	}

	switch argumentType {
	case RelationshipArgumentTypeLiteral:
		if value, ok := raw["value"]; !ok {
			return errors.New("field value in Argument is required for literal type")
		} else {
			result["value"] = value
		}
	default:
		name, err := getStringValueByKey(raw, "name")
		if err != nil {
			return fmt.Errorf("field name in Argument: %w", err)
		}

		if name == "" {
			return fmt.Errorf("field name in Argument is required for %s type", rawArgumentType)
		}

		result["name"] = name
	}

	*j = result

	return nil
}

// Type gets the type enum of the current type.
func (j RelationshipArgument) Type() (RelationshipArgumentType, error) {
	t, ok := j["type"]
	if !ok {
		return RelationshipArgumentType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseRelationshipArgumentType(raw)
		if err != nil {
			return RelationshipArgumentType(""), err
		}
		return v, nil
	case RelationshipArgumentType:
		return raw, nil
	default:
		return RelationshipArgumentType(""), fmt.Errorf("invalid Field type: %+v", t)
	}
}

// AsLiteral converts the instance to RelationshipArgumentLiteral.
func (j RelationshipArgument) AsLiteral() (*RelationshipArgumentLiteral, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != RelationshipArgumentTypeLiteral {
		return nil, fmt.Errorf("invalid RelationshipArgumentLiteral type; expected: %s, got: %s", RelationshipArgumentTypeLiteral, t)
	}

	value := j["value"]

	return &RelationshipArgumentLiteral{
		Type:  t,
		Value: value,
	}, nil
}

// AsVariable converts the instance to RelationshipArgumentVariable.
func (j RelationshipArgument) AsVariable() (*RelationshipArgumentVariable, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != RelationshipArgumentTypeVariable {
		return nil, fmt.Errorf("invalid RelationshipArgumentVariable type; expected: %s, got: %s", RelationshipArgumentTypeVariable, t)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("RelationshipArgumentVariable.name, %w", err)
	}

	if name == "" {
		return nil, errors.New("RelationshipArgumentVariable.name is required")
	}
	return &RelationshipArgumentVariable{
		Type: t,
		Name: name,
	}, nil
}

// AsColumn converts the instance to RelationshipArgumentColumn.
func (j RelationshipArgument) AsColumn() (*RelationshipArgumentColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != RelationshipArgumentTypeColumn {
		return nil, fmt.Errorf("invalid RelationshipArgumentTypeColumn type; expected: %s, got: %s", RelationshipArgumentTypeColumn, t)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("RelationshipArgumentColumn.name: %w", err)
	}

	if name == "" {
		return nil, errors.New("RelationshipArgumentColumn.name is required")
	}

	return &RelationshipArgumentColumn{
		Type: t,
		Name: name,
	}, nil
}

// Interface converts the comparison value to its generic interface.
func (j RelationshipArgument) Interface() RelationshipArgumentEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j RelationshipArgument) InterfaceT() (RelationshipArgumentEncoder, error) {
	ty, err := j.Type()
	switch ty {
	case RelationshipArgumentTypeLiteral:
		return j.AsLiteral()
	case RelationshipArgumentTypeVariable:
		return j.AsVariable()
	case RelationshipArgumentTypeColumn:
		return j.AsColumn()
	default:
		return nil, err
	}
}

// RelationshipArgumentLiteral represents the literal relationship argument.
type RelationshipArgumentLiteral struct {
	Type  RelationshipArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Value any                      `json:"value" yaml:"value" mapstructure:"value"`
}

// NewRelationshipArgumentLiteral creates a RelationshipArgumentLiteral instance.
func NewRelationshipArgumentLiteral(value any) *RelationshipArgumentLiteral {
	return &RelationshipArgumentLiteral{
		Type:  RelationshipArgumentTypeLiteral,
		Value: value,
	}
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentLiteral) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type":  j.Type,
		"value": j.Value,
	}
}

// RelationshipArgumentColumn represents the column relationship argument.
type RelationshipArgumentColumn struct {
	Type RelationshipArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Name string                   `json:"name" yaml:"name" mapstructure:"name"`
}

// NewRelationshipArgumentColumn creates a RelationshipArgumentColumn instance.
func NewRelationshipArgumentColumn(name string) *RelationshipArgumentColumn {
	return &RelationshipArgumentColumn{
		Type: RelationshipArgumentTypeLiteral,
		Name: name,
	}
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentColumn) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type": j.Type,
		"name": j.Name,
	}
}

// RelationshipArgumentVariable represents the variable relationship argument.
type RelationshipArgumentVariable struct {
	Type RelationshipArgumentType `json:"type" yaml:"type" mapstructure:"type"`
	Name string                   `json:"name" yaml:"name" mapstructure:"name"`
}

// NewRelationshipArgumentVariable creates a RelationshipArgumentVariable instance.
func NewRelationshipArgumentVariable(name string) *RelationshipArgumentVariable {
	return &RelationshipArgumentVariable{
		Type: RelationshipArgumentTypeVariable,
		Name: name,
	}
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentVariable) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type": j.Type,
		"name": j.Name,
	}
}

// FieldType represents a field type.
type FieldType string

const (
	FieldTypeColumn       FieldType = "column"
	FieldTypeRelationship FieldType = "relationship"
)

var enumValues_FieldType = []FieldType{
	FieldTypeColumn,
	FieldTypeRelationship,
}

// ParseFieldType parses a field type from string.
func ParseFieldType(input string) (FieldType, error) {
	result := FieldType(input)
	if !result.IsValid() {
		return FieldType(""), fmt.Errorf("failed to parse FieldType, expect one of %v, got %s", enumValues_FieldType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
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

// Field represents a field.
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
		return fmt.Errorf("field type in Field: %w", err)
	}

	results := map[string]any{
		"type": fieldType,
	}

	switch fieldType {
	case FieldTypeColumn:
		column, err := unmarshalStringFromJsonMap(raw, "column", true)
		if err != nil {
			return fmt.Errorf("field column in Field: %w", err)
		}

		results["column"] = column

		// decode fields
		var fields NestedField
		rawFields, ok := raw["fields"]
		if ok && !isNullJSON(rawFields) {
			if err = json.Unmarshal(rawFields, &fields); err != nil {
				return fmt.Errorf("field fields in Field: %w", err)
			}
			results["fields"] = fields
		}
		var arguments map[string]Argument
		rawArguments, ok := raw["arguments"]
		if ok && !isNullJSON(rawArguments) {
			if err = json.Unmarshal(rawArguments, &arguments); err != nil {
				return fmt.Errorf("field arguments in Field: %w", err)
			}
			results["arguments"] = arguments
		}
	case FieldTypeRelationship:
		relationship, err := unmarshalStringFromJsonMap(raw, "relationship", true)
		if err != nil {
			return fmt.Errorf("field relationship in Field: %w", err)
		}
		results["relationship"] = relationship

		rawQuery, ok := raw["query"]
		if !ok {
			return errors.New("field query in Field: required")
		}
		var query Query
		if err = json.Unmarshal(rawQuery, &query); err != nil {
			return fmt.Errorf("field query in Field: %w", err)
		}
		results["query"] = query

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in Field: required")
		}

		var arguments map[string]RelationshipArgument
		if err = json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in Field: %w", err)
		}
		results["arguments"] = arguments
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
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

// AsColumn tries to convert the current type to ColumnField.
func (j Field) AsColumn() (*ColumnField, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != FieldTypeColumn {
		return nil, fmt.Errorf("invalid Field type; expected %s, got %s", FieldTypeColumn, t)
	}

	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("ColumnField.column: %w", err)
	}

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

	rawArguments, ok := j["arguments"]
	if ok && !isNil(rawArguments) {
		arguments, ok := rawArguments.(map[string]Argument)
		if !ok {
			return nil, fmt.Errorf("invalid ColumnField.arguments type; expected map[string]Argument, got %+v", rawArguments)
		}

		result.Arguments = arguments
	}

	return result, nil
}

// AsRelationship tries to convert the current type to RelationshipField.
func (j Field) AsRelationship() (*RelationshipField, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != FieldTypeRelationship {
		return nil, fmt.Errorf("invalid Field type; expected %s, got %s", FieldTypeRelationship, t)
	}
	relationship, err := getStringValueByKey(j, "relationship")
	if err != nil {
		return nil, fmt.Errorf("RelationshipField.relationship: %w", err)
	}

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

// Interface converts the comparison value to its generic interface.
func (j Field) Interface() FieldEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
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

// ColumnField represents a column field.
type ColumnField struct {
	Type FieldType `json:"type" yaml:"type" mapstructure:"type"`
	// Column name
	Column string `json:"column" yaml:"column" mapstructure:"column"`
	// When the type of the column is a (possibly-nullable) array or object,
	// the caller can request a subset of the complete column data, by specifying fields to fetch here.
	// If omitted, the column data will be fetched in full.
	Fields NestedField `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields"`

	Arguments map[string]Argument `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"fields"`
}

// Encode converts the instance to raw Field.
func (f ColumnField) Encode() Field {
	r := Field{
		"type":   f.Type,
		"column": f.Column,
	}

	if len(f.Fields) > 0 {
		r["fields"] = f.Fields
	}
	if len(f.Arguments) > 0 {
		r["arguments"] = f.Arguments
	}
	return r
}

// NewColumnField creates a new ColumnField instance.
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

// NewColumnFieldWithArguments creates a new ColumnField instance with arguments.
func NewColumnFieldWithArguments(column string, fields NestedFieldEncoder, arguments map[string]Argument) *ColumnField {
	cf := NewColumnField(column, fields)
	cf.Arguments = arguments
	return cf
}

// RelationshipField represents a relationship field.
type RelationshipField struct {
	Type FieldType `json:"type" yaml:"type" mapstructure:"type"`
	// The relationship query
	Query Query `json:"query" yaml:"query" mapstructure:"query"`
	// The name of the relationship to follow for the subquery
	Relationship string `json:"relationship" yaml:"relationship" mapstructure:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
}

// Encode converts the instance to raw Field.
func (f RelationshipField) Encode() Field {
	return Field{
		"type":         f.Type,
		"query":        f.Query,
		"relationship": f.Relationship,
		"arguments":    f.Arguments,
	}
}

// NewRelationshipField creates a new RelationshipField instance.
func NewRelationshipField(query Query, relationship string, arguments map[string]RelationshipArgument) *RelationshipField {
	return &RelationshipField{
		Type:         FieldTypeRelationship,
		Query:        query,
		Relationship: relationship,
		Arguments:    arguments,
	}
}

// ComparisonTargetType represents comparison target enums.
type ComparisonTargetType string

const (
	ComparisonTargetTypeColumn               ComparisonTargetType = "column"
	ComparisonTargetTypeRootCollectionColumn ComparisonTargetType = "root_collection_column"
)

var enumValues_ComparisonTargetType = []ComparisonTargetType{
	ComparisonTargetTypeColumn,
	ComparisonTargetTypeRootCollectionColumn,
}

// ParseComparisonTargetType parses a comparison target type argument type from string.
func ParseComparisonTargetType(input string) (ComparisonTargetType, error) {
	result := ComparisonTargetType(input)
	if !result.IsValid() {
		return ComparisonTargetType(""), fmt.Errorf("failed to parse ComparisonTargetType, expect one of %v, got: %s", enumValues_ComparisonTargetType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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

// ComparisonTarget represents a comparison target object.
type ComparisonTarget struct {
	Type      ComparisonTargetType `json:"type" yaml:"type" mapstructure:"type"`
	Name      string               `json:"name" yaml:"name" mapstructure:"name"`
	Path      []PathElement        `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path"`
	FieldPath []string             `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// NewComparisonTargetColumn creates a ComparisonTarget with column type.
func NewComparisonTargetColumn(name string, fieldPath []string, path []PathElement) *ComparisonTarget {
	return &ComparisonTarget{
		Type:      ComparisonTargetTypeColumn,
		Name:      name,
		Path:      path,
		FieldPath: fieldPath,
	}
}

// NewComparisonTargetRootCollectionColumn creates a ComparisonTarget with root_collection_column type.
func NewComparisonTargetRootCollectionColumn(name string, fieldPath []string) *ComparisonTarget {
	return &ComparisonTarget{
		Type:      ComparisonTargetTypeRootCollectionColumn,
		Name:      name,
		FieldPath: fieldPath,
	}
}

// ExpressionType represents the filtering expression enums.
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

// ParseExpressionType parses an expression type argument type from string.
func ParseExpressionType(input string) (ExpressionType, error) {
	result := ExpressionType(input)
	if !result.IsValid() {
		return ExpressionType(""), fmt.Errorf("failed to parse ExpressionType, expect one of %v, got %s", enumValues_ExpressionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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

// ComparisonValueType represents a comparison value type enum.
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

// ParseComparisonValueType parses a comparison value type from string.
func ParseComparisonValueType(input string) (ComparisonValueType, error) {
	result := ComparisonValueType(input)
	if !result.IsValid() {
		return ComparisonValueType(""), fmt.Errorf("failed to parse ComparisonValueType, expect one of %v, got %s", enumValues_ComparisonValueType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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

// ComparisonValue represents a raw comparison value object with validation.
type ComparisonValue map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonValue) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	return j.FromValue(raw)
}

// FromValue decodes values from any map.
func (j *ComparisonValue) FromValue(input map[string]any) error {
	rawType, err := getStringValueByKey(input, "type")
	if err != nil {
		return fmt.Errorf("field type in ComparisonValue: %w", err)
	}

	ty, err := ParseComparisonValueType(rawType)
	if err != nil {
		return fmt.Errorf("field type in ComparisonValue: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	switch ty {
	case ComparisonValueTypeVariable:
		name, err := getStringValueByKey(input, "name")
		if err != nil {
			return fmt.Errorf("field name in ComparisonValue: %w", err)
		}

		if name == "" {
			return errors.New("field name in ComparisonValue is required for variable type")
		}

		result["name"] = name
	case ComparisonValueTypeColumn:
		rawColumn, ok := input["column"]
		if !ok {
			return errors.New("field column in ComparisonValue is required for column type")
		}

		var column ComparisonTarget
		if err := mapstructure.Decode(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in ComparisonValue: %w", err)
		}

		result["column"] = column
	case ComparisonValueTypeScalar:
		value, ok := input["value"]
		if !ok {
			return errors.New("field value in ComparisonValue is required for scalar type")
		}

		result["value"] = value
	}

	*j = result

	return nil
}

// GetType gets the type of comparison value.
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

// AsScalar tries to convert the comparison value to scalar.
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

// AsColumn tries to convert the comparison value to column.
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

// AsVariable tries to convert the comparison value to column.
func (cv ComparisonValue) AsVariable() (*ComparisonValueVariable, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != ComparisonValueTypeVariable {
		return nil, fmt.Errorf("invalid ComparisonValue type; expected %s, got %s", ComparisonValueTypeVariable, ty)
	}

	name, err := getStringValueByKey(cv, "name")
	if err != nil {
		return nil, fmt.Errorf("ComparisonValueVariable.name: %w", err)
	}

	if name == "" {
		return nil, errors.New("ComparisonValueVariable.name is required")
	}

	return &ComparisonValueVariable{
		Type: ty,
		Name: name,
	}, nil
}

// Interface converts the comparison value to its generic interface.
func (cv ComparisonValue) Interface() ComparisonValueEncoder {
	result, _ := cv.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
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

// ComparisonValueEncoder represents a comparison value encoder interface.
type ComparisonValueEncoder interface {
	Encode() ComparisonValue
}

// ComparisonValueColumn represents a comparison value with column type.
type ComparisonValueColumn struct {
	Type   ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Column ComparisonTarget    `json:"column" yaml:"column" mapstructure:"column"`
}

// NewComparisonValueColumn creates a new ComparisonValueColumn instance.
func NewComparisonValueColumn(column ComparisonTarget) *ComparisonValueColumn {
	return &ComparisonValueColumn{
		Type:   ComparisonValueTypeColumn,
		Column: column,
	}
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueColumn) Encode() ComparisonValue {
	return map[string]any{
		"type":   cv.Type,
		"column": cv.Column,
	}
}

// ComparisonValueScalar represents a comparison value with scalar type.
type ComparisonValueScalar struct {
	Type  ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Value any                 `json:"value" yaml:"value" mapstructure:"value"`
}

// NewComparisonValueScalar creates a new ComparisonValueScalar instance.
func NewComparisonValueScalar(value any) *ComparisonValueScalar {
	return &ComparisonValueScalar{
		Type:  ComparisonValueTypeScalar,
		Value: value,
	}
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueScalar) Encode() ComparisonValue {
	return map[string]any{
		"type":  cv.Type,
		"value": cv.Value,
	}
}

// ComparisonValueVariable represents a comparison value with variable type.
type ComparisonValueVariable struct {
	Type ComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Name string              `json:"name" yaml:"name" mapstructure:"name"`
}

// NewComparisonValueVariable creates a new ComparisonValueVariable instance.
func NewComparisonValueVariable(name string) *ComparisonValueVariable {
	return &ComparisonValueVariable{
		Type: ComparisonValueTypeVariable,
		Name: name,
	}
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueVariable) Encode() ComparisonValue {
	return map[string]any{
		"type": cv.Type,
		"name": cv.Name,
	}
}

// ExistsInCollectionType represents an exists in collection type enum.
type ExistsInCollectionType string

const (
	ExistsInCollectionTypeRelated          ExistsInCollectionType = "related"
	ExistsInCollectionTypeUnrelated        ExistsInCollectionType = "unrelated"
	ExistsInCollectionTypeNestedCollection ExistsInCollectionType = "nested_collection"
)

var enumValues_ExistsInCollectionType = []ExistsInCollectionType{
	ExistsInCollectionTypeRelated,
	ExistsInCollectionTypeUnrelated,
	ExistsInCollectionTypeNestedCollection,
}

// ParseExistsInCollectionType parses a comparison value type from string.
func ParseExistsInCollectionType(input string) (ExistsInCollectionType, error) {
	result := ExistsInCollectionType(input)
	if !result.IsValid() {
		return result, fmt.Errorf("failed to parse ExistsInCollectionType, expect one of %v, got %s", enumValues_ExistsInCollectionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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

// ExistsInCollection represents an Exists In Collection object.
type ExistsInCollection map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExistsInCollection) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	return j.FromValue(raw)
}

// FromValue decodes values from any map
func (j *ExistsInCollection) FromValue(input map[string]any) error {
	rawType, err := getStringValueByKey(input, "type")
	if err != nil {
		return fmt.Errorf("field type in ExistsInCollection: %w", err)
	}

	ty, err := ParseExistsInCollectionType(rawType)
	if err != nil {
		return fmt.Errorf("field type in ExistsInCollection: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	rawArguments, ok := input["arguments"]
	if ok {
		var arguments map[string]RelationshipArgument
		if err := mapstructure.Decode(rawArguments, &arguments); err != nil {
			return fmt.Errorf("field arguments in ExistsInCollection: %w", err)
		}

		result["arguments"] = arguments
	} else if ty != ExistsInCollectionTypeNestedCollection {
		return fmt.Errorf("field arguments in ExistsInCollection is required for %s type", ty)
	}

	switch ty {
	case ExistsInCollectionTypeRelated:
		relationship, err := getStringValueByKey(input, "relationship")
		if err != nil {
			return fmt.Errorf("field name in ExistsInCollection: %w", err)
		}

		if relationship == "" {
			return errors.New("field relationship in ExistsInCollection is required for related type")
		}

		result["relationship"] = relationship
	case ExistsInCollectionTypeUnrelated:
		collection, err := getStringValueByKey(input, "collection")
		if err != nil {
			return fmt.Errorf("field collection in ExistsInCollection: %w", err)
		}

		if collection == "" {
			return errors.New("field collection in ExistsInCollection is required for unrelated type")
		}

		result["collection"] = collection
	case ExistsInCollectionTypeNestedCollection:
		columnName, err := getStringValueByKey(input, "column_name")
		if err != nil {
			return fmt.Errorf("field column_name in ExistsInCollection: %w", err)
		}

		if columnName == "" {
			return errors.New("field column_name in ExistsInCollection is required for nested_collection type")
		}

		result["column_name"] = columnName

		rawFieldPath, ok := input["field_path"]
		if ok {
			var fieldPath []string

			if err := mapstructure.Decode(rawFieldPath, &fieldPath); err != nil {
				return fmt.Errorf("field field_path in ExistsInCollection: %w", err)
			}

			result["field_path"] = fieldPath
		}
	}

	*j = result

	return nil
}

// Type gets the type enum of the current type.
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

// AsRelated tries to convert the instance to related type.
func (j ExistsInCollection) AsRelated() (*ExistsInCollectionRelated, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExistsInCollectionTypeRelated {
		return nil, fmt.Errorf("invalid ExistsInCollection type; expected: %s, got: %s", ExistsInCollectionTypeRelated, t)
	}

	relationship, err := getStringValueByKey(j, "relationship")
	if err != nil {
		return nil, fmt.Errorf("ExistsInCollectionRelated.relationship: %w", err)
	}

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

// AsRelated tries to convert the instance to unrelated type.
func (j ExistsInCollection) AsUnrelated() (*ExistsInCollectionUnrelated, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExistsInCollectionTypeUnrelated {
		return nil, fmt.Errorf("invalid ExistsInCollection type; expected: %s, got: %s", ExistsInCollectionTypeUnrelated, t)
	}

	collection, err := getStringValueByKey(j, "collection")
	if err != nil {
		return nil, fmt.Errorf("ExistsInCollectionUnrelated.collection: %w", err)
	}

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

// AsNestedCollection tries to convert the instance to nested_collection type.
func (j ExistsInCollection) AsNestedCollection() (*ExistsInCollectionNestedCollection, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != ExistsInCollectionTypeNestedCollection {
		return nil, fmt.Errorf("invalid ExistsInCollection type; expected: %s, got: %s", ExistsInCollectionTypeNestedCollection, t)
	}

	columnName, err := getStringValueByKey(j, "column_name")
	if err != nil {
		return nil, fmt.Errorf("ExistsInCollectionNestedCollection.column_name: %w", err)
	}

	if columnName == "" {
		return nil, errors.New("ExistsInCollectionNestedCollection.column_name is required")
	}

	var args map[string]RelationshipArgument
	rawArgs, ok := j["arguments"]
	if ok && rawArgs != nil {
		args, ok = rawArgs.(map[string]RelationshipArgument)
		if !ok {
			return nil, fmt.Errorf("invalid ExistsInCollectionNestedCollection.arguments type; expected: map[string]RelationshipArgument, got: %+v", rawArgs)
		}
	}

	result := &ExistsInCollectionNestedCollection{
		Type:       t,
		ColumnName: columnName,
		Arguments:  args,
	}

	rawFieldPath, ok := j["field_path"]
	if ok && rawFieldPath != nil {
		fieldPath, ok := rawFieldPath.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid ExistsInCollectionNestedCollection.fieldPath type; expected: []string, got: %+v", rawArgs)
		}

		result.FieldPath = fieldPath
	}

	return result, nil
}

// Interface tries to convert the instance to the ExistsInCollectionEncoder interface.
func (j ExistsInCollection) Interface() ExistsInCollectionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to the ExistsInCollectionEncoder interface safely with explicit error.
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
	case ExistsInCollectionTypeNestedCollection:
		return j.AsNestedCollection()
	default:
		return nil, fmt.Errorf("invalid ExistsInCollection type: %s", t)
	}
}

// ExistsInCollectionEncoder abstracts the ExistsInCollection serialization interface.
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

// NewExistsInCollectionRelated creates an ExistsInCollectionRelated instance.
func NewExistsInCollectionRelated(relationship string, arguments map[string]RelationshipArgument) *ExistsInCollectionRelated {
	return &ExistsInCollectionRelated{
		Type:         ExistsInCollectionTypeRelated,
		Relationship: relationship,
		Arguments:    arguments,
	}
}

// Encode converts the instance to its raw type.
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

// NewExistsInCollectionUnrelated creates an ExistsInCollectionUnrelated instance.
func NewExistsInCollectionUnrelated(collection string, arguments map[string]RelationshipArgument) *ExistsInCollectionUnrelated {
	return &ExistsInCollectionUnrelated{
		Type:       ExistsInCollectionTypeUnrelated,
		Collection: collection,
		Arguments:  arguments,
	}
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionUnrelated) Encode() ExistsInCollection {
	return ExistsInCollection{
		"type":       ei.Type,
		"collection": ei.Collection,
		"arguments":  ei.Arguments,
	}
}

// ExistsInCollectionNestedCollection represents [nested collections] expression.
//
// [nested collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#nested-collections
type ExistsInCollectionNestedCollection struct {
	Type ExistsInCollectionType `json:"type" yaml:"type" mapstructure:"type"`
	// The name of column
	ColumnName string `json:"column_name" yaml:"column_name" mapstructure:"column_name"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments,omitempty" yaml:"arguments,omitempty" mapstructure:"arguments"`
	// Path to a nested collection via object columns
	FieldPath []string `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// NewExistsInCollectionNestedCollection creates an ExistsInCollectionNestedCollection instance.
func NewExistsInCollectionNestedCollection(columnName string, arguments map[string]RelationshipArgument, fieldPath []string) *ExistsInCollectionNestedCollection {
	return &ExistsInCollectionNestedCollection{
		Type:       ExistsInCollectionTypeNestedCollection,
		ColumnName: columnName,
		Arguments:  arguments,
		FieldPath:  fieldPath,
	}
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionNestedCollection) Encode() ExistsInCollection {
	result := ExistsInCollection{
		"type":        ei.Type,
		"column_name": ei.ColumnName,
	}
	if len(ei.Arguments) > 0 {
		result["arguments"] = ei.Arguments
	}
	if len(ei.FieldPath) > 0 {
		result["field_path"] = ei.FieldPath
	}
	return result
}

// Expression represents the query expression object.
type Expression map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Expression) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	return j.FromValue(raw)
}

// FromValue decodes values from any map
func (j *Expression) FromValue(input map[string]any) error {
	rawType, err := getStringValueByKey(input, "type")
	if err != nil {
		return fmt.Errorf("field type in Expression: %w", err)
	}

	ty, err := ParseExpressionType(rawType)
	if err != nil {
		return fmt.Errorf("field type in Expression: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	switch ty {
	case ExpressionTypeAnd, ExpressionTypeOr:
		rawExpressions, ok := input["expressions"]
		if !ok || rawExpressions == nil {
			return fmt.Errorf("field expressions in Expression is required for '%s' type", ty)
		}

		rawExpressionsArray, ok := rawExpressions.([]any)
		if !ok {
			return fmt.Errorf("field expressions in Expression: expected array, got %v", rawExpressions)
		}

		expressions := []Expression{}
		for i, rawItem := range rawExpressionsArray {
			if rawItem == nil {
				continue
			}

			itemMap, ok := rawItem.(map[string]any)
			if !ok {
				return fmt.Errorf("field expressions[%d] in Expression: expected array, got %v", i, rawExpressions)
			}

			if itemMap == nil {
				continue
			}

			expr := Expression{}
			if err := expr.FromValue(itemMap); err != nil {
				return fmt.Errorf("field expressions in Expression: %w", err)
			}

			expressions = append(expressions, expr)
		}

		result["expressions"] = expressions
	case ExpressionTypeNot:
		rawExpression, ok := input["expression"]
		if !ok || rawExpression == nil {
			return fmt.Errorf("field expressions in Expression is required for '%s' type", ty)
		}

		exprMap, ok := rawExpression.(map[string]any)
		if !ok {
			return fmt.Errorf("field expressions in Expression is required for '%s' type", ty)
		}

		var expression Expression
		if err := expression.FromValue(exprMap); err != nil {
			return fmt.Errorf("field expression in Expression: %w", err)
		}

		result["expression"] = expression
	case ExpressionTypeUnaryComparisonOperator, ExpressionTypeBinaryComparisonOperator:
		rawOperator, err := getStringValueByKey(input, "operator")
		if err != nil {
			return fmt.Errorf("field operator in Expression: %w", err)
		}

		if rawOperator == "" {
			return fmt.Errorf("field operator in Expression is required for '%s' type", ty)
		}

		result["operator"] = rawOperator

		rawColumn, ok := input["column"]
		if !ok {
			return fmt.Errorf("field column in Expression is required for '%s' type", ty)
		}

		var column ComparisonTarget
		if err := mapstructure.Decode(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Expression: %w", err)
		}

		result["column"] = column

		if ty != ExpressionTypeBinaryComparisonOperator {
			break
		}

		rawValue, ok := input["value"]
		if !ok || rawValue == nil {
			return fmt.Errorf("field value in Expression is required for '%s' type", ty)
		}

		rawValueMap, ok := rawValue.(map[string]any)
		if !ok {
			return fmt.Errorf("field value in Expression: expected map, got %v", rawValue)
		}

		var value ComparisonValue
		if err := value.FromValue(rawValueMap); err != nil {
			return fmt.Errorf("field value in Expression: %w", err)
		}

		result["value"] = value
	case ExpressionTypeExists:
		rawPredicate, ok := input["predicate"]
		if ok && rawPredicate != nil {
			rawPredicateMap, ok := rawPredicate.(map[string]any)
			if !ok {
				return fmt.Errorf("field predicate in Expression: expected map, got %v", rawPredicate)
			}

			predicate := Expression{}
			if err := predicate.FromValue(rawPredicateMap); err != nil {
				return fmt.Errorf("field predicate in Expression: %w", err)
			}

			result["predicate"] = predicate
		}

		rawInCollection, ok := input["in_collection"]
		if !ok || rawInCollection == nil {
			return fmt.Errorf("field in_collection in Expression is required for '%s' type", ty)
		}

		rawInCollectionMap, ok := rawInCollection.(map[string]any)
		if !ok || rawInCollectionMap == nil {
			return fmt.Errorf("field in_collection in Expression is required for '%s' type", ty)
		}

		var inCollection ExistsInCollection
		if err := inCollection.FromValue(rawInCollectionMap); err != nil {
			return fmt.Errorf("field in_collection in Expression: %w", err)
		}

		result["in_collection"] = inCollection
	}

	*j = result

	return nil
}

// Type gets the type enum of the current type.
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

// AsAnd tries to convert the instance to and type.
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
		if err := mapstructure.Decode(rawExpressions, &expressions); err != nil {
			return nil, fmt.Errorf("invalid ExpressionAnd.expression type; expected: []Expression, got: %v", rawExpressions)
		}
	}

	return &ExpressionAnd{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsOr tries to convert the instance to ExpressionOr instance.
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
		if err := mapstructure.Decode(rawExpressions, &expressions); err != nil {
			return nil, fmt.Errorf("invalid ExpressionOr.expression type; expected: []Expression, got: %v", rawExpressions)
		}
	}

	return &ExpressionOr{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsNot tries to convert the instance to ExpressionNot instance.
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
		if err := mapstructure.Decode(rawExpression, &expression); err != nil {
			return nil, fmt.Errorf("invalid ExpressionNot.expression type; expected: Expression, got: %+v", rawExpression)
		}
	}

	return &ExpressionNot{
		Type:       t,
		Expression: expression,
	}, nil
}

// AsUnaryComparisonOperator tries to convert the instance to ExpressionUnaryComparisonOperator instance.
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

// AsBinaryComparisonOperator tries to convert the instance to ExpressionBinaryComparisonOperator instance.
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

	operator, err := getStringValueByKey(j, "operator")
	if err != nil {
		return nil, fmt.Errorf("invalid ExpressionBinaryComparisonOperator.opeartor: %w", err)
	}

	return &ExpressionBinaryComparisonOperator{
		Type:     t,
		Operator: operator,
		Column:   column,
		Value:    value,
	}, nil
}

// AsExists tries to convert the instance to ExpressionExists instance.
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

// Interface tries to convert the instance to the ExpressionEncoder interface.
func (j Expression) Interface() ExpressionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to the ExpressionEncoder interface safely with explicit error.
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

// ExpressionEncoder abstracts the expression encoder interface.
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

// NewExpressionAnd creates an ExpressionAnd instance.
func NewExpressionAnd(expressions ...ExpressionEncoder) *ExpressionAnd {
	exprs := make([]Expression, len(expressions))
	for i, expr := range expressions {
		if expr == nil {
			continue
		}
		exprs[i] = expr.Encode()
	}
	return &ExpressionAnd{
		Type:        ExpressionTypeAnd,
		Expressions: exprs,
	}
}

// Encode converts the instance to a raw Expression.
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

// NewExpressionOr creates an ExpressionOr instance.
func NewExpressionOr(expressions ...ExpressionEncoder) *ExpressionOr {
	exprs := make([]Expression, len(expressions))
	for i, expr := range expressions {
		if expr == nil {
			continue
		}
		exprs[i] = expr.Encode()
	}
	return &ExpressionOr{
		Type:        ExpressionTypeOr,
		Expressions: exprs,
	}
}

// Encode converts the instance to a raw Expression.
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

// NewExpressionNot creates an ExpressionNot instance.
func NewExpressionNot(expression ExpressionEncoder) *ExpressionNot {
	result := &ExpressionNot{
		Type: ExpressionTypeNot,
	}
	if expression != nil {
		result.Expression = expression.Encode()
	}
	return result
}

// Encode converts the instance to a raw Expression.
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

// NewExpressionUnaryComparisonOperator creates an ExpressionUnaryComparisonOperator instance.
func NewExpressionUnaryComparisonOperator(column ComparisonTarget, operator UnaryComparisonOperator) *ExpressionUnaryComparisonOperator {
	return &ExpressionUnaryComparisonOperator{
		Type:     ExpressionTypeUnaryComparisonOperator,
		Column:   column,
		Operator: operator,
	}
}

// Encode converts the instance to a raw Expression.
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

// NewExpressionBinaryComparisonOperator creates an ExpressionBinaryComparisonOperator instance.
func NewExpressionBinaryComparisonOperator(column ComparisonTarget, operator string, value ComparisonValueEncoder) *ExpressionBinaryComparisonOperator {
	result := &ExpressionBinaryComparisonOperator{
		Type:     ExpressionTypeBinaryComparisonOperator,
		Column:   column,
		Operator: operator,
	}
	if value != nil {
		result.Value = value.Encode()
	}
	return result
}

// Encode converts the instance to a raw Expression.
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

// NewExpressionExists creates an ExpressionExists instance.
func NewExpressionExists(predicate ExpressionEncoder, inCollection ExistsInCollectionEncoder) *ExpressionExists {
	result := &ExpressionExists{
		Type:         ExpressionTypeExists,
		InCollection: inCollection.Encode(),
	}
	if predicate != nil {
		result.Predicate = predicate.Encode()
	}
	if inCollection != nil {
		result.InCollection = inCollection.Encode()
	}
	return result
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionExists) Encode() Expression {
	return Expression{
		"type":          exp.Type,
		"predicate":     exp.Predicate,
		"in_collection": exp.InCollection,
	}
}

// AggregateType represents an aggregate type.
type AggregateType string

const (
	// AggregateTypeColumnCount aggregates count the number of rows with non-null values in the specified columns.
	// If the distinct flag is set, then the count should only count unique non-null values of those columns,.
	AggregateTypeColumnCount AggregateType = "column_count"
	// AggregateTypeSingleColumn aggregates apply an aggregation function (as defined by the column's scalar type in the schema response) to a column.
	AggregateTypeSingleColumn AggregateType = "single_column"
	// AggregateTypeStarCount aggregates count all matched rows.
	AggregateTypeStarCount AggregateType = "star_count"
)

var enumValues_AggregateType = []AggregateType{
	AggregateTypeColumnCount,
	AggregateTypeSingleColumn,
	AggregateTypeStarCount,
}

// ParseAggregateType parses an aggregate type argument type from string.
func ParseAggregateType(input string) (AggregateType, error) {
	result := AggregateType(input)
	if !result.IsValid() {
		return AggregateType(""), fmt.Errorf("failed to parse AggregateType, expect one of %v, got %s", enumValues_AggregateType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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
		return fmt.Errorf("field type in Aggregate: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	var fieldPath []string
	if ty != AggregateTypeStarCount {
		rawFieldPath, ok := raw["field_path"]
		if ok {
			if err := json.Unmarshal(rawFieldPath, &fieldPath); err != nil {
				return fmt.Errorf("field field_path in Aggregate: %w", err)
			}
		}
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
			return fmt.Errorf("field column in Aggregate: %w", err)
		}
		result["column"] = column

		rawFunction, ok := raw["function"]
		if !ok {
			return errors.New("field function in Aggregate is required for single_column type")
		}
		var function string
		if err := json.Unmarshal(rawFunction, &function); err != nil {
			return fmt.Errorf("field function in Aggregate: %w", err)
		}
		result["function"] = function
		if fieldPath != nil {
			result["field_path"] = fieldPath
		}
	case AggregateTypeColumnCount:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in Aggregate is required for column_count type")
		}
		var column string
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in Aggregate: %w", err)
		}
		result["column"] = column

		rawDistinct, ok := raw["distinct"]
		if !ok {
			return errors.New("field distinct in Aggregate is required for column_count type")
		}
		var distinct bool
		if err := json.Unmarshal(rawDistinct, &distinct); err != nil {
			return fmt.Errorf("field distinct in Aggregate: %w", err)
		}
		result["distinct"] = distinct
		if fieldPath != nil {
			result["field_path"] = fieldPath
		}
	}

	*j = result
	return nil
}

func (j Aggregate) getFieldPath() ([]string, error) {
	rawFieldPath, ok := j["field_path"]
	if !ok {
		return nil, nil
	}
	fieldPath, ok := rawFieldPath.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid AggregateColumnCount.field_path type; expected string slice, got %+v", rawFieldPath)
	}
	return fieldPath, nil
}

// Type gets the type enum of the current type.
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

// AsStarCount tries to convert the instance to AggregateStarCount type.
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

// AsSingleColumn tries to convert the instance to AggregateSingleColumn type.
func (j Aggregate) AsSingleColumn() (*AggregateSingleColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateTypeSingleColumn {
		return nil, fmt.Errorf("invalid Aggregate type; expected: %s, got: %s", AggregateTypeSingleColumn, t)
	}

	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("AggregateSingleColumn.column: %w", err)
	}

	if column == "" {
		return nil, errors.New("AggregateSingleColumn.column is required")
	}

	function, err := getStringValueByKey(j, "function")
	if err != nil {
		return nil, fmt.Errorf("AggregateSingleColumn.function: %w", err)
	}

	if function == "" {
		return nil, errors.New("AggregateSingleColumn.function is required")
	}

	fieldPath, err := j.getFieldPath()
	if err != nil {
		return nil, err
	}

	return &AggregateSingleColumn{
		Type:      t,
		Column:    column,
		Function:  function,
		FieldPath: fieldPath,
	}, nil
}

// AsColumnCount tries to convert the instance to AggregateColumnCount type.
func (j Aggregate) AsColumnCount() (*AggregateColumnCount, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateTypeColumnCount {
		return nil, fmt.Errorf("invalid Aggregate type; expected: %s, got: %s", AggregateTypeColumnCount, t)
	}

	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("Aggregate.column: %w", err)
	}

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

	fieldPath, err := j.getFieldPath()
	if err != nil {
		return nil, err
	}

	return &AggregateColumnCount{
		Type:      t,
		Column:    column,
		Distinct:  distinct,
		FieldPath: fieldPath,
	}, nil
}

// Interface tries to convert the instance to AggregateEncoder interface.
func (j Aggregate) Interface() AggregateEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to AggregateEncoder interface safely with explicit error.
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

// AggregateEncoder abstracts the serialization interface for Aggregate.
type AggregateEncoder interface {
	Encode() Aggregate
}

// AggregateStarCount represents an aggregate object which counts all matched rows.
type AggregateStarCount struct {
	Type AggregateType `json:"type" mapstructure:"type"`
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateStarCount) Encode() Aggregate {
	return Aggregate{
		"type": ag.Type,
	}
}

// NewAggregateStarCount creates a new AggregateStarCount instance.
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
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateSingleColumn) Encode() Aggregate {
	result := Aggregate{
		"type":     ag.Type,
		"column":   ag.Column,
		"function": ag.Function,
	}

	if ag.FieldPath != nil {
		result["field_path"] = ag.FieldPath
	}

	return result
}

// NewAggregateSingleColumn creates a new AggregateSingleColumn instance.
func NewAggregateSingleColumn(column string, function string, fieldPath []string) *AggregateSingleColumn {
	return &AggregateSingleColumn{
		Type:      AggregateTypeSingleColumn,
		Column:    column,
		Function:  function,
		FieldPath: fieldPath,
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
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateColumnCount) Encode() Aggregate {
	result := Aggregate{
		"type":     ag.Type,
		"column":   ag.Column,
		"distinct": ag.Distinct,
	}
	if ag.FieldPath != nil {
		result["field_path"] = ag.FieldPath
	}

	return result
}

// NewAggregateColumnCount creates a new AggregateColumnCount instance.
func NewAggregateColumnCount(column string, distinct bool, fieldPath []string) *AggregateColumnCount {
	return &AggregateColumnCount{
		Type:      AggregateTypeColumnCount,
		Column:    column,
		Distinct:  distinct,
		FieldPath: fieldPath,
	}
}

// OrderByTargetType represents a ordering target type.
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

// ParseOrderByTargetType parses a ordering target type argument type from string.
func ParseOrderByTargetType(input string) (OrderByTargetType, error) {
	result := OrderByTargetType(input)
	if !result.IsValid() {
		return OrderByTargetType(""), fmt.Errorf("failed to parse OrderByTargetType, expect one of %v, got %s", enumValues_OrderByTargetType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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
		return fmt.Errorf("field type in OrderByTarget: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	rawPath, ok := raw["path"]
	if !ok {
		return fmt.Errorf("field path in OrderByTarget is required for `%s` type", ty)
	}
	var pathElem []PathElement
	if err := json.Unmarshal(rawPath, &pathElem); err != nil {
		return fmt.Errorf("field path in OrderByTarget: %w", err)
	}
	result["path"] = pathElem

	if ty == OrderByTargetTypeColumn || ty == OrderByTargetTypeSingleColumnAggregate {
		rawFieldPath, ok := raw["field_path"]
		var fieldPath []string
		if ok {
			if err := json.Unmarshal(rawFieldPath, &fieldPath); err != nil {
				return fmt.Errorf("field field_path in OrderByTarget: %w", err)
			}
			result["field_path"] = fieldPath
		}
	}

	switch ty {
	case OrderByTargetTypeColumn:
		rawName, ok := raw["name"]
		if !ok {
			return errors.New("field name in OrderByTarget is required for column type")
		}
		var name string
		if err := json.Unmarshal(rawName, &name); err != nil {
			return fmt.Errorf("field name in OrderByTarget: %w", err)
		}
		result["name"] = name

	case OrderByTargetTypeSingleColumnAggregate:
		rawColumn, ok := raw["column"]
		if !ok {
			return errors.New("field column in OrderByTarget is required for single_column_aggregate type")
		}
		var column string
		if err := json.Unmarshal(rawColumn, &column); err != nil {
			return fmt.Errorf("field column in OrderByTarget: %w", err)
		}
		result["column"] = column

		rawFunction, ok := raw["function"]
		if !ok {
			return errors.New("field function in OrderByTarget is required for single_column_aggregate type")
		}
		var function string
		if err := json.Unmarshal(rawFunction, &function); err != nil {
			return fmt.Errorf("field function in OrderByTarget: %w", err)
		}
		result["function"] = function
	case OrderByTargetTypeStarCountAggregate:
	}
	*j = result
	return nil
}

func (j OrderByTarget) getFieldPath() ([]string, error) {
	rawFieldPath, ok := j["field_path"]
	if !ok {
		return nil, nil
	}
	fieldPath, ok := rawFieldPath.([]string)
	if !ok {
		return nil, fmt.Errorf("invalid OrderByTarget.field_path type; expected string slice, got %+v", rawFieldPath)
	}
	return fieldPath, nil
}

func (j OrderByTarget) getPath() ([]PathElement, error) {
	rawPath, ok := j["path"]
	if !ok {
		return nil, errors.New("OrderByTarget.path is required")
	}
	p, ok := rawPath.([]PathElement)
	if !ok {
		return nil, fmt.Errorf("invalid OrderByTarget.path type; expected: []PathElement, got: %+v", rawPath)
	}
	return p, nil
}

// Type gets the type enum of the current type.
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

// AsColumn tries to convert the instance to OrderByColumn type.
func (j OrderByTarget) AsColumn() (*OrderByColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != OrderByTargetTypeColumn {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeColumn, t)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("OrderByColumn.name: %w", err)
	}

	if name == "" {
		return nil, errors.New("OrderByColumn.name is required")
	}

	p, err := j.getPath()
	if err != nil {
		return nil, err
	}

	fieldPath, err := j.getFieldPath()
	if err != nil {
		return nil, err
	}

	return &OrderByColumn{
		Type:      t,
		Name:      name,
		Path:      p,
		FieldPath: fieldPath,
	}, nil
}

// AsSingleColumnAggregate tries to convert the instance to OrderBySingleColumnAggregate type.
func (j OrderByTarget) AsSingleColumnAggregate() (*OrderBySingleColumnAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != OrderByTargetTypeSingleColumnAggregate {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeSingleColumnAggregate, t)
	}

	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("OrderBySingleColumnAggregate.column: %w", err)
	}

	if column == "" {
		return nil, errors.New("OrderBySingleColumnAggregate.column is required")
	}

	function, err := getStringValueByKey(j, "function")
	if function == "" {
		return nil, fmt.Errorf("OrderBySingleColumnAggregate.function: %w", err)
	}

	if function == "" {
		return nil, errors.New("OrderBySingleColumnAggregate.function is required")
	}

	p, err := j.getPath()
	if err != nil {
		return nil, err
	}

	fieldPath, err := j.getFieldPath()
	if err != nil {
		return nil, err
	}

	return &OrderBySingleColumnAggregate{
		Type:      t,
		Column:    column,
		Function:  function,
		Path:      p,
		FieldPath: fieldPath,
	}, nil
}

// AsStarCountAggregate tries to convert the instance to OrderByStarCountAggregate type.
func (j OrderByTarget) AsStarCountAggregate() (*OrderByStarCountAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != OrderByTargetTypeStarCountAggregate {
		return nil, fmt.Errorf("invalid OrderByTarget type; expected: %s, got: %s", OrderByTargetTypeStarCountAggregate, t)
	}

	p, err := j.getPath()
	if err != nil {
		return nil, err
	}
	return &OrderByStarCountAggregate{
		Type: t,
		Path: p,
	}, nil
}

// Interface tries to convert the instance to OrderByTargetEncoder interface.
func (j OrderByTarget) Interface() OrderByTargetEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to OrderByTargetEncoder interface safely with explicit error.
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

// OrderByTargetEncoder abstracts the serialization interface for OrderByTarget.
type OrderByTargetEncoder interface {
	Encode() OrderByTarget
}

// OrderByColumn represents an ordering object which compares the value in the selected column.
type OrderByColumn struct {
	Type OrderByTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// The name of the column
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Any relationships to traverse to reach this column
	Path []PathElement `json:"path" yaml:"path" mapstructure:"path"`
	// Any field path to a nested field within the column
	FieldPath []string `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// NewOrderByColumn creates an OrderByColumn instance.
func NewOrderByColumn(name string, path []PathElement, fieldPath []string) *OrderByColumn {
	return &OrderByColumn{
		Type:      OrderByTargetTypeColumn,
		Name:      name,
		FieldPath: fieldPath,
		Path:      path,
	}
}

// NewOrderByColumnName creates an OrderByColumn instance with column name only.
func NewOrderByColumnName(name string) *OrderByColumn {
	return NewOrderByColumn(name, []PathElement{}, nil)
}

// Encode converts the instance to raw OrderByTarget.
func (ob OrderByColumn) Encode() OrderByTarget {
	result := OrderByTarget{
		"type": ob.Type,
		"name": ob.Name,
		"path": ob.Path,
	}
	if ob.FieldPath != nil {
		result["field_path"] = ob.FieldPath
	}
	return result
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
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" yaml:"field_path,omitempty" mapstructure:"field_path"`
}

// NewOrderBySingleColumnAggregate creates an OrderBySingleColumnAggregate instance.
func NewOrderBySingleColumnAggregate(column string, function string, path []PathElement, fieldPath []string) *OrderBySingleColumnAggregate {
	return &OrderBySingleColumnAggregate{
		Type:      OrderByTargetTypeSingleColumnAggregate,
		Column:    column,
		Function:  function,
		Path:      path,
		FieldPath: fieldPath,
	}
}

// Encode converts the instance to raw OrderByTarget.
func (ob OrderBySingleColumnAggregate) Encode() OrderByTarget {
	result := OrderByTarget{
		"type":     ob.Type,
		"column":   ob.Column,
		"function": ob.Function,
		"path":     ob.Path,
	}
	if ob.FieldPath != nil {
		result["field_path"] = ob.FieldPath
	}
	return result
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

// NewOrderByStarCountAggregate creates an OrderByStarCountAggregate instance.
func NewOrderByStarCountAggregate(path []PathElement) *OrderByStarCountAggregate {
	return &OrderByStarCountAggregate{
		Type: OrderByTargetTypeStarCountAggregate,
		Path: path,
	}
}

// Encode converts the instance to raw OrderByTarget.
func (ob OrderByStarCountAggregate) Encode() OrderByTarget {
	return OrderByTarget{
		"type": ob.Type,
		"path": ob.Path,
	}
}

// NestedFieldType represents a nested field type enum.
type NestedFieldType string

const (
	NestedFieldTypeObject NestedFieldType = "object"
	NestedFieldTypeArray  NestedFieldType = "array"
)

var enumValues_NestedFieldType = []NestedFieldType{
	NestedFieldTypeObject,
	NestedFieldTypeArray,
}

// ParseNestedFieldType parses the type of nested field.
func ParseNestedFieldType(input string) (NestedFieldType, error) {
	result := NestedFieldType(input)
	if !result.IsValid() {
		return NestedFieldType(""), fmt.Errorf("failed to parse NestedFieldType, expect one of %v, got %s", enumValues_NestedFieldType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
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

// NestedField represents a nested field.
type NestedField map[string]any

// IsNil checks if the field is null or empty.
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
		return fmt.Errorf("field type in NestedField: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case NestedFieldTypeObject:
		rawFields, ok := raw["fields"]
		if !ok || isNullJSON(rawFields) {
			return errors.New("field fields in NestedField is required for object type")
		}
		var fields map[string]Field
		if err := json.Unmarshal(rawFields, &fields); err != nil {
			return fmt.Errorf("field fields in NestedField object: %w", err)
		}
		result["fields"] = fields
	case NestedFieldTypeArray:
		rawFields, ok := raw["fields"]
		if !ok || isNullJSON(rawFields) {
			return errors.New("field fields in NestedField is required for array type")
		}
		var fields NestedField
		if err := json.Unmarshal(rawFields, &fields); err != nil {
			return fmt.Errorf("field fields in NestedField array: %w", err)
		}
		result["fields"] = fields
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type.
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

// AsObject tries to convert the instance to NestedObject type.
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

// AsArray tries to convert the instance to NestedArray type.
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

// Interface tries to convert the instance to NestedFieldEncoder interface.
func (j NestedField) Interface() NestedFieldEncoder {
	result, _ := j.InterfaceT()
	return result
}

// Interface tries to convert the instance to NestedFieldEncoder interface safely with explicit error.
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

// NestedFieldEncoder abstracts the serialization interface for NestedField.
type NestedFieldEncoder interface {
	Encode() NestedField
}

// NestedObject presents a nested object field.
type NestedObject struct {
	Type   NestedFieldType  `json:"type" yaml:"type" mapstructure:"type"`
	Fields map[string]Field `json:"fields" yaml:"fields" mapstructure:"fields"`
}

// NewNestedObject create a new NestedObject instance.
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

// Encode converts the instance to raw NestedField.
func (ob NestedObject) Encode() NestedField {
	return NestedField{
		"type":   ob.Type,
		"fields": ob.Fields,
	}
}

// NestedArray presents a nested array field.
type NestedArray struct {
	Type   NestedFieldType `json:"type" yaml:"type" mapstructure:"type"`
	Fields NestedField     `json:"fields" yaml:"fields" mapstructure:"fields"`
}

// NewNestedArray create a new NestedArray instance.
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

// Encode converts the instance to raw NestedField.
func (ob NestedArray) Encode() NestedField {
	return NestedField{
		"type":   ob.Type,
		"fields": ob.Fields,
	}
}

// NewObjectType creates a new object type
func NewObjectType(fields ObjectTypeFields, foreignKeys ObjectTypeForeignKeys, description *string) ObjectType {
	if fields == nil {
		fields = ObjectTypeFields{}
	}
	if foreignKeys == nil {
		foreignKeys = ObjectTypeForeignKeys{}
	}

	return ObjectType{
		Fields:      fields,
		ForeignKeys: foreignKeys,
		Description: description,
	}
}
