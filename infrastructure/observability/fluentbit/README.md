# Fluent Bit Configuration

**Path:** `infrastructure/observability/fluentbit`

## Overview

Fluent Bit pipeline configuration for collecting and forwarding container/application logs.
In local stack, it forwards logs to Loki for querying in Grafana.

## Files

| File | Purpose |
| --- | --- |
| `fluent-bit.conf` | Input/filter/output pipeline definition |
| `parsers.conf` | Log parser configuration |

## Design Notes

- Keep labels aligned with tracing/log correlation strategy.
- Keep parser changes backward compatible with dashboards/queries.
- Avoid credentials in config files.

## Package Improvements

- Add sample log record before/after parsing in docs.
- Add validation script for config syntax.
- Add guidance for adding new log sources safely.
- Add alerting recommendations for dropped/failed log events.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
