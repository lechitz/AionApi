# ============================================================
#                                CODE QUALITY
# ============================================================

.PHONY: format lint lint-fix verify

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

verify: mocks graphql lint test test-cover test-ci test-clean
	@echo "✅  Verify passed successfully!"