# internal/admin/adapter/secondary/db

Database implementations of admin repositories.

## Purpose & Main Capabilities

- Persist and query admin-related data.
- Provide efficient access for role and status updates.

## Package Composition

- Repository implementations and mappers for admin entities.

## Flow (Where it comes from -> Where it goes)

Usecase -> admin repository -> db adapter -> database

## Why It Was Designed This Way

- Keep SQL and persistence details out of core.
- Align admin data access with migrations.

## Recommended Practices Visible Here

- Translate driver errors into semantic admin errors.
- Use safe metadata in db spans/logs.

## Differentials

- Admin-focused repository operations separated from general user repos.

## What Should NOT Live Here

- Business rules or HTTP details.
