"""Settings API router"""

from fastapi import APIRouter, Body
from pydantic import BaseModel
from typing import Optional, List

router = APIRouter()


class SettingsUpdate(BaseModel):
    notifications: Optional[dict] = None
    integrations: Optional[dict] = None
    preferences: Optional[dict] = None


@router.get("")
async def get_settings():
    """Get application settings"""
    return {
        "notifications": {
            "email": True,
            "slack": False,
            "webhook": False,
            "budgetAlerts": [80, 100],
            "dailySummary": True,
            "anomalyDetection": True
        },
        "integrations": {
            "aws": {"connected": True, "regions": ["us-east-1", "us-west-2"]},
            "azure": {"connected": False},
            "gcp": {"connected": False}
        },
        "preferences": {
            "currency": "USD",
            "dateFormat": "YYYY-MM-DD",
            "timezone": "UTC"
        }
    }


@router.put("")
async def update_settings(settings: SettingsUpdate):
    """Update application settings"""
    return {"message": "Settings updated", "settings": settings}


@router.get("/notifications")
async def get_notification_settings():
    """Get notification settings"""
    return {
        "channels": {
            "email": {"enabled": True, "address": "admin@prodory.com"},
            "slack": {"enabled": False, "webhook": None},
            "webhook": {"enabled": False, "url": None}
        },
        "alerts": {
            "budget": {"enabled": True, "thresholds": [80, 100]},
            "anomaly": {"enabled": True, "sensitivity": 2.5},
            "daily": {"enabled": True, "time": "09:00"},
            "weekly": {"enabled": False, "day": "monday", "time": "09:00"}
        }
    }


@router.put("/notifications")
async def update_notification_settings(settings: dict = Body(...)):
    """Update notification settings"""
    return {"message": "Notification settings updated"}
