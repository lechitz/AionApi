# internal/auth/core

Core authentication logic built around domain rules and ports.

## Purpose & Main Capabilities

- Enforce authentication policies and invariants.
- Orchestrate login/logout/validate/refresh flows.
- Return semantic errors for adapters to map.

## Package Composition

- `domain/`
  - Auth domain types and invariants.
- `ports/`
  - Input/output interfaces.
- `usecase/`
  - Auth usecase implementations.

## Flow (Where it comes from -> Where it goes)

Primary adapter -> input port -> usecase -> output ports

## Why It Was Designed This Way

- Keep auth rules independent of transport and infra.
- Preserve dependency inversion for security-critical logic.

## Recommended Practices Visible Here

- Context-first signatures with tracing.
- Claims validation and token sanitization in core.
- Semantic errors for auth failures.

## Differentials

- Centralized auth policy enforcement in core.

## What Should NOT Live Here

- HTTP handlers or cache implementations.
