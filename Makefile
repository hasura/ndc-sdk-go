VERSION ?= $(shell ./scripts/get-CalVer.sh)
PLUGINS_BRANCH ?= master
OUTPUT_DIR := _output

.PHONY: typegen
typegen:
	cd typegen && ./regenerate-schema.sh

.PHONY: typegen
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
	go build -o _output/ndc-go-sdk ./cmd/ndc-go-sdk
	
# build the build-codegen cli for all given platform/arch
.PHONY: build-codegen
ci-build-codegen: export CGO_ENABLED=0
ci-build-codegen: clean
	go run github.com/mitchellh/gox -ldflags '-X github.com/hasura/ndc-sdk-go/cmd/ndc-go-sdk/version.BuildVersion=$(VERSION) -s -w -extldflags "-static"' \
	-osarch="linux/amd64 darwin/amd64 windows/amd64 darwin/arm64" \
	-output="$(OUTPUT_DIR)/$(VERSION)/ndc-go-sdk-{{.OS}}-{{.Arch}}" \
	./cmd/ndc-go-sdk