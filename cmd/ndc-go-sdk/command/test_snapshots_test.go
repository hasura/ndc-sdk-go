package command

import (
	"testing"

	"github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/command/internal"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenTestSnapshots(t *testing.T) {
	assert.NoError(t, GenTestSnapshots(&GenTestSnapshotArguments{
		Schema:   "testdata/snapshots/schema.json",
		Dir:      "testdata/snapshots/generated",
		Seed:     utils.ToPtr(int64(10)),
		Depth:    10,
		Strategy: internal.WriteFileStrategyOverride,
	}))
}
