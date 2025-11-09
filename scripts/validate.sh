#!/bin/bash
# scripts/validate.sh - Validación completa pre-commit para KeeperCheky
# Este script ejecuta todas las verificaciones necesarias antes de hacer commit

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Unicode symbols
CHECK="✅"
CROSS="❌"
WARN="⚠️"
INFO="ℹ️"
ROCKET="🚀"
TOOLS="🔧"

# Error counter
ERRORS=0

echo -e "${BLUE}${ROCKET} Running KeeperCheky validation suite...${NC}\n"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

# Function to print section header
print_section() {
    echo -e "${BOLD}${YELLOW}$1${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}${CHECK} $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}${CROSS} $1${NC}"
    ERRORS=$((ERRORS + 1))
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}${WARN} $1${NC}"
}

# Function to print info
print_info() {
    echo -e "${CYAN}${INFO} $1${NC}"
}

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 1. GO FORMAT CHECK
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "📝 Checking Go code format..."

UNFORMATTED=$(gofmt -s -l . 2>&1 | grep -v '^vendor/' | grep '.go$' || true)
if [ -n "$UNFORMATTED" ]; then
    print_error "The following files need formatting:"
    echo "$UNFORMATTED" | while read -r file; do
        echo "    - $file"
    done
    echo ""
    print_info "💡 Run 'make lint-fix' or 'gofmt -s -w .' to fix"
    echo ""
else
    print_success "All files properly formatted"
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 2. GO VET
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "🔍 Running go vet..."

if go vet ./... 2>&1; then
    print_success "Go vet passed"
else
    print_error "Go vet found issues"
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 3. TESTS WITH COVERAGE
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "🧪 Running tests with coverage and race detector..."

if go test -v -race -coverprofile=coverage.out ./... 2>&1; then
    # Calculate coverage
    if [ -f coverage.out ]; then
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        print_success "Tests passed - Coverage: ${COVERAGE}"
        
        # Warn if coverage is low
        COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
        if (( $(echo "$COVERAGE_NUM < 50" | bc -l) )); then
            print_warning "Coverage is below 50% - consider adding more tests"
        fi
    else
        print_success "Tests passed"
    fi
else
    print_error "Tests failed"
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 4. BUILD CHECK
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "🔨 Checking build..."

if go build -o /tmp/keepercheky-test ./cmd/server 2>&1; then
    rm -f /tmp/keepercheky-test
    print_success "Build successful"
else
    print_error "Build failed"
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 5. GO MOD TIDY CHECK
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "📦 Checking go.mod and go.sum..."

# Create backup
cp go.mod go.mod.backup
cp go.sum go.sum.backup

# Run go mod tidy
go mod tidy 2>&1 > /dev/null

# Check for differences
if ! diff -q go.mod go.mod.backup > /dev/null 2>&1 || ! diff -q go.sum go.sum.backup > /dev/null 2>&1; then
    print_error "go.mod or go.sum needs tidying"
    print_info "💡 Run 'go mod tidy' and commit changes"
    echo ""
    print_info "Differences found:"
    diff -u go.mod.backup go.mod || true
    
    # Restore backup
    mv go.mod.backup go.mod
    mv go.sum.backup go.sum
else
    print_success "Dependencies are clean"
    rm go.mod.backup go.sum.backup
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 6. GOLANGCI-LINT (OPTIONAL)
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
print_section "🔎 Running golangci-lint (if installed)..."

if command -v golangci-lint &> /dev/null; then
    if golangci-lint run ./... 2>&1; then
        print_success "golangci-lint passed"
    else
        print_error "golangci-lint found issues"
    fi
else
    print_warning "golangci-lint not installed, skipping..."
    print_info "💡 Install: brew install golangci-lint (macOS) or see https://golangci-lint.run/usage/install/"
fi
echo ""

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# SUMMARY
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}${BOLD}"
    echo "╔════════════════════════════════════════════╗"
    echo "║  ✅ ALL VALIDATION CHECKS PASSED!         ║"
    echo "╚════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo -e "${BLUE}👍 Ready to commit and push${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}${BOLD}"
    echo "╔════════════════════════════════════════════╗"
    echo "║  ❌ VALIDATION FAILED ($ERRORS error(s))          ║"
    echo "╚════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo -e "${YELLOW}Please fix the issues above before committing.${NC}"
    echo ""
    echo -e "${CYAN}Quick fixes:${NC}"
    echo "  • Format issues:  make lint-fix"
    echo "  • Auto-fix all:   make check-and-fix"
    echo "  • Run tests:      make test"
    echo ""
    exit 1
fi
