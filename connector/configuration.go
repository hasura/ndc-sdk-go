package connector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
)

// ConfigurationServer provides configuration APIs of the connector
// that help validate and generate Hasura v3 metadata
type ConfigurationServer[RawConfiguration any, Configuration any, State any] struct {
	connector Connector[RawConfiguration, Configuration, State]
	logger    zerolog.Logger
}

// NewConfigurationServer creates a ConfigurationServer instance
func NewConfigurationServer[RawConfiguration any, Configuration any, State any](connector Connector[RawConfiguration, Configuration, State], options ...ServeOption) *ConfigurationServer[RawConfiguration, Configuration, State] {
	defaultOptions := defaultServeOptions()
	for _, opts := range options {
		opts(defaultOptions)
	}

	return &ConfigurationServer[RawConfiguration, Configuration, State]{
		connector: connector,
		logger:    defaultOptions.logger,
	}
}

// GetIndex implements a handler for the index endpoint, GET method.
// Returns an empty configuration of the connector
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) GetIndex(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	internal.WriteJson(w, http.StatusOK, cs.connector.MakeEmptyConfiguration())
}

// PostIndex implements a handler for the index endpoint, POST method.
// Take a raw configuration, update it where appropriate by connecting to the underlying data source, and otherwise return it as-is
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) PostIndex(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var rawConfig RawConfiguration
	if err := json.NewDecoder(r.Body).Decode(&rawConfig); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.BadRequestError(err.Error(), nil))
		return
	}

	conf, err := cs.connector.UpdateConfiguration(r.Context(), &rawConfig)
	if err != nil {
		writeError(w, err)
		return
	}
	internal.WriteJson(w, http.StatusOK, conf)
}

// GetSchema implements a handler for the /schema endpoint, GET method.
// Return jsonschema for the raw configuration for this connector
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) GetSchema(w http.ResponseWriter, r *http.Request) {
	internal.WriteJson(w, http.StatusOK, cs.connector.GetRawConfigurationSchema())
}

// Validate implements a handler for the /validate endpoint, POST method.
// that validates the raw configuration provided by the user
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) Validate(w http.ResponseWriter, r *http.Request) {
	var rawConfig RawConfiguration
	if err := json.NewDecoder(r.Body).Decode(&rawConfig); err != nil {
		internal.WriteJson(w, http.StatusBadRequest, schema.ErrorResponse{
			Message: "failed to decode json request body",
			Details: map[string]any{
				"cause": err.Error(),
			},
		})
		return
	}

	resolvedConfiguration, err := cs.connector.ValidateRawConfiguration(
		&rawConfig,
	)
	if err != nil {
		writeError(w, err)
		return
	}

	connectorSchema, err := cs.connector.GetSchema(resolvedConfiguration)
	if err != nil {
		writeError(w, err)
		return
	}

	capabilities := cs.connector.GetCapabilities(resolvedConfiguration)
	configurationBytes, err := json.Marshal(resolvedConfiguration)
	if err != nil {
		writeError(w, schema.InternalServerError(err.Error(), nil))
		return
	}

	internal.WriteJson(w, http.StatusOK, &schema.ValidateResponse{
		Schema:                *connectorSchema,
		Capabilities:          *capabilities,
		ResolvedConfiguration: string(configurationBytes),
	})
}

// Health implements a handler for /healthz endpoint.
// The endpoint has nothing to check, because the reference implementation does not need to connect to any other services.
// Therefore, once the reference implementation is running, it can always report a healthy status
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// ListenAndServe serves the configuration server with the standard http server.
// You can also replace this method with any router or web framework that is compatible with net/http.
func (cs *ConfigurationServer[RawConfiguration, Configuration, State]) ListenAndServe(port uint) error {
	router := internal.NewRouter(cs.logger)
	router.Use("/", http.MethodGet, cs.GetIndex)
	router.Use("/", http.MethodPost, cs.PostIndex)
	router.Use("/schema", http.MethodGet, cs.GetSchema)
	router.Use("/validate", http.MethodPost, cs.Validate)
	router.Use("/healthz", http.MethodGet, cs.Health)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.Build(),
	}

	cs.logger.Info().Msgf("Listening server on %s", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
