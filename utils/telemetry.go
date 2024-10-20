package utils

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
)

// JSONAttribute creates a OpenTelemetry attribute with JSON data
func JSONAttribute(key string, data any) attribute.KeyValue {
	switch d := data.(type) {
	case json.RawMessage:
		return attribute.String(key, string(d))
	default:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return attribute.String(key, err.Error())
		}
		return attribute.String(key, string(jsonBytes))
	}
}

// DebugJSONAttributes create OpenTelemetry attributes with JSON data.
// They are only visible on debug mode
func DebugJSONAttributes(data map[string]any, isDebug bool) []attribute.KeyValue {
	if !isDebug || len(data) == 0 {
		return []attribute.KeyValue{}
	}

	attrs := []attribute.KeyValue{}
	for k, v := range data {
		attrs = append(attrs, JSONAttribute(k, v))
	}
	return attrs
}
