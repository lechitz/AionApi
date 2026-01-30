# internal/shared

Cross-cutting shared assets used across multiple contexts without business rules.

## Purpose and Main Capabilities

- Centralize constants for headers, logging keys, tracing attributes, and context keys.
- Reduce magic strings and keep telemetry consistent.
- Provide safe, minimal shared types for context propagation.

## Package Composition

- `constants/`: claims, common keys, context keys, tracing keys.

## Flow (Where it comes from -> Where it goes)

Adapters/usecases -> shared constants -> logs/headers/traces

## Why It Was Designed This Way

- Avoid duplication and typos across contexts.
- Keep cross-cutting conventions stable and discoverable.
- Preserve domain isolation (no shared business logic).

## Recommended Practices Visible Here

- Prefer `commonkeys` and `tracingkeys` over ad-hoc strings.
- Use `ctxkeys` for context values, not header names.
- Keep new keys close to their domain area under `constants/`.

## Differentials

- Single source of truth for cross-cutting keys.

## What Should NOT Live Here

- Domain logic or usecases.
- Adapter implementations or infra clients.
