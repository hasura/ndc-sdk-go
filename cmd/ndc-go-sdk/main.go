package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cli struct {
	Init struct {
		Name     string `help:"Name of the connector." short:"n" required:""`
		Module   string `help:"Module name of the connector" short:"m" required:""`
		Output   string `help:"The location where source codes will be generated" short:"o" default:""`
		LogLevel string `help:"Log level." enum:"trace,debug,info,warn,error" default:"info"`
	} `cmd:"" help:"Initialize an NDC connector boilerplate. For example:\n ndc-go-sdk init -n example -m github.com/foo/example"`

	Generate struct {
		Path        string   `help:"The base path of the connector's source code" short:"p" default:"."`
		Directories []string `help:"Folders contain NDC operation functions" short:"d" default:"functions"`
		LogLevel    string   `help:"Log level." enum:"trace,debug,info,warn,error" default:"info"`
	} `cmd:"" help:"Generate schema and implementation for the connector from functions."`

	Version struct{} `cmd:"" help:"Print the CLI version."`
}

func main() {
	cmd := kong.Parse(&cli)
	switch cmd.Command() {
	case "init":
		setupGlobalLogger(cli.Init.LogLevel)
		log.Info().
			Str("name", cli.Init.Name).
			Str("module", cli.Init.Module).
			Str("output", cli.Init.Output).
			Msg("generating the NDC boilerplate...")
		if err := generateNewProject(cli.Init.Name, cli.Init.Module, cli.Init.Output); err != nil {
			log.Fatal().Err(err).Msg("failed to generate new project")
		}
		log.Info().Msg("generated successfully")
	case "generate":
		setupGlobalLogger(cli.Generate.LogLevel)
		log.Info().
			Str("path", cli.Generate.Path).
			Interface("directories", cli.Generate.Directories).
			Msg("generating connector schema...")
		moduleName, err := getModuleName(cli.Generate.Path)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get module name. The base path must contain a go.mod file")
		}

		if err = parseAndGenerateConnector(cli.Generate.Path, cli.Generate.Directories, moduleName); err != nil {
			log.Fatal().Err(err).Msg("failed to generate connector schema")
		}
		if err := execGoFormat("."); err != nil {
			log.Fatal().Err(err).Msg("failed to format code")
		}
		log.Info().Msg("generated successfully")
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
