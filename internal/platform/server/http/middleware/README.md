# Platform HTTP Middleware Layer

**Path:** `internal/platform/server/http/middleware`

## Purpose

This package owns cross-cutting HTTP middleware applied before bounded-context handlers.

## Current Runtime Usage

| Middleware | Current scope | What it does |
| --- | --- | --- |
| `requestid.New()` | global router and health endpoints | normalizes or generates `X-Request-ID`, stores it in context, echoes it in the response |
| `recovery.New(genericHandler)` | global router only | catches panics and delegates sanitized response handling to the generic handler |
| `cors.New()` | global router and health endpoints | applies the current browser-origin and credential policy |
| `servicetoken.New(cfg, log)` | GraphQL mount only | validates trusted S2S headers and optionally injects service-account user context |

## Effective Order In `composer.go`

1. `requestid.New()`
2. `recovery.New(genericHandler)`
3. `cors.New()`
4. `servicetoken.New(cfg, log)` only around the GraphQL mount

This order is the current truth and takes precedence over any older “recommended ordering” text elsewhere.

## Current Policy Notes

- CORS currently allows `http://localhost:5000` and `http://localhost:5173`, with credentials enabled.
- Request IDs must be UUIDs; invalid or oversized values are replaced.
- Health routes intentionally skip the global recovery and `otelhttp` wrappers.
- `servicetoken` is permissive when no service key header is present and becomes authoritative only when S2S headers are supplied.

## Boundary Rule

- This README is the current owner for `cors`, `recovery`, and `requestid` behavior.
- Those leaf READMEs are intentionally removed to reduce drift.
- `servicetoken` keeps its own README because it is a distinct trust boundary.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
