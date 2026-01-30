# Platform HTTP Middleware

**Folder:** `internal/platform/server/http/middleware`

## Purpose and Main Capabilities

- Centralize cross-cutting HTTP behavior (request ID, recovery, CORS, service token).
- Keep middleware consistent and framework-agnostic.
- Support observability and safe error handling across all routes.

## Package Composition

- `requestid/`: ensures `X-Request-ID` header and context value.
- `recovery/`: catches panics and delegates to generic recovery handler.
- `cors/`: configures CORS for allowed origins/methods/headers.
- `servicetoken/`: S2S authentication via `X-Service-Key`.

## Flow (Where it comes from -> Where it goes)

HTTP request -> middleware chain -> handler -> response

## Recommended Practices Visible Here

- Apply `recovery` as the outermost guard.
- Apply `requestid` early to enable log/trace correlation.
- Keep middleware stateless and transport-only.
- Avoid leaking internal errors to clients.

## What Should NOT Live Here

- Business logic or domain rules.
- Context-specific authorization (belongs in domain adapters).
