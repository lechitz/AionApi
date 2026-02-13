# Seed Helper Tool

**Path:** `hack/tools/seed-helper`

## Overview

This CLI generates local seed artifacts such as JWT tokens and bcrypt hashes.
It supports seed bootstrap flows used by Make targets and SQL seed scripts.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Seed env generation | Create local `.env` seed values |
| Token generation | Build JWT tokens for test users |
| Password hash generation | Produce bcrypt hashes compatible with runtime auth |

## Main Commands

| Command | Purpose |
| --- | --- |
| `generate-env` | Creates local seed env file (`infrastructure/db/seed/.env.local`) |
| `generate-token` | Generates JWT for a user ID |
| `generate-bcrypt` | Generates bcrypt hash for a plain password |

## Quick Run

```bash
make seed-helper
./bin/seed-helper generate-env 10
```

## Integration Points

- Used by Make targets (`seed-helper`, `seed-setup`, `seed-quick`).
- Values feed SQL seed scripts under `infrastructure/db/seed/`.

## Design Notes

- Keeps crypto/token generation aligned with backend libraries.
- Avoid committing generated artifacts or sensitive values.
- Keep this package strictly local-dev oriented.

## Package Improvements

- Add tests for command parsing and invalid argument handling.
- Add explicit output format docs for each command.
- Support writing to custom output path for CI/local sandbox runs.
- Add optional JSON output mode for automation scripts.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
