package scalar

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/hasura/ndc-sdk-go/utils"
)

// URL represents a URL string
//
// @scalar URL string.
type URL struct {
	*url.URL
}

// NewURL creates a URL instance.
func NewURL(rawURL string) (*URL, error) {
	u, err := parseURL(rawURL)
	if err != nil {
		return nil, err
	}
	return &URL{u}, nil
}

// ScalarName get the schema name of the scalar.
func (j URL) ScalarName() string {
	return "URL"
}

// Stringer implements fmt.Stringer interface.
func (u URL) String() string {
	if u.URL == nil {
		return ""
	}
	return u.URL.String()
}

// MarshalJSON implements json.Marshaler.
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *URL) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	value, err := parseURL(str)
	if err != nil {
		return err
	}
	u.URL = value

	return nil
}

// FromValue decode any value to URL.
func (u *URL) FromValue(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case URL:
		*u = v
	case *URL:
		if v != nil {
			*u = *v
		}
	case url.URL:
		u.URL = &v
	case *url.URL:
		u.URL = v
	default:
		str, err := utils.DecodeNullableString(value)
		if err == nil && str != nil {
			u.URL, err = parseURL(*str)
		}
		if err != nil {
			return fmt.Errorf("invalid url: %v", v)
		}
	}
	return nil
}

func parseURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, errors.New("invalid URL")
	}
	return url.Parse(rawURL)
}
