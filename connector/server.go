package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// ServerOptions presents the configuration object of the connector http server
type ServerOptions struct {
	OTLPConfig

	Configuration      string
	InlineConfig       bool
	ServiceTokenSecret string
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
		options.ServiceName = defaultOptions.serviceName
	}
	defaultOptions.logger.Debug().
		Any("otlp", options.OTLPConfig).
		Str("version", defaultOptions.version).
		Str("metrics_prefix", defaultOptions.metricsPrefix).
		Msg("initialize OpenTelemetry")

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.WithValue(context.TODO(), logContextKey, defaultOptions.logger), os.Interrupt)

	configuration, err := connector.ParseConfiguration(ctx, options.Configuration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(ctx, configuration, nil)
	if err != nil {
		return nil, err
	}

	telemetry, err := setupOTelSDK(ctx, &options.OTLPConfig, defaultOptions.version, defaultOptions.metricsPrefix, defaultOptions.logger)
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
			_ = writeJson(w, logger, http.StatusUnauthorized, schema.ErrorResponse{
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
	ctx, span := s.telemetry.Tracer.Start(
		otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
		"ndc_query",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, span, s.telemetry.queryCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)
	span.AddEvent("execute_query", trace.WithAttributes(
		collectionAttr,
	))
	execQueryCtx, execQuerySpan := s.telemetry.Tracer.Start(ctx, "ndc_execute_query")
	defer execQuerySpan.End()

	response, err := s.connector.Query(execQueryCtx, s.configuration, s.state, &body)

	if err != nil {
		status := writeError(w, logger, err)
		execQuerySpan.SetStatus(codes.Error, err.Error())
		execQuerySpan.RecordError(err)

		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(
			collectionAttr,
			failureStatusAttribute,
			httpStatusAttribute(status),
		))
		return
	}
	execQuerySpan.End()

	span.AddEvent("ndc_query_response")
	if err := writeJson(w, logger, http.StatusOK, response); err != nil {
		span.SetStatus(codes.Error, "failed to write response")
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "executed query successfully")
	}

	s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(collectionAttr, successStatusAttribute))
	// record latency for success requests only
	s.telemetry.queryLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// QueryExplain implements a handler for the /query/explain endpoint, POST method that explains a query by creating an execution plan.
func (s *Server[Configuration, State]) QueryExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(
		otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
		"ndc_explain_query",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	var body schema.QueryRequest
	if err := s.unmarshalBodyJSON(w, r, span, s.telemetry.queryExplainCounter, &body); err != nil {
		return
	}

	collectionAttr := attribute.String("collection", body.Collection)
	span.SetAttributes(collectionAttr)

	span.AddEvent("execute_query_plain", trace.WithAttributes(
		collectionAttr,
	))
	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "ndc_execute_plan")
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

	span.AddEvent("query_explain_response")
	if err := writeJson(w, logger, http.StatusOK, response); err != nil {
		span.SetStatus(codes.Error, "failed to write response")
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "explained query successfully")
	}
	s.telemetry.queryExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, collectionAttr))

	// record latency for success requests only
	s.telemetry.queryExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(collectionAttr))
}

// MutationExplain implements a handler for the /mutation/explain endpoint, POST method that explains a mutation by creating an execution plan.
func (s *Server[Configuration, State]) MutationExplain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(
		otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
		"ndc_mutation_explain",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, span, s.telemetry.mutationExplainCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	span.AddEvent("execute_mutation_plain", trace.WithAttributes(
		operationAttr,
	))
	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "ndc_execute_plan")
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

	span.AddEvent("mutation_explain_response")
	if err := writeJson(w, logger, http.StatusOK, response); err != nil {
		span.SetStatus(codes.Error, "failed to write response")
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "explained mutation successfully")
	}
	s.telemetry.mutationExplainCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationExplainLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// Mutation implements a handler for the /mutation endpoint, POST method that executes a mutation.
func (s *Server[Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	ctx, span := s.telemetry.Tracer.Start(
		otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
		"ndc_mutation",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	var body schema.MutationRequest
	if err := s.unmarshalBodyJSON(w, r, span, s.telemetry.mutationCounter, &body); err != nil {
		return
	}

	var operationNames []string
	for _, op := range body.Operations {
		operationNames = append(operationNames, op.Name)
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	span.AddEvent("execute_mutation", trace.WithAttributes(
		operationAttr,
	))
	execCtx, execSpan := s.telemetry.Tracer.Start(ctx, "ndc_execute_mutation")
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

	span.AddEvent("mutation_response")
	if err := writeJson(w, logger, http.StatusOK, response); err != nil {
		span.SetStatus(codes.Error, "failed to write response")
		span.RecordError(err)
	} else {
		span.SetStatus(codes.Ok, "executed mutation successfully")
	}

	s.telemetry.mutationCounter.Add(r.Context(), 1, metric.WithAttributes(successStatusAttribute, operationAttr))

	// record latency for success requests only
	s.telemetry.mutationLatencyHistogram.Record(r.Context(), time.Since(startTime).Seconds(), metric.WithAttributes(operationAttr))
}

// the common unmarshal json body method
func (s *Server[Configuration, State]) unmarshalBodyJSON(w http.ResponseWriter, r *http.Request, span trace.Span, counter metric.Int64Counter, body any) error {
	span.AddEvent("decode_body_json")
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeJson(w, GetLogger(r.Context()), http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})

		span.SetStatus(codes.Error, "json_decode_error")
		span.RecordError(err)

		counter.Add(r.Context(), 1, metric.WithAttributes(
			failureStatusAttribute,
			httpStatusAttribute(http.StatusUnprocessableEntity),
		))
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
	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) && s.options.PrometheusPort == nil {
		router.Use("/metrics", http.MethodGet, s.withAuth(promhttp.Handler().ServeHTTP))
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

	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) && s.options.PrometheusPort != nil {
		promServer := createPrometheusServer(*s.options.PrometheusPort)
		defer func() {
			_ = promServer.Shutdown(context.Background())
		}()
		go func() {
			s.logger.Info().Msgf("Listening prometheus server on %d", *s.options.PrometheusPort)
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
		s.logger.Info().Msg("received the quit signal, exiting...")
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
