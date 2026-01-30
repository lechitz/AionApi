# internal/adapter/primary/graphql/schema

Composable GraphQL schema definitions and module wiring.

## Package Composition

- `modules/`
  - Per-domain schema fragments.
- Schema root and shared primitives.

## Flow (Where it comes from -> Where it goes)

Schema modules -> gqlgen -> GraphQL API contracts

## Why It Was Designed This Way

- Keep schemas modular and owned by domains.
- Enable selective inclusion of features.
- Avoid a monolithic schema file.

## Recommended Practices Visible Here

- Keep modules self-contained and consistent with controllers.
- Update schema and resolvers together.
- Reuse shared directives and models.

## Differentials

- Modular schema composition across contexts.

## What Should NOT Live Here

- Resolver logic or business rules.
- Context-specific behavior not represented in schema.
