.PHONY: help dev dev-watch build test clean docker-build docker-run shell logs stop

# Default target
help:
	@echo "KeeperCheky - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development server with hot-reload"
	@echo "  make dev-watch    - Start with Docker Compose Watch (auto-rebuild)"
	@echo "  make logs         - Show development logs"
	@echo "  make shell        - Open shell in development container"
	@echo "  make stop         - Stop development server"
	@echo "  make clean-media  - Clean and recreate mock media library"
	@echo ""
	@echo "Build:"
	@echo "  make build        - Build production binary"
	@echo "  make docker-build - Build production Docker image (default: production target)"
	@echo "  make docker-build-dev - Build development Docker image (with hot-reload)"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo ""
	@echo "Validation (run before commit):"
	@echo "  make validate      - 🔍 Full validation (format, vet, test, lint)"
	@echo "  make validate-quick - ⚡ Quick validation (format, vet, test)"
	@echo "  make check-and-fix - 🔧 Auto-fix + validate (Go + Markdown)"
	@echo "  make lint-check    - Check code format"
	@echo "  make lint-fix      - Fix code format"
	@echo "  make markdown-fix  - 📝 Auto-fix Markdown formatting"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter (golangci-lint)"
	@echo "  make install-hooks - Install pre-commit git hook"
	@echo ""

# Development with hot-reload (Air + Docker Compose Watch)
dev:
	@echo "🚀 Starting development server with hot-reload..."
	@echo "📁 Creating volume directories..."
	@# Note: keepercheky-go-modules is managed by Docker, don't create it in host
	@mkdir -p volumes/radarr-config
	@mkdir -p volumes/sonarr-config
	@mkdir -p volumes/jellyfin-config
	@mkdir -p volumes/jellyseerr-config
	@mkdir -p volumes/qbittorrent-config
	@mkdir -p volumes/bazarr-config
	@mkdir -p volumes/jellystat-config
	@mkdir -p volumes/media-library/downloads
	@mkdir -p volumes/media-library/library
	@mkdir -p logs
	@echo "✅ Volume directories ready"
	@echo ""
	@echo "💡 Tip: Run './scripts/create-mock-media.sh' to create test media files"
	@echo "📝 Logs: logs/keepercheky-dev.log (auto-rotates at 1000 lines)"
	@echo ""
	@chmod +x scripts/log-with-rotation.sh
	@docker compose up --build --watch > /dev/null 2>&1 &
	@sleep 5
	@docker compose logs -f keepercheky 2>&1 | ./scripts/log-with-rotation.sh

# Development with Docker Compose Watch (Docker 28+)
dev-watch:
	@echo "🚀 Starting development server with Docker Compose Watch..."
	@docker compose watch

# Show development logs
logs:
	@docker compose logs -f keepercheky

# Open shell in development container
shell:
	@docker compose exec keepercheky sh

# Stop development server
stop:
	@docker compose down

# Stop and remove volumes
stop-clean:
	@echo "🧹 Stopping and cleaning volumes..."
	@docker compose down -v
	@echo "✅ Containers and volumes removed"

# Clean mock media library
clean-media:
	@echo "🧹 Cleaning mock media library..."
	@rm -rf volumes/media-library/downloads
	@rm -rf volumes/media-library/library
	@echo "✅ Media library cleaned"
	@echo "   Run './scripts/create-mock-media.sh' or 'make dev' to recreate it"

# Build production binary
build:
	@echo "🔨 Building production binary..."
	@CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/keepercheky ./cmd/server

# Build production Docker image
docker-build:
	@echo "🐳 Building production Docker image..."
	@docker build -t keepercheky:latest .
	@echo "✅ Production image built successfully"
	@echo "   - Uses multi-stage build with 'production' target (default)"
	@echo "   - Final image based on scratch (~25MB)"
	@echo "   - To run: make docker-run"

# Build development Docker image (for testing)
docker-build-dev:
	@echo "🐳 Building development Docker image..."
	@docker build --target=development -t keepercheky:dev .
	@echo "✅ Development image built successfully"
	@echo "   - Uses 'development' target with hot-reload"
	@echo "   - Based on golang:alpine with Air installed"

# Run production Docker image
docker-run:
	@echo "🚀 Running production Docker image..."
	@docker run -p 8000:8000 \
		-v $(PWD)/data:/data \
		-v $(PWD)/config:/config \
		keepercheky:latest

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "📈 Total coverage: " $$3}'

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...

# Run linter (golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
# 🔍 VALIDATION TARGETS - Run before committing
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# Validate all code (format, vet, test, lint) - RUN BEFORE COMMIT
validate: lint-fix markdown-fix
	@bash scripts/validate.sh

# Quick validation (format + vet + test) - Fast pre-commit check
validate-quick: lint-check vet test
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✅ Quick validation passed!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check code format (without modifying files)
lint-check:
	@echo "📝 Checking Go code format..."
	@OUTPUT=$$(gofmt -s -l .); \
	if [ -n "$$OUTPUT" ]; then \
		echo "❌ The following files need formatting:"; \
		echo "$$OUTPUT"; \
		echo ""; \
		echo "💡 Run 'make lint-fix' to fix automatically"; \
		exit 1; \
	fi
	@echo "✅ All files are properly formatted"

# Fix code format automatically
lint-fix:
	@echo "🔧 Fixing code format..."
	@gofmt -s -w .
	@echo "✅ Format applied successfully"

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Go vet passed"

# Check and fix common issues, then validate
check-and-fix: lint-fix mod-tidy markdown-fix validate-quick
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✅ All fixes applied and validated!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "👍 Ready to commit"

# Fix Markdown formatting issues
markdown-fix:
	@echo "📝 Fixing Markdown formatting..."
	@./scripts/fix-markdown.sh

# Tidy go modules
mod-tidy:
	@echo "📦 Tidying go modules..."
	@go mod tidy
	@echo "✅ Dependencies cleaned"

# Install git hooks (optional)
install-hooks:
	@echo "📎 Installing git pre-commit hook..."
	@mkdir -p .git/hooks
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Pre-commit hook installed"
	@echo "   Will run 'make validate' before each commit"
	@echo "   To skip: git commit --no-verify"

# Uninstall git hooks
uninstall-hooks:
	@echo "🗑️  Removing git pre-commit hook..."
	@rm -f .git/hooks/pre-commit
	@echo "✅ Hook removed"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ tmp/ coverage.out coverage.html
	@echo "⚠️  Note: Docker volumes in ./volumes/ are NOT deleted"
	@echo "   Run 'make clean-all' to also remove volume data"
	@echo "✅ Clean complete"

# Clean everything including volumes
clean-all:
	@echo "🧹 Cleaning everything (including volumes)..."
	@rm -rf bin/ tmp/ coverage.out coverage.html
	@docker compose down -v
	@rm -rf volumes/
	@echo "✅ Complete cleanup done"

# Initialize development environment
init:
	@echo "🔧 Initializing development environment..."
	@mkdir -p data config
	@mkdir -p volumes/keepercheky-go-modules
	@mkdir -p volumes/radarr-config
	@mkdir -p volumes/sonarr-config
	@mkdir -p volumes/jellyfin-config
	@mkdir -p volumes/jellyseerr-config
	@mkdir -p volumes/qbittorrent-config
	@mkdir -p volumes/bazarr-config
	@mkdir -p volumes/jellystat-config
	@mkdir -p volumes/media-library/library/movies
	@mkdir -p volumes/media-library/library/tv
	@mkdir -p volumes/media-library/downloads
	@echo "✅ Development environment initialized"
	@echo "🎬 Creating mock media library..."
	@./scripts/create-mock-media.sh
	@echo ""
	@echo "📁 Directory structure:"
	@echo "  ├── data/              (app data & database)"
	@echo "  ├── config/            (configuration files)"
	@echo "  └── volumes/           (Docker volume mounts)"
	@echo "      ├── keepercheky-go-modules/"
	@echo "      ├── radarr-config/"
	@echo "      ├── sonarr-config/"
	@echo "      ├── jellyfin-config/"
	@echo "      ├── jellyseerr-config/"
	@echo "      ├── qbittorrent-config/"
	@echo "      ├── bazarr-config/"
	@echo "      ├── jellystat-config/"
	@echo "      └── media-library/"
	@echo "          ├── library/"
	@echo "          │   ├── movies/"
	@echo "          │   └── tv/"
	@echo "          └── downloads/"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Run 'make dev' to start the development server"
	@echo "  2. Visit http://localhost:8000"
	@echo "  3. Check the documentation in docs/DEVELOPMENT.md"
