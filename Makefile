# ========================
# Global Variables
# ========================
APPLICATION_NAME := aion-api
HTTP_PORT := 5001
COMPOSE_FILE_DEV := docker-compose-dev.yaml
COMPOSE_FILE_PROD := docker-compose-prod.yaml

# ========================
# Help Section
# ========================
.PHONY: help
help:
	@echo ""
	@echo "\033[1;33m\033[1mAionApi - Developer CLI Help\033[0m"
	@echo ""

	@echo "\033[1;34m🔵 Docker Compose Commands:\033[0m"
	@echo "  \033[1;36mdev-up\033[0m     → Start development environment (resets DB)"
	@echo "  \033[1;36mdev-down\033[0m   → Stop development and remove volumes"
	@echo "  \033[1;36mdocker compose-prod-up\033[0m    → Start production (keeps DB)"
	@echo "  \033[1;36mdocker compose-prod-down\033[0m  → Stop production environment"
	@echo ""

	@echo "\033[1;34m🔵 Docker Build Commands:\033[0m"
	@echo "  \033[1;36mbuild-dev\033[0m          → Build dev image"
	@echo "  \033[1;36mdocker-build-prod\033[0m         → Build prod image"
	@echo "  \033[1;36mdev\033[0m                       → Build & start dev environment"
	@echo "  \033[1;36mprod\033[0m                      → Build & start prod environment"
	@echo ""

	@echo "\033[1;34m🔵 Docker Cleanup Commands:\033[0m"
	@echo "  \033[1;36mclean-dev\033[0m          → Clean dev containers, volumes, images"
	@echo "  \033[1;36mdocker-clean-prod\033[0m         → Clean prod containers, volumes, images"
	@echo "  \033[1;36mdocker-clean-all\033[0m          → Remove ALL containers, volumes, images"
	@echo ""

	@echo "\033[1;34m🔵 Testing Commands:\033[0m"
	@echo "  \033[1;36mtest-cover\033[0m                → Run tests and generate coverage report"
	@echo ""


# ========================
# Development Environment
# ========================
.PHONY: build-dev dev-up dev-down dev clean-dev

build-dev: clean-dev
	docker build -t $(APPLICATION_NAME):dev .

dev-up: dev-down
	@echo "Starting Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) up

dev-down:
	@echo "Stopping Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker compose -f $(COMPOSE_FILE_DEV) down -v

dev: clean-dev build-dev dev-up

clean-dev:
	@echo "\033[1;33m Cleaning dev containers...\033[0m"
	@containers=$$(docker ps -a --filter "name=dev" -q); \
	if [ -n "$$containers" ]; then \
		docker rm -f $$containers; \
	else \
		echo "No dev containers to remove."; \
	fi

	@echo "\033[1;33m Cleaning dev volumes...\033[0m"
	@volumes=$$(docker volume ls --filter "name=dev" -q); \
	if [ -n "$$volumes" ]; then \
		docker volume rm $$volumes; \
	else \
		echo "No dev volumes to remove."; \
	fi

	@echo "\033[1;33m Cleaning dev images...\033[0m"
	@images=$$(docker images --filter "reference=$(APPLICATION_NAME):dev" -q); \
	if [ -n "$$images" ]; then \
		docker rmi -f $$images || true; \
	else \
		echo "No dev images to remove."; \
	fi

# ========================
# Production Environment
# ========================
.PHONY: docker-build-prod docker-compose-prod-up docker-compose-prod-down prod docker-clean-prod

docker-build-prod: docker-clean-prod
	docker build -t $(APPLICATION_NAME):prod .

docker-compose-prod-up: docker-compose-prod-down
	docker-compose -f $(COMPOSE_FILE_PROD) up

docker-compose-prod-down:
	docker-compose -f $(COMPOSE_FILE_PROD) down -v

prod: docker-build-prod docker-compose-prod-up

docker-clean-prod:
	@containers=$$(docker ps -a --filter "name=prod" -q); \
	if [ -n "$$containers" ]; then \
		echo "Removing prod containers..."; \
		docker rm -f $$containers; \
	else \
		echo "No prod containers to remove."; \
	fi

	@volumes=$$(docker volume ls --filter "name=prod" -q); \
	if [ -n "$$volumes" ]; then \
		echo "Removing prod volumes..."; \
		docker volume rm $$volumes; \
	else \
		echo "No prod volumes to remove."; \
	fi

	@images=$$(docker images --filter "reference=*prod*" -q); \
	if [ -n "$$images" ]; then \
		echo "Removing prod images..."; \
		docker rmi -f $$images; \
	else \
		echo "No prod images to remove."; \
	fi

# ========================
# General Docker Commands
# ========================
.PHONY: docker-clean-all
docker-clean-all:
	@containers=$$(docker ps -a -q); \
	if [ -n "$$containers" ]; then \
		echo "Removing ALL containers..."; \
		docker rm -f $$containers; \
	else \
		echo "No containers to remove."; \
	fi

	@volumes=$$(docker volume ls -q); \
	if [ -n "$$volumes" ]; then \
		echo "Removing ALL volumes..."; \
		docker volume rm $$volumes; \
	else \
		echo "No volumes to remove."; \
	fi

	@images=$$(docker images -a -q); \
	if [ -n "$$images" ]; then \
		echo "Removing ALL images..."; \
		docker rmi -f $$images; \
	else \
		echo "No images to remove."; \
	fi

# ========================
# Go Testing Commands
# ========================
.PHONY: test
test:
	@echo "📦 Running unit tests..."
	go test ./... -v

.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage report..."
	go test ./... -coverprofile=coverage_tmp.out -v
	@echo "Filtering out mock files from coverage..."
	cat coverage_tmp.out | grep -v "Mock" > coverage.out
	@rm -f coverage_tmp.out
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out

.PHONY: test-ci
test-ci:
	@echo "Running CI tests with coverage output..."
	go test ./... -coverprofile=coverage.out -v

.PHONY: test-clean
test-clean:
	@echo "Cleaning up coverage reports..."
	@rm -f coverage.out coverage_tmp.out

# ========================
# HTML Test Report (go-test-html-report)
# ========================
.PHONY: test-html-report
test-html-report:
	@echo "🧪 Running tests and generating JSON output..."
	go test ./... -json > docs/coverage/report.json
	@echo "📄 Generating HTML report..."
	go-test-html-report -f docs/coverage/report.json -o docs/coverage/
	@echo "✅ HTML report generated at: docs/coverage/report.html"

# ========================
# Generate GraphQL
# ========================
.PHONY: graphql

graphql:
	cd adapters/primary/graph && go run github.com/99designs/gqlgen generate

# ========================
# Mock Generation Commands
# ========================
.PHONY: mocks
mocks:
	@echo "Generating mocks for output ports and usecases..."
	@mkdir -p tests/mocks/token tests/mocks/user tests/mocks/security tests/mocks/logger tests/mocks/category

	@echo "→ TokenStore"
	mockgen -source=internal/core/ports/output/cache/token.go \
		-destination=tests/mocks/token/mock_token_store.go \
		-package=tokenmocks

	@echo "→ TokenUsecase"
	mockgen -source=internal/core/usecase/token/token_usecase.go \
		-destination=tests/mocks/token/mock_token_usecase.go \
		-package=tokenmocks TokenUsecase

	@echo "→ UserStore"
	mockgen -source=internal/core/ports/output/db/user.go \
		-destination=tests/mocks/user/mock_user_store.go \
		-package=usermocks

	@echo "→ CategoryStore"
	mockgen -source=internal/core/ports/output/db/category.go \
		-destination=tests/mocks/category/mock_category_store.go \
		-package=categorymocks

	@echo "→ SecurityStore"
	mockgen -source=internal/core/ports/output/security/hasher.go \
		-destination=tests/mocks/security/mock_security_store.go \
		-package=securitymocks

	@echo "→ Logger"
	mockgen -source=internal/core/ports/output/logger/logger.go \
		-destination=tests/mocks/logger/mock_logger.go \
		-package=loggermocks

	@echo "✅ All mocks generated successfully."

