# internal/adapter

Shared adapters used across multiple contexts. This layer concentrates transport and infrastructure patterns to avoid duplication.

## Package Composition

- `primary/`
  - Shared entrypoint adapters (GraphQL infrastructure).
- `secondary/`
  - Shared output adapters (db, cache, crypto, tokens, logging, graphs).

## Flow (Where it comes from -> Where it goes)

External input -> primary adapters -> controllers/usecases

Usecases -> output ports -> secondary adapters -> external systems

## Why It Was Designed This Way

- Centralize cross-context adapter behavior.
- Keep domain contexts focused on business logic.
- Standardize error handling and observability.

## Recommended Practices Visible Here

- Adapters depend on ports; core never depends on adapters.
- Keep transport conventions centralized to avoid drift.
- Include tracing/logging without leaking sensitive data.

## Differentials

- Shared adapter patterns that enforce consistency.

## What Should NOT Live Here

- Domain rules or business decisions.
- Cross-context imports inside core.
- Infrastructure assets (Docker, migrations).
