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

	// Disable HTTP caching to always get fresh data
	client.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	client.SetHeader("Pragma", "no-cache")
	client.SetHeader("Expires", "0")

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

// sonarrSeries represents a TV series from Sonarr API with complete metadata.
type sonarrSeries struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	SortTitle         string    `json:"sortTitle"`
	Status            string    `json:"status"` // continuing, ended, upcoming
	Ended             bool      `json:"ended"`
	Overview          string    `json:"overview"`
	Network           string    `json:"network"`
	AirTime           string    `json:"airTime"` // HH:MM format
	Path              string    `json:"path"`
	Added             time.Time `json:"added"`
	Year              int       `json:"year"`
	Runtime           int       `json:"runtime"` // Minutes per episode
	TvdbID            int       `json:"tvdbId"`
	TvRageID          int       `json:"tvRageId"`
	TvMazeID          int       `json:"tvMazeId"`
	ImdbID            string    `json:"imdbId"`
	TitleSlug         string    `json:"titleSlug"`
	Certification     string    `json:"certification"` // TV-MA, TV-14, etc.
	Genres            []string  `json:"genres"`
	Tags              []int     `json:"tags"`
	QualityProfileID  int       `json:"qualityProfileId"`
	LanguageProfileID int       `json:"languageProfileId"`
	SeasonFolder      bool      `json:"seasonFolder"`
	Monitored         bool      `json:"monitored"`
	UseSceneNumbering bool      `json:"useSceneNumbering"`
	OriginalLanguage  string    `json:"originalLanguage"`
	FirstAired        time.Time `json:"firstAired"`
	LastAired         time.Time `json:"lastAired"`
	NextAiring        time.Time `json:"nextAiring"`
	PreviousAiring    time.Time `json:"previousAiring"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
		Statistics   struct {
			PreviousAiring    time.Time `json:"previousAiring"`
			EpisodeFileCount  int       `json:"episodeFileCount"`
			EpisodeCount      int       `json:"episodeCount"`
			TotalEpisodeCount int       `json:"totalEpisodeCount"`
			SizeOnDisk        int64     `json:"sizeOnDisk"`
			PercentOfEpisodes float64   `json:"percentOfEpisodes"`
		} `json:"statistics"`
	} `json:"seasons"`
	Statistics struct {
		SeasonCount       int     `json:"seasonCount"`
		EpisodeCount      int     `json:"episodeCount"`
		EpisodeFileCount  int     `json:"episodeFileCount"`
		TotalEpisodeCount int     `json:"totalEpisodeCount"`
		SizeOnDisk        int64   `json:"sizeOnDisk"`
		PercentOfEpisodes float64 `json:"percentOfEpisodes"`
	} `json:"statistics"`
	Images []struct {
		CoverType string `json:"coverType"` // poster, banner, fanart
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"images"`
}

// sonarrTag represents a tag from Sonarr API.
type sonarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// sonarrQueueItem represents an item in the download queue.
type sonarrQueueItem struct {
	ID                      int       `json:"id"`
	SeriesID                int       `json:"seriesId"`
	EpisodeID               int       `json:"episodeId"`
	Title                   string    `json:"title"`
	Size                    int64     `json:"size"`
	Sizeleft                int64     `json:"sizeleft"`
	Status                  string    `json:"status"`
	TrackedDownloadStatus   string    `json:"trackedDownloadStatus"`
	TrackedDownloadState    string    `json:"trackedDownloadState"`
	StatusMessages          []string  `json:"statusMessages"`
	DownloadID              string    `json:"downloadId"`
	Protocol                string    `json:"protocol"`
	DownloadClient          string    `json:"downloadClient"`
	Indexer                 string    `json:"indexer"`
	OutputPath              string    `json:"outputPath"`
	TimedOut                bool      `json:"timedOut"`
	EstimatedCompletionTime time.Time `json:"estimatedCompletionTime"`
}

// sonarrQueueResponse represents the queue response.
type sonarrQueueResponse struct {
	Page         int               `json:"page"`
	PageSize     int               `json:"pageSize"`
	TotalRecords int               `json:"totalRecords"`
	Records      []sonarrQueueItem `json:"records"`
}

// SonarrQueueItem represents a processed queue item for public API.
type SonarrQueueItem struct {
	ID                  int       `json:"id"`
	SeriesID            int       `json:"series_id"`
	EpisodeID           int       `json:"episode_id"`
	Title               string    `json:"title"`
	Size                int64     `json:"size"`
	SizeLeft            int64     `json:"size_left"`
	Progress            float64   `json:"progress"`
	Status              string    `json:"status"`
	DownloadStatus      string    `json:"download_status"`
	DownloadState       string    `json:"download_state"`
	Protocol            string    `json:"protocol"`
	DownloadClient      string    `json:"download_client"`
	Indexer             string    `json:"indexer"`
	TimedOut            bool      `json:"timed_out"`
	EstimatedCompletion time.Time `json:"estimated_completion"`
}

// sonarrHistoryItem represents a history entry.
type sonarrHistoryItem struct {
	ID          int    `json:"id"`
	EpisodeID   int    `json:"episodeId"`
	SeriesID    int    `json:"seriesId"`
	SourceTitle string `json:"sourceTitle"`
	Quality     struct {
		Quality struct {
			Name string `json:"name"`
		} `json:"quality"`
	} `json:"quality"`
	Date       time.Time `json:"date"`
	EventType  string    `json:"eventType"`
	DownloadID string    `json:"downloadId"`
}

// sonarrHistoryResponse represents the history response.
type sonarrHistoryResponse struct {
	Page         int                 `json:"page"`
	PageSize     int                 `json:"pageSize"`
	TotalRecords int                 `json:"totalRecords"`
	Records      []sonarrHistoryItem `json:"records"`
}

// SonarrHistoryItem represents a processed history item for public API.
type SonarrHistoryItem struct {
	ID          int       `json:"id"`
	EpisodeID   int       `json:"episode_id"`
	SeriesID    int       `json:"series_id"`
	SourceTitle string    `json:"source_title"`
	Quality     string    `json:"quality"`
	Date        time.Time `json:"date"`
	EventType   string    `json:"event_type"`
	DownloadID  string    `json:"download_id"`
}

// sonarrCalendarItem represents a calendar entry.
type sonarrCalendarItem struct {
	ID            int       `json:"id"`
	SeriesID      int       `json:"seriesId"`
	EpisodeFileID int       `json:"episodeFileId"`
	SeasonNumber  int       `json:"seasonNumber"`
	EpisodeNumber int       `json:"episodeNumber"`
	Title         string    `json:"title"`
	AirDate       string    `json:"airDate"`
	AirDateUtc    time.Time `json:"airDateUtc"`
	HasFile       bool      `json:"hasFile"`
	Monitored     bool      `json:"monitored"`
	Series        struct {
		Title string `json:"title"`
	} `json:"series"`
}

// SonarrCalendarItem represents a processed calendar item for public API.
type SonarrCalendarItem struct {
	ID            int       `json:"id"`
	SeriesID      int       `json:"series_id"`
	SeriesTitle   string    `json:"series_title"`
	SeasonNumber  int       `json:"season_number"`
	EpisodeNumber int       `json:"episode_number"`
	Title         string    `json:"title"`
	AirDate       string    `json:"air_date"`
	AirDateUtc    time.Time `json:"air_date_utc"`
	HasFile       bool      `json:"has_file"`
	Monitored     bool      `json:"monitored"`
}

// sonarrQualityProfile represents a quality profile.
type sonarrQualityProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SonarrQualityProfile represents a quality profile for public API.
type SonarrQualityProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// sonarrEpisode represents an episode from Sonarr API with complete metadata.
type sonarrEpisode struct {
	ID                         int       `json:"id"`
	SeriesID                   int       `json:"seriesId"`
	TvdbID                     int       `json:"tvdbId"`
	EpisodeFileID              int       `json:"episodeFileId"`
	SeasonNumber               int       `json:"seasonNumber"`
	EpisodeNumber              int       `json:"episodeNumber"`
	Title                      string    `json:"title"`
	AirDate                    string    `json:"airDate"` // YYYY-MM-DD format
	AirDateUtc                 time.Time `json:"airDateUtc"`
	Overview                   string    `json:"overview"`
	HasFile                    bool      `json:"hasFile"`
	Monitored                  bool      `json:"monitored"`
	AbsoluteEpisodeNumber      int       `json:"absoluteEpisodeNumber"`
	SceneAbsoluteEpisodeNumber int       `json:"sceneAbsoluteEpisodeNumber"`
	SceneEpisodeNumber         int       `json:"sceneEpisodeNumber"`
	SceneSeasonNumber          int       `json:"sceneSeasonNumber"`
	UnverifiedSceneNumbering   bool      `json:"unverifiedSceneNumbering"`
	Grabbed                    bool      `json:"grabbed"`
}

// SonarrEpisode represents a processed episode for public API.
type SonarrEpisode struct {
	ID                    int       `json:"id"`
	SeriesID              int       `json:"series_id"`
	TvdbID                int       `json:"tvdb_id"`
	EpisodeFileID         int       `json:"episode_file_id"`
	SeasonNumber          int       `json:"season_number"`
	EpisodeNumber         int       `json:"episode_number"`
	Title                 string    `json:"title"`
	AirDate               string    `json:"air_date"`
	AirDateUtc            time.Time `json:"air_date_utc"`
	Overview              string    `json:"overview"`
	HasFile               bool      `json:"has_file"`
	Monitored             bool      `json:"monitored"`
	AbsoluteEpisodeNumber int       `json:"absolute_episode_number"`
	Grabbed               bool      `json:"grabbed"`
}

// sonarrEpisodeFile represents an episode file from Sonarr API with complete metadata.
type sonarrEpisodeFile struct {
	ID           int       `json:"id"`
	SeriesID     int       `json:"seriesId"`
	SeasonNumber int       `json:"seasonNumber"`
	RelativePath string    `json:"relativePath"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	DateAdded    time.Time `json:"dateAdded"`
	SceneName    string    `json:"sceneName"`
	ReleaseGroup string    `json:"releaseGroup"`
	Quality      struct {
		Quality struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Source     string `json:"source"`
			Resolution int    `json:"resolution"`
		} `json:"quality"`
		Revision struct {
			Version  int  `json:"version"`
			Real     int  `json:"real"`
			IsRepack bool `json:"isRepack"`
		} `json:"revision"`
	} `json:"quality"`
	MediaInfo struct {
		AudioBitrate     int     `json:"audioBitrate"`
		AudioChannels    float64 `json:"audioChannels"`
		AudioCodec       string  `json:"audioCodec"`
		AudioLanguages   string  `json:"audioLanguages"`
		AudioStreamCount int     `json:"audioStreamCount"`
		VideoBitDepth    int     `json:"videoBitDepth"`
		VideoBitrate     int     `json:"videoBitrate"`
		VideoCodec       string  `json:"videoCodec"`
		VideoFps         float64 `json:"videoFps"`
		Resolution       string  `json:"resolution"`
		RunTime          string  `json:"runTime"`
		ScanType         string  `json:"scanType"`
		Subtitles        string  `json:"subtitles"`
	} `json:"mediaInfo"`
	QualityCutoffNotMet  bool `json:"qualityCutoffNotMet"`
	LanguageCutoffNotMet bool `json:"languageCutoffNotMet"`
}

// SonarrEpisodeFile represents a processed episode file for public API.
type SonarrEpisodeFile struct {
	ID               int       `json:"id"`
	SeriesID         int       `json:"series_id"`
	SeasonNumber     int       `json:"season_number"`
	RelativePath     string    `json:"relative_path"`
	Path             string    `json:"path"`
	Size             int64     `json:"size"`
	DateAdded        time.Time `json:"date_added"`
	SceneName        string    `json:"scene_name"`
	ReleaseGroup     string    `json:"release_group"`
	QualityName      string    `json:"quality_name"`
	AudioCodec       string    `json:"audio_codec"`
	AudioChannels    float64   `json:"audio_channels"`
	VideoCodec       string    `json:"video_codec"`
	Resolution       string    `json:"resolution"`
	VideoRange       string    `json:"video_range,omitempty"`
	RunTime          string    `json:"run_time"`
	AudioLanguages   string    `json:"audio_languages,omitempty"`
	Subtitles        string    `json:"subtitles,omitempty"`
	QualityCutoffMet bool      `json:"quality_cutoff_met"`
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

// GetQueue retrieves the current download queue from Sonarr.
func (c *SonarrClient) GetQueue(ctx context.Context) ([]SonarrQueueItem, error) {
	var queueResp sonarrQueueResponse

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&queueResp).
			SetQueryParam("pageSize", "100").
			SetQueryParam("includeUnknownSeriesItems", "false").
			Get("/api/v3/queue")

		if err != nil {
			return fmt.Errorf("failed to get queue: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrQueueItem, 0, len(queueResp.Records))
	for _, record := range queueResp.Records {
		progress := 0.0
		if record.Size > 0 {
			progress = float64(record.Size-record.Sizeleft) / float64(record.Size) * 100
			// Ensure progress is within 0-100 range (handle inconsistent data from Sonarr)
			if progress < 0 {
				progress = 0
			} else if progress > 100 {
				progress = 100
			}
		}

		items = append(items, SonarrQueueItem{
			ID:                  record.ID,
			SeriesID:            record.SeriesID,
			EpisodeID:           record.EpisodeID,
			Title:               record.Title,
			Size:                record.Size,
			SizeLeft:            record.Sizeleft,
			Progress:            progress,
			Status:              record.Status,
			DownloadStatus:      record.TrackedDownloadStatus,
			DownloadState:       record.TrackedDownloadState,
			Protocol:            record.Protocol,
			DownloadClient:      record.DownloadClient,
			Indexer:             record.Indexer,
			TimedOut:            record.TimedOut,
			EstimatedCompletion: record.EstimatedCompletionTime,
		})
	}

	c.logger.Info("Retrieved Sonarr queue",
		zap.Int("total_items", len(items)),
	)

	return items, nil
}

// GetHistory retrieves recent history from Sonarr.
func (c *SonarrClient) GetHistory(ctx context.Context, pageSize int) ([]SonarrHistoryItem, error) {
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var historyResp sonarrHistoryResponse

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&historyResp).
			SetQueryParam("pageSize", fmt.Sprintf("%d", pageSize)).
			SetQueryParam("sortKey", "date").
			SetQueryParam("sortDirection", "descending").
			Get("/api/v3/history")

		if err != nil {
			return fmt.Errorf("failed to get history: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrHistoryItem, 0, len(historyResp.Records))
	for _, record := range historyResp.Records {
		items = append(items, SonarrHistoryItem{
			ID:          record.ID,
			EpisodeID:   record.EpisodeID,
			SeriesID:    record.SeriesID,
			SourceTitle: record.SourceTitle,
			Quality:     record.Quality.Quality.Name,
			Date:        record.Date,
			EventType:   record.EventType,
			DownloadID:  record.DownloadID,
		})
	}

	c.logger.Info("Retrieved Sonarr history",
		zap.Int("total_items", len(items)),
	)

	return items, nil
}

// GetCalendar retrieves upcoming episodes from Sonarr.
func (c *SonarrClient) GetCalendar(ctx context.Context, startDate, endDate time.Time) ([]SonarrCalendarItem, error) {
	var calendarItems []sonarrCalendarItem

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&calendarItems).
			SetQueryParam("start", startDate.Format("2006-01-02")).
			SetQueryParam("end", endDate.Format("2006-01-02")).
			Get("/api/v3/calendar")

		if err != nil {
			return fmt.Errorf("failed to get calendar: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrCalendarItem, 0, len(calendarItems))
	for _, record := range calendarItems {
		items = append(items, SonarrCalendarItem{
			ID:            record.ID,
			SeriesID:      record.SeriesID,
			SeriesTitle:   record.Series.Title,
			SeasonNumber:  record.SeasonNumber,
			EpisodeNumber: record.EpisodeNumber,
			Title:         record.Title,
			AirDate:       record.AirDate,
			AirDateUtc:    record.AirDateUtc,
			HasFile:       record.HasFile,
			Monitored:     record.Monitored,
		})
	}

	c.logger.Info("Retrieved Sonarr calendar",
		zap.Int("total_items", len(items)),
		zap.String("start_date", startDate.Format("2006-01-02")),
		zap.String("end_date", endDate.Format("2006-01-02")),
	)

	return items, nil
}

// GetQualityProfiles retrieves available quality profiles from Sonarr.
func (c *SonarrClient) GetQualityProfiles(ctx context.Context) ([]SonarrQualityProfile, error) {
	var profiles []sonarrQualityProfile

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&profiles).
			Get("/api/v3/qualityprofile")

		if err != nil {
			return fmt.Errorf("failed to get quality profiles: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrQualityProfile, 0, len(profiles))
	for _, profile := range profiles {
		items = append(items, SonarrQualityProfile{
			ID:   profile.ID,
			Name: profile.Name,
		})
	}

	c.logger.Info("Retrieved Sonarr quality profiles",
		zap.Int("total_profiles", len(items)),
	)

	return items, nil
}

// GetEpisodes retrieves all episodes for a specific series.
func (c *SonarrClient) GetEpisodes(ctx context.Context, seriesID int) ([]SonarrEpisode, error) {
	var episodes []sonarrEpisode

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&episodes).
			SetQueryParam("seriesId", fmt.Sprintf("%d", seriesID)).
			Get("/api/v3/episode")

		if err != nil {
			return fmt.Errorf("failed to get episodes: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrEpisode, 0, len(episodes))
	for _, ep := range episodes {
		items = append(items, SonarrEpisode{
			ID:                    ep.ID,
			SeriesID:              ep.SeriesID,
			TvdbID:                ep.TvdbID,
			EpisodeFileID:         ep.EpisodeFileID,
			SeasonNumber:          ep.SeasonNumber,
			EpisodeNumber:         ep.EpisodeNumber,
			Title:                 ep.Title,
			AirDate:               ep.AirDate,
			AirDateUtc:            ep.AirDateUtc,
			Overview:              ep.Overview,
			HasFile:               ep.HasFile,
			Monitored:             ep.Monitored,
			AbsoluteEpisodeNumber: ep.AbsoluteEpisodeNumber,
			Grabbed:               ep.Grabbed,
		})
	}

	c.logger.Info("Retrieved Sonarr episodes",
		zap.Int("series_id", seriesID),
		zap.Int("total_episodes", len(items)),
	)

	return items, nil
}

// GetEpisodeFiles retrieves all episode files for a specific series with complete metadata.
func (c *SonarrClient) GetEpisodeFiles(ctx context.Context, seriesID int) ([]SonarrEpisodeFile, error) {
	var files []sonarrEpisodeFile

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&files).
			SetQueryParam("seriesId", fmt.Sprintf("%d", seriesID)).
			Get("/api/v3/episodefile")

		if err != nil {
			return fmt.Errorf("failed to get episode files: %w", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to public API model
	items := make([]SonarrEpisodeFile, 0, len(files))
	for _, file := range files {
		items = append(items, SonarrEpisodeFile{
			ID:               file.ID,
			SeriesID:         file.SeriesID,
			SeasonNumber:     file.SeasonNumber,
			RelativePath:     file.RelativePath,
			Path:             file.Path,
			Size:             file.Size,
			DateAdded:        file.DateAdded,
			SceneName:        file.SceneName,
			ReleaseGroup:     file.ReleaseGroup,
			QualityName:      file.Quality.Quality.Name,
			AudioCodec:       file.MediaInfo.AudioCodec,
			AudioChannels:    file.MediaInfo.AudioChannels,
			VideoCodec:       file.MediaInfo.VideoCodec,
			Resolution:       file.MediaInfo.Resolution,
			RunTime:          file.MediaInfo.RunTime,
			AudioLanguages:   file.MediaInfo.AudioLanguages,
			Subtitles:        file.MediaInfo.Subtitles,
			QualityCutoffMet: !file.QualityCutoffNotMet,
		})
	}

	c.logger.Info("Retrieved Sonarr episode files",
		zap.Int("series_id", seriesID),
		zap.Int("total_files", len(items)),
	)

	return items, nil
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
