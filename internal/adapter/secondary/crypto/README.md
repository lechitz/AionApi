# internal/adapter/secondary/crypto

Cryptography adapters for encrypt/decrypt and signing operations.

## Package Composition

- Crypto provider implementation and helpers.

## Flow (Where it comes from -> Where it goes)

Usecase -> crypto port -> crypto adapter -> cipher/signing ops

## Why It Was Designed This Way

- Keep crypto primitives centralized and consistent.
- Allow key management and defaults outside of core logic.

## Recommended Practices Visible Here

- Use secure defaults and document key rotation.
- Normalize crypto errors for callers.
- Never log secrets or plaintext.

## Differentials

- Consistent crypto wrapper with safe defaults.

## What Should NOT Live Here

- Business logic or transport concerns.
