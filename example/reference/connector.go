package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/swaggest/jsonschema-go"
)

type RawConfiguration struct{}
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

type State struct {
	Authors   map[int]Author
	Articles  map[int]Article
	Telemetry *connector.TelemetryState
}

func (s State) GetLatestArticle() *Article {
	if len(s.Articles) == 0 {
		return nil
	}

	articles := sortArticles(s.Articles, "id", true)

	return &articles[0]
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
func (mc *Connector) TryInitState(configuration *Configuration, metrics *connector.TelemetryState) (*State, error) {
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

	return &State{
		Authors:   authors,
		Articles:  articles,
		Telemetry: metrics,
	}, nil
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
						ArgumentType: schema.NewNamedType("String").Encode(),
					},
				},
			},
			"Int": schema.ScalarType{
				AggregateFunctions: schema.ScalarTypeAggregateFunctions{
					"max": schema.AggregateFunctionDefinition{
						ResultType: schema.NewNullableNamedType("Int").Encode(),
					},
					"min": schema.AggregateFunctionDefinition{
						ResultType: schema.NewNullableNamedType("Int").Encode(),
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
						Type:        schema.NewNamedType("Int").Encode(),
					},
					"title": schema.ObjectField{
						Description: schema.ToPtr("The article's title"),
						Type:        schema.NewNamedType("String").Encode(),
					},
					"author_id": schema.ObjectField{
						Description: schema.ToPtr("The article's author ID"),
						Type:        schema.NewNamedType("Int").Encode(),
					},
				},
			},
			"author": schema.ObjectType{
				Description: schema.ToPtr("An author"),
				Fields: schema.ObjectTypeFields{
					"id": schema.ObjectField{
						Description: schema.ToPtr("The author's primary key"),
						Type:        schema.NewNamedType("Int").Encode(),
					},
					"first_name": schema.ObjectField{
						Description: schema.ToPtr("The author's first name"),
						Type:        schema.NewNamedType("String").Encode(),
					},
					"last_name": schema.ObjectField{
						Description: schema.ToPtr("The author's last name"),
						Type:        schema.NewNamedType("String").Encode(),
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
						Type: schema.NewNamedType("Int").Encode(),
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
				ResultType:  schema.NewNullableNamedType("Int").Encode(),
			},
		},
		Procedures: []schema.ProcedureInfo{
			{
				Name:        "upsert_article",
				Description: schema.ToPtr("Insert or update an article"),
				Arguments: schema.ProcedureInfoArguments{
					"article": schema.ArgumentInfo{
						Description: schema.ToPtr("The article to insert or update"),
						Type:        schema.NewNamedType("article").Encode(),
					},
				},
				ResultType: schema.NewNullableNamedType("article").Encode(),
			},
		},
	}, nil
}

func (mc *Connector) Explain(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{
		Details: schema.ExplainResponseDetails{},
	}, nil
}

func (mc *Connector) Query(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (schema.QueryResponse, error) {
	var rows []schema.Row
	switch request.Collection {
	case "articles":
		rows = schema.ToRows(getMapValues(state.Articles))
	case "authors":
		rows = schema.ToRows(getMapValues(state.Authors))
	case "articles_by_author":
		authorIdArg, ok := request.Arguments["author_id"]
		if !ok {
			return nil, schema.BadGatewayError("missing argument author_id", nil)
		}

		for _, row := range state.Articles {
			switch authorIdArg.Type {
			case schema.ArgumentTypeLiteral:
				if fmt.Sprint(row.AuthorID) == fmt.Sprint(authorIdArg.Value) {
					rows = append(rows, row)
				}
			}
		}
	case "latest_article_id":
		latestArticle := state.GetLatestArticle()
		if latestArticle == nil {
			return nil, schema.BadRequestError("No available article", nil)
		}

		rows = []schema.Row{
			map[string]any{
				"__value": latestArticle.ID,
			},
		}
	default:
		return nil, schema.BadRequestError(fmt.Sprintf("invalid collection name %s", request.Collection), nil)
	}

	return schema.QueryResponse{
		{
			Rows:       rows,
			Aggregates: schema.RowSetAggregates{},
		},
	}, nil
}

func (mc *Connector) Mutation(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.MutationResponse, error) {

	operationResults := []schema.MutationOperationResults{}
	for _, operation := range request.Operations {
		results, err := executeMutationOperation(ctx, state, request.CollectionRelationships, &operation)
		if err != nil {
			return nil, err
		}
		operationResults = append(operationResults, *results)
	}

	return &schema.MutationResponse{
		OperationResults: operationResults,
	}, nil
}

func executeMutationOperation(ctx context.Context, state *State, collectionRelationship schema.MutationRequestCollectionRelationships, operation *schema.MutationOperation) (*schema.MutationOperationResults, error) {
	switch operation.Type {
	case schema.MutationOperationProcedure:
		return executeProcedure(ctx, state, collectionRelationship, operation)
	}

	return nil, schema.NotSupportedError(fmt.Sprintf("Unsupported operation type: %s", operation.Type), nil)
}

type UpsertArticleArguments struct {
	Article Article `json:"article"`
}

func executeProcedure(ctx context.Context, state *State, collectionRelationship schema.MutationRequestCollectionRelationships, operation *schema.MutationOperation) (*schema.MutationOperationResults, error) {
	switch operation.Name {
	case "upsert_article":
		var args UpsertArticleArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.BadRequestError(err.Error(), nil)
		}

		latestArticle := state.GetLatestArticle()
		if args.Article.ID <= 0 {
			if latestArticle == nil {
				args.Article.ID = 1
			} else {
				args.Article.ID = latestArticle.ID + 1
			}
		}
		state.Articles[args.Article.ID] = args.Article

		return &schema.MutationOperationResults{
			AffectedRows: 1,
			Returning:    []schema.Row{args.Article},
		}, nil
	default:
		return nil, schema.BadRequestError("unknown procedure", nil)
	}
}
