# Prodory Platform - Enterprise DevOps/FinOps Suite

A comprehensive suite of Podman-based enterprise DevOps and FinOps tools designed for modern cloud-native infrastructure management.

## Overview

Prodory Platform provides 7 integrated tools for managing cloud costs, monitoring infrastructure, provisioning Kubernetes clusters, and migrating workloads:

1. **AI FinOps Dashboard** - React-based cloud cost management and optimization
2. **Data FinOps Agent** - Python FastAPI backend for cost analysis and AI recommendations
3. **Cloud Sentinel** - Go-based multi-cloud monitoring and alerting
4. **Storage Autoscaler** - Kubernetes operator for PVC auto-scaling
5. **Kubernetes-in-a-Box** - Lightweight k3s cluster provisioning tool
6. **VMware to KubeVirt Migration** - VM migration tool for Kubernetes
7. **VM to Container Migration** - Node.js tool for containerizing VMs

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Prodory Platform                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │  FinOps         │  │  Cloud Sentinel │  │  K8s-in-a-Box   │             │
│  │  Dashboard      │  │  (Monitoring)   │  │  (Provisioning) │             │
│  │  (React)        │  │  (Go)           │  │  (Go)           │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
│           │                    │                    │                        │
│           └────────────────────┼────────────────────┘                        │
│                                │                                             │
│                     ┌──────────┴──────────┐                                  │
│                     │   Data FinOps Agent  │                                  │
│                     │   (Python/FastAPI)   │                                  │
│                     └──────────┬──────────┘                                  │
│                                │                                             │
│  ┌─────────────────┐  ┌────────┴────────┐  ┌─────────────────┐             │
│  │  Storage        │  │   Supabase      │  │  Migration      │             │
│  │  Autoscaler     │  │   (Backend)     │  │  Tools          │             │
│  │  (Python)       │  │                 │  │  (Go/Node.js)   │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites

- Docker/Podman 4.0+
- Kubernetes 1.25+ (for deployment)
- kubectl configured
- 8GB+ RAM available

### Local Development

```bash
# Clone the repository
git clone <repository-url>
cd prodory-platform

# Start all services with Docker Compose
docker-compose up -d

# Or use Podman
podman-compose up -d

# Access the services
# FinOps Dashboard: http://localhost:8080
# Data FinOps API: http://localhost:8081
# Cloud Sentinel: http://localhost:8083
```

### Kubernetes Deployment

```bash
# Deploy to Kubernetes
./scripts/deploy.sh kubernetes

# Check deployment status
kubectl get pods -n prodory

# Port forward for local access
kubectl port-forward svc/finops-dashboard 8080:80 -n prodory
```

## Services

### 1. AI FinOps Dashboard

React-based frontend for cloud cost management.

**Features:**
- Real-time cost dashboards
- Multi-cloud cost aggregation (AWS, Azure, GCP)
- AI-powered cost recommendations
- Budget management with alerts
- Cost anomaly detection
- Custom reports and exports

**Configuration:**
```yaml
environment:
  - VITE_API_URL=http://data-finops-agent:8081
  - VITE_SUPABASE_URL=${SUPABASE_URL}
  - VITE_SUPABASE_ANON_KEY=${SUPABASE_ANON_KEY}
```

**Access:** http://localhost:8080

### 2. Data FinOps Agent

Python FastAPI backend for cost analysis and optimization.

**Features:**
- Multi-cloud cost data ingestion
- AI/ML cost forecasting
- Anomaly detection
- Rightsizing recommendations
- Reserved instance planning
- Budget tracking

**API Endpoints:**
- `GET /health` - Health check
- `GET /dashboard/summary` - Dashboard summary
- `GET /costs/time-series` - Cost time series data
- `POST /recommendations/generate` - Generate recommendations
- `GET /budgets` - List budgets

**Configuration:**
```yaml
environment:
  - DATABASE_URL=${DATABASE_URL}
  - SUPABASE_URL=${SUPABASE_URL}
  - SUPABASE_KEY=${SUPABASE_KEY}
  - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
  - AZURE_CLIENT_ID=${AZURE_CLIENT_ID}
  - GCP_PROJECT_ID=${GCP_PROJECT_ID}
```

### 3. Cloud Sentinel

Go-based multi-cloud monitoring and alerting service.

**Features:**
- Prometheus metrics collection
- Multi-cloud monitoring (AWS CloudWatch, Azure Monitor, GCP Monitoring)
- Kubernetes metrics
- Custom alert rules
- Slack/PagerDuty/Email notifications
- Alert silencing
- Metric visualization

**API Endpoints:**
- `GET /health` - Health check
- `GET /api/v1/dashboard/summary` - Dashboard summary
- `GET /api/v1/metrics/prometheus` - Prometheus metrics
- `POST /api/v1/alerts` - Create alert
- `GET /api/v1/rules` - List alert rules

**Configuration:**
```yaml
environment:
  - SENTINEL_PROMETHEUS_URL=http://prometheus:9090
  - SENTINEL_SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL}
  - SENTINEL_PAGERDUTY_INTEGRATION_KEY=${PAGERDUTY_KEY}
```

### 4. Storage Autoscaler

Kubernetes operator for automatic PVC scaling.

**Features:**
- Automatic PVC scaling based on usage
- Configurable thresholds and scale factors
- Support for multiple storage classes
- Prometheus metrics export
- Safe scaling with data preservation
- Cooldown periods

**Installation:**
```bash
# Install CRD
kubectl apply -f services/storage-autoscaler/helm/crd.yaml

# Install operator
helm install storage-autoscaler services/storage-autoscaler/helm

# Create autoscaler configuration
cat <<EOF | kubectl apply -f -
apiVersion: prodory.io/v1
kind: StorageAutoscaler
metadata:
  name: default-autoscaler
spec:
  enabled: true
  policy:
    thresholdPercent: 80
    scaleFactor: 1.5
    maxSize: "1Ti"
    cooldownMinutes: 60
EOF
```

### 5. Kubernetes-in-a-Box

Go-based lightweight k3s cluster provisioning tool.

**Features:**
- Single-node and HA cluster deployment
- Multi-provider support (AWS, Azure, GCP, local)
- Automated TLS certificate management
- Integrated ingress controller
- Local storage provisioner
- Metrics server
- Add-on management

**Usage:**
```bash
# Build the CLI
cd services/kubernetes-in-a-box
go build -o kib ./cmd/kib

# Create a local cluster
./kib create --name local-dev --provider local

# Create AWS cluster with HA masters
./kib create --name production \
  --provider aws \
  --region us-west-2 \
  --masters 3 \
  --workers 3

# List clusters
./kib list

# Get kubeconfig
export KUBECONFIG=$(./kib kubeconfig local-dev)

# Delete cluster
./kib delete local-dev
```

### 6. VMware to KubeVirt Migration

Go-based tool for migrating VMware VMs to Kubernetes KubeVirt.

**Features:**
- VM inventory discovery from vCenter
- VM analysis and compatibility checking
- Hot and cold migration support
- Disk image conversion (VMDK to RAW/QCOW2)
- Network mapping
- Progress tracking
- Rollback support

**Usage:**
```bash
# Discover VMs in vCenter
./vmware-migrate discover \
  --vcenter-url https://vcenter.local \
  --vcenter-user admin

# Analyze a VM
./vmware-migrate analyze \
  --vcenter-url https://vcenter.local \
  --source-vm my-vm

# Migrate to KubeVirt
./vmware-migrate migrate \
  --vcenter-url https://vcenter.local \
  --source-vm my-vm \
  --target-name my-vm-kubevirt \
  --namespace default \
  --storage-class standard
```

### 7. VM to Container Migration

Node.js tool for containerizing virtual machines.

**Features:**
- VM analysis and discovery
- Application dependency detection
- Dockerfile generation
- Container image building
- Docker Compose generation
- Kubernetes manifest generation

**Usage:**
```bash
# Install dependencies
npm install

# Analyze a VM
node src/cli.js analyze --source /path/to/vm --type vmware

# Migrate to container
node src/cli.js migrate \
  --source /path/to/vm \
  --name my-container \
  --build \
  --generate-compose

# List detected services
node src/cli.js list-services --source /path/to/vm
```

## Database Schema

The platform uses Supabase (PostgreSQL) as the backend database. See `supabase/migrations/` for the complete schema.

**Key Tables:**
- `cloud_accounts` - Cloud provider account configurations
- `cost_data` - Cost and usage data
- `budgets` - Budget definitions and tracking
- `recommendations` - AI-generated cost recommendations
- `alerts` - Monitoring alerts
- `k8s_clusters` - Kubernetes cluster information
- `vm_migrations` - VM migration tracking

## Configuration

### Environment Variables

Create a `.env` file in the project root:

```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/prodory
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-service-key
SUPABASE_JWT_SECRET=your-jwt-secret

# Cloud Providers
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=us-east-1

AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret
AZURE_SUBSCRIPTION_ID=your-subscription-id

GCP_PROJECT_ID=your-project-id
GCP_SERVICE_ACCOUNT_KEY='{"type": "service_account", ...}'

# Notifications
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
PAGERDUTY_INTEGRATION_KEY=your-integration-key

# Monitoring
PROMETHEUS_URL=http://prometheus:9090
```

## Deployment

### Docker Compose (Development)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Kubernetes (Production)

```bash
# Create namespace and apply manifests
kubectl apply -f kubernetes/namespace.yaml
kubectl apply -f kubernetes/configmap.yaml
kubectl apply -f kubernetes/secrets.yaml
kubectl apply -f kubernetes/postgres.yaml
kubectl apply -f kubernetes/redis.yaml
kubectl apply -f kubernetes/data-finops-agent.yaml
kubectl apply -f kubernetes/finops-dashboard.yaml
kubectl apply -f kubernetes/cloud-sentinel.yaml
kubectl apply -f kubernetes/ingress.yaml

# Verify deployment
kubectl get pods -n prodory
kubectl get svc -n prodory
kubectl get ingress -n prodory
```

### Podman (Rootless)

```bash
# Generate Kubernetes manifests from compose
podman-compose -f docker-compose.yml convert > podman-k8s.yaml

# Deploy with Podman
podman play kube podman-k8s.yaml
```

## Monitoring

### Prometheus Metrics

All services expose Prometheus metrics:

- FinOps Dashboard: `http://localhost:8080/metrics`
- Data FinOps Agent: `http://localhost:8081/metrics`
- Cloud Sentinel: `http://localhost:8083/metrics`
- Storage Autoscaler: `http://localhost:8084/metrics`

### Grafana Dashboards

Import the provided dashboards from `monitoring/grafana/`:

- FinOps Dashboard
- Cloud Sentinel Overview
- Kubernetes Cluster Metrics

## API Documentation

### Data FinOps Agent

Swagger UI: http://localhost:8081/docs

### Cloud Sentinel

Swagger UI: http://localhost:8083/swagger/index.html

## Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run specific service tests
cd services/data-finops-agent && pytest
cd services/cloud-sentinel && go test ./...
```

## Troubleshooting

### Common Issues

**Issue:** Services fail to start with database connection errors

**Solution:**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check database exists
docker-compose exec postgres psql -U prodory -l

# Run migrations
make migrate
```

**Issue:** Cloud provider credentials not working

**Solution:**
- Verify credentials are correctly set in environment variables
- Check IAM permissions for the service accounts
- Verify network connectivity to cloud APIs

**Issue:** Storage Autoscaler not scaling PVCs

**Solution:**
```bash
# Check operator logs
kubectl logs -n prodory -l app=storage-autoscaler

# Verify StorageAutoscaler CR exists
kubectl get storageautoscalers

# Check PVC metrics are available
kubectl top pvc
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Support

- Documentation: https://docs.prodory.io
- Issues: https://github.com/prodory/platform/issues
- Email: support@prodory.io

## Roadmap

- [ ] Multi-region cost optimization
- [ ] Advanced AI forecasting models
- [ ] GitOps integration
- [ ] Service mesh support
- [ ] FinOps policy engine
- [ ] Carbon footprint tracking
