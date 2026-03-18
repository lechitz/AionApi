# Grafana Provisioning Assets

**Path:** `infrastructure/observability/grafana`

## Purpose

This package makes Grafana startup deterministic by provisioning datasources and dashboards from versioned files.

## Current Layout

| Folder | Purpose |
| --- | --- |
| `dashboards/` | checked-in dashboard JSON, including RED and HTTP request views |
| `datasources/` | Prometheus, Loki, and Jaeger datasource definitions |
| `provisioning/` | dashboard provider configuration loaded at container startup |

## Boundaries

- dashboard JSON and datasource YAML are code artifacts and should be reviewed like source
- datasource names must stay stable or existing dashboards will break
- this package provisions Grafana only; scrape, storage, and transport behavior live in the other observability components

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
