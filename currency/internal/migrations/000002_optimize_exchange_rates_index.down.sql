DROP INDEX IF EXISTS idx_exchange_rates_base_currency_date;

CREATE INDEX idx_exchange_rates_date_base_currency ON exchange_rates(date, base_currency);