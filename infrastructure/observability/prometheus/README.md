# Prometheus Scrape Config

**Path:** `infrastructure/observability/prometheus`

## Overview

Prometheus scrape configuration for metrics collection in local observability stack.
It scrapes OTel collector and other configured targets.

## Files

| File | Purpose |
| --- | --- |
| `prometheus.yml` | Scrape jobs, intervals, and target definitions |

## Design Notes

- Keep scrape intervals appropriate for environment scale.
- Keep target naming predictable for dashboard query reuse.
- Verify target health after config updates.

## Package Improvements

- Add example queries used by core dashboards.
- Add validation target for config syntax in CI.
- Add guidance for adding new scrape jobs safely.
- Add cardinality warning section for high-volume labels.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
