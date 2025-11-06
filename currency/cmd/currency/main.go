package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	currencyClient "github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"github.com/Sevn9/currency-screener/currency/internal/config"
	"github.com/Sevn9/currency-screener/currency/internal/handler"
	"github.com/Sevn9/currency-screener/currency/internal/services"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("main: currency microservice start")
	if err := run(); err != nil {
		log.Fatal("main: " + err.Error())
	}
	fmt.Println("main: currency microservice end")
}

func run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("main: Recovery:", r)
		}
	}()

	ctx := context.Background()

	logger, _ := zap.NewProduction()
	//cброс (flush) всех буферизованных записей логов во внешнее хранилище
	defer logger.Sync()

	//загрузка конфигов
	//example: go run main.go -config=.../currency-screener/currency/internal/config/config.yaml
	configPath := flag.String("config", "../../internal/config/config.yaml", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		logger.Fatal("main: error loading config",
			zap.Error(err))
		return err
	}

	//temp: запрос текущего курса
	currClient, err := currencyClient.NewCurrencyClient(cfg.PublicCurrencyApi, logger)

	if err != nil {
		logger.Fatal("main: error NewCurrencyClient create",
			zap.Error(err))
		return err
	}

	currClientTemp, er := currClient.GetCurrentRate(ctx)
	if er != nil {
		logger.Fatal("main: error GetCurrentRate create:",
			zap.Error(err))
		return err
	}
	_ = currClientTemp

	//added services
	svc := services.NewCurrencyService(currClient, logger)

	//конфигурируем gRPC сервер
	currencyServer := handler.NewCurrencyServer(svc, logger)

	// запускаем gRPC сервер
	stopGRPC, errCh, err := startGRPCServer(cfg, currencyServer, logger)
	if err != nil {
		logger.Fatal("main: error starting GRPC server:",
			zap.Error(err))
		return err
	}

	//корректное завершение
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop: // пришёл сигнал на выключение
		logger.Info("main: Shutting down gRPC server...")
		return stopGRPC(5 * time.Second)

	case serveErr := <-errCh: // сервер упал сам
		if serveErr != nil {
			return serveErr
		}
	}

	return err
}

func startGRPCServer(cfg *config.AppConfig, currencyServer *handler.CurrencyServer, logger *zap.Logger) (
	func(timeout time.Duration) error, <-chan error, error) {
	lis, err := net.Listen("tcp", ":"+cfg.Service.Port)
	if err != nil {
		return nil, nil, fmt.Errorf("main: failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	//todo: регистрация сервисов
	currency.RegisterCurrencyServiceServer(grpcServer, currencyServer)

	errCh := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("main: Recovery:", r)
			}
		}()

		logger.Info(
			// Статическое сообщение (message)
			"main: gRPC server is listening",
			// Структурированное поле (Field) для порта
			zap.String("port", cfg.Service.Port),
			zap.String("protocol", "grpc"),
		)

		if err := grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("main: filed to serve: %w", err)
		}
		close(errCh)
	}()

	gracefullyShutdownFunc := func(timeout time.Duration) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		doneCh := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(doneCh)
		}()

		select {
		case <-doneCh:
			logger.Info("main: gRPC server stopped gracefully")
		case <-ctx.Done():
			logger.Info("main: Graceful stop timed out, forcing stop")
			grpcServer.Stop()
		}

		_ = lis.Close()

		return nil

	}

	return gracefullyShutdownFunc, errCh, nil
}
