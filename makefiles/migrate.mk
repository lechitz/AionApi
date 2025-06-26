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
	@echo "Running all migrations (up)..."
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
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "❌ Migration name is required"; \
		exit 1; \
	fi; \
	$(MIGRATE_BIN) create -ext sql -dir "$(MIGRATION_PATH)" "$$name"