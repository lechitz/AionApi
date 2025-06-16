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
#                HELP & TOOLING SECTION
# ============================================================

.PHONY: help tools-install

help:
	@echo ""
	@echo "\033[48;5;235;33m====================================================================================\033[0m"
	@echo "\033[48;5;235;33m                              AION API - CLI COMMANDS                               \033[0m"
	@echo "\033[48;5;235;33m====================================================================================\033[0m"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ TOOLING ┃\033[0m"
	@echo ""
	@echo "    \033[1;37mtools-install        \033[0m    →  Install all development tools (goimports, golines, gofumpt, golangci-lint)"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ DOCKER ENVIRONMENT COMMANDS ┃\033[0m"
	@echo ""
	@echo "  \033[1;39m- [DEV]\033[0m"
	@echo "    \033[1;37mbuild-dev          \033[0m      →  Build the development Docker image"
	@echo "    \033[1;37mdev-up             \033[0m      →  Start the development environment (resets DB)"
	@echo "    \033[1;37mdev-down           \033[0m      →  Stop and remove dev environment containers/volumes"
	@echo "    \033[1;37mclean-dev          \033[0m      →  Clean all dev containers, volumes, and images"
	@echo ""
	@echo "  \033[1;39m- [PROD]\033[0m"
	@echo "    \033[1;37mbuild-prod         \033[0m      →  Build the production Docker image"
	@echo "    \033[1;37mprod-up            \033[0m      →  Start the production environment"
	@echo "    \033[1;37mprod-down          \033[0m      →  Stop and remove prod environment containers/volumes"
	@echo "    \033[1;37mclean-prod         \033[0m      →  Clean all prod containers, volumes, and images"
	@echo ""
	@echo "  \033[1;39m- [GENERAL]\033[0m"
	@echo "    \033[1;37mdocker-clean-all    \033[0m     →  Remove ALL Docker containers, volumes, and images"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ CODE GENERATION ┃\033[0m"
	@echo ""
	@echo "    \033[1;37mgraphql             \033[0m     →  Generate GraphQL files with gqlgen"
	@echo "    \033[1;37mmocks               \033[0m     →  Generate all GoMock mocks"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ CODE QUALITY ┃\033[0m"
	@echo ""
	@echo "    \033[1;37mformat               \033[0m    →  Format Go code using goimports/golines/gofumpt"
	@echo "    \033[1;37mlint                 \033[0m    →  Run golangci-lint (static code analysis)"
	@echo "    \033[1;37mlint-fix             \033[0m    →  Run golangci-lint with --fix (auto-fix where possible)"
	@echo "    \033[1;37mverify               \033[0m    →  Run full pre-commit pipeline (format, mocks, lint, tests, coverage, codegen)"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ MIGRATIONS ┃\033[0m"
	@echo ""
	@echo "    \033[1;37mmigrate-up           \033[0m    →  Run all migrations (up)"
	@echo "    \033[1;37mmigrate-down         \033[0m    →  Rollback the last migration"
	@echo "    \033[1;37mmigrate-force VERSION=X\033[0m  →  Force DB to specific version"
	@echo "    \033[1;37mmigrate-new          \033[0m    →  Create new migration (with prompt)"
	@echo ""
	@echo "\033[48;5;235;33m🔶 ┃ TESTING ┃\033[0m"
	@echo ""
	@echo "    \033[1;37mtest                 \033[0m    →  Run unit tests"
	@echo "    \033[1;37mtest-cover           \033[0m    →  Run tests with coverage report (excludes mocks)"
	@echo "    \033[1;37mtest-html-report     \033[0m    →  Generate HTML test report (requires go-test-html-report)"
	@echo ""
	@echo "\033[48;5;235;33m====================================================================================\033[0m"
	@echo ""

# ============================================================
#                		   TOOLING
# ============================================================

tools-install:
	@echo "📦 Installing development tools..."
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Tools installed successfully."

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

# ============================================================
#              		    CODE GENERATION
# ============================================================

.PHONY: graphql mocks

graphql:
	cd internal/adapters/primary/graph && go run github.com/99designs/gqlgen generate

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
#                		 CODE QUALITY
# ============================================================

.PHONY: format lint lint-fix verify

format:
	@echo "🎨 Running goimports..."
	goimports -w .
	@echo "🖊️  Running golines..."
	golines --max-len=170 --base-formatter=gofumpt -w .

lint: format
	@echo "🔎 Running golangci-lint check..."
	golangci-lint run --config=.golangci.yml ./...

lint-fix:
	@echo "🔧 Running golangci-lint with --fix..."
	golangci-lint run --fix --config=.golangci.yml ./...

verify: mocks graphql lint test test-cover test-ci test-clean
	@echo "✅ Verify passed successfully!"

# ============================================================
#                         MIGRATIONS
# ============================================================

.PHONY: migrate-up migrate-down migrate-force migrate-new

migrate-up:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "🔄 Running all migrations (up)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" up

migrate-down:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "↩️  Rolling back the last migration (1 step)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" down 1

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION not provided. Use 'make migrate-force VERSION=X'"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@if [ -z "$(MIGRATION_DB)" ]; then \
		echo "❌ MIGRATION_DB is not set. Use 'export MIGRATION_DB=...';"; \
		exit 1; \
	fi
	@echo "🚨 Forcing DB schema version to $(VERSION)..."
	@$(MIGRATE_BIN) -path "$(MIGRATION_PATH)" -database "$(MIGRATION_DB)" force "$(VERSION)"

migrate-new:
	@if [ -z "$(MIGRATE_BIN)" ]; then \
		echo "❌ 'migrate' CLI not found. Please install it: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi
	@read -p "🆕 Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Migration name is required"; \
		exit 1; \
	fi; \
	$(MIGRATE_BIN) create -ext sql -dir "$(MIGRATION_PATH)" "$$name"

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
#                 CONSOLIDATED .PHONY TARGETS
# ============================================================

.PHONY: \
	help tools-install \
	build-dev dev-up dev-down dev clean-dev \
	build-prod prod-up prod-down prod clean-prod \
	docker-clean-all \
	graphql mocks \
	format lint lint-fix verify \
	test test-cover test-html-report test-ci test-clean \
	migrate-up migrate-down migrate-force migrate-new
