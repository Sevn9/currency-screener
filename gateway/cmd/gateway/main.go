package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sevn9/currency-screener/gateway/internal/clients/auth"
	"github.com/Sevn9/currency-screener/gateway/internal/clients/redis"
	"github.com/Sevn9/currency-screener/gateway/internal/config"
	"github.com/Sevn9/currency-screener/gateway/internal/handler"
	"github.com/Sevn9/currency-screener/gateway/internal/middleware"
	"github.com/Sevn9/currency-screener/gateway/internal/repository"
	"github.com/Sevn9/currency-screener/gateway/internal/services"
	"github.com/Sevn9/currency-screener/pkg/grpc_client"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "currency_requests_total",
			Help: "Total number of requests handled by the currency service",
		},
		[]string{"method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "currency_request_duration_seconds",
			Help:    "Histogram of response times for requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	appUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "currency_service_uptime_seconds",
			Help: "Time since service start in seconds",
		},
	)
)

func init() {
	// metrics register
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(appUptime)
}

func main() {
	log.Println("main: gateway start")
	if err := run(); err != nil {
		log.Fatal("gateway main: " + err.Error())
	}
	log.Println("main: gateway server exited")
}

func run() error {
	logger, _ := zap.NewProduction()
	//flush all buffered log entries to external storage
	defer logger.Sync()

	//Loading configs
	configPath := flag.String("config", "../../internal/config/config.yaml", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		logger.Error("main: error loading config",
			zap.Error(err))
		return err
	}

	logger.Info("Gateway server initializing with config", zap.Any("config", cfg))

	// Added gin router
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// auth client
	authClient, err := auth.NewAuthClient(cfg.AuthApi)
	if err != nil {
		return fmt.Errorf("auth.NewAuthClient: %w", err)
	}

	resp, err := authClient.Ping()
	if err != nil {
		return fmt.Errorf("authClient.Ping: %w", err)
	}

	if resp != "pong" {
		return fmt.Errorf("auth client answered with invalid response: %w", err)
	}

	// currency client
	currencyClient, conn, err := grpc_client.NewCurrencyServiceClient(cfg.GrpcClientConfig.CurrencyServiceUrl)
	if err != nil {
		return fmt.Errorf("main: grpc_client.NewCurrencyServiceClient: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warn("main: Cannot close GRPC Client for auth service", zap.Error(err))
		}
	}()

	// redis client
	redisClient, err := redis.NewClient(
		cfg.RedisConfig.Host+cfg.RedisConfig.Port,
		cfg.RedisConfig.Password,
		cfg.RedisDbNums.CurrencyDb)
	if err != nil {
		return fmt.Errorf("main: redisClient create failed: %w", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Warn("main: Error closing Redis", zap.Error(err))
		}
	}()

	// add middleware
	authMiddleware := middleware.NewAuthorization(authClient, logger)

	//add repository
	userRepo := repository.NewUserRepository()
	redisRepo := repository.NewCurrencyRedisRepository(redisClient, 10*time.Minute)

	//add services
	authService := services.NewAuth(authClient, userRepo)
	currencyService := services.NewCurrency(currencyClient, redisRepo)

	// Route registration
	handler.RegisterRoutes(router, logger, authMiddleware, authService, currencyService)

	// Setting up the server
	srv := &http.Server{
		Addr:         cfg.Service.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// create metrics collector
	metricsCollector := middleware.NewMetricsMiddleware(
		requestCount,
		requestDuration,
		appUptime,
	)

	// create metrics server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	metricsSrv := &http.Server{
		Addr:    cfg.MetricsConfig.Port,
		Handler: metricsMux,
	}

	// Starting the server
	go func() {
		log.Printf("Starting Gateway server on %s", cfg.Service.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("main: error starting Gateway server:",
				zap.Error(err))
		}
	}()

	// Starting the server metrics Prometheus
	go func() {
		logger.Info("Starting Metrics server", zap.String("port", cfg.MetricsConfig.Port))
		if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Metrics server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	errCh := make(chan error, 1)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	startTime := time.Now()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			logger.Info("Shutting down servers...")

			srvCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// shutdown main server
			if err := srv.Shutdown(srvCtx); err != nil {
				logger.Error("main: Gateway server forced to shutdown", zap.Error(err))
			} else {
				logger.Info("main: Gateway server stopped gracefully")
			}

			metricsCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// shutdown metrics server
			if err := metricsSrv.Shutdown(metricsCtx); err != nil {
				logger.Error("main: Metrics server forced to shutdown", zap.Error(err))
			} else {
				logger.Info("main: Metrics server stopped gracefully")
			}

			return nil

		case serveErr := <-errCh: // panic or port cant access
			logger.Error("Server crashed", zap.Error(serveErr))
			return serveErr

		case <-ticker.C:
			uptime := time.Since(startTime).Seconds()
			metricsCollector.SetUptime(uptime)
		}
	}
}
