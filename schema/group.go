package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// GroupOrderByTargetType represents a type of GroupOrderByTarget.
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
		return GroupOrderByTargetType(
				"",
			), fmt.Errorf(
				"failed to parse GroupOrderByTargetType, expect one of %v, got %s",
				enumValues_GroupOrderByTargetType,
				input,
			)
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

// GroupOrderByTarget groups order by target.
type GroupOrderByTarget struct {
	inner GroupOrderByTargetInner
}

// GroupOrderByTargetInner abstracts a generic interface of GroupOrderByTarget.
type GroupOrderByTargetInner interface {
	Type() GroupOrderByTargetType
	ToMap() map[string]any
	Wrap() GroupOrderByTarget
}

// NewGroupOrderByTarget creates a new GroupOrderByTarget instance.
func NewGroupOrderByTarget[T GroupOrderByTargetInner](inner T) GroupOrderByTarget {
	return GroupOrderByTarget{
		inner: inner,
	}
}

// IsEmpty checks if the inner type is empty.
func (j GroupOrderByTarget) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j GroupOrderByTarget) Type() GroupOrderByTargetType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// Interface returns the inner interface.
func (j GroupOrderByTarget) Interface() GroupOrderByTargetInner {
	return j.inner
}

// MarshalJSON implements json.Marshaler interface.
func (j GroupOrderByTarget) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupOrderByTarget) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseGroupOrderByTargetType(rawType.Type)
	if err != nil {
		return fmt.Errorf("field type in GroupOrderByTarget: %w", err)
	}

	switch ty {
	case GroupOrderByTargetTypeDimension:
		var result GroupOrderByTargetDimension

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupOrderByTargetDimension: %w", err)
		}

		j.inner = &result
	case GroupOrderByTargetTypeAggregate:
		var result GroupOrderByTargetAggregate

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupOrderByTargetAggregate: %w", err)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported GroupOrderByTarget type: %s", ty)
	}

	return nil
}

// GroupOrderByTargetDimension represents a dimension object of GroupOrderByTarget.
type GroupOrderByTargetDimension struct {
	// The index of the dimension to order by, selected from the dimensions provided in the `Grouping` request.
	Index uint `json:"index" mapstructure:"index" yaml:"index"`
}

// NewGroupOrderByTargetDimension creates a GroupOrderByTargetDimension instance.
func NewGroupOrderByTargetDimension(index uint) *GroupOrderByTargetDimension {
	return &GroupOrderByTargetDimension{
		Index: index,
	}
}

// Type return the type name of the instance.
func (ob GroupOrderByTargetDimension) Type() GroupOrderByTargetType {
	return GroupOrderByTargetTypeDimension
}

// ToMap converts the instance to raw map.
func (ob GroupOrderByTargetDimension) ToMap() map[string]any {
	return map[string]any{
		"type":  ob.Type(),
		"index": ob.Index,
	}
}

// Encode converts the instance to raw OrderByTarget.
func (ob GroupOrderByTargetDimension) Wrap() GroupOrderByTarget {
	return NewGroupOrderByTarget(&ob)
}

// GroupOrderByTargetAggregate represents an aggregate object of GroupOrderByTarget.
type GroupOrderByTargetAggregate struct {
	// Aggregation method to apply.
	Aggregate Aggregate `json:"aggregate" mapstructure:"aggregate" yaml:"aggregate"`
}

// NewGroupOrderByTargetAggregate creates a GroupOrderByTargetAggregate instance.
func NewGroupOrderByTargetAggregate(aggregate AggregateEncoder) *GroupOrderByTargetAggregate {
	return &GroupOrderByTargetAggregate{
		Aggregate: aggregate.Encode(),
	}
}

// Type return the type name of the instance.
func (ob GroupOrderByTargetAggregate) Type() GroupOrderByTargetType {
	return GroupOrderByTargetTypeAggregate
}

// ToMap converts the instance to raw map.
func (ob GroupOrderByTargetAggregate) ToMap() map[string]any {
	return map[string]any{
		"type":      ob.Type(),
		"aggregate": ob.Aggregate,
	}
}

// Encode converts the instance to raw GroupOrderByTarget.
func (ob GroupOrderByTargetAggregate) Wrap() GroupOrderByTarget {
	return NewGroupOrderByTarget(&ob)
}

// GroupComparisonTargetType represents a type of GroupComparisonTarget.
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
		return GroupComparisonTargetType(
				"",
			), fmt.Errorf(
				"failed to parse GroupComparisonTargetType, expect one of %v, got %s",
				enumValues_GroupComparisonTargetType,
				input,
			)
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

// GroupComparisonTargetInner abstracts a generic interface of GroupComparisonTarget.
type GroupComparisonTargetInner interface {
	Type() GroupComparisonTargetType
	ToMap() map[string]any
	Wrap() GroupComparisonTarget
}

// GroupComparisonTarget represents an aggregate comparison target.
type GroupComparisonTarget struct {
	inner GroupComparisonTargetInner
}

// NewGroupComparisonTarget creates a new GroupComparisonTarget instance.
func NewGroupComparisonTarget[T GroupComparisonTargetInner](inner T) GroupComparisonTarget {
	return GroupComparisonTarget{
		inner: inner,
	}
}

// IsEmpty checks if the inner type is empty.
func (j GroupComparisonTarget) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j GroupComparisonTarget) Type() GroupComparisonTargetType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// MarshalJSON implements json.Marshaler interface.
func (j GroupComparisonTarget) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonTarget) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseGroupComparisonTargetType(rawType.Type)
	if err != nil {
		return fmt.Errorf("field type in GroupComparisonTarget: %w", err)
	}

	switch ty {
	case GroupComparisonTargetTypeAggregate:
		var result GroupComparisonTargetAggregate

		if err := json.Unmarshal(b, &result); err != nil {
			return fmt.Errorf("failed to unmarshal GroupComparisonTargetAggregate: %w", err)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported GroupComparisonTarget type: %s", ty)
	}

	return nil
}

// Interface converts the comparison value to its generic interface.
func (j GroupComparisonTarget) Interface() GroupComparisonTargetInner {
	return j.inner
}

// GroupComparisonTargetAggregate represents an aggregate object of GroupComparisonTarget.
type GroupComparisonTargetAggregate struct {
	// Aggregation method to apply.
	Aggregate Aggregate `json:"aggregate" mapstructure:"aggregate" yaml:"aggregate"`
}

// NewGroupComparisonTargetAggregate creates a GroupComparisonTargetAggregate instance.
func NewGroupComparisonTargetAggregate(aggregate AggregateEncoder) *GroupComparisonTargetAggregate {
	return &GroupComparisonTargetAggregate{
		Aggregate: aggregate.Encode(),
	}
}

// Type return the type name of the instance.
func (ob GroupComparisonTargetAggregate) Type() GroupComparisonTargetType {
	return GroupComparisonTargetTypeAggregate
}

// ToMap converts the instance to raw map.
func (ob GroupComparisonTargetAggregate) ToMap() map[string]any {
	return map[string]any{
		"type":      ob.Type(),
		"aggregate": ob.Aggregate,
	}
}

// Encode converts the instance to a raw GroupComparisonTarget.
func (ob GroupComparisonTargetAggregate) Wrap() GroupComparisonTarget {
	return NewGroupComparisonTarget(&ob)
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
		return GroupComparisonValueType(
				"",
			), fmt.Errorf(
				"failed to parse GroupComparisonValueType, expect one of %v, got %s",
				enumValues_GroupComparisonValueType,
				input,
			)
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

// GroupComparisonValueInner represents a group comparison value Inner interface.
type GroupComparisonValueInner interface {
	Type() GroupComparisonValueType
	ToMap() map[string]any
	Wrap() GroupComparisonValue
}

// GroupComparisonValue represents a group comparison value.
type GroupComparisonValue struct {
	inner GroupComparisonValueInner
}

// NewGroupComparisonValue creates a new GroupComparisonValue instance.
func NewGroupComparisonValue[T GroupComparisonValueInner](inner T) GroupComparisonValue {
	return GroupComparisonValue{
		inner: inner,
	}
}

// IsEmpty checks if the inner type is empty.
func (j GroupComparisonValue) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j GroupComparisonValue) Type() GroupComparisonValueType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// MarshalJSON implements json.Marshaler interface.
func (j GroupComparisonValue) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupComparisonValue) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	err := json.Unmarshal(b, &rawType)
	if err != nil {
		return err
	}

	ty, err := ParseGroupComparisonValueType(rawType.Type)
	if err != nil {
		return fmt.Errorf("field type in GroupComparisonValue: %w", err)
	}

	switch ty {
	case GroupComparisonValueTypeVariable:
		var result GroupComparisonValueVariable

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupComparisonValueVariable: %w", err)
		}

		j.inner = &result
	case GroupComparisonValueTypeScalar:
		var result GroupComparisonValueScalar

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupComparisonValueScalar: %w", err)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported GroupComparisonValue type: %s", ty)
	}

	return nil
}

// Interface converts the comparison value to its generic interface.
func (cv GroupComparisonValue) Interface() GroupComparisonValueInner {
	return cv.inner
}

// GroupComparisonValueScalar represents a group comparison value with scalar type.
type GroupComparisonValueScalar struct {
	Value any `json:"value" mapstructure:"value" yaml:"value"`
}

// NewGroupComparisonValueScalar creates a new GroupComparisonValueScalar instance.
func NewGroupComparisonValueScalar(value any) *GroupComparisonValueScalar {
	return &GroupComparisonValueScalar{
		Value: value,
	}
}

// Type return the type name of the instance.
func (cv GroupComparisonValueScalar) Type() GroupComparisonValueType {
	return GroupComparisonValueTypeScalar
}

// ToMap converts the instance to raw map.
func (cv GroupComparisonValueScalar) ToMap() map[string]any {
	return map[string]any{
		"type":  cv.Type(),
		"value": cv.Value,
	}
}

// Encode converts to the raw comparison value.
func (cv GroupComparisonValueScalar) Wrap() GroupComparisonValue {
	return NewGroupComparisonValue(&cv)
}

// GroupComparisonValueVariable represents a group comparison value with variable type.
type GroupComparisonValueVariable struct {
	Name string `json:"name" mapstructure:"name" yaml:"name"`
}

// NewGroupComparisonValueVariable creates a new GroupComparisonValueVariable instance.
func NewGroupComparisonValueVariable(name string) *GroupComparisonValueVariable {
	return &GroupComparisonValueVariable{
		Name: name,
	}
}

// Type return the type name of the instance.
func (cv GroupComparisonValueVariable) Type() GroupComparisonValueType {
	return GroupComparisonValueTypeVariable
}

// ToMap converts the instance to raw map.
func (cv GroupComparisonValueVariable) ToMap() map[string]any {
	return map[string]any{
		"type": cv.Type(),
		"name": cv.Name,
	}
}

// Encode converts to the raw comparison value.
func (cv GroupComparisonValueVariable) Wrap() GroupComparisonValue {
	return NewGroupComparisonValue(&cv)
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
		return GroupExpressionType(
				"",
			), fmt.Errorf(
				"failed to parse GroupExpressionType, expect one of %v, got %s",
				enumValues_GroupExpressionType,
				input,
			)
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

// GroupExpressionInner abstracts the expression Inner interface.
type GroupExpressionInner interface {
	Type() GroupExpressionType
	ToMap() map[string]any
	Wrap() GroupExpression
}

// GroupExpression represents a group expression.
type GroupExpression struct {
	inner GroupExpressionInner
}

// NewGroupExpression creates a new GroupExpression instance.
func NewGroupExpression[T GroupExpressionInner](inner T) GroupExpression {
	return GroupExpression{
		inner: inner,
	}
}

// IsEmpty checks if the inner type is empty.
func (j GroupExpression) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j GroupExpression) Type() GroupExpressionType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// MarshalJSON implements json.Marshaler interface.
func (j GroupExpression) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *GroupExpression) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseGroupExpressionType(rawType.Type)
	if err != nil {
		return fmt.Errorf("field type in GroupExpression: %w", err)
	}

	switch ty {
	case GroupExpressionTypeAnd:
		var result GroupExpressionAnd

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupExpressionAnd: %w", err)
		}

		j.inner = &result
	case GroupExpressionTypeOr:
		var result GroupExpressionOr

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupExpressionOr: %w", err)
		}

		j.inner = &result
	case GroupExpressionTypeNot:
		var result GroupExpressionNot

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupExpressionNot: %w", err)
		}

		j.inner = &result
	case GroupExpressionTypeUnaryComparisonOperator:
		var result GroupExpressionUnaryComparisonOperator

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal GroupExpressionUnaryComparisonOperator: %w", err)
		}

		j.inner = &result
	case GroupExpressionTypeBinaryComparisonOperator:
		var result GroupExpressionBinaryComparisonOperator

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal GroupExpressionBinaryComparisonOperator: %w",
				err,
			)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported GroupExpression type: %s", ty)
	}

	return nil
}

// Equal checks whether instances are the same.
func (j GroupExpression) Equal(value GroupExpression) bool {
	return reflect.DeepEqual(j, value)
}

// Interface tries to convert the instance to the GroupExpressionInner interface.
func (j GroupExpression) Interface() GroupExpressionInner {
	return j.inner
}

// GroupExpressionAnd is an object which represents the [conjunction of expressions]
//
// [conjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#conjunction-of-expressions
type GroupExpressionAnd struct {
	Expressions []GroupExpression `json:"expressions" mapstructure:"expressions" yaml:"expressions"`
}

// NewGroupExpressionAnd creates a GroupExpressionAnd instance.
func NewGroupExpressionAnd(expressions ...GroupExpressionInner) *GroupExpressionAnd {
	exprs := make([]GroupExpression, len(expressions))

	for i, expr := range expressions {
		if expr == nil {
			continue
		}

		exprs[i] = expr.Wrap()
	}

	return &GroupExpressionAnd{
		Expressions: exprs,
	}
}

// Type return the type name of the instance.
func (exp GroupExpressionAnd) Type() GroupExpressionType {
	return GroupExpressionTypeAnd
}

// ToMap converts the instance to raw map.
func (exp GroupExpressionAnd) ToMap() map[string]any {
	return map[string]any{
		"type":        exp.Type(),
		"expressions": exp.Expressions,
	}
}

// Encode converts the instance to a raw GroupExpression.
func (exp GroupExpressionAnd) Wrap() GroupExpression {
	return NewGroupExpression(&exp)
}

// GroupExpressionOr is an object which represents the [disjunction of expressions]
//
// [disjunction of expressions]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#disjunction-of-expressions
type GroupExpressionOr struct {
	Expressions []GroupExpression `json:"expressions" mapstructure:"expressions" yaml:"expressions"`
}

// NewGroupExpressionOr creates a GroupExpressionOr instance.
func NewGroupExpressionOr(expressions ...GroupExpressionInner) *GroupExpressionOr {
	exprs := make([]GroupExpression, len(expressions))

	for i, expr := range expressions {
		if expr == nil {
			continue
		}

		exprs[i] = expr.Wrap()
	}

	return &GroupExpressionOr{
		Expressions: exprs,
	}
}

// Type return the type name of the instance.
func (exp GroupExpressionOr) Type() GroupExpressionType {
	return GroupExpressionTypeOr
}

// ToMap converts the instance to raw map.
func (exp GroupExpressionOr) ToMap() map[string]any {
	return map[string]any{
		"type":        exp.Type(),
		"expressions": exp.Expressions,
	}
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionOr) Wrap() GroupExpression {
	return NewGroupExpression(&exp)
}

// GroupExpressionNot is an object which represents the [negation of an expression]
//
// [negation of an expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#negation
type GroupExpressionNot struct {
	Expression GroupExpression `json:"expression" mapstructure:"expression" yaml:"expression"`
}

// NewGroupExpressionNot creates a GroupExpressionNot instance.
func NewGroupExpressionNot[E GroupExpressionInner](expression E) *GroupExpressionNot {
	result := &GroupExpressionNot{
		Expression: expression.Wrap(),
	}

	return result
}

// Type return the type name of the instance.
func (exp GroupExpressionNot) Type() GroupExpressionType {
	return GroupExpressionTypeNot
}

// ToMap converts the instance to raw map.
func (exp GroupExpressionNot) ToMap() map[string]any {
	return map[string]any{
		"type":       exp.Type(),
		"expression": exp.Expression,
	}
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionNot) Wrap() GroupExpression {
	return NewGroupExpression(&exp)
}

// GroupExpressionUnaryComparisonOperator is an object which represents a [unary operator expression]
//
// [unary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type GroupExpressionUnaryComparisonOperator struct {
	Operator UnaryComparisonOperator `json:"operator" mapstructure:"operator" yaml:"operator"`
	Target   GroupComparisonTarget   `json:"target"   mapstructure:"target"   yaml:"target"`
}

// NewGroupExpressionUnaryComparisonOperator creates a GroupExpressionUnaryComparisonOperator instance.
func NewGroupExpressionUnaryComparisonOperator[T GroupComparisonTargetInner](
	target T,
	operator UnaryComparisonOperator,
) *GroupExpressionUnaryComparisonOperator {
	return &GroupExpressionUnaryComparisonOperator{
		Target:   target.Wrap(),
		Operator: operator,
	}
}

// Type return the type name of the instance.
func (exp GroupExpressionUnaryComparisonOperator) Type() GroupExpressionType {
	return GroupExpressionTypeUnaryComparisonOperator
}

// ToMap converts the instance to raw map.
func (exp GroupExpressionUnaryComparisonOperator) ToMap() map[string]any {
	return map[string]any{
		"type":     exp.Type(),
		"operator": exp.Operator,
		"target":   exp.Target,
	}
}

// Encode converts the instance to a raw Expression.
func (exp GroupExpressionUnaryComparisonOperator) Wrap() GroupExpression {
	return NewGroupExpression(&exp)
}

// GroupExpressionBinaryComparisonOperator is an object which represents an [binary operator expression]
//
// [binary operator expression]: https://hasura.github.io/ndc-spec/specification/queries/filtering.html?highlight=expression#unary-operators
type GroupExpressionBinaryComparisonOperator struct {
	Operator string                `json:"operator" mapstructure:"operator" yaml:"operator"`
	Target   GroupComparisonTarget `json:"target"   mapstructure:"target"   yaml:"target"`
	Value    GroupComparisonValue  `json:"value"    mapstructure:"value"    yaml:"value"`
}

// NewGroupExpressionBinaryComparisonOperator creates a GroupExpressionBinaryComparisonOperator instance.
func NewGroupExpressionBinaryComparisonOperator[T GroupComparisonTargetInner, V GroupComparisonValueInner](
	target T,
	operator string,
	value V,
) *GroupExpressionBinaryComparisonOperator {
	result := &GroupExpressionBinaryComparisonOperator{
		Target:   target.Wrap(),
		Operator: operator,
		Value:    value.Wrap(),
	}

	return result
}

// Type return the type name of the instance.
func (exp GroupExpressionBinaryComparisonOperator) Type() GroupExpressionType {
	return GroupExpressionTypeBinaryComparisonOperator
}

// ToMap converts the instance to raw map.
func (exp GroupExpressionBinaryComparisonOperator) ToMap() map[string]any {
	return map[string]any{
		"type":     exp.Type(),
		"operator": exp.Operator,
		"target":   exp.Target,
		"value":    exp.Value,
	}
}

// Encode converts the instance to a raw GroupExpression.
func (exp GroupExpressionBinaryComparisonOperator) Wrap() GroupExpression {
	return NewGroupExpression(&exp)
}
