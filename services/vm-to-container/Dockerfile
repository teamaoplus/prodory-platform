FROM node:18-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache \
    docker-cli \
    qemu-img \
    ca-certificates

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy application
COPY src/ ./src/

# Create non-root user
RUN addgroup -g 1000 vm2c && \
    adduser -u 1000 -G vm2c -s /bin/sh -D vm2c

# Make CLI executable
RUN chmod +x src/cli.js

USER vm2c

ENTRYPOINT ["node", "/app/src/cli.js"]
CMD ["--help"]
