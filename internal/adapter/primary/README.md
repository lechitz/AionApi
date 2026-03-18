# Primary Adapters Layer

**Path:** `internal/adapter/primary`

## Purpose

This layer owns shared inbound transport infrastructure.
Today that means the central GraphQL boundary consumed across multiple contexts.

## Current Surface

| Subpackage | Role |
| --- | --- |
| `graphql/` | gqlgen config, schema composition, shared resolvers, directives, and server bootstrap |

## Boundaries

- most REST handlers remain in the owning bounded context under `internal/<ctx>/adapter/primary/http`
- shared primary adapter code belongs here only when it coordinates multiple contexts
- resolvers and transport glue stay orchestration-only; use-case behavior remains in context services

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
