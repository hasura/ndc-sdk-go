package article

import (
	"context"

	"github.com/hasura/ndc-codegen-duplicated-func/types"
)

type GetArticlesArguments struct {
	Limit float64
}

type GetArticlesResult struct {
	ID   string `json:"id"`
	Name string
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
