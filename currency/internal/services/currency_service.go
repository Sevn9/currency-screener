package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/cache_repository"
	"github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"github.com/Sevn9/currency-screener/currency/internal/dto"
	"github.com/Sevn9/currency-screener/currency/internal/repository"
	"go.uber.org/zap"
)

// сurrencyService interface realization
type CurrencyService struct {
	cacheRepo      *cache_repository.MemoryRateStorage
	dbRepo         *repository.CurrencyPostgres
	currencyClient currency.CurrencyClient
	logger         *zap.Logger
}

func NewCurrencyService(cacheRepository *cache_repository.MemoryRateStorage, dbRepository *repository.CurrencyPostgres, client currency.CurrencyClient, logger *zap.Logger) *CurrencyService {
	return &CurrencyService{
		cacheRepo:      cacheRepository,
		dbRepo:         dbRepository,
		currencyClient: client,
		logger:         logger,
	}
}

func (c *CurrencyService) FetchAndSaveCurrencyRate(ctx context.Context, baseCurrency string) error {
	//get currency rate
	currClientTemp, err := c.currencyClient.GetCurrentRate(ctx)
	if err != nil {
		c.logger.Error("services: error GetCurrentRate create:",
			zap.Error(err))
		return fmt.Errorf("services: error GetCurrentRate create: %v ", err)
	}

	date, err := time.Parse("2006-01-02", currClientTemp.Date)

	if err != nil {
		return fmt.Errorf("services: failed to parse currency date: %v ", err)
	}

	//save currency rate to cache ram
	//c.cacheRepo.Save(ctx, date, baseCurrency, currClientTemp.Rub)
	c.logger.Info("currency save to cache", zap.Any("rates", currClientTemp.Rub))

	//save currency rate to db
	err = c.dbRepo.Save(ctx, date, baseCurrency, currClientTemp.Rub)

	if err != nil {
		c.logger.Error("error currency save to db", zap.Error(err))
		return fmt.Errorf("failed to save currency rates to database: %w", err)
	} else {
		c.logger.Info("currency save to db", zap.Any("rates", currClientTemp.Rub))
	}
	return nil
}

func (c *CurrencyService) GetCurrencyRateFromInterval(ctx context.Context, reqDto *dto.CurrencyRateRequestDTO) ([]dto.CurrencyRateResponseDTO, error) {

	//get from cache
	//result, err := c.cacheRepo.GetCurrencyRatesInInterval(ctx, reqDto)

	result, err := c.dbRepo.GetCurrencyRatesInInterval(ctx, reqDto)

	if err != nil {
		c.logger.Error("services: error get currency", zap.Error(err))
		return nil, fmt.Errorf("services: error get currency: %v ", err)
	}

	return result, nil
}
