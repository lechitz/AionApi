# infrastructure/db

Database assets for the AionAPI Postgres schema. This folder contains schema migrations and seed data used for local development and testing.

## Package Composition

- `migrations/`
  - Versioned schema changes (up/down) managed by golang-migrate.
- `seed/`
  - Deterministic seed datasets for dev/test.

## Flow (Where it comes from -> Where it goes)

Developer -> migrations/seed -> database schema and data

## Why It Was Designed This Way

- Keep schema evolution reproducible.
- Separate schema changes from data seeding.
- Support fast local setup without API dependency.

## Recommended Practices Visible Here

- Never edit applied migrations; add a new one.
- Keep names aligned with DB adapters in `internal/*/adapter/secondary/db`.
- Keep seeds small, safe, and representative.
- Use env vars for parameters (no hardcoded secrets).

## Common Commands

```bash
make migrate-up
make migrate-new
make seed-all
```

## What Should NOT Live Here

- Business logic or API behavior.
- One-off SQL without versioning.
- Sensitive production data.
