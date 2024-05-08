package command

import (
	"testing"

	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command/internal"
	"github.com/stretchr/testify/assert"
)

func TestGenTestSnapshots(t *testing.T) {
	assert.NoError(t, GenTestSnapshots(&GenTestSnapshotArguments{
		Schema:   "testdata/snapshots/schema.json",
		Dir:      "testdata/snapshots/generated",
		Depth:    10,
		Strategy: internal.WriteFileStrategyOverride,
	}))
}
