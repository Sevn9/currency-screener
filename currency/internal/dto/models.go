package dto

import (
	"time"

	"github.com/Sevn9/currency-screener/pkg/currency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CurrencyRateResponseDTO struct {
	Date time.Time
	Rate float32
}

type CurrencyRateRequestDTO struct {
	BaseCurrency   string
	TargetCurrency string
	DateFrom       time.Time
	DateTo         time.Time
}

type CurrencyResponseDTO struct {
	Currency string
	Rates    []CurrencyRateResponseDTO
}

const (
	DefaultBaseCurrency = "RUB"
)

func CurrencyRequestDTOFromProtobuf(req *currency.GetRateRequest, baseCurrency string) *CurrencyRateRequestDTO {
	return &CurrencyRateRequestDTO{
		BaseCurrency:   baseCurrency,
		TargetCurrency: req.Currency,
		DateFrom:       req.DateFrom.AsTime(),
		DateTo:         req.DateTo.AsTime(),
	}
}

func (dto *CurrencyResponseDTO) ToProtobuf() *currency.GetRateResponse {
	rateRecords := make([]*currency.RateRecord, 0, len(dto.Rates))
	for _, record := range dto.Rates {
		rateRecords = append(
			rateRecords, &currency.RateRecord{
				Date: timestamppb.New(record.Date),
				Rate: record.Rate,
			},
		)
	}

	return &currency.GetRateResponse{
		Currency: dto.Currency,
		Rates:    rateRecords,
	}
}
