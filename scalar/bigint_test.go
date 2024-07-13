package scalar

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBigInt(t *testing.T) {
	expected := "9999"
	value := NewBigInt(0)
	scalarName := value.ScalarName()

	if err := json.Unmarshal([]byte(expected), &value); err != nil {
		t.Errorf("failed to parse %s: %s", scalarName, err)
		t.FailNow()
	}
	if expected != value.String() {
		t.Errorf("expected: %s, got: %s", expected, value)
		t.FailNow()
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		t.Errorf("failed to encode %s: %s", scalarName, err)
		t.FailNow()
	}

	if expected != strings.Trim(string(rawValue), "\"") {
		t.Errorf("expected: %s, got: %s", expected, string(rawValue))
		t.FailNow()
	}

	var value1 BigInt
	if err := value1.FromValue(expected); err != nil {
		t.Errorf("failed to decode %s: %s", scalarName, err)
		t.FailNow()
	}

	if value1.String() != expected {
		t.Errorf("expected: %s, got: %s", expected, value1)
		t.FailNow()
	}

	if err := value1.FromValue(""); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	if err := json.Unmarshal([]byte("null"), &value1); err != nil {
		t.Errorf("expected nil, got error: %v", err)
		t.FailNow()
	}

	for _, b := range [][]byte{[]byte(""), []byte(`"abc"`)} {
		if err := json.Unmarshal(b, &value1); err == nil {
			t.Error("expected error, got nil")
			t.FailNow()
		}
	}
}
