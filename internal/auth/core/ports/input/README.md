# internal/auth/core/ports/input

Auth input port contracts.

## Purpose & Main Capabilities

- Define the auth service interface used by adapters.

## Package Composition

- Auth input port interfaces (login, logout, validate, refresh).

## Flow (Where it comes from -> Where it goes)

Adapter -> input port -> usecase

## Why It Was Designed This Way

- Keep adapters decoupled from concrete usecases.

## Recommended Practices Visible Here

- Context-first signatures.
- Domain types and semantic errors.

## Differentials

- Stable auth interface for multiple entrypoints.

## What Should NOT Live Here

- Implementations or transport DTOs.
