package main

import "github.com/hasura/ndc-sdk-go/connector"

func main() {
	if err := connector.Start[Configuration, State](
		&Connector{},
		connector.WithMetricsPrefix("ndc_ref"),
		connector.WithDefaultServiceName("ndc_ref"),
	); err != nil {
		panic(err)
	}
}
