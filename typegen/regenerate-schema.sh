#!/usr/bin/env bash
set -e -o pipefail

if ! command -v jq 2>&1 >/dev/null; then
  echo "please install jq"
  exit 1
fi

if ! command -v go-jsonschema > /dev/null; then
  go get github.com/atombender/go-jsonschema/...
  go install github.com/atombender/go-jsonschema@latest
  # cleanup unused json-patch packages
  go mod tidy
fi

if ! command -v json-patch 2>&1 >/dev/null; then
  go get github.com/evanphx/json-patch/cmd/json-patch
  go install github.com/evanphx/json-patch/cmd/json-patch
  # cleanup unused json-patch packages
  go mod tidy
fi

# download the schema json file from ndc-sdk-typescript repository and regenerate schema 
if [ ! -f ./schema.generated.json ]; then
  curl -L https://raw.githubusercontent.com/hasura/ndc-sdk-typescript/refs/tags/v8.0.0/src/schema/schema.generated.json -o schema.generated.json
fi

cat schema.generated.json | json-patch -p schema.patch.json > schema.patched.json

GENERATED_SCHEMA_GO="../schema/schema.generated.go"
SED_CMD="sed -i"
if [ "$(uname)" == "Darwin" ]; then
  SED_CMD="sed -i .bak"
fi

go-jsonschema \
  --package=github.com/hasura/ndc-sdk-go/schema \
  --output=$GENERATED_SCHEMA_GO \
  schema.patched.json

# patch some custom types because of the limitation of the generation tool
$SED_CMD 's/type Field interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type Argument interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type RelationshipArgument interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type MutationOperation interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type MutationRequestOperationsElem interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/MutationRequestOperationsElem/MutationOperation/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/QueryRequestArguments map\[string\]interface{}/QueryRequestArguments map[string]Argument/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type RowSetRowsElem map\[string\]interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/RowSetRowsElem/map[string]any/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type MutationOperationResultsReturningElem map\[string\]interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/MutationOperationResultsReturningElem/map[string]any/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type Expression interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ComparisonTarget interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type BinaryComparisonOperator interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/Where interface{}/Where Expression/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type Aggregate interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type OrderByTarget interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/QueryAggregates map\[string\]interface{}/QueryAggregates map[string]Aggregate/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/RelationshipArguments map\[string\]interface{}/RelationshipArguments map[string]RelationshipArgument/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/Predicate interface{}/Predicate Expression/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type OrderByElementTarget interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/OrderByElementTarget/OrderByTarget/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/PathElementArguments map\[string\]interface{}/PathElementArguments map[string]RelationshipArgument/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type Type interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/ResultType interface{}/ResultType Type/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/Type interface{}/Type Type/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/QueryFields map\[string\]interface{}/QueryFields map[string]Field/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ComparisonValue interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ExistsInCollection interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ComparisonOperatorDefinition interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type NestedField interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ScalarTypeComparisonOperators map\[string\]interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/ScalarTypeComparisonOperators/map[string]ComparisonOperatorDefinition/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type MutationOperationResults interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type MutationResponseOperationResultsElem interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/MutationResponseOperationResultsElem/MutationOperationResults/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type TypeRepresentation interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/Representation interface{}/Representation TypeRepresentation/g' $GENERATED_SCHEMA_GO
$SED_CMD '/type NestedFieldCapabilities struct {/,/}/s/OrderBy \*OrderBy/OrderBy interface{}/g' $GENERATED_SCHEMA_GO
$SED_CMD '/type Query struct {/,/}/s/OrderBy \interface{}/OrderBy *OrderBy/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/NestedFields interface{}/NestedFields NestedFieldCapabilities/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/plain.NestedFields = map\[string\]interface{}{}/plain.NestedFields = NestedFieldCapabilities{}/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/plain.Exists = map\[string\]interface{}{}/plain.Exists = ExistsCapabilities{}/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupOrderByTarget interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type AggregateFunctionDefinition interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ArrayComparison interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type Dimension interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupComparisonTarget interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupComparisonValue interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupExpression interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type RowFieldValue interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/\[\]GroupingDimensionsElem/[]Dimension/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupingDimensionsElem interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type GroupingAggregates map\[string\]interface{}/type GroupingAggregates map[string]Aggregate/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/order_by,omitempty/order_by/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/variables,omitempty/variables/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/order_by_aggregate,omitempty/order_by_aggregate/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/relation_comparisons,omitempty/relation_comparisons/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ScalarTypeAggregateFunctions map\[string\]interface{}/type ScalarTypeAggregateFunctions map[string]AggregateFunctionDefinition/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ExtractionFunctionDefinition interface{}//g' $GENERATED_SCHEMA_GO
$SED_CMD 's/type ScalarTypeExtractionFunctions map\[string\]interface{}/type ScalarTypeExtractionFunctions map[string]ExtractionFunctionDefinition/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/plain.ExtractionFunctions = map\[string\]interface{}{}/plain.ExtractionFunctions = map[string]ExtractionFunctionDefinition{}/g' $GENERATED_SCHEMA_GO
$SED_CMD 's/^.*DeleteThis .*$//g' $GENERATED_SCHEMA_GO

rm -f "$GENERATED_SCHEMA_GO.bak"

# format codes
gofmt -w -s ../

cat schema.patched.json | jq . > ../schema/schema.generated.json