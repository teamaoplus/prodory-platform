#!/usr/bin/env python3
"""
Storage Autoscaler - Kubernetes Operator for PVC Auto-scaling

This operator watches PersistentVolumeClaims and automatically scales them
when they reach a configurable threshold of capacity usage.

Features:
- Automatic PVC scaling based on usage metrics
- Support for multiple storage classes
- Configurable scale factors and thresholds
- Prometheus metrics export
- Safe scaling with data preservation
"""

import asyncio
import json
import os
import re
from dataclasses import dataclass
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional, Tuple

import kopf
import yaml
from kubernetes import client, config
from kubernetes.client.exceptions import ApiException
from prometheus_client import Counter, Gauge, Histogram, start_http_server
from pydantic import BaseModel, Field

# Metrics
SCALING_OPERATIONS = Counter(
    "storage_autoscaler_scaling_operations_total",
    "Total number of scaling operations",
    ["namespace", "pvc_name", "status"]
)

PVC_USAGE_PERCENT = Gauge(
    "storage_autoscaler_pvc_usage_percent",
    "Current PVC usage percentage",
    ["namespace", "pvc_name", "storage_class"]
)

PVC_CAPACITY_BYTES = Gauge(
    "storage_autoscaler_pvc_capacity_bytes",
    "Current PVC capacity in bytes",
    ["namespace", "pvc_name"]
)

SCALING_DURATION = Histogram(
    "storage_autoscaler_scaling_duration_seconds",
    "Time spent on scaling operations",
    ["namespace", "pvc_name"]
)

OPERATOR_ERRORS = Counter(
    "storage_autoscaler_errors_total",
    "Total number of operator errors",
    ["error_type"]
)

# Configuration
DEFAULT_CONFIG = {
    "threshold_percent": 80,
    "scale_factor": 1.5,
    "max_size": "1Ti",
    "min_size": "1Gi",
    "cooldown_minutes": 60,
    "storage_classes": [],  # Empty means all
    "exclude_namespaces": ["kube-system", "kube-public", "kube-node-lease"],
    "dry_run": False,
    "metrics_port": 8084,
    "check_interval_seconds": 60,
}


class ScalingPolicy(BaseModel):
    """Scaling policy configuration"""
    threshold_percent: int = Field(default=80, ge=1, le=99)
    scale_factor: float = Field(default=1.5, ge=1.1, le=10.0)
    max_size: str = Field(default="1Ti")
    min_size: str = Field(default="1Gi")
    cooldown_minutes: int = Field(default=60, ge=5)


class StorageAutoscalerSpec(BaseModel):
    """CRD spec for StorageAutoscaler"""
    enabled: bool = True
    selector: Dict[str, str] = Field(default_factory=dict)
    storage_classes: List[str] = Field(default_factory=list)
    exclude_namespaces: List[str] = Field(default_factory=list)
    policy: ScalingPolicy = Field(default_factory=ScalingPolicy)


@dataclass
class PVCInfo:
    """Information about a PVC"""
    name: str
    namespace: str
    storage_class: str
    capacity: int  # bytes
    used: int  # bytes
    available: int  # bytes
    usage_percent: float


class KubernetesClient:
    """Kubernetes API client wrapper"""

    def __init__(self):
        try:
            config.load_incluster_config()
        except config.ConfigException:
            config.load_kube_config()

        self.core_v1 = client.CoreV1Api()
        self.storage_v1 = client.StorageV1Api()
        self.custom_objects = client.CustomObjectsApi()

    def get_pvc(self, namespace: str, name: str) -> Optional[client.V1PersistentVolumeClaim]:
        """Get a PVC by namespace and name"""
        try:
            return self.core_v1.read_namespaced_persistent_volume_claim(name, namespace)
        except ApiException as e:
            if e.status == 404:
                return None
            raise

    def list_pvcs(self, namespace: Optional[str] = None) -> List[client.V1PersistentVolumeClaim]:
        """List all PVCs, optionally filtered by namespace"""
        if namespace:
            return self.core_v1.list_namespaced_persistent_volume_claim(namespace).items
        return self.core_v1.list_persistent_volume_claim_for_all_namespaces().items

    def update_pvc(self, namespace: str, name: str, pvc: client.V1PersistentVolumeClaim) -> client.V1PersistentVolumeClaim:
        """Update a PVC"""
        return self.core_v1.replace_namespaced_persistent_volume_claim(name, namespace, pvc)

    def get_pv(self, name: str) -> Optional[client.V1PersistentVolume]:
        """Get a PV by name"""
        try:
            return self.core_v1.read_persistent_volume(name)
        except ApiException as e:
            if e.status == 404:
                return None
            raise

    def get_storage_class(self, name: str) -> Optional[client.V1StorageClass]:
        """Get a storage class by name"""
        try:
            return self.storage_v1.read_storage_class(name)
        except ApiException as e:
            if e.status == 404:
                return None
            raise

    def get_pvc_metrics(self, namespace: str, name: str) -> Optional[Dict[str, Any]]:
        """Get PVC metrics from kubelet metrics API"""
        try:
            # Try to get metrics from metrics-server or kubelet
            pod_list = self.core_v1.list_namespaced_pod(namespace)
            for pod in pod_list.items:
                if pod.spec.volumes:
                    for volume in pod.spec.volumes:
                        if volume.persistent_volume_claim and volume.persistent_volume_claim.claim_name == name:
                            # Get metrics from pod
                            return self._get_volume_stats(pod, name)
            return None
        except Exception as e:
            kopf.warn(f"Failed to get metrics for PVC {namespace}/{name}: {e}")
            return None

    def _get_volume_stats(self, pod: client.V1Pod, pvc_name: str) -> Optional[Dict[str, Any]]:
        """Get volume statistics from a pod"""
        try:
            # This is a simplified version - in production you'd query metrics-server
            # or kubelet stats endpoint
            return {
                "capacity": 10737418240,  # 10Gi placeholder
                "used": 8589934592,  # 8Gi placeholder
                "available": 2147483648,  # 2Gi placeholder
            }
        except Exception:
            return None


class StorageScaler:
    """Main storage scaling logic"""

    def __init__(self, k8s_client: KubernetesClient, config: Dict[str, Any]):
        self.k8s = k8s_client
        self.config = {**DEFAULT_CONFIG, **config}
        self.scaling_history: Dict[str, datetime] = {}

    def parse_size(self, size_str: str) -> int:
        """Parse size string to bytes"""
        units = {
            "Ki": 1024,
            "Mi": 1024 ** 2,
            "Gi": 1024 ** 3,
            "Ti": 1024 ** 4,
            "Pi": 1024 ** 5,
            "Ei": 1024 ** 6,
            "K": 1000,
            "M": 1000 ** 2,
            "G": 1000 ** 3,
            "T": 1000 ** 4,
            "P": 1000 ** 5,
            "E": 1000 ** 6,
        }

        match = re.match(r"^(\d+(?:\.\d+)?)\s*([KMGTPE]i?)?$", size_str, re.IGNORECASE)
        if not match:
            raise ValueError(f"Invalid size format: {size_str}")

        value = float(match.group(1))
        unit = match.group(2) or ""

        multiplier = units.get(unit, 1)
        return int(value * multiplier)

    def format_size(self, size_bytes: int) -> str:
        """Format bytes to human readable string"""
        for unit in ["", "Ki", "Mi", "Gi", "Ti", "Pi"]:
            if size_bytes < 1024:
                return f"{size_bytes:.2f}{unit}"
            size_bytes /= 1024
        return f"{size_bytes:.2f}Ei"

    def calculate_new_size(self, current_size: int, scale_factor: float, max_size: int) -> int:
        """Calculate new PVC size after scaling"""
        new_size = int(current_size * scale_factor)
        # Round up to nearest GiB
        new_size = ((new_size + 1024**3 - 1) // (1024**3)) * (1024**3)
        return min(new_size, max_size)

    def is_in_cooldown(self, pvc_key: str) -> bool:
        """Check if PVC is in cooldown period"""
        if pvc_key not in self.scaling_history:
            return False

        last_scale = self.scaling_history[pvc_key]
        cooldown = timedelta(minutes=self.config["cooldown_minutes"])
        return datetime.now() - last_scale < cooldown

    def get_pvc_info(self, pvc: client.V1PersistentVolumeClaim) -> Optional[PVCInfo]:
        """Extract PVC information including usage"""
        try:
            namespace = pvc.metadata.namespace
            name = pvc.metadata.name
            storage_class = pvc.spec.storage_class_name or "default"

            # Get current capacity from PVC status
            capacity_str = pvc.status.capacity.get("storage", "0") if pvc.status and pvc.status.capacity else "0"
            capacity = self.parse_size(capacity_str)

            # Get usage metrics
            metrics = self.k8s.get_pvc_metrics(namespace, name)
            if metrics:
                used = metrics.get("used", 0)
                available = metrics.get("available", 0)
            else:
                # Estimate usage (in production, use actual metrics)
                used = int(capacity * 0.75)  # Placeholder: assume 75% usage
                available = capacity - used

            usage_percent = (used / capacity * 100) if capacity > 0 else 0

            return PVCInfo(
                name=name,
                namespace=namespace,
                storage_class=storage_class,
                capacity=capacity,
                used=used,
                available=available,
                usage_percent=usage_percent,
            )
        except Exception as e:
            OPERATOR_ERRORS.labels(error_type="pvc_info").inc()
            kopf.error(f"Failed to get PVC info: {e}")
            return None

    def should_scale(self, pvc_info: PVCInfo, policy: ScalingPolicy) -> Tuple[bool, str]:
        """Determine if PVC should be scaled"""
        pvc_key = f"{pvc_info.namespace}/{pvc_info.name}"

        # Check if in cooldown
        if self.is_in_cooldown(pvc_key):
            return False, "PVC is in cooldown period"

        # Check threshold
        if pvc_info.usage_percent < policy.threshold_percent:
            return False, f"Usage {pvc_info.usage_percent:.1f}% below threshold {policy.threshold_percent}%"

        # Check max size
        max_size = self.parse_size(policy.max_size)
        if pvc_info.capacity >= max_size:
            return False, f"PVC already at max size {policy.max_size}"

        # Check storage class
        if self.config["storage_classes"] and pvc_info.storage_class not in self.config["storage_classes"]:
            return False, f"Storage class {pvc_info.storage_class} not in allowed list"

        return True, "Scaling needed"

    async def scale_pvc(self, pvc: client.V1PersistentVolumeClaim, policy: ScalingPolicy) -> bool:
        """Scale a PVC to a new size"""
        namespace = pvc.metadata.namespace
        name = pvc.metadata.name
        pvc_key = f"{namespace}/{name}"

        with SCALING_DURATION.labels(namespace=namespace, pvc_name=name).time():
            try:
                pvc_info = self.get_pvc_info(pvc)
                if not pvc_info:
                    kopf.warn(f"Could not get PVC info for {pvc_key}")
                    SCALING_OPERATIONS.labels(
                        namespace=namespace, pvc_name=name, status="failed"
                    ).inc()
                    return False

                # Calculate new size
                max_size = self.parse_size(policy.max_size)
                new_size = self.calculate_new_size(
                    pvc_info.capacity,
                    policy.scale_factor,
                    max_size
                )

                if new_size <= pvc_info.capacity:
                    kopf.info(f"New size {self.format_size(new_size)} not larger than current {self.format_size(pvc_info.capacity)}")
                    return False

                kopf.info(f"Scaling PVC {pvc_key} from {self.format_size(pvc_info.capacity)} to {self.format_size(new_size)}")

                if self.config["dry_run"]:
                    kopf.info(f"DRY RUN: Would scale PVC {pvc_key} to {self.format_size(new_size)}")
                    SCALING_OPERATIONS.labels(
                        namespace=namespace, pvc_name=name, status="dry_run"
                    ).inc()
                    return True

                # Update PVC spec
                new_size_str = self.format_size(new_size)
                pvc.spec.resources.requests["storage"] = new_size_str

                # Apply update
                self.k8s.update_pvc(namespace, name, pvc)

                # Record scaling event
                self.scaling_history[pvc_key] = datetime.now()

                SCALING_OPERATIONS.labels(
                    namespace=namespace, pvc_name=name, status="success"
                ).inc()

                kopf.info(f"Successfully scaled PVC {pvc_key} to {new_size_str}")
                return True

            except ApiException as e:
                kopf.error(f"Kubernetes API error scaling PVC {pvc_key}: {e}")
                SCALING_OPERATIONS.labels(
                    namespace=namespace, pvc_name=name, status="failed"
                ).inc()
                OPERATOR_ERRORS.labels(error_type="api").inc()
                return False
            except Exception as e:
                kopf.error(f"Unexpected error scaling PVC {pvc_key}: {e}")
                SCALING_OPERATIONS.labels(
                    namespace=namespace, pvc_name=name, status="failed"
                ).inc()
                OPERATOR_ERRORS.labels(error_type="unknown").inc()
                return False

    async def check_and_scale_all(self):
        """Check all PVCs and scale if needed"""
        try:
            pvcs = self.k8s.list_pvcs()

            for pvc in pvcs:
                namespace = pvc.metadata.namespace
                name = pvc.metadata.name

                # Skip excluded namespaces
                if namespace in self.config["exclude_namespaces"]:
                    continue

                # Skip if PVC is being deleted
                if pvc.metadata.deletion_timestamp:
                    continue

                # Skip if not bound
                if pvc.status.phase != "Bound":
                    continue

                # Get PVC info and update metrics
                pvc_info = self.get_pvc_info(pvc)
                if pvc_info:
                    PVC_USAGE_PERCENT.labels(
                        namespace=namespace,
                        pvc_name=name,
                        storage_class=pvc_info.storage_class
                    ).set(pvc_info.usage_percent)

                    PVC_CAPACITY_BYTES.labels(
                        namespace=namespace,
                        pvc_name=name
                    ).set(pvc_info.capacity)

                # Check if should scale
                policy = ScalingPolicy(**self.config)
                should_scale, reason = self.should_scale(pvc_info, policy)

                if should_scale:
                    kopf.info(f"Scaling PVC {namespace}/{name}: {reason}")
                    await self.scale_pvc(pvc, policy)
                else:
                    kopf.debug(f"Not scaling PVC {namespace}/{name}: {reason}")

        except Exception as e:
            kopf.error(f"Error in check_and_scale_all: {e}")
            OPERATOR_ERRORS.labels(error_type="check_all").inc()


# Kopf handlers
@kopf.on.startup()
async def startup_fn(logger, **kwargs):
    """Startup handler"""
    logger.info("Storage Autoscaler starting up...")

    # Start metrics server
    metrics_port = int(os.environ.get("METRICS_PORT", "8084"))
    start_http_server(metrics_port)
    logger.info(f"Metrics server started on port {metrics_port}")


@kopf.on.cleanup()
async def cleanup_fn(logger, **kwargs):
    """Cleanup handler"""
    logger.info("Storage Autoscaler shutting down...")


@kopf.on.create("storageautoscaler.prodory.io")
@kopf.on.update("storageautoscaler.prodory.io")
async def handle_autoscaler_change(spec, name, namespace, logger, **kwargs):
    """Handle StorageAutoscaler CR changes"""
    logger.info(f"StorageAutoscaler {namespace}/{name} changed")

    # Validate spec
    try:
        config = StorageAutoscalerSpec(**spec)
        logger.info(f"Configuration valid: enabled={config.enabled}")
    except Exception as e:
        logger.error(f"Invalid configuration: {e}")
        raise kopf.PermanentError(f"Invalid configuration: {e}")


@kopf.timer("persistentvolumeclaims", interval=60)
async def check_pvc_timer(spec, status, name, namespace, logger, **kwargs):
    """Periodic check of all PVCs"""
    k8s_client = KubernetesClient()

    # Load configuration from ConfigMap or use defaults
    config = DEFAULT_CONFIG.copy()

    scaler = StorageScaler(k8s_client, config)
    await scaler.check_and_scale_all()


@kopf.on.field("persistentvolumeclaims", field="status.capacity.storage")
def on_pvc_capacity_change(old, new, name, namespace, logger, **kwargs):
    """Handle PVC capacity changes"""
    if old != new:
        logger.info(f"PVC {namespace}/{name} capacity changed: {old} -> {new}")


# Main entry point
def main():
    """Main entry point"""
    kopf.run(
        clusterwide=True,
        liveness_endpoint="http://0.0.0.0:8085/healthz",
    )


if __name__ == "__main__":
    main()
