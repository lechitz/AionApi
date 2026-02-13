# Primary Adapters Layer

**Path:** `internal/adapter/primary`

## Overview

Primary adapters expose usecases to external clients (GraphQL/HTTP).
They convert transport payloads into core input contracts and map responses back to transport formats.

## Subpackages

| Subpackage | Role |
| --- | --- |
| `graphql/` | Central GraphQL transport entrypoint, schema composition, directives, resolvers |

## Design Notes

- Keep resolvers/handlers thin and orchestration-only.
- Business rules must remain in context core/usecases.
- Shared transport conventions live here to avoid drift across contexts.

## Package Improvements

- Add short “transport boundary checklist” for new adapters.
- Add examples for error/status mapping conventions.
- Add explicit link map to context controllers.
- Add contract test guidance for primary adapter behavior.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
