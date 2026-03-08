"""Dashboard API router"""

from fastapi import APIRouter, Depends, Query
from typing import Optional
from datetime import datetime, timedelta

from app.services.cost_analyzer import CostAnalyzer
from app.services.ai_recommender import AIRecommender

router = APIRouter()


@router.get("")
async def get_dashboard_data(
    provider: Optional[str] = Query(None, description="Filter by cloud provider"),
    days: int = Query(30, ge=1, le=365, description="Number of days for trend")
):
    """Get dashboard overview data"""
    
    analyzer = CostAnalyzer()
    recommender = AIRecommender()
    
    # Get metrics
    metrics = await analyzer.get_summary_metrics(provider)
    
    # Get cost trend
    end_date = datetime.utcnow()
    start_date = end_date - timedelta(days=days)
    cost_trend = await analyzer.get_cost_trend(start_date, end_date, provider)
    
    # Get service breakdown
    service_breakdown = await analyzer.get_service_breakdown(provider)
    
    # Get AI recommendations count
    recommendations = await recommender.get_recommendations(provider)
    
    return {
        "metrics": metrics,
        "costTrend": cost_trend,
        "serviceBreakdown": service_breakdown,
        "recommendationsCount": len([r for r in recommendations if r["status"] == "pending"]),
        "generatedAt": datetime.utcnow().isoformat()
    }


@router.get("/metrics")
async def get_metrics(
    provider: Optional[str] = Query(None),
    metric_type: str = Query("all", enum=["all", "cost", "usage", "efficiency"])
):
    """Get detailed metrics"""
    analyzer = CostAnalyzer()
    return await analyzer.get_detailed_metrics(provider, metric_type)


@router.get("/alerts")
async def get_active_alerts(
    severity: Optional[str] = Query(None, enum=["critical", "warning", "info"]),
    limit: int = Query(10, ge=1, le=100)
):
    """Get active alerts"""
    analyzer = CostAnalyzer()
    return await analyzer.get_alerts(severity, limit)
