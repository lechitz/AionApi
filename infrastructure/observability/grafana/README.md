# infrastructure/observability/grafana

Grafana provisioning assets for the AionAPI observability stack. These files define dashboards and data sources so Grafana boots with a ready-to-use view.

## Package Composition

- `datasources/`
  - Prometheus datasource definition (defaults to `prometheus-dev:9090` in dev).
- `dashboards/`
  - Versioned dashboard JSON files (latency, errors, and domain KPIs).
- `provisioning/`
  - Providers that load datasources and dashboards on Grafana startup.

## Flow (Where it comes from -> Where it goes)

Grafana boot -> provisioning/ -> datasources + dashboards -> Grafana UI

## Why It Was Designed This Way

- Keep observability assets versioned alongside the codebase.
- Make local onboarding fast with preloaded dashboards.
- Avoid manual Grafana setup or clickops.

## Recommended Practices Visible Here

- Keep dashboards in JSON and review diffs like code.
- Update datasource URLs per environment (dev/staging/prod).
- Keep panel naming consistent with metric labels for searchability.
- Avoid embedding secrets in dashboard JSON.

## Differentials

- Auto-provisioned dashboards and datasources, no manual setup.
- Reproducible Grafana state across environments.

## What Should NOT Live Here

- Ad-hoc dashboards created only in the UI.
- Business logic or service configuration.
- Secrets, tokens, or private credentials.
