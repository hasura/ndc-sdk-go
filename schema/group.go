package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// GroupOrderByTargetType represents a type of GroupOrderByTarget
type GroupOrderByTargetType string

const (
	GroupOrderByTargetTypeDimension GroupOrderByTargetType = "dimension"
	GroupOrderByTargetTypeAggregate GroupOrderByTargetType = "aggregate"
)

var enumValues_GroupOrderByTargetType = []GroupOrderByTargetType{
	GroupOrderByTargetTypeDimension,
	GroupOrderByTargetTypeAggregate,
}

// ParseGroupOrderByTargetType parses a field type from string.
func ParseGroupOrderByTargetType(input string) (GroupOrderByTargetType, error) {
	result := GroupOrderByTargetType(input)
	if !result.IsValid() {
		return GroupOrderByTargetType(""), fmt.Errorf("failed to parse GroupOrderByTargetType, expect one of %v, got %s", enumValues_GroupOrderByTargetType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
func (j GroupOrderByTargetType) IsValid() bool {
	return slices.Contains(enumValues_GroupOrderByTargetType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupOrderByTargetType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseGroupOrderByTargetType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// GroupOrderByTargetEncoder abstracts a generic interface of GroupOrderByTarget
type GroupOrderByTargetEncoder interface {
	Encode() GroupOrderByTarget
}

// GroupOrderByTarget groups order by target
type GroupOrderByTarget map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupOrderByTarget) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty GroupOrderByTargetType

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in GroupOrderByTarget: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in GroupOrderByTarget: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case GroupOrderByTargetTypeDimension:
		rawIndex, ok := raw["index"]
		if !ok {
			return errors.New("field index in GroupOrderByTarget: required")
		}

		var index uint
		if err := json.Unmarshal(rawIndex, &index); err != nil {
			return fmt.Errorf("field index in GroupOrderByTarget: %w", err)
		}

		results["index"] = index
	case GroupOrderByTargetTypeAggregate:
		rawAggregate, ok := raw["aggregate"]
		if !ok {
			return errors.New("field aggregate in GroupOrderByTarget: required")
		}

		var aggregate Aggregate
		if err := aggregate.UnmarshalJSON(rawAggregate); err != nil {
			return fmt.Errorf("field aggregate in GroupOrderByTarget: %w", err)
		}

		results["aggregate"] = aggregate
	default:
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
func (j GroupOrderByTarget) Type() (GroupOrderByTargetType, error) {
	t, ok := j["type"]
	if !ok {
		return GroupOrderByTargetType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseGroupOrderByTargetType(raw)
		if err != nil {
			return GroupOrderByTargetType(""), err
		}
		return v, nil
	case GroupOrderByTargetType:
		return raw, nil
	default:
		return GroupOrderByTargetType(""), fmt.Errorf("invalid GroupOrderByTargetType type: %+v", t)
	}
}

// AsDimension tries to convert the current type to GroupOrderByTargetDimension.
func (j GroupOrderByTarget) AsDimension() (*GroupOrderByTargetDimension, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupOrderByTargetTypeDimension {
		return nil, fmt.Errorf("invalid GroupOrderByTarget type; expected %s, got %s", GroupOrderByTargetTypeDimension, t)
	}
	rawIndex, ok := j["index"]
	if !ok {
		return nil, errors.New("GroupOrderByTarget.index is required")
	}

	index, ok := rawIndex.(uint)
	if !ok {
		return nil, fmt.Errorf("invalid GroupOrderByTarget.index, expected uint, got %v", rawIndex)
	}

	result := &GroupOrderByTargetDimension{
		Type:  t,
		Index: index,
	}

	return result, nil
}

// AsAggregate tries to convert the current type to GroupOrderByTargetAggregate.
func (j GroupOrderByTarget) AsAggregate() (*GroupOrderByTargetAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupOrderByTargetTypeDimension {
		return nil, fmt.Errorf("invalid GroupOrderByTarget type; expected %s, got %s", GroupOrderByTargetTypeDimension, t)
	}
	rawAggregate, ok := j["aggregate"]
	if !ok {
		return nil, errors.New("GroupOrderByTargetAggregate.aggregate is required")
	}

	aggregate, ok := rawAggregate.(Aggregate)
	if !ok {
		return nil, fmt.Errorf("invalid GroupOrderByTargetAggregate.index, expected Aggregate, got %v", rawAggregate)
	}

	result := &GroupOrderByTargetAggregate{
		Type:      t,
		Aggregate: aggregate,
	}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j GroupOrderByTarget) Interface() GroupOrderByTargetEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j GroupOrderByTarget) InterfaceT() (GroupOrderByTargetEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case GroupOrderByTargetTypeDimension:
		return j.AsDimension()
	case GroupOrderByTargetTypeAggregate:
		return j.AsAggregate()
	default:
		return nil, fmt.Errorf("invalid GroupOrderByTarget type: %s", ty)
	}
}

// GroupOrderByTargetDimension represents a dimension object of GroupOrderByTarget
type GroupOrderByTargetDimension struct {
	Type GroupOrderByTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// The index of the dimension to order by, selected from the dimensions provided in the `Grouping` request.
	Index uint `json:"index" yaml:"index" mapstructure:"index"`
}

// NewGroupOrderByTargetDimension creates a GroupOrderByTargetDimension instance.
func NewGroupOrderByTargetDimension(index uint) *GroupOrderByTargetDimension {
	return &GroupOrderByTargetDimension{
		Type:  GroupOrderByTargetTypeDimension,
		Index: index,
	}
}

// Encode converts the instance to raw OrderByTarget.
func (ob GroupOrderByTargetDimension) Encode() GroupOrderByTarget {
	result := GroupOrderByTarget{
		"type":  ob.Type,
		"index": ob.Index,
	}
	return result
}

// GroupOrderByTargetAggregate represents an aggregate object of GroupOrderByTarget
type GroupOrderByTargetAggregate struct {
	Type GroupOrderByTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// Aggregation method to apply.
	Aggregate Aggregate `json:"aggregate" yaml:"aggregate" mapstructure:"aggregate"`
}

// NewGroupOrderByTargetAggregate creates a GroupOrderByTargetAggregate instance.
func NewGroupOrderByTargetAggregate(aggregate AggregateEncoder) *GroupOrderByTargetAggregate {
	return &GroupOrderByTargetAggregate{
		Type:      GroupOrderByTargetTypeAggregate,
		Aggregate: aggregate.Encode(),
	}
}

// Encode converts the instance to raw GroupOrderByTarget.
func (ob GroupOrderByTargetAggregate) Encode() GroupOrderByTarget {
	result := GroupOrderByTarget{
		"type":      ob.Type,
		"aggregate": ob.Aggregate,
	}
	return result
}

// GroupComparisonTargetType represents a type of GroupComparisonTarget
type GroupComparisonTargetType string

const (
	GroupComparisonTargetTypeAggregate GroupComparisonTargetType = "aggregate"
)

var enumValues_GroupComparisonTargetType = []GroupComparisonTargetType{
	GroupComparisonTargetTypeAggregate,
}

// ParseGroupComparisonTargetType parses a field type from string.
func ParseGroupComparisonTargetType(input string) (GroupComparisonTargetType, error) {
	result := GroupComparisonTargetType(input)
	if !result.IsValid() {
		return GroupComparisonTargetType(""), fmt.Errorf("failed to parse GroupComparisonTargetType, expect one of %v, got %s", enumValues_GroupComparisonTargetType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
func (j GroupComparisonTargetType) IsValid() bool {
	return slices.Contains(enumValues_GroupComparisonTargetType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonTargetType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseGroupComparisonTargetType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// GroupComparisonTargetEncoder abstracts a generic interface of GroupComparisonTarget
type GroupComparisonTargetEncoder interface {
	Encode() GroupComparisonTarget
}

// GroupComparisonTarget represents an aggregate comparison target.
type GroupComparisonTarget map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonTarget) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty GroupComparisonTargetType

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in GroupComparisonTarget: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in GroupComparisonTarget: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case GroupComparisonTargetTypeAggregate:
		rawAggregate, ok := raw["aggregate"]
		if !ok {
			return errors.New("field aggregate in GroupComparisonTarget: required")
		}

		var aggregate Aggregate
		if err := aggregate.UnmarshalJSON(rawAggregate); err != nil {
			return fmt.Errorf("field aggregate in GroupComparisonTarget: %w", err)
		}

		results["aggregate"] = aggregate
	default:
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
func (j GroupComparisonTarget) Type() (GroupComparisonTargetType, error) {
	t, ok := j["type"]
	if !ok {
		return GroupComparisonTargetType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseGroupComparisonTargetType(raw)
		if err != nil {
			return GroupComparisonTargetType(""), err
		}
		return v, nil
	case GroupComparisonTargetType:
		return raw, nil
	default:
		return GroupComparisonTargetType(""), fmt.Errorf("invalid GroupComparisonTarget type: %+v", t)
	}
}

// AsAggregate tries to convert the current type to GroupComparisonTargetAggregate.
func (j GroupComparisonTarget) AsAggregate() (*GroupComparisonTargetAggregate, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupComparisonTargetTypeAggregate {
		return nil, fmt.Errorf("invalid GroupComparisonTarget type; expected %s, got %s", GroupOrderByTargetTypeDimension, t)
	}
	rawAggregate, ok := j["aggregate"]
	if !ok {
		return nil, errors.New("GroupComparisonTargetAggregate.aggregate is required")
	}

	aggregate, ok := rawAggregate.(Aggregate)
	if !ok {
		return nil, fmt.Errorf("invalid GroupComparisonTargetAggregate.index, expected Aggregate, got %v", rawAggregate)
	}

	result := &GroupComparisonTargetAggregate{
		Type:      t,
		Aggregate: aggregate,
	}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j GroupComparisonTarget) Interface() GroupComparisonTargetEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j GroupComparisonTarget) InterfaceT() (GroupComparisonTargetEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case GroupComparisonTargetTypeAggregate:
		return j.AsAggregate()
	default:
		return nil, fmt.Errorf("invalid GroupComparisonTarget type: %s", ty)
	}
}

// GroupComparisonTargetAggregate represents an aggregate object of GroupComparisonTarget
type GroupComparisonTargetAggregate struct {
	Type GroupComparisonTargetType `json:"type" yaml:"type" mapstructure:"type"`
	// Aggregation method to apply.
	Aggregate Aggregate `json:"aggregate" yaml:"aggregate" mapstructure:"aggregate"`
}

// NewGroupComparisonTargetAggregate creates a GroupComparisonTargetAggregate instance.
func NewGroupComparisonTargetAggregate(aggregate AggregateEncoder) *GroupComparisonTargetAggregate {
	return &GroupComparisonTargetAggregate{
		Type:      GroupComparisonTargetTypeAggregate,
		Aggregate: aggregate.Encode(),
	}
}

// Encode converts the instance to a raw GroupComparisonTarget.
func (ob GroupComparisonTargetAggregate) Encode() GroupComparisonTarget {
	result := GroupComparisonTarget{
		"type":      ob.Type,
		"aggregate": ob.Aggregate,
	}
	return result
}

// GroupComparisonValueType represents a group comparison value type enum.
type GroupComparisonValueType string

const (
	GroupComparisonValueTypeScalar   GroupComparisonValueType = "scalar"
	GroupComparisonValueTypeVariable GroupComparisonValueType = "variable"
)

var enumValues_GroupComparisonValueType = []GroupComparisonValueType{
	GroupComparisonValueTypeScalar,
	GroupComparisonValueTypeVariable,
}

// ParseGroupComparisonValueType parses a group comparison value type from string.
func ParseGroupComparisonValueType(input string) (GroupComparisonValueType, error) {
	result := GroupComparisonValueType(input)
	if !result.IsValid() {
		return GroupComparisonValueType(""), fmt.Errorf("failed to parse GroupComparisonValueType, expect one of %v, got %s", enumValues_GroupComparisonValueType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j GroupComparisonValueType) IsValid() bool {
	return slices.Contains(enumValues_GroupComparisonValueType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonValueType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseGroupComparisonValueType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// GroupComparisonValueEncoder represents a group comparison value encoder interface.
type GroupComparisonValueEncoder interface {
	Encode() GroupComparisonValue
}

// GroupComparisonValue represents a group comparison value
type GroupComparisonValue map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonValue) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in GroupComparisonValue: required")
	}

	var ty GroupComparisonValueType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in GroupComparisonValue: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case GroupComparisonValueTypeVariable:
		rawName, ok := raw["name"]
		if !ok {
			return errors.New("field name in GroupComparisonValue is required for variable type")
		}
		var name string
		if err := json.Unmarshal(rawName, &name); err != nil {
			return fmt.Errorf("field name in GroupComparisonValue: %w", err)
		}
		result["name"] = name
	case GroupComparisonValueTypeScalar:
		rawValue, ok := raw["value"]
		if !ok {
			return errors.New("field value in GroupComparisonValue is required for scalar type")
		}
		var value any
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in GroupComparisonValue: %w", err)
		}
		result["value"] = value
	default:
	}
	*j = result
	return nil
}

// GetType gets the type of comparison value.
func (cv GroupComparisonValue) Type() (GroupComparisonValueType, error) {
	t, ok := cv["type"]
	if !ok {
		return GroupComparisonValueType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseGroupComparisonValueType(raw)
		if err != nil {
			return GroupComparisonValueType(""), err
		}
		return v, nil
	case GroupComparisonValueType:
		return raw, nil
	default:
		return GroupComparisonValueType(""), fmt.Errorf("invalid GroupComparisonValue type: %+v", t)
	}
}

// AsScalar tries to convert the comparison value to scalar.
func (cv GroupComparisonValue) AsScalar() (*GroupComparisonValueScalar, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != GroupComparisonValueTypeScalar {
		return nil, fmt.Errorf("invalid GroupComparisonValue type; expected %s, got %s", GroupComparisonValueTypeScalar, ty)
	}

	value, ok := cv["value"]
	if !ok {
		return nil, errors.New("GroupComparisonValueScalar.value is required")
	}

	return &GroupComparisonValueScalar{
		Type:  ty,
		Value: value,
	}, nil
}

// AsVariable tries to convert the comparison value to variable.
func (cv GroupComparisonValue) AsVariable() (*GroupComparisonValueVariable, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}
	if ty != GroupComparisonValueTypeVariable {
		return nil, fmt.Errorf("invalid GroupComparisonValue type; expected %s, got %s", GroupComparisonValueTypeVariable, ty)
	}

	name := getStringValueByKey(cv, "name")
	if name == "" {
		return nil, errors.New("GroupComparisonValueVariable.name is required")
	}
	return &GroupComparisonValueVariable{
		Type: ty,
		Name: name,
	}, nil
}

// Interface converts the comparison value to its generic interface.
func (cv GroupComparisonValue) Interface() GroupComparisonValueEncoder {
	result, _ := cv.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (cv GroupComparisonValue) InterfaceT() (GroupComparisonValueEncoder, error) {
	ty, err := cv.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case GroupComparisonValueTypeVariable:
		return cv.AsVariable()
	case GroupComparisonValueTypeScalar:
		return cv.AsScalar()
	default:
		return nil, fmt.Errorf("invalid GroupComparisonValue type: %s", ty)
	}
}

// GroupComparisonValueScalar represents a group comparison value with scalar type.
type GroupComparisonValueScalar struct {
	Type  GroupComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Value any                      `json:"value" yaml:"value" mapstructure:"value"`
}

// NewGroupComparisonValueScalar creates a new GroupComparisonValueScalar instance.
func NewGroupComparisonValueScalar(value any) *GroupComparisonValueScalar {
	return &GroupComparisonValueScalar{
		Type:  GroupComparisonValueTypeScalar,
		Value: value,
	}
}

// Encode converts to the raw comparison value.
func (cv GroupComparisonValueScalar) Encode() GroupComparisonValue {
	return map[string]any{
		"type":  cv.Type,
		"value": cv.Value,
	}
}

// GroupComparisonValueVariable represents a group comparison value with variable type.
type GroupComparisonValueVariable struct {
	Type GroupComparisonValueType `json:"type" yaml:"type" mapstructure:"type"`
	Name string                   `json:"name" yaml:"name" mapstructure:"name"`
}

// NewGroupComparisonValueVariable creates a new GroupComparisonValueVariable instance.
func NewGroupComparisonValueVariable(name string) *GroupComparisonValueVariable {
	return &GroupComparisonValueVariable{
		Type: GroupComparisonValueTypeVariable,
		Name: name,
	}
}

// Encode converts to the raw comparison value.
func (cv GroupComparisonValueVariable) Encode() GroupComparisonValue {
	return map[string]any{
		"type": cv.Type,
		"name": cv.Name,
	}
}

// GroupExpressionType represents the group expression filter enums.
type GroupExpressionType string

const (
	GroupExpressionTypeAnd                      GroupExpressionType = "and"
	GroupExpressionTypeOr                       GroupExpressionType = "or"
	GroupExpressionTypeNot                      GroupExpressionType = "not"
	GroupExpressionTypeUnaryComparisonOperator  GroupExpressionType = "unary_comparison_operator"
	GroupExpressionTypeBinaryComparisonOperator GroupExpressionType = "binary_comparison_operator"
)

var enumValues_GroupExpressionType = []GroupExpressionType{
	GroupExpressionTypeAnd,
	GroupExpressionTypeOr,
	GroupExpressionTypeNot,
	GroupExpressionTypeUnaryComparisonOperator,
	GroupExpressionTypeBinaryComparisonOperator,
}

// ParseGroupExpressionType parses an expression type argument type from string.
func ParseGroupExpressionType(input string) (GroupExpressionType, error) {
	result := GroupExpressionType(input)
	if !result.IsValid() {
		return GroupExpressionType(""), fmt.Errorf("failed to parse GroupExpressionType, expect one of %v, got %s", enumValues_GroupExpressionType, input)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j GroupExpressionType) IsValid() bool {
	return slices.Contains(enumValues_GroupExpressionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupExpressionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseGroupExpressionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// GroupExpressionEncoder abstracts the expression encoder interface.
type GroupExpressionEncoder interface {
	Encode() GroupExpression
}

// GroupExpression represents a group expression
type GroupExpression map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupExpression) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in GroupExpression: required")
	}

	var ty GroupExpressionType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in GroupExpression: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case GroupExpressionTypeAnd, GroupExpressionTypeOr:
		rawExpressions, ok := raw["expressions"]
		if !ok {
			return fmt.Errorf("field expressions in GroupExpression is required for '%s' type", ty)
		}
		var expressions []GroupExpression
		if err := json.Unmarshal(rawExpressions, &expressions); err != nil {
			return fmt.Errorf("field expressions in GroupExpression: %w", err)
		}
		result["expressions"] = expressions
	case GroupExpressionTypeNot:
		rawExpression, ok := raw["expression"]
		if !ok {
			return fmt.Errorf("field expressions in GroupExpression is required for '%s' type", ty)
		}
		var expression GroupExpression
		if err := json.Unmarshal(rawExpression, &expression); err != nil {
			return fmt.Errorf("field expression in GroupExpression: %w", err)
		}
		result["expression"] = expression
	case GroupExpressionTypeUnaryComparisonOperator:
		rawOperator, ok := raw["operator"]
		if !ok {
			return fmt.Errorf("field operator in GroupExpression is required for '%s' type", ty)
		}
		var operator UnaryComparisonOperator
		if err := json.Unmarshal(rawOperator, &operator); err != nil {
			return fmt.Errorf("field operator in GroupExpression: %w", err)
		}
		result["operator"] = operator

		rawTarget, ok := raw["target"]
		if !ok {
			return fmt.Errorf("field target in GroupExpression is required for '%s' type", ty)
		}
		var target GroupComparisonTarget
		if err := json.Unmarshal(rawTarget, &target); err != nil {
			return fmt.Errorf("field target in GroupExpression: %w", err)
		}
		result["target"] = target
	case GroupExpressionTypeBinaryComparisonOperator:
		rawOperator, ok := raw["operator"]
		if !ok {
			return fmt.Errorf("field operator in GroupExpression is required for '%s' type", ty)
		}
		var operator string
		if err := json.Unmarshal(rawOperator, &operator); err != nil {
			return fmt.Errorf("field operator in GroupExpression: %w", err)
		}

		if operator == "" {
			return fmt.Errorf("field operator in GroupExpression is required for '%s' type", ty)
		}
		result["operator"] = operator

		rawTarget, ok := raw["target"]
		if !ok {
			return fmt.Errorf("field target in GroupExpression is required for '%s' type", ty)
		}
		var target GroupComparisonTarget
		if err := json.Unmarshal(rawTarget, &target); err != nil {
			return fmt.Errorf("field target in GroupExpression: %w", err)
		}
		result["target"] = target

		rawValue, ok := raw["value"]
		if !ok {
			return fmt.Errorf("field value in GroupExpression is required for '%s' type", ty)
		}
		var value GroupComparisonValue
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in GroupExpression: %w", err)
		}
		result["value"] = value
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type.
func (j GroupExpression) Type() (GroupExpressionType, error) {
	t, ok := j["type"]
	if !ok {
		return GroupExpressionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseGroupExpressionType(raw)
		if err != nil {
			return GroupExpressionType(""), err
		}
		return v, nil
	case GroupExpressionType:
		return raw, nil
	default:
		return GroupExpressionType(""), fmt.Errorf("invalid GroupExpression type: %+v", t)
	}
}

// AsAnd tries to convert the instance to and type.
func (j GroupExpression) AsAnd() (*GroupExpressionAnd, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupExpressionTypeAnd {
		return nil, fmt.Errorf("invalid GroupExpression type; expected: %s, got: %s", GroupExpressionTypeAnd, t)
	}

	rawExpressions, ok := j["expressions"]
	if !ok {
		return nil, errors.New("ExpressionAnd.expressions is required")
	}
	expressions, ok := rawExpressions.([]GroupExpression)
	if !ok {
		return nil, fmt.Errorf("invalid ExpressionAnd.expressions type; expected: []GroupExpression, got: %+v", rawExpressions)
	}

	return &GroupExpressionAnd{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsOr tries to convert the instance to ExpressionOr instance.
func (j GroupExpression) AsOr() (*GroupExpressionOr, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupExpressionTypeOr {
		return nil, fmt.Errorf("invalid GroupExpression type; expected: %s, got: %s", GroupExpressionTypeOr, t)
	}

	rawExpressions, ok := j["expressions"]
	if !ok {
		return nil, errors.New("GroupExpressionOr.expression is required")
	}
	expressions, ok := rawExpressions.([]GroupExpression)
	if !ok {
		return nil, fmt.Errorf("invalid GroupExpressionOr.expression type; expected: []GroupExpression, got: %+v", rawExpressions)
	}

	return &GroupExpressionOr{
		Type:        t,
		Expressions: expressions,
	}, nil
}

// AsNot tries to convert the instance to ExpressionNot instance.
func (j GroupExpression) AsNot() (*GroupExpressionNot, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupExpressionTypeNot {
		return nil, fmt.Errorf("invalid GroupExpression type; expected: %s, got: %s", GroupExpressionTypeNot, t)
	}

	rawExpression, ok := j["expression"]
	if !ok {
		return nil, errors.New("ExpressionNot.expression is required")
	}
	expression, ok := rawExpression.(GroupExpression)
	if !ok {
		return nil, fmt.Errorf("invalid GroupExpressionNot.expression type; expected: GroupExpression, got: %+v", rawExpression)
	}

	return &GroupExpressionNot{
		Type:       t,
		Expression: expression,
	}, nil
}

// AsUnaryComparisonOperator tries to convert the instance to ExpressionUnaryComparisonOperator instance.
func (j GroupExpression) AsUnaryComparisonOperator() (*GroupExpressionUnaryComparisonOperator, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupExpressionTypeUnaryComparisonOperator {
		return nil, fmt.Errorf("invalid GroupExpression type; expected: %s, got: %s", GroupExpressionTypeUnaryComparisonOperator, t)
	}

	rawOperator, ok := j["operator"]
	if !ok {
		return nil, errors.New("GroupExpressionUnaryComparisonOperator.operator is required")
	}
	operator, ok := rawOperator.(UnaryComparisonOperator)
	if !ok {
		operatorStr, ok := rawOperator.(string)
		if !ok {
			return nil, fmt.Errorf("invalid GroupExpressionUnaryComparisonOperator.operator type; expected: UnaryComparisonOperator, got: %v", rawOperator)
		}

		operator = UnaryComparisonOperator(operatorStr)
	}

	rawTarget, ok := j["target"]
	if !ok {
		return nil, errors.New("GroupExpressionUnaryComparisonOperator.target is required")
	}
	target, ok := rawTarget.(GroupComparisonTarget)
	if !ok {
		return nil, fmt.Errorf("invalid GroupExpressionUnaryComparisonOperator.target type; expected: GroupComparisonTarget, got: %v", rawTarget)
	}

	return &GroupExpressionUnaryComparisonOperator{
		Type:     t,
		Operator: operator,
		Target:   target,
	}, nil
}

// AsBinaryComparisonOperator tries to convert the instance to ExpressionBinaryComparisonOperator instance.
func (j GroupExpression) AsBinaryComparisonOperator() (*GroupExpressionBinaryComparisonOperator, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != GroupExpressionTypeBinaryComparisonOperator {
		return nil, fmt.Errorf("invalid GroupExpression type; expected: %s, got: %s", GroupExpressionTypeBinaryComparisonOperator, t)
	}

	rawTarget, ok := j["target"]
	if !ok {
		return nil, errors.New("GroupExpressionBinaryComparisonOperator.target is required")
	}
	target, ok := rawTarget.(GroupComparisonTarget)
	if !ok {
		return nil, fmt.Errorf("invalid GroupExpressionBinaryComparisonOperator.target type; expected: GroupComparisonTarget, got: %+v", rawTarget)
	}

	rawValue, ok := j["value"]
	if !ok {
		return nil, errors.New("GroupExpressionBinaryComparisonOperator.value is required")
	}
	value, ok := rawValue.(GroupComparisonValue)
	if !ok {
		return nil, fmt.Errorf("invalid GroupExpressionBinaryComparisonOperator.value type; expected: GroupComparisonValue, got: %+v", rawValue)
	}

	return &GroupExpressionBinaryComparisonOperator{
		Type:     t,
		Operator: getStringValueByKey(j, "operator"),
		Target:   target,
		Value:    value,
	}, nil
}

// Interface tries to convert the instance to the GroupExpressionEncoder interface.
func (j GroupExpression) Interface() GroupExpressionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to the GroupExpressionEncoder interface safely with explicit error.
func (j GroupExpression) InterfaceT() (GroupExpressionEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	switch t {
	case GroupExpressionTypeAnd:
		return j.AsAnd()
	case GroupExpressionTypeOr:
		return j.AsOr()
	case GroupExpressionTypeNot:
		return j.AsNot()
	case GroupExpressionTypeUnaryComparisonOperator:
		return j.AsUnaryComparisonOperator()
	case GroupExpressionTypeBinaryComparisonOperator:
		return j.AsBinaryComparisonOperator()
	default:
		return nil, fmt.Errorf("invalid GroupExpression type: %s", t)
	}
}

// GroupExpressionAnd is an object which represents the [conjunction of expressions]
//
// [conjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#conjunction-of-expressions
type GroupExpressionAnd struct {
	Type        GroupExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expressions []GroupExpression   `json:"expressions" yaml:"expressions" mapstructure:"expressions"`
}

// NewGroupExpressionAnd creates a GroupExpressionAnd instance.
func NewGroupExpressionAnd(expressions ...GroupExpressionEncoder) *GroupExpressionAnd {
	exprs := make([]GroupExpression, len(expressions))
	for i, expr := range expressions {
		if expr == nil {
			continue
		}
		exprs[i] = expr.Encode()
	}
	return &GroupExpressionAnd{
		Type:        GroupExpressionTypeAnd,
		Expressions: exprs,
	}
}

// Encode converts the instance to a raw GroupExpression.
func (exp GroupExpressionAnd) Encode() GroupExpression {
	return GroupExpression{
		"type":        exp.Type,
		"expressions": exp.Expressions,
	}
}

// GroupExpressionOr is an object which represents the [disjunction of expressions]
//
// [disjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#disjunction-of-expressions
type GroupExpressionOr struct {
	Type        GroupExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expressions []GroupExpression   `json:"expressions" yaml:"expressions" mapstructure:"expressions"`
}

// NewGroupExpressionOr creates a GroupExpressionOr instance.
func NewGroupExpressionOr(expressions ...GroupExpressionEncoder) *GroupExpressionOr {
	exprs := make([]GroupExpression, len(expressions))
	for i, expr := range expressions {
		if expr == nil {
			continue
		}
		exprs[i] = expr.Encode()
	}
	return &GroupExpressionOr{
		Type:        GroupExpressionTypeOr,
		Expressions: exprs,
	}
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionOr) Encode() GroupExpression {
	return GroupExpression{
		"type":        exp.Type,
		"expressions": exp.Expressions,
	}
}

// GroupExpressionNot is an object which represents the [negation of an expression]
//
// [negation of an expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#negation
type GroupExpressionNot struct {
	Type       GroupExpressionType `json:"type" yaml:"type" mapstructure:"type"`
	Expression GroupExpression     `json:"expression" yaml:"expression" mapstructure:"expression"`
}

// NewGroupExpressionNot creates a GroupExpressionNot instance.
func NewGroupExpressionNot(expression GroupExpressionEncoder) *GroupExpressionNot {
	result := &GroupExpressionNot{
		Type: GroupExpressionTypeNot,
	}
	if expression != nil {
		result.Expression = expression.Encode()
	}
	return result
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionNot) Encode() GroupExpression {
	return GroupExpression{
		"type":       exp.Type,
		"expression": exp.Expression,
	}
}

// GroupExpressionUnaryComparisonOperator is an object which represents a [unary operator expression]
//
// [unary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type GroupExpressionUnaryComparisonOperator struct {
	Type     GroupExpressionType     `json:"type" yaml:"type" mapstructure:"type"`
	Operator UnaryComparisonOperator `json:"operator" yaml:"operator" mapstructure:"operator"`
	Target   GroupComparisonTarget   `json:"target" yaml:"target" mapstructure:"target"`
}

// NewGroupExpressionUnaryComparisonOperator creates a GroupExpressionUnaryComparisonOperator instance.
func NewGroupExpressionUnaryComparisonOperator(target GroupComparisonTarget, operator UnaryComparisonOperator) *GroupExpressionUnaryComparisonOperator {
	return &GroupExpressionUnaryComparisonOperator{
		Type:     GroupExpressionTypeUnaryComparisonOperator,
		Target:   target,
		Operator: operator,
	}
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionUnaryComparisonOperator) Encode() GroupExpression {
	return GroupExpression{
		"type":     exp.Type,
		"operator": exp.Operator,
		"target":   exp.Target,
	}
}

// GroupExpressionBinaryComparisonOperator is an object which represents an [binary operator expression]
//
// [binary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type GroupExpressionBinaryComparisonOperator struct {
	Type     GroupExpressionType   `json:"type" yaml:"type" mapstructure:"type"`
	Operator string                `json:"operator" yaml:"operator" mapstructure:"operator"`
	Target   GroupComparisonTarget `json:"target" yaml:"target" mapstructure:"target"`
	Value    GroupComparisonValue  `json:"value" yaml:"value" mapstructure:"value"`
}

// NewGroupExpressionBinaryComparisonOperator creates a GroupExpressionBinaryComparisonOperator instance.
func NewGroupExpressionBinaryComparisonOperator(target GroupComparisonTarget, operator string, value GroupComparisonValueEncoder) *GroupExpressionBinaryComparisonOperator {
	result := &GroupExpressionBinaryComparisonOperator{
		Type:     GroupExpressionTypeBinaryComparisonOperator,
		Target:   target,
		Operator: operator,
	}
	if value != nil {
		result.Value = value.Encode()
	}
	return result
}

// Encode converts the instance to a raw GroupExpression.
func (exp GroupExpressionBinaryComparisonOperator) Encode() GroupExpression {
	return GroupExpression{
		"type":     exp.Type,
		"operator": exp.Operator,
		"target":   exp.Target,
		"value":    exp.Value,
	}
}
