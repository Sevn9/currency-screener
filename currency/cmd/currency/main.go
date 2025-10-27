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

	"github.com/Sevn9/currency-screener/currency/internal/config"
	"github.com/Sevn9/currency-screener/currency/internal/handler"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("currency microservice start")
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("currency microservice end")
}

func run() error {
	//загрузка конфигов
	//example: go run main.go -config=.../currency-screener/currency/internal/config/config.yaml
	configPath := flag.String("config", "../../internal/config/config.yaml", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	//конфигурируем gRPC сервер
	currencyServer := handler.NewCurrencyServer()

	// запускаем gRPC сервер
	stopGRPC, errCh, err := startGRPCServer(cfg, currencyServer)
	if err != nil {
		log.Fatalf("Error starting GRPC server: %s", err)
		return err
	}

	//корректное завершение
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop: // пришёл сигнал на выключение
		log.Println("Shutting down gRPC server...")
		return stopGRPC(5 * time.Second)

	case serveErr := <-errCh: // сервер упал сам
		if serveErr != nil {
			return serveErr
		}
	}

	return err
}

func startGRPCServer(cfg *config.AppConfig, currencyServer handler.CurrencyServer) (func(timeout time.Duration) error, <-chan error, error) {
	lis, err := net.Listen("tcp", ":"+cfg.Service.Port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	//todo: регистрация сервисов

	errCh := make(chan error, 1)

	go func() {
		log.Printf("gRPC server is listening on :%s", cfg.Service.Port)

		if err := grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("filed to serve: %w", err)
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
			log.Println("gRPC server stopped gracefully")
		case <-ctx.Done():
			log.Println("Graceful stop timed out, forcing stop")
			grpcServer.Stop()
		}

		_ = lis.Close()

		return nil

	}

	return gracefullyShutdownFunc, errCh, nil
}
