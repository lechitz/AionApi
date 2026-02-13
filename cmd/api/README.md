# API Entrypoint (`cmd/api`)

**Path:** `cmd/api`

## Overview

This is the production application entrypoint.
It bootstraps the platform, loads dependencies, and starts HTTP/GraphQL servers.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| App bootstrap | Wire modules/dependencies via platform setup |
| Runtime startup | Initialize server, observability, config |
| Graceful shutdown | Coordinate lifecycle termination |

## Build and Run

```bash
go build -o bin/api ./cmd/api
make dev
# or
go run ./cmd/api
```

## Design Notes

- Keep this package orchestration-only.
- Business rules belong to bounded contexts in `internal/<ctx>`.
- Developer utilities should stay in `hack/`.

## Package Improvements

- Add startup sequence diagram (config -> providers -> server start).
- Add a short troubleshooting section for missing env/dependency failures.
- Add explicit command examples for local and containerized runs.
- Add pointers to runtime health/metrics endpoints.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
