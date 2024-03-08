package schema

import (
	"encoding/json"
	"errors"
	"fmt"
)

// UnmarshalJSONMap decodes FunctionInfo from a JSON map.
func (j *FunctionInfo) UnmarshalJSONMap(raw map[string]json.RawMessage) error {
	rawArguments, ok := raw["arguments"]
	var arguments FunctionInfoArguments
	if ok && !isNullJSON(rawArguments) {
		if err := json.Unmarshal(rawArguments, &arguments); err != nil {
			return fmt.Errorf("FunctionInfo.arguments: %s", err)
		}
	}

	rawName, ok := raw["name"]
	if !ok || isNullJSON(rawName) {
		return errors.New("FunctionInfo.name: required")
	}
	var name string
	if err := json.Unmarshal(rawName, &name); err != nil {
		return fmt.Errorf("FunctionInfo.name: %s", err)
	}
	if name == "" {
		return errors.New("FunctionInfo.name: required")
	}

	rawDescription, ok := raw["description"]
	var description *string
	if ok && !isNullJSON(rawDescription) {
		if err := json.Unmarshal(rawDescription, &description); err != nil {
			return fmt.Errorf("FunctionInfo.description: %s", err)
		}
	}

	rawResultType, ok := raw["result_type"]
	if !ok || isNullJSON(rawResultType) {
		return fmt.Errorf("FunctionInfo.result_type: required")
	}
	var resultType Type
	if err := json.Unmarshal(rawResultType, &resultType); err != nil {
		return fmt.Errorf("FunctionInfo.result_type: %s", err)
	}

	*j = FunctionInfo{
		Arguments:   arguments,
		Name:        name,
		ResultType:  resultType,
		Description: description,
	}
	return nil
}
