# Docker Hub Setup Guide

This guide explains how to set up GitHub Actions to automatically build and push container images to Docker Hub.

## Prerequisites

- Docker Hub account (https://hub.docker.com)
- GitHub repository with the Prodory Platform code

## Step 1: Create Docker Hub Repository

### Option A: Create via Docker Hub Web UI

1. Log in to https://hub.docker.com
2. Click "Create Repository"
3. Create repositories for each service:
   - `prodory/finops-dashboard`
   - `prodory/data-finops-agent`
   - `prodory/k8s-in-a-box`
   - `prodory/storage-autoscaler`
   - `prodory/cloud-sentinel`
   - `prodory/vmware-migration`
   - `prodory/vm-to-container`

### Option B: Create via Docker CLI

```bash
# Login to Docker Hub
docker login

# Create repositories (using Docker Hub API)
for repo in finops-dashboard data-finops-agent k8s-in-a-box storage-autoscaler cloud-sentinel vmware-migration vm-to-container; do
  curl -X POST \
    -H "Content-Type: application/json" \
    -u "your-username:your-password" \
    -d '{"name":"'$repo'","is_private":false}' \
    https://hub.docker.com/v2/repositories/prodory/
done
```

## Step 2: Create Docker Hub Access Token

1. Log in to https://hub.docker.com
2. Click your profile → Account Settings
3. Go to "Security" tab
4. Click "New Access Token"
5. Name: `GitHub Actions`
6. Permissions: `Read, Write, Delete`
7. Click "Generate"
8. **Copy the token immediately** (you won't see it again!)

## Step 3: Add Secrets to GitHub Repository

1. Go to your GitHub repository
2. Click "Settings" → "Secrets and variables" → "Actions"
3. Click "New repository secret"
4. Add the following secrets:

| Secret Name | Value | Description |
|-------------|-------|-------------|
| `DOCKERHUB_USERNAME` | Your Docker Hub username | Docker Hub login |
| `DOCKERHUB_TOKEN` | Your access token | From Step 2 |

### Adding Secrets via GitHub CLI

```bash
# Install GitHub CLI if not already installed
# https://cli.github.com/

# Login to GitHub
gh auth login

# Set secrets
gh secret set DOCKERHUB_USERNAME --body "your-dockerhub-username"
gh secret set DOCKERHUB_TOKEN --body "your-dockerhub-token"

# Verify secrets
gh secret list
```

## Step 4: Verify Workflow Files

Ensure the following files exist in your repository:

```
.github/workflows/
├── build-and-push.yml    # Automatic builds on push
├── build-pr.yml          # Build validation on PRs
└── manual-build.yml      # Manual trigger builds
```

## Step 5: Test the Setup

### Test 1: Push to Main Branch

```bash
# Make a small change
echo "# Test" >> README.md

# Commit and push
git add .
git commit -m "Test Docker Hub push"
git push origin main
```

Check the Actions tab in GitHub - you should see the workflow running.

### Test 2: Create a Tag

```bash
# Create a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will trigger the workflow with semantic version tags.

### Test 3: Manual Build

1. Go to GitHub repository → Actions tab
2. Click "Manual Build and Push" workflow
3. Click "Run workflow"
4. Select:
   - Service: `all` (or specific service)
   - Tag: (leave empty for auto)
   - Push to registry: ✅
   - Platforms: `linux/amd64,linux/arm64`
5. Click "Run workflow"

## Image Tagging Strategy

The workflows use the following tagging strategy:

| Event | Tags Created | Example |
|-------|--------------|---------|
| Push to `main` | `latest`, `sha-abc1234` | `prodory/finops-dashboard:latest` |
| Push to `develop` | `develop`, `sha-abc1234` | `prodory/finops-dashboard:develop` |
| Tag `v1.2.3` | `1.2.3`, `1.2`, `1`, `latest` | `prodory/finops-dashboard:1.2.3` |
| Pull Request | `pr-123` | `prodory/finops-dashboard:pr-123` |
| Manual | Custom or timestamp | `prodory/finops-dashboard:manual-20240115-120000` |

## Pulling Images

Once built, you can pull images from Docker Hub:

```bash
# Pull latest
docker pull prodory/finops-dashboard:latest
docker pull prodory/data-finops-agent:latest
docker pull prodory/k8s-in-a-box:latest
docker pull prodory/storage-autoscaler:latest
docker pull prodory/cloud-sentinel:latest
docker pull prodory/vmware-migration:latest
docker pull prodory/vm-to-container:latest

# Pull specific version
docker pull prodory/finops-dashboard:v1.0.0

# Pull for specific architecture
docker pull --platform linux/arm64 prodory/finops-dashboard:latest
```

## Using Images in docker-compose.yml

Update your `docker-compose.yml` to use the published images:

```yaml
version: '3.8'

services:
  finops-dashboard:
    image: prodory/finops-dashboard:latest
    ports:
      - "3000:80"
    environment:
      - REACT_APP_API_URL=http://localhost:8000

  data-finops-agent:
    image: prodory/data-finops-agent:latest
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql://prodory:password@postgres:5432/prodory

  k8s-in-a-box:
    image: prodory/k8s-in-a-box:latest
    ports:
      - "8080:8080"

  storage-autoscaler:
    image: prodory/storage-autoscaler:latest
    ports:
      - "8001:8000"

  cloud-sentinel:
    image: prodory/cloud-sentinel:latest
    ports:
      - "8081:8080"

  vmware-migration:
    image: prodory/vmware-migration:latest
    ports:
      - "3001:3000"

  vm-to-container:
    image: prodory/vm-to-container:latest
    ports:
      - "3002:3000"
```

## Using Images in Kubernetes

Update your Kubernetes manifests:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: finops-dashboard
spec:
  replicas: 2
  selector:
    matchLabels:
      app: finops-dashboard
  template:
    metadata:
      labels:
        app: finops-dashboard
    spec:
      containers:
        - name: finops-dashboard
          image: prodory/finops-dashboard:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 80
```

## Troubleshooting

### Issue: "denied: requested access to the resource is denied"

**Cause:** Invalid credentials or insufficient permissions

**Solution:**
1. Verify `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets
2. Ensure token has "Read, Write, Delete" permissions
3. Check repository exists on Docker Hub

### Issue: "repository not found"

**Cause:** Repository doesn't exist on Docker Hub

**Solution:**
```bash
# Create the repository manually on Docker Hub
# Or use Docker Hub API to create it
```

### Issue: "no space left on device"

**Cause:** GitHub Actions runner disk full

**Solution:**
- Enable build caching (already configured)
- Use smaller base images
- Clean up unnecessary layers

### Issue: Build takes too long

**Solution:**
- Build caching is enabled by default
- Use `cache-from` and `cache-to` in build-push-action
- Consider using GitHub Actions larger runners

## Security Best Practices

1. **Use Access Tokens, not passwords**
   - Tokens can be revoked independently
   - Tokens can have limited scope

2. **Rotate tokens regularly**
   - Set a reminder to rotate every 90 days
   - Update the `DOCKERHUB_TOKEN` secret

3. **Use private repositories for sensitive images**
   - Change `is_private` to `true` when creating repos

4. **Enable Docker Hub vulnerability scanning**
   - Go to repository → Settings → Security
   - Enable "Scan images for vulnerabilities"

## Monitoring Builds

### GitHub Actions Dashboard

View build status at:
```
https://github.com/YOUR_ORG/prodory-platform/actions
```

### Docker Hub Repository

View pushed images at:
```
https://hub.docker.com/r/prodory/finops-dashboard/tags
```

### Subscribe to Build Notifications

1. Go to repository Settings → Notifications
2. Add webhook for build events
3. Or use GitHub Actions to send notifications

## Advanced Configuration

### Multi-Registry Push

To push to multiple registries (Docker Hub + ECR + GCR):

```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}

- name: Login to Amazon ECR
  uses: aws-actions/amazon-ecr-login@v2

- name: Build and push
  uses: docker/build-push-action@v5
  with:
    push: true
    tags: |
      prodory/finops-dashboard:latest
      ${{ steps.login-ecr.outputs.registry }}/finops-dashboard:latest
```

### Build Matrix for Multiple Services

For faster builds, use a matrix strategy:

```yaml
strategy:
  matrix:
    service:
      - finops-dashboard
      - data-finops-agent
      - k8s-in-a-box
      - storage-autoscaler
      - cloud-sentinel
      - vmware-migration
      - vm-to-container
```

## Support

For issues with GitHub Actions:
- https://github.com/features/actions
- https://docs.github.com/en/actions

For issues with Docker Hub:
- https://hub.docker.com/support
- https://docs.docker.com/docker-hub/
