# Personas (AionAPI)

This folder defines the role-based agents used by AionAPI governance. Each persona has a clear scope, authority level, and decision boundaries.

## What You Will Find Here

- `aion.architect.md`  Architectural integrity and layer placement
- `aion.platform.md`   Operability, CI/CD, configs, and infra readiness
- `aion.reviewer.md`   Code quality and Shizen (simplicity) guardrails
- `aion.tester.md`     Test coverage and reliability focus
- `aion.developer.md`  Implementation within approved boundaries
- `aion.debugger.md`   Troubleshooting and incident response
- `aion.docs.md`       Documentation standards and package READMEs

## How To Use

1) Identify task type (feature, refactor, bug, tests, docs).
2) Pick the highest-authority persona that applies.
3) Follow its rules before writing code or docs.
4) Escalate when a decision crosses layers or affects operability.

## Key Principles

- Respect the authority hierarchy in `agents/main.md`.
- Keep scope tight: each persona exists to prevent drift.
- Favor clarity over cleverness.

## When To Add A New Persona

Create a new persona only if:
- The task category repeats across the repo,
- The existing personas do not cover it well,
- You can define clear boundaries and veto rules.

## Related Docs

- `agents/main.md`
- `agents/playbooks/`
- `agents/standards/`
