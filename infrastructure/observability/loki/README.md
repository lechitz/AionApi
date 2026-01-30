# infrastructure/observability/loki

Loki configuration assets for log ingestion and querying. This folder is the target for log storage and correlation in the observability stack.

## Package Composition

- (Reserved for Loki config files and storage settings)

## Flow (Where it comes from -> Where it goes)

AionAPI logs -> Fluent Bit -> Loki -> Grafana Explore / dashboards

## Why It Was Designed This Way

- Keep log storage configuration separate from collection.
- Enable consistent log querying across environments.
- Support trace correlation via shared labels.

## Recommended Practices Visible Here

- Align labels (service, env, trace_id) with Fluent Bit output.
- Keep retention/storage settings explicit per environment.
- Avoid committing storage credentials.

## Differentials

- Trace-log correlation built into label strategy.

## What Should NOT Live Here

- Log parsers or collection rules (Fluent Bit).
- Application logging code.
- Secrets or production storage credentials.
