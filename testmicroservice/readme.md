# Port Loader

A small Go service that streams a large `ports.json` file and upserts port records into an in-memory store. Designed for low memory usage and graceful shutdown.

## Requirements

- Go 1.22+
- Docker (optional)

## Project layout

```
cmd/portloader/          application entrypoint
internal/domain/port/    domain model and service
internal/parser/         streaming JSON reader
internal/storage/memory/ in-memory repository
```

## Run locally

```bash
make build
make run
```

Or directly:

```bash
go run ./cmd/portloader -file ports.json
```

Expected output:

```
level=INFO msg="loading ports" file=ports.json
level=INFO msg="ports loaded" processed=6 stored=5
```

The service keeps running until it receives `SIGTERM` or `SIGINT`, then shuts down gracefully.

Interrupt during load:

```bash
go run ./cmd/portloader -file ports.json
# press Ctrl+C while loading a large file
```

## Test

```bash
make test
```

Test coverage includes:

- unit tests for domain service validation
- integration tests for in-memory storage
- table-driven parser tests
- fixture-based test against `ports.json`

## Lint

```bash
make lint
```

## Docker

The image copies a pre-built binary instead of compiling inside Docker:

```bash
make docker-build
make docker-run
```

## Input format

`ports.json` is a JSON array. Each item uses UN/LOCODE as `id`:

```json
[
  {
    "id": "NLRTM",
    "name": "Rotterdam",
    "city": "Rotterdam",
    "country": "Netherlands",
    "coordinates": [4.4792, 51.9225]
  }
]
```

If the same `id` appears multiple times, the last record wins.

## Design notes

- JSON is parsed with `json.Decoder`, one record at a time.
- Business rules live in `internal/domain/port`.
- Storage is behind a repository interface for testability.
- `SIGKILL` cannot be handled; only `SIGTERM` and `SIGINT` trigger graceful shutdown.
