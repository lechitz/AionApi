# ============================================================
#                   GLOBAL VARIABLES & CONFIG
# ============================================================

APPLICATION_NAME := aion-api
HTTP_PORT        := 5001
COMPOSE_FILE_DEV := infra/docker/dev/docker-compose-dev.yaml
ENV_FILE_DEV     := infra/docker/dev/.env.dev
COMPOSE_FILE_PROD := infra/docker/prod/docker-compose-prod.yaml
ENV_FILE_PROD    := infra/docker/prod/.env.prod

# --- MIGRATION CONFIG ---
MIGRATION_PATH := infra/db/migrations
MIGRATION_DB   ?= $(DB_URL)
MIGRATE_BIN := $(shell command -v migrate 2> /dev/null)

# ============================================================
#                       HELP SECTION
# ============================================================

.PHONY: help
help:
	@echo ""
	@echo "\033[1;33m\033[1mAionApi - Developer CLI Help\033[0m"
	@echo ""
	@echo "\033[1;33m- Docker Compose Commands:\033[0m"
	@echo "  \033[1;36mbuild-dev\033[0m        → Build dev image"
	@echo "  \033[1;36mdev-up\033[0m           → Start dev environment (resets DB)"
	@echo "  \033[1;36mdev-down\033[0m         → Stop dev and remove volumes"
	@echo "  \033[1;36mclean-dev\033[0m        → Clean dev containers, volumes, images"
	@echo ""
	@echo "  \033[1;36mbuild-prod\033[0m       → Build prod image"
	@echo "  \033[1;36mprod-up\033[0m          → Start prod environment"
	@echo "  \033[1;36mprod-down\033[0m        → Stop prod and remove volumes"
	@echo "  \033[1;36mclean-prod\033[0m       → Clean prod containers, volumes, images"
	@echo ""
	@echo "  \033[1;36mdocker-clean-all\033[0m → Remove ALL containers, volumes, images"
	@echo ""
	@echo "\033[1;33m- Testing:\033[0m"
	@echo "  \033[1;36mtest\033[0m             → Run unit tests"
	@echo "  \033[1;36mtest-cover\033[0m       → Coverage (exclui mocks)"
	@echo "  \033[1;36mtest-html-report\033[0m → HTML report (requires go-test-html-report)"
	@echo ""
	@echo "\033[1;33m- Migrations (via migrate CLI):\033[0m"
	@echo "  \033[1;36mmigrate-up\033[0m       → Executa todas as migrations (up)"
	@echo "  \033[1;36mmigrate-down\033[0m     → Rollback última migration"
	@echo "  \033[1;36mmigrate-force VERSION=X\033[0m → Força DB para versão X"
	@echo "  \033[1;36mmigrate-new\033[0m      → Cria nova migration (prompt)"
	@echo ""
	@echo "\033[1;33m- Codegen:\033[0m"
	@echo "  \033[1;36mgraphql\033[0m          → Gera arquivos GraphQL via gqlgen"
	@echo "  \033[1;36mmocks\033[0m            → Gera todos os mocks GoMock"
	@echo ""

# ============================================================
#                DOCKER DEV ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-dev dev-up dev-down dev clean-dev

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

# ============================================================
#                DOCKER PROD ENVIRONMENT TARGETS
# ============================================================

.PHONY: build-prod prod-up prod-down prod clean-prod

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

# ============================================================
#                GENERAL DOCKER CLEANUP TARGET
# ============================================================

.PHONY: docker-clean-all
docker-clean-all:
	@echo "\033[1;33m[CLEAN-ALL]\033[0m Removing ALL containers, volumes, images..."
	@docker ps -a -q | xargs -r docker rm -f
	@docker volume ls -q | xargs -r docker volume rm
	@docker images -a -q | xargs -r docker rmi -f

# ============================================================
#                         TESTING
# ============================================================

.PHONY: test test-cover test-html-report test-ci test-clean

test:
	@echo "📦 Running unit tests..."
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
	@echo "🧪 Running tests and generating JSON output..."
	go test ./... -json > docs/coverage/report.json
	@echo "📄 Generating HTML report..."
	go-test-html-report -f docs/coverage/report.json -o docs/coverage/
	@echo "✅ HTML report generated at: docs/coverage/report.html"

test-ci:
	@echo "Running CI tests with coverage output..."
	go test ./... -coverprofile=coverage.out -v

test-clean:
	@echo "Cleaning up coverage reports..."
	@rm -f coverage.out coverage_tmp.out

# ============================================================
#                    GRAPHQL CODEGEN
# ============================================================

.PHONY: graphql

graphql:
	cd internal/adapters/primary/graph && go run github.com/99designs/gqlgen generate

# ============================================================
#                   MOCKS GENERATION (GOMOCK)
# ============================================================

.PHONY: mocks
mocks:
	@echo "Generating mocks for output ports and usecases..."
	@mkdir -p tests/mocks/token tests/mocks/user tests/mocks/security tests/mocks/logger tests/mocks/category

	@echo "→ TokenStore"
	mockgen -source=internal/core/ports/output/cache/token.go \
		-destination=tests/mocks/token/mock_token_store.go \
		-package=tokenmocks \
		-mock_names=Store=MockTokenStore

	@echo "→ TokenUsecase"
	mockgen -source=internal/core/usecase/token/token_usecase.go \
		-destination=tests/mocks/token/mock_token_usecase.go \
		-package=tokenmocks \
		-mock_names=Usecase=MockTokenUsecase

	@echo "→ UserStore"
	mockgen -source=internal/core/ports/output/db/user.go \
		-destination=tests/mocks/user/mock_user_store.go \
		-package=usermocks \
		-mock_names=UserStore=MockUserStore

	@echo "→ CategoryStore"
	mockgen -source=internal/core/ports/output/db/category.go \
		-destination=tests/mocks/category/mock_category_store.go \
		-package=categorymocks \
		-mock_names=CategoryStore=MockCategoryStore

	@echo "→ SecurityStore"
	mockgen -source=internal/core/ports/output/security/hasher.go \
		-destination=tests/mocks/security/mock_security_store.go \
		-package=securitymocks \
		-mock_names=Store=MockSecurityStore

	@echo "→ Logger"
	mockgen -source=internal/core/ports/output/logger/logger.go \
		-destination=tests/mocks/logger/mock_logger.go \
		-package=loggermocks \
		-mock_names=Logger=MockLogger

	@echo "✅ All mocks generated successfully."

# ============================================================
#                        MIGRATIONS
# ============================================================

.PHONY: migrate-up migrate-down migrate-force migrate-new

migrate-up:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export DB_URL=...' or adjust the Makefile."; \
		exit 1; \
	fi
	@echo "🔄 Running all migrations (up)..."
	@$(MIGRATE_BIN) -path $(MIGRATION_PATH) -database "$(MIGRATION_DB)" up

migrate-down:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export DB_URL=...' or adjust the Makefile."; \
		exit 1; \
	fi
	@echo "↩️  Rolling back the last migration (down)..."
	@$(MIGRATE_BIN) -path $(MIGRATION_PATH) -database "$(MIGRATION_DB)" down 1

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Please provide VERSION=X to force (e.g., make migrate-force VERSION=2)"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export DB_URL=...' or adjust the Makefile."; \
		exit 1; \
	fi
	@echo "🚨 Forcing DB schema version to $(VERSION)..."
	@$(MIGRATE_BIN) -path $(MIGRATION_PATH) -database "$(MIGRATION_DB)" force $(VERSION)

migrate-new:
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then echo "Migration name is required."; exit 1; fi; \
	migrate create -ext sql -dir $(MIGRATION_PATH) $$name
