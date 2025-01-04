package connector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelPrometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/log/global"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	traceapi "go.opentelemetry.io/otel/trace"
)

const (
	otlpDefaultGRPCPort = 4317
	otlpDefaultHTTPPort = 4318
	otlpCompressionNone = "none"
	otlpCompressionGzip = "gzip"
)

var sensitiveHeaderRegex = regexp.MustCompile(`auth|key|secret|token`)

type otlpProtocol string

const (
	otlpProtocolGRPC         otlpProtocol = "grpc"
	otlpProtocolHTTPProtobuf otlpProtocol = "http/protobuf"
)

// defines the type of OpenTelemetry metrics exporter.
type otelMetricsExporterType string

const (
	otelMetricsExporterNone       otelMetricsExporterType = "none"
	otelMetricsExporterOTLP       otelMetricsExporterType = "otlp"
	otelMetricsExporterPrometheus otelMetricsExporterType = "prometheus"
)

var (
	userVisibilityAttribute = traceapi.WithAttributes(attribute.String("internal.visibility", "user"))
	successStatusAttribute  = attribute.String("status", "success")
	failureStatusAttribute  = attribute.String("status", "failure")
)

// OTLPConfig contains configuration for OpenTelemetry exporter.
type OTLPConfig struct {
	ServiceName            string `help:"OpenTelemetry service name." env:"OTEL_SERVICE_NAME"`
	OtlpEndpoint           string `help:"OpenTelemetry receiver endpoint that is set as default for all types." env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtlpTracesEndpoint     string `help:"OpenTelemetry endpoint for traces." env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	OtlpMetricsEndpoint    string `help:"OpenTelemetry endpoint for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
	OtlpLogsEndpoint       string `help:"OpenTelemetry endpoint for logs." env:"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT"`
	OtlpInsecure           *bool  `help:"Disable LTS for OpenTelemetry exporters." env:"OTEL_EXPORTER_OTLP_INSECURE"`
	OtlpTracesInsecure     *bool  `help:"Disable LTS for OpenTelemetry traces exporter." env:"OTEL_EXPORTER_OTLP_TRACES_INSECURE"`
	OtlpMetricsInsecure    *bool  `help:"Disable LTS for OpenTelemetry metrics exporter." env:"OTEL_EXPORTER_OTLP_METRICS_INSECURE"`
	OtlpLogsInsecure       *bool  `help:"Disable LTS for OpenTelemetry logs exporter." env:"OTEL_EXPORTER_OTLP_LOGS_INSECURE"`
	OtlpProtocol           string `help:"OpenTelemetry receiver protocol for all types." env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	OtlpTracesProtocol     string `help:"OpenTelemetry receiver protocol for traces." env:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"`
	OtlpMetricsProtocol    string `help:"OpenTelemetry receiver protocol for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"`
	OtlpLogsProtocol       string `help:"OpenTelemetry receiver protocol for logs." env:"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL"`
	OtlpCompression        string `help:"Enable compression for OTLP exporters. Accept: none, gzip" enum:"none,gzip" env:"OTEL_EXPORTER_OTLP_COMPRESSION" default:"gzip"`
	OtlpTraceCompression   string `help:"Enable compression for OTLP traces exporter. Accept: none, gzip" enum:"none,gzip" env:"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION" default:"gzip"`
	OtlpMetricsCompression string `help:"Enable compression for OTLP metrics exporter. Accept: none, gzip" enum:"none,gzip" env:"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION" default:"gzip"`
	OtlpLogsCompression    string `help:"Enable compression for OTLP logs exporter. Accept: none, gzip" enum:"none,gzip" env:"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION" default:"gzip"`

	MetricsExporter  string `help:"Metrics export type. Accept: none, otlp, prometheus" enum:"none,otlp,prometheus" env:"OTEL_METRICS_EXPORTER" default:"none"`
	LogsExporter     string `help:"Logs export type. Accept: none, otlp" enum:"none,otlp" env:"OTEL_LOGS_EXPORTER" default:"none"`
	PrometheusPort   *uint  `help:"Prometheus port for the Prometheus HTTP server. Use /metrics endpoint of the connector server if empty" env:"OTEL_EXPORTER_PROMETHEUS_PORT"`
	DisableGoMetrics *bool  `help:"Disable internal Go and process metrics"`
}

// TelemetryState contains OpenTelemetry exporters and basic connector metrics.
type TelemetryState struct {
	*connectorMetrics

	Tracer   *Tracer
	Meter    metricapi.Meter
	Logger   *slog.Logger
	Shutdown func(context.Context) error
}

type connectorMetrics struct {
	queryCounter                    metricapi.Int64Counter
	mutationCounter                 metricapi.Int64Counter
	queryExplainCounter             metricapi.Int64Counter
	mutationExplainCounter          metricapi.Int64Counter
	queryLatencyHistogram           metricapi.Float64Histogram
	queryExplainLatencyHistogram    metricapi.Float64Histogram
	mutationExplainLatencyHistogram metricapi.Float64Histogram
	mutationLatencyHistogram        metricapi.Float64Histogram
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, config *OTLPConfig, serviceVersion, metricsPrefix string, logger *slog.Logger) (*TelemetryState, error) {
	state, err := SetupOTelExporters(ctx, config, serviceVersion, metricsPrefix, logger)
	if err != nil {
		return nil, err
	}

	return state, setupConnectorMetrics(state, metricsPrefix)
}

// SetupOTelExporters set up OpenTelemetry exporters from configuration.
func SetupOTelExporters(ctx context.Context, config *OTLPConfig, serviceVersion, metricsPrefix string, logger *slog.Logger) (*TelemetryState, error) {
	otel.SetLogger(logr.FromSlogHandler(logger.Handler()))
	otelDisabled := os.Getenv("OTEL_SDK_DISABLED") == "true"

	tracesEndpoint := utils.GetDefault(config.OtlpTracesEndpoint, config.OtlpEndpoint)
	metricsEndpoint := utils.GetDefault(config.OtlpMetricsEndpoint, config.OtlpEndpoint)

	// Set up resource.
	res := newResource(logger, config.ServiceName, serviceVersion)

	var traceProvider *trace.TracerProvider
	if !otelDisabled && tracesEndpoint != "" {
		endpoint, protocol, insecure, err := parseOTLPEndpoint(
			tracesEndpoint,
			utils.GetDefault(config.OtlpTracesProtocol, config.OtlpProtocol),
			utils.GetDefaultPtr(config.OtlpTracesInsecure, config.OtlpInsecure),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP traces endpoint: %w", err)
		}

		compressorStr, compressorInt, err := parseOTLPCompression(utils.GetDefault(config.OtlpTraceCompression, config.OtlpCompression))
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP traces compression: %w", err)
		}

		// Set up propagator.
		prop := newPropagator()
		otel.SetTextMapPropagator(prop)

		var traceExporter *otlptrace.Exporter

		if protocol == otlpProtocolGRPC {
			options := []otlptracegrpc.Option{
				otlptracegrpc.WithEndpoint(endpoint),
				otlptracegrpc.WithCompressor(compressorStr),
			}

			if insecure {
				options = append(options, otlptracegrpc.WithInsecure())
			}

			traceExporter, err = otlptracegrpc.New(ctx, options...)
			if err != nil {
				return nil, err
			}
		} else {
			options := []otlptracehttp.Option{
				otlptracehttp.WithEndpoint(endpoint),
				otlptracehttp.WithCompression(otlptracehttp.Compression(compressorInt)),
			}
			if insecure {
				options = append(options, otlptracehttp.WithInsecure())
			}

			traceExporter, err = otlptracehttp.New(ctx, options...)
			if err != nil {
				return nil, err
			}
		}

		traceProvider = trace.NewTracerProvider(
			trace.WithResource(res),
			trace.WithBatcher(traceExporter),
		)
	} else {
		traceProvider = trace.NewTracerProvider(trace.WithResource(res))
	}
	otel.SetTracerProvider(traceProvider)

	// configure metrics exporter
	metricsExporterType, err := parseOTELMetricsExporterType(config.MetricsExporter)
	if err != nil {
		return nil, err
	}

	metricOptions := []metric.Option{metric.WithResource(res)}

	if config.DisableGoMetrics != nil && !*config.DisableGoMetrics {
		// disable default process and go collector metrics
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		prometheus.Unregister(collectors.NewGoCollector())
	}

	switch metricsExporterType {
	case otelMetricsExporterPrometheus:
		// The exporter embeds a default OpenTelemetry Reader and
		// implements prometheus.Collector, allowing it to be used as
		// both a Reader and Collector.
		prometheusExporter, err := otelPrometheus.New()
		if err != nil {
			return nil, err
		}
		metricOptions = append(metricOptions, metric.WithReader(prometheusExporter))
	case otelMetricsExporterOTLP:
		if otelDisabled {
			break
		}

		if metricsEndpoint == "" {
			return nil, errors.New("OTLP endpoint is required for metrics exporter")
		}
		endpoint, protocol, insecure, err := parseOTLPEndpoint(
			metricsEndpoint,
			utils.GetDefault(config.OtlpMetricsProtocol, config.OtlpProtocol),
			utils.GetDefaultPtr(config.OtlpMetricsInsecure, config.OtlpInsecure),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP metrics endpoint: %w", err)
		}

		compressorStr, compressorInt, err := parseOTLPCompression(utils.GetDefault(config.OtlpMetricsCompression, config.OtlpCompression))
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP metrics compression: %w", err)
		}

		if protocol == otlpProtocolGRPC {
			options := []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithEndpoint(endpoint),
				otlpmetricgrpc.WithCompressor(compressorStr),
			}

			if insecure {
				options = append(options, otlpmetricgrpc.WithInsecure())
			}

			metricExporter, err := otlpmetricgrpc.New(ctx, options...)
			if err != nil {
				return nil, err
			}
			metricOptions = append(metricOptions, metric.WithReader(metric.NewPeriodicReader(metricExporter)))
		} else {
			options := []otlpmetrichttp.Option{
				otlpmetrichttp.WithEndpoint(endpoint),
				otlpmetrichttp.WithCompression(otlpmetrichttp.Compression(compressorInt)),
			}
			if insecure {
				options = append(options, otlpmetrichttp.WithInsecure())
			}

			metricExporter, err := otlpmetrichttp.New(ctx, options...)
			if err != nil {
				return nil, err
			}
			metricOptions = append(metricOptions, metric.WithReader(metric.NewPeriodicReader(metricExporter)))
		}
	}

	meterProvider := metric.NewMeterProvider(metricOptions...)
	otel.SetMeterProvider(meterProvider)

	// configure metrics exporter
	loggerProvider, err := newLoggerProvider(ctx, config, otelDisabled, res)
	if err != nil {
		return nil, err
	}
	global.SetLoggerProvider(loggerProvider)

	shutdownFunc := func(ctx context.Context) error {
		errorMsgs := []string{}
		if err := traceProvider.Shutdown(ctx); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}
		if err := loggerProvider.Shutdown(ctx); err != nil {
			errorMsgs = append(errorMsgs, err.Error())
		}

		if len(errorMsgs) > 0 {
			return errors.New(strings.Join(errorMsgs, ","))
		}

		return nil
	}

	otelLogger := slog.New(createLogHandler(config.ServiceName, logger, loggerProvider))
	state := &TelemetryState{
		Tracer:   &Tracer{traceProvider.Tracer(config.ServiceName, traceapi.WithSchemaURL(semconv.SchemaURL))},
		Meter:    meterProvider.Meter(config.ServiceName, metricapi.WithSchemaURL(semconv.SchemaURL)),
		Logger:   otelLogger,
		Shutdown: shutdownFunc,
	}

	return state, err
}

func newResource(logger *slog.Logger, serviceName, serviceVersion string) *resource.Resource {
	attrs := resource.NewWithAttributes(semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
	)
	res, err := resource.Merge(resource.Default(), attrs)
	if err != nil {
		logger.Warn(err.Error())
		return attrs
	}
	return res
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)),
	)
}

func setupConnectorMetrics(telemetry *TelemetryState, metricsPrefix string) error {
	if metricsPrefix != "" {
		metricsPrefix += "."
	}

	var err error
	meter := telemetry.Meter
	telemetry.connectorMetrics = &connectorMetrics{}

	telemetry.queryCounter, err = meter.Int64Counter(
		metricsPrefix+"query.total",
		metricapi.WithDescription("Total number of query requests"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationCounter, err = meter.Int64Counter(
		metricsPrefix+"mutation.total",
		metricapi.WithDescription("Total number of mutation requests"),
	)
	if err != nil {
		return err
	}

	telemetry.queryExplainCounter, err = meter.Int64Counter(
		metricsPrefix+"query.explain_total",
		metricapi.WithDescription("Total number of explain query requests"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationExplainCounter, err = meter.Int64Counter(
		metricsPrefix+"mutation.explain_total",
		metricapi.WithDescription("Total number of explain mutation requests"),
	)
	if err != nil {
		return err
	}

	telemetry.queryLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"query.total_time",
		metricapi.WithDescription("Total time taken to plan and execute a query, in seconds"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"mutation.total_time",
		metricapi.WithDescription("Total time taken to plan and execute a mutation, in seconds"),
	)
	if err != nil {
		return err
	}

	telemetry.queryExplainLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"query.explain_total_time",
		metricapi.WithDescription("Total time taken to plan and execute an explain query request, in seconds"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationExplainLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"mutation.explain_total_time",
		metricapi.WithDescription("Total time taken to plan and execute an explain mutation request, in seconds"),
	)

	return err
}

// Tracer is the wrapper of traceapi.Tracer with user visibility on Hasura Console.
type Tracer struct {
	traceapi.Tracer
}

var _ traceapi.Tracer = &Tracer{}

// NewTracer creates a new OpenTelemetry tracer.
func NewTracer(name string, opts ...traceapi.TracerOption) *Tracer {
	return &Tracer{
		Tracer: otel.Tracer(name, opts...),
	}
}

// Start creates a span and a context.Context containing the newly-created span
// with `internal.visibility: "user"` so that it shows up in the Hasura Console.
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...traceapi.SpanStartOption) (context.Context, traceapi.Span) {
	return t.Tracer.Start(ctx, spanName, append(opts, userVisibilityAttribute)...) //nolint:spancheck
}

// StartInternal creates a span and a context.Context containing the newly-created span.
// It won't show up in the Hasura Console.
func (t *Tracer) StartInternal(ctx context.Context, spanName string, opts ...traceapi.SpanStartOption) (context.Context, traceapi.Span) {
	return t.Tracer.Start(ctx, spanName, opts...) //nolint:spancheck
}

func httpStatusAttribute(code int) attribute.KeyValue {
	return attribute.Int("http_status", code)
}

func parseOTLPEndpoint(endpoint string, protocol string, insecurePtr *bool) (string, otlpProtocol, bool, error) {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	uri, err := url.Parse(endpoint)
	if err != nil {
		return "", otlpProtocol(""), false, err
	}

	insecure := utils.GetDefaultValuePtr(insecurePtr, uri.Scheme == "http")
	host := uri.Host
	if uri.Port() == "" {
		port := 443
		if insecure {
			port = 80
		}
		host = fmt.Sprintf("%s:%d", uri.Hostname(), port)
	}

	switch protocol {
	case string(otlpProtocolGRPC):
		return host, otlpProtocolGRPC, insecure, nil
	case string(otlpProtocolHTTPProtobuf):
		return host, otlpProtocol(protocol), insecure, nil
	case "":
		// auto detect via default OTLP port
		if uri.Port() == strconv.FormatInt(otlpDefaultHTTPPort, 10) {
			return host, otlpProtocol(protocol), insecure, nil
		}
		return host, otlpProtocolGRPC, insecure, nil
	default:
		return "", otlpProtocol(""), false, errors.New("invalid OTLP protocol " + protocol)
	}
}

func parseOTLPCompression(input string) (string, int, error) {
	switch input {
	case otlpCompressionGzip, "":
		return otlpCompressionGzip, int(otlptracehttp.GzipCompression), nil
	case otlpCompressionNone:
		return otlpCompressionNone, int(otlptracehttp.NoCompression), nil
	default:
		return "", 0, errors.New("invalid OTLP compression type, accept none, gzip only")
	}
}

func parseOTELMetricsExporterType(input string) (otelMetricsExporterType, error) {
	switch input {
	case string(otelMetricsExporterNone), "":
		return otelMetricsExporterNone, nil
	case string(otelMetricsExporterOTLP):
		return otelMetricsExporterOTLP, nil
	case string(otelMetricsExporterPrometheus):
		return otelMetricsExporterPrometheus, nil
	default:
		return otelMetricsExporterNone, errors.New("invalid metrics exporter type: " + input)
	}
}

// SetSpanHeaderAttributes sets header attributes to the otel span.
func SetSpanHeaderAttributes(span traceapi.Span, prefix string, httpHeaders http.Header, allowedHeaders ...string) {
	headers := NewTelemetryHeaders(httpHeaders, allowedHeaders...)

	for key, values := range headers {
		span.SetAttributes(attribute.StringSlice(prefix+strings.ToLower(key), values))
	}
}

// NewTelemetryHeaders creates a new header map with sensitive values masked.
func NewTelemetryHeaders(httpHeaders http.Header, allowedHeaders ...string) http.Header {
	result := http.Header{}

	if len(allowedHeaders) > 0 {
		for _, key := range allowedHeaders {
			value := httpHeaders.Get(key)

			if value == "" {
				continue
			}

			if IsSensitiveHeader(key) {
				result.Set(strings.ToLower(key), MaskString(value))
			} else {
				result.Set(strings.ToLower(key), value)
			}
		}

		return result
	}

	for key, headers := range httpHeaders {
		if len(headers) == 0 {
			continue
		}

		values := headers
		if IsSensitiveHeader(key) {
			values = make([]string, len(headers))
			for i, header := range headers {
				values[i] = MaskString(header)
			}
		}

		result[key] = values
	}

	return result
}

// IsSensitiveHeader checks if the header name is sensitive.
func IsSensitiveHeader(name string) bool {
	return sensitiveHeaderRegex.MatchString(strings.ToLower(name))
}

// MaskString masks the string value for security.
func MaskString(input string) string {
	inputLength := len(input)

	switch {
	case inputLength <= 6:
		return strings.Repeat("*", inputLength)
	case inputLength < 12:
		return input[0:1] + strings.Repeat("*", inputLength-1)
	default:
		return input[0:3] + strings.Repeat("*", 7) + fmt.Sprintf("(%d)", inputLength)
	}
}
