package command

import (
	"os"
	"time"

	"github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/command/internal"
	"github.com/rs/zerolog/log"
)

// UpdateArguments represent input arguments of the `update` command
type UpdateArguments internal.ConnectorGenerationArguments

// UpdateConnectorSchema updates connector schema
func UpdateConnectorSchema(args UpdateArguments, start time.Time) {
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

	if err = internal.ParseAndGenerateConnector(internal.ConnectorGenerationArguments(args), moduleName); err != nil {
		log.Fatal().Err(err).Msg("failed to generate connector schema")
	}
	if err := execGoFormat("."); err != nil {
		log.Fatal().Err(err).Msg("failed to format code")
	}
	log.Info().Str("exec_time", time.Since(start).Round(time.Millisecond).String()).
		Msg("generated successfully")
}
