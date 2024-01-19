package main

import (
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/swaggest/jsonschema-go"
)

//go:embed articles.csv
var csvArticles string

//go:embed authors.csv
var csvAuthors string

type RawConfiguration struct{}
type Configuration struct{}

type State struct {
	Authors  []schema.RowSetRowsElem
	Articles []schema.RowSetRowsElem
}

type Connector struct{}

func (mc *Connector) GetRawConfigurationSchema() *jsonschema.Schema {
	return nil
}
func (mc *Connector) MakeEmptyConfiguration() *RawConfiguration {
	return &RawConfiguration{}
}

func (mc *Connector) UpdateConfiguration(ctx context.Context, rawConfiguration *RawConfiguration) (*RawConfiguration, error) {
	return &RawConfiguration{}, nil
}
func (mc *Connector) ValidateRawConfiguration(rawConfiguration *RawConfiguration) (*Configuration, error) {
	return &Configuration{}, nil
}
func (mc *Connector) TryInitState(configuration *Configuration, metrics any) (*State, error) {
	articles, err := readArticles()

	if err != nil {
		return nil, connector.InternalServerError("failed to read articles from csv", map[string]any{
			"cause": err.Error(),
		})
	}

	authors, err := readAuthors()
	if err != nil {
		return nil, connector.InternalServerError("failed to read authors from csv", map[string]any{
			"cause": err.Error(),
		})
	}

	return &State{
		Authors:  authors,
		Articles: articles,
	}, nil
}

func (mc *Connector) FetchMetrics(ctx context.Context, configuration *Configuration, state *State) error {
	return nil
}

func (mc *Connector) HealthCheck(ctx context.Context, configuration *Configuration, state *State) error {
	return nil
}

func (mc *Connector) GetCapabilities(configuration *Configuration) *schema.CapabilitiesResponse {
	return &schema.CapabilitiesResponse{
		Versions: "^0.1.0",
		Capabilities: schema.Capabilities{
			Query: schema.QueryCapabilities{
				Aggregates: schema.LeafCapability{},
				Variables:  schema.LeafCapability{},
			},
			Relationships: schema.RelationshipCapabilities{
				OrderByAggregate:    schema.LeafCapability{},
				RelationComparisons: schema.LeafCapability{},
			},
		},
	}
}

func (mc *Connector) GetSchema(configuration *Configuration) (*schema.SchemaResponse, error) {
	return &schema.SchemaResponse{
		ScalarTypes: schema.SchemaResponseScalarTypes{
			"String": schema.ScalarType{
				AggregateFunctions: schema.ScalarTypeAggregateFunctions{},
				ComparisonOperators: schema.ScalarTypeComparisonOperators{
					"like": schema.ComparisonOperatorDefinition{
						ArgumentType: schema.NewNamedType("String"),
					},
				},
			},
			"Int": schema.ScalarType{
				AggregateFunctions: schema.ScalarTypeAggregateFunctions{
					"max": schema.AggregateFunctionDefinition{
						ResultType: schema.NewNullableNamedType("Int"),
					},
					"min": schema.AggregateFunctionDefinition{
						ResultType: schema.NewNullableNamedType("Int"),
					},
				},
			},
		},
		ObjectTypes: schema.SchemaResponseObjectTypes{
			"article": schema.ObjectType{
				Description: schema.ToPtr("An article"),
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Description: schema.ToPtr("The article's primary key"),
						Type:        schema.NewNamedType("Int"),
					},
					"title": schema.ObjectField{
						Description: schema.ToPtr("The article's title"),
						Type:        schema.NewNamedType("String"),
					},
					"author_id": schema.ObjectField{
						Description: schema.ToPtr("The article's author ID"),
						Type:        schema.NewNamedType("Int"),
					},
				},
			},
			"author": schema.ObjectType{
				Description: schema.ToPtr("An author"),
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Description: schema.ToPtr("The author's primary key"),
						Type:        schema.NewNamedType("Int"),
					},
					"first_name": schema.ObjectField{
						Description: schema.ToPtr("The author's first name"),
						Type:        schema.NewNamedType("String"),
					},
					"last_name": schema.ObjectField{
						Description: schema.ToPtr("The author's last name"),
						Type:        schema.NewNamedType("String"),
					},
				},
			},
		},
		Collections: []schema.CollectionInfo{
			{
				Name:        "articles",
				Description: schema.ToPtr("A collection of articles"),
				ForeignKeys: schema.CollectionInfoForeignKeys{
					"Article_AuthorID": schema.ForeignKeyConstraint{
						ForeignCollection: "authors",
						ColumnMapping: schema.ForeignKeyConstraintColumnMapping{
							"author_id": "id",
						},
					},
				},
				UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
					"ArticleByID": schema.UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
			},
			{
				Name:        "authors",
				Description: schema.ToPtr("A collection of authors"),
				ForeignKeys: schema.CollectionInfoForeignKeys{},
				UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
					"AuthorByID": schema.UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
			},
			{
				Name:        "articles_by_author",
				Description: schema.ToPtr("Articles parameterized by author"),
				Arguments: schema.CollectionInfoArguments{
					"author_id": schema.ArgumentInfo{
						Type: schema.NewNamedType("Int"),
					},
				},
				ForeignKeys: schema.CollectionInfoForeignKeys{},
				UniquenessConstraints: schema.CollectionInfoUniquenessConstraints{
					"AuthorByID": schema.UniquenessConstraint{
						UniqueColumns: []string{"id"},
					},
				},
			},
		},
		Functions: []schema.FunctionInfo{
			{
				Name:        "latest_article_id",
				Description: schema.ToPtr("Get the ID of the most recent article"),
				ResultType:  schema.NewNullableNamedType("Int"),
			},
		},
		Procedures: []schema.ProcedureInfo{
			{
				Name:        "upsert_article",
				Description: schema.ToPtr("Insert or update an article"),
				Arguments: schema.ProcedureInfoArguments{
					"article": schema.ArgumentInfo{
						Description: schema.ToPtr("The article to insert or update"),
						Type:        schema.NewNamedType("article"),
					},
				},
				ResultType: schema.NewNullableNamedType("article"),
			},
		},
	}, nil
}

func (mc *Connector) Explain(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}
func (mc *Connector) Mutation(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	return &schema.MutationResponse{
		OperationResults: []schema.MutationOperationResults{},
	}, nil
}

func (mc *Connector) Query(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.QueryResponse, error) {

	var rows []schema.RowSetRowsElem
	switch request.Collection {
	case "articles":
		rows = state.Articles
		break
	case "authors":
		rows = state.Authors
		break
	case "articles_by_author":
		authorIdArg, ok := request.Arguments["author_id"]
		if !ok {
			return nil, connector.BadGatewayError("missing argument author_id", nil)
		}

		for _, row := range state.Articles {
			switch authorIdArg.Type {
			case schema.ArgumentLiteral:
				// if row["author_id"] == authorIdArg.Value {
				rows = append(rows, row)
				// }
				break
			}
		}
		break
	case "latest_article_id":
		rows = []schema.RowSetRowsElem{
			{
				"__value": state.Articles[len(state.Articles)-1]["id"],
			},
		}
		break
	default:
		return nil, connector.BadRequestError(fmt.Sprintf("invalid collection name %s", request.Collection), nil)
	}

	return &schema.QueryResponse{
		{
			Rows:       rows,
			Aggregates: schema.RowSetAggregates{},
		},
	}, nil
}

func readArticles() ([]schema.RowSetRowsElem, error) {
	r := csv.NewReader(strings.NewReader(csvArticles))
	var results []schema.RowSetRowsElem
	// skip the title row
	_, err := r.Read()
	if err == io.EOF {
		return results, nil
	}
	if err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return nil, err
		}
		authorID, err := strconv.ParseInt(record[2], 10, 32)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]any{
			"id":        int(id),
			"title":     record[1],
			"author_id": int(authorID),
		})
	}

	return results, nil
}

func readAuthors() ([]schema.RowSetRowsElem, error) {
	r := csv.NewReader(strings.NewReader(csvArticles))
	var results []schema.RowSetRowsElem
	// skip the title row
	_, err := r.Read()
	if err == io.EOF {
		return results, nil
	}
	if err != nil {
		return nil, err
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id, err := strconv.ParseInt(record[0], 10, 32)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]any{
			"id":         int(id),
			"first_name": record[1],
			"last_name":  record[2],
		})
	}

	return results, nil
}
