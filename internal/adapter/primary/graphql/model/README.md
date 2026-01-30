# internal/adapter/primary/graphql/model

GraphQL transport models used by resolvers and controllers.

## Package Composition

- GraphQL-specific structs generated or maintained for schema alignment.

## Flow (Where it comes from -> Where it goes)

Domain entity -> model mapping -> GraphQL response

## Why It Was Designed This Way

- Keep transport shapes separate from domain entities.
- Allow safe exposure of fields with masking/transformations.

## Recommended Practices Visible Here

- Do not leak persistence types or DB tags.
- Keep mapping logic in adapters/controllers.
- Maintain consistency with schema definitions.

## Differentials

- Explicit transport models for safe exposure.

## What Should NOT Live Here

- Business rules or persistence concerns.
