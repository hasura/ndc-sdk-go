package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
)

// MarshalJSON implements json.Marshaler.
func (j RowSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.ToMap())
}

// ToMap encodes the struct to a value map.
func (j RowSet) ToMap() map[string]any {
	result := map[string]any{}

	if len(j.Aggregates) > 0 {
		result["aggregates"] = j.Aggregates
	}

	if j.Rows != nil {
		result["rows"] = j.Rows
	}

	if j.Groups != nil {
		groups := make([]map[string]any, len(j.Groups))

		for i, group := range j.Groups {
			groups[i] = group.ToMap()
		}

		result["groups"] = groups
	}

	return result
}

// ToMap encodes the struct to a value map.
func (j Group) ToMap() map[string]any {
	result := map[string]any{
		"aggregates": j.Aggregates,
		"dimensions": j.Dimensions,
	}

	return result
}

// FromValue decodes the raw object value to the instance.
func (pe *PathElement) FromValue(raw map[string]any) error {
	arguments, err := getRelationshipArgumentMapByKey(raw, "arguments")
	if err != nil {
		return fmt.Errorf("field arguments in PathElement: %w", err)
	}

	if arguments == nil {
		return errors.New("field arguments in PathElement is required")
	}

	pe.Arguments = arguments

	fieldPath, err := getStringSliceByKey(raw, "field_path")
	if err != nil {
		return fmt.Errorf("field field_path in PathElement: %w", err)
	}

	pe.FieldPath = fieldPath

	relationship, err := getStringValueByKey(raw, "relationship")
	if err != nil {
		return fmt.Errorf("field relationship in PathElement: %w", err)
	}

	pe.Relationship = relationship

	rawExpression, ok := raw["predicate"]
	if !ok {
		return nil
	}

	expression, ok := rawExpression.(Expression)
	if !ok {
		exprMap, ok := rawExpression.(map[string]any)
		if !ok {
			return errors.New("field predicate in PathElement is required")
		}

		if err := expression.FromValue(exprMap); err != nil {
			return fmt.Errorf("field predicate in PathElement: %w", err)
		}
	}

	pe.Predicate = expression

	return nil
}

// ArrayComparisonType represents a type of ArrayComparison.
type ArrayComparisonType string

const (
	// ArrayComparisonTypeContains check if the array contains the specified value. Only used if the 'query.nested_fields.filter_by.nested_arrays.contains' capability is supported.
	ArrayComparisonTypeContains ArrayComparisonType = "contains"
	// ArrayComparisonTypeIsEmpty check is the array is empty. Only used if the 'query.nested_fields.filter_by.nested_arrays.is_empty' capability is supported.
	ArrayComparisonTypeIsEmpty ArrayComparisonType = "is_empty"
)

var enumValues_ArrayComparisonType = []ArrayComparisonType{
	ArrayComparisonTypeContains,
	ArrayComparisonTypeIsEmpty,
}

// ParseArrayComparisonType parses an ArrayComparisonType from string.
func ParseArrayComparisonType(input string) (ArrayComparisonType, error) {
	result := ArrayComparisonType(input)
	if !result.IsValid() {
		return ArrayComparisonType(
				"",
			), fmt.Errorf(
				"failed to parse ArrayComparisonType, expect one of %v, got %s",
				enumValues_ArrayComparisonType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j ArrayComparisonType) IsValid() bool {
	return slices.Contains(enumValues_ArrayComparisonType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArrayComparisonType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseArrayComparisonType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// ArrayComparisonEncoder abstracts a generic interface of ArrayComparison.
type ArrayComparisonEncoder interface {
	Type() ArrayComparisonType
	Encode() ArrayComparison
}

// ArrayComparison represents an array comparison.
type ArrayComparison map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArrayComparison) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return j.FromValue(raw)
}

// FromValue maps the raw object value to the instance.
func (j *ArrayComparison) FromValue(raw map[string]any) error {
	rawType, err := getStringValueByKey(raw, "type")
	if err != nil {
		return fmt.Errorf("field type in ArrayComparison: %w", err)
	}

	ty, err := ParseArrayComparisonType(rawType)
	if err != nil {
		return fmt.Errorf("field type in ArrayComparison: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case ArrayComparisonTypeContains:
		ac, err := ArrayComparison(raw).asContains()
		if err != nil {
			return err
		}

		results = ac.Encode()
	default:
	}

	*j = results

	return nil
}

// Type gets the type enum of the current type.
func (j ArrayComparison) Type() (ArrayComparisonType, error) {
	t, ok := j["type"]
	if !ok {
		return ArrayComparisonType(""), errTypeRequired
	}

	switch raw := t.(type) {
	case string:
		v, err := ParseArrayComparisonType(raw)
		if err != nil {
			return ArrayComparisonType(""), err
		}

		return v, nil
	case ArrayComparisonType:
		return raw, nil
	default:
		return ArrayComparisonType(""), fmt.Errorf("invalid ArrayComparison type: %+v", t)
	}
}

// AsContains tries to convert the current type to ArrayComparisonContains.
func (j ArrayComparison) AsContains() (*ArrayComparisonContains, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ArrayComparisonTypeContains {
		return nil, fmt.Errorf(
			"invalid ArrayComparison type; expected %s, got %s",
			ArrayComparisonTypeContains,
			t,
		)
	}

	return j.asContains()
}

func (j ArrayComparison) asContains() (*ArrayComparisonContains, error) {
	rawValue, ok := j["value"]
	if !ok {
		return nil, errors.New("field value in GroupOrderByTarget is required")
	}

	value, ok := rawValue.(ComparisonValue)
	if !ok {
		rawValueMap, ok := rawValue.(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"field value in ArrayComparisonContains: expected object, got %v",
				rawValue,
			)
		}

		if err := value.FromValue(rawValueMap); err != nil {
			return nil, fmt.Errorf("field value in ArrayComparisonContains: %w", err)
		}
	}

	result := &ArrayComparisonContains{
		Value: value,
	}

	return result, nil
}

// AsIsEmpty tries to convert the current type to ArrayComparisonIsEmpty.
func (j ArrayComparison) AsIsEmpty() (*ArrayComparisonIsEmpty, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	if t != ArrayComparisonTypeIsEmpty {
		return nil, fmt.Errorf(
			"invalid ArrayComparison type; expected %s, got %s",
			ArrayComparisonTypeIsEmpty,
			t,
		)
	}

	result := &ArrayComparisonIsEmpty{}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j ArrayComparison) Interface() ArrayComparisonEncoder {
	result, _ := j.InterfaceT()

	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j ArrayComparison) InterfaceT() (ArrayComparisonEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case ArrayComparisonTypeContains:
		return j.AsContains()
	case ArrayComparisonTypeIsEmpty:
		return j.AsIsEmpty()
	default:
		return nil, fmt.Errorf("invalid GroupOrderByTarget type: %s", ty)
	}
}

// ArrayComparisonContains check if the array contains the specified value. Only used if the 'query.nested_fields.filter_by.nested_arrays.contains' capability is supported.
type ArrayComparisonContains struct {
	Value ComparisonValue `json:"value" mapstructure:"value" yaml:"value"`
}

// NewArrayComparisonContains creates an ArrayComparisonContains instance.
func NewArrayComparisonContains[T ComparisonValueEncoder](value T) *ArrayComparisonContains {
	return &ArrayComparisonContains{
		Value: value.Encode(),
	}
}

// Type return the type name of the instance.
func (j ArrayComparisonContains) Type() ArrayComparisonType {
	return ArrayComparisonTypeContains
}

// Encode converts the instance to raw ArrayComparison.
func (j ArrayComparisonContains) Encode() ArrayComparison {
	result := ArrayComparison{
		"type":  j.Type(),
		"value": j.Value,
	}

	return result
}

// ArrayComparisonIsEmpty checks if the array is empty.
// Only used if the 'query.nested_fields.filter_by.nested_arrays.is_empty' capability is supported.
type ArrayComparisonIsEmpty struct{}

// NewArrayComparisonIsEmpty creates an ArrayComparisonIsEmpty instance.
func NewArrayComparisonIsEmpty() *ArrayComparisonIsEmpty {
	return &ArrayComparisonIsEmpty{}
}

// Type return the type name of the instance.
func (j ArrayComparisonIsEmpty) Type() ArrayComparisonType {
	return ArrayComparisonTypeIsEmpty
}

// Encode converts the instance to raw ArrayComparison.
func (j ArrayComparisonIsEmpty) Encode() ArrayComparison {
	result := ArrayComparison{
		"type": j.Type(),
	}

	return result
}

// DimensionType represents a type of Dimension.
type DimensionType string

const (
	DimensionTypeColumn DimensionType = "column"
)

var enumValues_DimensionType = []DimensionType{
	DimensionTypeColumn,
}

// ParseDimensionType parses a DimensionType from string.
func ParseDimensionType(input string) (DimensionType, error) {
	result := DimensionType(input)
	if !result.IsValid() {
		return DimensionType(
				"",
			), fmt.Errorf(
				"failed to parse DimensionType, expect one of %v, got %s",
				enumValues_DimensionType,
				input,
			)
	}

	return result, nil
}

// IsValid checks if the value is invalid.
func (j DimensionType) IsValid() bool {
	return slices.Contains(enumValues_DimensionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *DimensionType) UnmarshalJSON(b []byte) error {
	var rawValue string

	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseDimensionType(rawValue)
	if err != nil {
		return err
	}

	*j = value

	return nil
}

// DimensionInner abstracts a generic interface of Dimension.
type DimensionInner interface {
	Type() DimensionType
	ToMap() map[string]any
	Wrap() Dimension
}

// Dimension represents a dimension object.
type Dimension struct {
	inner DimensionInner
}

// NewDimension creates a new Dimension instance.
func NewDimension[T DimensionInner](inner T) Dimension {
	return Dimension{
		inner: inner,
	}
}

// Equal checks whether instances are the same.
func (j Dimension) Equal(value Dimension) bool {
	return reflect.DeepEqual(j, value)
}

// IsEmpty checks if the inner type is empty.
func (j Dimension) IsEmpty() bool {
	return j.inner == nil
}

// Type gets the type enum of the current type.
func (j Dimension) Type() DimensionType {
	if j.inner != nil {
		return j.inner.Type()
	}

	return ""
}

// Interface converts the comparison value to its generic interface.
func (j Dimension) Interface() DimensionInner {
	return j.inner
}

// MarshalJSON implements json.Marshaler interface.
func (j Dimension) MarshalJSON() ([]byte, error) {
	if j.inner == nil {
		return json.Marshal(nil)
	}

	return json.Marshal(j.inner.ToMap())
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Dimension) UnmarshalJSON(b []byte) error {
	var rawType rawTypeStruct

	if err := json.Unmarshal(b, &rawType); err != nil {
		return err
	}

	ty, err := ParseDimensionType(rawType.Type)
	if err != nil {
		return fmt.Errorf("field type in Dimension: %w", err)
	}

	switch ty {
	case DimensionTypeColumn:
		var result DimensionColumn

		err := json.Unmarshal(b, &result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Dimension: %w", err)
		}

		j.inner = &result
	default:
		return fmt.Errorf("unsupported dimension type: %s", ty)
	}

	return nil
}

// DimensionColumn represents a dimension column.
type DimensionColumn struct {
	// Any (object) relationships to traverse to reach this column. Only non-empty if the 'relationships' capability is supported.
	Path []PathElement `json:"path"                 mapstructure:"path"        yaml:"path"`
	// The name of the column
	ColumnName string `json:"column_name"          mapstructure:"column_name" yaml:"column_name"`
	// Arguments to satisfy the column specified by 'column_name'
	Arguments map[string]Argument `json:"arguments,omitempty"  mapstructure:"arguments"   yaml:"arguments,omitempty"`
	// Path to a nested field within an object column.
	FieldPath []string `json:"field_path,omitempty" mapstructure:"field_path"  yaml:"field_path,omitempty"`
	// The name of the extraction function to apply to the selected value, if any.
	Extraction string `json:"extraction,omitempty" mapstructure:"extraction"  yaml:"extraction,omitempty"`
}

// NewDimensionColumn creates a DimensionColumn instance.
func NewDimensionColumn(columnName string, path []PathElement) *DimensionColumn {
	if path == nil {
		path = []PathElement{}
	}

	return &DimensionColumn{
		ColumnName: columnName,
		Path:       path,
	}
}

// WithFieldPath return a new column field with field_path set.
func (f DimensionColumn) WithFieldPath(fieldPath []string) *DimensionColumn {
	f.FieldPath = fieldPath

	return &f
}

// WithExtraction return a new column field with extraction set.
func (f DimensionColumn) WithExtraction(extraction string) *DimensionColumn {
	f.Extraction = extraction

	return &f
}

// WithArguments return a new column field with arguments set.
func (f DimensionColumn) WithArguments(arguments map[string]ArgumentEncoder) *DimensionColumn {
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
func (f DimensionColumn) WithArgument(key string, argument ArgumentEncoder) *DimensionColumn {
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
func (j DimensionColumn) Type() DimensionType {
	return DimensionTypeColumn
}

// ToMap converts the instance to a raw map.
func (j DimensionColumn) ToMap() map[string]any {
	result := map[string]any{
		"type":        j.Type(),
		"path":        j.Path,
		"column_name": j.ColumnName,
	}

	if j.Arguments != nil {
		result["arguments"] = j.Arguments
	}

	if j.FieldPath != nil {
		result["field_path"] = j.Extraction
	}

	if j.Extraction != "" {
		result["extraction"] = j.Extraction
	}

	return result
}

// Wrap returns the DimensionColumn wrapper.
func (j DimensionColumn) Wrap() Dimension {
	return NewDimension(&j)
}
