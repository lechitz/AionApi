# ============================================================
#                         SWAG DOCS
# ============================================================

SWAG_VERSION ?= v1.16.6
SWAG_BIN     ?= $(shell command -v swag || echo "$$(go env GOPATH)/bin/swag")
SWAG_OUT     ?= internal/platform/docs
SWAG_MAIN    ?= cmd/aion-api/main.go

SWAG_DIRS    ?= cmd/aion-api,internal

.PHONY: tools.s swag docs.gen docs.clean docs.check-dirty

tools.s:
	@echo ">> ensuring tools"

swag: tools.s
	@echo ">> installing swag $(SWAG_VERSION)"
	@GOBIN=$$(go env GOPATH)/bin go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

docs.gen: swag
	@echo ">> generating swagger docs"
	@$(SWAG_BIN) init \
		-g $(SWAG_MAIN) \
		-o $(SWAG_OUT) \
		--parseDependency \
		--parseInternal \
		-d $(SWAG_DIRS)

docs.clean:
	@echo ">> cleaning swagger docs"
	@rm -rf $(SWAG_OUT)

# Fail the build if documentation generation is out-of-date.
docs.check-dirty:
	@git diff --quiet -- $(SWAG_OUT) || (echo "Swagger docs out-of-date. Run 'make docs.gen'."; exit 1)
