# Developer Utilities (`hack`)

**Path:** `hack`

## Overview

This folder contains development-only tools and scripts.
It follows the common `hack/` convention for non-production operational helpers.

## Subfolders

| Folder | Responsibility |
| --- | --- |
| `tools/` | Go CLIs for seed and utility workflows |
| `dev/` | Shell scripts for diagnostics and local troubleshooting |

## Usage Examples

```bash
make seed-api-caller
make seed-helper
go run ./hack/tools/graph-projection-export --user-id 999
make record-projection-smoke
make realtime-record-smoke
make record-projection-page-smoke
make ingest-event-smoke
make outbox-diagnose
make event-backbone-gate-preflight
make event-backbone-gate
bash hack/dev/test-chat.sh
```

## Tool Matrix

| Tool | Make target | Purpose |
| --- | --- | --- |
| `record-projection-smoke` | `make record-projection-smoke` | Validate `record -> outbox -> kafka -> projection` end-to-end |
| `realtime-record-smoke` | `make realtime-record-smoke` | Validate `record_projection_changed` SSE delivery after the derived row is materialized |
| `record-projection-page-smoke` | `make record-projection-page-smoke` | Validate cursor pagination over `recordProjections` derived read model |
| `ingest-event-smoke` | `make ingest-event-smoke` | Validate `aion-ingest -> kafka` envelope publication |
| `outbox-diagnose` | `make outbox-diagnose` | Inspect outbox backlog, oldest pending age, and sample pending/failed rows |
| `event-backbone-gate-preflight` | `make event-backbone-gate-preflight` | Verify local repo and service prerequisites before the full gate |
| `event-backbone-gate` | `make event-backbone-gate` | Run the official v2 records gate across `AionApi` and `aionapi-dashboard` |

## Design Notes

- Keep this folder out of production image/runtime paths.
- Use it for reproducible local workflows and debugging support.
- Domain logic must remain in `internal/`, not in `hack/` scripts/tools.

## Package Improvements

- Add a top-level command matrix mapping each utility to its Make target.
- Add minimal contributor guidelines for adding new scripts/tools.
- Add shared script helper library for common logging/error output.
- Add quick links to per-tool READMEs (`tools/seed-caller`, `tools/seed-helper`).

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
