# Migrate from SDK v1 to v2

## Capabilities

- Replace the hard-coded NDC version to `schema.NDCVersion`. The SDK will automatically increase the version. 
- Left capabilities are type safe. Replace the `interface{}` with `&schema.LeafCapability{}`.

```go
import (
	"github.com/hasura/ndc-sdk-go/schema"
)

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
```

## Connector

### Add Close method

The `Connector` interface adds a new [Close](https://github.com/hasura/ndc-sdk-go/blob/83a469863391fa412e1abe0e8efb5952662752ce/connector/types.go#L106) method to handle the graceful shutdown of the connector. You need to add this method to satisfy the interface.

```go
// Close handles the graceful shutdown that cleans up the connector's state.
func (c *Connector) Close(state *types.State) error {
	return nil
}
```