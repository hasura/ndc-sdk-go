package utils

import (
	"log/slog"
	"testing"

	"gotest.tools/v3/assert"
)

func TestToPtrs(t *testing.T) {
	_, err := PointersToValues([]*string{ToPtr(""), nil})
	assertError(t, err, "element at 1 must not be nil")
	input, err := PointersToValues(ToPtrs([]string{"a", "b", "c"}))
	assert.NilError(t, err)
	expected, err := PointersToValues([]*string{ToPtr("a"), ToPtr("b"), ToPtr("c")})
	assert.NilError(t, err)

	assert.DeepEqual(t, expected, input)
}

func TestIsDebug(t *testing.T) {
	if IsDebug(slog.Default()) {
		t.Error("expected debug mode, got false")
	}
}

func TestGetSortedKeys(t *testing.T) {
	input := map[string]any{
		"b": "b",
		"c": "c",
		"a": "a",
		"d": "d",
	}
	expected := []string{"a", "b", "c", "d"}

	assert.DeepEqual(t, expected, GetSortedKeys(input))
}

func TestMergeMap(t *testing.T) {
	mapA := map[string]any{
		"a": 2,
		"c": 3,
	}
	mapB := map[string]any{
		"a": 1,
		"b": 2,
	}
	assert.DeepEqual(t, mapB, MergeMap(nil, mapB))
	assert.DeepEqual(t, map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}, MergeMap(mapA, mapB))
}

func TestParseIntMapFromString(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected map[string]int
		ErrorMsg string
	}{
		{
			Expected: map[string]int{},
		},
		{
			Input: "a=1;b=2;c=3",
			Expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			Input:    "a;b=2",
			ErrorMsg: "invalid int map string a;b=2, expected <key1>=<value1>;<key2>=<value2>",
		},
		{
			Input:    "a=c;b=2",
			ErrorMsg: "invalid integer value c in item a=c",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			result, err := ParseIntMapFromString(tc.Input)
			if tc.ErrorMsg != "" {
				assert.ErrorContains(t, err, tc.ErrorMsg)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, result, tc.Expected)
			}
		})
	}
}
