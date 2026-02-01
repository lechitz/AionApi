# GraphQL Documentation and Shared Queries
.PHONY: graphql.schema graphql.queries graphql.docs graphql.setup graphql.clean

ROOT_DIR := $(shell pwd)
DOCS_DIR := $(ROOT_DIR)/docs/graphql
SCHEMA_DIR := $(ROOT_DIR)/internal/adapter/primary/graphql/schema
SCHEMA_OUT := $(DOCS_DIR)/schema.graphql
QUERIES_DIR := $(ROOT_DIR)/graphql/queries

graphql.schema:
	@echo "📜 Exporting schema..."
	@mkdir -p "$(DOCS_DIR)"
	@echo "# Aion GraphQL Schema" > "$(SCHEMA_OUT)"
	@find "$(SCHEMA_DIR)" -name "*.graphqls" | sort | xargs cat >> "$(SCHEMA_OUT)"
	@echo "✅ Schema: $(SCHEMA_OUT)"

graphql.queries:
	@echo "📝 Creating queries..."
	@mkdir -p "$(QUERIES_DIR)/categories" "$(QUERIES_DIR)/tags" "$(QUERIES_DIR)/records"
	@printf 'query ListCategories { categories { id name description colorHex icon } }\n' > "$(QUERIES_DIR)/categories/list.graphql"
	@printf 'query ListTags { tags { id name categoryId icon } }\n' > "$(QUERIES_DIR)/tags/list.graphql"
	@printf 'query ListRecords($$limit: Int) { records(limit: $$limit) { id tagId } }\n' > "$(QUERIES_DIR)/records/list.graphql"
	@printf '# GraphQL Shared Queries\n' > "$(QUERIES_DIR)/README.md"
	@echo "✅ Queries: $(QUERIES_DIR)/"

graphql.docs: graphql.schema
	@printf '# GraphQL Documentation\n\nSchema: schema.graphql\nPlayground: http://localhost:5001/aion/api/v1/graphql/playground\n' > "$(DOCS_DIR)/README.md"
	@echo "✅ Docs ready"

graphql.clean:
	@rm -rf "$(DOCS_DIR)" "$(QUERIES_DIR)"
	@echo "🧹 Cleaned"

graphql.setup: graphql.schema graphql.queries graphql.docs
	@echo ""
	@echo "✅ GraphQL documentation setup complete!"
	@echo "   Schema:  $(SCHEMA_OUT)"
	@echo "   Queries: $(QUERIES_DIR)/"
	@echo ""
