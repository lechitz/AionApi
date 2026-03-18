# API Collections

**Path:** `docs/collections`

## Purpose

This folder publishes consumer-friendly request collections for manual QA and local integration checks.
The canonical REST contract still lives in `contracts/openapi/swagger.yaml`.

## Current Artifact

| Path | Purpose |
| --- | --- |
| `postman/AionApi.postman_collection.json` | Postman collection covering auth, user, admin, chat, GraphQL, and health flows |

## Usage

- import the collection into Postman
- set `{{baseURL}}` to the target API origin
- use collection variables or cookies for auth; never commit real tokens back into the file

## Boundaries

- update this artifact when consumer-facing REST flows change materially
- do not treat the collection as a source of truth over OpenAPI or runtime behavior
- keep secrets, personal environments, and local-only values outside the checked-in JSON

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
