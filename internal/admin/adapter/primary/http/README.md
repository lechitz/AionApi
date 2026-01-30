# internal/admin/adapter/primary/http

HTTP handlers for admin operations and transport mapping.

## Purpose & Main Capabilities

- Expose admin-only routes (roles, block/unblock, updates).
- Enforce elevated authorization and input validation.
- Map semantic errors to HTTP responses.

## Package Composition

- `dto/`
  - Admin HTTP request/response shapes.
- Handlers and routing helpers.

## Flow (Where it comes from -> Where it goes)

HTTP request -> handler -> core/usecase -> HTTP response

## Why It Was Designed This Way

- Centralize admin transport logic and policies.
- Keep core logic transport-agnostic.

## Recommended Practices Visible Here

- Validate admin claims early and log audit metadata.
- Keep DTO mapping inside adapters (see `dto/` README).

## Differentials

- Audit-friendly admin endpoints with strict boundary checks.

## What Should NOT Live Here

- Business rules or persistence code.
