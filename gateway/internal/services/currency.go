package services

import (
	"context"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/Sevn9/currency-screener/pkg/currency"
)

type CurrencyService struct {
	currencyClient currency.CurrencyServiceClient
}

func NewCurrency(currencyClient currency.CurrencyServiceClient) CurrencyService {
	return CurrencyService{
		currencyClient: currencyClient,
	}
}

func (svc *CurrencyService) GetCurrencyRates(
	ctx context.Context,
	request dto.ParsedCurrencyRequest) (*dto.CurrencyResponse, error) {
	return nil, nil
}
