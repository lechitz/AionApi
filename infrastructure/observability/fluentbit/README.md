# infrastructure/observability/fluentbit

Fluent Bit configuration for collecting and forwarding logs from the AionAPI runtime to Loki (or other sinks).

## Package Composition

- `fluent-bit.conf`
  - Inputs, filters, and Loki output for Docker logs.
- `parsers.conf`
  - JSON parsing for container log payloads.

## Flow (Where it comes from -> Where it goes)

Container logs -> Fluent Bit -> Loki -> Grafana

## Why It Was Designed This Way

- Keep log parsing and routing centralized.
- Support structured logging and correlation at the edge.
- Decouple log collection from storage.

## Recommended Practices Visible Here

- Parse JSON logs to keep fields queryable.
- Enrich with service/env/trace labels for correlation.
- Document any field mapping changes to keep queries stable.

## Differentials

- Correlation-ready log pipeline aligned with tracing labels.

## What Should NOT Live Here

- Loki storage config.
- Application logging code.
- Secrets or credentials.
