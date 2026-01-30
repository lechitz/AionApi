# internal/admin/core/ports/output

Admin output port contracts.

## Purpose & Main Capabilities

- Define repository and integration interfaces for admin usecases.

## Package Composition

- Admin output port interfaces.

## Flow (Where it comes from -> Where it goes)

Usecase -> output port -> adapter implementation

## Why It Was Designed This Way

- Keep infra details out of core.
- Allow swapping implementations without core changes.

## Recommended Practices Visible Here

- Keep interfaces focused on admin needs.
- Return semantic errors and domain types.

## Differentials

- Stable contracts for admin persistence and integrations.

## What Should NOT Live Here

- Concrete implementations or transport types.
