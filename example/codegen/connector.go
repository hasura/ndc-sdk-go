package main

import (
	"context"
	_ "embed"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/example/codegen/types"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/swaggest/jsonschema-go"
)

type Connector struct{}

func (connector *Connector) GetRawConfigurationSchema() *jsonschema.Schema {
	return nil
}
func (connector *Connector) MakeEmptyConfiguration() *types.RawConfiguration {
	return &types.RawConfiguration{}
}

func (connector *Connector) UpdateConfiguration(ctx context.Context, rawConfiguration *types.RawConfiguration) (*types.RawConfiguration, error) {
	return &types.RawConfiguration{}, nil
}
func (connector *Connector) ValidateRawConfiguration(rawConfiguration *types.RawConfiguration) (*types.Configuration, error) {
	return &types.Configuration{}, nil
}

func (connector *Connector) TryInitState(configuration *types.Configuration, metrics *connector.TelemetryState) (*types.State, error) {
	return &types.State{}, nil
}

func (mc *Connector) HealthCheck(ctx context.Context, configuration *types.Configuration, state *types.State) error {
	return nil
}

func (mc *Connector) GetCapabilities(configuration *types.Configuration) *schema.CapabilitiesResponse {
	return &schema.CapabilitiesResponse{
		Version: "^0.1.0",
		Capabilities: schema.Capabilities{
			Query: schema.QueryCapabilities{
				Variables: schema.LeafCapability{},
			},
			Mutation: schema.MutationCapabilities{},
		},
	}
}

func (mc *Connector) QueryExplain(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("query explain has not been supported yet", nil)
}

func (mc *Connector) MutationExplain(ctx context.Context, configuration *types.Configuration, state *types.State, request *schema.MutationRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("mutation explain has not been supported yet", nil)
}
