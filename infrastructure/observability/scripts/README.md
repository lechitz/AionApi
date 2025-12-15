# Observability Scripts

Automates applying and validating observability improvements (Prometheus with exemplars, Jaeger datasource, and provisioned RED dashboard).

## Main script
- `setup-improvements.sh`: validates config files, restarts the dev Docker stack, and checks health/endpoints (Prometheus, Grafana, Jaeger, and OTel).

## Prerequisites
- Docker Desktop/Engine running with access to the daemon.
- Docker Compose plugin (`docker compose`).
- `.env.dev` in `infrastructure/docker/environments/dev`.
- `curl`; `jq` is optional for parsing API responses.

## How to run
1. From the project root: `./infrastructure/observability/scripts/setup-improvements.sh`
2. Wait for the script to restart the stack and finish validations.
3. Open Grafana at `http://localhost:3000` (aion/aion) and load the **AionAPI - RED Metrics Dashboard (Professional)**.
