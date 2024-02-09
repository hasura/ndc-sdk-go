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