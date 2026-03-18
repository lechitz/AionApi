# Database Infrastructure

**Path:** `infrastructure/db`

## Purpose

`infrastructure/db` owns the Postgres lifecycle for `AionApi`.
It splits durable schema evolution from local and QA seed data.

## Current Areas

| Area | Responsibility |
| --- | --- |
| `migrations/` | ordered schema evolution for auth, records, dashboard, audit, and event-outbox tables |
| `seed/` | direct-SQL bootstrap data for local development and deterministic test baselines |

## Operational Use

```bash
make migrate-dev-up
make migrate-dev-status
make db-full
make seed-test
```

## Boundaries

- schema changes must land as migrations, not ad-hoc SQL edits
- seed data may assume a dev/test environment, but it must not redefine schema ownership
- API-based seed callers live under `hack/tools`; this folder owns the direct SQL side of local data provisioning

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
