package filesystem

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/carcheky/keepercheky/internal/models"
	"go.uber.org/zap"
)

// EnrichedFile represents a file with service metadata
type EnrichedFile struct {
	*FileEntry

	// Service flags
	InRadarr      bool
	InSonarr      bool
	InJellyfin    bool
	InJellyseerr  bool
	InJellystat   bool
	InQBittorrent bool

	// Service IDs
	RadarrID     *int
	SonarrID     *int
	JellyfinID   *string
	JellyseerrID *int

	// Torrent info
	TorrentHash     string
	TorrentCategory string
	TorrentTags     string
	TorrentState    string
	IsSeeding       bool
	SeedRatio       float64

	// Additional metadata
	Title       string
	Quality     string
	PosterURL   string
	Tags        []string
	Excluded    bool
	LastWatched *int64
	ModTime     int64 // Unix timestamp from FileEntry for easy access
}

// Enricher enriches file entries with service metadata
type Enricher struct {
	logger *zap.Logger
}

// NewEnricher creates a new enricher
func NewEnricher(logger *zap.Logger) *Enricher {
	return &Enricher{
		logger: logger,
	}
}

// EnrichWithRadarr enriches files with Radarr metadata
func (e *Enricher) EnrichWithRadarr(
	ctx context.Context,
	files map[string]*EnrichedFile,
	radarrMedia []*models.Media,
) int {
	e.logger.Info("Enriching with Radarr data",
		zap.Int("radarr_items", len(radarrMedia)),
	)

	enriched := 0

	// Create path map for fast lookup
	radarrByPath := make(map[string]*models.Media)
	for _, media := range radarrMedia {
		radarrByPath[media.FilePath] = media
	}

	// Match files with Radarr items
	for path, file := range files {
		// Try exact match first
		if radarrItem, found := radarrByPath[path]; found {
			e.applyRadarrData(file, radarrItem)
			enriched++
			continue
		}

		// Try matching with all hardlink paths (includes primary path)
		if file.IsHardlink {
			if radarrItem, matched := matchByHardlinks(radarrByPath, file.HardlinkPaths); matched {
				e.applyRadarrData(file, radarrItem)
				enriched++
				continue
			}
		}

		// Try fuzzy matching by directory name
		for radarrPath, radarrItem := range radarrByPath {
			if e.pathsMatch(path, radarrPath) {
				e.applyRadarrData(file, radarrItem)
				enriched++
				break
			}
		}
	}

	e.logger.Info("Radarr enrichment complete",
		zap.Int("enriched", enriched),
	)

	return enriched
}

// EnrichWithSonarr enriches files with Sonarr metadata
func (e *Enricher) EnrichWithSonarr(
	ctx context.Context,
	files map[string]*EnrichedFile,
	sonarrMedia []*models.Media,
) int {
	e.logger.Info("Enriching with Sonarr data",
		zap.Int("sonarr_items", len(sonarrMedia)),
	)

	enriched := 0

	// Create path map for fast lookup
	sonarrByPath := make(map[string]*models.Media)
	for _, media := range sonarrMedia {
		sonarrByPath[media.FilePath] = media
	}

	// Match files with Sonarr items
	for path, file := range files {
		// Try exact match
		if sonarrItem, found := sonarrByPath[path]; found {
			e.applySonarrData(file, sonarrItem)
			enriched++
			continue
		}

		// Try matching with all hardlink paths (includes primary path)
		if file.IsHardlink {
			if sonarrItem, matched := matchByHardlinks(sonarrByPath, file.HardlinkPaths); matched {
				e.applySonarrData(file, sonarrItem)
				enriched++
				continue
			}
		}

		// Try fuzzy matching
		for sonarrPath, sonarrItem := range sonarrByPath {
			if e.pathsMatch(path, sonarrPath) {
				e.applySonarrData(file, sonarrItem)
				enriched++
				break
			}
		}
	}

	e.logger.Info("Sonarr enrichment complete",
		zap.Int("enriched", enriched),
	)

	return enriched
}

// EnrichWithJellyfin enriches files with Jellyfin metadata
func (e *Enricher) EnrichWithJellyfin(
	ctx context.Context,
	files map[string]*EnrichedFile,
	jellyfinMedia []*models.Media,
) int {
	e.logger.Info("Enriching with Jellyfin data",
		zap.Int("jellyfin_items", len(jellyfinMedia)),
	)

	enriched := 0

	// Create path map
	jellyfinByPath := make(map[string]*models.Media)
	for _, media := range jellyfinMedia {
		jellyfinByPath[media.FilePath] = media
	}

	// Match files
	for path, file := range files {
		if jfItem, found := jellyfinByPath[path]; found {
			e.applyJellyfinData(file, jfItem)
			enriched++
			continue
		}

		// Try matching with all hardlink paths (includes primary path)
		if file.IsHardlink {
			if jfItem, matched := matchByHardlinks(jellyfinByPath, file.HardlinkPaths); matched {
				e.applyJellyfinData(file, jfItem)
				enriched++
				continue
			}
		}

		for jfPath, jfItem := range jellyfinByPath {
			if e.pathsMatch(path, jfPath) {
				e.applyJellyfinData(file, jfItem)
				enriched++
				break
			}
		}
	}

	e.logger.Info("Jellyfin enrichment complete",
		zap.Int("enriched", enriched),
	)

	return enriched
}

// EnrichWithQBittorrent enriches files with qBittorrent metadata
func (e *Enricher) EnrichWithQBittorrent(
	ctx context.Context,
	files map[string]*EnrichedFile,
	torrentMap map[string]*models.TorrentInfo,
) int {
	e.logger.Info("Enriching with qBittorrent data",
		zap.Int("torrents", len(torrentMap)),
	)

	enriched := 0

	for path, file := range files {
		// Try exact match first (fast path)
		if torrent, found := torrentMap[path]; found {
			e.applyTorrentData(file, torrent)
			enriched++
			continue
		}

		// Try matching with all hardlink paths (includes primary path)
		if file.IsHardlink {
			if torrent, matched := matchByHardlinks(torrentMap, file.HardlinkPaths); matched {
				e.applyTorrentData(file, torrent)
				enriched++
				continue
			}
		}

		// Try intelligent directory-based matching
		// qBittorrent torrentMap is indexed by content_path (directories)
		// We need to check if this file is contained within any torrent's directory
		if torrent := e.findTorrentByDirectory(path, file.HardlinkPaths, torrentMap); torrent != nil {
			e.applyTorrentData(file, torrent)
			enriched++
			continue
		}
	}

	e.logger.Info("qBittorrent enrichment complete",
		zap.Int("enriched", enriched),
	)

	return enriched
}

// findTorrentByDirectory finds a torrent by checking if the file path is within a torrent's directory.
// This handles the case where qBittorrent returns directory paths (content_path, save_path)
// but we have file paths from the filesystem scanner.
func (e *Enricher) findTorrentByDirectory(
	filePath string,
	hardlinkPaths []string,
	torrentMap map[string]*models.TorrentInfo,
) *models.TorrentInfo {
	// Build an index mapping normalized torrent directories to torrent objects
	// This avoids O(n*m) complexity when checking multiple paths against many torrents
	torrentDirIndex := make(map[string]*models.TorrentInfo, len(torrentMap))
	for torrentPath, torrent := range torrentMap {
		cleanTorrentPath := filepath.Clean(torrentPath)
		normalizedTorrentPath := ensureTrailingSlash(cleanTorrentPath)
		torrentDirIndex[normalizedTorrentPath] = torrent
	}

	// Collect all paths to check (file path + hardlink paths)
	pathsToCheck := []string{filePath}
	if len(hardlinkPaths) > 0 {
		pathsToCheck = append(pathsToCheck, hardlinkPaths...)
	}

	// For each path, check if it matches a torrent
	for _, checkPath := range pathsToCheck {
		// Clean the path to handle edge cases like double slashes
		cleanCheckPath := filepath.Clean(checkPath)

		// First, check for exact file path match (single-file torrent case)
		// This handles torrents where content_path points to the file itself
		normalizedCheckPath := ensureTrailingSlash(cleanCheckPath)
		if torrent, found := torrentDirIndex[normalizedCheckPath]; found {
			e.logger.Debug("Matched torrent by exact file path",
				zap.String("file_path", checkPath),
				zap.String("torrent_path", cleanCheckPath),
				zap.String("torrent_hash", torrent.Hash),
			)
			return torrent
		}

		// Check if the file's directory is within any torrent directory
		checkDir := filepath.Dir(cleanCheckPath)
		normalizedCheckDir := ensureTrailingSlash(checkDir)

		// Try to find a matching torrent directory by prefix
		for torrentDir, torrent := range torrentDirIndex {
			if strings.HasPrefix(normalizedCheckDir, torrentDir) {
				e.logger.Debug("Matched torrent by directory",
					zap.String("file_path", checkPath),
					zap.String("file_dir", checkDir),
					zap.String("torrent_path", strings.TrimSuffix(torrentDir, "/")),
					zap.String("torrent_hash", torrent.Hash),
				)
				return torrent
			}
		}
	}

	return nil
}

// Helper methods to apply service data

func (e *Enricher) applyRadarrData(file *EnrichedFile, radarr *models.Media) {
	file.InRadarr = true
	file.RadarrID = radarr.RadarrID
	if file.Title == "" {
		file.Title = radarr.Title
	}
	if file.Quality == "" {
		file.Quality = radarr.Quality
	}
	if file.PosterURL == "" {
		file.PosterURL = radarr.PosterURL
	}
	file.Excluded = radarr.Excluded
}

func (e *Enricher) applySonarrData(file *EnrichedFile, sonarr *models.Media) {
	file.InSonarr = true
	file.SonarrID = sonarr.SonarrID
	if file.Title == "" {
		file.Title = sonarr.Title
	}
	if file.Quality == "" {
		file.Quality = sonarr.Quality
	}
	if file.PosterURL == "" {
		file.PosterURL = sonarr.PosterURL
	}
	file.Excluded = sonarr.Excluded
}

func (e *Enricher) applyJellyfinData(file *EnrichedFile, jellyfin *models.Media) {
	file.InJellyfin = true
	file.JellyfinID = jellyfin.JellyfinID
	if file.Title == "" {
		file.Title = jellyfin.Title
	}
	if file.PosterURL == "" {
		file.PosterURL = jellyfin.PosterURL
	}
	if jellyfin.LastWatched != nil {
		timestamp := jellyfin.LastWatched.Unix()
		file.LastWatched = &timestamp
	}
}

// EnrichWithJellystat sets the InJellystat flag for files tracked by Jellystat
// Since Jellystat monitors all Jellyfin activity, files in Jellyfin are tracked by Jellystat
func (e *Enricher) EnrichWithJellystat(
	ctx context.Context,
	files map[string]*EnrichedFile,
	jellystatEnabled bool,
) int {
	if !jellystatEnabled {
		return 0
	}

	e.logger.Info("Enriching with Jellystat tracking")

	enriched := 0
	for _, file := range files {
		// Jellystat tracks all items in Jellyfin
		if file.InJellyfin {
			file.InJellystat = true
			enriched++
		}
	}

	e.logger.Info("Jellystat enrichment complete",
		zap.Int("files_tracked", enriched),
	)

	return enriched
}

func (e *Enricher) applyTorrentData(file *EnrichedFile, torrent *models.TorrentInfo) {
	file.InQBittorrent = true
	file.TorrentHash = torrent.Hash
	file.TorrentCategory = torrent.Category
	file.TorrentTags = torrent.Tags
	file.TorrentState = torrent.State
	file.IsSeeding = torrent.IsSeeding
	file.SeedRatio = torrent.Ratio
}

// pathsMatch performs fuzzy path matching
func (e *Enricher) pathsMatch(path1, path2 string) bool {
	// Normalize paths
	norm1 := strings.ToLower(strings.TrimSuffix(path1, "/"))
	norm2 := strings.ToLower(strings.TrimSuffix(path2, "/"))

	// Check if one contains the other
	if strings.Contains(norm1, norm2) || strings.Contains(norm2, norm1) {
		return true
	}

	// Check if base directories match
	base1 := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(norm1), ".", " "), "_", " "))
	base2 := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(norm2), ".", " "), "_", " "))

	return strings.Contains(base1, base2) || strings.Contains(base2, base1)
}

// matchByHardlinks attempts to match a file against a service's path map using the file's hardlink paths.
// Returns the matched item and true if found, zero value and false otherwise.
// If hardlinkPaths is empty, returns (zero value, false).
func matchByHardlinks[T any](pathMap map[string]T, hardlinkPaths []string) (T, bool) {
	for _, hlPath := range hardlinkPaths {
		if item, found := pathMap[hlPath]; found {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// ensureTrailingSlash ensures a path ends with a trailing slash for directory comparison.
// This is used to prevent false positives when matching directory paths.
func ensureTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}
