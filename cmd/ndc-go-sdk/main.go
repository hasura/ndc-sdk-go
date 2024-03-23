package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type GenerateArguments struct {
	Path        string   `help:"The base path of the connector's source code" short:"p" default:"."`
	Directories []string `help:"Folders contain NDC operation functions" short:"d" default:"functions"`
	Trace       string   `help:"Enable tracing and write to target file path."`
}

type NewArguments struct {
	Name    string `help:"Name of the connector." short:"n" required:""`
	Module  string `help:"Module name of the connector" short:"m" required:""`
	Version string `help:"The version of ndc-sdk-go."`
	Output  string `help:"The location where source codes will be generated" short:"o" default:""`
}

var cli struct {
	LogLevel string            `help:"Log level." enum:"debug,info,warn,error" default:"info"`
	New      NewArguments      `cmd:"" help:"Initialize an NDC connector boilerplate. For example:\n ndc-go-sdk new -n example -m github.com/foo/example"`
	Generate GenerateArguments `cmd:"" help:"Generate schema and implementation for the connector from functions."`

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
	case "generate":
		log.Info().
			Str("path", cli.Generate.Path).
			Interface("directories", cli.Generate.Directories).
			Msg("generating connector schema...")

		moduleName, err := getModuleName(cli.Generate.Path)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get module name. The base path must contain a go.mod file")
		}

		if err := os.Chdir(cli.Generate.Path); err != nil {
			log.Fatal().Err(err).Msg("")
		}

		if err = parseAndGenerateConnector(&cli.Generate, moduleName); err != nil {
			log.Fatal().Err(err).Msg("failed to generate connector schema")
		}
		if err := execGoFormat("."); err != nil {
			log.Fatal().Err(err).Msg("failed to format code")
		}
		log.Info().Str("exec_time", time.Since(start).Round(time.Millisecond).String()).
			Msg("generated successfully")
	case "version":
		_, _ = fmt.Print(version.BuildVersion)
	default:
		log.Fatal().Msgf("unknown command <%s>", cmd.Command())
	}
}

func setupGlobalLogger(level string) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse log level: %s", level)
	}
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Level(logLevel).Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
