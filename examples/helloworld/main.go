package main

import "github.com/hasura/ndc-sdk-go/connector"

func main() {
	if err := connector.Start[mockRawConfiguration, mockConfiguration, mockState](&mockConnector{}); err != nil {
		panic(err)
	}
}
