# Getting Started

Welcome to **AionApi**. This page takes you from zero to a working local environment with Docker, the database migrated, tests passing, and a few REST/GraphQL calls verified.

---

## TL;DR (fast path)

```bash
# 1) Clone and enter the project
git clone https://github.com/lechitz/AionApi.git
cd AionApi

# 2) (Optional, recommended) install local dev tools
make tools-install

# 3) Bring up the DEV stack (Docker + Postgres)
make dev

# 4) Apply database migrations
export MIGRATE_BIN="$(go env GOPATH)/bin/migrate"
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
export MIGRATION_PATH="infrastructure/db/migrations"
make migrate-up

# 5) (Optional) Seed sample data
make seed-all

# 6) Health check
curl -s http://localhost:8080/aion/health | jq
```

If things are ok, you should see a JSON payload containing the service name/version/environment.

---

## Prerequisites

* **Go** 1.24+
* **Docker** and **Docker Compose** (v2)
* **make** (GNU Make)
* (Optional) **jq** to pretty-print JSON on the CLI

> If you prefer not to install Go on your host (using only Docker), you still need the `migrate` CLI to run DB migrations via Make. Use `make tools-install` to install all expected tools under your `GOPATH/bin`.

---

## Repository layout (essentials)

```
AionApi/
├─ infrastructure/
│  ├─ db/
│  │  ├─ migrations/         # *.sql (schema & changesets)
│  │  └─ seed/               # sample data
│  └─ docker/
│     ├─ Dockerfile          # app build
│     └─ environments/
│        └─ dev/docker-compose-dev.yaml
├─ internal/                 # domain, ports & adapters (HTTP/GraphQL/DB/etc.)
├─ makefiles/                # grouped Make targets (migrate, docker, test, etc.)
├─ Makefile                  # includes makefiles/*
└─ README.md
```

---

## Environment configuration

For local development, `docker-compose-dev.yaml` already provides sensible defaults (Postgres on `localhost:5432`, API on `localhost:8080`).
If you need to customize values, create a `.env.dev` (optional):

```bash
# illustrative example — adjust to your needs
cp infrastructure/docker/environments/example/.env.example infrastructure/docker/environments/dev/.env.dev
# edit infrastructure/docker/environments/dev/.env.dev
```

> If you skip `.env.dev`, the compose defaults are enough to boot the stack.

---

## Bring up the DEV stack

```bash
make dev       # (= image build + docker compose up)
# or, if you built recently:
make dev-up
# to stop and remove volumes for this stack:
make dev-down
```

Expected services:

* **Postgres**: `localhost:5432` (db `aionapi`, user `aion`, pass `aion`)
* **API**: `localhost:8080` (REST base prefix: `/aion-api`, GraphQL: `/graphql`)

Health check:

```bash
curl -s http://localhost:8080/aion/health | jq
```

---

## Database: migrations & seeds

### Install `migrate` (once)

```bash
make tools-install
# ensures: migrate, golangci-lint, gqlgen, etc. in your $GOPATH/bin
```

### Environment variables expected by migration targets

```bash
export MIGRATE_BIN="$(go env GOPATH)/bin/migrate"
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
export MIGRATION_PATH="infrastructure/db/migrations"
```

### Apply/rollback migrations

```bash
make migrate-up        # apply all "up" migrations
make migrate-down      # rollback one step
# advanced usage:
make migrate-force VERSION=20250101010101  # force version (⚠️ beware)
```

### Seeds

```bash
# executes SQL seed scripts inside the DEV Postgres container
make seed-users
make seed-categories
make seed-all
```

---

## Tests, coverage & code quality

```bash
make test           # unit tests with -race
make test-cover     # coverage + HTML at tests/coverage/coverage.html
make test-html-report
make format         # goimports + golines (gofumpt)
make lint           # golangci-lint
make lint-fix       # attempt autofix
make verify         # local pipeline: graphql → mocks → lint → test → coverage
```

Open the coverage report in your browser:

```
tests/coverage/coverage.html
```

---

## Codegen (GraphQL & Mocks)

> These targets keep generated artifacts in sync with the source code.

```bash
# (re)generate GraphQL types/resolvers from *.graphqls
make graphql

# generate mocks (GoMock) for output-port interfaces (saved in tests/mocks/)
make mocks
# namespaced filenames to avoid basename collisions:
make mocks NAMESPACE=1
# generate only for a specific context:
make mocks CONTEXT=user
```

---

## Using the API (REST)

> REST base prefix: **`/aion-api`**

### Create user

```bash
curl -X POST http://localhost:8080/aion-api/v1/users/create \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "s3cret"
  }'
```

### Login (obtain token)

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/aion-api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"s3cret"}' | jq -r '.token')
echo "$TOKEN"
```

### List users (protected route)

```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/aion-api/v1/users/all | jq
```

> Other useful routes:
> `GET /aion-api/v1/users/{user_id}`
> `PUT /aion-api/v1/users/` (update self)
> `PUT /aion-api/v1/users/password`
> `DELETE /aion-api/v1/users/`

---

## Using the API (GraphQL)

* **Endpoint**: `http://localhost:8080/graphql`

List categories (query):

```bash
curl -s http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"query { listAllCategories { id name colorHex } }"}' | jq
```

Create category (mutation):

```bash
curl -s http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"mutation { createCategory(input:{name:\"Work\", colorHex:\"#3366FF\"}) { id name colorHex } }"}' | jq
```

> If your instance requires auth for GraphQL, include `Authorization: Bearer $TOKEN` in the headers.

---

## Observability (optional in DEV)

The platform is wired for **OpenTelemetry** (traces/metrics).
To export to a local Collector, configure environment variables like:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
export OTEL_SERVICE_NAME="AionApi"
export OTEL_SERVICE_VERSION="0.1.0"
```

Useful files:

* `infrastructure/observability/otel/otel-collector-config.yaml`
* `infrastructure/observability/grafana/` (dashboards)
* `infrastructure/observability/prometheus/prometheus.yml`

> You can extend your local stack to include these services via Docker Compose, as needed.

---

## Tips & Troubleshooting

* **`migrate: command not found`**
  Run `make tools-install` and ensure `$(go env GOPATH)/bin` is in your `PATH`.

* **Cannot connect to Postgres**
  Check running containers with `docker ps`, free port `5432`, and the DSN in `MIGRATION_DB`.

* **Port 8080 is already in use**
  Stop the service occupying it or change the HTTP port via env vars and restart.

* **401 on protected routes**
  Verify the `Authorization: Bearer <TOKEN>` header is present and the token is valid.

* **Unexpected 404/405**
  Remember the REST **base prefix** `/aion-api`.

---

## What’s next

* **Architecture**: high-level view (Ports & Adapters), conventions, and boundaries.
* **Platform**: Config, Router (port + chi), Observability.
* **API Reference**: details for each REST route and GraphQL operation.

> These pages will be published in the next sections of the documentation.

