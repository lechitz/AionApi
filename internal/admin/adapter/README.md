# internal/admin/adapter

Admin adapters for transport and persistence edges.

## Purpose & Main Capabilities

- Translate admin HTTP requests into domain commands.
- Map domain results to transport responses.
- Connect admin usecases to database implementations.

## Package Composition

- `primary/`
  - HTTP handlers for admin endpoints.
- `secondary/`
  - DB repositories for admin operations.

## Flow (Where it comes from -> Where it goes)

Admin request -> primary adapter -> core/usecase -> secondary adapter

## Why It Was Designed This Way

- Keep transport and persistence out of core logic.
- Allow infra changes without touching admin rules.

## Recommended Practices Visible Here

- Keep handlers thin and policy-aware.
- Translate infra errors to semantic errors in core.

## Differentials

- Admin adapters isolated from user-facing adapters.

## What Should NOT Live Here

- Business rules or role policy definitions.
