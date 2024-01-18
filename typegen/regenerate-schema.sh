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
