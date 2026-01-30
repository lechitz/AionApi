# internal/auth

Authentication and authorization domain for login/logout, token validation, and session/claims handling.

## Purpose and Main Capabilities

- Authenticate users and issue access/refresh tokens.
- Validate tokens and extract claims for downstream handlers.
- Support revocation checks via token storage.
- Provide HTTP middleware for protected routes.

## Package Composition

- `core/`: auth domain rules, ports, and usecases.
- `core/ports/input`: auth service contract used by handlers/middleware.
- `core/ports/output`: token provider and token store contracts.
- `core/usecase`: login/logout/refresh/validate orchestration.
- `adapter/primary/http`: HTTP handlers, middleware, and DTO mapping.
- `adapter/secondary`: cache-backed token store adapters.

## Flow (Where it comes from -> Where it goes)

HTTP request -> primary adapter -> input port -> usecase ->
output ports (token provider/store) -> secondary adapter -> cache -> response

## Diagram

![Auth Domain Flow](../../docs/diagram/images/internal-auth.svg)

Source: `../../docs/diagram/internal-auth.sequence.txt`

## How It Works (Concise)

- Handlers accept login/logout/refresh, validate input, and map DTOs to commands.
- Middleware extracts bearer tokens or cookies and validates claims.
- Usecases enforce auth policy and return semantic errors.
- Token store adapters check revocation and persist token references.

## Separation Inside the Bounded Context

- Core (domain, ports, usecases) is transport-agnostic and does not import adapters.
- Primary adapters own HTTP rules, cookies, headers, and DTO mapping.
- Secondary adapters isolate cache/infra details and error translation.
- Shared errors and constants come from `internal/shared`.

## Why It Was Designed This Way

- Keep auth policies isolated from transport and storage.
- Enforce consistent token validation and revocation checks.
- Provide a single source of truth for security flows.

## Recommended Practices Visible Here

- Sanitize bearer tokens and handle cookie fallback.
- Never log credentials or tokens.
- Use semantic errors for invalid credentials or token mismatch.
- Keep middleware thin: extract token, call input port, set context.

## Differentials

- Dual token sources (header and cookie) with centralized validation.
- Cache-backed token reference checks for revocation.

## What Should NOT Live Here

- Business logic unrelated to authentication.
- Direct database access from core.
- Transport-specific DTOs outside adapters.
- Cross-context imports or shared state leakage.
