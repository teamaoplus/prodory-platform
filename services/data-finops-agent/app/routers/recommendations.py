"""Recommendations API router"""

from fastapi import APIRouter, Path, Body
from typing import List, Optional
from pydantic import BaseModel

from app.services.ai_recommender import AIRecommender

router = APIRouter()


class RecommendationResponse(BaseModel):
    id: str
    title: str
    description: str
    category: str
    impact: str
    effort: str
    savings: float
    status: str
    resources: List[str]


class ApplyRecommendationRequest(BaseModel):
    dry_run: bool = True
    notes: Optional[str] = None


@router.get("", response_model=List[RecommendationResponse])
async def get_recommendations(
    provider: Optional[str] = None,
    category: Optional[str] = None,
    impact: Optional[str] = None,
    status: Optional[str] = "pending"
):
    """Get AI-generated cost optimization recommendations"""
    recommender = AIRecommender()
    return await recommender.get_recommendations(provider, category, impact, status)


@router.get("/{recommendation_id}")
async def get_recommendation_detail(
    recommendation_id: str = Path(..., description="Recommendation ID")
):
    """Get detailed information about a recommendation"""
    recommender = AIRecommender()
    return await recommender.get_recommendation_detail(recommendation_id)


@router.post("/{recommendation_id}/apply")
async def apply_recommendation(
    recommendation_id: str = Path(..., description="Recommendation ID"),
    request: ApplyRecommendationRequest = Body(...)
):
    """Apply a recommendation"""
    recommender = AIRecommender()
    return await recommender.apply_recommendation(recommendation_id, request.dry_run, request.notes)


@router.post("/{recommendation_id}/dismiss")
async def dismiss_recommendation(
    recommendation_id: str = Path(..., description="Recommendation ID"),
    reason: Optional[str] = Body(None)
):
    """Dismiss a recommendation"""
    recommender = AIRecommender()
    return await recommender.dismiss_recommendation(recommendation_id, reason)


@router.post("/generate")
async def generate_recommendations(
    provider: Optional[str] = None,
    force: bool = False
):
    """Manually trigger recommendation generation"""
    recommender = AIRecommender()
    return await recommender.generate_recommendations(provider, force)
