package cache_repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/dto"
)

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
	req *dto.CurrencyRateRequestDTO,
) ([]dto.CurrencyRateResponseDTO, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []dto.CurrencyRateResponseDTO

	dateFromUTC := normalizeDateToUTC(req.DateFrom)
	dateToUTC := normalizeDateToUTC(req.DateTo)

	// go by dates from DateFrom to DateTo
	for d := dateFromUTC; !d.After(dateToUTC); d = d.Add(24 * time.Hour) {

		// An additional check just in case: we make sure that d is always UTC.
		dInUTC := d.In(time.UTC)

		key := fmt.Sprintf("%s:%s", dInUTC.Format("2006-01-02"), req.BaseCurrency)

		currency, ok := s.data[key]
		if !ok {
			continue // there is no data on this date
		}

		rate, ok := currency.Rates[req.TargetCurrency]
		if !ok {
			continue // the required currency pair is not available on this date
		}

		result = append(result, dto.CurrencyRateResponseDTO{
			Date: dInUTC,
			Rate: float32(rate),
		})
	}

	return result, nil
}

func normalizeDateToUTC(t time.Time) time.Time {
	y, m, d := t.In(time.UTC).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
