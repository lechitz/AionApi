# ============================================================
#                          SEEDS
# ============================================================
.PHONY: seed-users seed-categories seed-all

POSTGRES_CONTAINER := postgres-dev
POSTGRES_USER := aion
POSTGRES_DB := aionapi

seed-users:
	@echo "Seeding users..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infra/db/seeds/user.sql

seed-categories:
	@echo "Seeding categories..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infra/db/seeds/category.sql

seed-all: seed-users seed-categories
	@echo "✅ All seeds applied."