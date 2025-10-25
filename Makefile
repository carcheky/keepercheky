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
	@echo ""
	@echo "Build:"
	@echo "  make build        - Build production binary"
	@echo "  make docker-build - Build production Docker image"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo ""

# Development with hot-reload (Air + Docker Compose Watch)
dev:
	@echo "🚀 Starting development server with hot-reload..."
	@echo "📁 Creating volume directories..."
	@mkdir -p volumes/keepercheky-go-modules
	@mkdir -p volumes/radarr-config
	@mkdir -p volumes/sonarr-config
	@mkdir -p volumes/jellyfin-config
	@mkdir -p volumes/jellyseerr-config
	@mkdir -p volumes/qbittorrent-config
	@mkdir -p volumes/bazarr-config
	@mkdir -p volumes/jellystat-config
	@mkdir -p volumes/jellystat-db
	@mkdir -p volumes/media-library/library/movies
	@mkdir -p volumes/media-library/library/tv
	@mkdir -p volumes/media-library/downloads
	@echo "✅ Volume directories ready"
	@docker compose -f docker-compose.dev.yml up --build --watch

# Development with Docker Compose Watch (Docker 28+)
dev-watch:
	@echo "🚀 Starting development server with Docker Compose Watch..."
	@docker compose -f docker-compose.dev.yml watch

# Show development logs
logs:
	@docker compose -f docker-compose.dev.yml logs -f keepercheky

# Open shell in development container
shell:
	@docker compose -f docker-compose.dev.yml exec keepercheky sh

# Stop development server
stop:
	@docker compose -f docker-compose.dev.yml down

# Stop and remove volumes
stop-clean:
	@echo "🧹 Stopping and cleaning volumes..."
	@docker compose -f docker-compose.dev.yml down -v
	@echo "✅ Containers and volumes removed"

# Build production binary
build:
	@echo "🔨 Building production binary..."
	@CGO_ENABLED=1 go build -ldflags="-w -s" -o bin/keepercheky ./cmd/server

# Build production Docker image
docker-build:
	@echo "🐳 Building production Docker image..."
	@docker build -t keepercheky:latest -f Dockerfile .

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
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

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
	@docker compose -f docker-compose.dev.yml down -v
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
	@mkdir -p volumes/jellystat-db
	@mkdir -p volumes/media-library/library/movies
	@mkdir -p volumes/media-library/library/tv
	@mkdir -p volumes/media-library/downloads
	@echo "✅ Development environment initialized"
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
	@echo "      ├── jellystat-db/"
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
