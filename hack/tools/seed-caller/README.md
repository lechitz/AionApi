# Seed Caller Tool

**Path:** `hack/tools/seed-caller`

## Overview

This CLI seeds data by calling the real API surface (login + GraphQL mutations).
It is intended for local smoke checks, repeatable data setup, and observability validation without direct DB writes.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Authenticated seeding | Login and execute GraphQL seeding flows |
| Multi-user load setup | Generate deterministic user-based data batches |
| Cleanup support | Optional clean modes for safe reruns |
| Diagnostics | Optional debug output for payload/troubleshooting |

## Main Artifacts

| File/Area | Purpose |
| --- | --- |
| `main.go` | CLI entrypoint and workflow wiring |
| HTTP/GraphQL workflow code | Login, optional user creation, category/tag/record seeding |
| Success log file | Persist successful run metadata |

## Runtime Flow

1. Read CLI/environment configuration.
2. Authenticate user against API.
3. Optionally auto-create user if login fails.
4. Execute GraphQL mutations for categories/tags/records.
5. Optionally clean previous data depending on flags.
6. Persist success run marker to configured log.

## Quick Run

```bash
make seed-api-caller
# or
go run ./hack/tools/seed-caller
```

## Key Environment Variables

| Variable | Description |
| --- | --- |
| `API_CALLER_HOST` | API base host |
| `API_CALLER_CONTEXT` / `API_CALLER_ROOT` | Base context and API root |
| `API_CALLER_GRAPHQL` | GraphQL path |
| `API_CALLER_USER`, `API_CALLER_PASS` | Login credentials |
| `API_CALLER_COUNT`, `API_CALLER_USER_PREFIX` | Multi-user batch controls |
| `API_CALLER_AUTO_CREATE` | Auto-create user when login fails |
| `API_CALLER_CLEAN`, `API_CALLER_ONLY_CLEAN` | Cleanup controls |
| `API_CALLER_DEBUG` | Enable request/response debug output |

## Design Notes

- This tool validates real API contracts and auth rules, not DB internals.
- Keep it dev-focused; it should not become production automation.
- Prefer Make targets for consistent usage across developers.

## Package Improvements

- Add table-driven tests for configuration parsing and cleanup behavior.
- Add rate-limit/backoff options for safer high-volume local runs.
- Emit structured run summary (JSON) for easier CI/local reporting.
- Add a dry-run mode for validating configuration without side effects.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
