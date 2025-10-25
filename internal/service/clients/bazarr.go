package clients

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// BazarrClient handles interactions with Bazarr API for subtitle management
type BazarrClient struct {
	client  *resty.Client
	baseURL string
	apiKey  string
	logger  zerolog.Logger
}

// SubtitleClient interface for subtitle management operations
type SubtitleClient interface {
	TestConnection(ctx context.Context) error
	GetMovieSubtitles(ctx context.Context, radarrID int) ([]SubtitleInfo, error)
	GetSeriesSubtitles(ctx context.Context, sonarrID int) ([]SubtitleInfo, error)
	CopySubtitlesBeforeDelete(ctx context.Context, mediaPath string) error
}

// SubtitleInfo contains subtitle information from Bazarr
type SubtitleInfo struct {
	Path     string
	Name     string
	Language string
	Forced   bool
	HI       bool // Hearing Impaired
}

// NewBazarrClient creates a new Bazarr client
func NewBazarrClient(baseURL, apiKey string, logger zerolog.Logger) *BazarrClient {
	client := resty.New()
	client.SetBaseURL(baseURL)
	client.SetHeader("X-API-Key", apiKey)
	client.SetTimeout(30)

	return &BazarrClient{
		client:  client,
		baseURL: baseURL,
		apiKey:  apiKey,
		logger:  logger.With().Str("client", "bazarr").Logger(),
	}
}

// TestConnection verifies connection to Bazarr
func (c *BazarrClient) TestConnection(ctx context.Context) error {
	c.logger.Debug().Msg("testing Bazarr connection")

	// Bazarr doesn't have a dedicated health endpoint, so we'll try to get system status
	resp, err := c.client.R().
		SetContext(ctx).
		Get("/api/system/status")

	if err != nil {
		c.logger.Error().Err(err).Msg("failed to connect to Bazarr")
		return fmt.Errorf("Bazarr connection failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		c.logger.Error().Int("status_code", resp.StatusCode()).Msg("Bazarr returned error status")
		return fmt.Errorf("Bazarr returned status %d", resp.StatusCode())
	}

	c.logger.Info().Msg("Bazarr connection successful")
	return nil
}

// GetMovieSubtitles retrieves subtitle information for a movie
func (c *BazarrClient) GetMovieSubtitles(ctx context.Context, radarrID int) ([]SubtitleInfo, error) {
	c.logger.Debug().Int("radarr_id", radarrID).Msg("getting movie subtitles")

	var result struct {
		Data []struct {
			Path     string `json:"path"`
			Name     string `json:"name"`
			Language string `json:"language"`
			Forced   bool   `json:"forced"`
			HI       bool   `json:"hearing_impaired"`
		} `json:"data"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("radarrid", fmt.Sprintf("%d", radarrID)).
		SetResult(&result).
		Get("/api/movies/subtitles")

	if err != nil {
		c.logger.Error().Err(err).Int("radarr_id", radarrID).Msg("failed to get movie subtitles")
		return nil, fmt.Errorf("failed to get movie subtitles: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Bazarr returned status %d", resp.StatusCode())
	}

	subtitles := make([]SubtitleInfo, len(result.Data))
	for i, sub := range result.Data {
		subtitles[i] = SubtitleInfo{
			Path:     sub.Path,
			Name:     sub.Name,
			Language: sub.Language,
			Forced:   sub.Forced,
			HI:       sub.HI,
		}
	}

	c.logger.Info().
		Int("radarr_id", radarrID).
		Int("subtitle_count", len(subtitles)).
		Msg("retrieved movie subtitles")

	return subtitles, nil
}

// GetSeriesSubtitles retrieves subtitle information for a series episode
func (c *BazarrClient) GetSeriesSubtitles(ctx context.Context, sonarrID int) ([]SubtitleInfo, error) {
	c.logger.Debug().Int("sonarr_id", sonarrID).Msg("getting series subtitles")

	var result struct {
		Data []struct {
			Path     string `json:"path"`
			Name     string `json:"name"`
			Language string `json:"language"`
			Forced   bool   `json:"forced"`
			HI       bool   `json:"hearing_impaired"`
		} `json:"data"`
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("episodeid", fmt.Sprintf("%d", sonarrID)).
		SetResult(&result).
		Get("/api/episodes/subtitles")

	if err != nil {
		c.logger.Error().Err(err).Int("sonarr_id", sonarrID).Msg("failed to get series subtitles")
		return nil, fmt.Errorf("failed to get series subtitles: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Bazarr returned status %d", resp.StatusCode())
	}

	subtitles := make([]SubtitleInfo, len(result.Data))
	for i, sub := range result.Data {
		subtitles[i] = SubtitleInfo{
			Path:     sub.Path,
			Name:     sub.Name,
			Language: sub.Language,
			Forced:   sub.Forced,
			HI:       sub.HI,
		}
	}

	c.logger.Info().
		Int("sonarr_id", sonarrID).
		Int("subtitle_count", len(subtitles)).
		Msg("retrieved series subtitles")

	return subtitles, nil
}

// CopySubtitlesBeforeDelete copies subtitles to a safe location before media deletion
// This is a placeholder - actual implementation would need filesystem access
func (c *BazarrClient) CopySubtitlesBeforeDelete(ctx context.Context, mediaPath string) error {
	c.logger.Info().Str("media_path", mediaPath).Msg("copying subtitles before deletion")

	// TODO: Implement actual subtitle copying logic
	// This would involve:
	// 1. Get subtitle files for the media path
	// 2. Copy them to a designated backup location
	// 3. Optionally compress them
	// 4. Log the operation

	c.logger.Warn().Msg("subtitle copying not yet implemented - placeholder")
	return nil
}
