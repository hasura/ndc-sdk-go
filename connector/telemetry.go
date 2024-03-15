package connector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/hasura/ndc-sdk-go/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelPrometheus "go.opentelemetry.io/otel/exporters/prometheus"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	traceapi "go.opentelemetry.io/otel/trace"
)

const (
	otlpDefaultGRPCPort = 4317
	otlpDefaultHTTPPort = 4318
	otlpCompressionNone = "none"
	otlpCompressionGzip = "gzip"
)

type otlpProtocol string

const (
	otlpProtocolGRPC         otlpProtocol = "grpc"
	otlpProtocolHTTPProtobuf otlpProtocol = "http/protobuf"
)

// defines the type of OpenTelemetry metrics exporter
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

// OTLPConfig contains configuration for OpenTelemetry exporter
type OTLPConfig struct {
	ServiceName            string `help:"OpenTelemetry service name." env:"OTEL_SERVICE_NAME"`
	OtlpEndpoint           string `help:"OpenTelemetry receiver endpoint that is set as default for all types." env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtlpTracesEndpoint     string `help:"OpenTelemetry endpoint for traces." env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	OtlpMetricsEndpoint    string `help:"OpenTelemetry endpoint for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
	OtlpInsecure           *bool  `help:"Disable LTS for OpenTelemetry exporters." env:"OTEL_EXPORTER_OTLP_INSECURE"`
	OtlpTracesInsecure     *bool  `help:"Disable LTS for OpenTelemetry traces exporter." env:"OTEL_EXPORTER_OTLP_TRACES_INSECURE"`
	OtlpMetricsInsecure    *bool  `help:"Disable LTS for OpenTelemetry metrics exporter." env:"OTEL_EXPORTER_OTLP_METRICS_INSECURE"`
	OtlpProtocol           string `help:"OpenTelemetry receiver protocol for all types." env:"OTEL_EXPORTER_OTLP_PROTOCOL"`
	OtlpTracesProtocol     string `help:"OpenTelemetry receiver protocol for traces." env:"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"`
	OtlpMetricsProtocol    string `help:"OpenTelemetry receiver protocol for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"`
	OtlpCompression        string `help:"Enable compression for OTLP exporters. Accept: none, gzip" env:"OTEL_EXPORTER_OTLP_COMPRESSION" default:"gzip"`
	OtlpTraceCompression   string `help:"Enable compression for OTLP traces exporter. Accept: none, gzip" env:"OTEL_EXPORTER_OTLP_TRACES_COMPRESSION"`
	OtlpMetricsCompression string `help:"Enable compression for OTLP metrics exporter. Accept: none, gzip." env:"OTEL_EXPORTER_OTLP_METRICS_COMPRESSION"`

	MetricsExporter string `help:"Metrics export type. Accept: none, otlp, prometheus." env:"OTEL_METRICS_EXPORTER" default:"none"`
	PrometheusPort  *uint  `help:"Prometheus port for the Prometheus HTTP server. Use /metrics endpoint of the connector server if empty" env:"OTEL_EXPORTER_PROMETHEUS_PORT"`
}

type TelemetryState struct {
	Tracer                          *Tracer
	Meter                           metricapi.Meter
	Shutdown                        func(context.Context) error
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

	otel.SetLogger(logr.FromSlogHandler(logger.Handler()))
	tracesEndpoint := utils.GetDefault(config.OtlpTracesEndpoint, config.OtlpEndpoint)
	metricsEndpoint := utils.GetDefault(config.OtlpMetricsEndpoint, config.OtlpEndpoint)

	// Set up resource.
	res, err := newResource(config.ServiceName, serviceVersion)
	if err != nil {
		return nil, err
	}

	var traceProvider *trace.TracerProvider
	if tracesEndpoint != "" {
		endpoint, protocol, insecure, err := parseOTLPEndpoint(
			tracesEndpoint,
			utils.GetDefault(config.OtlpTracesProtocol, config.OtlpProtocol),
			utils.GetDefaultPtr(config.OtlpTracesInsecure, config.OtlpInsecure),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP traces endpoint: %s", err)
		}

		compressorStr, compressorInt, err := parseOTLPCompression(utils.GetDefault(config.OtlpTraceCompression, config.OtlpCompression))
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP traces compression: %s", err)
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
			trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
		)
	} else {
		traceProvider = trace.NewTracerProvider()
	}
	otel.SetTracerProvider(traceProvider)

	// configure metrics exporter
	metricsExporterType, err := parseOTELMetricsExporterType(config.MetricsExporter)
	if err != nil {
		return nil, err
	}

	metricOptions := []metric.Option{metric.WithResource(res)}

	// disable default process and go collector metrics
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewGoCollector())

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
		if metricsEndpoint == "" {
			return nil, errors.New("OTLP endpoint is required for metrics exporter")
		}
		endpoint, protocol, insecure, err := parseOTLPEndpoint(
			metricsEndpoint,
			utils.GetDefault(config.OtlpMetricsProtocol, config.OtlpProtocol),
			utils.GetDefaultPtr(config.OtlpMetricsInsecure, config.OtlpInsecure),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP metrics endpoint: %s", err)
		}

		compressorStr, compressorInt, err := parseOTLPCompression(utils.GetDefault(config.OtlpMetricsCompression, config.OtlpCompression))
		if err != nil {
			return nil, fmt.Errorf("failed to parse OTLP metrics compression: %s", err)
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

	shutdownFunc := func(ctx context.Context) error {
		tErr := traceProvider.Shutdown(ctx)
		mErr := meterProvider.Shutdown(ctx)
		if tErr != nil || mErr != nil {
			return errors.New(strings.Join([]string{tErr.Error(), mErr.Error()}, ","))
		}
		return nil
	}

	state := &TelemetryState{
		Tracer:   &Tracer{traceProvider.Tracer(config.ServiceName)},
		Meter:    meterProvider.Meter(config.ServiceName),
		Shutdown: shutdownFunc,
	}

	err = setupMetrics(state, metricsPrefix)

	return state, err
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func setupMetrics(telemetry *TelemetryState, metricsPrefix string) error {
	if metricsPrefix != "" {
		metricsPrefix = metricsPrefix + "."
	}

	var err error
	meter := telemetry.Meter
	telemetry.queryCounter, err = meter.Int64Counter(
		fmt.Sprintf("%squery.total", metricsPrefix),
		metricapi.WithDescription("The total number of query requests"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationCounter, err = meter.Int64Counter(
		fmt.Sprintf("%smutation.total", metricsPrefix),
		metricapi.WithDescription("The total number of mutation requests"),
	)

	if err != nil {
		return err
	}

	telemetry.queryExplainCounter, err = meter.Int64Counter(
		fmt.Sprintf("%squery.explain_total", metricsPrefix),
		metricapi.WithDescription("The total number of explain query requests"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationExplainCounter, err = meter.Int64Counter(
		fmt.Sprintf("%smutation.explain_total", metricsPrefix),
		metricapi.WithDescription("The total number of explain mutation requests"),
	)
	if err != nil {
		return err
	}

	telemetry.queryLatencyHistogram, err = meter.Float64Histogram(
		fmt.Sprintf("%squery.total_time", metricsPrefix),
		metricapi.WithDescription("Total time taken to plan and execute a query, in seconds"),
	)
	if err != nil {
		return err
	}

	telemetry.mutationLatencyHistogram, err = meter.Float64Histogram(
		fmt.Sprintf("%smutation.total_time", metricsPrefix),
		metricapi.WithDescription("Total time taken to plan and execute a mutation, in seconds"),
	)
	if err != nil {
		return err
	}

	telemetry.queryExplainLatencyHistogram, err = meter.Float64Histogram(
		fmt.Sprintf("%squery.explain_total_time", metricsPrefix),
		metricapi.WithDescription("Total time taken to plan and execute an explain query request, in seconds"),
	)

	if err != nil {
		return err
	}

	telemetry.mutationExplainLatencyHistogram, err = meter.Float64Histogram(
		fmt.Sprintf("%smutation.explain_total_time", metricsPrefix),
		metricapi.WithDescription("Total time taken to plan and execute an explain mutation request, in seconds"),
	)

	return err
}

// Tracer is the wrapper of traceapi.Tracer with user visibility on Hasura Console
type Tracer struct {
	tracer traceapi.Tracer
}

// Start creates a span and a context.Context containing the newly-created span
// with `internal.visibility: "user"` so that it shows up in the Hasura Console.
func (t *Tracer) Start(ctx context.Context, spanName string, opts ...traceapi.SpanStartOption) (context.Context, traceapi.Span) {
	return t.tracer.Start(ctx, spanName, append(opts, userVisibilityAttribute)...)
}

// StartInternal creates a span and a context.Context containing the newly-created span.
// It won't show up in the Hasura Console
func (t *Tracer) StartInternal(ctx context.Context, spanName string, opts ...traceapi.SpanStartOption) (context.Context, traceapi.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

func httpStatusAttribute(code int) attribute.KeyValue {
	return attribute.Int("http_status", code)
}

func parseOTLPEndpoint(endpoint string, protocol string, insecurePtr *bool) (string, otlpProtocol, bool, error) {
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
		if uri.Port() == fmt.Sprint(otlpDefaultHTTPPort) {
			return host, otlpProtocol(protocol), insecure, nil
		}
		return host, otlpProtocolGRPC, insecure, nil
	default:
		return "", otlpProtocol(""), false, fmt.Errorf("invalid OTLP protocol %s", protocol)
	}
}

func parseOTLPCompression(input string) (string, int, error) {
	switch input {
	case otlpCompressionGzip, "":
		return otlpCompressionGzip, int(otlptracehttp.GzipCompression), nil
	case otlpCompressionNone:
		return otlpCompressionNone, int(otlptracehttp.NoCompression), nil
	default:
		return "", 0, fmt.Errorf("invalid OTLP compression type, accept none, gzip only")
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
		return otelMetricsExporterNone, fmt.Errorf("invalid metrics exporter type: %s", input)
	}
}
