# GraphQL Documentation Artifacts

**Path:** `docs/graphql`

## Overview

This folder stores generated GraphQL contract artifacts used by consumers and tooling.
It complements schema modules under `internal/adapter/primary/graphql/schema/modules/`.

## Files

| File | Purpose |
| --- | --- |
| `schema.graphql` | Flattened SDL representation of the current GraphQL schema |

## Regeneration Workflow

```bash
make graphql.schema
```

If your local setup includes introspection export targets, run them after server startup.

## Related Sources

- Schema modules: `internal/adapter/primary/graphql/schema/modules/`
- Shared queries: `contracts/graphql/queries/`

## Dashboard White Label Contract Notes

The schema now exposes dashboard layout and composition primitives:

- `dashboardViews`, `dashboardView`
- `dashboardWidgetCatalog`
- `suggestMetricDefinitions`
- `createDashboardView`, `setDefaultDashboardView`
- `upsertDashboardWidget`, `reorderDashboardWidgets`, `deleteDashboardWidget`
- `createMetricAndWidget`

## Design Notes

- Treat files here as generated artifacts.
- Regenerate after schema changes in the same PR.
- Keep docs artifacts aligned with consumer tooling expectations.

## Package Improvements

- Add CI guard for stale generated schema artifacts.
- Add optional introspection generation target with health checks.
- Add version metadata header for schema snapshots.
- Add compatibility note for codegen consumers.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
