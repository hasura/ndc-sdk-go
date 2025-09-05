package connector

import (
	"context"
	"log/slog"

	"github.com/hasura/ndc-sdk-go/v2/schema"
)

// Connector abstracts an interface with required methods for the [NDC Specification].
//
// [NDC Specification]: https://hasura.github.io/ndc-spec/specification/index.html
type Connector[Configuration any, State any] interface {
	// ParseConfiguration validates the configuration files provided by the user, returning a validated 'Configuration',
	// or throwing an error to prevents Connector startup.
	ParseConfiguration(ctx context.Context, configurationDir string) (*Configuration, error)

	// TryInitState initializes the connector's in-memory state.
	//
	// For example, any connection pools, prepared queries,
	// or other managed resources would be allocated here.
	//
	// In addition, this function should register any
	// connector-specific metrics with the metrics registry.
	TryInitState(
		ctx context.Context,
		configuration *Configuration,
		metrics *TelemetryState,
	) (*State, error)

	// HealthCheck checks the health of the connector.
	//
	// For example, this function should check that the connector
	// is able to reach its data source over the network.
	//
	// Should throw if the check fails, else resolve.
	HealthCheck(ctx context.Context, configuration *Configuration, state *State) error

	// GetCapabilities get the connector's capabilities.
	//
	// This function implements the [capabilities endpoint] from the NDC specification.
	//
	// This function should be synchronous.
	//
	// [capabilities endpoint]: https://hasura.github.io/ndc-spec/specification/capabilities.html
	GetCapabilities(configuration *Configuration) schema.CapabilitiesResponseMarshaler

	// GetSchema gets the connector's schema.
	//
	// This function implements the [schema endpoint] from the NDC specification.
	//
	// [schema endpoint]: https://hasura.github.io/ndc-spec/specification/schema/index.html
	GetSchema(
		ctx context.Context,
		configuration *Configuration,
		state *State,
	) (schema.SchemaResponseMarshaler, error)

	// QueryExplain explains a query by creating an execution plan.
	// This function implements the [explain endpoint] from the NDC specification.
	//
	// [explain endpoint]: https://hasura.github.io/ndc-spec/specification/explain.html
	QueryExplain(
		ctx context.Context,
		configuration *Configuration,
		state *State,
		request *schema.QueryRequest,
	) (*schema.ExplainResponse, error)

	// MutationExplain explains a mutation by creating an execution plan.
	// This function implements the [explain endpoint] from the NDC specification.
	//
	// [explain endpoint]: https://hasura.github.io/ndc-spec/specification/explain.html
	MutationExplain(
		ctx context.Context,
		configuration *Configuration,
		state *State,
		request *schema.MutationRequest,
	) (*schema.ExplainResponse, error)

	// Mutation executes a mutation.
	//
	// This function implements the [mutation endpoint] from the NDC specification.
	//
	// [mutation endpoint]: https://hasura.github.io/ndc-spec/specification/mutations/index.html
	Mutation(
		ctx context.Context,
		configuration *Configuration,
		state *State,
		request *schema.MutationRequest,
	) (*schema.MutationResponse, error)

	// Query executes a query.
	//
	// This function implements the [query endpoint] from the NDC specification.
	//
	// [query endpoint]: https://hasura.github.io/ndc-spec/specification/queries/index.html
	Query(
		ctx context.Context,
		configuration *Configuration,
		state *State,
		request *schema.QueryRequest,
	) (schema.QueryResponse, error)

	// Close handles the graceful shutdown that cleans up the connector's state.
	Close(state *State) error
}

// the common serve options for the server.
type serveOptions struct {
	logger          *slog.Logger
	logLevel        slog.Level
	metricsPrefix   string
	version         string
	serviceName     string
	withoutConfig   bool
	withoutRecovery bool
}

func defaultServeOptions() *serveOptions {
	return &serveOptions{
		logger:          slog.Default(),
		serviceName:     "hasura-ndc-go",
		version:         "0.1.0",
		withoutConfig:   false,
		withoutRecovery: false,
	}
}

// ServeOption abstracts a public interface to update server options.
type ServeOption func(*serveOptions)

// WithLogger sets a custom logger option.
func WithLogger(logger *slog.Logger) ServeOption {
	return func(so *serveOptions) {
		so.logger = logger
	}
}

// WithLoggerFunc sets a custom logger option with a constructor function.
func WithLoggerFunc(fn func(level slog.Level) *slog.Logger) ServeOption {
	return func(so *serveOptions) {
		so.logger = fn(so.logLevel)
	}
}

// WithMetricsPrefix sets the custom metrics prefix option.
func WithMetricsPrefix(prefix string) ServeOption {
	return func(so *serveOptions) {
		so.metricsPrefix = prefix
	}
}

// WithVersion sets the custom version option.
func WithVersion(version string) ServeOption {
	return func(so *serveOptions) {
		so.version = version
	}
}

// WithDefaultServiceName sets the default service name option.
func WithDefaultServiceName(name string) ServeOption {
	return func(so *serveOptions) {
		so.serviceName = name
	}
}

// WithoutRecovery disables recovery on panic.
func WithoutRecovery() ServeOption {
	return func(so *serveOptions) {
		so.withoutRecovery = true
	}
}

func withLogLevel(level slog.Level) ServeOption {
	return func(so *serveOptions) {
		so.logLevel = level
	}
}
