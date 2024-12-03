package function

import (
	"context"

	"github.com/hasura/ndc-codegen-function-only-test/types"
)

type SimpleResult struct {
	Reply string `json:"reply"`
}

func FunctionSimpleObject(ctx context.Context, state *types.State) (SimpleResult, error) {
	return SimpleResult{
		Reply: "Hello world",
	}, nil
}
