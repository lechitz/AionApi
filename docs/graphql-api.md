# GraphQL Guide

This page is the practical GraphQL entrypoint for AionApi consumers.

## Endpoints

| Endpoint | Purpose |
| --- | --- |
| `http://localhost:8080/graphql` | Main GraphQL endpoint |
| `http://localhost:8080/graphql/playground` | Native server playground (if enabled) |
| `graphql-playground.md` | Docs-hosted operation explorer (informative by default) |

## Authentication

Most operations require Bearer token authentication:

```http
Authorization: Bearer <JWT_TOKEN>
```

## Canonical Query Example

```graphql
query Categories {
  categories {
    id
    name
    description
    colorHex
    icon
  }
}
```

## Canonical Mutation Example

```graphql
mutation CreateCategory($input: CreateCategoryInput!) {
  createCategory(input: $input) {
    id
    name
    colorHex
  }
}
```

### Variables

```json
{
  "input": {
    "name": "Work",
    "description": "Work items",
    "colorHex": "#FF8B42",
    "icon": "briefcase"
  }
}
```

## Example via `curl`

```bash
curl -s http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query":"query Categories { categories { id name colorHex } }"}' | jq
```

## Complete Operation Catalog

Use the docs-hosted page for all queries and mutations:

- [GraphQL Playground](graphql-playground.md)

## Contract Artifacts

Use generated artifacts for tooling and integration:

- Flattened SDL: [`docs/graphql/schema.graphql`](graphql/schema.graphql)
- GraphQL docs/artifacts: [GraphQL Artifacts](graphql/README.md)
- Shared queries for consumers: `contracts/graphql/queries/`

## Regeneration Workflow

```bash
make graphql
make graphql.schema
```

## Best Practices

- Keep business logic in usecases, not in resolvers/controllers.
- Use context propagation in handlers and repositories.
- Keep schema modules cohesive by bounded context.
- Regenerate artifacts in the same PR as schema changes.

## Next Step

Use [REST API](rest-api.md) for endpoint-level REST exploration in the same docs portal.
