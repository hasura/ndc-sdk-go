package functions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
)

type GetArticlesArguments struct {
	BaseAuthor

	Limit float64
}

type GetArticlesResult struct {
	ID   string `json:"id"`
	Name types.Text
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

type CreateArticleArguments struct {
	Author struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"author"`
}

type CreateArticleResult struct {
	ID      uint           `json:"id"`
	Authors []types.Author `json:"authors"`
}

// CreateArticle
// @procedure create_article
func CreateArticle(ctx context.Context, state *types.State, arguments *CreateArticleArguments) (*CreateArticleResult, error) {
	return &CreateArticleResult{
		ID:      1,
		Authors: []types.Author{},
	}, nil
}

// Increase
// @procedure
func Increase(ctx context.Context, state *types.State) (int, error) {
	return 1, nil
}
