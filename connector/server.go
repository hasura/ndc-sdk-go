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
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	errConfigurationRequired = errors.New("Configuration is required")
)

type ServerOptions struct {
	Configuration      string
	InlineConfig       bool
	ServiceTokenSecret string
	OTLPEndpoint       string
	ServiceName        string
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

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	configuration, err := connector.ValidateRawConfiguration(&rawConfiguration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(configuration, nil)
	if err != nil {
		return nil, err
	}

	telemetry, err := setupOTelSDK(ctx, options.OTLPEndpoint, options.ServiceName, defaultOptions.version, defaultOptions.metricsPrefix)
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

func (s *Server[RawConfiguration, Configuration, State]) authorize(r *http.Request) error {
	if s.options.ServiceTokenSecret != "" && r.Header.Get("authorization") != fmt.Sprintf("Bearer %s", s.options.ServiceTokenSecret) {
		return schema.UnauthorizeError("Unauthorized", map[string]any{
			"cause": "Bearer token does not match.",
		})
	}

	return nil
}

func (s *Server[RawConfiguration, Configuration, State]) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	capacities := s.connector.GetCapabilities(s.configuration)
	internal.WriteJson(w, http.StatusOK, capacities)
}

func (s *Server[RawConfiguration, Configuration, State]) Health(w http.ResponseWriter, r *http.Request) {
	if err := s.connector.HealthCheck(r.Context(), s.configuration, s.state); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSchema implements a handler for the /schema endpoint, GET method.
func (cs *Server[RawConfiguration, Configuration, State]) GetSchema(w http.ResponseWriter, r *http.Request) {
	schemaResult, err := cs.connector.GetSchema(cs.configuration)
	if err != nil {
		writeError(w, err)
		return
	}

	internal.WriteJson(w, http.StatusOK, schemaResult)
}

func (cs *Server[RawConfiguration, Configuration, State]) Query(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx, span := cs.telemetry.Tracer.Start(r.Context(), "Query")
	defer span.End()
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
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
		cs.telemetry.queryCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	response, err := cs.connector.Query(ctx, cs.configuration, cs.state, &body)
	if err != nil {
		status := writeError(w, err)
		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		cs.telemetry.queryCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
	attributes := attribute.String("status", "success")
	span.SetAttributes(attributes)
	cs.telemetry.queryCounter.Add(ctx, 1, metric.WithAttributes(attributes))

	// record latency for success requests only
	cs.telemetry.queryLatencyHistogram.Record(ctx, time.Since(startTime).Seconds())
}

func (cs *Server[RawConfiguration, Configuration, State]) Explain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx, span := cs.telemetry.Tracer.Start(r.Context(), "Explain")
	defer span.End()
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
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
		cs.telemetry.explainCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	response, err := cs.connector.Explain(ctx, cs.configuration, cs.state, &body)
	if err != nil {
		status := writeError(w, err)
		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		cs.telemetry.explainCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
	attributes := attribute.String("status", "success")
	span.SetAttributes(attributes)
	cs.telemetry.explainCounter.Add(ctx, 1, metric.WithAttributes(attributes))

	// record latency for success requests only
	cs.telemetry.explainLatencyHistogram.Record(ctx, time.Since(startTime).Seconds())
}

func (cs *Server[RawConfiguration, Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	ctx, span := cs.telemetry.Tracer.Start(r.Context(), "Mutation")
	defer span.End()

	var body schema.MutationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
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
		cs.telemetry.mutationCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	response, err := cs.connector.Mutation(ctx, cs.configuration, cs.state, &body)
	if err != nil {
		status := writeError(w, err)
		attributes := []attribute.KeyValue{
			attribute.String("status", "failed"),
			attribute.String("reason", fmt.Sprintf("%d", status)),
		}
		span.SetAttributes(attributes...)
		cs.telemetry.mutationCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
	attributes := attribute.String("status", "success")
	span.SetAttributes(attributes)
	cs.telemetry.mutationCounter.Add(ctx, 1, metric.WithAttributes(attributes))

	// record latency for success requests only
	cs.telemetry.mutationLatencyHistogram.Record(ctx, time.Since(startTime).Seconds())
}

// ListenAndServe serves the configuration server with the standard http server.
// You can also replace this method with any router or web framework that is compatible with net/http.
func (cs *Server[RawConfiguration, Configuration, State]) ListenAndServe(port uint) error {
	defer cs.stop()
	defer cs.telemetry.Shutdown(context.Background())

	router := internal.NewRouter(cs.logger)
	router.Use("/capabilities", http.MethodGet, cs.GetCapabilities)
	router.Use("/schema", http.MethodGet, cs.GetSchema)
	router.Use("/query", http.MethodPost, cs.Query)
	router.Use("/explain", http.MethodPost, cs.Explain)
	router.Use("/mutation", http.MethodPost, cs.Mutation)
	router.Use("/healthz", http.MethodGet, cs.Health)
	router.Use("/metrics", http.MethodGet, promhttp.Handler().ServeHTTP)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context {
			return cs.context
		},
		Handler: router.Build(),
	}

	serverErr := make(chan error, 1)
	go func() {
		cs.logger.Info().Msgf("Listening server on %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interruption.
	select {
	case err := <-serverErr:
		// Error when starting HTTP server.
		return err
	case <-cs.context.Done():
		// Wait for first CTRL+C.
		cs.logger.Info().Msg("received the quit signal, exiting...")
		// Stop receiving signal notifications as soon as possible.
		cs.stop()
		// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
		return server.Shutdown(context.Background())
	}
}

func writeError(w http.ResponseWriter, err error) int {
	w.Header().Add("Content-Type", "application/json")
	var connectorError schema.ConnectorError
	if errors.As(err, &connectorError) {
		internal.WriteJson(w, connectorError.StatusCode(), connectorError)
		return connectorError.StatusCode()
	}

	internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
		Message: err.Error(),
	})

	return http.StatusBadRequest
}
