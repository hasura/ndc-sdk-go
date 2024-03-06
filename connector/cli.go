package connector

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

var cli struct {
	Serve struct {
		Configuration       string `help:"Configuration directory." env:"HASURA_CONFIGURATION_DIRECTORY"`
		Port                uint   `help:"Serve Port." env:"HASURA_CONNECTOR_PORT" default:"8080"`
		ServiceTokenSecret  string `help:"Service token secret." env:"HASURA_SERVICE_TOKEN_SECRET"`
		OtlpEndpoint        string `help:"OpenTelemetry receiver endpoint that is set as default for all types." env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
		OtlpTracesEndpoint  string `help:"OpenTelemetry endpoint for traces." env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
		OtlpInsecure        bool   `help:"Disable LTS for OpenTelemetry gRPC exporters." env:"OTEL_EXPORTER_OTLP_INSECURE"`
		OtlpMetricsEndpoint string `help:"OpenTelemetry endpoint for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
		ServiceName         string `help:"OpenTelemetry service name." env:"OTEL_SERVICE_NAME"`
		LogLevel            string `help:"Log level." env:"HASURA_LOG_LEVEL" enum:"trace,debug,info,warn,error" default:"info"`
	} `cmd:"" help:"Serve the NDC connector."`
}

// Starts the connector.
// Will read command line arguments or environment variables to determine runtime configuration.
//
// This should be the entrypoint of your connector
func Start[Configuration any, State any](connector Connector[Configuration, State], options ...ServeOption) error {
	cmd := kong.Parse(&cli)
	switch cmd.Command() {
	case "serve":
		logger, logLevel, err := initLogger(cli.Serve.LogLevel)
		if err != nil {
			return err
		}

		server, err := NewServer[Configuration, State](connector, &ServerOptions{
			Configuration:       cli.Serve.Configuration,
			ServiceTokenSecret:  cli.Serve.ServiceTokenSecret,
			OTLPEndpoint:        cli.Serve.OtlpEndpoint,
			OTLPInsecure:        cli.Serve.OtlpInsecure,
			OTLPTracesEndpoint:  cli.Serve.OtlpTracesEndpoint,
			OTLPMetricsEndpoint: cli.Serve.OtlpMetricsEndpoint,
			ServiceName:         cli.Serve.ServiceName,
		}, append([]ServeOption{WithLogger(logger), withLogLevel(logLevel)}, options...)...)
		if err != nil {
			return err
		}
		return server.ListenAndServe(cli.Serve.Port)
	default:
		return fmt.Errorf("unknown command <%s>", cmd.Command())
	}
}

func initLogger(logLevel string) (*slog.Logger, slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(strings.ToUpper(logLevel)))
	if err != nil {
		return nil, level, err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	return logger, level, nil
}
