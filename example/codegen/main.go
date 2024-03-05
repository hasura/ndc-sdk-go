package main

import (
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/connector"
)

func main() {
	if err := connector.Start[types.Configuration, types.State](
		&Connector{},
		connector.WithMetricsPrefix("codegen"),
		connector.WithDefaultServiceName("codegen"),
	); err != nil {
		panic(err)
	}
}
