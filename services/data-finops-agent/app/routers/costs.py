"""Costs API router"""

from fastapi import APIRouter, Query
from typing import Optional, List
from datetime import datetime, timedelta

from app.services.cost_analyzer import CostAnalyzer

router = APIRouter()


@router.get("/analysis")
async def get_cost_analysis(
    provider: Optional[str] = Query(None),
    service: Optional[str] = Query(None),
    start_date: Optional[datetime] = Query(None),
    end_date: Optional[datetime] = Query(None),
    period: str = Query("30d", regex="^(7d|30d|90d|1y)$"),
    group_by: str = Query("day", enum=["day", "week", "month"])
):
    """Get detailed cost analysis"""
    analyzer = CostAnalyzer()
    
    # Calculate date range if not provided
    if not end_date:
        end_date = datetime.utcnow()
    if not start_date:
        days_map = {"7d": 7, "30d": 30, "90d": 90, "1y": 365}
        start_date = end_date - timedelta(days=days_map.get(period, 30))
    
    daily_costs = await analyzer.get_daily_costs(start_date, end_date, provider, service, group_by)
    service_costs = await analyzer.get_service_costs(start_date, end_date, provider)
    
    return {
        "dailyCosts": daily_costs,
        "serviceCosts": service_costs,
        "period": {"start": start_date.isoformat(), "end": end_date.isoformat()},
        "filters": {"provider": provider, "service": service}
    }


@router.get("/by-resource")
async def get_costs_by_resource(
    provider: Optional[str] = Query(None),
    resource_type: Optional[str] = Query(None),
    limit: int = Query(50, ge=1, le=500)
):
    """Get costs grouped by resource"""
    analyzer = CostAnalyzer()
    return await analyzer.get_costs_by_resource(provider, resource_type, limit)


@router.get("/by-tag")
async def get_costs_by_tag(
    tag_key: str = Query(..., description="Tag key to group by"),
    provider: Optional[str] = Query(None)
):
    """Get costs grouped by tag"""
    analyzer = CostAnalyzer()
    return await analyzer.get_costs_by_tag(tag_key, provider)


@router.get("/forecast")
async def get_cost_forecast(
    provider: Optional[str] = Query(None),
    horizon_days: int = Query(30, ge=7, le=365),
    confidence_interval: float = Query(0.95, ge=0.8, le=0.99)
):
    """Get AI-powered cost forecast"""
    from app.services.forecaster import CostForecaster
    
    forecaster = CostForecaster()
    return await forecaster.forecast(provider, horizon_days, confidence_interval)


@router.get("/anomalies")
async def detect_anomalies(
    provider: Optional[str] = Query(None),
    sensitivity: float = Query(2.5, ge=1.0, le=5.0),
    days: int = Query(30, ge=7, le=90)
):
    """Detect cost anomalies"""
    analyzer = CostAnalyzer()
    return await analyzer.detect_anomalies(provider, sensitivity, days)
