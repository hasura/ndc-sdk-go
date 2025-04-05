package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

type Connector struct{}

func (mc *Connector) ParseConfiguration(ctx context.Context, rawConfiguration string) (*Configuration, error) {
	return &Configuration{}, nil
}

func (mc *Connector) TryInitState(ctx context.Context, configuration *Configuration, metrics *connector.TelemetryState) (*State, error) {
	articles, err := fetchArticles()
	if err != nil {
		return nil, schema.InternalServerError("failed to read articles", map[string]any{
			"cause": err.Error(),
		})
	}

	authors, err := fetchAuthors()
	if err != nil {
		return nil, schema.InternalServerError("failed to read authors", map[string]any{
			"cause": err.Error(),
		})
	}

	institutions, err := fetchInstitutions()
	if err != nil {
		return nil, schema.InternalServerError("failed to read institutions", map[string]any{
			"cause": err.Error(),
		})
	}

	countries, err := fetchCountries()
	if err != nil {
		return nil, schema.InternalServerError("failed to read institutions", map[string]any{
			"cause": err.Error(),
		})
	}

	return &State{
		Authors:      authors,
		Articles:     articles,
		Institutions: institutions,
		Countries:    countries,
		Telemetry:    metrics,
	}, nil
}

func (mc *Connector) HealthCheck(ctx context.Context, configuration *Configuration, state *State) error {
	return nil
}

func (mc *Connector) GetCapabilities(configuration *Configuration) schema.CapabilitiesResponseMarshaler {
	return capabilities
}

func (mc *Connector) GetSchema(ctx context.Context, configuration *Configuration, state *State) (schema.SchemaResponseMarshaler, error) {
	return ndcSchema, nil
}

func (mc *Connector) QueryExplain(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	if !slices.ContainsFunc(ndcSchema.Functions, func(f schema.FunctionInfo) bool {
		return f.Name == request.Collection
	}) && !slices.ContainsFunc(ndcSchema.Collections, func(f schema.CollectionInfo) bool {
		return f.Name == request.Collection
	}) {
		return nil, schema.UnprocessableContentError("invalid query name: "+request.Collection, nil)
	}

	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *Connector) MutationExplain(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.ExplainResponse, error) {
	if len(request.Operations) == 0 {
		return nil, schema.UnprocessableContentError("require at least 1 operation", nil)
	}

	if !slices.ContainsFunc(ndcSchema.Procedures, func(f schema.ProcedureInfo) bool {
		return f.Name == request.Operations[0].Name
	}) {
		return nil, schema.UnprocessableContentError("invalid mutation name: "+request.Operations[0].Name, nil)
	}

	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *Connector) Query(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (schema.QueryResponse, error) {
	variableSets := request.Variables
	if variableSets == nil {
		variableSets = []schema.QueryRequestVariablesElem{make(map[string]any)}
	}

	rowSets := make([]schema.RowSet, 0, len(variableSets))

	for _, variables := range variableSets {
		handler, err := NewQueryHandler(state, request, variables)
		if err != nil {
			return nil, err
		}

		rowSet, err := handler.Execute()
		if err != nil {
			return nil, err
		}

		rowSets = append(rowSets, *rowSet)
	}

	return rowSets, nil
}

func (mc *Connector) Mutation(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	operationResults := []schema.MutationOperationResults{}

	for _, operation := range request.Operations {
		results, err := executeMutationOperation(ctx, state, request.CollectionRelationships, operation)
		if err != nil {
			return nil, err
		}

		operationResults = append(operationResults, results)
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func executeMutationOperation(ctx context.Context, state *State, collectionRelationship schema.MutationRequestCollectionRelationships, operation schema.MutationOperation) (schema.MutationOperationResults, error) {
	switch operation.Type {
	case schema.MutationOperationProcedure:
		return executeProcedure(ctx, state, collectionRelationship, operation)
	default:
		return nil, schema.NotSupportedError(fmt.Sprintf("unsupported operation type: %s", operation.Type), nil)
	}
}

type UpsertArticleArguments struct {
	Article Article `json:"article"`
}

func executeProcedure(_ context.Context, state *State, collectionRelationships schema.MutationRequestCollectionRelationships, operation schema.MutationOperation) (schema.MutationOperationResults, error) {
	switch operation.Name {
	case "upsert_article":
		return executeUpsertArticle(state, operation.Arguments, operation.Fields, collectionRelationships)
	case "delete_articles":
		return executeDeleteArticles(state, operation.Arguments, operation.Fields, collectionRelationships)
	default:
		return nil, schema.UnprocessableContentError("unknown procedure", nil)
	}
}

func executeUpsertArticle(
	state *State,
	arguments json.RawMessage,
	fields schema.NestedField,
	collectionRelationships map[string]schema.Relationship,
) (schema.MutationOperationResults, error) {
	var args UpsertArticleArguments

	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	var oldRow *Article

	latestArticle := state.GetLatestArticle()
	if args.Article.ID <= 0 {
		if latestArticle == nil {
			args.Article.ID = 1
		} else {
			args.Article.ID = latestArticle.ID + 1
		}

		state.Articles = append(state.Articles, args.Article)
	} else {
		for i, article := range state.Articles {
			if article.ID == args.Article.ID {
				oldRow = utils.ToPtr(article)
				state.Articles[i] = args.Article

				break
			}
		}

		if oldRow == nil {
			state.Articles = append(state.Articles, args.Article)
		}
	}

	qh := &QueryHandler{
		state:     state,
		variables: map[string]any{},
	}

	returning, err := qh.evalNestedField(collectionRelationships, oldRow, fields)
	if err != nil {
		return nil, err
	}

	return schema.NewProcedureResult(returning).Encode(), nil
}

func executeDeleteArticles(
	state *State,
	arguments json.RawMessage,
	fields schema.NestedField,
	collectionRelationships map[string]schema.Relationship,
) (schema.MutationOperationResults, error) {
	var argumentData struct {
		Where schema.Expression `json:"where"`
	}

	if err := json.Unmarshal(arguments, &argumentData); err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	if len(argumentData.Where) == 0 {
		return nil, schema.UnprocessableContentError("Expected argument 'where'", nil)
	}

	var removed []map[string]any

	qh := &QueryHandler{
		state:     state,
		variables: map[string]any{},
	}

	for _, article := range state.Articles {
		encodedArticle, err := utils.EncodeObject(article)
		if err != nil {
			return nil, schema.InternalServerError(err.Error(), nil)
		}

		ok, err := qh.evalExpression(nil, argumentData.Where, nil, encodedArticle)
		if err != nil {
			return nil, err
		}

		if ok {
			removed = append(removed, encodedArticle)
		}
	}

	returning, err := qh.evalNestedField(collectionRelationships, removed, fields)
	if err != nil {
		return nil, err
	}

	return schema.NewProcedureResult(returning).Encode(), nil
}
