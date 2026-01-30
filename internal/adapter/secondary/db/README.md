# internal/adapter/secondary/db

Database adapters implementing persistence output ports for domain aggregates.

## Package Composition

- `postgres/`
  - Postgres-specific repository implementations.
- Mappers/models used for persistence translation.

## Flow (Where it comes from -> Where it goes)

Usecase -> repository port -> db adapter -> database

## Why It Was Designed This Way

- Separate domain entities from persistence models.
- Centralize transaction and error translation.

## Recommended Practices Visible Here

- Keep repositories thin; move mapping to mappers/models.
- Translate driver errors into semantic errors.
- Add spans/logs for db operations with safe metadata.

## Differentials

- Clear separation between domain entities and DB models.

## What Should NOT Live Here

- Business rules or transport DTOs.
