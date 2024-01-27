package connector

import (
	"context"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggest/jsonschema-go"
)

// Connector abstracts an interface with required methods for the [NDC Specification].
//
// [NDC Specification]: https://hasura.github.io/ndc-spec/specification/index.html
type Connector[RawConfiguration any, Configuration any, State any] interface {
	// Return jsonschema for the raw configuration for this connector.
	GetRawConfigurationSchema() *jsonschema.Schema

	// Return an empty raw configuration, to be manually filled in by the user to allow connection to the data source.
	// The exact shape depends on your connector's configuration. Example:
	//
	//   {
	//     "connection_string": "",
	//     "tables": []
	//   }
	MakeEmptyConfiguration() *RawConfiguration

	// Take a raw configuration, update it where appropriate by connecting to the underlying data source, and otherwise return it as-is
	// For example, if our configuration includes a list of tables, we may want to fetch an updated list from the data source.
	// This is also used to "hidrate" an "empty" configuration where a user has provided connection details and little else.
	UpdateConfiguration(ctx context.Context, rawConfiguration *RawConfiguration) (*RawConfiguration, error)

	// Validate the raw configuration provided by the user,
	// returning a configuration error or a validated [`Connector::Configuration`].
	ValidateRawConfiguration(rawConfiguration *RawConfiguration) (*Configuration, error)

	// Initialize the connector's in-memory state.
	//
	// For example, any connection pools, prepared queries,
	// or other managed resources would be allocated here.
	//
	// In addition, this function should register any
	// connector-specific metrics with the metrics registry.
	TryInitState(configuration *Configuration, metrics *TelemetryState) (*State, error)

	// Check the health of the connector.
	//
	// For example, this function should check that the connector
	// is able to reach its data source over the network.
	//
	// Should throw if the check fails, else resolve.
	HealthCheck(ctx context.Context, configuration *Configuration, state *State) error

	// Get the connector's capabilities.
	//
	// This function implements the [capabilities endpoint] from the NDC specification.
	//
	// This function should be synchronous.
	//
	// [capabilities endpoint]: https://hasura.github.io/ndc-spec/specification/capabilities.html
	GetCapabilities(configuration *Configuration) *schema.CapabilitiesResponse

	// Get the connector's schema.
	//
	// This function implements the [schema endpoint] from the NDC specification.
	//
	// [schema endpoint]: https://hasura.github.io/ndc-spec/specification/schema/index.html
	GetSchema(configuration *Configuration) (*schema.SchemaResponse, error)

	// Explain a query by creating an execution plan.
	// This function implements the [explain endpoint] from the NDC specification.
	//
	// [explain endpoint]: https://hasura.github.io/ndc-spec/specification/explain.html
	Explain(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.ExplainResponse, error)

	// Execute a mutation.
	//
	// This function implements the [mutation endpoint] from the NDC specification.
	//
	// [mutation endpoint]: https://hasura.github.io/ndc-spec/specification/mutations/index.html
	Mutation(ctx context.Context, configuration *Configuration, state *State, request *schema.MutationRequest) (*schema.MutationResponse, error)

	// Execute a query.
	//
	// This function implements the [query endpoint] from the NDC specification.
	//
	// [query endpoint]: https://hasura.github.io/ndc-spec/specification/queries/index.html
	Query(ctx context.Context, configuration *Configuration, state *State, request *schema.QueryRequest) (*schema.QueryResponse, error)
}

// the common serve options for the server
type serveOptions struct {
	logger        zerolog.Logger
	metricsPrefix string
	version       string
	serviceName   string
	withoutConfig bool
}

func defaultServeOptions() *serveOptions {
	return &serveOptions{
		logger:        log.Level(zerolog.GlobalLevel()),
		serviceName:   "ndc-go",
		version:       "0.1.0",
		withoutConfig: false,
	}
}

// ServeOption abstracts a public interface to update server options
type ServeOption func(*serveOptions)

// WithLogger sets a custom logger option
func WithLogger(logger zerolog.Logger) ServeOption {
	return func(so *serveOptions) {
		so.logger = logger
	}
}

// GetLogger gets the logger instance from context
func GetLogger(ctx context.Context) zerolog.Logger {
	return internal.GetLogger(ctx)
}

// WithMetricsPrefix sets the custom metrics prefix option
func WithMetricsPrefix(prefix string) ServeOption {
	return func(so *serveOptions) {
		so.metricsPrefix = prefix
	}
}

// WithVersion sets the custom version option
func WithVersion(version string) ServeOption {
	return func(so *serveOptions) {
		so.version = version
	}
}

// WithDefaultServiceName sets the default service name option
func WithDefaultServiceName(name string) ServeOption {
	return func(so *serveOptions) {
		so.serviceName = name
	}
}

// WithoutConfig makes the configuration flag optional
func WithoutConfig() ServeOption {
	return func(so *serveOptions) {
		so.withoutConfig = true
	}
}
