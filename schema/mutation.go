package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UnmarshalJSONMap decodes ProcedureInfo from a JSON map.
func (j *ProcedureInfo) UnmarshalJSONMap(raw map[string]json.RawMessage) error {
	rawArguments, ok := raw["arguments"]
	var arguments ProcedureInfoArguments
	if ok && !isNullJSON(rawArguments) {
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("ProcedureInfo.arguments: %w", err)
		}
	}

	rawName, ok := raw["name"]
	if !ok || isNullJSON(rawName) {
		return errors.New("ProcedureInfo.name: required")
	}
	var name string
	if err := json.Unmarshal(rawName, &name); err != nil {
		return fmt.Errorf("ProcedureInfo.name: %w", err)
	}
	if name == "" {
		return errors.New("ProcedureInfo.name: required")
	}

	rawDescription, ok := raw["description"]
	var description *string
	if ok && !isNullJSON(rawDescription) {
		if err := json.Unmarshal(rawDescription, &description); err != nil {
			return fmt.Errorf("ProcedureInfo.description: %w", err)
		}
	}

	rawResultType, ok := raw["result_type"]
	if !ok || isNullJSON(rawResultType) {
		return errors.New("ProcedureInfo.result_type: required")
	}
	var resultType Type
	if err := json.Unmarshal(rawResultType, &resultType); err != nil {
		return fmt.Errorf("ProcedureInfo.result_type: %w", err)
	}

	*j = ProcedureInfo{
		Arguments:   arguments,
		Name:        name,
		ResultType:  resultType,
		Description: description,
	}
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
		return nil, fmt.Errorf("failed to parse MutationOperationType, expect one of %v, got %s", []MutationOperationType{MutationOperationProcedure}, input)
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
	Type MutationOperationType `json:"type" yaml:"type" mapstructure:"type"`
	// The name of the operation
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// Any named procedure arguments
	Arguments json.RawMessage `json:"arguments" yaml:"arguments" mapstructure:"arguments"`
	// The fields to return from the result, or null to return everything
	Fields NestedField `json:"fields,omitempty" yaml:"fields,omitempty" mapstructure:"fields"`
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
		return fmt.Errorf("field type in MutationOperation: %w", err)
	}

	value := MutationOperation{
		Type: operationType,
	}

	switch value.Type {
	case MutationOperationProcedure:
		name, err := unmarshalStringFromJsonMap(raw, "name", true)
		if err != nil {
			return fmt.Errorf("field name in MutationOperation: %w", err)
		}
		value.Name = name

		rawArguments, ok := raw["arguments"]
		if !ok {
			return errors.New("field arguments in MutationOperation: required")
		}

		value.Arguments = rawArguments

		rawFields, ok := raw["fields"]
		if ok && !isNullJSON(rawFields) {
			var fields NestedField
			if err = json.Unmarshal(rawFields, &fields); err != nil {
				return fmt.Errorf("field fields in MutationOperation: %w", err)
			}
			value.Fields = fields
		}
	}

	*j = value
	return nil
}

// MutationOperationResults represent the result of mutation operation
type MutationOperationResults map[string]any

// UnmarshalJSON implements json.Unmarshaler.
func (j *MutationOperationResults) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	rawType, ok := raw["type"]
	if !ok {
		return errors.New("field type in MutationOperationResults: required")
	}

	var ty MutationOperationType
	if err := json.Unmarshal(rawType, &ty); err != nil {
		return fmt.Errorf("field type in MutationOperationResults: %w", err)
	}

	result := map[string]any{
		"type": ty,
	}
	switch ty {
	case MutationOperationProcedure:
		rawResult, ok := raw["result"]
		if !ok {
			return errors.New("field result in MutationOperationResults is required for procedure type")
		}
		var procedureResult any
		if err := json.Unmarshal(rawResult, &procedureResult); err != nil {
			return fmt.Errorf("field result in MutationOperationResults: %w", err)
		}
		result["result"] = procedureResult
	}
	*j = result
	return nil
}

// Type gets the type enum of the current type
func (j MutationOperationResults) Type() (MutationOperationType, error) {
	t, ok := j["type"]
	if !ok {
		return MutationOperationType(""), errTypeRequired
	}
	switch raw := t.(type) {
	case string:
		v, err := ParseMutationOperationType(raw)
		if err != nil {
			return MutationOperationType(""), err
		}
		return *v, nil
	case MutationOperationType:
		return raw, nil
	default:
		return MutationOperationType(""), fmt.Errorf("invalid type: %+v", t)
	}
}

// AsProcedure tries to convert the instance to ProcedureResult type
func (j MutationOperationResults) AsProcedure() (*ProcedureResult, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}
	if t != MutationOperationProcedure {
		return nil, fmt.Errorf("invalid type; expected: %s, got: %s", MutationOperationProcedure, t)
	}

	rawResult, ok := j["result"]
	if !ok {
		return nil, errors.New("ProcedureResult.result is required")
	}

	return &ProcedureResult{
		Type:   t,
		Result: rawResult,
	}, nil
}

// Interface tries to convert the instance to MutationOperationResultsEncoder interface
func (j MutationOperationResults) Interface() MutationOperationResultsEncoder {
	result, _ := j.InterfaceT()
	return result
}

// InterfaceT tries to convert the instance to MutationOperationResultsEncoder interface safely with explicit error
func (j MutationOperationResults) InterfaceT() (MutationOperationResultsEncoder, error) {
	t, err := j.Type()
	if err != nil {
		return nil, err
	}

	switch t {
	case MutationOperationProcedure:
		return j.AsProcedure()
	default:
		return nil, fmt.Errorf("invalid type: %s", t)
	}
}

// MutationOperationResultsEncoder abstracts the serialization interface for MutationOperationResults
type MutationOperationResultsEncoder interface {
	Encode() MutationOperationResults
}

// ProcedureResult represent the result of a procedure mutation operation
type ProcedureResult struct {
	Type   MutationOperationType `json:"type" yaml:"type" mapstructure:"type"`
	Result any                   `json:"result" yaml:"result" mapstructure:"result"`
}

// Encode encodes the struct to MutationOperationResults
func (pr ProcedureResult) Encode() MutationOperationResults {
	return MutationOperationResults{
		"type":   pr.Type,
		"result": pr.Result,
	}
}

// NewProcedureResult creates a MutationProcedureResult instance
func NewProcedureResult(result any) *ProcedureResult {
	return &ProcedureResult{
		Type:   MutationOperationProcedure,
		Result: result,
	}
}
