# Docker Infrastructure

**Path:** `infrastructure/docker`

## Overview

Docker build/runtime assets for local and production-like environments.
This package defines container images, compose profiles, and environment wiring.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `environments/` | Profile-specific compose/env definitions |
| root `Dockerfile` | Production-oriented API image build |

## Design Notes

- Keep environment-specific details inside profile folders.
- Keep production and development image concerns separated.
- Keep compose/runtime behavior reproducible via Make targets.

## Package Improvements

- Add profile matrix documenting service composition and ports.
- Add build cache strategy notes for faster local rebuilds.
- Add startup dependency ordering notes for local stack reliability.
- Add security checklist for production image hardening.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
