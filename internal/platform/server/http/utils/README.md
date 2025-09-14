# Handler Helpers (Shared)

**Folder:** `internal/shared/handlerhelpers`

## Responsibility

* Give HTTP/GQL controllers a **single, consistent way** to validate inputs and write responses.
* Centralize **error → HTTP mapping**, **tracing attributes**, and **structured logs** so handlers stay thin.
* Provide small **validation/parsing helpers** that are transport-agnostic.

## How it works

* Response helpers wrap `httpresponse` and `sharederrors`, attaching **OTel attributes** from `tracingkeys` and logging metadata with `logger.ContextLogger` (using keys from `commonkeys`).
* Validation helpers run at the **transport boundary** (controllers/DTOs), keeping the domain clean.

> No persistence/ORM or business logic here—only reusable, side-effect-light utilities.

## What’s inside

* `response.go` — helpers to standardize success/error responses and emit tracing/logs with canonical keys.
* `validation.go` — small, reusable validators (e.g., required fields checker) intended for DTO/boundary validation.

## Conventions

* **Validate early** (in DTOs/handlers) and return **semantic domain errors** (`sharederrors`) instead of ad-hoc strings.
* **Never** log sensitive payloads (passwords/tokens). Prefer **metadata** (IDs, counts, statuses).
* Always set OTel attributes like status code, request ID, and operation name via `tracingkeys` + `commonkeys`.

## Examples

### Required fields at the boundary

```go
if err := handlerhelpers.CheckRequiredFields(map[string]string{
    "username": req.Username,
    "password": req.Password,
}); err != nil {
    // Convert to a bad-request using your standard response writer
    // (e.g., httpresponse + sharederrors) and return.
    return
}
```

### Logging + tracing-friendly error responses

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ... decode/validate ...

    // On domain/validation error:
    // handlerhelpers will (a) map to the proper HTTP status,
    // (b) attach OTel span attributes (status_code, request_id, etc.),
    // (c) write a standardized error body.
    // Example usage mirrors your local helpers pattern.
    // handlerhelpers.WriteDomainError(ctx, w, err, h.Logger)
}
```

*(Use the exact helper(s) exposed in `response.go` per your project’s pattern.)*

## Design notes

* Keys come from `internal/shared/constants`:

    * `commonkeys` for log/HTTP fields,
    * `tracingkeys` for span attributes.
* Keep helpers **pure and small** so they are easy to test and reuse across adapters (HTTP & GraphQL).

## Testing hints

* Use `httptest.NewRecorder()` to assert:

    * status code mapping,
    * error envelope shape,
    * presence of request/trace IDs in headers (when applicable).
* With gomock, assert **metadata keys** rather than free-form strings:

```go
logger.EXPECT().Errorw(gomock.Any(), commonkeys.UserID, gomock.Any()).AnyTimes()
```

---

Keep handlers **thin**: decode + validate → call input port → map/return response. These helpers make the “thin” part effortless and consistent.
