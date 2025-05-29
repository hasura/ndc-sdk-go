package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"slices"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// ServerOptions presents the configuration object of the connector http server.
type ServerOptions struct {
	OTLPConfig
	HTTPServerConfig

	Configuration      string
	InlineConfig       bool
	ServiceTokenSecret string
}

// HTTPServerConfig the configuration of the HTTP server.
type HTTPServerConfig struct {
	ServerReadTimeout            time.Duration `env:"HASURA_SERVER_READ_TIMEOUT"              help:"Maximum duration for reading the entire request, including the body. A zero or negative value means there will be no timeout"`
	ServerReadHeaderTimeout      time.Duration `env:"HASURA_SERVER_READ_HEADER_TIMEOUT"       help:"Amount of time allowed to read request headers. If zero, the value of ReadTimeout is used"`
	ServerWriteTimeout           time.Duration `env:"HASURA_SERVER_WRITE_TIMEOUT"             help:"Maximum duration before timing out writes of the response. A zero or negative value means there will be no timeout"`
	ServerIdleTimeout            time.Duration `env:"HASURA_SERVER_IDLE_TIMEOUT"              help:"Maximum amount of time to wait for the next request when keep-alives are enabled. If zero, the value of ReadTimeout is used"`
	ServerMaxHeaderKilobytes     int           `env:"HASURA_SERVER_MAX_HEADER_KILOBYTES"      help:"Maximum number of kilobytes the server will read parsing the request header's keys and values, including the request line"    default:"1024"`
	ServerMaxBodyMegabytes       int           `env:"HASURA_SERVER_MAX_BODY_MEGABYTES"        help:"Maximum size of the request body in megabytes that the server accepts"                                                        default:"30"`
	ServerDefaultHTTPErrorStatus int           `env:"HASURA_SERVER_DEFAULT_HTTP_ERROR_STATUS" help:"Default HTTP status code if the error does not specify the explicit status"                                                   default:"422"  enum:"400,422,500"`
	ServerTLSCertFile            string        `env:"HASURA_SERVER_TLS_CERT_FILE"             help:"Path of the TLS certificate file"`
	ServerTLSKeyFile             string        `env:"HASURA_SERVER_TLS_KEY_FILE"              help:"Path of the TLS key file"`
}

// Validate checks valid configurations and set default values.
func (hsc *HTTPServerConfig) Validate() error {
	if hsc.ServerMaxBodyMegabytes <= 0 {
		hsc.ServerMaxBodyMegabytes = 30
	}

	allowedHttpErrorCodes := []int{
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	}

	if hsc.ServerDefaultHTTPErrorStatus == 0 {
		hsc.ServerDefaultHTTPErrorStatus = http.StatusUnprocessableEntity
	} else if !slices.Contains(allowedHttpErrorCodes, hsc.ServerDefaultHTTPErrorStatus) {
		return fmt.Errorf("invalid default server http error code, accepted one of %v, got: %d", allowedHttpErrorCodes, hsc.ServerDefaultHTTPErrorStatus)
	}

	return nil
}

// Server implements the [NDC API specification] for the connector
//
// [NDC API specification]: https://hasura.github.io/ndc-spec/specification/index.html
type Server[Configuration any, State any] struct {
	*serveOptions

	context              context.Context //nolint:containedctx
	stop                 context.CancelFunc
	connector            Connector[Configuration, State]
	state                *State
	configuration        *Configuration
	options              *ServerOptions
	telemetry            *TelemetryState
	ndcVersionConstraint *semver.Version
	maxBodySize          int64
}

// NewServer creates a Server instance.
func NewServer[Configuration any, State any](
	connector Connector[Configuration, State],
	options *ServerOptions,
	others ...ServeOption,
) (*Server[Configuration, State], error) {
	defaultOptions := defaultServeOptions()

	for _, opts := range others {
		opts(defaultOptions)
	}

	if options.ServiceName == "" {
		options.ServiceName = defaultOptions.serviceName
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	defaultOptions.logger.Debug(
		"initialize OpenTelemetry",
		slog.Any("otlp", options.OTLPConfig),
		slog.String("version", defaultOptions.version),
		slog.String("metrics_prefix", defaultOptions.metricsPrefix),
	)

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt)

	telemetry, err := setupOTelSDK(
		ctx,
		&options.OTLPConfig,
		defaultOptions.version,
		defaultOptions.metricsPrefix,
		defaultOptions.logger,
	)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, logContextKey, telemetry.Logger)

	configuration, err := connector.ParseConfiguration(ctx, options.Configuration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(ctx, configuration, telemetry)
	if err != nil {
		return nil, err
	}

	ndcVersionConstraint, err := semver.NewVersion(schema.NDCVersion)
	if err != nil {
		return nil, err
	}

	return &Server[Configuration, State]{
		context:              ctx,
		stop:                 stop,
		connector:            connector,
		state:                state,
		configuration:        configuration,
		options:              options,
		telemetry:            telemetry,
		serveOptions:         defaultOptions,
		ndcVersionConstraint: ndcVersionConstraint,
		maxBodySize:          int64(options.ServerMaxBodyMegabytes) * 1024 * 1024,
	}, nil
}

// GetCapabilities get the connector's capabilities. Implement a handler for the /capabilities endpoint, GET method.
func (s *Server[Configuration, State]) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())

	capabilities := s.connector.GetCapabilities(s.configuration)
	if capabilities == nil {
		s.writeError(w, logger, schema.InternalServerError("capabilities is empty", nil))

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
		s.writeError(w, logger, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSchema implements a handler for the /schema endpoint, GET method.
func (s *Server[Configuration, State]) GetSchema(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())

	schemaResult, err := s.connector.GetSchema(r.Context(), s.configuration, s.state)
	if err != nil {
		s.writeError(w, logger, err)

		return
	}

	if schemaResult == nil {
		s.writeError(w, logger, schema.InternalServerError("schema is empty", nil))

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
		status := s.writeError(w, logger, err)
		s.telemetry.queryCounter.Add(r.Context(), 1, metric.WithAttributes(
			collectionAttr,
			failureStatusAttribute,
			httpStatusAttribute(status),
		))

		return
	}

	execQuerySpan.End()

	writeJson(w, logger, http.StatusOK, response)

	s.telemetry.queryCounter.Add(
		r.Context(),
		1,
		metric.WithAttributes(collectionAttr, successStatusAttribute),
	)
	// record latency for success requests only
	s.telemetry.queryLatencyHistogram.Record(
		r.Context(),
		time.Since(startTime).Seconds(),
		metric.WithAttributes(collectionAttr),
	)
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
		status := s.writeError(w, logger, err)
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
	s.telemetry.queryExplainCounter.Add(
		r.Context(),
		1,
		metric.WithAttributes(successStatusAttribute, collectionAttr),
	)

	// record latency for success requests only
	s.telemetry.queryExplainLatencyHistogram.Record(
		r.Context(),
		time.Since(startTime).Seconds(),
		metric.WithAttributes(collectionAttr),
	)
}

// MutationExplain implements a handler for the /mutation/explain endpoint,
// POST method that explains a mutation by creating an execution plan.
func (s *Server[Configuration, State]) MutationExplain( //nolint:dupl
	w http.ResponseWriter,
	r *http.Request,
) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())

	var body schema.MutationRequest

	if err := s.unmarshalBodyJSON(w, r, s.telemetry.mutationExplainCounter, &body); err != nil {
		return
	}

	operationNames := make([]string, len(body.Operations))

	for i, op := range body.Operations {
		operationNames[i] = op.Name
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_plan")
	defer execSpan.End()

	response, err := s.connector.MutationExplain(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := s.writeError(w, logger, err)

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
	s.telemetry.mutationExplainCounter.Add(
		r.Context(),
		1,
		metric.WithAttributes(successStatusAttribute, operationAttr),
	)

	// record latency for success requests only
	s.telemetry.mutationExplainLatencyHistogram.Record(
		r.Context(),
		time.Since(startTime).Seconds(),
		metric.WithAttributes(operationAttr),
	)
}

// Mutation implements a handler for the /mutation endpoint, POST method that executes a mutation.
func (s *Server[Configuration, State]) Mutation( //nolint:dupl
	w http.ResponseWriter,
	r *http.Request,
) {
	startTime := time.Now()
	logger := GetLogger(r.Context())
	span := trace.SpanFromContext(r.Context())

	var body schema.MutationRequest

	if err := s.unmarshalBodyJSON(w, r, s.telemetry.mutationCounter, &body); err != nil {
		return
	}

	operationNames := make([]string, len(body.Operations))

	for i, op := range body.Operations {
		operationNames[i] = op.Name
	}

	operationAttr := attribute.String("operations", strings.Join(operationNames, ","))
	span.SetAttributes(operationAttr)

	execCtx, execSpan := s.telemetry.Tracer.Start(r.Context(), "ndc_execute_mutation")
	defer execSpan.End()

	response, err := s.connector.Mutation(execCtx, s.configuration, s.state, &body)
	if err != nil {
		status := s.writeError(w, logger, err)
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

	s.telemetry.mutationCounter.Add(
		r.Context(),
		1,
		metric.WithAttributes(successStatusAttribute, operationAttr),
	)

	// record latency for success requests only
	s.telemetry.mutationLatencyHistogram.Record(
		r.Context(),
		time.Since(startTime).Seconds(),
		metric.WithAttributes(operationAttr),
	)
}

// the common unmarshal json body method.
func (s *Server[Configuration, State]) unmarshalBodyJSON(
	w http.ResponseWriter,
	r *http.Request,
	counter metric.Int64Counter,
	body any,
) error {
	// Validate the max body size. In the worst scenario, if the Content-Length header doesn't exist,
	// the server will validate again after reading the body.
	if r.ContentLength > 0 {
		if err := s.validateMaxBodySize(w, r, counter, r.ContentLength); err != nil {
			return err
		}
	}

	requestReader := &sizeReader{
		reader: r.Body,
	}

	err := json.NewDecoder(requestReader).Decode(body)
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

		return err
	}

	return s.validateMaxBodySize(w, r, counter, requestReader.size)
}

func (s *Server[Configuration, State]) validateMaxBodySize(
	w http.ResponseWriter,
	r *http.Request,
	counter metric.Int64Counter,
	contentLength int64,
) error {
	if contentLength <= s.maxBodySize {
		return nil
	}

	err := fmt.Errorf("request body size exceeded %d MB(s)", s.options.ServerMaxBodyMegabytes)

	writeJson(w, GetLogger(r.Context()), http.StatusUnprocessableEntity, schema.ErrorResponse{
		Message: err.Error(),
		Details: map[string]any{},
	})

	counter.Add(r.Context(), 1, metric.WithAttributes(
		failureStatusAttribute,
		httpStatusAttribute(http.StatusUnprocessableEntity),
	))

	return err
}

func (s *Server[Configuration, State]) writeError(
	w http.ResponseWriter,
	logger *slog.Logger,
	err error,
) int {
	return writeError(w, logger, err, s.options.ServerDefaultHTTPErrorStatus)
}

func (s *Server[Configuration, State]) buildHandler() *http.ServeMux {
	middlewares := []http.HandlerFunc{s.withNDCVersionCheck}
	if s.options.ServiceTokenSecret != "" {
		middlewares = append(middlewares, s.withAuth)
	}

	router := newRouter(s.telemetry.Logger, s.telemetry, !s.withoutRecovery)

	router.Use(apiPathCapabilities, http.MethodGet, append(middlewares, s.GetCapabilities)...)
	router.Use(apiPathSchema, http.MethodGet, append(middlewares, s.GetSchema)...)
	router.Use(apiPathQuery, http.MethodPost, append(middlewares, s.Query)...)
	router.Use(apiPathQueryExplain, http.MethodPost, append(middlewares, s.QueryExplain)...)
	router.Use(apiPathMutationExplain, http.MethodPost, append(middlewares, s.MutationExplain)...)
	router.Use(apiPathMutation, http.MethodPost, append(middlewares, s.Mutation)...)
	router.Use(apiPathHealth, http.MethodGet, s.Health)

	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) &&
		s.options.PrometheusPort == nil {
		router.Use(apiPathMetrics, http.MethodGet, s.withAuth, promhttp.Handler().ServeHTTP)
	}

	return router.Build()
}

// BuildTestServer builds an http test server for testing purpose.
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
	if s.options.ServerMaxHeaderKilobytes > 0 {
		maxHeaderBytes = s.options.ServerMaxHeaderKilobytes * 1024
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
		var err error

		if s.options.ServerTLSCertFile != "" || s.options.ServerTLSKeyFile != "" {
			s.telemetry.Logger.Info("Listening server and serving TLS on " + server.Addr)
			err = server.ListenAndServeTLS(s.options.ServerTLSCertFile, s.options.ServerTLSKeyFile)
		} else {
			s.telemetry.Logger.Info("Listening server on " + server.Addr)
			err = server.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	if s.options.MetricsExporter == string(otelMetricsExporterPrometheus) &&
		s.options.PrometheusPort != nil {
		promServer := createPrometheusServer(*s.options.PrometheusPort)
		defer func() {
			_ = promServer.Shutdown(context.Background())
		}()

		go func() {
			s.telemetry.Logger.Info(
				fmt.Sprintf("Listening prometheus server on %d", *s.options.PrometheusPort),
			)

			if err := promServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
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
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
	}
}

// A reader wrapper to estimate the size when decoding with json.Decoder.
type sizeReader struct {
	size   int64
	reader io.Reader
}

func (sr *sizeReader) Read(p []byte) (int, error) {
	size, err := sr.reader.Read(p)
	sr.size += int64(size)

	return size, err
}
