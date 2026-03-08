FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY src/ ./

# Create non-root user
RUN useradd -m -u 1000 autoscaler && \
    chown -R autoscaler:autoscaler /app

USER autoscaler

# Expose ports
EXPOSE 8084 8085

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8085/healthz')" || exit 1

# Run the operator
CMD ["kopf", "run", "operator.py", "--verbose", "--all-namespaces"]
