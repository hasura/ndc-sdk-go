package functions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-codegen-test/types"
	"github.com/hasura/ndc-sdk-go/utils"
)

// CommentText
// @scalar CommentString
type CommentText struct {
	comment string
}

func (c CommentText) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.comment)
}

func (c *CommentText) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	c.comment = s

	return nil
}

func (ct *CommentText) FromValue(value any) (err error) {
	ct.comment, err = utils.DecodeString(value)
	return
}

// SomeEnum
// @enum foo, bar
type SomeEnum string

type GetArticlesArguments struct {
	Limit float64
}

type GetArticlesResult struct {
	ID   string `json:"id"`
	Name Text
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

// Increase
// @procedure
func Increase(ctx context.Context, state *types.State) (int, error) {
	return 1, nil
}
