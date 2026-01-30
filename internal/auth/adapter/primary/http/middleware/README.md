# internal/auth/adapter/primary/http/middleware

Auth middleware for protecting routes and injecting claims into context.

## Purpose & Main Capabilities

- Validate tokens on protected routes.
- Extract user ID and claims into context.
- Support bearer header and cookie tokens.

## Package Composition

- Auth middleware implementation and constants.

## Flow (Where it comes from -> Where it goes)

HTTP request -> auth middleware -> core/usecase -> next handler

## Why It Was Designed This Way

- Enforce auth before handlers run.
- Keep token validation centralized.

## Recommended Practices Visible Here

- Sanitize tokens before validation.
- Use context keys for downstream access.
- Never log raw tokens.

## Differentials

- Dual token source support with consistent validation.

## What Should NOT Live Here

- Business rules or persistence logic.
