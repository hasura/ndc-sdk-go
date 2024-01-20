package main

import "github.com/hasura/ndc-sdk-go/connector"

func main() {
	if err := connector.Start[RawConfiguration, Configuration, State](&Connector{}); err != nil {
		panic(err)
	}
}
