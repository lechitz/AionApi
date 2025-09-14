# Makefiles (child package)

**Folder:** `makefiles/`
**Used by:** the root `Makefile` (includes these targets)

## Responsibility

* Provide a **reusable toolbox** of `make` targets for local dev, Docker orchestration, DB migrations, code-gen (mocks/GraphQL), testing/coverage, seeding, and code quality.
* Keep **implementation** here and keep the root `Makefile` small (just variables + `include`).

## How it works

* The root `Makefile` exports a few env vars (paths, DSNs) and **includes** these files (e.g. `include makefiles/*.mk`).
* Targets are grouped by area (Docker, Migrations, Tooling, Codegen, Testing, Seeds, Quality) and are all **PHONY**.
* Some targets are guarded by required vars (they will explain what’s missing and exit with ❌).

---

## Targets (by area)

### Docker environment

* `build-dev` / `clean-dev` — build/cleanup the `:dev` image (`APPLICATION_NAME`).
* `dev-up` / `dev-down` / `dev` — compose **DEV** stack. `dev-up` loads `$(ENV_FILE_DEV)` and uses `$(COMPOSE_FILE_DEV)` (removes the `postgres` service once to rebuild volume).
* `build-prod` / `clean-prod` — build/cleanup the `:prod` image.
* `prod-up` / `prod-down` / `prod` — compose **PROD-like** stack from `$(ENV_FILE_PROD)` + `$(COMPOSE_FILE_PROD)`.
* `docker-clean-all` — ⚠️ removes **all** containers, volumes and images on your machine.

### Migrations (golang-migrate)

* `migrate-up` — apply all up migrations.
* `migrate-down` — roll back **one** step.
* `migrate-force VERSION=X` — force schema version to `X` (⚠️ dangerous).
* `migrate-new` — interactive: create a new `*.sql` pair under `$(MIGRATION_PATH)`.

**Requires**

* `MIGRATE_BIN` (path to `migrate` CLI)
* `MIGRATION_DB` (DSN: e.g., `postgres://user:pass@host:5432/db?sslmode=disable`)
* `MIGRATION_PATH` (e.g., `infrastructure/db/migrations`)

### Tooling

* `tools-install` — installs local tools: `gofumpt`, `golines`, `goimports`, `golangci-lint`, `fieldalignment`, `gotestsum`, `migrate`, `gqlgen`.

### Code generation

**Mocks**

* `mocks` — generate Go mocks from **ports/output** under all contexts **+** platform outputs into `tests/mocks/` (flat package `mocks`).

    * `NAMESPACE=1` → prefixes filenames with path (`/`→`__`) to avoid basename collisions.
    * `CONTEXT=user` → generate only for that context.
* `mocks-list` — print discovered sources/targets (audit).
* `clean_mocks` — remove `tests/mocks/`.
* Internally runs a **collision guard** in flat mode and fails fast with guidance.

**GraphQL**

* `graphql` — copies all `internal/**/adapter/primary/graphql/schema/*.graphqls` to `internal/platform/server/graph/schema/_modules/`, runs `gqlgen`, then `go mod tidy`.

### Testing & coverage

* `test` — run unit tests with `-race`.
* `test-cover` — run coverage, **filter out** `Mock` files, and produce `$(COVERAGE_DIR)/coverage.html`.
* `test-html-report` — run tests with JUnit XML via `gotestsum` → `$(COVERAGE_DIR)/junit-report.xml`.
* `test-ci` — coverage to `$(COVERAGE_DIR)/coverage.out` (no HTML).
* `test-clean` — clean coverage artifacts.

> Ensure `COVERAGE_DIR` exists (e.g., `coverage/`) or define it in the root `Makefile`.

### Seeds

* `seed-users` / `seed-categories` / `seed-all` — execute `psql` inside `$(POSTGRES_CONTAINER)` using seed files under `infrastructure/db/seeds/`.

Defaults:

```
POSTGRES_CONTAINER=postgres-dev
POSTGRES_USER=aion
POSTGRES_DB=aionapi
```

### Code quality & pipeline

* `format` — `goimports` + `golines (gofumpt)`.
* `lint` / `lint-fix` — `golangci-lint` (check/fix).
* `verify` — local pipeline: `graphql → mocks → lint → test → test-cover → test-ci → test-clean`.

---

## Required / useful variables (set in root `Makefile` or shell)

```make
# Images / compose
APPLICATION_NAME := aion-api
ENV_FILE_DEV     := infrastructure/docker/environments/dev.env
COMPOSE_FILE_DEV := infrastructure/docker/docker-compose.dev.yml
ENV_FILE_PROD    := infrastructure/docker/environments/prod.env
COMPOSE_FILE_PROD:= infrastructure/docker/docker-compose.yml

# Migrations
MIGRATE_BIN    := $(GOPATH)/bin/migrate
MIGRATION_DB   := postgres://aion:aion@localhost:5432/aionapi?sslmode=disable
MIGRATION_PATH := infrastructure/db/migrations

# Coverage
COVERAGE_DIR := coverage
```

Optional flags at call time:

```bash
make mocks NAMESPACE=1
make mocks CONTEXT=user
make migrate-force VERSION=20250101010101
```

---

## Include a pattern (root `Makefile`)

```make
# Root Makefile (example)
include makefiles/*.mk
```

---

## Safety notes

* `dev-down`/`prod-down` use `-v` → remove volumes of those stacks.
* `docker-clean-all` wipes **everything** (containers, volumes, images) — use with care.
* `migrate-force` changes a schema version without running migrations — know what you’re doing.

---

## Quickstart

```bash
# One-time tooling
make tools-install

# Dev stack
export ENV_FILE_DEV=...; export COMPOSE_FILE_DEV=...
make dev   # (= build-dev + dev-up)

# DB migrations
export MIGRATE_BIN=$(go env GOPATH)/bin/migrate
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
export MIGRATION_PATH=infrastructure/db/migrations
make migrate-up

# Codegen
make graphql
make mocks

# Tests & quality
make verify
```

This keeps ops/dev flows consistent across contexts while the root `Makefile` only wires variables and includes.
