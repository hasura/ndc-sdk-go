package main

import (
	"context"

	"github.com/hasura/ndc-codegen-empty-test/types"
	"github.com/hasura/ndc-sdk-go/v2/connector"
	"github.com/hasura/ndc-sdk-go/v2/schema"
)

var connectorCapabilities = schema.CapabilitiesResponse{
	Version: schema.NDCVersion,
	Capabilities: schema.Capabilities{
		Query: schema.QueryCapabilities{
			Variables: schema.LeafCapability{},
		},
		Mutation: schema.MutationCapabilities{},
	},
}

// Connector implements the SDK interface of NDC specification
type Connector struct{}

// ParseConfiguration validates the configuration files provided by the user, returning a validated 'Configuration',
// or throwing an error to prevents Connector startup.
func (connector *Connector) ParseConfiguration(ctx context.Context, configurationDir string) (*types.Configuration, error) {
	return &types.Configuration{}, nil
}

// TryInitState initializes the connector's in-memory state.
//
// For example, any connection pools, prepared queries,
// or other managed resources would be allocated here.
//
// In addition, this function should register any
// connector-specific metrics with the metrics registry.
func (connector *Connector) TryInitState(ctx context.Context, configuration *types.Configuration, metrics *connector.TelemetryState) (*types.State, error) {
	return &types.State{
		TelemetryState: metrics,
	}, nil
}

// HealthCheck checks the health of the connector.
//
// For example, this function should check that the connector
// is able to reach its data source over the network.
//
// Should throw if the check fails, else resolve.
func (mc *Connector) HealthCheck(ctx context.Context, configuration *types.Configuration, state *types.State) error {
	return nil
}

// GetCapabilities get the connector's capabilities.
func (mc *Connector) GetCapabilities(configuration *types.Configuration) schema.CapabilitiesResponseMarshaler {
	return connectorCapabilities
}

// QueryExplain explains a query by creating an execution plan.
func (mc *Connector) QueryExplain(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("query explain has not been supported yet", nil)
}

// MutationExplain explains a mutation by creating an execution plan.
func (mc *Connector) MutationExplain(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.MutationRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("mutation explain has not been supported yet", nil)
}
