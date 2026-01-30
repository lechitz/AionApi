# internal/auth/adapter/secondary/cache

Token storage adapter for auth sessions.

## Purpose & Main Capabilities

- Save, fetch, and delete token references by user ID.
- Enforce TTL for access and refresh tokens.

## Package Composition

- Cache-backed store implementation.

## Flow (Where it comes from -> Where it goes)

Usecase -> AuthStore -> cache backend

## Why It Was Designed This Way

- Enable token revocation and mismatch checks.
- Keep session storage independent of core logic.

## Recommended Practices Visible Here

- Use namespaced keys per user ID.
- Keep TTLs configurable.
- Avoid logging token values.

## Differentials

- Deterministic token reference storage.

## What Should NOT Live Here

- HTTP transport logic or domain rules.
