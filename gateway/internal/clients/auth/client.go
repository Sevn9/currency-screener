package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Sevn9/currency-screener/gateway/internal/config"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenGeneration      = errors.New("token generation failed")

	ErrTokenNotFound         = errors.New("token not found in header")
	ErrInvalidOrExpiredToken = errors.New("invalid signature or token expired")
)

const (
	pingPath     = "/ping"
	generatePath = "/generate"
	validatePath = "/validate"

	authorizationHeader = "Authorization"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewAuthClient(cfg config.AuthApiServiceConfig) (*Client, error) { // todo pass config
	parsedURL, err := url.Parse(cfg.AuthUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	return &Client{
		baseURL: parsedURL,
		httpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Duration(cfg.TimeoutSeconds),
		},
	}, nil
}

func (c *Client) Ping() (string, error) {
	relativePingPath, _ := url.Parse(pingPath)
	fullURL := *c.baseURL.ResolveReference(relativePingPath)

	resp, err := c.httpClient.Get(fullURL.String())
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(body), nil
}

func (c *Client) GenerateToken(ctx context.Context, login string) (string, error) {
	relativeGeneratePath, _ := url.Parse(generatePath)
	fullURL := *c.baseURL.ResolveReference(relativeGeneratePath)

	query := fullURL.Query()
	query.Set("login", login)
	fullURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("httpClient.Do: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}
		return string(bodyBytes), nil
	case http.StatusBadRequest:
		return "", fmt.Errorf("%w: bad request", ErrTokenGeneration)
	case http.StatusUnauthorized:
		return "", fmt.Errorf("%w: unauthorized", ErrInvalidCredentials)
	default:
		return "", fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) ValidateToken(ctx context.Context, token string) error {
	relativeValidatePath, _ := url.Parse(validatePath)
	fullURL := *c.baseURL.ResolveReference(relativeValidatePath)

	//create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	//add req token
	req.Header.Set(authorizationHeader, "Bearer "+token)

	//do req
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	errorMessage := string(body)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrTokenNotFound, errorMessage)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrInvalidOrExpiredToken, errorMessage)
	default:
		return fmt.Errorf("%w %d: %s", ErrUnexpectedStatusCode, resp.StatusCode, errorMessage)
	}
}
