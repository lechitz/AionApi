# cmd/ (Application Entrypoints)

This folder contains the main application binary for AionAPI. The code here should orchestrate, not implement domain logic.

## Structure

- `api/` - Main API server (GraphQL + HTTP)
  - Wires modules with Fx and configures Swagger metadata
  - **This is the only production application**

## Flow (Where it comes from -> Where it goes)

Operator → cmd/api → internal/platform/modules → adapters/usecases

The entrypoint builds the process graph and delegates real work to `internal/` layers.

## Development Tools

Development utilities have been moved to `/hack`:
- Go CLIs (seed-caller, seed-helper) → `hack/tools/`
- Bash scripts (test-*, force-insert-roles.sh) → `hack/dev/`

See `hack/README.md` for details.

## Why It Was Designed This Way

- Keep entrypoints thin and predictable
- Centralize configuration, lifecycle, and observability early
- Preserve clean boundaries: cmd never owns business rules
- Separate production code from development utilities

## Recommended Practices Visible Here

- main packages only orchestrate; no domain rules
- Configuration loaded via `internal/platform/config`
- Graceful shutdown with context timeout

## Differentials

- Swagger metadata applied at runtime from config
- Single production binary (clean separation from dev tools)

## What Should NOT Live Here

- Domain rules or validation
- Mapping of HTTP/GraphQL payloads
- Repository or external IO implementations
- Development tools (use `hack/` instead)

