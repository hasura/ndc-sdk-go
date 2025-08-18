package types

import (
	"encoding/json"
	"time"

	"github.com/hasura/ndc-sdk-go/v2/utils"
)

type Text string

// CommentText
// @scalar CommentString string
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

// AuthorStatus
// @enum active, inactive
type AuthorStatus string

type Author struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Tags      []string      `json:"tags,omitempty"`
	Status    *AuthorStatus `json:"status"`
	Author    *Author       `json:"author"`
}

type GetAuthorResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CustomHeadersResult[T any] struct {
	Headers  map[string]string `json:"headers,omitempty"`
	Response T
}
