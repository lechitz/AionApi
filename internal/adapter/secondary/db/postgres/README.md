# internal/adapter/secondary/db/postgres

Postgres-specific implementations of persistence ports.

## Package Composition

- Repository implementations and SQL queries.
- Mapping between domain entities and DB records.

## Flow (Where it comes from -> Where it goes)

Usecase -> repository port -> Postgres adapter -> Postgres

## Why It Was Designed This Way

- Keep SQL and connection details isolated.
- Align storage with migrations and schema evolution.

## Recommended Practices Visible Here

- Optimize queries using indexed columns from migrations.
- Document deviations from simple CRUD.
- Use integration tests for complex queries.

## Differentials

- Consistent Postgres mapping layer.

## What Should NOT Live Here

- Domain rules or API logic.
