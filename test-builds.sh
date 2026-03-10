#!/bin/bash
# Test all Docker builds locally

set -e

echo "=========================================="
echo "Testing Docker Builds Locally"
echo "=========================================="
echo ""

# Function to test a build
test_build() {
    local service=$1
    local context=$2
    local dockerfile=$3
    
    echo "------------------------------------------"
    echo "Building: $service"
    echo "------------------------------------------"
    
    if docker build -t "test-$service:latest" -f "$dockerfile" "$context"; then
        echo "✅ $service: BUILD SUCCESS"
        docker images "test-$service" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
        echo ""
        return 0
    else
        echo "❌ $service: BUILD FAILED"
        echo ""
        return 1
    fi
}

# Track failures
FAILED=0

# Test each service
test_build "finops-dashboard" "./services/finops-dashboard" "./services/finops-dashboard/Dockerfile" || FAILED=$((FAILED + 1))
test_build "data-finops-agent" "./services/data-finops-agent" "./services/data-finops-agent/Dockerfile" || FAILED=$((FAILED + 1))
test_build "k8s-in-a-box" "./services/kubernetes-in-a-box" "./services/kubernetes-in-a-box/Dockerfile" || FAILED=$((FAILED + 1))
test_build "storage-autoscaler" "./services/storage-autoscaler" "./services/storage-autoscaler/Dockerfile" || FAILED=$((FAILED + 1))
test_build "cloud-sentinel" "./services/cloud-sentinel" "./services/cloud-sentinel/Dockerfile" || FAILED=$((FAILED + 1))
test_build "vmware-migration" "./services/vmware-migration" "./services/vmware-migration/Dockerfile" || FAILED=$((FAILED + 1))
test_build "vm-to-container" "./services/vm-to-container" "./services/vm-to-container/Dockerfile" || FAILED=$((FAILED + 1))

echo "=========================================="
echo "Build Test Summary"
echo "=========================================="
if [ $FAILED -eq 0 ]; then
    echo "✅ All builds successful!"
    echo ""
    echo "Built images:"
    docker images | grep "^test-"
else
    echo "❌ $FAILED build(s) failed"
    exit 1
fi
