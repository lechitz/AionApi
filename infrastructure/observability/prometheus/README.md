# infrastructure/observability/prometheus

Prometheus configuration for scraping metrics emitted by the OTel Collector (and optionally other services).

## Package Composition

- `prometheus.yml`
  - Scrape targets and intervals for the dev stack.

## Flow (Where it comes from -> Where it goes)

OTel Collector -> Prometheus scrape -> Grafana dashboards

## Why It Was Designed This Way

- Keep metric scraping explicit and versioned.
- Separate data collection from visualization.
- Make dev observability deterministic.

## Recommended Practices Visible Here

- Keep scrape intervals reasonable for the environment (dev vs prod).
- Preserve label conventions used by Grafana dashboards.
- Validate targets after changes (`/targets` endpoint).

## Differentials

- Single-source scraping aligned with OTel metrics export.

## What Should NOT Live Here

- Grafana dashboards or provisioning configs.
- Application metrics code or business logic.
- Credentials or auth tokens.
