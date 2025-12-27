# ============================================================
#                   LOCAL DEVELOPMENT (HOT-RELOAD)
# ============================================================
# Run the API locally with hot-reload using Air
# Requires: air, docker-compose (for dependencies only)

.PHONY: dev-local dev-local-deps dev-local-stop air-install

# Check if air is installed
AIR_INSTALLED := $(shell command -v air 2> /dev/null)

air-install:
ifndef AIR_INSTALLED
	@echo "[AIR] Installing air for hot-reload..."
	@go install github.com/air-verse/air@latest
	@echo "✓ Air installed successfully"
else
	@echo "✓ Air is already installed"
endif

dev-local-deps:
	@echo "[DEV-LOCAL-DEPS] Starting infrastructure dependencies..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) up -d postgres redis otel-collector jaeger prometheus grafana
	@echo "✓ Dependencies started (postgres, redis, otel, jaeger, prometheus, grafana)"
	@echo "→ Waiting for services to be ready..."
	@sleep 5
	@echo "✓ Services should be ready now"

dev-local-stop:
	@echo "[DEV-LOCAL-STOP] Stopping local dev dependencies..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) stop postgres redis otel-collector jaeger prometheus grafana
	@echo "✓ Dependencies stopped (containers preserved for restart)"

dev-local-down:
	@echo "[DEV-LOCAL-DOWN] Removing local dev dependencies..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		docker compose -f $(COMPOSE_FILE_DEV) down
	@echo "✓ Dependencies removed"

dev-local: air-install dev-local-deps
	@echo "[DEV-LOCAL] Starting API with hot-reload..."
	@echo "(Air will watch for file changes and rebuild automatically)"
	@echo "(Press Ctrl+C to stop - dependencies will keep running)"
	@echo ""
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		DB_HOST=localhost \
		CACHE_ADDR=localhost:6379 \
		OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 \
		air

dev-local-full: dev-local-deps
	@echo "[DEV-LOCAL-FULL] Running without Air (normal go run)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && \
		DB_HOST=localhost \
		CACHE_ADDR=localhost:6379 \
		OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 \
		go run ./cmd/aion-api

