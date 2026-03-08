# Prodory Platform - Comprehensive User Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [AI FinOps Dashboard](#ai-finops-dashboard)
4. [Data FinOps Agent](#data-finops-agent)
5. [Kubernetes-in-a-Box](#kubernetes-in-a-box)
6. [Storage Autoscaler](#storage-autoscaler)
7. [Cloud Sentinel](#cloud-sentinel)
8. [VMware to OpenShift Migration](#vmware-to-openshift-migration)
9. [VM to Container Migration](#vm-to-container-migration)
10. [Integration & Workflows](#integration--workflows)
11. [Troubleshooting](#troubleshooting)
12. [Best Practices](#best-practices)

---

## Introduction

### What is Prodory Platform?

Prodory Platform is a comprehensive suite of 7 enterprise DevOps/FinOps tools designed to help organizations:
- Optimize cloud costs and resource utilization
- Automate Kubernetes deployments
- Manage storage efficiently
- Monitor security and compliance
- Migrate workloads from VMs to containers

### Platform Components Overview

| Component | Purpose | User Type |
|-----------|---------|-----------|
| AI FinOps Dashboard | Cloud cost visualization and optimization | Finance, DevOps |
| Data FinOps Agent | Automated cost data collection and analysis | DevOps, Finance |
| Kubernetes-in-a-Box | Simplified K8s cluster management | Platform Engineers |
| Storage Autoscaler | Intelligent storage capacity management | Storage Admins |
| Cloud Sentinel | Security monitoring and compliance | Security Teams |
| VMware to OpenShift Migration | VM to KubeVirt migration | Infrastructure Teams |
| VM to Container Migration | Traditional VM modernization | DevOps Teams |

---

## Getting Started

### First-Time Access

#### 1. Access the Platform

```
URL: https://prodory.your-domain.com
Default Admin: admin@prodory.local
Default Password: ChangeMeNow123!
```

#### 2. Initial Setup Checklist

- [ ] Change default admin password
- [ ] Configure email/SMTP settings
- [ ] Add cloud provider credentials
- [ ] Set up initial budgets
- [ ] Configure notification channels
- [ ] Review and accept security policies

#### 3. Dashboard Navigation

```
┌─────────────────────────────────────────────────────────────┐
│  Prodory Platform                                           │
├──────────┬──────────────────────────────────────────────────┤
│          │                                                  │
│  ┌────┐  │  ┌──────────────────────────────────────────┐  │
│  │🏠  │  │  │  Dashboard Overview                      │  │
│  │📊  │  │  │  - Cost Summary                          │  │
│  │⚙️  │  │  │  - Active Alerts                         │  │
│  │🔒  │  │  │  - Recent Recommendations                │  │
│  │📦  │  │  └──────────────────────────────────────────┘  │
│  │💾  │  │                                                  │
│  │🔄  │  │  ┌──────────────────────────────────────────┐  │
│  │☁️  │  │  │  Quick Actions                           │  │
│  └────┘  │  │  - Add Cloud Provider                      │  │
│          │  │  - Create Budget                           │  │
│  Modules │  │  - View Reports                            │  │
│          │  └──────────────────────────────────────────┘  │
└──────────┴──────────────────────────────────────────────────┘
```

---

## AI FinOps Dashboard

### Overview

The AI FinOps Dashboard provides intelligent cloud cost management with AI-powered recommendations for optimization.

### Accessing the Dashboard

1. Log in to Prodory Platform
2. Click "FinOps Dashboard" in the left sidebar
3. The main dashboard displays:
   - **Total Monthly Spend**: Current month's cloud costs
   - **Forecasted Spend**: AI-predicted end-of-month costs
   - **Potential Savings**: Total savings from pending recommendations
   - **Active Resources**: Number of cloud resources being monitored

### Understanding Key Metrics

| Metric | Description | Normal Range | Action Needed When |
|--------|-------------|--------------|-------------------|
| Total Monthly Spend | Sum of all cloud costs this month | Varies by org | Unexpected spikes |
| Spend Change | Percentage change vs last month | ±10% | >20% change |
| Forecasted Spend | AI prediction for month-end | Within budget | Exceeds budget |
| Potential Savings | Sum of pending recommendations | 10-30% of spend | <5% indicates efficiency |
| Cost per Resource | Average cost per active resource | Benchmark needed | Above benchmark |

### Cost Trend Analysis

#### Viewing Cost Trends

1. Navigate to **FinOps Dashboard → Cost Trends**
2. Select date range (default: Last 30 days)
3. Use filters:
   - **Provider**: AWS, Azure, GCP, or All
   - **Service**: EC2, S3, RDS, etc.
   - **Team/Department**: Filter by tags

#### Interpreting the Trend Chart

```
Cost Trend - Last 90 Days

$50K ┤                                    ╭─╮
     │                              ╭────╯  ╰──╮
$40K ┤                        ╭────╯            ╰──╮
     │                  ╭────╯                      ╰──
$30K ┤            ╭────╯
     │      ╭────╯
$20K ┤╭────╯
     │
$10K ┤
     └────┬────┬────┬────┬────┬────┬────┬────┬────┬────
         Jan  Feb  Mar  Apr  May  Jun  Jul  Aug  Sep

Legend:
─ Blue line: Actual spend
─ Dashed line: Forecasted spend
─ Red dots: Anomaly alerts
```

#### Anomaly Detection

The AI automatically detects:
- **Unusual spikes**: Costs exceeding 2 standard deviations
- **Gradual increases**: Sustained upward trends
- **Weekend activity**: Unexpected usage during off-hours
- **New services**: First-time usage of expensive services

### Service Breakdown Analysis

#### Viewing Service Costs

1. Go to **FinOps Dashboard → Service Breakdown**
2. View the pie chart showing cost distribution
3. Click any slice to drill down:
   - See individual resources
   - View usage patterns
   - Access recommendations

#### Understanding Service Categories

| Category | Services Included | Typical % of Bill |
|----------|-------------------|-------------------|
| Compute | EC2, VMs, Containers | 40-60% |
| Storage | S3, EBS, Blob, Disks | 15-25% |
| Database | RDS, DynamoDB, CosmosDB | 10-20% |
| Network | Data Transfer, CDN, Load Balancers | 5-15% |
| Other | Monitoring, Logging, Misc | 5-10% |

### AI Recommendations

#### Types of Recommendations

| Type | Description | Typical Savings | Risk Level |
|------|-------------|-----------------|------------|
| Rightsizing | Adjust instance sizes based on usage | 20-40% | Low |
| Reserved Capacity | Purchase RIs/SPs for steady workloads | 30-60% | Low |
| Storage Optimization | Delete unused volumes, enable tiering | 10-30% | Medium |
| Idle Resource Cleanup | Remove unused resources | 5-15% | Medium |
| Spot Instances | Use spot/preemptible for flexible workloads | 60-90% | High |
| Scheduling | Auto-shutdown dev/test during off-hours | 30-50% | Low |

#### Reviewing Recommendations

1. Navigate to **FinOps Dashboard → Recommendations**
2. Review the list of pending recommendations
3. Each recommendation shows:
   - **Impact**: High, Medium, Low
   - **Effort**: Easy, Medium, Complex
   - **Estimated Savings**: Monthly dollar amount
   - **Confidence**: AI confidence score (0-100%)

#### Applying a Recommendation

**Step-by-Step Process:**

1. **Select Recommendation**
   ```
   Click on recommendation card → View Details
   ```

2. **Review Affected Resources**
   ```
   ┌─────────────────────────────────────────┐
   │ Recommendation: Rightsize EC2 instances │
   ├─────────────────────────────────────────┤
   │                                         │
   │ Affected Resources (3):                 │
   │ ┌─────────────┬──────────┬────────────┐ │
   │ │ Resource    │ Current  │ Proposed   │ │
   │ ├─────────────┼──────────┼────────────┤ │
   │ │ i-abc123    │ m5.large │ m5.medium  │ │
   │ │ i-def456    │ m5.xlarge│ m5.large   │ │
   │ │ i-ghi789    │ c5.2xlarge│ c5.xlarge │ │
   │ └─────────────┴──────────┴────────────┘ │
   │                                         │
   │ Estimated Monthly Savings: $450.00      │
   │ Confidence Score: 87%                   │
   └─────────────────────────────────────────┘
   ```

3. **Run Dry Run (Optional but Recommended)**
   ```
   Click "Dry Run" → System simulates changes
   → Review simulation results
   → Check for any warnings
   ```

4. **Schedule or Apply**
   ```
   Option A: Apply Now
   Option B: Schedule for Maintenance Window
   Option C: Add to Change Request
   ```

5. **Monitor After Application**
   ```
   → Dashboard shows "Monitoring" status
   → Wait 24-48 hours for validation
   → Review actual vs predicted savings
   ```

#### Recommendation Status Tracking

| Status | Meaning | Next Action |
|--------|---------|-------------|
| Pending | Not yet reviewed | Review and decide |
| Under Review | Being evaluated | Wait or expedite |
| Approved | Ready to apply | Schedule application |
| Applied | Successfully implemented | Monitor results |
| Dismissed | Manually ignored | No action needed |
| Failed | Application failed | Review error and retry |

### Budget Management

#### Creating a Budget

**Step 1: Navigate to Budgets**
```
FinOps Dashboard → Budgets → Create New Budget
```

**Step 2: Configure Budget Parameters**

```yaml
Budget Name: Production Infrastructure
Amount: $50,000.00
Period: Monthly
Start Date: 2024-01-01
End Date: 2024-12-31 (Optional)

Scope:
  - Providers: AWS, Azure
  - Services: All Compute, Storage
  - Tags: Environment=Production
  - Teams: Platform Engineering

Alert Thresholds:
  - 50%: Info notification
  - 75%: Warning notification  
  - 90%: Critical notification
  - 100%: Budget exceeded alert

Notifications:
  - Email: team@company.com
  - Slack: #cloud-costs
  - PagerDuty: High priority at 90%
```

**Step 3: Set Advanced Options**

```yaml
Forecast Alerts:
  - Alert if forecast exceeds budget by: 10%
  - Forecast window: 7 days

Auto-Actions:
  - At 95%: Pause non-critical resource creation
  - At 100%: Require approval for new resources

Cost Anomaly Detection:
  - Enable AI anomaly detection: Yes
  - Sensitivity: Medium
  - Minimum alert amount: $100
```

#### Budget Status Indicators

| Status | Color | Description | Action |
|--------|-------|-------------|--------|
| On Track | Green | Spending within expected range | Continue monitoring |
| Warning | Yellow | Approaching alert threshold | Review spending |
| Critical | Orange | Near budget limit | Take immediate action |
| Exceeded | Red | Budget limit reached | Stop new spending |
| Forecasted Exceed | Purple | AI predicts overage | Plan adjustments |

#### Budget Reports

Generate detailed budget reports:

```
Budget Report - January 2024
═══════════════════════════════════════════════════════════════

Budget: Production Infrastructure
Budget Amount: $50,000.00
Actual Spend: $47,350.00 (94.7%)
Forecast: $51,200.00 (102.4% - OVER BUDGET)

Daily Breakdown:
┌─────────────┬──────────┬──────────┬──────────────┐
│ Date        │ Daily    │ Cumulative │ % of Budget │
├─────────────┼──────────┼──────────┼──────────────┤
│ Jan 1-7     │ $8,200   │ $8,200   │ 16.4%        │
│ Jan 8-14    │ $9,100   │ $17,300  │ 34.6%        │
│ Jan 15-21   │ $12,400  │ $29,700  │ 59.4%        │
│ Jan 22-28   │ $14,200  │ $43,900  │ 87.8%        │
│ Jan 29-31   │ $3,450   │ $47,350  │ 94.7%        │
└─────────────┴──────────┴──────────┴──────────────┘

Top Cost Drivers:
1. EC2 Instances: $18,500 (39.1%)
2. RDS Databases: $12,300 (26.0%)
3. S3 Storage: $8,200 (17.3%)
4. Data Transfer: $5,100 (10.8%)
5. Other: $3,250 (6.8%)

Recommendations Applied:
- Rightsized 12 EC2 instances: Saved $2,400
- Deleted 5 unused EBS volumes: Saved $350
- Enabled S3 Intelligent Tiering: Projected $500/month savings
```

### Reports & Analytics

#### Available Report Types

| Report | Description | Frequency | Recipients |
|--------|-------------|-----------|------------|
| Executive Summary | High-level cost overview | Monthly | C-suite |
| Detailed Cost Breakdown | Service-level costs | Weekly | Finance |
| Resource Utilization | Efficiency metrics | Weekly | Engineering |
| Reserved Capacity Analysis | RI/SP recommendations | Monthly | Finance |
| Cost Allocation | Chargeback by team/project | Monthly | All teams |
| Anomaly Report | Unusual spending patterns | Daily | DevOps |
| Forecast Report | Future cost predictions | Monthly | Planning |

#### Scheduling Reports

1. Go to **FinOps Dashboard → Reports → Scheduled Reports**
2. Click "Create Schedule"
3. Configure:
   ```yaml
   Report Type: Executive Summary
   Frequency: Monthly
   Day: 1st of month
   Time: 08:00 UTC
   Format: PDF
   
   Recipients:
     - cfo@company.com
     - cto@company.com
     - finance@company.com
   
   Include:
     - Cost summary
     - Month-over-month comparison
     - Top recommendations
     - Budget status
   ```

#### Exporting Data

Export options available:
- **CSV**: Raw data for analysis
- **Excel**: Formatted with charts
- **PDF**: Presentation-ready reports
- **JSON**: API integration
- **Parquet**: Big data analysis

---

## Data FinOps Agent

### Overview

The Data FinOps Agent automatically collects, processes, and analyzes cloud cost data from multiple providers.

### Connecting Cloud Providers

#### AWS Connection

**Step 1: Create IAM Role/User**

```bash
# Using AWS CLI
aws iam create-user --user-name prodory-finops

# Create and attach policy
cat > prodory-finops-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ce:GetCostAndUsage",
        "ce:GetCostForecast",
        "ce:GetUsageForecast",
        "ce:GetReservationUtilization",
        "ce:GetSavingsPlansUtilization",
        "ec2:DescribeInstances",
        "ec2:DescribeVolumes",
        "ec2:DescribeSnapshots",
        "rds:DescribeDBInstances",
        "s3:ListAllMyBuckets",
        "s3:GetBucketLocation",
        "cloudwatch:GetMetricData",
        "organizations:ListAccounts"
      ],
      "Resource": "*"
    }
  ]
}
EOF

aws iam put-user-policy \
  --user-name prodory-finops \
  --policy-name ProdoryFinOpsPolicy \
  --policy-document file://prodory-finops-policy.json

# Create access keys
aws iam create-access-key --user-name prodory-finops
```

**Step 2: Add to Prodory Platform**

1. Navigate to **Settings → Cloud Providers → Add Provider**
2. Select "AWS"
3. Enter credentials:
   ```yaml
   Provider Name: AWS Production
   Access Key ID: AKIA...
   Secret Access Key: ********
   Default Region: us-east-1
   Account ID: 123456789012 (Optional)
   ```
4. Click "Test Connection"
5. If successful, click "Save & Sync"

#### Azure Connection

**Step 1: Create Service Principal**

```bash
# Login to Azure
az login

# Create service principal
az ad sp create-for-rbac \
  --name prodory-finops \
  --role "Cost Management Reader" \
  --scopes /subscriptions/YOUR_SUBSCRIPTION_ID

# Output will include:
# - appId (Client ID)
# - password (Client Secret)
# - tenant (Tenant ID)
```

**Step 2: Grant Billing Access**

```bash
# Grant billing access for EA accounts
az role assignment create \
  --assignee APP_ID \
  --role "Billing Reader" \
  --scope /providers/Microsoft.Billing/billingAccounts/YOUR_BILLING_ACCOUNT
```

**Step 3: Add to Prodory Platform**

1. Navigate to **Settings → Cloud Providers → Add Provider**
2. Select "Azure"
3. Enter credentials:
   ```yaml
   Provider Name: Azure Production
   Client ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
   Client Secret: ********
   Tenant ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
   Subscription ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
   ```

#### GCP Connection

**Step 1: Create Service Account**

```bash
# Create service account
gcloud iam service-accounts create prodory-finops \
  --display-name "Prodory FinOps"

# Grant billing viewer role
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member "serviceAccount:prodory-finops@PROJECT_ID.iam.gserviceaccount.com" \
  --role "roles/billing.viewer"

# Grant BigQuery data viewer (for detailed costs)
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member "serviceAccount:prodory-finops@PROJECT_ID.iam.gserviceaccount.com" \
  --role "roles/bigquery.dataViewer"

# Create and download key
gcloud iam service-accounts keys create prodory-gcp-key.json \
  --iam-account prodory-finops@PROJECT_ID.iam.gserviceaccount.com
```

**Step 2: Add to Prodory Platform**

1. Navigate to **Settings → Cloud Providers → Add Provider**
2. Select "GCP"
3. Upload the JSON key file or paste contents
4. Enter Project ID

### Data Synchronization

#### Sync Schedule

| Data Type | Frequency | Retention |
|-----------|-----------|-----------|
| Cost Data | Hourly | 24 months |
| Resource Inventory | Every 6 hours | Current only |
| Usage Metrics | Every hour | 90 days |
| Recommendations | Daily | 12 months |

#### Manual Sync Trigger

```
Settings → Cloud Providers → [Provider] → Actions → Sync Now
```

#### Sync Status Monitoring

```
┌─────────────────────────────────────────────────────────────┐
│ Data Sync Status                                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ AWS Production                    Last Sync: 23 mins ago   │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ✅ Healthy    │
│                                                             │
│ Azure Production                  Last Sync: 1 hour ago    │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ✅ Healthy    │
│                                                             │
│ GCP Production                    Last Sync: 45 mins ago   │
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ⚠️ Warning    │
│   Note: Some BigQuery tables not accessible                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Data Quality & Validation

#### Data Validation Checks

The agent performs automatic validation:

| Check | Description | Action on Failure |
|-------|-------------|-------------------|
| Completeness | All expected data received | Retry with backoff |
| Accuracy | Costs match provider totals | Flag for review |
| Timeliness | Data within expected timeframe | Alert operators |
| Consistency | No duplicate or missing records | Automatic dedupe |

#### Handling Data Discrepancies

If you notice discrepancies:

1. **Check Sync Status**
   ```
   Settings → Cloud Providers → [Provider] → Sync History
   ```

2. **Compare with Provider Console**
   - AWS: Cost Explorer
   - Azure: Cost Management
   - GCP: Cloud Billing

3. **Request Data Reconciliation**
   ```
   Settings → Cloud Providers → [Provider] → Actions → Reconcile Data
   ```

---

## Kubernetes-in-a-Box

### Overview

Kubernetes-in-a-Box provides simplified Kubernetes cluster deployment and management for development, testing, and production workloads.

### Creating a New Cluster

#### Step 1: Launch Cluster Wizard

```
Kubernetes-in-a-Box → Create Cluster
```

#### Step 2: Select Cluster Type

| Type | Use Case | Specifications |
|------|----------|---------------|
| Development | Local development, testing | 1-3 nodes, 4GB RAM/node |
| Testing | CI/CD, integration tests | 3 nodes, 8GB RAM/node |
| Staging | Pre-production validation | 3-5 nodes, 16GB RAM/node |
| Production | Production workloads | 5+ nodes, 32GB RAM/node |
| GPU | ML/AI workloads | GPU-enabled nodes |

#### Step 3: Configure Cluster

```yaml
Cluster Name: prod-cluster-01
Environment: Production
Region: us-east-1

Node Configuration:
  Master Nodes: 3 (HA configuration)
  Worker Nodes: 5
  
  Master Node Specs:
    Instance Type: c5.2xlarge
    CPU: 8 cores
    Memory: 16 GB
    Storage: 100 GB SSD
    
  Worker Node Specs:
    Instance Type: m5.xlarge
    CPU: 4 cores
    Memory: 16 GB
    Storage: 200 GB SSD

Networking:
  VPC: Create new (10.0.0.0/16)
  Pod CIDR: 10.244.0.0/16
  Service CIDR: 10.96.0.0/12
  
Add-ons:
  - Ingress Controller: NGINX
  - Monitoring: Prometheus + Grafana
  - Logging: Fluent Bit + Elasticsearch
  - Service Mesh: Istio (optional)
  - Cert Manager: Enabled
```

#### Step 4: Review and Create

```
┌─────────────────────────────────────────────────────────────┐
│ Cluster Configuration Summary                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Name: prod-cluster-01                                      │
│ Environment: Production                                    │
│ Region: us-east-1                                          │
│                                                             │
│ Nodes:                                                     │
│   - 3 Master (c5.2xlarge)                                  │
│   - 5 Worker (m5.xlarge)                                   │
│                                                             │
│ Estimated Cost: $1,250/month                               │
│ Creation Time: ~15 minutes                                 │
│                                                             │
│ [Create Cluster]  [Save as Template]  [Cancel]            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Managing Clusters

#### Cluster Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│ prod-cluster-01                              [Actions ▼]   │
│ Status: ✅ Running    Uptime: 45 days                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Resources                    Health                        │
│  ┌─────────────────────┐     ┌─────────────────────┐       │
│  │ Nodes: 8/8          │     │ API Server: ✅      │       │
│  │ Pods: 127/500       │     │ etcd: ✅            │       │
│  │ CPU: 45%            │     │ Scheduler: ✅       │       │
│  │ Memory: 62%         │     │ Controller: ✅      │       │
│  │ Storage: 38%        │     │ DNS: ✅             │       │
│  └─────────────────────┘     └─────────────────────┘       │
│                                                             │
│  Recent Events                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 10:23 - Pod nginx-7d8c9b5f-x2a4 started            │   │
│  │ 09:45 - Node worker-3 joined cluster               │   │
│  │ 08:12 - HorizontalPodAutoscaler scaled deployment  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Scaling Operations

**Scale Nodes:**
```bash
# Via UI
Cluster → Actions → Scale Nodes

# Via CLI
prodory k8s scale prod-cluster-01 --workers 8
```

**Scale Workloads:**
```yaml
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

#### Accessing Clusters

**Download Kubeconfig:**
```
Cluster → Actions → Download Kubeconfig
```

**Use with kubectl:**
```bash
# Set KUBECONFIG environment variable
export KUBECONFIG=~/Downloads/prod-cluster-01-kubeconfig.yaml

# Verify connection
kubectl get nodes
kubectl get pods --all-namespaces
```

**Web Terminal:**
```
Cluster → Terminal
# Opens in-browser kubectl shell
```

### Application Deployment

#### Deploy from UI

1. Navigate to **Cluster → Applications → Deploy**
2. Choose deployment method:
   - **Helm Chart**: Select from catalog
   - **YAML Manifest**: Paste or upload
   - **Container Image**: Specify image and settings

#### Deploy from Helm Catalog

```
Applications → Deploy → Helm Chart

Select Chart:
  [x] nginx-ingress
  [ ] cert-manager
  [ ] prometheus
  [ ] grafana
  [ ] elasticsearch
  [ ] kafka

Configure Values:
  controller.replicaCount: 2
  controller.service.type: LoadBalancer

Target Namespace: ingress-nginx
Create Namespace: Yes
```

#### Deploy Custom Application

```yaml
# Application Manifest
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-application
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-application
  template:
    metadata:
      labels:
        app: my-application
    spec:
      containers:
        - name: app
          image: myregistry/myapp:v1.2.3
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: my-application
  namespace: production
spec:
  selector:
    app: my-application
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-application
  namespace: production
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt"
spec:
  tls:
    - hosts:
        - myapp.company.com
      secretName: myapp-tls
  rules:
    - host: myapp.company.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-application
                port:
                  number: 80
```

### Monitoring & Logging

#### Built-in Monitoring

Access monitoring dashboards:
```
Cluster → Monitoring → Dashboards

Available Dashboards:
- Cluster Overview
- Node Metrics
- Pod Metrics
- Network Traffic
- Storage Usage
```

#### Log Access

```
Cluster → Logs

Search: error OR failed
Namespace: production
Time Range: Last 1 hour

Results:
[2024-01-15 10:23:45] [ERROR] Database connection failed
[2024-01-15 10:23:47] [ERROR] Retry attempt 1/3
```

---

## Storage Autoscaler

### Overview

Storage Autoscaler automatically manages storage capacity across your infrastructure, preventing outages from disk full conditions while optimizing costs.

### Supported Storage Types

| Type | Platforms | Auto-Scaling |
|------|-----------|--------------|
| EBS Volumes | AWS | Yes |
| Azure Disks | Azure | Yes |
| Persistent Disks | GCP | Yes |
| Local Storage | On-prem | Monitoring only |
| NFS/SMB | Any | Monitoring only |

### Setting Up Storage Monitoring

#### Step 1: Connect Storage Provider

```
Storage Autoscaler → Settings → Add Provider

Provider: AWS
Region: us-east-1
Credentials: [Select from Cloud Providers]
```

#### Step 2: Configure Auto-Scaling Policies

```yaml
Policy Name: Production Storage Policy

Scaling Rules:
  Scale Up Trigger:
    - When disk usage > 80%
    - Increase size by: 20%
    - Minimum increment: 10 GB
    - Maximum size: 1000 GB
    - Cooldown period: 6 hours
    
  Scale Down Trigger:
    - When disk usage < 30% for 7 days
    - Decrease size to: 150% of used space
    - Minimum size: 20 GB
    - Cooldown period: 24 hours

Notifications:
  - Alert at 70%: Warning
  - Alert at 85%: Critical
  - Alert at 95%: Emergency

Cost Optimization:
  - Enable storage tiering: Yes
  - Move to cheaper tier after: 30 days
  - Compress old data: Yes
```

#### Step 3: Select Volumes to Monitor

```
Storage Autoscaler → Volumes → Select Volumes

[x] vol-0a1b2c3d (Production Database) - 500 GB
[x] vol-0e5f6g7h (Application Logs) - 200 GB  
[ ] vol-0i9j0k1l (Backup Storage) - 1000 GB
[x] vol-0m2n3o4p (User Uploads) - 300 GB

Apply Policy: Production Storage Policy
```

### Monitoring Storage Usage

#### Storage Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│ Storage Overview                                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Total Capacity: 5.2 TB                                     │
│  Used: 3.8 TB (73%)                                         │
│  Available: 1.4 TB                                          │
│  Monthly Cost: $420                                         │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Usage Trend (30 days)                               │   │
│  │                                                     │   │
│  │ 100% ┤                                    ╭────╮   │   │
│  │  80% ┤                          ╭────────╯    │   │   │
│  │  60% ┤              ╭───────────╯              │   │   │
│  │  40% ┤  ╭──────────╯                           │   │   │
│  │  20% ┤──╯                                      │   │   │
│  │   0% ┼────┬────┬────┬────┬────┬────┬────┬────┤   │   │
│  │       W1   W2   W3   W4                              │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  Volumes Requiring Attention:                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ⚠️  vol-0a1b2c3d - 87% full (Action needed)        │   │
│  │ ✅ vol-0e5f6g7h - 45% full                         │   │
│  │ ⚠️  vol-0m2n3o4p - 82% full (Scaling scheduled)    │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Manual Storage Operations

#### Expand Volume Manually

```
Storage Autoscaler → Volumes → [Volume] → Actions → Expand

Current Size: 500 GB
New Size: 750 GB

[Confirm Expansion]

Note: Volume expansion may take 5-15 minutes.
Application restart may be required for some filesystems.
```

#### Create Snapshot

```
Storage Autoscaler → Volumes → [Volume] → Actions → Create Snapshot

Snapshot Name: pre-upgrade-backup-2024-01-15
Description: Backup before application upgrade
Retention: 30 days

[Create Snapshot]
```

---

## Cloud Sentinel

### Overview

Cloud Sentinel provides continuous security monitoring, compliance checking, and threat detection across your cloud infrastructure.

### Security Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│ Cloud Sentinel - Security Overview                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Security Score: 87/100                    [View Details]  │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━                  │
│                                                             │
│  ┌───────────────┐  ┌───────────────┐  ┌───────────────┐  │
│  │ 🔴 Critical   │  │ 🟠 High       │  │ 🟡 Medium     │  │
│  │      2        │  │      5        │  │     12        │  │
│  └───────────────┘  └───────────────┘  └───────────────┘  │
│                                                             │
│  Compliance Status                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ SOC 2:    ✅ Compliant (98%)                        │   │
│  │ ISO 27001: ⚠️  3 findings                           │   │
│  │ PCI DSS:  🔴 7 critical findings                    │   │
│  │ HIPAA:    ✅ Compliant (95%)                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  Recent Security Events                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 10:23 - Unauthorized API call detected (Blocked)   │   │
│  │ 09:45 - S3 bucket policy change detected           │   │
│  │ 08:12 - New IAM user created                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Compliance Frameworks

#### Supported Frameworks

| Framework | Description | Checks |
|-----------|-------------|--------|
| CIS Benchmarks | Center for Internet Security | 200+ |
| SOC 2 | Service Organization Control | 150+ |
| ISO 27001 | Information Security Standard | 180+ |
| PCI DSS | Payment Card Industry | 100+ |
| HIPAA | Healthcare Compliance | 80+ |
| GDPR | Data Protection | 70+ |
| NIST 800-53 | Security Controls | 250+ |
| Custom | Your own policies | Unlimited |

#### Running Compliance Scan

```
Cloud Sentinel → Compliance → Run Scan

Select Frameworks:
  [x] CIS AWS Foundations
  [x] SOC 2 Type II
  [ ] ISO 27001
  [x] PCI DSS

Scope:
  - AWS Account: Production
  - Azure Subscription: All
  - GCP Project: Production

Schedule:
  [ ] Run now
  [x] Schedule: Daily at 02:00 UTC

[Start Scan]
```

#### Reviewing Findings

```
┌─────────────────────────────────────────────────────────────┐
│ Compliance Finding Details                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ ID: CIS-1.4                                                 │
│ Severity: 🔴 Critical                                       │
│ Framework: CIS AWS Foundations v1.5                         │
│                                                             │
│ Title:                                                      │
│ Ensure access keys are rotated every 90 days               │
│                                                             │
│ Description:                                                │
│ Access keys should be rotated regularly to minimize the    │
│ impact of compromised credentials.                         │
│                                                             │
│ Affected Resources:                                         │
│ ┌─────────────────────────────────────────────────────┐    │
│ │ User              │ Access Key Age │ Last Rotated   │    │
│ ├───────────────────┼────────────────┼────────────────┤    │
│ │ deploy-user       │ 127 days       │ 2023-09-10     │    │
│ │ backup-service    │ 156 days       │ 2023-08-12     │    │
│ └─────────────────────────────────────────────────────┘    │
│                                                             │
│ Remediation:                                                │
│ 1. Create new access key for each user                     │
│ 2. Update applications to use new key                      │
│ 3. Disable old access key                                  │
│ 4. Delete old access key after 7 days                      │
│                                                             │
│ [Mark as Exception]  [Create Ticket]  [Auto-Remediate]    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Threat Detection

#### Real-time Alerts

```
Cloud Sentinel → Threat Detection → Alerts

Filter: All Open Alerts

┌─────────────────────────────────────────────────────────────┐
│ 🔴 CRITICAL: Unusual API Activity                          │
│ Time: 2024-01-15 10:23:45 UTC                              │
│ Account: AWS Production (123456789012)                     │
│                                                             │
│ Description:                                                │
│ Detected 1,247 DeleteBucket API calls in 5 minutes from    │
│ IP 203.0.113.45. Normal rate is <10 per hour.              │
│                                                             │
│ Affected Resources:                                         │
│ - s3://company-backups                                     │
│ - s3://customer-data-prod                                  │
│                                                             │
│ Recommended Actions:                                        │
│ 1. Review IAM user 'backup-script' activity                │
│ 2. Verify IP 203.0.113.45 is authorized                    │
│ 3. Consider revoking access if unauthorized                │
│                                                             │
│ [Investigate]  [Acknowledge]  [Escalate]  [Auto-Remediate]│
└─────────────────────────────────────────────────────────────┘
```

#### Custom Detection Rules

```yaml
# Custom Detection Rule
name: Unusual Database Access
severity: High
description: Alert when database is accessed from new IP

trigger:
  event: rds:Connect
  condition: |
    source_ip NOT IN (
      SELECT ip FROM approved_ips 
      WHERE service = 'rds'
    )

actions:
  - type: alert
    channels: [email, slack]
  - type: block
    duration: 1 hour
    require_approval: true
  - type: create_ticket
    system: jira
    project: SEC
```

---

## VMware to OpenShift Migration

### Overview

Migrate VMware virtual machines to OpenShift Virtualization (KubeVirt) for a unified container and VM platform.

### Migration Planning

#### Step 1: Discovery

```
VMware Migration → Discovery → Scan vCenter

vCenter Server: vcenter.company.com
Credentials: [Select from vault]

Scanning... Found 47 VMs

┌─────────────────────────────────────────────────────────────┐
│ Discovered VMs                                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Name              │ OS           │ CPU │ RAM   │ Disk │    │
├───────────────────┼──────────────┼─────┼───────┼──────┤    │
│ web-server-01     │ Ubuntu 22.04 │ 4   │ 8 GB  │ 100G │    │
│ db-server-01      │ RHEL 8       │ 8   │ 32 GB │ 500G │    │
│ app-server-01     │ Windows 2019 │ 4   │ 16 GB │ 200G │    │
│ ...               │ ...          │ ... │ ...   │ ...  │    │
│                                                             │
│ [Select All]  [Select Compatible]  [Import Selected]       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Step 2: Compatibility Assessment

```
VMware Migration → Assessment → Run Assessment

Selected VMs: 15

Assessment Results:
┌─────────────────────────────────────────────────────────────┐
│ Compatibility Report                                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ ✅ Compatible (12 VMs)                                      │
│   - All Linux VMs with standard kernels                    │
│   - No specialized hardware dependencies                   │
│                                                             │
│ ⚠️  Requires Changes (2 VMs)                               │
│   - app-server-03: Remove VMware Tools                     │
│   - db-server-02: Update network drivers                   │
│                                                             │
│ ❌ Not Compatible (1 VM)                                    │
│   - legacy-app-01: Requires specialized PCI device         │
│     Recommendation: Keep on VMware or refactor             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Step 3: Migration Plan

```yaml
Migration Plan: Production VMs Batch 1

Source: VMware vCenter (vcenter.company.com)
Target: OpenShift Cluster (prod-openshift)
Namespace: migrated-vms

VMs to Migrate:
  1. web-server-01
     - Downtime window: 2 hours
     - Cutover: Saturday 02:00 UTC
     - Rollback plan: Keep VM for 48 hours
     
  2. app-server-01
     - Downtime window: 4 hours
     - Cutover: Saturday 06:00 UTC
     - Rollback plan: Keep VM for 48 hours

Pre-Migration Tasks:
  - [ ] Notify stakeholders
  - [ ] Create VM snapshots
  - [ ] Configure OpenShift networking
  - [ ] Test target storage class

Post-Migration Tasks:
  - [ ] Verify VM functionality
  - [ ] Update DNS records
  - [ ] Configure monitoring
  - [ ] Decommission source VM after 48h
```

### Executing Migration

#### Warm Migration (Minimal Downtime)

```
VMware Migration → Execute → Start Warm Migration

Migration: web-server-01
Status: In Progress

Phase 1: Initial Sync
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 100%
  Full VM disk copied to OpenShift

Phase 2: Incremental Sync
  ━━━━━━━━━━━━━━━━━━━━━━━━░░░░░░░░░░░░░░░░ 45%
  Syncing changes since initial sync
  Estimated completion: 15 minutes

Phase 3: Final Sync & Cutover
  ⏸️  Waiting for maintenance window
  Scheduled: Saturday 02:00 UTC

[Pause]  [Cancel]  [Force Cutover Now]
```

#### Post-Migration Verification

```
VMware Migration → Verification → Run Checks

VM: web-server-01
Status: ✅ Migration Successful

Verification Results:
┌─────────────────────────────────────────────────────────────┐
│ Check                          │ Status │ Details           │
├────────────────────────────────┼────────┼───────────────────┤
│ VM Boot                        │ ✅     │ Booted in 45s     │
│ Network Connectivity           │ ✅     │ All interfaces up │
│ SSH/Remote Access              │ ✅     | Authentication OK │
│ Application Services           │ ✅     │ All services running│
│ Data Integrity                 │ ✅     │ Checksums match   │
│ Performance Baseline           | ✅     │ Within 5% of original│
└────────────────────────────────┴────────┴───────────────────┘

[Complete Migration]  [Rollback]  [Schedule Decommission]
```

---

## VM to Container Migration

### Overview

Modernize traditional virtual machines by migrating them to containerized workloads.

### Migration Assessment

#### Application Analysis

```
VM to Container → Assessment → Analyze VM

VM: app-server-01 (192.168.1.100)
Analysis Status: Complete

┌─────────────────────────────────────────────────────────────┐
│ Application Analysis Report                                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ Detected Applications:                                      │
│ ┌─────────────────────────────────────────────────────┐    │
│ │ Application    │ Type      │ Containerizable │ Effort│    │
│ ├────────────────┼───────────┼─────────────────┼───────┤    │
│ │ Nginx          │ Web Server│ ✅ Easy         │ Low   │    │
│ │ Node.js App    │ App Server│ ✅ Easy         │ Low   │    │
│ │ PostgreSQL     │ Database  │ ⚠️  Medium      │ Med   │    │
│ │ Redis          │ Cache     │ ✅ Easy         │ Low   │    │
│ │ Custom Perl    │ Legacy    │ ❌ Hard         │ High  │    │
│ └────────────────┴───────────┴─────────────────┴───────┘    │
│                                                             │
│ Migration Strategy Recommendation:                          │
│ - Containerize: Nginx, Node.js, Redis                      │
│ - Keep as VM: PostgreSQL (use managed service instead)     │
│ - Refactor: Custom Perl script (modernize to Python/Go)    │
│                                                             │
│ Estimated Effort: 2-3 weeks                                 │
│ Expected Benefits:                                          │
│ - 60% reduction in infrastructure costs                    │
│ - Improved deployment speed (minutes vs hours)             │
│ - Better resource utilization                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Containerization Process

#### Step 1: Generate Dockerfile

```dockerfile
# Generated Dockerfile for app-server-01
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application code
COPY . .

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Change ownership
RUN chown -R nodejs:nodejs /app
USER nodejs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

# Start application
CMD ["node", "server.js"]
```

#### Step 2: Generate Kubernetes Manifests

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-server-01
  labels:
    app: app-server-01
spec:
  replicas: 3
  selector:
    matchLabels:
      app: app-server-01
  template:
    metadata:
      labels:
        app: app-server-01
    spec:
      containers:
        - name: app
          image: registry.company.com/app-server-01:v1.0.0
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: "production"
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: database-url
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 3000
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: app-server-01
spec:
  selector:
    app: app-server-01
  ports:
    - port: 80
      targetPort: 3000
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-server-01
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt"
spec:
  tls:
    - hosts:
        - app-server-01.company.com
      secretName: app-server-01-tls
  rules:
    - host: app-server-01.company.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app-server-01
                port:
                  number: 80
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-server-01
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app-server-01
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

#### Step 3: Build and Deploy

```bash
# Build container image
docker build -t registry.company.com/app-server-01:v1.0.0 .

# Push to registry
docker push registry.company.com/app-server-01:v1.0.0

# Deploy to Kubernetes
kubectl apply -f k8s/

# Verify deployment
kubectl get pods -l app=app-server-01
kubectl get svc app-server-01
kubectl get ingress app-server-01
```

---

## Integration & Workflows

### Cross-Module Workflows

#### Cost Optimization Workflow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Data FinOps │────▶│  AI FinOps   │────▶│  Kubernetes  │
│    Agent     │     │  Dashboard   │     │   in-a-Box   │
└──────────────┘     └──────────────┘     └──────────────┘
       │                    │                    │
       │ Collect cost data  │ Analyze & recommend│ Apply changes
       │                    │                    │
       ▼                    ▼                    ▼
  Cloud APIs            AI Engine            Auto-scaling
```

**Workflow Steps:**

1. **Data Collection** (Data FinOps Agent)
   - Hourly sync from cloud providers
   - Store in time-series database

2. **Analysis** (AI FinOps Dashboard)
   - AI analyzes usage patterns
   - Generates recommendations
   - Creates cost forecasts

3. **Action** (Kubernetes-in-a-Box)
   - Rightsize workloads
   - Adjust resource requests/limits
   - Enable auto-scaling

#### Security & Compliance Workflow

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Cloud      │────▶│   Cloud      │────▶│  Kubernetes  │
│   Sentinel   │     │   Sentinel   │     │   in-a-Box   │
│  (Detect)    │     │  (Remediate) │     │  (Enforce)   │
└──────────────┘     └──────────────┘     └──────────────┘
```

**Workflow Steps:**

1. **Detection**
   - Continuous security scanning
   - Threat detection
   - Compliance checking

2. **Remediation**
   - Auto-remediate where safe
   - Create tickets for manual review
   - Alert security team

3. **Enforcement**
   - Apply security policies
   - Block non-compliant deployments
   - Maintain audit trail

### API Integration

#### REST API Usage

```bash
# Get authentication token
curl -X POST https://api.prodory.local/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "username": "api-user",
    "password": "your-password"
  }'

# Use token for API calls
TOKEN="your-jwt-token"

# Get cost summary
curl https://api.prodory.local/finops/costs/summary \
  -H "Authorization: Bearer $TOKEN"

# Get recommendations
curl https://api.prodory.local/finops/recommendations \
  -H "Authorization: Bearer $TOKEN"

# Apply recommendation
curl -X POST https://api.prodory.local/finops/recommendations/123/apply \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "dry_run": false,
    "scheduled_time": null
  }'
```

#### Webhook Integration

Configure webhooks for external notifications:

```yaml
# Webhook configuration
webhooks:
  - name: slack-alerts
    url: https://hooks.slack.com/services/YOUR/WEBHOOK/URL
    events:
      - budget.threshold_exceeded
      - security.critical_alert
      - recommendation.high_impact
    
  - name: jira-tickets
    url: https://jira.company.com/rest/api/2/issue
    auth:
      type: basic
      username: prodory-bot
      password: ${JIRA_API_TOKEN}
    events:
      - security.finding.detected
      - migration.failed
```

---

## Troubleshooting

### Common Issues & Solutions

#### Issue: Dashboard Not Loading

**Symptoms:**
- Blank page after login
- Spinning loader indefinitely
- Error messages in browser console

**Solutions:**

1. **Clear Browser Cache**
   ```
   Chrome: Ctrl+Shift+R (Windows) / Cmd+Shift+R (Mac)
   Firefox: Ctrl+F5
   Safari: Cmd+Option+R
   ```

2. **Check Browser Console**
   ```
   Press F12 → Console tab
   Look for red error messages
   ```

3. **Try Different Browser**
   - Supported: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+

4. **Check Network Connectivity**
   ```bash
   ping prodory.your-domain.com
   curl -I https://prodory.your-domain.com
   ```

#### Issue: Cloud Provider Sync Failed

**Symptoms:**
- "Last sync: Never" or old timestamp
- Missing cost data
- Error notifications

**Solutions:**

1. **Verify Credentials**
   ```
   Settings → Cloud Providers → [Provider] → Test Connection
   ```

2. **Check IAM Permissions**
   - AWS: Verify Cost Explorer access
   - Azure: Check Billing Reader role
   - GCP: Confirm billing.viewer permission

3. **Review Error Logs**
   ```
   Settings → Cloud Providers → [Provider] → View Logs
   ```

4. **Manual Sync**
   ```
   Settings → Cloud Providers → [Provider] → Actions → Sync Now
   ```

#### Issue: Recommendations Not Appearing

**Symptoms:**
- Empty recommendations list
- "No recommendations found" message

**Solutions:**

1. **Check Data Availability**
   - Minimum 7 days of data required
   - Verify cloud providers connected

2. **Enable AI Analysis**
   ```
   Settings → AI Configuration → Enable Analysis
   ```

3. **Review AI Service Status**
   ```
   Settings → System Status → AI Services
   ```

#### Issue: Kubernetes Cluster Creation Failed

**Symptoms:**
- Cluster stuck in "Creating" state
- Error during provisioning

**Solutions:**

1. **Check Resource Quotas**
   - Verify cloud provider quotas
   - Check available IPs in VPC

2. **Review Events**
   ```
   Kubernetes-in-a-Box → [Cluster] → Events
   ```

3. **Check Logs**
   ```
   Kubernetes-in-a-Box → [Cluster] → Logs
   ```

4. **Retry with Different Configuration**
   - Try different region
   - Use smaller instance types
   - Reduce node count

### Getting Support

#### Support Channels

| Channel | Response Time | Best For |
|---------|---------------|----------|
| Documentation | Immediate | Self-service help |
| Community Forum | 24-48 hours | General questions |
| Email Support | 4 hours | Technical issues |
| Phone Support | 1 hour | Critical issues |
| On-site Support | Scheduled | Complex deployments |

#### Creating a Support Ticket

When creating a ticket, include:

1. **Issue Summary**: Brief description
2. **Steps to Reproduce**: Detailed steps
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Screenshots**: Visual evidence
6. **Logs**: Relevant log excerpts
7. **Environment**: Version, browser, OS

---

## Best Practices

### Daily Operations

#### Morning Checklist

- [ ] Review dashboard for overnight alerts
- [ ] Check budget status
- [ ] Review new recommendations
- [ ] Verify all cloud providers synced

#### Weekly Review

- [ ] Analyze cost trends
- [ ] Review applied recommendations impact
- [ ] Check compliance status
- [ ] Review security findings
- [ ] Update budgets if needed

#### Monthly Tasks

- [ ] Generate executive reports
- [ ] Review Reserved Instance utilization
- [ ] Audit user access
- [ ] Review and update policies
- [ ] Plan for next month's budgets

### Security Best Practices

1. **Use Strong Passwords**
   - Minimum 12 characters
   - Mix of uppercase, lowercase, numbers, symbols
   - Use password manager

2. **Enable Multi-Factor Authentication**
   ```
   Settings → Security → Enable 2FA
   ```

3. **Regular Access Review**
   ```
   Settings → Users → Review Access
   - Remove inactive users
   - Update role assignments
   - Audit API keys
   ```

4. **Rotate Credentials**
   - Cloud provider keys: Every 90 days
   - API tokens: Every 180 days
   - Service account passwords: Every 90 days

5. **Monitor Audit Logs**
   ```
   Settings → Audit Logs → Review
   - Look for unusual activity
   - Track configuration changes
   ```

### Cost Optimization Best Practices

1. **Set Up Budgets Early**
   - Create budgets before spending grows
   - Set conservative thresholds
   - Use forecast alerts

2. **Review Recommendations Weekly**
   - Don't ignore low-impact recommendations
   - Small savings add up
   - Track actual vs predicted savings

3. **Tag Resources Consistently**
   ```yaml
   Required Tags:
     - Environment (prod/staging/dev)
     - Team/Owner
     - Project
     - Cost Center
   ```

4. **Use Reserved Capacity**
   - Identify steady-state workloads
   - Purchase RIs/SPs for 1-3 years
   - Monitor utilization

5. **Clean Up Regularly**
   - Delete unused resources
   - Remove old snapshots
   - Archive old logs

---

## Appendix

### Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `?` | Show keyboard shortcuts |
| `g d` | Go to Dashboard |
| `g f` | Go to FinOps |
| `g k` | Go to Kubernetes |
| `g s` | Go to Storage |
| `g c` | Go to Cloud Sentinel |
| `/` | Search |
| `Esc` | Close modal/cancel |

### Glossary

| Term | Definition |
|------|------------|
| FinOps | Cloud Financial Management practice |
| RI | Reserved Instances (AWS) |
| SP | Savings Plans (AWS) |
| CUD | Committed Use Discounts (GCP) |
| Rightsizing | Adjusting resource sizes to match usage |
| KubeVirt | Kubernetes virtualization solution |
| HPA | Horizontal Pod Autoscaler |
| PVC | Persistent Volume Claim |

### Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-01-15 | Initial release |

---

*For additional support, contact: support@prodory.local*
