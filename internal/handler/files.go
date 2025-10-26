package handler

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"github.com/carcheky/keepercheky/pkg/filesystem"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// FilesHandler handles file listing and management
type FilesHandler struct {
	mediaRepo   *repository.MediaRepository
	config      *config.Config
	syncService *service.SyncService
	logger      *zap.Logger
}

// NewFilesHandler creates a new files handler
func NewFilesHandler(mediaRepo *repository.MediaRepository, cfg *config.Config, syncService *service.SyncService, logger *zap.Logger) *FilesHandler {
	return &FilesHandler{
		mediaRepo:   mediaRepo,
		config:      cfg,
		syncService: syncService,
		logger:      logger,
	}
}

// PathInfo represents monitored path information from external services
type PathInfo struct {
	Service string `json:"service"` // radarr, sonarr, jellyfin, qbittorrent
	Type    string `json:"type"`    // library, download
	Path    string `json:"path"`
	Label   string `json:"label"` // Human-readable label
}

// UserViewInfo represents user viewing information
type UserViewInfo struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	PlayCount  int    `json:"play_count"`
	LastPlayed string `json:"last_played"`
}

// MediaFileInfo represents media file information for display
type MediaFileInfo struct {
	ID            uint   `json:"id" gorm:"column:id"`
	Title         string `json:"title" gorm:"column:title"`
	Type          string `json:"type" gorm:"column:type"`
	FilePath      string `json:"file_path" gorm:"column:file_path"`
	Size          int64  `json:"size" gorm:"column:size"`
	PosterURL     string `json:"poster_url" gorm:"column:poster_url"`
	Quality       string `json:"quality" gorm:"column:quality"`
	IsHardlink    bool   `json:"is_hardlink" gorm:"column:is_hardlink"`
	HardlinkPaths string `json:"hardlink_paths" gorm:"column:hardlink_paths"`
	PrimaryPath   string `json:"primary_path" gorm:"column:primary_path"`

	// Service flags
	InRadarr      bool `json:"in_radarr" gorm:"column:in_radarr"`
	InSonarr      bool `json:"in_sonarr" gorm:"column:in_sonarr"`
	InJellyfin    bool `json:"in_jellyfin" gorm:"column:in_jellyfin"`
	InJellyseerr  bool `json:"in_jellyseerr" gorm:"column:in_jellyseerr"`
	InQBittorrent bool `json:"in_qbittorrent" gorm:"column:in_qbittorrent"`

	// Service IDs
	RadarrID     *int    `json:"radarr_id" gorm:"column:radarr_id"`
	SonarrID     *int    `json:"sonarr_id" gorm:"column:sonarr_id"`
	JellyfinID   *string `json:"jellyfin_id" gorm:"column:jellyfin_id"`
	JellyseerrID *int    `json:"jellyseerr_id" gorm:"column:jellyseerr_id"`

	// Torrent info
	TorrentHash     string  `json:"torrent_hash" gorm:"column:torrent_hash"`
	TorrentCategory string  `json:"torrent_category" gorm:"column:torrent_category"`
	TorrentState    string  `json:"torrent_state" gorm:"column:torrent_state"`
	TorrentTags     string  `json:"torrent_tags" gorm:"column:torrent_tags"`
	IsSeeding       bool    `json:"is_seeding" gorm:"column:is_seeding"`
	SeedRatio       float64 `json:"seed_ratio" gorm:"column:seed_ratio"`

	// Viewing information (not stored in DB, calculated on-the-fly)
	HasBeenWatched  bool           `json:"has_been_watched" gorm:"-"`
	WatchedByUsers  []UserViewInfo `json:"watched_by_users" gorm:"-"`
	TotalPlayCount  int            `json:"total_play_count" gorm:"-"`
	LastWatchedDate string         `json:"last_watched_date" gorm:"-"`

	// Metadata
	Tags     []string `json:"tags" gorm:"-"`
	Excluded bool     `json:"excluded" gorm:"column:excluded"`
}

// RenderFilesPage renders the files listing page
func (h *FilesHandler) RenderFilesPage(c *fiber.Ctx) error {
	// Get all media sorted by:
	// 1. InQBittorrent DESC (downloads first)
	// 2. InJellyfin DESC (library items)
	// 3. FilePath ASC (alphabetical)
	ctx := context.Background()

	// Get service paths first (for scanning)
	servicePaths := h.getServicePaths(ctx)

	// Scan filesystem directly from those paths
	scannedFiles, err := h.scanFilesystemFromPaths(servicePaths)
	if err != nil {
		h.logger.Error("Failed to scan filesystem",
			zap.Error(err),
		)
		return c.Status(500).Render("pages/files", fiber.Map{
			"Title": "Archivos del Sistema",
			"Files": []MediaFileInfo{},
			"Paths": servicePaths,
		}, "layouts/main")
	}

	h.logger.Info("Files scanned from filesystem",
		zap.Int("total_files", len(scannedFiles)),
	)

	return c.Render("pages/files", fiber.Map{
		"Title": "Archivos del Sistema",
		"Files": scannedFiles,
		"Paths": servicePaths,
	}, "layouts/main")
}

// getServicePaths retrieves configured paths from external services
func (h *FilesHandler) getServicePaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

	// Get qBittorrent save paths FIRST (origin/source)
	if h.config.Clients.QBittorrent.Enabled && h.syncService.GetQBittorrentClient() != nil {
		qbittorrentPaths := h.getQBittorrentPaths(ctx)
		paths = append(paths, qbittorrentPaths...)
	}

	// Get Jellyfin library paths (destination)
	if h.config.Clients.Jellyfin.Enabled && h.syncService.GetJellyfinClient() != nil {
		jellyfinPaths := h.getJellyfinPaths(ctx)
		paths = append(paths, jellyfinPaths...)
	}

	// Get Radarr root folders
	if h.config.Clients.Radarr.Enabled && h.syncService.GetRadarrClient() != nil {
		radarrPaths := h.getRadarrPaths(ctx)
		paths = append(paths, radarrPaths...)
	}

	// Get Sonarr root folders
	if h.config.Clients.Sonarr.Enabled && h.syncService.GetSonarrClient() != nil {
		sonarrPaths := h.getSonarrPaths(ctx)
		paths = append(paths, sonarrPaths...)
	}

	return paths
}

// getRadarrPaths retrieves root folder paths from Radarr
func (h *FilesHandler) getRadarrPaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

	client := h.syncService.GetRadarrClient()
	if client == nil {
		return paths
	}

	// Access internal resty client via reflection/type assertion
	// We'll use a simpler approach: create a method in the client or access via config
	// For now, let's use the config-based approach

	// Since we can't easily access the internal HTTP client,
	// we'll rely on media files to discover paths for now
	// TODO: Add GetRootFolders() method to RadarrClient interface

	return paths
}

// getSonarrPaths retrieves root folder paths from Sonarr
func (h *FilesHandler) getSonarrPaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

	client := h.syncService.GetSonarrClient()
	if client == nil {
		return paths
	}

	// TODO: Add GetRootFolders() method to SonarrClient interface

	return paths
}

// getJellyfinPaths retrieves library folder paths from Jellyfin
func (h *FilesHandler) getJellyfinPaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

	client := h.syncService.GetJellyfinClient()
	if client == nil {
		return paths
	}

	// Try to get virtual folders (library root paths) from Jellyfin
	// This requires type assertion to access the GetVirtualFolders method
	jellyfinClient, ok := client.(*clients.JellyfinClient)
	if !ok {
		h.logger.Warn("Could not type assert JellyfinClient")
		return paths
	}

	// Get library virtual folders
	folders, err := jellyfinClient.GetVirtualFolders(ctx)
	if err != nil {
		h.logger.Error("Failed to get Jellyfin virtual folders",
			zap.Error(err),
		)
		return paths
	}

	// Extract library paths from virtual folders
	for _, folder := range folders {
		for _, location := range folder.Locations {
			label := fmt.Sprintf("📚 Jellyfin: %s", folder.Name)
			paths = append(paths, PathInfo{
				Service: "jellyfin",
				Type:    "library",
				Path:    location,
				Label:   label,
			})
		}
	}

	return paths
} // getQBittorrentPaths retrieves download paths from qBittorrent
func (h *FilesHandler) getQBittorrentPaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

	client := h.syncService.GetQBittorrentClient()
	if client == nil {
		h.logger.Warn("qBittorrent client is not configured")
		return paths
	}

	// Get qBittorrent preferences (configuration) for default save path (completed downloads)
	prefs, err := client.GetPreferences(ctx)
	if err != nil {
		h.logger.Error("Failed to get qBittorrent preferences",
			zap.Error(err),
		)
		return paths
	}

	// Add completed downloads path (ExportDirFin) if available, otherwise SavePath
	completedPath := prefs.ExportDirFin
	if completedPath == "" {
		completedPath = prefs.SavePath
	}

	if completedPath != "" {
		paths = append(paths, PathInfo{
			Service: "qbittorrent",
			Type:    "download",
			Path:    completedPath,
			Label:   "⬇️ qBittorrent: Descargas completadas",
		})
		h.logger.Info("Added qBittorrent completed downloads path",
			zap.String("path", completedPath),
		)
	}

	return paths
}

// mergePaths combines service paths and media paths, removing duplicates
func (h *FilesHandler) mergePaths(servicePaths, mediaPaths []PathInfo) []PathInfo {
	pathMap := make(map[string]PathInfo)

	// Add service paths first (higher priority)
	for _, p := range servicePaths {
		// Normalize path for comparison
		normalizedPath := filepath.Clean(p.Path)
		key := p.Service + ":" + normalizedPath
		pathMap[key] = p
	}

	// Add media paths (only if not already present)
	for _, p := range mediaPaths {
		normalizedPath := filepath.Clean(p.Path)
		key := p.Service + ":" + normalizedPath
		if _, exists := pathMap[key]; !exists {
			pathMap[key] = p
		}
	}

	// Convert map to slice
	result := make([]PathInfo, 0, len(pathMap))
	for _, path := range pathMap {
		result = append(result, path)
	}

	return result
}

// extractPathsFromMedia extracts unique directory paths from media files
func (h *FilesHandler) extractPathsFromMedia(media []MediaFileInfo) []PathInfo {
	pathMap := make(map[string]PathInfo)

	for _, m := range media {
		if m.FilePath == "" {
			continue
		}

		// Get parent directory
		dir := m.FilePath
		// Remove filename, keep directory path
		for i := len(dir) - 1; i >= 0; i-- {
			if dir[i] == '/' {
				dir = dir[:i]
				break
			}
		}

		// Determine service and type based on flags
		var service, pathType, label string

		// Check for qBittorrent downloads
		if m.InQBittorrent && m.TorrentCategory != "" {
			service = "qbittorrent"
			pathType = "download"
			label = "⬇️ qBittorrent: " + m.TorrentCategory
			key := service + ":" + dir
			if _, exists := pathMap[key]; !exists {
				pathMap[key] = PathInfo{
					Service: service,
					Type:    pathType,
					Path:    dir,
					Label:   label,
				}
			}
		}

		// Check for Radarr movies
		if m.InRadarr && m.Type == "movie" {
			service = "radarr"
			pathType = "library"
			label = "🎬 Radarr: Películas"
			key := service + ":" + dir
			if _, exists := pathMap[key]; !exists {
				pathMap[key] = PathInfo{
					Service: service,
					Type:    pathType,
					Path:    dir,
					Label:   label,
				}
			}
		}

		// Check for Sonarr series
		if m.InSonarr && (m.Type == "series" || m.Type == "episode") {
			service = "sonarr"
			pathType = "library"
			label = "📺 Sonarr: Series"
			key := service + ":" + dir
			if _, exists := pathMap[key]; !exists {
				pathMap[key] = PathInfo{
					Service: service,
					Type:    pathType,
					Path:    dir,
					Label:   label,
				}
			}
		}

		// Check for Jellyfin library
		if m.InJellyfin {
			service = "jellyfin"
			pathType = "library"
			if m.Type == "movie" {
				label = "📚 Jellyfin: Biblioteca de Películas"
			} else if m.Type == "series" || m.Type == "episode" {
				label = "📚 Jellyfin: Biblioteca de Series"
			} else {
				label = "📚 Jellyfin: Biblioteca"
			}
			key := service + ":" + dir
			if _, exists := pathMap[key]; !exists {
				pathMap[key] = PathInfo{
					Service: service,
					Type:    pathType,
					Path:    dir,
					Label:   label,
				}
			}
		}
	}

	// Convert map to slice
	paths := make([]PathInfo, 0, len(pathMap))
	for _, path := range pathMap {
		paths = append(paths, path)
	}

	return paths
}

// scanFilesystemFromPaths scans the filesystem from service paths and returns unified files
func (h *FilesHandler) scanFilesystemFromPaths(servicePaths []PathInfo) ([]MediaFileInfo, error) {
	ctx := context.Background()

	// Extract unique root paths and categorize them
	var rootPaths []string
	var libraryPaths []string
	var downloadPaths []string

	pathSet := make(map[string]bool)

	for _, pathInfo := range servicePaths {
		if pathSet[pathInfo.Path] {
			continue
		}
		pathSet[pathInfo.Path] = true
		rootPaths = append(rootPaths, pathInfo.Path)

		if pathInfo.Type == "library" {
			libraryPaths = append(libraryPaths, pathInfo.Path)
		} else if pathInfo.Type == "download" {
			downloadPaths = append(downloadPaths, pathInfo.Path)
		}
	}

	if len(rootPaths) == 0 {
		h.logger.Warn("No paths to scan")
		return []MediaFileInfo{}, nil
	}

	// Create scanner with options
	scanOptions := filesystem.ScanOptions{
		RootPaths:       rootPaths,
		LibraryPaths:    libraryPaths,
		DownloadPaths:   downloadPaths,
		VideoExtensions: []string{".mkv", ".mp4", ".avi", ".m4v", ".ts", ".m2ts", ".wmv", ".flv", ".webm"},
		MinSize:         50 * 1024 * 1024, // 50MB minimum
		SkipHidden:      true,
	}

	scanner := filesystem.NewScanner(scanOptions, h.logger)

	// Scan filesystem
	fileEntries, err := scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("filesystem scan failed: %w", err)
	}

	// Convert FileEntry map to EnrichedFile map for enrichment
	enrichedFiles := make(map[string]*filesystem.EnrichedFile)
	for path, entry := range fileEntries {
		enrichedFiles[path] = &filesystem.EnrichedFile{
			FileEntry: entry,
			ModTime:   entry.ModTime,
		}
	}

	// Enrich with service data
	h.enrichWithServices(ctx, enrichedFiles)

	// Convert to MediaFileInfo, grouping hardlinks
	mediaFiles := make([]MediaFileInfo, 0, len(enrichedFiles))
	processedInodes := make(map[uint64]bool)

	for _, enrichedFile := range enrichedFiles {
		// Skip if we already processed this inode (hardlink group)
		if enrichedFile.IsHardlink && processedInodes[enrichedFile.Inode] {
			continue
		}

		// Mark inode as processed
		if enrichedFile.IsHardlink {
			processedInodes[enrichedFile.Inode] = true
		}

		// Use enriched title if available, otherwise infer from filename
		title := enrichedFile.Title
		if title == "" {
			title = h.inferTitleFromPath(enrichedFile.PrimaryPath)
		}

		// Build hardlink paths string
		hardlinkPathsStr := ""
		if enrichedFile.IsHardlink && len(enrichedFile.HardlinkPaths) > 0 {
			hardlinkPathsStr = strings.Join(enrichedFile.HardlinkPaths, "|")
		}

		// Convert tags
		tags := enrichedFile.Tags
		if tags == nil {
			tags = []string{}
		}

		mediaFile := MediaFileInfo{
			ID:              0, // No database ID since this is direct filesystem scan
			Title:           title,
			Type:            enrichedFile.MediaType,
			FilePath:        enrichedFile.PrimaryPath,
			Size:            enrichedFile.Size,
			PosterURL:       enrichedFile.PosterURL,
			Quality:         enrichedFile.Quality,
			IsHardlink:      enrichedFile.IsHardlink,
			HardlinkPaths:   hardlinkPathsStr,
			PrimaryPath:     enrichedFile.PrimaryPath,
			InRadarr:        enrichedFile.InRadarr,
			InSonarr:        enrichedFile.InSonarr,
			InJellyfin:      enrichedFile.InJellyfin,
			InJellyseerr:    enrichedFile.InJellyseerr,
			InQBittorrent:   enrichedFile.InQBittorrent,
			RadarrID:        enrichedFile.RadarrID,
			SonarrID:        enrichedFile.SonarrID,
			JellyfinID:      enrichedFile.JellyfinID,
			JellyseerrID:    enrichedFile.JellyseerrID,
			TorrentHash:     enrichedFile.TorrentHash,
			TorrentCategory: enrichedFile.TorrentCategory,
			TorrentState:    enrichedFile.TorrentState,
			TorrentTags:     enrichedFile.TorrentTags,
			IsSeeding:       enrichedFile.IsSeeding,
			SeedRatio:       enrichedFile.SeedRatio,
			Tags:            tags,
			Excluded:        enrichedFile.Excluded,
			HasBeenWatched:  false,
			WatchedByUsers:  []UserViewInfo{},
			TotalPlayCount:  0,
			LastWatchedDate: "",
		}

		// Add watching info if available
		if enrichedFile.LastWatched != nil {
			mediaFile.HasBeenWatched = true
			// TODO: Get detailed user watching info from Jellyfin/Jellystat
		}

		mediaFiles = append(mediaFiles, mediaFile)
	}

	h.logger.Info("Filesystem scan complete",
		zap.Int("total_entries", len(fileEntries)),
		zap.Int("unique_files", len(mediaFiles)),
		zap.Int("hardlink_groups", len(processedInodes)),
	)

	return mediaFiles, nil
}

// enrichWithServices enriches files with data from all configured services
func (h *FilesHandler) enrichWithServices(ctx context.Context, files map[string]*filesystem.EnrichedFile) {
	enricher := filesystem.NewEnricher(h.logger)

	// Enrich with Radarr
	if h.config.Clients.Radarr.Enabled && h.syncService.GetRadarrClient() != nil {
		h.logger.Info("Enriching with Radarr data")
		radarrMedia, err := h.syncService.GetRadarrClient().GetLibrary(ctx)
		if err != nil {
			h.logger.Error("Failed to get Radarr library", zap.Error(err))
		} else {
			count := enricher.EnrichWithRadarr(ctx, files, radarrMedia)
			h.logger.Info("Radarr enrichment complete", zap.Int("enriched", count))
		}
	}

	// Enrich with Sonarr
	if h.config.Clients.Sonarr.Enabled && h.syncService.GetSonarrClient() != nil {
		h.logger.Info("Enriching with Sonarr data")
		sonarrMedia, err := h.syncService.GetSonarrClient().GetLibrary(ctx)
		if err != nil {
			h.logger.Error("Failed to get Sonarr library", zap.Error(err))
		} else {
			count := enricher.EnrichWithSonarr(ctx, files, sonarrMedia)
			h.logger.Info("Sonarr enrichment complete", zap.Int("enriched", count))
		}
	}

	// Enrich with Jellyfin
	if h.config.Clients.Jellyfin.Enabled && h.syncService.GetJellyfinClient() != nil {
		h.logger.Info("Enriching with Jellyfin data")
		jellyfinMedia, err := h.syncService.GetJellyfinClient().GetLibrary(ctx)
		if err != nil {
			h.logger.Error("Failed to get Jellyfin library", zap.Error(err))
		} else {
			count := enricher.EnrichWithJellyfin(ctx, files, jellyfinMedia)
			h.logger.Info("Jellyfin enrichment complete", zap.Int("enriched", count))
		}
	}

	// Enrich with qBittorrent
	if h.config.Clients.QBittorrent.Enabled && h.syncService.GetQBittorrentClient() != nil {
		h.logger.Info("Enriching with qBittorrent data")
		torrentMap, err := h.syncService.GetQBittorrentClient().GetAllTorrentsMap(ctx)
		if err != nil {
			h.logger.Error("Failed to get qBittorrent torrents", zap.Error(err))
		} else {
			count := enricher.EnrichWithQBittorrent(ctx, files, torrentMap)
			h.logger.Info("qBittorrent enrichment complete", zap.Int("enriched", count))
		}
	}
}

// inferTitleFromPath extracts a title from the file path
func (h *FilesHandler) inferTitleFromPath(path string) string {
	// Get filename without extension
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	nameWithoutExt := strings.TrimSuffix(base, ext)

	// Clean up common patterns
	// Remove quality tags like [1080p], (BluRay), etc.
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, ".", " ")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "_", " ")

	// Simple cleanup
	title := strings.TrimSpace(nameWithoutExt)

	// Limit length
	if len(title) > 80 {
		title = title[:80] + "..."
	}

	return title
}

// GetFilesAPI returns files as JSON for API access
func (h *FilesHandler) GetFilesAPI(c *fiber.Ctx) error {
	var media []MediaFileInfo

	err := h.mediaRepo.GetDB().
		Table("media").
		Select(`
			id, title, type, file_path, size, poster_url, quality,
			is_hardlink, hardlink_paths, primary_path,
			in_radarr, in_sonarr, in_jellyfin, in_jellyseerr, in_q_bittorrent,
			radarr_id, sonarr_id, jellyfin_id, jellyseerr_id,
			torrent_hash, torrent_category, torrent_state, torrent_tags,
			is_seeding, seed_ratio, excluded
		`).
		Order("in_q_bittorrent DESC, in_jellyfin DESC, file_path ASC").
		Find(&media).Error

	if err != nil {
		h.logger.Error("Failed to query media from database", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Enrich with viewing information if needed
	// TODO: Add user viewing info from Jellyfin/Jellystat when available

	// Get last sync time from settings
	var lastSyncSetting struct {
		Value string `gorm:"column:value"`
	}
	lastSyncTime := ""
	err = h.mediaRepo.GetDB().
		Table("settings").
		Select("value").
		Where("key = ?", "last_files_sync").
		First(&lastSyncSetting).Error

	if err == nil {
		lastSyncTime = lastSyncSetting.Value
	}

	return c.JSON(fiber.Map{
		"files":     media,
		"total":     len(media),
		"last_sync": lastSyncTime,
	})
}
