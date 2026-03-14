# System Design

AionApi follows Hexagonal (Ports & Adapters) architecture with clear bounded contexts and strict dependency direction.

## Architecture Snapshot

| Principle | Rule |
| --- | --- |
| Dependency rule | Core (`domain`, `ports`, `usecase`) does not import adapters or infrastructure |
| Context isolation | Contexts do not import each other directly |
| Thin adapters | Transport and persistence layers translate only |
| Semantic errors | Core returns semantic errors; adapters map to HTTP/GraphQL |
| Observability | Handlers, usecases, and repositories create spans and propagate context |

## Project Layout

```text
cmd/                      # Entrypoints
internal/
  <context>/
    core/
      domain/
      ports/
      usecase/
    adapter/
      primary/            # HTTP / GraphQL
      secondary/          # DB / cache / providers
  platform/               # Bootstrap, server, config, observability
  shared/                 # Shared constants, errors, response helpers
infrastructure/           # Docker, migrations, observability stack
docs/                     # MkDocs portal
```

## Request Lifecycle

### REST

1. Request enters primary adapter (`handler`).
2. Handler validates input and starts tracing span.
3. Handler calls input port (usecase) with `context.Context`.
4. Usecase orchestrates business logic via output ports.
5. Secondary adapters execute IO with context and map infra errors.
6. Handler maps semantic errors to HTTP response envelope.

### GraphQL

1. Resolver/controller receives operation.
2. Controller maps DTO/input and calls usecase.
3. Usecase returns domain result or semantic error.
4. GraphQL adapter maps output to schema model and response errors.

## Layer Responsibilities

| Layer | Responsibility | Must avoid |
| --- | --- | --- |
| Domain | Entities and invariants | Infra imports, transport types |
| Usecase | Business orchestration | SQL/GORM/HTTP details |
| Primary adapter | Validation, mapping, response | Business rules |
| Secondary adapter | Infrastructure IO and error translation | Domain orchestration |
| Platform | Wiring, config, server setup | Domain decisions |

## Error Mapping

| Semantic error | HTTP status |
| --- | --- |
| `validation` | 400 |
| `unauthorized` | 401 |
| `forbidden` | 403 |
| `not_found` | 404 |
| `conflict` | 409 |
| `internal` | 500 |

## Testing Strategy

- Usecases: table-driven unit tests with mocks.
- Adapters: mapping and boundary behavior tests.
- Contract/tooling: regenerate GraphQL and mocks when ports/schema change.
- Quality gate: `make verify` before merge.

## Decision Checklist (Before Coding)

- Is logic in the correct layer?
- Is context propagated to all boundaries?
- Are errors semantic and safely mapped?
- Are spans/logs present at adapter/usecase/repository boundaries?

## Next Step

Read [Platform Runtime](platform.md) for server composition, configuration, and observability runtime details.
