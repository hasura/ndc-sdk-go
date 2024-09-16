# codegen Connector

## Get started

Start the connector server at http://localhost:8080

```go
go run . serve
```

See [NDC Go SDK](https://github.com/hasura/ndc-sdk-go) for more information and [the generation tool](https://github.com/hasura/ndc-sdk-go/tree/main/cmd/hasura-ndc-go) for command documentation.

## Development

Generate codes with the following command:

```sh
go run ../../cmd/hasura-ndc-go update --schema-format go
```
