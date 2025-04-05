package connector

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// ServeCommandArguments contains argument flags of the serve command.
type ServeCommandArguments struct {
	OTLPConfig
	HTTPServerConfig

	Configuration      string `env:"HASURA_CONFIGURATION_DIRECTORY" help:"Configuration directory"`
	Port               uint   `env:"HASURA_CONNECTOR_PORT"          help:"Serve Port"              default:"8080"`
	ServiceTokenSecret string `env:"HASURA_SERVICE_TOKEN_SECRET"    help:"Service token secret"`
}

// ServeCLI is used for CLI argument binding.
type ServeCLI struct {
	LogLevel string                `default:"info" enum:"trace,debug,info,warn,error" env:"HASURA_LOG_LEVEL" help:"Log level."`
	Serve    ServeCommandArguments `                                                                         help:"Serve the NDC connector." cmd:""`
}

// GetServeCLI returns the inner serve cli.
func (cli *ServeCLI) GetServeCLI() *ServeCLI {
	return cli
}

// Execute executes the command.
func (cli *ServeCLI) Execute(ctx context.Context, command string) error {
	return fmt.Errorf("unknown command <%s>", command)
}

// ConnectorCLI abstracts the connector CLI so NDC authors can extend it.
type ConnectorCLI interface {
	GetServeCLI() *ServeCLI
	Execute(ctx context.Context, command string) error
}

// Starts the connector.
// Will read command line arguments or environment variables to determine runtime configuration.
//
// This should be the entrypoint of your connector.
func Start[Configuration any, State any](
	connector Connector[Configuration, State],
	options ...ServeOption,
) error {
	var cli ServeCLI

	return StartCustom(&cli, connector, options...)
}

// Starts the connector with custom CLI.
// Will read command line arguments or environment variables to determine runtime configuration.
//
// This should be the entrypoint of your connector.
func StartCustom[Configuration any, State any](
	cli ConnectorCLI,
	connector Connector[Configuration, State],
	options ...ServeOption,
) error {
	cmd := kong.Parse(cli, kong.UsageOnError())
	serveCLI := cli.GetServeCLI()
	command := cmd.Command()

	logger, logLevel, err := initLogger(serveCLI.LogLevel)
	if err != nil {
		return err
	}

	switch command {
	case "serve":
		logger.Debug("Start to serve the connector", slog.Any("arguments", serveCLI.Serve))

		server, err := NewServer(connector, &ServerOptions{
			Configuration:      serveCLI.Serve.Configuration,
			ServiceTokenSecret: serveCLI.Serve.ServiceTokenSecret,
			OTLPConfig:         serveCLI.Serve.OTLPConfig,
			HTTPServerConfig:   serveCLI.Serve.HTTPServerConfig,
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

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	return logger, level, nil
}
