package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
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
type Server[Configuration any, State any] struct {
	*serveOptions

	context       context.Context
	stop          context.CancelFunc
	connector     Connector[Configuration, State]
	state         *State
	configuration *Configuration
	options       *ServerOptions
	telemetry     *TelemetryState
}

// NewServer creates a Server instance
func NewServer[Configuration any, State any](connector Connector[Configuration, State], options *ServerOptions, others ...ServeOption) (*Server[Configuration, State], error) {
	defaultOptions := defaultServeOptions()
	for _, opts := range others {
		opts(defaultOptions)
	}

	configuration, err := connector.ParseConfiguration(options.Configuration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(configuration, nil)
	if err != nil {
		return nil, err
	}

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defaultOptions.logger.Debug(
		"initialize OpenTelemetry",
		slog.String("endpoint", options.OTLPEndpoint),
		slog.String("traces_endpoint", options.OTLPTracesEndpoint),
		slog.String("metrics_endpoint", options.OTLPMetricsEndpoint),
		slog.String("service_name", options.ServiceName),
		slog.String("version", defaultOptions.version),
		slog.String("metrics_prefix", defaultOptions.metricsPrefix),
	)

	if options.ServiceName == "" {
		options.ServiceName = defaultOptions.serviceName
	}

	telemetry, err := setupOTelSDK(ctx, options, defaultOptions.version, defaultOptions.metricsPrefix, defaultOptions.logger)
	if err != nil {
		return nil, err
	}

	return &Server[Configuration, State]{
		context:       ctx,
		stop:          stop,
		connector:     connector,
		state:         state,
		configuration: configuration,
		options:       options,
		telemetry:     telemetry,
		serveOptions:  defaultOptions,
	}, nil
}

func (s *Server[Configuration, State]) withAuth(handler http.HandlerFunc) http.HandlerFunc {

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
				failureStatusAttribute,
				httpStatusAttribute(http.StatusUnauthorized),
			))
			return
		}

		handler(w, r)
	}
}

// GetCapabilities get the connector's capabilities. Implement a handler for the /capabilities endpoint, GET method.
func (s *Server[Configuration, State]) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	capabilities := s.connector.GetCapabilities(s.configuration)
	writeJson(w, logger, http.StatusOK, capabilities)
}

// Health checks the health of the connector. Implement a handler for the /health endpoint, GET method.
func (s *Server[Configuration, State]) Health(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	if err := s.connector.HealthCheck(r.Context(), s.configuration, s.state); err != nil {
		writeError(w, logger, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSchema implements a handler for the /schema endpoint, GET method.
func (s *Server[Configuration, State]) GetSchema(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	schemaResult, err := s.connector.GetSchema(s.configuration)
	if err != nil {
		writeError(w, logger, err)
		return
	}

	writeJson(w, logger, http.StatusOK, schemaResult)
}

// Query implements a handler for the /query endpoint, POST method that executes a query.
func (s *Server[Configuration, State]) Query(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Query", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, ctx, span, s.telemetry.queryCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)

	execQueryCtx, execQuerySpan := s.telemetry.Tracer.Start(ctx, "Execute Query")
	defer execQuerySpan.End()

	response, err := s.connector.Query(execQueryCtx, s.configuration, s.state, &body)

	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			failureStatusAttribute,
			httpStatusAttribute(status),
		}
		span.SetAttributes(append(statusAttributes, attribute.String("reason", err.Error()))...)
		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(append(statusAttributes, collectionAttr)...))
		return
	}
	execQuerySpan.End()

	span.SetAttributes(successStatusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()

	s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(collectionAttr, successStatusAttribute))
	// record latency for success requests only
	s.telemetry.queryLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// QueryExplain implements a handler for the /query/explain endpoint, POST method that explains a query by creating an execution plan.
func (s *Server[Configuration, State]) QueryExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Query Explain", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, ctx, span, s.telemetry.queryExplainCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Explain")
	defer execSpan.End()

	response, err := s.connector.QueryExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			failureStatusAttribute,
			httpStatusAttribute(status),
		}
		span.SetAttributes(append(statusAttributes, attribute.String("reason", err.Error()))...)
		s.telemetry.queryExplainCounter.Add(r.Context(), 1, metric.WithAttributes(append(statusAttributes, collectionAttr)...))
		return
	}
	execSpan.End()

	span.SetAttributes(successStatusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()
	s.telemetry.queryExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, collectionAttr))

	// record latency for success requests only
	s.telemetry.queryExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// MutationExplain implements a handler for the /mutation/explain endpoint, POST method that explains a mutation by creating an execution plan.
func (s *Server[Configuration, State]) MutationExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Mutation Explain", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, ctx, span, s.telemetry.mutationExplainCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Explain")
	defer execSpan.End()

	response, err := s.connector.MutationExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			failureStatusAttribute,
			httpStatusAttribute(status),
		}
		span.SetAttributes(append(statusAttributes, attribute.String("reason", err.Error()))...)
		s.telemetry.mutationExplainCounter.Add(r.Context(), 1, metric.WithAttributes(append(statusAttributes, operationAttr)...))
		return
	}
	execSpan.End()

	span.SetAttributes(successStatusAttribute)

	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()
	s.telemetry.mutationExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// Mutation implements a handler for the /mutation endpoint, POST method that executes a mutation.
func (s *Server[Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(r.Context(), "Mutation", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, ctx, span, s.telemetry.mutationCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "Execute Mutation")
	defer execSpan.End()
	response, err := s.connector.Mutation(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		statusAttributes := []attribute.KeyValue{
			failureStatusAttribute,
			httpStatusAttribute(status),
		}
		span.SetAttributes(append(statusAttributes, attribute.String("reason", err.Error()))...)
		s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(append(statusAttributes, operationAttr)...))
		return
	}
	execSpan.End()

	span.SetAttributes(successStatusAttribute)
	_, responseSpan := s.telemetry.Tracer.Start(ctx, "Response")
	writeJson(w, logger, http.StatusOK, response)
	responseSpan.End()

	s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// the common unmarshal json body method
func (s *Server[Configuration, State]) unmarshalBodyJSON(w http.ResponseWriter, r *http.Request, ctx context.Context, span trace.Span, counter metric.Int64Counter, body any) error {
	_, decodeSpan := s.telemetry.Tracer.Start(ctx, "Decode JSON Body")
	defer decodeSpan.End()
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, GetLogger(r.Context()), http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		attributes := []attribute.KeyValue{
			failureStatusAttribute,
			httpStatusAttribute(http.StatusBadRequest),
		}
		span.SetAttributes(append(attributes, attribute.String("status", "json_decode"))...)
		counter.Add(r.Context(), 1, metric.WithAttributes(attributes...))
		return err
	}

	return nil
}

func (s *Server[Configuration, State]) buildHandler() *http.ServeMux {
	router := newRouter(s.logger, !s.withoutRecovery)
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

// BuildTestServer builds an http test server for testing purpose
func (s *Server[Configuration, State]) BuildTestServer() *httptest.Server {
	_ = s.telemetry.Shutdown(context.Background())
	return httptest.NewServer(s.buildHandler())
}

// ListenAndServe serves the configuration server with the standard http server.
// You can also replace this method with any router or web framework that is compatible with net/http.
func (s *Server[Configuration, State]) ListenAndServe(port uint) error {
	defer s.stop()
	defer func() {
		if err := s.telemetry.Shutdown(context.Background()); err != nil {
			s.logger.Error(
				"failed to shutdown OpenTelemetry",
				slog.Any("error", err),
			)
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
		s.logger.Info(fmt.Sprintf("Listening server on %s", server.Addr))
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
		s.logger.Info("received the quit signal, exiting...")
		// Stop receiving signal notifications as soon as possible.
		s.stop()
		// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
		return server.Shutdown(context.Background())
	}
}
