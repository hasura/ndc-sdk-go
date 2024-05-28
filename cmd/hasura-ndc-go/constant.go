package main

import "regexp"

const (
	connectorOutputFile   = "connector.generated.go"
	schemaOutputFile      = "schema.generated.json"
	typeMethodsOutputFile = "types.generated.go"
)

type nativeScalarPackageConfig struct {
	PackageName string
	Pattern     *regexp.Regexp
}

var fieldNameRegex = regexp.MustCompile(`[^\w]`)

var nativeScalarPackages map[string]nativeScalarPackageConfig = map[string]nativeScalarPackageConfig{
	"scalar": {
		PackageName: "github.com/hasura/ndc-sdk-go/scalar",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*github\.com\/hasura\/ndc-sdk-go\/scalar\.`),
	},
	"json": {
		PackageName: "encoding/json",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*encoding\/json\.`),
	},
	"uuid": {
		PackageName: "github.com/google/uuid",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*github\.com\/google\/uuid\.`),
	},
	"time": {
		PackageName: "time",
		Pattern:     regexp.MustCompile(`^(\[\]|\*)*time\.`),
	},
}

const textBlockErrorCheck = `
    if err != nil {
		  return err
    }
`

const textBlockErrorCheck2 = `
    if err != nil {
      return nil, err
    }
`
