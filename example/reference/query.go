package main

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

type QueryHandler struct {
	state     *State
	request   *schema.QueryRequest
	arguments map[string]any
	variables map[string]any
}

func NewQueryHandler(
	state *State,
	request *schema.QueryRequest,
	variables map[string]any,
) (*QueryHandler, error) {
	arguments, err := utils.ResolveArgumentVariables(request.Arguments, variables)
	if err != nil {
		return nil, err
	}

	return &QueryHandler{
		state:     state,
		request:   request,
		arguments: arguments,
		variables: variables,
	}, nil
}

func (qh *QueryHandler) Execute() (*schema.RowSet, error) {
	coll, err := qh.getCollectionByName(qh.request.Collection, qh.arguments)
	if err != nil {
		return nil, err
	}

	return qh.executeQuery(qh.request.CollectionRelationships, &qh.request.Query, nil, coll, false)
}

func (qh *QueryHandler) executeQuery(
	collectionRelationships map[string]schema.Relationship,
	query *schema.Query,
	root []map[string]any,
	collection []map[string]any,
	skipMappingFields bool,
) (*schema.RowSet, error) {
	sorted, err := qh.sortCollection(collectionRelationships, collection, query.OrderBy)
	if err != nil {
		return nil, err
	}

	filtered := sorted

	if len(query.Predicate) > 0 {
		filtered = []map[string]any{}

		for _, item := range sorted {
			ok, err := qh.evalExpression(
				collectionRelationships,
				query.Predicate,
				append(root, item),
				item,
			)
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
		aggValue, err := qh.evalAggregate(aggregate, paginated)
		if err != nil {
			return nil, err
		}

		aggregates[aggKey] = aggValue
	}

	rows := paginated

	if !skipMappingFields {
		rows = make([]map[string]any, 0)

		for _, item := range paginated {
			row, err := qh.evalRow(query.Fields, collectionRelationships, item)
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

	if query.Groups != nil {
		groups, err := qh.evalGroups(collectionRelationships, query.Groups, paginated)
		if err != nil {
			return nil, err
		}

		result.Groups = groups
	}

	if len(rows) > 0 || (len(result.Aggregates) == 0 && len(result.Groups) == 0) {
		result.Rows = rows
	}

	return result, nil
}

// Reference: https://github.com/hasura/ndc-spec/blob/50d6d618a1415cd0fc2d7710f42b1a04051529a8/ndc-reference/bin/reference/main.rs#L1123
func (qh *QueryHandler) getCollectionByName(
	collectionName string,
	arguments map[string]any,
) ([]map[string]any, error) {
	var rows []map[string]any

	switch collectionName {
	// collections
	case "articles":
		for _, item := range qh.state.Articles {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	case "authors":
		for _, item := range qh.state.Authors {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	case "institutions":
		for _, item := range qh.state.Institutions {
			row, err := utils.EncodeObject(item)
			if err != nil {
				return nil, err
			}

			rows = append(rows, row)
		}
	case "countries":
		for _, item := range qh.state.Countries {
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

		for _, row := range qh.state.Articles {
			if isEqual(row.AuthorID, authorId) {
				r, err := utils.EncodeObject(row)
				if err != nil {
					return nil, err
				}

				rows = append(rows, r)
			}
		}
	// function
	case "latest_article_id":
		var latestID *int

		latestArticle := qh.state.GetLatestArticle()

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
				"__value": qh.state.GetLatestArticle(),
			},
		}, nil
	default:
		return nil, schema.UnprocessableContentError("invalid collection name "+collectionName, nil)
	}

	return rows, nil
}

func (qh *QueryHandler) sortCollection(
	collectionRelationships map[string]schema.Relationship,
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
			if inserted {
				newResults = append(newResults, other)

				continue
			}

			ordering, err := qh.evalOrderBy(collectionRelationships, orderBy, other, itemToInsert)
			if err != nil {
				return nil, err
			}

			if ordering > 0 ||
				// workaround to make the expected result matches the Rust reference connector
				(ordering == 0 && len(orderBy.Elements) > 0 && orderBy.Elements[0].OrderDirection == schema.OrderDirectionDesc) {
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

// Reference: https://github.com/hasura/ndc-spec/blob/50d6d618a1415cd0fc2d7710f42b1a04051529a8/ndc-reference/bin/reference/main.rs#L1996
func (qh *QueryHandler) evalOrderBy(
	collectionRelationships map[string]schema.Relationship,
	orderBy *schema.OrderBy,
	t1 map[string]any,
	t2 map[string]any,
) (int, error) {
	ordering := 0

	for _, orderElem := range orderBy.Elements {
		v1, err := qh.evalOrderByElement(collectionRelationships, orderElem, t1)
		if err != nil {
			return 0, err
		}

		v2, err := qh.evalOrderByElement(collectionRelationships, orderElem, t2)
		if err != nil {
			return 0, err
		}

		switch orderElem.OrderDirection {
		case schema.OrderDirectionAsc:
			ordering, err = compare(v1, v2)
		case schema.OrderDirectionDesc:
			ordering, err = compare(v2, v1)
		}

		if err != nil {
			return 0, err
		}

		if ordering != 0 {
			return ordering, nil
		}
	}

	return ordering, nil
}

// Reference: https://github.com/hasura/ndc-spec/blob/50d6d618a1415cd0fc2d7710f42b1a04051529a8/ndc-reference/bin/reference/main.rs#L2046
func (qh *QueryHandler) evalOrderByElement(
	collectionRelationships map[string]schema.Relationship,
	element schema.OrderByElement,
	item map[string]any,
) (any, error) {
	switch target := element.Target.Interface().(type) {
	case *schema.OrderByColumn:
		return qh.evalColumnAtPath(collectionRelationships, item, target.Path, target.Name, target.Arguments, target.FieldPath)
	case *schema.OrderByAggregate:
		rows, err := qh.evalPath(collectionRelationships, target.Path, []map[string]any{item})
		if err != nil {
			return nil, err
		}

		return qh.evalAggregate(target.Aggregate, rows)
	default:
		return nil, schema.UnprocessableContentError("invalid order by field", map[string]any{
			"value": element.Target,
		})
	}
}

// Reference: https://github.com/hasura/ndc-spec/blob/50d6d618a1415cd0fc2d7710f42b1a04051529a8/ndc-reference/bin/reference/main.rs#L2150
func (qh *QueryHandler) evalPath(
	collectionRelationships map[string]schema.Relationship,
	path []schema.PathElement,
	item []map[string]any,
) ([]map[string]any, error) {
	var err error

	result := item

	for _, pathElem := range path {
		relationshipName := pathElem.Relationship

		relationship, ok := collectionRelationships[relationshipName]
		if !ok {
			return nil, schema.UnprocessableContentError(
				"invalid relationship name in path: "+relationshipName,
				nil,
			)
		}

		result, err = qh.evalPathElement(
			collectionRelationships,
			&relationship,
			pathElem.Arguments,
			result,
			pathElem.FieldPath,
			pathElem.Predicate,
		)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (qh *QueryHandler) evalPathElement(
	collectionRelationships map[string]schema.Relationship,
	relationship *schema.Relationship,
	arguments map[string]schema.RelationshipArgument,
	source []map[string]any,
	fieldPath []string,
	predicate schema.Expression,
) ([]map[string]any, error) {
	var matchingRows []map[string]any

	allArguments := make(map[string]any)

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
	for _, row := range source {
		srcRow, err := evalRowFieldPath(fieldPath, row)
		if err != nil {
			return nil, err
		}

		for argName, arg := range relationship.Arguments {
			relValue, err := qh.evalRelationshipArgument(srcRow, arg)
			if err != nil {
				return nil, err
			}

			allArguments[argName] = relValue
		}

		for argName, arg := range arguments {
			if _, ok := allArguments[argName]; ok {
				return nil, schema.UnprocessableContentError(
					"duplicate argument name: "+argName,
					nil,
				)
			}

			relValue, err := qh.evalRelationshipArgument(srcRow, arg)
			if err != nil {
				return nil, err
			}

			allArguments[argName] = relValue
		}

		targetRows, err := qh.getCollectionByName(relationship.TargetCollection, allArguments)
		if err != nil {
			return nil, err
		}

		for _, targetRow := range targetRows {
			ok, err := qh.evalColumnMapping(relationship, srcRow, targetRow)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}

			if predicate != nil {
				ok, err := qh.evalExpression(collectionRelationships, predicate, nil, targetRow)
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

func (qh *QueryHandler) evalComparisonValue(
	collectionRelationships map[string]schema.Relationship,
	comparisonValue schema.ComparisonValue,
	scopes []map[string]any,
	item map[string]any,
) ([]any, error) {
	comparisonValueT, err := comparisonValue.InterfaceT()
	if err != nil {
		return nil, err
	}

	switch compValue := comparisonValueT.(type) {
	case *schema.ComparisonValueColumn:
		scope := item

		if compValue.Scope != nil && *compValue.Scope > 0 {
			index := len(scopes) - 1 - int(*compValue.Scope)

			if index < 0 {
				return nil, fmt.Errorf("named scope %d is invalid", *compValue.Scope)
			}

			scope = scopes[index]
		}

		items, err := qh.evalPath(collectionRelationships, compValue.Path, []map[string]any{scope})
		if err != nil {
			return nil, err
		}

		var results []any

		for _, item := range items {
			result, err := qh.evalColumnFieldPath(item, compValue.Name, compValue.FieldPath, compValue.Arguments)
			if err != nil {
				return nil, err
			}

			results = append(results, result)
		}

		return results, nil
	case *schema.ComparisonValueScalar:
		return []any{compValue.Value}, nil
	case *schema.ComparisonValueVariable:
		if len(qh.variables) == 0 {
			return nil, schema.UnprocessableContentError("invalid variable name: "+compValue.Name, nil)
		}

		val, ok := qh.variables[compValue.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid variable name: "+compValue.Name, nil)
		}

		return []any{val}, nil
	default:
		return nil, schema.UnprocessableContentError("invalid comparison value", map[string]any{
			"value": comparisonValueT,
		})
	}
}

func (qh *QueryHandler) evalExpression(
	collectionRelationships map[string]schema.Relationship,
	expr schema.Expression,
	scopes []map[string]any,
	item map[string]any,
) (bool, error) {
	exprT, err := expr.InterfaceT()
	if err != nil {
		return false, err
	}

	switch expression := exprT.(type) {
	case *schema.ExpressionAnd:
		for _, exp := range expression.Expressions {
			ok, err := qh.evalExpression(collectionRelationships, exp, scopes, item)
			if err != nil || !ok {
				return false, err
			}
		}

		return true, nil
	case *schema.ExpressionOr:
		for _, exp := range expression.Expressions {
			ok, err := qh.evalExpression(collectionRelationships, exp, scopes, item)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}
		}

		return false, nil
	case *schema.ExpressionNot:
		ok, err := qh.evalExpression(collectionRelationships, expression.Expression, scopes, item)
		if err != nil {
			return false, err
		}

		return !ok, nil
	case *schema.ExpressionUnaryComparisonOperator:
		switch expression.Operator {
		case schema.UnaryComparisonOperatorIsNull:
			value, err := qh.evalComparisonTarget(collectionRelationships, expression.Column, item)
			if err != nil {
				return false, err
			}

			return utils.IsNil(value), nil
		default:
			return false, schema.UnprocessableContentError(fmt.Sprintf("invalid unary comparison operator: %s", expression.Operator), nil)
		}
	case *schema.ExpressionBinaryComparisonOperator:
		leftValues, err := qh.evalComparisonTarget(collectionRelationships, expression.Column, item)
		if err != nil {
			return false, err
		}

		rightValueSets, err := qh.evalComparisonValue(collectionRelationships, expression.Value, scopes, item)
		if err != nil {
			return false, err
		}

		return evalComparisonOperator(expression.Operator, leftValues, rightValueSets)
	case *schema.ExpressionArrayComparison:
		leftVal, err := qh.evalComparisonTarget(collectionRelationships, expression.Column, item)
		if err != nil {
			return false, err
		}

		return qh.evalArrayComparison(collectionRelationships, leftVal, expression.Comparison, scopes, item)
	case *schema.ExpressionExists:
		query := &schema.Query{
			Predicate: expression.Predicate,
		}

		collection, err := qh.evalInCollection(collectionRelationships, item, expression.InCollection)
		if err != nil {
			return false, err
		}

		rowSet, err := qh.executeQuery(collectionRelationships, query, scopes, collection, true)
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

func (qh *QueryHandler) evalRelationshipArgument(
	row map[string]any,
	argument schema.RelationshipArgument,
) (any, error) {
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
		variable, ok := qh.variables[arg.Name]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid variable name: "+arg.Name, nil)
		}

		return variable, nil
	default:
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}
}

func (qh *QueryHandler) evalInCollection(
	collectionRelationships map[string]schema.Relationship,
	item map[string]any,
	inCollection schema.ExistsInCollection,
) ([]map[string]any, error) {
	switch inCol := inCollection.Interface().(type) {
	case *schema.ExistsInCollectionRelated:
		relationship, ok := collectionRelationships[inCol.Relationship]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid in collection relationship: "+inCol.Relationship, nil)
		}

		source := []map[string]any{item}

		return qh.evalPathElement(collectionRelationships, &relationship, inCol.Arguments, source, inCol.FieldPath, nil)
	case *schema.ExistsInCollectionUnrelated:
		arguments := make(map[string]any)

		for key, relArg := range inCol.Arguments {
			argValue, err := qh.evalRelationshipArgument(item, relArg)
			if err != nil {
				return nil, err
			}

			arguments[key] = argValue
		}

		return qh.getCollectionByName(inCol.Collection, arguments)
	case *schema.ExistsInCollectionNestedCollection:
		value, err := qh.evalColumnFieldPath(item, inCol.ColumnName, inCol.FieldPath, inCol.Arguments)
		if err != nil {
			return nil, err
		}

		result, err := utils.EncodeObjects(value)
		if err != nil {
			return nil, fmt.Errorf("nested collection must be an array of objects, got %v", value)
		}

		return result, nil
	case *schema.ExistsInCollectionNestedScalarCollection:
		value, err := qh.evalColumnFieldPath(item, inCol.ColumnName, inCol.FieldPath, inCol.Arguments)
		if err != nil {
			return nil, err
		}

		valueArray, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("nested collection must be an array, got %v", value)
		}

		wrappedArrayValues := make([]map[string]any, len(valueArray))

		for i, item := range valueArray {
			wrappedArrayValues[i] = map[string]any{
				"__value": item,
			}
		}

		return wrappedArrayValues, nil
	default:
		return nil, schema.UnprocessableContentError("invalid in collection field", map[string]any{
			"value": inCollection,
		})
	}
}

func (qh *QueryHandler) evalRow(
	fields map[string]schema.Field,
	collectionRelationships map[string]schema.Relationship,
	item map[string]any,
) (map[string]any, error) {
	if len(fields) == 0 {
		return nil, nil
	}

	row := make(map[string]any)

	for fieldName, field := range fields {
		fieldValue, err := qh.evalField(collectionRelationships, field, item)
		if err != nil {
			return nil, err
		}

		row[fieldName] = fieldValue
	}

	return row, nil
}

func (qh *QueryHandler) evalNestedField(
	collectionRelationships map[string]schema.Relationship,
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

		return qh.evalRow(nf.Fields, collectionRelationships, fullRow)
	case *schema.NestedArray:
		array, err := utils.EncodeObjects(value)
		if err != nil {
			return nil, err
		}

		result := []any{}

		for _, item := range array {
			val, err := qh.evalNestedField(collectionRelationships, item, nf.Fields)
			if err != nil {
				return nil, err
			}

			result = append(result, val)
		}

		return result, nil
	case *schema.NestedCollection:
		collection, err := utils.EncodeObjects(value)
		if err != nil {
			return nil, err
		}

		rowSet, err := qh.executeQuery(collectionRelationships, &nf.Query, nil, collection, false)
		if err != nil {
			return nil, err
		}

		return rowSet, nil
	default:
		return nil, schema.UnprocessableContentError("invalid nested field", map[string]any{
			"value": nestedField,
		})
	}
}

func (qh *QueryHandler) evalField(
	collectionRelationships map[string]schema.Relationship,
	field schema.Field,
	row map[string]any,
) (any, error) {
	switch f := field.Interface().(type) {
	case *schema.ColumnField:
		value, err := qh.evalColumn(row, f.Column, f.Arguments)
		if err != nil {
			return nil, schema.UnprocessableContentError(err.Error(), nil)
		}

		if len(f.Fields) == 0 {
			return value, nil
		}

		return qh.evalNestedField(collectionRelationships, value, f.Fields)
	case *schema.RelationshipField:
		relationship, ok := collectionRelationships[f.Relationship]
		if !ok {
			return nil, schema.UnprocessableContentError("invalid relationship name "+f.Relationship, nil)
		}

		collection, err := qh.evalPathElement(collectionRelationships, &relationship, f.Arguments, []map[string]any{row}, nil, nil)
		if err != nil {
			return nil, err
		}

		return qh.executeQuery(collectionRelationships, &f.Query, nil, collection, false)
	default:
		return nil, schema.UnprocessableContentError("invalid field", map[string]any{
			"value": field,
		})
	}
}

func (qh *QueryHandler) evalDimensions(
	collectionRelationships map[string]schema.Relationship,
	row map[string]any,
	dimensions []schema.Dimension,
) ([]any, error) {
	results := make([]any, len(dimensions))

	for i, dimension := range dimensions {
		result, err := qh.evalDimension(collectionRelationships, row, dimension)
		if err != nil {
			return nil, err
		}

		results[i] = result
	}

	return results, nil
}

func (qh *QueryHandler) evalDimension(
	collectionRelationships map[string]schema.Relationship,
	row map[string]any,
	dimension schema.Dimension,
) (any, error) {
	dim := dimension.Interface()
	if dim == nil {
		return nil, errors.New(" nil query handler")
	}

	switch d := dim.(type) {
	case *schema.DimensionColumn:
		value, err := qh.evalColumnAtPath(collectionRelationships, row, d.Path, d.ColumnName, d.Arguments, d.FieldPath)
		if err != nil {
			return nil, err
		}

		return evalExtraction(d.Extraction, value)
	default:
		return nil, fmt.Errorf("unsupported dimension type %v", dim)
	}
}

// Reference: https://github.com/hasura/ndc-spec/blob/50d6d618a1415cd0fc2d7710f42b1a04051529a8/ndc-reference/bin/reference/main.rs#L2116
func (qh *QueryHandler) evalColumnAtPath(
	collectionRelationships map[string]schema.Relationship,
	row map[string]any,
	pathElems []schema.PathElement,
	name string,
	rawArguments map[string]schema.Argument,
	fieldPath []string,
) (any, error) {
	rows, err := qh.evalPath(collectionRelationships, pathElems, []map[string]any{row})
	if err != nil {
		return nil, err
	}

	if len(rows) > 1 {
		return nil, schema.UnprocessableContentError(
			"path elements used in sorting and grouping cannot yield multiple rows",
			nil,
		)
	}

	if len(rows) == 0 || rows[0] == nil {
		return nil, nil
	}

	return qh.evalColumnFieldPath(rows[0], name, fieldPath, rawArguments)
}

func (qh *QueryHandler) evalGroups(
	collectionRelationships map[string]schema.Relationship,
	grouping *schema.Grouping,
	paginated []map[string]any,
) ([]schema.Group, error) {
	chunks := make([]Chunk, 0)

L:
	for _, row := range paginated {
		dimensions, err := qh.evalDimensions(collectionRelationships, row, grouping.Dimensions)
		if err != nil {
			return nil, err
		}

		for i, chunk := range chunks {
			if reflect.DeepEqual(chunk.Dimensions, dimensions) {
				chunk.Rows = append(chunk.Rows, row)
				chunks[i] = chunk

				continue L
			}
		}

		chunks = append(chunks, Chunk{
			Dimensions: dimensions,
			Rows:       []map[string]any{row},
		})
	}

	sorted, err := qh.groupSort(chunks, grouping.OrderBy)
	if err != nil {
		return nil, err
	}

	groups := []schema.Group{}

	for _, chunk := range sorted {
		isValid, err := qh.evalGroupExpression(grouping.Predicate, chunk.Rows)
		if err != nil {
			return nil, err
		}

		if !isValid {
			continue
		}

		group := schema.Group{
			Dimensions: chunk.Dimensions,
			Aggregates: make(schema.GroupAggregates),
		}

		for aggName, rawAgg := range grouping.Aggregates {
			agg, err := qh.evalAggregate(rawAgg, chunk.Rows)
			if err != nil {
				return nil, err
			}

			group.Aggregates[aggName] = agg
		}

		groups = append(groups, group)
	}

	results := paginate(groups, grouping.Limit, grouping.Offset)

	return results, nil
}

func (qh *QueryHandler) groupSort(groups []Chunk, orderBy *schema.GroupOrderBy) ([]Chunk, error) {
	if orderBy == nil {
		return groups, nil
	}

	results := []Chunk{}

	for _, itemToInsert := range groups {
		index := 0

		for _, other := range results {
			cmpResult, err := qh.evalGroupOrderBy(orderBy, other, itemToInsert)
			if err != nil {
				return nil, schema.UnprocessableContentError(err.Error(), nil)
			}

			if cmpResult < 0 {
				break
			}

			index++
		}

		if index >= len(results) {
			results = append(results, itemToInsert)

			continue
		}

		results = append(results[:index], append([]Chunk{itemToInsert}, results[index:]...)...)
	}

	return results, nil
}

func (qh *QueryHandler) evalGroupOrderBy(
	orderBy *schema.GroupOrderBy,
	t1 Chunk,
	t2 Chunk,
) (int, error) {
	for _, elem := range orderBy.Elements {
		v1, err := qh.evalGroupOrderByElement(elem, t1)
		if err != nil {
			return 0, err
		}

		v2, err := qh.evalGroupOrderByElement(elem, t2)
		if err != nil {
			return 0, err
		}

		comparedResult, err := compare(v1, v2)
		if err != nil {
			return 0, err
		}

		if comparedResult == 0 {
			continue
		}

		if elem.OrderDirection == schema.OrderDirectionDesc {
			return comparedResult, nil
		}

		return -1 * comparedResult, nil
	}

	return 0, nil
}

func (qh *QueryHandler) evalGroupOrderByElement(
	element schema.GroupOrderByElement,
	group Chunk,
) (any, error) {
	rawTarget := element.Target.Interface()
	if rawTarget == nil {
		return nil, errors.New("nil target")
	}

	switch target := rawTarget.(type) {
	case *schema.GroupOrderByTargetDimension:
		if len(group.Dimensions) <= int(target.Index)-1 {
			return nil, fmt.Errorf("dimension index %d out of range", target.Index)
		}

		return group.Dimensions[target.Index], nil
	case *schema.GroupOrderByTargetAggregate:
		return qh.evalAggregate(target.Aggregate, group.Rows)
	default:
		return nil, fmt.Errorf("unsupported target type %s", rawTarget.Type())
	}
}

func (qh *QueryHandler) evalAggregate(
	aggregate schema.Aggregate,
	paginated []map[string]any,
) (any, error) {
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

func (qh *QueryHandler) evalGroupExpression(
	expression *schema.GroupExpression,
	rows []map[string]any,
) (bool, error) {
	if expression == nil {
		return true, nil
	}

	rawExpr := expression.Interface()
	if rawExpr == nil {
		return true, nil
	}

	switch expr := rawExpr.(type) {
	case *schema.GroupExpressionAnd:
		for _, e := range expr.Expressions {
			ok, err := qh.evalGroupExpression(&e, rows)
			if err != nil {
				return false, err
			}

			if !ok {
				return false, nil
			}
		}

		return true, nil
	case *schema.GroupExpressionOr:
		for _, e := range expr.Expressions {
			ok, err := qh.evalGroupExpression(&e, rows)
			if err != nil {
				return false, err
			}

			if ok {
				return true, nil
			}
		}

		return false, nil
	case *schema.GroupExpressionNot:
		ok, err := qh.evalGroupExpression(&expr.Expression, rows)
		if err != nil {
			return false, err
		}

		return !ok, nil
	case *schema.GroupExpressionBinaryComparisonOperator:
		leftVal, err := qh.evalGroupComparisonTarget(expr.Target, rows)
		if err != nil {
			return false, err
		}

		rightVal, err := qh.evalGroupComparisonValue(expr.Value)
		if err != nil {
			return false, err
		}

		return evalComparisonOperator(expr.Operator, leftVal, []any{rightVal})
	case *schema.GroupExpressionUnaryComparisonOperator:
		switch expr.Operator {
		case schema.UnaryComparisonOperatorIsNull:
			value, err := qh.evalGroupComparisonTarget(expr.Target, rows)
			if err != nil {
				return false, err
			}

			return utils.IsNil(value), nil
		default:
			return false, fmt.Errorf("unsupported unary comparison operator %s", expr.Operator)
		}

	default:
		return false, fmt.Errorf("unsupported group expression type %s", rawExpr.Type())
	}
}

func (qh *QueryHandler) evalGroupComparisonValue(
	comparisonValue schema.GroupComparisonValue,
) (any, error) {
	cmpValueT := comparisonValue.Interface()
	if cmpValueT == nil {
		return nil, errors.New("nil group comparison value")
	}

	switch cmpValue := cmpValueT.(type) {
	case *schema.GroupComparisonValueScalar:
		return cmpValue.Value, nil
	case *schema.GroupComparisonValueVariable:
		v, ok := qh.variables[cmpValue.Name]
		if !ok {
			return nil, fmt.Errorf("invalid comparison value, variable name '%s' does not exist", cmpValue.Name)
		}

		return v, nil
	default:
		return nil, fmt.Errorf("unsupported comparison value %s", cmpValueT.Type())
	}
}

func (qh *QueryHandler) evalGroupComparisonTarget(
	target schema.GroupComparisonTarget,
	rows []map[string]any,
) (any, error) {
	targetT := target.Interface()
	if targetT == nil {
		return nil, errors.New("nil group comparison target")
	}

	switch t := targetT.(type) {
	case *schema.GroupComparisonTargetAggregate:
		return qh.evalAggregate(t.Aggregate, rows)
	default:
		return nil, fmt.Errorf("unsupported group comparison target %s", targetT.Type())
	}
}

func (qh *QueryHandler) evalArrayComparison(
	collectionRelationships map[string]schema.Relationship,
	leftVal any,
	comparison schema.ArrayComparison,
	scopes []map[string]any,
	item map[string]any,
) (bool, error) {
	leftValArray, ok := leftVal.([]any)
	if !ok {
		return false, errors.New("column used in array comparison is not an array")
	}

	comparisonT, err := comparison.InterfaceT()
	if err != nil {
		return false, err
	}

	switch comp := comparisonT.(type) {
	case *schema.ArrayComparisonContains:
		rightVals, err := qh.evalComparisonValue(collectionRelationships, comp.Value, scopes, item)
		if err != nil {
			return false, err
		}

		for _, rightVal := range rightVals {
			if slices.Contains(leftValArray, rightVal) {
				return true, nil
			}
		}

		return false, nil
	case *schema.ArrayComparisonIsEmpty:
		return len(leftValArray) == 0, nil
	default:
		return false, fmt.Errorf("unsupported array comparison %s", comparisonT.Type())
	}
}

func (qh *QueryHandler) evalColumn(
	row map[string]any,
	columnName string,
	arguments map[string]schema.Argument,
) (any, error) {
	column, ok := row[columnName]
	if !ok {
		return nil, fmt.Errorf("invalid column name %s", columnName)
	}

	if utils.IsNil(column) {
		return nil, nil
	}

	values, ok := column.([]any)
	if !ok {
		return column, nil
	}

	if len(values) == 0 {
		return values, nil
	}

	limitArgument, ok := arguments["limit"]
	if !ok {
		return nil, schema.UnprocessableContentError(
			"Expected argument 'limit' in column "+columnName,
			nil,
		)
	}

	rawLimitArg, err := utils.ResolveArgument(limitArgument, qh.variables)
	if err != nil {
		return nil, err
	}

	limitArg, err := utils.DecodeNullableInt[int](rawLimitArg)
	if err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	if limitArg != nil && *limitArg < len(values) {
		values = values[:*limitArg]
	}

	return values, nil
}

func (qh *QueryHandler) evalColumnFieldPath(
	row map[string]any,
	columnName string,
	fieldPath []string,
	arguments map[string]schema.Argument,
) (any, error) {
	columnValue, err := qh.evalColumn(row, columnName, arguments)
	if err != nil {
		return nil, err
	}

	return evalFieldPath(fieldPath, columnValue)
}

func (qh *QueryHandler) evalComparisonTarget(
	collectionRelationships map[string]schema.Relationship,
	target schema.ComparisonTarget,
	item map[string]any,
) (any, error) {
	targetT, err := target.InterfaceT()
	if err != nil {
		return nil, err
	}

	switch t := targetT.(type) {
	case *schema.ComparisonTargetColumn:
		return qh.evalColumnFieldPath(item, t.Name, t.FieldPath, t.Arguments)
	case *schema.ComparisonTargetAggregate:
		rows, err := qh.evalPath(collectionRelationships, t.Path, []map[string]any{item})
		if err != nil {
			return nil, err
		}

		return qh.evalAggregate(t.Aggregate, rows)
	default:
		return nil, schema.UnprocessableContentError(fmt.Sprintf("invalid comparison target type: %s", t.Type()), nil)
	}
}

func (qh *QueryHandler) evalColumnMapping(
	relationship *schema.Relationship,
	srcRow map[string]any,
	tgtRow map[string]any,
) (bool, error) {
	for srcColumn, targetColumnPath := range relationship.ColumnMapping {
		srcValue, err := qh.evalColumn(srcRow, srcColumn, nil)
		if err != nil {
			return false, err
		}

		var targetColumn string

		targetRow := tgtRow

		switch len(targetColumnPath) {
		case 0:
			return false, errors.New(
				"relationship column mapping target column path were empty for column " + srcColumn,
			)
		case 1:
			targetColumn = targetColumnPath[0]
		default:
			targetColumn = targetColumnPath[len(targetColumnPath)-1]

			targetRow, err = evalRowFieldPath(targetColumnPath[:len(targetColumnPath)-1], tgtRow)
			if err != nil {
				return false, err
			}
		}

		targetValue, err := qh.evalColumn(targetRow, targetColumn, nil)
		if err != nil {
			return false, err
		}

		if !isEqual(srcValue, targetValue) {
			return false, nil
		}
	}

	return true, nil
}
