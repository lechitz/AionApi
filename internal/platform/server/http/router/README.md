# Platform HTTP Router

**Folder:** `internal/platform/server/http/router`

## Purpose and Main Capabilities

- Hold concrete router adapters that implement `ports.Router`.
- Allow swapping the routing engine without touching context code.
- Keep routing details isolated from domain adapters.

## Package Composition

- `chi/`: chi-based implementation of `ports.Router` (current adapter).

## Flow (Where it comes from -> Where it goes)

HTTP composer -> router adapter -> middleware -> handlers

## Recommended Practices Visible Here

- Contexts depend only on `internal/platform/server/http/ports`.
- Use `Mount` for composite handlers (e.g., GraphQL).
- Keep router adapters thin and focused on wiring.
- Swap routers by changing the composer import, not context code.

## What Should NOT Live Here

- Context-specific route registration.
- Business logic or transport DTOs.
