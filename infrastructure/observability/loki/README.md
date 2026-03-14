# Loki Configuration

**Path:** `infrastructure/observability/loki`

## Overview

Loki storage/query configuration for centralized log retention and exploration.
Logs are typically ingested via Fluent Bit and consumed in Grafana.

## Files

| File | Purpose |
| --- | --- |
| `loki.yaml` | Loki single-binary/server configuration |

## Design Notes

- Keep retention and storage config explicit by environment profile.
- Align labels with ingestion pipeline for searchable logs.
- Separate storage concerns from parser/collection concerns.

## Package Improvements

- Add retention policy rationale and expected volume guidance.
- Add startup verification checklist (`ready`, ingestion path, query path).
- Add notes for scaling beyond local single-binary mode.
- Add sanity tests for label cardinality controls.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
