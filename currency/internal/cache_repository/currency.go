package cache_repository

import (
	"context"
	"fmt"
	"sync"
	"time"
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

type Currency struct {
	Date         time.Time
	BaseCurrency string
	Rates        map[string]float64
}

type MemoryRateStorage struct {
	mu   sync.RWMutex
	data map[string]Currency // key = "2025-11-08:rub"
}

func NewMemoryRateStorage() *MemoryRateStorage {
	return &MemoryRateStorage{
		data: make(map[string]Currency),
	}
}

func (s *MemoryRateStorage) Save(
	ctx context.Context,
	date time.Time,
	baseCurrency string,
	rates map[string]float64) {
	key := fmt.Sprintf("%s:%s", date.Format("2006-01-02"), baseCurrency)

	s.mu.Lock()
	s.data[key] = Currency{Date: date, BaseCurrency: baseCurrency, Rates: rates}
	s.mu.Unlock()
}

func (s *MemoryRateStorage) GetCurrencyRatesInInterval(
	ctx context.Context,
	req CurrencyRateRequestDTO,
) ([]CurrencyRateResponseDTO, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []CurrencyRateResponseDTO

	// go by dates from DateFrom to DateTo
	for d := req.DateFrom; !d.After(req.DateTo); d = d.Add(24 * time.Hour) {

		key := fmt.Sprintf("%s:%s", d.Format("2006-01-02"), req.BaseCurrency)

		currency, ok := s.data[key]
		if !ok {
			continue // there is no data on this date
		}

		rate, ok := currency.Rates[req.TargetCurrency]
		if !ok {
			continue // the required currency pair is not available on this date
		}

		result = append(result, CurrencyRateResponseDTO{
			Date: d,
			Rate: float32(rate),
		})
	}

	return result, nil
}
