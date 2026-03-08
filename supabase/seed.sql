-- Seed data for Prodory Platform
-- This file contains sample data for development and testing

-- Sample cloud accounts (for testing only - use demo credentials)
INSERT INTO cloud_accounts (user_id, provider, account_id, account_name, credentials, is_active)
VALUES 
    ('00000000-0000-0000-0000-000000000000', 'aws', '123456789012', 'Production AWS', '{"role_arn": "arn:aws:iam::123456789012:role/FinOpsRole"}', true),
    ('00000000-0000-0000-0000-000000000000', 'azure', 'sub-12345', 'Production Azure', '{"subscription_id": "sub-12345"}', true),
    ('00000000-0000-0000-0000-000000000000', 'gcp', 'proj-12345', 'Production GCP', '{"project_id": "proj-12345"}', true);

-- Sample budgets
INSERT INTO budgets (user_id, name, amount, currency, period, start_date, alert_thresholds)
VALUES 
    ('00000000-0000-0000-0000-000000000000', 'Monthly Cloud Budget', 50000.00, 'USD', 'monthly', '2024-01-01', ARRAY[50, 80, 100]),
    ('00000000-0000-0000-0000-000000000000', 'Q1 Development Budget', 25000.00, 'USD', 'quarterly', '2024-01-01', ARRAY[75, 90, 100]);

-- Sample alert rules
INSERT INTO alert_rules (user_id, name, description, query, severity, threshold, operator)
VALUES 
    ('00000000-0000-0000-0000-000000000000', 'High CPU Usage', 'CPU usage exceeds 80%', 'cpu_usage_percent > 80', 'warning', 80, 'gt'),
    ('00000000-0000-0000-0000-000000000000', 'High Memory Usage', 'Memory usage exceeds 90%', 'memory_usage_percent > 90', 'critical', 90, 'gt'),
    ('00000000-0000-0000-0000-000000000000', 'Cost Spike', 'Daily cost exceeds budget', 'daily_cost > 2000', 'warning', 2000, 'gt');

-- Sample notification channels
INSERT INTO notification_channels (user_id, name, type, config, enabled)
VALUES 
    ('00000000-0000-0000-0000-000000000000', 'DevOps Slack', 'slack', '{"webhook_url": "https://hooks.slack.com/services/xxx"}', true),
    ('00000000-0000-0000-0000-000000000000', 'On-Call Email', 'email', '{"recipients": ["oncall@company.com"]}', true);

-- Sample recommendations
INSERT INTO recommendations (user_id, type, category, service_name, title, description, estimated_savings, effort, impact, confidence_score)
VALUES 
    ('00000000-0000-0000-0000-000000000000', 'rightsizing', 'compute', 'EC2', 'Right-size over-provisioned instances', 'Resize EC2 instances based on actual usage', 1250.00, 'low', 'high', 95.5),
    ('00000000-0000-0000-0000-000000000000', 'reserved_instances', 'compute', 'EC2', 'Purchase Reserved Instances', 'Convert On-Demand to Reserved Instances', 5000.00, 'medium', 'high', 88.0),
    ('00000000-0000-0000-0000-000000000000', 'storage_optimization', 'storage', 'S3', 'Enable S3 Intelligent Tiering', 'Move infrequently accessed data to cheaper tiers', 350.00, 'low', 'medium', 92.0);
