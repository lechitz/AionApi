# internal/adapter/secondary/contextlogger

Context-aware logger adapter that enriches logs with trace and request metadata.

## Package Composition

- Logger implementation compatible with `logger.ContextLogger` port.

## Flow (Where it comes from -> Where it goes)

Context -> logger adapter -> structured log output

## Why It Was Designed This Way

- Standardize log metadata across the system.
- Enable trace correlation without manual plumbing.

## Recommended Practices Visible Here

- Keep log schema aligned with the observability stack.
- Avoid PII and secrets in logs.
- Include only stable identifiers (trace_id, request_id).

## Differentials

- Automatic context enrichment for logs.

## What Should NOT Live Here

- Business logic or transport concerns.
