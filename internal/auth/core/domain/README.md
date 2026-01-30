# internal/auth/core/domain

Domain types and invariants used by the auth usecases.

## Purpose & Main Capabilities

- Define auth-related domain types and constraints.
- Keep auth domain pure and transport-agnostic.

## Package Composition

- Domain structs and validation helpers.

## Flow (Where it comes from -> Where it goes)

Usecases -> domain types -> validation/invariants

## Why It Was Designed This Way

- Keep domain logic free of infra and transport concerns.

## Recommended Practices Visible Here

- Pure Go domain (stdlib only).
- No persistence or transport tags.

## Differentials

- Clean domain boundary for security-sensitive logic.

## What Should NOT Live Here

- Adapters, repositories, or HTTP logic.
