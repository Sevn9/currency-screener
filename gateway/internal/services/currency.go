package services

import (
	"context"
	"fmt"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrencyService struct {
	currencyClient currency.CurrencyServiceClient
}

func NewCurrency(currencyClient currency.CurrencyServiceClient) CurrencyService {
	return CurrencyService{
		currencyClient: currencyClient,
	}
}

func (s *CurrencyService) GetCurrencyRates(
	ctx context.Context,
	request dto.ParsedCurrencyRequest) (*dto.CurrencyResponse, error) {

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
	return resp, nil
}
