#!/usr/bin/env bash
# Check if introspection is properly configured after `make graphql`
# This script should be run after regenerating GraphQL code
set -euo pipefail

GENERATED_FILE="internal/adapter/primary/graphql/generated.go"

echo "🔍 Checking GraphQL introspection configuration..."

# Check if DisableIntrospection field exists in executableSchema
if ! grep -q "DisableIntrospection.*bool" "$GENERATED_FILE"; then
    echo ""
    echo "❌ ERROR: generated.go is missing DisableIntrospection field!"
    echo ""
    echo "After running 'make graphql', you must manually restore introspection support:"
    echo ""
    echo "1. Add to executableSchema struct (around line 206):"
    echo "   DisableIntrospection  bool"
    echo ""
    echo "2. Add to NewExecutableSchema function (around line 31):"
    echo "   DisableIntrospection:  false, // Introspection enabled for LLM tool discovery"
    echo ""
    echo "See project docs/scripts for introspection patch workflow."
    echo ""
    exit 1
fi

# Check if field is initialized
if ! grep -q "DisableIntrospection.*false" "$GENERATED_FILE"; then
    echo ""
    echo "⚠️  WARNING: DisableIntrospection field exists but may not be initialized correctly"
    echo ""
fi

echo "✅ Introspection configuration check passed!"
exit 0
