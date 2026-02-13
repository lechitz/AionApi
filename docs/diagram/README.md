# Architecture Diagrams

**Path:** `docs/diagram`

## Overview

This folder stores architecture and flow diagrams in two formats:
- editable source (`*.sequence.txt`)
- rendered output (`images/*.svg`)

## Diagram Catalog

| Source file | Rendered SVG | Purpose |
| --- | --- | --- |
| `cmd-aion-api.sequence.txt` | `images/cmd-aion-api.svg` | API entrypoint flow |
| `cmd-api-seed-caller.sequence.txt` | `images/cmd-api-seed-caller.svg` | Seed caller flow |
| `cmd-seed-helper.sequence.txt` | `images/cmd-seed-helper.svg` | Seed helper flow |
| `internal-admin.sequence.txt` | `images/internal-admin.svg` | Admin context flow |
| `internal-auth.sequence.txt` | `images/internal-auth.svg` | Auth context flow |
| `internal-adapter-primary-graphql.sequence.txt` | `images/internal-adapter-primary-graphql.svg` | Primary GraphQL adapter flow |
| `internal-platform.sequence.txt` | `images/internal-platform.svg` | Platform bootstrap flow |
| `internal-platform-server.sequence.txt` | `images/internal-platform-server.svg` | HTTP server routing flow |
| `internal-platform-httpclient.sequence.txt` | `images/internal-platform-httpclient.svg` | HTTP client flow |

## Editing Workflow

1. Edit `*.sequence.txt` source.
2. Render/export SVG (e.g., sequencediagram.org).
3. Replace corresponding file under `images/`.
4. Commit source and rendered artifact together.

## Package Improvements

- Add diagram naming conventions for new contexts/features.
- Add lightweight quality checklist (actors, direction, error path).
- Add script/target to validate missing source↔svg pairs.
- Add timestamp/version metadata in diagram headers for traceability.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
