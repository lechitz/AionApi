# ========================
# Global Variables
# ========================
APPLICATION_NAME := aion-api
PORT := 5001
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
	@echo "  \033[1;36mdocker-compose-dev-up\033[0m     → Start development environment (resets DB)"
	@echo "  \033[1;36mdocker-compose-dev-down\033[0m   → Stop development and remove volumes"
	@echo "  \033[1;36mdocker-compose-prod-up\033[0m    → Start production (keeps DB)"
	@echo "  \033[1;36mdocker-compose-prod-down\033[0m  → Stop production environment"
	@echo ""

	@echo "\033[1;34m🔵 Docker Build Commands:\033[0m"
	@echo "  \033[1;36mdocker-build-dev\033[0m          → Build dev image"
	@echo "  \033[1;36mdocker-build-prod\033[0m         → Build prod image"
	@echo "  \033[1;36mdev\033[0m                       → Build & start dev environment"
	@echo "  \033[1;36mprod\033[0m                      → Build & start prod environment"
	@echo ""

	@echo "\033[1;34m🔵 Docker Cleanup Commands:\033[0m"
	@echo "  \033[1;36mdocker-clean-dev\033[0m          → Clean dev containers, volumes, images"
	@echo "  \033[1;36mdocker-clean-prod\033[0m         → Clean prod containers, volumes, images"
	@echo "  \033[1;36mdocker-clean-all\033[0m          → Remove ALL containers, volumes, images"
	@echo ""

	@echo "\033[1;34m🔵 Testing Commands:\033[0m"
	@echo "  \033[1;36mtest-cover\033[0m                → Run tests and generate coverage report"
	@echo ""


# ========================
# Development Environment
# ========================
.PHONY: docker-build-dev docker-compose-dev-up docker-compose-dev-down dev docker-clean-dev

docker-build-dev: docker-clean-dev
	docker build -t $(APPLICATION_NAME):dev .

docker-compose-dev-up: docker-compose-dev-down
	@echo "Starting Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) up

docker-compose-dev-down:
	@echo "Stopping Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) down -v

dev: docker-clean-dev docker-build-dev docker-compose-dev-up

docker-clean-dev:
	@containers=$$(docker ps -a --filter "name=dev" -q); \
	if [ -n "$$containers" ]; then \
		echo "Removing dev containers..."; \
		docker rm -f $$containers; \
	else \
		echo "No dev containers to remove."; \
	fi

	@volumes=$$(docker volume ls --filter "name=dev" -q); \
	if [ -n "$$volumes" ]; then \
		echo "Removing dev volumes..."; \
		docker volume rm $$volumes; \
	else \
		echo "No dev volumes to remove."; \
	fi

	@images=$$(docker images --filter "reference=*dev*" -q); \
	if [ -n "$$images" ]; then \
		echo "Removing dev images..."; \
		docker rmi -f $$images; \
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
	go test ./... -json > internal/docs/coverage/report.json
	@echo "📄 Generating HTML report..."
	go-test-html-report -f internal/docs/coverage/report.json -o internal/docs/coverage/
	@echo "✅ HTML report generated at: internal/docs/coverage/report.html"

# ========================
# Mock Generation Commands
# ========================
.PHONY: mocks
mocks:
	@echo "Generating mocks for output ports..."
	@mkdir -p tests/mocks

	@echo "→ Token Repository"
	mockgen -source=internal/core/ports/output/cache/token.go -destination=tests/mocks/mock_token_repository.go -package=mocks

	@echo "→ User Repository"
	mockgen -source=internal/core/ports/output/db/user.go -destination=tests/mocks/mock_user_repository.go -package=mocks

	@echo "→ Security Hasher"
	mockgen -source=internal/core/ports/output/security/hasher.go -destination=tests/mocks/mock_security_store.go -package=mocks

	@echo "→ Redis Client"
	mockgen -source=internal/infrastructure/cache/redis.go -destination=tests/mocks/mock_redis_client.go -package=mocks

	@echo "→ Auth Service"
	mockgen -source=internal/core/ports/input/http/auth.go -destination=tests/mocks/mock_auth_service.go -package=mocks

	@echo "→ User Service"
	mockgen -source=internal/core/ports/input/http/user.go -destination=tests/mocks/mock_user_service.go -package=mocks

	@echo "All mocks generated successfully."