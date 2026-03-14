#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

fail=0

is_ignored_link() {
  local link="$1"
  [[ "$link" =~ ^https?:// ]] && return 0
  [[ "$link" =~ ^mailto: ]] && return 0
  [[ "$link" =~ ^# ]] && return 0
  return 1
}

normalize_link() {
  local raw="$1"
  raw="${raw%%#*}"
  raw="${raw%%\?*}"
  raw="${raw%\"}"
  raw="${raw#\"}"
  raw="${raw%\'}"
  raw="${raw#\'}"
  raw="${raw%% *}"
  printf '%s' "$raw"
}

check_file() {
  local file="$1"
  local dir
  dir="$(dirname "$file")"

  while IFS= read -r token; do
    local link
    link="${token#*\(}"
    link="${link%\)}"
    link="$(normalize_link "$link")"

    [[ -z "$link" ]] && continue
    is_ignored_link "$link" && continue

    local target
    if [[ "$link" = /* ]]; then
      target="$ROOT_DIR$link"
    else
      target="$dir/$link"
    fi

    local resolved
    if ! resolved="$(realpath -m "$target" 2>/dev/null)"; then
      echo "BROKEN: $file -> $link"
      fail=1
      continue
    fi

    if [[ "$resolved" != "$ROOT_DIR" && "$resolved" != "$ROOT_DIR/"* ]]; then
      echo "OUTSIDE-REPO: $file -> $link"
      fail=1
      continue
    fi

    if [[ ! -e "$resolved" ]]; then
      echo "BROKEN: $file -> $link"
      fail=1
    fi
  done < <(grep -oE '!{0,1}\[[^][]*\]\([^)]*\)' "$file" || true)
}

while IFS= read -r md; do
  check_file "$md"
done < <(git ls-files '*.md')

if [[ "$fail" -ne 0 ]]; then
  echo "Markdown link check failed."
  exit 1
fi

echo "Markdown link check passed."
