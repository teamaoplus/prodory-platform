module cloud-sentinel

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/common v0.45.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.17.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
	github.com/lib/pq v1.10.9
	github.com/robfig/cron/v3 v3.0.1
	github.com/slack-go/slack v0.12.3
	github.com/aws/aws-sdk-go-v2 v1.23.0
	github.com/aws/aws-sdk-go-v2/config v1.25.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.30.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery v1.1.0
	cloud.google.com/go/monitoring v1.16.3
	google.golang.org/api v0.151.0
	github.com/supabase-community/supabase-go v0.0.4
	k8s.io/client-go v0.28.4
	k8s.io/api v0.28.4
	k8s.io/apimachinery v0.28.4
)
