# internal/admin/adapter/secondary/db/repository

Admin repository implementations for database access.

## Purpose & Main Capabilities

- Execute admin-specific queries and mutations.
- Translate driver errors to semantic errors.

## Package Composition

- Repository structs and methods for admin storage.

## Flow (Where it comes from -> Where it goes)

Usecase -> repository -> DB driver

## Why It Was Designed This Way

- Keep data access logic isolated and testable.
- Align repository behavior with admin usecases.

## Recommended Practices Visible Here

- Keep queries scoped to admin needs.
- Add observability around DB calls.

## Differentials

- Admin-specific repositories isolated from general user repos.

## What Should NOT Live Here

- Business rules or transport DTOs.
