#!/bin/bash

set -eo pipefail

trap 'printf "\nkilling process..." && kill $serverPID' EXIT

cd "$(dirname "${BASH_SOURCE[0]}")"
mkdir -p ./dist

if [ ! -f ./dist/ndc-test ]; then
  curl -L https://github.com/hasura/ndc-spec/releases/download/v0.1.6/ndc-test-x86_64-unknown-linux-gnu -o ./dist/ndc-test
  chmod +x ./dist/ndc-test
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

./dist/ndc-test test --endpoint http://localhost:8080
