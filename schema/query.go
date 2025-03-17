package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// MarshalJSON implements json.Marshaler.
func (j RowSet) MarshalJSON() ([]byte, error) {
	result := map[string]any{}
	if len(j.Aggregates) > 0 {
		result["aggregates"] = j.Aggregates
	}
	if j.Rows != nil {
		result["rows"] = j.Rows
	}

	return json.Marshal(result)
}

// UnmarshalJSONMap decodes FunctionInfo from a JSON map.
func (j *FunctionInfo) UnmarshalJSONMap(raw map[string]json.RawMessage) error {
	rawArguments, ok := raw["arguments"]
	var arguments FunctionInfoArguments
	if ok && !isNullJSON(rawArguments) {
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("FunctionInfo.arguments: %w", err)
		}
	}

	rawName, ok := raw["name"]
	if !ok || isNullJSON(rawName) {
		return errors.New("FunctionInfo.name: required")
	}
	var name string
	if err := json.Unmarshal(rawName, &name); err != nil {
		return fmt.Errorf("FunctionInfo.name: %w", err)
	}
	if name == "" {
		return errors.New("FunctionInfo.name: required")
	}

	rawDescription, ok := raw["description"]
	var description *string
	if ok && !isNullJSON(rawDescription) {
		if err := json.Unmarshal(rawDescription, &description); err != nil {
			return fmt.Errorf("FunctionInfo.description: %w", err)
		}
	}

	rawResultType, ok := raw["result_type"]
	if !ok || isNullJSON(rawResultType) {
		return errors.New("FunctionInfo.result_type: required")
	}
	var resultType Type
	if err := json.Unmarshal(rawResultType, &resultType); err != nil {
		return fmt.Errorf("FunctionInfo.result_type: %w", err)
	}

	*j = FunctionInfo{
		Arguments:   arguments,
		Name:        name,
		ResultType:  resultType,
		Description: description,
	}
	return nil
}

// ArrayComparisonType represents a type of ArrayComparison
type ArrayComparisonType string

const (
	// Check if the array contains the specified value. Only used if the 'query.nested_fields.filter_by.nested_arrays.contains' capability is supported.
	ArrayComparisonTypeContains ArrayComparisonType = "contains"
	// Check is the array is empty. Only used if the 'query.nested_fields.filter_by.nested_arrays.is_empty' capability is supported.
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
		return ArrayComparisonType(""), fmt.Errorf("failed to parse ArrayComparisonType, expect one of %v, got %s", enumValues_ArrayComparisonType, input)
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

// ArrayComparisonEncoder abstracts a generic interface of ArrayComparison
type ArrayComparisonEncoder interface {
	Encode() ArrayComparison
}

// ArrayComparison represents an array comparison
type ArrayComparison map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArrayComparison) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty ArrayComparisonType
	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in ArrayComparison: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in ArrayComparison: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case ArrayComparisonTypeContains:
		rawValue, ok := raw["value"]
		if !ok {
			return errors.New("field value in ArrayComparison: required")
		}

		var value ComparisonValue
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in ArrayComparison: %w", err)
		}

		results["value"] = value
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
		return nil, fmt.Errorf("invalid ArrayComparison type; expected %s, got %s", ArrayComparisonTypeContains, t)
	}
	rawValue, ok := j["value"]
	if !ok {
		return nil, errors.New("GroupOrderByTarget.index is required")
	}

	value, ok := rawValue.(ComparisonValue)
	if !ok {
		return nil, fmt.Errorf("invalid ArrayComparisonContains.value, expected ComparisonValue, got %v", rawValue)
	}

	result := &ArrayComparisonContains{
		Type:  t,
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
		return nil, fmt.Errorf("invalid ArrayComparison type; expected %s, got %s", ArrayComparisonTypeIsEmpty, t)
	}

	result := &ArrayComparisonIsEmpty{
		Type: t,
	}

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
	Type  ArrayComparisonType `json:"type" yaml:"type" mapstructure:"type"`
	Value ComparisonValue     `json:"value" yaml:"value" mapstructure:"value"`
}

// NewArrayComparisonContains creates an ArrayComparisonContains instance.
func NewArrayComparisonContains(value ComparisonValueEncoder) *ArrayComparisonContains {
	return &ArrayComparisonContains{
		Type:  ArrayComparisonTypeContains,
		Value: value.Encode(),
	}
}

// Encode converts the instance to raw ArrayComparison.
func (j ArrayComparisonContains) Encode() ArrayComparison {
	result := ArrayComparison{
		"type":  j.Type,
		"value": j.Value,
	}
	return result
}

// ArrayComparisonIsEmpty checks if the array is empty. Only used if the 'query.nested_fields.filter_by.nested_arrays.is_empty' capability is supported.
type ArrayComparisonIsEmpty struct {
	Type  ArrayComparisonType `json:"type" yaml:"type" mapstructure:"type"`
	Value ComparisonValue     `json:"value" yaml:"value" mapstructure:"value"`
}

// NewArrayComparisonIsEmpty creates an ArrayComparisonIsEmpty instance.
func NewArrayComparisonIsEmpty() *ArrayComparisonIsEmpty {
	return &ArrayComparisonIsEmpty{
		Type: ArrayComparisonTypeIsEmpty,
	}
}

// Encode converts the instance to raw ArrayComparison.
func (j ArrayComparisonIsEmpty) Encode() ArrayComparison {
	result := ArrayComparison{
		"type": j.Type,
	}
	return result
}

// DimensionType represents a type of Dimension
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
		return DimensionType(""), fmt.Errorf("failed to parse DimensionType, expect one of %v, got %s", enumValues_DimensionType, input)
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

// DimensionEncoder abstracts a generic interface of Dimension
type DimensionEncoder interface {
	Encode() Dimension
}

// Dimension represents a dimension object
type Dimension map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *Dimension) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty DimensionType
	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in DimensionColumn: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in DimensionColumn: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case DimensionTypeColumn:
		rawValue, ok := raw["value"]
		if !ok {
			return errors.New("field value in ArrayComparison: required")
		}

		var value ComparisonValue
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return fmt.Errorf("field value in ArrayComparison: %w", err)
		}

		results["value"] = value
	default:
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
func (j Dimension) Type() (DimensionType, error) {
	t, ok := j["type"]
	if !ok {
		return DimensionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseDimensionType(raw)
		if err != nil {
			return DimensionType(""), err
		}
		return v, nil
	case DimensionType:
		return raw, nil
	default:
		return DimensionType(""), fmt.Errorf("invalid Dimension type: %+v", t)
	}
}

// AsColumn tries to convert the current type to DimensionColumn.
func (j Dimension) AsColumn() (*DimensionColumn, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != DimensionTypeColumn {
		return nil, fmt.Errorf("invalid Dimension type; expected %s, got %s", DimensionTypeColumn, t)
	}
	columnName, err := getStringValueByKey(j, "column_name")
	if err != nil {
		return nil, fmt.Errorf("DimensionColumn.column_name: %w", err)
	}

	if columnName == "" {
		return nil, errors.New("DimensionColumn.column_name is required")
	}

	result := &DimensionColumn{
		Type:       t,
		ColumnName: columnName,
	}
	rawPath, ok := j["path"]
	if !ok {
		return nil, errors.New("DimensionColumn.path is required")
	}
	if rawPath != nil {
		p, ok := rawPath.([]PathElement)
		if !ok {
			return nil, fmt.Errorf("invalid DimensionColumn.path, expected []PathElement, got %v", rawPath)
		}
		result.Path = p
	}

	rawArguments, ok := j["arguments"]
	if ok && rawArguments != nil {
		arguments, ok := rawArguments.(map[string]Argument)
		if !ok {
			return nil, fmt.Errorf("invalid DimensionColumn.arguments, expected map[string]Argument, got %v", rawArguments)
		}
		result.Arguments = arguments
	}

	rawFieldPath, ok := j["field_path"]
	if ok && rawFieldPath != nil {
		fieldPath, ok := rawFieldPath.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid DimensionColumn.field_path, expected []string, got %v", rawFieldPath)
		}
		result.FieldPath = fieldPath
	}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j Dimension) Interface() DimensionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j Dimension) InterfaceT() (DimensionEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case DimensionTypeColumn:
		return j.AsColumn()
	default:
		return nil, fmt.Errorf("invalid Dimension type: %s", ty)
	}
}

// DimensionColumn represents a dimension column
type DimensionColumn struct {
	Type DimensionType `json:"type" yaml:"type" mapstructure:"type"`
	// Any (object) relationships to traverse to reach this column. Only non-empty if the 'relationships' capability is supported.
	Path []PathElement `json:"path" yaml:"path" mapstructure:"path"`
	// The name of the column
	ColumnName string `json:"column_name" yaml:"column_name" mapstructure:"column_name"`
	// Arguments to satisfy the column specified by 'column_name'
	Arguments map[string]Argument `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	// Path to a nested field within an object column
	FieldPath []string `json:"field_path" yaml:"field_path" mapstructure:"field_path"`
}

// NewDimensionColumn creates a DimensionColumn instance.
func NewDimensionColumn(columnName string, arguments map[string]Argument, fieldPath []string, path []PathElement) *DimensionColumn {
	return &DimensionColumn{
		Type:       DimensionTypeColumn,
		ColumnName: columnName,
		Arguments:  arguments,
		Path:       path,
		FieldPath:  fieldPath,
	}
}

// Encode converts the instance to a raw Dimension.
func (j DimensionColumn) Encode() Dimension {
	result := Dimension{
		"type":        j.Type,
		"path":        j.Path,
		"column_name": j.ColumnName,
		"arguments":   j.Arguments,
		"field_path":  j.FieldPath,
	}
	return result
}
