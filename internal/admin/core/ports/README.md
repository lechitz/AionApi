# internal/admin/core/ports

Admin ports define the contracts between core and adapters.

## Purpose & Main Capabilities

- Declare input ports implemented by usecases.
- Declare output ports implemented by adapters.

## Package Composition

- `input/`
  - Admin service interfaces.
- `output/`
  - Repository and integration interfaces.

## Flow (Where it comes from -> Where it goes)

Primary adapter -> input port -> usecase -> output port -> secondary adapter

## Why It Was Designed This Way

- Enforce dependency inversion for admin logic.
- Make integrations swappable.

## Recommended Practices Visible Here

- Keep ports small and domain-focused.
- Avoid leaking transport or infra types.

## Differentials

- Clear boundary between admin core and adapters.

## What Should NOT Live Here

- Implementations or adapter wiring.
