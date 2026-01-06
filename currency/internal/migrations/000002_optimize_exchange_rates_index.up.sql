DROP INDEX IF EXISTS idx_exchange_rates_date_base_currency;

CREATE INDEX idx_exchange_rates_base_currency_date ON exchange_rates(base_currency, date);