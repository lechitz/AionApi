# internal/auth/core/ports/output

Auth output port contracts.

## Purpose & Main Capabilities

- Define token provider and store interfaces for auth usecases.

## Package Composition

- Output port interfaces for token generation/validation and token storage.

## Flow (Where it comes from -> Where it goes)

Usecase -> output port -> adapter implementation

## Why It Was Designed This Way

- Keep infra details out of core.
- Allow swapping token providers or stores.

## Recommended Practices Visible Here

- Return semantic errors to core.
- Keep interfaces minimal and auth-focused.

## Differentials

- Explicit contracts for security dependencies.

## What Should NOT Live Here

- Concrete implementations or transport types.
