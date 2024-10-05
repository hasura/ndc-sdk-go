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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// ServerOptions presents the configuration object of the connector http server
type ServerOptions struct {
	OTLPConfig
	HTTPServerConfig

	Configuration      string
	InlineConfig       bool
	ServiceTokenSecret string
}

// HTTPServerConfig the configuration of the HTTP server
type HTTPServerConfig struct {
	ServerReadTimeout        time.Duration `help:"The maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout" env:"HASURA_SERVER_READ_TIMEOUT"`
	ServerReadHeaderTimeout  time.Duration `help:"The amount of time allowed to read request headers. If zero, the value of ReadTimeout is used" env:"HASURA_SERVER_READ_HEADER_TIMEOUT"`
	ServerWriteTimeout       time.Duration `help:"The maximum duration before timing out writes of the response. A zero or negative value means there will be no timeout" env:"HASURA_SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout        time.Duration `help:"The maximum amount of time to wait for the next request when keep-alives are enabled. If zero, the value of ReadTimeout is used" env:"HASURA_SERVER_IDLE_TIMEOUT"`
	ServerMaxHeaderKilobytes int           `help:"The maximum number of kilobytes the server will read parsing the request header's keys and values, including the request line" default:"1024" env:"HASURA_SERVER_MAX_HEADER_KILOBYTES"`
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
	if options.ServiceName == "" {
		options.OTLPConfig.ServiceName = defaultOptions.serviceName
	}
	defaultOptions.logger.Debug(
		"initialize OpenTelemetry",
		slog.Any("otlp", options.OTLPConfig),
		slog.String("version", defaultOptions.version),
		slog.String("metrics_prefix", defaultOptions.metricsPrefix),
	)

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.WithValue(context.TODO(), logContextKey, defaultOptions.logger), os.Interrupt)

	configuration, err := connector.ParseConfiguration(ctx, options.Configuration)
	if err != nil {
		return nil, err
	}

	telemetry, err := setupOTelSDK(ctx, &options.OTLPConfig, defaultOptions.version, defaultOptions.metricsPrefix, defaultOptions.logger)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(ctx, configuration, telemetry)
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
	if capabilities == nil {
		writeError(w, logger, schema.InternalServerError("capabilities is empty", nil))
		return
	}
	writeJsonFunc(w, logger, http.StatusOK, func() ([]byte, error) {
		return capabilities.MarshalCapabilitiesJSON()
	})
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
	schemaResult, err := s.connector.GetSchema(r.Context(), s.configuration, s.state)
	if err != nil {
		writeError(w, logger, err)
		return
	}
	if schemaResult == nil {
		writeError(w, logger, schema.InternalServerError("schema is empty", nil))
		return
	}

	writeJsonFunc(w, logger, http.StatusOK, func() ([]byte, error) {
		return schemaResult.MarshalSchemaJSON()
	})
}

// Query implements a handler for the /query endpoint, POST method that executes a query.
func (s *Server[Configuration, State]) Query(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())
	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, s.telemetry.queryCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)
	execQueryCtx, execQuerySpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_query")
	defer execQuerySpan.End()

	response, err := s.connector.Query(execQueryCtx, s.configuration, s.state, &body)

	if err != nil {
		status := writeError(w, logger, err)
		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(
			collectionAttr,
			failureStatusAttribute,
			httpStatusAttribute(status),
		))
		return
	}
	execQuerySpan.End()

	writeJson(w, logger, http.StatusOK, response)

	s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(collectionAttr, successStatusAttribute))
	// record latency for success requests only
	s.telemetry.queryLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// QueryExplain implements a handler for the /query/explain endpoint, POST method that explains a query by creating an execution plan.
func (s *Server[Configuration, State]) QueryExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())

	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, s.telemetry.queryExplainCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_plan")
	defer execSpan.End()

	response, err := s.connector.QueryExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		s.telemetry.queryExplainCounter.Add(r.Context(), 1, metric.WithAttributes(
			failureStatusAttribute,
			httpStatusAttribute(status),
			collectionAttr,
		))
		return
	}
	execSpan.End()

	writeJson(w, logger, http.StatusOK, response)
	s.telemetry.queryExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, collectionAttr))

	// record latency for success requests only
	s.telemetry.queryExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// MutationExplain implements a handler for the /mutation/explain endpoint, POST method that explains a mutation by creating an execution plan.
func (s *Server[Configuration, State]) MutationExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())
	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, s.telemetry.mutationExplainCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_plan")
	defer execSpan.End()

	response, err := s.connector.MutationExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.telemetry.mutationExplainCounter.Add(r.Context(), 1, metric.WithAttributes(
			failureStatusAttribute,
			httpStatusAttribute(status),
			operationAttr,
		))
		return
	}
	execSpan.End()

	writeJson(w, logger, http.StatusOK, response)
	s.telemetry.mutationExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// Mutation implements a handler for the /mutation endpoint, POST method that executes a mutation.
func (s *Server[Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())

	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, s.telemetry.mutationCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_mutation")
	defer execSpan.End()
	response, err := s.connector.Mutation(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := writeError(w, logger, err)
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(
			failureStatusAttribute,
			httpStatusAttribute(status),
			operationAttr,
		))
		return
	}
	execSpan.End()

	writeJson(w, logger, http.StatusOK, response)

	s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// the common unmarshal json body method
func (s *Server[Configuration, State]) unmarshalBodyJSON(w http.ResponseWriter, r *http.Request, counter metric.Int64Counter, body any) error {
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		writeJson(w, GetLogger(r.Context()), http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		counter.Add(r.Context(), 1, metric.WithAttributes(
			failureStatusAttribute,
			httpStatusAttribute(http.StatusUnprocessableEntity),
		))
	}

	return err
}

func (s *Server[Configuration, State]) buildHandler() *http.ServeMux {
	router := newRouter(s.logger, s.telemetry, !s.withoutRecovery)
	router.Use(apiPathCapabilities, http.MethodGet, s.withAuth(s.GetCapabilities))
	router.Use(apiPathSchema, http.MethodGet, s.withAuth(s.GetSchema))
	router.Use(apiPathQuery, http.MethodPost, s.withAuth(s.Query))
	router.Use(apiPathQueryExplain, http.MethodPost, s.withAuth(s.QueryExplain))
	router.Use(apiPathMutationExplain, http.MethodPost, s.withAuth(s.MutationExplain))
	router.Use(apiPathMutation, http.MethodPost, s.withAuth(s.Mutation))
	router.Use(apiPathHealth, http.MethodGet, s.Health)
	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) && s.options.PrometheusPort == nil {
		router.Use(apiPathMetrics, http.MethodGet, s.withAuth(promhttp.Handler().ServeHTTP))
	}

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

	maxHeaderBytes := http.DefaultMaxHeaderBytes
	if s.options.HTTPServerConfig.ServerMaxHeaderKilobytes > 0 {
		maxHeaderBytes = s.options.HTTPServerConfig.ServerMaxHeaderKilobytes * 1024
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context {
			return s.context
		},
		Handler:           s.buildHandler(),
		ReadTimeout:       s.options.ServerReadTimeout,
		ReadHeaderTimeout: s.options.ServerReadHeaderTimeout,
		WriteTimeout:      s.options.ServerWriteTimeout,
		IdleTimeout:       s.options.ServerIdleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	serverErr := make(chan error, 1)
	go func() {
		s.logger.Info(fmt.Sprintf("Listening server on %s", server.Addr))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) && s.options.PrometheusPort != nil {
		promServer := createPrometheusServer(*s.options.PrometheusPort)
		defer func() {
			_ = promServer.Shutdown(context.Background())
		}()
		go func() {
			s.logger.Info(fmt.Sprintf("Listening prometheus server on %d", *s.options.PrometheusPort))
			if err := promServer.ListenAndServe(); err != http.ErrServerClosed {
				serverErr <- err
			}
		}()
	}

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

func createPrometheusServer(port uint) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}
