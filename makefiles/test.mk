# ============================================================
#                         TESTING
# ============================================================

.PHONY: test test-cover test-html-report test-ci test-clean

test:
	@echo "Running unit tests with race detector..."
	go test ./... -v -race

# Run tests with coverage, filter mocks, and generate HTML coverage report
test-cover:
	@echo "Running tests with coverage report..."
	go test ./... -race -coverprofile=coverage_tmp.out -v
	@echo "Filtering out mock files from coverage..."
	grep -v "Mock" coverage_tmp.out > coverage.out
	@rm -f coverage_tmp.out
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html

# Generate JUnit XML test report via gotestsum
test-html-report:
	@echo "Running tests and generating JUnit XML report..."
	gotestsum --junitfile tests/coverage/junit-report.xml -- -race ./...
	@echo "✅ JUnit report generated at tests/coverage/junit-report.xml"

# CI target: tests with coverage but no HTML UI
test-ci:
	@echo "Running CI tests with coverage output..."
	go test ./... -race -coverprofile=coverage.out -v

# Cleanup coverage artifacts and test reports
test-clean:
	@echo "Cleaning up coverage reports and test artifacts..."
	rm -f coverage.out coverage_tmp.out coverage.html tests/coverage/junit-report.xml
