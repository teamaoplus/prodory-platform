# Prodory Platform - Comprehensive Administrator Guide

## Table of Contents
1. [Architecture Overview](#architecture-overview)
2. [Installation & Deployment](#installation--deployment)
3. [Configuration Management](#configuration-management)
4. [Manual Operations](#manual-operations)
5. [Kubernetes Deployment](#kubernetes-deployment)
6. [Podman Deployment](#podman-deployment)
7. [Database Management](#database-management)
8. [Security Configuration](#security-configuration)
9. [Monitoring & Alerting](#monitoring--alerting)
10. [Backup & Recovery](#backup--recovery)
11. [Upgrade Procedures](#upgrade-procedures)
12. [Troubleshooting](#troubleshooting)
13. [API Reference](#api-reference)
14. [Custom Modifications](#custom-modifications)

---

## Architecture Overview

### System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Prodory Platform                                │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Load Balancer (NGINX)                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────┼─────────────────────────────────┐     │
│  │                                 │                                 │     │
│  ▼                                 ▼                                 ▼     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  AI FinOps   │  │ Data FinOps  │  │  K8s-in-a-   │  │   Storage    │   │
│  │  Dashboard   │  │    Agent     │  │     Box      │  │  Autoscaler  │   │
│  │   (React)    │  │   (Python)   │  │    (Go)      │  │   (Python)   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                      │
│  │   Cloud      │  │  VMware to   │  │  VM to       │                      │
│  │   Sentinel   │  │  OpenShift   │  │  Container   │                      │
│  │    (Go)      │  │   (Node.js)  │  │   (Node.js)  │                      │
│  └──────────────┘  └──────────────┘  └──────────────┘                      │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         Shared Services                              │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐            │   │
│  │  │PostgreSQL│  │  Redis   │  │ RabbitMQ │  │   MinIO  │            │   │
│  │  │(Primary) │  │ (Cache)  │  │ (Queue)  │  │(Object S3)│           │   │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Service Details

| Service | Port | Language | Database | Purpose |
|---------|------|----------|----------|---------|
| FinOps Dashboard | 3000 | React/TypeScript | - | Web UI |
| Data FinOps Agent | 8000 | Python/FastAPI | PostgreSQL | Cost data collection |
| K8s-in-a-Box API | 8080 | Go/Gin | PostgreSQL | K8s management |
| Storage Autoscaler | 8001 | Python/FastAPI | PostgreSQL | Storage management |
| Cloud Sentinel | 8081 | Go/Gin | PostgreSQL | Security monitoring |
| VMware Migration | 3001 | Node.js/Express | PostgreSQL | VM migration |
| VM2Container | 3002 | Node.js/Express | PostgreSQL | Containerization |

### Network Diagram

```
Internet
    │
    ▼
┌─────────────┐
│   CDN/WAF   │ (Cloudflare/AWS CloudFront)
└─────────────┘
    │
    ▼
┌─────────────┐
│   Ingress   │ (NGINX/Kubernetes Ingress)
│  Controller │
└─────────────┘
    │
    ├──▶ FinOps Dashboard (3000)
    ├──▶ Data FinOps Agent (8000)
    ├──▶ K8s-in-a-Box (8080)
    ├──▶ Storage Autoscaler (8001)
    ├──▶ Cloud Sentinel (8081)
    ├──▶ VMware Migration (3001)
    └──▶ VM2Container (3002)
    
Internal Services:
    ├──▶ PostgreSQL (5432)
    ├──▶ Redis (6379)
    ├──▶ RabbitMQ (5672)
    └──▶ MinIO (9000)
```

---

## Installation & Deployment

### Prerequisites

#### Hardware Requirements

| Environment | CPU | RAM | Storage | Network |
|-------------|-----|-----|---------|---------|
| Development | 4 cores | 8 GB | 50 GB SSD | 100 Mbps |
| Testing | 8 cores | 16 GB | 100 GB SSD | 1 Gbps |
| Production | 16+ cores | 32+ GB | 500 GB SSD | 10 Gbps |
| HA Production | 32+ cores | 64+ GB | 1 TB SSD | 10 Gbps |

#### Software Requirements

```bash
# Required software versions
Docker: 24.0+
Docker Compose: 2.20+
Podman: 4.0+
Podman Compose: 1.0+
Kubernetes: 1.27+
Helm: 3.12+
kubectl: 1.27+
```

#### Operating System Support

| OS | Version | Support Level |
|----|---------|---------------|
| Ubuntu | 22.04 LTS | Full |
| RHEL | 8.x, 9.x | Full |
| CentOS Stream | 9 | Full |
| Debian | 11, 12 | Community |
| Amazon Linux | 2023 | Full |

### Environment Setup

#### Step 1: Install Docker/Podman

**Ubuntu/Debian:**
```bash
# Install Docker
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

**RHEL/CentOS:**
```bash
# Install Podman
sudo dnf -y install podman podman-compose

# Enable Podman socket
sudo systemctl enable --now podman.socket
```

#### Step 2: Install Kubernetes Tools

```bash
# Install kubectl
curl -LO "https://dl.k8s/release/$(curl -L -s https://dl.k8s/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Verify installations
kubectl version --client
helm version
```

#### Step 3: Clone Repository

```bash
# Clone the repository
git clone https://github.com/your-org/prodory-platform.git
cd prodory-platform

# Checkout specific version (optional)
git checkout v1.0.0
```

---

## Configuration Management

### Environment Variables

Create a comprehensive `.env` file:

```bash
# Copy example environment file
cp .env.example .env

# Edit with your values
nano .env
```

#### Complete Environment Configuration

```bash
# =============================================================================
# PRODORY PLATFORM - ENVIRONMENT CONFIGURATION
# =============================================================================

# -----------------------------------------------------------------------------
# GENERAL SETTINGS
# -----------------------------------------------------------------------------
ENVIRONMENT=production
DEBUG=false
LOG_LEVEL=INFO
TIMEZONE=UTC

# -----------------------------------------------------------------------------
# DATABASE CONFIGURATION
# -----------------------------------------------------------------------------
# Primary PostgreSQL Database
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_NAME=prodory
DATABASE_USER=prodory
DATABASE_PASSWORD=your-secure-password-here
DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}

# Connection Pool Settings
DATABASE_POOL_SIZE=20
DATABASE_MAX_OVERFLOW=30
DATABASE_POOL_TIMEOUT=30

# Read Replica (Optional)
DATABASE_REPLICA_HOST=
DATABASE_REPLICA_URL=

# -----------------------------------------------------------------------------
# REDIS CONFIGURATION
# -----------------------------------------------------------------------------
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0
REDIS_URL=redis://:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/${REDIS_DB}

# Redis Cache Settings
REDIS_CACHE_TTL=3600
REDIS_SESSION_TTL=86400

# -----------------------------------------------------------------------------
# SUPABASE CONFIGURATION (Required for Production)
# -----------------------------------------------------------------------------
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here
SUPABASE_SERVICE_KEY=your-service-role-key-here
SUPABASE_JWT_SECRET=your-jwt-secret-here

# -----------------------------------------------------------------------------
# CLOUD PROVIDER CREDENTIALS
# -----------------------------------------------------------------------------
# AWS Configuration
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_DEFAULT_REGION=us-east-1
AWS_ACCOUNT_ID=123456789012

# Azure Configuration
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret
AZURE_TENANT_ID=your-tenant-id
AZURE_SUBSCRIPTION_ID=your-subscription-id

# GCP Configuration
GOOGLE_APPLICATION_CREDENTIALS=/app/config/gcp-key.json
GCP_PROJECT_ID=your-project-id
GCP_BILLING_ACCOUNT=your-billing-account

# -----------------------------------------------------------------------------
# SECURITY SETTINGS
# -----------------------------------------------------------------------------
# JWT Configuration
JWT_SECRET=your-32-character-jwt-secret-here
JWT_ALGORITHM=HS256
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=30
JWT_REFRESH_TOKEN_EXPIRE_DAYS=7

# Encryption
ENCRYPTION_KEY=your-32-byte-encryption-key-here

# CORS Settings
CORS_ORIGINS=https://prodory.your-domain.com,https://admin.prodory.your-domain.com
CORS_ALLOW_CREDENTIALS=true

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# -----------------------------------------------------------------------------
# EMAIL CONFIGURATION
# -----------------------------------------------------------------------------
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=notifications@your-domain.com
SMTP_PASSWORD=your-app-password
SMTP_TLS=true
SMTP_FROM=Prodory Platform <notifications@your-domain.com>

# -----------------------------------------------------------------------------
# NOTIFICATION SETTINGS
# -----------------------------------------------------------------------------
# Slack Integration
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
SLACK_CHANNEL=#alerts

# PagerDuty Integration
PAGERDUTY_SERVICE_KEY=your-service-key
PAGERDUTY_API_TOKEN=your-api-token

# Microsoft Teams
TEAMS_WEBHOOK_URL=https://your-org.webhook.office.com/...

# -----------------------------------------------------------------------------
# AI/ML CONFIGURATION
# -----------------------------------------------------------------------------
# OpenAI API (for AI recommendations)
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4

# Claude API (alternative)
ANTHROPIC_API_KEY=sk-ant-...
ANTHROPIC_MODEL=claude-3-opus-20240229

# Local AI Model (optional)
LOCAL_AI_ENABLED=false
LOCAL_AI_URL=http://localhost:11434

# -----------------------------------------------------------------------------
# MONITORING & OBSERVABILITY
# -----------------------------------------------------------------------------
# Prometheus
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090

# Grafana
GRAFANA_ENABLED=true
GRAFANA_PORT=3001
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your-grafana-password

# Jaeger (Distributed Tracing)
JAEGER_ENABLED=true
JAEGER_AGENT_HOST=jaeger
JAEGER_AGENT_PORT=6831

# Sentry (Error Tracking)
SENTRY_DSN=https://...@sentry.io/...
SENTRY_ENVIRONMENT=production

# -----------------------------------------------------------------------------
# STORAGE CONFIGURATION
# -----------------------------------------------------------------------------
# MinIO/S3
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET_NAME=prodory-storage
MINIO_USE_SSL=false
MINIO_REGION=us-east-1

# Backup Storage
BACKUP_S3_BUCKET=prodory-backups
BACKUP_S3_REGION=us-east-1
BACKUP_RETENTION_DAYS=30

# -----------------------------------------------------------------------------
# KUBERNETES CONFIGURATION
# -----------------------------------------------------------------------------
KUBECONFIG_PATH=/app/config/kubeconfig
K8S_DEFAULT_NAMESPACE=prodory
K8S_INGRESS_CLASS=nginx
K8S_CERT_MANAGER_ISSUER=letsencrypt-prod

# -----------------------------------------------------------------------------
# FEATURE FLAGS
# -----------------------------------------------------------------------------
FEATURE_AI_RECOMMENDATIONS=true
FEATURE_AUTO_REMEDIATION=false
FEATURE_ADVANCED_REPORTING=true
FEATURE_COST_ALLOCATION=true
FEATURE_SECURITY_SCANNING=true

# -----------------------------------------------------------------------------
# PERFORMANCE TUNING
# -----------------------------------------------------------------------------
WORKER_PROCESSES=4
WORKER_THREADS=8
MAX_CONNECTIONS=1000
REQUEST_TIMEOUT=30

# -----------------------------------------------------------------------------
# CUSTOM SETTINGS
# -----------------------------------------------------------------------------
# Add your custom settings below
```

### Configuration Files

#### Database Configuration

```yaml
# config/database.yaml
production:
  primary:
    host: postgres
    port: 5432
    database: prodory
    user: prodory
    password: ${DATABASE_PASSWORD}
    pool:
      min: 5
      max: 20
      idle_timeout: 300
    ssl:
      enabled: true
      mode: require
      
  replicas:
    - host: postgres-replica-1
      port: 5432
      weight: 100
      
  migrations:
    directory: ./migrations
    table: schema_migrations
    
  backup:
    enabled: true
    schedule: "0 2 * * *"  # Daily at 2 AM
    retention: 30
    s3_bucket: prodory-backups
```

#### Cache Configuration

```yaml
# config/cache.yaml
redis:
  clusters:
    default:
      host: redis
      port: 6379
      password: ${REDIS_PASSWORD}
      db: 0
      
    sessions:
      host: redis
      port: 6379
      password: ${REDIS_PASSWORD}
      db: 1
      
    cache:
      host: redis
      port: 6379
      password: ${REDIS_PASSWORD}
      db: 2
      ttl: 3600
      
  sentinel:
    enabled: false
    master_name: mymaster
    nodes:
      - host: sentinel-1
        port: 26379
      - host: sentinel-2
        port: 26379
```

---

## Manual Operations

### Starting Services Manually

#### Start Individual Services

```bash
# Start PostgreSQL
docker run -d \
  --name prodory-postgres \
  -e POSTGRES_USER=prodory \
  -e POSTGRES_PASSWORD=your-password \
  -e POSTGRES_DB=prodory \
  -v postgres_data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:15-alpine

# Start Redis
docker run -d \
  --name prodory-redis \
  -e REDIS_PASSWORD=your-password \
  -v redis_data:/data \
  -p 6379:6379 \
  redis:7-alpine \
  redis-server --requirepass your-password

# Start RabbitMQ
docker run -d \
  --name prodory-rabbitmq \
  -e RABBITMQ_DEFAULT_USER=prodory \
  -e RABBITMQ_DEFAULT_PASS=your-password \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3-management
```

#### Start Application Services

```bash
# Build and start Data FinOps Agent
cd services/data-finops-agent
docker build -t prodory/data-finops-agent:latest .
docker run -d \
  --name data-finops-agent \
  --env-file ../../.env \
  -p 8000:8000 \
  prodory/data-finops-agent:latest

# Build and start FinOps Dashboard
cd services/finops-dashboard
docker build -t prodory/finops-dashboard:latest .
docker run -d \
  --name finops-dashboard \
  --env-file ../../.env \
  -p 3000:80 \
  prodory/finops-dashboard:latest
```

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop individual service
docker stop data-finops-agent
docker rm data-finops-agent

# Stop with data removal
docker-compose down -v
```

### Viewing Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f data-finops-agent

# View last 100 lines
docker-compose logs --tail=100 data-finops-agent

# View logs with timestamps
docker-compose logs -f -t data-finops-agent

# View logs since specific time
docker-compose logs -f --since=2024-01-15T10:00:00
```

### Health Checks

```bash
# Check all services
docker-compose ps

# Check individual service health
curl http://localhost:8000/health
curl http://localhost:3000/health

# Detailed health check
curl http://localhost:8000/health/detailed
```

---

## Kubernetes Deployment

### Pre-Deployment Checklist

- [ ] Kubernetes cluster version 1.27+
- [ ] kubectl configured and authenticated
- [ ] Helm 3.12+ installed
- [ ] Ingress controller installed (NGINX recommended)
- [ ] cert-manager installed (for TLS)
- [ ] Storage class configured
- [ ] Container registry accessible

### Step 1: Create Namespace

```bash
# Create namespace
kubectl create namespace prodory

# Set as default namespace
kubectl config set-context --current --namespace=prodory

# Label namespace for monitoring
kubectl label namespace prodory monitoring=enabled
```

### Step 2: Create Secrets

```bash
# Create secrets from environment file
kubectl create secret generic prodory-secrets \
  --from-env-file=.env \
  --namespace=prodory

# Create individual secrets
kubectl create secret generic database-credentials \
  --from-literal=username=prodory \
  --from-literal=password=your-password \
  --namespace=prodory

kubectl create secret generic cloud-credentials \
  --from-literal=aws-access-key=AKIA... \
  --from-literal=aws-secret-key=... \
  --namespace=prodory

# Create TLS secret
kubectl create secret tls prodory-tls \
  --cert=cert.pem \
  --key=key.pem \
  --namespace=prodory
```

### Step 3: Deploy Database

```yaml
# kubernetes/postgres.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: prodory
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: database-credentials
                  key: username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: database-credentials
                  key: password
            - name: POSTGRES_DB
              value: prodory
            - name: PGDATA
              value: /var/lib/postgresql/data/pgdata
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
          livenessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - prodory
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - prodory
            initialDelaySeconds: 5
            periodSeconds: 5
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: standard
        resources:
          requests:
            storage: 50Gi
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: prodory
spec:
  selector:
    app: postgres
  ports:
    - port: 5432
      targetPort: 5432
  type: ClusterIP
```

```bash
# Deploy PostgreSQL
kubectl apply -f kubernetes/postgres.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s

# Verify deployment
kubectl get pods -l app=postgres
kubectl get svc postgres
```

### Step 4: Deploy Redis

```yaml
# kubernetes/redis.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  namespace: prodory
spec:
  serviceName: redis
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          command:
            - redis-server
            - --requirepass
            - $(REDIS_PASSWORD)
            - --appendonly
            - "yes"
          ports:
            - containerPort: 6379
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: prodory-secrets
                  key: REDIS_PASSWORD
          volumeMounts:
            - name: redis-storage
              mountPath: /data
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "1Gi"
              cpu: "500m"
  volumeClaimTemplates:
    - metadata:
        name: redis-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: standard
        resources:
          requests:
            storage: 10Gi
---
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: prodory
spec:
  selector:
    app: redis
  ports:
    - port: 6379
      targetPort: 6379
```

### Step 5: Deploy Application Services

```yaml
# kubernetes/data-finops-agent.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: data-finops-agent
  namespace: prodory
  labels:
    app: data-finops-agent
spec:
  replicas: 2
  selector:
    matchLabels:
      app: data-finops-agent
  template:
    metadata:
      labels:
        app: data-finops-agent
    spec:
      containers:
        - name: data-finops-agent
          image: prodory/data-finops-agent:latest
          ports:
            - containerPort: 8000
          envFrom:
            - secretRef:
                name: prodory-secrets
          env:
            - name: DATABASE_HOST
              value: postgres
            - name: REDIS_HOST
              value: redis
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8000
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8000
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: data-finops-agent
  namespace: prodory
spec:
  selector:
    app: data-finops-agent
  ports:
    - port: 8000
      targetPort: 8000
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: data-finops-agent
  namespace: prodory
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: data-finops-agent
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

```bash
# Deploy all services
kubectl apply -f kubernetes/

# Verify deployments
kubectl get deployments
kubectl get pods
kubectl get services
```

### Step 6: Configure Ingress

```yaml
# kubernetes/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: prodory-ingress
  namespace: prodory
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - prodory.your-domain.com
        - api.prodory.your-domain.com
      secretName: prodory-tls
  rules:
    - host: prodory.your-domain.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: finops-dashboard
                port:
                  number: 80
    - host: api.prodory.your-domain.com
      http:
        paths:
          - path: /finops
            pathType: Prefix
            backend:
              service:
                name: data-finops-agent
                port:
                  number: 8000
          - path: /k8s
            pathType: Prefix
            backend:
              service:
                name: k8s-in-a-box
                port:
                  number: 8080
          - path: /storage
            pathType: Prefix
            backend:
              service:
                name: storage-autoscaler
                port:
                  number: 8001
          - path: /security
            pathType: Prefix
            backend:
              service:
                name: cloud-sentinel
                port:
                  number: 8081
```

```bash
# Apply ingress
kubectl apply -f kubernetes/ingress.yaml

# Verify ingress
kubectl get ingress
kubectl describe ingress prodory-ingress
```

### Step 7: Verify Deployment

```bash
# Check all resources
kubectl get all -n prodory

# Check pod status
kubectl get pods -n prodory -o wide

# Check service endpoints
kubectl get endpoints -n prodory

# Test API endpoints
curl https://api.prodory.your-domain.com/finops/health
curl https://api.prodory.your-domain.com/k8s/health

# View logs
kubectl logs -f deployment/data-finops-agent -n prodory
```

---

## Podman Deployment

### Podman-Specific Configuration

#### Install Podman Compose

```bash
# RHEL/CentOS/Fedora
sudo dnf install podman-compose

# Ubuntu/Debian
sudo apt install podman-compose

# Verify installation
podman-compose version
```

#### Create Podman Network

```bash
# Create custom network
podman network create prodory-network

# Verify network
podman network ls
podman network inspect prodory-network
```

### Podman Compose Configuration

Create `podman-compose.yml`:

```yaml
version: '3.8'

networks:
  prodory-network:
    external: true

volumes:
  postgres_data:
  redis_data:
  minio_data:

services:
  postgres:
    image: docker.io/library/postgres:15-alpine
    container_name: prodory-postgres
    networks:
      - prodory-network
    environment:
      POSTGRES_USER: prodory
      POSTGRES_PASSWORD: ${DATABASE_PASSWORD}
      POSTGRES_DB: prodory
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U prodory"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: docker.io/library/redis:7-alpine
    container_name: prodory-redis
    networks:
      - prodory-network
    command: >
      sh -c "redis-server --requirepass $${REDIS_PASSWORD} --appendonly yes"
    environment:
      REDIS_PASSWORD: ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  data-finops-agent:
    image: localhost/prodory/data-finops-agent:latest
    container_name: prodory-data-finops-agent
    networks:
      - prodory-network
    environment:
      - DATABASE_URL=postgresql://prodory:${DATABASE_PASSWORD}@postgres:5432/prodory
      - REDIS_URL=redis://:${REDIS_PASSWORD}@redis:6379/0
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
    env_file:
      - .env
    ports:
      - "8000:8000"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    # Podman-specific security options
    security_opt:
      - label=disable
    cap_drop:
      - ALL
    cap_add:
      - CHOWN
      - SETGID
      - SETUID

  finops-dashboard:
    image: localhost/prodory/finops-dashboard:latest
    container_name: prodory-finops-dashboard
    networks:
      - prodory-network
    environment:
      - REACT_APP_API_URL=http://localhost:8000
    ports:
      - "3000:80"
    depends_on:
      - data-finops-agent
    restart: unless-stopped
    security_opt:
      - label=disable
```

### Start Podman Services

```bash
# Start all services
podman-compose up -d

# View logs
podman-compose logs -f

# Check status
podman ps
podman pod ps
```

### Podman Systemd Integration

Generate systemd service files for auto-start:

```bash
# Generate systemd unit files
podman generate systemd --new --name prodory-postgres > ~/.config/systemd/user/prodory-postgres.service
podman generate systemd --new --name prodory-redis > ~/.config/systemd/user/prodory-redis.service
podman generate systemd --new --name prodory-data-finops-agent > ~/.config/systemd/user/prodory-data-finops-agent.service

# Enable and start services
systemctl --user daemon-reload
systemctl --user enable prodory-postgres
systemctl --user enable prodory-redis
systemctl --user enable prodory-data-finops-agent
systemctl --user start prodory-postgres
systemctl --user start prodory-redis
systemctl --user start prodory-data-finops-agent

# Check service status
systemctl --user status prodory-*
```

---

## Database Management

### Database Schema

#### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'viewer',
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Cloud providers table
CREATE TABLE cloud_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- aws, azure, gcp
    name VARCHAR(255) NOT NULL,
    credentials JSONB NOT NULL,
    settings JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'pending',
    last_sync TIMESTAMP,
    sync_error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Cost data table
CREATE TABLE cost_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id UUID REFERENCES cloud_providers(id) ON DELETE CASCADE,
    service VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    resource_type VARCHAR(100),
    cost DECIMAL(15,4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    usage_quantity DECIMAL(15,4),
    usage_unit VARCHAR(50),
    tags JSONB DEFAULT '{}',
    region VARCHAR(100),
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Budgets table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    period VARCHAR(50) NOT NULL, -- daily, weekly, monthly, quarterly, yearly
    start_date DATE NOT NULL,
    end_date DATE,
    alert_thresholds INTEGER[] DEFAULT ARRAY[50, 75, 90, 100],
    filters JSONB DEFAULT '{}',
    notifications JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Recommendations table
CREATE TABLE recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID REFERENCES cloud_providers(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    impact VARCHAR(20) NOT NULL, -- high, medium, low
    effort VARCHAR(20) NOT NULL, -- easy, medium, complex
    estimated_savings DECIMAL(15,2),
    actual_savings DECIMAL(15,2),
    status VARCHAR(50) DEFAULT 'pending',
    resources JSONB DEFAULT '[]',
    confidence DECIMAL(3,2),
    ai_analysis JSONB,
    applied_at TIMESTAMP,
    applied_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Audit log table
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_cost_data_timestamp ON cost_data(timestamp);
CREATE INDEX idx_cost_data_provider ON cost_data(provider_id);
CREATE INDEX idx_cost_data_service ON cost_data(service);
CREATE INDEX idx_recommendations_status ON recommendations(status);
CREATE INDEX idx_recommendations_user ON recommendations(user_id);
CREATE INDEX idx_audit_log_user ON audit_log(user_id);
CREATE INDEX idx_audit_log_created ON audit_log(created_at);
```

### Database Migrations

#### Using Alembic (Python)

```bash
# Install Alembic
pip install alembic

# Initialize migrations
alembic init migrations

# Create new migration
alembic revision -m "add user preferences"

# Edit migration file
# migrations/versions/xxx_add_user_preferences.py

# Apply migrations
alembic upgrade head

# Rollback migration
alembic downgrade -1

# View current version
alembic current

# View history
alembic history
```

#### Manual Migration Example

```python
# migrations/versions/001_add_user_preferences.py
"""Add user preferences table

Revision ID: 001
Revises: 
Create Date: 2024-01-15 10:00:00

"""
from alembic import op
import sqlalchemy as sa

# revision identifiers
revision = '001'
down_revision = None
branch_labels = None
depends_on = None

def upgrade():
    op.create_table(
        'user_preferences',
        sa.Column('id', sa.UUID(), nullable=False),
        sa.Column('user_id', sa.UUID(), nullable=False),
        sa.Column('preferences', sa.JSON(), nullable=False, default={}),
        sa.Column('created_at', sa.TIMESTAMP(), nullable=False, server_default=sa.func.now()),
        sa.Column('updated_at', sa.TIMESTAMP(), nullable=False, server_default=sa.func.now()),
        sa.ForeignKeyConstraint(['user_id'], ['users.id'], ondelete='CASCADE'),
        sa.PrimaryKeyConstraint('id'),
        sa.UniqueConstraint('user_id')
    )
    
    op.create_index('idx_user_preferences_user', 'user_preferences', ['user_id'])

def downgrade():
    op.drop_index('idx_user_preferences_user', table_name='user_preferences')
    op.drop_table('user_preferences')
```

### Database Backup & Restore

#### Automated Backup Script

```bash
#!/bin/bash
# backup-database.sh

set -e

# Configuration
DB_HOST="${DATABASE_HOST:-localhost}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-prodory}"
DB_USER="${DATABASE_USER:-prodory}"
DB_PASSWORD="${DATABASE_PASSWORD}"
BACKUP_DIR="${BACKUP_DIR:-/backups}"
S3_BUCKET="${BACKUP_S3_BUCKET:-prodory-backups}"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"

# Create backup filename
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/prodory_backup_${TIMESTAMP}.sql.gz"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Perform backup
echo "Starting backup at $(date)"
PGPASSWORD="$DB_PASSWORD" pg_dump \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d "$DB_NAME" \
  --verbose \
  --format=custom \
  | gzip > "$BACKUP_FILE"

# Verify backup
if [ -f "$BACKUP_FILE" ]; then
    FILESIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo "Backup completed: $BACKUP_FILE ($FILESIZE)"
else
    echo "Backup failed!"
    exit 1
fi

# Upload to S3
if [ -n "$S3_BUCKET" ]; then
    echo "Uploading to S3..."
    aws s3 cp "$BACKUP_FILE" "s3://$S3_BUCKET/database/"
    echo "Upload completed"
fi

# Clean up old backups
echo "Cleaning up backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "prodory_backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Clean up old S3 backups
if [ -n "$S3_BUCKET" ]; then
    echo "Cleaning up S3 backups older than $RETENTION_DAYS days..."
    aws s3 ls "s3://$S3_BUCKET/database/" | \
    while read -r line; do
        FILE_DATE=$(echo "$line" | awk '{print $1}')
        FILE_NAME=$(echo "$line" | awk '{print $4}')
        FILE_AGE=$(( ( $(date +%s) - $(date -d "$FILE_DATE" +%s) ) / 86400 ))
        if [ "$FILE_AGE" -gt "$RETENTION_DAYS" ]; then
            aws s3 rm "s3://$S3_BUCKET/database/$FILE_NAME"
        fi
    done
fi

echo "Backup process completed at $(date)"
```

#### Restore from Backup

```bash
#!/bin/bash
# restore-database.sh

set -e

# Configuration
DB_HOST="${DATABASE_HOST:-localhost}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-prodory}"
DB_USER="${DATABASE_USER:-prodory}"
DB_PASSWORD="${DATABASE_PASSWORD}"
BACKUP_FILE="$1"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# Download from S3 if needed
if [[ "$BACKUP_FILE" == s3://* ]]; then
    echo "Downloading from S3..."
    LOCAL_FILE="/tmp/$(basename "$BACKUP_FILE")"
    aws s3 cp "$BACKUP_FILE" "$LOCAL_FILE"
    BACKUP_FILE="$LOCAL_FILE"
fi

# Verify backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Confirm restore
echo "WARNING: This will replace the current database!"
echo "Backup file: $BACKUP_FILE"
read -p "Are you sure? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled"
    exit 0
fi

# Create temporary database
echo "Creating temporary database..."
PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d postgres \
  -c "CREATE DATABASE prodory_restore;"

# Restore to temporary database
echo "Restoring backup..."
if [[ "$BACKUP_FILE" == *.gz ]]; then
    gunzip -c "$BACKUP_FILE" | PGPASSWORD="$DB_PASSWORD" pg_restore \
      -h "$DB_HOST" \
      -p "$DB_PORT" \
      -U "$DB_USER" \
      -d prodory_restore \
      --verbose \
      --no-owner \
      --no-privileges
else
    PGPASSWORD="$DB_PASSWORD" pg_restore \
      -h "$DB_HOST" \
      -p "$DB_PORT" \
      -U "$DB_USER" \
      -d prodory_restore \
      --verbose \
      --no-owner \
      --no-privileges \
      "$BACKUP_FILE"
fi

# Swap databases
echo "Swapping databases..."
PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d postgres \
  -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$DB_NAME';"

PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d postgres \
  -c "ALTER DATABASE $DB_NAME RENAME TO ${DB_NAME}_old;"

PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d postgres \
  -c "ALTER DATABASE prodory_restore RENAME TO $DB_NAME;"

# Drop old database
PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d postgres \
  -c "DROP DATABASE ${DB_NAME}_old;"

echo "Restore completed successfully!"
```

---

## Security Configuration

### TLS/SSL Configuration

#### Generate Self-Signed Certificates (Development)

```bash
# Generate private key
openssl genrsa -out server.key 2048

# Generate certificate signing request
openssl req -new -key server.key -out server.csr \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=prodory.local"

# Generate self-signed certificate
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Combine into PEM file
cat server.crt server.key > server.pem
```

#### Let's Encrypt Certificates (Production)

```bash
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d prodory.your-domain.com -d api.prodory.your-domain.com

# Certificates will be at:
# /etc/letsencrypt/live/prodory.your-domain.com/fullchain.pem
# /etc/letsencrypt/live/prodory.your-domain.com/privkey.pem

# Auto-renewal (certbot sets up cron automatically)
sudo certbot renew --dry-run
```

### Network Policies

#### Kubernetes Network Policy

```yaml
# kubernetes/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: prodory-default-deny
  namespace: prodory
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-ingress-nginx
  namespace: prodory
spec:
  podSelector:
    matchLabels:
      app: finops-dashboard
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: ingress-nginx
      ports:
        - protocol: TCP
          port: 80
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-database-access
  namespace: prodory
spec:
  podSelector:
    matchLabels:
      app: postgres
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: data-finops-agent
        - podSelector:
            matchLabels:
              app: k8s-in-a-box
      ports:
        - protocol: TCP
          port: 5432
```

### Pod Security Standards

```yaml
# kubernetes/pod-security.yaml
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
  namespace: prodory
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    runAsGroup: 1000
    fsGroup: 1000
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: app
      image: prodory/app:latest
      securityContext:
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        capabilities:
          drop:
            - ALL
      resources:
        limits:
          memory: "512Mi"
          cpu: "500m"
        requests:
          memory: "256Mi"
          cpu: "250m"
      volumeMounts:
        - name: tmp
          mountPath: /tmp
        - name: cache
          mountPath: /cache
  volumes:
    - name: tmp
      emptyDir: {}
    - name: cache
      emptyDir: {}
```

---

## Monitoring & Alerting

### Prometheus Configuration

```yaml
# config/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

rule_files:
  - /etc/prometheus/rules/*.yml

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'data-finops-agent'
    static_configs:
      - targets: ['data-finops-agent:8000']
    metrics_path: /metrics

  - job_name: 'k8s-in-a-box'
    static_configs:
      - targets: ['k8s-in-a-box:8080']
    metrics_path: /metrics

  - job_name: 'cloud-sentinel'
    static_configs:
      - targets: ['cloud-sentinel:8081']
    metrics_path: /metrics

  - job_name: 'postgres-exporter'
    static_configs:
      - targets: ['postgres-exporter:9187']

  - job_name: 'redis-exporter'
    static_configs:
      - targets: ['redis-exporter:9121']
```

### Alerting Rules

```yaml
# config/alert-rules.yml
groups:
  - name: prodory-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 5% for {{ $labels.service }}"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "95th percentile latency is above 1s for {{ $labels.service }}"

      - alert: DatabaseConnectionsHigh
        expr: pg_stat_activity_count > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High database connection count"
          description: "Database has {{ $value }} active connections"

      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Low disk space"
          description: "Disk space is below 10% on {{ $labels.device }}"

      - alert: MemoryUsageHigh
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is above 90%"

      - alert: PodCrashLooping
        expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Pod is crash looping"
          description: "Pod {{ $labels.pod }} in namespace {{ $labels.namespace }} is restarting frequently"
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Prodory Platform Overview",
    "tags": ["prodory"],
    "timezone": "UTC",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{service}}"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "{{service}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile - {{service}}"
          }
        ]
      },
      {
        "title": "Database Connections",
        "type": "singlestat",
        "targets": [
          {
            "expr": "pg_stat_activity_count"
          }
        ]
      }
    ]
  }
}
```

---

## Backup & Recovery

### Full System Backup

```bash
#!/bin/bash
# full-backup.sh

set -e

BACKUP_DIR="/backups/$(date +%Y%m%d_%H%M%S)"
S3_BUCKET="prodory-backups"

mkdir -p "$BACKUP_DIR"

echo "=== Starting Full System Backup ==="

# Backup database
echo "[1/5] Backing up database..."
./backup-database.sh

# Backup configuration
echo "[2/5] Backing up configuration..."
tar -czf "$BACKUP_DIR/config.tar.gz" \
  .env \
  config/ \
  kubernetes/ \
  docker-compose.yml

# Backup application data
echo "[3/5] Backing up application data..."
docker run --rm -v prodory_minio_data:/data -v "$BACKUP_DIR":/backup alpine \
  tar -czf /backup/minio_data.tar.gz -C /data .

# Backup Kubernetes resources
echo "[4/5] Backing up Kubernetes resources..."
kubectl get all -n prodory -o yaml > "$BACKUP_DIR/k8s-resources.yaml"
kubectl get configmap -n prodory -o yaml >> "$BACKUP_DIR/k8s-resources.yaml"
kubectl get secret -n prodory -o yaml >> "$BACKUP_DIR/k8s-resources.yaml"

# Create backup manifest
echo "[5/5] Creating backup manifest..."
cat > "$BACKUP_DIR/MANIFEST.json" << EOF
{
  "backup_date": "$(date -Iseconds)",
  "version": "$(git describe --tags --always)",
  "components": [
    "database",
    "configuration",
    "application_data",
    "kubernetes_resources"
  ],
  "files": [
    "database_backup.sql.gz",
    "config.tar.gz",
    "minio_data.tar.gz",
    "k8s-resources.yaml"
  ]
}
EOF

# Compress full backup
echo "Compressing backup..."
tar -czf "$BACKUP_DIR.tar.gz" -C "$(dirname "$BACKUP_DIR")" "$(basename "$BACKUP_DIR")"

# Upload to S3
echo "Uploading to S3..."
aws s3 cp "$BACKUP_DIR.tar.gz" "s3://$S3_BUCKET/full/"

# Cleanup
echo "Cleaning up..."
rm -rf "$BACKUP_DIR" "$BACKUP_DIR.tar.gz"

echo "=== Backup Completed ==="
echo "Backup location: s3://$S3_BUCKET/full/$(basename "$BACKUP_DIR").tar.gz"
```

### Disaster Recovery

```bash
#!/bin/bash
# disaster-recovery.sh

set -e

BACKUP_FILE="$1"
S3_BUCKET="prodory-backups"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file_or_s3_path>"
    exit 1
fi

echo "=== Starting Disaster Recovery ==="

# Download from S3 if needed
if [[ "$BACKUP_FILE" == s3://* ]]; then
    echo "[1/6] Downloading backup from S3..."
    LOCAL_FILE="/tmp/$(basename "$BACKUP_FILE")"
    aws s3 cp "$BACKUP_FILE" "$LOCAL_FILE"
    BACKUP_FILE="$LOCAL_FILE"
fi

# Extract backup
echo "[2/6] Extracting backup..."
EXTRACT_DIR="/tmp/prodory-restore-$(date +%s)"
mkdir -p "$EXTRACT_DIR"
tar -xzf "$BACKUP_FILE" -C "$EXTRACT_DIR"
BACKUP_DIR="$EXTRACT_DIR/$(ls "$EXTRACT_DIR")"

# Restore database
echo "[3/6] Restoring database..."
./restore-database.sh "$BACKUP_DIR/database_backup.sql.gz"

# Restore configuration
echo "[4/6] Restoring configuration..."
tar -xzf "$BACKUP_DIR/config.tar.gz" -C /

# Restore application data
echo "[5/6] Restoring application data..."
docker run --rm -v prodory_minio_data:/data -v "$BACKUP_DIR":/backup alpine \
  tar -xzf /backup/minio_data.tar.gz -C /data

# Restore Kubernetes resources
echo "[6/6] Restoring Kubernetes resources..."
kubectl apply -f "$BACKUP_DIR/k8s-resources.yaml"

# Cleanup
rm -rf "$EXTRACT_DIR"
if [[ "$BACKUP_FILE" == /tmp/* ]]; then
    rm -f "$BACKUP_FILE"
fi

echo "=== Disaster Recovery Completed ==="
echo "Please verify all services are running correctly"
```

---

## Upgrade Procedures

### Pre-Upgrade Checklist

- [ ] Review release notes
- [ ] Backup current system
- [ ] Test upgrade in staging environment
- [ ] Notify users of maintenance window
- [ ] Prepare rollback plan

### Rolling Update Procedure

```bash
#!/bin/bash
# upgrade.sh

VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "=== Starting Upgrade to $VERSION ==="

# Backup before upgrade
echo "Creating pre-upgrade backup..."
./full-backup.sh

# Update images
echo "Updating Docker images..."
docker-compose pull

# Run database migrations
echo "Running database migrations..."
docker-compose run --rm data-finops-agent alembic upgrade head

# Rolling restart
echo "Performing rolling restart..."
docker-compose up -d --no-deps data-finops-agent
sleep 30

docker-compose up -d --no-deps finops-dashboard
sleep 30

docker-compose up -d --no-deps k8s-in-a-box
sleep 30

docker-compose up -d --no-deps storage-autoscaler
sleep 30

docker-compose up -d --no-deps cloud-sentinel
sleep 30

docker-compose up -d --no-deps vmware-migration
sleep 30

docker-compose up -d --no-deps vm2container

# Verify upgrade
echo "Verifying upgrade..."
docker-compose ps

echo "=== Upgrade Completed ==="
```

### Rollback Procedure

```bash
#!/bin/bash
# rollback.sh

BACKUP_FILE="$1"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

echo "=== Starting Rollback ==="
echo "WARNING: This will restore the system to the backup state!"
read -p "Are you sure? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Rollback cancelled"
    exit 0
fi

# Stop all services
echo "Stopping services..."
docker-compose down

# Restore from backup
echo "Restoring from backup..."
./disaster-recovery.sh "$BACKUP_FILE"

# Start services
echo "Starting services..."
docker-compose up -d

echo "=== Rollback Completed ==="
```

---

## Troubleshooting

### Common Issues

#### Issue: Database Connection Failed

**Symptoms:**
- Application logs show connection errors
- Services fail to start

**Diagnosis:**
```bash
# Check PostgreSQL status
docker-compose ps postgres
kubectl get pods -l app=postgres

# Test connection
docker-compose exec postgres pg_isready -U prodory

# Check logs
docker-compose logs postgres
kubectl logs -l app=postgres
```

**Solutions:**

1. **Check credentials**
   ```bash
   # Verify environment variables
   echo $DATABASE_PASSWORD
   
   # Update secrets
   kubectl create secret generic database-credentials \
     --from-literal=password=newpassword \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

2. **Restart PostgreSQL**
   ```bash
   docker-compose restart postgres
   kubectl rollout restart deployment/postgres
   ```

3. **Check disk space**
   ```bash
   df -h
   docker system df
   ```

#### Issue: High Memory Usage

**Symptoms:**
- System slowdown
- OOM kills
- High swap usage

**Diagnosis:**
```bash
# Check memory usage
free -h
docker stats
kubectl top pods

# Find memory-intensive processes
ps aux --sort=-%mem | head -20
```

**Solutions:**

1. **Increase memory limits**
   ```yaml
   # Update docker-compose.yml
   services:
     data-finops-agent:
       deploy:
         resources:
           limits:
             memory: 4G
   ```

2. **Restart services**
   ```bash
   docker-compose restart
   ```

3. **Scale down replicas**
   ```bash
   kubectl scale deployment data-finops-agent --replicas=1
   ```

#### Issue: Pod Stuck in Pending

**Symptoms:**
- Pod status shows "Pending"
- No events in describe output

**Diagnosis:**
```bash
# Check pod status
kubectl get pods -o wide
kubectl describe pod <pod-name>

# Check node resources
kubectl describe node

# Check PVC status
kubectl get pvc
```

**Solutions:**

1. **Check resource quotas**
   ```bash
   kubectl get resourcequota
   kubectl describe resourcequota
   ```

2. **Verify storage class**
   ```bash
   kubectl get storageclass
   kubectl get pvc
   ```

3. **Check node selectors**
   ```bash
   kubectl get nodes --show-labels
   ```

### Debug Commands

```bash
# Get shell into pod
kubectl exec -it <pod-name> -- /bin/sh

# Check environment variables
kubectl exec <pod-name> -- env | sort

# Test network connectivity
kubectl exec <pod-name> -- nc -zv postgres 5432

# Check DNS resolution
kubectl exec <pod-name> -- nslookup kubernetes.default

# View resource usage
kubectl top pod <pod-name>
kubectl top node

# Get detailed pod info
kubectl get pod <pod-name> -o yaml

# Check events
kubectl get events --sort-by='.lastTimestamp'
```

---

## API Reference

### Authentication

```bash
# Get access token
curl -X POST https://api.prodory.local/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your-password"
  }'

# Response
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "bearer",
  "expires_in": 1800
}

# Use token
curl https://api.prodory.local/finops/costs \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### API Endpoints

#### FinOps API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/finops/health` | GET | Health check |
| `/finops/costs` | GET | Get cost data |
| `/finops/costs/summary` | GET | Cost summary |
| `/finops/recommendations` | GET | List recommendations |
| `/finops/recommendations/{id}/apply` | POST | Apply recommendation |
| `/finops/budgets` | GET/POST | List/Create budgets |
| `/finops/reports` | GET | Generate reports |

#### Kubernetes API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/k8s/health` | GET | Health check |
| `/k8s/clusters` | GET/POST | List/Create clusters |
| `/k8s/clusters/{id}` | GET/PUT/DELETE | Cluster operations |
| `/k8s/clusters/{id}/scale` | POST | Scale cluster |
| `/k8s/namespaces` | GET | List namespaces |
| `/k8s/deployments` | GET/POST | List/Create deployments |

#### Storage API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/storage/health` | GET | Health check |
| `/storage/volumes` | GET | List volumes |
| `/storage/volumes/{id}/expand` | POST | Expand volume |
| `/storage/volumes/{id}/snapshot` | POST | Create snapshot |
| `/storage/policies` | GET/POST | List/Create policies |

---

## Custom Modifications

### Adding a New Service

1. **Create service directory**
   ```bash
   mkdir services/my-new-service
   cd services/my-new-service
   ```

2. **Create Dockerfile**
   ```dockerfile
   FROM python:3.11-slim
   WORKDIR /app
   COPY requirements.txt .
   RUN pip install -r requirements.txt
   COPY . .
   CMD ["python", "main.py"]
   ```

3. **Add to docker-compose.yml**
   ```yaml
   my-new-service:
     build: ./services/my-new-service
     ports:
       - "8002:8000"
     environment:
       - DATABASE_URL=${DATABASE_URL}
     depends_on:
       - postgres
   ```

4. **Add to Kubernetes**
   ```yaml
   # kubernetes/my-new-service.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: my-new-service
   spec:
     replicas: 2
     selector:
       matchLabels:
         app: my-new-service
     template:
       metadata:
         labels:
           app: my-new-service
       spec:
         containers:
           - name: my-new-service
             image: prodory/my-new-service:latest
             ports:
               - containerPort: 8000
   ```

### Custom Alert Rules

Add to `config/alert-rules.yml`:

```yaml
groups:
  - name: custom-alerts
    rules:
      - alert: CustomMetricThreshold
        expr: your_custom_metric > threshold
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Custom alert triggered"
```

### Custom Dashboards

Add Grafana dashboard JSON to `config/grafana/dashboards/`:

```bash
# Create dashboard directory
mkdir -p config/grafana/dashboards

# Add your dashboard JSON
cp my-dashboard.json config/grafana/dashboards/

# Update Grafana config to load dashboards
```

---

## Appendix

### Environment Variable Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `REDIS_URL` | Yes | - | Redis connection string |
| `JWT_SECRET` | Yes | - | JWT signing secret |
| `AWS_ACCESS_KEY_ID` | No | - | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | No | - | AWS secret key |
| `OPENAI_API_KEY` | No | - | OpenAI API key |
| `SENTRY_DSN` | No | - | Sentry error tracking |

### File Locations

| File | Location | Description |
|------|----------|-------------|
| `.env` | Project root | Environment variables |
| `docker-compose.yml` | Project root | Docker Compose config |
| `kubernetes/` | Project root | K8s manifests |
| `config/` | Project root | Configuration files |
| `logs/` | Project root | Application logs |
| `backups/` | Project root | Backup files |

### Useful Commands Quick Reference

```bash
# Docker Compose
docker-compose up -d              # Start all services
docker-compose down               # Stop all services
docker-compose logs -f            # View logs
docker-compose ps                 # List services
docker-compose build              # Build images
docker-compose pull               # Pull latest images

# Kubernetes
kubectl get all -n prodory        # List all resources
kubectl get pods -o wide          # List pods with details
kubectl logs -f <pod>             # View pod logs
kubectl exec -it <pod> -- sh      # Shell into pod
kubectl apply -f <file>           # Apply manifest
kubectl delete -f <file>          # Delete manifest
kubectl port-forward <pod> 8080   # Port forward

# Database
pg_dump -h localhost -U prodory prodory > backup.sql
psql -h localhost -U prodory prodory < backup.sql

# Redis
redis-cli -a password ping
redis-cli -a password info
```

---

*For additional support, contact: admin@prodory.local*
*Documentation Version: 1.0.0*
*Last Updated: 2024-01-15*
