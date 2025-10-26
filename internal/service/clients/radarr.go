package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// RadarrClient implements the MediaClient interface for Radarr.
type RadarrClient struct {
	client  *resty.Client
	baseURL string
	apiKey  string
	logger  *zap.Logger
}

// NewRadarrClient creates a new Radarr client.
func NewRadarrClient(config ClientConfig, logger *zap.Logger) *RadarrClient {
	client := resty.New()
	client.SetBaseURL(config.BaseURL)
	client.SetHeader("X-Api-Key", config.APIKey)
	client.SetTimeout(config.Timeout)

	if config.Timeout == 0 {
		client.SetTimeout(DefaultTimeout)
	}

	return &RadarrClient{
		client:  client,
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		logger:  logger,
	}
}

// radarrSystemStatus represents the system status response from Radarr API.
type radarrSystemStatus struct {
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

// RadarrSystemInfo representa toda la información del sistema de Radarr
type RadarrSystemInfo struct {
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

// radarrMovie represents a movie from Radarr API.
type radarrMovie struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	SizeOnDisk int64     `json:"sizeOnDisk"`
	Added      time.Time `json:"added"`
	Quality    struct {
		Quality struct {
			Name string `json:"name"`
		} `json:"quality"`
	} `json:"quality"`
	Tags    []int `json:"tags"`
	HasFile bool  `json:"hasFile"`
	Images  []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
}

// radarrTag represents a tag from Radarr API.
type radarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// TestConnection verifies the connection to Radarr.
func (c *RadarrClient) TestConnection(ctx context.Context) error {
	return c.callWithRetry(ctx, func() error {
		var status radarrSystemStatus
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

		c.logger.Info("Radarr connection successful",
			zap.String("version", status.Version),
			zap.String("url", c.baseURL),
		)

		return nil
	})
}

// GetSystemInfo retrieves complete system information from Radarr.
func (c *RadarrClient) GetSystemInfo(ctx context.Context) (*RadarrSystemInfo, error) {
	var status radarrSystemStatus

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
	info := &RadarrSystemInfo{
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

	c.logger.Info("Retrieved Radarr system info",
		zap.String("version", info.Version),
		zap.String("os", info.OS),
		zap.String("runtime", info.Runtime),
	)

	return info, nil
}

// GetLibrary retrieves all movies from Radarr.
func (c *RadarrClient) GetLibrary(ctx context.Context) ([]*models.Media, error) {
	var radarrMovies []radarrMovie

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&radarrMovies).
			Get("/api/v3/movie")

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

	// Convert Radarr movies to internal Media model
	mediaList := make([]*models.Media, 0, len(radarrMovies))
	for _, movie := range radarrMovies {
		if !movie.HasFile {
			continue // Skip movies without files
		}

		media := c.convertToMedia(&movie)
		mediaList = append(mediaList, media)
	}

	c.logger.Info("Retrieved Radarr library",
		zap.Int("total_movies", len(radarrMovies)),
		zap.Int("with_files", len(mediaList)),
	)

	return mediaList, nil
}

// GetItem retrieves a specific movie from Radarr.
func (c *RadarrClient) GetItem(ctx context.Context, id int) (*models.Media, error) {
	var movie radarrMovie

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&movie).
			SetPathParam("id", fmt.Sprintf("%d", id)).
			Get("/api/v3/movie/{id}")

		if err != nil {
			return fmt.Errorf("failed to get movie: %w", err)
		}

		if resp.StatusCode() == 404 {
			return fmt.Errorf("movie not found: %d", id)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c.convertToMedia(&movie), nil
}

// DeleteItem removes a movie from Radarr.
func (c *RadarrClient) DeleteItem(ctx context.Context, id int, deleteFiles bool) error {
	return c.callWithRetry(ctx, func() error {
		deleteFilesStr := "false"
		if deleteFiles {
			deleteFilesStr = "true"
		}

		resp, err := c.client.R().
			SetContext(ctx).
			SetPathParam("id", fmt.Sprintf("%d", id)).
			SetQueryParam("deleteFiles", deleteFilesStr).
			SetQueryParam("addImportExclusion", "false").
			Delete("/api/v3/movie/{id}")

		if err != nil {
			return fmt.Errorf("failed to delete movie: %w", err)
		}

		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		c.logger.Info("Deleted movie from Radarr",
			zap.Int("movie_id", id),
			zap.Bool("deleted_files", deleteFiles),
		)

		return nil
	})
}

// GetTags retrieves all tags from Radarr.
func (c *RadarrClient) GetTags(ctx context.Context) ([]models.Tag, error) {
	var radarrTags []radarrTag

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&radarrTags).
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
	tags := make([]models.Tag, 0, len(radarrTags))
	for _, tag := range radarrTags {
		tags = append(tags, models.Tag{
			ServiceType: "radarr",
			ServiceID:   tag.ID,
			Label:       tag.Label,
		})
	}

	return tags, nil
}

// convertToMedia converts a Radarr movie to internal Media model.
func (c *RadarrClient) convertToMedia(movie *radarrMovie) *models.Media {
	media := &models.Media{
		Title:     movie.Title,
		Type:      "movie",
		FilePath:  movie.Path,
		Size:      movie.SizeOnDisk,
		AddedDate: movie.Added,
		RadarrID:  &movie.ID,
		Quality:   movie.Quality.Quality.Name,
	}

	// Extract poster URL
	for _, image := range movie.Images {
		if image.CoverType == "poster" {
			media.PosterURL = image.URL
			break
		}
	}

	// Convert tag IDs to strings (we'll need to map these to tag labels later)
	if len(movie.Tags) > 0 {
		media.Tags = make([]string, 0, len(movie.Tags))
		for _, tagID := range movie.Tags {
			media.Tags = append(media.Tags, fmt.Sprintf("radarr_tag_%d", tagID))
		}
	}

	return media
}

// callWithRetry executes a function with retry logic.
func (c *RadarrClient) callWithRetry(ctx context.Context, fn func() error) error {
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
			c.logger.Warn("Radarr API call failed, retrying",
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
