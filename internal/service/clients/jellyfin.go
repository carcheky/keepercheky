package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// JellyfinClient implements the StreamingClient interface for Jellyfin.
type JellyfinClient struct {
	client  *resty.Client
	baseURL string
	apiKey  string
	logger  *zap.Logger
}

// NewJellyfinClient creates a new Jellyfin client.
func NewJellyfinClient(config ClientConfig, logger *zap.Logger) *JellyfinClient {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	client.SetHeader("X-Emby-Token", config.APIKey) // Jellyfin uses X-Emby-Token header
	client.SetTimeout(config.Timeout)

	if config.Timeout == 0 {
		client.SetTimeout(DefaultTimeout)
	}

	return &JellyfinClient{
		client:  client,
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		logger:  logger,
	}
}

// jellyfinSystemInfo represents the system info response from Jellyfin API.
type jellyfinSystemInfo struct {
	Version                    string `json:"Version"`
	ID                         string `json:"Id"`
	OperatingSystem            string `json:"OperatingSystem"`
	OperatingSystemDisplayName string `json:"OperatingSystemDisplayName"`
	ServerName                 string `json:"ServerName"`
	LocalAddress               string `json:"LocalAddress"`
	ProductName                string `json:"ProductName"`
}

// JellyfinSystemInfo representa toda la información del sistema de Jellyfin
type JellyfinSystemInfo struct {
	Version       string `json:"version"`
	ServerID      string `json:"server_id"`
	ServerName    string `json:"server_name"`
	ProductName   string `json:"product_name"`
	OS            string `json:"os"`
	OSDisplayName string `json:"os_display_name"`
	LocalAddress  string `json:"local_address"`
}

// jellyfinItem represents a media item from Jellyfin API.
type jellyfinItem struct {
	ID          string    `json:"Id"`
	Name        string    `json:"Name"`
	Type        string    `json:"Type"` // "Movie", "Series", "Episode", etc.
	Path        string    `json:"Path"`
	DateCreated time.Time `json:"DateCreated"`
	UserData    struct {
		PlayCount        int       `json:"PlayCount"`
		IsFavorite       bool      `json:"IsFavorite"`
		LastPlayedDate   time.Time `json:"LastPlayedDate"`
		PlaybackPosition int64     `json:"PlaybackPositionTicks"`
	} `json:"UserData"`
	MediaSources []struct {
		Size int64 `json:"Size"`
	} `json:"MediaSources"`
	ImageTags map[string]string `json:"ImageTags"`
}

// jellyfinItemsResponse represents the response from items endpoint.
type jellyfinItemsResponse struct {
	Items      []jellyfinItem `json:"Items"`
	TotalCount int            `json:"TotalRecordCount"`
}

// TestConnection verifies the connection to Jellyfin.
func (c *JellyfinClient) TestConnection(ctx context.Context) error {
	return c.callWithRetry(ctx, func() error {
		var info jellyfinSystemInfo
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&info).
			Get("/System/Info/Public")

		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		c.logger.Info("Jellyfin connection successful",
			zap.String("version", info.Version),
			zap.String("server_id", info.ID),
			zap.String("url", c.baseURL),
		)

		return nil
	})
}

// GetSystemInfo retrieves complete system information from Jellyfin.
func (c *JellyfinClient) GetSystemInfo(ctx context.Context) (*JellyfinSystemInfo, error) {
	var info jellyfinSystemInfo

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&info).
			Get("/System/Info/Public")

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

	// Convertir a nuestro modelo
	systemInfo := &JellyfinSystemInfo{
		Version:       info.Version,
		ServerID:      info.ID,
		ServerName:    info.ServerName,
		ProductName:   info.ProductName,
		OS:            info.OperatingSystem,
		OSDisplayName: info.OperatingSystemDisplayName,
		LocalAddress:  info.LocalAddress,
	}

	c.logger.Info("Retrieved Jellyfin system info",
		zap.String("version", systemInfo.Version),
		zap.String("server_name", systemInfo.ServerName),
		zap.String("os", systemInfo.OS),
	)

	return systemInfo, nil
}

// GetLibrary retrieves all media items from Jellyfin.
func (c *JellyfinClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	var response jellyfinItemsResponse

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&response).
			SetQueryParams(map[string]string{
				"IncludeItemTypes": "Movie,Series",
				"Recursive":        "true",
				"Fields":           "Path,DateCreated,MediaSources,UserData",
			}).
			Get("/Items")

		if err != nil {
			return fmt.Errorf("failed to get library: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert Jellyfin items to internal Media model
	mediaList := make([]*models.Media, 0, len(response.Items))
	for _, item := range response.Items {
		media := c.convertToMedia(&item)
		mediaList = append(mediaList, media)
	}

	c.logger.Info("Retrieved Jellyfin library",
		zap.Int("total_items", response.TotalCount),
		zap.Int("retrieved", len(mediaList)),
	)

	return mediaList, nil
}

// GetPlaybackInfo retrieves playback information for media.
func (c *JellyfinClient) GetPlaybackInfo(ctx context.Context, mediaID string) (*models.PlaybackInfo, error) {
	var item jellyfinItem

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&item).
			SetPathParam("id", mediaID).
			SetQueryParam("Fields", "UserData").
			Get("/Items/{id}")

		if err != nil {
			return fmt.Errorf("failed to get playback info: %w", err)
		}

		if resp.StatusCode() == 404 {
			return fmt.Errorf("item not found: %s", mediaID)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	playbackInfo := &models.PlaybackInfo{
		MediaID:     item.ID,
		LastPlayed:  item.UserData.LastPlayedDate,
		PlayCount:   item.UserData.PlayCount,
		IsFavorite:  item.UserData.IsFavorite,
		PlaybackPos: item.UserData.PlaybackPosition,
	}

	return playbackInfo, nil
}

// DeleteItem removes a media item from Jellyfin.
func (c *JellyfinClient) DeleteItem(ctx context.Context, id string) error {
	return c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetPathParam("id", id).
			Delete("/Items/{id}")

		if err != nil {
			return fmt.Errorf("failed to delete item: %w", err)
		}

		if resp.StatusCode() != 204 && resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		c.logger.Info("Deleted item from Jellyfin",
			zap.String("item_id", id),
		)

		return nil
	})
}

// convertToMedia converts a Jellyfin item to internal Media model.
func (c *JellyfinClient) convertToMedia(item *jellyfinItem) *models.Media {
	mediaType := "movie"
	if item.Type == "Series" {
		mediaType = "series"
	}

	var size int64
	if len(item.MediaSources) > 0 {
		size = item.MediaSources[0].Size
	}

	media := &models.Media{
		Title:      item.Name,
		Type:       mediaType,
		FilePath:   item.Path,
		Size:       size,
		AddedDate:  item.DateCreated,
		JellyfinID: &item.ID,
	}

	// Set last watched if available
	if !item.UserData.LastPlayedDate.IsZero() {
		media.LastWatched = &item.UserData.LastPlayedDate
	}

	// Extract poster URL if available
	if primaryTag, ok := item.ImageTags["Primary"]; ok {
		media.PosterURL = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s",
			c.baseURL, item.ID, primaryTag)
	}

	// Add favorite tag
	if item.UserData.IsFavorite {
		media.Tags = []string{"jellyfin_favorite"}
	}

	return media
}

// callWithRetry executes a function with retry logic.
func (c *JellyfinClient) callWithRetry(ctx context.Context, fn func() error) error {
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
			c.logger.Warn("Jellyfin API call failed, retrying",
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
