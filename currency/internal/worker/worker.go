package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/config"
	"github.com/Sevn9/currency-screener/currency/internal/services"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Currency struct {
	cron            *cron.Cron
	currencyService *services.CurrencyService
	schedule        string
	baseCurrency    string
	targetCurrency  string
	logger          *zap.Logger
	timeout         int
}

func NewCurrency(
	cfg config.WorkerConfig,
	currencyService *services.CurrencyService,
	cron *cron.Cron,
	logger *zap.Logger,
) *Currency {
	return &Currency{
		cron:            cron,
		currencyService: currencyService,
		schedule:        cfg.Schedule,
		baseCurrency:    cfg.CurrencyPair.BaseCurrency,
		targetCurrency:  cfg.CurrencyPair.TargetCurrency,
		logger:          logger,
		timeout:         cfg.Timeout,
	}
}

func (w *Currency) StartFetchingCurrencyRates() error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(w.timeout))
		defer cancel()

		err := w.currencyService.FetchAndSaveCurrencyRate(ctx, w.baseCurrency)

		if err != nil {
			w.logger.Error(
				"worker: Failed to fetch currency rate immediately on startup",
				zap.Time("timestamp", time.Now()),
				zap.Error(err),
			)
		}
	}()

	_, err := w.cron.AddFunc(
		w.schedule, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(w.timeout))
			defer cancel()

			err := w.currencyService.FetchAndSaveCurrencyRate(ctx, w.baseCurrency)
			if err != nil {
				w.logger.Error(
					"worker: Failed to fetch currency rate on scheduled run",
					zap.Time("timestamp", time.Now()),
					zap.Error(err),
					zap.String("schedule", w.schedule),
				)
			}
		},
	)

	if err != nil {
		return fmt.Errorf("worker: Cron.AddFunc: %w", err)
	}

	w.cron.Start()

	return nil
}

func (w *Currency) Stop() {
	w.cron.Stop()
}
