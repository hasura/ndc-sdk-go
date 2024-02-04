.PHONY: typegen
typegen:
	cd typegen && ./regenerate-schema.sh

.PHONY: typegen
format:
	gofmt -w -s .


.PHONY: lint
lint:
	@(./scripts/check_installed.sh golangci-lint "golangci-lint: https://golangci-lint.run/usage/install/" && \
	golangci-lint run )