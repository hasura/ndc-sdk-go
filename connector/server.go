package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	errConfigurationRequired = errors.New("Configuration is required")
)

// ServerOptions presents the configuration object of the connector http server
type ServerOptions struct {
	Configuration       string
	InlineConfig        bool
	ServiceTokenSecret  string
	OTLPEndpoint        string
	OTLPInsecure        bool
	OTLPTracesEndpoint  string
	OTLPMetricsEndpoint string
	ServiceName         string
}

// Server implements the [NDC API specification] for the connector
//
// [NDC API specification]: https://hasura.github.io/ndc-spec/specification/index.html
type Server[RawConfiguration any, Configuration any, State any] struct {
	context       context.Context
	stop          context.CancelFunc
	connector     Connector[RawConfiguration, Configuration, State]
	state         *State
	configuration *Configuration
	options       *ServerOptions
	telemetry     *TelemetryState
	logger        zerolog.Logger
}

// NewServer creates a Server instance
func NewServer[RawConfiguration any, Configuration any, State any](connector Connector[RawConfiguration, Configuration, State], options *ServerOptions, others ...ServeOption) (*Server[RawConfiguration, Configuration, State], error) {
	defaultOptions := defaultServeOptions()
	for _, opts := range others {
		opts(defaultOptions)
	}

	var rawConfiguration RawConfiguration
	if !defaultOptions.withoutConfig {
		if options.Configuration == "" {
			return nil, errConfigurationRequired
		}

		configBytes := []byte(options.Configuration)
		if !options.InlineConfig {
			var err error
			configBytes, err = os.ReadFile(options.Configuration)
			if err != nil {
				return nil, fmt.Errorf("Invalid configuration provided: %s", err)
			}

			if len(configBytes) == 0 {
				return nil, errConfigurationRequired
			}
		}

		if err := json.Unmarshal(configBytes, &rawConfiguration); err != nil {
			return nil, fmt.Errorf("Invalid configuration provided: %s", err)
		}
	}

	configuration, err := connector.ValidateRawConfiguration(&rawConfiguration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(configuration, nil)
	if err != nil {
		return nil, err
	}

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defaultOptions.logger.Debug().
		Str("endpoint", options.OTLPEndpoint).
		Str("traces_endpoint", options.OTLPTracesEndpoint).
		Str("metrics_endpoint", options.OTLPMetricsEndpoint).
		Str("service_name", options.ServiceName).
		Str("version", defaultOptions.version).
		Str("metrics_prefix", defaultOptions.metricsPrefix).
		Msg("initialize OpenTelemetry")

	if options.ServiceName == "" {
		options.ServiceName = defaultOptions.serviceName
	}

	telemetry, err := setupOTelSDK(ctx, options, defaultOptions.version, defaultOptions.metricsPrefix)
	if err != nil {
		return nil, err
	}

	return &Server[RawConfiguration, Configuration, State]{
		context:       ctx,
		stop:          stop,
		connector:     connector,
		state:         state,
		configuration: configuration,
		options:       options,
		telemetry:     telemetry,
		logger:        defaultOptions.logger,
	}, nil
}

func (s *Server[RawConfiguration, Configuration, State]) withAuth(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r.Context())
		if s.options.ServiceTokenSecret != "" && r.Header.Get("authorization") != fmt.Sprintf("Bearer %s", s.options.ServiceTokenSecret) {
			writeJson(w, logger, http.StatusUnauthorized, schema.ErrorResponse{
				Message: "Unauthorized",
				Details: map[string]any{
					"cause": "Bearer token does not match.",
				},
			})

			s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("status", "failed"),
				attribute.String("reason", "unauthorized"),
			))
			return
		}

		handler(w, r)
	}
}

// GetCapabilities get the connector's capabilities. Implement a handler for the /capabilities endpoint, GET method.
func (s *Server[RawConfiguration, Configuration, State]) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	capabilities := s.connector.GetCapabilities(s.configuration)
	writeJson(w, logger, http.StatusOK, capabilities)
}

// Health checks the health of the connector. Implement a handler for the /health endpoint, GET method.
func (s *Server[RawConfiguration, Configuration, State]) Health(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	if err := s.connector.HealthCheck(r.Context(), s.configuration, s.state); err != nil {
		writeError(w, logger, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSchema implements a handler for the /schema endpoint, GET method.
func (s *Server[RawConfiguration, Configuration, State]) GetSchema(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	schemaResult, err := s.connector.GetSchema(s.configuration)
	if err != nil {
		writeError(w, logger, err)
		return
	}

	writeJson(w, logger, http.StatusOK, schemaResult)
}

// Query implements a handler for the /query endpoint, POST method that executes a query.
func (s *Server[RawConfiguration, Configuration, State]) Query(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Query", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	attributes := []attribute.KeyValue{}
	_, decodeSpan := s.telemetry.Tracer.Start(ctx, "Decode JSON Body")
	defer decodeSpan.End()
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", "json_decode"),
		}
		span.SetAttributes(attributes...)
		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return
	}
	decodeSpan.End()

	collectionAttr := attribute.String("collection", body.Collection)
	attributes = append(attributes, collectionAttr)
	span.SetAttributes(attributes...)
	execQueryCtx, execQuerySpan := s.telemetry.Tracer.Start(ctx, "Execute Query")
	defer execQuerySpan.End()

	response, err := s.connector.Query(execQueryCtx, s.configuration, s.state, &body)

	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(statusAttributes...)
		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttributes...)...))
		return
	}
	execQuerySpan.End()

	statusAttribute := attribute.String("status", "success")
	span.SetAttributes(statusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()

	s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttribute)...))
	// record latency for success requests only
	s.telemetry.queryLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// QueryExplain implements a handler for the /query/explain endpoint, POST method that explains a query by creating an execution plan.
func (s *Server[RawConfiguration, Configuration, State]) QueryExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Query Explain", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	attributes := []attribute.KeyValue{}
	_, decodeSpan := s.telemetry.Tracer.Start(ctx, "Decode JSON Body")
	defer decodeSpan.End()
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", "json_decode"),
		}
		span.SetAttributes(attributes...)
		s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return
	}
	decodeSpan.End()

	collectionAttr := attribute.String("collection", body.Collection)
	attributes = append(attributes, collectionAttr)
	span.SetAttributes(attributes...)
	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Explain")
	defer execSpan.End()

	response, err := s.connector.QueryExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttributes...)...))
		return
	}
	execSpan.End()

	statusAttribute := attribute.String("status", "success")
	span.SetAttributes(statusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()
	s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttribute)...))

	// record latency for success requests only
	s.telemetry.explainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// MutationExplain implements a handler for the /mutation/explain endpoint, POST method that explains a mutation by creating an execution plan.
func (s *Server[RawConfiguration, Configuration, State]) MutationExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Mutation Explain", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	attributes := []attribute.KeyValue{}
	_, decodeSpan := s.telemetry.Tracer.Start(ctx, "Decode JSON Body")
	defer decodeSpan.End()
	var body schema.MutationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", "json_decode"),
		}
		span.SetAttributes(attributes...)
		s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return
	}
	decodeSpan.End()

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}
	collectionAttr := attribute.String("operations", strings.Join(operationNames, ","))
	attributes = append(attributes, collectionAttr)
	span.SetAttributes(attributes...)
	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Explain")
	defer execSpan.End()

	response, err := s.connector.MutationExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttributes...)...))
		return
	}
	execSpan.End()

	statusAttribute := attribute.String("status", "success")
	span.SetAttributes(statusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()
	s.telemetry.explainCounter.Add(r.Context(), 1, metric.WithAttributes(append(attributes, statusAttribute)...))

	// record latency for success requests only
	s.telemetry.explainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// Mutation implements a handler for the /mutation endpoint, POST method that executes a mutation.
func (s *Server[RawConfiguration, Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Mutation", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	_, decodeSpan := s.telemetry.Tracer.Start(ctx, "Decode JSON Body")
	defer decodeSpan.End()
	var body schema.MutationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", "json_decode"),
		}
		span.SetAttributes(attributes...)
		s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return
	}
	decodeSpan.End()

	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Mutation")
	defer execSpan.End()
	response, err := s.connector.Mutation(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return
	}
	execSpan.End()

	attributes := attribute.String("status", "success")
	span.SetAttributes(attributes)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()

	s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(attributes))

	// record latency for success requests only
	s.telemetry.mutationLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds())
}

func (s *Server[RawConfiguration, Configuration, State]) buildHandler() *http.ServeMux {
	router := newRouter(s.logger)
	router.Use("/capabilities", http.MethodGet, s.withAuth(s.GetCapabilities))
	router.Use("/schema", http.MethodGet, s.withAuth(s.GetSchema))
	router.Use("/query", http.MethodPost, s.withAuth(s.Query))
	router.Use("/query/explain", http.MethodPost, s.withAuth(s.QueryExplain))
	router.Use("/mutation/explain", http.MethodPost, s.withAuth(s.MutationExplain))
	router.Use("/mutation", http.MethodPost, s.withAuth(s.Mutation))
	router.Use("/health", http.MethodGet, s.Health)
	router.Use("/metrics", http.MethodGet, s.withAuth(promhttp.Handler().ServeHTTP))

	return router.Build()
}

// ListenAndServe serves the configuration server with the standard http server.
// You can also replace this method with any router or web framework that is compatible with net/http.
func (s *Server[RawConfiguration, Configuration, State]) ListenAndServe(port uint) error {
	defer s.stop()
	defer func() {
		if err := s.telemetry.Shutdown(context.Background()); err != nil {
			s.logger.Error().Err(err).Msg("failed to shutdown OpenTelemetry")
		}
	}()

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context {
			return s.context
		},
		Handler: s.buildHandler(),
	}

	serverErr := make(chan error, 1)
	go func() {
		s.logger.Info().Msgf("Listening server on %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interruption.
	select {
	case err := <-serverErr:
		// Error when starting HTTP server.
		return err
	case <-s.context.Done():
		// Wait for first CTRL+C.
		s.logger.Info().Msg("received the quit signal, exiting...")
		// Stop receiving signal notifications as soon as possible.
		s.stop()
		// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
		return server.Shutdown(context.Background())
	}
}
