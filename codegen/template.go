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

//go:embed templates/new
var initTemplateFS embed.FS

const (
	templateNewPath = "templates/new"
)

func generateNewProject(name string, moduleName string, srcPath string) error {
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

	if err := parseAndGenerateConnector(srcPath, []string{"functions", "types"}, moduleName); err != nil {
		return err
	}

	return execGoModTidy(srcPath)
}

func generateNewProjectFiles(name string, moduleName string, srcPath string) error {

	return fs.WalkDir(initTemplateFS, templateNewPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

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
	cmd := exec.Command("go", "mod", "tidy")
	if basePath != "" {
		cmd.Dir = basePath
	}
	out, err := cmd.Output()
	log.Info().Str("result", string(out)).Msg("go mod tidy")
	return err
}

func getModuleName(basePath string) (string, error) {
	goModBytes, err := os.ReadFile(path.Join(basePath, "go.mod"))
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(goModBytes), nil
}
