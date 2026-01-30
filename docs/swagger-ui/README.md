# docs/swagger-ui

Static Swagger UI bundle used to explore the generated OpenAPI contract under `docs/swagger`.

## Package Composition

- `index.html`
  - Points to the generated swagger file.
- `swagger-ui-bundle.js`, `swagger-ui.css`
  - Swagger UI distribution files.
- `custom.css`
  - Local style overrides.

## Flow (Where it comes from -> Where it goes)

Generated OpenAPI -> Swagger UI -> API exploration

## Why It Was Designed This Way

- Keep API exploration static and easy to host.
- Allow custom styling without rebuilding Swagger UI.

## Recommended Practices Visible Here

- Regenerate OpenAPI before browsing the UI.
- Keep relative paths stable for GitHub Pages.
- Validate compatibility when bumping Swagger UI versions.

## Differentials (Rare but Valuable)

- Fully static UI; no runtime dependencies.

## What Should NOT Live Here

- Generated OpenAPI files (keep in `docs/swagger`).
- App runtime assets.
