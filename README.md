# Native Data Connector SDK for Go

This SDK is mostly analogous to the Rust SDK, except where necessary.

All functions of the Connector interface are analogous to their Rust counterparts, with the addition of `GetRawConfigurationSchema` which does exactly what it sounds like.


## Using this SDK

The SDK exports a `Start` function, which takes a `connector` object, that is an object that implements the `Connector` interface defined in [connector/types.go](connector/types.go)

This function should be your starting point.

A connector can thus start likes so:

```go
import "github.com/hasura/ndc-sdk-go/connector"

type RawConfiguration struct{ ... }
type Configuration struct { ... }
type State struct { ... }
type Connector struct { ... }

/* implementation of the Connector interface removed for brevity */

func main() {
  
	if err := connector.Start[RawConfiguration, Configuration, State](&Connector{}); err != nil {
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
      --configuration=STRING            Configuration file path ($CONFIGURATION).
      --inline-config                   Inline JSON string or configuration file? ($INLINE_CONFIG)
      --port=8100                       Serve Port ($PORT).
      --service-token-secret=STRING     Service token secret ($SERVICE_TOKEN_SECRET).
      --otlp-endpoint=STRING            OpenTelemetry receiver endpoint that is set as default for all types ($OTLP_ENDPOINT).
      --otlp-traces-endpoint=STRING     OpenTelemetry endpoint for traces ($OTLP_TRACES_ENDPOINT).
      --otlp-insecure                   Disable LTS for OpenTelemetry gRPC exporters ($OTLP_INSECURE).
      --otlp-metrics-endpoint=STRING    OpenTelemetry endpoint for metrics ($OTLP_METRICS_ENDPOINT).
      --service-name=STRING             OpenTelemetry service name ($OTEL_SERVICE_NAME).
      --log-level="info"                Log level ($LOG_LEVEL).

  configuration serve
    Serve the NDC configuration service.
		
    Flags:
      --port=8100           Serve Port ($PORT).
      --log-level="info"    Log level ($LOG_LEVEL).
```

Please refer to the [NDC Spec](https://hasura.github.io/ndc-spec/) for details on implementing the Connector interface, or see [examples](./examples).  

## Regenerating Schema Types

The NDC spec types are borrowed from ndc-sdk-typescript that are generated from the NDC Spec Rust types.
Then the Go types are generated from that JSON Schema document into `./schema/schema.generated.go`.

In order to regenerate the types, run

```
> cd typegen && ./regenerate-schema.sh
```
