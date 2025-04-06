package article

import (
	"context"
	"time"

	"github.com/hasura/ndc-codegen-duplicated-proc/types"
)

type CreateArticleArguments struct {
	Author struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"author"`
}

type Author struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateArticleResult struct {
	ID      uint     `json:"id"`
	Authors []Author `json:"authors"`
}

// CreateArticle
// @procedure create_article
func CreateArticle(ctx context.Context, state *types.State, arguments *CreateArticleArguments) (*CreateArticleResult, error) {
	return &CreateArticleResult{
		ID:      1,
		Authors: []Author{},
	}, nil
}
