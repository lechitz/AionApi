# Changelog - 2026-02-01

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
query_path = Path("../AionApi/graphql/queries/categories/list.graphql")
```

**aionapi-dashboard (TypeScript):**
```typescript
const query = readFileSync('../AionApi/graphql/queries/categories/list.graphql');
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
