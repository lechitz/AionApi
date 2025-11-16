# ============================================================
#                          SEEDS
# ============================================================
.PHONY: seed-users seed-categories seed-all seed-tags seed-records

POSTGRES_CONTAINER := postgres-dev
POSTGRES_USER := aion
POSTGRES_DB := aionapi

seed-users:
	@echo "Seeding users..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/user.sql

seed-categories:
	@echo "Seeding categories..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/category.sql

seed-tags:
	@echo "Seeding tags..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/tags.sql

seed-records:
	@echo "Seeding records..."
	@docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < infrastructure/db/seed/records.sql

seed-all: seed-users seed-categories seed-tags seed-records
	@echo "✅ All seeds applied."