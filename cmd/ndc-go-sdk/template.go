package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
)

//go:embed all:* templates/new
var initTemplateFS embed.FS

//go:embed templates/connector/connector.go.tmpl
var connectorTemplateStr string
var connectorTemplate *template.Template

const (
	templateNewPath = "templates/new"
)

func init() {
	var err error
	connectorTemplate, err = template.New(connectorOutputFile).Parse(connectorTemplateStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse connector template: %s", err))
	}
}

func generateNewProject(name string, moduleName string, srcPath string, silent bool) error {
	if srcPath == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		srcPath = path.Join(p, name)
	}
	if err := os.MkdirAll(srcPath, 0755); err != nil {
		return err
	}

	if err := generateNewProjectFiles(name, moduleName, srcPath); err != nil {
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

	if err := execGoGetUpdate(".", "github.com/hasura/ndc-sdk-go"); err != nil {
		if silent {
			return nil
		}
		return err
	}

	log.Info().Msg("generating connector functions...")
	if err := parseAndGenerateConnector(srcPath, []string{"functions", "types"}, moduleName); err != nil {
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

func generateNewProjectFiles(name string, moduleName string, srcPath string) error {

	return fs.WalkDir(initTemplateFS, templateNewPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		log.Debug().Msgf("%s", filePath)
		if filePath == templateNewPath {
			return nil
		}

		targetPath := path.Join(srcPath, strings.TrimPrefix(filePath, fmt.Sprintf("%s/", templateNewPath)))
		if d.IsDir() {
			return os.Mkdir(targetPath, 0755)
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
			"Name":   name,
			"Module": moduleName,
		})
		if err != nil {
			return err
		}
		return w.Flush()
	})
}

func execGoModTidy(basePath string) error {
	return execCommand(basePath, "go", "mod", "tidy")
}

func execGoFormat(basePath string) error {
	return execCommand(basePath, "gofmt", "-w", "-s", ".")
}

func execGoGetUpdate(basePath string, moduleName string) error {
	return execCommand(basePath, "go", "get", "-u", moduleName)
}

func execCommand(basePath string, commandName string, args ...string) error {
	cmd := exec.Command(commandName, args...)
	if basePath != "" {
		cmd.Dir = basePath
	}
	out, err := cmd.Output()
	l := log.Debug()
	if err != nil {
		l = log.Error().Err(err)
	}
	l.Strs("args", args).Str("result", string(out)).Msg(commandName)
	return err
}

func getModuleName(basePath string) (string, error) {
	goModBytes, err := os.ReadFile(path.Join(basePath, "go.mod"))
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(goModBytes), nil
}
