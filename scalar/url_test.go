package scalar

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestURL(t *testing.T) {
	expectedStr := "http://localhost:8080/foo?bar=baz"
	expected, err := NewURL(expectedStr)
	if err != nil {
		t.Fatalf("failed to parsed expected URL")
	}

	if _, err := NewURL(""); err == nil {
		t.Fatal("expected error, got nil")
	}
	rawJSON := fmt.Sprintf(`"%s"`, expected)
	var value URL
	scalarName := value.ScalarName()
	if value.String() != "" {
		t.Fatalf("expected empty string, got %s", value.String())
	}

	if err := json.Unmarshal([]byte(`""`), &value); err == nil {
		t.Fatal("expected error, got nil")
	}

	if err := json.Unmarshal([]byte(rawJSON), &value); err != nil {
		t.Fatalf("failed to parse %s: %s", scalarName, err)
	}
	if expected.String() != value.String() {
		t.Fatalf("expected: %s, got: %s", expected, value)
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("failed to encode %s: %s", scalarName, err)
	}

	if expectedStr != strings.Trim(string(rawValue), "\"") {
		t.Fatalf("expected: %s, got: %s", expectedStr, string(rawValue))
	}

	var value1 URL
	if err := value1.FromValue(expectedStr); err != nil {
		t.Fatalf("failed to decode %s: %s", scalarName, err)
	}

	if value1.String() != expectedStr {
		t.Fatalf("expected: %s, got: %s", expected, value1)
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

	if err := json.Unmarshal([]byte("0"), &value1); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	var anyStr any = "example.com"
	if err := value.FromValue(any(&anyStr)); err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}
	if value.String() != fmt.Sprint(anyStr) {
		t.Fatalf("expected %s, got: %s", anyStr, value.String())
	}

	if err := value1.FromValue(value); err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}
	if value1.String() != fmt.Sprint(anyStr) {
		t.Fatalf("expected %s, got: %s", anyStr, value1.String())
	}

	var value2 URL
	if err := value2.FromValue(&value); err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}
	if value2.String() != fmt.Sprint(anyStr) {
		t.Fatalf("expected %s, got: %s", anyStr, value2.String())
	}

	rawURL, _ := url.Parse("https://google.com")
	if err := value1.FromValue(rawURL); err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}
	if value1.String() != rawURL.String() {
		t.Fatalf("expected %s, got: %s", rawURL.String(), value1.String())
	}
	if err := value2.FromValue(*rawURL); err != nil {
		t.Fatalf("expected nil error, got: %s", err)
	}
	if value2.String() != rawURL.String() {
		t.Fatalf("expected %s, got: %s", rawURL.String(), value2.String())
	}
}
