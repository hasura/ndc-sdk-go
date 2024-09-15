package types

import "github.com/hasura/ndc-sdk-go/connector"

// Configuration contains required settings for the connector.
type Configuration struct{}

// State is the global state which is shared for every connector request.
type State struct {
	*connector.TelemetryState
}
