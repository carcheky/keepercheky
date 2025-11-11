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

	// Disable HTTP caching to always get fresh data
	client.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	client.SetHeader("Pragma", "no-cache")
	client.SetHeader("Expires", "0")

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

// radarrMovie represents a movie from Radarr API with complete metadata.
type radarrMovie struct {
	// Basic Info
	ID            int    `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	SortTitle     string `json:"sortTitle"`
	CleanTitle    string `json:"cleanTitle"`
	Year          int    `json:"year"`
	SecondaryYear *int   `json:"secondaryYear"`

	// Status and Monitoring
	Status              string `json:"status"` // tba, announced, inCinemas, released, deleted
	Overview            string `json:"overview"`
	Monitored           bool   `json:"monitored"`
	MinimumAvailability string `json:"minimumAvailability"` // tba, announced, inCinemas, released, preDB
	IsAvailable         bool   `json:"isAvailable"`

	// Dates
	InCinemas       *time.Time `json:"inCinemas"`
	PhysicalRelease *time.Time `json:"physicalRelease"`
	DigitalRelease  *time.Time `json:"digitalRelease"`
	Added           time.Time  `json:"added"`

	// File Info
	Path             string `json:"path"`
	FolderName       string `json:"folderName"`
	SizeOnDisk       int64  `json:"sizeOnDisk"`
	HasFile          bool   `json:"hasFile"`
	MovieFileID      int    `json:"movieFileId"`
	QualityProfileID int    `json:"qualityProfileId"`

	// Technical Details
	Runtime int `json:"runtime"` // minutes

	// External IDs
	IMDbID    string `json:"imdbId"`
	TMDbID    int    `json:"tmdbId"`
	TitleSlug string `json:"titleSlug"`

	// Media Info
	Website          string   `json:"website"`
	RemotePoster     string   `json:"remotePoster"`
	YouTubeTrailerID string   `json:"youTubeTrailerId"`
	Studio           string   `json:"studio"`
	Certification    string   `json:"certification"` // G, PG, PG-13, R, NC-17
	Genres           []string `json:"genres"`
	Tags             []int    `json:"tags"`
	Popularity       float64  `json:"popularity"`

	// Images
	Images []struct {
		CoverType string `json:"coverType"` // poster, fanart, banner, logo
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"images"`

	// Alternate Titles
	AlternateTitles []struct {
		SourceType      string `json:"sourceType"`
		MovieMetadataID int    `json:"movieMetadataId"`
		Title           string `json:"title"`
		ID              int    `json:"id"`
	} `json:"alternateTitles"`

	// Ratings
	Ratings struct {
		IMDb *struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"imdb"`
		TMDb *struct {
			Votes int     `json:"votes"`
			Value float64 `json:"value"`
			Type  string  `json:"type"`
		} `json:"tmdb"`
		Metacritic *struct {
			Votes int    `json:"votes"`
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"metacritic"`
		RottenTomatoes *struct {
			Votes int    `json:"votes"`
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"rottenTomatoes"`
	} `json:"ratings"`

	// Movie File (complete details)
	MovieFile *struct {
		MovieID      int       `json:"movieId"`
		RelativePath string    `json:"relativePath"`
		Path         string    `json:"path"`
		Size         int64     `json:"size"`
		DateAdded    time.Time `json:"dateAdded"`
		SceneName    string    `json:"sceneName"`
		IndexerFlags int       `json:"indexerFlags"`
		Quality      struct {
			Quality struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Source     string `json:"source"`     // bluray, webdl, webrip, hdtv
				Resolution int    `json:"resolution"` // 480, 720, 1080, 2160
				Modifier   string `json:"modifier"`   // remux, brdisk, regional, none
			} `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		CustomFormatScore int `json:"customFormatScore"`
		CustomFormats     []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"customFormats"`
		MediaInfo *struct {
			AudioBitrate          int     `json:"audioBitrate"`
			AudioChannels         float64 `json:"audioChannels"`
			AudioCodec            string  `json:"audioCodec"`
			AudioLanguages        string  `json:"audioLanguages"`
			AudioStreamCount      int     `json:"audioStreamCount"`
			VideoBitDepth         int     `json:"videoBitDepth"`
			VideoBitrate          int     `json:"videoBitrate"`
			VideoCodec            string  `json:"videoCodec"`
			VideoFps              float64 `json:"videoFps"`
			VideoDynamicRange     string  `json:"videoDynamicRange"`     // SDR, HDR
			VideoDynamicRangeType string  `json:"videoDynamicRangeType"` // HDR10, DolbyVision
			Resolution            string  `json:"resolution"`
			RunTime               string  `json:"runTime"`
			ScanType              string  `json:"scanType"` // Progressive, Interlaced
			Subtitles             string  `json:"subtitles"`
		} `json:"mediaInfo"`
		QualityCutoffNotMet bool `json:"qualityCutoffNotMet"`
		Languages           []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"languages"`
		ReleaseGroup string `json:"releaseGroup"`
		Edition      string `json:"edition"`
		ID           int    `json:"id"`
	} `json:"movieFile"`

	// Collection
	Collection *struct {
		Title               string `json:"title"`
		TMDbID              int    `json:"tmdbId"`
		Monitored           bool   `json:"monitored"`
		QualityProfileID    int    `json:"qualityProfileId"`
		SearchOnAdd         bool   `json:"searchOnAdd"`
		MinimumAvailability string `json:"minimumAvailability"`
		Images              []struct {
			CoverType string `json:"coverType"`
			URL       string `json:"url"`
		} `json:"images"`
		Added time.Time `json:"added"`
		ID    int       `json:"id"`
	} `json:"collection"`
}

// radarrTag represents a tag from Radarr API.
type radarrTag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// radarrQueueItem represents an item in the download queue.
type radarrQueueItem struct {
	ID                      int       `json:"id"`
	MovieID                 int       `json:"movieId"`
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

// radarrQueueResponse represents the queue response.
type radarrQueueResponse struct {
	Page         int               `json:"page"`
	PageSize     int               `json:"pageSize"`
	TotalRecords int               `json:"totalRecords"`
	Records      []radarrQueueItem `json:"records"`
}

// RadarrQueueItem represents a processed queue item for public API.
type RadarrQueueItem struct {
	ID                  int       `json:"id"`
	MovieID             int       `json:"movie_id"`
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

// radarrHistoryItem represents a history entry.
type radarrHistoryItem struct {
	ID          int    `json:"id"`
	MovieID     int    `json:"movieId"`
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

// radarrHistoryResponse represents the history response.
type radarrHistoryResponse struct {
	Page         int                 `json:"page"`
	PageSize     int                 `json:"pageSize"`
	TotalRecords int                 `json:"totalRecords"`
	Records      []radarrHistoryItem `json:"records"`
}

// RadarrHistoryItem represents a processed history item for public API.
type RadarrHistoryItem struct {
	ID          int       `json:"id"`
	MovieID     int       `json:"movie_id"`
	SourceTitle string    `json:"source_title"`
	Quality     string    `json:"quality"`
	Date        time.Time `json:"date"`
	EventType   string    `json:"event_type"`
	DownloadID  string    `json:"download_id"`
}

// radarrCalendarItem represents a calendar entry.
type radarrCalendarItem struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	InCinemas       time.Time `json:"inCinemas"`
	PhysicalRelease time.Time `json:"physicalRelease"`
	DigitalRelease  time.Time `json:"digitalRelease"`
	Year            int       `json:"year"`
	HasFile         bool      `json:"hasFile"`
	Monitored       bool      `json:"monitored"`
}

// RadarrCalendarItem represents a processed calendar item for public API.
type RadarrCalendarItem struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	InCinemas       time.Time `json:"in_cinemas"`
	PhysicalRelease time.Time `json:"physical_release"`
	DigitalRelease  time.Time `json:"digital_release"`
	Year            int       `json:"year"`
	HasFile         bool      `json:"has_file"`
	Monitored       bool      `json:"monitored"`
}

// radarrQualityProfile represents a quality profile.
type radarrQualityProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RadarrQualityProfile represents a quality profile for public API.
type RadarrQualityProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// RadarrMovieMetadata contains complete Radarr movie metadata for public API.
// This struct exposes all available fields from the Radarr API for use in handlers and UI.
type RadarrMovieMetadata struct {
	// Basic Info
	ID            int    `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title,omitempty"`
	SortTitle     string `json:"sort_title,omitempty"`
	Year          int    `json:"year,omitempty"`

	// Status and Monitoring
	Status              string `json:"status,omitempty"` // tba, announced, inCinemas, released, deleted
	Overview            string `json:"overview,omitempty"`
	Monitored           bool   `json:"monitored"`
	MinimumAvailability string `json:"minimum_availability,omitempty"` // tba, announced, inCinemas, released, preDB
	IsAvailable         bool   `json:"is_available"`

	// Dates
	InCinemas       string `json:"in_cinemas,omitempty"`
	PhysicalRelease string `json:"physical_release,omitempty"`
	DigitalRelease  string `json:"digital_release,omitempty"`
	DateAdded       string `json:"date_added,omitempty"`

	// File Info
	Path       string `json:"path,omitempty"`
	FolderName string `json:"folder_name,omitempty"`
	SizeOnDisk int64  `json:"size_on_disk"`
	HasFile    bool   `json:"has_file"`

	// Technical Details
	Runtime int `json:"runtime,omitempty"` // minutes

	// External IDs and Links
	IMDbID           string `json:"imdb_id,omitempty"`
	TMDbID           int    `json:"tmdb_id,omitempty"`
	Website          string `json:"website,omitempty"`
	YouTubeTrailerID string `json:"youtube_trailer_id,omitempty"`

	// Media Info
	Studio        string   `json:"studio,omitempty"`
	Certification string   `json:"certification,omitempty"` // G, PG, PG-13, R, NC-17
	Genres        []string `json:"genres,omitempty"`
	Tags          []string `json:"tags,omitempty"` // Tag labels (not IDs)
	Popularity    float64  `json:"popularity,omitempty"`

	// Images
	PosterURL string   `json:"poster_url,omitempty"`
	FanartURL string   `json:"fanart_url,omitempty"`
	BannerURL string   `json:"banner_url,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`

	// Alternate Titles
	AlternateTitles []string `json:"alternate_titles,omitempty"`

	// Ratings
	IMDbRating          float64 `json:"imdb_rating,omitempty"`
	IMDbVotes           int     `json:"imdb_votes,omitempty"`
	TMDbRating          float64 `json:"tmdb_rating,omitempty"`
	TMDbVotes           int     `json:"tmdb_votes,omitempty"`
	MetacriticScore     int     `json:"metacritic_score,omitempty"`
	RottenTomatoesScore int     `json:"rotten_tomatoes_score,omitempty"`

	// Quality Info
	QualityName       string `json:"quality_name,omitempty"`
	QualitySource     string `json:"quality_source,omitempty"`     // bluray, webdl, webrip, hdtv
	QualityResolution int    `json:"quality_resolution,omitempty"` // 480, 720, 1080, 2160
	QualityModifier   string `json:"quality_modifier,omitempty"`   // remux, brdisk, regional, none
	QualityVersion    int    `json:"quality_version,omitempty"`
	IsRepack          bool   `json:"is_repack"`

	// Custom Formats
	CustomFormatScore int      `json:"custom_format_score,omitempty"`
	CustomFormats     []string `json:"custom_formats,omitempty"`

	// Movie File MediaInfo
	MediaInfo *struct {
		// Audio
		AudioCodec       string  `json:"audio_codec,omitempty"`
		AudioChannels    float64 `json:"audio_channels,omitempty"`
		AudioBitrate     int     `json:"audio_bitrate,omitempty"`
		AudioLanguages   string  `json:"audio_languages,omitempty"`
		AudioStreamCount int     `json:"audio_stream_count,omitempty"`

		// Video
		VideoCodec            string  `json:"video_codec,omitempty"`
		VideoBitrate          int     `json:"video_bitrate,omitempty"`
		VideoFps              float64 `json:"video_fps,omitempty"`
		VideoBitDepth         int     `json:"video_bit_depth,omitempty"`
		VideoDynamicRange     string  `json:"video_dynamic_range,omitempty"`      // SDR, HDR
		VideoDynamicRangeType string  `json:"video_dynamic_range_type,omitempty"` // HDR10, DolbyVision
		Resolution            string  `json:"resolution,omitempty"`               // 1920x1080, 3840x2160
		ScanType              string  `json:"scan_type,omitempty"`                // Progressive, Interlaced

		// Other
		RunTime   string `json:"run_time,omitempty"`
		Subtitles string `json:"subtitles,omitempty"`
	} `json:"media_info,omitempty"`

	// Additional File Info
	SceneName    string   `json:"scene_name,omitempty"`
	ReleaseGroup string   `json:"release_group,omitempty"`
	Edition      string   `json:"edition,omitempty"` // Extended, Director's Cut, etc.
	Languages    []string `json:"languages,omitempty"`

	// Collection
	CollectionTitle  string `json:"collection_title,omitempty"`
	CollectionTMDbID int    `json:"collection_tmdb_id,omitempty"`
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
	// Determine quality name - check movieFile first, fallback to movie level
	qualityName := ""
	if movie.MovieFile != nil {
		qualityName = movie.MovieFile.Quality.Quality.Name
	}

	media := &models.Media{
		Title:     movie.Title,
		Type:      "movie",
		FilePath:  movie.Path,
		Size:      movie.SizeOnDisk,
		AddedDate: movie.Added,
		RadarrID:  &movie.ID,
		Quality:   qualityName,
	}

	// Extract poster URL and other images
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

// ConvertToMetadata converts a Radarr movie to complete RadarrMovieMetadata.
// This function exposes ALL available fields from the Radarr API for use in handlers.
func (c *RadarrClient) ConvertToMetadata(movie *radarrMovie) *RadarrMovieMetadata {
	metadata := &RadarrMovieMetadata{
		ID:                  movie.ID,
		Title:               movie.Title,
		OriginalTitle:       movie.OriginalTitle,
		SortTitle:           movie.SortTitle,
		Year:                movie.Year,
		Status:              movie.Status,
		Overview:            movie.Overview,
		Monitored:           movie.Monitored,
		MinimumAvailability: movie.MinimumAvailability,
		IsAvailable:         movie.IsAvailable,
		Path:                movie.Path,
		FolderName:          movie.FolderName,
		SizeOnDisk:          movie.SizeOnDisk,
		HasFile:             movie.HasFile,
		Runtime:             movie.Runtime,
		IMDbID:              movie.IMDbID,
		TMDbID:              movie.TMDbID,
		Website:             movie.Website,
		YouTubeTrailerID:    movie.YouTubeTrailerID,
		Studio:              movie.Studio,
		Certification:       movie.Certification,
		Genres:              movie.Genres,
		Popularity:          movie.Popularity,
	}

	// Format dates
	if movie.InCinemas != nil && !movie.InCinemas.IsZero() {
		metadata.InCinemas = movie.InCinemas.Format("2006-01-02")
	}
	if movie.PhysicalRelease != nil && !movie.PhysicalRelease.IsZero() {
		metadata.PhysicalRelease = movie.PhysicalRelease.Format("2006-01-02")
	}
	if movie.DigitalRelease != nil && !movie.DigitalRelease.IsZero() {
		metadata.DigitalRelease = movie.DigitalRelease.Format("2006-01-02")
	}
	if !movie.Added.IsZero() {
		metadata.DateAdded = movie.Added.Format("2006-01-02")
	}

	// Extract image URLs
	var imageURLs []string
	for _, image := range movie.Images {
		switch image.CoverType {
		case "poster":
			metadata.PosterURL = image.URL
		case "fanart":
			metadata.FanartURL = image.URL
		case "banner":
			metadata.BannerURL = image.URL
		}
		if image.URL != "" {
			imageURLs = append(imageURLs, image.URL)
		}
	}
	metadata.ImageURLs = imageURLs

	// Extract alternate titles
	if len(movie.AlternateTitles) > 0 {
		altTitles := make([]string, 0, len(movie.AlternateTitles))
		for _, alt := range movie.AlternateTitles {
			if alt.Title != "" {
				altTitles = append(altTitles, alt.Title)
			}
		}
		metadata.AlternateTitles = altTitles
	}

	// Extract ratings
	if movie.Ratings.IMDb != nil {
		metadata.IMDbRating = movie.Ratings.IMDb.Value
		metadata.IMDbVotes = movie.Ratings.IMDb.Votes
	}
	if movie.Ratings.TMDb != nil {
		metadata.TMDbRating = movie.Ratings.TMDb.Value
		metadata.TMDbVotes = movie.Ratings.TMDb.Votes
	}
	if movie.Ratings.Metacritic != nil {
		metadata.MetacriticScore = movie.Ratings.Metacritic.Value
	}
	if movie.Ratings.RottenTomatoes != nil {
		metadata.RottenTomatoesScore = movie.Ratings.RottenTomatoes.Value
	}

	// Extract quality and file info
	if movie.MovieFile != nil {
		metadata.QualityName = movie.MovieFile.Quality.Quality.Name
		metadata.QualitySource = movie.MovieFile.Quality.Quality.Source
		metadata.QualityResolution = movie.MovieFile.Quality.Quality.Resolution
		metadata.QualityModifier = movie.MovieFile.Quality.Quality.Modifier
		metadata.QualityVersion = movie.MovieFile.Quality.Revision.Version
		metadata.IsRepack = movie.MovieFile.Quality.Revision.IsRepack

		metadata.CustomFormatScore = movie.MovieFile.CustomFormatScore
		if len(movie.MovieFile.CustomFormats) > 0 {
			cfNames := make([]string, 0, len(movie.MovieFile.CustomFormats))
			for _, cf := range movie.MovieFile.CustomFormats {
				cfNames = append(cfNames, cf.Name)
			}
			metadata.CustomFormats = cfNames
		}

		metadata.SceneName = movie.MovieFile.SceneName
		metadata.ReleaseGroup = movie.MovieFile.ReleaseGroup
		metadata.Edition = movie.MovieFile.Edition

		if len(movie.MovieFile.Languages) > 0 {
			langs := make([]string, 0, len(movie.MovieFile.Languages))
			for _, lang := range movie.MovieFile.Languages {
				langs = append(langs, lang.Name)
			}
			metadata.Languages = langs
		}

		// Extract MediaInfo
		if movie.MovieFile.MediaInfo != nil {
			metadata.MediaInfo = &struct {
				AudioCodec            string  `json:"audio_codec,omitempty"`
				AudioChannels         float64 `json:"audio_channels,omitempty"`
				AudioBitrate          int     `json:"audio_bitrate,omitempty"`
				AudioLanguages        string  `json:"audio_languages,omitempty"`
				AudioStreamCount      int     `json:"audio_stream_count,omitempty"`
				VideoCodec            string  `json:"video_codec,omitempty"`
				VideoBitrate          int     `json:"video_bitrate,omitempty"`
				VideoFps              float64 `json:"video_fps,omitempty"`
				VideoBitDepth         int     `json:"video_bit_depth,omitempty"`
				VideoDynamicRange     string  `json:"video_dynamic_range,omitempty"`
				VideoDynamicRangeType string  `json:"video_dynamic_range_type,omitempty"`
				Resolution            string  `json:"resolution,omitempty"`
				ScanType              string  `json:"scan_type,omitempty"`
				RunTime               string  `json:"run_time,omitempty"`
				Subtitles             string  `json:"subtitles,omitempty"`
			}{
				AudioCodec:            movie.MovieFile.MediaInfo.AudioCodec,
				AudioChannels:         movie.MovieFile.MediaInfo.AudioChannels,
				AudioBitrate:          movie.MovieFile.MediaInfo.AudioBitrate,
				AudioLanguages:        movie.MovieFile.MediaInfo.AudioLanguages,
				AudioStreamCount:      movie.MovieFile.MediaInfo.AudioStreamCount,
				VideoCodec:            movie.MovieFile.MediaInfo.VideoCodec,
				VideoBitrate:          movie.MovieFile.MediaInfo.VideoBitrate,
				VideoFps:              movie.MovieFile.MediaInfo.VideoFps,
				VideoBitDepth:         movie.MovieFile.MediaInfo.VideoBitDepth,
				VideoDynamicRange:     movie.MovieFile.MediaInfo.VideoDynamicRange,
				VideoDynamicRangeType: movie.MovieFile.MediaInfo.VideoDynamicRangeType,
				Resolution:            movie.MovieFile.MediaInfo.Resolution,
				ScanType:              movie.MovieFile.MediaInfo.ScanType,
				RunTime:               movie.MovieFile.MediaInfo.RunTime,
				Subtitles:             movie.MovieFile.MediaInfo.Subtitles,
			}
		}
	}

	// Extract collection info
	if movie.Collection != nil {
		metadata.CollectionTitle = movie.Collection.Title
		metadata.CollectionTMDbID = movie.Collection.TMDbID
	}

	// Note: Tags are stored as IDs in the movie struct
	// To get tag labels, we'd need to fetch them separately via GetTags()
	// For now, we'll store them as formatted strings
	if len(movie.Tags) > 0 {
		tagStrings := make([]string, 0, len(movie.Tags))
		for _, tagID := range movie.Tags {
			tagStrings = append(tagStrings, fmt.Sprintf("tag_%d", tagID))
		}
		metadata.Tags = tagStrings
	}

	return metadata
}

// GetQueue retrieves the current download queue from Radarr.
func (c *RadarrClient) GetQueue(ctx context.Context) ([]RadarrQueueItem, error) {
	var queueResp radarrQueueResponse

	err := c.callWithRetry(ctx, func() error {
		resp, err := c.client.R().
			SetContext(ctx).
			SetResult(&queueResp).
			SetQueryParam("pageSize", "100").
			SetQueryParam("includeUnknownMovieItems", "false").
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
	items := make([]RadarrQueueItem, 0, len(queueResp.Records))
	for _, record := range queueResp.Records {
		progress := 0.0
		if record.Size > 0 {
			progress = float64(record.Size-record.Sizeleft) / float64(record.Size) * 100
			// Ensure progress is within 0-100 range (handle inconsistent data from Radarr)
			if progress < 0 {
				progress = 0
			} else if progress > 100 {
				progress = 100
			}
		}

		items = append(items, RadarrQueueItem{
			ID:                  record.ID,
			MovieID:             record.MovieID,
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

	c.logger.Info("Retrieved Radarr queue",
		zap.Int("total_items", len(items)),
	)

	return items, nil
}

// GetHistory retrieves recent history from Radarr.
func (c *RadarrClient) GetHistory(ctx context.Context, pageSize int) ([]RadarrHistoryItem, error) {
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var historyResp radarrHistoryResponse

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
	items := make([]RadarrHistoryItem, 0, len(historyResp.Records))
	for _, record := range historyResp.Records {
		items = append(items, RadarrHistoryItem{
			ID:          record.ID,
			MovieID:     record.MovieID,
			SourceTitle: record.SourceTitle,
			Quality:     record.Quality.Quality.Name,
			Date:        record.Date,
			EventType:   record.EventType,
			DownloadID:  record.DownloadID,
		})
	}

	c.logger.Info("Retrieved Radarr history",
		zap.Int("total_items", len(items)),
	)

	return items, nil
}

// GetCalendar retrieves upcoming movies from Radarr.
func (c *RadarrClient) GetCalendar(ctx context.Context, startDate, endDate time.Time) ([]RadarrCalendarItem, error) {
	var calendarItems []radarrCalendarItem

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
	items := make([]RadarrCalendarItem, 0, len(calendarItems))
	for _, record := range calendarItems {
		items = append(items, RadarrCalendarItem{
			ID:              record.ID,
			Title:           record.Title,
			InCinemas:       record.InCinemas,
			PhysicalRelease: record.PhysicalRelease,
			DigitalRelease:  record.DigitalRelease,
			Year:            record.Year,
			HasFile:         record.HasFile,
			Monitored:       record.Monitored,
		})
	}

	c.logger.Info("Retrieved Radarr calendar",
		zap.Int("total_items", len(items)),
		zap.String("start_date", startDate.Format("2006-01-02")),
		zap.String("end_date", endDate.Format("2006-01-02")),
	)

	return items, nil
}

// GetQualityProfiles retrieves available quality profiles from Radarr.
func (c *RadarrClient) GetQualityProfiles(ctx context.Context) ([]RadarrQualityProfile, error) {
	var profiles []radarrQualityProfile

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
	items := make([]RadarrQualityProfile, 0, len(profiles))
	for _, profile := range profiles {
		items = append(items, RadarrQualityProfile{
			ID:   profile.ID,
			Name: profile.Name,
		})
	}

	c.logger.Info("Retrieved Radarr quality profiles",
		zap.Int("total_profiles", len(items)),
	)

	return items, nil
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
