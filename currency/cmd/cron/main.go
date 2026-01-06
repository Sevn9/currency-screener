package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/Sevn9/currency-screener/currency/internal/cache_repository"
	currencyClient "github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"github.com/Sevn9/currency-screener/currency/internal/config"
	"github.com/Sevn9/currency-screener/currency/internal/db"
	"github.com/Sevn9/currency-screener/currency/internal/repository"
	"github.com/Sevn9/currency-screener/currency/internal/services"
	"github.com/Sevn9/currency-screener/currency/internal/worker"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("main: currency cron start")
	if err := run(); err != nil {
		log.Fatal("cron main: " + err.Error())
	}
	fmt.Println("main: currency cron end")
}

func run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("main: Recovery:", r)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//create logger
	logger, _ := zap.NewProduction()
	//cброс (flush) всех буферизованных записей логов во внешнее хранилище
	defer logger.Sync()

	//загрузка конфигов
	//example: go run main.go -config=.../currency-screener/currency/internal/config/config.yaml
	configPath := flag.String("config", "../../internal/config/config.yaml", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)

	if err != nil {
		logger.Error("main: error loading config",
			zap.Error(err))
		return err
	}

	//added db connection
	pool, err := db.NewPgxPool(cfg.PostgresDb)
	if err != nil {
		logger.Error("main: failed to connect to pgx db",
			zap.Error(err))
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("main: failed ping to pgx db",
			zap.Error(err))
		return err
	}

	//added cache_repository
	cacheRepo := cache_repository.NewMemoryRateStorage()

	//added db repository
	dbRepo, err := repository.NewCurrencyPostgres(pool)

	if err != nil {
		logger.Error("main: failed create db",
			zap.Error(err))
		return err
	}

	//added clients
	currClient, err := currencyClient.NewCurrencyClient(cfg.PublicCurrencyApi, logger)

	if err != nil {
		logger.Error("main: error NewCurrencyClient create",
			zap.Error(err))
		return fmt.Errorf("main: error NewCurrencyClient create: %v", err)
	}

	//added services
	svc := services.NewCurrencyService(cacheRepo, dbRepo, currClient, logger)

	//added cron jobs
	cronJob := cron.New()

	//created workers
	currencyWorker := worker.NewCurrency(cfg.Worker, svc, cronJob, logger)

	if err := currencyWorker.StartFetchingCurrencyRates(); err != nil {
		return fmt.Errorf("error start fetching currency rates: %v", err)
	}

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	currencyWorker.Stop()

	return nil
}
