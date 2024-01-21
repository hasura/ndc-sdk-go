package connector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelPrometheus "go.opentelemetry.io/otel/exporters/prometheus"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	traceapi "go.opentelemetry.io/otel/trace"
)

type TelemetryState struct {
	Tracer                   traceapi.Tracer
	Meter                    metricapi.Meter
	Shutdown                 func(context.Context) error
	queryCounter             metricapi.Int64Counter
	mutationCounter          metricapi.Int64Counter
	explainCounter           metricapi.Int64Counter
	queryLatencyHistogram    metricapi.Float64Histogram
	explainLatencyHistogram  metricapi.Float64Histogram
	mutationLatencyHistogram metricapi.Float64Histogram
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context, endpoint, serviceName, serviceVersion, metricsPrefix string) (*TelemetryState, error) {

	var traceProvider *trace.TracerProvider
	if endpoint != "" {
		// Set up resource.
		res, err := newResource(serviceName, serviceVersion)
		if err != nil {
			return nil, err
		}

		// Set up propagator.
		prop := newPropagator()
		otel.SetTextMapPropagator(prop)

		// Set up trace exporter.
		traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(endpoint))
		if err != nil {
			return nil, err
		}

		bsp := trace.NewBatchSpanProcessor(traceExporter)
		traceProvider = trace.NewTracerProvider(
			trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
			trace.WithResource(res),
			trace.WithSpanProcessor(bsp),
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
	metricExporter, err := otelPrometheus.New()
	if err != nil {
		return nil, err
	}

	metricOptions := []metric.Option{
		metric.WithReader(metricExporter),
	}

	if endpoint != "" {
		httpMetricExporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(endpoint))
		if err != nil {
			return nil, err
		}
		metricOptions = append(metricOptions, metric.WithReader(metric.NewPeriodicReader(httpMetricExporter)))
	}

	meterProvider := metric.NewMeterProvider(metricOptions...)
	otel.SetMeterProvider(meterProvider)

	shutdownFunc := func(ctx context.Context) error {
		return errors.Join(
			traceProvider.Shutdown(ctx),
			meterProvider.Shutdown(ctx),
		)
	}

	state := &TelemetryState{
		Tracer:   traceProvider.Tracer(serviceName),
		Meter:    meterProvider.Meter(serviceName),
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

	telemetry.explainCounter, err = meter.Int64Counter(
		fmt.Sprintf("%sexplain.total", metricsPrefix),
		metricapi.WithDescription("The total number of explain requests"),
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

	telemetry.explainLatencyHistogram, err = meter.Float64Histogram(
		fmt.Sprintf("%sexplain.total_time", metricsPrefix),
		metricapi.WithDescription("Total time taken to plan and execute an explain request, in seconds"),
	)

	return err
}
