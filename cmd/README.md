# Application Entrypoints (`cmd`)

**Path:** `cmd`

## Overview

This layer contains executable application entrypoints.
Current canonical runtime entrypoint is `cmd/api`.

## Structure

| Folder | Role |
| --- | --- |
| `api/` | Main API bootstrap and runtime startup |

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Bootstrap | Initialize platform dependencies and wiring |
| Runtime | Start HTTP/GraphQL servers and middlewares |
| Lifecycle | Handle graceful shutdown and cleanup |

## Build and Run

```bash
go build -o bin/api ./cmd/api
go run ./cmd/api
# local stack with dependencies
make dev
```

## Design Notes

- Keep entrypoints orchestration-only.
- Domain rules belong to bounded contexts under `internal/<ctx>/core`.
- Dev-only scripts/tools should live under `hack/`.

## Package Improvements

- Add startup sequence diagram (config -> deps -> server).
- Add troubleshooting for missing env/infra dependencies.
- Add explicit local vs container run matrix.
- Add link to health/metrics endpoints for runtime checks.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
