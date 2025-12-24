# Quickstart

## Installation

```
brew install go
```

## Testing

```
go test ./...
```

## Fuzzing

```
go test -run=^$ -fuzz=^FuzzNewFieldElement$ -fuzztime=10s ./pkg/fields
```
