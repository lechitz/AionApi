# internal/adapter/primary/graphql

Shared GraphQL infrastructure: schema, directives, and transport models used by all contexts.

## Package Composition

- `schema/`
  - Schema definitions and module composition.
- `directives/`
  - Cross-cutting rules at schema level (auth, validation, tracing).
- `model/`
  - GraphQL transport models.

## Flow (Where it comes from -> Where it goes)

Schema + directives -> gqlgen -> resolvers/controllers -> usecases

## Why It Was Designed This Way

- Keep schema modular and composable.
- Centralize cross-cutting GraphQL behavior.
- Avoid duplication across contexts.

## Recommended Practices Visible Here

- Keep modules in `schema/modules` for clear ownership.
- Use directives instead of duplicating policy in resolvers.
- Align models with schema and controller mappings.

## Differentials

- Modular schema + directive-driven policies.

## What Should NOT Live Here

- Domain logic or persistence.
- Context-specific controllers (they live under each context).
