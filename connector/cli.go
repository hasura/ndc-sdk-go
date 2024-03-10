package connector

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// ServeCommandArguments contains argument flags of the serve command
type ServeCommandArguments struct {
	Configuration       string `help:"Configuration directory." env:"HASURA_CONFIGURATION_DIRECTORY"`
	Port                uint   `help:"Serve Port." env:"HASURA_CONNECTOR_PORT" default:"8080"`
	ServiceTokenSecret  string `help:"Service token secret." env:"HASURA_SERVICE_TOKEN_SECRET"`
	OtlpEndpoint        string `help:"OpenTelemetry receiver endpoint that is set as default for all types." env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtlpTracesEndpoint  string `help:"OpenTelemetry endpoint for traces." env:"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"`
	OtlpInsecure        bool   `help:"Disable LTS for OpenTelemetry gRPC exporters." env:"OTEL_EXPORTER_OTLP_INSECURE"`
	OtlpMetricsEndpoint string `help:"OpenTelemetry endpoint for metrics." env:"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"`
	ServiceName         string `help:"OpenTelemetry service name." env:"OTEL_SERVICE_NAME"`
}

// ServeCLI is used for CLI argument binding
type ServeCLI struct {
	LogLevel string                `help:"Log level." env:"HASURA_LOG_LEVEL" enum:"trace,debug,info,warn,error" default:"info"`
	Serve    ServeCommandArguments `cmd:"" help:"Serve the NDC connector."`
}

// GetServeCLI returns the inner serve cli
func (cli *ServeCLI) GetServeCLI() *ServeCLI {
	return cli
}

// Execute executes the command
func (cli *ServeCLI) Execute(ctx context.Context, command string) error {
	return fmt.Errorf("unknown command <%s>", command)
}

// ConnectorCLI abstracts the connector CLI so NDC authors can extend it
type ConnectorCLI interface {
	GetServeCLI() *ServeCLI
	Execute(ctx context.Context, command string) error
}

// Starts the connector.
// Will read command line arguments or environment variables to determine runtime configuration.
//
// This should be the entrypoint of your connector
func Start[Configuration any, State any](connector Connector[Configuration, State], options ...ServeOption) error {
	var cli ServeCLI
	return StartCustom[Configuration, State](&cli, connector)
}

// Starts the connector with custom CLI.
// Will read command line arguments or environment variables to determine runtime configuration.
//
// This should be the entrypoint of your connector
func StartCustom[Configuration any, State any](cli ConnectorCLI, connector Connector[Configuration, State], options ...ServeOption) error {
	cmd := kong.Parse(cli, kong.UsageOnError())
	serveCLI := cli.GetServeCLI()
	command := cmd.Command()
	logger, logLevel, err := initLogger(serveCLI.LogLevel)
	if err != nil {
		return err
	}
	switch command {
	case "serve":
		server, err := NewServer[Configuration, State](connector, &ServerOptions{
			Configuration:       serveCLI.Serve.Configuration,
			ServiceTokenSecret:  serveCLI.Serve.ServiceTokenSecret,
			OTLPEndpoint:        serveCLI.Serve.OtlpEndpoint,
			OTLPInsecure:        serveCLI.Serve.OtlpInsecure,
			OTLPTracesEndpoint:  serveCLI.Serve.OtlpTracesEndpoint,
			OTLPMetricsEndpoint: serveCLI.Serve.OtlpMetricsEndpoint,
			ServiceName:         serveCLI.Serve.ServiceName,
		}, append([]ServeOption{WithLogger(logger), withLogLevel(logLevel)}, options...)...)
		if err != nil {
			return err
		}
		return server.ListenAndServe(serveCLI.Serve.Port)
	default:
		ctx := context.WithValue(context.Background(), logContextKey, logger)
		return cli.Execute(ctx, command)
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
