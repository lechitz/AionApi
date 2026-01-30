# docs/assets

Static assets used by documentation pages and Swagger UI. These files are not bundled into the runtime binary.

## Package Composition

- `stylesheets/`
  - Theme overrides for Swagger UI and documentation pages.
- Images and other static resources referenced by docs.

## Flow (Where it comes from -> Where it goes)

Docs pages -> assets -> rendered docs / Swagger UI

## Why It Was Designed This Way

- Keep visual overrides alongside docs.
- Avoid build steps for docs assets.

## Recommended Practices Visible Here

- Keep assets lightweight and versioned.
- Use stable relative paths for GitHub Pages.
- Document visual changes that affect Swagger UI.

## Differentials (Rare but Valuable)

- Zero build pipeline for docs assets.

## What Should NOT Live Here

- Runtime assets for the API.
- Large binaries or build outputs.
