package handler

import (
	"fmt"
	"path/filepath"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// MoviesHandler handles movie-related requests
type MoviesHandler struct {
	cfg            *config.Config
	jellyfinClient *clients.JellyfinClient
	logger         *zap.Logger
}

// NewMoviesHandler creates a new movies handler
func NewMoviesHandler(
	cfg *config.Config,
	jellyfinClient *clients.JellyfinClient,
	logger *zap.Logger,
) *MoviesHandler {
	return &MoviesHandler{
		cfg:            cfg,
		jellyfinClient: jellyfinClient,
		logger:         logger,
	}
}

// MovieFile represents a movie file with its metadata
type MovieFile struct {
	Path      string                 `json:"path"`
	Title     string                 `json:"title"`
	Year      int                    `json:"year,omitempty"`
	Size      int64                  `json:"size"`
	Extension string                 `json:"extension"`
	Directory string                 `json:"directory"`
	PosterURL string                 `json:"poster_url,omitempty"`
	Jellyfin  *JellyfinMovieMetadata `json:"jellyfin,omitempty"`
}

// JellyfinMovieMetadata contains Jellyfin-specific metadata
type JellyfinMovieMetadata struct {
	ID              string            `json:"id"`
	Overview        string            `json:"overview,omitempty"`
	CommunityRating float64           `json:"community_rating,omitempty"`
	CriticRating    float64           `json:"critic_rating,omitempty"`
	OfficialRating  string            `json:"official_rating,omitempty"`
	PlayCount       int               `json:"play_count"`
	LastPlayed      string            `json:"last_played,omitempty"`
	DateAdded       string            `json:"date_added,omitempty"`
	ProductionYear  int               `json:"production_year,omitempty"`
	PremiereDate    string            `json:"premiere_date,omitempty"`
	Runtime         int               `json:"runtime,omitempty"` // Duración en minutos
	Genres          []string          `json:"genres,omitempty"`
	Studios         []string          `json:"studios,omitempty"`
	Actors          []string          `json:"actors,omitempty"`
	Directors       []string          `json:"directors,omitempty"`
	Writers         []string          `json:"writers,omitempty"`
	IMDbID          string            `json:"imdb_id,omitempty"`
	TMDbID          string            `json:"tmdb_id,omitempty"`
	Taglines        []string          `json:"taglines,omitempty"`
	VideoStreams    []VideoStreamInfo `json:"video_streams,omitempty"`
	AudioStreams    []AudioStreamInfo `json:"audio_streams,omitempty"`
	SubtitleStreams []SubtitleInfo    `json:"subtitle_streams,omitempty"`
}

// VideoStreamInfo contains video stream information
type VideoStreamInfo struct {
	Codec      string `json:"codec"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	VideoRange string `json:"video_range,omitempty"` // SDR, HDR, etc.
	IsDefault  bool   `json:"is_default"`
}

// AudioStreamInfo contains audio stream information
type AudioStreamInfo struct {
	Codec         string `json:"codec"`
	Language      string `json:"language"`
	DisplayTitle  string `json:"display_title"`
	Channels      int    `json:"channels,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
	IsDefault     bool   `json:"is_default"`
}

// SubtitleInfo contains subtitle information
type SubtitleInfo struct {
	Language        string `json:"language"`
	DisplayLanguage string `json:"display_language"`
	Title           string `json:"title,omitempty"`
	IsDefault       bool   `json:"is_default"`
	IsForced        bool   `json:"is_forced"`
}

// RenderMoviesPage renders the movies page
func (h *MoviesHandler) RenderMoviesPage(c *fiber.Ctx) error {
	return c.Render("pages/movies", fiber.Map{
		"Title": "Películas",
	})
}

// GetMovies retrieves all movies from Jellyfin library paths and enriches with Jellyfin metadata
func (h *MoviesHandler) GetMovies(c *fiber.Ctx) error {
	ctx := c.Context()

	// Check if Jellyfin is enabled
	if !h.cfg.Clients.Jellyfin.Enabled {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Jellyfin no está configurado",
		})
	}

	h.logger.Info("Getting movies from Jellyfin")

	// Get movies directly from Jellyfin library with full metadata
	jellyfinItems, err := h.jellyfinClient.GetLibraryItems(ctx)
	if err != nil {
		h.logger.Error("Failed to get library from Jellyfin", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener biblioteca de Jellyfin",
		})
	}

	// Get virtual folders for library paths info
	virtualFolders, err := h.jellyfinClient.GetVirtualFolders(ctx)
	if err != nil {
		h.logger.Warn("Failed to get virtual folders", zap.Error(err))
	}

	var libraryPaths []string
	for _, folder := range virtualFolders {
		libraryPaths = append(libraryPaths, folder.Locations...)
	}

	// Convert Jellyfin items to MovieFile format
	movies := h.convertJellyfinItemsToMovies(jellyfinItems)

	h.logger.Info("Movies retrieval complete",
		zap.Int("movie_count", len(movies)),
	)

	return c.JSON(fiber.Map{
		"movies":        movies,
		"count":         len(movies),
		"library_paths": libraryPaths,
	})
}

// convertJellyfinItemsToMovies converts raw Jellyfin items to MovieFile format with full metadata
func (h *MoviesHandler) convertJellyfinItemsToMovies(jellyfinItems []clients.JellyfinItem) []MovieFile {
	var movies []MovieFile
	var seriesCount, episodeCount, seasonCount, otherCount int

	for _, item := range jellyfinItems {
		// Only process movies (not series, episodes, or seasons)
		if item.Type != "Movie" {
			h.logger.Debug("Filtering out non-movie item",
				zap.String("title", item.Name),
				zap.String("type", item.Type),
			)

			if item.Type == "Series" {
				seriesCount++
			} else if item.Type == "Episode" {
				episodeCount++
			} else if item.Type == "Season" {
				seasonCount++
			} else {
				otherCount++
			}
			continue
		}

		// Extract year from filename or use default
		year := h.extractYear(item.Name)

		// Get file size from MediaSources
		var size int64
		if len(item.MediaSources) > 0 {
			size = item.MediaSources[0].Size
		}

		// Build poster URL if available
		var posterURL string
		if primaryTag, ok := item.ImageTags["Primary"]; ok {
			posterURL = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s",
				h.jellyfinClient.GetBaseURL(), item.ID, primaryTag)
		}

		movie := MovieFile{
			Path:      item.Path,
			Title:     item.Name,
			Year:      year,
			Size:      size,
			Extension: filepath.Ext(item.Path),
			Directory: filepath.Dir(item.Path),
			PosterURL: posterURL,
		}

		// Add COMPLETE Jellyfin metadata from the item
		movie.Jellyfin = &JellyfinMovieMetadata{
			ID:              item.ID,
			Overview:        item.Overview,
			CommunityRating: item.CommunityRating,
			CriticRating:    item.CriticRating,
			OfficialRating:  item.OfficialRating,
			PlayCount:       item.UserData.PlayCount,
			ProductionYear:  item.ProductionYear,
			Genres:          item.Genres,
			Taglines:        item.Taglines,
		}

		// Runtime en minutos (convertir de ticks)
		if item.RunTimeTicks > 0 {
			movie.Jellyfin.Runtime = int(item.RunTimeTicks / 600000000) // 1 minuto = 600,000,000 ticks
		}

		// Studios
		for _, studio := range item.Studios {
			movie.Jellyfin.Studios = append(movie.Jellyfin.Studios, studio.Name)
		}

		// Cast y Crew
		var actors, directors, writers []string
		for _, person := range item.People {
			switch person.Type {
			case "Actor":
				actors = append(actors, person.Name)
			case "Director":
				directors = append(directors, person.Name)
			case "Writer":
				writers = append(writers, person.Name)
			}
		}
		movie.Jellyfin.Actors = actors
		movie.Jellyfin.Directors = directors
		movie.Jellyfin.Writers = writers

		// Provider IDs (IMDb, TMDb, etc.)
		if imdbID, ok := item.ProviderIds["Imdb"]; ok {
			movie.Jellyfin.IMDbID = imdbID
		}
		if tmdbID, ok := item.ProviderIds["Tmdb"]; ok {
			movie.Jellyfin.TMDbID = tmdbID
		}

		// Media Streams
		for _, stream := range item.MediaStreams {
			switch stream.Type {
			case "Video":
				movie.Jellyfin.VideoStreams = append(movie.Jellyfin.VideoStreams, VideoStreamInfo{
					Codec:      stream.Codec,
					Width:      stream.Width,
					Height:     stream.Height,
					VideoRange: stream.VideoRange,
					IsDefault:  stream.IsDefault,
				})
			case "Audio":
				movie.Jellyfin.AudioStreams = append(movie.Jellyfin.AudioStreams, AudioStreamInfo{
					Codec:         stream.Codec,
					Language:      stream.Language,
					DisplayTitle:  stream.DisplayTitle,
					Channels:      stream.Channels,
					ChannelLayout: stream.ChannelLayout,
					IsDefault:     stream.IsDefault,
				})
			case "Subtitle":
				movie.Jellyfin.SubtitleStreams = append(movie.Jellyfin.SubtitleStreams, SubtitleInfo{
					Language:        stream.Language,
					DisplayLanguage: stream.DisplayLanguage,
					Title:           stream.Title,
					IsDefault:       stream.IsDefault,
					IsForced:        stream.IsForced,
				})
			}
		}

		// Add date information
		if !item.DateCreated.IsZero() {
			movie.Jellyfin.DateAdded = item.DateCreated.Format("2006-01-02")
		}

		if !item.PremiereDate.IsZero() {
			movie.Jellyfin.PremiereDate = item.PremiereDate.Format("2006-01-02")
		}

		if !item.UserData.LastPlayedDate.IsZero() {
			movie.Jellyfin.LastPlayed = item.UserData.LastPlayedDate.Format("2006-01-02 15:04")
		}

		movies = append(movies, movie)
	}

	h.logger.Info("Converted Jellyfin items to movies",
		zap.Int("total_items", len(jellyfinItems)),
		zap.Int("movies", len(movies)),
		zap.Int("series_filtered", seriesCount),
		zap.Int("episodes_filtered", episodeCount),
		zap.Int("seasons_filtered", seasonCount),
		zap.Int("other_filtered", otherCount),
	)

	return movies
}

// extractYear tries to extract year from movie title
func (h *MoviesHandler) extractYear(title string) int {
	// Look for (YYYY) pattern
	for i := 0; i < len(title)-5; i++ {
		if title[i] == '(' && i+5 < len(title) && title[i+5] == ')' {
			yearStr := title[i+1 : i+5]
			// Check if it's a valid year (1900-2099)
			if len(yearStr) == 4 {
				var year int
				_, err := fmt.Sscanf(yearStr, "%d", &year)
				if err == nil && year >= 1900 && year <= 2099 {
					return year
				}
			}
		}
	}
	return 0
}

// formatBytes formats bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
