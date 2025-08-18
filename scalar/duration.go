package scalar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hasura/ndc-sdk-go/v2/utils"
	"github.com/prometheus/common/model"
)

// @scalar Duration.
type Duration struct {
	time.Duration
}

// NewDuration creates a Duration instance.
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

// DurationString wraps the scalar implementation for duration with string.
type DurationString struct {
	time.Duration
}

// NewDurationString creates a DurationString instance.
func NewDurationString(value time.Duration) DurationString {
	return DurationString{
		Duration: value,
	}
}

// ScalarName get the schema name of the scalar.
func (d DurationString) ScalarName() string {
	return "DurationString"
}

// Stringer implements fmt.Stringer interface.
func (d DurationString) String() string {
	return d.Duration.String()
}

// MarshalJSON implements json.Marshaler.
func (d DurationString) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DurationString) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("failed to unmarshal duration from string: %w", err)
	}

	dur, err := model.ParseDuration(value)
	if err != nil {
		return err
	}

	d.Duration = time.Duration(dur)

	return nil
}

// FromValue decode any value to Duration.
func (d *DurationString) FromValue(value any) error {
	result, err := utils.DecodeNullableDuration(value)
	if err != nil {
		return fmt.Errorf("invalid duration value: %w", err)
	}

	if result != nil {
		d.Duration = *result
	}

	return nil
}

// DurationInt64 wraps the scalar implementation for duration with int64.
type DurationInt64 struct {
	time.Duration
}

// NewDurationInt64 creates a DurationInt64 instance.
func NewDurationInt64(value time.Duration) DurationInt64 {
	return DurationInt64{
		Duration: value,
	}
}

// ScalarName get the schema name of the scalar.
func (d DurationInt64) ScalarName() string {
	return "DurationInt64"
}

// Stringer implements fmt.Stringer interface.
func (d DurationInt64) String() string {
	return d.Duration.String()
}

// MarshalJSON implements json.Marshaler.
func (d DurationInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(d.Duration))
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DurationInt64) UnmarshalJSON(b []byte) error {
	var value int64
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("failed to unmarshal duration from string: %w", err)
	}

	d.Duration = time.Duration(value)

	return nil
}

// FromValue decode any value to Duration.
func (d *DurationInt64) FromValue(value any) error {
	result, err := utils.DecodeNullableDuration(value)
	if err != nil {
		return fmt.Errorf("invalid duration value: %w", err)
	}

	if result != nil {
		d.Duration = *result
	}

	return nil
}
