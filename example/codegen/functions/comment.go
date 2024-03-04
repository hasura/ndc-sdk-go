package functions

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/utils"
)

// @scalar CommentString
type CommentText struct {
	comment string
}

func (ct *CommentText) FromValue(value any) (err error) {
	ct.comment, err = utils.DecodeString(value)
	return
}

type GetArticlesArguments struct {
	Uint        uint
	Limit       float64
	Comment     CommentText
	NullableStr *string `json:"nullable_str"`
}

type GetArticlesResult struct {
	ID     string `json:"id"`
	Name   Text
	Author struct {
		ID        uuid.UUID  `json:"id"`
		Decimal   complex128 `json:"decimal"`
		CreatedAt time.Time  `json:"created_at"`
	} `json:"author"`
	Location *struct {
		Long int
		Lat  int
	}
	Comments []struct {
		Content string `json:"content"`
	}
}

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
		ID        uuid.UUID  `json:"id"`
		Decimal   complex128 `json:"decimal"`
		CreatedAt time.Time  `json:"created_at"`
	} `json:"author"`
}

type Author struct {
	ID        string        `json:"id"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
}

type CreateArticleResult struct {
	ID      uint     `json:"id"`
	Authors []Author `json:"authors"`
}

// @procedure create_article
func CreateArticle(ctx context.Context, state *types.State, arguments *CreateArticleArguments) (*CreateArticleResult, error) {
	return &CreateArticleResult{
		ID:      1,
		Authors: []Author{},
	}, nil
}

// @procedure
func Increase(ctx context.Context, state *types.State) (int, error) {
	return 0, nil
}
