package command

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestPatchConnectorContent(t *testing.T) {
	fixture := `
	var connectorCapabilities = schema.CapabilitiesResponse{
		Version: "0.1.6",
		Capabilities: schema.Capabilities{
			Query: schema.QueryCapabilities{
				Variables:    schema.LeafCapability{},
				NestedFields: schema.NestedFieldCapabilities{},
			},
			Mutation: schema.MutationCapabilities{},
		},
	}`

	expected := `
	var connectorCapabilities = schema.CapabilitiesResponse{
		Version: schema.NDCVersion,
		Capabilities: schema.Capabilities{
			Query: schema.QueryCapabilities{
				Variables:    &schema.LeafCapability{},
				NestedFields: schema.NestedFieldCapabilities{},
			},
			Mutation: schema.MutationCapabilities{},
		},
	}

// Close handles the graceful shutdown that cleans up the connector's state.
func (c *Connector) Close(state *types.State) error {
	return nil
}`

	ucc := upgradeConnectorCommand{}
	result, changed := ucc.patchConnectorContent([]byte(fixture))
	assert.Assert(t, changed)
	assert.Equal(t, expected, result)
}

func TestPatchImportSdkV2Content(t *testing.T) {
	fixture := `
	import (
		"github.com/hasura/ndc-sdk-go/schema"
		"github.com/hasura/ndc-sdk-go/utils"
	)`

	expected := `
	import (
		"github.com/hasura/ndc-sdk-go/v2/schema"
		"github.com/hasura/ndc-sdk-go/v2/utils"
	)`

	ucc := upgradeConnectorCommand{}
	result, changed := ucc.patchImportSdkV2Content([]byte(fixture))
	assert.Assert(t, changed)
	assert.Equal(t, expected, result)
}
