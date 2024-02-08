package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hasura/ndc-sdk-go/schema"
)

func (c *Connector) GetSchema(configuration *Configuration) (*schema.SchemaResponse, error) {
	return &schema.SchemaResponse{}, nil
}

func (c *Connector) Query(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (schema.QueryResponse, error) {
	requestVars := request.Variables
	if len(requestVars) == 0 {
		requestVars = []schema.QueryRequestVariablesElem{{}}
	}

	var rowSets []schema.RowSet

	for _, requestVar := range requestVars {
		rawResult, err := execQuery(ctx, configuration, state, request, requestVar)
		if err != nil {
			return nil, err
		}
		result, err := PruneFields(request.Query.Fields, rawResult)
		if err != nil {
			return nil, err
		}
		rowSets = append(rowSets, schema.RowSet{
			Aggregates: schema.RowSetAggregates{},
			Rows: []schema.Row{
				map[string]any{
					"__value": result,
				},
			},
		})
	}

	return rowSets, nil
}

func (c *Connector) Mutation(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	operationResults := make([]schema.MutationOperationResults, 0, len(request.Operations))

	for _, operation := range request.Operations {
		rawResult, err := execProcedure(ctx, configuration, state, request, &operation)
		if err != nil {
			return nil, err
		}
		result, err := PruneFields(operation.Fields, rawResult)
		if err != nil {
			return nil, err
		}
		operationResults = append(operationResults, schema.MutationOperationResults{
			AffectedRows: 1,
			Returning: []schema.Row{
				map[string]any{
					"__value": result,
				},
			},
		})
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func execQuery(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest, variables map[string]any) (any, error) {

	switch request.Collection {
	case "hello":
		args, err := ResolveArguments[HelloArguments](request.Arguments, variables)
		if err != nil {
			return nil, err
		}
		return Hello(ctx, state, args)
	default:
		return nil, fmt.Errorf("unsupported query: %s", request.Collection)
	}
}

func execProcedure(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest, operation *schema.MutationOperation) (any, error) {
	switch operation.Name {
	case "create_author":
		var args CreateAuthorArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		return CreateAuthor(ctx, state, &args)
	default:
		return nil, fmt.Errorf("unsupported procedure operation: %s", operation.Name)
	}
}

// PruneFields prune unnecessary fields from selection
func PruneFields(fields map[string]schema.Field, result any) (any, error) {
	if len(fields) == 0 {
		return result, nil
	}

	if result == nil {
		return nil, errors.New("expected object fields, got nil")
	}

	var outputMap map[string]any
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &outputMap,
		TagName: "json",
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(result); err != nil {
		return nil, err
	}

	output := make(map[string]any)
	for key, field := range fields {
		f, err := field.Interface()
		switch fi := f.(type) {
		case schema.ColumnField:
			if col, ok := outputMap[fi.Column]; ok {
				output[fi.Column] = col
			} else {
				output[fi.Column] = nil
			}
		case schema.RelationshipField:
			return nil, fmt.Errorf("unsupported relationship field,  %s", key)
		default:
			return nil, err
		}
	}

	return output, nil
}

// ResolveArguments resolve variables in arguments and map them to struct
func ResolveArguments[R any](arguments map[string]schema.Argument, variables map[string]any) (*R, error) {
	resolvedArgs, err := ResolveArgumentVariables(arguments, variables)
	if err != nil {
		return nil, err
	}

	var result R

	if err = mapstructure.Decode(resolvedArgs, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResolveArgumentVariables resolve variables in arguments if exist
func ResolveArgumentVariables(arguments map[string]schema.Argument, variables map[string]any) (map[string]any, error) {
	results := make(map[string]any)
	for key, arg := range arguments {
		switch arg.Type {
		case schema.ArgumentTypeLiteral:
			results[key] = arg.Value
		case schema.ArgumentTypeVariable:
			value, ok := variables[arg.Name]
			if !ok {
				return nil, fmt.Errorf("variable %s not found", arg.Name)
			}
			results[key] = value
		default:
			return nil, fmt.Errorf("unsupported argument type: %s", arg.Type)
		}
	}

	return results, nil
}
