package schema

import (
	"encoding/json"
	"testing"

	"github.com/hasura/ndc-sdk-go/internal"
)

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
								"count": {
									"type": "star_count"
								},
								"min_id": {
										"type": "single_column",
										"column": "id",
										"function": "min",
										"field_path": ["foo"]
								},
								"bar_count": {
										"type": "column_count",
										"column": "id",
										"distinct": true,
										"field_path": ["bar"]
								}
						}
				},
				"collection_relationships": {}
		}`),
			QueryRequest{
				Collection: "articles",
				Query: Query{
					Aggregates: QueryAggregates{
						"count":     NewAggregateStarCount().Encode(),
						"min_id":    NewAggregateSingleColumn("id", "min", []string{"foo"}).Encode(),
						"bar_count": NewAggregateColumnCount("id", true, []string{"bar"}).Encode(),
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
		{
			"predicate_with_nested_field_eq",
			[]byte(`{
				"collection": "institutions",
				"arguments": {},
				"query": {
						"fields": {
								"id": {
										"type": "column",
										"column": "id"
								},
								"location": {
										"type": "column",
										"column": "location",
										"fields": {
												"type": "object",
												"fields": {
														"city": {
																"type": "column",
																"column": "city"
														},
														"campuses": {
																"type": "column",
																"column": "campuses",
																"arguments": {
																		"limit": {
																				"type": "literal",
																				"value": null
																		}
																}
														}
												}
										}
								}
						},
						"predicate": {
								"type": "binary_comparison_operator",
								"column": {
										"type": "column",
										"name": "location",
										"field_path": ["city"],
										"path": []
								},
								"operator": "eq",
								"value": {
										"type": "scalar",
										"value": "London"
								}
						}
				},
				"collection_relationships": {}
		}`),
			QueryRequest{
				Collection: "institutions",
				Query: Query{
					Fields: QueryFields{
						"id": NewColumnField("id", nil).Encode(),
						"location": NewColumnField("location", NewNestedObject(map[string]FieldEncoder{
							"city": NewColumnField("city", nil),
							"campuses": NewColumnFieldWithArguments("campuses", nil, map[string]Argument{
								"limit": *NewLiteralArgument(nil),
							}),
						})).Encode(),
					},
					Predicate: NewExpressionBinaryComparisonOperator(*NewComparisonTargetColumn("location", []string{"city"}, []PathElement{}), "eq", NewComparisonValueScalar("London")).Encode(),
				},
				CollectionRelationships: QueryRequestCollectionRelationships{},
			},
		},
		{
			"order_by_column",
			[]byte(`{
				"collection": "articles",
				"arguments": {},
				"query": {
						"fields": {
								"title": {
										"type": "column",
										"column": "title"
								}
						},
						"order_by": {
								"elements": [
										{
												"target": {
														"type": "column",
														"name": "title",
														"path": []
												},
												"order_direction": "desc"
										}
								]
						}
				},
				"collection_relationships": {}
		}`),
			QueryRequest{
				Collection: "articles",
				Arguments:  QueryRequestArguments{},
				Query: Query{
					Fields: QueryFields{
						"title": NewColumnField("title", nil).Encode(),
					},
					OrderBy: &OrderBy{
						Elements: []OrderByElement{
							{
								OrderDirection: OrderDirectionDesc,
								Target:         NewOrderByColumnName("title").Encode(),
							},
						},
					},
				},
				CollectionRelationships: QueryRequestCollectionRelationships{},
			},
		},
		{
			"order_by_nested_field",
			[]byte(`{
				"collection": "institutions",
				"arguments": {},
				"query": {
						"fields": {
								"location": {
										"type": "column",
										"column": "location",
										"fields": {
												"type": "object",
												"fields": {
														"country": {
																"type": "column",
																"column": "country"
														}
												}
										}
								}
						},
						"order_by": {
								"elements": [
										{
												"target": {
														"type": "column",
														"name": "location",
														"field_path": ["country"],
														"path": []
												},
												"order_direction": "asc"
										}
								]
						}
				},
				"collection_relationships": {}
		}`),
			QueryRequest{
				Collection:              "institutions",
				Arguments:               QueryRequestArguments{},
				CollectionRelationships: QueryRequestCollectionRelationships{},
				Query: Query{
					Fields: QueryFields{
						"location": NewColumnField("location", NewNestedObject(map[string]FieldEncoder{
							"country": NewColumnField("country", nil),
						})).Encode(),
					},
					OrderBy: &OrderBy{
						Elements: []OrderByElement{
							{
								OrderDirection: OrderDirectionAsc,
								Target:         NewOrderByColumn("location", []PathElement{}, []string{"country"}).Encode(),
							},
						},
					},
				},
			},
		},
		{
			"order_by_aggregate",
			[]byte(`{
				"collection": "authors",
				"arguments": {},
				"query": {
						"fields": {
								"articles_aggregate": {
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
						},
						"order_by": {
								"elements": [
										{
												"order_direction": "desc",
												"target": {
														"type": "star_count_aggregate",
														"path": [
																{
																		"arguments": {},
																		"relationship": "author_articles",
																		"predicate": {
																				"type": "and",
																				"expressions": []
																		}
																}
														]
												}
										}
								]
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
				Arguments:  QueryRequestArguments{},
				Query: Query{
					Fields: QueryFields{
						"articles_aggregate": NewRelationshipField(Query{
							Aggregates: QueryAggregates{
								"count": NewAggregateStarCount().Encode(),
							},
						}, "author_articles", map[string]RelationshipArgument{}).Encode(),
					},
					OrderBy: &OrderBy{
						Elements: []OrderByElement{
							{
								OrderDirection: OrderDirectionDesc,
								Target: NewOrderByStarCountAggregate([]PathElement{
									{
										Arguments:    PathElementArguments{},
										Relationship: "author_articles",
										Predicate:    NewExpressionAnd().Encode(),
									},
								}).Encode(),
							},
						},
					},
				},
				CollectionRelationships: QueryRequestCollectionRelationships{
					"author_articles": Relationship{
						Arguments: RelationshipArguments{},
						ColumnMapping: RelationshipColumnMapping{
							"id": "author_id",
						},
						RelationshipType: RelationshipTypeArray,
						TargetCollection: "articles",
					},
				},
			},
		},

		{
			"order_by_aggregate_function",
			[]byte(`{
				"collection": "authors",
				"arguments": {},
				"query": {
						"fields": {
								"articles_aggregate": {
										"type": "relationship",
										"arguments": {},
										"relationship": "author_articles",
										"query": {
												"aggregates": {
														"max_id": {
																"type": "single_column",
																"column": "id",
																"function": "max"
														}
												}
										}
								}
						},
						"order_by": {
								"elements": [
										{
												"order_direction": "asc",
												"target": {
														"type": "single_column_aggregate",
														"column": "id",
														"function": "max",
														"path": [
																{
																		"arguments": {},
																		"relationship": "author_articles",
																		"predicate": {
																				"type": "and",
																				"expressions": []
																		}
																}
														],
														"field_path": []
												}
										}
								]
						}
				},
				"collection_relationships": {
						"author_articles": {
								"arguments": {},
								"column_mapping": {
										"id": "author_id"
								},
								"relationship_type": "array",
								"source_collection_or_type": "author",
								"target_collection": "articles"
						}
				}
		}`),
			QueryRequest{
				Collection: "authors",
				Arguments:  QueryRequestArguments{},
				Query: Query{
					Fields: QueryFields{
						"articles_aggregate": NewRelationshipField(Query{
							Aggregates: QueryAggregates{
								"max_id": NewAggregateSingleColumn("id", "max", nil).Encode(),
							},
						}, "author_articles", map[string]RelationshipArgument{}).Encode(),
					},
					OrderBy: &OrderBy{
						Elements: []OrderByElement{
							{
								OrderDirection: OrderDirectionAsc,
								Target: NewOrderBySingleColumnAggregate("id", "max", []PathElement{
									{
										Arguments:    PathElementArguments{},
										Relationship: "author_articles",
										Predicate:    NewExpressionAnd().Encode(),
									},
								}, []string{}).Encode(),
							},
						},
					},
				},
				CollectionRelationships: QueryRequestCollectionRelationships{
					"author_articles": Relationship{
						Arguments: RelationshipArguments{},
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
