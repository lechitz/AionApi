# internal/admin/core

Core admin business logic built around domain models and ports.

## Purpose & Main Capabilities

- Enforce admin policies and invariants.
- Orchestrate privileged workflows via ports.
- Return semantic errors for adapters to map.

## Package Composition

- `domain/`
  - Admin domain entities and rules.
- `ports/`
  - Input/output contracts.
- `usecase/`
  - Admin usecase implementations.

## Flow (Where it comes from -> Where it goes)

Primary adapter -> input port -> usecase -> output port

## Why It Was Designed This Way

- Keep admin logic independent of transport and persistence.
- Preserve strict policy enforcement in one place.

## Recommended Practices Visible Here

- Usecases depend only on ports.
- Semantic errors define expected outcomes.

## Differentials

- Privileged policy enforcement centralized in core.

## What Should NOT Live Here

- Transport DTOs or persistence models.
