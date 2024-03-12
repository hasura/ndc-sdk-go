package schema

import (
	"encoding/json"
	"fmt"
)

// SchemaResponseMarshaler abstract the response for /schema handler
type SchemaResponseMarshaler interface {
	MarshalSchemaJSON() ([]byte, error)
}

// MarshalSchemaJSON encodes the NDC schema response to JSON
func (j SchemaResponse) MarshalSchemaJSON() ([]byte, error) {
	return json.Marshal(j)
}

// RawSchemaResponse represents a NDC schema response with pre-encoded raw bytes
type RawSchemaResponse struct {
	data []byte
}

// NewRawSchemaResponse creates and validate a RawSchemaResponse instance
func NewRawSchemaResponse(data []byte) (*RawSchemaResponse, error) {
	// try to decode the response to ensure type-safe
	var resp SchemaResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to validate SchemaResponse from raw input: %s", err)
	}
	return NewRawSchemaResponseUnsafe(data), nil
}

// NewRawSchemaResponse creates a RawSchemaResponse instance from raw bytes without validation
func NewRawSchemaResponseUnsafe(data []byte) *RawSchemaResponse {
	return &RawSchemaResponse{
		data: data,
	}
}

// MarshalSchemaJSON encodes the NDC schema response to JSON
func (j RawSchemaResponse) MarshalSchemaJSON() ([]byte, error) {
	return j.data, nil
}

// CapabilitiesResponseMarshaler abstract the response for /capabilities handler
type CapabilitiesResponseMarshaler interface {
	MarshalCapabilitiesJSON() ([]byte, error)
}

// MarshalCapabilitiesJSON encodes the NDC schema response to JSON
func (j CapabilitiesResponse) MarshalCapabilitiesJSON() ([]byte, error) {
	return json.Marshal(j)
}

// RawCapabilitiesResponse represents a NDC capabilities response with pre-encoded raw bytes
type RawCapabilitiesResponse struct {
	data []byte
}

// NewRawCapabilitiesResponse creates and validate a RawSchemaResponse instance
func NewRawCapabilitiesResponse(data []byte) (*RawCapabilitiesResponse, error) {
	// try to decode the response to ensure type-safe
	var resp CapabilitiesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to validate CapabilitiesResponse from raw input: %s", err)
	}
	return NewRawCapabilitiesResponseUnsafe(data), nil
}

// NewRawCapabilitiesResponseUnsafe creates a RawSchemaResponse instance from raw bytes without validation
func NewRawCapabilitiesResponseUnsafe(data []byte) *RawCapabilitiesResponse {
	return &RawCapabilitiesResponse{
		data: data,
	}
}

// MarshalCapabilitiesJSON encodes the NDC schema response to JSON
func (j RawCapabilitiesResponse) MarshalCapabilitiesJSON() ([]byte, error) {
	return j.data, nil
}
