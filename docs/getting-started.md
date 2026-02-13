# Getting Started

This guide takes you from clone to a healthy local AionApi environment with migrations and baseline validation.

## Prerequisites

| Tool | Required version | Notes |
| --- | --- | --- |
| Go | 1.24+ | Required for local commands and tooling |
| Docker | Recent | Compose v2 required |
| GNU Make | Recent | Primary task runner |
| jq | Optional | Better JSON output in terminal |

## 1) Clone and Install Tools

```bash
git clone https://github.com/lechitz/AionApi.git
cd AionApi
make tools-install
```

## 2) Start Local Stack

```bash
make dev
```

Core services expected in dev mode:

- API: `http://localhost:8080`
- Postgres: `localhost:5432`
- Localstack: `localhost:4566`

## 3) Apply Migrations

```bash
export MIGRATE_BIN="$(go env GOPATH)/bin/migrate"
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
export MIGRATION_PATH="infrastructure/db/migrations"
make migrate-up
```

## 4) Optional Seeds

```bash
make seed-all
```

## 5) Validate Health

```bash
curl -s http://localhost:8080/aion/health | jq
```

## Development Commands

| Command | Purpose |
| --- | --- |
| `make dev` | Build and start dev stack |
| `make dev-fast` | Start quickly without full rebuild |
| `make dev-down` | Stop stack, keep volumes |
| `make dev-clean` | Stop stack and remove volumes |
| `make rebuild-api` | Force rebuild only API service |
| `make rebuild-chat` | Force rebuild only chat service |
| `make rebuild-dashboard` | Force rebuild only dashboard service |
| `make verify` | Local quality pipeline |

## Hot Reload (Dev)

`make dev` is configured for hot reload across services in the dev compose profile:

- API (Go): Air reloads on `.go` changes.
- Chat (Python): Uvicorn reloads on `.py` changes.
- Dashboard (TypeScript): Vite HMR updates UI on `.ts/.tsx/.css` changes.

Expected behavior:

- Most code edits do not require full container rebuild.
- Use targeted rebuild commands only when Dockerfile/dependency layers changed.

Useful logs:

```bash
make logs-api
make logs-chat
make logs-dashboard
```

## Codegen and Tests

```bash
make graphql
make mocks
make test
make test-cover
make lint
```

## API Quick Checks

### REST

```bash
curl -s http://localhost:8080/aion-api/v1/users/all
```

### GraphQL

```bash
curl -s http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"query { listAllCategories { id name colorHex } }"}' | jq
```

## Troubleshooting

- Migration errors: verify `MIGRATION_DB` and that Postgres is ready.
- Missing tools: rerun `make tools-install` and check `$(go env GOPATH)/bin` in `PATH`.
- Port conflicts: stop local services already using 8080/5432/4566.
- Hot reload not updating:
  - check service logs (`make logs-api`, `make logs-chat`, `make logs-dashboard`)
  - ensure source volumes are mounted in dev compose profile
  - run targeted rebuild (`make rebuild-api`, `make rebuild-chat`, or `make rebuild-dashboard`)

## Next Step

Continue with [System Design](architecture.md) for architecture and boundaries.
