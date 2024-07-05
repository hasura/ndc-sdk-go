package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command"
	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// UpdateArguments represent input arguments of the `update` command
type UpdateArguments struct {
	Path         string   `help:"The path of the root directory where the go.mod file is present" short:"p" env:"HASURA_PLUGIN_CONNECTOR_CONTEXT_PATH" default:"."`
	ConnectorDir string   `help:"The directory where the connector.go file is placed" default:"."`
	PackageTypes string   `help:"The name of types package where the State struct is in"`
	Directories  []string `help:"Folders contain NDC operation functions" short:"d"`
	Trace        string   `help:"Enable tracing and write to target file path."`
}

type NewArguments struct {
	Name    string `help:"Name of the connector." short:"n" required:""`
	Module  string `help:"Module name of the connector" short:"m" required:""`
	Version string `help:"The version of ndc-sdk-go."`
	Output  string `help:"The location where source codes will be generated" short:"o" default:""`
}

var cli struct {
	LogLevel string          `help:"Log level." enum:"debug,info,warn,error,DEBUG,INFO,WARN,ERROR" env:"HASURA_PLUGIN_LOG_LEVEL" default:"info"`
	New      NewArguments    `cmd:"" help:"Initialize an NDC connector boilerplate. For example:\n hasura-ndc-go new -n example -m github.com/foo/example"`
	Update   UpdateArguments `cmd:"" help:"Generate schema and implementation for the connector from functions."`
	Generate UpdateArguments `cmd:"" help:"(deprecated) The alias of the 'update' command."`
	Test     struct {
		Snapshots command.GenTestSnapshotArguments `cmd:"" help:"Generate test snapshots."`
	} `cmd:"" help:"Test helpers."`

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

		if err := generateNewProject(&cli.New, false); err != nil {
			log.Fatal().Err(err).Msg("failed to generate new project")
		}
		log.Info().Str("exec_time", time.Since(start).Round(time.Second).String()).
			Msg("generated successfully")
	case "update":
		execUpdate(&cli.Update, start)
	case "generate":
		execUpdate(&cli.Generate, start)
	case "test snapshots":
		log.Info().
			Str("endpoint", cli.Test.Snapshots.Endpoint).
			Str("path", cli.Test.Snapshots.Schema).
			Interface("dir", cli.Test.Snapshots.Dir).
			Msg("generating test snapshots...")

		if err := command.GenTestSnapshots(&cli.Test.Snapshots); err != nil {
			log.Fatal().Err(err).Msg("failed to generate test snapshots")
		}
	case "version":
		_, _ = fmt.Print(version.BuildVersion)
	default:
		log.Fatal().Msgf("unknown command <%s>", cmd.Command())
	}
}

func execUpdate(args *UpdateArguments, start time.Time) {
	log.Info().
		Str("path", args.Path).
		Str("connector_dir", args.ConnectorDir).
		Str("package_types", args.PackageTypes).
		Msg("generating connector schema...")

	moduleName, err := getModuleName(args.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get module name. The base path must contain a go.mod file")
	}

	if err := os.Chdir(args.Path); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	if err = parseAndGenerateConnector(args, moduleName); err != nil {
		log.Fatal().Err(err).Msg("failed to generate connector schema")
	}
	if err := execGoFormat("."); err != nil {
		log.Fatal().Err(err).Msg("failed to format code")
	}
	log.Info().Str("exec_time", time.Since(start).Round(time.Millisecond).String()).
		Msg("generated successfully")
}

func setupGlobalLogger(level string) {
	logLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse log level: %s", level)
	}
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Level(logLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
