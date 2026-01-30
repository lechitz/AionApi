# internal/adapter/secondary

Secondary adapters implement output ports for external integrations (db, cache, crypto, tokens, logging, graphs).

## Purpose and Main Capabilities

- Bridge usecases to infrastructure without leaking infra types into core.
- Provide swappable implementations for external systems.
- Translate infra errors into semantic errors expected by the core.
- Centralize resilience and observability at IO boundaries.

## Package Composition

- `db/`: persistence adapters (SQL databases).
- `cache/`: cache adapters (e.g., Redis).
- `graph/`: graph store adapters (e.g., Gremlin).
- `crypto/`: encryption/decryption helpers behind output ports.
- `hasher/`: hashing adapters for password or token material.
- `token/`: token issuance/validation adapters.
- `contextlogger/`: logger adapter used by output ports that require structured logging.

## Flow (Where it comes from -> Where it goes)

Usecase -> output port -> secondary adapter -> external system

## Why It Was Designed This Way

- Isolate infrastructure details from core logic.
- Standardize error translation and resilience policies.
- Keep external integrations swappable.
- Encourage consistent IO boundaries across contexts.

## Recommended Practices Visible Here

- Translate infra errors to semantic errors in core.
- Keep observability at adapter boundaries.
- Apply timeouts/retries where needed.
- Keep adapters thin: map inputs/outputs and delegate to infra clients.

## Differentials

- Consistent integration patterns across domains.
- Port-first adapters that keep Clean Architecture boundaries intact.

## What Should NOT Live Here

- Business rules or domain logic.
- Transport-layer DTOs.
- Input ports or usecase orchestration (those belong in core).
