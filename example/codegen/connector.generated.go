package main

import (
	"context"
	"encoding/json"
	"fmt"

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
		result, err := schema.EvalColumnFields(request.Query.Fields, rawResult)
		if err != nil {
			return nil, err
		}
		rowSets = append(rowSets, schema.RowSet{
			Aggregates: schema.RowSetAggregates{},
			Rows: []map[string]any{
				{
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
		result, err := schema.EvalColumnFields(nil, rawResult)
		if err != nil {
			return nil, err
		}
		operationResults = append(operationResults, schema.NewProcedureResult(result).Encode())
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func execQuery(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest, variables map[string]any) (any, error) {

	switch request.Collection {
	case "hello":
		args, err := schema.ResolveArguments[HelloArguments](request.Arguments, variables)
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
