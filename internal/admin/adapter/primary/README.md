# internal/admin/adapter/primary

Primary adapters for admin entrypoints (HTTP).

## Purpose & Main Capabilities

- Validate and authorize admin requests.
- Map HTTP payloads to admin commands.
- Return transport-safe responses and errors.

## Package Composition

- `http/`
  - Admin HTTP handlers and routing helpers.

## Flow (Where it comes from -> Where it goes)

HTTP request -> admin handler -> core/usecase

## Why It Was Designed This Way

- Keep admin transport concerns centralized.
- Enforce policy checks at the boundary.

## Recommended Practices Visible Here

- Apply strict auth/roles validation before calling core.
- Keep handlers thin; let usecases own rules.

## Differentials

- Dedicated admin HTTP boundary.

## What Should NOT Live Here

- Business rules or persistence logic.
