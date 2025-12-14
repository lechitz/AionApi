# ============================================================
#                          SEEDS
# ============================================================
.PHONY: seed-users seed-categories seed-all seed-tags seed-records seed-user1-all seed-everybody seed-clean-users seed-clean-categories seed-clean-tags seed-clean-records seed-clean-all seed-helper seed-setup seed-quick seed-api-caller seed-api-caller-bootstrap seed-api-caller-clean

POSTGRES_CONTAINER := postgres-dev
POSTGRES_USER := aion
POSTGRES_DB := aionapi

SEED_DEFAULT_PASSWORD := testpassword123
SEED_DEFAULT_PASSWORD_HASH := $2a$10$BIv0nYxelFEGDods46gtuuIpGH8NCThM1frbbhG5Ro/UqQ80ziwXS

# Build the seed-helper tool
seed-helper:
	@echo "🔨 Building seed-helper..."
	@go build -o bin/seed-helper ./cmd/seed-helper

# Generate .env.local with all seed variables (interactive setup)
seed-setup: seed-helper
	@echo "🚀 Setting up seed environment..."
	@read -p "Number of users to generate (default 10): " count; \
	count=$${count:-10}; \
	./bin/seed-helper generate-env $$count
	@echo ""
	@echo "✅ Setup complete! Now you can run: make seed-quick"

# Quick seed using .env.local (must run seed-setup first)
seed-quick:
	@echo "🌱 Quick seeding with .env.local..."
	@if [ ! -f infrastructure/db/seed/.env.local ]; then \
		echo "❌ .env.local not found. Run 'make seed-setup' first."; \
		exit 1; \
	fi
	@export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
	$(MAKE) seed-all-local

# Seeds require USER_TOKEN_TEST (bcrypt hash) — see infrastructure/db/seed/README.md
seed-users:
	@echo "Seeding users..."
	@if [ -z "$${USER_TOKEN_TEST:-}" ]; then \
		echo "USER_TOKEN_TEST not set; using default hash for password '$(SEED_DEFAULT_PASSWORD)'"; \
		USER_TOKEN_TEST="$(SEED_DEFAULT_PASSWORD_HASH)"; \
	fi; \
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v user_seed_password_hash="$${USER_TOKEN_TEST}" < infrastructure/db/seed/user.sql

seed-users-local:
	@echo "Seeding users (local .env if present)..."
	@if [ -z "$${SEED_USER_COUNT:-}" ]; then \
		if [ -f infrastructure/db/seed/.env.local ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
		elif [ -f infrastructure/db/seed/.env ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env | xargs); \
		fi; \
	fi; \
	test -n "$${USER_TOKEN_TEST:-}" || echo "USER_TOKEN_TEST not set (will attempt to generate if SEED_USER_COUNT is set)"; \
	if [ -n "$${SEED_USER_COUNT:-}" ]; then \
		echo "Generating $${SEED_USER_COUNT} users in-database using pgcrypto (password from DEV_PASSWORD or default)..."; \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v seed_count="$${SEED_USER_COUNT}" -v user_seed_password_plain="$${DEV_PASSWORD:-$(SEED_DEFAULT_PASSWORD)}" < infrastructure/db/seed/user_generate.sql; \
	else \
		if [ -z "$${USER_TOKEN_TEST:-}" ]; then \
			echo "USER_TOKEN_TEST not set; using default hash for password '$(SEED_DEFAULT_PASSWORD)'"; \
			USER_TOKEN_TEST="$(SEED_DEFAULT_PASSWORD_HASH)"; \
		fi; \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v user_seed_password_hash="$${USER_TOKEN_TEST}" < infrastructure/db/seed/user.sql; \
	fi

seed-categories:
	@echo "Seeding categories..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/category.sql

seed-categories-local:
	@echo "Seeding categories (dynamic generation)..."
	@if [ -z "$${SEED_USER_COUNT:-}" ]; then \
		if [ -f infrastructure/db/seed/.env.local ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
		elif [ -f infrastructure/db/seed/.env ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env | xargs); \
		fi; \
	fi; \
	if [ -n "$${SEED_USER_COUNT:-}" ]; then \
		echo "Generating categories for $${SEED_USER_COUNT} users..."; \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v seed_count="$${SEED_USER_COUNT}" < infrastructure/db/seed/category_generate.sql; \
	else \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/category.sql; \
	fi

seed-tags:
	@echo "Seeding tags..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/tags.sql

seed-tags-local:
	@echo "Seeding tags (dynamic generation)..."
	@if [ -z "$${SEED_USER_COUNT:-}" ]; then \
		if [ -f infrastructure/db/seed/.env.local ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
		elif [ -f infrastructure/db/seed/.env ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env | xargs); \
		fi; \
	fi; \
	if [ -n "$${SEED_USER_COUNT:-}" ]; then \
		echo "Generating tags for $${SEED_USER_COUNT} users..."; \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v seed_count="$${SEED_USER_COUNT}" < infrastructure/db/seed/tags_generate.sql; \
	else \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/tags.sql; \
	fi

seed-records:
	@echo "Seeding records..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/records.sql

seed-records-local:
	@echo "Seeding records (dynamic generation)..."
	@if [ -z "$${SEED_USER_COUNT:-}" ]; then \
		if [ -f infrastructure/db/seed/.env.local ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env.local | xargs); \
		elif [ -f infrastructure/db/seed/.env ]; then \
			export $$(grep -v '^#' infrastructure/db/seed/.env | xargs); \
		fi; \
	fi; \
	if [ -n "$${SEED_USER_COUNT:-}" ]; then \
		echo "Generating records for $${SEED_USER_COUNT} users (7 days of data)..."; \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v seed_count="$${SEED_USER_COUNT}" -v days=7 < infrastructure/db/seed/records_generate.sql; \
	else \
		docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/records.sql; \
	fi

seed-all: seed-users seed-categories seed-tags seed-records
	@echo "✅ All seeds applied."

seed-all-local: seed-users-local seed-categories-local seed-tags-local seed-records-local
	@echo "✅ All seeds applied (local)."

# Convenience target: seeds the full dataset for the default user (id=1).
# Keeps existing order to satisfy FKs.
seed-user1-all: seed-users seed-categories seed-tags seed-records
	@echo "✅ User 1 dataset applied."

# Alias to seed everything available; easier mnemonic than seed-all.
seed-everybody: seed-all
	@echo "✅ Everyone seeded."

seed-api-caller:
	@echo "🌐 Seeding via API (HTTP/GraphQL)..."
	@go run ./cmd/api-seed-caller

seed-api-caller-bootstrap:
	@echo "🌐 Seeding via API (bootstrap: cria usuário se necessário)..."
	@API_CALLER_AUTO_CREATE=true go run ./cmd/api-seed-caller

seed-api-caller-clean:
	@echo "🧹 Limpando via API (soft delete de records, sem criar nada)..."
	@API_CALLER_CLEAN=true API_CALLER_ONLY_CLEAN=true go run ./cmd/api-seed-caller

N ?= 1
ifdef n
N := $(n)
endif

seed-caller:
	@echo "🌐 Seeding via API (multi-user: count=$(N))..."
	@API_CALLER_COUNT=$(N) API_CALLER_AUTO_CREATE=true go run ./cmd/api-seed-caller

seed-api-caller-many: seed-caller

# Populate N users (default 10) with categories/tags/records via SQL generators (deterministic, generic naming).
POPULATE ?= 10
populate:
	@echo "🧹 Cleaning tables and populating $(N) users (password=$(SEED_DEFAULT_PASSWORD))..."
	@$(MAKE) seed-clean-all
	@SEED_USER_COUNT=$(N) DEV_PASSWORD=$(SEED_DEFAULT_PASSWORD) $(MAKE) seed-users-local seed-categories-local seed-tags-local seed-records-local
	@echo "✅ Populate completed for $(N) users."

# --- Clean helpers (dev-only): truncate seeded tables safely and reset IDs ---
# NOTE: Intended for local dev. These use TRUNCATE ... RESTART IDENTITY CASCADE.
seed-clean-users:
	@echo "🧹 Truncating users (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.users RESTART IDENTITY CASCADE;"

seed-clean-categories:
	@echo "🧹 Truncating tag_categories (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.tag_categories RESTART IDENTITY CASCADE;"

seed-clean-tags:
	@echo "🧹 Truncating tags (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.tags RESTART IDENTITY CASCADE;"

seed-clean-records:
	@echo "🧹 Truncating records (dev only)..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "TRUNCATE aion_api.records RESTART IDENTITY CASCADE;"

seed-clean-all: seed-clean-records seed-clean-tags seed-clean-categories seed-clean-users
	@echo "✅ All seeded tables truncated (dev only)."
