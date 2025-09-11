# Admin Usecases (Core)

**Folder:** `internal/admin/core/usecase`

## Responsibility

* Implement the **Admin** domain logic behind **input ports** (e.g., `AdminService`).
* Orchestrate policies (RBAC/admin-only), validations, and coordination with **output ports**.
* Return **domain entities** and **semantic errors**; never leak transport/infra details.

## How it works

* Lives under `internal/admin/core/usecase` and **implements** interfaces from `internal/admin/core/ports/input`.
* Depends **only on ports** (DIP): concrete adapters (DB/cache/http/etc.) are injected by bootstrap.
* Every method accepts a `context.Context` for cancellation, deadlines, and **OpenTelemetry** propagation.
* Uses domain entities/VOs; mapping to/from persistence models is done in **secondary adapters**.
* Emits structured logs through the platform logger port and sets **canonical OTel attributes**.

## Typical flow

1. Receive validated input via an input port method.
2. Enforce **admin policies** (role/permissions), invariants, and cross-entity rules.
3. Query/update state through **output ports** (repositories, cache, feature flags, etc.).
4. Aggregate the result into domain types and return (or return a well-typed error).

## Example use cases (suggested)

* **User administration:** list/filter users, update roles, soft-delete/restore accounts.
* **Operational insights:** read aggregated counters/health from repositories.
* **Feature flags / switches:** toggle flags, scope to tenants/users, read effective config.
* **Audit helpers:** produce records/events for sensitive admin mutations (via an output port).

## Ports & dependencies (typical)

* `AdminRepository` (read/write admin projections & stats).
* `UserRepository` (admin actions over users).
* `FeatureFlagStore` (optional).
* `Auth/ZAuth` (to verify caller is admin).
* `Logger` (context logger) and **OTel** tracer.
* Optionally: `EventBus`/`Publisher` for audit streams.

> The service must speak only through ports. No ORM/HTTP/Redis code here.

## Error semantics

Prefer **semantic** errors over ad-hoc strings:

* `ErrUnauthorized`, `ErrForbidden` (RBAC)
* `ErrNotFound`, `ErrAlreadyExists`
* `ValidationError{field, reason}`
* `ErrUnexpected` for unknown failures

These are later mapped to transport codes by adapters.

## Tracing & logging

* Start an **OTel span** per use case (e.g., `SpanUpdateUserAsAdmin`).
* Add attributes: `operation`, `admin_id`, `target_user_id`, `status`.
* On error: `span.RecordError(err)` and set `codes.Error`; log **metadata only** (no PII/secrets).

## Design notes

* Keep policies **inside** the usecase layer (not handlers or repositories).
* Use **UTC** timestamps; keep side effects small and explicit.
* Avoid N+1 by exposing repository methods that fit admin views (aggregations/pagination).
* Remain **transport-agnostic** (HTTP/GraphQL/CLI shouldn’t leak in).

## Testing hints

* **Table-driven** unit tests with **mocked output ports**.
* Cover:

    * RBAC failures (non-admin, missing claims).
    * Validation errors (bad filters, invalid role transitions).
    * Happy paths (including pagination & sorting).
    * Error propagation from repositories/feature-flag store.
* Assert **no infra coupling**: only ports are called; no SQL/HTTP in tests.
