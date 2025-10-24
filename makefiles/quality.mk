# ============================================================
#                                CODE QUALITY
# ============================================================

.PHONY: format lint lint-fix verify verify-ci

format:
	@echo "Running goimports..."
	goimports -w .
	@echo "Running golines..."
	golines --max-len=170 --base-formatter=gofumpt -w .

lint: format
	@echo "Running golangci-lint check..."
	golangci-lint run --config=.golangci.yml ./...

lint-fix:
	@echo "Running golangci-lint with --fix..."
	golangci-lint run --fix --config=.golangci.yml ./...

# Local pre-commit verify:
# - DOES NOT modify your working tree
# - Validates Swagger artifacts by generating to a temp dir and diffing
# - Then runs linters and tests
verify: graphql mocks docs.validate lint test test-cover test-ci test-clean
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  Verify passed successfully!"

# CI-style verify (stricter, enforces committed artifacts)
verify-ci: docs.gen docs.check-dirty lint test
	@echo "Running test checks..."
	@$(MAKE) -s test-checks
	@echo "✅  CI verify passed!"
