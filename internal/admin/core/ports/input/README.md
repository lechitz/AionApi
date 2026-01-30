# internal/admin/core/ports/input

Admin input port contracts.

## Purpose & Main Capabilities

- Define the admin service interface consumed by adapters.

## Package Composition

- Admin input port interfaces.

## Flow (Where it comes from -> Where it goes)

Adapter -> input port -> usecase

## Why It Was Designed This Way

- Keep adapters decoupled from concrete usecases.

## Recommended Practices Visible Here

- Keep method signatures context-first.
- Use domain types and semantic errors.

## Differentials

- Stable admin interface for multiple entrypoints.

## What Should NOT Live Here

- Implementations or transport DTOs.
