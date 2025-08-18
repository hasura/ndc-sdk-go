package command

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var packageV2UpgradeRegex = regexp.MustCompile(`github.com/hasura/ndc-sdk-go/[^v2]`)

const (
	sdkPackageV1 = "github.com/hasura/ndc-sdk-go"
	sdkPackageV2 = "github.com/hasura/ndc-sdk-go/v2"
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

	// if the github.com/hasura/ndc-sdk-go/v2 package exists in go.mod,
	// skip the migration.
	isChanged, err := ucc.patchGoMod()
	if !isChanged || err != nil {
		return err
	}

	return ucc.patchConnectorFile()
}

func (ucc upgradeConnectorCommand) patchGoMod() (bool, error) {
	goModFilePath := filepath.Join(ucc.BasePath, "go.mod")

	goModContent, err := os.ReadFile(goModFilePath)
	if err != nil {
		return false, err
	}

	contentStr := string(goModContent)

	if strings.Contains(contentStr, sdkPackageV2) {
		return false, nil
	}

	contentStr = strings.ReplaceAll(contentStr, sdkPackageV1, sdkPackageV2)

	return true, os.WriteFile(goModFilePath, []byte(contentStr), 0o664)
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

	if err := os.WriteFile(connectorFilePath, []byte(contentStr), 0o664); err != nil {
		return err
	}

	return ucc.patchImportSdkV2Files()
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

func (ucc upgradeConnectorCommand) patchImportSdkV2Content(originalContent []byte) (string, bool) {
	var isChanged bool

	contentStr := string(originalContent)

	if packageV2UpgradeRegex.MatchString(contentStr) {
		isChanged = true
		contentStr = strings.ReplaceAll(contentStr, sdkPackageV1, sdkPackageV2)
	}

	return contentStr, isChanged
}

func (ucc upgradeConnectorCommand) patchImportSdkV2Files() error {
	return filepath.WalkDir(ucc.BasePath, func(fileName string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		log.Debug().Msg(fileName)

		if !strings.HasSuffix(fileName, ".go") {
			return nil
		}

		filePath := filepath.Join(ucc.BasePath, fileName)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		newContent, isChanged := ucc.patchImportSdkV2Content(fileContent)
		if !isChanged {
			return nil
		}

		return os.WriteFile(filePath, []byte(newContent), 0o664)
	})
}
