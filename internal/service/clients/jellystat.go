package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// JellystatClient handles interactions with Jellystat API for playback statistics
type JellystatClient struct {
	client  *resty.Client
	baseURL string
	token   string
	logger  zerolog.Logger
}

// PlaybackStatsClient interface for playback statistics operations
type PlaybackStatsClient interface {
	TestConnection(ctx context.Context) error
	GetItemPlaybackStats(ctx context.Context, itemID string, hours int) (*PlaybackStats, error)
	GetRecentlyWatched(ctx context.Context, userID string, limit int) ([]WatchedItem, error)
	GetLibraryStats(ctx context.Context, libraryID string) (*LibraryStats, error)
	GetUserActivity(ctx context.Context, userID string, days int) (*UserActivity, error)
}

// PlaybackStats contains playback statistics for a media item
type PlaybackStats struct {
	Plays                  int64 `json:"plays"`
	TotalPlaybackDuration  int64 `json:"total_playback_duration"`
	LastActivityDate       time.Time
}

// WatchedItem represents a recently watched item
type WatchedItem struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Type                  string    `json:"type"`
	SeriesName            string    `json:"series_name,omitempty"`
	SeasonNumber          int       `json:"season_number,omitempty"`
	EpisodeNumber         int       `json:"episode_number,omitempty"`
	ActivityDateInserted  time.Time `json:"activity_date_inserted"`
	PlaybackDuration      int64     `json:"playback_duration"`
	UserName              string    `json:"user_name"`
}

// LibraryStats contains statistics for a library
type LibraryStats struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Plays                 int64     `json:"plays"`
	TotalPlaybackDuration int64     `json:"total_playback_duration"`
	LastActivity          time.Time `json:"last_activity"`
	ItemCount             int       `json:"item_count"`
}

// UserActivity contains user activity statistics
type UserActivity struct {
	UserID                string    `json:"user_id"`
	UserName              string    `json:"user_name"`
	Plays                 int64     `json:"plays"`
	TotalPlaybackDuration int64     `json:"total_playback_duration"`
	LastActivity          time.Time `json:"last_activity"`
}

// NewJellystatClient creates a new Jellystat client
func NewJellystatClient(baseURL, token string, logger zerolog.Logger) *JellystatClient {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	client.SetTimeout(30 * time.Second)

	return &JellystatClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
		logger:  logger.With().Str("client", "jellystat").Logger(),
	}
}

// TestConnection verifies connection to Jellystat
func (c *JellystatClient) TestConnection(ctx context.Context) error {
	c.logger.Debug().Msg("testing Jellystat connection")

	// Try to get libraries as a connection test
	resp, err := c.client.R().
		SetContext(ctx).
		Get("/api/getLibraries")

	if err != nil {
		c.logger.Error().Err(err).Msg("failed to connect to Jellystat")
		return fmt.Errorf("Jellystat connection failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		c.logger.Error().Int("status_code", resp.StatusCode()).Msg("Jellystat returned error status")
		return fmt.Errorf("Jellystat returned status %d", resp.StatusCode())
	}

	c.logger.Info().Msg("Jellystat connection successful")
	return nil
}

// GetItemPlaybackStats retrieves playback statistics for a specific item
func (c *JellystatClient) GetItemPlaybackStats(ctx context.Context, itemID string, hours int) (*PlaybackStats, error) {
	c.logger.Debug().
		Str("item_id", itemID).
		Int("hours", hours).
		Msg("getting item playback stats")

	var result struct {
		Plays                 int64  `json:"Plays"`
		TotalPlaybackDuration int64  `json:"total_playback_duration"`
		LastActivityDate      string `json:"LastActivityDate"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"hours":  hours,
			"itemid": itemID,
		}).
		SetResult(&result).
		Post("/stats/getGlobalItemStats")

	if err != nil {
		c.logger.Error().Err(err).Str("item_id", itemID).Msg("failed to get item stats")
		return nil, fmt.Errorf("failed to get item stats: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jellystat returned status %d", resp.StatusCode())
	}

	stats := &PlaybackStats{
		Plays:                 result.Plays,
		TotalPlaybackDuration: result.TotalPlaybackDuration,
	}

	if result.LastActivityDate != "" {
		if t, err := time.Parse(time.RFC3339, result.LastActivityDate); err == nil {
			stats.LastActivityDate = t
		}
	}

	c.logger.Info().
		Str("item_id", itemID).
		Int64("plays", stats.Plays).
		Int64("total_duration", stats.TotalPlaybackDuration).
		Msg("retrieved item playback stats")

	return stats, nil
}

// GetRecentlyWatched retrieves recently watched items for a user
func (c *JellystatClient) GetRecentlyWatched(ctx context.Context, userID string, limit int) ([]WatchedItem, error) {
	c.logger.Debug().
		Str("user_id", userID).
		Int("limit", limit).
		Msg("getting recently watched items")

	var result []struct {
		ID                   string `json:"Id"`
		Name                 string `json:"Name"`
		Type                 string `json:"Type"`
		SeriesName           string `json:"SeriesName"`
		SeasonNumber         int    `json:"SeasonNumber"`
		EpisodeNumber        int    `json:"EpisodeNumber"`
		ActivityDateInserted string `json:"ActivityDateInserted"`
		PlaybackDuration     int64  `json:"PlaybackDuration"`
		UserName             string `json:"UserName"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("libraryid", "all").
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetResult(&result).
		Get("/api/getRecentlyAdded")

	if err != nil {
		c.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get recently watched")
		return nil, fmt.Errorf("failed to get recently watched: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jellystat returned status %d", resp.StatusCode())
	}

	watched := make([]WatchedItem, 0, len(result))
	for _, item := range result {
		w := WatchedItem{
			ID:               item.ID,
			Name:             item.Name,
			Type:             item.Type,
			SeriesName:       item.SeriesName,
			SeasonNumber:     item.SeasonNumber,
			EpisodeNumber:    item.EpisodeNumber,
			PlaybackDuration: item.PlaybackDuration,
			UserName:         item.UserName,
		}

		if item.ActivityDateInserted != "" {
			if t, err := time.Parse(time.RFC3339, item.ActivityDateInserted); err == nil {
				w.ActivityDateInserted = t
			}
		}

		watched = append(watched, w)
	}

	c.logger.Info().
		Str("user_id", userID).
		Int("count", len(watched)).
		Msg("retrieved recently watched items")

	return watched, nil
}

// GetLibraryStats retrieves statistics for a library
func (c *JellystatClient) GetLibraryStats(ctx context.Context, libraryID string) (*LibraryStats, error) {
	c.logger.Debug().Str("library_id", libraryID).Msg("getting library stats")

	var result []struct {
		ID                    string `json:"Id"`
		Name                  string `json:"Name"`
		Plays                 int64  `json:"Plays"`
		TotalPlaybackDuration int64  `json:"total_playback_duration"`
		LastActivity          string `json:"LastActivity"`
		LibraryCount          int    `json:"Library_Count"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&result).
		Get("/api/getLibraryOverview")

	if err != nil {
		c.logger.Error().Err(err).Str("library_id", libraryID).Msg("failed to get library stats")
		return nil, fmt.Errorf("failed to get library stats: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jellystat returned status %d", resp.StatusCode())
	}

	// Find the library in the results
	for _, lib := range result {
		if lib.ID == libraryID {
			stats := &LibraryStats{
				ID:                    lib.ID,
				Name:                  lib.Name,
				Plays:                 lib.Plays,
				TotalPlaybackDuration: lib.TotalPlaybackDuration,
				ItemCount:             lib.LibraryCount,
			}

			if lib.LastActivity != "" {
				// LastActivity comes as an interval string from PostgreSQL
				// For now, we'll set it to current time if not parseable
				stats.LastActivity = time.Now()
			}

			c.logger.Info().
				Str("library_id", libraryID).
				Str("library_name", stats.Name).
				Int64("plays", stats.Plays).
				Msg("retrieved library stats")

			return stats, nil
		}
	}

	return nil, fmt.Errorf("library %s not found", libraryID)
}

// GetUserActivity retrieves activity statistics for a user
func (c *JellystatClient) GetUserActivity(ctx context.Context, userID string, days int) (*UserActivity, error) {
	c.logger.Debug().
		Str("user_id", userID).
		Int("days", days).
		Msg("getting user activity")

	// Jellystat uses getMostActiveUsers endpoint
	var result []struct {
		UserID                string `json:"UserId"`
		UserName              string `json:"UserName"`
		Plays                 int64  `json:"Plays"`
		TotalPlaybackDuration int64  `json:"total_playback_duration"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"days": days,
		}).
		SetResult(&result).
		Post("/stats/getMostActiveUsers")

	if err != nil {
		c.logger.Error().Err(err).Str("user_id", userID).Msg("failed to get user activity")
		return nil, fmt.Errorf("failed to get user activity: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Jellystat returned status %d", resp.StatusCode())
	}

	// Find the user in the results
	for _, user := range result {
		if user.UserID == userID {
			activity := &UserActivity{
				UserID:                user.UserID,
				UserName:              user.UserName,
				Plays:                 user.Plays,
				TotalPlaybackDuration: user.TotalPlaybackDuration,
				LastActivity:          time.Now(), // Endpoint doesn't return this, use current time
			}

			c.logger.Info().
				Str("user_id", userID).
				Str("user_name", activity.UserName).
				Int64("plays", activity.Plays).
				Msg("retrieved user activity")

			return activity, nil
		}
	}

	// If user not found in most active, return zero stats
	c.logger.Warn().Str("user_id", userID).Msg("user not found in active users list")
	return &UserActivity{
		UserID:                userID,
		Plays:                 0,
		TotalPlaybackDuration: 0,
		LastActivity:          time.Now(),
	}, nil
}
