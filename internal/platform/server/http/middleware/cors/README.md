# HTTP CORS Middleware

**Path:** `internal/platform/server/http/middleware/cors`

## Overview

This package provides the platform-level CORS middleware used by the HTTP server.
It defines which frontend origins can call the API and enables cookie-based cross-origin requests.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Origin policy | Allow only known frontend origins |
| Method policy | Restrict accepted HTTP methods |
| Header policy | Restrict allowed request headers |
| Credential support | Enable cookies/credentials for browser requests |

## Files

| File | Purpose |
| --- | --- |
| `cors_middleware.go` | Exposes `New()` middleware configured with `go-chi/cors` options |

## Public API Reference

| Function | Returns | Description |
| --- | --- | --- |
| `New()` | `func(http.Handler) http.Handler` | CORS middleware preconfigured for Aion API frontend integration |

## Current CORS Policy

| Setting | Value |
| --- | --- |
| `AllowedOrigins` | `http://localhost:5000`, `http://localhost:5173` |
| `AllowedMethods` | `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS` |
| `AllowedHeaders` | `Accept`, `Authorization`, `Content-Type`, `X-CSRF-Token` |
| `ExposedHeaders` | `Set-Cookie` |
| `AllowCredentials` | `true` |
| `MaxAge` | `300` seconds |

## Usage

```go
r.Use(cors.New())
```

## Design Notes

- Keep CORS policy centralized to avoid divergence across routes.
- Use explicit origins because `AllowCredentials` is enabled.
- Keep transport policy here; do not mix authorization/business rules into middleware configuration.

## Package Improvements

- Move allowed origins to platform configuration (`env`) to support multiple environments without code changes.
- Add unit tests for middleware behavior (`OPTIONS` preflight, allowed origin, blocked origin).
- Evaluate whether exposing `Set-Cookie` is still required for all consumers or should be narrowed.
- Consider explicit local constants for CORS values to simplify future policy reviews.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
