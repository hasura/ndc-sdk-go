package schema

import (
	"encoding/json"
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
)

func TestSchemaResponse(t *testing.T) {
	rawInput := []byte(`{
		"scalar_types": {
			"Int": {
				"aggregate_functions": {
					"max": {
						"result_type": {
							"type": "nullable",
							"underlying_type": {
								"type": "named",
								"name": "Int"
							}
						}
					},
					"min": {
						"result_type": {
							"type": "nullable",
							"underlying_type": {
								"type": "named",
								"name": "Int"
							}
						}
					}
				},
				"comparison_operators": {}
			},
			"String": {
				"aggregate_functions": {},
				"comparison_operators": {
					"like": {
						"argument_type": {
							"type": "named",
							"name": "String"
						}
					}
				}
			}
		},
		"object_types": {
			"article": {
				"description": "An article",
				"fields": {
					"author_id": {
						"description": "The article's author ID",
						"type": {
							"type": "named",
							"name": "Int"
						}
					},
					"id": {
						"description": "The article's primary key",
						"type": {
							"type": "named",
							"name": "Int"
						}
					},
					"title": {
						"description": "The article's title",
						"type": {
							"type": "named",
							"name": "String"
						}
					}
				}
			},
			"author": {
				"description": "An author",
				"fields": {
					"first_name": {
						"description": "The author's first name",
						"type": {
							"type": "named",
							"name": "String"
						}
					},
					"id": {
						"description": "The author's primary key",
						"type": {
							"type": "named",
							"name": "Int"
						}
					},
					"last_name": {
						"description": "The author's last name",
						"type": {
							"type": "named",
							"name": "String"
						}
					}
				}
			}
		},
		"collections": [
			{
				"name": "articles",
				"description": "A collection of articles",
				"arguments": {},
				"type": "article",
				"uniqueness_constraints": {
					"ArticleByID": {
						"unique_columns": [
							"id"
						]
					}
				},
				"foreign_keys": {
					"Article_AuthorID": {
						"column_mapping": {
							"author_id": "id"
						},
						"foreign_collection": "authors"
					}
				}
			},
			{
				"name": "authors",
				"description": "A collection of authors",
				"arguments": {},
				"type": "author",
				"uniqueness_constraints": {
					"AuthorByID": {
						"unique_columns": [
							"id"
						]
					}
				},
				"foreign_keys": {}
			},
			{
				"name": "articles_by_author",
				"description": "Articles parameterized by author",
				"arguments": {
					"author_id": {
						"type": {
							"type": "named",
							"name": "Int"
						}
					}
				},
				"type": "article",
				"uniqueness_constraints": {},
				"foreign_keys": {}
			}
		],
		"functions": [
			{
				"name": "latest_article_id",
				"description": "Get the ID of the most recent article",
				"arguments": {},
				"result_type": {
					"type": "nullable",
					"underlying_type": {
						"type": "named",
						"name": "Int"
					}
				}
			}
		],
		"procedures": [
			{
				"name": "upsert_article",
				"description": "Insert or update an article",
				"arguments": {
					"article": {
						"description": "The article to insert or update",
						"type": {
							"type": "named",
							"name": "article"
						}
					}
				},
				"result_type": {
					"type": "nullable",
					"underlying_type": {
						"type": "named",
						"name": "article"
					}
				}
			}
		]
	}`)

	var resp SchemaResponse
	if err := json.Unmarshal(rawInput, &resp); err != nil {
		t.Errorf("failed to decode SchemaResponse: %s", err)
		t.FailNow()
	}

	if intScalar, ok := resp.ScalarTypes["Int"]; !ok {
		t.Error("Int scalar in SchemaResponse: required")
		t.FailNow()
	} else if len(intScalar.ComparisonOperators) != 0 {
		t.Errorf("Int scalar in SchemaResponse: expected no comparison operator, got: %+v", intScalar.ComparisonOperators)
		t.FailNow()
	} else if !internal.DeepEqual(intScalar.AggregateFunctions, ScalarTypeAggregateFunctions{
		"max": AggregateFunctionDefinition{
			ResultType: NewNullableNamedType("Int").Encode(),
		},
		"min": AggregateFunctionDefinition{
			ResultType: NewNullableNamedType("Int").Encode(),
		},
	}) {
		t.Errorf("Int scalar in SchemaResponse: expected equal aggregate functions; %+v", intScalar.AggregateFunctions)
		t.FailNow()
	} else if _, err := intScalar.AggregateFunctions["max"].ResultType.AsNullable(); err != nil {
		t.Errorf("Int scalar in SchemaResponse: expected aggregate function max, got error: %s", err)
		t.FailNow()
	}

	if stringScalar, ok := resp.ScalarTypes["String"]; !ok {
		t.Error("String scalar in SchemaResponse: required")
		t.FailNow()
	} else if len(stringScalar.AggregateFunctions) != 0 {
		t.Errorf("Int scalar in SchemaResponse: expected no aggregate function, got: %+v", stringScalar.AggregateFunctions)
		t.FailNow()
	} else if !internal.DeepEqual(stringScalar.ComparisonOperators, ScalarTypeComparisonOperators{
		"like": ComparisonOperatorDefinition{
			ArgumentType: NewNamedType("String").Encode(),
		},
	}) {
		t.Errorf("String scalar in SchemaResponse: expected equal comparison operators; %+v", stringScalar.ComparisonOperators)
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
						"first_name": NewColumnField("first_name").Encode(),
						"last_name":  NewColumnField("last_name").Encode(),
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
