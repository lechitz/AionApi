# Platform HTTP Utils

**Folder:** `internal/platform/server/http/utils`

## Purpose and Main Capabilities

- Provide shared HTTP helpers for responses, errors, and cookies.
- Keep handlers thin and consistent across contexts.
- Centralize response envelopes and error mapping.

## Package Composition

- `httpresponse/`: JSON response writers and error mapping.
- `sharederrors/`: shared semantic error helpers for HTTP mapping.
- `cookies/`: auth/session cookie helpers driven by config.

## Flow (Where it comes from -> Where it goes)

Handler -> utils (httpresponse/sharederrors/cookies) -> HTTP response

## Recommended Practices Visible Here

- Use semantic errors (`sharederrors`) and let `httpresponse` map to HTTP codes.
- Avoid leaking sensitive data in logs or responses.
- Keep cookie flags centralized and config-driven.

## What Should NOT Live Here

- Domain rules or persistence logic.
- Context-specific DTOs or handlers.
