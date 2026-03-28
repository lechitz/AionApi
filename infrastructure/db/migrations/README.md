# Database Migrations

**Path:** `infrastructure/db/migrations`

## Purpose

This folder is the canonical schema history for `aion-api`, managed through `golang-migrate`.
Each migration pair must preserve the ordered evolution of the live database contract.

## Current Shape

| Artifact | Responsibility |
| --- | --- |
| `000001` - `000019` | current ordered migration chain |
| `*.up.sql` | forward schema changes |
| `*.down.sql` | rollback path for the matching migration |

## Common Commands

```bash
make migrate-dev-up
make migrate-dev-down
make migrate-dev-status
make migrate-dev-reset
make migrate-new
```

## Boundaries

- never rewrite an already-applied migration in shared environments
- keep forward and rollback files paired
- application code and seed scripts must adapt to the migration chain, not bypass it

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
