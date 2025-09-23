package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// SchemaResponseMarshaler abstract the response for /schema handler.
type SchemaResponseMarshaler interface {
	MarshalSchemaJSON() ([]byte, error)
}

// MarshalSchemaJSON encodes the NDC schema response to JSON.
func (j SchemaResponse) MarshalSchemaJSON() ([]byte, error) {
	return json.Marshal(j)
}

// RawSchemaResponse represents a NDC schema response with pre-encoded raw bytes.
type RawSchemaResponse struct {
	data []byte
}

// NewRawSchemaResponse creates and validate a RawSchemaResponse instance.
func NewRawSchemaResponse(data []byte) (*RawSchemaResponse, error) {
	// try to decode the response to ensure type-safe
	var resp SchemaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to validate SchemaResponse from raw input: %w", err)
	}

	return NewRawSchemaResponseUnsafe(data), nil
}

// NewRawSchemaResponseUnsafe creates a RawSchemaResponse instance from raw bytes without validation.
func NewRawSchemaResponseUnsafe(data []byte) *RawSchemaResponse {
	return &RawSchemaResponse{
		data: data,
	}
}

// MarshalSchemaJSON encodes the NDC schema response to JSON.
func (j RawSchemaResponse) MarshalSchemaJSON() ([]byte, error) {
	return j.data, nil
}

// CapabilitiesResponseMarshaler abstract the response for /capabilities handler.
type CapabilitiesResponseMarshaler interface {
	MarshalCapabilitiesJSON() ([]byte, error)
}

// MarshalCapabilitiesJSON encodes the NDC schema response to JSON.
func (j CapabilitiesResponse) MarshalCapabilitiesJSON() ([]byte, error) {
	return json.Marshal(j)
}

// RawCapabilitiesResponse represents a NDC capabilities response with pre-encoded raw bytes.
type RawCapabilitiesResponse struct {
	data []byte
}

// NewRawCapabilitiesResponse creates and validate a RawSchemaResponse instance.
func NewRawCapabilitiesResponse(data []byte) (*RawCapabilitiesResponse, error) {
	// try to decode the response to ensure type-safe
	var resp CapabilitiesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to validate CapabilitiesResponse from raw input: %w", err)
	}

	return NewRawCapabilitiesResponseUnsafe(data), nil
}

// NewRawCapabilitiesResponseUnsafe creates a RawSchemaResponse instance from raw bytes without validation.
func NewRawCapabilitiesResponseUnsafe(data []byte) *RawCapabilitiesResponse {
	return &RawCapabilitiesResponse{
		data: data,
	}
}

// MarshalCapabilitiesJSON encodes the NDC schema response to JSON.
func (j RawCapabilitiesResponse) MarshalCapabilitiesJSON() ([]byte, error) {
	return j.data, nil
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

// MarshalJSON implements json.Marshaler.
func (j ObjectType) MarshalJSON() ([]byte, error) {
	if len(j.Fields) == 0 {
		return nil, errors.New("fields in ObjectType must not be empty")
	}

	result := map[string]any{
		"description": j.Description,
		"fields":      j.Fields,
	}

	if j.ForeignKeys == nil {
		result["foreign_keys"] = ObjectTypeForeignKeys{}
	} else {
		result["foreign_keys"] = j.ForeignKeys
	}

	return json.Marshal(result)
}
