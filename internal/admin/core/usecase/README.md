# internal/admin/core/usecase

Admin usecase implementations that enforce privileged workflows.

## Purpose & Main Capabilities

- Role management and admin-only mutations.
- Policy enforcement and semantic error handling.
- Coordination with admin repositories.

## Package Composition

- Usecase implementations and constants.
- Tests for admin workflows.

## Flow (Where it comes from -> Where it goes)

Input port -> usecase -> output port -> adapter

## Why It Was Designed This Way

- Centralize admin policy enforcement.
- Keep transport and storage concerns out of core.

## Recommended Practices Visible Here

- Use context-first signatures and spans.
- Return semantic errors, not driver errors.
- Keep usecases small and policy-focused.

## Differentials

- Privileged policy enforcement kept inside usecases.

## What Should NOT Live Here

- HTTP handlers or persistence code.
