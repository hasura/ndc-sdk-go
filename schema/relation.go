package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// JoinType represents a join type.
type JoinType string

const (
	// Left join. Only used when the capability `relational_query.join.join_types.left` is supported.
	LeftJoin JoinType = "left"
	// Right join. Only used when the capability `relational_query.join.join_types.right` is supported.
	RightJoin JoinType = "right"
	// Inner join. Only used when the capability `relational_query.join.join_types.inner` is supported.
	InnerJoin JoinType = "inner"
	// Full join. Only used when the capability `relational_query.join.join_types.full` is supported.
	FullJoin JoinType = "full"
	// Left anti join. Only used when the capability `relational_query.join.join_types.left_anti` is supported.
	LeftAntiJoin JoinType = "left_anti"
	// Left semi join. Only used when the capability `relational_query.join.join_types.left_semi` is supported.
	LeftSemiJoin JoinType = "left_semi"
	// Right anti join. Only used when the capability `relational_query.join.join_types.right_anti` is supported.
	RightAntiJoin JoinType = "right_anti"
	// Right semi join. Only used when the capability `relational_query.join.join_types.right_semi` is supported.
	RightSemiJoin JoinType = "right_semi"
)

var enumValues_JoinType = []JoinType{
	LeftJoin,
	RightJoin,
	InnerJoin,
	FullJoin,
	LeftAntiJoin,
	LeftSemiJoin,
	RightAntiJoin,
	RightSemiJoin,
}

// ParseJoinType parses a join type from string.
func ParseJoinType(input string) (JoinType, error) {
	result := JoinType(input)
	if !result.IsValid() {
		return JoinType(
				"",
			), fmt.Errorf(
				"failed to parse JoinType, expect one of %v, got %s",
				enumValues_JoinType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j JoinType) IsValid() bool {
	return slices.Contains(enumValues_JoinType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *JoinType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseJoinType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// RelationType represents a relation type enum.
type RelationType string

const (
	RelationTypeFrom      RelationType = "from"
	RelationTypePaginate  RelationType = "paginate"
	RelationTypeProject   RelationType = "project"
	RelationTypeFilter    RelationType = "filter"
	RelationTypeSort      RelationType = "sort"
	RelationTypeJoin      RelationType = "join"
	RelationTypeAggregate RelationType = "aggregate"
	RelationTypeWindow    RelationType = "window"
	RelationTypeUnion     RelationType = "union"
)

var enumValues_RelationType = []RelationType{
	RelationTypeFrom,
	RelationTypePaginate,
	RelationTypeProject,
	RelationTypeFilter,
	RelationTypeSort,
	RelationTypeJoin,
	RelationTypeAggregate,
	RelationTypeWindow,
	RelationTypeUnion,
}

// ParseRelationType parses a relation type from string.
func ParseRelationType(input string) (RelationType, error) {
	result := RelationType(input)
	if !result.IsValid() {
		return RelationType(
				"",
			), fmt.Errorf(
				"failed to parse RelationType, expect one of %v, got %s",
				enumValues_RelationType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j RelationType) IsValid() bool {
	return slices.Contains(enumValues_RelationType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseRelationType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// Relation is provided by reference to a relation.
type Relation struct {
	inner RelationInner
}

// NewRelation creates the relation instance.
func NewRelation[R RelationInner](inner R) Relation {
	return Relation{
		inner: inner,
	}
}

// RelationInner abstracts the interface for Relation.
type RelationInner interface {
	Type() RelationType
	ToMap() map[string]any
	Wrap() Relation
}

// IsEmpty checks if the inner type is empty.
func (j Relation) IsEmpty() bool {
	return j.inner == nil
}

// MarshalJSON implements json.Marshaler interface.
func (j Relation) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Relation) UnmarshalJSON(b []byte) error {
	var rawTypeObject rawTypeStruct

	if err := json.Unmarshal(b, &rawTypeObject); err != nil {
		return err
	}

	ty, err := ParseRelationType(rawTypeObject.Type)
	if err != nil {
		return fmt.Errorf("field type in Relation: %w", err)
	}

	switch ty {
	case RelationTypeAggregate:
		var result RelationAggregate

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationAggregate: %w", err)
		}

		j.inner = &result
	case RelationTypeFilter:
		var result RelationFilter

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationFilter: %w", err)
		}

		j.inner = &result
	case RelationTypeFrom:
		var result RelationFrom

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationFrom: %w", err)
		}

		j.inner = &result
	case RelationTypeJoin:
		var result RelationJoin

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationJoin: %w", err)
		}

		j.inner = &result
	case RelationTypePaginate:
		var result RelationPaginate

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationPaginate: %w", err)
		}

		j.inner = &result
	case RelationTypeSort:
		var result RelationSort

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationSort: %w", err)
		}

		j.inner = &result
	case RelationTypeUnion:
		var result RelationUnion

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationUnion: %w", err)
		}

		j.inner = &result
	case RelationTypeWindow:
		var result RelationWindow

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationWindow: %w", err)
		}

		j.inner = &result
	case RelationTypeProject:
		var result RelationProject

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to decode RelationProject: %w", err)
		}

		j.inner = &result
	default:
		return errors.New("invalid relation type: " + string(ty))
	}

	return nil
}

// Type gets the type enum of the current type.
func (j Relation) Type() RelationType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// Interface tries to convert the instance to RelationInner interface.
func (j Relation) Interface() RelationInner {
	return j.inner
}

// RelationFrom represents a from relation.
type RelationFrom struct {
	Collection string                       `json:"collection"          mapstructure:"collection" yaml:"collection"`
	Columns    []string                     `json:"columns"             mapstructure:"columns"    yaml:"columns"`
	Arguments  map[string]RelationalLiteral `json:"arguments,omitempty" mapstructure:"arguments"  yaml:"arguments,omitempty"`
}

// NewRelationFrom creates a RelationFrom instance.
func NewRelationFrom(collection string, columns []string) *RelationFrom {
	return &RelationFrom{
		Collection: collection,
		Columns:    columns,
	}
}

// WithArguments return a new column field with arguments set.
func (f RelationFrom) WithArguments(arguments map[string]RelationalLiteralInner) *RelationFrom {
	if arguments == nil {
		f.Arguments = nil

		return &f
	}

	args := make(map[string]RelationalLiteral)

	for key, arg := range arguments {
		if arg == nil {
			continue
		}

		args[key] = arg.Wrap()
	}

	f.Arguments = args

	return &f
}

// WithArgument return a new column field with an arguments set.
func (f RelationFrom) WithArgument(key string, argument RelationalLiteralInner) *RelationFrom {
	if argument == nil {
		delete(f.Arguments, key)
	} else {
		if f.Arguments == nil {
			f.Arguments = make(map[string]RelationalLiteral)
		}

		f.Arguments[key] = argument.Wrap()
	}

	return &f
}

// Type return the type name of the instance.
func (j RelationFrom) Type() RelationType {
	return RelationTypeFrom
}

// Encode returns the relation wrapper.
func (j RelationFrom) Wrap() Relation {
	return NewRelation(&j)
}

// ToMap converts the instance to raw Relation.
func (j RelationFrom) ToMap() map[string]any {
	result := map[string]any{
		"type":       j.Type(),
		"collection": j.Collection,
		"columns":    j.Columns,
	}

	if j.Arguments != nil {
		result["arguments"] = j.Arguments
	}

	return result
}

// RelationPaginate represents a paginate relation.
type RelationPaginate struct {
	Input Relation `json:"input" mapstructure:"input" yaml:"input"`
	Fetch *uint64  `json:"fetch" mapstructure:"fetch" yaml:"fetch"`
	Skip  uint64   `json:"skip"  mapstructure:"skip"  yaml:"skip"`
}

// NewRelationPaginate creates a RelationPaginate instance.
func NewRelationPaginate[R RelationInner](input R, fetch *uint64, skip uint64) *RelationPaginate {
	return &RelationPaginate{
		Input: NewRelation(input),
		Fetch: fetch,
		Skip:  skip,
	}
}

// Type return the type name of the instance.
func (j RelationPaginate) Type() RelationType {
	return RelationTypePaginate
}

// ToMap converts the instance to raw Field.
func (j RelationPaginate) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"input": j.Input,
		"fetch": j.Fetch,
		"skip":  j.Skip,
	}
}

// Encode returns the relation wrapper.
func (j RelationPaginate) Wrap() Relation {
	return NewRelation(&j)
}

// RelationProject represents a project relation.
type RelationProject struct {
	Input Relation               `json:"input" mapstructure:"input" yaml:"input"`
	Exprs []RelationalExpression `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationProject creates a RelationProject instance.
func NewRelationProject[R RelationInner](
	input R,
	expressions []RelationalExpressionInner,
) *RelationProject {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr == nil {
			continue
		}

		exprs = append(exprs, expr.Wrap())
	}

	return &RelationProject{
		Input: Relation{inner: input},
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationProject) Type() RelationType {
	return RelationTypeProject
}

// ToMap converts the instance to raw Field.
func (j RelationProject) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"input": j.Input,
		"exprs": j.Exprs,
	}
}

// Encode returns the relation wrapper.
func (j RelationProject) Wrap() Relation {
	return NewRelation(&j)
}

// RelationFilter represents a filter relation.
type RelationFilter struct {
	Input     Relation             `json:"input"     mapstructure:"input"     yaml:"input"`
	Predicate RelationalExpression `json:"predicate" mapstructure:"predicate" yaml:"predicate"`
}

// NewRelationFilter creates a RelationFilter instance.
func NewRelationFilter[R RelationInner, P RelationalExpressionInner](
	input R,
	predicate P,
) *RelationFilter {
	return &RelationFilter{
		Input:     NewRelation(input),
		Predicate: predicate.Wrap(),
	}
}

// Type return the type name of the instance.
func (j RelationFilter) Type() RelationType {
	return RelationTypeFilter
}

// ToMap converts the instance to raw Field.
func (j RelationFilter) ToMap() map[string]any {
	return map[string]any{
		"type":      j.Type(),
		"input":     j.Input,
		"predicate": j.Predicate,
	}
}

// Encode returns the relation wrapper.
func (j RelationFilter) Wrap() Relation {
	return NewRelation(&j)
}

// RelationSort represents a sort relation.
type RelationSort struct {
	Input Relation `json:"input" mapstructure:"input" yaml:"input"`
	Exprs []Sort   `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationSort creates a RelationSort instance.
func NewRelationSort[R RelationInner](input R, exprs []Sort) *RelationSort {
	return &RelationSort{
		Input: NewRelation(input),
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationSort) Type() RelationType {
	return RelationTypeSort
}

// ToMap converts the instance to raw Field.
func (j RelationSort) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"input": j.Input,
		"exprs": j.Exprs,
	}
}

// Encode returns the relation wrapper.
func (j RelationSort) Wrap() Relation {
	return NewRelation(&j)
}

// RelationJoin represents a join relation.
type RelationJoin struct {
	Left     Relation `json:"left"      mapstructure:"left"      yaml:"left"`
	Right    Relation `json:"right"     mapstructure:"right"     yaml:"right"`
	On       []JoinOn `json:"on"        mapstructure:"on"        yaml:"on"`
	JoinType JoinType `json:"join_type" mapstructure:"join_type" yaml:"join_type"`
}

// NewRelationJoin creates a RelationJoin instance.
func NewRelationJoin[L RelationInner, R RelationInner](
	left L,
	right R,
	on []JoinOn,
	joinType JoinType,
) *RelationJoin {
	return &RelationJoin{
		Left:     NewRelation(left),
		Right:    NewRelation(right),
		On:       on,
		JoinType: joinType,
	}
}

// Type return the type name of the instance.
func (j RelationJoin) Type() RelationType {
	return RelationTypeJoin
}

// ToMap converts the instance to raw Field.
func (j RelationJoin) ToMap() map[string]any {
	return map[string]any{
		"type":      j.Type(),
		"left":      j.Left,
		"right":     j.Right,
		"on":        j.On,
		"join_type": j.JoinType,
	}
}

// Encode returns the relation wrapper.
func (j RelationJoin) Wrap() Relation {
	return NewRelation(&j)
}

// RelationAggregate represents an aggregate relation.
type RelationAggregate struct {
	Input Relation `json:"input"      mapstructure:"input"      yaml:"input"`
	// Only non-empty if the 'relational_query.aggregate.group_by' capability is supported.
	GroupBy    []RelationalExpression `json:"group_by"   mapstructure:"group_by"   yaml:"group_by"`
	Aggregates []RelationalExpression `json:"aggregates" mapstructure:"aggregates" yaml:"aggregates"`
}

// NewRelationAggregate creates a RelationAggregate instance.
func NewRelationAggregate[R RelationInner](
	input R,
	groupBy []RelationalExpressionInner,
	aggregates []RelationalExpressionInner,
) *RelationAggregate {
	groups := []RelationalExpression{}
	aggs := []RelationalExpression{}

	for _, g := range groupBy {
		if g != nil {
			groups = append(groups, g.Wrap())
		}
	}

	for _, agg := range aggregates {
		if agg != nil {
			aggs = append(aggs, agg.Wrap())
		}
	}

	return &RelationAggregate{
		Input:      NewRelation(input),
		GroupBy:    groups,
		Aggregates: aggs,
	}
}

// Type return the type name of the instance.
func (j RelationAggregate) Type() RelationType {
	return RelationTypeAggregate
}

// ToMap converts the instance to raw Field.
func (j RelationAggregate) ToMap() map[string]any {
	return map[string]any{
		"type":       j.Type(),
		"input":      j.Input,
		"group_by":   j.GroupBy,
		"aggregates": j.Aggregates,
	}
}

// Encode returns the relation wrapper.
func (j RelationAggregate) Wrap() Relation {
	return NewRelation(&j)
}

// RelationWindow represents a window relation.
type RelationWindow struct {
	Input Relation               `json:"input" mapstructure:"input" yaml:"input"`
	Exprs []RelationalExpression `json:"exprs" mapstructure:"exprs" yaml:"exprs"`
}

// NewRelationWindow creates a RelationWindow instance.
func NewRelationWindow[R RelationInner](
	input R,
	expressions []RelationalExpressionInner,
) *RelationWindow {
	exprs := []RelationalExpression{}

	for _, expr := range expressions {
		if expr != nil {
			exprs = append(exprs, expr.Wrap())
		}
	}

	return &RelationWindow{
		Input: NewRelation(input),
		Exprs: exprs,
	}
}

// Type return the type name of the instance.
func (j RelationWindow) Type() RelationType {
	return RelationTypeWindow
}

// ToMap converts the instance to raw Field.
func (j RelationWindow) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"input": j.Input,
		"exprs": j.Exprs,
	}
}

// Encode returns the relation wrapper.
func (j RelationWindow) Wrap() Relation {
	return NewRelation(&j)
}

// RelationUnion represents a union relation.
type RelationUnion struct {
	Relations []Relation `json:"relations" mapstructure:"relations" yaml:"relations"`
}

// NewRelationUnion creates a RelationUnion instance.
func NewRelationUnion(relations []RelationInner) *RelationUnion {
	rels := []Relation{}

	for _, rel := range relations {
		if rel == nil {
			continue
		}

		rels = append(rels, NewRelation(rel))
	}

	return &RelationUnion{
		Relations: rels,
	}
}

// Type return the type name of the instance.
func (j RelationUnion) Type() RelationType {
	return RelationTypeUnion
}

// ToMap converts the instance to raw Field.
func (j RelationUnion) ToMap() map[string]any {
	return map[string]any{
		"type":      j.Type(),
		"relations": j.Relations,
	}
}

// Encode returns the relation wrapper.
func (j RelationUnion) Wrap() Relation {
	return NewRelation(&j)
}

// RelationalLiteralType represents a relation literal type enum.
type RelationalLiteralType string

const (
	RelationalLiteralTypeNull                 RelationalLiteralType = "null"
	RelationalLiteralTypeBoolean              RelationalLiteralType = "boolean"
	RelationalLiteralTypeString               RelationalLiteralType = "string"
	RelationalLiteralTypeInt8                 RelationalLiteralType = "int8"
	RelationalLiteralTypeInt16                RelationalLiteralType = "int16"
	RelationalLiteralTypeInt32                RelationalLiteralType = "int32"
	RelationalLiteralTypeInt64                RelationalLiteralType = "int64"
	RelationalLiteralTypeUint8                RelationalLiteralType = "uint8"
	RelationalLiteralTypeUint16               RelationalLiteralType = "uint16"
	RelationalLiteralTypeUint32               RelationalLiteralType = "uint32"
	RelationalLiteralTypeUint64               RelationalLiteralType = "uint64"
	RelationalLiteralTypeFloat32              RelationalLiteralType = "float32"
	RelationalLiteralTypeFloat64              RelationalLiteralType = "float64"
	RelationalLiteralTypeDecimal128           RelationalLiteralType = "decimal128"
	RelationalLiteralTypeDecimal256           RelationalLiteralType = "decimal256"
	RelationalLiteralTypeDate32               RelationalLiteralType = "date32"
	RelationalLiteralTypeDate64               RelationalLiteralType = "date64"
	RelationalLiteralTypeTime32Second         RelationalLiteralType = "time32_second"
	RelationalLiteralTypeTime32Millisecond    RelationalLiteralType = "time32_millisecond"
	RelationalLiteralTypeTime64Microsecond    RelationalLiteralType = "time64_microsecond"
	RelationalLiteralTypeTime64Nanosecond     RelationalLiteralType = "time64_nanosecond"
	RelationalLiteralTypeTimestampSecond      RelationalLiteralType = "timestamp_second"
	RelationalLiteralTypeTimestampMillisecond RelationalLiteralType = "timestamp_millisecond"
	RelationalLiteralTypeTimestampMicrosecond RelationalLiteralType = "timestamp_microsecond"
	RelationalLiteralTypeTimestampNanosecond  RelationalLiteralType = "timestamp_nanosecond"
	RelationalLiteralTypeDurationSecond       RelationalLiteralType = "duration_second"
	RelationalLiteralTypeDurationMillisecond  RelationalLiteralType = "duration_millisecond"
	RelationalLiteralTypeDurationMicrosecond  RelationalLiteralType = "duration_microsecond"
	RelationalLiteralTypeDurationNanosecond   RelationalLiteralType = "duration_nanosecond"
	RelationalLiteralTypeInterval             RelationalLiteralType = "interval"
)

var enumValues_RelationalLiteralType = []RelationalLiteralType{
	RelationalLiteralTypeNull,
	RelationalLiteralTypeBoolean,
	RelationalLiteralTypeString,
	RelationalLiteralTypeInt8,
	RelationalLiteralTypeInt16,
	RelationalLiteralTypeInt32,
	RelationalLiteralTypeInt64,
	RelationalLiteralTypeUint8,
	RelationalLiteralTypeUint16,
	RelationalLiteralTypeUint32,
	RelationalLiteralTypeUint64,
	RelationalLiteralTypeFloat32,
	RelationalLiteralTypeFloat64,
	RelationalLiteralTypeDecimal128,
	RelationalLiteralTypeDecimal256,
	RelationalLiteralTypeDate32,
	RelationalLiteralTypeDate64,
	RelationalLiteralTypeTime32Second,
	RelationalLiteralTypeTime32Millisecond,
	RelationalLiteralTypeTime64Microsecond,
	RelationalLiteralTypeTime64Nanosecond,
	RelationalLiteralTypeTimestampSecond,
	RelationalLiteralTypeTimestampMillisecond,
	RelationalLiteralTypeTimestampMicrosecond,
	RelationalLiteralTypeTimestampNanosecond,
	RelationalLiteralTypeDurationSecond,
	RelationalLiteralTypeDurationMillisecond,
	RelationalLiteralTypeDurationMicrosecond,
	RelationalLiteralTypeDurationNanosecond,
	RelationalLiteralTypeInterval,
}

// ParseRelationType parses a relation type from string.
func ParseRelationalLiteralType(input string) (RelationalLiteralType, error) {
	result := RelationalLiteralType(input)
	if !result.IsValid() {
		return RelationalLiteralType(
				"",
			), fmt.Errorf(
				"failed to parse RelationType, expect one of %v, got %s",
				enumValues_RelationType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j RelationalLiteralType) IsValid() bool {
	return slices.Contains(enumValues_RelationalLiteralType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationalLiteralType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseRelationalLiteralType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// RelationalLiteral is provided by reference to a relation literal.
type RelationalLiteral struct {
	inner RelationalLiteralInner
}

// NewRelationalLiteral creates a RelationalLiteral instance.
func NewRelationalLiteral[T RelationalLiteralInner](inner T) RelationalLiteral {
	return RelationalLiteral{
		inner: inner,
	}
}

// RelationalLiteralInner abstracts the interface for Relation.
type RelationalLiteralInner interface {
	Type() RelationalLiteralType
	ToMap() map[string]any
	Wrap() RelationalLiteral
}

// IsEmpty checks if the inner type is empty.
func (j RelationalLiteral) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j RelationalLiteral) Type() RelationalLiteralType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// Interface tries to convert the instance to AggregateInner interface.
func (j RelationalLiteral) Interface() RelationalLiteralInner {
	return j.inner
}

// MarshalJSON implements json.Marshaler interface.
func (j RelationalLiteral) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RelationalLiteral) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseRelationalLiteralType(rawType.Type)
	if err != nil {
		return fmt.Errorf("invalid RelationalLiteral type: %w", err)
	}

	switch ty {
	case RelationalLiteralTypeNull:
		j.inner = &RelationalLiteralNull{}
	case RelationalLiteralTypeBoolean:
		var inner RelationalLiteralNull

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralNull: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDate32:
		var inner RelationalLiteralDate32

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDate32: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDate64:
		var inner RelationalLiteralDate64

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDate64: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDecimal128:
		var inner RelationalLiteralDecimal128

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDecimal128: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDecimal256:
		var inner RelationalLiteralDecimal256

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDecimal256: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDurationMicrosecond:
		var inner RelationalLiteralDurationMicrosecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDurationMicrosecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDurationMillisecond:
		var inner RelationalLiteralDurationMillisecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDurationMillisecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDurationNanosecond:
		var inner RelationalLiteralDurationNanosecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDurationNanosecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeDurationSecond:
		var inner RelationalLiteralDurationSecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralDurationSecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeFloat32:
		var inner RelationalLiteralFloat32

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralFloat32: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeFloat64:
		var inner RelationalLiteralFloat64

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralFloat64: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeInt8:
		var inner RelationalLiteralInt8

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralInt8: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeInt16:
		var inner RelationalLiteralInt16

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralInt16: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeInt32:
		var inner RelationalLiteralInt32

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralInt32: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeInt64:
		var inner RelationalLiteralInt64

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralInt64: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeInterval:
		var inner RelationalLiteralInterval

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralInterval: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeString:
		var inner RelationalLiteralString

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralString: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTime32Millisecond:
		var inner RelationalLiteralTime32Millisecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTime32Millisecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTime32Second:
		var inner RelationalLiteralTime32Second

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTime32Second: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTime64Microsecond:
		var inner RelationalLiteralTime64Microsecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTime64Microsecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTime64Nanosecond:
		var inner RelationalLiteralTime64Nanosecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTime64Nanosecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTimestampMicrosecond:
		var inner RelationalLiteralTimestampMicrosecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTimestampMicrosecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTimestampMillisecond:
		var inner RelationalLiteralTimestampMillisecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTimestampMillisecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTimestampNanosecond:
		var inner RelationalLiteralTimestampNanosecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTimestampNanosecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeTimestampSecond:
		var inner RelationalLiteralTimestampSecond

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralTimestampSecond: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeUint8:
		var inner RelationalLiteralUint8

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralUint8: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeUint16:
		var inner RelationalLiteralUint16

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralUint16: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeUint32:
		var inner RelationalLiteralUint32

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralUint32: %w", err)
		}

		j.inner = &inner
	case RelationalLiteralTypeUint64:
		var inner RelationalLiteralUint64

		err := json.Unmarshal(b, &inner)
		if err != nil {
			return fmt.Errorf("failed to decode RelationalLiteralUint64: %w", err)
		}

		j.inner = &inner
	default:
		return fmt.Errorf("unsupported RelationalLiteral type: %s", ty)
	}

	return nil
}

// RelationalLiteralNull represents a RelationalLiteral null.
type RelationalLiteralNull struct{}

// NewRelationalLiteralNull creates a RelationalLiteralNull instance.
func NewRelationalLiteralNull() *RelationalLiteralNull {
	return &RelationalLiteralNull{}
}

// Type return the type name of the instance.
func (j RelationalLiteralNull) Type() RelationalLiteralType {
	return RelationalLiteralTypeNull
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralNull) ToMap() map[string]any {
	return map[string]any{
		"type": j.Type(),
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralNull) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralBoolean represents a RelationalLiteral boolean.
type RelationalLiteralBoolean struct {
	Value bool `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralBoolean creates a RelationalLiteralBoolean instance.
func NewRelationalLiteralBoolean(value bool) *RelationalLiteralBoolean {
	return &RelationalLiteralBoolean{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralBoolean) Type() RelationalLiteralType {
	return RelationalLiteralTypeBoolean
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralBoolean) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralBoolean) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralString represents a RelationalLiteral with utf-8 encoded string.
type RelationalLiteralString struct {
	Value string `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralString creates a RelationalLiteralString instance.
func NewRelationalLiteralString(value string) *RelationalLiteralString {
	return &RelationalLiteralString{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralString) Type() RelationalLiteralType {
	return RelationalLiteralTypeString
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralString) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralString) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralInt8 represents a RelationalLiteral with signed 8-bit int.
type RelationalLiteralInt8 struct {
	Value int8 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralInt8 creates a RelationalLiteralInt8 instance.
func NewRelationalLiteralInt8(value int8) *RelationalLiteralInt8 {
	return &RelationalLiteralInt8{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralInt8) Type() RelationalLiteralType {
	return RelationalLiteralTypeInt8
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralInt8) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralInt8) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralInt16 represents a RelationalLiteral with signed 16-bit int.
type RelationalLiteralInt16 struct {
	Value int16 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralInt16 creates a RelationalLiteralInt16 instance.
func NewRelationalLiteralInt16(value int16) *RelationalLiteralInt16 {
	return &RelationalLiteralInt16{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralInt16) Type() RelationalLiteralType {
	return RelationalLiteralTypeInt16
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralInt16) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralInt16) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralInt32 represents a RelationalLiteral with signed 32-bit int.
type RelationalLiteralInt32 struct {
	Value int32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralInt32 creates a RelationalLiteralInt32 instance.
func NewRelationalLiteralInt32(value int32) *RelationalLiteralInt32 {
	return &RelationalLiteralInt32{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralInt32) Type() RelationalLiteralType {
	return RelationalLiteralTypeInt32
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralInt32) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralInt32) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralInt64 represents a RelationalLiteral with signed 64-bit int.
type RelationalLiteralInt64 struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralInt64 creates a RelationalLiteralInt64 instance.
func NewRelationalLiteralInt64(value int64) *RelationalLiteralInt64 {
	return &RelationalLiteralInt64{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralInt64) Type() RelationalLiteralType {
	return RelationalLiteralTypeInt64
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralInt64) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralInt64) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralUint8 represents a RelationalLiteral with unsigned 8-bit int.
type RelationalLiteralUint8 struct {
	Value uint8 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralUint8 creates a RelationalLiteralUint8 instance.
func NewRelationalLiteralUint8(value uint8) *RelationalLiteralUint8 {
	return &RelationalLiteralUint8{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralUint8) Type() RelationalLiteralType {
	return RelationalLiteralTypeUint8
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralUint8) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralUint8) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralUint16 represents a RelationalLiteral with unsigned 16-bit int.
type RelationalLiteralUint16 struct {
	Value uint16 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralUint16 creates a RelationalLiteralUint16 instance.
func NewRelationalLiteralUint16(value uint16) *RelationalLiteralUint16 {
	return &RelationalLiteralUint16{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralUint16) Type() RelationalLiteralType {
	return RelationalLiteralTypeUint16
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralUint16) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralUint16) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralUint32 represents a RelationalLiteral with unsigned 32-bit int.
type RelationalLiteralUint32 struct {
	Value uint32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralUint32 creates a RelationalLiteralUint32 instance.
func NewRelationalLiteralUint32(value uint32) *RelationalLiteralUint32 {
	return &RelationalLiteralUint32{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralUint32) Type() RelationalLiteralType {
	return RelationalLiteralTypeUint32
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralUint32) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralUint32) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralUint64 represents a RelationalLiteral with unsigned 64-bit int.
type RelationalLiteralUint64 struct {
	Value uint64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralUint64 creates a RelationalLiteralUint64 instance.
func NewRelationalLiteralUint64(value uint64) *RelationalLiteralUint64 {
	return &RelationalLiteralUint64{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralUint64) Type() RelationalLiteralType {
	return RelationalLiteralTypeUint64
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralUint64) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralUint64) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralFloat32 represents a RelationalLiteral with unsigned 32-bit float.
type RelationalLiteralFloat32 struct {
	Value float32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralFloat32 creates a RelationalLiteralFloat32 instance.
func NewRelationalLiteralFloat32(value float32) *RelationalLiteralFloat32 {
	return &RelationalLiteralFloat32{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralFloat32) Type() RelationalLiteralType {
	return RelationalLiteralTypeFloat32
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralFloat32) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralFloat32) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralFloat64 represents a RelationalLiteral with unsigned 64-bit float.
type RelationalLiteralFloat64 struct {
	Value float64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralFloat64 creates a RelationalLiteralFloat64 instance.
func NewRelationalLiteralFloat64(value float64) *RelationalLiteralFloat64 {
	return &RelationalLiteralFloat64{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralFloat64) Type() RelationalLiteralType {
	return RelationalLiteralTypeFloat64
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralFloat64) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralFloat64) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDecimal128 represents a RelationalLiteral with unsigned 128-bit decimal.
type RelationalLiteralDecimal128 struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
	Scale int8  `json:"scale" mapstructure:"scale" yaml:"scale"`
	Spec  uint8 `json:"prec"  mapstructure:"prec"  yaml:"prec"`
}

// NewRelationalLiteralDecimal128 creates a RelationalLiteralDecimal128 instance.
func NewRelationalLiteralDecimal128(
	value int64,
	scale int8,
	spec uint8,
) *RelationalLiteralDecimal128 {
	return &RelationalLiteralDecimal128{
		Value: value,
		Scale: scale,
		Spec:  spec,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDecimal128) Type() RelationalLiteralType {
	return RelationalLiteralTypeDecimal128
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDecimal128) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
		"scale": j.Scale,
		"prec":  j.Spec,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDecimal128) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDecimal256 represents a RelationalLiteral with unsigned 256-bit decimal.
type RelationalLiteralDecimal256 struct {
	Value string `json:"value" mapstructure:"value" yaml:"value"`
	Scale int8   `json:"scale" mapstructure:"scale" yaml:"scale"`
	Spec  uint8  `json:"prec"  mapstructure:"prec"  yaml:"prec"`
}

// NewRelationalLiteralDecimal256 creates a RelationalLiteralDecimal256 instance.
func NewRelationalLiteralDecimal256(
	value string,
	scale int8,
	spec uint8,
) *RelationalLiteralDecimal256 {
	return &RelationalLiteralDecimal256{
		Value: value,
		Scale: scale,
		Spec:  spec,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDecimal256) Type() RelationalLiteralType {
	return RelationalLiteralTypeDecimal256
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDecimal256) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
		"scale": j.Scale,
		"prec":  j.Spec,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDecimal256) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDate32 represents a RelationalLiteral
// with Date stored as a signed 32bit int days since UNIX epoch 1970-01-01.
type RelationalLiteralDate32 struct {
	Value int32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDate32 creates a RelationalLiteralDate32 instance.
func NewRelationalLiteralDate32(value int32) *RelationalLiteralDate32 {
	return &RelationalLiteralDate32{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDate32) Type() RelationalLiteralType {
	return RelationalLiteralTypeDate32
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDate32) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDate32) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDate64 represents a RelationalLiteral
// with Date stored as a signed 64-bit int days since UNIX epoch 1970-01-01.
type RelationalLiteralDate64 struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDate64 creates a RelationalLiteralDate64 instance.
func NewRelationalLiteralDate64(value int64) *RelationalLiteralDate64 {
	return &RelationalLiteralDate64{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDate64) Type() RelationalLiteralType {
	return RelationalLiteralTypeDate64
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDate64) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDate64) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTime32Second represents a RelationalLiteral
// with Time stored as a signed 32-bit int as seconds since midnight.
type RelationalLiteralTime32Second struct {
	Value int32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTime32Second creates a RelationalLiteralTime32Second instance.
func NewRelationalLiteralTime32Second(value int32) *RelationalLiteralTime32Second {
	return &RelationalLiteralTime32Second{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTime32Second) Type() RelationalLiteralType {
	return RelationalLiteralTypeTime32Second
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTime32Second) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTime32Second) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTime32Millisecond represents a RelationalLiteral
// with Time stored as a signed 32-bit int as seconds since midnight.
type RelationalLiteralTime32Millisecond struct {
	Value int32 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTime32Millisecond creates a RelationalLiteralTime32Millisecond instance.
func NewRelationalLiteralTime32Millisecond(value int32) *RelationalLiteralTime32Millisecond {
	return &RelationalLiteralTime32Millisecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTime32Millisecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTime32Millisecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTime32Millisecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTime32Millisecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTime64Microsecond represents a RelationalLiteral
// with Time stored as a signed 64-bit int as microseconds since midnight.
type RelationalLiteralTime64Microsecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTime64Microsecond creates a RelationalLiteralTime64Microsecond instance.
func NewRelationalLiteralTime64Microsecond(value int64) *RelationalLiteralTime64Microsecond {
	return &RelationalLiteralTime64Microsecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTime64Microsecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTime64Microsecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTime64Microsecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTime64Microsecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTime64Nanosecond represents a RelationalLiteral
// with Time stored as a signed 64bit int as nanoseconds since midnight.
type RelationalLiteralTime64Nanosecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTime64Nanosecond creates a RelationalLiteralTime64Nanosecond instance.
func NewRelationalLiteralTime64Nanosecond(value int64) *RelationalLiteralTime64Nanosecond {
	return &RelationalLiteralTime64Nanosecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTime64Nanosecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTime64Nanosecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTime64Nanosecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTime64Nanosecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTimestampSecond represents a RelationalLiteral
// with Timestamp Second.
type RelationalLiteralTimestampSecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTimestampSecond creates a RelationalLiteralTimestampSecond instance.
func NewRelationalLiteralTimestampSecond(value int64) *RelationalLiteralTimestampSecond {
	return &RelationalLiteralTimestampSecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTimestampSecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTimestampSecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTimestampSecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTimestampSecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTimestampMillisecond represents a RelationalLiteral
// with Timestamp Milliseconds.
type RelationalLiteralTimestampMillisecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTimestampMillisecond creates a RelationalLiteralTimestampMillisecond instance.
func NewRelationalLiteralTimestampMillisecond(value int64) *RelationalLiteralTimestampMillisecond {
	return &RelationalLiteralTimestampMillisecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTimestampMillisecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTimestampMillisecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTimestampMillisecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTimestampMillisecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTimestampMicrosecond represents a RelationalLiteral
// with Timestamp Microseconds.
type RelationalLiteralTimestampMicrosecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTimestampMicrosecond creates a RelationalLiteralTimestampMicrosecond instance.
func NewRelationalLiteralTimestampMicrosecond(value int64) *RelationalLiteralTimestampMicrosecond {
	return &RelationalLiteralTimestampMicrosecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTimestampMicrosecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTimestampMicrosecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTimestampMicrosecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTimestampMicrosecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralTimestampNanosecond represents a RelationalLiteral
// with Timestamp Nanoseconds.
type RelationalLiteralTimestampNanosecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralTimestampNanosecond creates a RelationalLiteralTimestampNanosecond instance.
func NewRelationalLiteralTimestampNanosecond(value int64) *RelationalLiteralTimestampNanosecond {
	return &RelationalLiteralTimestampNanosecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralTimestampNanosecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeTimestampNanosecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralTimestampNanosecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralTimestampNanosecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDurationSecond represents a RelationalLiteral
// with Duration in seconds.
type RelationalLiteralDurationSecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDurationSecond creates a RelationalLiteralDurationSecond instance.
func NewRelationalLiteralDurationSecond(value int64) *RelationalLiteralDurationSecond {
	return &RelationalLiteralDurationSecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDurationSecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeDurationSecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDurationSecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDurationSecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDurationMillisecond represents a RelationalLiteral
// with Duration in milliseconds.
type RelationalLiteralDurationMillisecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDurationMillisecond creates a RelationalLiteralDurationMillisecond instance.
func NewRelationalLiteralDurationMillisecond(value int64) *RelationalLiteralDurationMillisecond {
	return &RelationalLiteralDurationMillisecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDurationMillisecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeDurationMillisecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDurationMillisecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDurationMillisecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDurationMicrosecond represents a RelationalLiteral
// with Duration in microseconds.
type RelationalLiteralDurationMicrosecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDurationMicrosecond creates a RelationalLiteralDurationMicrosecond instance.
func NewRelationalLiteralDurationMicrosecond(value int64) *RelationalLiteralDurationMicrosecond {
	return &RelationalLiteralDurationMicrosecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDurationMicrosecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeDurationMicrosecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDurationMicrosecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDurationMicrosecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralDurationNanosecond represents a RelationalLiteral
// with Duration in nanoseconds.
type RelationalLiteralDurationNanosecond struct {
	Value int64 `json:"value" mapstructure:"value" yaml:"value"`
}

// NewRelationalLiteralDurationNanosecond creates a RelationalLiteralDurationNanosecond instance.
func NewRelationalLiteralDurationNanosecond(value int64) *RelationalLiteralDurationNanosecond {
	return &RelationalLiteralDurationNanosecond{
		Value: value,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralDurationNanosecond) Type() RelationalLiteralType {
	return RelationalLiteralTypeDurationNanosecond
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralDurationNanosecond) ToMap() map[string]any {
	return map[string]any{
		"type":  j.Type(),
		"value": j.Value,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralDurationNanosecond) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}

// RelationalLiteralInterval represents a RelationalLiteral
// with Interval represented as months, days, and nanoseconds.
type RelationalLiteralInterval struct {
	Months      int32 `json:"months"      mapstructure:"months"      yaml:"months"`
	Days        int32 `json:"days"        mapstructure:"days"        yaml:"days"`
	Nanoseconds int64 `json:"nanoseconds" mapstructure:"nanoseconds" yaml:"nanoseconds"`
}

// NewRelationalLiteralInterval creates a RelationalLiteralInterval instance.
func NewRelationalLiteralInterval(
	months int32,
	days int32,
	nanoseconds int64,
) *RelationalLiteralInterval {
	return &RelationalLiteralInterval{
		Months:      months,
		Days:        days,
		Nanoseconds: nanoseconds,
	}
}

// Type return the type name of the instance.
func (j RelationalLiteralInterval) Type() RelationalLiteralType {
	return RelationalLiteralTypeInterval
}

// ToMap converts the instance to raw Field.
func (j RelationalLiteralInterval) ToMap() map[string]any {
	return map[string]any{
		"type":        j.Type(),
		"months":      j.Months,
		"days":        j.Days,
		"nanoseconds": j.Nanoseconds,
	}
}

// Encode returns the relation wrapper.
func (j RelationalLiteralInterval) Wrap() RelationalLiteral {
	return NewRelationalLiteral(&j)
}
