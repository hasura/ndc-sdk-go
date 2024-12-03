package procedure

import (
	"context"

	"github.com/hasura/ndc-codegen-function-only-test/types"
)

func ProcedureCreateDemo(ctx context.Context, state *types.State) (*types.DemographyResult, error) {
	return &types.DemographyResult{}, nil
}
