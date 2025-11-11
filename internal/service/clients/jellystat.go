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

// JellystatSystemInfo represents complete system information from Jellystat.
type JellystatSystemInfo struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// JellystatStatistics represents general statistics from Jellystat.
type JellystatStatistics struct {
	Days     int `json:"days"`
	Movies   int `json:"movies"`
	Episodes int `json:"episodes"`
	Songs    int `json:"songs"`
	Total    int `json:"total"`
}

// ViewsByLibraryType represents views aggregated by library type.
type ViewsByLibraryType struct {
	Music   int `json:"music"`
	Movie   int `json:"movie"`
	Episode int `json:"episode"`
	Book    int `json:"book"`
}

// UserActivity represents user activity statistics.
type UserActivity struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	TotalPlays   int    `json:"total_plays"`
	TotalMinutes int    `json:"total_minutes"`
}

// JellystatLibraryStats represents statistics for a specific library in Jellystat.
type JellystatLibraryStats struct {
	LibraryID    string `json:"library_id"`
	LibraryName  string `json:"library_name"`
	TotalItems   int    `json:"total_items"`
	TotalPlays   int    `json:"total_plays"`
	TotalMinutes int    `json:"total_minutes"`
}

// ItemPlaybackStats represents playback statistics for a specific item.
type ItemPlaybackStats struct {
	ItemID            string    `json:"item_id"`
	ItemName          string    `json:"item_name"`
	ItemType          string    `json:"item_type"`
	LibraryName       string    `json:"library_name"`
	PlayCount         int       `json:"play_count"`
	UniqueUsers       int       `json:"unique_users"`
	TotalMinutes      int       `json:"total_minutes"`
	AverageCompletion float64   `json:"average_completion"`
	DateAdded         time.Time `json:"date_added,omitempty"`
	DateFirstPlayed   time.Time `json:"date_first_played,omitempty"`
}

// PlaybackHistory represents a single playback event in the history.
type PlaybackHistory struct {
	PlaybackID   string    `json:"playback_id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	ItemID       string    `json:"item_id"`
	ItemName     string    `json:"item_name"`
	ItemType     string    `json:"item_type"`
	PlayedAt     time.Time `json:"played_at"`
	PlayDuration int       `json:"play_duration"` // in seconds
	PlayMethod   string    `json:"play_method"`
	Device       string    `json:"device"`
}

// ActiveSession represents an active playback session.
type ActiveSession struct {
	SessionID      string `json:"session_id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	Device         string `json:"device"`
	NowPlayingItem *struct {
		ItemID   string `json:"item_id"`
		ItemName string `json:"item_name"`
		ItemType string `json:"item_type"`
	} `json:"now_playing_item"`
	PlaybackPosition int64  `json:"playback_position"` // in ticks
	PlayMethod       string `json:"play_method"`
	IsPaused         bool   `json:"is_paused"`
}

// HistoryParams represents parameters for querying playback history.
type HistoryParams struct {
	UserID string
	ItemID string
	Days   int
	Limit  int
	Offset int
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

// GetStatistics retrieves general statistics from Jellystat.
func (c *JellystatClient) GetStatistics(ctx context.Context, days int) (*JellystatStatistics, error) {
	var stats struct {
		Days     int `json:"days"`
		Movies   int `json:"movies"`
		Episodes int `json:"episodes"`
		Songs    int `json:"songs"`
	}

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&stats).
			SetQueryParam("days", fmt.Sprintf("%d", days)).
			Get("/api/statistics")

		if err != nil {
			return fmt.Errorf("failed to get statistics: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &JellystatStatistics{
		Days:     stats.Days,
		Movies:   stats.Movies,
		Episodes: stats.Episodes,
		Songs:    stats.Songs,
		Total:    stats.Movies + stats.Episodes + stats.Songs,
	}

	c.logger.Info("Retrieved Jellystat statistics",
		zap.Int("days", days),
		zap.Int("movies", result.Movies),
		zap.Int("episodes", result.Episodes),
		zap.Int("songs", result.Songs),
		zap.Int("total", result.Total),
	)

	return result, nil
}

// GetViewsByLibraryType retrieves views aggregated by library type.
func (c *JellystatClient) GetViewsByLibraryType(ctx context.Context, days int) (*ViewsByLibraryType, error) {
	var views ViewsByLibraryType

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&views).
			SetQueryParam("days", fmt.Sprintf("%d", days)).
			Get("/api/stats/getViewsByLibraryType")

		if err != nil {
			return fmt.Errorf("failed to get views by library type: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat views by library type",
		zap.Int("days", days),
		zap.Int("movies", views.Movie),
		zap.Int("episodes", views.Episode),
		zap.Int("music", views.Music),
		zap.Int("books", views.Book),
	)

	return &views, nil
}

// GetUserActivity retrieves user activity statistics.
func (c *JellystatClient) GetUserActivity(ctx context.Context, days int) ([]UserActivity, error) {
	var activities []UserActivity

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&activities).
			SetQueryParam("days", fmt.Sprintf("%d", days)).
			Get("/api/stats/getUserActivity")

		if err != nil {
			return fmt.Errorf("failed to get user activity: %w", err)
		}

		if resp.StatusCode() != 200 {
			// If endpoint doesn't exist, return empty array
			if resp.StatusCode() == 404 {
				c.logger.Warn("User activity endpoint not available")
				return nil
			}
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat user activity",
		zap.Int("days", days),
		zap.Int("user_count", len(activities)),
	)

	return activities, nil
}

// GetLibraryStats retrieves statistics for all libraries.
func (c *JellystatClient) GetLibraryStats(ctx context.Context, days int) ([]JellystatLibraryStats, error) {
	var stats []JellystatLibraryStats

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&stats).
			SetQueryParam("days", fmt.Sprintf("%d", days)).
			Get("/api/stats/getLibraryStats")

		if err != nil {
			return fmt.Errorf("failed to get library stats: %w", err)
		}

		if resp.StatusCode() != 200 {
			// If endpoint doesn't exist, return empty array
			if resp.StatusCode() == 404 {
				c.logger.Warn("Library stats endpoint not available")
				return nil
			}
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat library stats",
		zap.Int("days", days),
		zap.Int("library_count", len(stats)),
	)

	return stats, nil
}

// GetActiveSessions retrieves currently active playback sessions.
func (c *JellystatClient) GetActiveSessions(ctx context.Context) ([]ActiveSession, error) {
	var sessions []ActiveSession

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&sessions).
			Get("/api/sessions")

		if err != nil {
			return fmt.Errorf("failed to get active sessions: %w", err)
		}

		if resp.StatusCode() == 404 {
			c.logger.Warn("Sessions endpoint not available")
			return nil
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat active sessions",
		zap.Int("session_count", len(sessions)),
	)

	return sessions, nil
}

// GetPlaybackHistory retrieves playback history with optional filters.
func (c *JellystatClient) GetPlaybackHistory(ctx context.Context, params HistoryParams) ([]PlaybackHistory, error) {
	var history []PlaybackHistory

	err := c.callWithRetry(ctx, func() error {
		req := c.client.R().
			SetContext(ctx).
			SetResult(&history)

		// Add query parameters
		if params.UserID != "" {
			req.SetQueryParam("user_id", params.UserID)
		}
		if params.ItemID != "" {
			req.SetQueryParam("item_id", params.ItemID)
		}
		if params.Days > 0 {
			req.SetQueryParam("days", fmt.Sprintf("%d", params.Days))
		}
		if params.Limit > 0 {
			req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Offset > 0 {
			req.SetQueryParam("offset", fmt.Sprintf("%d", params.Offset))
		}

		resp, err := req.Get("/api/history")

		if err != nil {
			return fmt.Errorf("failed to get playback history: %w", err)
		}

		if resp.StatusCode() == 404 {
			c.logger.Warn("History endpoint not available")
			return nil
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat playback history",
		zap.Int("record_count", len(history)),
		zap.String("user_id", params.UserID),
		zap.String("item_id", params.ItemID),
	)

	return history, nil
}

// GetMostWatchedItems retrieves the most watched items.
func (c *JellystatClient) GetMostWatchedItems(ctx context.Context, days, limit int, itemType string) ([]ItemPlaybackStats, error) {
	var items []ItemPlaybackStats

	err := c.callWithRetry(ctx, func() error {
		req := c.client.R().
			SetContext(ctx).
			SetResult(&items).
			SetQueryParam("days", fmt.Sprintf("%d", days))

		if limit > 0 {
			req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
		}
		if itemType != "" {
			req.SetQueryParam("type", itemType)
		}

		resp, err := req.Get("/api/items/most-watched")

		if err != nil {
			return fmt.Errorf("failed to get most watched items: %w", err)
		}

		if resp.StatusCode() == 404 {
			c.logger.Warn("Most watched items endpoint not available")
			return nil
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat most watched items",
		zap.Int("days", days),
		zap.Int("item_count", len(items)),
		zap.String("type", itemType),
	)

	return items, nil
}

// GetRecentlyAddedItems retrieves recently added items to Jellyfin.
func (c *JellystatClient) GetRecentlyAddedItems(ctx context.Context, limit int, itemType string) ([]ItemPlaybackStats, error) {
	var items []ItemPlaybackStats

	err := c.callWithRetry(ctx, func() error {
		req := c.client.R().
			SetContext(ctx).
			SetResult(&items)

		if limit > 0 {
			req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
		}
		if itemType != "" {
			req.SetQueryParam("type", itemType)
		}

		resp, err := req.Get("/api/items/recently-added")

		if err != nil {
			return fmt.Errorf("failed to get recently added items: %w", err)
		}

		if resp.StatusCode() == 404 {
			c.logger.Warn("Recently added items endpoint not available")
			return nil
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	c.logger.Info("Retrieved Jellystat recently added items",
		zap.Int("item_count", len(items)),
		zap.String("type", itemType),
	)

	return items, nil
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
