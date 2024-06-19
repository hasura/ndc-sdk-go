package functions

import (
	"context"

	"github.com/hasura/ndc-codegen-subdir-test/types"
)

type GetArticlesArguments struct {
	Limit float64
}

type GetArticlesResult struct {
	ID string `json:"id"`
}

// GetArticles
// @function
func GetArticles(ctx context.Context, state *types.State, arguments *GetArticlesArguments) ([]GetArticlesResult, error) {
	return []GetArticlesResult{
		{
			ID:   "1",
			Name: "Article 1",
		},
	}, nil
}
