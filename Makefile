# ========================
# Global Variables
# ========================
include .env
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

docker-build-dev:
	docker build -t $(APPLICATION_NAME):dev .

docker-compose-dev-up: docker-compose-dev-down
	docker-compose -f $(COMPOSE_FILE_DEV) rm -f -v postgres
	docker-compose -f $(COMPOSE_FILE_DEV) up

docker-compose-dev-down:
	docker-compose -f $(COMPOSE_FILE_DEV) down -v

docker-build-run-dev: docker-build-dev docker-compose-dev-up

docker-clean-dev:
	docker rm -f $(shell docker ps -a --filter "name=dev" -q) || true
	docker volume rm $(shell docker volume ls --filter "name=dev" -q) || true
	docker rmi -f $(shell docker images --filter "reference=*dev*" -q) || true

# ========================
# Production Environment
# ========================
.PHONY: docker-build-prod docker-compose-prod-up docker-compose-prod-down docker-build-run-prod docker-clean-prod

docker-build-prod:
	docker build -t $(APPLICATION_NAME):prod .

docker-compose-prod-up: docker-compose-prod-down
	docker-compose -f $(COMPOSE_FILE_PROD) up

docker-compose-prod-down:
	docker-compose -f $(COMPOSE_FILE_PROD) down -v

docker-build-run-prod: docker-build-prod docker-compose-prod-up

docker-clean-prod:
	docker rm -f $(shell docker ps -a --filter "name=prod" -q) || true
	docker volume rm $(shell docker volume ls --filter "name=prod" -q) || true
	docker rmi -f $(shell docker images --filter "reference=*prod*" -q) || true

# ========================
# General Docker Commands
# ========================
.PHONY: docker-clean-all
docker-clean-all:
	# Remove all containers:
	docker rm -f $(shell docker ps -a -q) || true
	# Remove all volumes:
	docker volume rm $(shell docker volume ls -q) || true
	# Remove all images:
	docker rmi -f $(shell docker images -a -q) || true

# ========================
# Go Testing Commands
# ========================
.PHONY: test-cover
test-cover:
	go test ./... -coverprofile=coverage_tmp.out
	cat coverage_tmp.out | grep -v "Mock" > coverage.out
	rm -f coverage_tmp.out
	go tool cover -html=coverage.out
