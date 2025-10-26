package service

import (
	"context"
	"fmt"
	"time"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"github.com/carcheky/keepercheky/pkg/logger"
	"go.uber.org/zap"
)

// SyncService handles synchronization of media from external services.
type SyncService struct {
	mediaRepo         *repository.MediaRepository
	radarrClient      clients.MediaClient
	sonarrClient      clients.MediaClient
	jellyfinClient    clients.StreamingClient
	jellyseerrClient  clients.RequestClient
	qbittorrentClient *clients.QBittorrentClient
	logger            *zap.Logger
	config            *config.Config
}

// NewSyncService creates a new sync service.
func NewSyncService(
	mediaRepo *repository.MediaRepository,
	appLogger *logger.Logger,
	cfg *config.Config,
) *SyncService {
	// Get the underlying zap.Logger from our Logger wrapper
	zapLogger := appLogger.Desugar()

	svc := &SyncService{
		mediaRepo: mediaRepo,
		logger:    zapLogger,
		config:    cfg,
	}

	// Initialize clients based on configuration
	if cfg.Clients.Radarr.Enabled {
		svc.radarrClient = clients.NewRadarrClient(
			clients.ClientConfig{
				BaseURL: cfg.Clients.Radarr.URL,
				APIKey:  cfg.Clients.Radarr.APIKey,
				Timeout: 30 * time.Second,
			},
			zapLogger,
		)
	}

	if cfg.Clients.Sonarr.Enabled {
		svc.sonarrClient = clients.NewSonarrClient(
			clients.ClientConfig{
				BaseURL: cfg.Clients.Sonarr.URL,
				APIKey:  cfg.Clients.Sonarr.APIKey,
				Timeout: 30 * time.Second,
			},
			zapLogger,
		)
	}

	if cfg.Clients.Jellyfin.Enabled {
		svc.jellyfinClient = clients.NewJellyfinClient(
			clients.ClientConfig{
				BaseURL: cfg.Clients.Jellyfin.URL,
				APIKey:  cfg.Clients.Jellyfin.APIKey,
				Timeout: 30 * time.Second,
			},
			zapLogger,
		)
	}

	if cfg.Clients.Jellyseerr.Enabled {
		svc.jellyseerrClient = clients.NewJellyseerrClient(
			clients.ClientConfig{
				BaseURL: cfg.Clients.Jellyseerr.URL,
				APIKey:  cfg.Clients.Jellyseerr.APIKey,
				Timeout: 30 * time.Second,
			},
			zapLogger,
		)
	}

	if cfg.Clients.QBittorrent.Enabled {
		svc.qbittorrentClient = clients.NewQBittorrentClient(
			cfg.Clients.QBittorrent.URL,
			cfg.Clients.QBittorrent.Username,
			cfg.Clients.QBittorrent.Password,
			zapLogger,
		)
	}

	return svc
}

// SyncAll synchronizes media from all configured services.
func (s *SyncService) SyncAll(ctx context.Context) error {
	s.logger.Info("🔄 Starting FULL SYNC from all services (this will replace ALL media)")

	// CRITICAL: Delete ALL existing media to ensure clean sync
	s.logger.Info("🗑️  Clearing existing media database...")
	if err := s.mediaRepo.DeleteAll(); err != nil {
		s.logger.Error("❌ Failed to clear media database", zap.Error(err))
		return fmt.Errorf("failed to clear database: %w", err)
	}
	s.logger.Info("✅ Database cleared, starting fresh sync")

	// Invalidate cache on all services before syncing
	s.invalidateAllCaches(ctx)

	// Map to track media by unique key (title + year) for merging
	mediaMap := make(map[string]*models.Media)

	// 1. Sync from Radarr (movies - primary source)
	if s.radarrClient != nil {
		media, err := s.syncRadarr(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Radarr", zap.Error(err))
		} else {
			for _, m := range media {
				key := m.Title // Could use title+year for better matching
				mediaMap[key] = m
			}
		}
	}

	// 2. Sync from Sonarr (series - primary source)
	if s.sonarrClient != nil {
		media, err := s.syncSonarr(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Sonarr", zap.Error(err))
		} else {
			for _, m := range media {
				key := m.Title
				mediaMap[key] = m
			}
		}
	}

	// 3. Sync from Jellyfin (enrichment + additional media)
	if s.jellyfinClient != nil {
		if err := s.syncJellyfin(ctx, mediaMap); err != nil {
			s.logger.Error("Failed to sync Jellyfin", zap.Error(err))
		}
	}

	// 4. Save all media to database in batch
	s.logger.Info("💾 Saving media to database",
		zap.Int("total_items", len(mediaMap)),
	)

	savedCount := 0
	errorCount := 0

	for _, media := range mediaMap {
		// Since we cleared the database, we can use Create instead of CreateOrUpdate
		if err := s.mediaRepo.Create(media); err != nil {
			s.logger.Error("Failed to save media",
				zap.String("title", media.Title),
				zap.Error(err),
			)
			errorCount++
		} else {
			savedCount++
		}
	}

	// 4.5. Enrich with seeding status from qBittorrent
	if s.qbittorrentClient != nil {
		if err := s.enrichWithSeedingStatus(ctx, mediaMap); err != nil {
			s.logger.Error("Failed to enrich with seeding status", zap.Error(err))
		}
	}

	// NOTE: No need to cleanup orphaned media since we cleared everything before sync

	s.logger.Info("✅ Sync completed successfully",
		zap.Int("total_synced", len(mediaMap)),
		zap.Int("saved", savedCount),
		zap.Int("errors", errorCount),
	)

	return nil
}

// SyncProgress representa un mensaje de progreso durante la sincronización
type SyncProgress struct {
	Step    string `json:"step"`
	Message string `json:"message"`
	Status  string `json:"status"` // "processing", "success", "error"
	Data    any    `json:"data,omitempty"`
}

// SyncAllWithProgress synchronizes media from all configured services with progress reporting.
func (s *SyncService) SyncAllWithProgress(ctx context.Context, progressChan chan<- SyncProgress) error {
	s.logger.Info("🔄 Starting FULL SYNC with progress reporting")

	progressChan <- SyncProgress{
		Step:    "clear_db",
		Message: "🗑️  Limpiando base de datos existente...",
		Status:  "processing",
	}
	s.logger.Info("Sent progress: clear_db")

	// CRITICAL: Delete ALL existing media to ensure clean sync
	if err := s.mediaRepo.DeleteAll(); err != nil {
		s.logger.Error("❌ Failed to clear media database", zap.Error(err))
		return fmt.Errorf("failed to clear database: %w", err)
	}

	progressChan <- SyncProgress{
		Step:    "clear_db_complete",
		Message: "✅ Base de datos limpiada",
		Status:  "success",
	}
	s.logger.Info("Sent progress: clear_db_complete")

	// Invalidate cache on all services before syncing
	progressChan <- SyncProgress{
		Step:    "invalidate_cache",
		Message: "🔄 Invalidando cachés de servicios...",
		Status:  "processing",
	}
	s.invalidateAllCaches(ctx)

	// Map to track media by unique key (title + year) for merging
	mediaMap := make(map[string]*models.Media)

	// 1. Sync from Radarr (movies - primary source)
	if s.radarrClient != nil {
		progressChan <- SyncProgress{
			Step:    "sync_radarr",
			Message: "🎬 Sincronizando películas desde Radarr...",
			Status:  "processing",
		}

		media, err := s.radarrClient.GetLibrary(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Radarr", zap.Error(err))
			progressChan <- SyncProgress{
				Step:    "sync_radarr_error",
				Message: fmt.Sprintf("⚠️ Error al sincronizar Radarr: %v", err),
				Status:  "processing",
			}
		} else {
			for _, m := range media {
				key := m.Title
				mediaMap[key] = m
			}
			progressChan <- SyncProgress{
				Step:    "sync_radarr_complete",
				Message: fmt.Sprintf("✅ Radarr: %d películas obtenidas", len(media)),
				Status:  "success",
			}
		}
	}

	// 2. Sync from Sonarr (series - primary source)
	if s.sonarrClient != nil {
		progressChan <- SyncProgress{
			Step:    "sync_sonarr",
			Message: "📺 Sincronizando series desde Sonarr...",
			Status:  "processing",
		}

		media, err := s.sonarrClient.GetLibrary(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Sonarr", zap.Error(err))
			progressChan <- SyncProgress{
				Step:    "sync_sonarr_error",
				Message: fmt.Sprintf("⚠️ Error al sincronizar Sonarr: %v", err),
				Status:  "processing",
			}
		} else {
			for _, m := range media {
				key := m.Title
				mediaMap[key] = m
			}
			progressChan <- SyncProgress{
				Step:    "sync_sonarr_complete",
				Message: fmt.Sprintf("✅ Sonarr: %d series obtenidas", len(media)),
				Status:  "success",
			}
		}
	}

	// 3. Sync from Jellyfin (enrichment + additional media)
	if s.jellyfinClient != nil {
		progressChan <- SyncProgress{
			Step:    "sync_jellyfin",
			Message: "🎥 Sincronizando desde Jellyfin...",
			Status:  "processing",
		}

		if err := s.syncJellyfin(ctx, mediaMap); err != nil {
			s.logger.Error("Failed to sync Jellyfin", zap.Error(err))
			progressChan <- SyncProgress{
				Step:    "sync_jellyfin_error",
				Message: fmt.Sprintf("⚠️ Error al sincronizar Jellyfin: %v", err),
				Status:  "processing",
			}
		} else {
			progressChan <- SyncProgress{
				Step:    "sync_jellyfin_complete",
				Message: "✅ Jellyfin sincronizado",
				Status:  "success",
			}
		}
	}

	// 4. Save all media to database in batch
	progressChan <- SyncProgress{
		Step:    "save_db",
		Message: fmt.Sprintf("💾 Guardando %d elementos en base de datos...", len(mediaMap)),
		Status:  "processing",
	}

	savedCount := 0
	errorCount := 0

	for _, media := range mediaMap {
		// Since we cleared the database, we can use Create instead of CreateOrUpdate
		if err := s.mediaRepo.Create(media); err != nil {
			s.logger.Error("Failed to save media",
				zap.String("title", media.Title),
				zap.Error(err),
			)
			errorCount++
			continue
		}
		savedCount++
	}

	progressChan <- SyncProgress{
		Step:    "save_db_complete",
		Message: fmt.Sprintf("✅ Guardados: %d elementos (%d errores)", savedCount, errorCount),
		Status:  "success",
	}

	// 4.5. Enrich with seeding status from qBittorrent
	if s.qbittorrentClient != nil {
		progressChan <- SyncProgress{
			Step:    "enrich_torrents",
			Message: "🌱 Enriqueciendo con estado de torrents...",
			Status:  "processing",
		}

		if err := s.enrichWithSeedingStatus(ctx, mediaMap); err != nil {
			s.logger.Error("Failed to enrich with seeding status", zap.Error(err))
			progressChan <- SyncProgress{
				Step:    "enrich_torrents_error",
				Message: fmt.Sprintf("⚠️ Error al obtener estado de torrents: %v", err),
				Status:  "processing",
			}
		} else {
			progressChan <- SyncProgress{
				Step:    "enrich_torrents_complete",
				Message: "✅ Estado de torrents actualizado",
				Status:  "success",
			}
		}
	}

	s.logger.Info("✅ Sync completed successfully",
		zap.Int("total_synced", len(mediaMap)),
		zap.Int("saved", savedCount),
		zap.Int("errors", errorCount),
	)

	return nil
}

// syncRadarr syncs movies from Radarr.
func (s *SyncService) syncRadarr(ctx context.Context) ([]*models.Media, error) {
	s.logger.Info("Syncing from Radarr")

	media, err := s.radarrClient.GetLibrary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Radarr library: %w", err)
	}

	s.logger.Info("Radarr sync complete",
		zap.Int("movies", len(media)),
	)

	return media, nil
}

// syncSonarr syncs TV series from Sonarr.
func (s *SyncService) syncSonarr(ctx context.Context) ([]*models.Media, error) {
	s.logger.Info("Syncing from Sonarr")

	media, err := s.sonarrClient.GetLibrary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Sonarr library: %w", err)
	}

	s.logger.Info("Sonarr sync complete",
		zap.Int("series", len(media)),
	)

	return media, nil
}

// syncJellyfin syncs media from Jellyfin and merges with existing data.
func (s *SyncService) syncJellyfin(ctx context.Context, mediaMap map[string]*models.Media) error {
	s.logger.Info("Syncing from Jellyfin")

	jellyfinMedia, err := s.jellyfinClient.GetLibrary(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Jellyfin library: %w", err)
	}

	newFromJellyfin := 0
	enriched := 0
	skipped := 0
	totalJellyfinEpisodes := 0

	for _, jfMedia := range jellyfinMedia {
		// Skip if no valid data
		if jfMedia == nil || jfMedia.Title == "" {
			skipped++
			continue
		}

		// Log episode count for series
		if jfMedia.Type == "series" {
			totalJellyfinEpisodes += jfMedia.EpisodeFileCount
			s.logger.Info("Jellyfin series found",
				zap.String("title", jfMedia.Title),
				zap.Int("episodes", jfMedia.EpisodeFileCount),
			)
		}

		key := jfMedia.Title

		if existingMedia, exists := mediaMap[key]; exists {
			// Media already exists from Radarr/Sonarr, enrich it with Jellyfin data

			// Always update JellyfinID if available
			if jfMedia.JellyfinID != nil {
				existingMedia.JellyfinID = jfMedia.JellyfinID
			}

			// También copiar playback info si está disponible
			if jfMedia.LastWatched != nil {
				existingMedia.LastWatched = jfMedia.LastWatched
			}

			// Si Jellyfin tiene más episodios descargados que Sonarr, usar el de Jellyfin
			// Esto puede pasar si Jellyfin tiene episodios que Sonarr no rastrea
			if jfMedia.Type == "series" {
				if jfMedia.EpisodeFileCount > existingMedia.EpisodeFileCount {
					s.logger.Info("Updating episode count from Jellyfin",
						zap.String("title", existingMedia.Title),
						zap.Int("sonarr_count", existingMedia.EpisodeFileCount),
						zap.Int("jellyfin_count", jfMedia.EpisodeFileCount),
					)
					existingMedia.EpisodeFileCount = jfMedia.EpisodeFileCount
				} else {
					s.logger.Info("Keeping Sonarr episode count (higher or equal)",
						zap.String("title", existingMedia.Title),
						zap.Int("sonarr_count", existingMedia.EpisodeFileCount),
						zap.Int("jellyfin_count", jfMedia.EpisodeFileCount),
					)
				}
			}

			enriched++
		} else {
			// Media only exists in Jellyfin, add it to the map
			mediaMap[key] = jfMedia
			newFromJellyfin++
		}
	}

	s.logger.Info("Jellyfin sync complete",
		zap.Int("total_items", len(jellyfinMedia)),
		zap.Int("enriched", enriched),
		zap.Int("new_from_jellyfin", newFromJellyfin),
		zap.Int("skipped", skipped),
		zap.Int("total_jellyfin_episodes", totalJellyfinEpisodes),
	)

	return nil
}

// updateJellyfinPlayback updates playback information from Jellyfin for a single media item.
func (s *SyncService) updateJellyfinPlayback(ctx context.Context, media *models.Media) error {
	if media.JellyfinID == nil {
		return nil
	}

	playbackInfo, err := s.jellyfinClient.GetPlaybackInfo(ctx, *media.JellyfinID)
	if err != nil {
		return fmt.Errorf("failed to get playback info: %w", err)
	}

	if !playbackInfo.LastPlayed.IsZero() {
		media.LastWatched = &playbackInfo.LastPlayed
	}

	return nil
}

// TestConnection tests connection to a specific service.
func (s *SyncService) TestConnection(ctx context.Context, service string) error {
	switch service {
	case "radarr":
		if s.radarrClient == nil {
			return fmt.Errorf("Radarr not configured")
		}
		return s.radarrClient.TestConnection(ctx)

	case "sonarr":
		if s.sonarrClient == nil {
			return fmt.Errorf("Sonarr not configured")
		}
		return s.sonarrClient.TestConnection(ctx)

	case "jellyfin":
		if s.jellyfinClient == nil {
			return fmt.Errorf("Jellyfin not configured")
		}
		return s.jellyfinClient.TestConnection(ctx)

	case "jellyseerr":
		if s.jellyseerrClient == nil {
			return fmt.Errorf("Jellyseerr not configured")
		}
		return s.jellyseerrClient.TestConnection(ctx)

	case "qbittorrent":
		if s.qbittorrentClient == nil {
			return fmt.Errorf("qBittorrent not configured")
		}
		return s.qbittorrentClient.TestConnection(ctx)

	default:
		return fmt.Errorf("unknown service: %s", service)
	}
}

// GetRadarrSystemInfo returns complete system information from Radarr.
func (s *SyncService) GetRadarrSystemInfo(ctx context.Context) (*clients.RadarrSystemInfo, error) {
	if s.radarrClient == nil {
		return nil, fmt.Errorf("Radarr not configured")
	}

	// Type assertion to access RadarrClient specific methods
	radarrClient, ok := s.radarrClient.(*clients.RadarrClient)
	if !ok {
		return nil, fmt.Errorf("invalid Radarr client type")
	}

	return radarrClient.GetSystemInfo(ctx)
}

// GetSonarrSystemInfo returns complete system information from Sonarr.
func (s *SyncService) GetSonarrSystemInfo(ctx context.Context) (*clients.SonarrSystemInfo, error) {
	if s.sonarrClient == nil {
		return nil, fmt.Errorf("Sonarr not configured")
	}

	// Type assertion to access SonarrClient specific methods
	sonarrClient, ok := s.sonarrClient.(*clients.SonarrClient)
	if !ok {
		return nil, fmt.Errorf("invalid Sonarr client type")
	}

	return sonarrClient.GetSystemInfo(ctx)
}

// GetJellyfinSystemInfo returns complete system information from Jellyfin.
func (s *SyncService) GetJellyfinSystemInfo(ctx context.Context) (*clients.JellyfinSystemInfo, error) {
	if s.jellyfinClient == nil {
		return nil, fmt.Errorf("Jellyfin not configured")
	}

	// Type assertion to access JellyfinClient specific methods
	jellyfinClient, ok := s.jellyfinClient.(*clients.JellyfinClient)
	if !ok {
		return nil, fmt.Errorf("invalid Jellyfin client type")
	}

	return jellyfinClient.GetSystemInfo(ctx)
}

// GetJellyseerrSystemInfo returns complete system information from Jellyseerr.
func (s *SyncService) GetJellyseerrSystemInfo(ctx context.Context) (*clients.JellyseerrSystemInfo, error) {
	if s.jellyseerrClient == nil {
		return nil, fmt.Errorf("Jellyseerr not configured")
	}

	// Type assertion to access JellyseerrClient specific methods
	jellyseerrClient, ok := s.jellyseerrClient.(*clients.JellyseerrClient)
	if !ok {
		return nil, fmt.Errorf("invalid Jellyseerr client type")
	}

	return jellyseerrClient.GetSystemInfo(ctx)
}

// GetQBittorrentSystemInfo returns complete system information from qBittorrent.
func (s *SyncService) GetQBittorrentSystemInfo(ctx context.Context) (*clients.QBittorrentSystemInfo, error) {
	if s.qbittorrentClient == nil {
		return nil, fmt.Errorf("qBittorrent not configured")
	}

	return s.qbittorrentClient.GetSystemInfo(ctx)
}

// GetRadarrClient returns the Radarr client instance.
func (s *SyncService) GetRadarrClient() clients.MediaClient {
	return s.radarrClient
}

// GetSonarrClient returns the Sonarr client instance.
func (s *SyncService) GetSonarrClient() clients.MediaClient {
	return s.sonarrClient
}

// GetJellyfinClient returns the Jellyfin client instance.
func (s *SyncService) GetJellyfinClient() clients.StreamingClient {
	return s.jellyfinClient
}

// cleanupOrphanedMedia removes media from database that no longer exists in any service.
func (s *SyncService) cleanupOrphanedMedia(ctx context.Context, currentMedia map[string]*models.Media) error {
	s.logger.Info("Starting orphaned media cleanup")

	// Get all media from database
	allMedia, err := s.mediaRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all media: %w", err)
	}

	// Build a set of valid IDs from current sync
	validIDs := make(map[uint]bool)
	for _, media := range currentMedia {
		if media.ID > 0 {
			validIDs[media.ID] = true
		}
	}

	// Track which media to delete
	var toDelete []uint

	for _, media := range allMedia {
		// Check if media exists in current sync
		if validIDs[media.ID] {
			continue
		}

		// Media not in current sync - check if it still exists in ANY service
		existsInService := false

		// Check Radarr
		if media.RadarrID != nil && s.radarrClient != nil {
			if _, err := s.radarrClient.GetItem(ctx, *media.RadarrID); err == nil {
				existsInService = true
			}
		}

		// Check Sonarr
		if !existsInService && media.SonarrID != nil && s.sonarrClient != nil {
			if _, err := s.sonarrClient.GetItem(ctx, *media.SonarrID); err == nil {
				existsInService = true
			}
		}

		// Check Jellyfin
		if !existsInService && media.JellyfinID != nil && s.jellyfinClient != nil {
			if _, err := s.jellyfinClient.GetPlaybackInfo(ctx, *media.JellyfinID); err == nil {
				existsInService = true
			}
		}

		// If doesn't exist in any service, mark for deletion
		if !existsInService {
			toDelete = append(toDelete, media.ID)
			s.logger.Info("Found orphaned media",
				zap.Uint("id", media.ID),
				zap.String("title", media.Title),
				zap.Int("radarr_id", ptrIntValue(media.RadarrID)),
				zap.Int("sonarr_id", ptrIntValue(media.SonarrID)),
				zap.String("jellyfin_id", ptrStringValue(media.JellyfinID)),
			)
		}
	}

	// Delete orphaned media
	deletedCount := 0
	for _, id := range toDelete {
		if err := s.mediaRepo.Delete(id); err != nil {
			s.logger.Error("Failed to delete orphaned media",
				zap.Uint("id", id),
				zap.Error(err),
			)
		} else {
			deletedCount++
		}
	}

	s.logger.Info("Orphaned media cleanup complete",
		zap.Int("total_checked", len(allMedia)),
		zap.Int("orphaned_found", len(toDelete)),
		zap.Int("deleted", deletedCount),
	)

	return nil
}

// enrichWithSeedingStatus enriches media with seeding status from qBittorrent.
func (s *SyncService) enrichWithSeedingStatus(ctx context.Context, mediaMap map[string]*models.Media) error {
	s.logger.Info("Enriching media with seeding status from qBittorrent")

	enrichedCount := 0
	errorCount := 0

	for _, media := range mediaMap {
		// Skip if media has no file path
		if media.FilePath == "" {
			continue
		}

		// Try to get torrent info by file path
		torrentInfo, err := s.qbittorrentClient.GetTorrentByPath(ctx, media.FilePath)
		if err != nil {
			// Not an error if torrent doesn't exist - media might not be from torrent
			continue
		}

		// Update media with torrent info
		media.TorrentHash = torrentInfo.Hash
		media.IsSeeding = torrentInfo.IsSeeding
		media.SeedRatio = torrentInfo.Ratio

		// Save updated media
		if err := s.mediaRepo.CreateOrUpdate(media); err != nil {
			s.logger.Error("Failed to update media with seeding status",
				zap.String("title", media.Title),
				zap.Error(err),
			)
			errorCount++
		} else {
			enrichedCount++
		}
	}

	s.logger.Info("Seeding status enrichment complete",
		zap.Int("enriched", enrichedCount),
		zap.Int("errors", errorCount),
	)

	return nil
}

// invalidateAllCaches invalidates caches on all configured services to ensure fresh data.
func (s *SyncService) invalidateAllCaches(ctx context.Context) {
	s.logger.Info("Invalidating caches on all services before sync")

	// Invalidate Jellyfin cache
	if s.jellyfinClient != nil {
		jellyfinClient, ok := s.jellyfinClient.(*clients.JellyfinClient)
		if ok {
			if err := jellyfinClient.InvalidateCache(ctx); err != nil {
				s.logger.Warn("Failed to invalidate Jellyfin cache",
					zap.Error(err),
				)
			}
		}
	}

	// Note: Radarr and Sonarr already have Cache-Control headers in HTTP client
	// so their requests are never cached. No additional action needed.

	s.logger.Info("Cache invalidation complete")
}

// Helper functions for logging pointer values
func ptrIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func ptrStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
