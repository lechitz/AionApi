# REST API

This page keeps REST documentation in the same AionDocs experience while preserving full Swagger interactivity.

## Quick Access

- Embedded explorer (below)
- Direct page: <https://lechitz.github.io/aion-api/swagger-ui/>
- Raw OpenAPI: <https://raw.githubusercontent.com/lechitz/aion-api/main/swagger/swagger.yaml>

## Authentication

Most protected endpoints require Bearer token:

```http
Authorization: Bearer <JWT_TOKEN>
```

## Interactive Explorer

<div style="border:1px solid #cfd8e3;border-radius:12px;overflow:hidden;height:78vh;background:#fff;">
  <iframe
    src="../swagger-ui/"
    title="aion-api Swagger UI"
    style="width:100%;height:100%;border:0;"
    loading="lazy">
  </iframe>
</div>

## Notes

- If the embedded frame is blocked by browser policy, open the direct Swagger URL.
- Contract updates come from `swagger/swagger.yaml`.
- Use `make docs-verify` to validate site consistency before publishing.
