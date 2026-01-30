# Platform HTTP — Generic

**Folder:** `internal/platform/server/http/generic`
**Subpackages:** `dto/`, `handler/`

## Purpose and Main Capabilities

- Provide platform-wide HTTP handlers (health, 404/405, error, recovery).
- Keep cross-cutting responses out of domain adapters.
- Ensure consistent envelopes, logs, and tracing.

## Package Composition

- `handler/`: generic handlers (health, not found, method not allowed, error, recovery).
- `dto/`: transport-only DTOs (e.g., health response).

## Flow (Where it comes from -> Where it goes)

HTTP composer -> generic handlers -> httpresponse/sharederrors -> response

## How it works

- `New(logger, generalCfg)` builds a `*handler.Handler`.
- The composer wires:
  - `GET /health` -> `Handler.HealthCheck`
  - `SetNotFound`, `SetMethodNotAllowed`, `SetError`
  - Recovery middleware uses `Handler.RecoveryHandler`
- Responses use `httpresponse` for consistent JSON envelopes.

## Endpoints / Behaviors

- `/health` returns service metadata (`name`, `env`, `version`, `timestamp`).
- 404/405 return standardized JSON errors.
- Error/Recovery return safe 500 responses with correlation IDs.

## Observability

- Spans per operation (see `0_generic_handler_constants.go`).
- Logs include request ID and metadata only.

## What Should NOT Live Here

- Domain logic or context-specific behavior.
