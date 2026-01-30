# internal/auth/adapter/primary/http/handler

Auth HTTP handlers for login, logout, refresh, and session responses.

## Purpose & Main Capabilities

- Expose auth endpoints over HTTP.
- Validate inputs and map DTOs to usecases.
- Set auth cookies and map semantic errors to HTTP.

## Package Composition

- Handler implementations and route registration.
- Tests for handler behavior.

## Flow (Where it comes from -> Where it goes)

HTTP request -> handler -> core/usecase -> HTTP response

## Why It Was Designed This Way

- Keep transport rules out of auth core.
- Centralize HTTP response mapping for auth.

## Recommended Practices Visible Here

- Keep handlers thin and deterministic.
- Never log credentials or tokens.
- Use shared response helpers for consistency.

## Differentials

- Auth-specific HTTP boundary with cookie handling.

## What Should NOT Live Here

- Business rules or persistence code.
