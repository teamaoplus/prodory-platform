package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"cloud-sentinel/internal/models"
)

var (
	metricsCollected = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sentinel_metrics_collected_total",
			Help: "Total number of metrics collected",
		},
		[]string{"source"},
	)

	collectionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sentinel_collection_duration_seconds",
			Help:    "Duration of metric collection operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"source"},
	)

	collectionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sentinel_collection_errors_total",
			Help: "Total number of collection errors",
		},
		[]string{"source"},
	)
)

func init() {
	prometheus.MustRegister(metricsCollected, collectionDuration, collectionErrors)
}

// Collector handles metric collection from various sources
type Collector struct {
	log          *logrus.Logger
	prometheusURL string
	httpClient   *http.Client
	cache        *metricCache
}

// metricCache provides simple caching for metrics
type metricCache struct {
	data      map[string]cacheEntry
	ttl       time.Duration
}

type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// NewCollector creates a new metrics collector
func NewCollector(log *logrus.Logger) *Collector {
	return &Collector{
		log:          log,
		prometheusURL: "http://prometheus:9090",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: &metricCache{
			data: make(map[string]cacheEntry),
			ttl:  30 * time.Second,
		},
	}
}

// SetPrometheusURL sets the Prometheus URL
func (c *Collector) SetPrometheusURL(url string) {
	c.prometheusURL = url
}

// QueryPrometheus executes a PromQL query against Prometheus
func (c *Collector) QueryPrometheus(ctx context.Context, query string) ([]models.MetricResult, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("prometheus").Observe(time.Since(start).Seconds())
	}()

	u, err := url.Parse(c.prometheusURL + "/api/v1/query")
	if err != nil {
		collectionErrors.WithLabelValues("prometheus").Inc()
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("time", fmt.Sprintf("%d", time.Now().Unix()))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		collectionErrors.WithLabelValues("prometheus").Inc()
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		collectionErrors.WithLabelValues("prometheus").Inc()
		return nil, fmt.Errorf("failed to query Prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		collectionErrors.WithLabelValues("prometheus").Inc()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("prometheus returned status %d: %s", resp.StatusCode, string(body))
	}

	var result prometheusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		collectionErrors.WithLabelValues("prometheus").Inc()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		collectionErrors.WithLabelValues("prometheus").Inc()
		return nil, fmt.Errorf("prometheus query failed: %s", result.Error)
	}

	metricsCollected.WithLabelValues("prometheus").Inc()
	return c.convertPrometheusResult(result.Data), nil
}

// QueryPrometheusRange executes a range query against Prometheus
func (c *Collector) QueryPrometheusRange(ctx context.Context, query string, start, end time.Time, step string) ([]models.MetricResult, error) {
	startTime := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("prometheus_range").Observe(time.Since(startTime).Seconds())
	}()

	u, err := url.Parse(c.prometheusURL + "/api/v1/query_range")
	if err != nil {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", step)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		return nil, fmt.Errorf("failed to query Prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("prometheus returned status %d: %s", resp.StatusCode, string(body))
	}

	var result prometheusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		collectionErrors.WithLabelValues("prometheus_range").Inc()
		return nil, fmt.Errorf("prometheus query failed: %s", result.Error)
	}

	metricsCollected.WithLabelValues("prometheus_range").Inc()
	return c.convertPrometheusResult(result.Data), nil
}

// GetAWSMetrics retrieves AWS CloudWatch metrics
func (c *Collector) GetAWSMetrics(ctx context.Context, region string) ([]models.Metric, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("aws").Observe(time.Since(start).Seconds())
	}()

	// Placeholder for AWS CloudWatch integration
	// In production, use aws-sdk-go-v2/service/cloudwatch
	c.log.Debugf("Collecting AWS metrics for region: %s", region)

	metrics := []models.Metric{
		{
			Name:      "aws_ec2_cpu_utilization",
			Labels:    map[string]string{"region": region, "service": "ec2"},
			Value:     45.2,
			Timestamp: time.Now(),
			Unit:      "Percent",
		},
		{
			Name:      "aws_rds_free_storage",
			Labels:    map[string]string{"region": region, "service": "rds"},
			Value:     1024.5,
			Timestamp: time.Now(),
			Unit:      "Gigabytes",
		},
	}

	metricsCollected.WithLabelValues("aws").Add(float64(len(metrics)))
	return metrics, nil
}

// GetAzureMetrics retrieves Azure Monitor metrics
func (c *Collector) GetAzureMetrics(ctx context.Context) ([]models.Metric, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("azure").Observe(time.Since(start).Seconds())
	}()

	// Placeholder for Azure Monitor integration
	c.log.Debug("Collecting Azure metrics")

	metrics := []models.Metric{
		{
			Name:      "azure_vm_cpu_percentage",
			Labels:    map[string]string{"service": "vm"},
			Value:     32.5,
			Timestamp: time.Now(),
			Unit:      "Percent",
		},
	}

	metricsCollected.WithLabelValues("azure").Add(float64(len(metrics)))
	return metrics, nil
}

// GetGCPMetrics retrieves GCP Monitoring metrics
func (c *Collector) GetGCPMetrics(ctx context.Context, projectID string) ([]models.Metric, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("gcp").Observe(time.Since(start).Seconds())
	}()

	// Placeholder for GCP Monitoring integration
	c.log.Debugf("Collecting GCP metrics for project: %s", projectID)

	metrics := []models.Metric{
		{
			Name:      "gcp_compute_cpu_utilization",
			Labels:    map[string]string{"project": projectID, "service": "compute"},
			Value:     28.7,
			Timestamp: time.Now(),
			Unit:      "Percent",
		},
	}

	metricsCollected.WithLabelValues("gcp").Add(float64(len(metrics)))
	return metrics, nil
}

// GetKubernetesMetrics retrieves Kubernetes metrics
func (c *Collector) GetKubernetesMetrics(ctx context.Context) ([]models.Metric, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("kubernetes").Observe(time.Since(start).Seconds())
	}()

	// Query kube-state-metrics and metrics-server via Prometheus
	queries := []string{
		"kube_node_status_condition{condition=\"Ready\",status=\"true\"}",
		"kube_pod_status_phase{phase=\"Running\"}",
		"container_cpu_usage_seconds_total",
		"container_memory_usage_bytes",
	}

	var allMetrics []models.Metric
	for _, query := range queries {
		results, err := c.QueryPrometheus(ctx, query)
		if err != nil {
			c.log.Warnf("Failed to query %s: %v", query, err)
			continue
		}

		for _, result := range results {
			if len(result.Values) > 0 {
				allMetrics = append(allMetrics, models.Metric{
					Name:      result.Metric["__name__"],
					Labels:    result.Metric,
					Value:     result.Values[len(result.Values)-1].Value,
					Timestamp: time.Now(),
				})
			}
		}
	}

	metricsCollected.WithLabelValues("kubernetes").Add(float64(len(allMetrics)))
	return allMetrics, nil
}

// GetCustomMetrics retrieves custom application metrics
func (c *Collector) GetCustomMetrics(ctx context.Context) ([]models.Metric, error) {
	start := time.Now()
	defer func() {
		collectionDuration.WithLabelValues("custom").Observe(time.Since(start).Seconds())
	}()

	// Query custom metrics from Prometheus
	queries := []string{
		"up",
		"http_requests_total",
		"http_request_duration_seconds",
	}

	var allMetrics []models.Metric
	for _, query := range queries {
		results, err := c.QueryPrometheus(ctx, query)
		if err != nil {
			c.log.Warnf("Failed to query %s: %v", query, err)
			continue
		}

		for _, result := range results {
			if len(result.Values) > 0 {
				allMetrics = append(allMetrics, models.Metric{
					Name:      result.Metric["__name__"],
					Labels:    result.Metric,
					Value:     result.Values[len(result.Values)-1].Value,
					Timestamp: time.Now(),
				})
			}
		}
	}

	metricsCollected.WithLabelValues("custom").Add(float64(len(allMetrics)))
	return allMetrics, nil
}

// GetDashboardSummary returns summary metrics for the dashboard
func (c *Collector) GetDashboardSummary(ctx context.Context) (*models.DashboardSummary, error) {
	summary := &models.DashboardSummary{
		SourcesConnected: []string{},
		AlertTrend:       []models.AlertTrendPoint{},
		TopAlertingRules: []models.RuleAlertCount{},
	}

	// Query active alerts from Prometheus
	alertResults, err := c.QueryPrometheus(ctx, "ALERTS{alertstate=\"firing\"}")
	if err != nil {
		c.log.Warnf("Failed to query alerts: %v", err)
	} else {
		summary.ActiveAlerts = int64(len(alertResults))
	}

	// Check connected sources
	if _, err := c.QueryPrometheus(ctx, "up"); err == nil {
		summary.SourcesConnected = append(summary.SourcesConnected, "prometheus")
	}

	return summary, nil
}

// convertPrometheusResult converts Prometheus response to MetricResult
func (c *Collector) convertPrometheusResult(data prometheusData) []models.MetricResult {
	var results []models.MetricResult

	switch data.ResultType {
	case "vector":
		for _, r := range data.Result {
			result := models.MetricResult{
				Metric: r.Metric,
				Values: []models.MetricValue{
					{
						Timestamp: r.Value[0].(float64),
						Value:     parseFloat(r.Value[1]),
					},
				},
			}
			results = append(results, result)
		}
	case "matrix":
		for _, r := range data.Result {
			result := models.MetricResult{
				Metric: r.Metric,
				Values: make([]models.MetricValue, 0, len(r.Values)),
			}
			for _, v := range r.Values {
				result.Values = append(result.Values, models.MetricValue{
					Timestamp: v[0].(float64),
					Value:     parseFloat(v[1]),
				})
			}
			results = append(results, result)
		}
	}

	return results
}

// prometheusResponse represents Prometheus API response
type prometheusResponse struct {
	Status string         `json:"status"`
	Data   prometheusData `json:"data"`
	Error  string         `json:"error,omitempty"`
}

// prometheusData represents Prometheus data
type prometheusData struct {
	ResultType string             `json:"resultType"`
	Result     []prometheusResult `json:"result"`
}

// prometheusResult represents a single Prometheus result
type prometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values [][]interface{}   `json:"values"`
}

// parseFloat parses interface to float64
func parseFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}
