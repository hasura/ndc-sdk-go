# Observability services

This project sets up an observability stack that is ready to be used with examples.  

## Getting Started

```bash
docker-compose up -d
```

| Service                             | Endpoint                      |
| ----------------------------------- | ----------------------------- |
| Prometheus                          | http://localhost:9090         |
| Jaeger                              | http://localhost:16686        |
| OpenTelemetry collector (HTTP)      | http://localhost:4318         |
| OpenTelemetry collector (metrics)   | http://localhost:8889/metrics |

## Using with connectors

Add the `--otlp-traces-endpoint http://localhost:14318 --otlp-metrics-endpoint http://localhost:4318` flag when running connector serve. For example:

```sh
go run ./example/reference serve --service-name ndc-csv --otlp-traces-endpoint http://localhost:14318 --otlp-metrics-endpoint http://localhost:4318 --inline-config --configuration {}
```