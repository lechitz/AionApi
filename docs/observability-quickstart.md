# Quickstart: Apply Observability Improvements

Use this quick guide to load the new configs (Prometheus with exemplars, Jaeger datasource, RED dashboard) into the local stack.

## Steps
1. Prereqs: Docker running, `docker compose` plugin installed, `.env.dev` present in `infrastructure/docker/environments/dev`.
2. From the project root, run:
   ```bash
   ./infrastructure/observability/scripts/setup-improvements.sh
   ```
3. Wait for the script to restart the stack and complete automated checks (container health, Prometheus/Grafana/Jaeger endpoints, and datasource/dashboard provisioning).
4. Open Grafana at `http://localhost:3000` (login: `aion/aion`) and open **AionAPI - RED Metrics Dashboard (Professional)**.
5. Quick validation:
   - Panels show metrics (latency, error rate, throughput).
   - Click a point/bar → “View Traces in Jaeger” opens Jaeger search with filters pre-filled.
6. (Optional) Generate traffic to populate metrics:
   ```bash
   for i in {1..50}; do curl -s http://localhost:5001/aion/health > /dev/null; done
   ```
