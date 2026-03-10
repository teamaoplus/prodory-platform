# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Copy source code first (needed for go mod tidy)
COPY . .

# Generate go.sum and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloud-sentinel .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/cloud-sentinel .

# Copy swagger docs if they exist
COPY --from=builder /app/docs ./docs 2>/dev/null || true

# Create non-root user
RUN addgroup -g 1000 -S sentinel && \
    adduser -u 1000 -S sentinel -G sentinel

USER sentinel

# Expose port
EXPOSE 8083

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8083/health || exit 1

# Run the application
CMD ["./cloud-sentinel"]
