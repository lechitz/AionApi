# internal/auth/adapter/secondary

Secondary adapters for auth infrastructure dependencies.

## Purpose & Main Capabilities

- Store and retrieve token references for revocation checks.
- Implement auth output ports against concrete backends.

## Package Composition

- `cache/`
  - Token storage adapter.

## Flow (Where it comes from -> Where it goes)

Usecase -> auth store -> cache backend

## Why It Was Designed This Way

- Keep infra details out of auth core.
- Allow cache provider changes without core updates.

## Recommended Practices Visible Here

- TTL-based storage for access/refresh tokens.
- Avoid logging token values.

## Differentials

- Cache-backed token reference checks.

## What Should NOT Live Here

- Business rules or HTTP logic.
