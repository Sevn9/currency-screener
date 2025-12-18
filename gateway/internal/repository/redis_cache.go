package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sevn9/currency-screener/gateway/internal/dto"
	"github.com/go-redis/redis/v8"
)

type CurrencyRedisRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCurrencyRedisRepository(client *redis.Client, ttl time.Duration) *CurrencyRedisRepository {
	return &CurrencyRedisRepository{
		client: client,
		ttl:    ttl,
	}
}

// generateKey
// example: "rates:RUB:2024-01-01:2024-01-31"
func (r *CurrencyRedisRepository) generateKey(req dto.ParsedCurrencyRequest) string {
	return fmt.Sprintf("rates:%s:%s:%s",
		req.Currency,
		req.DateFrom.Format("2006-01-02"),
		req.DateTo.Format("2006-01-02"),
	)
}

func (r *CurrencyRedisRepository) Get(ctx context.Context, req dto.ParsedCurrencyRequest) (*dto.CurrencyResponse, error) {
	key := r.generateKey(req)

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // cache empty
		}
		return nil, err
	}

	var response dto.CurrencyResponse
	if err := json.Unmarshal(val, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	return &response, nil
}

func (r *CurrencyRedisRepository) Set(ctx context.Context, req dto.ParsedCurrencyRequest, resp *dto.CurrencyResponse) error {
	key := r.generateKey(req)

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	return r.client.Set(ctx, key, data, r.ttl).Err()
}
