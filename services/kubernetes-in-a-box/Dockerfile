# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Generate go.sum and download dependencies
RUN go mod tidy && go mod download

# Copy source code
COPY . .

# Build the CLI
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kib ./cmd/kib

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install required tools
RUN apk --no-cache add \
    ca-certificates \
    curl \
    openssh-client \
    bash \
    jq

# Install k3sup
RUN curl -sLS https://get.k3sup.dev | sh && \
    mv k3sup /usr/local/bin/

# Install kubectl
RUN curl -LO "https://dl.k8s/release/$(curl -L -s https://dl.k8s/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Copy binary from builder
COPY --from=builder /app/kib /usr/local/bin/

# Create non-root user
RUN addgroup -g 1000 -S kib && \
    adduser -u 1000 -S kib -G kib

USER kib

# Create config directory
RUN mkdir -p /home/kib/.kib

ENTRYPOINT ["kib"]
CMD ["--help"]
