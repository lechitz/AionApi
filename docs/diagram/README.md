# Architecture Diagrams

All diagrams live in this folder as:
- SVG images under `images/`
- Editable sources as `*.sequence.txt` (sequencediagram.org compatible)

## Available Diagrams

- `cmd-aion-api.sequence.txt`
  - Image: `images/cmd-aion-api.svg`
  - Purpose: Entry-point flow for the AionAPI process
- `cmd-api-seed-caller.sequence.txt`
  - Image: `images/cmd-api-seed-caller.svg`
  - Purpose: API seeding via authenticated endpoints
- `cmd-seed-helper.sequence.txt`
  - Image: `images/cmd-seed-helper.svg`
  - Purpose: Local seed env, tokens, and bcrypt hashes

## How to View

1) Open `images/*.svg` in your browser or on GitHub.
2) If you need to edit a diagram, open the matching `*.sequence.txt` file in a text editor.
3) Paste into https://sequencediagram.org/ and export as SVG.

## Syntax (Short)

title Example Flow

participant "User" as U
participant "Service" as S

U->S: Call endpoint
S->U: Response

Notes:
- Use `title` for the diagram title.
- Define participants with readable names.
- Keep flows short and focused.
