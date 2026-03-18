# User Bounded Context

**Path:** `internal/user`

## Purpose

`internal/user` owns account lifecycle, public registration, profile/avatar management, password changes, and soft deletion.

## Current HTTP Surface

### Public routes

- `POST /user/create`
- `POST /user/avatar/upload`
- `POST /registration/start`
- `PUT /registration/{registration_id}/profile`
- `PUT /registration/{registration_id}/avatar`
- `POST /registration/{registration_id}/complete`

### Authenticated routes

- `GET /user/all`
- `GET /user/me`
- `GET /user/{user_id}`
- `PUT /user/`
- `DELETE /user/avatar`
- `PUT /user/password`
- `DELETE /user/`

## Runtime Contract

- password updates refresh the auth cookie/token on success
- cache layers must never store password hashes or raw passwords
- registration is a multi-step public flow separate from the authenticated profile-update surface
- avatar upload/removal is owned here even when backed by external object storage

## Boundaries

- identity, password, and registration rules stay in core usecases
- transport adapters own request decoding, auth context extraction, and cookie refresh wiring
- auth/session semantics collaborate with `internal/auth`, but profile and account ownership stay here

## Related Docs

- [`../auth/README.md`](../auth/README.md)
- [`../platform/server/http/utils/cookies/README.md`](../platform/server/http/utils/cookies/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
