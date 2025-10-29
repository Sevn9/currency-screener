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
)

type CurrencyClient struct {
	baseUrl    *url.URL
	httpClient *http.Client
}

type ResponceCurrencyRate struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

func NewCurrencyClient(cfg config.PublicCurrencyAPIConfig) (CurrencyClient, error) {
	baseUrl, err := url.Parse(cfg.ApiUrl)

	if err != nil {
		return CurrencyClient{}, fmt.Errorf("invalid base URL: %w", err)
	}

	httpClient := http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}

	return CurrencyClient{
		baseUrl:    baseUrl,
		httpClient: &httpClient,
	}, nil
}

func (c *CurrencyClient) GetCurrentRate(ctx context.Context) {
	currency := "rub"
	//currencies endpoint url
	url, _ := url.Parse(fmt.Sprintf("/v1/currencies/%s.json", strings.ToLower(currency)))
	//add base path to url
	fullUrl := c.baseUrl.ResolveReference(url).String()

	fmt.Println("fullUrl: ", fullUrl)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)

	if err != nil {
		fmt.Println("Error create request: ", err)
	}

	resp, err := c.httpClient.Do(request)

	if err != nil {
		fmt.Println("Error responce: ", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("received non-200 response code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("reading responce body error: %v\n", err)
		return
	}
	fmt.Println("Answer:", string(bodyBytes))

	var rateResponse ResponceCurrencyRate
	if err := json.Unmarshal(bodyBytes, &rateResponse); err != nil {
		fmt.Println("failed to unmarshal response: ", err)
	}

}
