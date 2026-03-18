# Loki Configuration

**Path:** `infrastructure/observability/loki`

## Purpose

This package configures Loki as the log store queried by Grafana.

## Current File

| File | Responsibility |
| --- | --- |
| `loki.yaml` | single-binary Loki config with filesystem storage and 7-day retention |

## Boundaries

- ingestion labels must stay compatible with Fluent Bit output and Grafana queries
- this config targets the current local and prod-like single-node setup, not a scaled Loki deployment
- parser behavior belongs upstream in Fluent Bit, not here

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
