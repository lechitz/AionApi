# ============================================================
#                           CODE GENERATION
# ============================================================

.PHONY: graphql mocks

OUTPUT_PKG := github.com/lechitz/AionApi/internal/core/ports/output

graphql:
	@echo "Generating GraphQL code with gqlgen..."
	cd ../internal/adapters/primary/graph && \
	go run github.com/99designs/gqlgen generate
	@go mod tidy
	@echo "✅  GraphQL code generated successfully."

mocks:
	@echo "Generating mocks for ports and services..."
	@mkdir -p tests/mocks

	# ---- Token ----
	@echo "→ Output: TokenProvider"
	mockgen -destination=tests/mocks/token_provider_mock.go \
	  -package=mocks \
	  -mock_names TokenProvider=TokenProvider \
	  $(OUTPUT_PKG) TokenProvider

	@echo "→ Output: TokenStore"
	mockgen -destination=tests/mocks/token_store_mock.go \
	  -package=mocks \
	  -mock_names TokenStore=TokenStore \
	  $(OUTPUT_PKG) TokenStore

	# ---- User ----
	@echo "→ Output: UserRepository"
	mockgen -destination=tests/mocks/user_repository_mock.go \
	  -package=mocks \
	  -mock_names UserRepository=UserRepository \
	  $(OUTPUT_PKG) UserRepository

	# ---- Category ----
	@echo "→ Output: CategoryRepository"
	mockgen -destination=tests/mocks/category_repository_mock.go \
	  -package=mocks \
	  -mock_names CategoryRepository=CategoryRepository \
	  $(OUTPUT_PKG) CategoryRepository

	# ---- Hasher ----
	@echo "→ Output: Hasher"
	mockgen -destination=tests/mocks/hasher_mock.go \
	  -package=mocks \
	  -mock_names Hasher=Hasher \
	  $(OUTPUT_PKG) Hasher

	# ---- Logger ----
	@echo "→ Output: ContextLogger"
	mockgen -destination=tests/mocks/logger_mock.go \
	  -package=mocks \
	  -mock_names ContextLogger=ContextLogger \
	  $(OUTPUT_PKG) ContextLogger

	@echo "✅  All mocks generated successfully."
