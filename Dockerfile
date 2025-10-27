# Production Dockerfile - Multi-stage build for minimal final image
# Target: <25MB image, <50MB RAM usage

# Build stage
FROM golang:1.25-alpine AS builder

# Build arguments for versioning
ARG VERSION=dev
ARG COMMIT_SHA=unknown

# Platform-specific build arguments (automatically set by buildx)
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && \
    go mod verify

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY web/ ./web/

# Build binary with optimizations
# - CGO_ENABLED=1: Required for SQLite (but using musl for static linking)
# - TARGETOS/TARGETARCH: Set by buildx for multi-platform builds
# - -ldflags="-w -s": Strip debug information (-w) and symbol table (-s)
# - -trimpath: Remove file system paths from binary
# - -X: Inject version information at build time
# - -linkmode external: Use external linker for CGO
# - -extldflags '-static': Force static linking
RUN apk add --no-cache gcc musl-dev && \
    CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -linkmode external -extldflags '-static' -X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA}" \
    -trimpath \
    -o /app/bin/keepercheky \
    ./cmd/server

# Verify binary exists and is executable
RUN ls -lh /app/bin/keepercheky && file /app/bin/keepercheky

# Final stage - Use scratch for absolute minimal image
FROM scratch

# Copy certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/bin/keepercheky /keepercheky

# Copy web assets
COPY --from=builder /app/web /web

# Create non-root user (note: in scratch, this is just metadata)
USER 65534:65534

# Expose port
EXPOSE 8000

# Health check - Disabled: requires running server with database
# Use docker-compose healthcheck with curl/wget instead
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD ["/keepercheky", "healthcheck"]

# Set environment
ENV KEEPERCHEKY_APP_ENVIRONMENT=production \
    KEEPERCHEKY_SERVER_PORT=8000 \
    KEEPERCHEKY_SERVER_HOST=0.0.0.0

# Run
ENTRYPOINT ["/keepercheky"]
