# internal/auth/adapter

Auth adapters connect the core auth logic to HTTP transport and infrastructure dependencies.

## Purpose & Main Capabilities

- Expose auth endpoints and middleware.
- Map HTTP requests to auth usecases.
- Store and retrieve token references.

## Package Composition

- `primary/`
  - HTTP handlers and middleware.
- `secondary/`
  - Token storage adapters (cache).

## Flow (Where it comes from -> Where it goes)

HTTP request -> primary adapter -> core/usecase -> secondary adapter

## Why It Was Designed This Way

- Keep transport and infra concerns out of core.
- Standardize auth boundary behavior.

## Recommended Practices Visible Here

- Keep handlers thin and policy-aware.
- Never log credentials or tokens.

## Differentials

- Dedicated auth boundary with middleware for protected routes.

## What Should NOT Live Here

- Business rules outside authentication.
- Domain logic inside handlers.
