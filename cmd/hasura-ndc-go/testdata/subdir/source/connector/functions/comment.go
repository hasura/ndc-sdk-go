package functions

import (
	"context"

	example "github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-codegen-subdir-test/types"
)

type GetArticlesArguments struct {
	Limit float64
	example.Author
}

type GetArticlesResult struct {
	ID string `json:"id"`
	example.Author
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
