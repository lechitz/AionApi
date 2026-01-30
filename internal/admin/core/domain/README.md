# internal/admin/core/domain

Admin domain types and invariants.

## Purpose & Main Capabilities

- Define admin-specific entities and role concepts.
- Express invariants for privileged actions.

## Package Composition

- Domain structs, value objects, and validation helpers.

## Flow (Where it comes from -> Where it goes)

Usecases -> domain entities -> validation/invariants

## Why It Was Designed This Way

- Keep domain rules transport-agnostic.
- Make admin policies explicit and reusable.

## Recommended Practices Visible Here

- Pure Go domain (stdlib only).
- No persistence or transport tags.

## Differentials

- Explicit admin domain boundary.

## What Should NOT Live Here

- Adapters, repositories, or HTTP logic.
