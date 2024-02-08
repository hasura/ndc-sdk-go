#!/usr/bin/env bash
set -eu -o pipefail

if [ ! -f ./go-jsonschema ]; then
  wget -qO- https://github.com/omissis/go-jsonschema/releases/download/v0.14.1/go-jsonschema_Linux_x86_64.tar.gz | tar xvz
  chmod +x go-jsonschema
fi

# download the schema json file from ndc-sdk-typescript repository and regerate schema 
wget https://raw.githubusercontent.com/hasura/ndc-sdk-typescript/main/src/schema/schema.generated.json

./go-jsonschema \
  --package=github.com/hasura/ndc-sdk-go/schema \
  --output=../schema/schema.generated.go \
  schema.generated.json

# finally move that json file to schema folder
mv schema.generated.json ../schema/schema.generated.json

# patch some custom types because of the limitation of the generation tool
sed -i 's/type Field interface{}//g' ../schema/schema.generated.go
sed -i 's/type Argument interface{}//g' ../schema/schema.generated.go
sed -i 's/type RelationshipArgument interface{}//g' ../schema/schema.generated.go
sed -i 's/type MutationOperation interface{}//g' ../schema/schema.generated.go
sed -i 's/type MutationRequestOperationsElem interface{}//g' ../schema/schema.generated.go
sed -i 's/MutationRequestOperationsElem/MutationOperation/g' ../schema/schema.generated.go
sed -i 's/QueryRequestArguments map\[string\]interface{}/QueryRequestArguments map[string]Argument/g' ../schema/schema.generated.go
sed -i 's/RowSetRowsElem map\[string\]interface{}/Row any/g' ../schema/schema.generated.go
sed -i 's/RowSetRowsElem/Row/g' ../schema/schema.generated.go
sed -i 's/type MutationOperationResultsReturningElem map\[string\]interface{}//g' ../schema/schema.generated.go
sed -i 's/MutationOperationResultsReturningElem/Row/g' ../schema/schema.generated.go
sed -i 's/Query interface{}/Query Query/g' ../schema/schema.generated.go
sed -i 's/OrderBy interface{}/OrderBy *OrderBy/g' ../schema/schema.generated.go
sed -i 's/type Expression interface{}//g' ../schema/schema.generated.go
sed -i 's/type ComparisonTarget interface{}//g' ../schema/schema.generated.go
sed -i 's/type BinaryComparisonOperator interface{}//g' ../schema/schema.generated.go
sed -i 's/Where interface{}/Where Expression/g' ../schema/schema.generated.go
sed -i 's/type Aggregate interface{}//g' ../schema/schema.generated.go
sed -i 's/type OrderByTarget interface{}//g' ../schema/schema.generated.go
sed -i 's/QueryAggregates map\[string\]interface{}/QueryAggregates map[string]Aggregate/g' ../schema/schema.generated.go
sed -i 's/Predicate interface{}/Predicate Expression/g' ../schema/schema.generated.go
sed -i 's/type OrderByElementTarget interface{}//g' ../schema/schema.generated.go
sed -i 's/OrderByElementTarget/OrderByTarget/g' ../schema/schema.generated.go
sed -i 's/PathElementArguments map\[string\]interface{}/PathElementArguments map[string]RelationshipArgument/g' ../schema/schema.generated.go
sed -i 's/type Type interface{}//g' ../schema/schema.generated.go
sed -i 's/ResultType interface{}/ResultType Type/g' ../schema/schema.generated.go
sed -i 's/Type interface{}/Type Type/g' ../schema/schema.generated.go
sed -i 's/QueryFields map\[string\]interface{}/QueryFields map[string]Field/g' ../schema/schema.generated.go
sed -i 's/type ComparisonValue interface{}//g' ../schema/schema.generated.go
sed -i 's/type ExistsInCollection interface{}//g' ../schema/schema.generated.go
sed -i 's/type ComparisonOperatorDefinition interface{}//g' ../schema/schema.generated.go
sed -i 's/type NestedField interface{}//g' ../schema/schema.generated.go
sed -i 's/type ScalarTypeComparisonOperators map\[string\]interface{}//g' ../schema/schema.generated.go
sed -i 's/ScalarTypeComparisonOperators/map[string]ComparisonOperatorDefinition/g' ../schema/schema.generated.go

# format codes
gofmt -w -s ../