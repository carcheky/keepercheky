package handler

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/internal/service/clients"
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

// MediaFileInfo represents media file information for display
type MediaFileInfo struct {
	ID              uint    `json:"id" gorm:"column:id"`
	Title           string  `json:"title" gorm:"column:title"`
	Type            string  `json:"type" gorm:"column:type"`
	FilePath        string  `json:"file_path" gorm:"column:file_path"`
	Size            int64   `json:"size" gorm:"column:size"`
	PosterURL       string  `json:"poster_url" gorm:"column:poster_url"`
	IsHardlink      bool    `json:"is_hardlink" gorm:"column:is_hardlink"`
	HardlinkPaths   string  `json:"hardlink_paths" gorm:"column:hardlink_paths"`
	PrimaryPath     string  `json:"primary_path" gorm:"column:primary_path"`
	InRadarr        bool    `json:"in_radarr" gorm:"column:in_radarr"`
	InSonarr        bool    `json:"in_sonarr" gorm:"column:in_sonarr"`
	InJellyfin      bool    `json:"in_jellyfin" gorm:"column:in_jellyfin"`
	InJellyseerr    bool    `json:"in_jellyseerr" gorm:"column:in_jellyseerr"`
	InQBittorrent   bool    `json:"in_qbittorrent" gorm:"column:in_qbittorrent"`
	TorrentHash     string  `json:"torrent_hash" gorm:"column:torrent_hash"`
	TorrentCategory string  `json:"torrent_category" gorm:"column:torrent_category"`
	TorrentState    string  `json:"torrent_state" gorm:"column:torrent_state"`
	IsSeeding       bool    `json:"is_seeding" gorm:"column:is_seeding"`
	SeedRatio       float64 `json:"seed_ratio" gorm:"column:seed_ratio"`
}

// RenderFilesPage renders the files listing page
func (h *FilesHandler) RenderFilesPage(c *fiber.Ctx) error {
	// Get all media sorted by:
	// 1. InQBittorrent DESC (downloads first)
	// 2. InJellyfin DESC (library items)
	// 3. FilePath ASC (alphabetical)
	var media []MediaFileInfo

	// Query with custom ordering: downloads first, then library, then alphabetical
	err := h.mediaRepo.GetDB().
		Table("media").
		Order("in_q_bittorrent DESC, in_jellyfin DESC, file_path ASC").
		Find(&media).Error

	if err != nil {
		return c.Status(500).SendString("Error loading files")
	}

	// Get paths from external services (Radarr, Sonarr, Jellyfin, qBittorrent)
	ctx := context.Background()
	servicePaths := h.getServicePaths(ctx)

	// Also extract paths from existing media files
	mediaPaths := h.extractPathsFromMedia(media)

	// Merge both path sources (service configs have priority)
	allPaths := h.mergePaths(servicePaths, mediaPaths)

	return c.Render("pages/files", fiber.Map{
		"Title": "Archivos del Sistema",
		"Files": media,
		"Paths": allPaths,
	}, "layouts/main")
}

// getServicePaths retrieves configured paths from external services
func (h *FilesHandler) getServicePaths(ctx context.Context) []PathInfo {
	var paths []PathInfo

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

	// Get Jellyfin library paths
	if h.config.Clients.Jellyfin.Enabled && h.syncService.GetJellyfinClient() != nil {
		jellyfinPaths := h.getJellyfinPaths(ctx)
		paths = append(paths, jellyfinPaths...)
	}

	// Get qBittorrent save paths
	if h.config.Clients.QBittorrent.Enabled && h.syncService.GetQBittorrentClient() != nil {
		qbittorrentPaths := h.getQBittorrentPaths(ctx)
		paths = append(paths, qbittorrentPaths...)
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

	// Get all torrents and extract unique SavePath directories
	torrentMap, err := client.GetAllTorrentsMap(ctx)
	if err != nil {
		h.logger.Error("Failed to get qBittorrent torrents for paths",
			zap.Error(err),
		)
		return paths
	}

	h.logger.Info("Retrieved qBittorrent torrents for paths",
		zap.Int("torrent_count", len(torrentMap)),
	)

	// Use a map to deduplicate SavePath directories
	savePathSet := make(map[string]bool)

	for _, torrent := range torrentMap {
		// Use SavePath as the base download directory
		if torrent.SavePath != "" {
			savePathSet[torrent.SavePath] = true
			h.logger.Debug("Found SavePath from qBittorrent",
				zap.String("save_path", torrent.SavePath),
				zap.String("category", torrent.Category),
			)
		}
	}

	h.logger.Info("Unique SavePaths from qBittorrent",
		zap.Int("unique_paths", len(savePathSet)),
	)

	// Convert to PathInfo slice
	for savePath := range savePathSet {
		label := "⬇️ qBittorrent: Descargas"
		paths = append(paths, PathInfo{
			Service: "qbittorrent",
			Type:    "download",
			Path:    savePath,
			Label:   label,
		})
		h.logger.Info("Added qBittorrent path",
			zap.String("path", savePath),
			zap.String("label", label),
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

// GetFilesAPI returns files as JSON for API access
func (h *FilesHandler) GetFilesAPI(c *fiber.Ctx) error {
	var media []struct {
		ID        uint   `json:"id" gorm:"column:id"`
		Title     string `json:"title" gorm:"column:title"`
		Type      string `json:"type" gorm:"column:type"`
		FilePath  string `json:"file_path" gorm:"column:file_path"`
		Size      int64  `json:"size" gorm:"column:size"`
		PosterURL string `json:"poster_url" gorm:"column:poster_url"`

		// Filesystem
		IsHardlink    bool   `json:"is_hardlink" gorm:"column:is_hardlink"`
		HardlinkPaths string `json:"hardlink_paths" gorm:"column:hardlink_paths"`
		PrimaryPath   string `json:"primary_path" gorm:"column:primary_path"`

		// Service flags
		InRadarr      bool `json:"in_radarr" gorm:"column:in_radarr"`
		InSonarr      bool `json:"in_sonarr" gorm:"column:in_sonarr"`
		InJellyfin    bool `json:"in_jellyfin" gorm:"column:in_jellyfin"`
		InJellyseerr  bool `json:"in_jellyseerr" gorm:"column:in_jellyseerr"`
		InQBittorrent bool `json:"in_qbittorrent" gorm:"column:in_qbittorrent"`

		// Torrent info
		TorrentHash     string  `json:"torrent_hash" gorm:"column:torrent_hash"`
		TorrentCategory string  `json:"torrent_category" gorm:"column:torrent_category"`
		TorrentState    string  `json:"torrent_state" gorm:"column:torrent_state"`
		IsSeeding       bool    `json:"is_seeding" gorm:"column:is_seeding"`
		SeedRatio       float64 `json:"seed_ratio" gorm:"column:seed_ratio"`
	}

	err := h.mediaRepo.GetDB().
		Table("media").
		Order("in_q_bittorrent DESC, in_jellyfin DESC, file_path ASC").
		Find(&media).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"files": media,
		"total": len(media),
	})
}
