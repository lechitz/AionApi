# HTTP Utils — Cookies

**Folder:** `internal/platform/server/http/utils/cookies`

## Purpose and Main Capabilities

- Standardize auth/session cookie creation and clearing.
- Centralize security flags and TTL based on platform config.
- Avoid duplicated cookie logic across handlers.

## How it works

- Helpers read cookie settings from `internal/platform/config`.
- Cookie names reuse constants (e.g., `commonkeys.AuthTokenCookieName`).
- Writes `Set-Cookie` with `HttpOnly`, `Secure`, `SameSite`, and explicit TTL.

## Usage (pattern)

```go
// cookies.SetAuthCookie(w, token, cfg)
// cookies.ClearAuthCookie(w, cfg)
```

## Conventions

- Never log raw cookie values.
- Avoid wildcard origins when using credentials.

## What Should NOT Live Here

- Domain logic or handler orchestration.
