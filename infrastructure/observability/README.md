# Observability Infrastructure

**Path:** `infrastructure/observability`

## Overview

Infrastructure assets for logs, metrics, and traces.
This package composes collector, storage, scraping, and visualization layers.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `otel/` | Telemetry ingestion and routing |
| `prometheus/` | Metrics scraping |
| `loki/` | Log storage/query backend |
| `fluentbit/` | Log collection/forwarding |
| `grafana/` | Dashboards and datasource provisioning |

## Design Notes

- Keep each observability concern isolated by component.
- Use in-repo config for deterministic local environments.
- Align naming conventions across logs, metrics, and traces.

## Package Improvements

- Add end-to-end observability startup verification runbook.
- Add architecture diagram of signal flow by component.
- Add environment-specific override strategy documentation.
- Add minimum supported tool versions (Grafana/Prometheus/Loki/OTel).

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
