# Database Seed Scripts

**Path:** `infrastructure/db/seed`

## Purpose

This folder owns direct SQL seed artifacts for local development, QA, and deterministic reset flows.
It is optimized for fast Postgres bootstrap, not for production data loading.

## Current Files

| Pattern | Purpose |
| --- | --- |
| `roles.sql`, `user_roles.sql`, `admin_user.sql` | role and admin bootstrap |
| `*_generate.sql` | generated categories, tags, users, and records |
| `test_*.sql` | scenario-oriented datasets for timeline and demo coverage |
| `.env.example` | helper env template for local seed tooling |

## Common Flows

```bash
make seed-essential
make seed-test
make seed-clean-all
make db-full
```

`db-full` remains the fastest path to a realistic local profile with roles, admin, taxonomy, dashboard tables, and high-volume records.

## Boundaries

- keep seed data representative but disposable
- do not let seed scripts redefine schema ownership or hide migration gaps
- API-driven seed callers and synthetic generation flows belong under `hack/tools` and complement, rather than replace, these SQL assets

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
