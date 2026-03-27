#!/usr/bin/env bash
set -euo pipefail

# Cross-repo local regression gate focused on draft-card flow.
# Must be executed from AionApi make target or directly from this repo.

API_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AION_ROOT="$(cd "${API_DIR}/.." && pwd)"
DASHBOARD_DIR="${AION_ROOT}/aion-web"
CHAT_DIR="${AION_ROOT}/aion-chat"
GOCACHE_DIR="${API_DIR}/.cache/go-build"

pass() { printf 'PASS: %s\n' "$1"; }
fail() { printf 'FAIL: %s\n' "$1" >&2; }

run_step() {
  local label="$1"
  shift
  printf '\n== %s ==\n' "$label"
  if "$@"; then
    pass "$label"
    return 0
  fi
  fail "$label"
  return 1
}

assert_dir() {
  local path="$1"
  local label="$2"
  if [[ ! -d "$path" ]]; then
    echo "Missing directory for ${label}: ${path}" >&2
    exit 1
  fi
}

assert_dir "$DASHBOARD_DIR" "aion-web"
assert_dir "$CHAT_DIR" "aion-chat"
mkdir -p "$GOCACHE_DIR"

run_step "Dashboard tests (DC-08/DC-09 partial automation)" \
  bash -lc "cd \"$DASHBOARD_DIR\" && npm test"

run_step "aion-chat ui_action/quick_add unit tests" \
  bash -lc "cd \"$CHAT_DIR\" && ./.venv/bin/pytest -q tests/unit/contexts/chat_interaction/adapter/secondary/test_chat_service.py -k \"ui_action or quick_add or mixed_accept_cancel_with_edits_sequence\""

run_step "aion-chat draft time-range/implicit-action regression tests" \
  bash -lc "cd \"$CHAT_DIR\" && ./.venv/bin/pytest -q tests/unit/contexts/chat_interaction/core/usecase/test_mutation_draft_flow.py"

run_step "AionApi chat HTTP handler tests" \
  bash -lc "cd \"$API_DIR\" && GOCACHE=\"$GOCACHE_DIR\" go test ./internal/chat/adapter/primary/http/handler/..."

run_step "AionApi ui_action extraction contract test" \
  bash -lc "cd \"$API_DIR\" && GOCACHE=\"$GOCACHE_DIR\" go test ./internal/chat/core/usecase -run \"TestExtractUIActionMetadata\""

printf '\nAll regression gate checks passed.\n'
