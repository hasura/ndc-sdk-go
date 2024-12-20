# Observability services

This project sets up an observability stack that is ready to be used with examples.

## Getting Started

```bash
docker-compose up -d
```

| Service                           | Endpoint                      |
| --------------------------------- | ----------------------------- |
| Prometheus                        | http://localhost:9090         |
| Jaeger Dashboard                  | http://localhost:16686        |
| Jaeger OpenTelemetry (gRPC)       | localhost:14317               |
| Jaeger OpenTelemetry (HTTP)       | http://localhost:14318        |
| OpenTelemetry collector (gRPC)    | localhost:4317                |
| OpenTelemetry collector (HTTP)    | http://localhost:4318         |
| OpenTelemetry collector (metrics) | http://localhost:8889/metrics |

## Using with connectors

Add `--otlp-xxx` flags or environment variables when running connector serve. For example:

```sh
OTEL_EXPORTER_OTLP_HEADERS="Authorization=Bearer randomtoken" \
OTEL_RESOURCE_ATTRIBUTES="foo=bar" \
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 \
OTEL_LOGS_EXPORTER=otlp \
OTEL_METRICS_EXPORTER=otlp \
OTEL_EXPORTER_OTLP_INSECURE=true \
  go run ./example/reference serve
```
