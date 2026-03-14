# AionAPI — Project Board Reference

**Board:** https://github.com/users/lechitz/projects/1
**Repository:** https://github.com/lechitz/AionApi
**Last updated:** 2026-03-14

This document is the canonical reference for how the AionAPI project board is structured. It is intended for both human contributors and AI agents (Codex, Copilot) that need to understand the issue hierarchy, labeling system, naming conventions, and board views before creating or triaging issues.

---

## 1. Issue Hierarchy

The board follows a strict **four-level hierarchy**:

```
Release
  └─ Master Epic
        └─ Epic
              └─ User Story / Feature / Refactor / Docs / Tests / Bugfix / DevOps / Database
```

### Level 0 — Release

An umbrella card that tracks alignment across **all Aion repositories** for a named release milestone.

- Label: `Release`
- Title format: `Release vX.Y.Z`
- One per release cycle, spans multiple master-epics
- Must reference all repos involved: AionApi, aion-chat, AionApi-dashboard, Aion-ingest, Aion-streams
- Does **not** decompose into sub-issues directly; sub-issues link back to it

Current: `#206 Release v1.0.0` — OPEN

---

### Level 1 — Master Epic

A strategic umbrella that groups a **domain pillar** for a release cycle. Master Epics do not contain implementation details; they aggregate epics.

- Label: `master-epic` (always combined with a domain label: `domain`, `infra`, or `docs`)
- Title format: `` `[ㅤ<Name> - X.0ㅤ]` `` (uses Unicode padding characters for visual alignment in the board)
- Assignee: `lechitz`

| Issue | Title | Domain Label | Status |
|-------|-------|-------------|--------|
| #88   | `[ㅤDocumentation - 1.0ㅤ]`            | `docs`   | closed |
| #123  | `[ㅤBusiness Rules & Entities - 1.0ㅤ]` | `domain` | closed |
| #124  | `[ㅤInfrastructure - 1.0ㅤ]`            | `infra`  | closed |

---

### Level 2 — Epic

A scoped grouping for a **specific module or concern** that has its own lifecycle. Epics aggregate user stories and implementation tasks.

- Label: `epic`
- Title format: `` `[ㅤ<Module Name>ㅤ]` `` (uses Unicode padding characters)
- Assignee: `lechitz`

| Issue | Title | Status |
|-------|-------|--------|
| #86   | `[ㅤAuthentication & Authorizationㅤ]` | closed |
| #87   | `[ㅤObservabilityㅤ]` | closed |
| #89   | `[ㅤUser Moduleㅤ]` | closed |
| #95   | `[ㅤDevOps Infrastructure & CI/CDㅤ]` | closed |
| #113  | `[ㅤTAGs Moduleㅤ]` | closed |
| #114  | `[ㅤCategories Moduleㅤ]` | closed |
| #116  | `[ㅤProfessional Diary Moduleㅤ]` | closed |
| #117  | `[ㅤIntentions Moduleㅤ]` | closed |
| #118  | `[ㅤMoods Moduleㅤ]` | closed |
| #119  | `[ㅤEnergy Moduleㅤ]` | closed |
| #120  | `[ㅤPersonal Diary Moduleㅤ]` | closed |
| #121  | `[ㅤWater Intake Moduleㅤ]` | closed |
| #122  | `[ㅤDay TAG Summary Moduleㅤ]` | closed |
| #125  | `[ㅤDatabaseㅤ]` | closed |
| #128  | `[ㅤCrosscutting Concernsㅤ]` | closed |
| #192  | `[ㅤRecord Moduleㅤ]` | closed |
| #195  | `[ㅤChat-IAㅤ]` | closed |

---

### Level 3 — Leaf Issues

Leaf issues are the actual work items. Each belongs to one of the following types:

| Label | Template File | Purpose | Title Convention |
|-------|--------------|---------|-----------------|
| `user story` | `user-story.md` | Flow-level milestone; describes a complete capability | `` `Category` - Title `` |
| `feature` | `feature.md` | New functionality implementation task | Short imperative description |
| `refactor` | `refactor.md` | Code improvement without behavior change | Short imperative description |
| `docs` | `docs.md` | Documentation creation or update | Short imperative description |
| `tests` | `tests.md` | Test coverage task | Short imperative description |
| `bugfix` | `bugfixes.md` | Bug fix | Short imperative description |
| `devops` | `devops.md` | CI/CD, pipeline, infrastructure task | Short imperative description |
| `database` | `database.md` | Schema, migration, index change | Short imperative description |

---

## 2. Labels Reference

All labels used in the board:

| Label | Category | Description |
|-------|----------|-------------|
| `Release` | Milestone | Cross-repo release umbrella |
| `master-epic` | Hierarchy | Strategic domain pillar grouping epics |
| `epic` | Hierarchy | Module/concern grouping for a bounded scope |
| `user story` | Work | Complete flow milestone (acceptance-criteria driven) |
| `feature` | Work | New implementation task |
| `refactor` | Work | Code improvement task |
| `docs` | Work | Documentation task |
| `tests` | Work | Test coverage task |
| `bugfix` | Work | Bug fix |
| `devops` | Work | CI/CD / pipeline / infra ops |
| `database` | Work | Schema / migration task |
| `domain` | Qualifier | Used on master-epics that cover business domain |
| `infra` | Qualifier | Used on master-epics that cover infrastructure |

> **Note:** `domain` and `infra` are qualifier labels — they always appear alongside `master-epic`, never alone.

---

## 3. Title Naming Conventions

### Master Epic

```
`[ㅤ<Pillar Name> - X.0ㅤ]`
```

The Unicode character `ㅤ` (U+3164, Hangul Filler) is used as padding for visual alignment in the board sidebar. Always include it before and after the text.

### Epic

```
`[ㅤ<Module Name>ㅤ]`
```

Same Unicode padding pattern. No version suffix.

### User Story

```
`Category` - Short description of the complete flow
```

Category is the bounded context or module name (e.g., `User`, `Auth`, `Tag`, `Categories`).

### Leaf Issues (feature, refactor, etc.)

Free text imperative title, optionally prefixed with the context:

```
`Context` - Short imperative description
```

or for cross-cutting issues, just a short description without a prefix.

---

## 4. Issue Body Structure

Each issue type has a canonical template in `.github/ISSUE_TEMPLATE/`.

### Epic body sections

1. **General Information** — Category (`epic`), Priority
2. **Description** — Scope, goal, alignment, expected value
3. **Impact** — Positive impacts and considerations
4. **Supplementary Materials** — References, architecture guides

### User Story body sections

1. **General Information** — Category (`user story`), Priority
2. **Description** — Scope, business logic, dependencies
3. **Execution Plan** — Modules involved, key decisions
4. **Impacts** — Benefits and side effects
5. **Supplementary Materials** — Links to specs, mockups, related issues

### Feature / Refactor / Bugfix body sections

1. **General Information** — Category, Priority
2. **Feature / Refactor / Bug Details** — What and why
3. **Acceptance Criteria** — Checkbox list
4. **Execution Plan** — Steps, strategy, rollback
5. **Impacts** — Positive and negative

### Docs body sections

1. **General Information** — Category, Priority
2. **Documentation Details** — Scope, target audience
3. **Acceptance Criteria** — Checklist
4. **Supplementary Materials** — Guidelines, style guide refs

---

## 5. Project Board Views

The board at https://github.com/users/lechitz/projects/1 is a **GitHub Projects v2** board with multiple views. The primary views are described below.

### View: Backlog (Table)

- **Type:** Table
- **Sort:** By issue number (descending) or creation date
- **Grouped by:** None (flat list)
- **Visible fields:** Issue number, title, status, labels, assignees, milestone/iteration
- **Purpose:** Full backlog overview — all issues regardless of type

### View: Board (Kanban)

- **Type:** Board
- **Columns (Status field):**
  - `No Status` — Not yet triaged
  - `Backlog` — Triaged, not yet started
  - `In Progress` — Active development
  - `In Review` — PR open / under review
  - `Done` — Closed / merged
- **Purpose:** Sprint-level tracking of active work

### View: Epics (Table)

- **Type:** Table
- **Filter:** `label:epic OR label:master-epic OR label:Release`
- **Grouped by:** Labels
- **Purpose:** Strategic overview of all epics and master-epics, used for release planning

### View: Roadmap (Roadmap)

- **Type:** Roadmap
- **Grouped by:** Labels (epic / master-epic)
- **Date fields:** Created at → Closed at
- **Purpose:** Timeline visualization of delivery cycles

---

## 6. Epic-to-Module Mapping

Each epic maps to a bounded context or platform concern in the codebase:

| Epic | Code Path | Domain |
|------|-----------|--------|
| User Module (#89) | `internal/user/` | Core user entity, registration, profile |
| Authentication & Authorization (#86) | `internal/auth/`, `internal/token/` | Auth flow, JWT, Keycloak |
| Categories Module (#114) | `internal/category/` | Tag category management |
| TAGs Module (#113) | `internal/tag/` | Habit tags, soft-delete, usage |
| Record Module (#192) | `internal/record/` | Daily record, projections |
| Chat-IA (#195) | `internal/chat/` (+ aion-chat repo) | AI chat, voice, history |
| Observability (#87) | `internal/platform/` (OTel, Prometheus, Loki, Jaeger) | Tracing, metrics, logging |
| DevOps Infrastructure & CI/CD (#95) | `.github/workflows/`, `Makefile`, `infrastructure/` | CI, lint, tests |
| Database (#125) | `infrastructure/db/migrations/` | Schema, migrations |
| Crosscutting Concerns (#128) | `internal/shared/`, `internal/platform/` | Errors, middleware, logging |
| Professional Diary Module (#116) | `internal/professional_diary/` (planned) | Work diary |
| Intentions Module (#117) | `internal/intentions/` (planned) | Daily intentions |
| Moods Module (#118) | `internal/moods/` (planned) | Mood tracking |
| Energy Module (#119) | `internal/energy/` (planned) | Energy score |
| Personal Diary Module (#120) | `internal/personal_diary/` (planned) | Personal diary |
| Water Intake Module (#121) | `internal/water_intake/` (planned) | Hydration logs |
| Day TAG Summary Module (#122) | `internal/day_tag_summary/` (planned) | Habit aggregation per day |

---

## 7. Release Alignment Map

```
Release v1.0.0 (#206)
  ├─ Master Epic: Infrastructure - 1.0 (#124)
  │    ├─ Epic: Database (#125)
  │    ├─ Epic: Observability (#87)
  │    ├─ Epic: DevOps Infrastructure & CI/CD (#95)
  │    └─ Epic: Crosscutting Concerns (#128)
  │
  ├─ Master Epic: Business Rules & Entities - 1.0 (#123)
  │    ├─ Epic: User Module (#89)
  │    ├─ Epic: Authentication & Authorization (#86)
  │    ├─ Epic: Categories Module (#114)
  │    ├─ Epic: TAGs Module (#113)
  │    ├─ Epic: Record Module (#192)
  │    ├─ Epic: Chat-IA (#195)
  │    ├─ Epic: Professional Diary Module (#116)
  │    ├─ Epic: Intentions Module (#117)
  │    ├─ Epic: Moods Module (#118)
  │    ├─ Epic: Energy Module (#119)
  │    ├─ Epic: Personal Diary Module (#120)
  │    ├─ Epic: Water Intake Module (#121)
  │    └─ Epic: Day TAG Summary Module (#122)
  │
  └─ Master Epic: Documentation - 1.0 (#88)
```

---

## 8. Issue Creation Workflow

When creating a new issue, follow this order:

1. **Identify the level.** Is this a whole new module (Epic), a complete flow (User Story), or a specific task (Feature/Refactor/etc.)?
2. **Select the correct template** from `.github/ISSUE_TEMPLATE/`.
3. **Apply the correct label(s).** Never mix work-level labels (e.g., `feature` + `refactor`). A qualifier label (`domain`, `infra`) is only valid alongside `master-epic`.
4. **Use the naming convention** matching the type (see Section 3).
5. **Link the parent.** In the issue body, reference the parent Epic or User Story under "Supplementary Materials". GitHub Projects v2 does not natively enforce parent-child relationships, so explicit `#issue-number` references in the body are the linking mechanism.
6. **Add to the board.** After creation, move the issue into the correct column in the Board view.
7. **Set Priority.** Use the `Priority` field in the board (Low / Medium / High / Critical).

---

## 9. Automation Notes (Codex / Copilot Integration)

When an AI agent creates or updates issues, it must:

- Use Unicode padding `ㅤ` (U+3164) in Epic and Master Epic titles
- Apply exactly one hierarchy-level label (never `epic` + `master-epic` together)
- Follow the body structure from the corresponding template
- Cross-reference parent issues by number in the body
- Assign `lechitz` as the assignee for all issues
- Set the `Release` label only for cross-repo release umbrella issues

When an AI agent closes an issue, it must:
- Set `state_reason: completed` (not `not_planned`) for normal closure
- Reference the closing PR or commit in a comment

---

## Related Files

- `.github/ISSUE_TEMPLATE/` — Issue body templates
- `AGENTS.md` — AI agent rules and architecture guide
- `docs/architecture.md` — System design reference
- `agents/personas/` — Agent role definitions
- `agents/playbooks/` — Step-by-step workflows
