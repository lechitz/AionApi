# hack/ (Developer Utilities)

Development-only tools and scripts. **NOT included in production builds.**

Following Kubernetes convention: `hack/` contains experimental/dev utilities.

## Structure

- `tools/` - Go CLIs for dev/test (seed-caller, seed-helper)
- `dev/` - Bash scripts for troubleshooting (Ollama setup, DB roles, tests)

## Usage

Via Makefile (recommended):
```bash
make seed-api-caller    # runs hack/tools/seed-caller
make hash-gen PASS=...  # inline bcrypt (no binary)
bash hack/dev/test-chat.sh
```

Direct:
```bash
go run ./hack/tools/seed-caller
bash hack/dev/force-insert-roles.sh
```

## Not for Production

This folder is ignored in:
- `.dockerignore` (not copied to images)
- CI release builds (only cmd/api is built)

## Why "hack/"?

This naming follows the Kubernetes convention where `hack/` contains:
- Build automation scripts
- Development tools
- Test utilities
- Code generation helpers

It's called "hack" because these are utilities that "hack around" during development,
not production code.
