package connector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/ndc-sdk-go/v2/schema"
	"github.com/hasura/ndc-sdk-go/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	headerContentType string = "Content-Type"
	contentTypeJson   string = "application/json"
)

const (
	apiPathCapabilities    = "/capabilities"
	apiPathSchema          = "/schema"
	apiPathQuery           = "/query"
	apiPathQueryExplain    = "/query/explain"
	apiPathMutation        = "/mutation"
	apiPathMutationExplain = "/mutation/explain"
	apiPathHealth          = "/health"
	apiPathMetrics         = "/metrics"
)

var allowedTraceEndpoints = map[string]string{
	apiPathQuery:           "ndc_query",
	apiPathQueryExplain:    "ndc_query_explain",
	apiPathMutation:        "ndc_mutation",
	apiPathMutationExplain: "ndc_mutation_explain",
}

var debugApiPaths = []string{apiPathMetrics, apiPathHealth}

// implements a simple router to reuse for both configuration and connector servers.
type router struct {
	routes          map[string]map[string][]http.HandlerFunc
	logger          *slog.Logger
	telemetry       *TelemetryState
	recoveryEnabled bool
}

func newRouter(logger *slog.Logger, telemetry *TelemetryState, enableRecovery bool) *router {
	return &router{
		routes:          make(map[string]map[string][]http.HandlerFunc),
		logger:          logger,
		telemetry:       telemetry,
		recoveryEnabled: enableRecovery,
	}
}

func (rt *router) Use(path string, method string, handlers ...http.HandlerFunc) {
	if _, ok := rt.routes[path]; !ok {
		rt.routes[path] = make(map[string][]http.HandlerFunc)
	}

	rt.routes[path][method] = handlers
}

func (rt *router) Build() *http.ServeMux {
	mux := http.NewServeMux()

	for path, handlers := range rt.routes {
		handler := rt.createHandleFunc(handlers)
		mux.HandleFunc(path, handler)
	}

	return mux
}

func (rt *router) createHandleFunc( //nolint:gocognit
	handlers map[string][]http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		isDebug := rt.logger.Enabled(r.Context(), slog.LevelDebug)
		requestLogData := map[string]any{
			"url":            r.URL.String(),
			"method":         r.Method,
			"remote_address": r.RemoteAddr,
		}

		ctx := r.Context()
		//lint:ignore SA1012 possible to set nil
		span := trace.SpanFromContext(nil) //nolint:contextcheck,staticcheck
		spanName, spanOk := allowedTraceEndpoints[strings.ToLower(r.URL.Path)]

		if spanOk {
			ctx, span = rt.telemetry.Tracer.Start(
				otel.GetTextMapPropagator().
					Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
				spanName,
				trace.WithSpanKind(trace.SpanKindServer),
			)
		}

		defer span.End()

		requestID := getRequestID(r, span)

		// Add HTTP semantic attributes to the server span
		// See: https://opentelemetry.io/docs/specs/semconv/http/http-spans/#http-server-semantic-conventions
		span.SetAttributes(
			attribute.String("http.request.method", strings.ToUpper(r.Method)),
			attribute.String("url.path", r.URL.Path),
			attribute.String("url.scheme", r.URL.Scheme),
			attribute.Int64("http.request.body.size", r.ContentLength),
		)

		if isDebug {
			requestLogData["headers"] = r.Header

			if spanOk {
				SetSpanHeaderAttributes(span, "http.request.header", r.Header)
			}

			if r.Body != nil {
				bodyStr, err := rt.debugRequestBody(w, r, span)
				if err != nil {
					rt.logger.Error("failed to read request",
						slog.String("request_id", requestID),
						slog.Duration("latency", time.Since(startTime)),
						slog.Any("request", requestLogData),
						slog.Any("error", err),
					)

					return
				}

				requestLogData["body"] = bodyStr
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

					span.SetAttributes(
						attribute.Int("http.response.status_code", http.StatusInternalServerError),
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

		handlers, ok := handlers[r.Method]
		if !ok {
			http.NotFound(w, r)
			rt.logger.Error(
				"handler not found",
				slog.String("request_id", requestID),
				slog.Duration("latency", time.Since(startTime)),
				slog.Any("request", requestLogData),
				slog.Any("response", map[string]any{
					"status": http.StatusNotFound,
				}),
			)

			span.SetAttributes(attribute.Int("http.response.status_code", http.StatusNotFound))
			span.SetStatus(codes.Error, fmt.Sprintf("path %s is not found", r.URL.RequestURI()))

			return
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut ||
			r.Method == http.MethodPatch {
			contentType := r.Header.Get(headerContentType)
			if contentType != contentTypeJson {
				err := schema.ErrorResponse{
					Message: fmt.Sprintf(
						"Invalid content type %s, accept %s only",
						contentType,
						contentTypeJson,
					),
				}

				writeJson(w, rt.logger, http.StatusUnprocessableEntity, err)

				rt.logger.Error(
					"invalid content type",
					slog.String("request_id", requestID),
					slog.Duration("latency", time.Since(startTime)),
					slog.Any("request", requestLogData),
					slog.Any("response", map[string]any{
						"status": http.StatusUnprocessableEntity,
						"body":   err,
					}),
				)

				span.SetAttributes(
					attribute.Int("http.response.status_code", http.StatusUnprocessableEntity),
				)
				span.SetStatus(codes.Error, "invalid content type: "+contentType)

				return
			}
		}

		logger := rt.logger.With(slog.String("request_id", requestID))
		req := r.WithContext(NewContextLogger(ctx, logger))
		writer := newCustomResponseWriter(r, w)

		for _, h := range handlers {
			h(writer, req)

			// stop the loop if the status code was written
			if writer.statusCode > 0 {
				break
			}
		}

		span.SetAttributes(attribute.Int("http.response.status_code", writer.statusCode))
		span.SetAttributes(attribute.Int("http.response.body.size", writer.bodyLength))
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

		SetSpanHeaderAttributes(span, "http.response.header", w.Header())

		if writer.statusCode >= http.StatusBadRequest {
			logger.Error(
				http.StatusText(writer.statusCode),
				slog.Duration("latency", time.Since(startTime)),
				slog.Any("request", requestLogData),
				slog.Any("response", responseLogData),
			)

			span.SetStatus(codes.Error, http.StatusText(writer.statusCode))

			return
		}

		printSuccess := logger.Info

		if slices.Contains(debugApiPaths, r.URL.Path) {
			printSuccess = logger.Debug
		}

		printSuccess(
			"success",
			slog.Duration("latency", time.Since(startTime)),
			slog.Any("request", requestLogData),
			slog.Any("response", responseLogData),
		)
		span.SetStatus(codes.Ok, "success")
	}
}

func (rt *router) debugRequestBody(
	w http.ResponseWriter,
	r *http.Request,
	span trace.Span,
) (string, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		span.SetAttributes(
			attribute.Int("http.response.status_code", http.StatusUnprocessableEntity),
		)

		writeJson(w, rt.logger, http.StatusUnprocessableEntity, schema.ErrorResponse{
			Message: "failed to read request",
			Details: map[string]any{
				"cause": err,
			},
		})

		span.SetStatus(codes.Error, "read_request_body_failure")
		span.RecordError(err)

		return "", err
	}

	bodyStr := string(bodyBytes)
	span.SetAttributes(attribute.String("request.body", bodyStr))

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyStr, nil
}

func getRequestID(r *http.Request, span trace.Span) string {
	requestID := r.Header.Get("x-request-id")
	if requestID != "" {
		return requestID
	}

	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}

	return uuid.NewString()
}

// writeJsonFunc writes raw bytes data with a json encoding callback.
func writeJsonFunc(
	w http.ResponseWriter,
	logger *slog.Logger,
	statusCode int,
	encodeFunc func() ([]byte, error),
) {
	w.Header().Set(headerContentType, contentTypeJson)

	jsonBytes, err := encodeFunc()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		if _, err := fmt.Fprintf(w, `{"message": "%s"}`, http.StatusText(http.StatusInternalServerError)); err != nil {
			logger.Error("failed to write response", slog.Any("error", err))
		}

		return
	}

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonBytes); err != nil {
		logger.Error("failed to write response", slog.Any("error", err))
	}
}

// writeJson writes response data with json encode.
func writeJson(w http.ResponseWriter, logger *slog.Logger, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)

		return
	}

	writeJsonFunc(w, logger, statusCode, func() ([]byte, error) {
		return json.Marshal(body)
	})
}

func writeError(w http.ResponseWriter, logger *slog.Logger, err error, defaultHttpStatus int) int {
	w.Header().Add("Content-Type", "application/json")

	var connectorErrorPtr *schema.ConnectorError
	if errors.As(err, &connectorErrorPtr) {
		writeJson(w, logger, connectorErrorPtr.StatusCode(), connectorErrorPtr)

		return connectorErrorPtr.StatusCode()
	}

	var errorResponse schema.ErrorResponse
	if errors.As(err, &errorResponse) {
		writeJson(w, logger, defaultHttpStatus, errorResponse)

		return http.StatusInternalServerError
	}

	var errorResponsePtr *schema.ErrorResponse
	if errors.As(err, &errorResponsePtr) {
		writeJson(w, logger, defaultHttpStatus, errorResponsePtr)

		return http.StatusInternalServerError
	}

	writeJson(w, logger, defaultHttpStatus, schema.ErrorResponse{
		Message: err.Error(),
	})

	return http.StatusInternalServerError
}
