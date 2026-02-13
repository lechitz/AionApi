# Makefile Modules

**Path:** `makefiles`

## Overview

Modular Makefile fragments included by the root `Makefile`.
They centralize commands for build, environment, migrations, codegen, tests, and quality workflows.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Command modularization | Split make targets by concern/domain |
| Developer workflow | Expose reproducible local commands |
| CI parity | Keep local verify steps aligned with pipeline checks |

## Typical Command Areas

- Docker/runtime
- Migrations/seeds
- Codegen (GraphQL/mocks)
- Tests/coverage
- Lint/format/verify

## Design Notes

- Keep root `Makefile` thin; logic belongs in modules.
- Keep target naming predictable and discoverable.
- Keep environment variable requirements documented.

## Package Improvements

- Add generated command reference table from make metadata.
- Add dependency map between major targets (`verify`, `test`, `graphql`, etc.).
- Add target stability policy for CI-consumed commands.
- Add troubleshooting section for common tooling failures.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
