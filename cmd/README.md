# cmd (Entrypoints)

This folder contains the binary entrypoints (main packages) for AionAPI. Each subfolder builds a standalone CLI or server binary. The code here should orchestrate, not implement domain logic.

## Package Composition

- `aion-api/`
  - Main API server entrypoint (HTTP + GraphQL).
  - Wires modules with Fx and configures Swagger metadata.
- `api-seed-caller/`
  - CLI that seeds via API endpoints (auth + GraphQL).
- `seed-helper/`
  - CLI that generates seed tokens, bcrypt hashes, and local seed env files.

## Flow (Where it comes from -> Where it goes)

Operator -> cmd/* -> internal platform/modules -> adapters/usecases

Each entrypoint builds the process graph and delegates real work to `internal/` layers.

## Why It Was Designed This Way

- Keep entrypoints thin and predictable.
- Centralize configuration, lifecycle, and observability early.
- Preserve clean boundaries: cmd never owns business rules.

## Recommended Practices Visible Here

- main packages only orchestrate; no domain rules.
- Configuration loaded via `internal/platform/config`.
- Graceful shutdown with context timeout.

## Differentials (Rare but Valuable)

- Swagger metadata applied at runtime from config.
- Seed tooling uses the same auth libraries as production.

## Run Locally

```bash
make dev
go build -o bin/aion-api ./cmd/aion-api
APP_ENV=development ./bin/aion-api
```

## What Should NOT Live Here

- Domain rules or validation.
- Mapping of HTTP/GraphQL payloads.
- Repository or external IO implementations.
