# Observability

## OpenTelemetry

OpenTelemetry exporter is disabled by default unless one of `--otlp-endpoint`, `--otlp-traces-endpoint` or `--otlp-metrics-endpoint` argument is set. By default, the SDK treats port `4318` as HTTP protocol and gRPC protocol for others.

If `--otlp-*insecure` flags are not set, the SDK can also detect TLS connections via http(s).

Other configurations are inherited from the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go). Currently the SDK supports `traces` and `metrics`. See [Environment Variable Specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) and [OTLP Exporter Configuration](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/).

## Metrics

The SDK supports OTLP and Prometheus metrics exporters that is enabled by `--metrics-exporter` (`OTEL_METRICS_EXPORTER`) flag. Supported values:

- `none` (default): disable the exporter
- `otlp`: OTLP exporter
- `prometheus`: Prometheus exporter via the `/metrics` endpoint.

Prometheus exporter is served in the same HTTP server with the connector. If you want to run it on another port, configure the `--prometheus-port` (`OTEL_EXPORTER_PROMETHEUS_PORT`) flag.

### Built-in metrics

| Name                                             | Type      | Description                                                                  |
| ------------------------------------------------ | --------- | ---------------------------------------------------------------------------- |
| _(<prefix\>)_\_query_total                       | Counter   | Total number of query requests                                               |
| _(<prefix\>)_\_query_total_time                  | Histogram | Total time taken to plan and execute a query                                 |
| _(<prefix\>)_\_query_explain_total               | Counter   | Total number of explain query requests                                       |
| _(<prefix\>)_\_query_explain_total_time          | Histogram | Total time taken to plan and execute an explain query request, in seconds    |
| _(<prefix\>)_\_mutation_total                    | Counter   | Total number of mutation requests                                            |
| _(<prefix\>)_\_query_mutation_total_time         | Histogram | Total time taken to plan and execute a mutation request, in seconds          |
| _(<prefix\>)_\_mutation_explain_total            | Counter   | Total number of explain mutation requests                                    |
| _(<prefix\>)_\_query_mutation_explain_total_time | Histogram | Total time taken to plan and execute an explain mutation request, in seconds |

The prefix is empty by default. You can set the prefix for your connector by `WithMetricsPrefix` option.

## Logging

NDC Go SDK uses the standard [log/slog](https://pkg.go.dev/log/slog) that provides highly customizable and structured logging. By default, the logger is printed in JSON format and configurable level with `--log-level` (HASURA_LOG_LEVEL) flag. You can also replace it with different logging libraries that can wrap the `slog.Handler` interface, and set the logger with the `WithLogger` or `WithLoggerFunc` option.
