# internal/admin/adapter/secondary

Secondary adapters for admin persistence and external integrations.

## Purpose & Main Capabilities

- Implement admin repositories against concrete storage.
- Translate database errors into semantic errors.

## Package Composition

- `db/`
  - Database-backed repositories for admin usecases.

## Flow (Where it comes from -> Where it goes)

Usecase -> repository port -> db adapter -> database

## Why It Was Designed This Way

- Keep storage details isolated from admin logic.
- Allow swaps or refactors without core changes.

## Recommended Practices Visible Here

- Keep repository methods focused on admin needs.
- Add observability to db boundaries.

## Differentials

- Admin-specific repositories separated from user-facing repos.

## What Should NOT Live Here

- Business rules or transport DTOs.
