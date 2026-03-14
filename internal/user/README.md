# User Bounded Context

**Path:** `internal/user`

## Overview

User domain for account lifecycle, profile updates, password changes, and soft deletion.
It centralizes identity-related business rules and security-sensitive flows.

## Typical Responsibilities

| Area | Responsibility |
| --- | --- |
| Account lifecycle | Create/read/update/delete user operations |
| Credential management | Password validation/hash/update flows |
| Identity constraints | Username/email uniqueness and normalization |
| Session interplay | Token/store coordination for sensitive user changes |

## Design Notes

- Keep user rules and security checks in core usecases.
- Keep transport concerns in primary adapters.
- Keep storage/hash/token integrations behind output ports.

## Package Improvements

- Add end-to-end flow diagrams for create/update-password/delete.
- Add tests for race/conflict scenarios on unique fields.
- Add clear policy for PII-safe logging at boundaries.
- Add compatibility notes for auth/session invalidation behavior.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
