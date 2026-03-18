# GraphQL Documentation Artifacts

**Path:** `docs/graphql`

## Overview

This folder stores generated GraphQL artifacts used by consumers and tooling.
It complements, but does not replace, the live schema modules under `internal/adapter/primary/graphql/schema/modules/`
and the shared operation set under `contracts/graphql/`.

## Files

| File | Purpose |
| --- | --- |
| `schema.graphql` | Flattened SDL snapshot generated from the current schema modules |

## Regeneration Workflow

```bash
make graphql.schema graphql.queries graphql.manifest
```

If your local setup includes introspection export targets, run them after server startup.

## Related Sources

- Schema modules: `internal/adapter/primary/graphql/schema/modules/`
- Shared operations: `contracts/graphql/queries/` and `contracts/graphql/mutations/`
- Contract manifest: `contracts/graphql/manifest.json`
- Ownership map: `.github/DOCUMENTATION_OWNERSHIP.md`

## Public Contract Notes

These artifacts currently include the main read and dashboard surfaces consumed by the workspace:

- legacy record reads such as `recordsLatest`
- derived projection reads such as `recordProjectionById`, `recordProjections`, and `recordProjectionsLatest`
- dashboard reads such as `dashboardSnapshot`, `insightFeed`, `analyticsSeries`, `metricDefinitions`, `dashboardViews`, `dashboardView`, `dashboardWidgetCatalog`, and `suggestMetricDefinitions`
- dashboard mutations such as `createDashboardView`, `setDefaultDashboardView`, `upsertDashboardWidget`, `reorderDashboardWidgets`, `deleteDashboardWidget`, and `createMetricAndWidget`

For the widget system specifically, the published GraphQL contract is
intentionally coarse for v1:

- `dashboardWidgetCatalog` exposes canonical widget types, coarse sizes, and
  large-card limits
- widget records expose persisted `configJson`, but `AionApi` currently treats
  the richer visual layout grammar inside that JSON as a dashboard-owned concern
  rather than a server-validated schema

Consumers should not assume the backend already owns minimum/recommended grid
sizes, free-placement rules, or Radar widget scope semantics.

Consumer-facing query documents live under `contracts/graphql/` and should remain aligned with the schema snapshot here.

## Insight And Analytics Governance

`insightFeed` and `analyticsSeries` remain the canonical v1 intelligence surface used across:

- `Radar`
- canonical `Analytics`
- dashboard `INSIGHT_FEED` widgets
- MCP read tools in `aion-chat`

For compatibility policy and ownership, refer to:

- `contracts/graphql/queries/README.md`
- `aion-docs/planning/v1/adr/adr-005-insight-contract-policy.md`
- `aion-docs/planning/v1/reference/dashboard-backend-consumption.md`

## Design Notes

- Treat files here as generated artifacts.
- Regenerate after schema changes in the same PR.
- Keep this folder aligned with consumer tooling expectations, but treat `contracts/graphql/` as the reusable public operation surface.
- If this folder drifts from live modules or shared operations, the live modules win.

## Package Improvements

- Add CI guard for stale generated schema artifacts.
- Add optional introspection generation target with health checks.
- Add version metadata header for schema snapshots.
- Add compatibility note for codegen consumers.
- Expose a lightweight diff check between `schema.graphql` and published shared operations.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
