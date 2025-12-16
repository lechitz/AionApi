# ============================================================
#                DOCKER ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-dev dev-up dev-down dev dev-clean clean-dev
.PHONY: build-prod prod-up prod-down prod clean-prod
.PHONY: docker-clean-all

APPLICATION_NAME := aion-api

build-dev: clean-dev
	@echo "\033[1;36m[BUILD-DEV]\033[0m Building DEV image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):dev .

dev-up: dev-down
	@echo "\033[1;36m[DEV-UP]\033[0m Starting DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-down:
	@echo "\033[1;36m[DEV-DOWN]\033[0m Stopping DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v

dev: dev-clean

dev-logs: build-dev dev-up
	@echo "\033[1;36m[DEV-LOGS]\033[0m Starting DEV environment (foreground, all logs)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-clean: build-dev
	@echo "\033[1;36m[DEV-CLEAN]\033[0m Starting DEV (detached, following aion-api only)..."
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v || true
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres || true
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up -d
	@echo "\033[1;32m✓\033[0m Services started in background. Following aion-api logs..."
	@echo "\033[1;33m→\033[0m Use 'docker compose -f $(COMPOSE_FILE_DEV) logs -f' to see all logs"
	@echo "\033[1;33m→\033[0m Use 'Ctrl+C' to stop following logs (containers will keep running)"
	@echo "\033[1;33m→\033[0m Use 'make dev-down' to stop all services\n"
	@trap 'echo ""; echo "Logs interrupted. Containers still running. Use \"make dev-down\" to stop."; exit 0' INT; \
		export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) logs -f aion-api

clean-dev:
	@echo "\033[1;33m[CLEAN-DEV]\033[0m Cleaning DEV containers, volumes, images..."
	@echo "\033[1;90m→ Stopping and removing compose stack...\033[0m"
	@export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v --remove-orphans || true
	@echo "\033[1;90m→ Removing dev image...\033[0m"
	@docker images --filter "reference=$(APPLICATION_NAME):dev" -q | xargs -r docker rmi -f || true
	@echo "\033[1;32m✓ Cleanup complete\033[0m"

build-prod: clean-prod
	@echo "\033[1;36m[BUILD-PROD]\033[0m Building PROD image..."
	DOCKER_BUILDKIT=1 docker build --progress=plain --build-arg BUILD_LDFLAGS="-s -w" -f infrastructure/docker/Dockerfile -t $(APPLICATION_NAME):prod .

prod-up: prod-down
	@echo "\033[1;36m[PROD-UP]\033[0m Starting PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) up

prod-down:
	@echo "\033[1;36m[PROD-DOWN]\033[0m Stopping PROD environment..."
	export $$(cat $(ENV_FILE_PROD) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_PROD) down -v

prod: build-prod prod-up

clean-prod:
	@echo "\033[1;33m[CLEAN-PROD]\033[0m Cleaning PROD containers, volumes, images..."
	@docker ps -a --filter "name=prod" -q | xargs -r docker rm -f
	@docker volume ls --filter "name=prod" -q | xargs -r docker volume rm
	@docker images --filter "reference=*prod*" -q | xargs -r docker rmi -f

docker-clean-all:
	@echo "\033[1;33m[CLEAN-ALL]\033[0m Removing ALL containers, volumes, images..."
	@docker ps -a -q | xargs -r docker rm -f
	@docker volume ls -q | xargs -r docker volume rm
	@docker images -a -q | xargs -r docker rmi -f
