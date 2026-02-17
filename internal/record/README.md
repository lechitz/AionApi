# Record Bounded Context

**Path:** `internal/record`

## Overview

Record domain for user event/diary entries.
Handles record lifecycle and query flows with strict user-scoped validation.

## Typical Responsibilities

| Area | Responsibility |
| --- | --- |
| Record lifecycle | Create/read/update/soft-delete records |
| Query operations | Date/tag/category/user-based retrieval |
| Domain validation | Ensure required relationships and field constraints |
| Dashboard semantics | Metric definitions, goal templates, and white-label layout (views/widgets) |

## Design Notes

- Keep record invariants in core usecases.
- Keep adapters mapping-only.
- Keep persistence concerns in secondary adapters.
- Dashboard white-label layout rules are enforced in usecase layer (including large-card limits).

## Package Improvements

- Add query operation matrix with expected filters/pagination behavior.
- Add tests for date/timezone boundary conditions.
- Add relation consistency checks for category/tag references.
- Add explicit notes for soft-delete semantics and recovery expectations.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
