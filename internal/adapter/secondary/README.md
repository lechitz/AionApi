# internal/adapter/secondary

Secondary adapters implement output ports for external integrations (db, cache, crypto, tokens, logging, graphs).

## Package Composition

- `db/`, `cache/`, `crypto/`, `hasher/`, `token/`, `contextlogger/`, `graph/`

## Flow (Where it comes from -> Where it goes)

Usecase -> output port -> secondary adapter -> external system

## Why It Was Designed This Way

- Isolate infrastructure details from core logic.
- Standardize error translation and resilience policies.
- Keep external integrations swappable.

## Recommended Practices Visible Here

- Translate infra errors to semantic errors in core.
- Keep observability at adapter boundaries.
- Apply timeouts/retries where needed.

## Differentials

- Consistent integration patterns across domains.

## What Should NOT Live Here

- Business rules or domain logic.
- Transport-layer DTOs.
