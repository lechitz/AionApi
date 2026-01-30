# internal/admin/adapter/primary/http/dto

HTTP DTOs for admin endpoints.

## Purpose & Main Capabilities

- Define request/response shapes for admin HTTP handlers.

## Package Composition

- Request and response structs used by handlers.

## Flow (Where it comes from -> Where it goes)

HTTP payload -> DTO -> handler -> usecase

## Why It Was Designed This Way

- Keep transport shapes separate from domain entities.

## Recommended Practices Visible Here

- Validate and map in adapters, not in core.
- Avoid exposing internal fields.

## Differentials

- Explicit admin transport contracts.

## What Should NOT Live Here

- Business rules or persistence models.
