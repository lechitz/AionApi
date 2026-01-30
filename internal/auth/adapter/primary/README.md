# internal/auth/adapter/primary

Primary adapters for auth entrypoints (HTTP handlers and middleware).

## Purpose & Main Capabilities

- Handle login/logout/refresh HTTP endpoints.
- Protect routes via auth middleware.
- Map errors to HTTP responses.

## Package Composition

- `http/`
  - Handlers, middleware, and DTOs.

## Flow (Where it comes from -> Where it goes)

HTTP request -> handler/middleware -> core/usecase

## Why It Was Designed This Way

- Keep auth transport logic centralized.
- Enforce security checks at the boundary.

## Recommended Practices Visible Here

- Validate inputs before calling core.
- Inject user/claims into context for downstream handlers.

## Differentials

- Dedicated middleware for auth with context propagation.

## What Should NOT Live Here

- Persistence or token generation logic.
