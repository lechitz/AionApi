# cmd/api (Main API Server)

The primary entrypoint for the Aion API application (production binary).

## What It Does

- Wires all modules via Uber Fx (dependency injection)
- Configures Swagger metadata at runtime
- Boots HTTP server (GraphQL + REST endpoints)
- Initializes observability (OTel tracing, logging)
- Manages graceful shutdown

## Build

```bash
# Production build:
go build -o bin/api ./cmd/api

# Docker build:
make build-dev
make build-prod
```

## Run

```bash
# Docker (recommended - full stack with hot reload):
make dev

# Local with Air hot reload (API on host, deps in Docker):
make dev-local

# Direct:
go run ./cmd/api
```

## Configuration

All config via environment variables (see `infrastructure/docker/environments/dev/.env.dev`).

Key vars:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `CACHE_ADDR` (Redis)
- `OTEL_EXPORTER_OTLP_ENDPOINT`
- `SECRET_KEY` (JWT signing)

## Architecture

This is a **thin orchestrator** following Hexagonal Architecture:
- No business logic here (lives in `internal/<ctx>/core/usecase/`)
- No infrastructure details (wired via `internal/platform/`)
- Just composition and bootstrap

## Development Tools

For seeding, testing, utilities, see `hack/`:
- `hack/tools/seed-caller/` - Seed via API
- `hack/tools/seed-helper/` - Generate JWT/bcrypt
- `hack/dev/test-*.sh` - Test scripts
