# Shared GraphQL Operations

**Path:** `contracts/graphql`

## Overview

This folder stores reusable GraphQL operation documents aligned with the backend schema.
It is intended to be the shared contract surface consumed by Aion clients.

## Structure

| Folder | Scope |
| --- | --- |
| `manifest.json` | Deterministic index of shared operations and checksums |
| `queries/categories/` | Category read operations |
| `queries/tags/` | Tag read operations |
| `queries/records/` | Record read operations |
| `queries/chat/` | Chat read operations |
| `queries/user/` | User read operations |
| `mutations/categories/` | Category write operations |
| `mutations/tags/` | Tag write operations |
| `mutations/records/` | Record write operations |

## Query Inventory

### Categories
- `queries/categories/list.graphql`
- `queries/categories/by-id.graphql`
- `queries/categories/by-name.graphql`

### Tags
- `queries/tags/list.graphql`
- `queries/tags/by-id.graphql`
- `queries/tags/by-name.graphql`
- `queries/tags/by-category-id.graphql`

### Records
- `queries/records/list.graphql`
- `queries/records/by-id.graphql`
- `queries/records/latest.graphql`
- `queries/records/by-tag.graphql`
- `queries/records/by-category.graphql`
- `queries/records/by-day.graphql`
- `queries/records/until.graphql`
- `queries/records/between.graphql`
- `queries/records/search.graphql`
- `queries/records/stats.graphql`

### Chat
- `queries/chat/history.graphql`
- `queries/chat/context.graphql`
- `queries/chat/data-pack.graphql`

### User
- `queries/user/stats.graphql`

## Mutation Inventory

### Categories
- `mutations/categories/create.graphql`
- `mutations/categories/update.graphql`
- `mutations/categories/delete.graphql`

### Tags
- `mutations/tags/create.graphql`
- `mutations/tags/update.graphql`
- `mutations/tags/delete.graphql`

### Records
- `mutations/records/create.graphql`
- `mutations/records/update.graphql`
- `mutations/records/delete.graphql`
- `mutations/records/delete-all.graphql`

## Notes

- Keep operation names stable for observability and client cache behavior.
- Keep selection sets consistent across clients unless a consumer needs a narrower shape.
- Validate operation documents against current schema in CI.
- Regenerate contract files + manifest with: `make graphql.queries graphql.manifest`.
- Validate schema compatibility with: `make graphql.validate`.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
