# ============================================================
#                         TESTING
# ============================================================

.PHONY: test test-cover test-html-report test-ci test-clean test-checks

# Execute unit tests
test:
	@echo "Running unit tests with race detector..."
	go test ./... -v -race

# Run tests with coverage, filter mocks, and generate HTML coverage report
test-cover:
	@echo "Running tests with coverage report..."
	go test ./... -race -coverprofile=$(COVERAGE_DIR)/coverage_tmp.out -v
	@echo "Filtering out mock files from coverage..."
	grep -v "Mock" $(COVERAGE_DIR)/coverage_tmp.out > $(COVERAGE_DIR)/coverage.out
	@rm -f $(COVERAGE_DIR)/coverage_tmp.out
	@echo "Generating HTML coverage report..."
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

# Generate JUnit XML test report via gotestsum
test-html-report:
	@echo "Running tests and generating JUnit XML report..."
	gotestsum --junitfile $(COVERAGE_DIR)/junit-report.xml -- -race ./...
	@echo "✅ JUnit report generated at $(COVERAGE_DIR)/junit-report.xml"

# CI target: tests with coverage but no HTML UI
test-ci:
	@echo "Running CI tests with coverage output..."
	go test ./... -race -coverprofile=$(COVERAGE_DIR)/coverage.out -v

# Check tests for common anti-patterns and fail early with actionable message
# Currently checks for uses of context.Background() inside _test.go files.
# Prefer t.Context() (or suite.T().Context()) so test cancellations/timeouts propagate.
test-checks:
	@echo "Checking tests for discouraged patterns..."
	@matches=$$(grep -R --line-number "context.Background()" --include="*_test.go" . \
		--exclude-dir=.git --exclude-dir=vendor --exclude-dir=tests || true); \
	if [ -n "$$matches" ]; then \
		echo "Found context.Background() usages in tests. Prefer using t.Context() or the suite's context. Matches:"; \
		echo "$$matches"; \
		exit 1; \
	else \
		echo "No discouraged test patterns found."; \
	fi

# Cleanup coverage artifacts and test reports
test-clean:
	@echo "Cleaning up coverage reports and test artifacts..."
	@rm -f \
		$(COVERAGE_DIR)/coverage.out \
		$(COVERAGE_DIR)/coverage_tmp.out \
		$(COVERAGE_DIR)/coverage.html \
		$(COVERAGE_DIR)/junit-report.xml
	@echo "✅ Cleanup complete!"
