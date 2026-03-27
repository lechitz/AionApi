#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DASHBOARD_DIR="${ROOT_DIR}/../aion-web"

check_path() {
  local path="$1"
  if [[ ! -d "${path}" ]]; then
    echo "missing required path: ${path}" >&2
    exit 1
  fi
}

check_http() {
  local name="$1"
  local url="$2"
  if ! curl -fsS -m 5 "${url}" >/dev/null; then
    echo "required service unavailable: ${name} (${url})" >&2
    exit 1
  fi
}

check_cmd() {
  local cmd="$1"
  if ! command -v "${cmd}" >/dev/null 2>&1; then
    echo "required command not found: ${cmd}" >&2
    exit 1
  fi
}

check_path "${ROOT_DIR}"
check_path "${DASHBOARD_DIR}"
check_cmd curl
check_cmd npm

check_http "AionApi" "http://localhost:5001/aion/health"
check_http "aion-ingest" "http://localhost:8091/health"
check_http "aion-streams-admin" "http://localhost:8092/health"

echo "event backbone gate preflight passed"
