# Audit Context

**Path:** `internal/audit`

## Purpose

`internal/audit` owns immutable action-event persistence and the internal diagnostics query surface for those events.

## Current Surface

| Surface | Current contract |
| --- | --- |
| `core/ports/input.Service.WriteEvent` | validates and persists one immutable audit event |
| `core/ports/input.Service.ListEvents` | returns filtered audit events for diagnostics |
| HTTP `GET /audit/events` | authenticated diagnostics endpoint; self-scope by default, `user_id` cross-user queries only for admin callers |
| Storage | `aion_api.audit_action_events` |

## Current Producers

- `chat` emits UI-action audit events through `WriteEvent`
- other contexts may publish audit events, but they must treat persistence as non-blocking from a business-outcome perspective

## Invariants

- Audit persistence must not change producer business outcomes.
- Payloads remain allow-listed and redacted.
- Events are append-only in application behavior.
- Query semantics are diagnostic, not product-facing.

## Boundaries

- `audit` owns the immutable event log, not the business workflows that generated it.
- The read surface is HTTP-only in the current runtime.
- Consumers should not depend on audit writes for transactional guarantees.

## Related Docs

- [`../chat/README.md`](../chat/README.md)
- [`../platform/server/http/README.md`](../platform/server/http/README.md)

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
