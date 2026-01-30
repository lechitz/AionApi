# internal/user

User lifecycle domain (create, query, update profile/password, soft delete) exposed via HTTP/GraphQL.

## Purpose and Main Capabilities

- Create and manage users with validation and normalization in core.
- Enforce uniqueness for username/email.
- Handle password hashing and token issuance via ports.
- Provide soft delete and safe listing/lookup.

## Package Composition

- `core/`: user entities, ports, and usecases.
- `core/ports/input`: user service interface for adapters.
- `core/ports/output`: repository, hasher, token provider, auth store.
- `core/usecase`: Create, GetByID, GetByUsername, ListAll, Update, UpdatePassword, SoftDelete.
- `adapter/primary`: HTTP/GraphQL controllers and DTO mapping.
- `adapter/secondary`: db repositories, storage, hasher, token/cache adapters.

## Flow (Where it comes from -> Where it goes)

HTTP/GraphQL request -> primary adapter -> input port -> usecase ->
output ports -> secondary adapters -> database/cache/storage -> response

## How It Works (Concise)

- Create: validate and normalize fields, check uniqueness, hash password, persist user.
- Update profile: accept partial updates, reject empty changes, persist updates.
- Update password: verify current password, re-hash, update, issue/store new tokens.
- Soft delete: revoke tokens in auth store, mark user as deleted in repository.

## Separation Inside the Bounded Context

- Core is transport-agnostic and owns validation and security rules.
- Primary adapters map DTOs and handle HTTP/GraphQL specifics.
- Secondary adapters isolate infrastructure (db, storage, cache, hashing).

## Why It Was Designed This Way

- Keep security-sensitive rules centralized in core.
- Make adapters thin and replaceable.
- Preserve consistent behavior across transports.

## Recommended Practices Visible Here

- Never log PII, hashes, or tokens; log IDs only.
- Return semantic errors for adapters to map (e.g., username in use).
- Use OTel spans per operation with standard attributes.

## Differentials

- Password and token handling via explicit output ports.

## What Should NOT Live Here

- Business logic in adapters.
- Infra types inside core.
- Cross-context imports.
