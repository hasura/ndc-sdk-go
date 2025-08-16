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

