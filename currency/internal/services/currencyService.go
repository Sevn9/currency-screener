package services

import (
	"github.com/Sevn9/currency-screener/currency/internal/clients/currency"
	"go.uber.org/zap"
)

// сurrencyService interface realization
type currencyService struct {
	currencyClient currency.CurrencyClient
	logger         *zap.Logger
}

func NewCurrencyService(client currency.CurrencyClient, logger *zap.Logger) *currencyService {
	return &currencyService{
		currencyClient: client,
		logger:         logger,
	}
}
