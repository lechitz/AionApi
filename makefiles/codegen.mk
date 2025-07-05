# ============================================================
#                           CODE GENERATION
# ============================================================

.PHONY: graphql mocks

graphql:
	cd internal/adapters/primary/graph && go run github.com/99designs/gqlgen generate

mocks:
	@echo "Generating mocks for output ports and usecases..."
	@mkdir -p tests/mocks/token tests/mocks/user tests/mocks/security tests/mocks/logger tests/mocks/category
	@echo "→ TokenStore"
	mockgen -source=internal/core/ports/output/token_output.go \
	  -destination=tests/mocks/token/mock_token_store.go \
	  -package=tokenmocks \
	  -mock_names=Store=MockTokenStore
	@echo "→ UserStore"
	mockgen -source=internal/core/ports/output/user_output.go \
	  -destination=tests/mocks/user/mock_user_store.go \
	  -package=usermocks \
	  -mock_names=UserStore=MockUserStore
	@echo "→ CategoryStore"
	mockgen -source=internal/core/ports/output/category_output.go \
	  -destination=tests/mocks/category/mock_category_store.go \
	  -package=categorymocks \
	  -mock_names=CategoryStore=MockCategoryStore
	@echo "→ SecurityStore"
	mockgen -source=internal/core/ports/output/hasher_output.go \
	  -destination=tests/mocks/security/mock_security_store.go \
	  -package=securitymocks \
	  -mock_names=Store=MockSecurityStore
	@echo "→ Logger"
	mockgen -source=internal/core/ports/output/logger_output.go \
	  -destination=tests/mocks/logger/mock_logger.go \
	  -package=loggermocks \
	  -mock_names=Logger=MockLogger
	@echo "✅  All mocks generated successfully."
