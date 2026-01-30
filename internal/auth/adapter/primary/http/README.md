# internal/auth/adapter/primary/http

HTTP transport for auth endpoints and middleware.

## Purpose & Main Capabilities

- Expose login, logout, and refresh endpoints.
- Protect routes with auth middleware.
- Map HTTP DTOs to core commands.

## Package Composition

- `handler/`
  - Auth HTTP handlers.
- `middleware/`
  - Auth middleware for protected routes.
- `dto/`
  - Request/response DTOs.

## Flow (Where it comes from -> Where it goes)

HTTP request -> handler/middleware -> core/usecase -> HTTP response

## Why It Was Designed This Way

- Centralize auth transport behavior.
- Keep core logic transport-agnostic.

## Recommended Practices Visible Here

- Sanitize bearer tokens and support cookie fallback.
- Keep DTO mapping inside adapters.

## Differentials

- Dual token source support (header + cookie).

## What Should NOT Live Here

- Business rules or persistence code.
