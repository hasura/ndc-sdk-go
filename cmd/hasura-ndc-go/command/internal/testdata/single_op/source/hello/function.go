package hello

import (
	"context"

	"github.com/hasura/ndc-codegen-function-only-test/types"
)

func FunctionHello(ctx context.Context, state *types.State) (string, error) {
	return "test", nil
}
