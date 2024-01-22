# Observability services

This project sets up an observability stack that is ready to be used with examples.  

## Getting Started

```bash
docker-compose up -d
```

| Service                             | Endpoint                      |
| ----------------------------------- | ----------------------------- |
| Prometheus                          | http://localhost:9090         |
| Jaeger Dashboard                    | http://localhost:16686        |
| Jaeger OpenTelemetry (gRPC)         | localhost:14317               |
| Jaeger OpenTelemetry (HTTP)         | http://localhost:14318        |
| OpenTelemetry collector (gRPC)      | localhost:4317                |
| OpenTelemetry collector (HTTP)      | http://localhost:4318         |
| OpenTelemetry collector (metrics)   | http://localhost:8889/metrics |

## Using with connectors

Add the `--otlp-traces-endpoint localhost:14317 --otlp-metrics-endpoint localhost:4317 --otlp-insecure` flag when running connector serve. For example:

```sh
go run ./example/reference serve --otlp-traces-endpoint localhost:14317 --otlp-metrics-endpoint localhost:4317 --otlp-insecure --inline-config --configuration {}
```