package command

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog/log"
)

// UpgradeArguments represent input arguments of the upgrade command.
type UpgradeArguments struct {
	Path         string `default:"." env:"HASURA_PLUGIN_CONNECTOR_CONTEXT_PATH" help:"The path of the root directory where the go.mod file is present" short:"p"`
	ConnectorDir string `default:"."                                            help:"The directory where the connector.go file is placed"`
}

type upgradeConnectorCommand struct {
	BasePath string
}

// UpgradeConnector upgrades the connector.
func UpgradeConnector(args UpgradeArguments) error {
	log.Info().
		Str("path", args.Path).
		Str("connector_dir", args.ConnectorDir).
		Msg("upgrade the connector")

	srcPath := filepath.Join(args.Path, args.ConnectorDir)

	if srcPath == "" || !filepath.IsAbs(srcPath) {
		p, err := os.Getwd()
		if err != nil {
			return err
		}

		srcPath = filepath.Join(p, srcPath)
	}

	ucc := upgradeConnectorCommand{
		BasePath: srcPath,
	}

	return ucc.patchConnectorFile()
}

func (ucc upgradeConnectorCommand) patchConnectorFile() error {
	connectorFilePath := filepath.Join(ucc.BasePath, "connector.go")

	originalContent, err := os.ReadFile(connectorFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("can not find the connector.go file")
		}

		return err
	}

	contentStr, isChanged := ucc.patchConnectorContent(originalContent)
	if !isChanged {
		return nil
	}

	return os.WriteFile(connectorFilePath, []byte(contentStr), 0o664)
}

func (ucc upgradeConnectorCommand) patchConnectorContent(originalContent []byte) (string, bool) {
	var isChanged bool

	contentStr := string(originalContent)
	versionRegexp := regexp.MustCompile(`Version:[\s\t]+"0\.1\.\d"`)
	varCapsRegexp := regexp.MustCompile(`Variables:[\s\t]+schema.LeafCapability\{\}`)

	if versionRegexp.MatchString(contentStr) {
		isChanged = true
		contentStr = versionRegexp.ReplaceAllString(contentStr, "Version: schema.NDCVersion")
	}

	if varCapsRegexp.MatchString(contentStr) {
		isChanged = true
		contentStr = varCapsRegexp.ReplaceAllString(
			contentStr,
			"Variables:    &schema.LeafCapability{}",
		)
	}

	return contentStr, isChanged
}
