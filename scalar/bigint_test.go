package scalar

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBigInt(t *testing.T) {
	expected := "9999"
	value := NewBigInt(0)
	if err := json.Unmarshal([]byte(expected), &value); err != nil {
		t.Errorf("failed to parse bigint: %s", err)
		t.FailNow()
	}
	if expected != value.String() {
		t.Errorf("expected: %s, got: %s", expected, value)
		t.FailNow()
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		t.Errorf("failed to encode bigint: %s", err)
		t.FailNow()
	}

	if expected != strings.Trim(string(rawValue), "\"") {
		t.Errorf("expected: %s, got: %s", expected, string(rawValue))
		t.FailNow()
	}

	var value1 BigInt
	if err := value1.FromValue(expected); err != nil {
		t.Errorf("failed to decode bigint: %s", err)
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
	if err := json.Unmarshal([]byte(""), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}
}
