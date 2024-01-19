package schema

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/swaggest/jsonschema-go"
)

//go:embed schema.generated.json
var rawNdcSchema string

var _ndcSchema *jsonschema.Schema

func getNdcSchema() *jsonschema.Schema {
	if _ndcSchema != nil {
		return _ndcSchema
	}

	var inputSchema jsonschema.Schema
	if err := json.Unmarshal([]byte(rawNdcSchema), &inputSchema); err != nil {
		panic(fmt.Errorf("failed to decode NDC json schema: %s", err))
	}

	_ndcSchema = &inputSchema

	return _ndcSchema
}

func schemaForType(typeName string) *jsonschema.Schema {
	ref := fmt.Sprintf("#/definitions/%s", typeName)
	ndcSchema := getNdcSchema()
	return &jsonschema.Schema{
		Schema:      ndcSchema.Schema,
		Ref:         &ref,
		Definitions: ndcSchema.Definitions,
	}
}

var (
	CapabilitiesResponseSchema = schemaForType("CapabilitiesResponse")
	SchemaResponseSchema       = schemaForType("SchemaResponse")
	QueryRequestSchema         = schemaForType("QueryRequest")
	QueryResponseSchema        = schemaForType("QueryResponse")
	ExplainResponseSchema      = schemaForType("ExplainResponse")
	MutationRequestSchema      = schemaForType("MutationRequest")
	MutationResponseSchema     = schemaForType("MutationResponse")
	ErrorResponseSchema        = schemaForType("ErrorResponse")
	ValidateResponseSchema     = schemaForType("ValidateResponse")
)

func getStringValueByKey(collection map[string]any, key string) string {
	if collection == nil {
		return ""
	}

	anyValue, ok := collection[key]
	if !ok || anyValue == nil {
		return ""
	}

	if arg, ok := anyValue.(string); ok {
		return arg
	}

	return ""
}

// ToPtr converts a value to its pointer
func ToPtr[V any](value V) *V {
	return &value
}
