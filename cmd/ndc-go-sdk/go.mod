module github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk

go 1.19

require (
	github.com/alecthomas/kong v0.8.1
	github.com/fatih/structtag v1.2.0
	github.com/hasura/ndc-sdk-go v0.0.0-20240218161048-ce1f9dfc50bc
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/mod v0.15.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/swaggest/jsonschema-go v0.3.64 // indirect
	github.com/swaggest/refl v1.3.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hasura/ndc-sdk-go => ../../
