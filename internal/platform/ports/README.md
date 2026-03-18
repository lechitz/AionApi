# Platform Output Ports

**Path:** `internal/platform/ports`

## Purpose

This folder defines shared output-port contracts reused across the application.

## Current Output Ports

| Package | Role |
| --- | --- |
| `output/cache` | cache operations and lifecycle |
| `output/db` | database abstraction and transaction support |
| `output/hasher` | password hashing |
| `output/httpclient` | outbound HTTP requests with instrumentation |
| `output/keygen` | key generation helpers |
| `output/logger` | structured logger contract |

## Boundaries

- keep interfaces small and technology-agnostic
- ports exist to decouple core logic from vendor implementations
- adding a new port here means the contract is shared across multiple boundaries, not just one context

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
