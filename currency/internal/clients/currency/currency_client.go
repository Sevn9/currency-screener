package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sevn9/currency-screener/currency/internal/config"
	"go.uber.org/zap"
)

type CurrencyClient struct {
	baseUrl    *url.URL
	httpClient *http.Client
	logger     *zap.Logger
}

type ResponceCurrencyRate struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

func NewCurrencyClient(cfg config.PublicCurrencyAPIConfig, logger *zap.Logger) (CurrencyClient, error) {
	baseUrl, err := url.Parse(cfg.ApiUrl)

	if err != nil {
		logger.Error("currency: invalid base URL",
			zap.Error(err))
		return CurrencyClient{}, fmt.Errorf("invalid base URL: %w", err)
	}

	httpClient := http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}

	return CurrencyClient{
		baseUrl:    baseUrl,
		httpClient: &httpClient,
		logger:     logger,
	}, nil
}

func (c *CurrencyClient) GetCurrentRate(ctx context.Context) (ResponceCurrencyRate, error) {
	ctxTimeout, close := context.WithTimeout(ctx, 10*time.Second)
	defer close()

	currency := "rub"
	//currencies endpoint url
	url, _ := url.Parse(fmt.Sprintf("/v1/currencies/%s.json", strings.ToLower(currency)))
	//add base path to url
	fullUrl := c.baseUrl.ResolveReference(url).String()

	fmt.Println("currency: fullUrl: ", fullUrl)

	request, err := http.NewRequestWithContext(ctxTimeout, http.MethodGet, fullUrl, nil)

	if err != nil {
		c.logger.Error("currency: Error create request",
			zap.Error(err))
		return ResponceCurrencyRate{}, fmt.Errorf("error create request: %w", err)
	}

	resp, err := c.httpClient.Do(request)

	if err != nil {
		c.logger.Error("currency: Error responce",
			zap.Error(err))
		return ResponceCurrencyRate{}, fmt.Errorf("error responce: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("received non-200 response code: %d", resp.StatusCode)
		return ResponceCurrencyRate{}, fmt.Errorf("received non-200 response code:Error responce: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("currency: reading responce body error",
			zap.Error(err))
		return ResponceCurrencyRate{}, fmt.Errorf("reading responce body error: %w", err)
	}
	fmt.Println("Answer:", string(bodyBytes))

	var rateResponse ResponceCurrencyRate
	if err := json.Unmarshal(bodyBytes, &rateResponse); err != nil {
		c.logger.Error("currency: failed to unmarshal response",
			zap.Error(err))
		return ResponceCurrencyRate{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return rateResponse, nil
}
