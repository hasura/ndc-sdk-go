package main

import (
	"github.com/hasura/ndc-sdk-go/connector"
	"github.com/hasura/ndc-sdk-go/example/codegen/types"
)

func main() {
	if err := connector.Start[types.RawConfiguration, types.Configuration, types.State](
		&Connector{},
		connector.WithMetricsPrefix("ndc_codegen"),
		connector.WithDefaultServiceName("ndc_codegen"),
		connector.WithoutConfig(),
	); err != nil {
		panic(err)
	}
}
