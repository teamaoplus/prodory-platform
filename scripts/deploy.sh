#!/bin/bash
# ============================================
# Prodory Platform - Deployment Script
# ============================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl not found. Please install kubectl."
        exit 1
    fi
    
    # Check Docker/Podman
    if command -v docker &> /dev/null; then
        CONTAINER_RUNTIME="docker"
    elif command -v podman &> /dev/null; then
        CONTAINER_RUNTIME="podman"
    else
        log_error "Neither Docker nor Podman found. Please install one of them."
        exit 1
    fi
    
    log_success "Using container runtime: $CONTAINER_RUNTIME"
}

# Build images
build_images() {
    log_info "Building container images..."
    
    # Build Data FinOps Agent
    log_info "Building data-finops-agent..."
    $CONTAINER_RUNTIME build -t prodory/data-finops-agent:latest \
        -f services/data-finops-agent/Dockerfile \
        services/data-finops-agent/
    
    # Build FinOps Dashboard
    log_info "Building finops-dashboard..."
    $CONTAINER_RUNTIME build -t prodory/finops-dashboard:latest \
        -f services/finops-dashboard/Dockerfile \
        services/finops-dashboard/
    
    log_success "Images built successfully"
}

# Deploy to Kubernetes
deploy_kubernetes() {
    log_info "Deploying to Kubernetes..."
    
    # Create namespace
    kubectl apply -f kubernetes/namespace.yaml
    
    # Apply configmaps and secrets
    kubectl apply -f kubernetes/configmap.yaml
    kubectl apply -f kubernetes/secret.yaml
    
    # Deploy databases
    kubectl apply -f kubernetes/postgres.yaml
    kubectl apply -f kubernetes/redis.yaml
    
    # Wait for databases
    log_info "Waiting for databases to be ready..."
    kubectl wait --for=condition=ready pod -l app=postgres -n prodory --timeout=120s
    kubectl wait --for=condition=ready pod -l app=redis -n prodory --timeout=60s
    
    # Deploy applications
    kubectl apply -f kubernetes/data-finops-agent.yaml
    kubectl apply -f kubernetes/finops-dashboard.yaml
    
    # Apply RBAC
    kubectl apply -f kubernetes/rbac.yaml
    
    # Apply ingress (optional)
    if kubectl get namespace ingress-nginx &> /dev/null; then
        log_info "Applying ingress configuration..."
        kubectl apply -f kubernetes/ingress.yaml
    else
        log_warning "NGINX Ingress Controller not found. Skipping ingress deployment."
    fi
    
    log_success "Deployment complete!"
}

# Deploy with Podman Compose
deploy_podman() {
    log_info "Deploying with Podman Compose..."
    
    # Check if podman-compose is installed
    if ! command -v podman-compose &> /dev/null; then
        log_error "podman-compose not found. Please install it."
        exit 1
    fi
    
    # Start services
    podman-compose up -d
    
    log_success "Deployment complete!"
}

# Deploy with Docker Compose
deploy_docker() {
    log_info "Deploying with Docker Compose..."
    
    # Start services
    docker-compose up -d
    
    log_success "Deployment complete!"
}

# Check deployment status
check_status() {
    log_info "Checking deployment status..."
    
    echo ""
    echo "=== Pods ==="
    kubectl get pods -n prodory
    
    echo ""
    echo "=== Services ==="
    kubectl get svc -n prodory
    
    echo ""
    echo "=== Ingress ==="
    kubectl get ingress -n prodory 2>/dev/null || echo "No ingress configured"
}

# Main
deployment_type=${1:-"kubernetes"}

case $deployment_type in
    "kubernetes"|"k8s")
        check_prerequisites
        build_images
        deploy_kubernetes
        check_status
        ;;
    "podman")
        check_prerequisites
        build_images
        deploy_podman
        ;;
    "docker")
        check_prerequisites
        build_images
        deploy_docker
        ;;
    *)
        echo "Usage: $0 [kubernetes|podman|docker]"
        echo ""
        echo "Options:"
        echo "  kubernetes  Deploy to Kubernetes cluster (default)"
        echo "  podman      Deploy using Podman Compose"
        echo "  docker      Deploy using Docker Compose"
        exit 1
        ;;
esac

log_success "Prodory Platform deployed successfully!"
echo ""
echo "Access your services:"
echo "  - Dashboard: http://localhost:3000 (or your ingress URL)"
echo "  - API: http://localhost:8000/docs"
echo ""
echo "For help, see: docs/ADMIN_GUIDE.md"
