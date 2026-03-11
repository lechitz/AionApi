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
bash hack/dev/test-chat.sh
```

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
