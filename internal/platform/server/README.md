# Platform Server Layer

**Path:** `internal/platform/server`

## Overview

Server-level platform composition layer.
Currently focused on HTTP stack composition, routing abstraction, middleware chain, and generic handlers.

## Subpackages

| Subpackage | Role |
| --- | --- |
| `http/` | HTTP server composition and transport infrastructure |

## Design Notes

- Keep this layer transport/platform-focused.
- Context modules register endpoints via platform ports, not router-specific APIs.
- Detailed HTTP behavior is documented in `internal/platform/server/http/README.md`.

## Package Improvements

- Add transport roadmap section if additional protocols are introduced.
- Add startup/shutdown lifecycle contract notes.
- Add links to integration tests for composed server behavior.
- Add extension guide for new server transports.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
