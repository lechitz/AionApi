# ============================================================
#                           CODE GENERATION
# ============================================================
# Purpose:
#   - Generate Go mocks from ports under:
#       internal/<context>/core/ports/output
#       internal/platform/ports/output
#   - Place ALL generated mocks FLAT into: tests/mocks/
#     (no nested folders; a single "mocks" package).
#
# Key features:
#   - Flat filename mode by default (uses file basenames).
#   - Collision guard: fails if two sources share the same basename.
#   - Optional namespaced mode (NAMESPACE=1) to avoid collisions by
#     prefixing filenames with the relative path (slashes -> "__").
#   - Context filter (CONTEXT=<name>) to generate only one context.
#
# Requirements:
#   - bash available (we rely on pipefail)
#   - Go toolchain installed
#   - mockgen installed:
#       go install go.uber.org/mock/mockgen@latest
#
# Common usage:
#   make mocks                 # generate all mocks (flat)
#   make mocks NAMESPACE=1     # generate all mocks with namespaced filenames
#   make mocks CONTEXT=user    # generate only "user" context
#   make clean_mocks           # remove generated mocks
#   make mocks-list            # audit which sources/targets will be used
#   make graphql               # run gqlgen and tidy modules

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.SILENT:

.PHONY: graphql mocks clean_mocks verify_mockgen mocks-list

# Repository root (works from subdirectories too)
ROOT_DIR := $(shell git rev-parse --show-toplevel 2>/dev/null || pwd)

# Output (FLAT): all mocks live directly under tests/mocks/
MOCKS_DIR := $(ROOT_DIR)/tests/mocks
MOCKS_PKG := mocks

# --------------------------------------------------------------------
# Discover source directories (absolute paths)
# --------------------------------------------------------------------
# Context outputs: internal/<context>/core/ports/output
OUTPUT_PORTS_DIRS_ALL := $(shell find "$(ROOT_DIR)/internal" -type d -path "*/core/ports/output" 2>/dev/null)

# Optional context filter: make mocks CONTEXT=user
ifeq ($(strip $(CONTEXT)),)
  OUTPUT_PORTS_DIRS := $(OUTPUT_PORTS_DIRS_ALL)
else
  OUTPUT_PORTS_DIRS := $(shell find "$(ROOT_DIR)/internal/$(CONTEXT)" -type d -path "*/core/ports/output" 2>/dev/null)
endif

# Platform outputs (cross-cutting)
PLATFORM_OUTPUT_DIR := $(ROOT_DIR)/internal/platform/ports/output

# --------------------------------------------------------------------
# Collect source files (*.go, excluding *_test.go) — absolute paths
# --------------------------------------------------------------------
MOCK_SOURCES := \
	$(shell \
		[ -n "$(OUTPUT_PORTS_DIRS)" ] && find $(OUTPUT_PORTS_DIRS) -type f -name "*.go" ! -name "*_test.go" 2>/dev/null; \
		[ -d "$(PLATFORM_OUTPUT_DIR)" ] && find "$(PLATFORM_OUTPUT_DIR)" -type f -name "*.go" ! -name "*_test.go" 2>/dev/null; \
	)

# --------------------------------------------------------------------
# Target filename strategies
# --------------------------------------------------------------------
# FLAT (default): use basenames only (e.g., user_output.go -> user_output_mock.go)
BASENAMES     := $(notdir $(MOCK_SOURCES))
FLAT_TARGETS  := $(addprefix $(MOCKS_DIR)/,$(patsubst %.go,%_mock.go,$(BASENAMES)))

# NAMESPACED (NAMESPACE=1): prefix with relative path without "internal/",
# replacing "/" with "__" (e.g., user/core/ports/output/user_output.go
# -> user__core__ports__output__user_output_mock.go).
REL_SOURCES          := $(patsubst $(ROOT_DIR)/%,%,$(MOCK_SOURCES))
REL_NO_INTERNAL      := $(patsubst internal/%,%,$(REL_SOURCES))
NAMESPACED_FILENAMES := $(patsubst %.go,%_mock.go,$(subst /,__,$(REL_NO_INTERNAL)))
NAMESPACED_TARGETS   := $(addprefix $(MOCKS_DIR)/,$(NAMESPACED_FILENAMES))

# Choose targets set
ifeq ($(strip $(NAMESPACE)),1)
  GENERATED_MOCKS := $(NAMESPACED_TARGETS)
else
  GENERATED_MOCKS := $(FLAT_TARGETS)
endif

# --------------------------------------------------------------------
# GraphQL (gqlgen)
# --------------------------------------------------------------------

# Centralized GraphQL adapter directory
GRAPH_DIR := internal/adapter/primary/graphql

# Optional list of schemas (for logging)
GRAPHQL_SOURCES := $(shell find "$(GRAPH_DIR)/schema" -type f -name "*.graphqls" 2>/dev/null)

graphql:
	@echo "Generating GraphQL code with gqlgen..."
	@echo "Schemas found:"; if [ -n "$(GRAPHQL_SOURCES)" ]; then printf "  %s\n" $(GRAPHQL_SOURCES); else echo "  (no .graphqls found)"; fi
	cd "$(GRAPH_DIR)" && go run github.com/99designs/gqlgen generate
	cd "$(ROOT_DIR)" && go mod tidy
	@echo "✅  GraphQL code generated successfully."

# --------------------------------------------------------------------
# Tooling checks
# --------------------------------------------------------------------
verify_mockgen:
	@if ! command -v mockgen >/dev/null 2>&1 ; then \
		echo "❌ 'mockgen' not found."; \
		echo "   Install it with: go install go.uber.org/mock/mockgen@latest"; \
		exit 1; \
	fi
	@mockgen -version || true

# --------------------------------------------------------------------
# Audit (print sources and targets)
# --------------------------------------------------------------------
mocks-list:
	@echo "ROOT_DIR: $(ROOT_DIR)"
	@echo
	@echo "OUTPUT_PORTS_DIRS:"; if [ -n "$(OUTPUT_PORTS_DIRS)" ]; then printf "  %s\n" $(OUTPUT_PORTS_DIRS); else echo "  (none)"; fi
	@echo "PLATFORM_OUTPUT_DIR: $(PLATFORM_OUTPUT_DIR)"
	@echo
	@echo "MOCKS_DIR: $(MOCKS_DIR)"
	@echo "NAMESPACE=$(NAMESPACE)"
	@echo "CONTEXT=$(CONTEXT)"
	@echo
	@echo "MOCK_SOURCES (abs):"; if [ -n "$(MOCK_SOURCES)" ]; then printf "  %s\n" $(MOCK_SOURCES); else echo "  (empty)"; fi
	@if [ "$(strip $(NAMESPACE))" != "1" ]; then \
		echo; \
		echo "FLAT_TARGETS:"; \
		if [ -n "$(FLAT_TARGETS)" ]; then printf "  %s\n" $(FLAT_TARGETS); else echo "  (empty)"; fi; \
	else \
		echo; \
		echo "REL_NO_INTERNAL:"; \
		if [ -n "$(REL_NO_INTERNAL)" ]; then printf "  %s\n" $(REL_NO_INTERNAL); else echo "  (empty)"; fi; \
		echo "NAMESPACED_TARGETS:"; \
		if [ -n "$(NAMESPACED_TARGETS)" ]; then printf "  %s\n" $(NAMESPACED_TARGETS); else echo "  (empty)"; fi; \
	fi

# --------------------------------------------------------------------
# Duplicate-basename guard for FLAT mode
# --------------------------------------------------------------------
define DUP_CHECK_SCRIPT
b=($(BASENAMES)); \
declare -A seen; dup=0; \
for f in "$${b[@]}"; do \
  if [[ -n "$${seen[$$f]}" ]]; then dup=1; echo "❌ Duplicate basename detected: $$f"; fi; \
  seen[$$f]=1; \
done; \
if [[ $$dup -eq 1 ]]; then \
  echo ""; \
  echo "⚠️  Collision in FLAT mode. Options:"; \
  echo "   - Rename conflicting files/interfaces for unique basenames; or"; \
  echo "   - Run with namespacing: make mocks NAMESPACE=1"; \
  exit 2; \
fi
endef

# --------------------------------------------------------------------
# Main target
# --------------------------------------------------------------------
mocks: verify_mockgen
	@if [ "$(strip $(NAMESPACE))" != "1" ]; then \
		bash -c '$(DUP_CHECK_SCRIPT)'; \
	fi
	@$(MAKE) --no-print-directory $(GENERATED_MOCKS)
	@if [ -z "$(GENERATED_MOCKS)" ]; then \
		echo "ℹ️  No eligible files found under ports/output."; \
	else \
		echo "✅  All mocks generated successfully at: $(MOCKS_DIR)"; \
	fi

# --------------------------------------------------------------------
# Per-file generation rules (expanded dynamically)
#   - FLAT:    tests/mocks/<basename>_mock.go
#   - NS mode: tests/mocks/<path__with__double_underscores>_mock.go
# --------------------------------------------------------------------
define GEN_RULE
$(1): $(2)
	mkdir -p "$(MOCKS_DIR)"
	mockgen -source="$$(realpath $(2))" -destination="$(1)" -package="$(MOCKS_PKG)"
endef

# Expand rules for FLAT targets
ifneq ($(strip $(NAMESPACE)),1)
$(foreach i,$(MOCK_SOURCES),\
  $(eval $(call GEN_RULE,$(MOCKS_DIR)/$(patsubst %.go,%_mock.go,$(notdir $(i))),$(i))) \
)
else
# Expand rules for NAMESPACED targets
$(foreach i,$(MOCK_SOURCES),\
  $(eval $(call GEN_RULE,$(MOCKS_DIR)/$(patsubst %.go,%_mock.go,$(subst /,__,$(patsubst $(ROOT_DIR)/internal/%,%,$(patsubst $(ROOT_DIR)/%,%,$(i))))),$(i))) \
)
endif

# --------------------------------------------------------------------
# Cleanup
# --------------------------------------------------------------------
clean_mocks:
	@if [ -d "$(MOCKS_DIR)" ]; then \
		rm -rf "$(MOCKS_DIR)"; \
		echo "🧹  Cleaned: $(MOCKS_DIR)"; \
	else \
		echo "ℹ️  Nothing to clean at $(MOCKS_DIR)"; \
	fi
