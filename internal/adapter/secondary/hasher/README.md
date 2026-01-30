# internal/adapter/secondary/hasher

Hashing adapters for passwords and integrity checks.

## Package Composition

- Hasher implementation compatible with the hasher port.

## Flow (Where it comes from -> Where it goes)

Usecase -> hasher port -> hashing adapter -> hash/compare

## Why It Was Designed This Way

- Centralize hashing policy and cost parameters.
- Prevent leaking crypto details into core logic.

## Recommended Practices Visible Here

- Use secure algorithms and constant-time compares.
- Keep cost factors documented and configurable.
- Avoid logging sensitive inputs.

## Differentials

- Secure defaults and consistent hashing policy.

## What Should NOT Live Here

- Authentication business rules.
