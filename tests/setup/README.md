# Test Setup Helpers

**Path:** `tests/setup`

## Overview

Shared test builders and helpers for unit testing service/usecase layers.
These helpers reduce boilerplate around gomock setup and default fixtures.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Suite builders | Build ready-to-use test suites per domain service |
| Mock wiring | Create and expose required mocks consistently |
| Shared fixtures | Provide default entities/helpers for tests |

## Usage Pattern

```go
suite := setup.UserServiceTest(t)
defer suite.Ctrl.Finish()
```

## Design Notes

- Keep helpers deterministic and focused.
- Favor explicit fixture builders over hidden global state.
- Keep suite APIs stable across refactors.

## Package Improvements

- Add helper index table (builder -> mocks -> SUT).
- Add extension guide for adding new service test suites.
- Add stricter naming conventions for suite fields.
- Add examples covering error-path assertions.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
