package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cloud-sentinel/internal/alerts"
	"cloud-sentinel/internal/metrics"
	"cloud-sentinel/internal/models"
)

// Dashboard handlers

// GetDashboardSummary returns dashboard summary
func GetDashboardSummary(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		summary, err := collector.GetDashboardSummary(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

// GetMetricsOverview returns metrics overview
func GetMetricsOverview(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		overview := gin.H{
			"sources": []string{"prometheus", "aws", "azure", "gcp", "kubernetes"},
			"metrics_collected": 0,
			"last_collection": time.Now(),
		}
		c.JSON(http.StatusOK, overview)
	}
}

// GetActiveAlerts returns active alerts
func GetActiveAlerts(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		alerts := manager.GetAlerts("firing", "", "", 50)
		c.JSON(http.StatusOK, gin.H{"alerts": alerts})
	}
}

// Metrics handlers

// GetAWSMetrics returns AWS metrics
func GetAWSMetrics(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		region := c.DefaultQuery("region", "us-east-1")
		metrics, err := collector.GetAWSMetrics(c.Request.Context(), region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// GetAzureMetrics returns Azure metrics
func GetAzureMetrics(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, err := collector.GetAzureMetrics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// GetGCPMetrics returns GCP metrics
func GetGCPMetrics(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.DefaultQuery("project", "")
		metrics, err := collector.GetGCPMetrics(c.Request.Context(), projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// GetKubernetesMetrics returns Kubernetes metrics
func GetKubernetesMetrics(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, err := collector.GetKubernetesMetrics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// GetCustomMetrics returns custom metrics
func GetCustomMetrics(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, err := collector.GetCustomMetrics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"metrics": metrics})
	}
}

// CreateCustomMetric creates a custom metric
func CreateCustomMetric(collector *metrics.Collector) gin.HandlerFunc {
	return func(c *gin.Context) {
		var metric models.Metric
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metric.Timestamp = time.Now()
		c.JSON(http.StatusCreated, metric)
	}
}

// Alert handlers

// GetAlerts returns alerts with filtering
func GetAlerts(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		severity := c.Query("severity")
		source := c.Query("source")
		limit := 100

		alerts := manager.GetAlerts(status, severity, source, limit)
		c.JSON(http.StatusOK, gin.H{
			"alerts": alerts,
			"total":  len(alerts),
		})
	}
}

// GetAlert returns a single alert
func GetAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		alert, err := manager.GetAlert(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, alert)
	}
}

// CreateAlert creates a new alert
func CreateAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var alert models.Alert
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		created, err := manager.CreateAlert(c.Request.Context(), &alert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

// UpdateAlert updates an alert
func UpdateAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updates models.Alert
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := manager.UpdateAlert(c.Request.Context(), id, &updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// DeleteAlert deletes an alert
func DeleteAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := manager.DeleteAlert(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// AcknowledgeAlert acknowledges an alert
func AcknowledgeAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			User string `json:"user"`
		}
		c.ShouldBindJSON(&req)

		alert, err := manager.AcknowledgeAlert(c.Request.Context(), id, req.User)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, alert)
	}
}

// ResolveAlert resolves an alert
func ResolveAlert(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		alert, err := manager.ResolveAlert(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, alert)
	}
}

// Alert Rule handlers

// GetAlertRules returns all alert rules
func GetAlertRules(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		enabledOnly := c.Query("enabled") == "true"
		rules := manager.GetAlertRules(enabledOnly)
		c.JSON(http.StatusOK, gin.H{
			"rules": rules,
			"total": len(rules),
		})
	}
}

// GetAlertRule returns a single alert rule
func GetAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := manager.GetAlertRule(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rule)
	}
}

// CreateAlertRule creates a new alert rule
func CreateAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rule models.AlertRule
		if err := c.ShouldBindJSON(&rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		created, err := manager.CreateAlertRule(c.Request.Context(), &rule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

// UpdateAlertRule updates an alert rule
func UpdateAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updates models.AlertRule
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := manager.UpdateAlertRule(c.Request.Context(), id, &updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// DeleteAlertRule deletes an alert rule
func DeleteAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := manager.DeleteAlertRule(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// EnableAlertRule enables an alert rule
func EnableAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := manager.EnableAlertRule(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rule)
	}
}

// DisableAlertRule disables an alert rule
func DisableAlertRule(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := manager.DisableAlertRule(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rule)
	}
}

// Notification Channel handlers

// GetNotificationChannels returns all notification channels
func GetNotificationChannels(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		channels := manager.GetNotificationChannels()
		c.JSON(http.StatusOK, gin.H{
			"channels": channels,
			"total":    len(channels),
		})
	}
}

// CreateNotificationChannel creates a new notification channel
func CreateNotificationChannel(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var channel models.NotificationChannel
		if err := c.ShouldBindJSON(&channel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		created, err := manager.CreateNotificationChannel(c.Request.Context(), &channel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

// UpdateNotificationChannel updates a notification channel
func UpdateNotificationChannel(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updates models.NotificationChannel
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := manager.UpdateNotificationChannel(c.Request.Context(), id, &updates)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// DeleteNotificationChannel deletes a notification channel
func DeleteNotificationChannel(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := manager.DeleteNotificationChannel(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// TestNotificationChannel tests a notification channel
func TestNotificationChannel(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := manager.TestNotificationChannel(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Test notification sent successfully"})
	}
}

// Silence handlers

// GetSilences returns active silences
func GetSilences(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		silences := manager.GetSilences()
		c.JSON(http.StatusOK, gin.H{
			"silences": silences,
			"total":    len(silences),
		})
	}
}

// CreateSilence creates a new silence
func CreateSilence(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var silence models.Silence
		if err := c.ShouldBindJSON(&silence); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		created, err := manager.CreateSilence(c.Request.Context(), &silence)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, created)
	}
}

// DeleteSilence deletes a silence
func DeleteSilence(manager *alerts.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := manager.DeleteSilence(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

// Report handlers

// GetReports returns generated reports
func GetReports() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"reports": []models.Report{},
			"total":   0,
		})
	}
}

// GenerateReport generates a new report
func GenerateReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name   string    `json:"name"`
			Type   string    `json:"type"`
			Format string    `json:"format"`
			Start  time.Time `json:"start"`
			End    time.Time `json:"end"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		report := models.Report{
			ID:          "report-" + time.Now().Format("20060102150405"),
			Name:        req.Name,
			Type:        req.Type,
			Format:      req.Format,
			Status:      "generating",
			StartTime:   req.Start,
			EndTime:     req.End,
			GeneratedAt: time.Now(),
		}

		c.JSON(http.StatusAccepted, report)
	}
}

// DownloadReport downloads a report
func DownloadReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "Download not yet implemented for report: " + id,
		})
	}
}

// Settings handlers

// GetSettings returns application settings
func GetSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		settings := models.Settings{
			ID:                   "default",
			DefaultAlertInterval: 60,
			MaxAlertsPerMinute:   100,
			AlertRetentionDays:   30,
			EnableAutoResolve:    true,
			AutoResolveAfter:     24 * time.Hour,
			Theme:                "dark",
			Timezone:             "UTC",
			DateFormat:           "2006-01-02 15:04:05",
		}
		c.JSON(http.StatusOK, settings)
	}
}

// UpdateSettings updates application settings
func UpdateSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings models.Settings
		if err := c.ShouldBindJSON(&settings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		settings.UpdatedAt = time.Now()
		c.JSON(http.StatusOK, settings)
	}
}

// Discovery handlers

// RunDiscovery runs resource discovery
func RunDiscovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Providers []string `json:"providers"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message":   "Discovery started",
			"providers": req.Providers,
		})
	}
}

// GetDiscoveredResources returns discovered resources
func GetDiscoveredResources() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"resources": []models.CloudResource{},
			"total":     0,
		})
	}
}
