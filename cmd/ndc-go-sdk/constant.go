package main

const (
	connectorOutputFile   = "connector.generated.go"
	schemaOutputFile      = "schema.generated.json"
	typeMethodsOutputFile = "types.generated.go"
	googleUuidPackageName = "github.com/google/uuid"
	sdkScalarPackageName  = "github.com/hasura/ndc-sdk-go/scalar"
)

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
