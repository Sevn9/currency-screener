package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/Sevn9/currency-screener/gateway/internal/repository"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrencyService struct {
	currencyClient currency.CurrencyServiceClient
	cache          *repository.CurrencyRedisRepository
}

func NewCurrency(
	currencyClient currency.CurrencyServiceClient,
	cache *repository.CurrencyRedisRepository) CurrencyService {
	return CurrencyService{
		currencyClient: currencyClient,
		cache:          cache,
	}
}

func (s *CurrencyService) GetCurrencyRates(
	ctx context.Context,
	request dto.ParsedCurrencyRequest) (*dto.CurrencyResponse, error) {

	// reading from cache
	cachedResp, err := s.cache.Get(ctx, request)
	if err != nil {
		log.Printf("Error reading from cache: %v", err)
	}
	if cachedResp != nil {
		return cachedResp, nil
	}

	pbResp, err := s.currencyClient.GetRates(
		ctx, &currency.GetRateRequest{
			Currency: request.Currency,
			DateFrom: timestamppb.New(request.DateFrom),
			DateTo:   timestamppb.New(request.DateTo),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("currencyClient.GetRate: %s", err)
	}

	resp := &dto.CurrencyResponse{
		Currency: pbResp.GetCurrency(),
		Rates:    make([]dto.CurrencyRate, 0, len(pbResp.Rates)),
	}

	for _, rate := range pbResp.Rates {
		resp.Rates = append(
			resp.Rates, dto.CurrencyRate{
				Rate: rate.Rate,
				Date: rate.Date.AsTime(),
			},
		)
	}

	// save to cache
	if err := s.cache.Set(ctx, request, resp); err != nil {
		log.Printf("Error saving to cache: %v", err)
	}

	return resp, nil
}
