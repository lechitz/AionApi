# Grafana Provisioning Assets

**Path:** `infrastructure/observability/grafana`

## Overview

Grafana provisioning files for dashboards and datasources.
These assets make Grafana startup deterministic across local environments.

## Structure

| Folder | Purpose |
| --- | --- |
| `dashboards/` | Versioned dashboard JSON files |
| `datasources/` | Datasource provisioning files |
| `provisioning/` | Dashboard provider configuration |

## Design Notes

- Treat dashboard JSON as versioned code artifacts.
- Keep datasource naming stable to avoid dashboard breakage.
- Do not store secrets in dashboard definitions.

## Package Improvements

- Add dashboard ownership map and review policy.
- Add JSON validation/lint in CI.
- Add changelog for major dashboard redesigns.
- Add screenshot snapshots for baseline visual regression.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
