package main

import "github.com/hasura/ndc-sdk-go/connector"

func main() {
	if err := connector.Start[RawConfiguration, Configuration, State](
		&Connector{},
		connector.WithMetricsPrefix("ndc_codegen"),
		connector.WithDefaultServiceName("ndc_codegen"),
		connector.WithoutConfig(),
	); err != nil {
		panic(err)
	}
}
