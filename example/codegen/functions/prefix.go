package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-codegen-example/types/arguments"
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
		Error: fmt.Errorf("unknown error"),
	}, nil
}

// A create author argument
type CreateAuthorArguments struct {
	BaseAuthor
}
type CreateAuthorsArguments struct {
	Authors []CreateAuthorArguments
}

// A create author result
type CreateAuthorResult struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ProcedureCreateAuthor creates an author
func ProcedureCreateAuthor(ctx context.Context, state *types.State, arguments *CreateAuthorArguments) (*CreateAuthorResult, error) {
	return &CreateAuthorResult{
		ID:   1,
		Name: arguments.Name,
	}, nil
}

// ProcedureCreateAuthors creates a list of authors
func ProcedureCreateAuthors(ctx context.Context, state *types.State, arguments *CreateAuthorsArguments) ([]CreateAuthorResult, error) {
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

func FunctionGetTypes(ctx context.Context, state *types.State, arguments *arguments.GetTypesArguments) (*arguments.GetTypesArguments, error) {
	return arguments, nil
}

type BaseAuthor struct {
	Name string `json:"name"`
}

type GetAuthorArguments struct {
	*BaseAuthor

	ID string `json:"id"`
}

type GetAuthorResult struct {
	*CreateAuthorResult

	Disabled bool `json:"disabled"`
}

func FunctionGetAuthor(ctx context.Context, state *types.State, arguments *GetAuthorArguments) (*GetAuthorResult, error) {
	return &GetAuthorResult{
		CreateAuthorResult: &CreateAuthorResult{
			ID:   1,
			Name: arguments.Name,
		},
		Disabled: false,
	}, nil
}

func FunctionGetAuthor2(ctx context.Context, state *types.State, arguments *GetAuthorArguments) (types.GetAuthorResult, error) {
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

func FunctionGetCustomHeaders(ctx context.Context, state *types.State, arguments *GetCustomHeadersArguments[BaseAuthor, int]) (GetCustomHeadersResult[HelloResult, int64], error) {
	if arguments.Headers == nil {
		arguments.Headers = make(map[string]string)
	}
	arguments.Headers["X-Test-ResponseHeader"] = "I set this in the code"
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
func FunctionGetGenericWithoutDecodingMethod(ctx context.Context, state *types.State, arguments *GetCustomHeadersArguments[arguments.GetCustomHeadersInput, int]) (GetCustomHeadersResult[HelloResult, int64], error) {
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
