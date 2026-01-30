# internal/auth/core/usecase

Auth usecase implementations for login, logout, validate, and refresh.

## Purpose & Main Capabilities

- Validate credentials and issue tokens.
- Verify token integrity and claims.
- Coordinate token storage for revocation checks.

## Package Composition

- Usecase implementations and tests.

## Flow (Where it comes from -> Where it goes)

Input port -> usecase -> output ports -> adapters

## Why It Was Designed This Way

- Centralize auth policy enforcement.
- Keep storage and JWT details behind ports.

## Recommended Practices Visible Here

- Sanitize bearer tokens and normalize claims.
- Never log credentials or raw tokens.
- Use semantic errors for invalid credentials or token mismatch.

## Differentials

- Token validation with cache-backed reference checks.

## What Should NOT Live Here

- HTTP handlers or cache/JWT implementations.
