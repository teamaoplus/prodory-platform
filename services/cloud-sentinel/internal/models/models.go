package models

import (
	"time"
)

// Alert represents a monitoring alert
type Alert struct {
	ID          string                 `json:"id" db:"id"`
	RuleID      string                 `json:"rule_id" db:"rule_id"`
	Title       string                 `json:"title" db:"title"`
	Description string                 `json:"description" db:"description"`
	Severity    AlertSeverity          `json:"severity" db:"severity"`
	Status      AlertStatus            `json:"status" db:"status"`
	Source      string                 `json:"source" db:"source"`
	Labels      map[string]string      `json:"labels" db:"labels"`
	Value       float64                `json:"value" db:"value"`
	Threshold   float64                `json:"threshold" db:"threshold"`
	StartedAt   time.Time              `json:"started_at" db:"started_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty" db:"resolved_at"`
	AcknowledgedAt *time.Time          `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	AcknowledgedBy string             `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	RunbookURL  string                 `json:"runbook_url,omitempty" db:"runbook_url"`
	Annotations map[string]string      `json:"annotations" db:"annotations"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// AlertSeverity represents the severity of an alert
type AlertSeverity string

const (
	SeverityCritical AlertSeverity = "critical"
	SeverityWarning  AlertSeverity = "warning"
	SeverityInfo     AlertSeverity = "info"
)

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	StatusFiring      AlertStatus = "firing"
	StatusResolved    AlertStatus = "resolved"
	StatusAcknowledged AlertStatus = "acknowledged"
	StatusSuppressed  AlertStatus = "suppressed"
)

// AlertRule represents a rule for generating alerts
type AlertRule struct {
	ID          string            `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Enabled     bool              `json:"enabled" db:"enabled"`
	Query       string            `json:"query" db:"query"`
	Duration    time.Duration     `json:"duration" db:"duration"`
	Severity    AlertSeverity     `json:"severity" db:"severity"`
	Labels      map[string]string `json:"labels" db:"labels"`
	Annotations map[string]string `json:"annotations" db:"annotations"`
	Threshold   float64           `json:"threshold" db:"threshold"`
	Operator    string            `json:"operator" db:"operator"` // gt, lt, eq, ne
	Source      string            `json:"source" db:"source"`     // prometheus, cloudwatch, azure_monitor, gcp_monitoring
	Interval    time.Duration     `json:"interval" db:"interval"`
	RunbookURL  string            `json:"runbook_url" db:"runbook_url"`
	Group       string            `json:"group" db:"group"`
	ChannelIDs  []string          `json:"channel_ids" db:"channel_ids"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
	LastEvaluatedAt *time.Time    `json:"last_evaluated_at,omitempty" db:"last_evaluated_at"`
}

// NotificationChannel represents a channel for sending alerts
type NotificationChannel struct {
	ID        string                 `json:"id" db:"id"`
	Name      string                 `json:"name" db:"name"`
	Type      ChannelType            `json:"type" db:"type"`
	Config    map[string]interface{} `json:"config" db:"config"`
	Enabled   bool                   `json:"enabled" db:"enabled"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// ChannelType represents the type of notification channel
type ChannelType string

const (
	ChannelSlack     ChannelType = "slack"
	ChannelEmail     ChannelType = "email"
	ChannelPagerDuty ChannelType = "pagerduty"
	ChannelWebhook   ChannelType = "webhook"
	ChannelTeams     ChannelType = "teams"
	ChannelSMS       ChannelType = "sms"
)

// Silence represents a suppression of alerts
type Silence struct {
	ID          string            `json:"id" db:"id"`
	Matchers    map[string]string `json:"matchers" db:"matchers"`
	StartsAt    time.Time         `json:"starts_at" db:"starts_at"`
	EndsAt      time.Time         `json:"ends_at" db:"ends_at"`
	CreatedBy   string            `json:"created_by" db:"created_by"`
	Comment     string            `json:"comment" db:"comment"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

// Metric represents a time-series metric
type Metric struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Unit      string            `json:"unit,omitempty"`
}

// MetricQuery represents a query for metrics
type MetricQuery struct {
	Query     string    `json:"query"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Step      string    `json:"step"`
}

// MetricResult represents the result of a metric query
type MetricResult struct {
	Metric map[string]string `json:"metric"`
	Values []MetricValue     `json:"values"`
}

// MetricValue represents a single metric value
type MetricValue struct {
	Timestamp float64 `json:"timestamp"`
	Value     float64 `json:"value"`
}

// DashboardSummary represents a summary for the dashboard
type DashboardSummary struct {
	ActiveAlerts        int64                  `json:"active_alerts"`
	CriticalAlerts      int64                  `json:"critical_alerts"`
	WarningAlerts       int64                  `json:"warning_alerts"`
	ResolvedToday       int64                  `json:"resolved_today"`
	TotalRules          int64                  `json:"total_rules"`
	EnabledRules        int64                  `json:"enabled_rules"`
	MetricsCollected    int64                  `json:"metrics_collected"`
	SourcesConnected    []string               `json:"sources_connected"`
	RecentAlerts        []Alert                `json:"recent_alerts"`
	AlertTrend          []AlertTrendPoint      `json:"alert_trend"`
	TopAlertingRules    []RuleAlertCount       `json:"top_alerting_rules"`
}

// AlertTrendPoint represents a point in the alert trend
type AlertTrendPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

// RuleAlertCount represents alert count per rule
type RuleAlertCount struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
	Count    int64  `json:"count"`
}

// CloudResource represents a discovered cloud resource
type CloudResource struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Provider     string            `json:"provider"`
	Region       string            `json:"region"`
	Account      string            `json:"account"`
	Labels       map[string]string `json:"labels"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	LastSeenAt   time.Time         `json:"last_seen_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// Report represents a generated report
type Report struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Format      string    `json:"format" db:"format"`
	Status      string    `json:"status" db:"status"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	GeneratedAt time.Time `json:"generated_at" db:"generated_at"`
	GeneratedBy string    `json:"generated_by" db:"generated_by"`
	FileURL     string    `json:"file_url" db:"file_url"`
	Size        int64     `json:"size" db:"size"`
	Parameters  map[string]interface{} `json:"parameters" db:"parameters"`
}

// Settings represents application settings
type Settings struct {
	ID                    string                 `json:"id" db:"id"`
	DefaultAlertInterval  int                    `json:"default_alert_interval" db:"default_alert_interval"`
	MaxAlertsPerMinute    int                    `json:"max_alerts_per_minute" db:"max_alerts_per_minute"`
	AlertRetentionDays    int                    `json:"alert_retention_days" db:"alert_retention_days"`
	EnableAutoResolve     bool                   `json:"enable_auto_resolve" db:"enable_auto_resolve"`
	AutoResolveAfter      time.Duration          `json:"auto_resolve_after" db:"auto_resolve_after"`
	Theme                 string                 `json:"theme" db:"theme"`
	Timezone              string                 `json:"timezone" db:"timezone"`
	DateFormat            string                 `json:"date_format" db:"date_format"`
	NotificationDefaults  map[string]interface{} `json:"notification_defaults" db:"notification_defaults"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
	UpdatedBy             string                 `json:"updated_by" db:"updated_by"`
}
