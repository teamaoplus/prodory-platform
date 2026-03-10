# VMware Migration Service - Dockerfile
# Multi-stage build for Go application

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Copy source code first (needed for go mod tidy)
COPY . .

# Generate go.sum and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vmware-migrate ./cmd/vmware-migrate

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S prodory && \
    adduser -u 1000 -S prodory -G prodory

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/vmware-migrate /app/

# Change ownership
RUN chown -R prodory:prodory /app

# Switch to non-root user
USER prodory

# Expose ports
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3001/health || exit 1

# Run the application
CMD ["./vmware-migrate", "server"]
