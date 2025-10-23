package utils

import (
	"log/slog"
	"testing"

	"gotest.tools/v3/assert"
)

func TestToPtrs(t *testing.T) {
	_, err := PointersToValues([]*string{ToPtr(""), nil})
	assert.ErrorContains(t, err, "element at 1 must not be nil")
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
