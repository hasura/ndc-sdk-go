# Native Data Connector SDK for Go

This SDK is mostly analogous to the Rust SDK, except where necessary.

All functions of the Connector interface are analogous to their Rust counterparts.

## Features

The SDK fully supports [NDC Specification v0.1.0](https://github.com/hasura/ndc-spec/tree/v0.1.0) and [Connector Deployment spec](https://github.com/hasura/ndc-hub/blob/main/rfcs/0000-deployment.md) with following features:

- Connector HTTP server
- Authentication
- Observability with OpenTelemetry and Prometheus

## Prerequisites

- Go 1.21+

> Downgrade to SDK v0.x If you are using Go v1.19+

## Quick start

Check out the [generation tool](cmd/ndc-go-sdk) to quickly setup and develop data connectors.

## Using this SDK

The SDK exports a `Start` function, which takes a `connector` object, that is an object that implements the `Connector` interface defined in [connector/types.go](connector/types.go)

This function should be your starting point.

A connector can thus start likes so:

```go
import "github.com/hasura/ndc-sdk-go/connector"

type Configuration struct { ... }
type State struct { ... }
type Connector struct { ... }

/* implementation of the Connector interface removed for brevity */

func main() {

	if err := connector.Start[Configuration, State](&Connector{}); err != nil {
		panic(err)
	}
}
```

The `Start` function create a CLI application with following commands:

```sh
Commands:
  serve
    Serve the NDC connector.

    Flags:
      --configuration=STRING            Configuration directory ($HASURA_CONFIGURATION_DIRECTORY).
      --port=8080                       Serve Port ($HASURA_CONNECTOR_PORT).
      --service-token-secret=STRING     Service token secret ($HASURA_SERVICE_TOKEN_SECRET).
      --otlp-endpoint=STRING            OpenTelemetry receiver endpoint that is set as default for all types ($OTEL_EXPORTER_OTLP_ENDPOINT).
      --otlp-traces-endpoint=STRING     OpenTelemetry endpoint for traces ($OTEL_EXPORTER_OTLP_TRACES_ENDPOINT).
      --otlp-insecure                   Disable LTS for OpenTelemetry gRPC exporters ($OTEL_EXPORTER_OTLP_INSECURE).
      --otlp-metrics-endpoint=STRING    OpenTelemetry endpoint for metrics ($OTEL_EXPORTER_OTLP_METRICS_ENDPOINT).
      --service-name=STRING             OpenTelemetry service name ($OTEL_SERVICE_NAME).
      --log-level="info"                Log level ($HASURA_LOG_LEVEL).
```

Please refer to the [NDC Spec](https://hasura.github.io/ndc-spec/) for details on implementing the Connector interface, or see [examples](./example).

## Observability

### OpenTelemetry

OpenTelemetry exporter is disabled by default unless one of `--otlp-endpoint`, `--otlp-traces-endpoint` or `--otlp-metrics-endpoint` argument is set. The SDK automatically detects either HTTP or gRPC protocol by the URL scheme. For example:

- `http://localhost:4318`: HTTP
- `localhost:4317`: gRPC

The SDK can also detect TLS connections via http(s). However, if you want to disable TLS for gRPC, you must add `--otlp-insecure` the flag.

Other configurations are inherited from the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go). Currently the SDK supports `traces` and `metrics`. See [Environment Variable Specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) and [OTLP Exporter Configuration](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/).

### Prometheus

Prometheus metrics are exported via the `/metrics` endpoint.

### Logging

NDC Go SDK uses the standard [log/slog](https://pkg.go.dev/log/slog) that provides highly customizable and structured logging. By default, the logger is printed in JSON format and configurable level with `--log-level` (HASURA_LOG_LEVEL) flag. You also can replace it with different logging libraries that can wrap the `slog.Handler` interface, and set the logger with the `WithLogger` or `WithLoggerFunc` option.

## Regenerating Schema Types

The NDC spec types are borrowed from ndc-sdk-typescript that are generated from the NDC Spec Rust types.
Then the Go types are generated from that JSON Schema document into `./schema/schema.generated.go`.

In order to regenerate the types, run

```
> cd typegen && ./regenerate-schema.sh
```
