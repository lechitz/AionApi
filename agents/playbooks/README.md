# Playbooks (aion-api)

This folder contains step-by-step workflows used by aion-api governance. Playbooks define **how** to execute common tasks safely and consistently.

## What You Will Find Here

- `feature-request-workflow.md`  Mandatory 5-stage governance pipeline
- `add-usecase.md`               How to add a new usecase
- `add-endpoint.md`              How to add a new GraphQL mutation/query
- `add-context.md`               How to add a new bounded context
- `repository-patterns.md`       Advanced repository patterns
- `cache-patterns.md`            Redis caching strategies
- `migrations-and-seeds.md`      Database migrations and seeds
- `add-middleware.md`            HTTP middleware and GraphQL directives

## How To Use

1) Identify the task type (feature, refactor, docs, tests).
2) Locate the matching playbook.
3) Follow steps in order; do not skip stages unless explicitly allowed.
4) If a step fails, use the playbook recovery guidance.

## Principles

- Playbooks are the source of truth for repeatable workflows.
- Consistency beats speed; follow the pipeline.
- Escalate when a task crosses architectural boundaries.

## Related Docs

- `agents/main.md`
- `agents/personas/`
- `agents/standards/`
