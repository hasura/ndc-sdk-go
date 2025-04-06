package functions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-codegen-example/types/arguments"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
)

// A foo scalar
type ScalarFoo struct {
	bar string
}

func (c ScalarFoo) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.bar)
}

// FromValue decodes the scalar from an unknown value
func (s *ScalarFoo) FromValue(value any) error {
	s.bar = fmt.Sprint(value)
	return nil
}

// A hello result
type HelloResult struct {
	ID    uuid.UUID  `json:"id"`
	Num   int        `json:"num"`
	Text  types.Text `json:"text"`
	Foo   ScalarFoo  `json:"foo"`
	Error error      `json:"error"`
}

// FunctionHello sends a hello message
func FunctionHello(ctx context.Context, state *types.State) (*HelloResult, error) {
	return &HelloResult{
		ID:    uuid.MustParse("c67407a0-a825-49e3-9d4a-8265df1490a7"),
		Num:   1,
		Text:  "world",
		Error: errors.New("unknown error"),
	}, nil
}

// A create author argument
type CreateAuthorArguments struct {
	BaseAuthor

	Where schema.Expression `json:"where" ndc:"predicate=Author"`
}

type CreateAuthorsArguments struct {
	Authors []CreateAuthorArguments
}

// A create author result
type CreateAuthorResult struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	Where     schema.Expression `json:"where"`
}

// ProcedureCreateAuthor creates an author
func ProcedureCreateAuthor(
	ctx context.Context,
	state *types.State,
	arguments *CreateAuthorArguments,
) (*CreateAuthorResult, error) {
	selection := utils.CommandSelectionFieldFromContext(ctx)
	if len(selection) == 0 {
		return nil, errors.New("expected not-null selection field, got null")
	}

	return &CreateAuthorResult{
		ID:    1,
		Name:  arguments.Name,
		Where: arguments.Where,
	}, nil
}

// ProcedureCreateAuthors creates a list of authors
func ProcedureCreateAuthors(
	ctx context.Context,
	state *types.State,
	arguments *CreateAuthorsArguments,
) ([]CreateAuthorResult, error) {
	results := make([]CreateAuthorResult, len(arguments.Authors))

	for i, arg := range arguments.Authors {
		results[i] = CreateAuthorResult{
			ID:   i,
			Name: arg.Name,
		}
	}

	return results, nil
}

// FunctionGetBool return an scalar boolean
func FunctionGetBool(ctx context.Context, state *types.State) (bool, error) {
	return true, nil
}

// FunctionGetInts return a slice of scalar ints
func FunctionGetInts(ctx context.Context, state *types.State) ([]*int, error) {
	return []*int{utils.ToPtr(1), utils.ToPtr(2), utils.ToPtr(3)}, nil
}

func FunctionGetTypes(
	ctx context.Context,
	state *types.State,
	arguments *arguments.GetTypesArguments,
) (*arguments.GetTypesArguments, error) {
	selection := utils.CommandSelectionFieldFromContext(ctx)
	if len(selection) == 0 {
		return nil, errors.New("expected not-null selection field, got null")
	}

	return arguments, nil
}

type BaseAuthor struct {
	Name string `json:"name"`
}

type GetAuthorArguments struct {
	*BaseAuthor

	ID    string            `json:"id"`
	Where schema.Expression `json:"where" ndc:"predicate=Author"`
}

type GetAuthorResult struct {
	*CreateAuthorResult

	Where    schema.Expression `json:"where"`
	Disabled bool              `json:"disabled"`
}

func FunctionGetAuthor(
	ctx context.Context,
	state *types.State,
	arguments *GetAuthorArguments,
) (*GetAuthorResult, error) {
	return &GetAuthorResult{
		CreateAuthorResult: &CreateAuthorResult{
			ID:   1,
			Name: arguments.Name,
		},
		Where:    arguments.Where,
		Disabled: false,
	}, nil
}

func FunctionGetAuthor2(
	ctx context.Context,
	state *types.State,
	arguments *GetAuthorArguments,
) (types.GetAuthorResult, error) {
	return types.GetAuthorResult{
		ID:   1,
		Name: arguments.Name,
	}, nil
}

// generate schema and methods for generic types
type GetCustomHeadersArguments[T any, S int | int8 | int16 | int32 | int64] struct {
	Headers map[string]string         `json:"headers"`
	Input   *T                        `json:"input"`
	Other   *GetCustomHeadersOther[S] `json:"other"`
}

type GetCustomHeadersOther[T int | int8 | int16 | int32 | int64] struct {
	Value T `json:"value"`
}

type GetCustomHeadersResult[T any, O int | int8 | int16 | int32 | int64] struct {
	Headers  map[string]string `json:"headers"`
	Response T
	Other    *GetCustomHeadersOther[O] `json:"other"`
}

func FunctionGetCustomHeaders(
	ctx context.Context,
	state *types.State,
	arguments *GetCustomHeadersArguments[BaseAuthor, int],
) (GetCustomHeadersResult[HelloResult, int64], error) {
	if arguments.Headers == nil {
		arguments.Headers = make(map[string]string)
	}

	arguments.Headers["X-Test-ResponseHeader"] = arguments.Input.Name
	result := HelloResult{}

	if arguments.Input != nil {
		result.Text = types.Text(arguments.Input.Name)
	}

	if arguments.Other != nil {
		result.Num = arguments.Other.Value
	}

	return GetCustomHeadersResult[HelloResult, int64]{
		Headers:  arguments.Headers,
		Response: result,
	}, nil
}

// the generic type's method doesn't allow type parameters from other packages
// so the FromValue method isn't generated
func FunctionGetGenericWithoutDecodingMethod(
	ctx context.Context,
	state *types.State,
	arguments *GetCustomHeadersArguments[arguments.GetCustomHeadersInput, int],
) (GetCustomHeadersResult[HelloResult, int64], error) {
	if arguments.Headers == nil {
		arguments.Headers = make(map[string]string)
	}

	arguments.Headers["X-Test-ResponseHeader"] = "I set this in the code"
	result := HelloResult{}

	if arguments.Input != nil {
		result.ID = arguments.Input.ID
		result.Num = arguments.Input.Num
	}

	return GetCustomHeadersResult[HelloResult, int64]{
		Headers:  arguments.Headers,
		Response: result,
	}, nil
}

func ProcedureDoCustomHeaders(
	ctx context.Context,
	state *types.State,
	arguments *GetCustomHeadersArguments[*[]BaseAuthor, int],
) (*types.CustomHeadersResult[[]*BaseAuthor], error) {
	resp := []*BaseAuthor{}

	if arguments.Input != nil && *arguments.Input != nil {
		for _, v := range **arguments.Input {
			resp = append(resp, &v)
		}
	}

	result := &types.CustomHeadersResult[[]*BaseAuthor]{
		Headers:  arguments.Headers,
		Response: resp,
	}

	if result.Headers == nil {
		result.Headers = map[string]string{}
	}

	result.Headers["X-Test-ResponseHeader"] = "I set this in the code"

	return result, nil
}
