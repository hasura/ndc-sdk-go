package connector

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/hasura/gotel"
	"github.com/hasura/gotel/otelutils"
	"go.opentelemetry.io/otel/attribute"
	metricapi "go.opentelemetry.io/otel/metric"
	traceapi "go.opentelemetry.io/otel/trace"
)

var (
	successStatusAttribute = attribute.String("status", "success")
	failureStatusAttribute = attribute.String("status", "failure")
)

// TelemetryState contains OpenTelemetry exporters and basic connector metrics.
type TelemetryState struct {
	*connectorMetrics
	*gotel.OTelExporters
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
func setupOTelSDK(
	ctx context.Context,
	config *gotel.OTLPConfig,
	serviceVersion, metricsPrefix string,
	logger *slog.Logger,
) (*TelemetryState, error) {
	state, err := gotel.SetupOTelExporters(ctx, config, serviceVersion, logger)
	if err != nil {
		return nil, err
	}

	connectorMetrics, err := setupConnectorMetrics(state, metricsPrefix)
	if err != nil {
		return nil, err
	}

	return &TelemetryState{
		OTelExporters:    state,
		connectorMetrics: connectorMetrics,
	}, nil
}

func setupConnectorMetrics(
	telemetry *gotel.OTelExporters,
	metricsPrefix string,
) (*connectorMetrics, error) {
	if metricsPrefix != "" {
		metricsPrefix += "."
	}

	var err error

	meter := telemetry.Meter
	connectorMetrics := &connectorMetrics{}

	connectorMetrics.queryCounter, err = meter.Int64Counter(
		metricsPrefix+"query.total",
		metricapi.WithDescription("Total number of query requests"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.mutationCounter, err = meter.Int64Counter(
		metricsPrefix+"mutation.total",
		metricapi.WithDescription("Total number of mutation requests"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.queryExplainCounter, err = meter.Int64Counter(
		metricsPrefix+"query.explain_total",
		metricapi.WithDescription("Total number of explain query requests"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.mutationExplainCounter, err = meter.Int64Counter(
		metricsPrefix+"mutation.explain_total",
		metricapi.WithDescription("Total number of explain mutation requests"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.queryLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"query.total_time",
		metricapi.WithDescription("Total time taken to plan and execute a query, in seconds"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.mutationLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"mutation.total_time",
		metricapi.WithDescription("Total time taken to plan and execute a mutation, in seconds"),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.queryExplainLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"query.explain_total_time",
		metricapi.WithDescription(
			"Total time taken to plan and execute an explain query request, in seconds",
		),
	)
	if err != nil {
		return nil, err
	}

	connectorMetrics.mutationExplainLatencyHistogram, err = meter.Float64Histogram(
		metricsPrefix+"mutation.explain_total_time",
		metricapi.WithDescription(
			"Total time taken to plan and execute an explain mutation request, in seconds",
		),
	)

	return connectorMetrics, err
}

func httpStatusAttribute(code int) attribute.KeyValue {
	return attribute.Int("http_status", code)
}

// SetSpanHeaderAttributes sets header attributes to the otel span.
func SetSpanHeaderAttributes(
	span traceapi.Span,
	prefix string,
	httpHeaders http.Header,
	allowedHeaders ...string,
) {
	headers := otelutils.NewTelemetryHeaders(httpHeaders, allowedHeaders...)

	otelutils.SetSpanHeaderAttributes(span, prefix, headers)
}
