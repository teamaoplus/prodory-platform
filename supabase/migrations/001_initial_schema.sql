-- Prodory Platform Initial Schema
-- This migration creates the core tables for all services

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- FinOps Dashboard Tables
-- ============================================

-- Cloud Provider Accounts
CREATE TABLE cloud_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('aws', 'azure', 'gcp')),
    account_id VARCHAR(255) NOT NULL,
    account_name VARCHAR(255),
    credentials JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, provider, account_id)
);

-- Cost Data
CREATE TABLE cost_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES cloud_accounts(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    resource_id VARCHAR(500),
    region VARCHAR(100),
    cost_amount DECIMAL(15, 4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    usage_quantity DECIMAL(15, 4),
    usage_unit VARCHAR(50),
    charge_type VARCHAR(50), -- OnDemand, Reserved, Spot
    date DATE NOT NULL,
    tags JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Budgets
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'quarterly', 'yearly')),
    start_date DATE NOT NULL,
    end_date DATE,
    alert_thresholds INTEGER[] DEFAULT ARRAY[50, 80, 100],
    alert_channels JSONB DEFAULT '[]',
    is_active BOOLEAN DEFAULT true,
    current_spend DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Cost Anomalies
CREATE TABLE cost_anomalies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES cloud_accounts(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    expected_cost DECIMAL(15, 4) NOT NULL,
    actual_cost DECIMAL(15, 4) NOT NULL,
    difference_amount DECIMAL(15, 4) NOT NULL,
    difference_percent DECIMAL(10, 2) NOT NULL,
    anomaly_date DATE NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'investigating', 'resolved', 'false_positive')),
    description TEXT,
    root_cause TEXT,
    resolution TEXT,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- Recommendations
CREATE TABLE recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    account_id UUID REFERENCES cloud_accounts(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- rightsizing, reserved_instances, spot_instances, etc.
    category VARCHAR(100) NOT NULL,
    service_name VARCHAR(255),
    resource_id VARCHAR(500),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    current_state JSONB,
    recommended_state JSONB,
    estimated_savings DECIMAL(15, 4),
    savings_currency VARCHAR(3) DEFAULT 'USD',
    effort VARCHAR(20) CHECK (effort IN ('low', 'medium', 'high')),
    impact VARCHAR(20) CHECK (impact IN ('low', 'medium', 'high')),
    confidence_score DECIMAL(5, 2),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'implemented', 'dismissed', 'snoozed')),
    implemented_at TIMESTAMP WITH TIME ZONE,
    actual_savings DECIMAL(15, 4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- Cloud Sentinel Tables
-- ============================================

-- Monitored Resources
CREATE TABLE monitored_resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(500) NOT NULL,
    resource_name VARCHAR(500),
    region VARCHAR(100),
    account_id VARCHAR(255),
    labels JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, provider, resource_id)
);

-- Alert Rules
CREATE TABLE alert_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    query TEXT NOT NULL,
    duration_minutes INTEGER DEFAULT 5,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
    labels JSONB DEFAULT '{}',
    annotations JSONB DEFAULT '{}',
    threshold DECIMAL(15, 4),
    operator VARCHAR(10) CHECK (operator IN ('gt', 'lt', 'eq', 'ne')),
    source VARCHAR(50) DEFAULT 'prometheus',
    check_interval_seconds INTEGER DEFAULT 60,
    runbook_url TEXT,
    notification_channels UUID[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Alerts
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES alert_rules(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'firing' CHECK (status IN ('firing', 'acknowledged', 'resolved', 'suppressed')),
    source VARCHAR(100),
    labels JSONB DEFAULT '{}',
    value DECIMAL(15, 4),
    threshold DECIMAL(15, 4),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by UUID REFERENCES auth.users(id),
    runbook_url TEXT,
    annotations JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Notification Channels
CREATE TABLE notification_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('email', 'slack', 'pagerduty', 'webhook', 'teams')),
    config JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Metric Data
CREATE TABLE metric_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    labels JSONB DEFAULT '{}',
    value DECIMAL(15, 4) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    source VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for metric data (if TimescaleDB is available)
-- SELECT create_hypertable('metric_data', 'timestamp');

-- ============================================
-- Kubernetes-in-a-Box Tables
-- ============================================

-- Clusters
CREATE TABLE k8s_clusters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    region VARCHAR(100),
    version VARCHAR(50),
    k3s_version VARCHAR(50),
    masters INTEGER DEFAULT 1,
    workers INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'creating',
    endpoint TEXT,
    kubeconfig TEXT,
    ssh_key_path TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, name)
);

-- Cluster Nodes
CREATE TABLE k8s_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cluster_id UUID REFERENCES k8s_clusters(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('master', 'worker')),
    instance_type VARCHAR(100),
    public_ip INET,
    private_ip INET,
    status VARCHAR(50) DEFAULT 'pending',
    labels JSONB DEFAULT '{}',
    taints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- Migration Tables
-- ============================================

-- VM Migrations
CREATE TABLE vm_migrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN ('vmware', 'virtualbox', 'hyperv', 'raw')),
    source_vm_name VARCHAR(255) NOT NULL,
    source_vm_id VARCHAR(500),
    target_type VARCHAR(50) NOT NULL CHECK (target_type IN ('kubevirt', 'container')),
    target_name VARCHAR(255) NOT NULL,
    target_namespace VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'analyzing', 'converting', 'transferring', 'completed', 'failed', 'rolled_back')),
    analysis_report JSONB,
    disk_size_gb INTEGER,
    transferred_bytes BIGINT DEFAULT 0,
    progress_percent INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- Indexes
-- ============================================

-- Cost data indexes
CREATE INDEX idx_cost_data_account_date ON cost_data(account_id, date);
CREATE INDEX idx_cost_data_service ON cost_data(service_name);
CREATE INDEX idx_cost_data_resource ON cost_data(resource_id);

-- Alert indexes
CREATE INDEX idx_alerts_user_status ON alerts(user_id, status);
CREATE INDEX idx_alerts_rule ON alerts(rule_id);
CREATE INDEX idx_alerts_started ON alerts(started_at);

-- Metric indexes
CREATE INDEX idx_metrics_user_name ON metric_data(user_id, metric_name);
CREATE INDEX idx_metrics_timestamp ON metric_data(timestamp);

-- Resource indexes
CREATE INDEX idx_resources_user_provider ON monitored_resources(user_id, provider);
CREATE INDEX idx_resources_type ON monitored_resources(resource_type);

-- ============================================
-- Row Level Security Policies
-- ============================================

-- Enable RLS on all tables
ALTER TABLE cloud_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE cost_data ENABLE ROW LEVEL SECURITY;
ALTER TABLE budgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE cost_anomalies ENABLE ROW LEVEL SECURITY;
ALTER TABLE recommendations ENABLE ROW LEVEL SECURITY;
ALTER TABLE monitored_resources ENABLE ROW LEVEL SECURITY;
ALTER TABLE alert_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE alerts ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_channels ENABLE ROW LEVEL SECURITY;
ALTER TABLE metric_data ENABLE ROW LEVEL SECURITY;
ALTER TABLE k8s_clusters ENABLE ROW LEVEL SECURITY;
ALTER TABLE k8s_nodes ENABLE ROW LEVEL SECURITY;
ALTER TABLE vm_migrations ENABLE ROW LEVEL SECURITY;

-- Create policies
CREATE POLICY "Users can only see their own cloud accounts"
    ON cloud_accounts FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own cost data"
    ON cost_data FOR ALL USING (auth.uid() = (SELECT user_id FROM cloud_accounts WHERE id = account_id));

CREATE POLICY "Users can only see their own budgets"
    ON budgets FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own anomalies"
    ON cost_anomalies FOR ALL USING (auth.uid() = (SELECT user_id FROM cloud_accounts WHERE id = account_id));

CREATE POLICY "Users can only see their own recommendations"
    ON recommendations FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own resources"
    ON monitored_resources FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own alert rules"
    ON alert_rules FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own alerts"
    ON alerts FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own notification channels"
    ON notification_channels FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own metrics"
    ON metric_data FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own clusters"
    ON k8s_clusters FOR ALL USING (auth.uid() = user_id);

CREATE POLICY "Users can only see their own cluster nodes"
    ON k8s_nodes FOR ALL USING (auth.uid() = (SELECT user_id FROM k8s_clusters WHERE id = cluster_id));

CREATE POLICY "Users can only see their own migrations"
    ON vm_migrations FOR ALL USING (auth.uid() = user_id);

-- ============================================
-- Functions and Triggers
-- ============================================

-- Update timestamp function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create update triggers
CREATE TRIGGER update_cloud_accounts_updated_at BEFORE UPDATE ON cloud_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recommendations_updated_at BEFORE UPDATE ON recommendations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notification_channels_updated_at BEFORE UPDATE ON notification_channels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_k8s_clusters_updated_at BEFORE UPDATE ON k8s_clusters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_k8s_nodes_updated_at BEFORE UPDATE ON k8s_nodes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vm_migrations_updated_at BEFORE UPDATE ON vm_migrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
