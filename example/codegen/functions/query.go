package functions

import (
	"context"

	"github.com/hasura/ndc-sdk-go/example/codegen/types"
)

// A hello argument
type HelloArguments struct {
	Num int    `json:"num"`
	Str string `json:"str"`
}

// A hello result
type HelloResult struct {
	Num int    `json:"num"`
	Str string `json:"str"`
}

func FunctionGetScalar(ctx context.Context, state *types.State, arguments *HelloArguments) (int, error) {
	return 0, nil
}

// Send hello message
func FunctionHello(ctx context.Context, state *types.State, arguments *HelloArguments) (*HelloResult, error) {
	return &HelloResult{
		Num: 1,
		Str: "world",
	}, nil
}

// A create author argument
type CreateAuthorArguments struct {
	Name string `json:"name"`
}

// A create author result
type CreateAuthorResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// A procedure example
func ProcedureCreateAuthor(ctx context.Context, state *types.State, arguments *CreateAuthorArguments) (*CreateAuthorResult, error) {
	return &CreateAuthorResult{
		ID:   1,
		Name: arguments.Name,
	}, nil
}
