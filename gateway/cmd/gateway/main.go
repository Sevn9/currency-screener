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
	"github.com/Sevn9/currency-screener/gateway/internal/config"
	"github.com/Sevn9/currency-screener/gateway/internal/handler"
	"github.com/Sevn9/currency-screener/gateway/internal/middleware"
	"github.com/Sevn9/currency-screener/gateway/internal/repository"
	"github.com/Sevn9/currency-screener/gateway/internal/services"
	"github.com/Sevn9/currency-screener/pkg/grpc_client"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
		//todo: up auth services
		//return fmt.Errorf("authClient.Ping: %w", err)
	}

	if resp != "pong" {
		//return fmt.Errorf("auth client answered with invalid response: %w", err)
	}

	// currency client
	currencyClient, conn, err := grpc_client.NewCurrencyServiceClient(cfg.GrpcClientConfig.CurrencyServiceUrl)
	if err != nil {
		return fmt.Errorf("grpc_client.NewCurrencyServiceClient: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warn("Cannot close GRPC Client for auth service", zap.Error(err))
		}
	}()

	// add middleware
	authMiddleware := middleware.NewAuthorization(authClient, logger)

	//add repository
	userRepo := repository.NewUserRepository()

	//add services
	authService := services.NewAuth(authClient, userRepo)
	currencyService := services.NewCurrency(currencyClient)

	// Route registration
	handler.RegisterRoutes(router, logger, authMiddleware, authService, currencyService)

	// Setting up the server
	srv := &http.Server{
		Addr:         cfg.Service.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Starting the server
	go func() {
		log.Printf("Starting Gateway server on %s", cfg.Service.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("srv.ListenAndServe: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	return nil
}
