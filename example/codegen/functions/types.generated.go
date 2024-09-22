// Code generated by github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go, DO NOT EDIT.
package functions

import (
	"context"
	"encoding/json"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-codegen-example/types/arguments"
	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"slices"
)

var functions_Decoder = utils.NewDecoder()

// FromValue decodes values from map
func (j *GetArticlesArguments) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObject(&j.BaseAuthor, input)
	if err != nil {
		return err
	}
	j.Limit, err = utils.GetFloat[float64](input, "Limit")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *GetAuthorArguments) FromValue(input map[string]any) error {
	var err error
	j.BaseAuthor = new(BaseAuthor)
	err = functions_Decoder.DecodeObject(j.BaseAuthor, input)
	if err != nil {
		return err
	}
	j.ID, err = utils.GetString(input, "id")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *BaseAuthor) FromValue(input map[string]any) error {
	var err error
	j.Name, err = utils.GetString(input, "name")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *CreateArticleResult) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObjectValue(&j.Authors, input, "authors")
	if err != nil {
		return err
	}
	j.ID, err = utils.GetUint[uint](input, "id")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *CreateAuthorArguments) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObject(&j.BaseAuthor, input)
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *CreateAuthorResult) FromValue(input map[string]any) error {
	var err error
	j.CreatedAt, err = utils.GetDateTime(input, "created_at")
	if err != nil {
		return err
	}
	j.ID, err = utils.GetInt[int](input, "id")
	if err != nil {
		return err
	}
	j.Name, err = utils.GetString(input, "name")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *GetArticlesResult) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObjectValue(&j.Name, input, "Name")
	if err != nil {
		return err
	}
	j.ID, err = utils.GetString(input, "id")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *GetAuthorResult) FromValue(input map[string]any) error {
	var err error
	j.CreateAuthorResult = new(CreateAuthorResult)
	err = functions_Decoder.DecodeObject(j.CreateAuthorResult, input)
	if err != nil {
		return err
	}
	j.Disabled, err = utils.GetBoolean(input, "disabled")
	if err != nil {
		return err
	}
	return nil
}

// FromValue decodes values from map
func (j *HelloResult) FromValue(input map[string]any) error {
	var err error
	err = functions_Decoder.DecodeObjectValue(&j.Error, input, "error")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.Foo, input, "foo")
	if err != nil {
		return err
	}
	j.ID, err = utils.GetUUID(input, "id")
	if err != nil {
		return err
	}
	j.Num, err = utils.GetInt[int](input, "num")
	if err != nil {
		return err
	}
	err = functions_Decoder.DecodeObjectValue(&j.Text, input, "text")
	if err != nil {
		return err
	}
	return nil
}

// ToMap encodes the struct to a value map
func (j BaseAuthor) ToMap() map[string]any {
	r := make(map[string]any)
	r["name"] = j.Name

	return r
}

// ToMap encodes the struct to a value map
func (j CreateArticleResult) ToMap() map[string]any {
	r := make(map[string]any)
	j_Authors := make([]any, len(j.Authors))
	for i, j_Authors_v := range j.Authors {
		j_Authors[i] = j_Authors_v
	}
	r["authors"] = j_Authors
	r["id"] = j.ID

	return r
}

// ToMap encodes the struct to a value map
func (j CreateAuthorArguments) ToMap() map[string]any {
	r := make(map[string]any)
	r = utils.MergeMap(r, j.BaseAuthor.ToMap())

	return r
}

// ToMap encodes the struct to a value map
func (j CreateAuthorResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["created_at"] = j.CreatedAt
	r["id"] = j.ID
	r["name"] = j.Name

	return r
}

// ToMap encodes the struct to a value map
func (j GetArticlesResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["Name"] = j.Name
	r["id"] = j.ID

	return r
}

// ToMap encodes the struct to a value map
func (j GetAuthorResult) ToMap() map[string]any {
	r := make(map[string]any)
	if j.CreateAuthorResult != nil {
		r = utils.MergeMap(r, (*j.CreateAuthorResult).ToMap())
	}
	r["disabled"] = j.Disabled

	return r
}

// ToMap encodes the struct to a value map
func (j HelloResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["error"] = j.Error
	r["foo"] = j.Foo
	r["id"] = j.ID
	r["num"] = j.Num
	r["text"] = j.Text

	return r
}

// ScalarName get the schema name of the scalar
func (j ScalarFoo) ScalarName() string {
	return "Foo"
}

// DataConnectorHandler implements the data connector handler
type DataConnectorHandler struct{}

// QueryExists check if the query name exists
func (dch DataConnectorHandler) QueryExists(name string) bool {
	return slices.Contains(enumValues_FunctionName, name)
}
func (dch DataConnectorHandler) Query(ctx context.Context, state *types.State, request *schema.QueryRequest, rawArgs map[string]any) (*schema.RowSet, error) {
	if !dch.QueryExists(request.Collection) {
		return nil, utils.ErrHandlerNotfound
	}
	queryFields, err := utils.EvalFunctionSelectionFieldValue(request)
	if err != nil {
		return nil, schema.UnprocessableContentError(err.Error(), nil)
	}

	result, err := dch.execQuery(ctx, state, request, queryFields, rawArgs)
	if err != nil {
		return nil, err
	}

	return &schema.RowSet{
		Aggregates: schema.RowSetAggregates{},
		Rows: []map[string]any{
			{
				"__value": result,
			},
		},
	}, nil
}

func (dch DataConnectorHandler) execQuery(ctx context.Context, state *types.State, request *schema.QueryRequest, queryFields schema.NestedField, rawArgs map[string]any) (any, error) {
	span := trace.SpanFromContext(ctx)
	logger := connector.GetLogger(ctx)
	switch request.Collection {
	case "getAuthor":

		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args GetAuthorArguments
		if parseErr := args.FromValue(rawArgs); parseErr != nil {
			return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
				"cause": parseErr.Error(),
			})
		}

		connector_addSpanEvent(span, logger, "execute_function", map[string]any{
			"arguments": args,
		})
		rawResult, err := FunctionGetAuthor(ctx, state, &args)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnObject(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil

	case "getBool":

		if len(queryFields) > 0 {
			return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
		}
		return FunctionGetBool(ctx, state)

	case "getTypes":

		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args arguments.GetTypesArguments
		if parseErr := args.FromValue(rawArgs); parseErr != nil {
			return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
				"cause": parseErr.Error(),
			})
		}

		connector_addSpanEvent(span, logger, "execute_function", map[string]any{
			"arguments": args,
		})
		rawResult, err := FunctionGetTypes(ctx, state, &args)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnObject(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil

	case "hello":

		selection, err := queryFields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		rawResult, err := FunctionHello(ctx, state)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}

		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnObject(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil

	case "getArticles":

		selection, err := queryFields.AsArray()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
				"cause": err.Error(),
			})
		}
		var args GetArticlesArguments
		if parseErr := args.FromValue(rawArgs); parseErr != nil {
			return nil, schema.UnprocessableContentError("failed to resolve arguments", map[string]any{
				"cause": parseErr.Error(),
			})
		}

		connector_addSpanEvent(span, logger, "execute_function", map[string]any{
			"arguments": args,
		})
		rawResult, err := GetArticles(ctx, state, &args)
		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, schema.UnprocessableContentError("expected not null result", nil)
		}

		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)
		if err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, utils.ErrHandlerNotfound
	}
}

var enumValues_FunctionName = []string{"getAuthor", "getBool", "getTypes", "hello", "getArticles"}

// MutationExists check if the mutation name exists
func (dch DataConnectorHandler) MutationExists(name string) bool {
	return slices.Contains(enumValues_ProcedureName, name)
}
func (dch DataConnectorHandler) Mutation(ctx context.Context, state *types.State, operation *schema.MutationOperation) (schema.MutationOperationResults, error) {
	span := trace.SpanFromContext(ctx)
	logger := connector.GetLogger(ctx)
	connector_addSpanEvent(span, logger, "validate_request", map[string]any{
		"operations_name": operation.Name,
	})

	switch operation.Name {
	case "create_article":

		selection, err := operation.Fields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args CreateArticleArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.UnprocessableContentError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		span.AddEvent("execute_procedure")
		rawResult, err := CreateArticle(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}
		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnObject(selection, rawResult)

		if err != nil {
			return nil, err
		}
		return schema.NewProcedureResult(result).Encode(), nil

	case "increase":

		if len(operation.Fields) > 0 {
			return nil, schema.UnprocessableContentError("cannot evaluate selection fields for scalar", nil)
		}
		span.AddEvent("execute_procedure")
		result, err := Increase(ctx, state)
		if err != nil {
			return nil, err
		}
		return schema.NewProcedureResult(result).Encode(), nil

	case "createAuthor":

		selection, err := operation.Fields.AsObject()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be object", map[string]any{
				"cause": err.Error(),
			})
		}
		var args CreateAuthorArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.UnprocessableContentError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		span.AddEvent("execute_procedure")
		rawResult, err := ProcedureCreateAuthor(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, nil
		}
		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnObject(selection, rawResult)

		if err != nil {
			return nil, err
		}
		return schema.NewProcedureResult(result).Encode(), nil

	case "createAuthors":

		selection, err := operation.Fields.AsArray()
		if err != nil {
			return nil, schema.UnprocessableContentError("the selection field type must be array", map[string]any{
				"cause": err.Error(),
			})
		}
		var args CreateAuthorsArguments
		if err := json.Unmarshal(operation.Arguments, &args); err != nil {
			return nil, schema.UnprocessableContentError("failed to decode arguments", map[string]any{
				"cause": err.Error(),
			})
		}
		span.AddEvent("execute_procedure")
		rawResult, err := ProcedureCreateAuthors(ctx, state, &args)

		if err != nil {
			return nil, err
		}

		if rawResult == nil {
			return nil, schema.UnprocessableContentError("expected not null result", nil)
		}
		connector_addSpanEvent(span, logger, "evaluate_response_selection", map[string]any{
			"raw_result": rawResult,
		})
		result, err := utils.EvalNestedColumnArrayIntoSlice(selection, rawResult)

		if err != nil {
			return nil, err
		}
		return schema.NewProcedureResult(result).Encode(), nil

	default:
		return nil, utils.ErrHandlerNotfound
	}
}

var enumValues_ProcedureName = []string{"create_article", "increase", "createAuthor", "createAuthors"}

func connector_addSpanEvent(span trace.Span, logger *slog.Logger, name string, data map[string]any, options ...trace.EventOption) {
	logger.Debug(name, slog.Any("data", data))
	attrs := utils.DebugJSONAttributes(data, utils.IsDebug(logger))
	span.AddEvent(name, append(options, trace.WithAttributes(attrs...))...)
}
