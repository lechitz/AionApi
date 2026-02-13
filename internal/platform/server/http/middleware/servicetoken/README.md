# HTTP Service Token Middleware

**Path:** `internal/platform/server/http/middleware/servicetoken`

## Overview

This package provides service-to-service (S2S) authentication through HTTP headers.
It validates a trusted service key and can optionally inject a user identity into request context.

## Package Scope

| Area | Responsibility |
| --- | --- |
| S2S key validation | Validate `X-Service-Key` against configured service key |
| Optional user impersonation | Parse `X-Service-User-Id` and inject into context when valid |
| Context enrichment | Set `ctxkeys.ServiceAccount` and optional `ctxkeys.UserID` |
| Fast rejection | Return `401 Unauthorized` on invalid/misconfigured service key |

## Files

| File | Purpose |
| --- | --- |
| `service_token.go` | Middleware implementation, header constants, and logging messages |

## Public API Reference

| Function / Constant | Description |
| --- | --- |
| `New(cfg, log)` | Returns middleware that validates S2S calls when service key header is present |
| `HeaderServiceKey` | Header name for service credential: `X-Service-Key` |
| `HeaderServiceUser` | Optional user ID header: `X-Service-User-Id` |
| `ErrServiceTokenInvalid` | Error message for unauthorized S2S attempts |

## Runtime Behavior

1. Read `X-Service-Key` from request.
2. If header is absent, request is passed through unchanged.
3. If header exists but configured service key is empty, respond `401`.
4. If header exists and does not match configured key, respond `401`.
5. On success, set `ctxkeys.ServiceAccount = true`.
6. If `X-Service-User-Id` is present and valid `uint64`, set `ctxkeys.UserID`.
7. Continue request with enriched context.

## Headers

| Header | Required | Description |
| --- | --- | --- |
| `X-Service-Key` | Required for S2S mode | Shared secret to authenticate trusted service calls |
| `X-Service-User-Id` | Optional | User context injected into request when valid |

## Configuration

| Config key | Description |
| --- | --- |
| `cfg.AionChat.ServiceKey` | Expected value for `X-Service-Key` |

Example:

```env
AION_CHAT_SERVICE_KEY=your-secret-key-here
```

## Usage

```go
r.GroupWith(servicetoken.New(cfg, log), func(sr ports.Router) {
    sr.Mount(cfg.ServerGraphql.Path, graphqlHandler)
})
```

## Context Contract

| Context key | Type | Set when |
| --- | --- | --- |
| `ctxkeys.ServiceAccount` | `bool` | Valid service key was provided |
| `ctxkeys.UserID` | `uint64` | Optional service user header is present and parseable |

## Design Notes

- Middleware is permissive by default when no service key header is provided.
- S2S authentication is transport-level behavior; domain authorization remains outside this package.
- Shared constants are used for log keys and context keys to avoid string drift.

## Package Improvements

- Add nil-safe logging guards (`if log != nil`) to avoid panic risk when middleware is initialized without a logger.
- Consider replacing `http.Error` plain-text body with centralized JSON error writer for response consistency.
- Add unit tests for: missing key pass-through, empty configured key, invalid key, valid key with/without valid user ID.
- Evaluate timing-safe comparison for service key validation if threat model requires stronger side-channel resistance.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
