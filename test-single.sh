#!/bin/bash
# Test a single service build

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: ./test-single.sh <service-name>"
    echo ""
    echo "Available services:"
    ls -1 services/
    exit 1
fi

DOCKERFILE="services/$SERVICE/Dockerfile"
CONTEXT="services/$SERVICE"

if [ ! -f "$DOCKERFILE" ]; then
    echo "❌ Dockerfile not found: $DOCKERFILE"
    exit 1
fi

echo "Building: $SERVICE"
echo "Dockerfile: $DOCKERFILE"
echo "Context: $CONTEXT"
echo ""

docker build --progress=plain -t "test-$SERVICE:latest" -f "$DOCKERFILE" "$CONTEXT" 2>&1

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Build successful!"
    echo ""
    docker images "test-$SERVICE" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
    echo ""
    echo "To test run:"
    echo "  docker run --rm -p 8080:80 test-$SERVICE:latest"
else
    echo ""
    echo "❌ Build failed!"
    exit 1
fi
