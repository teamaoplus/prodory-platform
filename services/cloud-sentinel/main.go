package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "cloud-sentinel/docs"
	"cloud-sentinel/internal/alerts"
	"cloud-sentinel/internal/api"
	"cloud-sentinel/internal/config"
	"cloud-sentinel/internal/metrics"
	"cloud-sentinel/internal/scheduler"
)

// @title Cloud Sentinel API
// @version 1.0.0
// @description Multi-cloud monitoring and alerting service with Prometheus integration
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@prodory.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8083
// @BasePath /api/v1
// @schemes http https

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set log level
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	log.Info("Starting Cloud Sentinel...")

	// Initialize metrics collector
	metricsCollector := metrics.NewCollector(log)

	// Initialize alert manager
	alertManager, err := alerts.NewManager(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize alert manager: %v", err)
	}

	// Initialize scheduler
	sched := scheduler.New(cfg, log, metricsCollector, alertManager)
	sched.Start()
	defer sched.Stop()

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(requestLogger(log))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "cloud-sentinel",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC(),
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	apiV1 := router.Group("/api/v1")
	{
		// Dashboard routes
		dashboard := apiV1.Group("/dashboard")
		{
			dashboard.GET("/summary", api.GetDashboardSummary(metricsCollector))
			dashboard.GET("/metrics", api.GetMetricsOverview(metricsCollector))
			dashboard.GET("/alerts", api.GetActiveAlerts(alertManager))
		}

		// Metrics routes
		metricsGroup := apiV1.Group("/metrics")
		{
			metricsGroup.GET("/aws", api.GetAWSMetrics(metricsCollector))
			metricsGroup.GET("/azure", api.GetAzureMetrics(metricsCollector))
			metricsGroup.GET("/gcp", api.GetGCPMetrics(metricsCollector))
			metricsGroup.GET("/kubernetes", api.GetKubernetesMetrics(metricsCollector))
			metricsGroup.GET("/custom", api.GetCustomMetrics(metricsCollector))
			metricsGroup.POST("/custom", api.CreateCustomMetric(metricsCollector))
		}

		// Alerts routes
		alertsGroup := apiV1.Group("/alerts")
		{
			alertsGroup.GET("", api.GetAlerts(alertManager))
			alertsGroup.GET("/:id", api.GetAlert(alertManager))
			alertsGroup.POST("", api.CreateAlert(alertManager))
			alertsGroup.PUT("/:id", api.UpdateAlert(alertManager))
			alertsGroup.DELETE("/:id", api.DeleteAlert(alertManager))
			alertsGroup.POST("/:id/acknowledge", api.AcknowledgeAlert(alertManager))
			alertsGroup.POST("/:id/resolve", api.ResolveAlert(alertManager))
		}

		// Alert rules routes
		rules := apiV1.Group("/rules")
		{
			rules.GET("", api.GetAlertRules(alertManager))
			rules.GET("/:id", api.GetAlertRule(alertManager))
			rules.POST("", api.CreateAlertRule(alertManager))
			rules.PUT("/:id", api.UpdateAlertRule(alertManager))
			rules.DELETE("/:id", api.DeleteAlertRule(alertManager))
			rules.POST("/:id/enable", api.EnableAlertRule(alertManager))
			rules.POST("/:id/disable", api.DisableAlertRule(alertManager))
		}

		// Notification channels routes
		channels := apiV1.Group("/channels")
		{
			channels.GET("", api.GetNotificationChannels(alertManager))
			channels.POST("", api.CreateNotificationChannel(alertManager))
			channels.PUT("/:id", api.UpdateNotificationChannel(alertManager))
			channels.DELETE("/:id", api.DeleteNotificationChannel(alertManager))
			channels.POST("/:id/test", api.TestNotificationChannel(alertManager))
		}

		// Silences routes
		silences := apiV1.Group("/silences")
		{
			silences.GET("", api.GetSilences(alertManager))
			silences.POST("", api.CreateSilence(alertManager))
			silences.DELETE("/:id", api.DeleteSilence(alertManager))
		}

		// Reports routes
		reports := apiV1.Group("/reports")
		{
			reports.GET("", api.GetReports())
			reports.POST("/generate", api.GenerateReport())
			reports.GET("/:id/download", api.DownloadReport())
		}

		// Settings routes
		settings := apiV1.Group("/settings")
		{
			settings.GET("", api.GetSettings())
			settings.PUT("", api.UpdateSettings())
		}

		// Discovery routes
		discovery := apiV1.Group("/discovery")
		{
			discovery.POST("/scan", api.RunDiscovery())
			discovery.GET("/resources", api.GetDiscoveredResources())
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Infof("Cloud Sentinel started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Cloud Sentinel...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Cloud Sentinel stopped")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func requestLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
		}).Info("Request")
	}
}
