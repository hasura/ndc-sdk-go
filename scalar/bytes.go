package scalar

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hasura/ndc-sdk-go/utils"
)

// Bytes wraps the scalar implementation for bytes with base64 encoding
//
// @scalar Bytes string.
type Bytes struct {
	data []byte
}

// NewBytes create a Bytes instance.
func NewBytes(data []byte) *Bytes {
	return &Bytes{data: data}
}

// ParseBytes parses bytes from base64 string.
func ParseBytes(value string) (*Bytes, error) {
	bs, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return &Bytes{data: bs}, nil
}

// ScalarName get the schema name of the scalar.
func (bs Bytes) ScalarName() string {
	return "Bytes"
}

// Bytes get the inner bytes value.
func (bs Bytes) Bytes() []byte {
	return bs.data
}

// Bytes get the inner bytes length.
func (bs Bytes) Len() int {
	return len(bs.data)
}

// String implements fmt.Stringer interface.
func (bs Bytes) String() string {
	return string(bs.data)
}

// MarshalJSON implements json.Marshaler.
func (bs Bytes) MarshalJSON() ([]byte, error) {
	str := base64.StdEncoding.EncodeToString(bs.data)

	return json.Marshal(str)
}

// UnmarshalJSON implements json.Unmarshaler.
func (bs *Bytes) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	result, err := ParseBytes(value)
	if err != nil {
		return err
	}

	*bs = *result

	return nil
}

// FromValue decode any value to d Date.
func (bs *Bytes) FromValue(value any) error {
	sValue, err := utils.DecodeNullableString(value)
	if err != nil {
		return err
	}

	if sValue == nil {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(*sValue)
	if err != nil {
		return err
	}

	bs.data = data

	return nil
}
