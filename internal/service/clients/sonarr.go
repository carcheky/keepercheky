package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// SonarrClient implements the MediaClient interface for Sonarr.
type SonarrClient struct {
	client  *resty.Client
	baseURL string
	apiKey  string
	logger  *zap.Logger
}

// NewSonarrClient creates a new Sonarr client.
func NewSonarrClient(config ClientConfig, logger *zap.Logger) *SonarrClient {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	client.SetHeader("X-Api-Key", config.APIKey)
	client.SetTimeout(config.Timeout)

	if config.Timeout == 0 {
		client.SetTimeout(DefaultTimeout)
	}

	return &SonarrClient{
		client:  client,
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		logger:  logger,
	}
}

// sonarrSystemStatus represents the system status response from Sonarr API.
type sonarrSystemStatus struct {
	Version           string `json:"version"`
	BuildTime         string `json:"buildTime"`
	IsDebug           bool   `json:"isDebug"`
	IsProduction      bool   `json:"isProduction"`
	IsAdmin           bool   `json:"isAdmin"`
	IsUserInteractive bool   `json:"isUserInteractive"`
	StartupPath       string `json:"startupPath"`
	AppData           string `json:"appData"`
	OsName            string `json:"osName"`
	OsVersion         string `json:"osVersion"`
	IsMonoRuntime     bool   `json:"isMonoRuntime"`
	IsMono            bool   `json:"isMono"`
	IsLinux           bool   `json:"isLinux"`
	IsOsx             bool   `json:"isOsx"`
	IsWindows         bool   `json:"isWindows"`
	Mode              string `json:"mode"`
	Branch            string `json:"branch"`
	Authentication    string `json:"authentication"`
	SqliteVersion     string `json:"sqliteVersion"`
	UrlBase           string `json:"urlBase"`
	RuntimeVersion    string `json:"runtimeVersion"`
	RuntimeName       string `json:"runtimeName"`
}

// SonarrSystemInfo representa toda la información del sistema de Sonarr
type SonarrSystemInfo struct {
	Version        string `json:"version"`
	BuildTime      string `json:"build_time"`
	Branch         string `json:"branch"`
	OS             string `json:"os"`
	OSVersion      string `json:"os_version"`
	Runtime        string `json:"runtime"`
	RuntimeVersion string `json:"runtime_version"`
	IsDebug        bool   `json:"is_debug"`
	IsProduction   bool   `json:"is_production"`
	Authentication string `json:"authentication"`
	URLBase        string `json:"url_base"`
	StartupPath    string `json:"startup_path"`
	AppData        string `json:"app_data"`
	SqliteVersion  string `json:"sqlite_version"`
}

// sonarrSeries represents a TV series from Sonarr API.
type sonarrSeries struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	Added      time.Time `json:"added"`
	Tags       []int     `json:"tags"`
	Statistics struct {
		SeasonCount       int     `json:"seasonCount"`
		EpisodeCount      int     `json:"episodeCount"`
		EpisodeFileCount  int     `json:"episodeFileCount"`
		TotalEpisodeCount int     `json:"totalEpisodeCount"`
		SizeOnDisk        int64   `json:"sizeOnDisk"`
		PercentOfEpisodes float64 `json:"percentOfEpisodes"`
	} `json:"statistics"`
	Images []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
}

// sonarrEpisodeFile represents an episode file from Sonarr API.
type sonarrEpisodeFile struct {
	ID       int    `json:"id"`
	SeriesID int    `json:"seriesId"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Quality  struct {
		Quality struct {
			Name string `json:"name"`
		} `json:"quality"`
	} `json:"quality"`
}

// sonarrTag represents a tag from Sonarr API.
type sonarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// TestConnection verifies the connection to Sonarr.
func (c *SonarrClient) TestConnection(ctx context.Context) error {
	return c.callWithRetry(ctx, func() error {
		var status sonarrSystemStatus
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&status).
			Get("/api/v3/system/status")

		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		c.logger.Info("Sonarr connection successful",
			zap.String("version", status.Version),
			zap.String("url", c.baseURL),
		)

		return nil
	})
}

// GetSystemInfo retrieves complete system information from Sonarr.
func (c *SonarrClient) GetSystemInfo(ctx context.Context) (*SonarrSystemInfo, error) {
	var status sonarrSystemStatus

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&status).
			Get("/api/v3/system/status")

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

	// Construir el nombre del OS
	osName := status.OsName
	if status.IsLinux {
		osName = "Linux"
	} else if status.IsWindows {
		osName = "Windows"
	} else if status.IsOsx {
		osName = "macOS"
	}

	// Convertir a nuestro modelo
	info := &SonarrSystemInfo{
		Version:        status.Version,
		BuildTime:      status.BuildTime,
		Branch:         status.Branch,
		OS:             osName,
		OSVersion:      status.OsVersion,
		Runtime:        status.RuntimeName,
		RuntimeVersion: status.RuntimeVersion,
		IsDebug:        status.IsDebug,
		IsProduction:   status.IsProduction,
		Authentication: status.Authentication,
		URLBase:        status.UrlBase,
		StartupPath:    status.StartupPath,
		AppData:        status.AppData,
		SqliteVersion:  status.SqliteVersion,
	}

	c.logger.Info("Retrieved Sonarr system info",
		zap.String("version", info.Version),
		zap.String("os", info.OS),
		zap.String("runtime", info.Runtime),
	)

	return info, nil
}

// GetLibrary retrieves all TV series from Sonarr.
func (c *SonarrClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	var sonarrSeries []sonarrSeries

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&sonarrSeries).
			Get("/api/v3/series")

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

	// Convert Sonarr series to internal Media model
	mediaList := make([]*models.Media, 0, len(sonarrSeries))
	for _, series := range sonarrSeries {
		if series.Statistics.EpisodeFileCount == 0 {
			continue // Skip series without files
		}

		media := c.convertToMedia(&series)
		mediaList = append(mediaList, media)
	}

	c.logger.Info("Retrieved Sonarr library",
		zap.Int("total_series", len(sonarrSeries)),
		zap.Int("with_files", len(mediaList)),
	)

	return mediaList, nil
}

// GetItem retrieves a specific TV series from Sonarr.
func (c *SonarrClient) GetItem(ctx context.Context, id int) (*models.Media, error) {
	var series sonarrSeries

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&series).
			SetPathParam("id", fmt.Sprintf("%d", id)).
			Get("/api/v3/series/{id}")

		if err != nil {
			return fmt.Errorf("failed to get series: %w", err)
		}

		if resp.StatusCode() == 404 {
			return fmt.Errorf("series not found: %d", id)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c.convertToMedia(&series), nil
}

// DeleteItem removes a TV series from Sonarr.
func (c *SonarrClient) DeleteItem(ctx context.Context, id int, deleteFiles bool) error {
	return c.callWithRetry(ctx, func() error {
		deleteFilesStr := "false"
		if deleteFiles {
			deleteFilesStr = "true"
		}

		resp, err := c.client.R().
			SetContext(ctx).
			SetPathParam("id", fmt.Sprintf("%d", id)).
			SetQueryParam("deleteFiles", deleteFilesStr).
			SetQueryParam("addImportListExclusion", "false").
			Delete("/api/v3/series/{id}")

		if err != nil {
			return fmt.Errorf("failed to delete series: %w", err)
		}

		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		c.logger.Info("Deleted series from Sonarr",
			zap.Int("series_id", id),
			zap.Bool("deleted_files", deleteFiles),
		)

		return nil
	})
}

// GetTags retrieves all tags from Sonarr.
func (c *SonarrClient) GetTags(ctx context.Context) ([]models.Tag, error) {
	var sonarrTags []sonarrTag

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&sonarrTags).
			Get("/api/v3/tag")

		if err != nil {
			return fmt.Errorf("failed to get tags: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to internal Tag model
	tags := make([]models.Tag, 0, len(sonarrTags))
	for _, tag := range sonarrTags {
		tags = append(tags, models.Tag{
			ServiceType: "sonarr",
			ServiceID:   tag.ID,
			Label:       tag.Label,
		})
	}

	return tags, nil
}

// convertToMedia converts a Sonarr series to internal Media model.
func (c *SonarrClient) convertToMedia(series *sonarrSeries) *models.Media {
	media := &models.Media{
		Title:            series.Title,
		Type:             "series",
		FilePath:         series.Path,
		Size:             series.Statistics.SizeOnDisk,
		AddedDate:        series.Added,
		SonarrID:         &series.ID,
		EpisodeCount:     series.Statistics.TotalEpisodeCount,
		SeasonCount:      series.Statistics.SeasonCount,
		EpisodeFileCount: series.Statistics.EpisodeFileCount,
		// Quality will be determined from episode files
		Quality: "Unknown",
	}

	// Extract poster URL
	for _, image := range series.Images {
		if image.CoverType == "poster" {
			media.PosterURL = image.URL
			break
		}
	}

	// Convert tag IDs to strings
	if len(series.Tags) > 0 {
		media.Tags = make([]string, 0, len(series.Tags))
		for _, tagID := range series.Tags {
			media.Tags = append(media.Tags, fmt.Sprintf("sonarr_tag_%d", tagID))
		}
	}

	return media
}

// callWithRetry executes a function with retry logic.
func (c *SonarrClient) callWithRetry(ctx context.Context, fn func() error) error {
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
			c.logger.Warn("Sonarr API call failed, retrying",
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
