# internal/adapter/secondary/token

Token adapters for session and service tokens.

## Package Composition

- Token providers for issuing and validating tokens.

## Flow (Where it comes from -> Where it goes)

Usecase -> token port -> token adapter -> crypto/claims validation

## Why It Was Designed This Way

- Centralize token policies and lifetimes.
- Avoid leaking crypto details to the core.

## Recommended Practices Visible Here

- Keep claims aligned with `internal/shared/constants`.
- Avoid logging token values.
- Document rotation and revocation strategy.

## Differentials

- Consistent token handling across contexts.

## What Should NOT Live Here

- Authentication business rules.
- Transport/session storage logic.
