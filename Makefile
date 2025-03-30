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
	@echo "Available commands:"
	@echo ""
	@echo "Docker Compose Commands:"
	@echo "  docker-compose-dev-up     Start the development environment (resetting the database)."
	@echo "  docker-compose-dev-down   Stop the development environment and remove volumes."
	@echo "  docker-compose-prod-up    Start the production environment (database is not reset)."
	@echo "  docker-compose-prod-down  Stop the production environment."
	@echo ""
	@echo "Docker Build Commands:"
	@echo "  docker-build-dev          Build the development Docker image."
	@echo "  docker-build-prod         Build the production Docker image."
	@echo "  docker-build-run-dev      Build and start the development environment."
	@echo "  docker-build-run-prod     Build and start the production environment."
	@echo ""
	@echo "Docker Cleanup Commands:"
	@echo "  docker-clean-dev          Clean development containers, volumes, and images."
	@echo "  docker-clean-prod         Clean production containers, volumes, and images."
	@echo "  docker-clean-all          Clean all containers, volumes, and images."
	@echo ""
	@echo "Testing Commands:"
	@echo "  test-cover                Run tests and generate a coverage report."

# ========================
# Development Environment
# ========================
.PHONY: docker-build-dev docker-compose-dev-up docker-compose-dev-down docker-build-run-dev docker-clean-dev

docker-build-dev: docker-clean-dev
	docker build -t $(APPLICATION_NAME):dev .

docker-compose-dev-up: docker-compose-dev-down
	@echo "Starting Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) up

docker-compose-dev-down:
	@echo "Stopping Dev Environment..."
	export $$(cat .env.dev | grep -v '^#' | xargs) && docker-compose -f $(COMPOSE_FILE_DEV) down -v

docker-build-run-dev: docker-clean-dev docker-build-dev docker-compose-dev-up

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
.PHONY: docker-build-prod docker-compose-prod-up docker-compose-prod-down docker-build-run-prod docker-clean-prod

docker-build-prod: docker-clean-prod
	docker build -t $(APPLICATION_NAME):prod .

docker-compose-prod-up: docker-compose-prod-down
	docker-compose -f $(COMPOSE_FILE_PROD) up

docker-compose-prod-down:
	docker-compose -f $(COMPOSE_FILE_PROD) down -v

docker-build-run-prod: docker-build-prod docker-compose-prod-up

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
# Mock Generation Commands
# ========================
.PHONY: generate-mocks
generate-mocks:
	@echo "Generating mocks for User use cases..."
	mockgen -source=internal/core/usecase/user/create_user.go -destination=tests/mocks/user/mock_user_creator.go -package=mocks
	mockgen -source=internal/core/usecase/user/update_user.go -destination=tests/mocks/user/mock_user_updater.go -package=mocks
	mockgen -source=internal/core/usecase/user/delete_user.go -destination=tests/mocks/user/mock_user_deleter.go -package=mocks
	mockgen -source=internal/core/usecase/user/get_user.go -destination=tests/mocks/user/mock_user_retriever.go -package=mocks
	@echo "Mocks generated successfully!"