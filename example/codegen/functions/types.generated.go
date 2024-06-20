// Code generated by github.com/hasura/ndc-sdk-go/codegen, DO NOT EDIT.
package functions

import (
	"github.com/hasura/ndc-sdk-go/utils"
)

var functions_Decoder = utils.NewDecoder()

// FromValue decodes values from map
func (j *GetArticlesArguments) FromValue(input map[string]any) error {
	var err error
	j.Limit, err = utils.GetFloat[float64](input, "Limit")
	if err != nil {
		return err
	}
	return nil
}

// ToMap encodes the struct to a value map
func (j CreateArticleResult) ToMap() map[string]any {
	r := make(map[string]any)
	j_Authors := make([]map[string]any, len(j.Authors))
	for i, j_Authors_v := range j.Authors {
		j_Authors[i] = utils.EncodeMap(j_Authors_v)
	}
	r["authors"] = j_Authors
	r["id"] = j.ID

	return r
}

// ToMap encodes the struct to a value map
func (j CreateAuthorResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["created_at"] = j.CreatedAt
	r["id"] = j.ID
	r["name"] = j.Name

	return r
}

// ToMap encodes the struct to a value map
func (j GetArticlesResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["id"] = j.ID
	r["Name"] = j.Name

	return r
}

// ToMap encodes the struct to a value map
func (j HelloResult) ToMap() map[string]any {
	r := make(map[string]any)
	r["foo"] = j.Foo
	r["id"] = j.ID
	r["num"] = j.Num
	r["text"] = j.Text

	return r
}

// ScalarName get the schema name of the scalar
func (j ScalarFoo) ScalarName() string {
	return "Foo"
}
