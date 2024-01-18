package connector

import (
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/swaggest/jsonschema-go"
)

type mockRawConfiguration struct{}
type mockConfiguration struct{}
type mockState struct{}
type mockConnector struct{}

func (mc *mockConnector) GetRawConfigurationSchema() *jsonschema.Schema {
	return nil
}
func (mc *mockConnector) MakeEmptyConfiguration() *mockRawConfiguration {
	return &mockRawConfiguration{}
}

func (mc *mockConnector) UpdateConfiguration(rawConfiguration *mockRawConfiguration) (*mockRawConfiguration, error) {
	return &mockRawConfiguration{}, nil
}
func (mc *mockConnector) ValidateRawConfiguration(rawConfiguration *mockRawConfiguration) (*mockConfiguration, error) {
	return &mockConfiguration{}, nil
}
func (mc *mockConnector) TryInitState(configuration *mockConfiguration, metrics any) *mockState {
	return &mockState{}
}

func (mc *mockConnector) FetchMetrics(configuration *mockConfiguration, state *mockState) error {
	return nil
}

func (mc *mockConnector) HealthCheck(configuration *mockConfiguration, state *mockState) error {
	return nil
}

func (mc *mockConnector) GetCapabilities(configuration *mockConfiguration) (*schema.CapabilitiesResponse, error) {
	return &schema.CapabilitiesResponse{}, nil
}

func (mc *mockConnector) GetSchema(configuration *mockConfiguration) (*schema.SchemaResponse, error) {
	return &schema.SchemaResponse{}, nil
}
func (mc *mockConnector) Explain(configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (*schema.ExplainResponse, error) {
	return &schema.ExplainResponse{}, nil
}
func (mc *mockConnector) Mutation(configuration *mockConfiguration, state *mockState, request *schema.MutationRequest) (*schema.MutationResponse, error) {
	return &schema.MutationResponse{}, nil
}
func (mc *mockConnector) Query(configuration *mockConfiguration, state *mockState, request *schema.QueryRequest) (*schema.QueryResponse, error) {
	return &schema.QueryResponse{}, nil
}
