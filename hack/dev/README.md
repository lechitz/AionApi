# Development Scripts

**Path:** `hack/dev`

## Overview

This folder contains operational helper scripts for local diagnostics and troubleshooting.
Scripts here are intended for manual dev workflows, not production automation.

## Script Inventory

| Script | Purpose |
| --- | --- |
| `check-and-setup-ollama.sh` | Ensure Ollama runtime/model availability |
| `diagnostico-ollama.sh` | End-to-end Ollama + Aion Chat diagnostics |
| `force-insert-roles.sh` | Local DB role bootstrap/repair |
| `test-aion-chat.sh` | Functional chat flow validation |
| `test-chat.sh` | Basic chat API health/sample checks |
| `test-gpu.sh` | GPU readiness/performance checks |

## Usage Pattern

```bash
bash hack/dev/test-chat.sh
```

## Design Notes

- Keep scripts idempotent where possible.
- Prefer environment variables over hardcoded secrets.
- Promote stable routines into Make targets when repeated often.

## Package Improvements

- Add script-level help (`-h/--help`) for input expectations.
- Standardize shell strict mode (`set -euo pipefail`) across scripts.
- Add a small script matrix (prerequisites + expected output) in this README.
- Add lint checks (`shellcheck`) to CI for this folder.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
