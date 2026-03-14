# Platform HTTP Middleware Layer

**Path:** `internal/platform/server/http/middleware`

## Overview

This package groups platform-level HTTP middlewares that run before context handlers.
It centralizes cross-cutting transport concerns such as correlation, panic safety, CORS policy, and optional service-to-service authentication.

## Subpackages

| Subpackage | Primary responsibility |
| --- | --- |
| `cors/` | Browser cross-origin policy and credentials support |
| `recovery/` | Panic interception and delegation to generic recovery handler |
| `requestid/` | Request correlation ID normalization and propagation |
| `servicetoken/` | Optional S2S header-based authentication and context enrichment |

## Middleware Chain Intent

| Concern | Why it exists |
| --- | --- |
| Reliability | Prevent panics from crashing request execution |
| Traceability | Guarantee request-level correlation key across logs/traces/headers |
| Security boundaries | Enforce transport-level policies before adapter logic |
| Consistency | Apply shared behavior once, not per handler |

## Recommended Ordering

1. `recovery.New(...)`
2. `requestid.New()`
3. `cors.New()`
4. `servicetoken.New(cfg, log)` (scoped to routes that need S2S)

## Example Wiring

```go
r.Use(recovery.New(genericHandler))
r.Use(requestid.New())
r.Use(cors.New())

r.GroupWith(servicetoken.New(cfg, log), func(sr ports.Router) {
    sr.Mount(cfg.ServerGraphql.Path, graphqlHandler)
})
```

## Design Notes

- Middlewares in this layer must stay transport/platform-focused.
- Business rules and domain authorization must remain in context adapters/usecases.
- Subpackage READMEs contain implementation-level details; this README is intentionally high-level.

## Package Improvements

- Add an integration test matrix that validates middleware ordering and combined behavior (panic + request-id + error response).
- Document which routes are intentionally outside `servicetoken` scope to avoid security ambiguity.
- Consider a single composer example showing global vs route-scoped middleware in `internal/platform/server/http` docs.
- Add a short troubleshooting section for common local issues (CORS origin mismatch, missing service key, missing request ID).

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
