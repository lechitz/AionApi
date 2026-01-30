# internal

Application code for AionAPI, organized by bounded contexts and the shared platform layer (Hexagonal/Clean Architecture).

## Package Composition

- `<context>/`
  - Domain-focused modules (auth, user, chat, etc.).
- `adapter/`
  - Shared primary adapters (GraphQL infrastructure) and shared secondary adapters (db/cache/crypto).
- `platform/`
  - Cross-cutting runtime wiring: config, server, observability.
- `shared/`
  - Shared constants, errors, and utilities.

## Flow (Where it comes from -> Where it goes)

Transport (HTTP/GraphQL) -> adapter/controller -> core/usecase -> output ports -> adapters

## Why It Was Designed This Way

- Preserve dependency inversion and domain purity.
- Keep contexts isolated and easy to evolve.
- Centralize cross-cutting concerns in `platform`.

## Recommended Practices Visible Here

- Usecases depend only on input/output ports.
- Adapters map DTOs and handle transport concerns.
- Semantic errors and OTel spans at all boundaries.

## Differentials

- Strict bounded-context isolation.
- Central GraphQL routing with context controllers.

## What Should NOT Live Here

- Infrastructure assets (Docker, migrations).
- Cross-context imports or shared state.
- Business rules inside adapters.
