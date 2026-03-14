# Audit Context

The `audit` bounded context owns immutable action-event persistence for operational and compliance diagnostics.

## Scope (V1)
- Write immutable audit events.
- Query audit events internally by trace/draft/user/status/time.
- Primary HTTP read endpoint available for internal diagnostics:
  - `GET /audit/events`
  - authenticated user scope by default
  - cross-user query via `user_id` allowed only for `admin` role

## Invariants
- Audit persistence must not change business operation outcomes in producer contexts.
- Payload must stay redacted/allow-listed.
- Events are append-only in application behavior.

## Current Integration
- `chat` publishes UI-action audit events via `audit` usecase (`WriteEvent`).
- `audit` persistence is handled by DB secondary adapter (`aion_api.audit_action_events`).
