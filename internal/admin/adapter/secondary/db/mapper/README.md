# internal/admin/adapter/secondary/db/mapper

Mapping helpers between domain entities and DB models for admin persistence.

## Purpose & Main Capabilities

- Convert admin domain entities to DB records and back.

## Package Composition

- Mapper functions for admin entities.

## Flow (Where it comes from -> Where it goes)

Domain entity -> mapper -> DB model -> database

## Why It Was Designed This Way

- Keep mapping logic separate from repositories.
- Avoid leaking persistence details into core.

## Recommended Practices Visible Here

- Keep mappings explicit and testable.
- Preserve domain invariants during conversion.

## Differentials

- Dedicated mapping layer for admin persistence.

## What Should NOT Live Here

- Business rules or query logic.
