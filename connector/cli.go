package connector

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cli struct {
	Serve struct {
		Configuration      string `help:"Configuration file path." env:"CONFIGURATION"`
		InlineConfig       bool   `help:"Inline JSON string or configuration file?" env:"INLINE_CONFIG"`
		Port               uint   `help:"Serve Port." env:"PORT" default:"8100"`
		ServiceTokenSecret string `help:"Service token secret." env:"SERVICE_TOKEN_SECRET"`
		OltpEndpoint       string `help:"OpenTelemetry receiver endpoint." env:"OTLP_ENDPOINT"`
		ServiceName        string `help:"OpenTelemetry service name." env:"OTEL_SERVICE_NAME"`
		LogLevel           string `help:"Log level." env:"LOG_LEVEL" enum:"trace,debug,info,warn,error" default:"info"`
	} `cmd:"" help:"Serve the NDC connector."`

	Configuration struct {
		Serve struct {
			Port     int    `help:"Serve Port." env:"PORT" default:"8100"`
			LogLevel string `help:"Log level." env:"LOG_LEVEL" default:"info"`
		} `cmd:"" help:"Serve the NDC configuration service."`
	} `cmd:"" help:"Configuration helpers."`
}

// Starts the connector.
// Will read runtime flags or environment variables to determine startup mode.
//
// This should be the entrypoint of your connector
func Start[RawConfiguration any, Configuration any, State any](connector Connector[RawConfiguration, Configuration, State], options ...ServeOption) error {
	cmd := kong.Parse(&cli)
	switch cmd.Command() {
	case "serve":
		logger, err := initLogger(cli.Serve.LogLevel)
		if err != nil {
			return err
		}

		server, err := NewServer[RawConfiguration, Configuration, State](connector, &ServerOptions{
			Configuration:      cli.Serve.Configuration,
			InlineConfig:       cli.Serve.InlineConfig,
			ServiceTokenSecret: cli.Serve.ServiceTokenSecret,
			OTLPEndpoint:       cli.Serve.OltpEndpoint,
			ServiceName:        cli.Serve.ServiceName,
		}, append(options, WithLogger(*logger))...)
		if err != nil {
			return err
		}
		return server.ListenAndServe(cli.Serve.Port)

	case "configuration serve":
		logger, err := initLogger(cli.Serve.LogLevel)
		if err != nil {
			return err
		}

		server := NewConfigurationServer[RawConfiguration, Configuration, State](connector, WithLogger(*logger))
		return server.ListenAndServe(cli.Serve.Port)
	default:
		return fmt.Errorf("Unknown command <%s>", cmd.Command())
	}
}

func initLogger(logLevel string) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(level)
	logger := log.Level(level)

	return &logger, nil
}
