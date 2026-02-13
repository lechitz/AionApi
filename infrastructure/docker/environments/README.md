# Docker Environment Profiles

**Path:** `infrastructure/docker/environments`

## Overview

Environment profiles for Docker-based runtime setups.
This folder defines profile-specific compose and env configurations (dev/prod/personal/template).

## Profiles

| Profile | Purpose |
| --- | --- |
| `dev/` | Shared local development profile |
| `prod/` | Production-like compose profile |
| `example/` | Template baseline for new profiles |
| `my/` | Personal isolated local profile |

## Design Notes

- Profile READMEs document implementation details per environment.
- Keep profile-specific secrets outside repository or in ignored files.
- Keep profile names and make targets aligned for discoverability.

## Package Improvements

- Add profile capability matrix (services, ports, persistence, hot-reload).
- Add script to validate required env vars across all profiles.
- Add explicit docs on when to use `dev` vs `my` profile.
- Add centralized list of active compose files/entrypoints.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
