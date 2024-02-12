package schema

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
)

func TestSchemaResponse(t *testing.T) {
	// reuse test fixtures from ndc-reference
	httpResp, err := http.Get("https://raw.githubusercontent.com/hasura/ndc-spec/main/ndc-reference/tests/schema/expected.json")
	if err != nil {
		t.Errorf("failed to fetch schema: %s", err.Error())
		t.FailNow()
	}

	var resp SchemaResponse
	if err = json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode schema json: %s", err.Error())
		t.FailNow()
	}

	expected := SchemaResponse{
		ScalarTypes: SchemaResponseScalarTypes{
			"Int": ScalarType{
				AggregateFunctions: ScalarTypeAggregateFunctions{
					"max": AggregateFunctionDefinition{
						ResultType: NewNullableNamedType("Int").Encode(),
					},
					"min": AggregateFunctionDefinition{
						ResultType: NewNullableNamedType("Int").Encode(),
					},
				},
				ComparisonOperators: map[string]ComparisonOperatorDefinition{
					"eq": NewComparisonOperatorEqual().Encode(),
					"in": NewComparisonOperatorIn().Encode(),
				},
			},
			"String": {
				AggregateFunctions: ScalarTypeAggregateFunctions{},
				ComparisonOperators: map[string]ComparisonOperatorDefinition{
					"eq":   NewComparisonOperatorEqual().Encode(),
					"in":   NewComparisonOperatorIn().Encode(),
					"like": NewComparisonOperatorCustom(NewNamedType("String")).Encode(),
				},
			},
		},
		ObjectTypes: SchemaResponseObjectTypes{
			"article": ObjectType{
				Description: ToPtr("An article"),
				Fields: ObjectTypeFields{
					"author_id": ObjectField{
						Description: ToPtr("The article's author ID"),
						Type:        NewNamedType("Int").Encode(),
					},
					"id": {
						Description: ToPtr("The article's primary key"),
						Type:        NewNamedType("Int").Encode(),
					},
					"title": {
						Description: ToPtr("The article's title"),
						Type:        NewNamedType("String").Encode(),
					},
				},
			},
			"author": ObjectType{
				Description: ToPtr("An author"),
				Fields: ObjectTypeFields{
					"first_name": {
						Description: ToPtr("The author's first name"),
						Type:        NewNamedType("String").Encode(),
					},
					"id": {
						Description: ToPtr("The author's primary key"),
						Type:        NewNamedType("Int").Encode(),
					},
					"last_name": {
						Description: ToPtr("The author's last name"),
						Type:        NewNamedType("String").Encode(),
					},
				},
			},
			"institution": ObjectType{
				Description: ToPtr("An institution"),
				Fields: ObjectTypeFields{
					"departments": ObjectField{
						Description: ToPtr("The institution's departments"),
						Type:        NewArrayType(NewNamedType("String")).Encode(),
					},
					"id": ObjectField{
						Description: ToPtr("The institution's primary key"),
						Type:        NewNamedType("Int").Encode(),
					},
					"location": ObjectField{
						Description: ToPtr("The institution's location"),
						Type:        NewNamedType("location").Encode(),
					},
					"name": ObjectField{
						Description: ToPtr("The institution's name"),
						Type:        NewNamedType("String").Encode(),
					},
					"staff": ObjectField{
						Description: ToPtr("The institution's staff"),
						Type:        NewArrayType(NewNamedType("staff_member")).Encode(),
					},
				},
			},
			"location": ObjectType{
				Description: ToPtr("A location"),
				Fields: ObjectTypeFields{
					"campuses": ObjectField{
						Description: ToPtr("The location's campuses"),
						Type:        NewArrayType(NewNamedType("String")).Encode(),
					},
					"city": ObjectField{
						Description: ToPtr("The location's city"),
						Type:        NewNamedType("String").Encode(),
					},
					"country": ObjectField{
						Description: ToPtr("The location's country"),
						Type:        NewNamedType("String").Encode(),
					},
				},
			},
			"staff_member": ObjectType{
				Description: ToPtr("A staff member"),
				Fields: ObjectTypeFields{
					"first_name": ObjectField{
						Description: ToPtr("The staff member's first name"),
						Type:        NewNamedType("String").Encode(),
					},
					"last_name": ObjectField{
						Description: ToPtr("The staff member's last name"),
						Type:        NewNamedType("String").Encode(),
					},
					"specialities": ObjectField{
						Description: ToPtr("The staff member's specialities"),
						Type:        NewArrayType(NewNamedType("String")).Encode(),
					},
				},
			},
		},
		Collections: []CollectionInfo{
			{
				Name:        "articles",
				Description: ToPtr("A collection of articles"),
				Arguments:   CollectionInfoArguments{},
				Type:        "article",
				UniquenessConstraints: CollectionInfoUniquenessConstraints{
					"ArticleByID": UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
				ForeignKeys: CollectionInfoForeignKeys{
					"Article_AuthorID": ForeignKeyConstraint{
						ColumnMapping: ForeignKeyConstraintColumnMapping{
							"author_id": "id",
						},
						ForeignCollection: "authors",
					},
				},
			},
			{
				Name:        "authors",
				Description: ToPtr("A collection of authors"),
				Arguments:   CollectionInfoArguments{},
				Type:        "author",
				UniquenessConstraints: CollectionInfoUniquenessConstraints{
					"AuthorByID": UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
				ForeignKeys: CollectionInfoForeignKeys{},
			},
			{
				Name:        "institutions",
				Description: ToPtr("A collection of institutions"),
				Arguments:   CollectionInfoArguments{},
				Type:        "institution",
				UniquenessConstraints: CollectionInfoUniquenessConstraints{
					"InstitutionByID": UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
				ForeignKeys: CollectionInfoForeignKeys{},
			},
			{
				Name:        "articles_by_author",
				Description: ToPtr("Articles parameterized by author"),
				Arguments: CollectionInfoArguments{
					"author_id": ArgumentInfo{
						Type: NewNamedType("Int").Encode(),
					},
				},
				Type:                  "article",
				UniquenessConstraints: CollectionInfoUniquenessConstraints{},
				ForeignKeys:           CollectionInfoForeignKeys{},
			},
		},
		Functions: []FunctionInfo{
			{
				Name:        "latest_article_id",
				Description: ToPtr("Get the ID of the most recent article"),
				Arguments:   FunctionInfoArguments{},
				ResultType:  NewNullableNamedType("Int").Encode(),
			},
		},
		Procedures: []ProcedureInfo{
			{
				Name:        "upsert_article",
				Description: ToPtr("Insert or update an article"),
				Arguments: ProcedureInfoArguments{
					"article": ArgumentInfo{
						Description: ToPtr("The article to insert or update"),
						Type:        NewNamedType("article").Encode(),
					},
				},
				ResultType: NewNullableNamedType("article").Encode(),
			},
			{
				Name:        "delete_articles",
				Description: ToPtr("Delete articles which match a predicate"),
				Arguments: ProcedureInfoArguments{
					"where": ArgumentInfo{
						Description: ToPtr("The predicate"),
						Type:        NewPredicateType("article").Encode(),
					},
				},
				ResultType: NewArrayType(NewNamedType("article")).Encode(),
			},
		},
	}

	if !internal.DeepEqual(expected.Collections, resp.Collections) {
		t.Errorf("collections in SchemaResponse: unexpected equality;\nexpected:	%+v,\n got:	%+v\n", expected.Collections, resp.Collections)
		t.FailNow()
	}

	if !internal.DeepEqual(expected.Functions, resp.Functions) {
		t.Errorf("functions in SchemaResponse: unexpected equality;\nexpected:	%+v,\n got:	%+v\n", expected.Functions, resp.Functions)
		t.FailNow()
	}

	if !internal.DeepEqual(expected.Procedures, resp.Procedures) {
		t.Errorf("procedures in SchemaResponse: unexpected equality;\nexpected:	%+v,\n got:	%+v\n", expected.Procedures, resp.Procedures)
		t.FailNow()
	}
}

func TestQueryRequest(t *testing.T) {
	testCases := []struct {
		name     string
		rawJson  []byte
		expected QueryRequest
	}{
		{
			"aggregate_function",
			[]byte(`{
				"collection": "articles",
				"arguments": {},
				"query": {
						"aggregates": {
								"min_id": {
										"type": "single_column",
										"column": "id",
										"function": "min"
								},
								"max_id": {
										"type": "single_column",
										"column": "id",
										"function": "max"
								}
						}
				},
				"collection_relationships": {}
		}`),
			QueryRequest{
				Collection: "articles",
				Query: Query{
					Aggregates: QueryAggregates{
						"min_id": NewAggregateSingleColumn("id", "min").Encode(),
						"max_id": NewAggregateSingleColumn("id", "max").Encode(),
					},
				},
			},
		},
		{
			"authors_with_article_aggregate",
			[]byte(`{
				"collection": "authors",
				"arguments": {},
				"query": {
						"fields": {
								"first_name": {
										"type": "column",
										"column": "first_name"
								},
								"last_name": {
										"type": "column",
										"column": "last_name"
								},
								"articles": {
										"type": "relationship",
										"arguments": {},
										"relationship": "author_articles",
										"query": {
												"aggregates": {
														"count": {
																"type": "star_count"
														}
												}
										}
								}
						}
				},
				"collection_relationships": {
						"author_articles": {
								"arguments": {},
								"column_mapping": {
										"id": "author_id"
								},
								"relationship_type": "array",
								"target_collection": "articles"
						}
				}
		}`),
			QueryRequest{
				Collection: "authors",
				Query: Query{
					Fields: QueryFields{
						"first_name": NewColumnField("first_name", nil).Encode(),
						"last_name":  NewColumnField("last_name", nil).Encode(),
						"articles": NewRelationshipField(
							Query{
								Aggregates: QueryAggregates{
									"count": NewAggregateStarCount().Encode(),
								},
							},
							"author_articles",
							map[string]RelationshipArgument{},
						).Encode(),
					},
				},
				CollectionRelationships: QueryRequestCollectionRelationships{
					"author_articles": Relationship{
						ColumnMapping: RelationshipColumnMapping{
							"id": "author_id",
						},
						RelationshipType: RelationshipTypeArray,
						TargetCollection: "articles",
						Arguments:        RelationshipArguments{},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var q QueryRequest
			if err := json.Unmarshal(tc.rawJson, &q); err != nil {
				t.Errorf("failed to decode: %s", err)
				t.FailNow()
			}

			if !internal.DeepEqual(tc.expected.Collection, q.Collection) {
				t.Errorf("unexpected collection equality; expected: %+v, got: %+v", tc.expected.Collection, q.Collection)
				t.FailNow()
			}

			if !internal.DeepEqual(tc.expected.Query, q.Query) {
				t.Errorf("unexpected Query equality;\n expected: %+v,\n got: %+v\n", tc.expected.Query, q.Query)
				t.FailNow()
			}

			if !internal.DeepEqual(tc.expected.CollectionRelationships, q.CollectionRelationships) {
				t.Errorf("unexpected CollectionRelationships equality;\n expected: %+v,\n got: %+v\n", tc.expected.CollectionRelationships, q.CollectionRelationships)
				t.FailNow()
			}
		})
	}
}
