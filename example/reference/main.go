package main

import (
	"context"
	"fmt"

	"github.com/hasura/ndc-sdk-go/connector"
)

func main() {
	var cli CLI
	if err := connector.StartCustom[Configuration, State](
		&cli,
		&Connector{},
		connector.WithMetricsPrefix("ndc_ref"),
		connector.WithDefaultServiceName("ndc_ref"),
	); err != nil {
		panic(err)
	}
}

type CLI struct {
	connector.ServeCLI
	Version struct{} `cmd:"" help:"Print the version."`
}

func (cli *CLI) Execute(ctx context.Context, command string) error {
	logger := connector.GetLogger(ctx)
	switch command {
	case "version":
		logger.Info().Msg("v0.1.0")
		return nil
	default:
		return fmt.Errorf("unknown command <%s>", command)
	}
}
