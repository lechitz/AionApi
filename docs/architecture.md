# Architecture — AionApi

This document describes the architectural principles, component responsibilities, and runtime flows in AionApi. It is intended as a technical reference for maintainers and for guiding future extensions.

> Summary
>
> - Hexagonal (Ports & Adapters) architecture: core business logic is independent of transport and infrastructure.
> - Primary adapters (HTTP / GraphQL) are thin; usecases contain business rules and orchestrations.
> - Observability-first: OpenTelemetry tracing, Prometheus metrics, and Grafana dashboards.
> - Consistent error model, standard response envelope, and strict use of context propagation.

---

## High-level layout

Top-level layout (simplified):

```
AionApi/
├─ cmd/                      # application entrypoints
├─ internal/                 # bounded contexts and platform (bootstrap, server, config)
│  ├─ <bounded-context>/     # auth, user, category, tag, habit, admin ...
│  │  ├─ core/               # business logic
|  |  │  ├─ domain/          # entities, value objects
│  │  │  ├─ ports/           # input/output interfaces
│  │  │  └─ usecase/         # application services
│  │  ├─ adapter/            # primary (http/graphql) + secondary (db/cache/token)
│  └─ platform/              # config, bootstrap, server, observability
├─ infrastructure/           # docker, migrations, observability resources
├─ docs/                     # MkDocs site
├─ swagger/                  # OpenAPI artifacts
└─ tests/                    # test helpers and generated mocks
```

Bounded contexts include `auth`, `user`, `category`, `tag`, `habit`, and `admin`. Each context defines its own domain model, ports, and adapters.

---

## Design principles

- Explicit contracts: the core defines input and output ports (interfaces) so implementations can vary without changing business code.
- Single responsibility: adapters are responsible only for translation and transport concerns; usecases handle business rules.
- Testability: core logic is pure Go and easy to unit-test with table-driven tests + mocks for ports.
- Observability: every handler, usecase and repository opens a span and adds structured attributes.
- Config-driven timeouts and limits: no hard-coded timeouts; read from configuration with safe defaults.

---

## Components and responsibilities

- Core (Usecases)
  - Located at `internal/<ctx>/core/usecase`.
  - Implement application flows, validations, and orchestrations. Return domain values or typed semantic errors.

- Primary adapters (HTTP / GraphQL)
  - Located at `internal/<ctx>/adapter/primary`.
  - Decode requests, validate DTOs, start traces/spans, call usecases, map domain → transport response.

- Secondary adapters (DB, Cache, Token, Logger)
  - Located at `internal/<ctx>/adapter/secondary`.
  - Implement output ports declared by the core. They keep infra-specific code isolated (GORM, redis client, external SDKs).

- Platform layer
  - `internal/platform` contains bootstrap wiring, server setup (HTTP / GraphQL), shared middleware, config loader, and observability wiring.

- Shared packages
  - `internal/shared/*` holds shared helpers: `sharederrors`, `httpresponse`, `constants`, and test helpers.

---

## Request lifecycle (detailed)

### REST flow (example: Update password)
1. HTTP request arrives at handler (`adapter/primary/http`).
2. Handler parses the request and validates the DTO.
3. Handler starts a trace span and sets common attributes (route, method, user_id if available).
4. Handler calls the usecase input port with `context.Context`.
5. Usecase executes domain logic, calling output ports (repositories, caches, token provider).
6. Output adapters perform context-aware IO (e.g. `db.WithContext(ctx)`) and map infra errors to semantic domain errors.
7. Usecase returns domain values or a semantic error.
8. Handler maps results to a standardized JSON envelope (`internal/shared/httpresponse`) and ends the span.

### GraphQL flow (example: Create Category)
- Resolver → GraphQL handler → DTO mapping → usecase → repository → domain mapping → resolver returns GraphQL model.
- GraphQL middleware should also start a top-level span and propagate context to resolvers.

---

## Error model and HTTP mapping

- The repository uses typed semantic errors (validation, not_found, conflict, unauthorized, internal, etc.) defined in `internal/shared/sharederrors`.
- The adapter maps these semantic errors to HTTP status codes consistently:
  - validation → 400
  - unauthorized → 401
  - forbidden → 403
  - not_found → 404
  - conflict → 409
  - internal → 500
- All responses use a consistent envelope shape to simplify clients and monitoring.

---

## Persistence and migrations

- Migrations are SQL files under `infrastructure/db/migrations` and applied with the `migrate` CLI via Make targets (`make migrate-up`).
- Repositories implement explicit mapping between DB models and domain models to keep the domain clean of ORM types.
- Repositories always use context-aware DB operations and return translated semantic errors.
- Soft deletes: the application ensures reads exclude logically deleted rows and updates return not-found when no rows were affected.

---

## Authentication & security

- Authentication usecase issues access and refresh tokens through an `AuthProvider` output port.
- Tokens may be stored / referenced in a cache (`AuthStore`) for revocation support.
- HTTP authentication middleware validates `Authorization: Bearer <token>` and populates context with `user_id` and claims.
- Secrets must never be logged (the logger and helpers ensure sensitive fields are redacted).
- Follow OWASP guidance for JWT usage, cookie flags and CORS policies.

---

## Observability

- Tracing: OpenTelemetry spans are created at the handler and propagated into usecases and repositories. Naming convention: `aionapi.<context>.<component>` (e.g. `aionapi.user.usecase`).
- Metrics: Prometheus instrumentation and scrape config live under `infrastructure/observability/prometheus`.
- Dashboards & logs: Grafana dashboards and Fluent Bit/Loki configs are under `infrastructure/observability`.

Suggested local env to enable OTLP export:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
export OTEL_SERVICE_NAME="AionApi"
export OTEL_SERVICE_VERSION="0.1.0"
```

---

## Testing strategy

- Unit tests: focus on usecases; use gomock-generated mocks for output ports (`tests/mocks`), and helper suites in `tests/setup`.
- Integration tests: exercise adapters and real infra (run against Dockerized Postgres/Redis). Use a dedicated test DB and cleanup strategy.
- End-to-end / smoke tests: run against the full `make dev` stack and validate common flows (health, login, create user, basic GraphQL ops).
- Coverage: `make test-cover` produces coverage artifacts in `tests/coverage/`.

---

## Code generation & developer tooling

- GraphQL: schema fragments live in contexts; `make graphql` collects them and runs `gqlgen`.
- Mocks: `make mocks` runs `mockgen` and stores generated mocks under `tests/mocks`.
- Quality: `make lint`, `make lint-fix`, `make format` and `golangci-lint` configuration are part of the developer workflow.

---

## Adding a new bounded context — checklist

1. Create `internal/<ctx>/core/domain` for entities and value objects.
2. Define input/output ports under `internal/<ctx>/core/ports`.
3. Implement usecases in `internal/<ctx>/core/usecase`.
4. Add primary adapters (`adapter/primary/http` and/or `graphql`) for transport.
5. Add secondary adapters (`adapter/secondary/db`, `cache`, `token`) implementing the output ports.
6. Wire concrete adapters in `internal/platform/bootstrap`.
7. Mount routes in `internal/platform/server/http/composer.go`.
8. Add tests and generate mocks (`make mocks`).
9. Document the context in `docs/` and update `mkdocs.yml` nav if needed.

---

## Deployment and operational notes

- Use environment-specific config files for production vs staging vs local development.  
- Keep secrets in a secure vault or secrets manager (do not commit to repo).  
- CI jobs should run `make lint`, `make test`, and `make docs.gen` (if docs changes are present).  
- Consider blue/green or rolling deployment patterns for production services.

---

## Diagram (Mermaid)

You can embed a Mermaid diagram into MkDocs (if plugin enabled). Example:

```mermaid
flowchart LR
  Client -->|HTTP/GraphQL| Handler[Handler (primary adapter)]
  Handler -->|calls| Usecase[Usecase (core)]
  Usecase --> Repo[Repository (secondary adapter)]
  Usecase --> AuthProvider[AuthProvider (secondary)]
  Repo --> DB[(Postgres)]
  AuthProvider --> Cache[(Redis)]
```

---

## Conventions & style

- Always propagate `context.Context` and honor cancellations and deadlines.  
- Keep handlers thin: validation, tracing, DTO mapping, usecase invocation, response mapping.  
- Keep usecases pure from infra concerns; use small interfaces for side effects.  
- Avoid magic strings — centralize keys in `internal/shared/constants`.

---

