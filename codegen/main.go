package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
)

var cli struct {
	Init struct {
		Name    string `help:"Name of the connector." short:"n" required:""`
		Package string `help:"Package name of the connector" short:"p" required:""`
		Output  string `help:"The location where source codes will be generated" short:"o" default:""`
	} `cmd:"" help:"Initialize a NDC connector boilerplate."`

	Generate struct {
		Path        string   `help:"The base path of the connector's source code" short:"p" default:"."`
		Directories []string `help:"Folders contain NDC operation functions" short:"d" default:"functions,types"`
	} `cmd:"" help:"Generate schema and implementation for the connector from functions."`
}

func main() {
	cmd := kong.Parse(&cli)
	switch cmd.Command() {
	case "init":
		log.Println("init", cli.Init)
		if err := generateNewProject(cli.Init.Name, cli.Init.Package, cli.Init.Output); err != nil {
			panic(err)
		}
	case "generate":
		log.Println("generate", cli.Generate)
		sm, err := parseRawConnectorSchemaFromGoCode(cli.Generate.Path, cli.Generate.Directories)
		if err != nil {
			panic(err)
		}
		if err = generateConnector(sm, cli.Generate.Path); err != nil {
			panic(err)
		}

	default:
		panic(fmt.Errorf("unknown command <%s>", cmd.Command()))
	}
}
