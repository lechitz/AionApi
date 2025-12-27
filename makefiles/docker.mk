# ============================================================
#                DOCKER ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-dev dev-up dev-down dev dev-clean clean-dev
.PHONY: build-prod prod-up prod-down prod clean-prod
.PHONY: docker-clean-all

APPLICATION_NAME := aion-api

build-dev: clean-dev
	@echo "[BUILD-DEV] Building DEV image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):dev .

dev-up: dev-down
	@echo "[DEV-UP] Starting DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-down:
	@echo "[DEV-DOWN] Stopping DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v

dev: build-dev
	@echo "[DEV] Starting DEV environment (detached)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v || true
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d
	@echo "✓ Services started in background"
	@echo "→ Use 'make dev-logs' or 'make dev-attach' to view logs"
	@echo "→ Use 'make dev-down' to stop all services"
	@echo ""
	@echo "ℹ️  Quick commands:"
	@echo "   make dev-attach   → Attach to aion-api logs (without rebuild)"
	@echo "   make dev-logs     → Show all services logs"
	@echo "   make migrate-up   → Run database migrations"
	@echo "   make seed-caller n=1 → Seed data via API"

dev-attach:
	@echo "[DEV-ATTACH] Attaching to aion-api logs..."
	@echo "(Press Ctrl+C to detach - containers will keep running)"
	@echo ""
	@trap 'echo ""; echo "✓ Detached. Containers still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-api

dev-logs:
	@echo "[DEV-LOGS] Showing all services logs..."
	@echo "(Press Ctrl+C to stop - containers will keep running)"
	@echo ""
	@trap 'echo ""; echo "✓ Stopped viewing logs. Containers still running."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f

dev-clean: clean-dev dev

clean-dev:
	@echo "[CLEAN-DEV] Cleaning DEV containers, volumes, images..."
	@echo "→ Stopping and removing compose stack..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v --remove-orphans || true
	@echo "→ Removing dev image..."
	@docker images --filter "reference=$(APPLICATION_NAME):dev" -q | xargs -r docker rmi -f || true
	@echo "✓ Cleanup complete"

build-prod: clean-prod
	@echo "[BUILD-PROD] Building PROD image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="-s -w" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):prod .

prod-up: prod-down
	@echo "[PROD-UP] Starting PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) up

prod-down:
	@echo "[PROD-DOWN] Stopping PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) down -v

prod: build-prod prod-up

clean-prod:
	@echo "[CLEAN-PROD] Cleaning PROD containers, volumes, images..."
	@docker ps -a --filter "name=prod" -q | xargs -r docker rm -f
	@docker volume ls --filter "name=prod" -q | xargs -r docker volume rm
	@docker images --filter "reference=*prod*" -q | xargs -r docker rmi -f

docker-clean-all:
	@echo "[CLEAN-ALL] Removing ALL containers, volumes, images..."
	@docker ps -a -q | xargs -r docker rm -f
	@docker volume ls -q | xargs -r docker volume rm
	@docker images -a -q | xargs -r docker rmi -f
