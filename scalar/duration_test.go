package scalar

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestDuration(t *testing.T) {
	expected := NewDuration(time.Hour).String()
	rawJSON := fmt.Sprintf(`"%s"`, expected)
	var value Duration
	scalarName := value.ScalarName()

	assert.Equal(t, "0s", (Duration{}).String())
	if err := json.Unmarshal([]byte(rawJSON), &value); err != nil {
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

	var value1 Duration
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

	if err := value1.FromValue(nil); err != nil {
		t.Errorf("expected no error, got %v", err)
		t.FailNow()
	}

	if err := json.Unmarshal([]byte(`"abc"`), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	if err := json.Unmarshal([]byte("abc"), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}
}
