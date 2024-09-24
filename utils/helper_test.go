package utils

import (
	"log/slog"
	"testing"
)

func TestToPtrs(t *testing.T) {
	assertDeepEqual(t, []*string{ToPtr("a"), ToPtr("b"), ToPtr("c")}, ToPtrs([]string{"a", "b", "c"}))
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

	assertDeepEqual(t, expected, GetSortedKeys(input))
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
	assertDeepEqual(t, mapB, MergeMap(nil, mapB))
	assertDeepEqual(t, map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}, MergeMap(mapA, mapB))
}
