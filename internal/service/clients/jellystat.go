package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// JellystatClient implements a client for Jellystat statistics service.
type JellystatClient struct {
	client  *resty.Client
	baseURL string
	apiKey  string
	logger  *zap.Logger
}

// NewJellystatClient creates a new Jellystat client.
func NewJellystatClient(config ClientConfig, logger *zap.Logger) *JellystatClient {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	client.SetHeader("x-api-token", config.APIKey)
	client.SetTimeout(config.Timeout)

	if config.Timeout == 0 {
		client.SetTimeout(DefaultTimeout)
	}

	return &JellystatClient{
		client:  client,
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		logger:  logger,
	}
}

// jellystatSystemInfo represents the system info response from Jellystat API.
// nolint:unused // Reserved for future use
type jellystatSystemInfo struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// JellystatSystemInfo represents complete system information from Jellystat.
type JellystatSystemInfo struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// TestConnection verifies the connection to Jellystat.
func (c *JellystatClient) TestConnection(ctx context.Context) error {
	return c.callWithRetry(ctx, func() error {
		// Try the getconfig endpoint which is publicly available
		resp, err := c.client.R().
			SetContext(ctx).
			Get("/api/getconfig")

		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), string(resp.Body()))
		}

		c.logger.Info("Jellystat connection successful",
			zap.String("url", c.baseURL),
		)

		return nil
	})
}

// GetSystemInfo retrieves complete system information from Jellystat.
func (c *JellystatClient) GetSystemInfo(ctx context.Context) (*JellystatSystemInfo, error) {
	var info struct {
		JFHOST       string `json:"JF_HOST"`
		APPUSER      string `json:"APP_USER"`
		RequireLogin bool   `json:"REQUIRE_LOGIN"`
		IsJellyfin   bool   `json:"IS_JELLYFIN"`
	}

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&info).
			Get("/api/getconfig")

		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to our model
	systemInfo := &JellystatSystemInfo{
		Version: "unknown", // Jellystat doesn't expose version in getconfig
		Status:  "connected",
	}

	c.logger.Info("Retrieved Jellystat system info",
		zap.String("host", info.JFHOST),
		zap.String("status", systemInfo.Status),
	)

	return systemInfo, nil
}

// callWithRetry executes a function with retry logic.
func (c *JellystatClient) callWithRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	backoff := RetryDelay

	for i := 0; i < MaxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if i < MaxRetries-1 {
			c.logger.Warn("Jellystat API call failed, retrying",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", MaxRetries),
				zap.Error(err),
				zap.Duration("backoff", backoff),
			)

			select {
			case <-time.After(backoff):
				backoff *= 2 // Exponential backoff
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}
