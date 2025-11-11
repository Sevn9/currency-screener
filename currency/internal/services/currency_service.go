package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/cache_repository"
	"github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"github.com/Sevn9/currency-screener/currency/internal/dto"
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
		return fmt.Errorf("services: error GetCurrentRate create: %v ", er)
	}

	date, err := time.Parse("2006-01-02", currClientTemp.Date)

	if err != nil {
		return fmt.Errorf("services: failed to parse currency date: %v ", err)
	}

	//save currency rate to cache
	c.cacheRepo.Save(ctx, date, baseCurrency, currClientTemp.Rub)
	c.logger.Info("currency save to cache", zap.Any("rates", currClientTemp.Rub))

	//todo: save currency rate to db

	return nil
}

func (c *currencyService) GetCurrencyRateFromInterval(ctx context.Context, reqDto *dto.CurrencyRateRequestDTO) ([]dto.CurrencyRateResponseDTO, error) {

	//get from cache
	result, err := c.cacheRepo.GetCurrencyRatesInInterval(ctx, reqDto)

	if err != nil {
		return nil, fmt.Errorf("failed to parse currency date: %v ", err)
	}

	return result, nil
}
