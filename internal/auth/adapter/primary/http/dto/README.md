# internal/auth/adapter/primary/http/dto

HTTP DTOs for auth endpoints.

## Purpose & Main Capabilities

- Define request/response shapes for auth handlers.

## Package Composition

- Login, logout, refresh request/response structs.

## Flow (Where it comes from -> Where it goes)

HTTP payload -> DTO -> handler -> usecase

## Why It Was Designed This Way

- Keep transport shapes separate from domain entities.

## Recommended Practices Visible Here

- Validate in handlers, not in core.
- Avoid leaking internal fields.

## Differentials

- Explicit auth transport contracts.

## What Should NOT Live Here

- Business rules or persistence models.
