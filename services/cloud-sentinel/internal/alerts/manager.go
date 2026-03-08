package alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"cloud-sentinel/internal/config"
	"cloud-sentinel/internal/models"
)

var (
	alertsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sentinel_alerts_total",
			Help: "Total number of alerts",
		},
		[]string{"severity", "status"},
	)

	notificationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sentinel_notifications_total",
			Help: "Total number of notifications sent",
		},
		[]string{"channel", "status"},
	)
)

func init() {
	prometheus.MustRegister(alertsTotal, notificationsTotal)
}

// Manager handles alert management
type Manager struct {
	config           *config.Config
	log              *logrus.Logger
	alerts           map[string]*models.Alert
	rules            map[string]*models.AlertRule
	channels         map[string]*models.NotificationChannel
	silences         map[string]*models.Silence
	mutex            sync.RWMutex
	slackClient      *slack.Client
	httpClient       *http.Client
	alertHistory     *AlertHistory
}

// AlertHistory maintains alert history
type AlertHistory struct {
	alerts []models.Alert
	mutex  sync.RWMutex
	limit  int
}

// NewManager creates a new alert manager
func NewManager(cfg *config.Config, log *logrus.Logger) (*Manager, error) {
	m := &Manager{
		config:   cfg,
		log:      log,
		alerts:   make(map[string]*models.Alert),
		rules:    make(map[string]*models.AlertRule),
		channels: make(map[string]*models.NotificationChannel),
		silences: make(map[string]*models.Silence),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		alertHistory: &AlertHistory{
			alerts: make([]models.Alert, 0, 1000),
			limit:  1000,
		},
	}

	// Initialize Slack client if configured
	if cfg.SlackBotToken != "" {
		m.slackClient = slack.New(cfg.SlackBotToken)
	}

	// Load default rules
	m.loadDefaultRules()

	// Load default channels
	m.loadDefaultChannels()

	return m, nil
}

// CreateAlert creates a new alert
func (m *Manager) CreateAlert(ctx context.Context, alert *models.Alert) (*models.Alert, error) {
	// Check if alert is silenced
	if m.isSilenced(alert) {
		m.log.Debugf("Alert is silenced: %s", alert.Title)
		return nil, fmt.Errorf("alert is silenced")
	}

	alert.ID = uuid.New().String()
	alert.Status = models.StatusFiring
	alert.StartedAt = time.Now()
	alert.CreatedAt = time.Now()
	alert.UpdatedAt = time.Now()

	m.mutex.Lock()
	m.alerts[alert.ID] = alert
	m.mutex.Unlock()

	// Add to history
	m.alertHistory.add(*alert)

	// Update Prometheus metrics
	alertsTotal.WithLabelValues(string(alert.Severity), string(alert.Status)).Inc()

	// Send notifications
	go m.sendNotifications(alert)

	m.log.Infof("Alert created: %s (severity: %s)", alert.Title, alert.Severity)
	return alert, nil
}

// GetAlert retrieves an alert by ID
func (m *Manager) GetAlert(id string) (*models.Alert, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	alert, exists := m.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert not found: %s", id)
	}

	return alert, nil
}

// GetAlerts retrieves alerts with optional filtering
func (m *Manager) GetAlerts(status string, severity string, source string, limit int) []models.Alert {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var alerts []models.Alert
	for _, alert := range m.alerts {
		if status != "" && string(alert.Status) != status {
			continue
		}
		if severity != "" && string(alert.Severity) != severity {
			continue
		}
		if source != "" && alert.Source != source {
			continue
		}
		alerts = append(alerts, *alert)
	}

	// Sort by started_at descending
	for i, j := 0, len(alerts)-1; i < j; i, j = i+1, j-1 {
		alerts[i], alerts[j] = alerts[j], alerts[i]
	}

	if limit > 0 && len(alerts) > limit {
		alerts = alerts[:limit]
	}

	return alerts
}

// UpdateAlert updates an alert
func (m *Manager) UpdateAlert(ctx context.Context, id string, updates *models.Alert) (*models.Alert, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	alert, exists := m.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert not found: %s", id)
	}

	if updates.Title != "" {
		alert.Title = updates.Title
	}
	if updates.Description != "" {
		alert.Description = updates.Description
	}
	if updates.Severity != "" {
		alert.Severity = updates.Severity
	}
	if updates.Annotations != nil {
		alert.Annotations = updates.Annotations
	}

	alert.UpdatedAt = time.Now()

	m.log.Infof("Alert updated: %s", alert.Title)
	return alert, nil
}

// DeleteAlert deletes an alert
func (m *Manager) DeleteAlert(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.alerts[id]; !exists {
		return fmt.Errorf("alert not found: %s", id)
	}

	delete(m.alerts, id)
	m.log.Infof("Alert deleted: %s", id)
	return nil
}

// AcknowledgeAlert acknowledges an alert
func (m *Manager) AcknowledgeAlert(ctx context.Context, id string, user string) (*models.Alert, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	alert, exists := m.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert not found: %s", id)
	}

	now := time.Now()
	alert.Status = models.StatusAcknowledged
	alert.AcknowledgedAt = &now
	alert.AcknowledgedBy = user
	alert.UpdatedAt = now

	m.log.Infof("Alert acknowledged: %s by %s", alert.Title, user)
	return alert, nil
}

// ResolveAlert resolves an alert
func (m *Manager) ResolveAlert(ctx context.Context, id string) (*models.Alert, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	alert, exists := m.alerts[id]
	if !exists {
		return nil, fmt.Errorf("alert not found: %s", id)
	}

	now := time.Now()
	alert.Status = models.StatusResolved
	alert.ResolvedAt = &now
	alert.UpdatedAt = now

	m.log.Infof("Alert resolved: %s", alert.Title)
	return alert, nil
}

// CreateAlertRule creates a new alert rule
func (m *Manager) CreateAlertRule(ctx context.Context, rule *models.AlertRule) (*models.AlertRule, error) {
	rule.ID = uuid.New().String()
	rule.Enabled = true
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	m.mutex.Lock()
	m.rules[rule.ID] = rule
	m.mutex.Unlock()

	m.log.Infof("Alert rule created: %s", rule.Name)
	return rule, nil
}

// GetAlertRule retrieves an alert rule by ID
func (m *Manager) GetAlertRule(id string) (*models.AlertRule, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	rule, exists := m.rules[id]
	if !exists {
		return nil, fmt.Errorf("alert rule not found: %s", id)
	}

	return rule, nil
}

// GetAlertRules retrieves all alert rules
func (m *Manager) GetAlertRules(enabledOnly bool) []models.AlertRule {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var rules []models.AlertRule
	for _, rule := range m.rules {
		if enabledOnly && !rule.Enabled {
			continue
		}
		rules = append(rules, *rule)
	}

	return rules
}

// UpdateAlertRule updates an alert rule
func (m *Manager) UpdateAlertRule(ctx context.Context, id string, updates *models.AlertRule) (*models.AlertRule, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rule, exists := m.rules[id]
	if !exists {
		return nil, fmt.Errorf("alert rule not found: %s", id)
	}

	if updates.Name != "" {
		rule.Name = updates.Name
	}
	if updates.Description != "" {
		rule.Description = updates.Description
	}
	if updates.Query != "" {
		rule.Query = updates.Query
	}
	if updates.Duration > 0 {
		rule.Duration = updates.Duration
	}
	if updates.Severity != "" {
		rule.Severity = updates.Severity
	}
	if updates.Threshold != 0 {
		rule.Threshold = updates.Threshold
	}
	if updates.Operator != "" {
		rule.Operator = updates.Operator
	}
	if updates.Interval > 0 {
		rule.Interval = updates.Interval
	}
	if updates.RunbookURL != "" {
		rule.RunbookURL = updates.RunbookURL
	}
	if updates.Group != "" {
		rule.Group = updates.Group
	}
	if updates.ChannelIDs != nil {
		rule.ChannelIDs = updates.ChannelIDs
	}
	if updates.Labels != nil {
		rule.Labels = updates.Labels
	}
	if updates.Annotations != nil {
		rule.Annotations = updates.Annotations
	}

	rule.UpdatedAt = time.Now()

	m.log.Infof("Alert rule updated: %s", rule.Name)
	return rule, nil
}

// DeleteAlertRule deletes an alert rule
func (m *Manager) DeleteAlertRule(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.rules[id]; !exists {
		return fmt.Errorf("alert rule not found: %s", id)
	}

	delete(m.rules, id)
	m.log.Infof("Alert rule deleted: %s", id)
	return nil
}

// EnableAlertRule enables an alert rule
func (m *Manager) EnableAlertRule(id string) (*models.AlertRule, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rule, exists := m.rules[id]
	if !exists {
		return nil, fmt.Errorf("alert rule not found: %s", id)
	}

	rule.Enabled = true
	rule.UpdatedAt = time.Now()

	m.log.Infof("Alert rule enabled: %s", rule.Name)
	return rule, nil
}

// DisableAlertRule disables an alert rule
func (m *Manager) DisableAlertRule(id string) (*models.AlertRule, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rule, exists := m.rules[id]
	if !exists {
		return nil, fmt.Errorf("alert rule not found: %s", id)
	}

	rule.Enabled = false
	rule.UpdatedAt = time.Now()

	m.log.Infof("Alert rule disabled: %s", rule.Name)
	return rule, nil
}

// CreateNotificationChannel creates a notification channel
func (m *Manager) CreateNotificationChannel(ctx context.Context, channel *models.NotificationChannel) (*models.NotificationChannel, error) {
	channel.ID = uuid.New().String()
	channel.Enabled = true
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	m.mutex.Lock()
	m.channels[channel.ID] = channel
	m.mutex.Unlock()

	m.log.Infof("Notification channel created: %s (%s)", channel.Name, channel.Type)
	return channel, nil
}

// GetNotificationChannels retrieves all notification channels
func (m *Manager) GetNotificationChannels() []models.NotificationChannel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var channels []models.NotificationChannel
	for _, channel := range m.channels {
		channels = append(channels, *channel)
	}

	return channels
}

// UpdateNotificationChannel updates a notification channel
func (m *Manager) UpdateNotificationChannel(ctx context.Context, id string, updates *models.NotificationChannel) (*models.NotificationChannel, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	channel, exists := m.channels[id]
	if !exists {
		return nil, fmt.Errorf("notification channel not found: %s", id)
	}

	if updates.Name != "" {
		channel.Name = updates.Name
	}
	if updates.Config != nil {
		channel.Config = updates.Config
	}
	if updates.Enabled != channel.Enabled {
		channel.Enabled = updates.Enabled
	}

	channel.UpdatedAt = time.Now()

	m.log.Infof("Notification channel updated: %s", channel.Name)
	return channel, nil
}

// DeleteNotificationChannel deletes a notification channel
func (m *Manager) DeleteNotificationChannel(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.channels[id]; !exists {
		return fmt.Errorf("notification channel not found: %s", id)
	}

	delete(m.channels, id)
	m.log.Infof("Notification channel deleted: %s", id)
	return nil
}

// TestNotificationChannel tests a notification channel
func (m *Manager) TestNotificationChannel(ctx context.Context, id string) error {
	m.mutex.RLock()
	channel, exists := m.channels[id]
	m.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("notification channel not found: %s", id)
	}

	testAlert := &models.Alert{
		Title:       "Test Alert",
		Description: "This is a test alert from Cloud Sentinel",
		Severity:    models.SeverityInfo,
		Source:      "test",
	}

	return m.sendToChannel(testAlert, channel)
}

// CreateSilence creates a silence
func (m *Manager) CreateSilence(ctx context.Context, silence *models.Silence) (*models.Silence, error) {
	silence.ID = uuid.New().String()
	silence.CreatedAt = time.Now()

	m.mutex.Lock()
	m.silences[silence.ID] = silence
	m.mutex.Unlock()

	m.log.Infof("Silence created by %s", silence.CreatedBy)
	return silence, nil
}

// GetSilences retrieves active silences
func (m *Manager) GetSilences() []models.Silence {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	now := time.Now()
	var silences []models.Silence
	for _, silence := range m.silences {
		if silence.EndsAt.After(now) {
			silences = append(silences, *silence)
		}
	}

	return silences
}

// DeleteSilence deletes a silence
func (m *Manager) DeleteSilence(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.silences[id]; !exists {
		return fmt.Errorf("silence not found: %s", id)
	}

	delete(m.silences, id)
	m.log.Infof("Silence deleted: %s", id)
	return nil
}

// isSilenced checks if an alert matches any active silence
func (m *Manager) isSilenced(alert *models.Alert) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	now := time.Now()
	for _, silence := range m.silences {
		if silence.EndsAt.Before(now) {
			continue
		}

		// Check if alert matches silence matchers
		matches := true
		for key, pattern := range silence.Matchers {
			var value string
			switch key {
			case "alertname":
				value = alert.Title
			case "severity":
				value = string(alert.Severity)
			case "source":
				value = alert.Source
			default:
				value = alert.Labels[key]
			}

			if !matchPattern(value, pattern) {
				matches = false
				break
			}
		}

		if matches {
			return true
		}
	}

	return false
}

// sendNotifications sends notifications to configured channels
func (m *Manager) sendNotifications(alert *models.Alert) {
	// Find rule for this alert
	var rule *models.AlertRule
	if alert.RuleID != "" {
		rule, _ = m.GetAlertRule(alert.RuleID)
	}

	// Get channel IDs from rule or use defaults
	var channelIDs []string
	if rule != nil && len(rule.ChannelIDs) > 0 {
		channelIDs = rule.ChannelIDs
	} else {
		// Use all enabled channels
		for id, ch := range m.channels {
			if ch.Enabled {
				channelIDs = append(channelIDs, id)
			}
		}
	}

	// Send to each channel
	for _, channelID := range channelIDs {
		channel, exists := m.channels[channelID]
		if !exists || !channel.Enabled {
			continue
		}

		if err := m.sendToChannel(alert, channel); err != nil {
			m.log.Errorf("Failed to send notification to %s: %v", channel.Name, err)
			notificationsTotal.WithLabelValues(string(channel.Type), "failed").Inc()
		} else {
			notificationsTotal.WithLabelValues(string(channel.Type), "success").Inc()
		}
	}
}

// sendToChannel sends an alert to a specific channel
func (m *Manager) sendToChannel(alert *models.Alert, channel *models.NotificationChannel) error {
	switch channel.Type {
	case models.ChannelSlack:
		return m.sendSlackNotification(alert, channel)
	case models.ChannelEmail:
		return m.sendEmailNotification(alert, channel)
	case models.ChannelPagerDuty:
		return m.sendPagerDutyNotification(alert, channel)
	case models.ChannelWebhook:
		return m.sendWebhookNotification(alert, channel)
	default:
		return fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
}

// sendSlackNotification sends a Slack notification
func (m *Manager) sendSlackNotification(alert *models.Alert, channel *models.NotificationChannel) error {
	if m.slackClient == nil {
		return fmt.Errorf("Slack client not configured")
	}

	webhookURL, _ := channel.Config["webhook_url"].(string)
	if webhookURL == "" {
		webhookURL = m.config.SlackWebhookURL
	}

	color := "#36a64f" // green
	switch alert.Severity {
	case models.SeverityCritical:
		color = "#ff0000" // red
	case models.SeverityWarning:
		color = "#ff9900" // orange
	}

	attachment := slack.Attachment{
		Color:      color,
		Title:      alert.Title,
		Text:       alert.Description,
		Footer:     "Cloud Sentinel",
		Timestamp:  alert.StartedAt.Unix(),
		Fields: []slack.AttachmentField{
			{
				Title: "Severity",
				Value: string(alert.Severity),
				Short: true,
			},
			{
				Title: "Source",
				Value: alert.Source,
				Short: true,
			},
		},
	}

	if webhookURL != "" {
		msg := &slack.WebhookMessage{
			Attachments: []slack.Attachment{attachment},
		}
		return slack.PostWebhook(webhookURL, msg)
	}

	channelName, _ := channel.Config["channel"].(string)
	if channelName == "" {
		channelName = m.config.SlackDefaultChannel
	}

	_, _, err := m.slackClient.PostMessage(
		channelName,
		slack.MsgOptionAttachments(attachment),
	)

	return err
}

// sendEmailNotification sends an email notification
func (m *Manager) sendEmailNotification(alert *models.Alert, channel *models.NotificationChannel) error {
	// Placeholder for email implementation
	// In production, use net/smtp or email service
	m.log.Debugf("Would send email for alert: %s", alert.Title)
	return nil
}

// sendPagerDutyNotification sends a PagerDuty notification
func (m *Manager) sendPagerDutyNotification(alert *models.Alert, channel *models.NotificationChannel) error {
	integrationKey, _ := channel.Config["integration_key"].(string)
	if integrationKey == "" {
		integrationKey = m.config.PagerDutyIntegrationKey
	}

	if integrationKey == "" {
		return fmt.Errorf("PagerDuty integration key not configured")
	}

	payload := map[string]interface{}{
		"routing_key":  integrationKey,
		"event_action": "trigger",
		"dedup_key":    alert.ID,
		"payload": map[string]interface{}{
			"summary":   alert.Title,
			"severity":  strings.ToLower(string(alert.Severity)),
			"source":    alert.Source,
			"component": alert.Labels["service"],
			"custom_details": map[string]interface{}{
				"description": alert.Description,
				"value":       alert.Value,
				"threshold":   alert.Threshold,
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := m.httpClient.Post(
		"https://events.pagerduty.com/v2/enqueue",
		"application/json",
		strings.NewReader(string(jsonPayload)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("pagerduty returned status %d", resp.StatusCode)
	}

	return nil
}

// sendWebhookNotification sends a webhook notification
func (m *Manager) sendWebhookNotification(alert *models.Alert, channel *models.NotificationChannel) error {
	webhookURL, _ := channel.Config["url"].(string)
	if webhookURL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	payload := map[string]interface{}{
		"id":          alert.ID,
		"title":       alert.Title,
		"description": alert.Description,
		"severity":    alert.Severity,
		"status":      alert.Status,
		"source":      alert.Source,
		"value":       alert.Value,
		"threshold":   alert.Threshold,
		"started_at":  alert.StartedAt,
		"labels":      alert.Labels,
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := m.httpClient.Post(
		webhookURL,
		"application/json",
		strings.NewReader(string(jsonPayload)),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// loadDefaultRules loads default alert rules
func (m *Manager) loadDefaultRules() {
	defaults := []models.AlertRule{
		{
			Name:        "High CPU Usage",
			Description: "CPU usage exceeds 80% for 5 minutes",
			Query:       "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100) > 80",
			Duration:    5 * time.Minute,
			Severity:    models.SeverityWarning,
			Threshold:   80,
			Operator:    "gt",
			Source:      "prometheus",
			Interval:    1 * time.Minute,
			Group:       "infrastructure",
			Labels: map[string]string{
				"team": "sre",
			},
			Annotations: map[string]string{
				"summary":     "High CPU usage detected",
				"description": "CPU usage is above 80% on {{ $labels.instance }}",
				"runbook_url": "https://wiki.internal/runbooks/high-cpu",
			},
		},
		{
			Name:        "High Memory Usage",
			Description: "Memory usage exceeds 90% for 5 minutes",
			Query:       "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 90",
			Duration:    5 * time.Minute,
			Severity:    models.SeverityCritical,
			Threshold:   90,
			Operator:    "gt",
			Source:      "prometheus",
			Interval:    1 * time.Minute,
			Group:       "infrastructure",
			Labels: map[string]string{
				"team": "sre",
			},
			Annotations: map[string]string{
				"summary":     "High memory usage detected",
				"description": "Memory usage is above 90% on {{ $labels.instance }}",
				"runbook_url": "https://wiki.internal/runbooks/high-memory",
			},
		},
		{
			Name:        "Disk Space Low",
			Description: "Disk space usage exceeds 85%",
			Query:       "(node_filesystem_size_bytes - node_filesystem_avail_bytes) / node_filesystem_size_bytes * 100 > 85",
			Duration:    5 * time.Minute,
			Severity:    models.SeverityWarning,
			Threshold:   85,
			Operator:    "gt",
			Source:      "prometheus",
			Interval:    5 * time.Minute,
			Group:       "infrastructure",
			Labels: map[string]string{
				"team": "sre",
			},
			Annotations: map[string]string{
				"summary":     "Low disk space",
				"description": "Disk usage is above 85% on {{ $labels.instance }}",
			},
		},
		{
			Name:        "Pod Crash Looping",
			Description: "Pod is restarting frequently",
			Query:       "rate(kube_pod_container_status_restarts_total[15m]) > 0",
			Duration:    5 * time.Minute,
			Severity:    models.SeverityCritical,
			Threshold:   0,
			Operator:    "gt",
			Source:      "prometheus",
			Interval:    1 * time.Minute,
			Group:       "kubernetes",
			Labels: map[string]string{
				"team": "platform",
			},
			Annotations: map[string]string{
				"summary":     "Pod crash looping",
				"description": "Pod {{ $labels.pod }} in namespace {{ $labels.namespace }} is crash looping",
			},
		},
		{
			Name:        "Service Down",
			Description: "Service endpoint is not responding",
			Query:       "up == 0",
			Duration:    1 * time.Minute,
			Severity:    models.SeverityCritical,
			Threshold:   0,
			Operator:    "eq",
			Source:      "prometheus",
			Interval:    30 * time.Second,
			Group:       "availability",
			Labels: map[string]string{
				"team": "sre",
			},
			Annotations: map[string]string{
				"summary":     "Service is down",
				"description": "Service {{ $labels.job }} on {{ $labels.instance }} is down",
			},
		},
	}

	for _, rule := range defaults {
		rule.ID = uuid.New().String()
		rule.Enabled = true
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		m.rules[rule.ID] = &rule
	}

	m.log.Infof("Loaded %d default alert rules", len(defaults))
}

// loadDefaultChannels loads default notification channels
func (m *Manager) loadDefaultChannels() {
	// Slack channel (if configured)
	if m.config.SlackWebhookURL != "" || m.config.SlackBotToken != "" {
		channel := &models.NotificationChannel{
			Name:    "Default Slack",
			Type:    models.ChannelSlack,
			Enabled: true,
			Config: map[string]interface{}{
				"webhook_url": m.config.SlackWebhookURL,
				"channel":     m.config.SlackDefaultChannel,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		channel.ID = uuid.New().String()
		m.channels[channel.ID] = channel
		m.log.Info("Loaded default Slack notification channel")
	}

	// PagerDuty channel (if configured)
	if m.config.PagerDutyIntegrationKey != "" {
		channel := &models.NotificationChannel{
			Name:    "Default PagerDuty",
			Type:    models.ChannelPagerDuty,
			Enabled: true,
			Config: map[string]interface{}{
				"integration_key": m.config.PagerDutyIntegrationKey,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		channel.ID = uuid.New().String()
		m.channels[channel.ID] = channel
		m.log.Info("Loaded default PagerDuty notification channel")
	}
}

// matchPattern checks if a value matches a pattern (supports wildcards)
func matchPattern(value, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.Contains(pattern, "*") {
		// Simple wildcard matching
		prefix := strings.Split(pattern, "*")[0]
		suffix := strings.Split(pattern, "*")[1]
		return strings.HasPrefix(value, prefix) && strings.HasSuffix(value, suffix)
	}
	return value == pattern
}

// add adds an alert to history
func (h *AlertHistory) add(alert models.Alert) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.alerts = append(h.alerts, alert)
	if len(h.alerts) > h.limit {
		h.alerts = h.alerts[len(h.alerts)-h.limit:]
	}
}

// GetRecent returns recent alerts from history
func (h *AlertHistory) GetRecent(limit int) []models.Alert {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if limit > len(h.alerts) {
		limit = len(h.alerts)
	}

	result := make([]models.Alert, limit)
	copy(result, h.alerts[len(h.alerts)-limit:])

	// Reverse to get newest first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}
