# Bootstrap (Domain Contract)

**Path:** `internal/platform/bootstrap`

This package now only defines the `AppDependencies` contract used by the HTTP/GraphQL composers and the Fx wiring. Domain wiring lives in `internal/platform/fxapp`.

## Responsibilities
- Expose a single dependency bundle (`*AppDependencies`) consumed by transports.
- Remain DI-agnostic so primary adapters don't depend on Fx internals.

## Notes
- The former `InitializeDependencies` implementation was replaced by Fx modules under `internal/platform/fxapp`.
- Extend `AppDependencies` if new domain services need to be exposed to transports.
