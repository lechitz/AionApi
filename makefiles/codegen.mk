# ============================================================
#                           CODE GENERATION
# ============================================================

.PHONY: graphql mocks

graphql:
	@echo "Generating GraphQL code with gqlgen..."
	cd internal/adapters/primary/graph && \
	go run github.com/99designs/gqlgen generate
	@go mod tidy
	@echo "✅  GraphQL code generated successfully."

mocks:
	@echo "Generating mocks for output ports and usecases..."
	@mkdir -p tests/mocks
	@echo "→ TokenStore"
	mockgen -source=internal/core/ports/output/token_output.go \
	  -destination=tests/mocks/mock_token_store.go \
	  -package=mocks \
	  -mock_names=Store=MockTokenStore
	@echo "→ TokenUsecase"
	mockgen -source=internal/core/ports/input/token_input.go \
      -destination=tests/mocks/mock_token_usecase.go \
      -package=mocks \
      -mock_names=TokenService=MockTokenUsecase
	@echo "→ UserStore"
	mockgen -source=internal/core/ports/output/user_output.go \
	  -destination=tests/mocks/mock_user_store.go \
	  -package=mocks \
	  -mock_names=UserStore=MockUserStore
	@echo "→ CategoryStore"
	mockgen -source=internal/core/ports/output/category_output.go \
	  -destination=tests/mocks/mock_category_store.go \
	  -package=mocks \
	  -mock_names=CategoryStore=MockCategoryStore
	@echo "→ SecurityStore"
	mockgen -source=internal/core/ports/output/password_hasher_output.go \
	  -destination=tests/mocks/mock_password_hasher_store.go \
	  -package=mocks \
	  -mock_names=Store=MockPasswordHasher
	@echo "→ Logger"
	mockgen -source=internal/core/ports/output/logger_output.go \
	  -destination=tests/mocks/mock_logger.go \
	  -package=mocks \
	  -mock_names=ContextLogger=MockLogger
	@echo "✅  All mocks generated successfully."
