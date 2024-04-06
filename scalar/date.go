package scalar

import (
	"encoding/json"
	"time"

	"github.com/hasura/ndc-sdk-go/utils"
)

const (
	dateFormat = "2006-01-02"
)

// Date wraps the scalar implementation for date representation string
//
// @scalar Date date
type Date struct {
	time.Time
}

// Stringer implements fmt.Stringer interface.
func (d Date) String() string {
	return d.Format(dateFormat)
}

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	date, err := time.Parse(dateFormat, value)
	if err != nil {
		return err
	}
	d.Time = date

	return nil
}

// FromValue decode any value to d Date.
func (d *Date) FromValue(value any) error {
	sValue, err := utils.DecodeNullableString(value)
	if err != nil {
		return err
	}
	if sValue == nil {
		return nil
	}

	date, err := time.Parse(dateFormat, *sValue)
	if err != nil {
		return err
	}
	d.Time = date

	return nil
}
