# Observability Infrastructure

**Path:** `infrastructure/observability`

## Purpose

This folder owns the repo-local telemetry configs mounted by the Docker profiles.
It covers traces, metrics, logs, and Grafana provisioning for the local and prod-like stacks.

## Current Signal Flow

1. Aion services emit OTLP telemetry to the collector.
2. `otel/` exports traces to Jaeger and metrics to a Prometheus scrape endpoint.
3. `prometheus/` scrapes the collector.
4. `fluentbit/` tails Docker logs and forwards them to `loki/`.
5. `grafana/` provisions datasources and dashboards over Prometheus, Loki, and Jaeger.

## Current Areas

| Area | Responsibility |
| --- | --- |
| `otel/` | telemetry ingestion and routing |
| `prometheus/` | metrics scraping |
| `loki/` | log storage/query backend |
| `fluentbit/` | log collection and forwarding |
| `grafana/` | datasource and dashboard provisioning |
| `scripts/` | ad-hoc operator helpers; not the canonical source of truth |

## Boundaries

- compose profiles and config files are canonical; helper scripts are secondary
- keep telemetry wiring deterministic and versioned in-repo
- if a dashboard or query depends on a label or exporter, update the relevant config in the same change

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
