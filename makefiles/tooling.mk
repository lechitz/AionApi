# ============================================================Ajuste
#                		   TOOLING
# ============================================================

tools-install:
	@echo "Installing development tools..."
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/99designs/gqlgen@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@v1.16.6
	@echo "✅  Tools installed successfully."

# Check for required development tools and fail if any are missing.
.PHONY: tools.check
tools.check:
	@echo "Checking required development tools..."
	@missing=0; \
	for cmd in goimports golines gofumpt golangci-lint mockgen migrate swag gqlgen gotestsum fieldalignment; do \
		if ! command -v $$cmd >/dev/null 2>&1; then \
			echo " - $$cmd: MISSING"; missing=1; \
		else \
			echo " - $$cmd: present"; \
		fi; \
	done; \
	if [ $$missing -eq 1 ]; then \
		echo ""; \
		echo "One or more development tools are missing. Run 'make tools-install' to install them or install manually."; \
		exit 1; \
	fi; \
	echo "All required development tools are present.";

.PHONY: graph-projection-export
GRAPH_PROJECTION_USER_ID ?= 999
GRAPH_PROJECTION_WINDOW ?= WINDOW_30D
GRAPH_PROJECTION_TIMEZONE ?= America/Sao_Paulo

graph-projection-export:
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		DB_HOST=localhost \
		go run ./hack/tools/graph-projection-export \
			--user-id "$(GRAPH_PROJECTION_USER_ID)" \
			--window "$(GRAPH_PROJECTION_WINDOW)" \
			--timezone "$(GRAPH_PROJECTION_TIMEZONE)" \
			$(if $(DATE),--date "$(DATE)",) \
			$(if $(CATEGORY_ID),--category-id "$(CATEGORY_ID)",) \
			$(if $(TAG_IDS),--tag-ids "$(TAG_IDS)",) \
			$(if $(OUTPUT),--output "$(OUTPUT)",)
