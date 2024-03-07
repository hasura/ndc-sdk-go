package main

import (
	"github.com/hasura/ndc-codegen-example/types"
	"github.com/hasura/ndc-sdk-go/connector"
)

// Start the connector server at http://localhost:8080
//
//	go run . serve
//
// See [NDC Go SDK] for more information.
//
// [NDC Go SDK]: https://github.com/hasura/ndc-sdk-go
func main() {
	if err := connector.Start[types.Configuration, types.State](
		&Connector{},
		connector.WithMetricsPrefix("codegen"),
		connector.WithDefaultServiceName("codegen"),
	); err != nil {
		panic(err)
	}
}
