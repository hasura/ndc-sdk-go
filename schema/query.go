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

// AggregateFunctionDefinitionType represents a type of GroupOrderByTarget
type AggregateFunctionDefinitionType string

const (
	AggregateFunctionDefinitionTypeMin     AggregateFunctionDefinitionType = "min"
	AggregateFunctionDefinitionTypeMax     AggregateFunctionDefinitionType = "max"
	AggregateFunctionDefinitionTypeSum     AggregateFunctionDefinitionType = "sum"
	AggregateFunctionDefinitionTypeAverage AggregateFunctionDefinitionType = "average"
	AggregateFunctionDefinitionTypeCustom  AggregateFunctionDefinitionType = "custom"
)

var enumValues_AggregateFunctionDefinitionType = []AggregateFunctionDefinitionType{
	AggregateFunctionDefinitionTypeMin,
	AggregateFunctionDefinitionTypeMax,
	AggregateFunctionDefinitionTypeSum,
	AggregateFunctionDefinitionTypeAverage,
	AggregateFunctionDefinitionTypeCustom,
}

// ParseAggregateFunctionDefinitionType parses a AggregateFunctionDefinitionType from string.
func ParseAggregateFunctionDefinitionType(input string) (AggregateFunctionDefinitionType, error) {
	result := AggregateFunctionDefinitionType(input)
	if !result.IsValid() {
		return AggregateFunctionDefinitionType(""), fmt.Errorf("failed to parse AggregateFunctionDefinitionType, expect one of %v, got %s", enumValues_AggregateFunctionDefinitionType, input)
	}
	return result, nil
}

// IsValid checks if the value is invalid.
func (j AggregateFunctionDefinitionType) IsValid() bool {
	return slices.Contains(enumValues_AggregateFunctionDefinitionType, j)
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateFunctionDefinitionType) UnmarshalJSON(b []byte) error {
	var rawValue string
	if err := json.Unmarshal(b, &rawValue); err != nil {
		return err
	}

	value, err := ParseAggregateFunctionDefinitionType(rawValue)
	if err != nil {
		return err
	}

	*j = value
	return nil
}

// AggregateFunctionDefinitionEncoder abstracts a generic interface of AggregateFunctionDefinition
type AggregateFunctionDefinitionEncoder interface {
	Encode() AggregateFunctionDefinition
}

// AggregateFunctionDefinition represents the definition of an aggregation function on a scalar type
type AggregateFunctionDefinition map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *AggregateFunctionDefinition) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var ty AggregateFunctionDefinitionType

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in AggregateFunctionDefinition: required")
	}
	err := json.Unmarshal(rawType, &ty)
	if err != nil {
		return fmt.Errorf("field type in AggregateFunctionDefinition: %w", err)
	}

	results := map[string]any{
		"type": ty,
	}

	switch ty {
	case AggregateFunctionDefinitionTypeSum, AggregateFunctionDefinitionTypeAverage:
		rawResultType, ok := raw["result_type"]
		if !ok || len(rawResultType) == 0 {
			return errors.New("field result_type in AggregateFunctionDefinition: required")
		}

		var resultType string
		if err := json.Unmarshal(rawResultType, &resultType); err != nil {
			return fmt.Errorf("field result_type in AggregateFunctionDefinition: %w", err)
		}

		results["result_type"] = resultType
	case AggregateFunctionDefinitionTypeCustom:
		rawResultType, ok := raw["result_type"]
		if !ok || len(rawResultType) == 0 {
			return errors.New("field result_type in AggregateFunctionDefinition: required")
		}

		var resultType Type
		if err := resultType.UnmarshalJSON(rawResultType); err != nil {
			return fmt.Errorf("field result_type in AggregateFunctionDefinition: %w", err)
		}

		results["result_type"] = resultType
	}

	*j = results
	return nil
}

// Type gets the type enum of the current type.
func (j AggregateFunctionDefinition) Type() (AggregateFunctionDefinitionType, error) {
	t, ok := j["type"]
	if !ok {
		return AggregateFunctionDefinitionType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseAggregateFunctionDefinitionType(raw)
		if err != nil {
			return AggregateFunctionDefinitionType(""), err
		}
		return v, nil
	case AggregateFunctionDefinitionType:
		return raw, nil
	default:
		return AggregateFunctionDefinitionType(""), fmt.Errorf("invalid AggregateFunctionDefinition type: %+v", t)
	}
}

// AsMin tries to convert the current type to AggregateFunctionDefinitionMin.
func (j AggregateFunctionDefinition) AsMin() (*AggregateFunctionDefinitionMin, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeMin {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeMin, t)
	}

	result := &AggregateFunctionDefinitionMin{
		Type: t,
	}

	return result, nil
}

// AsMax tries to convert the current type to AggregateFunctionDefinitionMax.
func (j AggregateFunctionDefinition) AsMax() (*AggregateFunctionDefinitionMax, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeMax {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeMax, t)
	}

	result := &AggregateFunctionDefinitionMax{
		Type: t,
	}

	return result, nil
}

// AsSum tries to convert the current type to AggregateFunctionDefinitionSum.
func (j AggregateFunctionDefinition) AsSum() (*AggregateFunctionDefinitionSum, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeSum {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeSum, t)
	}

	resultType := getStringValueByKey(j, "result_type")
	if resultType == "" {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionSum: required")
	}
	result := &AggregateFunctionDefinitionSum{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// AsAverage tries to convert the current type to AggregateFunctionDefinitionAverage.
func (j AggregateFunctionDefinition) AsAverage() (*AggregateFunctionDefinitionAverage, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeAverage {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeAverage, t)
	}

	resultType := getStringValueByKey(j, "result_type")
	if resultType == "" {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionAverage: required")
	}
	result := &AggregateFunctionDefinitionAverage{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// AsCustom tries to convert the current type to AggregateFunctionDefinitionCustom.
func (j AggregateFunctionDefinition) AsCustom() (*AggregateFunctionDefinitionCustom, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != AggregateFunctionDefinitionTypeCustom {
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type; expected %s, got %s", AggregateFunctionDefinitionTypeCustom, t)
	}

	rawResultType, ok := j["result_type"]
	if !ok {
		return nil, errors.New("field result_type in AggregateFunctionDefinitionCustom: required")
	}

	resultType, ok := rawResultType.(Type)
	if !ok {
		return nil, fmt.Errorf("invalid result_type in AggregateFunctionDefinitionCustom, expected Type, got %v", rawResultType)
	}
	result := &AggregateFunctionDefinitionCustom{
		Type:       t,
		ResultType: resultType,
	}

	return result, nil
}

// Interface converts the comparison value to its generic interface.
func (j AggregateFunctionDefinition) Interface() AggregateFunctionDefinitionEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT converts the comparison value to its generic interface safely with explicit error.
func (j AggregateFunctionDefinition) InterfaceT() (AggregateFunctionDefinitionEncoder, error) {
	ty, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch ty {
	case AggregateFunctionDefinitionTypeMin:
		return j.AsMin()
	case AggregateFunctionDefinitionTypeMax:
		return j.AsMax()
	case AggregateFunctionDefinitionTypeSum:
		return j.AsSum()
	case AggregateFunctionDefinitionTypeAverage:
		return j.AsAverage()
	case AggregateFunctionDefinitionTypeCustom:
		return j.AsCustom()
	default:
		return nil, fmt.Errorf("invalid AggregateFunctionDefinition type: %s", ty)
	}
}

// AggregateFunctionDefinitionMin represents a min aggregate function definition
type AggregateFunctionDefinitionMin struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewAggregateFunctionDefinitionMin creates an AggregateFunctionDefinitionMin instance.
func NewAggregateFunctionDefinitionMin() *AggregateFunctionDefinitionMin {
	return &AggregateFunctionDefinitionMin{
		Type: AggregateFunctionDefinitionTypeMin,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionMin) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type": j.Type,
	}
	return result
}

// AggregateFunctionDefinitionMax represents a max aggregate function definition
type AggregateFunctionDefinitionMax struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
}

// NewAggregateFunctionDefinitionMax creates an AggregateFunctionDefinitionMax instance.
func NewAggregateFunctionDefinitionMax() *AggregateFunctionDefinitionMax {
	return &AggregateFunctionDefinitionMax{
		Type: AggregateFunctionDefinitionTypeMax,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionMax) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type": j.Type,
	}
	return result
}

// AggregateFunctionDefinitionAverage represents an average aggregate function definition
type AggregateFunctionDefinitionAverage struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar type of the result of this function, which should have the type representation Float64
	ResultType string `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionAverage creates an AggregateFunctionDefinitionAverage instance.
func NewAggregateFunctionDefinitionAverage(resultType string) *AggregateFunctionDefinitionAverage {
	return &AggregateFunctionDefinitionAverage{
		Type:       AggregateFunctionDefinitionTypeAverage,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionAverage) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
}

// AggregateFunctionDefinitionSum represents a sum aggregate function definition
type AggregateFunctionDefinitionSum struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar type of the result of this function, which should have one of the type representations Int64 or Float64, depending on whether this function is defined on a scalar type with an integer or floating-point representation, respectively.
	ResultType string `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionSum creates an AggregateFunctionDefinitionSum instance.
func NewAggregateFunctionDefinitionSum(resultType string) *AggregateFunctionDefinitionSum {
	return &AggregateFunctionDefinitionSum{
		Type:       AggregateFunctionDefinitionTypeSum,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionSum) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
}

// AggregateFunctionDefinitionCustom represents a sum aggregate function definition
type AggregateFunctionDefinitionCustom struct {
	Type AggregateFunctionDefinitionType `json:"type" yaml:"type" mapstructure:"type"`
	// The scalar or object type of the result of this function.
	ResultType Type `json:"result_type" yaml:"result_type" mapstructure:"result_type"`
}

// NewAggregateFunctionDefinitionCustom creates an AggregateFunctionDefinitionCustom instance.
func NewAggregateFunctionDefinitionCustom(resultType Type) *AggregateFunctionDefinitionCustom {
	return &AggregateFunctionDefinitionCustom{
		Type:       AggregateFunctionDefinitionTypeCustom,
		ResultType: resultType,
	}
}

// Encode converts the instance to raw AggregateFunctionDefinition.
func (j AggregateFunctionDefinitionCustom) Encode() AggregateFunctionDefinition {
	result := AggregateFunctionDefinition{
		"type":        j.Type,
		"result_type": j.ResultType,
	}
	return result
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
	columnName := getStringValueByKey(j, "column_name")
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
