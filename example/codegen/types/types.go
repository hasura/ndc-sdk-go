package types

import (
	"encoding/json"
	"time"

	"github.com/hasura/ndc-sdk-go/utils"
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

type Author struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
	Author    *Author   `json:"author"`
}
