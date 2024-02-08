package main

import (
	"encoding/json"
	"flag"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"strings"
)

func main() {
	filePath := flag.String("path", "", "Source file path")
	flag.Parse()
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, *filePath, func(fi fs.FileInfo) bool {
		return !fi.IsDir() && !strings.Contains(fi.Name(), "generated")
	}, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sm, err := parseRawConnectorSchema(pkgs)
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(sm.Schema(), "", "  ")
	if err != nil {
		panic(err)
	}

	log.Println(string(bytes))
}
