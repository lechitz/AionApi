# Internal Application Layer

**Path:** `internal`

## Overview

Main application code organized by bounded contexts and cross-cutting platform/shared layers.
Architecture follows ports-and-adapters with clear dependency direction.

## Main Areas

| Area | Role |
| --- | --- |
| Context modules (`admin`, `auth`, `category`, `chat`, `record`, `tag`, `user`) | Domain-specific core + adapters |
| `adapter/` | Shared primary/secondary adapter infrastructure |
| `platform/` | Runtime configuration, DI, observability, server composition |
| `shared/` | Cross-cutting constants and lightweight shared contracts |

## Design Notes

- Preserve context isolation and dependency rule.
- Keep domain logic in core/usecase layers.
- Keep transport and infrastructure concerns in adapters/platform.

## Package Improvements

- Add architecture map linking each context to its exposed adapters.
- Add dependency direction examples with allowed/forbidden imports.
- Add checklist for adding a new bounded context.
- Add quick references to key subpackage READMEs.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
