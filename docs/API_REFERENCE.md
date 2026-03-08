# Prodory Platform - API Reference

## Base URL

```
Production: https://api.aoplus.in
Local: http://localhost:8000
```

## Authentication

All API requests require authentication using a Bearer token:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" https://api.aoplus.in/dashboard
```

## Endpoints

### Health Checks

#### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "service": "data-finops-agent",
  "timestamp": "2024-01-30T10:00:00Z",
  "system": {
    "cpu_percent": 15.2,
    "memory_percent": 45.8,
    "disk_percent": 32.1
  }
}
```

#### GET /ready
Readiness check for Kubernetes.

**Response:**
```json
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "cache": "ok",
    "cloud_providers": "ok"
  }
}
```

---

### Dashboard

#### GET /dashboard
Get dashboard overview data.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider (aws, azure, gcp) |
| days | integer | No | Number of days for trend (default: 30) |

**Response:**
```json
{
  "metrics": {
    "totalSpend": 45230,
    "spendChange": 12.5,
    "forecastedSpend": 52100,
    "savings": 8500,
    "resources": 142,
    "alerts": 3
  },
  "costTrend": [
    {"date": "Jan", "cost": 38000},
    {"date": "Feb", "cost": 39500},
    {"date": "Mar", "cost": 41200}
  ],
  "serviceBreakdown": [
    {"name": "Compute", "value": 45},
    {"name": "Storage", "value": 25}
  ],
  "recommendationsCount": 8,
  "generatedAt": "2024-01-30T10:00:00Z"
}
```

#### GET /dashboard/metrics
Get detailed metrics.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider |
| metric_type | string | No | Type of metrics (all, cost, usage, efficiency) |

---

### Cost Analysis

#### GET /costs/analysis
Get detailed cost analysis.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider |
| service | string | No | Filter by service |
| start_date | datetime | No | Start date (ISO 8601) |
| end_date | datetime | No | End date (ISO 8601) |
| period | string | No | Period (7d, 30d, 90d, 1y) |
| group_by | string | No | Grouping (day, week, month) |

**Response:**
```json
{
  "dailyCosts": [
    {"date": "01", "cost": 1200},
    {"date": "02", "cost": 1350}
  ],
  "serviceCosts": [
    {"service": "EC2", "cost": 5200, "change": 12},
    {"service": "S3", "cost": 1800, "change": -5}
  ],
  "period": {
    "start": "2024-01-01T00:00:00Z",
    "end": "2024-01-30T00:00:00Z"
  },
  "filters": {
    "provider": "aws",
    "service": null
  }
}
```

#### GET /costs/forecast
Get AI-powered cost forecast.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider |
| horizon_days | integer | No | Forecast horizon (7-365, default: 30) |
| confidence_interval | float | No | Confidence interval (0.8-0.99, default: 0.95) |

**Response:**
```json
{
  "forecast": [
    {"date": "2024-02-01", "forecast": 1520, "lower": 1400, "upper": 1640},
    {"date": "2024-02-02", "forecast": 1530, "lower": 1410, "upper": 1650}
  ],
  "summary": {
    "totalForecast": 45900,
    "averageDaily": 1530,
    "trend": "increasing",
    "trendPercentage": 8.5,
    "confidence": 0.95
  },
  "insights": [
    "Costs are expected to increase by 8.5% over the next 30 days"
  ],
  "generatedAt": "2024-01-30T10:00:00Z"
}
```

#### GET /costs/anomalies
Detect cost anomalies.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider |
| sensitivity | float | No | Anomaly sensitivity (1.0-5.0, default: 2.5) |
| days | integer | No | Analysis period (7-90, default: 30) |

**Response:**
```json
[
  {
    "date": "2024-01-15",
    "expected": 1200,
    "actual": 2800,
    "difference": 1600,
    "severity": "high"
  }
]
```

---

### Recommendations

#### GET /recommendations
Get AI-generated cost optimization recommendations.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | No | Filter by cloud provider |
| category | string | No | Filter by category |
| impact | string | No | Filter by impact (high, medium, low) |
| status | string | No | Filter by status (pending, applied, dismissed) |

**Response:**
```json
[
  {
    "id": "rec-1",
    "title": "Rightsize EC2 instances",
    "description": "5 instances are running at average 15% CPU...",
    "category": "Compute",
    "impact": "high",
    "effort": "low",
    "savings": 1250,
    "status": "pending",
    "resources": ["i-0a1b2c3d", "i-1b2c3d4e"],
    "confidence": 0.92
  }
]
```

#### POST /recommendations/{id}/apply
Apply a recommendation.

**Request Body:**
```json
{
  "dry_run": true,
  "notes": "Scheduled maintenance window"
}
```

**Response:**
```json
{
  "success": true,
  "dryRun": true,
  "message": "Dry run completed successfully",
  "estimatedSavings": 1250,
  "affectedResources": 5
}
```

---

### Budgets

#### GET /budgets
Get all budgets.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| status | string | No | Filter by status (active, exceeded, warning) |
| provider | string | No | Filter by cloud provider |

**Response:**
```json
[
  {
    "id": "1",
    "name": "Production Infrastructure",
    "amount": 15000,
    "spent": 14200,
    "period": "monthly",
    "alerts": [80, 100],
    "status": "warning",
    "createdAt": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /budgets
Create a new budget.

**Request Body:**
```json
{
  "name": "Development Team",
  "amount": 5000,
  "period": "monthly",
  "alerts": [80, 100],
  "filters": {
    "tags": {
      "Environment": "dev"
    }
  }
}
```

**Response:**
```json
{
  "id": "new-budget-id",
  "name": "Development Team",
  "amount": 5000,
  "period": "monthly",
  "alerts": [80, 100],
  "status": "active",
  "createdAt": "2024-01-30T10:00:00Z"
}
```

---

### Reports

#### GET /reports
Get all reports.

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| type | string | No | Filter by report type |
| status | string | No | Filter by status |

**Response:**
```json
[
  {
    "id": "1",
    "name": "Monthly Cost Summary - January 2024",
    "type": "cost-summary",
    "createdAt": "2024-02-01",
    "size": "2.4 MB",
    "status": "ready"
  }
]
```

#### POST /reports
Generate a new report.

**Request Body:**
```json
{
  "type": "cost-summary",
  "params": {
    "startDate": "2024-01-01",
    "endDate": "2024-01-31",
    "provider": "aws"
  },
  "schedule": null
}
```

**Response:**
```json
{
  "id": "new-report-id",
  "name": "cost-summary Report",
  "type": "cost-summary",
  "status": "generating",
  "createdAt": "2024-01-30T10:00:00Z"
}
```

---

### Cloud Providers

#### GET /providers
Get all connected cloud providers.

**Response:**
```json
[
  {
    "id": "aws-1",
    "provider": "aws",
    "name": "AWS Production",
    "status": "connected",
    "lastSync": "2024-01-30T10:00:00Z",
    "accounts": ["123456789012"],
    "regions": ["us-east-1", "us-west-2"]
  }
]
```

#### POST /providers/connect
Connect a cloud provider.

**Request Body (AWS):**
```json
{
  "provider": "aws",
  "credentials": {
    "accessKeyId": "AKIA...",
    "secretAccessKey": "...",
    "region": "us-east-1"
  }
}
```

**Response:**
```json
{
  "message": "AWS connected successfully",
  "provider": "aws",
  "status": "connected"
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "detail": "Invalid request parameters",
  "errors": {
    "field": "error message"
  }
}
```

### 401 Unauthorized
```json
{
  "detail": "Authentication required"
}
```

### 403 Forbidden
```json
{
  "detail": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "detail": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "detail": "Internal server error"
}
```

---

## Rate Limiting

API requests are rate limited:
- 1000 requests per hour per API key
- 100 requests per minute per IP

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1706611200
```

---

## Pagination

List endpoints support pagination:

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number |
| limit | integer | 50 | Items per page (max: 500) |

**Response:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 150,
    "pages": 3
  }
}
```
