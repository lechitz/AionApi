# Platform HTTP Router Adapters

**Path:** `internal/platform/server/http/router`

## Purpose

This package contains concrete router adapters behind the framework-agnostic `ports.Router` contract.

## Current Adapter

`chi/` is the active implementation.

It currently supports:

- global middleware via `Use(...)`
- grouped routes via `Group(...)` and `GroupWith(...)`
- handler mounts via `Mount(...)`
- method helpers for `GET`, `POST`, `PUT`, and `DELETE`
- custom `NotFound` and `MethodNotAllowed` handlers
- optional `SetError(...)` callback storage

## Important Nuance

`SetError(...)` does not cause chi to invoke a centralized error callback automatically.
The current adapter only stores the callback so the port surface stays consistent with platform expectations.

In practice:

- panic handling is done by middleware
- 404/405 behavior is done by explicit fallback handlers
- transport error mapping is done in handlers/utilities

## Boundaries

- Bounded contexts must depend on `ports.Router`, never on `chi` directly.
- Any new adapter must preserve the same method and grouping semantics expected by current registrars.
- Business logic and response mapping do not belong here.

## Validate

```bash
go test ./internal/platform/server/http/router/...
```

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
