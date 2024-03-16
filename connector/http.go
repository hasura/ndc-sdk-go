package connector

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/hasura/ndc-sdk-go/internal"
	"github.com/hasura/ndc-sdk-go/schema"
	"github.com/hasura/ndc-sdk-go/utils"
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
	logger          *slog.Logger
	telemetry       *TelemetryState
	recoveryEnabled bool
}

func newRouter(logger *slog.Logger, telemetry *TelemetryState, enableRecovery bool) *router {
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
			isDebug := rt.logger.Enabled(context.Background(), slog.LevelDebug)
			requestID := getRequestID(r)
			requestLogData := map[string]any{
				"url":            r.URL.String(),
				"method":         r.Method,
				"remote_address": r.RemoteAddr,
			}

			ctx := r.Context()
			//lint:ignore SA1012 possible to set nil
			span := trace.SpanFromContext(nil) //nolint:all
			spanName, spanOk := allowedTraceEndpoints[strings.ToLower(r.URL.Path)]
			if spanOk {
				ctx, span = rt.telemetry.Tracer.Start(
					otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
					spanName,
					trace.WithSpanKind(trace.SpanKindServer),
				)
			}
			defer span.End()

			if isDebug {
				requestLogData["headers"] = r.Header
				if spanOk {
					setSpanHeadersAttributes(span, r.Header, isDebug)
				}
				if r.Body != nil {
					bodyBytes, err := io.ReadAll(r.Body)
					if err != nil {
						rt.logger.Error("failed to read request",
							slog.String("request_id", requestID),
							slog.Duration("latency", time.Since(startTime)),
							slog.Any("request", requestLogData),
							slog.Any("error", err),
						)

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
						rt.logger.Error(
							"internal server error",
							slog.String("request_id", requestID),
							slog.Duration("latency", time.Since(startTime)),
							slog.Any("request", requestLogData),
							slog.Any("error", err),
							slog.String("stacktrace", stack),
						)

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
				rt.logger.Error(
					"handler not found",
					slog.String("request_id", requestID),
					slog.Duration("latency", time.Since(startTime)),
					slog.Any("request", requestLogData),
					slog.Any("response", map[string]any{
						"status": 404,
					}),
				)
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

					rt.logger.Error(
						"invalid content type",
						slog.String("request_id", requestID),
						slog.Duration("latency", time.Since(startTime)),
						slog.Any("request", requestLogData),
						slog.Any("response", map[string]any{
							"status": 422,
							"body":   err,
						}),
					)
					span.SetStatus(codes.Error, fmt.Sprintf("invalid content type: %s", contentType))
					return
				}
			}

			logger := rt.logger.With(slog.String("request_id", requestID))
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
				logger.Error(
					http.StatusText(writer.statusCode),
					slog.Duration("latency", time.Since(startTime)),
					slog.Any("request", requestLogData),
					slog.Any("response", responseLogData),
				)
				span.SetStatus(codes.Error, http.StatusText(writer.statusCode))
			} else {
				logger.Info(
					"success",
					slog.Duration("latency", time.Since(startTime)),
					slog.Any("request", requestLogData),
					slog.Any("response", responseLogData),
				)
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
func writeJsonFunc(w http.ResponseWriter, logger *slog.Logger, statusCode int, encodeFunc func() ([]byte, error)) {
	w.Header().Set(headerContentType, contentTypeJson)
	jsonBytes, err := encodeFunc()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, http.StatusText(http.StatusInternalServerError)))); err != nil {
			logger.Error("failed to write response", slog.Any("error", err))
		}
		return
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonBytes); err != nil {
		logger.Error("failed to write response", slog.Any("error", err))
	}
}

// writeJson writes response data with json encode
func writeJson(w http.ResponseWriter, logger *slog.Logger, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)
		return
	}

	writeJsonFunc(w, logger, statusCode, func() ([]byte, error) {
		return json.Marshal(body)
	})
}

// GetLogger gets the logger instance from context
func GetLogger(ctx context.Context) *slog.Logger {
	value := ctx.Value(logContextKey)
	if value != nil {
		if logger, ok := value.(*slog.Logger); ok {
			return logger
		}
	}
	return slog.Default()
}

func writeError(w http.ResponseWriter, logger *slog.Logger, err error) int {
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
