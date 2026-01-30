# internal/auth/core/ports

Auth ports define the contracts between core and adapters.

## Purpose & Main Capabilities

- Declare input ports for auth usecases.
- Declare output ports for token providers and storage.

## Package Composition

- `input/`
  - Auth service interfaces.
- `output/`
  - Token provider and store interfaces.

## Flow (Where it comes from -> Where it goes)

Adapter -> input port -> usecase -> output port -> adapter

## Why It Was Designed This Way

- Enforce dependency inversion for auth flows.
- Keep integrations swappable.

## Recommended Practices Visible Here

- Keep ports focused and context-first.
- Avoid leaking transport types.

## Differentials

- Clear contracts for security-critical integrations.

## What Should NOT Live Here

- Implementations or wiring.
