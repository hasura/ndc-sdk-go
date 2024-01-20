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
