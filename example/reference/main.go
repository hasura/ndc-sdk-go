package main

import "github.com/hasura/ndc-sdk-go/connector"

func main() {
	if err := connector.Start[RawConfiguration, Configuration, State](
		&Connector{},
		connector.WithMetricsPrefix("ndc_ref"),
		connector.WithDefaultServiceName("ndc_ref"),
		connector.WithoutConfig(),
	); err != nil {
		panic(err)
	}
}
