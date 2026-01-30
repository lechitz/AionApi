# internal/adapter/primary/graphql/schema/modules

Per-domain GraphQL schema modules, composed into the root schema.

## Package Composition

- `*.graphqls`
  - Domain-specific queries, mutations, and types.

## Flow (Where it comes from -> Where it goes)

Domain schema module -> schema composition -> gqlgen -> resolvers/controllers

## Why It Was Designed This Way

- Enable modular ownership of GraphQL contracts.
- Keep each domain schema isolated and reusable.

## Recommended Practices Visible Here

- Keep modules aligned with domain usecases and controllers.
- Document any directive dependencies.
- Update schema and resolver mappings together.

## Differentials

- Modular schema contracts per domain.

## What Should NOT Live Here

- Business logic or resolver code.
- Types unrelated to a domain module.
