VERSION ?= $(shell date +"%Y%m%d")
OUTPUT_DIR := _output
ROOT_DIR := $(shell pwd)

.PHONY: typegen
typegen:
	cd typegen && ./regenerate-schema.sh

.PHONY: format
format:
	gofmt -w -s .


.PHONY: test-sdk
test-sdk:
	go test -v -race -timeout 3m ./...

.PHONY: test-hasura-ndc-go
test-hasura-ndc-go:
	cd cmd/hasura-ndc-go && \
		go test -v -race ./...

.PHONY: test-example-codegen
test-example-codegen:
	cd example/codegen && \
		go test -v -race ./...

.PHONY: test
test: test-sdk test-hasura-ndc-go test-example-codegen

.PHONY: go-tidy
go-tidy:
	go mod tidy
	cd $(ROOT_DIR)/cmd/hasura-ndc-go && go mod tidy
	cd $(ROOT_DIR)/cmd/hasura-ndc-go/command/internal/testdata/basic/source && go mod tidy
	cd $(ROOT_DIR)/cmd/hasura-ndc-go/command/internal/testdata/empty/source && go mod tidy
	cd $(ROOT_DIR)/example/codegen && go mod tidy
	
# Install golangci-lint tool to run lint locally
# https://golangci-lint.run/usage/install
.PHONY: lint
lint:
	golangci-lint run --fix
	cd cmd/hasura-ndc-go && golangci-lint run --fix
	
# clean the output directory and generated tests snapshots
.PHONY: clean
clean:
	rm -rf "$(OUTPUT_DIR)"
	rm -f cmd/hasura-ndc-go/command/internal/testdata/*/source/connector.generated.go
	rm -f cmd/hasura-ndc-go/command/internal/testdata/*/source/**/connector.generated.go
	rm -f cmd/hasura-ndc-go/command/internal/testdata/*/source/schema.generated.json
	rm -f cmd/hasura-ndc-go/command/internal/testdata/*/source/**/schema.generated.json
	rm -f cmd/hasura-ndc-go/command/internal/testdata/*/source/**/types.generated.go
	rm -rf cmd/hasura-ndc-go/testdata/**/testdata

.PHONY: build-codegen
build-codegen:
	cd ./cmd/hasura-ndc-go && go build -o ../../_output/hasura-ndc-go .
	
# build the build-codegen cli for all given platform/arch
.PHONY: ci-build-codegen
ci-build-codegen: export CGO_ENABLED=0
ci-build-codegen: clean
	cd ./cmd/hasura-ndc-go && \
	go get github.com/mitchellh/gox && \
	go run github.com/mitchellh/gox -ldflags '-X github.com/hasura/ndc-sdk-go/cmd/hasura-ndc-go/version.BuildVersion=$(VERSION) -s -w -extldflags "-static"' \
		-osarch="linux/amd64 linux/arm64 darwin/amd64 windows/amd64 darwin/arm64" \
		-output="../../$(OUTPUT_DIR)/hasura-ndc-go-{{.OS}}-{{.Arch}}" \
		.