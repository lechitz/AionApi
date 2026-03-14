# HTTP Cookie Utilities

**Path:** `internal/platform/server/http/utils/cookies`

## Overview

This package centralizes authentication cookie behavior for HTTP adapters.
It defines how access/refresh cookies are created, cleared, and extracted from incoming requests using platform configuration.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Access cookie lifecycle | Set and clear auth token cookie (`commonkeys.AuthTokenCookieName`) |
| Refresh cookie lifecycle | Set and clear refresh token cookie (`refresh_token`) |
| Token extraction | Read access/refresh values from request cookies |
| Security flags | Apply `HttpOnly`, `Secure`, `SameSite`, and cookie path/domain from config |

## Files

| File | Purpose |
| --- | --- |
| `cookies.go` | Cookie set/clear/extract helpers and SameSite/Secure mapping |

## Public API Reference

### Cookie Write Helpers

| Function | Behavior |
| --- | --- |
| `SetAuthCookie(w, token, cfg)` | Writes auth cookie using configured path/domain/samesite/max-age |
| `ClearAuthCookie(w, cfg)` | Expires auth cookie immediately (`MaxAge=-1`, `Expires=Unix(0,0)`) |
| `SetRefreshCookie(w, token, cfg)` | Writes refresh cookie with `MaxAge = cfg.MaxAge * 7` |
| `ClearRefreshCookie(w, cfg)` | Expires refresh cookie immediately |

### Cookie Read Helpers

| Function | Returns |
| --- | --- |
| `ExtractAuthToken(r)` | Access token value from auth cookie |
| `ExtractRefreshToken(r)` | Refresh token value from `refresh_token` cookie |

### Internal Mapping Helpers

| Helper | Purpose |
| --- | --- |
| `mapSameSite(sameSite string)` | Maps `Strict`/`Lax`/`None` to `http.SameSite*` |
| `secureFlag(cfg)` | Resolves secure flag from `config.CookieConfig` |

## Cookie Policy Summary

| Property | Auth cookie | Refresh cookie |
| --- | --- | --- |
| `HttpOnly` | `true` | `true` |
| `Secure` | Config-driven | Config-driven |
| `SameSite` | Config-driven on set, `Strict` on clear | Config-driven on set, `Strict` on clear |
| `Path` | `cfg.Path` | `cfg.Path` |
| `Domain` | `cfg.Domain` | `cfg.Domain` |
| `MaxAge` | `cfg.MaxAge` | `cfg.MaxAge * 7` |

## Usage Example

```go
cookies.SetAuthCookie(w, accessToken, cfg.Cookie)
cookies.SetRefreshCookie(w, refreshToken, cfg.Cookie)

// Later in request flow:
authToken, err := cookies.ExtractAuthToken(r)
if err != nil {
    // handle missing/invalid auth cookie
}
_ = authToken
```

## Design Notes

- Keep cookie mechanics in one package to avoid duplicated security options in handlers.
- Keep token naming stable to preserve client compatibility.
- Use config-driven secure behavior for local development and production parity.

## Package Improvements

- Replace hardcoded refresh cookie name (`"refresh_token"`) with a shared constant in `internal/shared/constants/commonkeys`.
- Add focused unit tests for `mapSameSite`, extraction edge cases, and cookie attributes written by each helper.
- Consider explicit validation/log-free handling when `SameSite=None` and `Secure=false` to avoid invalid browser behavior.
- Evaluate configurability for refresh cookie lifetime instead of fixed multiplier (`cfg.MaxAge * 7`).

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
