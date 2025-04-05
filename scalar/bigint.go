package scalar

import (
	"encoding/json"
	"strconv"

	"github.com/hasura/ndc-sdk-go/utils"
)

// BigInt wraps the scalar implementation for big integer,
// with string representation
//
// @scalar BigInt string.
type BigInt int64

// NewBigInt creates a BigInt instance.
func NewBigInt(value int64) BigInt {
	return BigInt(value)
}

// ScalarName get the schema name of the scalar.
func (bi BigInt) ScalarName() string {
	return "BigInt"
}

// Stringer implements fmt.Stringer interface.
func (bi BigInt) String() string {
	return strconv.FormatInt(int64(bi), 10)
}

// MarshalJSON implements json.Marshaler.
func (bi BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(bi), 10))
}

// UnmarshalJSON implements json.Unmarshaler.
func (bi *BigInt) UnmarshalJSON(b []byte) error {
	var value any
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	iValue, err := utils.DecodeNullableInt[int64](value)
	if err != nil {
		return err
	}

	if iValue == nil {
		return nil
	}

	*bi = BigInt(*iValue)

	return nil
}

// FromValue decode any value to Int64.
func (bi *BigInt) FromValue(value any) error {
	iValue, err := utils.DecodeNullableInt[int64](value)
	if err != nil {
		return err
	}

	if iValue != nil {
		*bi = BigInt(*iValue)
	}

	return nil
}
