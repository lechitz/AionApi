# internal/record

Record domain for diary/habit entries, exposed via GraphQL/HTTP with strict core rules.

## Purpose and Main Capabilities

- Create, list, retrieve, update, and soft delete records.
- Enforce required fields and relationships (user/category/tag).
- Keep record rules centralized in core.
- Provide consistent semantic errors for transport mapping.

## Package Composition

- `core/`: record entities, ports, and usecases.
- `core/ports/input`: record service interface for adapters.
- `core/ports/output`: repository and related lookups.
- `core/usecase`: create/list/get/update/soft delete flows.
- `adapter/primary`: GraphQL/HTTP controllers and DTO mapping.
- `adapter/secondary/db`: record persistence adapters and mappers.

## Flow (Where it comes from -> Where it goes)

GraphQL/HTTP request -> primary adapter -> input port -> usecase ->
output port -> db adapter -> database -> response

## How It Works (Concise)

- Adapters map transport inputs to commands and open spans.
- Usecases validate invariants, build partial updates, and orchestrate repositories.
- Soft delete preserves history and avoids hard deletes.

## Why It Was Designed This Way

- Maintain consistent rules across different transports.
- Keep ORM and persistence details out of core.
- Enable safe evolution of record fields and relationships.

## Recommended Practices Visible Here

- Validate relationships before persist when required.
- Use semantic errors and avoid leaking driver errors.
- Emit spans with `record_id`, `user_id`, `category_id` metadata only.

## Differentials

- Soft delete by default to preserve historical integrity.

## What Should NOT Live Here

- Business logic inside adapters.
- Transport DTOs inside core.
- Cross-context imports.
