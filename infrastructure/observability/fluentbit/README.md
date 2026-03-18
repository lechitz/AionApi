# Fluent Bit Configuration

**Path:** `infrastructure/observability/fluentbit`

## Purpose

This package collects container logs from the Docker host and forwards them to Loki.

## Current Flow

| File | Responsibility |
| --- | --- |
| `fluent-bit.conf` | tail Docker JSON logs, merge Docker metadata, add `APP_ENV`, ship to Loki |
| `parsers.conf` | Docker log parsing rules |

## Boundaries

- labels and fields added here must stay compatible with Grafana and Loki queries
- this package owns log collection and forwarding only; retention and query behavior live in Loki
- keep secrets and per-machine overrides out of the checked-in config

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
