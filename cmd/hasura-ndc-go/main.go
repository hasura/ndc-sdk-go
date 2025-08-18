package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/hasura/ndc-sdk-go/v2/cmd/hasura-ndc-go/command"
	"github.com/hasura/ndc-sdk-go/v2/cmd/hasura-ndc-go/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cli struct {
	LogLevel string                   `default:"info" enum:"debug,info,warn,error,DEBUG,INFO,WARN,ERROR" env:"HASURA_PLUGIN_LOG_LEVEL" help:"Log level."`
	New      command.NewArguments     `cmd:"" help:"Initialize an NDC connector boilerplate. For example:\n hasura-ndc-go new -n example -m github.com/foo/example"`
	Update   command.UpdateArguments  `cmd:"" help:"Generate schema and implementation for the connector from functions."`
	Upgrade  command.UpgradeArguments `cmd:"" help:"Patch source codes to upgrade the connector"`
	Generate struct {
		Snapshots command.GenTestSnapshotArguments `cmd:"" help:"Generate test snapshots."`
	} `cmd:"" help:"Generator helpers."`

	Version struct{} `cmd:"" help:"Print the CLI version."`
}

func main() {
	cmd := kong.Parse(&cli, kong.UsageOnError())
	start := time.Now()

	setupGlobalLogger(cli.LogLevel)

	switch cmd.Command() {
	case "new":
		if cli.New.Version == "" {
			cli.New.Version = version.BuildVersion
		}

		log.Info().
			Str("name", cli.New.Name).
			Str("module", cli.New.Module).
			Str("output", cli.New.Output).
			Str("version", cli.New.Version).
			Msg("generating the NDC boilerplate...")

		if err := command.GenerateNewProject(&cli.New, false); err != nil {
			log.Fatal().Err(err).Msg("failed to generate new project")
		}

		log.Info().Str("exec_time", time.Since(start).Round(time.Second).String()).
			Msg("generated successfully")
	case "update":
		command.UpdateConnectorSchema(cli.Update, start)
	case "upgrade":
		err := command.UpgradeConnector(cli.Upgrade)
		if err != nil {
			log.Error().Err(err).
				Msg("failed to upgrade the connector. Check out the migration guide at https://github.com/hasura/ndc-sdk-go/blob/main/docs/migrate-v1-to-v2.md and update manually")
		}
	case "generate snapshots":
		log.Info().
			Str("endpoint", cli.Generate.Snapshots.Endpoint).
			Str("path", cli.Generate.Snapshots.Schema).
			Interface("dir", cli.Generate.Snapshots.Dir).
			Msg("generating test snapshots...")

		if err := command.GenTestSnapshots(&cli.Generate.Snapshots); err != nil {
			log.Fatal().Err(err).Msg("failed to generate test snapshots")
		}
	case "version":
		_, _ = fmt.Fprint(os.Stderr, version.BuildVersion)
	default:
		log.Fatal().Msgf("unknown command <%s>", cmd.Command())
	}
}

func setupGlobalLogger(level string) {
	logLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse log level: %s", level)
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Level(logLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
