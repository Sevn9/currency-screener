package handler

import (
	"github.com/Sevn9/currency-screener/pkg/currency"
	"go.uber.org/zap"
)

// setup Services for CurrencyServer
type CurrencyService interface {
}

// setup CurrencyServer
type CurrencyServer struct {
	currency.UnimplementedCurrencyServiceServer
	service CurrencyService
	logger  *zap.Logger
}

func NewCurrencyServer(svc CurrencyService, logger *zap.Logger) *CurrencyServer {
	return &CurrencyServer{
		service: svc,
		logger:  logger,
	}
}
