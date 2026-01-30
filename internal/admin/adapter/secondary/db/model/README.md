# internal/admin/adapter/secondary/db/model

Database models used by admin repositories.

## Purpose & Main Capabilities

- Define persistence shapes for admin-related data.

## Package Composition

- DB structs aligned with migrations.

## Flow (Where it comes from -> Where it goes)

Mapper -> DB model -> repository -> database

## Why It Was Designed This Way

- Keep DB schemas explicit and localized.
- Separate storage concerns from domain.

## Recommended Practices Visible Here

- Keep models aligned with migrations.
- Avoid embedding business rules in models.

## Differentials

- Clear separation of admin persistence types.

## What Should NOT Live Here

- Domain logic or transport DTOs.
