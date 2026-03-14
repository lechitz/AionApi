# Tag Bounded Context

**Path:** `internal/tag`

## Overview

Tag domain for user-scoped labeling tied to categories and records.
Provides tag lifecycle operations with uniqueness and ownership constraints.

## Typical Responsibilities

| Area | Responsibility |
| --- | --- |
| Tag lifecycle | Create/read/update/soft-delete tags |
| Domain constraints | Enforce per-user uniqueness and category relation |
| Adapter exposure | Provide GraphQL/HTTP-facing controller operations |

## Design Notes

- Keep validation and uniqueness in core.
- Keep resolver/controller mapping separate from domain logic.
- Keep semantic error contracts stable for transport mapping.

## Package Improvements

- Add uniqueness conflict behavior table.
- Add tests for category reassignment/update edge cases.
- Add guidance for icon/metadata normalization.
- Add notes for tag usage impact on record queries.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
