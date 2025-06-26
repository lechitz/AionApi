# ============================================================
#                         TESTING
# ============================================================

.PHONY: test test-cover test-html-report test-ci test-clean

test:
	@echo "Running unit tests..."
	go test ./... -v

test-cover:
	@echo "Running tests with coverage report..."
	go test ./... -coverprofile=coverage_tmp.out -v
	@echo "Filtering out mock files from coverage..."
	cat coverage_tmp.out | grep -v "Mock" > coverage.out
	@rm -f coverage_tmp.out
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out

test-html-report:
	@echo "Running tests and generating JSON output..."
	go test ./... -json > docs/coverage/report.json
	@echo "Generating HTML report..."
	go-test-html-report -f docs/coverage/report.json -o docs/coverage/
	@echo "✅ HTML report generated at: docs/coverage/report.html"

test-ci:
	@echo "Running CI tests with coverage output..."
	go test ./... -coverprofile=coverage.out -v

test-clean:
	@echo "Cleaning up coverage reports..."
	@rm -f coverage.out coverage_tmp.out