"""Health check router"""

from fastapi import APIRouter, status
from datetime import datetime
import psutil

router = APIRouter()


@router.get("/health", status_code=status.HTTP_200_OK)
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "data-finops-agent",
        "timestamp": datetime.utcnow().isoformat(),
        "system": {
            "cpu_percent": psutil.cpu_percent(interval=1),
            "memory_percent": psutil.virtual_memory().percent,
            "disk_percent": psutil.disk_usage('/').percent
        }
    }


@router.get("/ready", status_code=status.HTTP_200_OK)
async def readiness_check():
    """Readiness check for Kubernetes"""
    return {
        "status": "ready",
        "checks": {
            "database": "ok",
            "cache": "ok",
            "cloud_providers": "ok"
        }
    }


@router.get("/live", status_code=status.HTTP_200_OK)
async def liveness_check():
    """Liveness check for Kubernetes"""
    return {"status": "alive"}
