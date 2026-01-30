# internal/admin

Administrative domain for privileged operations separated from end-user flows.

## Purpose and Main Capabilities

- Enforce admin-only actions (roles, block/unblock, privileged updates).
- Apply elevated authorization and auditing policies at a dedicated boundary.
- Provide a clean separation between admin behavior and user-facing flows.
- Expose admin operations via HTTP without leaking infra or core types.

## Package Composition

- `core/`: domain models, ports, and usecases (pure business logic).
- `core/ports/input`: admin service contracts implemented by usecases.
- `core/ports/output`: repository/integration contracts implemented by adapters.
- `core/usecase`: privileged workflows and policy enforcement.
- `adapter/primary/http`: admin HTTP handlers and DTO mapping.
- `adapter/secondary/db`: admin repositories and database mapping.

## Flow (Where it comes from -> Where it goes)

Admin HTTP request -> primary adapter -> input port -> usecase ->
output port -> secondary adapter -> database -> response

## Diagram

![Admin Domain Flow](../../docs/diagram/images/internal-admin.svg)

Source: `../../docs/diagram/internal-admin.sequence.txt`

## How It Works (Concise)

- HTTP handlers validate input and admin claims, open spans, and call input ports.
- Usecases enforce admin policy, orchestrate repositories, and return semantic errors.
- Secondary adapters translate infra errors and map persistence models.
- Adapters map domain results into transport responses (no domain leakage).

## Separation Inside the Bounded Context

- Core (domain, ports, usecases) never imports adapters or infrastructure.
- Primary adapters translate HTTP <-> core and own DTOs and validation.
- Secondary adapters implement output ports and isolate persistence details.
- Shared concerns (auth, logging, constants, errors) come from `internal/shared` and platform layers.

## Why It Was Designed This Way

- Keep admin permissions isolated from user routes.
- Enforce strict policies at dedicated boundaries.
- Allow auditing and trace correlation for sensitive actions.

## Recommended Practices Visible Here

- Validate admin claims in adapters and re-check in core.
- Emit audit-ready logs with trace correlation.
- Reuse shared ports when it avoids duplicating rules.
- Map driver/infra errors to semantic admin errors in adapters.
- Keep handlers thin: decode, validate, call usecase, map response.

## Differentials

- Dedicated admin boundary with explicit role documentation.
- Clear separation between transport, policy, and persistence within the context.

## What Should NOT Live Here

- End-user behavior or generic CRUD.
- UI-specific logic or transport DTOs outside adapters.
- Cross-context imports or shared state leakage.
