# Native Data Connector SDK for Go

This SDK is mostly analogous to the Rust SDK, except where necessary.

All functions of the Connector interface are analogous to their Rust counterparts.

## Features

The SDK fully supports [NDC Specification v0.1.1](https://hasura.github.io/ndc-spec/specification/changelog.html#011) and [Connector Deployment spec](https://github.com/hasura/ndc-hub/blob/main/rfcs/0000-deployment.md) with following features:

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
      --configuration=STRING               Configuration directory ($HASURA_CONFIGURATION_DIRECTORY).
      --port=8080                          Serve Port ($HASURA_CONNECTOR_PORT).
      --service-token-secret=STRING        Service token secret ($HASURA_SERVICE_TOKEN_SECRET).
      --service-name=STRING                OpenTelemetry service name ($OTEL_SERVICE_NAME).
      --otlp-endpoint=STRING               OpenTelemetry receiver endpoint that is set as default for all types ($OTEL_EXPORTER_OTLP_ENDPOINT).
      --metrics-exporter="none"            Metrics export type. Accept: none, otlp, prometheus ($OTEL_METRICS_EXPORTER).
      --prometheus-port=PROMETHEUS-PORT    Prometheus port for the Prometheus HTTP server. Use /metrics endpoint of the connector server if empty ($OTEL_EXPORTER_PROMETHEUS_PORT)
      ...
```

Please refer to the [NDC Spec](https://hasura.github.io/ndc-spec/) for details on implementing the Connector interface, or see [examples](./example).

## Observability

### OpenTelemetry

OpenTelemetry exporter is disabled by default unless one of `--otlp-endpoint`, `--otlp-traces-endpoint` or `--otlp-metrics-endpoint` argument is set. By default, the SDK treats port `4318` as HTTP protocol and gRPC protocol for others.

If `--otlp-*insecure` flags are not set, the SDK can also detect TLS connections via http(s).

Other configurations are inherited from the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go). Currently the SDK supports `traces` and `metrics`. See [Environment Variable Specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) and [OTLP Exporter Configuration](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/).

### Metrics

The SDK supports OTLP and Prometheus metrics exporters that is enabled by `--metrics-exporter` (`OTEL_METRICS_EXPORTER`) flag. Supported values:

- `none` (default): disable the exporter
- `otlp`: OTLP exporter
- `prometheus`: Prometheus exporter via the `/metrics` endpoint.

Prometheus exporter is served in the same HTTP server with the connector. If you want to run it on another port, configure the `--prometheus-port` (`OTEL_EXPORTER_PROMETHEUS_PORT`) flag.

#### Built-in metrics

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

### Logging

NDC Go SDK uses the standard [log/slog](https://pkg.go.dev/log/slog) that provides highly customizable and structured logging. By default, the logger is printed in JSON format and configurable level with `--log-level` (HASURA_LOG_LEVEL) flag. You also can replace it with different logging libraries that can wrap the `slog.Handler` interface, and set the logger with the `WithLogger` or `WithLoggerFunc` option.

## Customize the CLI

The SDK uses [Kong](https://github.com/alecthomas/kong), a lightweight command-line parser to implement the CLI interface.

The default CLI already implements the `serve` command, so you don't need to do anything. However, it's also easy to extend if you want to add more custom commands.

The SDK abstracts an interface for the CLI that requires embedding the base `ServeCLI` command and can execute other commands.

```go
type ConnectorCLI interface {
	GetServeCLI() *ServeCLI
	Execute(ctx context.Context, command string) error
}
```

And use the `StartCustom` function to start the CLI.

See the [custom CLI example](./example/reference/main.go) in the reference connector.

## Regenerating Schema Types

The NDC spec types are borrowed from ndc-sdk-typescript that are generated from the NDC Spec Rust types.
Then the Go types are generated from that JSON Schema document into `./schema/schema.generated.go`.

In order to regenerate the types, download the json schema and run:

```
> cd typegen && ./regenerate-schema.sh
```
