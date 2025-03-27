package main

import (
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

var capabilities = schema.CapabilitiesResponse{
	Version: schema.SchemaVersion,
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
		Mutation: schema.MutationCapabilities{
			Explain:       &schema.LeafCapability{},
			Transactional: &schema.LeafCapability{},
		},
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
		"Int": schema.ScalarType{
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{
				"max": schema.NewAggregateFunctionDefinitionMax().Encode(),
				"min": schema.NewAggregateFunctionDefinitionMin().Encode(),
			},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq": schema.NewComparisonOperatorEqual().Encode(),
				"in": schema.NewComparisonOperatorIn().Encode(),
			},
			Representation: schema.NewTypeRepresentationInt32().Encode(),
		},
		"String": {
			AggregateFunctions: schema.ScalarTypeAggregateFunctions{},
			ComparisonOperators: map[string]schema.ComparisonOperatorDefinition{
				"eq":   schema.NewComparisonOperatorEqual().Encode(),
				"in":   schema.NewComparisonOperatorIn().Encode(),
				"like": schema.NewComparisonOperatorCustom(schema.NewNamedType("String")).Encode(),
			},
			Representation: schema.NewTypeRepresentationString().Encode(),
		},
	},
	ObjectTypes: schema.SchemaResponseObjectTypes{
		"article": schema.NewObjectType(
			schema.ObjectTypeFields{
				"author_id": schema.ObjectField{
					Description: utils.ToPtr("The article's author ID"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"id": {
					Description: utils.ToPtr("The article's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"title": {
					Description: utils.ToPtr("The article's title"),
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
			utils.ToPtr("An article"),
		),
		"author": schema.NewObjectType(
			schema.ObjectTypeFields{
				"first_name": {
					Description: utils.ToPtr("The author's first name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"id": {
					Description: utils.ToPtr("The author's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"last_name": {
					Description: utils.ToPtr("The author's last name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			},
			nil,
			utils.ToPtr("An author"),
		),
		"institution": schema.NewObjectType(
			schema.ObjectTypeFields{
				"departments": schema.ObjectField{
					Description: utils.ToPtr("The institution's departments"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
				},
				"id": schema.ObjectField{
					Description: utils.ToPtr("The institution's primary key"),
					Type:        schema.NewNamedType("Int").Encode(),
				},
				"location": schema.ObjectField{
					Description: utils.ToPtr("The institution's location"),
					Type:        schema.NewNamedType("location").Encode(),
				},
				"name": schema.ObjectField{
					Description: utils.ToPtr("The institution's name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"staff": schema.ObjectField{
					Description: utils.ToPtr("The institution's staff"),
					Type:        schema.NewArrayType(schema.NewNamedType("staff_member")).Encode(),
				},
			},
			nil,
			utils.ToPtr("An institution"),
		),
		"location": schema.NewObjectType(
			schema.ObjectTypeFields{
				"campuses": schema.ObjectField{
					Description: utils.ToPtr("The location's campuses"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
				},
				"city": schema.ObjectField{
					Description: utils.ToPtr("The location's city"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"country": schema.ObjectField{
					Description: utils.ToPtr("The location's country"),
					Type:        schema.NewNamedType("String").Encode(),
				},
			}, nil,
			utils.ToPtr("A location"),
		),
		"staff_member": schema.NewObjectType(
			schema.ObjectTypeFields{
				"first_name": schema.ObjectField{
					Description: utils.ToPtr("The staff member's first name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"last_name": schema.ObjectField{
					Description: utils.ToPtr("The staff member's last name"),
					Type:        schema.NewNamedType("String").Encode(),
				},
				"specialities": schema.ObjectField{
					Description: utils.ToPtr("The staff member's specialities"),
					Type:        schema.NewArrayType(schema.NewNamedType("String")).Encode(),
				},
			},
			nil,
			utils.ToPtr("A staff member"),
		),
	},
	Collections: []schema.CollectionInfo{
		{
			Name:        "articles",
			Description: utils.ToPtr("A collection of articles"),
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
			Description: utils.ToPtr("A collection of authors"),
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
			Description: utils.ToPtr("A collection of institutions"),
			Arguments:   schema.CollectionInfoArguments{},
			Type:        "institution",
			UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
				"InstitutionByID": schema.UniquenessConstraint{
					UniqueColumns: []string{"id"},
				},
			},
		},
		{
			Name:        "articles_by_author",
			Description: utils.ToPtr("Articles parameterized by author"),
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
			Description: utils.ToPtr("Get the ID of the most recent article"),
			Arguments:   schema.FunctionInfoArguments{},
			ResultType:  schema.NewNullableNamedType("Int").Encode(),
		},
		{
			Name:        "latest_article",
			Description: utils.ToPtr("Get the most recent article"),
			Arguments:   schema.FunctionInfoArguments{},
			ResultType:  schema.NewNullableNamedType("article").Encode(),
		},
	},
	Procedures: []schema.ProcedureInfo{
		{
			Name:        "upsert_article",
			Description: utils.ToPtr("Insert or update an article"),
			Arguments: schema.ProcedureInfoArguments{
				"article": schema.ArgumentInfo{
					Description: utils.ToPtr("The article to insert or update"),
					Type:        schema.NewNamedType("article").Encode(),
				},
			},
			ResultType: schema.NewNullableNamedType("article").Encode(),
		},
		{
			Name:        "delete_articles",
			Description: utils.ToPtr("Delete articles which match a predicate"),
			Arguments: schema.ProcedureInfoArguments{
				"where": schema.ArgumentInfo{
					Description: utils.ToPtr("The predicate"),
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
}
