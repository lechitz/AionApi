# Application Entrypoints (`cmd`)

**Path:** `cmd`

## Purpose

`cmd` owns long-lived process entrypoints only.
The current repo ships two binaries:

| Folder | Binary | Role |
| --- | --- | --- |
| `api/` | `aion-api` | HTTP, GraphQL, health, and runtime bootstrap |
| `outbox-publisher/` | `aion-api-outbox-publisher` | background publisher for durable outbox rows |

## Current Flow

- each command keeps `main.go` minimal and delegates to `runWithDeps`
- bootstrap timeout parsing stays local to the entrypoint
- the Fx composition root lives under `internal/platform/fxapp`
- durable config, DB wiring, HTTP server, and Kafka adapters stay outside `cmd`

## Boundaries

- do not place domain, repository, or transport logic in `cmd`
- add a new subfolder here only when a new binary/process exists
- dev and lab tools belong under `hack/`, not beside runtime entrypoints

## Validate

```bash
go run ./cmd/api
go run ./cmd/outbox-publisher
make dev
make logs-api-publisher
```

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
