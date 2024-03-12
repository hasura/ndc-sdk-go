// Code generated by github.com/hasura/ndc-sdk-go/codegen, DO NOT EDIT.
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/hasura/ndc-codegen-example/functions"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

//go:embed schema.generated.json
var rawSchema []byte
var schemaResponse *schema.RawSchemaResponse

func init() {
	var err error

	schemaResponse, err = schema.NewRawSchemaResponse(rawSchema)
	if err != nil {
		panic(err)
	}
}

// GetSchema gets the connector's schema.
func (c *Connector) GetSchema(ctx context.Context, configuration *types.Configuration, _ *types.State) (schema.SchemaResponseMarshaler, error) {
	return schemaResponse, nil
}

// Query executes a query.
func (c *Connector) Query(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.QueryRequest) (schema.QueryResponse, error) {
	valueField, err := utils.EvalFunctionSelectionFieldValue(request)
	if err != nil {
		return nil, schema.BadRequestError(err.Error(), nil)
	}
	requestVars := request.Variables
	if len(requestVars) == 0 {
		requestVars = []schema.QueryRequestVariablesElem{make(schema.QueryRequestVariablesElem)}
	}

	rowSets := make([]schema.RowSet, len(requestVars))
	for i, requestVar := range requestVars {
		result, err := execQuery(ctx, state, request, valueField, requestVar)
		if err != nil {
			return nil, err
		}
		rowSets[i] = schema.RowSet{
			Aggregates: schema.RowSetAggregates{},
			Rows: []map[string]any{
				{
					"__value": result,
				},
			},
		}
	}

	return rowSets, nil
}

// Mutation executes a mutation.
func (c *Connector) Mutation(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	operationResults := make([]schema.MutationOperationResults, len(request.Operations))

	for i, operation := range request.Operations {
		switch operation.Type {
		case schema.MutationOperationProcedure:
			result, err := execProcedure(ctx, state, &operation)
			if err != nil {
				return nil, err
			}
			operationResults[i] = result
		default:
			return nil, schema.BadRequestError(fmt.Sprintf("invalid operation type: %s", operation.Type), nil)
		}
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func execQuery(ctx context.Context, state *types.State, request *schema.QueryRequest, queryFields schema.NestedField, variables map[string]any) (any, error) {

	switch request.Collection {
	case "getBool":
		if len(queryFields) > 0 {
			return nil, schema.BadRequestError("cannot evaluate selection fields for scalar", nil)
		}
		return functions.FunctionGetBool(ctx, state)
	case "getTypes":
		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		rawArgs, err := utils.ResolveArgumentVariables(request.Arguments, variables)
		if err != nil {
			return nil, schema.BadRequestError("failed to resolve argument variables", map[string]any{
				"cause": err.Error(),
			})
		}

		var args functions.GetTypesArguments
		if err = args.FromValue(rawArgs); err != nil {
			return nil, schema.BadRequestError("failed to resolve arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.FunctionGetTypes(ctx, state, &args)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		result, err := utils.EvalNestedColumnObject(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "hello":
		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.FunctionHello(ctx, state)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		result, err := utils.EvalNestedColumnObject(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil
	case "getArticles":
		selection, err := queryFields.AsArray()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be array", map[string]any{
				"cause": err.Error(),
			})
		}
		rawArgs, err := utils.ResolveArgumentVariables(request.Arguments, variables)
		if err != nil {
			return nil, schema.BadRequestError("failed to resolve argument variables", map[string]any{
				"cause": err.Error(),
			})
		}

		var args functions.GetArticlesArguments
		if err = args.FromValue(rawArgs); err != nil {
			return nil, schema.BadRequestError("failed to resolve arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.GetArticles(ctx, state, &args)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, schema.BadRequestError("expected not null result", nil)
		}

		result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, schema.BadRequestError(fmt.Sprintf("unsupported query: %s", request.Collection), nil)
	}
}

func execProcedure(ctx context.Context, state *types.State, operation *schema.MutationOperation) (schema.MutationOperationResults, error) {

	var result any
	switch operation.Name {
	case "create_article":
		selection, err := operation.Fields.AsObject()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args functions.CreateArticleArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.CreateArticle(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		result, err = utils.EvalNestedColumnObject(selection, rawResult)

		if err != nil {
			return nil, err
		}
	case "increase":
		if len(operation.Fields) > 0 {
			return nil, schema.BadRequestError("cannot evaluate selection fields for scalar", nil)
		}
		var err error
		result, err = functions.Increase(ctx, state)
		if err != nil {
			return nil, err
		}
	case "createAuthor":
		selection, err := operation.Fields.AsObject()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args functions.CreateAuthorArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.ProcedureCreateAuthor(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		result, err = utils.EvalNestedColumnObject(selection, rawResult)

		if err != nil {
			return nil, err
		}
	case "createAuthors":
		selection, err := operation.Fields.AsArray()
		if err != nil {
			return nil, schema.BadRequestError("the selection field type must be array", map[string]any{
				"cause": err.Error(),
			})
		}
		var args functions.CreateAuthorsArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := functions.ProcedureCreateAuthors(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, schema.BadRequestError("expected not null result", nil)
		}

		result, err = utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)

		if err != nil {
			return nil, err
		}

	default:
		return nil, schema.BadRequestError(fmt.Sprintf("unsupported procedure operation: %s", operation.Name), nil)
	}

	return schema.NewProcedureResult(result).Encode(), nil
}
