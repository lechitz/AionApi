# Secondary Adapters Layer

**Path:** `internal/adapter/secondary`

## Purpose

This layer holds shared outbound implementations reused across the application.

## Current Packages

| Package | Role |
| --- | --- |
| `contextlogger/` | structured logger implementation |
| `crypto/` | key generation helpers |
| `hasher/` | bcrypt password hashing |
| `token/` | JWT token provider and expiry helpers |

## Boundaries

- context-owned DB, cache, HTTP, Kafka, and storage adapters stay inside the owning bounded context
- packages here implement reusable platform or security concerns shared across contexts
- business semantics must remain in core/usecase layers and input/output ports

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
