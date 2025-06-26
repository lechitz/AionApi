# ============================================================
#                DOCKER ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-dev dev-up dev-down dev clean-dev
.PHONY: build-prod prod-up prod-down prod clean-prod
.PHONY: docker-clean-all

build-dev: clean-dev
	@echo "\033[1;36m[BUILD-DEV]\033[0m Building DEV image..."
	docker build -t $(APPLICATION_NAME):dev .

dev-up: dev-down
	@echo "\033[1;36m[DEV-UP]\033[0m Starting DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-down:
	@echo "\033[1;36m[DEV-DOWN]\033[0m Stopping DEV environment..."
	export $$(cat $(ENV_FILE_DEV) | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v

dev: build-dev dev-up

clean-dev:
	@echo "\033[1;33m[CLEAN-DEV]\033[0m Cleaning DEV containers, volumes, images..."
	@docker ps -a --filter "name=dev" -q | xargs -r docker rm -f
	@docker volume ls --filter "name=dev" -q | xargs -r docker volume rm
	@docker images --filter "reference=$(APPLICATION_NAME):dev" -q | xargs -r docker rmi -f

build-prod: clean-prod
	@echo "\033[1;36m[BUILD-PROD]\033[0m Building PROD image..."
	docker build -t $(APPLICATION_NAME):prod .

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