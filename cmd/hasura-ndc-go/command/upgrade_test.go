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
	}`

	ucc := upgradeConnectorCommand{}
	result, changed := ucc.patchConnectorContent([]byte(fixture))
	assert.Assert(t, changed)
	assert.Equal(t, expected, result)
}
