"""
Data FinOps Agent - Main Application
AI-powered cloud cost optimization and analysis API
"""

from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from typing import Optional, List
import logging

from app.config import settings
from app.routers import (
    dashboard,
    costs,
    recommendations,
    budgets,
    reports,
    settings as settings_router,
    providers,
    health
)
from app.services.scheduler import start_scheduler, stop_scheduler

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan handler"""
    # Startup
    logger.info("Starting Data FinOps Agent...")
    await start_scheduler()
    yield
    # Shutdown
    logger.info("Shutting down Data FinOps Agent...")
    await stop_scheduler()


# Create FastAPI application
app = FastAPI(
    title="Data FinOps Agent",
    description="AI-powered cloud cost optimization and analysis API",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# Include routers
app.include_router(health.router, tags=["Health"])
app.include_router(dashboard.router, prefix="/dashboard", tags=["Dashboard"])
app.include_router(costs.router, prefix="/costs", tags=["Costs"])
app.include_router(recommendations.router, prefix="/recommendations", tags=["Recommendations"])
app.include_router(budgets.router, prefix="/budgets", tags=["Budgets"])
app.include_router(reports.router, prefix="/reports", tags=["Reports"])
app.include_router(settings_router.router, prefix="/settings", tags=["Settings"])
app.include_router(providers.router, prefix="/providers", tags=["Providers"])


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "name": "Data FinOps Agent",
        "version": "1.0.0",
        "status": "running",
        "docs": "/docs"
    }


@app.exception_handler(Exception)
async def global_exception_handler(request, exc):
    """Global exception handler"""
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={"detail": "Internal server error"}
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=8000,
        reload=settings.DEBUG
    )
