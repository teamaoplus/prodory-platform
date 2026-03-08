package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"cloud-sentinel/internal/alerts"
	"cloud-sentinel/internal/config"
	"cloud-sentinel/internal/metrics"
	"cloud-sentinel/internal/models"
)

// Scheduler handles scheduled tasks for Cloud Sentinel
type Scheduler struct {
	config    *config.Config
	log       *logrus.Logger
	cron      *cron.Cron
	collector *metrics.Collector
	manager   *alerts.Manager
	tasks     map[string]cron.EntryID
}

// New creates a new scheduler
func New(cfg *config.Config, log *logrus.Logger, collector *metrics.Collector, manager *alerts.Manager) *Scheduler {
	return &Scheduler{
		config:    cfg,
		log:       log,
		cron:      cron.New(cron.WithSeconds()),
		collector: collector,
		manager:   manager,
		tasks:     make(map[string]cron.EntryID),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.log.Info("Starting scheduler...")

	// Schedule alert rule evaluation
	s.scheduleRuleEvaluation()

	// Schedule metric collection
	s.scheduleMetricCollection()

	// Schedule alert cleanup
	s.scheduleAlertCleanup()

	// Schedule health checks
	s.scheduleHealthChecks()

	s.cron.Start()
	s.log.Info("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.log.Info("Stopping scheduler...")
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.log.Info("Scheduler stopped")
}

// scheduleRuleEvaluation schedules alert rule evaluation
func (s *Scheduler) scheduleRuleEvaluation() {
	// Evaluate rules every minute
	id, err := s.cron.AddFunc("0 * * * * *", func() {
		s.evaluateRules()
	})
	if err != nil {
		s.log.Errorf("Failed to schedule rule evaluation: %v", err)
		return
	}
	s.tasks["rule_evaluation"] = id
	s.log.Debug("Scheduled rule evaluation")
}

// scheduleMetricCollection schedules metric collection
func (s *Scheduler) scheduleMetricCollection() {
	// Collect metrics every 30 seconds
	id, err := s.cron.AddFunc("*/30 * * * * *", func() {
		s.collectMetrics()
	})
	if err != nil {
		s.log.Errorf("Failed to schedule metric collection: %v", err)
		return
	}
	s.tasks["metric_collection"] = id
	s.log.Debug("Scheduled metric collection")
}

// scheduleAlertCleanup schedules alert cleanup
func (s *Scheduler) scheduleAlertCleanup() {
	// Clean up old alerts daily at midnight
	id, err := s.cron.AddFunc("0 0 0 * * *", func() {
		s.cleanupAlerts()
	})
	if err != nil {
		s.log.Errorf("Failed to schedule alert cleanup: %v", err)
		return
	}
	s.tasks["alert_cleanup"] = id
	s.log.Debug("Scheduled alert cleanup")
}

// scheduleHealthChecks schedules health checks
func (s *Scheduler) scheduleHealthChecks() {
	// Health check every 5 minutes
	id, err := s.cron.AddFunc("0 */5 * * * *", func() {
		s.runHealthChecks()
	})
	if err != nil {
		s.log.Errorf("Failed to schedule health checks: %v", err)
		return
	}
	s.tasks["health_check"] = id
	s.log.Debug("Scheduled health checks")
}

// evaluateRules evaluates all enabled alert rules
func (s *Scheduler) evaluateRules() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	rules := s.manager.GetAlertRules(true)
	s.log.Debugf("Evaluating %d alert rules", len(rules))

	for _, rule := range rules {
		go s.evaluateRule(ctx, &rule)
	}
}

// evaluateRule evaluates a single alert rule
func (s *Scheduler) evaluateRule(ctx context.Context, rule *models.AlertRule) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Errorf("Panic evaluating rule %s: %v", rule.Name, r)
		}
	}()

	var results []models.MetricResult
	var err error

	switch rule.Source {
	case "prometheus":
		results, err = s.collector.QueryPrometheus(ctx, rule.Query)
	default:
		s.log.Warnf("Unknown rule source: %s", rule.Source)
		return
	}

	if err != nil {
		s.log.Errorf("Failed to evaluate rule %s: %v", rule.Name, err)
		return
	}

	// Check if any result violates the threshold
	for _, result := range results {
		if len(result.Values) == 0 {
			continue
		}

		value := result.Values[len(result.Values)-1].Value
		violated := false

		switch rule.Operator {
		case "gt":
			violated = value > rule.Threshold
		case "lt":
			violated = value < rule.Threshold
		case "eq":
			violated = value == rule.Threshold
		case "ne":
			violated = value != rule.Threshold
		}

		if violated {
			s.createAlertFromRule(rule, result, value)
		}
	}
}

// createAlertFromRule creates an alert from a rule violation
func (s *Scheduler) createAlertFromRule(rule *models.AlertRule, result models.MetricResult, value float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build alert title and description from annotations
	title := rule.Annotations["summary"]
	if title == "" {
		title = rule.Name
	}

	description := rule.Annotations["description"]
	if description == "" {
		description = fmt.Sprintf("%s: value %.2f violates threshold %.2f", rule.Name, value, rule.Threshold)
	}

	// Replace template variables in description
	for key, val := range result.Metric {
		description = replaceAll(description, "{{ $labels."+key+" }}", val)
		description = replaceAll(description, "{{ $value }}", fmt.Sprintf("%.2f", value))
	}

	alert := &models.Alert{
		RuleID:      rule.ID,
		Title:       title,
		Description: description,
		Severity:    rule.Severity,
		Source:      rule.Source,
		Labels:      result.Metric,
		Value:       value,
		Threshold:   rule.Threshold,
		RunbookURL:  rule.RunbookURL,
		Annotations: rule.Annotations,
	}

	_, err := s.manager.CreateAlert(ctx, alert)
	if err != nil {
		s.log.Debugf("Alert not created (may be silenced or duplicate): %v", err)
	}
}

// collectMetrics collects metrics from all sources
func (s *Scheduler) collectMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Collect from Prometheus
	if _, err := s.collector.GetCustomMetrics(ctx); err != nil {
		s.log.Warnf("Failed to collect custom metrics: %v", err)
	}

	// Collect from Kubernetes
	if _, err := s.collector.GetKubernetesMetrics(ctx); err != nil {
		s.log.Warnf("Failed to collect Kubernetes metrics: %v", err)
	}

	// Collect from cloud providers if configured
	if s.config.AWSAccessKeyID != "" {
		for _, region := range s.config.GetAWSRegions() {
			if _, err := s.collector.GetAWSMetrics(ctx, region); err != nil {
				s.log.Warnf("Failed to collect AWS metrics for %s: %v", region, err)
			}
		}
	}

	if s.config.AzureClientID != "" {
		if _, err := s.collector.GetAzureMetrics(ctx); err != nil {
			s.log.Warnf("Failed to collect Azure metrics: %v", err)
		}
	}

	if s.config.GCPProjectID != "" {
		if _, err := s.collector.GetGCPMetrics(ctx, s.config.GCPProjectID); err != nil {
			s.log.Warnf("Failed to collect GCP metrics: %v", err)
		}
	}
}

// cleanupAlerts removes old resolved alerts
func (s *Scheduler) cleanupAlerts() {
	s.log.Debug("Cleaning up old alerts")

	// Get all alerts
	alerts := s.manager.GetAlerts("", "", "", 0)

	retentionDays := s.config.AlertRetentionDays
	if retentionDays == 0 {
		retentionDays = 30
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	for _, alert := range alerts {
		// Delete resolved alerts older than retention period
		if alert.Status == models.StatusResolved && alert.ResolvedAt != nil && alert.ResolvedAt.Before(cutoff) {
			if err := s.manager.DeleteAlert(alert.ID); err != nil {
				s.log.Warnf("Failed to delete old alert %s: %v", alert.ID, err)
			}
		}

		// Auto-resolve old firing alerts if enabled
		if s.config.EnableAutoResolve && alert.Status == models.StatusFiring {
			autoResolveAfter := s.config.AutoResolveAfter
			if autoResolveAfter == 0 {
				autoResolveAfter = 24 * time.Hour
			}

			if time.Since(alert.StartedAt) > autoResolveAfter {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, err := s.manager.ResolveAlert(ctx, alert.ID)
				cancel()

				if err != nil {
					s.log.Warnf("Failed to auto-resolve alert %s: %v", alert.ID, err)
				} else {
					s.log.Infof("Auto-resolved alert %s after %v", alert.ID, autoResolveAfter)
				}
			}
		}
	}
}

// runHealthChecks runs health checks on configured sources
func (s *Scheduler) runHealthChecks() {
	s.log.Debug("Running health checks")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check Prometheus connectivity
	if _, err := s.collector.QueryPrometheus(ctx, "up"); err != nil {
		s.log.Warnf("Prometheus health check failed: %v", err)
	}

	// Additional health checks can be added here
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	for {
		newStr := ""
		for i := 0; i < len(s); {
			if i+len(old) <= len(s) && s[i:i+len(old)] == old {
				newStr += new
				i += len(old)
			} else {
				newStr += string(s[i])
				i++
			}
		}
		if newStr == s {
			return s
		}
		s = newStr
	}
}
