# Native Data Connector SDK for Go

This SDK is mostly analogous to the Rust SDK, except where necessary.

All functions of the Connector interface are analogous to their Rust counterparts.

## Features

The SDK fully supports [NDC Specification v0.1.0](https://github.com/hasura/ndc-spec/tree/v0.1.0) and [Connector Deployment spec](https://github.com/hasura/ndc-hub/blob/main/rfcs/0000-deployment.md) with following features:

- Connector HTTP server
- Authentication
- Observability with OpenTelemetry and Prometheus

## Quick start

Checkout the [generation tool](cmd/ndc-go-sdk) to quickly setup and develop data connectors.

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
      --service-name=STRING                OpenTelemetry service name ($OTEL_SERVICE_NAME).
      --otlp-endpoint=STRING               OpenTelemetry receiver endpoint that is set as default for all types ($OTEL_EXPORTER_OTLP_ENDPOINT).
      --otlp-traces-endpoint=STRING        OpenTelemetry endpoint for traces ($OTEL_EXPORTER_OTLP_TRACES_ENDPOINT).
      --otlp-metrics-endpoint=STRING       OpenTelemetry endpoint for metrics ($OTEL_EXPORTER_OTLP_METRICS_ENDPOINT).
      --otlp-insecure                      Disable LTS for OpenTelemetry exporters ($OTEL_EXPORTER_OTLP_INSECURE).
      --otlp-traces-insecure               Disable LTS for OpenTelemetry traces exporter ($OTEL_EXPORTER_OTLP_TRACES_INSECURE).
      --otlp-metrics-insecure              Disable LTS for OpenTelemetry metrics exporter ($OTEL_EXPORTER_OTLP_METRICS_INSECURE).
      --otlp-protocol=STRING               OpenTelemetry receiver protocol for all types ($OTEL_EXPORTER_OTLP_PROTOCOL).
      --otlp-traces-protocol=STRING        OpenTelemetry receiver protocol for traces ($OTEL_EXPORTER_OTLP_TRACES_PROTOCOL).
      --otlp-metrics-protocol=STRING       OpenTelemetry receiver protocol for metrics ($OTEL_EXPORTER_OTLP_METRICS_PROTOCOL).
      --otlp-compression="gzip"            Enable compression for OTLP exporters. Accept: none, gzip ($OTEL_EXPORTER_OTLP_COMPRESSION)
      --otlp-trace-compression="gzip"      Enable compression for OTLP traces exporter. Accept: none, gzip ($OTEL_EXPORTER_OTLP_TRACES_COMPRESSION)
      --otlp-metrics-compression="gzip"    Enable compression for OTLP metrics exporter. Accept: none, gzip ($OTEL_EXPORTER_OTLP_METRICS_COMPRESSION).
      --metrics-exporter="none"            Metrics export type. Accept: none, otlp, prometheus ($OTEL_METRICS_EXPORTER).
      --prometheus-port=PROMETHEUS-PORT    Prometheus port for the Prometheus HTTP server. Use /metrics endpoint of the connector server if empty ($OTEL_EXPORTER_PROMETHEUS_PORT)
      --configuration=STRING               Configuration directory ($HASURA_CONFIGURATION_DIRECTORY).
      --port=8080                          Serve Port ($HASURA_CONNECTOR_PORT).
      --service-token-secret=STRING        Service token secret ($HASURA_SERVICE_TOKEN_SECRET).
```

Please refer to the [NDC Spec](https://hasura.github.io/ndc-spec/) for details on implementing the Connector interface, or see [examples](./example).

## Observability

### OpenTelemetry

OpenTelemetry exporter is disabled by default unless one of `--otlp-endpoint`, `--otlp-traces-endpoint` or `--otlp-metrics-endpoint` argument is set. By default, the SDK treats port `4318` as HTTP protocol and gRPC protocol for others.

If `--otlp-*insecure` flags are not set, the SDK can also detect TLS connections via http(s).

Other configurations are inherited from the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go). See [Environment Variable Specification](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) and [OTLP Exporter Configuration](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/).

### Metrics

The SDK supports OTLP and Prometheus metrics exporters that is enabled by `--metrics-exporter` (`OTEL_METRICS_EXPORTER`) flag. Supported values:

- `none` (default): disable the exporter
- `otlp`: OTLP exporter
- `prometheus`: Prometheus exporter via the `/metrics` endpoint.

Prometheus exporter is served in the same HTTP server with the connector. If you want to run it on another port, configure the `--prometheus-port` (`OTEL_EXPORTER_PROMETHEUS_PORT`) flag.

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

In order to regenerate the types, run

```
> cd typegen && ./regenerate-schema.sh
```
