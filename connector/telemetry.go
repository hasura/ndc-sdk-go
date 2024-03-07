package connector

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-logr/zerologr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/rs/zerolog"
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
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	traceapi "go.opentelemetry.io/otel/trace"
)

var (
	userVisibilityAttribute = traceapi.WithAttributes(attribute.String("internal.visibility", "user"))
	successStatusAttribute  = attribute.String("status", "success")
	failureStatusAttribute  = attribute.String("status", "failure")
)

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
func setupOTelSDK(ctx context.Context, serverOptions *ServerOptions, serviceVersion, metricsPrefix string, logger zerolog.Logger) (*TelemetryState, error) {

	otel.SetLogger(zerologr.New(&logger))
	tracesEndpoint := serverOptions.OTLPTracesEndpoint
	if tracesEndpoint == "" {
		tracesEndpoint = serverOptions.OTLPEndpoint
	}
	metricsEndpoint := serverOptions.OTLPMetricsEndpoint
	if metricsEndpoint == "" {
		metricsEndpoint = serverOptions.OTLPEndpoint
	}

	// Set up resource.
	res, err := newResource(serverOptions.ServiceName, serviceVersion)
	if err != nil {
		return nil, err
	}

	var traceProvider *trace.TracerProvider
	if tracesEndpoint != "" {
		// Set up propagator.
		prop := newPropagator()
		otel.SetTextMapPropagator(prop)

		var traceExporter *otlptrace.Exporter

		// use grpc protocol by default if the scheme is empty
		if !strings.HasPrefix(tracesEndpoint, "http://") && !strings.HasPrefix(tracesEndpoint, "https://") {
			options := []otlptracegrpc.Option{
				otlptracegrpc.WithEndpoint(tracesEndpoint),
				otlptracegrpc.WithCompressor("gzip"),
			}

			if serverOptions.OTLPInsecure {
				options = append(options, otlptracegrpc.WithInsecure())
			}

			traceExporter, err = otlptracegrpc.New(ctx, options...)
			if err != nil {
				return nil, err
			}
		} else {
			// Set up trace exporter.
			endpointURL, err := url.Parse(tracesEndpoint)
			if err != nil {
				return nil, err
			}

			options := []otlptracehttp.Option{
				otlptracehttp.WithEndpoint(endpointURL.Host),
				otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
			}
			if endpointURL.Scheme == "http" {
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

	// disable default process and go collector metrics
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewGoCollector())

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	prometheusExporter, err := otelPrometheus.New()
	if err != nil {
		return nil, err
	}

	metricOptions := []metric.Option{
		metric.WithResource(res),
		metric.WithReader(prometheusExporter),
	}

	if metricsEndpoint != "" {
		// use grpc protocol by default if the scheme is empty
		if !strings.HasPrefix(metricsEndpoint, "http://") && !strings.HasPrefix(metricsEndpoint, "https://") {
			options := []otlpmetricgrpc.Option{
				otlpmetricgrpc.WithEndpoint(metricsEndpoint),
				otlpmetricgrpc.WithCompressor("gzip"),
			}

			if serverOptions.OTLPInsecure {
				options = append(options, otlpmetricgrpc.WithInsecure())
			}

			metricExporter, err := otlpmetricgrpc.New(ctx, options...)
			if err != nil {
				return nil, err
			}
			metricOptions = append(metricOptions, metric.WithReader(metric.NewPeriodicReader(metricExporter)))
		} else {

			endpointURL, err := url.Parse(metricsEndpoint)
			if err != nil {
				return nil, err
			}
			options := []otlpmetrichttp.Option{
				otlpmetrichttp.WithEndpoint(endpointURL.Host),
				otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
			}
			if endpointURL.Scheme == "http" {
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
		if tErr == nil && mErr == nil {
			return nil
		}
		var errorMsgs []string
		if tErr != nil {
			errorMsgs = append(errorMsgs, tErr.Error())
		}
		if mErr != nil {
			errorMsgs = append(errorMsgs, mErr.Error())
		}
		return errors.New(strings.Join(errorMsgs, ", "))
	}

	state := &TelemetryState{
		Tracer:   &Tracer{traceProvider.Tracer(serverOptions.ServiceName)},
		Meter:    meterProvider.Meter(serverOptions.ServiceName),
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
