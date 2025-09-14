# Architecture

This document explains how AionApi is organized and how requests move through the system. It’s a practical guide for contributors who want to understand the Hexagonal (Ports & Adapters) structure, observability, config, and testing strategy.

> TL;DR
>
> - **Ports & Adapters** with clean, testable boundaries  
> - **Primary adapters** (HTTP/GraphQL) stay thin; **Usecases** hold business rules; **Secondary adapters** isolate infra  
> - **Observability-first**: OpenTelemetry traces + structured logs, Prometheus metrics, Grafana dashboards  
> - **Consistency**: shared error types, response helpers, constants, and middleware

---

## 1) High-level layout

```

cmd/                    # app entrypoint
infrastructure/         # docker, migrations, otel/prometheus/grafana, loki/fluentbit
internal/ <context>/            # bounded contexts (auth, user, category, tag, admin)
adapter/
primary/          # inbound adapters (http / graphql)
secondary/        # outbound adapters (db, cache, token, logger)
core/               # domain + ports + usecases (business logic)
platform/             # cross-cutting platform: config, bootstrap, server, observability
shared/               # reusable helpers: constants, responses, validation, errors
makefiles/              # grouped Make targets used by the root Makefile
tests/                  # unit test suites, mocks, fixtures, coverage

```

**Bounded contexts currently present:** `auth`, `user`, `category`, `tag`, `admin`.

---

## 2) Hexagonal (Ports & Adapters)

**Why:** decouple business logic from frameworks and vendors; make the core testable and stable.

- **Core / Usecases** (`internal/*/core/usecase`):  
  Implements input-port interfaces. Pure go, no HTTP/ORM/Redis types. Accepts `context.Context`, returns domain or semantic errors.
- **Primary Adapters** (`internal/*/adapter/primary`):  
  Transport-facing (HTTP or GraphQL). Decode/validate → call input ports → map domain → response. **No business rules** here.
- **Secondary Adapters** (`internal/*/adapter/secondary`):  
  Implement output-port interfaces (DB/Cache/Token/Logger). Hide GORM/Redis/JWT behind small interfaces.

**Shared contracts & helpers** live under:
- `internal/shared/sharederrors` — semantic errors (validation, not_found, conflict, unauthorized, internal, …)
- `internal/shared/httpresponse` — opinionated JSON envelope + error→HTTP mapping
- `internal/shared/constants` — keys for logs/headers/claims/tracing
- `internal/shared/handlerhelpers` — boundary validation + tracing-friendly responses

---

## 3) Request lifecycles

### 3.1 REST (example: Update user password)

1. **HTTP route** (primary adapter)  
   `internal/user/adapter/primary/http/handler/update_password_user_handler.go`
2. **DTO validation**  
   `internal/user/adapter/primary/http/dto/update_password_dto.go`
3. **Call usecase**  
   `internal/user/core/usecase/update_password.go`
4. **Usecase orchestration**  
   - `repo.GetByID` → `hasher.Compare` → `hasher.Hash` → `repo.Update`  
   - `tokenProvider.Generate` → `authStore.Save`
5. **Response mapping**  
   `internal/shared/httpresponse` writes a standardized JSON and status code
6. **Observability**  
   Span created in handler & usecase; canonical attributes added (user_id, operation, http_status)

### 3.2 GraphQL (example: Create category)

1. **Resolver calls handler**  
   `internal/category/adapter/primary/graphql/resolver/*` → handler `Create`
2. **Handler** maps GraphQL input → domain input, starts span  
   `internal/category/adapter/primary/graphql/handler/create.go`
3. **Usecase** applies rules and calls repo via output ports  
   `internal/category/core/usecase/create.go`
4. **Repo** persists with GORM and maps DB ↔ domain  
   `internal/category/adapter/secondary/db/repository/create.go`
5. **Handler** maps domain → GraphQL model and returns to resolver

---

## 4) Platform layer

- **Config** (`internal/platform/config`)  
  Loads & validates env vars, normalizes paths/timeouts, may generate a **dev-only secret** if missing.
- **Bootstrap** (`internal/platform/bootstrap`)  
  Wires concrete adapters (DB, cache, token, hasher, logger) and constructs services (usecases).
- **HTTP server** (`internal/platform/server/http`)  
  - **Router port** + **chi adapter**  
  - **Generic handlers**: `/health`, 404/405, error/recovery  
  - **Middlewares**: `requestid`, `recovery`  
  - **Composer** mounts all context routes and GraphQL
- **GraphQL server** (`internal/platform/server/graph`)  
  Gathers schema modules from contexts; `make graphql` runs gqlgen and assembles the executable schema.
- **Observability** (`internal/platform/observability`)  
  OTel tracer/meter providers (OTLP over HTTP), resource attributes, and helper utilities.

---

## 5) Persistence & data rules

- **ORM**: GORM models live in `internal/<context>/adapter/secondary/db/model`.  
  Mapping functions in `mapper/` translate DB ↔ domain; the domain never sees GORM structs.
- **Repositories**:  
  Implement CRUD with **context-aware** DB calls (`db.WithContext(ctx)`), set span attributes, and translate driver errors to `sharederrors`.
- **Soft delete**:  
  Reads must exclude `deleted_at`; updates return not-found when rows affected = 0.
- **Migrations**:  
  SQL files in `infrastructure/db/migrations`. Use `make migrate-up` / `make migrate-down`.  
  DSN provided via `MIGRATION_DB` (see `makefiles/migrate.mk`).

---

## 6) Security & Auth (HTTP)

- **Login**: `auth` usecase verifies credentials, issues token (`AuthProvider`), stores reference in cache (`AuthStore`).  
- **Middleware** (primary HTTP adapter): validates `Authorization: Bearer <token>`, resolves `user_id`, injects into `context`.  
- **Cookies**: helpers in `internal/platform/server/http/helpers/httpresponse` to set secure/HTTPOnly cookies when applicable.  
- **Never log secrets** (passwords, raw tokens, cookie values).

---

## 7) Observability

- **Tracing**:  
  Consistent tracer names (e.g., `aionapi.user.usecase`, `aionapi.user.repository`, `aionapi.graphql.handler`).  
  Each handler/usecase/repo opens a span, sets canonical attributes (operation, ids, status), records errors.
- **Metrics**:  
  Prometheus scrape via `infrastructure/observability/prometheus/prometheus.yml`.  
- **Dashboards**:  
  Grafana provisioning + a sample dashboard under `infrastructure/observability/grafana/dashboards/`.
- **Logs**:  
  Structured with contextual keys from `shared/constants/commonkeys`. Fluent Bit + Loki optional.

---

## 8) Errors & responses

- **Domain errors** come from `shared/sharederrors`.  
- **HTTP mapping** is centralized in `shared/httpresponse`:  
  - `validation` → 400  
  - `unauthorized` → 401  
  - `forbidden` → 403  
  - `not_found` → 404  
  - `conflict` → 409  
  - `rate_limited` → 429  
  - `internal` → 500
- **Consistency**: All adapters return the same envelope shape.

---

## 9) Concurrency & context

- Always pass `context.Context` through handlers, usecases, and repos.  
- Use deadlines/timeouts from config (HTTP/GraphQL).  
- **Never** block without honoring ctx cancellation.

---

## 10) Testing strategy

- **Unit tests** target **usecases** with port **mocks** (gomock).  
  Suites in `tests/setup` and generated mocks in `tests/mocks`.
- **Transport tests** exercise handlers with fake services.  
- **Coverage**: `make test-cover` → `tests/coverage/coverage.html`.  
- **GraphQL codegen** + **mocks** are part of the local pipeline: `make verify`.

---

## 11) Code generation & tooling

- **GraphQL**: `make graphql` copies `*.graphqls` from contexts into the platform graph module and runs `gqlgen` + `go mod tidy`.
- **Mocks**: `make mocks` generates gomock doubles for all output ports under `tests/mocks/` (flat or namespaced mode).
- **Quality**: `make format`, `make lint`, `make lint-fix`.
- **Dev/Prod stacks**: `make dev`, `make prod` orchestrate Docker Compose with env files under `infrastructure/docker/environments/`.

---

## 12) Adding a new bounded context (checklist)

1. **Core**  
   - `internal/<ctx>/core/domain` (entities/VOs)  
   - `internal/<ctx>/core/ports` (input/output)  
   - `internal/<ctx>/core/usecase` (service + validations + errors)
2. **Adapters**  
   - Primary: `http` and/or `graphql` (DTOs, handlers, resolver, thin mapping)  
   - Secondary: `db` (model/mapper/repository), `cache`, or any external port implementation
3. **Platform wiring**  
   - Add repos/services in `internal/platform/bootstrap`  
   - Mount routes in `internal/platform/server/http/composer.go` (and/or GraphQL handler)
4. **Schema / codegen**  
   - If GraphQL, drop `schema/<ctx>.graphqls` and run `make graphql`
5. **Tests**  
   - Add gomock expectations & table tests for usecases  
   - Add handler tests with fake services
6. **Docs**  
   - Create `docs/<ctx>.md` and link from `mkdocs.yml`

---

## 13) Style & conventions

- **Primary adapters**: thin, transport-only, start spans, use shared helpers, never touch ORM/cache.  
- **Usecases**: pure business rules; input/output ports only; return domain + semantic errors.  
- **Secondary adapters**: context-aware DB/Cache calls; map infra errors to domain errors; log metadata only.  
- **Shared constants**: no magic strings (claims, headers, log keys, trace attrs).

---

## 14) Glossary

- **Primary adapter**: inbound interface (HTTP/GraphQL) that calls the domain via input ports.  
- **Usecase (core)**: application/business logic that coordinates ports.  
- **Secondary adapter**: outbound implementation for repositories or external services.  
- **Port**: interface owned by the domain that expresses what it needs (output) or offers (input).  
- **Semantic error**: a typed error category (`validation`, `not_found`, …) that maps to consistent transports.

```
