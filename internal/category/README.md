# Category Bounded Context

**Path:** `internal/category`

## Purpose

`internal/category` owns user-scoped category lifecycle and lookup flows.

## Current Surface

Current transport exposure is GraphQL-driven through the category controller surface:

- create category
- update category
- soft-delete category
- get by id
- get by name
- list all categories for the authenticated user

There is no dedicated REST category surface in the current runtime.

## Runtime Contract

- core usecases own create/read/update/soft-delete semantics
- cache adapters support lookup by id, by name, and list results
- DB adapters remain the authority for persistence and ownership checks
- domain output carries `usageCount` and `lastUsedAt`, which are consumed by higher-level product surfaces

## Boundaries

- Category rules stay user-scoped and must not leak cross-user data.
- Relationship semantics with tags and records belong here and in the collaborating repositories/adapters, not in transport code.
- Transport controllers only map GraphQL types and user context into category commands.

## Related Docs

- [`../tag/README.md`](../tag/README.md)
- [`../record/README.md`](../record/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
