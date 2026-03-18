# Prometheus Scrape Config

**Path:** `infrastructure/observability/prometheus`

## Purpose

This package defines Prometheus scraping for the local observability stack.

## Current File

| File | Responsibility |
| --- | --- |
| `prometheus.yml` | scrapes the OTEL collector metrics endpoint every 5 seconds |

## Boundaries

- new scrape targets should only be added when they are intentionally part of the local operator surface
- metric naming and labels remain owned by the emitting services and OTEL pipeline
- dashboards that depend on new metrics should ship with scrape changes in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
