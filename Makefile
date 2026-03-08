# ============================================
# Prodory Platform - Makefile
# ============================================

.PHONY: help build deploy test clean

# Default target
help:
	@echo "Prodory Platform - Available Commands"
	@echo ""
	@echo "  make build          Build all container images"
	@echo "  make deploy         Deploy to Kubernetes"
	@echo "  make deploy-podman  Deploy with Podman Compose"
	@echo "  make deploy-docker  Deploy with Docker Compose"
	@echo "  make test           Run all tests"
	@echo "  make test-backend   Run backend tests only"
	@echo "  make test-frontend  Run frontend tests only"
	@echo "  make clean          Clean up resources"
	@echo "  make logs           View logs"
	@echo "  make status         Check deployment status"
	@echo "  make stop           Stop all services"

# Build all images
build:
	@echo "Building container images..."
	docker build -t prodory/data-finops-agent:latest \
		-f services/data-finops-agent/Dockerfile \
		services/data-finops-agent/
	docker build -t prodory/finops-dashboard:latest \
		-f services/finops-dashboard/Dockerfile \
		services/finops-dashboard/

# Deploy to Kubernetes
deploy:
	@echo "Deploying to Kubernetes..."
	./scripts/deploy.sh kubernetes

# Deploy with Podman Compose
deploy-podman:
	@echo "Deploying with Podman Compose..."
	podman-compose up -d

# Deploy with Docker Compose
deploy-docker:
	@echo "Deploying with Docker Compose..."
	docker-compose up -d

# Run all tests
test: test-backend test-frontend

# Run backend tests
test-backend:
	@echo "Running backend tests..."
	cd services/data-finops-agent && \
		python -m pytest tests/ -v --cov=app --cov-report=html

# Run frontend tests
test-frontend:
	@echo "Running frontend tests..."
	cd services/finops-dashboard && \
		npm test -- --coverage

# Clean up resources
clean:
	@echo "Cleaning up..."
	kubectl delete namespace prodory --ignore-not-found=true
	docker-compose down -v --remove-orphans || true
	podman-compose down -v --remove-orphans || true

# View logs
logs:
	@echo "Viewing logs..."
	kubectl logs -f -n prodory -l app=data-finops-agent

# Check status
status:
	@echo "Checking deployment status..."
	kubectl get pods -n prodory
	kubectl get svc -n prodory

# Stop all services
stop:
	@echo "Stopping all services..."
	kubectl delete -f kubernetes/ --ignore-not-found=true
	docker-compose down || true
	podman-compose down || true

# Development commands
dev-backend:
	@echo "Starting backend in development mode..."
	cd services/data-finops-agent && \
		uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

dev-frontend:
	@echo "Starting frontend in development mode..."
	cd services/finops-dashboard && \
		npm run dev

# Database migrations
migrate:
	@echo "Running database migrations..."
	cd services/data-finops-agent && \
		alembic upgrade head

migrate-create:
	@echo "Creating new migration..."
	cd services/data-finops-agent && \
		alembic revision --autogenerate -m "$(message)"

# Code quality
lint:
	@echo "Running linters..."
	cd services/data-finops-agent && \
		flake8 app/ && \
		black app/ --check && \
		isort app/ --check-only
	cd services/finops-dashboard && \
		npm run lint

format:
	@echo "Formatting code..."
	cd services/data-finops-agent && \
		black app/ && \
		isort app/
	cd services/finops-dashboard && \
		npm run format

# Backup and restore
backup:
	@echo "Creating backup..."
	./scripts/backup.sh

restore:
	@echo "Restoring from backup..."
	./scripts/restore.sh $(file)
