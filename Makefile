.PHONY: typegen
typegen:
	cd typegen && ./regenerate-schema.sh

.PHONY: typegen
format:
	gofmt -w -s .