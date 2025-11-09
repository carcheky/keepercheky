package handler

import (
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// SeasonInfo represents a season with its episodes
type SeasonInfo struct {
	SeasonNumber int           `json:"season_number"`
	EpisodeCount int           `json:"episode_count"`
	TotalSize    int64         `json:"total_size"`
	Episodes     []EpisodeInfo `json:"episodes"`
}

// EpisodeInfo represents an individual episode
type EpisodeInfo struct {
	EpisodeNumber int      `json:"episode_number"`
	Title         string   `json:"title"`
	FilePath      string   `json:"file_path"`
	Size          int64    `json:"size"`
	Quality       string   `json:"quality"`
	Versions      []string `json:"versions,omitempty"` // Otras versiones del mismo episodio
}

// SeriesInfo represents a TV series with seasons
type SeriesInfo struct {
	SeriesTitle  string          `json:"series_title"`
	TotalSize    int64           `json:"total_size"`
	SeasonCount  int             `json:"season_count"`
	EpisodeCount int             `json:"episode_count"`
	Seasons      []SeasonInfo    `json:"seasons"`
	PosterURL    string          `json:"poster_url"`
	PrimaryPath  string          `json:"primary_path"` // Ruta principal (en Jellyfin)
	Files        []MediaFileInfo `json:"files"`        // Todos los archivos relacionados
	Metadata     SeriesMetadata  `json:"metadata"`
}

// SeriesMetadata contains metadata about the series
type SeriesMetadata struct {
	InJellyfin    bool    `json:"in_jellyfin"`
	InSonarr      bool    `json:"in_sonarr"`
	InQBittorrent bool    `json:"in_qbittorrent"`
	SonarrID      *int    `json:"sonarr_id,omitempty"`
	JellyfinID    *string `json:"jellyfin_id,omitempty"`
}

// MovieInfo represents a movie with possible multiple versions
type MovieInfo struct {
	Title         string          `json:"title"`
	TotalSize     int64           `json:"total_size"`
	PrimaryFile   MediaFileInfo   `json:"primary_file"`             // Versión principal (en Jellyfin)
	OtherVersions []MediaFileInfo `json:"other_versions,omitempty"` // Otras versiones
	Metadata      MovieMetadata   `json:"metadata"`
}

// MovieMetadata contains metadata about the movie
type MovieMetadata struct {
	InJellyfin    bool    `json:"in_jellyfin"`
	InRadarr      bool    `json:"in_radarr"`
	InQBittorrent bool    `json:"in_qbittorrent"`
	RadarrID      *int    `json:"radarr_id,omitempty"`
	JellyfinID    *string `json:"jellyfin_id,omitempty"`
}

// OrganizedFilesResponse represents the organized response
type OrganizedFilesResponse struct {
	Series     []SeriesInfo `json:"series"`
	Movies     []MovieInfo  `json:"movies"`
	TotalCount int          `json:"total_count"`
}

// Regular expressions for parsing file names (compiled once for performance)
var (
	// Matches patterns like "S01E01", "s01e01"
	episodeRegex = regexp.MustCompile(`(?i)[Ss](\d+)[Ee](\d+)`)
	// Alternative pattern: 1x01, 1X01
	altEpisodeRegex = regexp.MustCompile(`(?i)(\d+)[xX](\d+)`)
	// Quality and release tags to remove
	qualityRegex = regexp.MustCompile(`(?i)(1080p|720p|2160p|4k|BluRay|WEB-DL|HDTV|x264|x265|HEVC).*`)
	// Multiple spaces cleanup
	spaceRegex = regexp.MustCompile(`\s+`)
)

// parseSeasonEpisode extracts season and episode numbers from filename
func parseSeasonEpisode(filename string) (season int, episode int, found bool) {
	// Try standard SxxExx format first
	matches := episodeRegex.FindStringSubmatch(filename)
	if len(matches) == 3 {
		season, _ = strconv.Atoi(matches[1])
		episode, _ = strconv.Atoi(matches[2])
		return season, episode, true
	}

	// Try alternative XxY format
	matches = altEpisodeRegex.FindStringSubmatch(filename)
	if len(matches) == 3 {
		season, _ = strconv.Atoi(matches[1])
		episode, _ = strconv.Atoi(matches[2])
		return season, episode, true
	}

	return 0, 0, false
}

// extractSeriesName extracts series name from file path
func extractSeriesName(path string) string {
	// Get the base name
	base := filepath.Base(path)

	// Remove extension
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	// Find episode pattern and truncate before it
	episodeMatch := episodeRegex.FindStringIndex(nameWithoutExt)
	if episodeMatch != nil {
		nameWithoutExt = nameWithoutExt[:episodeMatch[0]]
	} else {
		altEpisodeMatch := altEpisodeRegex.FindStringIndex(nameWithoutExt)
		if altEpisodeMatch != nil {
			nameWithoutExt = nameWithoutExt[:altEpisodeMatch[0]]
		}
	}

	// Clean up
	name := strings.TrimSpace(nameWithoutExt)
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Remove common quality/release tags (case insensitive)
	name = qualityRegex.ReplaceAllString(name, "")

	// Clean up multiple spaces
	name = spaceRegex.ReplaceAllString(name, " ")

	return strings.TrimSpace(name)
}

// GetOrganizedFilesAPI returns files organized by series/seasons and movies
func (h *FilesHandler) GetOrganizedFilesAPI(c *fiber.Ctx) error {
	// Get pagination parameters
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("perPage", 25)
	tab := c.Query("tab", "")
	typeFilter := c.Query("type", "")       // "movie", "series"
	search := c.Query("search", "")         // Search query
	serviceFilter := c.Query("service", "") // Filter by service

	// Validate parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 25
	}
	if perPage > 100 {
		perPage = 100
	}

	// Build query for all media files
	query := h.mediaRepo.GetDB().
		Table("media").
		Select(`
			id, title, type, file_path, size, poster_url, quality,
			is_hardlink, hardlink_paths, primary_path,
			in_radarr, in_sonarr, in_jellyfin, in_jellyseerr, in_jellystat, in_q_bittorrent,
			radarr_id, sonarr_id, jellyfin_id, jellyseerr_id, jellystat_id,
			torrent_hash, torrent_category, torrent_state, torrent_tags,
			is_seeding, seed_ratio, excluded
		`)

	// Apply search filter (minimum 3 characters for performance)
	if search != "" && len(strings.TrimSpace(search)) >= 3 {
		searchPattern := "%" + search + "%"

		// Use case-insensitive search
		dialectName := h.mediaRepo.GetDB().Dialector.Name()
		if dialectName == "postgres" {
			query = query.Where(
				"title ILIKE ? OR file_path ILIKE ? OR torrent_hash ILIKE ?",
				searchPattern, searchPattern, searchPattern,
			)
		} else {
			query = query.Where(
				"title LIKE ? OR file_path LIKE ? OR torrent_hash LIKE ?",
				searchPattern, searchPattern, searchPattern,
			)
		}
	}

	// Apply type filter
	if typeFilter != "" {
		switch typeFilter {
		case "movie":
			query = query.Where("type = ?", "movie")
		case "series":
			// Include both series and episode types for series filter
			query = query.Where("type IN (?)", []string{"series", "episode"})
		}
	}

	// Apply service filter
	if serviceFilter != "" {
		switch serviceFilter {
		case "qbittorrent":
			query = query.Where("in_q_bittorrent = ?", true)
		case "radarr":
			query = query.Where("in_radarr = ?", true)
		case "sonarr":
			query = query.Where("in_sonarr = ?", true)
		case "jellyfin":
			query = query.Where("in_jellyfin = ?", true)
		case "orphan":
			query = query.Where("in_q_bittorrent = ? AND in_radarr = ? AND in_sonarr = ? AND in_jellyfin = ?",
				true, false, false, false)
		}
	}

	// Apply tab filtering (same as regular files API)
	switch tab {
	case "attention":
		query = query.Where("in_q_bittorrent = ? AND in_jellyfin = ? AND in_radarr = ? AND in_sonarr = ?",
			true, false, false, false)
	case "critical":
		query = query.Where("in_q_bittorrent = ? AND (torrent_state = ? OR torrent_state = ?)",
			true, "error", "missingFiles")
	case "hardlinks":
		query = query.Where("is_hardlink = ?", true)
	case "unwatched":
		query = query.Where("in_jellyfin = ?", true)
	case "healthy":
		query = query.Where("in_jellyfin = ? AND (in_radarr = ? OR in_sonarr = ?)",
			true, true, true).
			Where("(torrent_state IS NULL OR torrent_state NOT IN (?, ?))", "error", "missingFiles")
	}

	// Get all matching files (we'll organize them in memory)
	var allFiles []MediaFileInfo
	err := query.Find(&allFiles).Error
	if err != nil {
		h.logger.Error("Failed to query media", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query media"})
	}

	// Organize files
	organized := h.organizeFiles(allFiles)

	// Apply pagination to organized results
	totalSeries := len(organized.Series)
	totalMovies := len(organized.Movies)
	totalCount := totalSeries + totalMovies

	// For simplicity, paginate across both series and movies
	startIdx := (page - 1) * perPage
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + perPage

	response := OrganizedFilesResponse{
		Series:     []SeriesInfo{},
		Movies:     []MovieInfo{},
		TotalCount: totalCount,
	}

	// Slice series if needed
	if startIdx < totalSeries {
		seriesEnd := min(endIdx, totalSeries)
		response.Series = organized.Series[startIdx:seriesEnd]

		// If there's room for movies
		if endIdx > totalSeries {
			moviesStart := 0
			moviesEnd := min(endIdx-totalSeries, totalMovies)
			response.Movies = organized.Movies[moviesStart:moviesEnd]
		}
	} else if startIdx < totalCount {
		// Only movies
		moviesStart := startIdx - totalSeries
		moviesEnd := min(endIdx-totalSeries, totalMovies)
		response.Movies = organized.Movies[moviesStart:moviesEnd]
	}

	h.logger.Info("Organized files API request",
		zap.Int("total_files", len(allFiles)),
		zap.Int("total_series", totalSeries),
		zap.Int("total_movies", totalMovies),
		zap.Int("page", page),
		zap.Int("perPage", perPage),
		zap.Int("returned_series", len(response.Series)),
		zap.Int("returned_movies", len(response.Movies)),
		zap.String("tab", tab),
		zap.String("type_filter", typeFilter),
	)

	return c.JSON(fiber.Map{
		"data":       response,
		"page":       page,
		"perPage":    perPage,
		"totalCount": totalCount,
	})
}

// organizeFiles groups files into series and movies with hierarchical structure
func (h *FilesHandler) organizeFiles(files []MediaFileInfo) OrganizedFilesResponse {
	seriesMap := make(map[string]*SeriesInfo)
	movieMap := make(map[string]*MovieInfo)

	// Count how many files we process
	seriesCount := 0
	movieCount := 0

	// Group files
	for _, file := range files {
		isSeries := false

		// Method 1: Check explicit file type first
		if file.Type == "series" || file.Type == "episode" {
			isSeries = true
		} else if file.Type == "movie" {
			isSeries = false
		} else {
			// Method 2: If type is not set, use intelligent heuristics

			// Check if filename contains episode patterns (S01E01, 1x01, etc.)
			_, _, hasEpisodePattern := parseSeasonEpisode(file.FilePath)
			if hasEpisodePattern {
				isSeries = true
			} else if file.InSonarr {
				// If it's in Sonarr, it's definitely a series
				isSeries = true
			} else if file.InRadarr {
				// If it's in Radarr, it's definitely a movie
				isSeries = false
			} else if strings.Contains(strings.ToLower(file.FilePath), "/tv/") ||
				strings.Contains(strings.ToLower(file.FilePath), "/series/") ||
				strings.Contains(strings.ToLower(file.FilePath), "/shows/") {
				// Path-based heuristic: common TV show directory names
				isSeries = true
			} else if strings.Contains(strings.ToLower(file.FilePath), "/movies/") ||
				strings.Contains(strings.ToLower(file.FilePath), "/films/") {
				// Path-based heuristic: common movie directory names
				isSeries = false
			} else {
				// Default: if no indicators, treat as movie
				// (movies are more common as standalone files)
				isSeries = false
			}
		}

		if isSeries {
			h.addToSeries(seriesMap, file)
			seriesCount++
		} else {
			h.addToMovies(movieMap, file)
			movieCount++
		}
	}

	h.logger.Info("Organized files by type",
		zap.Int("total_files", len(files)),
		zap.Int("files_to_series", seriesCount),
		zap.Int("files_to_movies", movieCount),
		zap.Int("unique_series", len(seriesMap)),
		zap.Int("unique_movies", len(movieMap)),
	)

	// Convert maps to slices and sort
	series := make([]SeriesInfo, 0, len(seriesMap))
	for _, s := range seriesMap {
		// Sort seasons
		sort.Slice(s.Seasons, func(i, j int) bool {
			return s.Seasons[i].SeasonNumber < s.Seasons[j].SeasonNumber
		})

		// Sort episodes within each season
		for i := range s.Seasons {
			sort.Slice(s.Seasons[i].Episodes, func(a, b int) bool {
				return s.Seasons[i].Episodes[a].EpisodeNumber < s.Seasons[i].Episodes[b].EpisodeNumber
			})
		}

		series = append(series, *s)
	}
	sort.Slice(series, func(i, j int) bool {
		return series[i].SeriesTitle < series[j].SeriesTitle
	})

	movies := make([]MovieInfo, 0, len(movieMap))
	for _, m := range movieMap {
		movies = append(movies, *m)
	}
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Title < movies[j].Title
	})

	return OrganizedFilesResponse{
		Series: series,
		Movies: movies,
	}
}

// addToSeries adds a file to the series map
func (h *FilesHandler) addToSeries(seriesMap map[string]*SeriesInfo, file MediaFileInfo) {
	// Extract series name and season/episode info
	seriesName := file.Title
	if seriesName == "" {
		seriesName = extractSeriesName(file.FilePath)
	}

	season, episode, found := parseSeasonEpisode(file.FilePath)
	if !found {
		// If we can't parse season/episode, treat as a series-level file
		season = 0
		episode = 0
	}

	// Get or create series
	if seriesMap[seriesName] == nil {
		seriesMap[seriesName] = &SeriesInfo{
			SeriesTitle: seriesName,
			Seasons:     []SeasonInfo{},
			Files:       []MediaFileInfo{},
			PosterURL:   file.PosterURL,
			Metadata: SeriesMetadata{
				InJellyfin:    file.InJellyfin,
				InSonarr:      file.InSonarr,
				InQBittorrent: file.InQBittorrent,
				SonarrID:      file.SonarrID,
				JellyfinID:    file.JellyfinID,
			},
		}
		// Prioritize Jellyfin path as primary
		if file.InJellyfin {
			seriesMap[seriesName].PrimaryPath = file.FilePath
		}
	}

	series := seriesMap[seriesName]
	series.Files = append(series.Files, file)
	series.TotalSize += file.Size

	// Update primary path if this file is in Jellyfin and we don't have one yet
	if file.InJellyfin && series.PrimaryPath == "" {
		series.PrimaryPath = file.FilePath
	}

	// Add to appropriate season
	var targetSeason *SeasonInfo
	for i := range series.Seasons {
		if series.Seasons[i].SeasonNumber == season {
			targetSeason = &series.Seasons[i]
			break
		}
	}

	if targetSeason == nil {
		series.Seasons = append(series.Seasons, SeasonInfo{
			SeasonNumber: season,
			Episodes:     []EpisodeInfo{},
		})
		targetSeason = &series.Seasons[len(series.Seasons)-1]
		series.SeasonCount++
	}

	// Add episode
	episodeInfo := EpisodeInfo{
		EpisodeNumber: episode,
		Title:         file.Title,
		FilePath:      file.FilePath,
		Size:          file.Size,
		Quality:       file.Quality,
	}

	targetSeason.Episodes = append(targetSeason.Episodes, episodeInfo)
	targetSeason.EpisodeCount++
	targetSeason.TotalSize += file.Size
	series.EpisodeCount++
}

// addToMovies adds a file to the movies map
func (h *FilesHandler) addToMovies(movieMap map[string]*MovieInfo, file MediaFileInfo) {
	movieTitle := file.Title
	if movieTitle == "" {
		movieTitle = h.inferTitleFromPath(file.FilePath)
	}

	// Get or create movie
	if movieMap[movieTitle] == nil {
		primaryFile := file
		otherVersions := []MediaFileInfo{}

		movieMap[movieTitle] = &MovieInfo{
			Title:         movieTitle,
			TotalSize:     file.Size,
			PrimaryFile:   primaryFile,
			OtherVersions: otherVersions,
			Metadata: MovieMetadata{
				InJellyfin:    file.InJellyfin,
				InRadarr:      file.InRadarr,
				InQBittorrent: file.InQBittorrent,
				RadarrID:      file.RadarrID,
				JellyfinID:    file.JellyfinID,
			},
		}
	} else {
		// This is another version of the same movie
		movie := movieMap[movieTitle]
		movie.TotalSize += file.Size

		// Prioritize Jellyfin version as primary
		if file.InJellyfin && !movie.PrimaryFile.InJellyfin {
			// Swap: current primary becomes other version
			movie.OtherVersions = append(movie.OtherVersions, movie.PrimaryFile)
			movie.PrimaryFile = file
		} else {
			movie.OtherVersions = append(movie.OtherVersions, file)
		}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
