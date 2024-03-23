VERSION ?= $(shell date +"%Y%m%d")
OUTPUT_DIR := _output

.PHONY: typegen
typegen:
	cd typegen && ./regenerate-schema.sh

.PHONY: format
format:
	gofmt -w -s .

.PHONY: test
test:
	go test -v -race -timeout 3m ./...

# Install golangci-lint tool to run lint locally
# https://golangci-lint.run/usage/install
.PHONY: lint
lint:
	golangci-lint run

# clean the output directory
.PHONY: clean
clean:
	rm -rf "$(OUTPUT_DIR)"

.PHONY: build-codegen
build-codegen:
	cd ./cmd/ndc-go-sdk && go build -o ../../_output/hasura-ndc-go .
	
# build the build-codegen cli for all given platform/arch
.PHONY: ci-build-codegen
ci-build-codegen: export CGO_ENABLED=0
ci-build-codegen: clean
	cd ./cmd/ndc-go-sdk && \
	go get github.com/mitchellh/gox && \
	go run github.com/mitchellh/gox -ldflags '-X github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/version.BuildVersion=$(VERSION) -s -w -extldflags "-static"' \
		-osarch="linux/amd64 linux/arm64 darwin/amd64 windows/amd64 darwin/arm64" \
		-output="../../$(OUTPUT_DIR)/hasura-ndc-go-{{.OS}}-{{.Arch}}" \
		.