package scalar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBytes(t *testing.T) {
	plain := "Hello world"
	expected := NewBytes([]byte(plain))
	expectedB64 := base64.StdEncoding.EncodeToString([]byte(plain))

	rawJSON := fmt.Sprintf(`"%s"`, expectedB64)
	var value Bytes
	if err := json.Unmarshal([]byte(rawJSON), &value); err != nil {
		t.Errorf("failed to parse Bytes: %s", err)
		t.FailNow()
	}
	if plain != value.String() {
		t.Errorf("expected: %s, got: %s", expected, value)
		t.FailNow()
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		t.Errorf("failed to encode Bytes: %s", err)
		t.FailNow()
	}

	if string(expectedB64) != strings.Trim(string(rawValue), "\"") {
		t.Errorf("expected: %s, got: %s", expected, string(rawValue))
		t.FailNow()
	}

	var value1 Bytes

	if err := value1.FromValue(nil); err != nil {
		t.Errorf("expected nil, got error: %s", err)
		t.FailNow()
	}

	if err := value1.FromValue(expectedB64); err != nil {
		t.Errorf("expected nil, got error: %s", err)
		t.FailNow()
	}

	if value1.String() != plain {
		t.Errorf("expected: %s, got: %s", expected, value1)
		t.FailNow()
	}

	if err := value1.FromValue(0); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	if err := value1.FromValue("abc"); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}

	for _, b := range [][]byte{[]byte(""), []byte("0"), []byte(`"abc"`)} {
		if err := json.Unmarshal(b, &value1); err == nil {
			t.Error("expected error, got nil")
			t.FailNow()
		}
	}

	if _, err := ParseBytes("abc"); err == nil {
		t.Error("expected error, got nil")
		t.FailNow()
	}
}
