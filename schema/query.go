package schema

import (
	"encoding/json"
	"errors"
	"fmt"
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
