"""Budgets API router"""

from fastapi import APIRouter, Path, Body, Query
from typing import List, Optional
from pydantic import BaseModel
from datetime import datetime

router = APIRouter()


class BudgetCreate(BaseModel):
    name: str
    amount: float
    period: str  # daily, weekly, monthly, quarterly, yearly
    alerts: List[int]  # Percentage thresholds
    filters: Optional[dict] = None


class BudgetUpdate(BaseModel):
    name: Optional[str] = None
    amount: Optional[float] = None
    alerts: Optional[List[int]] = None
    status: Optional[str] = None


@router.get("")
async def get_budgets(
    status: Optional[str] = Query(None, enum=["active", "exceeded", "warning"]),
    provider: Optional[str] = None
):
    """Get all budgets"""
    # Mock data - replace with database query
    return [
        {
            "id": "1",
            "name": "Production Infrastructure",
            "amount": 15000,
            "spent": 14200,
            "period": "monthly",
            "alerts": [80, 100],
            "status": "warning",
            "createdAt": datetime.utcnow().isoformat()
        },
        {
            "id": "2",
            "name": "Development Environment",
            "amount": 5000,
            "spent": 3200,
            "period": "monthly",
            "alerts": [80, 100],
            "status": "active",
            "createdAt": datetime.utcnow().isoformat()
        }
    ]


@router.post("")
async def create_budget(budget: BudgetCreate):
    """Create a new budget"""
    return {
        "id": "new-budget-id",
        "name": budget.name,
        "amount": budget.amount,
        "period": budget.period,
        "alerts": budget.alerts,
        "status": "active",
        "createdAt": datetime.utcnow().isoformat()
    }


@router.get("/{budget_id}")
async def get_budget(
    budget_id: str = Path(..., description="Budget ID")
):
    """Get budget details"""
    return {
        "id": budget_id,
        "name": "Production Infrastructure",
        "amount": 15000,
        "spent": 14200,
        "period": "monthly",
        "alerts": [80, 100],
        "status": "warning",
        "history": []
    }


@router.put("/{budget_id}")
async def update_budget(
    budget_id: str = Path(..., description="Budget ID"),
    budget: BudgetUpdate = Body(...)
):
    """Update a budget"""
    return {"message": "Budget updated", "id": budget_id}


@router.delete("/{budget_id}")
async def delete_budget(
    budget_id: str = Path(..., description="Budget ID")
):
    """Delete a budget"""
    return {"message": "Budget deleted", "id": budget_id}


@router.get("/{budget_id}/alerts")
async def get_budget_alerts(
    budget_id: str = Path(..., description="Budget ID"),
    limit: int = Query(10, ge=1, le=100)
):
    """Get alert history for a budget"""
    return [
        {
            "id": "alert-1",
            "message": "Budget reached 80%",
            "timestamp": datetime.utcnow().isoformat(),
            "severity": "warning"
        }
    ]
