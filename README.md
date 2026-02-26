# relint

`relint` is a custom Go multichecker that runs this repository's analyzers.

## Requirements

- Go 1.26+

## Build

```bash
go build -o relint .
```

## Run

```bash
go run . ./...
```

Or run the built binary:

```bash
./relint ./...
```

## Test

```bash
go test ./...
```
