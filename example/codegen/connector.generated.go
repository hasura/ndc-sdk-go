// Code generated by github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go, DO NOT EDIT.
package main

import (
	"context"

	"fmt"
	"os"
	"strconv"
	"sync"

	"encoding/json"
	"github.com/hasura/ndc-codegen-example/functions"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
)

var loadGlobalEnvOnce sync.Once
var schemaResponse *schema.RawSchemaResponse
var connectorQueryHandlers = []ConnectorQueryHandler{functions.DataConnectorHandler{}}
var connectorMutationHandlers = []ConnectorMutationHandler{functions.DataConnectorHandler{}}

// ConnectorQueryHandler abstracts the connector query handler
type ConnectorQueryHandler interface {
	Query(ctx context.Context, state *types.State, request *schema.QueryRequest, arguments map[string]any) (*schema.RowSet, error)
}

// ConnectorMutationHandler abstracts the connector mutation handler
type ConnectorMutationHandler interface {
	Mutation(ctx context.Context, state *types.State, request *schema.MutationOperation) (schema.MutationOperationResults, error)
}

func init() {
	rawSchema, err := json.Marshal(GetConnectorSchema())
	if err != nil {
		panic(err)
	}
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
	if len(connectorQueryHandlers) == 0 {
		return nil, schema.UnprocessableContentError(fmt.Sprintf("unsupported query: %s", request.Collection), nil)
	}
	requestVars := request.Variables
	if len(requestVars) == 0 {
		requestVars = []schema.QueryRequestVariablesElem{make(schema.QueryRequestVariablesElem)}
	}

	if getGlobalEnvironments().queryConcurrencyLimit <= 1 || len(request.Variables) <= 1 {
		return c.execQuerySync(ctx, state, request, requestVars)
	}
	return c.execQueryAsync(ctx, state, request, requestVars)
}

func (c *Connector) execQuerySync(ctx context.Context, state *types.State, request *schema.QueryRequest, requestVars []schema.QueryRequestVariablesElem) (schema.QueryResponse, error) {
	rowSets := make([]schema.RowSet, len(requestVars))
	for i, requestVar := range requestVars {
		result, err := c.execQuery(ctx, state, request, requestVar, i)
		if err != nil {
			return nil, err
		}
		rowSets[i] = *result
	}

	return rowSets, nil
}

func (c *Connector) execQueryAsync(ctx context.Context, state *types.State, request *schema.QueryRequest, requestVars []schema.QueryRequestVariablesElem) (schema.QueryResponse, error) {
	rowSets := make([]schema.RowSet, len(requestVars))
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(getGlobalEnvironments().queryConcurrencyLimit)

	for i, requestVar := range requestVars {
		func(index int, vars schema.QueryRequestVariablesElem) {
			eg.Go(func() error {
				result, err := c.execQuery(ctx, state, request, vars, index)
				if err != nil {
					return err
				}
				rowSets[index] = *result
				return nil
			})
		}(i, requestVar)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return rowSets, nil
}

func (c *Connector) execQuery(ctx context.Context, state *types.State, request *schema.QueryRequest, variables map[string]any, index int) (*schema.RowSet, error) {
	ctx, span := state.Tracer.Start(ctx, fmt.Sprintf("Execute Function %d", index))
	defer span.End()

	rawArgs, err := utils.ResolveArgumentVariables(request.Arguments, variables)
	if err != nil {
		span.SetStatus(codes.Error, "failed to resolve argument variables")
		span.RecordError(err)
		return nil, schema.UnprocessableContentError("failed to resolve argument variables", map[string]any{
			"cause": err.Error(),
		})
	}

	for _, handler := range connectorQueryHandlers {
		result, err := handler.Query(ctx, state, request, rawArgs)
		if err == nil {
			return result, nil
		}

		if err != utils.ErrHandlerNotfound {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to execute function %d", index))
			span.RecordError(err)
			return nil, err
		}
	}

	errorMsg := fmt.Sprintf("unsupported query: %s", request.Collection)
	span.SetStatus(codes.Error, errorMsg)
	return nil, schema.UnprocessableContentError(errorMsg, nil)
}

// Mutation executes a mutation.
func (c *Connector) Mutation(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	if len(connectorMutationHandlers) == 0 {
		return nil, schema.UnprocessableContentError("unsupported mutation", nil)
	}

	if len(request.Operations) <= 1 || getGlobalEnvironments().mutationConcurrencyLimit <= 1 {
		return c.execMutationSync(ctx, state, request)
	}

	return c.execMutationAsync(ctx, state, request)
}

func (c *Connector) execMutationSync(ctx context.Context, state *types.State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	operationResults := make([]schema.MutationOperationResults, len(request.Operations))
	for i, operation := range request.Operations {
		result, err := c.execMutation(ctx, state, operation, i)
		if err != nil {
			return nil, err
		}
		operationResults[i] = result
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func (c *Connector) execMutationAsync(ctx context.Context, state *types.State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	operationResults := make([]schema.MutationOperationResults, len(request.Operations))
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(getGlobalEnvironments().mutationConcurrencyLimit)

	for i, operation := range request.Operations {
		func(index int, op schema.MutationOperation) {
			eg.Go(func() error {
				result, err := c.execMutation(ctx, state, op, index)
				if err != nil {
					return err
				}
				operationResults[index] = result
				return nil
			})
		}(i, operation)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func (c *Connector) execMutation(ctx context.Context, state *types.State, operation schema.MutationOperation, index int) (schema.MutationOperationResults, error) {
	ctx, span := state.Tracer.Start(ctx, fmt.Sprintf("Execute Procedure %d", index))
	defer span.End()

	span.SetAttributes(
		attribute.String("operation.type", string(operation.Type)),
		attribute.String("operation.name", string(operation.Name)),
	)

	switch operation.Type {
	case schema.MutationOperationProcedure:
		result, err := c.execProcedure(ctx, state, &operation)
		if err != nil {
			span.SetStatus(codes.Error, fmt.Sprintf("failed to execute procedure %d", index))
			span.RecordError(err)
			return nil, err
		}
		return result, nil
	default:
		errorMsg := fmt.Sprintf("invalid operation type: %s", operation.Type)
		span.SetStatus(codes.Error, errorMsg)
		return nil, schema.UnprocessableContentError(errorMsg, nil)
	}
}

func (c *Connector) execProcedure(ctx context.Context, state *types.State, operation *schema.MutationOperation) (schema.MutationOperationResults, error) {
	for _, handler := range connectorMutationHandlers {
		result, err := handler.Mutation(ctx, state, operation)
		if err == nil {
			return result, nil
		}
		if err != utils.ErrHandlerNotfound {
			return nil, err
		}
	}

	return nil, schema.UnprocessableContentError(fmt.Sprintf("unsupported procedure operation: %s", operation.Name), nil)
}

type globalEnvironments struct {
	queryConcurrencyLimit    int
	mutationConcurrencyLimit int
}

var _globalEnvironments = globalEnvironments{}

func initGlobalEnvironments() {
	rawQueryConcurrencyLimit := os.Getenv("QUERY_CONCURRENCY_LIMIT")
	if rawQueryConcurrencyLimit != "" {
		limit, err := strconv.ParseInt(rawQueryConcurrencyLimit, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("QUERY_CONCURRENCY_LIMIT: invalid integer <%s>", rawQueryConcurrencyLimit))
		}
		_globalEnvironments.queryConcurrencyLimit = int(limit)
	}

	rawMutationConcurrencyLimit := os.Getenv("MUTATION_CONCURRENCY_LIMIT")
	if rawMutationConcurrencyLimit != "" {
		limit, err := strconv.ParseInt(rawMutationConcurrencyLimit, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("MUTATION_CONCURRENCY_LIMIT: invalid integer <%s>", rawMutationConcurrencyLimit))
		}
		_globalEnvironments.mutationConcurrencyLimit = int(limit)
	}
}

func getGlobalEnvironments() *globalEnvironments {
	loadGlobalEnvOnce.Do(initGlobalEnvironments)
	return &_globalEnvironments
}
