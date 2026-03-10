# Local Docker Build Testing Guide

This guide explains how to test Docker builds locally before pushing to GitHub.

## Quick Start

### Test a Single Service

```bash
# Navigate to project directory
cd prodory-platform

# Test one service
./test-single.sh finops-dashboard
```

### Test All Services

```bash
# Test all services
./test-builds.sh
```

## Manual Testing Commands

### 1. Test FinOps Dashboard (React)

```bash
# Build
docker build -t test-finops-dashboard:latest \
  -f services/finops-dashboard/Dockerfile \
  services/finops-dashboard

# Run and test
docker run --rm -p 3000:80 test-finops-dashboard:latest

# Open browser: http://localhost:3000

# Check health
curl http://localhost:3000/health
```

### 2. Test Data FinOps Agent (Python)

```bash
# Build
docker build -t test-data-finops-agent:latest \
  -f services/data-finops-agent/Dockerfile \
  services/data-finops-agent

# Run and test
docker run --rm -p 8000:8000 \
  -e DATABASE_URL="postgresql://test:test@host.docker.internal:5432/test" \
  test-data-finops-agent:latest

# Check health
curl http://localhost:8000/health
```

### 3. Test K8s-in-a-Box (Go)

```bash
# Build
docker build -t test-k8s-in-a-box:latest \
  -f services/kubernetes-in-a-box/Dockerfile \
  services/kubernetes-in-a-box

# Run and test
docker run --rm test-k8s-in-a-box:latest --help

# Test version
docker run --rm test-k8s-in-a-box:latest version
```

### 4. Test Storage Autoscaler (Python)

```bash
# Build
docker build -t test-storage-autoscaler:latest \
  -f services/storage-autoscaler/Dockerfile \
  services/storage-autoscaler

# Run and test
docker run --rm -p 8084:8084 test-storage-autoscaler:latest

# Check health
curl http://localhost:8084/healthz
```

### 5. Test Cloud Sentinel (Go)

```bash
# Build
docker build -t test-cloud-sentinel:latest \
  -f services/cloud-sentinel/Dockerfile \
  services/cloud-sentinel

# Run and test
docker run --rm -p 8083:8083 test-cloud-sentinel:latest

# Check health
curl http://localhost:8083/health
```

### 6. Test VMware Migration (Go)

```bash
# Build
docker build -t test-vmware-migration:latest \
  -f services/vmware-migration/Dockerfile \
  services/vmware-migration

# Run and test
docker run --rm test-vmware-migration:latest --help

# Test version
docker run --rm test-vmware-migration:latest version
```

### 7. Test VM to Container (Node.js)

```bash
# Build
docker build -t test-vm-to-container:latest \
  -f services/vm-to-container/Dockerfile \
  services/vm-to-container

# Run and test
docker run --rm test-vm-to-container:latest --help
```

## Using Podman Instead of Docker

If you're using Podman:

```bash
# Replace 'docker' with 'podman' in all commands
podman build -t test-finops-dashboard:latest \
  -f services/finops-dashboard/Dockerfile \
  services/finops-dashboard

podman run --rm -p 3000:80 test-finops-dashboard:latest
```

## Multi-Architecture Builds (Optional)

Test multi-arch builds locally (requires Docker Buildx):

```bash
# Create buildx builder
docker buildx create --name multiarch --use
docker buildx inspect --bootstrap

# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t test-finops-dashboard:latest \
  -f services/finops-dashboard/Dockerfile \
  services/finops-dashboard
```

## Debugging Failed Builds

### View Build Logs

```bash
# Build with detailed output
docker build --progress=plain \
  -t test-service:latest \
  -f services/service-name/Dockerfile \
  services/service-name 2>&1 | tee build.log
```

### Interactive Debugging

```bash
# Build up to a specific stage
docker build --target builder \
  -t test-debug:latest \
  -f services/service-name/Dockerfile \
  services/service-name

# Shell into the build container
docker run --rm -it test-debug:latest /bin/sh
```

### Common Issues

#### Issue: "Cannot connect to Docker daemon"

```bash
# Start Docker service
sudo systemctl start docker

# Or use Podman (daemonless)
podman build ...
```

#### Issue: "No space left on device"

```bash
# Clean up Docker cache
docker system prune -a

# Clean build cache
docker builder prune -f
```

#### Issue: Slow builds

```bash
# Use BuildKit for faster builds
export DOCKER_BUILDKIT=1

# Enable layer caching
docker build --cache-from test-service:latest ...
```

## Pre-Push Checklist

Before pushing to GitHub, verify:

- [ ] All services build successfully locally
- [ ] Images are reasonably sized (< 500MB preferred)
- [ ] Health checks work correctly
- [ ] No sensitive data in images
- [ ] Multi-arch builds work (if needed)

## Automated Testing Script

```bash
#!/bin/bash
# Full test suite

set -e

echo "=== Starting Full Test Suite ==="

# Build all images
./test-builds.sh

# Test each image
echo ""
echo "=== Testing Built Images ==="

# Test finops-dashboard
docker run --rm -d --name test-dashboard -p 3000:80 test-finops-dashboard:latest
sleep 5
curl -f http://localhost:3000/health && echo "✅ Dashboard OK" || echo "❌ Dashboard FAILED"
docker stop test-dashboard

# Test data-finops-agent (requires DB, just check it starts)
docker run --rm test-data-finops-agent:latest --help && echo "✅ Agent OK" || echo "❌ Agent FAILED"

# Test k8s-in-a-box
docker run --rm test-k8s-in-a-box:latest version && echo "✅ K8s OK" || echo "❌ K8s FAILED"

# Test cloud-sentinel
docker run --rm -d --name test-sentinel -p 8083:8083 test-cloud-sentinel:latest
sleep 5
curl -f http://localhost:8083/health && echo "✅ Sentinel OK" || echo "❌ Sentinel FAILED"
docker stop test-sentinel

# Test storage-autoscaler
docker run --rm test-storage-autoscaler:latest --help && echo "✅ Storage OK" || echo "❌ Storage FAILED"

# Test vmware-migration
docker run --rm test-vmware-migration:latest version && echo "✅ VMware OK" || echo "❌ VMware FAILED"

# Test vm-to-container
docker run --rm test-vm-to-container:latest --help && echo "✅ VM2C OK" || echo "❌ VM2C FAILED"

echo ""
echo "=== Test Suite Complete ==="
```

## Clean Up Test Images

```bash
# Remove all test images
docker images | grep "^test-" | awk '{print $3}' | xargs docker rmi -f

# Or use this command
docker rmi $(docker images --filter "reference=test-*" -q) -f 2>/dev/null || true
```

## Next Steps

After local testing passes:

1. Commit changes:
   ```bash
   git add .
   git commit -m "Fix Dockerfiles - tested locally"
   ```

2. Push to GitHub:
   ```bash
   git push origin main
   ```

3. Monitor GitHub Actions builds at:
   ```
   https://github.com/YOUR_ORG/prodory-platform/actions
   ```
