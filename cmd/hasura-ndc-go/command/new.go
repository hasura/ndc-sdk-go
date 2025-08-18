package command

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hasura/ndc-sdk-go/v2/cmd/hasura-ndc-go/command/internal"
	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
)

//go:embed all:internal/templates/new/* internal/templates/new
var initTemplateFS embed.FS

const (
	templateNewPath = "internal/templates/new"
)

// NewArguments input arguments for the new command.
type NewArguments struct {
	Name    string `help:"Name of the connector."                            required:"" short:"n"`
	Module  string `help:"Module name of the connector"                      required:"" short:"m"`
	Version string `help:"The version of ndc-sdk-go"`
	Output  string `help:"The location where source codes will be generated"             short:"o" default:""`
}

// GenerateNewProject generates a new project boilerplate.
func GenerateNewProject(args *NewArguments, silent bool) error {
	srcPath := args.Output
	if srcPath == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}

		srcPath = filepath.Join(p, args.Name)
	}

	if err := os.MkdirAll(srcPath, 0o755); err != nil {
		return err
	}

	if err := generateNewProjectFiles(args, srcPath); err != nil {
		return err
	}

	if err := os.Chdir(srcPath); err != nil {
		return err
	}

	log.Info().Msg("downloading dependencies...")

	if err := execGoModTidy(""); err != nil {
		if silent {
			return nil
		}

		return err
	}

	log.Info().Msg("generating connector functions...")

	if err := internal.ParseAndGenerateConnector(internal.ConnectorGenerationArguments{
		Path:        ".",
		Directories: []string{"functions", "types"},
	}, args.Module); err != nil {
		return err
	}

	// reverify dependencies after generated
	if err := execGoModTidy(""); err != nil {
		if silent {
			return nil
		}

		return err
	}

	err := execGoFormat(".")
	if err != nil && silent {
		return nil
	}

	return err
}

func generateNewProjectFiles(args *NewArguments, srcPath string) error {
	return fs.WalkDir(
		initTemplateFS,
		templateNewPath,
		func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			log.Debug().Msg(filePath)

			if filePath == templateNewPath {
				return nil
			}

			targetPath := filepath.Join(srcPath, strings.TrimPrefix(filePath, templateNewPath+"/"))
			if d.IsDir() {
				return os.Mkdir(targetPath, 0o755)
			}

			fileTemplate, err := template.ParseFS(initTemplateFS, filePath)
			if err != nil {
				return err
			}

			targetPath = strings.TrimSuffix(targetPath, ".tmpl")

			f, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			defer func() {
				_ = f.Close()
			}()

			w := bufio.NewWriter(f)

			err = fileTemplate.Execute(w, map[string]any{
				"Name":    args.Name,
				"Module":  args.Module,
				"Version": args.Version,
			})
			if err != nil {
				return err
			}

			return w.Flush()
		},
	)
}

func execGetLatestSDK(basePath string) error {
	return execCommand(basePath, "go", "get", "github.com/hasura/ndc-sdk-go/v2@v2")
}

func execGoModTidy(basePath string) error {
	return execCommand(basePath, "go", "mod", "tidy")
}

func execGoFormat(basePath string) error {
	return execCommand(basePath, "gofmt", "-w", "-s", ".")
}

func execCommand(basePath string, commandName string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), commandName, args...)
	if basePath != "" {
		cmd.Dir = basePath
	}

	l := log.With().Strs("args", args).Str("command", commandName).Logger()

	var errBuf bytes.Buffer

	cmd.Stderr = &errBuf

	out, err := cmd.Output()
	if err != nil {
		l.Error().Err(errors.New(errBuf.String())).Msg(err.Error())
	} else {
		l.Debug().Str("logs", errBuf.String()).Msg(string(out))
	}

	return err
}

func getModuleName(basePath string) (string, error) {
	goModBytes, err := os.ReadFile(path.Join(basePath, "go.mod"))
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(goModBytes), nil
}
