#!/bin/bash
# scripts/pre-commit.sh - Git pre-commit hook
# Automatically runs validation before each commit
#
# Installation: make install-hooks
# Skip hook:    git commit --no-verify

echo "🔍 Running pre-commit validation..."
echo ""

# Run validation script
if ! make validate-quick; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "❌ Pre-commit validation failed. Commit aborted."
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Options:"
    echo "  1. Fix the issues and try again"
    echo "  2. Run 'make check-and-fix' to auto-fix"
    echo "  3. Skip hook with: git commit --no-verify (NOT RECOMMENDED)"
    echo ""
    exit 1
fi

echo ""
echo "✅ Pre-commit validation passed. Proceeding with commit..."
echo ""
