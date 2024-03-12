module github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk

go 1.19

require (
	github.com/alecthomas/kong v0.8.1
	github.com/fatih/structtag v1.2.0
	github.com/hasura/ndc-sdk-go v0.1.3
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/mod v0.15.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hasura/ndc-sdk-go => ../../
