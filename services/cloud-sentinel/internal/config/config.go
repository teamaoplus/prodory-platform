package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for Cloud Sentinel
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	// Database
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	// Supabase
	SupabaseURL    string `mapstructure:"SUPABASE_URL"`
	SupabaseKey    string `mapstructure:"SUPABASE_KEY"`
	SupabaseJWT    string `mapstructure:"SUPABASE_JWT_SECRET"`

	// Prometheus
	PrometheusURL string `mapstructure:"PROMETHEUS_URL"`

	// AWS
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`

	// Azure
	AzureTenantID       string `mapstructure:"AZURE_TENANT_ID"`
	AzureClientID       string `mapstructure:"AZURE_CLIENT_ID"`
	AzureClientSecret   string `mapstructure:"AZURE_CLIENT_SECRET"`
	AzureSubscriptionID string `mapstructure:"AZURE_SUBSCRIPTION_ID"`

	// GCP
	GCPProjectID          string `mapstructure:"GCP_PROJECT_ID"`
	GCPServiceAccountKey  string `mapstructure:"GCP_SERVICE_ACCOUNT_KEY"`

	// Kubernetes
	KubernetesConfigPath string `mapstructure:"KUBERNETES_CONFIG_PATH"`
	InCluster            bool   `mapstructure:"IN_CLUSTER"`

	// Alerting
	DefaultAlertInterval int `mapstructure:"DEFAULT_ALERT_INTERVAL"`
	MaxAlertsPerMinute   int `mapstructure:"MAX_ALERTS_PER_MINUTE"`
	AlertRetentionDays   int `mapstructure:"ALERT_RETENTION_DAYS"`

	// Slack
	SlackWebhookURL   string `mapstructure:"SLACK_WEBHOOK_URL"`
	SlackBotToken     string `mapstructure:"SLACK_BOT_TOKEN"`
	SlackDefaultChannel string `mapstructure:"SLACK_DEFAULT_CHANNEL"`

	// PagerDuty
	PagerDutyIntegrationKey string `mapstructure:"PAGERDUTY_INTEGRATION_KEY"`

	// Email
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`

	// Webhook
	WebhookURL string `mapstructure:"WEBHOOK_URL"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	viper.SetEnvPrefix("SENTINEL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8083")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("PROMETHEUS_URL", "http://prometheus:9090")
	viper.SetDefault("DEFAULT_ALERT_INTERVAL", 60)
	viper.SetDefault("MAX_ALERTS_PER_MINUTE", 100)
	viper.SetDefault("ALERT_RETENTION_DAYS", 30)
	viper.SetDefault("IN_CLUSTER", false)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// IsCloudProviderEnabled checks if any cloud provider is configured
func (c *Config) IsCloudProviderEnabled() bool {
	return c.AWSAccessKeyID != "" || c.AzureClientID != "" || c.GCPProjectID != ""
}

// GetAWSRegions returns list of AWS regions to monitor
func (c *Config) GetAWSRegions() []string {
	if c.AWSRegion != "" {
		return []string{c.AWSRegion}
	}
	return []string{"us-east-1", "us-west-2", "eu-west-1"}
}
