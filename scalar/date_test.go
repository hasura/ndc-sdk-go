package scalar

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestDate(t *testing.T) {
	expected := NewDate(2014, 04, 1).String()
	rawJSON := fmt.Sprintf(`"%s"`, expected)
	var value Date
	if err := json.Unmarshal([]byte(rawJSON), &value); err != nil {
		t.Errorf("failed to parse Date: %s", err)
		t.FailNow()
	}
	if expected != value.String() {
		t.Errorf("expected: %s, got: %s", expected, value)
		t.FailNow()
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		t.Errorf("failed to encode Date: %s", err)
		t.FailNow()
	}

	if expected != strings.Trim(string(rawValue), "\"") {
		t.Errorf("expected: %s, got: %s", expected, string(rawValue))
		t.FailNow()
	}

	var value1 Date
	if err := value1.FromValue(expected); err != nil {
		t.Errorf("failed to decode Date: %s", err)
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

	if err := value1.FromValue(0); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	if err := value1.FromValue(nil); err != nil {
		t.Errorf("expected no error, got %v", err)
		t.FailNow()
	}

	if err := json.Unmarshal([]byte(`"abc"`), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	if err := json.Unmarshal([]byte("0"), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}
	if _, err := ParseDate(""); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}
}
