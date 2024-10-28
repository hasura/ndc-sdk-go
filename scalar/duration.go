package scalar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hasura/ndc-sdk-go/utils"
)

// Decimal wraps the scalar implementation for duration,
// with string or unix time integer in nanoseconds
//
// @scalar Duration
type Duration struct {
	time.Duration
}

// NewDuration creates a Duration instance
func NewDuration(value time.Duration) Duration {
	return Duration{
		Duration: value,
	}
}

// ScalarName get the schema name of the scalar.
func (d Duration) ScalarName() string {
	return "Duration"
}

// Stringer implements fmt.Stringer interface.
func (d Duration) String() string {
	return d.Duration.String()
}

// MarshalJSON implements json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var value any
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("failed to unmarshal duration: %w", err)
	}

	return d.FromValue(value)
}

// FromValue decode any value to Duration.
func (d *Duration) FromValue(value any) error {
	result, err := utils.DecodeNullableDuration(value)
	if err != nil {
		return fmt.Errorf("invalid duration value: %w", err)
	}
	if result != nil {
		d.Duration = *result
	}

	return nil
}
