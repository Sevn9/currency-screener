package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Currency struct {
	Date         time.Time
	BaseCurrency string
	Rates        map[string]float64
}

type CurrencyPostgres struct {
	DB *pgxpool.Pool
}

func NewCurrencyPostgres(dbPool *pgxpool.Pool) (*CurrencyPostgres, error) {
	if dbPool == nil {
		return nil, fmt.Errorf("database pool cannot be nil")
	}

	return &CurrencyPostgres{
		DB: dbPool,
	}, nil
}

// Close db connection pull
func (r *CurrencyPostgres) Close() {
	r.DB.Close()
}

func (s *CurrencyPostgres) Save(
	ctx context.Context,
	date time.Time,
	baseCurrency string,
	rates map[string]float64,
) error {
	ratesJSON, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("failed to marshal currency rates: %w", err)
	}

	_, err = s.DB.Exec(
		ctx,
		`INSERT INTO exchange_rates (date, base_currency, currency_rates)
		 VALUES ($1, $2, $3)`,
		date, baseCurrency, ratesJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to save exchange rates: %w", err)
	}

	return nil
}

func (s *CurrencyPostgres) GetCurrencyRatesInInterval(
	ctx context.Context,
	req *dto.CurrencyRateRequestDTO,
) ([]dto.CurrencyRateResponseDTO, error) {

	query := `
		SELECT date, (currency_rates ->> $1)::float
		FROM exchange_rates
		WHERE date BETWEEN $2 AND $3
		  AND base_currency = $4
		ORDER BY date ASC;
	`

	rows, err := s.DB.Query(
		ctx,
		query,
		req.TargetCurrency,
		req.DateFrom.Format("2006-01-02"),
		req.DateTo.Format("2006-01-02"),
		req.BaseCurrency,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query exchange rates: %w", err)
	}
	defer rows.Close()

	var result []dto.CurrencyRateResponseDTO
	for rows.Next() {
		var date time.Time
		var rate float32
		if err := rows.Scan(&date, &rate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, dto.CurrencyRateResponseDTO{
			Date: date.UTC(),
			Rate: rate,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return result, nil
}
