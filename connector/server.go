package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
)

var (
	errConfigurationRequired = errors.New("Configuration is required")
)

type ServerOptions struct {
	Configuration      string
	ConfigurationJSON  string
	ServiceTokenSecret string
	OTLPEndpoint       string
	ServiceName        string
	Logger             zerolog.Logger
}

// Server implements the [NDC API specification] for the connector
//
// [NDC API specification]: https://hasura.github.io/ndc-spec/specification/index.html
type Server[RawConfiguration any, Configuration any, State any] struct {
	connector     Connector[RawConfiguration, Configuration, State]
	state         *State
	configuration *Configuration
	options       *ServerOptions
}

// NewServer creates a Server instance
func NewServer[RawConfiguration any, Configuration any, State any](connector Connector[RawConfiguration, Configuration, State], options *ServerOptions) (*Server[RawConfiguration, Configuration, State], error) {
	var rawConfiguration RawConfiguration
	configBytes := []byte(options.ConfigurationJSON)
	if options.ConfigurationJSON == "" {
		if options.Configuration == "" {
			return nil, errConfigurationRequired
		}
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

	configuration, err := connector.ValidateRawConfiguration(&rawConfiguration)
	if err != nil {
		return nil, err
	}

	state, err := connector.TryInitState(configuration, nil)
	if err != nil {
		return nil, err
	}

	return &Server[RawConfiguration, Configuration, State]{
		connector:     connector,
		state:         state,
		configuration: configuration,
		options:       options,
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

func (s *Server[RawConfiguration, Configuration, State]) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: return Prometheus metrics
	if err := s.connector.FetchMetrics(r.Context(), s.configuration, s.state); err != nil {
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
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})
		return
	}

	response, err := cs.connector.Query(r.Context(), cs.configuration, cs.state, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
}

func (cs *Server[RawConfiguration, Configuration, State]) Explain(w http.ResponseWriter, r *http.Request) {
	var body schema.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})
		return
	}

	response, err := cs.connector.Explain(r.Context(), cs.configuration, cs.state, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
}

func (cs *Server[RawConfiguration, Configuration, State]) Mutation(w http.ResponseWriter, r *http.Request) {
	var body schema.MutationRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})
		return
	}

	response, err := cs.connector.Mutation(r.Context(), cs.configuration, cs.state, &body)
	if err != nil {
		writeError(w, err)
		return
	}

	internal.WriteJson(w, http.StatusOK, response)
}

// ListenAndServe serves the configuration server with the standard http server.
// You can also replace this method with any router or web framework that is compatible with net/http.
func (cs *Server[RawConfiguration, Configuration, State]) ListenAndServe(port uint) error {
	router := internal.NewRouter(cs.options.Logger)
	router.Use("/capabilities", http.MethodGet, cs.GetCapabilities)
	router.Use("/schema", http.MethodGet, cs.GetSchema)
	router.Use("/query", http.MethodPost, cs.Query)
	router.Use("/explain", http.MethodPost, cs.Explain)
	router.Use("/mutation", http.MethodPost, cs.Mutation)
	router.Use("/healthz", http.MethodGet, cs.Health)
	router.Use("/metrics", http.MethodGet, cs.GetMetrics)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.Build(),
	}

	cs.options.Logger.Info().Msgf("Listening server on %s", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	var connectorError schema.ConnectorError
	if errors.As(err, &connectorError) {
		internal.WriteJson(w, connectorError.StatusCode(), connectorError)
		return
	}

	internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
		Message: err.Error(),
	})
}
