package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type serverContextKey string

const (
	logContextKey     serverContextKey = "hasura-log"
	headerContentType string           = "Content-Type"
	contentTypeJson   string           = "application/json"
)

// define a custom response write to capture response information for logging
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (cw *customResponseWriter) WriteHeader(statusCode int) {
	cw.statusCode = statusCode
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *customResponseWriter) Write(body []byte) (int, error) {
	cw.body = body
	return cw.ResponseWriter.Write(body)
}

// implements a simple router to reuse for both configuration and connector servers
type router struct {
	routes          map[string]map[string]http.HandlerFunc
	logger          zerolog.Logger
	recoveryEnabled bool
}

func newRouter(logger zerolog.Logger, enableRecovery bool) *router {
	return &router{
		routes:          make(map[string]map[string]http.HandlerFunc),
		logger:          logger,
		recoveryEnabled: enableRecovery,
	}
}

func (rt *router) Use(path string, method string, handler http.HandlerFunc) {
	if _, ok := rt.routes[path]; !ok {
		rt.routes[path] = make(map[string]http.HandlerFunc)
	}
	rt.routes[path][method] = handler
}

func (rt *router) Build() *http.ServeMux {
	mux := http.NewServeMux()

	handleFunc := func(handlers map[string]http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			isDebug := rt.logger.GetLevel() <= zerolog.DebugLevel
			requestID := getRequestID(r)
			requestLogData := map[string]any{
				"url":            r.URL.String(),
				"method":         r.Method,
				"remote_address": r.RemoteAddr,
			}

			if isDebug {
				requestLogData["headers"] = r.Header
				if r.Body != nil {
					bodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						rt.logger.Error().
							Str("request_id", requestID).
							Dur("latency", time.Since(startTime)).
							Interface("request", requestLogData).
							Interface("error", err).
							Msg("failed to read request")

						_ = writeJson(w, rt.logger, http.StatusUnprocessableEntity, schema.ErrorResponse{
							Message: "failed to read request",
							Details: map[string]any{
								"cause": err,
							},
						})
						return
					}

					requestLogData["body"] = string(bodyBytes)
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			// recover from panic
			if rt.recoveryEnabled {
				defer func() {
					if err := recover(); err != nil {
						rt.logger.Error().
							Str("request_id", requestID).
							Dur("latency", time.Since(startTime)).
							Interface("request", requestLogData).
							Interface("error", err).
							Str("stacktrace", string(debug.Stack())).
							Msg("internal server error")

						_ = writeJson(w, rt.logger, http.StatusInternalServerError, schema.ErrorResponse{
							Message: "internal server error",
							Details: map[string]any{
								"cause": err,
							},
						})
					}
				}()
			}

			h, ok := handlers[r.Method]
			if !ok {
				http.NotFound(w, r)
				rt.logger.Error().
					Str("request_id", requestID).
					Dur("latency", time.Since(startTime)).
					Interface("request", requestLogData).
					Interface("response", map[string]any{
						"status": 404,
					}).
					Msg("")
				return
			}

			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				contentType := r.Header.Get(headerContentType)
				if contentType != contentTypeJson {
					err := schema.ErrorResponse{
						Message: fmt.Sprintf("Invalid content type %s, accept %s only", contentType, contentTypeJson),
					}
					_ = writeJson(w, rt.logger, http.StatusUnprocessableEntity, err)

					rt.logger.Error().
						Str("request_id", requestID).
						Dur("latency", time.Since(startTime)).
						Interface("request", requestLogData).
						Interface("response", map[string]any{
							"status": 404,
							"body":   err,
						}).
						Msg("")
					return
				}
			}

			logger := rt.logger.With().Str("request_id", requestID).Logger()
			req := r.WithContext(context.WithValue(r.Context(), logContextKey, logger))
			writer := &customResponseWriter{ResponseWriter: w}
			h(writer, req)

			responseLogData := map[string]any{
				"status": writer.statusCode,
			}
			if isDebug || writer.statusCode >= 400 {
				responseLogData["headers"] = writer.Header()
				if len(writer.body) > 0 {
					responseLogData["body"] = string(writer.body)
				}
			}

			if writer.statusCode >= 400 {
				logger.Error().
					Dur("latency", time.Since(startTime)).
					Interface("request", requestLogData).
					Interface("response", responseLogData).Msg("")
			} else {
				logger.Info().
					Dur("latency", time.Since(startTime)).
					Interface("request", requestLogData).
					Interface("response", responseLogData).Msg("")
			}
		}
	}

	for path, handlers := range rt.routes {
		handler := handleFunc(handlers)
		mux.HandleFunc(path, handler)
	}

	return mux
}

func getRequestID(r *http.Request) string {
	requestID := r.Header.Get("x-request-id")
	if requestID == "" {
		requestID = internal.GenRandomString(16)
	}
	return requestID
}

// writeJsonFunc writes raw bytes data with a json encoding callback
func writeJsonFunc(w http.ResponseWriter, logger zerolog.Logger, statusCode int, encodeFunc func() ([]byte, error)) error {
	w.Header().Set(headerContentType, contentTypeJson)
	jsonBytes, err := encodeFunc()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, wErr := w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, http.StatusText(http.StatusInternalServerError))))
		if wErr != nil {
			logger.Error().Err(err).Msg("failed to write response")
		}
		return err
	}
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonBytes)
	if err != nil {
		logger.Error().Err(err).Msg("failed to write response")
	}
	return err
}

// writeJson writes response data with json encode
func writeJson(w http.ResponseWriter, logger zerolog.Logger, statusCode int, body any) error {
	if body == nil {
		w.WriteHeader(statusCode)
		return nil
	}

	return writeJsonFunc(w, logger, statusCode, func() ([]byte, error) {
		return json.Marshal(body)
	})
}

// GetLogger gets the logger instance from context
func GetLogger(ctx context.Context) zerolog.Logger {
	value := ctx.Value(logContextKey)
	if value != nil {
		if logger, ok := value.(zerolog.Logger); ok {
			return logger
		}
	}

	return log.Level(zerolog.GlobalLevel())
}

func writeError(w http.ResponseWriter, logger zerolog.Logger, err error) int {
	w.Header().Add("Content-Type", "application/json")

	var connectorErrorPtr *schema.ConnectorError
	if errors.As(err, &connectorErrorPtr) {
		_ = writeJson(w, logger, connectorErrorPtr.StatusCode(), connectorErrorPtr)
		return connectorErrorPtr.StatusCode()
	}

	var errorResponse schema.ErrorResponse
	if errors.As(err, &errorResponse) {
		_ = writeJson(w, logger, http.StatusBadRequest, errorResponse)
		return http.StatusBadRequest
	}

	var errorResponsePtr *schema.ErrorResponse
	if errors.As(err, &errorResponsePtr) {
		_ = writeJson(w, logger, http.StatusBadRequest, errorResponsePtr)
		return http.StatusBadRequest
	}

	_ = writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
		Message: err.Error(),
	})

	return http.StatusBadRequest
}
