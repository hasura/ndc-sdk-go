module github.com/hasura/ndc-codegen-subdir-test

go 1.21

require (
	github.com/hasura/ndc-sdk-go v1.4.1
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	golang.org/x/sync v0.8.0
	github.com/hasura/ndc-codegen-example v1.2.5
)

replace github.com/hasura/ndc-sdk-go => ../../../../../../../

replace github.com/hasura/ndc-codegen-example => ../../../../../../../example/codegen
