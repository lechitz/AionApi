# Category Bounded Context

**Path:** `internal/category`

## Overview

Category management domain for user-scoped classification entities.
Supports CRUD-style operations and transport exposure through primary adapters.

## Typical Responsibilities

| Area | Responsibility |
| --- | --- |
| Category lifecycle | Create/read/update/soft-delete category entities |
| Validation | Enforce required fields and user scoping |
| Adapter integration | Expose operations through GraphQL/HTTP controllers |

## Design Notes

- Keep invariants and uniqueness checks in core.
- Keep transport mapping outside usecases.
- Preserve semantic error usage for adapter mapping.

## Package Improvements

- Add operation table with expected semantic errors.
- Add tests for update/soft-delete edge cases.
- Add explicit field normalization guidelines.
- Add interaction notes with `tag`/`record` relations.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
