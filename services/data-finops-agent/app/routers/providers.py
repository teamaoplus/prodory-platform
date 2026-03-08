"""Cloud Providers API router"""

from fastapi import APIRouter, Path, Body
from pydantic import BaseModel
from typing import Optional, Dict

router = APIRouter()


class AWSConnectionRequest(BaseModel):
    access_key_id: str
    secret_access_key: str
    region: str = "us-east-1"


class AzureConnectionRequest(BaseModel):
    subscription_id: str
    tenant_id: str
    client_id: str
    client_secret: str


class GCPConnectionRequest(BaseModel):
    project_id: str
    credentials_json: str


@router.get("")
async def get_connected_providers():
    """Get all connected cloud providers"""
    return [
        {
            "id": "aws-1",
            "provider": "aws",
            "name": "AWS Production",
            "status": "connected",
            "lastSync": "2024-01-30T10:00:00Z",
            "accounts": ["123456789012"],
            "regions": ["us-east-1", "us-west-2"]
        },
        {
            "id": "azure-1",
            "provider": "azure",
            "name": "Azure Subscription",
            "status": "disconnected",
            "lastSync": None,
            "accounts": [],
            "regions": []
        },
        {
            "id": "gcp-1",
            "provider": "gcp",
            "name": "GCP Project",
            "status": "disconnected",
            "lastSync": None,
            "accounts": [],
            "regions": []
        }
    ]


@router.post("/connect")
async def connect_provider(
    provider: str = Path(..., enum=["aws", "azure", "gcp"]),
    credentials: Dict = Body(...)
):
    """Connect a cloud provider"""
    return {
        "message": f"{provider.upper()} connected successfully",
        "provider": provider,
        "status": "connected"
    }


@router.post("/{provider_id}/sync")
async def sync_provider(
    provider_id: str = Path(..., description="Provider connection ID")
):
    """Manually sync provider data"""
    return {
        "message": "Sync started",
        "providerId": provider_id,
        "status": "syncing"
    }


@router.delete("/{provider_id}")
async def disconnect_provider(
    provider_id: str = Path(..., description="Provider connection ID")
):
    """Disconnect a cloud provider"""
    return {
        "message": "Provider disconnected",
        "providerId": provider_id
    }


@router.get("/{provider_id}/accounts")
async def get_provider_accounts(
    provider_id: str = Path(..., description="Provider connection ID")
):
    """Get accounts for a provider"""
    return [
        {
            "id": "123456789012",
            "name": "Production",
            "status": "active"
        },
        {
            "id": "098765432109",
            "name": "Development",
            "status": "active"
        }
    ]


@router.get("/{provider_id}/services")
async def get_provider_services(
    provider_id: str = Path(..., description="Provider connection ID")
):
    """Get available services for a provider"""
    return [
        {"id": "ec2", "name": "EC2", "category": "Compute"},
        {"id": "s3", "name": "S3", "category": "Storage"},
        {"id": "rds", "name": "RDS", "category": "Database"},
        {"id": "lambda", "name": "Lambda", "category": "Serverless"}
    ]
