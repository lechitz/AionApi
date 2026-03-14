# ============================================================
#                          API CALLS
# ============================================================
.PHONY: call-login call-health call-users call-me call-chat call-graphql

API_HOST      ?= http://localhost:5001
API_CONTEXT   ?= /aion
API_ROOT      ?= /api/v1
API_BASE      := $(API_CONTEXT)$(API_ROOT)
API_URL       := $(API_HOST)$(API_BASE)

HEALTH_ROUTE  ?= /health
HEALTH_URL    := $(API_HOST)$(API_CONTEXT)$(HEALTH_ROUTE)
GRAPHQL_PATH  ?= /graphql

USER          ?= user
PASS          ?= user
TOKEN_FILE    ?= .cache/aion/token

# curl options can be overridden if needed (e.g., -v for verbose)
CURL_OPTS     ?= -s

# Internal helper: resolve Authorization header from $TOKEN or token file (if present).
define WITH_AUTH
TOKEN_VALUE="$(TOKEN)"; \
if [ -z "$$TOKEN_VALUE" ] && [ -f "$(TOKEN_FILE)" ]; then TOKEN_VALUE=$$(cat "$(TOKEN_FILE)"); fi; \
AUTH_ARG=""; [ -n "$$TOKEN_VALUE" ] && AUTH_ARG="-H \"Authorization: Bearer $$TOKEN_VALUE\""; \
curl $(CURL_OPTS) $$AUTH_ARG
endef

call-login:
	@echo "🔑 Logging in as '$(USER)' at $(API_URL)/auth/login"
	@response=$$(curl $(CURL_OPTS) -X POST "$(API_URL)/auth/login" \
		-H "Content-Type: application/json" \
		-d '{"username":"$(USER)","password":"$(PASS)"}'); \
	echo "$$response"; \
	if [ "$(SAVE_TOKEN)" = "true" ]; then \
		token=$$(printf '%s' "$$response" | python3 -c 'import json,sys; data=json.load(sys.stdin); print(data.get("token",""))' || true); \
		if [ -n "$$token" ]; then \
			mkdir -p $(dir $(TOKEN_FILE)); \
			printf '%s' "$$token" > "$(TOKEN_FILE)"; \
			printf '\n(Token saved to %s)\n' "$(TOKEN_FILE)"; \
		else \
			echo "⚠️  Could not parse token from response."; \
		fi; \
	fi

call-health:
	@echo "🩺 GET $(HEALTH_URL)"
	@$(WITH_AUTH) "$(HEALTH_URL)"

call-users:
	@echo "👥 GET $(API_URL)/user/all"
	@$(WITH_AUTH) "$(API_URL)/user/all"

call-me:
	@echo "🙋 GET $(API_URL)/user/me"
	@$(WITH_AUTH) "$(API_URL)/user/me"

call-chat:
	@if [ -z "$(MESSAGE)" ]; then echo "Usage: make call-chat MESSAGE='hi there' [TOKEN=...]"; exit 1; fi
	@echo "💬 POST $(API_URL)/chat"
	@$(WITH_AUTH) -H "Content-Type: application/json" \
		-d '{"message":"$(MESSAGE)"}' \
		-X POST "$(API_URL)/chat"

# GraphQL call. Provide QUERY or QUERY_FILE. Example:
# make call-graphql QUERY='query { categories { id name } }'
# make call-graphql QUERY_FILE=./tmp/query.graphql
call-graphql:
	@echo "📡 POST $(API_URL)$(GRAPHQL_PATH)"
	@query_payload=""; \
	if [ -n "$(QUERY_FILE)" ]; then \
		if [ ! -f "$(QUERY_FILE)" ]; then echo "Query file not found: $(QUERY_FILE)"; exit 1; fi; \
		query_payload=$$(cat "$(QUERY_FILE)"); \
	elif [ -n "$(QUERY)" ]; then \
		query_payload='$(QUERY)'; \
	else \
		echo "Usage: make call-graphql QUERY='query { __typename }' [TOKEN=...]"; exit 1; \
	fi; \
	$(WITH_AUTH) -H "Content-Type: application/json" \
		-d "$$(printf '{"'\"'\"'query'\"'\"':'"'"'%s'"'"','"'\"'\"'variables'\"'\"':{}}' "$$query_payload")" \
		-X POST "$(API_URL)$(GRAPHQL_PATH)"
