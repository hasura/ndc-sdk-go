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
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type serverContextKey string

const (
	logContextKey     serverContextKey = "hasura-log"
	headerContentType string           = "Content-Type"
	contentTypeJson   string           = "application/json"
)

var allowedTraceEndpoints = map[string]string{
	"/query":            "ndc_query",
	"/query/explain":    "ndc_query_explain",
	"/mutation":         "ndc_mutation",
	"/mutation/explain": "ndc_mutation_explain",
}

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
	telemetry       *TelemetryState
	recoveryEnabled bool
}

func newRouter(logger zerolog.Logger, telemetry *TelemetryState, enableRecovery bool) *router {
	return &router{
		routes:          make(map[string]map[string]http.HandlerFunc),
		logger:          logger,
		telemetry:       telemetry,
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

			ctx := r.Context()
			span := trace.SpanFromContext(nil)
			spanName, spanOk := allowedTraceEndpoints[strings.ToLower(r.URL.Path)]
			if spanOk {
				ctx, span = rt.telemetry.Tracer.Start(
					otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
					spanName,
					trace.WithSpanKind(trace.SpanKindServer),
				)
			}
			defer span.End()

			log.Printf("is debug: %t", isDebug)
			if isDebug {
				requestLogData["headers"] = r.Header
				if spanOk {
					setSpanHeadersAttributes(span, r.Header, isDebug)
				}
				if r.Body != nil {
					bodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						rt.logger.Error().
							Str("request_id", requestID).
							Dur("latency", time.Since(startTime)).
							Interface("request", requestLogData).
							Interface("error", err).
							Msg("failed to read request")

						writeJson(w, rt.logger, http.StatusUnprocessableEntity, schema.ErrorResponse{
							Message: "failed to read request",
							Details: map[string]any{
								"cause": err,
							},
						})

						span.SetStatus(codes.Error, "read_request_body_failure")
						span.RecordError(err)
						return
					}

					bodyStr := string(bodyBytes)
					span.SetAttributes(attribute.String("request.body", bodyStr))
					requestLogData["body"] = bodyStr
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			// recover from panic
			if rt.recoveryEnabled {
				defer func() {
					if err := recover(); err != nil {
						stack := string(debug.Stack())
						rt.logger.Error().
							Str("request_id", requestID).
							Dur("latency", time.Since(startTime)).
							Interface("request", requestLogData).
							Interface("error", err).
							Str("stacktrace", stack).
							Msg("internal server error")

						writeJson(w, rt.logger, http.StatusInternalServerError, schema.ErrorResponse{
							Message: "internal server error",
							Details: map[string]any{
								"cause": err,
							},
						})

						span.SetAttributes(utils.JSONAttribute("error", err))
						span.SetAttributes(attribute.String("stacktrace", stack))
						span.SetStatus(codes.Error, "panic")
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
				span.SetStatus(codes.Error, fmt.Sprintf("path %s is not found", r.URL.RequestURI()))
				return
			}

			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				contentType := r.Header.Get(headerContentType)
				if contentType != contentTypeJson {
					err := schema.ErrorResponse{
						Message: fmt.Sprintf("Invalid content type %s, accept %s only", contentType, contentTypeJson),
					}
					writeJson(w, rt.logger, http.StatusUnprocessableEntity, err)

					rt.logger.Error().
						Str("request_id", requestID).
						Dur("latency", time.Since(startTime)).
						Interface("request", requestLogData).
						Interface("response", map[string]any{
							"status": 422,
							"body":   err,
						}).
						Msg("")
					span.SetStatus(codes.Error, fmt.Sprintf("invalid content type: %s", contentType))
					return
				}
			}

			logger := rt.logger.With().Str("request_id", requestID).Logger()
			req := r.WithContext(context.WithValue(ctx, logContextKey, logger))
			writer := &customResponseWriter{ResponseWriter: w}
			h(writer, req)

			responseLogData := map[string]any{
				"status": writer.statusCode,
			}
			if isDebug || writer.statusCode >= 400 {
				responseLogData["headers"] = writer.Header()
				if len(writer.body) > 0 {
					responseLogData["body"] = string(writer.body)
					span.SetAttributes(attribute.String("response.body", string(writer.body)))
				}
			}
			setSpanHeadersAttributes(span, w.Header(), isDebug)

			if writer.statusCode >= 400 {
				logger.Error().
					Dur("latency", time.Since(startTime)).
					Interface("request", requestLogData).
					Interface("response", responseLogData).Msg("")
				span.SetStatus(codes.Error, http.StatusText(writer.statusCode))
			} else {
				logger.Info().
					Dur("latency", time.Since(startTime)).
					Interface("request", requestLogData).
					Interface("response", responseLogData).Msg("")
				span.SetStatus(codes.Ok, "success")
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
func writeJson(w http.ResponseWriter, logger zerolog.Logger, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)
		return
	}

	writeJsonFunc(w, logger, statusCode, func() ([]byte, error) {
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
		writeJson(w, logger, connectorErrorPtr.StatusCode(), connectorErrorPtr)
		return connectorErrorPtr.StatusCode()
	}

	var errorResponse schema.ErrorResponse
	if errors.As(err, &errorResponse) {
		writeJson(w, logger, http.StatusBadRequest, errorResponse)
		return http.StatusBadRequest
	}

	var errorResponsePtr *schema.ErrorResponse
	if errors.As(err, &errorResponsePtr) {
		writeJson(w, logger, http.StatusBadRequest, errorResponsePtr)
		return http.StatusBadRequest
	}

	writeJson(w, logger, http.StatusBadRequest, schema.ErrorResponse{
		Message: err.Error(),
	})

	return http.StatusBadRequest
}
