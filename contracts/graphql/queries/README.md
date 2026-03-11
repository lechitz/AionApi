# Shared GraphQL Operations

**Path:** `contracts/graphql`

## Overview

This folder stores reusable GraphQL operation documents aligned with the backend schema.
It is intended to be the shared contract surface consumed by Aion clients.

For v1 insight and analytics work, this folder is not optional documentation.
It is part of the canonical consumer contract and must remain aligned with:

1. schema modules under `internal/adapter/primary/graphql/schema/modules/`
2. generated server artifacts
3. downstream typed consumers in dashboard/chat layers

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

### Dashboard
- `queries/dashboard/snapshot.graphql`
- `queries/dashboard/insight-feed.graphql`
- `queries/dashboard/analytics-series.graphql`
- `queries/dashboard/metric-definitions.graphql`
- `queries/dashboard/views.graphql`
- `queries/dashboard/view.graphql`
- `queries/dashboard/widget-catalog.graphql`
- `queries/dashboard/suggest-metric-definitions.graphql`

## v1 Insight and Analytics Contract Notes

Canonical operations:

- `InsightFeed`
- `AnalyticsSeries`

Current scope model shared by both operations:

- `window`
- optional `date`
- optional `timezone`
- optional `categoryId`
- optional `tagIds`

Current v1 restriction:

- `AnalyticsSeries` is intentionally narrow and currently centered on `records.count`

Consumer compatibility rules:

- additive changes are preferred
- field removal or meaning changes require explicit coordinated updates across `AionApi`, `aionapi-dashboard`, and `aion-chat`
- the first insight returned by `InsightFeed` is the dominant insight for that scope
- consumers may humanize wording, but they must not reinterpret business meaning

Governance reference:

- `/Aion/notes/v1-0-0/v1-gov-04-insight-api-contract-policy.md`

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
- Treat `queries/dashboard/insight-feed.graphql` and `queries/dashboard/analytics-series.graphql` as backend-owned public contracts for v1 surfaces such as `Radar`, `Analytics`, and MCP read tools.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
