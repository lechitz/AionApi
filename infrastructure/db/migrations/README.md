# infrastructure/db/migrations

Database migrations managed by golang-migrate. This folder is the source of truth for schema evolution.

## Package Composition

- `*.up.sql`
  - Forward migrations (apply changes).
- `*.down.sql`
  - Rollback migrations (revert changes).

## Flow (Where it comes from -> Where it goes)

Developer -> migration files -> migrate tool -> database schema

## Why It Was Designed This Way

- Keep schema changes versioned and reproducible.
- Enable forward and rollback paths for each change.
- Match production tooling (golang-migrate).

## Recommended Practices Visible Here

- Pair every `up` with a matching `down` migration.
- Use 6-digit ordered versions.
- Keep changes small and reversible.

## Commands

```bash
make migrate-up
make migrate-dev-up
make migrate-down
make migrate-new
make migrate-force VERSION=6
```

## Environment Variables

```bash
MIGRATION_DB=postgres://user:password@host:port/database?sslmode=disable
export MIGRATION_DB="postgres://aion:aion123@localhost:5432/aionapi?sslmode=disable"
```

## What Should NOT Live Here

- Seed data (use `infrastructure/db/seed`).
- Ad-hoc SQL without migrations.
- Irreversible changes without a clear rollback.
