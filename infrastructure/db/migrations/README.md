# Database Migrations

**Path:** `infrastructure/db/migrations`

## Overview

Versioned SQL migrations managed by `golang-migrate`.
This folder is the canonical source for schema evolution.

## Scope

| Artifact | Responsibility |
| --- | --- |
| `*.up.sql` | Forward schema changes |
| `*.down.sql` | Rollback path for matching forward migration |

## Workflow

1. Create migration pair.
2. Apply with migration tooling.
3. Validate schema/application compatibility.
4. Keep rollback path consistent.

## Common Commands

```bash
make migrate-up
make migrate-down
make migrate-new
make migrate-force VERSION=6
```

## Design Notes

- Keep migrations immutable after execution in shared environments.
- Prefer small, focused migration steps.
- Always include rollback strategy when feasible.

## Package Improvements

- Add CI migration smoke test against disposable Postgres.
- Add lint/check for missing down migration pairs.
- Add migration naming convention guide in this README.
- Add preflight checklist for destructive schema changes.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
