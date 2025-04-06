package functions

import (
	"context"
	"math"
	"time"

	"github.com/hasura/ndc-codegen-example/types"
)

type GetArticlesArguments struct {
	BaseAuthor

	AuthorID int
	Limit    float64
}

type GetArticlesResult struct {
	ID   uint `json:"id"`
	Name types.Text
}

// GetArticles
// @function
func GetArticles(
	ctx context.Context,
	state *types.State,
	arguments *GetArticlesArguments,
) ([]GetArticlesResult, error) {
	if arguments.AuthorID > 0 {
		time.Sleep(time.Duration(math.Min(float64(arguments.AuthorID), 2)) * time.Second)
	}

	return []GetArticlesResult{
		{
			ID:   uint(arguments.AuthorID),
			Name: types.Text(arguments.Name),
		},
	}, nil
}

type CreateArticleArguments struct {
	Author struct {
		ID        int       `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"author"`
}

type CreateArticleResult struct {
	ID      uint           `json:"id"`
	Authors []types.Author `json:"authors"`
}

// CreateArticle
// @procedure create_article
func CreateArticle(
	ctx context.Context,
	state *types.State,
	arguments *CreateArticleArguments,
) (CreateArticleResult, error) {
	if arguments.Author.ID > 0 {
		time.Sleep(
			time.Duration(math.Min(float64(arguments.Author.ID), 2)) * 500 * time.Millisecond,
		)
	}

	return CreateArticleResult{
		ID:      uint(arguments.Author.ID),
		Authors: []types.Author{},
	}, nil
}

// Increase
// @procedure
func Increase(ctx context.Context, state *types.State) (int, error) {
	return 1, nil
}
