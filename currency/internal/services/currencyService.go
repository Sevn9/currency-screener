package services

import "github.com/Sevn9/currency-screener/currency/internal/clients/currency"

type CurrencyService struct {
	currencyClient currency.CurrencyClient
}
