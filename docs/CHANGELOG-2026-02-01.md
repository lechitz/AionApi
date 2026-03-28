# Changelog - 2026-02-01

## Audit Context + Chat Integration (2026-02-26)

### Added

- New `audit` bounded context with:
  - Domain/entity and filter model for action events.
  - Input/output ports and use case service (`WriteEvent`, `ListEvents`).
  - DB secondary adapter (model, mapper, repository with `save`/`list`).
  - Primary HTTP adapter for `GET /aion/api/v1/audit/events`.
- Structured security log for cross-user audit reads:
  - `actor_user_id`, `target_user_id`, `trace_id`, `draft_id`.
- New mocks and tests for audit use case, repository, and handler.

### Changed

- Chat flow now emits audit action events through `audit.Service` on relevant message-processing outcomes.
- Dependency wiring updated to compose and inject `AuditService` into chat use case.
- HTTP composer now registers audit endpoints when dependencies are available.
- Documentation updated:
  - `internal/audit/README.md`
  - `internal/chat/README.md`
  - `internal/platform/server/http/README.md`

### Tests and Coverage

- `go test -cover ./internal/audit/...`
  - handler: 84.3%
  - repository: 97.7%
  - usecase: 91.2%
- `go test -cover ./internal/chat/core/usecase`
  - usecase: 96.8%

## Auth Session Compatibility (2026-02-26)

### Changed

- `GET /aion/api/v1/auth/session` now accepts bearer token in `Authorization` header, while keeping cookie fallback for browser compatibility.
- Added tests for bearer-token session validation path.

### Tests and Coverage

- `go test -cover ./internal/auth/adapter/primary/http/handler`
  - handler: 93.5%

## GraphQL Documentation System

### Added

**New Make Targets (makefiles/graphql.mk):**
- `make graphql.schema` - Export merged SDL schema from all .graphqls files
- `make graphql.queries` - Generate shared query templates for clients
- `make graphql.docs` - Generate complete documentation
- `make graphql.setup` - Complete setup (schema + queries + docs)
- `make graphql.clean` - Clean generated artifacts

**Generated Artifacts:**
- `docs/graphql/schema.graphql` - Complete SDL schema (231 lines, auto-generated)
- `docs/graphql/README.md` - Documentation with playground links and usage examples
- `graphql/queries/` - Shared queries directory:
  - `categories/list.graphql` - List all categories query
  - `categories/create.graphql` - Create category mutation
  - `tags/list.graphql` - List all tags query
  - `tags/create.graphql` - Create tag mutation
  - `records/list.graphql` - List records query
  - `records/create.graphql` - Create record mutation
  - `README.md` - Queries documentation

### Benefits

- **Single Source of Truth:** Schema exported automatically from .graphqls sources
- **Shared Queries:** Same queries used by dashboard, aion-chat, and future CLI
- **Professional Structure:** Similar to Swagger docs at `docs/swagger/`
- **Easy Regeneration:** `make graphql.schema` rebuilds everything
- **Prepared for Scale:** Ready for GraphQL Codegen integration when needed

### Usage

```bash
# Initial setup
make graphql.setup

# After modifying .graphqls files
make graphql.schema

# Clean and rebuild
make graphql.clean && make graphql.setup
```

### Integration

**aion-chat (Python):**
```python
query_path = Path("../aion-api/graphql/queries/categories/list.graphql")
```

**aion-web (TypeScript):**
```typescript
const query = readFileSync('../aion-api/graphql/queries/categories/list.graphql');
```

### Future Roadmap

- [ ] Setup GraphQL Codegen for TypeScript (when 50+ queries)
- [ ] Setup ariadne-codegen for Python
- [ ] Add `make graphql.introspect` with running server check
- [ ] CI/CD validation (queries against schema)

---

**Commit:** d1151ad  
**Branch:** feature/observability-tracing-fixes  
**Files Changed:** 7 files, 278+ lines
