package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

type Configuration struct{}

type Article struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
}

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type InstitutionLocation struct {
	CountryID int      `json:"country_id"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Campuses  []string `json:"campuses"`
}

type InstitutionStaff struct {
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	Specialities  []string `json:"specialities"`
	BornCountryID int      `json:"born_country_id"`
}

type Institution struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Location    InstitutionLocation `json:"location"`
	Staff       []InstitutionStaff  `json:"staff"`
	Departments []string            `json:"departments"`
}

type State struct {
	Authors      []Author
	Articles     []Article
	Institutions []Institution
	Telemetry    *connector.TelemetryState
}

func (s *State) GetLatestArticle() *Article {
	if len(s.Articles) == 0 {
		return nil
	}

	var latestArticle Article
	for _, article := range s.Articles {
		if latestArticle.ID < article.ID {
			latestArticle = article
		}
	}

	return &latestArticle
}

type Connector struct{}

func (mc *Connector) ParseConfiguration(ctx context.Context, rawConfiguration string) (*Configuration, error) {
	return &Configuration{}, nil
}

func (mc *Connector) TryInitState(ctx context.Context, configuration *Configuration, metrics *connector.TelemetryState) (*State, error) {
	articles, err := readArticles()
	if err != nil {
		return nil, schema.InternalServerError("failed to read articles from csv", map[string]any{
			"cause": err.Error(),
		})
	}

	authors, err := readAuthors()
	if err != nil {
		return nil, schema.InternalServerError("failed to read authors from csv", map[string]any{
			"cause": err.Error(),
		})
	}

	institutions, err := readInstitutions()
	if err != nil {
		return nil, schema.InternalServerError("failed to read institutions from json", map[string]any{
			"cause": err.Error(),
		})
	}

	return &State{
		Authors:      authors,
		Articles:     articles,
		Institutions: institutions,
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
		rowSet, err := executeQueryWithVariables(request.Collection, request.Arguments, request.CollectionRelationships, &request.Query, variables, state)
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

	returning, err := evalNestedField(collectionRelationships, nil, state, oldRow, fields)
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
	for _, article := range state.Articles {
		encodedArticle, err := utils.EncodeObject(article)
		if err != nil {
			return nil, schema.InternalServerError(err.Error(), nil)
		}

		ok, err := evalExpression(nil, nil, state, argumentData.Where, encodedArticle, encodedArticle)
		if err != nil {
			return nil, err
		}
		if ok {
			removed = append(removed, encodedArticle)
		}
	}

	returning, err := evalNestedField(collectionRelationships, nil, state, removed, fields)
	if err != nil {
		return nil, err
	}

	return schema.NewProcedureResult(returning).Encode(), nil
}

func executeQueryWithVariables(
	collection string,
	arguments map[string]schema.Argument,
	collectionRelationships map[string]schema.Relationship,
	query *schema.Query,
	variables map[string]any,
	state *State,
) (*schema.RowSet, error) {
	argumentValues, err := utils.ResolveArgumentVariables(arguments, variables)
	if err != nil {
		return nil, err
	}

	coll, err := getCollectionByName(collection, argumentValues, state)
	if err != nil {
		return nil, err
	}
	return executeQuery(collectionRelationships, variables, state, query, nil, coll, false)
}

func evalAggregate(aggregate schema.Aggregate, paginated []map[string]any) (any, error) {
	switch agg := aggregate.Interface().(type) {
	case *schema.AggregateStarCount:
		return len(paginated), nil
	case *schema.AggregateColumnCount:
		var values []string
		for _, value := range paginated {
			v, ok := value[agg.Column]
			if !ok {
				return nil, schema.UnprocessableContentError("invalid column name: "+agg.Column, nil)
			}
			if v == nil {
				continue
			}
			values = append(values, fmt.Sprint(v))
		}
		if !agg.Distinct {
			return len(values), nil
		}
		distinctValue := make(map[string]bool)
		for _, v := range values {
			distinctValue[v] = true
		}
		return len(distinctValue), nil
	case *schema.AggregateSingleColumn:
		var values []any
		for _, value := range paginated {
			v, ok := value[agg.Column]
			if !ok {
				return nil, schema.UnprocessableContentError("invalid column name: "+agg.Column, nil)
			}
			if v == nil {
				continue
			}
			values = append(values, v)
		}
		return evalAggregateFunction(agg.Function, values)
	default:
		return nil, schema.UnprocessableContentError("invalid aggregate field", map[string]any{
			"value": aggregate,
		})
	}
}

func evalAggregateFunction(function string, values []any) (*int, error) {
	if len(values) == 0 {
		return nil, nil
	}

	var intValues []int
	for _, value := range values {
		switch v := value.(type) {
		case int:
			intValues = append(intValues, v)
		case int16:
			intValues = append(intValues, int(v))
		case int32:
			intValues = append(intValues, int(v))
		case int64:
			intValues = append(intValues, int(v))
		default:
			return nil, schema.UnprocessableContentError(fmt.Sprintf("%s: column is not an integer, got %+v", function, reflect.ValueOf(v).Kind()), nil)
		}
	}

	sort.Ints(intValues)

	switch function {
	case "min":
		return &intValues[0], nil
	case "max":
		return &intValues[len(intValues)-1], nil
	default:
		return nil, schema.UnprocessableContentError(function+": invalid aggregation function", nil)
	}
}

func executeQuery(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	query *schema.Query,
	root map[string]any,
	collection []map[string]any,
	skipMappingFields bool,
) (*schema.RowSet, error) {
	sorted, err := sortCollection(collectionRelationships, variables, state, collection, query.OrderBy)
	if err != nil {
		return nil, err
	}

	filtered := sorted
	if len(query.Predicate) > 0 {
		filtered = []map[string]any{}
		for _, item := range sorted {
			rootItem := root
			if rootItem == nil {
				rootItem = item
			}
			ok, err := evalExpression(collectionRelationships, variables, state, query.Predicate, rootItem, item)
			if err != nil {
				return nil, err
			}
			if ok {
				filtered = append(filtered, item)
			}
		}
	}

	paginated := paginate(filtered, query.Limit, query.Offset)
	aggregates := make(map[string]any)

	for aggKey, aggregate := range query.Aggregates {
		aggValue, err := evalAggregate(aggregate, paginated)
		if err != nil {
			return nil, err
		}
		aggregates[aggKey] = aggValue
	}

	rows := paginated
	if !skipMappingFields {
		rows = make([]map[string]any, 0)
		for _, item := range paginated {
			row, err := evalRow(query.Fields, collectionRelationships, variables, state, item)
			if err != nil {
				return nil, err
			}
			if row != nil {
				rows = append(rows, row)
			}
		}
	}

	result := &schema.RowSet{
		Aggregates: aggregates,
	}
	if len(rows) > 0 || len(result.Aggregates) == 0 {
		result.Rows = rows
	}
	return result, nil
}

func sortCollection(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	collection []map[string]any,
	orderBy *schema.OrderBy,
) ([]map[string]any, error) {
	if orderBy == nil || len(orderBy.Elements) == 0 {
		return collection, nil
	}

	var results []map[string]any
	for _, itemToInsert := range collection {
		if len(results) == 0 {
			results = append(results, itemToInsert)
			continue
		}
		inserted := false
		newResults := []map[string]any{}
		for _, other := range results {
			ordering, err := evalOrderBy(collectionRelationships, variables, state, orderBy, other, itemToInsert)
			if err != nil {
				return nil, err
			}
			if ordering > 0 {
				newResults = append(newResults, itemToInsert, other)
				inserted = true
			} else {
				newResults = append(newResults, other)
			}
		}

		if !inserted {
			newResults = append(newResults, itemToInsert)
		}

		results = newResults
	}
	return results, nil
}

func evalOrderBy(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	orderBy *schema.OrderBy,
	t1 map[string]any,
	t2 map[string]any,
) (int, error) {
	ordering := 0
	for _, orderElem := range orderBy.Elements {
		v1, err := evalOrderByElement(collectionRelationships, variables, state, orderElem, t1)
		if err != nil {
			return 0, err
		}
		v2, err := evalOrderByElement(collectionRelationships, variables, state, orderElem, t2)
		if err != nil {
			return 0, err
		}
		switch orderElem.OrderDirection {
		case schema.OrderDirectionAsc:
			// FIXME: compose ordering
			ordering, err = compare(v1, v2)
			if err != nil {
				return 0, err
			}
		case schema.OrderDirectionDesc:
			ordering, err = compare(v2, v1)
			if err != nil {
				return 0, err
			}
		}
		if ordering != 0 {
			return ordering, nil
		}
	}

	return ordering, nil
}

func evalOrderByElement(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	element schema.OrderByElement,
	item map[string]any,
) (any, error) {
	switch target := element.Target.Interface().(type) {
	case *schema.OrderByColumn:
		return evalOrderByColumn(collectionRelationships, variables, state, item, target.Path, target.Name)
	case *schema.OrderByAggregate:
		rows, err := evalPath(collectionRelationships, variables, state, target.Path, item)
		if err != nil {
			return nil, err
		}

		return evalAggregate(target.Aggregate, rows)
	default:
		return nil, schema.UnprocessableContentError("invalid order by field", map[string]any{
			"value": element.Target,
		})
	}
}

func evalOrderByColumn(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	item map[string]any,
	path []schema.PathElement,
	name string,
) (any, error) {
	rows, err := evalPath(collectionRelationships, variables, state, path, item)
	if err != nil {
		return nil, err
	}
	if len(rows) > 1 {
		return nil, schema.UnprocessableContentError("expected one path value only", nil)
	}
	if len(rows) == 0 || rows[0] == nil {
		return nil, nil
	}
	value, ok := rows[0][name]
	if !ok {
		return nil, schema.UnprocessableContentError("invalid column name: "+name, nil)
	}
	return value, nil
}

func evalInCollection(
	collectionRelationships map[string]schema.Relationship,
	item map[string]any,
	variables map[string]any,
	state *State,
	inCollection schema.ExistsInCollection,
) ([]map[string]any, error) {
	switch inCol := inCollection.Interface().(type) {
	case *schema.ExistsInCollectionRelated:
		relationship, ok := collectionRelationships[inCol.Relationship]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid in collection relationship: "+inCol.Relationship, nil)
		}
		source := []map[string]any{item}
		return evalPathElement(collectionRelationships, variables, state, &relationship, inCol.Arguments, source, nil)
	case *schema.ExistsInCollectionUnrelated:
		arguments := make(map[string]any)
		for key, relArg := range inCol.Arguments {
			argValue, err := evalRelationshipArgument(variables, item, relArg)
			if err != nil {
				return nil, err
			}
			arguments[key] = argValue
		}
		return getCollectionByName(inCol.Collection, arguments, state)
	default:
		return nil, schema.UnprocessableContentError("invalid in collection field", map[string]any{
			"value": inCollection,
		})
	}
}

func evalRow(fields map[string]schema.Field, collectionRelationships map[string]schema.Relationship, variables map[string]any, state *State, item map[string]any) (map[string]any, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	row := make(map[string]any)
	for fieldName, field := range fields {
		fieldValue, err := evalField(collectionRelationships, variables, state, field, item)
		if err != nil {
			return nil, err
		}
		row[fieldName] = fieldValue
	}

	return row, nil
}

func evalNestedField(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	value any,
	nestedField schema.NestedField,
) (any, error) {
	if utils.IsNil(value) {
		return value, nil
	}
	switch nf := nestedField.Interface().(type) {
	case *schema.NestedObject:
		fullRow, err := utils.EncodeObject(value)
		if err != nil {
			return nil, schema.UnprocessableContentError(fmt.Sprintf("expected object, got %s", reflect.ValueOf(value).Kind()), nil)
		}

		return evalRow(nf.Fields, collectionRelationships, variables, state, fullRow)
	case *schema.NestedArray:
		array, err := utils.EncodeObjects(value)
		if err != nil {
			return nil, err
		}

		result := []any{}
		for _, item := range array {
			val, err := evalNestedField(collectionRelationships, variables, state, item, nf.Fields)
			if err != nil {
				return nil, err
			}
			result = append(result, val)
		}
		return result, nil
	default:
		return nil, schema.UnprocessableContentError("invalid nested field", map[string]any{
			"value": nestedField,
		})
	}
}

func evalField(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	field schema.Field,
	row map[string]any,
) (any, error) {
	switch f := field.Interface().(type) {
	case *schema.ColumnField:
		value, ok := row[f.Column]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid column name: "+f.Column, nil)
		}
		if len(f.Fields) == 0 {
			return value, nil
		}
		return evalNestedField(collectionRelationships, variables, state, value, f.Fields)
	case *schema.RelationshipField:
		relationship, ok := collectionRelationships[f.Relationship]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid relationship name "+f.Relationship, nil)
		}

		collection, err := evalPathElement(collectionRelationships, variables, state, &relationship, f.Arguments, []map[string]any{row}, nil)
		if err != nil {
			return nil, err
		}

		return executeQuery(collectionRelationships, variables, state, &f.Query, nil, collection, false)

	default:
		return nil, schema.UnprocessableContentError("invalid field", map[string]any{
			"value": field,
		})
	}
}

func evalPathElement(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	relationship *schema.Relationship,
	arguments map[string]schema.RelationshipArgument,
	source []map[string]any,
	predicate schema.Expression,
) ([]map[string]any, error) {
	allArguments := make(map[string]any)
	var matchingRows []map[string]any

	// Note: Join strategy
	//
	// Rows can be related in two ways: 1) via a column mapping, and
	// 2) via collection arguments. Because collection arguments can be computed
	// using the columns on the source side of a relationship, in general
	// we need to compute the target collection once for each source row.
	// This join strategy can result in some target rows appearing in the
	// resulting row set more than once, if two source rows are both related
	// to the same target row.
	//
	// In practice, this is not an issue, either because a) the relationship
	// is computed in the course of evaluating a predicate, and all predicates are
	// implicitly or explicitly existentially quantified, or b) if the
	// relationship is computed in the course of evaluating an ordering, the path
	// should consist of all object relationships, and possibly terminated by a
	// single array relationship, so there should be no double counting.
	for _, srcRow := range source {
		for argName, arg := range relationship.Arguments {
			relValue, err := evalRelationshipArgument(variables, srcRow, arg)
			if err != nil {
				return nil, err
			}
			allArguments[argName] = relValue
		}
		for argName, arg := range arguments {
			if _, ok := allArguments[argName]; ok {
				return nil, schema.UnprocessableContentError("duplicate argument name: "+argName, nil)
			}
			relValue, err := evalRelationshipArgument(variables, srcRow, arg)
			if err != nil {
				return nil, err
			}
			allArguments[argName] = relValue
		}

		targetRows, err := getCollectionByName(relationship.TargetCollection, allArguments, state)
		if err != nil {
			return nil, err
		}

		for _, targetRow := range targetRows {
			ok, err := evalColumnMapping(relationship, srcRow, targetRow)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}

			if predicate != nil {
				ok, err := evalExpression(collectionRelationships, variables, state, predicate, targetRow, targetRow)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}
			}
			matchingRows = append(matchingRows, targetRow)
		}
	}

	return matchingRows, nil
}

func evalRelationshipArgument(variables map[string]any, row map[string]any, argument schema.RelationshipArgument) (any, error) {
	argT, err := argument.InterfaceT()
	switch arg := argT.(type) {
	case *schema.RelationshipArgumentColumn:
		value, ok := row[arg.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid column name: "+arg.Name, nil)
		}
		return value, nil
	case *schema.RelationshipArgumentLiteral:
		return arg.Value, nil
	case *schema.RelationshipArgumentVariable:
		variable, ok := variables[arg.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid variable name: "+arg.Name, nil)
		}
		return variable, nil
	default:
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}
}

func getCollectionByName(collectionName string, arguments map[string]any, state *State) ([]map[string]any, error) {
	var rows []map[string]any
	switch collectionName {
	// function
	case "latest_article_id":
		latestArticle := state.GetLatestArticle()
		var latestID *int
		if latestArticle != nil {
			latestID = &latestArticle.ID
		}
		return []map[string]any{
			{
				"__value": latestID,
			},
		}, nil
	case "latest_article":
		return []map[string]any{
			{
				"__value": state.GetLatestArticle(),
			},
		}, nil
		// collections
	case "articles":
		for _, item := range state.Articles {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		}
	case "authors":
		for _, item := range state.Authors {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		}
	case "institutions":
		for _, item := range state.Institutions {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		}
	case "articles_by_author":
		authorId, ok := arguments["author_id"]
		if !ok {
			return nil, schema.UnprocessableContentError("missing argument author_id", nil)
		}

		for _, row := range state.Articles {
			if strconv.Itoa(row.AuthorID) == fmt.Sprint(authorId) {
				r, err := utils.EncodeObject(row)
				if err != nil {
					return nil, err
				}
				rows = append(rows, r)
			}
		}
	default:
		return nil, schema.UnprocessableContentError("invalid collection name "+collectionName, nil)
	}

	return rows, nil
}

func evalComparisonValue(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	comparisonValue schema.ComparisonValue,
	root map[string]any,
	item map[string]any,
) ([]any, error) {
	switch compValue := comparisonValue.Interface().(type) {
	case *schema.ComparisonValueColumn:
		return evalComparisonTarget(collectionRelationships, variables, state, &compValue.Column, root, item)
	case *schema.ComparisonValueScalar:
		return []any{compValue.Value}, nil
	case *schema.ComparisonValueVariable:
		if len(variables) == 0 {
			return nil, schema.UnprocessableContentError("invalid variable name: "+compValue.Name, nil)
		}
		val, ok := variables[compValue.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid variable name: "+compValue.Name, nil)
		}
		return []any{val}, nil
	default:
		return nil, schema.UnprocessableContentError("invalid comparison value", map[string]any{
			"value": comparisonValue,
		})
	}
}

func evalComparisonTarget(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	target *schema.ComparisonTarget,
	root map[string]any,
	item map[string]any,
) ([]any, error) {
	switch target.Type {
	case schema.ComparisonTargetTypeColumn:
		rows, err := evalPath(collectionRelationships, variables, state, target.Path, item)
		if err != nil {
			return nil, err
		}
		var result []any
		for _, row := range rows {
			value, ok := row[target.Name]
			if !ok {
				return nil, schema.UnprocessableContentError("invalid comparison target column name: "+target.Name, nil)
			}
			result = append(result, value)
		}
		return result, nil
	case schema.ComparisonTargetTypeRootCollectionColumn:
		value, ok := root[target.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid comparison target column name: "+target.Name, nil)
		}
		return []any{value}, nil
	default:
		return nil, schema.UnprocessableContentError(fmt.Sprintf("invalid comparison target type: %s", target.Type), nil)
	}
}

func evalPath(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	path []schema.PathElement,
	item map[string]any,
) ([]map[string]any, error) {
	var err error
	result := []map[string]any{item}

	for _, pathElem := range path {
		relationshipName := pathElem.Relationship
		relationship, ok := collectionRelationships[relationshipName]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid relationship name in path: "+relationshipName, nil)
		}
		result, err = evalPathElement(collectionRelationships, variables, state, &relationship, pathElem.Arguments, result, pathElem.Predicate)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalExpression(
	collectionRelationships map[string]schema.Relationship,
	variables map[string]any,
	state *State,
	expr schema.Expression,
	root map[string]any,
	item map[string]any,
) (bool, error) {
	exprT, err := expr.InterfaceT()
	if err != nil {
		return false, err
	}

	switch expression := exprT.(type) {
	case *schema.ExpressionAnd:
		for _, exp := range expression.Expressions {
			ok, err := evalExpression(collectionRelationships, variables, state, exp, root, item)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil
	case *schema.ExpressionOr:
		for _, exp := range expression.Expressions {
			ok, err := evalExpression(collectionRelationships, variables, state, exp, root, item)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	case *schema.ExpressionNot:
		ok, err := evalExpression(collectionRelationships, variables, state, expression.Expression, root, item)
		if err != nil {
			return false, err
		}
		return !ok, nil
	case *schema.ExpressionUnaryComparisonOperator:
		switch expression.Operator {
		case schema.UnaryComparisonOperatorIsNull:
			values, err := evalComparisonTarget(collectionRelationships, variables, state, &expression.Column, root, item)
			if err != nil {
				return false, err
			}
			for _, val := range values {
				if val == nil {
					return true, nil
				}
			}
			return false, nil
		default:
			return false, schema.UnprocessableContentError(fmt.Sprintf("invalid unary comparison operator: %s", expression.Operator), nil)
		}
	case *schema.ExpressionBinaryComparisonOperator:
		switch expression.Operator {
		case "eq":
			leftValues, err := evalComparisonTarget(collectionRelationships, variables, state, &expression.Column, root, item)
			if err != nil {
				return false, err
			}
			rightValues, err := evalComparisonValue(collectionRelationships, variables, state, expression.Value, root, item)
			if err != nil {
				return false, err
			}

			for _, leftVal := range leftValues {
				for _, rightVal := range rightValues {
					// TODO: coalesce equality
					if leftVal == rightVal || fmt.Sprint(leftVal) == fmt.Sprint(rightVal) {
						return true, nil
					}
				}
			}
			return false, nil
		case "like":
			columnValues, err := evalComparisonTarget(collectionRelationships, variables, state, &expression.Column, root, item)
			if err != nil {
				return false, err
			}
			regexValues, err := evalComparisonValue(collectionRelationships, variables, state, expression.Value, root, item)
			if err != nil {
				return false, err
			}

			for _, columnValue := range columnValues {
				columnStr, ok := columnValue.(string)
				if !ok {
					return false, schema.UnprocessableContentError(fmt.Sprintf("value of column %s is not a string, got %+v", expression.Column, columnValue), nil)
				}
				for _, rawRegex := range regexValues {
					regexStr, ok := rawRegex.(string)
					if !ok {
						return false, schema.UnprocessableContentError(fmt.Sprintf("invalid regular expression, got %+v", rawRegex), nil)
					}

					regex, err := regexp.Compile(regexStr)
					if err != nil {
						return false, schema.UnprocessableContentError(fmt.Sprintf("invalid regular expression: %s", err), nil)
					}

					if regex.MatchString(columnStr) {
						return true, nil
					}
				}
			}

			return false, nil
		case "in":
			leftValues, err := evalComparisonTarget(collectionRelationships, variables, state, &expression.Column, root, item)
			if err != nil {
				return false, err
			}
			rightValueSets, err := evalComparisonValue(collectionRelationships, variables, state, expression.Value, root, item)
			if err != nil {
				return false, err
			}
			for _, rightValueSet := range rightValueSets {
				rightValues, ok := rightValueSet.([]any)
				if !ok {
					return false, schema.UnprocessableContentError(fmt.Sprintf("expected array, got %+v", rightValueSet), nil)
				}
				for _, leftVal := range leftValues {
					for _, rightVal := range rightValues {
						// TODO: coalesce equality
						if leftVal == rightVal || fmt.Sprint(leftVal) == fmt.Sprint(rightVal) {
							return true, nil
						}
					}
				}
			}
			return false, nil
		default:
			return false, schema.UnprocessableContentError("invalid comparison operator: "+expression.Operator, nil)
		}
	case *schema.ExpressionExists:
		query := &schema.Query{
			Predicate: expression.Predicate,
		}
		collection, err := evalInCollection(collectionRelationships, item, variables, state, expression.InCollection)
		if err != nil {
			return false, err
		}

		rowSet, err := executeQuery(collectionRelationships, variables, state, query, root, collection, true)
		if err != nil {
			return false, err
		}

		return len(rowSet.Rows) > 0, nil
	default:
		return false, schema.UnprocessableContentError("invalid expression", map[string]any{
			"value": expr,
		})
	}
}

func evalColumnMapping(relationship *schema.Relationship, srcRow map[string]any, target map[string]any) (bool, error) {
	for srcColumn, targetColumn := range relationship.ColumnMapping {
		srcValue, ok := srcRow[srcColumn]
		if !ok {
			return false, schema.UnprocessableContentError("source column does not exist: "+srcColumn, nil)
		}
		targetValue, ok := target[strings.Join(targetColumn, ".")]
		if !ok {
			return false, schema.UnprocessableContentError(fmt.Sprintf("target column does not exist: %v", targetColumn), nil)
		}
		if srcValue != targetValue {
			return false, nil
		}
	}
	return true, nil
}
