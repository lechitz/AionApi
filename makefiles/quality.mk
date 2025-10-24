# ============================================================
#                                CODE QUALITY
# ============================================================

.PHONY: format lint lint-fix verify verify-ci

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

lint: format
	@echo "Running golangci-lint check..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=.golangci.yml ./... || { \
			echo "golangci-lint failed. If this is due to unknown linters, try installing a compatible version:"; \
			echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0"; \
			exit 1; \
		}; \
	else \
		echo "warning: 'golangci-lint' not found, skipping lint (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0)"; \
	fi

lint-fix:
	@echo "Running golangci-lint with --fix..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix --config=.golangci.yml ./... || true; \
	else \
		echo "warning: 'golangci-lint' not found, skipping lint-fix (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0)"; \
	fi

# Local pre-commit verify:
# - DOES NOT modify your working tree
# - Validates Swagger artifacts by generating to a temp dir and diffing
# - Then runs linters and tests
verify: graphql mocks docs.validate lint test test-cover test-ci test-clean
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  Verify passed successfully!"

# CI-style verify (stricter, enforces committed artifacts)
verify-ci: tools.check docs.gen docs.check-dirty lint test
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  CI verify passed!"
