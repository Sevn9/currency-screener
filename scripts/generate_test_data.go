package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "currency1"
	password = "secret123"
	dbname   = "currency_db"
)

func main() {
	// connection for pgx
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the DB using pgxpool: %v", err)
	}
	defer dbPool.Close()

	// test connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := generateTestData(ctx, dbPool); err != nil {
		log.Fatalf("Failed to generate test data: %v", err)
	}

	fmt.Println("Test data generated successfully using pgxpool.")
}

// for *pgxpool.Pool
func generateTestData(ctx context.Context, dbPool *pgxpool.Pool) error {
	baseCurrency := "RUB"
	targetCurrencies := []string{"usd", "eur", "gbp", "jpy", "cny"}

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// dbPool.BeginTx for transaction
	tx, err := dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	insertStmt := `INSERT INTO exchange_rates (date, base_currency, currency_rates) VALUES ($1, $2, $3)`

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		ratesMap := make(map[string]float64)
		for _, targetCurrency := range targetCurrencies {
			ratesMap[targetCurrency] = generateRandomRate()
		}

		ratesJSON, err := json.Marshal(ratesMap)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		_, err = tx.Exec(ctx, insertStmt, date, baseCurrency, ratesJSON)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to insert test data: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func generateRandomRate() float64 {
	return 0.01 + (0.1-0.01)*rand.Float64()
}
