# ============================================================
#                                CODE QUALITY
# ============================================================

GO_CACHE := $(CURDIR)/.cache/go-build
GOLANGCI_CACHE := $(CURDIR)/.cache/golangci-lint

.PHONY: format lint lint-fix fieldalignment fieldalignment-report go-check verify verify-ci

# Run goimports and golines to format code
format:
	@echo "Running goimports..."
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "warning: 'goimports' not found, skipping (install with: go install golang.org/x/tools/cmd/goimports@latest)"; \
	fi
	@echo "Running golines..."
	@if command -v golines >/dev/null 2>&1; then \
		golines --max-len=170 --base-formatter=gofumpt -w .; \
	else \
		echo "warning: 'golines' not found, skipping (install with: go install github.com/segmentio/golines@latest)"; \
	fi
	@echo "Finished formatting checks."

# Run golangci-lint checks
lint: format
	@echo "Running golangci-lint check..."
	@mkdir -p $(GO_CACHE) $(GOLANGCI_CACHE)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		GOCACHE=$(GO_CACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_CACHE) golangci-lint run --config=.golangci.yml ./... || { \
			echo "golangci-lint failed. If this is due to unknown linters, try installing a compatible version:"; \
			echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0"; \
			exit 1; \
		}; \
	else \
		echo "warning: 'golangci-lint' not found, skipping lint (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0)"; \
	fi

# Check struct field alignment for memory efficiency (strict - fails on issues)
fieldalignment:
	@if command -v fieldalignment >/dev/null 2>&1; then \
		echo "Running fieldalignment on critical paths..."; \
		mkdir -p $(GO_CACHE); \
		GOCACHE=$(GO_CACHE) fieldalignment ./internal/platform/config ./internal/chat/core/domain || exit 1; \
	else \
		echo "warning: 'fieldalignment' not found, skipping (install with: make tools-install)"; \
	fi

# Check struct field alignment across entire repo (advisory - reports but doesn't fail)
fieldalignment-report:
	@if command -v fieldalignment >/dev/null 2>&1; then \
		echo "Running fieldalignment across entire repository..."; \
		mkdir -p $(GO_CACHE); \
		GOCACHE=$(GO_CACHE) fieldalignment ./... || echo "⚠️  fieldalignment found issues (advisory only)"; \
	else \
		echo "warning: 'fieldalignment' not found, skipping (install with: make tools-install)"; \
	fi

# Fast sanity check to catch broken imports or modules before heavier steps
go-check:
	@echo "Checking if all Go packages compile (imports/modules)..."
	@mkdir -p $(GO_CACHE)
	@GOCACHE=$(GO_CACHE) go test ./... -run=^$ -count=0 >/dev/null

# Auto-fix lint issues where possible
lint-fix:
	@echo "Running golangci-lint with --fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --config=.golangci.yml ./... || true; \
	else \
		echo "warning: 'golangci-lint' not found, skipping lint-fix (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0)"; \
	fi

# General verify (checks code quality, but does not enforce committed artifacts)
verify: go-check lint fieldalignment graphql mocks docs.validate test test-cover-detail test-ci test-clean
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  Verify passed successfully!"

# CI-style verify (stricter, enforces committed artifacts)
verify-ci: tools.check docs.gen docs.check-dirty lint fieldalignment test
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  CI verify passed!"
