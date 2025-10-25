# Production Dockerfile - Multi-stage build for minimal final image
# Target: <25MB image, <50MB RAM usage

# Build stage
FROM golang:1.23-alpine AS builder

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
# - CGO_ENABLED=0: Static binary (no C dependencies)
# - -ldflags="-w -s": Strip debug information (-w) and symbol table (-s)
# - -trimpath: Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /app/bin/keepercheky \
    ./cmd/server

# Verify binary
RUN /app/bin/keepercheky --version 2>&1 || true

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

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/keepercheky", "healthcheck"]

# Set environment
ENV KEEPERCHEKY_APP_ENVIRONMENT=production \
    KEEPERCHEKY_SERVER_PORT=8000 \
    KEEPERCHEKY_SERVER_HOST=0.0.0.0

# Run
ENTRYPOINT ["/keepercheky"]
