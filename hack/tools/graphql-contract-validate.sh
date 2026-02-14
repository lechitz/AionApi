#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SCHEMA_DIR="${ROOT_DIR}/internal/adapter/primary/graphql/schema/modules"
CONTRACTS_DIR="${ROOT_DIR}/contracts/graphql"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

SCHEMA_OPS="${TMP_DIR}/schema_ops.txt"
CONTRACT_OPS="${TMP_DIR}/contract_ops.txt"

awk '
  /extend type Query/ {mode="Q"; next}
  /extend type Mutation/ {mode="M"; next}
  mode!="" && /^}/ {mode=""; next}
  mode!="" && /^[[:space:]]*[A-Za-z_][A-Za-z0-9_]*([(:])/ {
    line=$0
    sub(/^[[:space:]]*/, "", line)
    name=line
    sub(/[(:].*/, "", name)
    if (mode == "Q") {
      print "Q\t" name
    } else {
      print "M\t" name
    }
  }
' "${SCHEMA_DIR}"/*.graphqls | sort -u > "${SCHEMA_OPS}"

while IFS= read -r abs_file; do
  kind_raw="$(sed -E 's/#.*$//' "${abs_file}" | tr '\n' ' ' | sed -E 's/^[[:space:]]*([A-Za-z]+).*/\1/' | tr '[:upper:]' '[:lower:]')"
  case "${kind_raw}" in
    query) kind="Q" ;;
    mutation) kind="M" ;;
    *)
      echo "Unknown operation type in ${abs_file}" >&2
      exit 1
      ;;
  esac

  root_field="$(tr '\n' ' ' < "${abs_file}" | sed -E 's/[[:space:]]+/ /g' | sed -E 's/^[^{]*\{[[:space:]]*//' | sed -E 's/^([A-Za-z_][A-Za-z0-9_]*).*/\1/')"
  if [[ -z "${root_field}" ]]; then
    echo "Failed to parse root field in ${abs_file}" >&2
    exit 1
  fi

  printf '%s\t%s\n' "${kind}" "${root_field}" >> "${CONTRACT_OPS}"
done < <(find "${CONTRACTS_DIR}/queries" "${CONTRACTS_DIR}/mutations" -type f -name '*.graphql' | sort)

sort -u -o "${CONTRACT_OPS}" "${CONTRACT_OPS}"

MISSING_IN_CONTRACTS="${TMP_DIR}/missing_in_contracts.txt"
EXTRA_IN_CONTRACTS="${TMP_DIR}/extra_in_contracts.txt"

comm -23 "${SCHEMA_OPS}" "${CONTRACT_OPS}" > "${MISSING_IN_CONTRACTS}" || true
comm -13 "${SCHEMA_OPS}" "${CONTRACT_OPS}" > "${EXTRA_IN_CONTRACTS}" || true

if [[ -s "${MISSING_IN_CONTRACTS}" ]]; then
  echo "❌ Missing GraphQL operations in contracts:" >&2
  sed 's/^/  - /' "${MISSING_IN_CONTRACTS}" >&2
  exit 1
fi

if [[ -s "${EXTRA_IN_CONTRACTS}" ]]; then
  echo "❌ Unknown GraphQL operations in contracts (not present in schema):" >&2
  sed 's/^/  - /' "${EXTRA_IN_CONTRACTS}" >&2
  exit 1
fi

echo "✅ GraphQL contracts validated against schema"
