package handler

import (
	"context"
	"fmt"

	"github.com/Sevn9/currency-screener/currency/internal/dto"
	"github.com/Sevn9/currency-screener/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpc proto server function realization
func (c *CurrencyServer) GetRates(ctx context.Context, request *currency.GetRateRequest) (*currency.GetRateResponse, error) {

	reqDTO := dto.CurrencyRequestDTOFromProtobuf(request, dto.DefaultBaseCurrency)

	rates, err := c.service.GetCurrencyRateFromInterval(ctx, reqDTO)
	if err != nil {
		return nil, fmt.Errorf("service.GetCurrencyRatesInInterval: %w", err)
	}

	rateRecords := make([]*currency.RateRecord, len(rates))
	for i, rate := range rates {
		rateRecords[i] = &currency.RateRecord{
			Date: timestamppb.New(rate.Date),
			Rate: rate.Rate,
		}
	}

	return &currency.GetRateResponse{
		Currency: reqDTO.TargetCurrency,
		Rates:    rateRecords,
	}, nil
}
