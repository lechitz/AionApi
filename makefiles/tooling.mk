# ============================================================
#                		   TOOLING
# ============================================================

tools-install:
	@echo "Installing development tools..."
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/99designs/gqlgen@latest
	@echo "✅  Tools installed successfully."