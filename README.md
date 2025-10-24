# AionApi

This repository hosts AionApi — a modular backend service implemented in Go that exposes both REST and GraphQL APIs for habit and diary management. The project follows a Ports & Adapters (Hexagonal) architecture and is designed for testability, observability, and iterative development.


---

## Quick links
- Documentation site: [AionApi - Github Pages](https://lechitz.github.io/AionApi/)

> Visit the live [AionApi - Swagger UI](https://lechitz.github.io/AionApi/swagger-ui/) to interact with the API and try requests.
---

## Overview

Aion is a habit management system focused on helping users organize, track and analyze routines to improve physical, mental and emotional well-being. The API provides endpoints for user management, authentication, categories, tags and diary entries (personal and professional).

### Technology stack
- Go 1.24
- chi router for HTTP
- gqlgen for GraphQL
- GORM for PostgreSQL integration
- Redis for caching / session management
- Docker & Docker Compose for local environments
- OpenTelemetry (OTel) for traces and metrics
- Jaeger / Prometheus / Grafana for observability
- zap for structured logging

### Key features (current / planned)
- Habit and diary entry management (personal and diary contexts)
- Tagging and category system with GraphQL support
- Authentication (access and refresh tokens) and session invalidation
- Database migrations and seed data for reproducible dev environments
- Observability (traces, metrics, dashboards)
- Developer-friendly tooling: codegen, mocks, linters and formatting

---

## Project management

This repository is organized using a public [AionApi - GitHub Projects](https://github.com/users/lechitz/projects/1) where tasks, issues, and epics are tracked. The board provides visibility into ongoing work and completed milestones, keeping development structured and transparent.

---

## Installation

### Prerequisites
- Go 1.24 or newer
- Git
- Docker & Docker Compose (for containerized development)
- Make (GNU Make)

Clone and prepare
```bash
git clone git@github.com:lechitz/AionApi.git
cd AionApi
```
Install development tools (recommended)
```bash
make -f makefiles/tooling.mk tools-install
```
Download modules
```bash
go mod download
```

---

## Configuration

Copy the example environment and edit for local development:
```bash
cp infrastructure/docker/environments/example/.env.example infrastructure/docker/environments/dev/.env.dev
# edit infrastructure/docker/environments/dev/.env.dev
```
Start the development environment (Docker)
```bash
make dev
```
Run migrations (example)
```bash
export MIGRATE_BIN="$(go env GOPATH)/bin/migrate"
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
export MIGRATION_PATH="infrastructure/db/migrations"
make migrate-up
```
Seed sample data (optional)
```bash
make seed-all
```

---

## Development & common commands

Run formatting and linting
```bash
make lint        # format + golangci-lint checks
make lint-fix    # attempt autofix
```
Run tests and coverage
```bash
make test
make test-cover  # generates coverage report in tests/coverage/
```
Code generation
```bash
make graphql  # gqlgen codegen
make mocks    # generate gomock mocks in tests/mocks/
```
Build the server binary
```bash
go build -o bin/aion-api ./cmd/aion-api
```
Run the built server (example)
```bash
export APP_ENV=development
export DATABASE_URL=postgres://aion:aion@localhost:5432/aionapi?sslmode=disable
./bin/aion-api
```

---

## API summary

REST base prefix: `/aion-api/v1`
- `GET  /aion/health` — service health
- `POST /aion-api/v1/user/create` — create user
- `GET  /aion-api/v1/user/all` — list users
- `GET  /aion-api/v1/user/{user_id}` — get user by ID
- `PUT  /aion-api/v1/user` — update user
- `PUT  /aion-api/v1/user/password` — update logged user's password
- `DELETE /aion-api/v1/user` — soft-delete logged user
- `POST /aion-api/v1/auth/login` — login and obtain tokens
- `POST /aion-api/v1/auth/logout` — invalidate session

GraphQL endpoint (example): `/aion-api/v1/graphql`
- Queries: `GetAllCategories`, `GetCategoryByID`, `GetCategoryByName`, `GetAllTags`, `GetTagByID`
- Mutations: `CreateCategory`, `CreateTag`, `UpdateCategory`, `SoftDeleteCategory`

For full API spec, consult `swagger/swagger.yaml` and the generated JSON (`swagger/swagger.json`).

---

## Architecture (preview)

The codebase follows a Hexagonal architecture (Ports & Adapters): business logic (usecases) is isolated from transport and infrastructure. Primary adapters (HTTP/GraphQL) are thin layers that map requests to input ports and format responses. Secondary adapters implement persistence, cache and external integrations behind defined interfaces.

For a detailed architecture overview, see: [AionApi - Architecture](https://lechitz.github.io/AionApi/architecture/)

---

## Observability & monitoring

The platform integrates OpenTelemetry. Configuration for local preview is available under `infrastructure/observability/` and includes collector config, Prometheus and Grafana dashboards. To enable local OTLP exporter, set:
```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
export OTEL_SERVICE_NAME="AionApi"
export OTEL_SERVICE_VERSION="0.1.0"
```



---

## License

This project is available under the MIT License — see the `LICENSE` file for details.

---

