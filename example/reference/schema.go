package main

import (
	"github.com/hasura/ndc-sdk-go/v2/schema"
)

var capabilities = schema.CapabilitiesResponse{
	// the reference connector in the Rust SDK doesn't have the 'v' prefix
	Version: schema.NDCVersion,
	Capabilities: schema.Capabilities{
		Query: schema.QueryCapabilities{
			Aggregates: &schema.AggregateCapabilities{
				FilterBy: &schema.LeafCapability{},
				GroupBy: &schema.GroupByCapabilities{
					Filter:   &schema.LeafCapability{},
					Order:    &schema.LeafCapability{},
					Paginate: &schema.LeafCapability{},
				},
			},
			Variables: &schema.LeafCapability{},
			NestedFields: schema.NestedFieldCapabilities{
				FilterBy: &schema.NestedFieldFilterByCapabilities{
					NestedArrays: &schema.NestedArrayFilterByCapabilities{
						Contains: &schema.LeafCapability{},
						IsEmpty:  &schema.LeafCapability{},
					},
				},
				OrderBy:           &schema.LeafCapability{},
				Aggregates:        &schema.LeafCapability{},
				NestedCollections: &schema.LeafCapability{},
			},
			Exists: schema.ExistsCapabilities{
				NamedScopes:             &schema.LeafCapability{},
				NestedCollections:       &schema.LeafCapability{},
				NestedScalarCollections: &schema.LeafCapability{},
				Unrelated:               &schema.LeafCapability{},
			},
		},
		Mutation: schema.MutationCapabilities{},
		Relationships: &schema.RelationshipCapabilities{
			OrderByAggregate:    &schema.LeafCapability{},
			RelationComparisons: &schema.LeafCapability{},
			Nested: &schema.NestedRelationshipCapabilities{
				Array:     &schema.LeafCapability{},
				Filtering: &schema.LeafCapability{},
				Ordering:  &schema.LeafCapability{},
			},
		},
	},
}

var ndcSchema = schema.SchemaResponse{
	ScalarTypes: schema.SchemaResponseScalarTypes{
		"Date": {
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq": schema.NewComparisonOperatorEqual().Encode(),
				"in": schema.NewComparisonOperatorIn().Encode(),
			},
			ExtractionFunctions: schema.ScalarTypeExtractionFunctions{
				"day":   schema.NewExtractionFunctionDefinitionDay("Int").Encode(),
				"month": schema.NewExtractionFunctionDefinitionMonth("Int").Encode(),
				"year":  schema.NewExtractionFunctionDefinitionYear("Int").Encode(),
			},
			Representation: schema.TypeRepresentation{
				"type": schema.TypeRepresentationType("date"),
			},
		},
		"Float": {
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"avg": schema.NewAggregateFunctionDefinitionAverage("Float").Encode(),
				"max": schema.NewAggregateFunctionDefinitionMax().Encode(),
				"min": schema.NewAggregateFunctionDefinitionMin().Encode(),
				"sum": schema.NewAggregateFunctionDefinitionSum("Float").Encode(),
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq":  schema.NewComparisonOperatorEqual().Encode(),
				"gt":  schema.NewComparisonOperatorGreaterThan().Encode(),
				"gte": schema.NewComparisonOperatorGreaterThanOrEqual().Encode(),
				"in":  schema.NewComparisonOperatorIn().Encode(),
				"lt":  schema.NewComparisonOperatorLessThan().Encode(),
				"lte": schema.NewComparisonOperatorLessThanOrEqual().Encode(),
			},
			ExtractionFunctions: schema.ScalarTypeExtractionFunctions{},
			Representation:      schema.NewTypeRepresentationFloat64().Encode(),
		},
		"Int": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"avg": schema.NewAggregateFunctionDefinitionAverage("Float").Encode(),
				"max": schema.NewAggregateFunctionDefinitionMax().Encode(),
				"min": schema.NewAggregateFunctionDefinitionMin().Encode(),
				"sum": schema.NewAggregateFunctionDefinitionSum("Int64").Encode(),
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq":  schema.NewComparisonOperatorEqual().Encode(),
				"gt":  schema.NewComparisonOperatorGreaterThan().Encode(),
				"gte": schema.NewComparisonOperatorGreaterThanOrEqual().Encode(),
				"in":  schema.NewComparisonOperatorIn().Encode(),
				"lt":  schema.NewComparisonOperatorLessThan().Encode(),
				"lte": schema.NewComparisonOperatorLessThanOrEqual().Encode(),
			},
			Representation: schema.NewTypeRepresentationInt32().Encode(),
		},
		"Int64": {
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"avg": {
					"result_type": string("Float"),
					"type":        schema.AggregateFunctionDefinitionType("average"),
				},
				"max": {"type": schema.AggregateFunctionDefinitionType("max")},
				"min": {"type": schema.AggregateFunctionDefinitionType("min")},
				"sum": {
					"result_type": string("Int64"),
					"type":        schema.AggregateFunctionDefinitionType("sum"),
				},
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq":  schema.NewComparisonOperatorEqual().Encode(),
				"gt":  schema.NewComparisonOperatorGreaterThan().Encode(),
				"gte": schema.NewComparisonOperatorGreaterThanOrEqual().Encode(),
				"in":  schema.NewComparisonOperatorIn().Encode(),
				"lt":  schema.NewComparisonOperatorLessThan().Encode(),
				"lte": schema.NewComparisonOperatorLessThanOrEqual().Encode(),
			},
			ExtractionFunctions: schema.ScalarTypeExtractionFunctions{},
			Representation: schema.TypeRepresentation{
				"type": schema.TypeRepresentationType("int64"),
			},
		},
		"String": {
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"max": schema.NewAggregateFunctionDefinitionMax().Encode(),
				"min": schema.NewAggregateFunctionDefinitionMin().Encode(),
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"contains":   schema.NewComparisonOperatorContains().Encode(),
				"ends_with":  schema.NewComparisonOperatorEndsWith().Encode(),
				"eq":         schema.NewComparisonOperatorEqual().Encode(),
				"gt":         schema.NewComparisonOperatorGreaterThan().Encode(),
				"gte":        schema.NewComparisonOperatorGreaterThanOrEqual().Encode(),
				"icontains":  schema.NewComparisonOperatorContainsInsensitive().Encode(),
				"iends_with": schema.NewComparisonOperatorEndsWithInsensitive().Encode(),
				"in":         schema.NewComparisonOperatorIn().Encode(),
				"like": schema.NewComparisonOperatorCustom(schema.NewNamedType("String")).
					Encode(),
				"lt":           schema.NewComparisonOperatorLessThan().Encode(),
				"lte":          schema.NewComparisonOperatorLessThanOrEqual().Encode(),
				"starts_with":  schema.NewComparisonOperatorStartsWith().Encode(),
				"istarts_with": schema.NewComparisonOperatorStartsWithInsensitive().Encode(),
			},
			Representation: schema.NewTypeRepresentationString().Encode(),
		},
	},
	ObjectTypes: schema.SchemaResponseObjectTypes{
		"article": schema.NewObjectType(
			schema.ObjectTypeFields{
				"author_id": schema.ObjectField{
					Description: new("The article's author ID"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"id": {
					Description: new("The article's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"published_date": {
					Description: new("The article's date of publication"),
					Type:        schema.NewNamedType("Date").Encode(),
				},
				"title": {
					Description: new("The article's title"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			},
			schema.ObjectTypeForeignKeys{
				"Article_AuthorID": schema.ForeignKeyConstraint{
					ColumnMapping: schema.ForeignKeyConstraintColumnMapping{
						"author_id": []string{"id"},
					},
					ForeignCollection: "authors",
				},
			},
			new("An article"),
		),
		"author": schema.NewObjectType(
			schema.ObjectTypeFields{
				"first_name": {
					Description: new("The author's first name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"id": {
					Description: new("The author's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"last_name": {
					Description: new("The author's last name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			},
			nil,
			new("An author"),
		),
		"city": {
			Description: new("A city"),
			Fields: schema.ObjectTypeFields{
				"name": {
					Description: new("The institution's name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			},
		},
		"country": {
			Description: new("A country"),
			Fields: schema.ObjectTypeFields{
				"area_km2": {
					Description: new("The country's area size in square kilometers"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"cities": {
					Arguments: schema.ObjectFieldArguments{
						"limit": {
							Type: schema.NewNullableNamedType("Int").Encode(),
						},
					},
					Description: new("The cities in the country"),
					Type:        schema.NewArrayType(schema.NewNamedType("city")).Encode(),
				},
				"id": {
					Description: new("The country's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"name": {
					Description: new("The country's name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			},
			ForeignKeys: schema.ObjectTypeForeignKeys{},
		},
		"institution": schema.NewObjectType(
			schema.ObjectTypeFields{
				"departments": schema.ObjectField{
					Arguments: schema.ObjectFieldArguments{
						"limit": {
							Type: schema.NewNullableNamedType("Int").Encode(),
						},
					},
					Description: new("The institution's departments"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
				},
				"id": schema.ObjectField{
					Description: new("The institution's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"location": schema.ObjectField{
					Description: new("The institution's location"),
					Type:        schema.NewNamedType("location").Encode(),
				},
				"name": schema.ObjectField{
					Description: new("The institution's name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"staff": schema.ObjectField{
					Description: new("The institution's staff"),
					Type:        schema.NewArrayType(schema.NewNamedType("staff_member")).Encode(),
					Arguments: schema.ObjectFieldArguments{
						"limit": {
							Type: schema.NewNullableNamedType("Int").Encode(),
						},
					},
				},
			},
			nil,
			new("An institution"),
		),
		"location": schema.NewObjectType(
			schema.ObjectTypeFields{
				"campuses": schema.ObjectField{
					Description: new("The location's campuses"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
					Arguments: schema.ObjectFieldArguments{
						"limit": {
							Type: schema.NewNullableNamedType("Int").Encode(),
						},
					},
				},
				"city": schema.ObjectField{
					Description: new("The location's city"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"country": schema.ObjectField{
					Description: new("The location's country"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"country_id": {
					Description: new("The location's country ID"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
			},
			schema.ObjectTypeForeignKeys{
				"Location_CountryID": {
					ColumnMapping: schema.ForeignKeyConstraintColumnMapping{
						"country_id": {"id"},
					},
					ForeignCollection: "countries",
				},
			},
			new("A location"),
		),
		"staff_member": schema.NewObjectType(
			schema.ObjectTypeFields{
				"born_country_id": {
					Description: new("The ID of the country the staff member was born in"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"first_name": schema.ObjectField{
					Description: new("The staff member's first name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"last_name": schema.ObjectField{
					Description: new("The staff member's last name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"specialities": schema.ObjectField{
					Description: new("The staff member's specialities"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
					Arguments: schema.ObjectFieldArguments{
						"limit": {
							Type: schema.NewNullableNamedType("Int").Encode(),
						},
					},
				},
			},
			schema.ObjectTypeForeignKeys{
				"Staff_BornCountryID": {
					ColumnMapping: schema.ForeignKeyConstraintColumnMapping{
						"born_country_id": {"id"},
					},
					ForeignCollection: "countries",
				},
			},
			new("A staff member"),
		),
	},
	Collections: []schema.CollectionInfo{
		{
			Name:        "articles",
			Description: new("A collection of articles"),
			Arguments:   schema.CollectionInfoArguments{},
			Type:        "article",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"ArticleByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
		{
			Name:        "authors",
			Description: new("A collection of authors"),
			Arguments:   schema.CollectionInfoArguments{},
			Type:        "author",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"AuthorByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
		{
			Name:        "institutions",
			Description: new("A collection of institutions"),
			Arguments:   schema.CollectionInfoArguments{},
			Type:        "institution",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"InstitutionByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
		{
			Name:        "countries",
			Description: new("A collection of countries"),
			Arguments:   schema.CollectionInfoArguments{},
			Type:        "country",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"CountryByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
		{
			Name:        "articles_by_author",
			Description: new("Articles parameterized by author"),
			Arguments: schema.CollectionInfoArguments{
				"author_id": schema.ArgumentInfo{
					Type: schema.NewNamedType("Int").Encode(),
				},
			},
			Type:                  "article",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{},
		},
	},
	Functions: []schema.FunctionInfo{
		{
			Name:        "latest_article_id",
			Description: new("Get the ID of the most recent article"),
			Arguments:   schema.FunctionInfoArguments{},
			ResultType:  schema.NewNullableNamedType("Int").Encode(),
		},
		{
			Name:        "latest_article",
			Description: new("Get the most recent article"),
			Arguments:   schema.FunctionInfoArguments{},
			ResultType:  schema.NewNullableNamedType("article").Encode(),
		},
	},
	Procedures: []schema.ProcedureInfo{
		{
			Name:        "upsert_article",
			Description: new("Insert or update an article"),
			Arguments: schema.ProcedureInfoArguments{
				"article": schema.ArgumentInfo{
					Description: new("The article to insert or update"),
					Type:        schema.NewNamedType("article").Encode(),
				},
			},
			ResultType: schema.NewNullableNamedType("article").Encode(),
		},
		{
			Name:        "delete_articles",
			Description: new("Delete articles which match a predicate"),
			Arguments: schema.ProcedureInfoArguments{
				"where": schema.ArgumentInfo{
					Description: new("The predicate"),
					Type:        schema.NewPredicateType("article").Encode(),
				},
			},
			ResultType: schema.NewArrayType(schema.NewNamedType("article")).Encode(),
		},
	},
	Capabilities: &schema.CapabilitySchemaInfo{
		Query: &schema.QueryCapabilitiesSchemaInfo{
			Aggregates: &schema.AggregateCapabilitiesSchemaInfo{
				CountScalarType: "Int",
			},
		},
	},
	RequestArguments: &schema.RequestLevelArguments{
		MutationArguments:        make(schema.RequestLevelArgumentsMutationArguments),
		QueryArguments:           make(schema.RequestLevelArgumentsQueryArguments),
		RelationalQueryArguments: make(schema.RequestLevelArgumentsRelationalQueryArguments),
	},
}
