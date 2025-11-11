package handler

import (
	"context"

	"github.com/Sevn9/currency-screener/currency/internal/dto"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"go.uber.org/zap"
)

// setup Services methods for CurrencyServer
type CurrencyService interface {
	GetCurrencyRateFromInterval(ctx context.Context, reqDto *dto.CurrencyRateRequestDTO) ([]dto.CurrencyRateResponseDTO, error)
	FetchAndSaveCurrencyRate(ctx context.Context, baseCurrency string) error
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
