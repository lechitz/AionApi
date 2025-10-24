# ============================================================
#                   GLOBAL VARIABLES & CONFIG
# ============================================================

APPLICATION_NAME := aion-api

COMPOSE_FILE_DEV  := infrastructure/docker/environments/dev/docker-compose-dev.yaml
ENV_FILE_DEV      := infrastructure/docker/environments/dev/.env.dev
COMPOSE_FILE_PROD := infrastructure/docker/environments/prod/docker-compose-prod.yaml
ENV_FILE_PROD     := infrastructure/docker/environments/prod/.env.prod

COVERAGE_DIR = tests/coverage

# --- MIGRATION CONFIG ---
MIGRATION_PATH := infrastructure/db/migrations
MIGRATION_DB   ?= $(DB_URL)
MIGRATE_BIN    := $(shell command -v migrate 2> /dev/null)

# ============================================================
#                HELP & TOOLING SECTION
# ============================================================

.PHONY: all help tools-install tools.check

all: help

help:
	@echo ""
	@echo ""
	@echo "\033[48;5;235;33m┃==================================================================================================================┃\033[0m"
	@echo "\033[48;5;235;33m┃                                            \033[1mAION API - CLI COMMANDS                                               ┃\033[0m"
	@echo "\033[48;5;235;33m┃==================================================================================================================┃\033[0m"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mTOOLING ┃\033[0m"
	@echo ""
	@echo "    \033[0m tools-install        \033[1;37m    →  Install all development tools (goimports, golines, gofumpt, golangci-lint)"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mDOCKER ENVIRONMENT COMMANDS ┃\033[0m"
	@echo ""
	@echo "  \033[1;37m- [DEV]\033[0m"
	@echo ""
	@echo "    \033[0m build-dev          \033[1;37m      →  Build the development Docker image"
	@echo "    \033[0m dev-up             \033[1;37m      →  Start the development environment (resets DB)"
	@echo "    \033[0m dev-down           \033[1;37m      →  Stop and remove dev environment containers/volumes"
	@echo "    \033[0m clean-dev          \033[1;37m      →  Clean all dev containers, volumes, and images"
	@echo ""
	@echo "  \033[1;37m- [PROD]\033[0m"
	@echo ""
	@echo "    \033[0m build-prod         \033[1;37m      →  Build the production Docker image"
	@echo "    \033[0m prod-up            \033[1;37m      →  Start the production environment"
	@echo "    \033[0m prod-down          \033[1;37m      →  Stop and remove prod environment containers/volumes"
	@echo "    \033[0m clean-prod         \033[1;37m      →  Clean all prod containers, volumes, and images"
	@echo ""
	@echo "  \033[1;37m- [GENERAL]\033[0m"
	@echo ""
	@echo "    \033[0m docker-clean-all    \033[1;37m     →  Remove ALL Docker containers, volumes, and images"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mCODE GENERATION ┃\033[0m"
	@echo ""
	@echo "    \033[0m graphql             \033[1;37m     →  Generate GraphQL files with gqlgen"
	@echo "    \033[0m mocks               \033[1;37m     →  Generate all GoMock mocks"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mCODE QUALITY ┃\033[0m"
	@echo ""
	@echo "    \033[0m format               \033[1;37m    →  Format Go code using goimports/golines/gofumpt"
	@echo "    \033[0m lint                 \033[1;37m    →  Run golangci-lint (static code analysis)"
	@echo "    \033[0m lint-fix             \033[1;37m    →  Run golangci-lint with --fix (auto-fix where possible)"
	@echo "    \033[0m verify               \033[1;37m    →  Run full pre-commit pipeline (format, mocks, lint, tests, coverage, codegen)"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mMIGRATIONS ┃\033[0m"
	@echo ""
	@echo "    \033[0m migrate-up           \033[1;37m    →  Run all migrations (up)"
	@echo "    \033[0m migrate-down         \033[1;37m    →  Rollback the last migration"
	@echo "    \033[0m migrate-force VERSION=\033[1;32mX \033[1;37m →  Force DB to specific version"
	@echo "    \033[0m migrate-new          \033[1;37m    →  Create new migration (with prompt)"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mSEEDS ┃\033[0m"
	@echo ""
	@echo "    \033[0m seed-users           \033[1;37m    →  Seed the users table"
	@echo "    \033[0m seed-categories      \033[1;37m    →  Seed the categories table"
	@echo "    \033[0m seed-all             \033[1;37m    →  Run all seed scripts"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mTESTING ┃\033[0m"
	@echo ""
	@echo "    \033[0m test                 \033[1;37m    →  Run unit tests"
	@echo "    \033[0m test-cover           \033[1;37m    →  Run tests with coverage report (excludes mocks)"
	@echo "    \033[0m test-html-report     \033[1;37m    →  Generate HTML test report (requires go-test-html-report)"
	@echo ""
	@echo ""
	@echo " 🔶 \033[48;5;235;33m┃ \033[1mAPI DOCS (SWAGGER) ┃\033[0m"
	@echo ""
	@echo "    \033[0m docs.gen            \033[1;37m     →  Generate Swagger artifacts (docs.go, swagger.json/yaml)"
	@echo "    \033[0m docs.check-dirty    \033[1;37m     →  Fail if Swagger artifacts are out-of-date"
	@echo "    \033[0m docs.clean          \033[1;37m     →  Remove generated Swagger artifacts"
	@echo ""
	@echo ""
	@echo "\033[48;5;235;33m┃==================================================================================================================┃\033[0m"
	@echo ""

# ============================================================
#                 CONSOLIDATED .PHONY TARGETS
# ============================================================

-include makefiles/*.mk

.PHONY: graphql mocks docs.gen docs.validate docs.check-dirty lint test test-cover test-ci test-clean

# Short aliases
.PHONY: install-tools
install-tools: tools-install

.PHONY: \
	help tools-install tools.check \
	build-dev dev-up dev-down dev clean-dev \
	build-prod prod-up prod-down prod clean-prod \
	docker-clean-all \
	graphql mocks \
	format lint lint-fix verify \
	test test-cover test-html-report test-ci test-clean \
	migrate-up migrate-down migrate-force migrate-new \
	docs.gen docs.check-dirty docs.clean docs.validate

docs-serve:
	@.venv-docs/bin/mkdocs serve

docs-build:
	@.venv-docs/bin/mkdocs build
