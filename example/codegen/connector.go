package main

import (
	"context"
	_ "embed"

	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/swaggest/jsonschema-go"
)

type RawConfiguration struct{}
type Configuration struct{}

type State struct{}
type Connector struct{}

func (connector *Connector) GetRawConfigurationSchema() *jsonschema.Schema {
	return nil
}
func (connector *Connector) MakeEmptyConfiguration() *RawConfiguration {
	return &RawConfiguration{}
}

func (connector *Connector) UpdateConfiguration(ctx context.Context, rawConfiguration *RawConfiguration) (*RawConfiguration, error) {
	return &RawConfiguration{}, nil
}
func (connector *Connector) ValidateRawConfiguration(rawConfiguration *RawConfiguration) (*Configuration, error) {
	return &Configuration{}, nil
}

func (connector *Connector) TryInitState(configuration *Configuration, metrics *connector.TelemetryState) (*State, error) {

	return &State{}, nil
}

func (mc *Connector) HealthCheck(ctx context.Context, configuration *Configuration, state *State) error {
	return nil
}

func (mc *Connector) GetCapabilities(configuration *Configuration) *schema.CapabilitiesResponse {
	return &schema.CapabilitiesResponse{
		Version: "^0.1.0",
		Capabilities: schema.Capabilities{
			Query: schema.QueryCapabilities{
				Variables: schema.LeafCapability{},
			},
		},
	}
}

func (mc *Connector) QueryExplain(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("query explain has not been supported yet", nil)
}

func (mc *Connector) MutationExplain(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.ExplainResponse, error) {
	return nil, schema.NotSupportedError("mutation explain has not been supported yet", nil)
}
