# HTTP Response Helpers (Shared)

**Folder:** `internal/shared/httpresponse`

## Responsibility

* Give all controllers a **single place** to format HTTP responses.
* **Encode JSON** consistently and set the right headers (e.g., `Content-Type`).
* **Map domain/semantic errors** (`sharederrors`) to HTTP status codes in a uniform way.
* Attach **minimal, useful metadata** (e.g., timestamps; request IDs when available) to error responses for operability.

## How it works

* Exposes small helpers to:

    * write **success bodies** (JSON),
    * write **error bodies** from `sharederrors` (standard envelope + status),
    * set common headers.
* Integrates with `logger.ContextLogger` for **structured logs** and uses `commonkeys` to keep log/field names consistent.
* Treats domain errors as **first-class**, avoiding ad-hoc strings in handlers.

> Keep handlers **thin**: decode/validate → call input port → pass result/error to `httpresponse`.

## Error mapping (recommended)

`sharederrors` → HTTP status (typical mapping):

* `ErrInvalidInput`, validation issues → **400 Bad Request**
* `ErrUnauthorized` → **401 Unauthorized**
* `ErrForbidden` → **403 Forbidden**
* `ErrNotFound` → **404 Not Found**
* `ErrConflict` / uniqueness violations → **409 Conflict**
* `ErrTooManyRequests` → **429 Too Many Requests**
* `ErrInternal` / unknown → **500 Internal Server Error**

The helper centralizes this mapping so every adapter returns the **same** status and body shape.

## Conventions

* Always return **JSON**; set `Content-Type: application/json; charset=utf-8`.
* **Never** leak secrets (passwords, tokens, raw auth headers) into logs or bodies.
* Prefer **semantic domain errors** from `sharederrors`; avoid building HTTP status logic inside handlers.
* Include metadata that helps operability (e.g., `request_id`, `timestamp`) but avoid PII.

## Typical usage

### Success response

```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // ...decode + call input port...
    out := struct {
        ID   uint64 `json:"id"`
        Name string `json:"name"`
    }{ID: created.ID, Name: created.Name}

    // Example helper: write a 201 JSON response.
    // httpresponse.WriteJSON(w, http.StatusCreated, out)
}
```

### Error response (uniform mapping)

```go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    // ...decode + call input port...
    if err != nil {
        // Single call that:
        // 1) maps sharederrors -> status,
        // 2) emits a standardized error envelope,
        // 3) logs with ContextLogger (metadata only).
        // httpresponse.WriteError(r.Context(), w, h.Logger, err)
        return
    }
}
```

> The exact function names may differ in your implementation; the pattern above shows **where** and **how** to use this package.

## Envelope (illustrative)

Successful:

```json
{
  "data": { ... },
  "timestamp": "2024-06-16T09:02:00Z"
}
```

Error:

```json
{
  "error": {
    "code": "not_found",
    "message": "resource not found"
  },
  "meta": {
    "request_id": "...",
    "timestamp": "2024-06-16T09:02:00Z"
  }
}
```

*(Field names are representative—the package standardizes them so all adapters look alike.)*

## Testing hints

* Use `httptest.NewRecorder()` to assert:

    * correct **status code** for each `sharederrors` type,
    * **Content-Type** header is JSON,
    * error envelope fields (`error.code`, `meta.request_id`, etc.) are present.
* Table-driven tests make it easy to verify all mappings (400/401/403/404/409/429/500).

---

Centralizing response formatting and error mapping here keeps adapters **small, predictable, and easy to change** without touching domain code.
