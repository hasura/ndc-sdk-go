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
		return ArgumentType(
				"",
			), fmt.Errorf(
				"failed to parse ArgumentType, expect one of %v, got %s",
				enumValues_ArgumentType,
				input,
			)
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

// Argument is provided by reference to a variable or as a literal value.
type Argument map[string]any

// ArgumentEncoder abstracts the interface for Argument.
type ArgumentEncoder interface {
	Type() ArgumentType
	Encode() Argument
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Argument) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue parses values from the raw object.
func (j *Argument) FromValue(raw map[string]any) error {
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
		value, ok := raw["value"]
		if !ok {
			return errors.New("field value in Argument is required for literal type")
		}

		arg["value"] = value
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
		return "", errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseArgumentType(raw)
		if err != nil {
			return "", err
		}

		return v, nil
	case ArgumentType:
		return raw, nil
	default:
		return "", fmt.Errorf("invalid Field type: %+v", t)
	}
}

// AsLiteral converts the instance to ArgumentLiteral.
func (j Argument) AsLiteral() (*ArgumentLiteral, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ArgumentTypeLiteral {
		return nil, fmt.Errorf(
			"invalid ArgumentLiteral type; expected: %s, got: %s",
			ArgumentTypeLiteral,
			t,
		)
	}

	value := j["value"]

	return &ArgumentLiteral{
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
		return nil, fmt.Errorf(
			"invalid ArgumentVariable type; expected: %s, got: %s",
			ArgumentTypeVariable,
			t,
		)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("ArgumentVariable.name: %w", err)
	}

	if name == "" {
		return nil, errors.New("ArgumentVariable.name is required")
	}

	return &ArgumentVariable{
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
	Value any `json:"value" mapstructure:"value" yaml:"value"`
}

// NewArgumentLiteral creates an argument with a literal value.
func NewArgumentLiteral(value any) *ArgumentLiteral {
	return &ArgumentLiteral{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j ArgumentLiteral) Type() ArgumentType {
	return ArgumentTypeLiteral
}

// Encode converts the instance to raw Field.
func (j ArgumentLiteral) Encode() Argument {
	return Argument{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// ArgumentVariable represents the variable argument.
type ArgumentVariable struct {
	Name string `json:"name" mapstructure:"name" yaml:"name"`
}

// NewArgumentVariable creates an argument with a variable name.
func NewArgumentVariable(name string) *ArgumentVariable {
	return &ArgumentVariable{
		Name: name,
	}
}

// Type return the type name of the instance.
func (j ArgumentVariable) Type() ArgumentType {
	return ArgumentTypeVariable
}

// Encode converts the instance to raw Field.
func (j ArgumentVariable) Encode() Argument {
	return Argument{
		"type": j.Type(),
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
		return RelationshipArgumentType(
				"",
			), fmt.Errorf(
				"failed to parse RelationshipArgumentType, expect one of %v, got %s",
				enumValues_RelationshipArgumentType,
				input,
			)
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
	Type() RelationshipArgumentType
	Encode() RelationshipArgument
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationshipArgument) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue decodes the raw object value to the instance.
func (j *RelationshipArgument) FromValue(raw map[string]any) error {
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
		argument, err := RelationshipArgument(raw).asLiteral()
		if err != nil {
			return err
		}

		result = argument.Encode()
	case RelationshipArgumentTypeColumn:
		argument, err := RelationshipArgument(raw).asColumn()
		if err != nil {
			return err
		}

		result = argument.Encode()
	case RelationshipArgumentTypeVariable:
		argument, err := RelationshipArgument(raw).asVariable()
		if err != nil {
			return err
		}

		result = argument.Encode()
	default:
	}

	*j = result

	return nil
}

// Type gets the type enum of the current type.
func (j RelationshipArgument) Type() (RelationshipArgumentType, error) {
	t, ok := j["type"]
	if !ok {
		return "", errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseRelationshipArgumentType(raw)
		if err != nil {
			return "", err
		}

		return v, nil
	case RelationshipArgumentType:
		return raw, nil
	default:
		return "", fmt.Errorf("invalid RelationshipArgument type: %+v", t)
	}
}

// AsLiteral converts the instance to RelationshipArgumentLiteral.
func (j RelationshipArgument) AsLiteral() (*RelationshipArgumentLiteral, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != RelationshipArgumentTypeLiteral {
		return nil, fmt.Errorf(
			"invalid RelationshipArgumentLiteral type; expected: %s, got: %s",
			RelationshipArgumentTypeLiteral,
			t,
		)
	}

	return j.asLiteral()
}

func (j RelationshipArgument) asLiteral() (*RelationshipArgumentLiteral, error) {
	value, ok := j["value"]
	if !ok {
		return nil, errors.New("field value in RelationshipArgumentLiteral is required")
	}

	return &RelationshipArgumentLiteral{
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
		return nil, fmt.Errorf(
			"invalid RelationshipArgumentVariable type; expected: %s, got: %s",
			RelationshipArgumentTypeVariable,
			t,
		)
	}

	return j.asVariable()
}

func (j RelationshipArgument) asVariable() (*RelationshipArgumentVariable, error) {
	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("field name in RelationshipArgumentVariable: %w", err)
	}

	if name == "" {
		return nil, errors.New("field name in RelationshipArgumentVariable is required")
	}

	return &RelationshipArgumentVariable{
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
		return nil, fmt.Errorf(
			"invalid RelationshipArgumentTypeColumn type; expected: %s, got: %s",
			RelationshipArgumentTypeColumn,
			t,
		)
	}

	return j.asColumn()
}

func (j RelationshipArgument) asColumn() (*RelationshipArgumentColumn, error) {
	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("field name in RelationshipArgumentColumn: %w", err)
	}

	if name == "" {
		return nil, errors.New("field name in RelationshipArgumentColumn is required")
	}

	return &RelationshipArgumentColumn{
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
	Value any `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationshipArgumentLiteral creates a RelationshipArgumentLiteral instance.
func NewRelationshipArgumentLiteral(value any) *RelationshipArgumentLiteral {
	return &RelationshipArgumentLiteral{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationshipArgumentLiteral) Type() RelationshipArgumentType {
	return RelationshipArgumentTypeLiteral
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentLiteral) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// RelationshipArgumentColumn represents the column relationship argument.
type RelationshipArgumentColumn struct {
	Name string `json:"name" mapstructure:"name" yaml:"name"`
}

// NewRelationshipArgumentColumn creates a RelationshipArgumentColumn instance.
func NewRelationshipArgumentColumn(name string) *RelationshipArgumentColumn {
	return &RelationshipArgumentColumn{
		Name: name,
	}
}

// Type return the type name of the instance.
func (j RelationshipArgumentColumn) Type() RelationshipArgumentType {
	return RelationshipArgumentTypeColumn
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentColumn) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type": RelationshipArgumentTypeColumn,
		"name": j.Name,
	}
}

// RelationshipArgumentVariable represents the variable relationship argument.
type RelationshipArgumentVariable struct {
	Name string `json:"name" mapstructure:"name" yaml:"name"`
}

// NewRelationshipArgumentVariable creates a RelationshipArgumentVariable instance.
func NewRelationshipArgumentVariable(name string) *RelationshipArgumentVariable {
	return &RelationshipArgumentVariable{
		Name: name,
	}
}

// Type return the type name of the instance.
func (j RelationshipArgumentVariable) Type() RelationshipArgumentType {
	return RelationshipArgumentTypeVariable
}

// Encode converts the instance to raw Field.
func (j RelationshipArgumentVariable) Encode() RelationshipArgument {
	return RelationshipArgument{
		"type": j.Type(),
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
		return FieldType(
				"",
			), fmt.Errorf(
				"failed to parse FieldType, expect one of %v, got %s",
				enumValues_FieldType,
				input,
			)
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

// FieldEncoder abstracts the serialization interface for Field family.
type FieldEncoder interface {
	Type() FieldType
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
		return "", errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseFieldType(raw)
		if err != nil {
			return "", err
		}

		return v, nil
	case FieldType:
		return raw, nil
	default:
		return "", fmt.Errorf("invalid Field type: %+v", t)
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
		Column: column,
	}

	rawFields, ok := j["fields"]
	if ok && !isNil(rawFields) {
		fields, ok := rawFields.(NestedField)
		if !ok {
			return nil, fmt.Errorf(
				"invalid ColumnField.fields type; expected NestedField, got %+v",
				rawFields,
			)
		}

		result.Fields = fields
	}

	rawArguments, ok := j["arguments"]
	if ok && !isNil(rawArguments) {
		arguments, ok := rawArguments.(map[string]Argument)
		if !ok {
			return nil, fmt.Errorf(
				"invalid ColumnField.arguments type; expected map[string]Argument, got %+v",
				rawArguments,
			)
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
		return nil, fmt.Errorf(
			"invalid RelationshipField.query type; expected Query, got %+v",
			rawQuery,
		)
	}

	arguments, err := getRelationshipArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in RelationshipField: %w", err)
	}

	if arguments == nil {
		return nil, errors.New("field arguments in RelationshipField is required")
	}

	return &RelationshipField{
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
	// Column name
	Column string `json:"column"           mapstructure:"column" yaml:"column"`
	// When the type of the column is a (possibly-nullable) array or object,
	// the caller can request a subset of the complete column data, by specifying fields to fetch here.
	// If omitted, the column data will be fetched in full.
	Fields NestedField `json:"fields,omitempty" mapstructure:"fields" yaml:"fields,omitempty"`

	Arguments map[string]Argument `json:"arguments,omitempty" mapstructure:"fields" yaml:"arguments,omitempty"`
}

// NewColumnField creates a new ColumnField instance.
func NewColumnField(column string) *ColumnField {
	return &ColumnField{
		Column: column,
	}
}

// WithNestedField return a new column field with nested fields set.
func (f ColumnField) WithNestedField(fields NestedFieldEncoder) *ColumnField {
	if fields != nil {
		f.Fields = fields.Encode()
	}

	return &f
}

// WithArguments return a new column field with arguments set.
func (f ColumnField) WithArguments(arguments map[string]ArgumentEncoder) *ColumnField {
	if arguments == nil {
		f.Arguments = nil

		return &f
	}

	args := make(map[string]Argument)

	for key, arg := range arguments {
		if arg == nil {
			continue
		}

		args[key] = arg.Encode()
	}

	f.Arguments = args

	return &f
}

// WithArgument return a new column field with an arguments set.
func (f ColumnField) WithArgument(key string, argument ArgumentEncoder) *ColumnField {
	if argument == nil {
		delete(f.Arguments, key)
	} else {
		if f.Arguments == nil {
			f.Arguments = make(map[string]Argument)
		}

		f.Arguments[key] = argument.Encode()
	}

	return &f
}

// Type return the type name of the instance.
func (f ColumnField) Type() FieldType {
	return FieldTypeColumn
}

// Encode converts the instance to raw Field.
func (f ColumnField) Encode() Field {
	r := Field{
		"type":   f.Type(),
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

// RelationshipField represents a relationship field.
type RelationshipField struct {
	// The relationship query
	Query Query `json:"query"        mapstructure:"query"        yaml:"query"`
	// The name of the relationship to follow for the subquery
	Relationship string `json:"relationship" mapstructure:"relationship" yaml:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments"    mapstructure:"arguments"    yaml:"arguments"`
}

// NewRelationshipField creates a new RelationshipField instance.
func NewRelationshipField(
	query Query,
	relationship string,
	arguments map[string]RelationshipArgument,
) *RelationshipField {
	return &RelationshipField{
		Query:        query,
		Relationship: relationship,
		Arguments:    arguments,
	}
}

// Type return the type name of the instance.
func (f RelationshipField) Type() FieldType {
	return FieldTypeRelationship
}

// Encode converts the instance to raw Field.
func (f RelationshipField) Encode() Field {
	return Field{
		"type":         f.Type(),
		"query":        f.Query,
		"relationship": f.Relationship,
		"arguments":    f.Arguments,
	}
}

// ComparisonTargetType represents comparison target enums.
type ComparisonTargetType string

const (
	ComparisonTargetTypeColumn    ComparisonTargetType = "column"
	ComparisonTargetTypeAggregate ComparisonTargetType = "aggregate"
)

var enumValues_ComparisonTargetType = []ComparisonTargetType{
	ComparisonTargetTypeColumn,
	ComparisonTargetTypeAggregate,
}

// ParseComparisonTargetType parses a comparison target type argument type from string.
func ParseComparisonTargetType(input string) (ComparisonTargetType, error) {
	result := ComparisonTargetType(input)
	if !result.IsValid() {
		return ComparisonTargetType(
				"",
			), fmt.Errorf(
				"failed to parse ComparisonTargetType, expect one of %v, got: %s",
				enumValues_ComparisonTargetType,
				input,
			)
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
type ComparisonTarget map[string]any

// ComparisonTargetEncoder abstracts the serialization interface for ComparisonTarget family.
type ComparisonTargetEncoder interface {
	Type() ComparisonTargetType
	Encode() ComparisonTarget
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ComparisonTarget) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue parses values from the raw object.
func (j *ComparisonTarget) FromValue(raw map[string]any) error {
	rawFieldType, err := getStringValueByKey(raw, "type")
	if err != nil {
		return fmt.Errorf("field type in ComparisonTarget: %w", err)
	}

	fieldType, err := ParseComparisonTargetType(rawFieldType)
	if err != nil {
		return fmt.Errorf("field type in ComparisonTarget: %w", err)
	}

	results := map[string]any{
		"type": fieldType,
	}

	switch fieldType {
	case ComparisonTargetTypeColumn:
		ct, err := ComparisonTarget(raw).asColumn()
		if err != nil {
			return err
		}

		results = ct.Encode()
	case ComparisonTargetTypeAggregate:
		ct, err := ComparisonTarget(raw).asAggregate()
		if err != nil {
			return err
		}

		results = ct.Encode()
	}

	*j = results

	return nil
}

// Type gets the type enum of the current type.
func (j ComparisonTarget) Type() (ComparisonTargetType, error) {
	t, ok := j["type"]
	if !ok {
		return ComparisonTargetType(""), errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseComparisonTargetType(raw)
		if err != nil {
			return ComparisonTargetType(""), err
		}

		return v, nil
	case ComparisonTargetType:
		return raw, nil
	default:
		return ComparisonTargetType(""), fmt.Errorf("invalid ComparisonTarget type: %+v", t)
	}
}

// AsColumn tries to convert the current type to ComparisonTargetColumn.
func (j ComparisonTarget) AsColumn() (*ComparisonTargetColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ComparisonTargetTypeColumn {
		return nil, fmt.Errorf(
			"invalid ComparisonTarget type; expected %s, got %s",
			ComparisonTargetTypeColumn,
			t,
		)
	}

	return j.asColumn()
}

func (j ComparisonTarget) asColumn() (*ComparisonTargetColumn, error) {
	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("field name in ComparisonTargetColumn: %w", err)
	}

	if name == "" {
		return nil, errors.New("field name in ComparisonTargetColumn is required")
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf("field field_path in ComparisonTargetColumn: %w", err)
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in ComparisonTargetColumn: %w", err)
	}

	result := &ComparisonTargetColumn{
		Name:      name,
		FieldPath: fieldPath,
		Arguments: arguments,
	}

	return result, nil
}

// AsAggregate tries to convert the current type to ComparisonTargetAggregate.
func (j ComparisonTarget) AsAggregate() (*ComparisonTargetAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ComparisonTargetTypeAggregate {
		return nil, fmt.Errorf(
			"invalid ComparisonTarget type; expected %s, got %s",
			ComparisonTargetTypeAggregate,
			t,
		)
	}

	return j.asAggregate()
}

func (j ComparisonTarget) asAggregate() (*ComparisonTargetAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ComparisonTargetTypeAggregate {
		return nil, fmt.Errorf(
			"invalid ComparisonTarget type; expected %s, got %s",
			ComparisonTargetTypeAggregate,
			t,
		)
	}

	result := &ComparisonTargetAggregate{}

	rawAggregate, ok := j["aggregate"]
	if !ok || rawAggregate == nil {
		return nil, errors.New("field aggregate in ComparisonTargetAggregate is required")
	}

	aggregate, ok := rawAggregate.(Aggregate)
	if !ok {
		rawAggregateMap, ok := rawAggregate.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field aggregate in ComparisonTargetAggregate: expected object, got %v",
				rawAggregate,
			)
		}

		if err = aggregate.FromValue(rawAggregateMap); err != nil {
			return nil, fmt.Errorf("field aggregate in ComparisonTargetAggregate: %w", err)
		}
	}

	result.Aggregate = aggregate

	pathElem, err := getPathElementByKey(j, "path")
	if err != nil {
		return nil, fmt.Errorf("field path in ComparisonTargetAggregate: %w", err)
	}

	if pathElem == nil {
		return nil, errors.New("field path in ComparisonTargetAggregate is required")
	}

	result.Path = pathElem

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j ComparisonTarget) Interface() ComparisonTargetEncoder {
	result, _ := j.InterfaceT()

	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j ComparisonTarget) InterfaceT() (ComparisonTargetEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case ComparisonTargetTypeColumn:
		return j.AsColumn()
	case ComparisonTargetTypeAggregate:
		return j.AsAggregate()
	default:
		return nil, fmt.Errorf("invalid ComparisonTarget type: %s", ty)
	}
}

// ComparisonTargetColumn represents a comparison targets a column.
type ComparisonTargetColumn struct {
	// The name of the column
	Name string `json:"name"                 mapstructure:"name"       yaml:"name"`
	// Arguments to satisfy the column specified by 'name'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"  yaml:"arguments,omitempty"`
	// Path to a nested field within an object column. Only non-empty if the 'query.nested_fields.filter_by' capability is supported.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path" yaml:"field_path,omitempty"`
}

// NewComparisonTargetColumn creates a ComparisonTarget with column type.
func NewComparisonTargetColumn(name string) *ComparisonTargetColumn {
	return &ComparisonTargetColumn{
		Name: name,
	}
}

// WithFieldPath returns a new instance with field_path set.
func (f ComparisonTargetColumn) WithFieldPath(fieldPath []string) *ComparisonTargetColumn {
	f.FieldPath = fieldPath

	return &f
}

// WithArguments return a new instance with arguments set.
func (f ComparisonTargetColumn) WithArguments(
	arguments map[string]ArgumentEncoder,
) *ComparisonTargetColumn {
	if arguments == nil {
		f.Arguments = nil

		return &f
	}

	args := make(map[string]Argument)

	for key, arg := range arguments {
		if arg == nil {
			continue
		}

		args[key] = arg.Encode()
	}

	f.Arguments = args

	return &f
}

// WithArgument return a new instance with an arguments set.
func (f ComparisonTargetColumn) WithArgument(
	key string,
	argument ArgumentEncoder,
) *ComparisonTargetColumn {
	if argument == nil {
		delete(f.Arguments, key)
	} else {
		if f.Arguments == nil {
			f.Arguments = make(map[string]Argument)
		}

		f.Arguments[key] = argument.Encode()
	}

	return &f
}

// Type return the type name of the instance.
func (f ComparisonTargetColumn) Type() ComparisonTargetType {
	return ComparisonTargetTypeColumn
}

// Encode converts the instance to raw Field.
func (f ComparisonTargetColumn) Encode() ComparisonTarget {
	r := ComparisonTarget{
		"type": f.Type(),
		"name": f.Name,
	}

	if f.Arguments != nil {
		r["arguments"] = f.Arguments
	}

	if f.FieldPath != nil {
		r["field_path"] = f.FieldPath
	}

	return r
}

// ComparisonTargetAggregate represents a comparison that targets the result of aggregation.
// Only used if the 'query.aggregates.filter_by' capability is supported.
type ComparisonTargetAggregate struct {
	// The aggregation method to use
	Aggregate Aggregate `json:"aggregate" mapstructure:"aggregate" yaml:"aggregate"`
	// Non-empty collection of relationships to traverse
	Path []PathElement `json:"path"      mapstructure:"path"      yaml:"path"`
}

// NewComparisonTargetAggregate creates a ComparisonTargetAggregate instance.
func NewComparisonTargetAggregate[T AggregateEncoder](
	aggregate T,
	path []PathElement,
) *ComparisonTargetAggregate {
	return &ComparisonTargetAggregate{
		Aggregate: aggregate.Encode(),
		Path:      path,
	}
}

// Type return the type name of the instance.
func (f ComparisonTargetAggregate) Type() ComparisonTargetType {
	return ComparisonTargetTypeAggregate
}

// Encode converts the instance to raw Field.
func (f ComparisonTargetAggregate) Encode() ComparisonTarget {
	r := ComparisonTarget{
		"type":      f.Type(),
		"aggregate": f.Aggregate,
		"path":      f.Path,
	}

	return r
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
		return ComparisonValueType(
				"",
			), fmt.Errorf(
				"failed to parse ComparisonValueType, expect one of %v, got %s",
				enumValues_ComparisonValueType,
				input,
			)
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
func (cv *ComparisonValue) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return cv.FromValue(raw)
}

// FromValue decodes values from any map.
func (cv *ComparisonValue) FromValue(input map[string]any) error {
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
		variable, err := ComparisonValue(input).asVariable()
		if err != nil {
			return err
		}

		result = variable.Encode()
	case ComparisonValueTypeColumn:
		column, err := ComparisonValue(input).asColumn()
		if err != nil {
			return err
		}

		result = column.Encode()
	case ComparisonValueTypeScalar:
		value, ok := input["value"]
		if !ok {
			return errors.New("field value in ComparisonValue is required for scalar type")
		}

		result["value"] = value
	}

	*cv = result

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
		return nil, fmt.Errorf(
			"invalid ComparisonValue type; expected %s, got %s",
			ComparisonValueTypeScalar,
			ty,
		)
	}

	value, ok := cv["value"]
	if !ok {
		return nil, errors.New("ComparisonValueScalar.value is required")
	}

	return &ComparisonValueScalar{
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
		return nil, fmt.Errorf(
			"invalid ComparisonValue type; expected %s, got %s",
			ComparisonValueTypeColumn,
			ty,
		)
	}

	return cv.asColumn()
}

// AsColumn tries to convert the comparison value to column.
func (cv ComparisonValue) asColumn() (*ComparisonValueColumn, error) {
	name, err := getStringValueByKey(cv, "name")
	if err != nil {
		return nil, fmt.Errorf("ComparisonValueColumn.name: %w", err)
	}

	pathElem, err := getPathElementByKey(cv, "path")
	if err != nil {
		return nil, fmt.Errorf("field 'path' in ComparisonValueColumn: %w", err)
	}

	if pathElem == nil {
		return nil, errors.New("field 'path' in ComparisonValueColumn is required")
	}

	fieldPath, err := getStringSliceByKey(cv, "field_path")
	if err != nil {
		return nil, fmt.Errorf("field field_path in ComparisonValueColumn: %w", err)
	}

	arguments, err := getArgumentMapByKey(cv, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in ComparisonValueColumn: %w", err)
	}

	result := &ComparisonValueColumn{
		Name:      name,
		Path:      pathElem,
		FieldPath: fieldPath,
		Arguments: arguments,
	}

	rawScope, ok := cv["scope"]
	if ok && !isNil(rawScope) {
		scope, ok := rawScope.(*uint)
		if !ok {
			scope = new(uint)

			if err := mapstructure.Decode(rawScope, scope); err != nil {
				return nil, fmt.Errorf(
					"invalid scope in ComparisonValue; expected *uint, got %v",
					rawScope,
				)
			}
		}

		result.Scope = scope
	}

	return result, nil
}

// AsVariable tries to convert the comparison value to column.
func (cv ComparisonValue) AsVariable() (*ComparisonValueVariable, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}

	if ty != ComparisonValueTypeVariable {
		return nil, fmt.Errorf(
			"invalid ComparisonValue type; expected %s, got %s",
			ComparisonValueTypeVariable,
			ty,
		)
	}

	return cv.asVariable()
}

func (cv ComparisonValue) asVariable() (*ComparisonValueVariable, error) {
	name, err := getStringValueByKey(cv, "name")
	if err != nil {
		return nil, fmt.Errorf("field name in ComparisonValueVariable: %w", err)
	}

	if name == "" {
		return nil, errors.New("field name in ComparisonValueVariable is required")
	}

	return &ComparisonValueVariable{
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
	Type() ComparisonValueType
	Encode() ComparisonValue
}

// The value to compare against should be drawn from another column.
type ComparisonValueColumn struct {
	// The name of the column
	Name string `json:"name"                 mapstructure:"name"       yaml:"name"`
	// Any relationships to traverse to reach this column.
	// Only non-empty if the 'relationships.relation_comparisons' is supported.
	Path []PathElement `json:"path"                 mapstructure:"path"       yaml:"path"`
	// Arguments to satisfy the column specified by 'name'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"  yaml:"arguments,omitempty"`
	// Path to a nested field within an object column. Only non-empty if the 'query.nested_fields.filter_by' capability is supported.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path" yaml:"field_path,omitempty"`
	// The scope in which this column exists, identified by an top-down index into the stack of scopes.
	// The stack grows inside each `Expression::Exists`, so scope 0 (the default) refers to the current collection,
	// and each subsequent index refers to the collection outside its predecessor's immediately enclosing `Expression::Exists` expression.
	// Only used if the 'query.exists.named_scopes' capability is supported.
	Scope *uint `json:"scope,omitempty"      mapstructure:"scope"      yaml:"scope,omitempty"`
}

// NewComparisonValueColumn creates a new ComparisonValueColumn instance.
func NewComparisonValueColumn(
	name string,
	path []PathElement,
	arguments map[string]Argument,
	fieldPath []string,
	scope *uint,
) *ComparisonValueColumn {
	return &ComparisonValueColumn{
		Name:      name,
		Path:      path,
		Arguments: arguments,
		FieldPath: fieldPath,
		Scope:     scope,
	}
}

// Type return the type name of the instance.
func (cv ComparisonValueColumn) Type() ComparisonValueType {
	return ComparisonValueTypeColumn
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueColumn) Encode() ComparisonValue {
	result := map[string]any{
		"type": cv.Type(),
		"name": cv.Name,
		"path": cv.Path,
	}

	if cv.Arguments != nil {
		result["arguments"] = cv.Arguments
	}

	if cv.FieldPath != nil {
		result["field_path"] = cv.FieldPath
	}

	if cv.Scope != nil {
		result["scope"] = cv.Scope
	}

	return result
}

// ComparisonValueScalar represents a comparison value with scalar type.
type ComparisonValueScalar struct {
	Value any `json:"value" mapstructure:"value" yaml:"value"`
}

// NewComparisonValueScalar creates a new ComparisonValueScalar instance.
func NewComparisonValueScalar(value any) *ComparisonValueScalar {
	return &ComparisonValueScalar{
		Value: value,
	}
}

// Type return the type name of the instance.
func (cv ComparisonValueScalar) Type() ComparisonValueType {
	return ComparisonValueTypeScalar
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueScalar) Encode() ComparisonValue {
	return map[string]any{
		"type":  cv.Type(),
		"value": cv.Value,
	}
}

// ComparisonValueVariable represents a comparison value with variable type.
type ComparisonValueVariable struct {
	Name string `json:"name" mapstructure:"name" yaml:"name"`
}

// NewComparisonValueVariable creates a new ComparisonValueVariable instance.
func NewComparisonValueVariable(name string) *ComparisonValueVariable {
	return &ComparisonValueVariable{
		Name: name,
	}
}

// Type return the type name of the instance.
func (cv ComparisonValueVariable) Type() ComparisonValueType {
	return ComparisonValueTypeVariable
}

// Encode converts to the raw comparison value.
func (cv ComparisonValueVariable) Encode() ComparisonValue {
	return map[string]any{
		"type": cv.Type(),
		"name": cv.Name,
	}
}

// ExistsInCollectionType represents an exists in collection type enum.
type ExistsInCollectionType string

const (
	ExistsInCollectionTypeRelated                ExistsInCollectionType = "related"
	ExistsInCollectionTypeUnrelated              ExistsInCollectionType = "unrelated"
	ExistsInCollectionTypeNestedCollection       ExistsInCollectionType = "nested_collection"
	ExistsInCollectionTypeNestedScalarCollection ExistsInCollectionType = "nested_scalar_collection"
)

var enumValues_ExistsInCollectionType = []ExistsInCollectionType{
	ExistsInCollectionTypeRelated,
	ExistsInCollectionTypeUnrelated,
	ExistsInCollectionTypeNestedCollection,
	ExistsInCollectionTypeNestedScalarCollection,
}

// ParseExistsInCollectionType parses a comparison value type from string.
func ParseExistsInCollectionType(input string) (ExistsInCollectionType, error) {
	result := ExistsInCollectionType(input)
	if !result.IsValid() {
		return result, fmt.Errorf(
			"failed to parse ExistsInCollectionType, expect one of %v, got %s",
			enumValues_ExistsInCollectionType,
			input,
		)
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

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue decodes values from any map.
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

	switch ty {
	case ExistsInCollectionTypeRelated:
		eic, err := ExistsInCollection(input).asRelated()
		if err != nil {
			return err
		}

		result = eic.Encode()
	case ExistsInCollectionTypeUnrelated:
		eic, err := ExistsInCollection(input).asUnrelated()
		if err != nil {
			return err
		}

		result = eic.Encode()
	case ExistsInCollectionTypeNestedCollection:
		eic, err := ExistsInCollection(input).asNestedCollection()
		if err != nil {
			return err
		}

		result = eic.Encode()
	case ExistsInCollectionTypeNestedScalarCollection:
		eic, err := ExistsInCollection(input).asNestedScalarCollection()
		if err != nil {
			return err
		}

		result = eic.Encode()
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
		return nil, fmt.Errorf(
			"invalid ExistsInCollection type; expected: %s, got: %s",
			ExistsInCollectionTypeRelated,
			t,
		)
	}

	return j.asRelated()
}

func (j ExistsInCollection) asRelated() (*ExistsInCollectionRelated, error) {
	relationship, err := getStringValueByKey(j, "relationship")
	if err != nil {
		return nil, fmt.Errorf("ExistsInCollectionRelated.relationship: %w", err)
	}

	if relationship == "" {
		return nil, errors.New("ExistsInCollectionRelated.relationship is required")
	}

	arguments, err := getRelationshipArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in ExistsInCollectionRelated: %w", err)
	}

	if arguments == nil {
		return nil, errors.New("field arguments in ExistsInCollectionRelated is required")
	}

	result := &ExistsInCollectionRelated{
		Relationship: relationship,
		Arguments:    arguments,
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf(
			"field field_path in ExistsInCollectionNestedScalarCollection: %w",
			err,
		)
	}

	if fieldPath != nil {
		result.FieldPath = fieldPath
	}

	return result, nil
}

// AsRelated tries to convert the instance to unrelated type.
func (j ExistsInCollection) AsUnrelated() (*ExistsInCollectionUnrelated, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExistsInCollectionTypeUnrelated {
		return nil, fmt.Errorf(
			"invalid ExistsInCollection type; expected: %s, got: %s",
			ExistsInCollectionTypeUnrelated,
			t,
		)
	}

	return j.asUnrelated()
}

func (j ExistsInCollection) asUnrelated() (*ExistsInCollectionUnrelated, error) {
	collection, err := getStringValueByKey(j, "collection")
	if err != nil {
		return nil, fmt.Errorf("ExistsInCollectionUnrelated.collection: %w", err)
	}

	if collection == "" {
		return nil, errors.New("ExistsInCollectionUnrelated.collection is required")
	}

	arguments, err := getRelationshipArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in ExistsInCollectionUnrelated: %w", err)
	}

	if arguments == nil {
		return nil, errors.New("field arguments in ExistsInCollectionUnrelated is required")
	}

	return &ExistsInCollectionUnrelated{
		Collection: collection,
		Arguments:  arguments,
	}, nil
}

// AsNestedCollection tries to convert the instance to nested_collection type.
func (j ExistsInCollection) AsNestedCollection() (*ExistsInCollectionNestedCollection, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExistsInCollectionTypeNestedCollection {
		return nil, fmt.Errorf(
			"invalid ExistsInCollection type; expected: %s, got: %s",
			ExistsInCollectionTypeNestedCollection,
			t,
		)
	}

	return j.asNestedCollection()
}

func (j ExistsInCollection) asNestedCollection() (*ExistsInCollectionNestedCollection, error) {
	columnName, err := getStringValueByKey(j, "column_name")
	if err != nil {
		return nil, fmt.Errorf("field column_name ExistsInCollectionNestedCollection: %w", err)
	}

	if columnName == "" {
		return nil, errors.New(
			"field column_name in ExistsInCollectionNestedCollection is required",
		)
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in ExistsInCollectionNestedCollection: %w", err)
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf("field field_path in ExistsInCollectionNestedCollection: %w", err)
	}

	result := &ExistsInCollectionNestedCollection{
		ColumnName: columnName,
		Arguments:  arguments,
		FieldPath:  fieldPath,
	}

	return result, nil
}

// AsNestedScalarCollection tries to convert the instance to nested_scalar_collection type.
func (j ExistsInCollection) AsNestedScalarCollection() (*ExistsInCollectionNestedScalarCollection, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExistsInCollectionTypeNestedScalarCollection {
		return nil, fmt.Errorf(
			"invalid ExistsInCollection type; expected: %s, got: %s",
			ExistsInCollectionTypeNestedScalarCollection,
			t,
		)
	}

	return j.asNestedScalarCollection()
}

func (j ExistsInCollection) asNestedScalarCollection() (*ExistsInCollectionNestedScalarCollection, error) {
	columnName, err := getStringValueByKey(j, "column_name")
	if err != nil {
		return nil, fmt.Errorf(
			"field column_name in ExistsInCollectionNestedScalarCollection: %w",
			err,
		)
	}

	if columnName == "" {
		return nil, errors.New(
			"field column_name in ExistsInCollectionNestedScalarCollection is required",
		)
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf(
			"field arguments in ExistsInCollectionNestedScalarCollection: %w",
			err,
		)
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf(
			"field field_path in ExistsInCollectionNestedScalarCollection: %w",
			err,
		)
	}

	result := &ExistsInCollectionNestedScalarCollection{
		ColumnName: columnName,
		Arguments:  arguments,
		FieldPath:  fieldPath,
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
	case ExistsInCollectionTypeNestedScalarCollection:
		return j.AsNestedScalarCollection()
	default:
		return nil, fmt.Errorf("invalid ExistsInCollection type: %s", t)
	}
}

// ExistsInCollectionEncoder abstracts the ExistsInCollection serialization interface.
type ExistsInCollectionEncoder interface {
	Type() ExistsInCollectionType
	Encode() ExistsInCollection
}

// ExistsInCollectionRelated represents [Related collections] that are related to the original collection by a relationship in the collection_relationships field of the top-level QueryRequest.
//
// [Related collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#related-collections
type ExistsInCollectionRelated struct {
	Relationship string `json:"relationship"         mapstructure:"relationship" yaml:"relationship"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments"            mapstructure:"arguments"    yaml:"arguments"`
	// Path to a nested field within an object column that must be navigated before the relationship is navigated.
	// Only non-empty if the 'relationships.nested.filtering' capability is supported.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path"   yaml:"field_path,omitempty"`
}

// NewExistsInCollectionRelated creates an ExistsInCollectionRelated instance.
func NewExistsInCollectionRelated(
	relationship string,
	arguments map[string]RelationshipArgument,
) *ExistsInCollectionRelated {
	return &ExistsInCollectionRelated{
		Relationship: relationship,
		Arguments:    arguments,
	}
}

// Type return the type name of the instance.
func (ei ExistsInCollectionRelated) Type() ExistsInCollectionType {
	return ExistsInCollectionTypeRelated
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionRelated) Encode() ExistsInCollection {
	result := ExistsInCollection{
		"type":         ei.Type(),
		"relationship": ei.Relationship,
		"arguments":    ei.Arguments,
	}

	if ei.FieldPath != nil {
		result["field_path"] = ei.FieldPath
	}

	return result
}

// ExistsInCollectionUnrelated represents [unrelated collections].
//
// [unrelated collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#unrelated-collections
type ExistsInCollectionUnrelated struct {
	// The name of a collection
	Collection string `json:"collection" mapstructure:"collection" yaml:"collection"`
	// Values to be provided to any collection arguments
	Arguments map[string]RelationshipArgument `json:"arguments"  mapstructure:"arguments"  yaml:"arguments"`
}

// NewExistsInCollectionUnrelated creates an ExistsInCollectionUnrelated instance.
func NewExistsInCollectionUnrelated(
	collection string,
	arguments map[string]RelationshipArgument,
) *ExistsInCollectionUnrelated {
	return &ExistsInCollectionUnrelated{
		Collection: collection,
		Arguments:  arguments,
	}
}

// Type return the type name of the instance.
func (ei ExistsInCollectionUnrelated) Type() ExistsInCollectionType {
	return ExistsInCollectionTypeUnrelated
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionUnrelated) Encode() ExistsInCollection {
	return ExistsInCollection{
		"type":       ei.Type(),
		"collection": ei.Collection,
		"arguments":  ei.Arguments,
	}
}

// ExistsInCollectionNestedCollection represents [nested collections] expression.
//
// [nested collections]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=exists#nested-collections
type ExistsInCollectionNestedCollection struct {
	// The name of column
	ColumnName string `json:"column_name"          mapstructure:"column_name" yaml:"column_name"`
	// Values to be provided to any collection arguments
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"   yaml:"arguments,omitempty"`
	// Path to a nested collection via object columns
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path"  yaml:"field_path,omitempty"`
}

// NewExistsInCollectionNestedCollection creates an ExistsInCollectionNestedCollection instance.
func NewExistsInCollectionNestedCollection(
	columnName string,
	arguments map[string]Argument,
	fieldPath []string,
) *ExistsInCollectionNestedCollection {
	return &ExistsInCollectionNestedCollection{
		ColumnName: columnName,
		Arguments:  arguments,
		FieldPath:  fieldPath,
	}
}

// Type return the type name of the instance.
func (ei ExistsInCollectionNestedCollection) Type() ExistsInCollectionType {
	return ExistsInCollectionTypeNestedCollection
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionNestedCollection) Encode() ExistsInCollection {
	result := ExistsInCollection{
		"type":        ei.Type(),
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

// ExistsInCollectionNestedScalarCollection Specifies a column that contains a nested array of scalars.
// The array will be brought into scope of the nested expression where each element becomes an object with one '__value' column that contains the element value.
//
//	Only used if the 'query.exists.nested_scalar_collections' capability is supported.
type ExistsInCollectionNestedScalarCollection struct {
	// The name of column
	ColumnName string `json:"column_name"          mapstructure:"column_name" yaml:"column_name"`
	// Values to be provided to any collection arguments
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"   yaml:"arguments,omitempty"`
	// Path to a nested collection via object columns
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path"  yaml:"field_path,omitempty"`
}

// NewExistsInCollectionNestedScalarCollection creates an ExistsInCollectionNestedScalarCollection instance.
func NewExistsInCollectionNestedScalarCollection(
	columnName string,
	arguments map[string]Argument,
	fieldPath []string,
) *ExistsInCollectionNestedScalarCollection {
	return &ExistsInCollectionNestedScalarCollection{
		ColumnName: columnName,
		Arguments:  arguments,
		FieldPath:  fieldPath,
	}
}

// Type return the type name of the instance.
func (ei ExistsInCollectionNestedScalarCollection) Type() ExistsInCollectionType {
	return ExistsInCollectionTypeNestedScalarCollection
}

// Encode converts the instance to its raw type.
func (ei ExistsInCollectionNestedScalarCollection) Encode() ExistsInCollection {
	result := ExistsInCollection{
		"type":        ei.Type(),
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

// ExpressionType represents the filtering expression enums.
type ExpressionType string

const (
	ExpressionTypeAnd                      ExpressionType = "and"
	ExpressionTypeOr                       ExpressionType = "or"
	ExpressionTypeNot                      ExpressionType = "not"
	ExpressionTypeUnaryComparisonOperator  ExpressionType = "unary_comparison_operator"
	ExpressionTypeBinaryComparisonOperator ExpressionType = "binary_comparison_operator"
	ExpressionTypeArrayComparison          ExpressionType = "array_comparison"
	ExpressionTypeExists                   ExpressionType = "exists"
)

var enumValues_ExpressionType = []ExpressionType{
	ExpressionTypeAnd,
	ExpressionTypeOr,
	ExpressionTypeNot,
	ExpressionTypeUnaryComparisonOperator,
	ExpressionTypeBinaryComparisonOperator,
	ExpressionTypeArrayComparison,
	ExpressionTypeExists,
}

// ParseExpressionType parses an expression type argument type from string.
func ParseExpressionType(input string) (ExpressionType, error) {
	result := ExpressionType(input)
	if !result.IsValid() {
		return ExpressionType(
				"",
			), fmt.Errorf(
				"failed to parse ExpressionType, expect one of %v, got %s",
				enumValues_ExpressionType,
				input,
			)
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

// Expression represents the query expression object.
type Expression map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Expression) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue decodes values from any map.
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
		expressions, err := Expression(input).asExpressions()
		if err != nil {
			return err
		}

		result["expressions"] = expressions
	case ExpressionTypeNot:
		expr, err := Expression(input).asNot()
		if err != nil {
			return err
		}

		result = expr.Encode()
	case ExpressionTypeUnaryComparisonOperator:
		expr, err := Expression(input).asUnaryComparisonOperator()
		if err != nil {
			return err
		}

		result = expr.Encode()
	case ExpressionTypeBinaryComparisonOperator:
		expr, err := Expression(input).asBinaryComparisonOperator()
		if err != nil {
			return err
		}

		result = expr.Encode()
	case ExpressionTypeExists:
		expr, err := Expression(input).asExists()
		if err != nil {
			return err
		}

		result = expr.Encode()
	case ExpressionTypeArrayComparison:
		expr, err := Expression(input).asArrayComparison()
		if err != nil {
			return err
		}

		result = expr.Encode()
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
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeAnd,
			t,
		)
	}

	expressions, err := j.asExpressions()
	if err != nil {
		return nil, err
	}

	return &ExpressionAnd{
		Expressions: expressions,
	}, nil
}

func (j Expression) asExpressions() ([]Expression, error) {
	rawExpressions, ok := j["expressions"]
	if !ok || rawExpressions == nil {
		return nil, errors.New("field expressions in Expression is required")
	}

	expressions, ok := rawExpressions.([]Expression)
	if ok {
		return expressions, nil
	}

	rawExpressionsArray, ok := rawExpressions.([]any)
	if !ok {
		return nil, fmt.Errorf(
			"field expressions in Expression: expected array, got %v",
			rawExpressions,
		)
	}

	expressions = []Expression{}

	for i, rawItem := range rawExpressionsArray {
		if rawItem == nil {
			continue
		}

		itemMap, ok := rawItem.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field expressions[%d] in Expression: expected array, got %v",
				i,
				rawExpressions,
			)
		}

		if itemMap == nil {
			continue
		}

		expr := Expression{}
		if err := expr.FromValue(itemMap); err != nil {
			return nil, fmt.Errorf("field expressions in Expression: %w", err)
		}

		expressions = append(expressions, expr)
	}

	return expressions, nil
}

// AsOr tries to convert the instance to ExpressionOr instance.
func (j Expression) AsOr() (*ExpressionOr, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExpressionTypeOr {
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeOr,
			t,
		)
	}

	expressions, err := j.asExpressions()
	if err != nil {
		return nil, err
	}

	return &ExpressionOr{
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
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeNot,
			t,
		)
	}

	return j.asNot()
}

func (j Expression) asNot() (*ExpressionNot, error) {
	rawExpression, ok := j["expression"]
	if !ok {
		return nil, errors.New("field expressions in ExpressionNot is required")
	}

	expression, ok := rawExpression.(Expression)
	if !ok {
		exprMap, ok := rawExpression.(map[string]any)
		if !ok {
			return nil, errors.New("field expression in ExpressionNot is required")
		}

		if err := expression.FromValue(exprMap); err != nil {
			return nil, fmt.Errorf("field expression in Expression: %w", err)
		}
	}

	return &ExpressionNot{
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
		return nil, fmt.Errorf(
			"invalid ExpressionUnaryComparisonOperator type; expected: %s, got: %s",
			ExpressionTypeUnaryComparisonOperator,
			t,
		)
	}

	return j.asUnaryComparisonOperator()
}

func (j Expression) asUnaryComparisonOperator() (*ExpressionUnaryComparisonOperator, error) {
	rawOperator, ok := j["operator"]
	if !ok {
		return nil, errors.New("ExpressionUnaryComparisonOperator.operator is required")
	}

	operator, ok := rawOperator.(UnaryComparisonOperator)
	if !ok {
		operatorStr, ok := rawOperator.(string)
		if !ok || operatorStr == "" {
			return nil, fmt.Errorf(
				"invalid ExpressionUnaryComparisonOperator.operator type; expected: UnaryComparisonOperator, got: %v",
				rawOperator,
			)
		}

		operator = UnaryComparisonOperator(operatorStr)

		if !slices.Contains(enumValues_UnaryComparisonOperator, operator) {
			return nil, fmt.Errorf(
				"invalid operator in ExpressionUnaryComparisonOperator, expect one of %v",
				enumValues_UnaryComparisonOperator,
			)
		}
	}

	rawColumn, ok := j["column"]
	if !ok {
		return nil, errors.New(
			"field column in ExpressionUnaryComparisonOperator.column is required",
		)
	}

	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		rawColumnMap, ok := rawColumn.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field column in ExpressionUnaryComparisonOperator: expected object, got %v",
				rawColumn,
			)
		}

		if err := column.FromValue(rawColumnMap); err != nil {
			return nil, fmt.Errorf("field column in ExpressionUnaryComparisonOperator: %w", err)
		}
	}

	return &ExpressionUnaryComparisonOperator{
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
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeBinaryComparisonOperator,
			t,
		)
	}

	return j.asBinaryComparisonOperator()
}

func (j Expression) asBinaryComparisonOperator() (*ExpressionBinaryComparisonOperator, error) {
	rawColumn, ok := j["column"]
	if !ok {
		return nil, errors.New("ExpressionBinaryComparisonOperator.column is required")
	}

	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		rawColumnMap, ok := rawColumn.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field column in ExpressionBinaryComparisonOperator: expected object, got %v",
				rawColumn,
			)
		}

		if err := column.FromValue(rawColumnMap); err != nil {
			return nil, fmt.Errorf("field column in ExpressionBinaryComparisonOperator: %w", err)
		}
	}

	rawValue, ok := j["value"]
	if !ok {
		return nil, errors.New("field value in ExpressionBinaryComparisonOperator is required")
	}

	value, ok := rawValue.(ComparisonValue)
	if !ok {
		rawValueMap, ok := rawValue.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field value in ExpressionBinaryComparisonOperator: expected map, got %v",
				rawValue,
			)
		}

		value = ComparisonValue{}
		if err := value.FromValue(rawValueMap); err != nil {
			return nil, fmt.Errorf("field value in Expression: %w", err)
		}
	}

	operator, err := getStringValueByKey(j, "operator")
	if err != nil {
		return nil, fmt.Errorf("field operator in ExpressionBinaryComparisonOperator: %w", err)
	}

	if operator == "" {
		return nil, errors.New("field operator in ExpressionBinaryComparisonOperator is required")
	}

	return &ExpressionBinaryComparisonOperator{
		Operator: operator,
		Column:   column,
		Value:    value,
	}, nil
}

// AsArrayComparison tries to convert the instance to ExpressionArrayComparison instance.
func (j Expression) AsArrayComparison() (*ExpressionArrayComparison, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExpressionTypeArrayComparison {
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeArrayComparison,
			t,
		)
	}

	return j.asArrayComparison()
}

func (j Expression) asArrayComparison() (*ExpressionArrayComparison, error) {
	rawColumn, ok := j["column"]
	if !ok || rawColumn == nil {
		return nil, errors.New("field column in ExpressionArrayComparison is required")
	}

	result := &ExpressionArrayComparison{}

	column, ok := rawColumn.(ComparisonTarget)
	if !ok {
		rawColumnMap, ok := rawColumn.(map[string]any)
		if !ok || rawColumnMap == nil {
			return nil, fmt.Errorf(
				"field column in ExpressionArrayComparison: expected object, got %v",
				rawColumn,
			)
		}

		if err := column.FromValue(rawColumnMap); err != nil {
			return nil, fmt.Errorf("field column in ExpressionArrayComparison: %w", err)
		}
	}

	result.Column = column

	rawComparison, ok := j["comparison"]
	if !ok || rawComparison == nil {
		return nil, errors.New("field comparison in ExpressionArrayComparison is required")
	}

	comparison, ok := rawComparison.(ArrayComparison)
	if !ok {
		rawComparisonMap, ok := rawComparison.(map[string]any)
		if !ok || rawComparisonMap == nil {
			return nil, fmt.Errorf(
				"field comparison in ExpressionArrayComparison: expected object, got %v",
				rawComparison,
			)
		}

		if err := comparison.FromValue(rawComparisonMap); err != nil {
			return nil, fmt.Errorf("field comparison in Expression: %w", err)
		}
	}

	result.Comparison = comparison

	return result, nil
}

// AsExists tries to convert the instance to ExpressionExists instance.
func (j Expression) AsExists() (*ExpressionExists, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ExpressionTypeExists {
		return nil, fmt.Errorf(
			"invalid Expression type; expected: %s, got: %s",
			ExpressionTypeExists,
			t,
		)
	}

	return j.asExists()
}

func (j Expression) asExists() (*ExpressionExists, error) {
	rawInCollection, ok := j["in_collection"]
	if !ok || rawInCollection == nil {
		return nil, errors.New("field in_collection in Expression is required")
	}

	inCollection, ok := rawInCollection.(ExistsInCollection)
	if !ok {
		rawInCollectionMap, ok := rawInCollection.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field in_collection in Expression: expected object, got %v",
				rawInCollection,
			)
		}

		if err := inCollection.FromValue(rawInCollectionMap); err != nil {
			return nil, fmt.Errorf("field in_collection in Expression: %w", err)
		}
	}

	result := &ExpressionExists{
		InCollection: inCollection,
	}

	rawPredicate, ok := j["predicate"]
	if ok && rawPredicate != nil {
		predicate, ok := rawPredicate.(Expression)
		if !ok {
			rawPredicateMap, ok := rawPredicate.(map[string]any)
			if !ok {
				return nil, fmt.Errorf(
					"field predicate in Expression: expected map, got %v",
					rawPredicate,
				)
			}

			if err := predicate.FromValue(rawPredicateMap); err != nil {
				return nil, fmt.Errorf("field predicate in Expression: %w", err)
			}
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
	case ExpressionTypeArrayComparison:
		return j.AsArrayComparison()
	default:
		return nil, fmt.Errorf("invalid Expression type: %s", t)
	}
}

// ExpressionEncoder abstracts the expression encoder interface.
type ExpressionEncoder interface {
	Type() ExpressionType
	Encode() Expression
}

// ExpressionAnd is an object which represents the [conjunction of expressions]
//
// [conjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#conjunction-of-expressions
type ExpressionAnd struct {
	Expressions []Expression `json:"expressions" mapstructure:"expressions" yaml:"expressions"`
}

// NewExpressionAnd creates an ExpressionAnd instance.
func NewExpressionAnd(expressions ...ExpressionEncoder) *ExpressionAnd {
	exprs := []Expression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Encode())
		}
	}

	return &ExpressionAnd{
		Expressions: exprs,
	}
}

// Type return the type name of the instance.
func (exp ExpressionAnd) Type() ExpressionType {
	return ExpressionTypeAnd
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionAnd) Encode() Expression {
	return Expression{
		"type":        exp.Type(),
		"expressions": exp.Expressions,
	}
}

// ExpressionOr is an object which represents the [disjunction of expressions]
//
// [disjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#disjunction-of-expressions
type ExpressionOr struct {
	Expressions []Expression `json:"expressions" mapstructure:"expressions" yaml:"expressions"`
}

// NewExpressionOr creates an ExpressionOr instance.
func NewExpressionOr(expressions ...ExpressionEncoder) *ExpressionOr {
	exprs := []Expression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Encode())
		}
	}

	return &ExpressionOr{
		Expressions: exprs,
	}
}

// Type return the type name of the instance.
func (exp ExpressionOr) Type() ExpressionType {
	return ExpressionTypeOr
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionOr) Encode() Expression {
	return Expression{
		"type":        exp.Type(),
		"expressions": exp.Expressions,
	}
}

// ExpressionNot is an object which represents the [negation of an expression]
//
// [negation of an expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#negation
type ExpressionNot struct {
	Expression Expression `json:"expression" mapstructure:"expression" yaml:"expression"`
}

// NewExpressionNot creates an ExpressionNot instance.
func NewExpressionNot[E ExpressionEncoder](expression E) *ExpressionNot {
	return &ExpressionNot{
		Expression: expression.Encode(),
	}
}

// Type return the type name of the instance.
func (exp ExpressionNot) Type() ExpressionType {
	return ExpressionTypeNot
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionNot) Encode() Expression {
	return Expression{
		"type":       exp.Type(),
		"expression": exp.Expression,
	}
}

// ExpressionUnaryComparisonOperator is an object which represents a [unary operator expression]
//
// [unary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionUnaryComparisonOperator struct {
	Operator UnaryComparisonOperator `json:"operator" mapstructure:"operator" yaml:"operator"`
	Column   ComparisonTarget        `json:"column"   mapstructure:"column"   yaml:"column"`
}

// NewExpressionUnaryComparisonOperator creates an ExpressionUnaryComparisonOperator instance.
func NewExpressionUnaryComparisonOperator[C ComparisonTargetEncoder](
	column C,
	operator UnaryComparisonOperator,
) *ExpressionUnaryComparisonOperator {
	return &ExpressionUnaryComparisonOperator{
		Column:   column.Encode(),
		Operator: operator,
	}
}

// Type return the type name of the instance.
func (exp ExpressionUnaryComparisonOperator) Type() ExpressionType {
	return ExpressionTypeUnaryComparisonOperator
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionUnaryComparisonOperator) Encode() Expression {
	return Expression{
		"type":     exp.Type(),
		"operator": exp.Operator,
		"column":   exp.Column,
	}
}

// ExpressionBinaryComparisonOperator is an object which represents an [binary operator expression]
//
// [binary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type ExpressionBinaryComparisonOperator struct {
	Operator string           `json:"operator" mapstructure:"operator" yaml:"operator"`
	Column   ComparisonTarget `json:"column"   mapstructure:"column"   yaml:"column"`
	Value    ComparisonValue  `json:"value"    mapstructure:"value"    yaml:"value"`
}

// NewExpressionBinaryComparisonOperator creates an ExpressionBinaryComparisonOperator instance.
func NewExpressionBinaryComparisonOperator[T ComparisonTargetEncoder, V ComparisonValueEncoder](
	column T,
	operator string,
	value V,
) *ExpressionBinaryComparisonOperator {
	result := &ExpressionBinaryComparisonOperator{
		Column:   column.Encode(),
		Operator: operator,
		Value:    value.Encode(),
	}

	return result
}

// Type return the type name of the instance.
func (exp ExpressionBinaryComparisonOperator) Type() ExpressionType {
	return ExpressionTypeBinaryComparisonOperator
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionBinaryComparisonOperator) Encode() Expression {
	return Expression{
		"type":     exp.Type(),
		"operator": exp.Operator,
		"column":   exp.Column,
		"value":    exp.Value,
	}
}

// ExpressionArrayComparison is comparison against a nested array column.
// Only used if the 'query.nested_fields.filter_by.nested_arrays' capability is supported.
type ExpressionArrayComparison struct {
	Column     ComparisonTarget `json:"column"     mapstructure:"column"     yaml:"column"`
	Comparison ArrayComparison  `json:"comparison" mapstructure:"comparison" yaml:"comparison"`
}

// NewExpressionArrayComparison creates an ExpressionArrayComparison instance.
func NewExpressionArrayComparison[C ComparisonTargetEncoder, A ArrayComparisonEncoder](
	column C,
	comparison A,
) *ExpressionArrayComparison {
	return &ExpressionArrayComparison{
		Column:     column.Encode(),
		Comparison: comparison.Encode(),
	}
}

// Type return the type name of the instance.
func (exp ExpressionArrayComparison) Type() ExpressionType {
	return ExpressionTypeArrayComparison
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionArrayComparison) Encode() Expression {
	return Expression{
		"type":       exp.Type(),
		"column":     exp.Column,
		"comparison": exp.Comparison,
	}
}

// ExpressionExists is an object which represents an [EXISTS expression]
//
// [EXISTS expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#exists-expressions
type ExpressionExists struct {
	Predicate    Expression         `json:"predicate"     mapstructure:"predicate"     yaml:"predicate"`
	InCollection ExistsInCollection `json:"in_collection" mapstructure:"in_collection" yaml:"in_collection"`
}

// NewExpressionExists creates an ExpressionExists instance.
func NewExpressionExists[E ExpressionEncoder, I ExistsInCollectionEncoder](
	predicate E,
	inCollection I,
) *ExpressionExists {
	return &ExpressionExists{
		InCollection: inCollection.Encode(),
		Predicate:    predicate.Encode(),
	}
}

// Type return the type name of the instance.
func (exp ExpressionExists) Type() ExpressionType {
	return ExpressionTypeExists
}

// Encode converts the instance to a raw Expression.
func (exp ExpressionExists) Encode() Expression {
	return Expression{
		"type":          exp.Type(),
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
		return AggregateType(
				"",
			), fmt.Errorf(
				"failed to parse AggregateType, expect one of %v, got %s",
				enumValues_AggregateType,
				input,
			)
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
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue parses values from a raw object.
func (j *Aggregate) FromValue(raw map[string]any) error {
	rawType, err := getStringValueByKey(raw, "type")
	if err != nil {
		return fmt.Errorf("field type in Aggregate: %w", err)
	}

	ty, err := ParseAggregateType(rawType)
	if err != nil {
		return fmt.Errorf("field type in Aggregate: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}

	switch ty {
	case AggregateTypeStarCount:
	case AggregateTypeSingleColumn:
		agg, err := Aggregate(raw).asSingleColumn()
		if err != nil {
			return err
		}

		result = agg.Encode()
	case AggregateTypeColumnCount:
		agg, err := Aggregate(raw).asColumnCount()
		if err != nil {
			return err
		}

		result = agg.Encode()
	}

	*j = result

	return nil
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
		return nil, fmt.Errorf(
			"invalid Aggregate type; expected: %s, got: %s",
			AggregateTypeStarCount,
			t,
		)
	}

	return &AggregateStarCount{}, nil
}

// AsSingleColumn tries to convert the instance to AggregateSingleColumn type.
func (j Aggregate) AsSingleColumn() (*AggregateSingleColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != AggregateTypeSingleColumn {
		return nil, fmt.Errorf(
			"invalid Aggregate type; expected: %s, got: %s",
			AggregateTypeSingleColumn,
			t,
		)
	}

	return j.asSingleColumn()
}

func (j Aggregate) asSingleColumn() (*AggregateSingleColumn, error) {
	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("field column in AggregateSingleColumn: %w", err)
	}

	if column == "" {
		return nil, errors.New("field column in AggregateSingleColumn is required")
	}

	function, err := getStringValueByKey(j, "function")
	if err != nil {
		return nil, fmt.Errorf("field function in AggregateSingleColumn: %w", err)
	}

	if function == "" {
		return nil, errors.New("field function in AggregateSingleColumn is required")
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in AggregateSingleColumn: %w", err)
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf("field field_path in AggregateSingleColumn: %w", err)
	}

	result := &AggregateSingleColumn{
		Column:    column,
		Function:  function,
		Arguments: arguments,
		FieldPath: fieldPath,
	}

	return result, nil
}

// AsColumnCount tries to convert the instance to AggregateColumnCount type.
func (j Aggregate) AsColumnCount() (*AggregateColumnCount, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != AggregateTypeColumnCount {
		return nil, fmt.Errorf(
			"invalid Aggregate type; expected: %s, got: %s",
			AggregateTypeColumnCount,
			t,
		)
	}

	return j.asColumnCount()
}

func (j Aggregate) asColumnCount() (*AggregateColumnCount, error) {
	column, err := getStringValueByKey(j, "column")
	if err != nil {
		return nil, fmt.Errorf("field column in AggregateColumnCount: %w", err)
	}

	if column == "" {
		return nil, errors.New("field column in AggregateColumnCount is required")
	}

	rawDistinct, ok := j["distinct"]
	if !ok || rawDistinct == nil {
		return nil, errors.New("field distinct in AggregateColumnCount is required")
	}

	distinct, ok := rawDistinct.(bool)
	if !ok {
		return nil, fmt.Errorf(
			"invalid distinct in AggregateColumnCount; expected bool, got %+v",
			rawDistinct,
		)
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, fmt.Errorf("field arguments in AggregateColumnCount: %w", err)
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, fmt.Errorf("field field_path in AggregateColumnCount: %w", err)
	}

	return &AggregateColumnCount{
		Column:    column,
		Distinct:  distinct,
		FieldPath: fieldPath,
		Arguments: arguments,
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
	Type() AggregateType
	Encode() Aggregate
}

// AggregateStarCount represents an aggregate object which counts all matched rows.
type AggregateStarCount struct{}

// NewAggregateStarCount creates a new AggregateStarCount instance.
func NewAggregateStarCount() *AggregateStarCount {
	return &AggregateStarCount{}
}

// Type return the type name of the instance.
func (ag AggregateStarCount) Type() AggregateType {
	return AggregateTypeStarCount
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateStarCount) Encode() Aggregate {
	return Aggregate{
		"type": ag.Type(),
	}
}

// AggregateSingleColumn represents an aggregate object which applies an aggregation function (as defined by the column's scalar type in the schema response) to a column.
type AggregateSingleColumn struct {
	// The column to apply the aggregation function to
	Column string `json:"column"               mapstructure:"column"     yaml:"column"`
	// Single column aggregate function name.
	Function string `json:"function"             mapstructure:"function"   yaml:"function"`
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path" yaml:"field_path,omitempty"`
	// Arguments to satisfy the column specified by 'column'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"  yaml:"arguments,omitempty"`
}

// NewAggregateSingleColumn creates a new AggregateSingleColumn instance.
func NewAggregateSingleColumn(
	column string,
	function string,
	fieldPath []string,
	arguments map[string]Argument,
) *AggregateSingleColumn {
	return &AggregateSingleColumn{
		Column:    column,
		Function:  function,
		FieldPath: fieldPath,
		Arguments: arguments,
	}
}

// Type return the type name of the instance.
func (ag AggregateSingleColumn) Type() AggregateType {
	return AggregateTypeSingleColumn
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateSingleColumn) Encode() Aggregate {
	result := Aggregate{
		"type":     ag.Type(),
		"column":   ag.Column,
		"function": ag.Function,
	}

	if ag.FieldPath != nil {
		result["field_path"] = ag.FieldPath
	}

	if ag.Arguments != nil {
		result["arguments"] = ag.Arguments
	}

	return result
}

// AggregateColumnCount represents an aggregate object which count the number of rows with non-null values in the specified columns.
// If the distinct flag is set, then the count should only count unique non-null values of those columns.
type AggregateColumnCount struct {
	// The column to apply the aggregation function to
	Column string `json:"column"               mapstructure:"column"     yaml:"column"`
	// Whether or not only distinct items should be counted.
	Distinct bool `json:"distinct"             mapstructure:"distinct"   yaml:"distinct"`
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path" yaml:"field_path,omitempty"`
	// Arguments to satisfy the column specified by 'column'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"  yaml:"arguments,omitempty"`
}

// NewAggregateColumnCount creates a new AggregateColumnCount instance.
func NewAggregateColumnCount(
	column string,
	distinct bool,
	fieldPath []string,
	arguments map[string]Argument,
) *AggregateColumnCount {
	return &AggregateColumnCount{
		Column:    column,
		Distinct:  distinct,
		FieldPath: fieldPath,
		Arguments: arguments,
	}
}

// Type return the type name of the instance.
func (ag AggregateColumnCount) Type() AggregateType {
	return AggregateTypeColumnCount
}

// Encode converts the instance to raw Aggregate.
func (ag AggregateColumnCount) Encode() Aggregate {
	result := Aggregate{
		"type":     ag.Type(),
		"column":   ag.Column,
		"distinct": ag.Distinct,
	}

	if ag.FieldPath != nil {
		result["field_path"] = ag.FieldPath
	}

	if ag.Arguments != nil {
		result["arguments"] = ag.Arguments
	}

	return result
}

// OrderByTargetType represents a ordering target type.
type OrderByTargetType string

const (
	OrderByTargetTypeColumn    OrderByTargetType = "column"
	OrderByTargetTypeAggregate OrderByTargetType = "aggregate"
)

var enumValues_OrderByTargetType = []OrderByTargetType{
	OrderByTargetTypeColumn,
	OrderByTargetTypeAggregate,
}

// ParseOrderByTargetType parses a ordering target type argument type from string.
func ParseOrderByTargetType(input string) (OrderByTargetType, error) {
	result := OrderByTargetType(input)
	if !result.IsValid() {
		return OrderByTargetType(
				"",
			), fmt.Errorf(
				"failed to parse OrderByTargetType, expect one of %v, got %s",
				enumValues_OrderByTargetType,
				input,
			)
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

		rawFieldPath, ok := raw["field_path"]
		if ok {
			var fieldPath []string

			if err := json.Unmarshal(rawFieldPath, &fieldPath); err != nil {
				return fmt.Errorf("field field_path in OrderByTarget: %w", err)
			}

			result["field_path"] = fieldPath
		}

		rawArguments, ok := raw["arguments"]
		if ok {
			var arguments map[string]Argument

			if err := json.Unmarshal(rawArguments, &arguments); err != nil {
				return fmt.Errorf("field arguments in OrderByTarget: %w", err)
			}

			result["arguments"] = arguments
		}
	case OrderByTargetTypeAggregate:
		rawAggregate, ok := raw["aggregate"]
		if !ok {
			return errors.New("field aggregate in OrderByTarget is required for the aggregate type")
		}

		var aggregate Aggregate
		if err := json.Unmarshal(rawAggregate, &aggregate); err != nil {
			return fmt.Errorf("field aggregate in OrderByTarget: %w", err)
		}

		result["aggregate"] = aggregate
	}

	*j = result

	return nil
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
		return nil, fmt.Errorf(
			"invalid OrderByTarget type; expected: %s, got: %s",
			OrderByTargetTypeColumn,
			t,
		)
	}

	name, err := getStringValueByKey(j, "name")
	if err != nil {
		return nil, fmt.Errorf("field name in OrderByColumn: %w", err)
	}

	if name == "" {
		return nil, errors.New("field name in OrderByColumn is required")
	}

	pathElem, err := getPathElementByKey(j, "path")
	if err != nil {
		return nil, fmt.Errorf("field 'path' in OrderByColumn: %w", err)
	}

	if pathElem == nil {
		return nil, errors.New("field 'path' in OrderByColumn is required")
	}

	fieldPath, err := getStringSliceByKey(j, "field_path")
	if err != nil {
		return nil, err
	}

	arguments, err := getArgumentMapByKey(j, "arguments")
	if err != nil {
		return nil, err
	}

	result := &OrderByColumn{
		Name:      name,
		Path:      pathElem,
		FieldPath: fieldPath,
		Arguments: arguments,
	}

	return result, nil
}

// AsAggregate tries to convert the instance to OrderByAggregate type.
func (j OrderByTarget) AsAggregate() (*OrderByAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != OrderByTargetTypeAggregate {
		return nil, fmt.Errorf(
			"invalid OrderByTarget type; expected: %s, got: %s",
			OrderByTargetTypeAggregate,
			t,
		)
	}

	pathElem, err := getPathElementByKey(j, "path")
	if err != nil {
		return nil, fmt.Errorf("field 'path' in OrderByTarget: %w", err)
	}

	if pathElem == nil {
		return nil, errors.New("field 'path' in OrderByTarget is required")
	}

	rawAggregate, ok := j["aggregate"]
	if !ok {
		return nil, errors.New("invalid OrderByTarget aggregate; aggregate is required")
	}

	aggregate, ok := rawAggregate.(Aggregate)
	if !ok {
		rawAggregateMap, ok := rawAggregate.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"invalid OrderByTarget aggregate; expected map, got %v",
				rawAggregate,
			)
		}

		aggregate = Aggregate{}

		if err := aggregate.FromValue(rawAggregateMap); err != nil {
			return nil, fmt.Errorf("invalid OrderByTarget aggregate: %w", err)
		}
	}

	return &OrderByAggregate{
		Path:      pathElem,
		Aggregate: aggregate,
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
	case OrderByTargetTypeAggregate:
		return j.AsAggregate()
	default:
		return nil, fmt.Errorf("invalid OrderByTarget type: %s", t)
	}
}

// OrderByTargetEncoder abstracts the serialization interface for OrderByTarget.
type OrderByTargetEncoder interface {
	Type() OrderByTargetType
	Encode() OrderByTarget
}

// OrderByColumn represents an ordering object which compares the value in the selected column.
type OrderByColumn struct {
	// The name of the column
	Name string `json:"name"                 mapstructure:"name"       yaml:"name"`
	// Any relationships to traverse to reach this column
	Path []PathElement `json:"path"                 mapstructure:"path"       yaml:"path"`
	// Any field path to a nested field within the column
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path" yaml:"field_path,omitempty"`
	// Arguments to satisfy the column specified by 'name'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"  yaml:"arguments,omitempty"`
}

// NewOrderByColumn creates an OrderByColumn instance.
func NewOrderByColumn(name string, path []PathElement) *OrderByColumn {
	if path == nil {
		path = []PathElement{}
	}

	return &OrderByColumn{
		Name: name,
		Path: path,
	}
}

// WithFieldPath returns a new instance with field_path set.
func (f OrderByColumn) WithFieldPath(fieldPath []string) *OrderByColumn {
	f.FieldPath = fieldPath

	return &f
}

// WithArguments return a new instance with arguments set.
func (f OrderByColumn) WithArguments(
	arguments map[string]ArgumentEncoder,
) *OrderByColumn {
	if arguments == nil {
		f.Arguments = nil

		return &f
	}

	args := make(map[string]Argument)

	for key, arg := range arguments {
		if arg == nil {
			continue
		}

		args[key] = arg.Encode()
	}

	f.Arguments = args

	return &f
}

// WithArgument return a new instance with an arguments set.
func (f OrderByColumn) WithArgument(
	key string,
	argument ArgumentEncoder,
) *OrderByColumn {
	if argument == nil {
		delete(f.Arguments, key)
	} else {
		if f.Arguments == nil {
			f.Arguments = make(map[string]Argument)
		}

		f.Arguments[key] = argument.Encode()
	}

	return &f
}

// Type return the type name of the instance.
func (ob OrderByColumn) Type() OrderByTargetType {
	return OrderByTargetTypeColumn
}

// Encode converts the instance to raw OrderByTarget.
func (ob OrderByColumn) Encode() OrderByTarget {
	result := OrderByTarget{
		"type": ob.Type(),
		"name": ob.Name,
		"path": ob.Path,
	}

	if ob.FieldPath != nil {
		result["field_path"] = ob.FieldPath
	}

	if ob.Arguments != nil {
		result["arguments"] = ob.Arguments
	}

	return result
}

// OrderByAggregate The ordering is performed over the result of an aggregation.
// Only used if the 'relationships.order_by_aggregate' capability is supported.
type OrderByAggregate struct {
	// The aggregation method to use.
	Aggregate Aggregate `json:"aggregate" mapstructure:"aggregate" yaml:"aggregate"`
	// Non-empty collection of relationships to traverse. Only non-empty if the 'relationships' capability is supported.
	// 'PathElement.field_path' will only be non-empty if the 'relationships.nested.ordering' capability is supported.
	Path []PathElement `json:"path"      mapstructure:"path"      yaml:"path"`
}

// NewOrderByAggregate creates an OrderByAggregate instance.
func NewOrderByAggregate[A AggregateEncoder](aggregate A, path []PathElement) *OrderByAggregate {
	return &OrderByAggregate{
		Aggregate: aggregate.Encode(),
		Path:      path,
	}
}

// Type return the type name of the instance.
func (ob OrderByAggregate) Type() OrderByTargetType {
	return OrderByTargetTypeAggregate
}

// Encode converts the instance to raw OrderByTarget.
func (ob OrderByAggregate) Encode() OrderByTarget {
	result := OrderByTarget{
		"type":      ob.Type(),
		"aggregate": ob.Aggregate,
		"path":      ob.Path,
	}

	return result
}

// NestedFieldType represents a nested field type enum.
type NestedFieldType string

const (
	NestedFieldTypeObject     NestedFieldType = "object"
	NestedFieldTypeArray      NestedFieldType = "array"
	NestedFieldTypeCollection NestedFieldType = "collection"
)

var enumValues_NestedFieldType = []NestedFieldType{
	NestedFieldTypeObject,
	NestedFieldTypeArray,
	NestedFieldTypeCollection,
}

// ParseNestedFieldType parses the type of nested field.
func ParseNestedFieldType(input string) (NestedFieldType, error) {
	result := NestedFieldType(input)
	if !result.IsValid() {
		return NestedFieldType(
				"",
			), fmt.Errorf(
				"failed to parse NestedFieldType, expect one of %v, got %s",
				enumValues_NestedFieldType,
				input,
			)
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
	case NestedFieldTypeCollection:
		rawQuery, ok := raw["query"]

		if !ok || isNullJSON(rawQuery) {
			return errors.New("field query in NestedField is required for the collection type")
		}

		var query Query

		if err := json.Unmarshal(rawQuery, &query); err != nil {
			return fmt.Errorf("field query in NestedField collection: %w", err)
		}

		result["query"] = query
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
		return nil, fmt.Errorf(
			"invalid NestedField type; expected: %s, got: %s",
			NestedFieldTypeObject,
			t,
		)
	}

	rawFields, ok := j["fields"]
	if !ok {
		return nil, errors.New("NestedObject.fields is required")
	}

	fields, ok := rawFields.(map[string]Field)
	if !ok {
		return nil, fmt.Errorf(
			"invalid NestedObject.fields type; expected: map[string]Field, got: %+v",
			rawFields,
		)
	}

	return &NestedObject{
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
		return nil, fmt.Errorf(
			"invalid NestedField type; expected: %s, got: %s",
			NestedFieldTypeArray,
			t,
		)
	}

	rawFields, ok := j["fields"]
	if !ok {
		return nil, errors.New("NestedArray.fields is required")
	}

	fields, ok := rawFields.(NestedField)
	if !ok {
		return nil, fmt.Errorf(
			"invalid NestedArray.fields type; expected: NestedField, got: %+v",
			rawFields,
		)
	}

	return &NestedArray{
		Fields: fields,
	}, nil
}

// AsCollection tries to convert the instance to NestedCollection type.
func (j NestedField) AsCollection() (*NestedCollection, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != NestedFieldTypeCollection {
		return nil, fmt.Errorf(
			"invalid NestedField type; expected: %s, got: %s",
			NestedFieldTypeCollection,
			t,
		)
	}

	rawQuery, ok := j["query"]
	if !ok {
		return nil, errors.New("field query in NestedCollection is required")
	}

	query, ok := rawQuery.(Query)
	if !ok {
		return nil, fmt.Errorf(
			"invalid field query in NestedCollection; expected: Query, got: %+v",
			rawQuery,
		)
	}

	return &NestedCollection{
		Query: query,
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
	case NestedFieldTypeCollection:
		return j.AsCollection()
	default:
		return nil, fmt.Errorf("invalid NestedField type: %s", t)
	}
}

// NestedFieldEncoder abstracts the serialization interface for NestedField.
type NestedFieldEncoder interface {
	Type() NestedFieldType
	Encode() NestedField
}

// NestedObject presents a nested object field.
type NestedObject struct {
	Fields map[string]Field `json:"fields" mapstructure:"fields" yaml:"fields"`
}

// NewNestedObject create a new NestedObject instance.
func NewNestedObject(fields map[string]FieldEncoder) *NestedObject {
	fieldMap := make(map[string]Field)

	for k, v := range fields {
		if v == nil {
			continue
		}

		fieldMap[k] = v.Encode()
	}

	return &NestedObject{
		Fields: fieldMap,
	}
}

// Type return the type name of the instance.
func (ob NestedObject) Type() NestedFieldType {
	return NestedFieldTypeObject
}

// Encode converts the instance to raw NestedField.
func (ob NestedObject) Encode() NestedField {
	return NestedField{
		"type":   ob.Type(),
		"fields": ob.Fields,
	}
}

// NestedArray presents a nested array field.
type NestedArray struct {
	Fields NestedField `json:"fields" mapstructure:"fields" yaml:"fields"`
}

// NewNestedArray create a new NestedArray instance.
func NewNestedArray[T NestedFieldEncoder](fields T) *NestedArray {
	return &NestedArray{
		Fields: fields.Encode(),
	}
}

// Type return the type name of the instance.
func (ob NestedArray) Type() NestedFieldType {
	return NestedFieldTypeArray
}

// Encode converts the instance to raw NestedField.
func (ob NestedArray) Encode() NestedField {
	return NestedField{
		"type":   ob.Type(),
		"fields": ob.Fields,
	}
}

// NestedCollection presents a nested collection field.
type NestedCollection struct {
	Query Query `json:"query" mapstructure:"query" yaml:"query"`
}

// NewNestedCollection create a new NestedCollection instance.
func NewNestedCollection(query Query) *NestedCollection {
	return &NestedCollection{
		Query: query,
	}
}

// Type return the type name of the instance.
func (ob NestedCollection) Type() NestedFieldType {
	return NestedFieldTypeCollection
}

// Encode converts the instance to raw NestedField.
func (ob NestedCollection) Encode() NestedField {
	return NestedField{
		"type":  ob.Type(),
		"query": ob.Query,
	}
}

// NewObjectType creates a new object type.
func NewObjectType(
	fields ObjectTypeFields,
	foreignKeys ObjectTypeForeignKeys,
	description *string,
) ObjectType {
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
