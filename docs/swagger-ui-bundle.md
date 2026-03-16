# Swagger UI Bundle Reference

This page documents how the static Swagger UI bundle is published in GitHub Pages.

## Location

- Static bundle directory: `docs/swagger-ui/`
- Entry file: `docs/swagger-ui/index.html`
- OpenAPI source URL used by UI: `../swagger/swagger.yaml` (published as `/swagger/swagger.yaml` on GitHub Pages)

## Files in the Bundle

| File | Purpose |
| --- | --- |
| `index.html` | Swagger UI entrypoint and runtime config |
| `swagger-ui-bundle.js` | Upstream Swagger UI JavaScript bundle |
| `swagger-ui.css` | Upstream Swagger UI styles |
| `custom.css` | Aion visual customization (dark palette, branding) |

## Maintenance Workflow

1. Regenerate/update OpenAPI contract (`swagger/swagger.yaml`).
2. Validate interactive docs at `/swagger-ui/`.
3. If Swagger UI assets are upgraded, retest authorization, execute button, and response rendering.
4. Keep custom styles minimal and focused on branding/usability.

## Quality Checklist

- OpenAPI URL resolves from GitHub Pages.
- Authorization persists when enabled.
- Endpoint groups are readable in dark mode.
- Try-it-out controls and response blocks remain accessible.
