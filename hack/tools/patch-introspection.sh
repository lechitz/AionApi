#!/usr/bin/env bash
# Post-process script to fix introspection after `make graphql`
# This patches the generated.go file to enable introspection for aion-chat
# 
# The issue: gqlgen generates code that checks ec.DisableIntrospection,
# but graphql.OperationContext already provides this field.
# Solution: Simply comment out or remove the introspection checks.

set -euo pipefail

GENERATED_FILE="internal/adapter/primary/graphql/generated.go"

echo "🔧 Patching generated.go for introspection support..."

if [[ ! -f "$GENERATED_FILE" ]]; then
    echo "❌ File not found: $GENERATED_FILE"
    exit 1
fi

# Backup original
cp "$GENERATED_FILE" "${GENERATED_FILE}.bak"

# Strategy: Replace the introspection checks to always allow introspection
# Change: if ec.DisableIntrospection { ... } → if false { ... }
sed -i 's/if ec\.DisableIntrospection {/if false { \/\/ Introspection always enabled for aion-chat/g' "$GENERATED_FILE"

# gqlgen may embed module path segments into generated identifier names using
# rune separators. If the module path contains a hyphen (e.g. aion-api), the
# resulting identifiers become invalid Go symbols. Normalize only that segment
# in mangled identifiers while preserving import paths.
sed -i 's/ᚋaion-apiᚋ/ᚋaionapiᚋ/g; s/ᚋaion-apiᚐ/ᚋaionapiᚐ/g; s/ᚋaion-apiᚄ/ᚋaionapiᚄ/g' "$GENERATED_FILE"

# Add warning comment at top
sed -i '2 a\
//\
// ⚠️  AUTO-PATCHED FOR INTROSPECTION:\
// This file was automatically patched by hack/tools/patch-introspection.sh\
// to enable GraphQL introspection required by aion-chat LangChain integration.\
// Introspection checks (if ec.DisableIntrospection) are disabled (if false).\
// gqlgen hyphen-safe patch also normalizes mangled identifiers for aion-api.\
// DO NOT manually edit - changes will be overwritten by make graphql + auto-patch.\
' "$GENERATED_FILE"

# Verify patch worked
if grep -q "if false { // Introspection always enabled" "$GENERATED_FILE"; then
    echo "✅ Introspection patch applied successfully!"
    rm "${GENERATED_FILE}.bak"
    exit 0
else
    echo "❌ Patch failed! Restoring backup..."
    mv "${GENERATED_FILE}.bak" "$GENERATED_FILE"
    exit 1
fi
