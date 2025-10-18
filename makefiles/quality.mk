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
	@echo "✅  Verify passed successfully!"

# CI-style verify (stricter, enforces committed artifacts)
verify-ci: docs.gen docs.check-dirty lint test
	@echo "✅  CI verify passed!"

sensitive-strings:
	@echo "Checking for sensitive literal strings (this may report false positives)..."
	@grep -R --line-number -E "(refresh_token|\bAuthorization\b|\btoken\b)" . \
		--exclude-dir=.git --exclude-dir=tests --exclude-dir=vendor --exclude-dir=.venv --exclude-dir=.venv-docs --exclude-dir=node_modules \
		--exclude=**/*_test.go --exclude=**/*.md --exclude=**/migrations/** || true
	@echo "If the previous command printed matches, review them. To fail CI on matches, run:\n\tmake sensitive-strings-ci"

sensitive-strings-ci:
	@echo "Running sensitive string check (CI - will fail on matches)"
	@matches=$$(grep -R --line-number -E "(refresh_token|\bAuthorization\b|\btoken\b)" . --exclude-dir=.git --exclude-dir=tests --exclude-dir=vendor --exclude-dir=.venv --exclude-dir=.venv-docs --exclude-dir=node_modules --exclude=**/*_test.go --exclude=**/*.md --exclude=**/migrations/** || true); \
	if [ -n "$$matches" ]; then echo "Sensitive strings found:"; echo "$$matches"; exit 1; else echo "No sensitive strings detected."; fi
