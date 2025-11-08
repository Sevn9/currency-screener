package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/cache_repository"
	"github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"go.uber.org/zap"
)

// сurrencyService interface realization
type currencyService struct {
	cacheRepo      *cache_repository.MemoryRateStorage
	currencyClient currency.CurrencyClient
	logger         *zap.Logger
}

func NewCurrencyService(cacheRepository *cache_repository.MemoryRateStorage, client currency.CurrencyClient, logger *zap.Logger) *currencyService {
	return &currencyService{
		cacheRepo:      cacheRepository,
		currencyClient: client,
		logger:         logger,
	}
}

func (c *currencyService) FetchAndSaveCurrencyRate(ctx context.Context, baseCurrency string) error {
	//get currency rate
	currClientTemp, er := c.currencyClient.GetCurrentRate(ctx)
	if er != nil {
		c.logger.Error("services: error GetCurrentRate create:",
			zap.Error(er))
		return er
	}

	date, err := time.Parse("2006-01-02", currClientTemp.Date)

	if err != nil {
		return fmt.Errorf("failed to parse currency date: %v ", err)
	}

	//save currency rate to cache
	c.cacheRepo.Save(ctx, date, baseCurrency, currClientTemp.Rub)

	//todo: save currency rate to db

	return nil
}

func (c *currencyService) GetCurrencyRateFromInterval(ctx context.Context) {

	//todo: get from cache

	return
}
