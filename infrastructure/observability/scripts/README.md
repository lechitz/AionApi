# infrastructure/observability/scripts

Automation scripts for applying and validating observability improvements in the dev stack.

## Package Composition

- `setup-improvements.sh`
  - Validates configs, restarts the Docker stack, and checks endpoints.

## Flow (Where it comes from -> Where it goes)

Developer run -> script validation -> docker compose restart -> observability endpoints

## Why It Was Designed This Way

- Make observability setup repeatable.
- Reduce manual steps and drift.
- Provide fast feedback on health and wiring.

## Recommended Practices Visible Here

- Keep scripts idempotent and safe to re-run.
- Validate endpoints before declaring success.
- Avoid embedding credentials in scripts.

## Differentials

- Automated validation of Grafana/Prometheus/Jaeger/OTel wiring.

## What Should NOT Live Here

- Production deployment logic.
- Environment secrets or tokens.
- Application runtime logic.

## How to Run

```bash
./infrastructure/observability/scripts/setup-improvements.sh
```
