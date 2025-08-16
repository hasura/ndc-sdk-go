#!/bin/bash

set -eo pipefail

trap 'printf "\nkilling process..." && kill $serverPID' EXIT

cd "$(dirname "${BASH_SOURCE[0]}")"
mkdir -p ./dist

NDC_TEST_VERSION=v0.2.10
NDC_TEST_PATH=./dist/ndc-test

if [ ! -f ./tmp/ndc-test ]; then
  if [ "$(uname -m)" == "arm64" ]; then
    curl -L https://github.com/hasura/ndc-spec/releases/download/$NDC_TEST_VERSION/ndc-test-aarch64-apple-darwin -o $NDC_TEST_PATH
  elif [ $(uname) == "Darwin" ]; then
    curl -L https://github.com/hasura/ndc-spec/releases/download/$NDC_TEST_VERSION/ndc-test-x86_64-apple-darwin -o $NDC_TEST_PATH
  else
    curl -L https://github.com/hasura/ndc-spec/releases/download/$NDC_TEST_VERSION/ndc-test-x86_64-unknown-linux-gnu -o $NDC_TEST_PATH
  fi
  chmod +x $NDC_TEST_PATH
fi

http_wait() {
  printf "$1:\t "
  for i in {1..120};
  do
    local code="$(curl -s -o /dev/null -m 2 -w '%{http_code}' $1)"
    if [ "$code" != "200" ]; then
      printf "."
      sleep 1
    else
      printf "\r\033[K$1:\t ${GREEN}OK${NC}\n"
      return 0
    fi
  done
  printf "\n${RED}ERROR${NC}: cannot connect to $1.\n"
  exit 1
}

go build -o ./dist/ndc-codegen-example .
./dist/ndc-codegen-example serve > /dev/null 2>&1 &
serverPID=$!

http_wait http://localhost:8080/health

$NDC_TEST_PATH test --endpoint http://localhost:8080
