// Package service provides business logic services.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SyncService handles synchronization of media from external services.
type SyncService struct {
	mediaRepo        *repository.MediaRepository
	radarrClient     clients.MediaClient
	sonarrClient     clients.MediaClient
	jellyfinClient   clients.StreamingClient
	jellyseerrClient clients.RequestClient
	logger           *zap.Logger
	config           *config.Config
}

// NewSyncService creates a new sync service.
func NewSyncService(
	mediaRepo *repository.MediaRepository,
	logger *zap.Logger,
	cfg *config.Config,
) *SyncService {
	svc := &SyncService{
		mediaRepo: mediaRepo,
		logger:    logger,
		config:    cfg,
	}

	// Initialize clients based on configuration
	if cfg.Services.Radarr.Enabled {
		svc.radarrClient = clients.NewRadarrClient(
			clients.ClientConfig{
				BaseURL: cfg.Services.Radarr.URL,
				APIKey:  cfg.Services.Radarr.APIKey,
				Timeout: 30 * time.Second,
			},
			logger,
		)
	}

	if cfg.Services.Sonarr.Enabled {
		svc.sonarrClient = clients.NewSonarrClient(
			clients.ClientConfig{
				BaseURL: cfg.Services.Sonarr.URL,
				APIKey:  cfg.Services.Sonarr.APIKey,
				Timeout: 30 * time.Second,
			},
			logger,
		)
	}

	if cfg.Services.Jellyfin.Enabled {
		svc.jellyfinClient = clients.NewJellyfinClient(
			clients.ClientConfig{
				BaseURL: cfg.Services.Jellyfin.URL,
				APIKey:  cfg.Services.Jellyfin.APIKey,
				Timeout: 30 * time.Second,
			},
			logger,
		)
	}

	if cfg.Services.Jellyseerr.Enabled {
		svc.jellyseerrClient = clients.NewJellyseerrClient(
			clients.ClientConfig{
				BaseURL: cfg.Services.Jellyseerr.URL,
				APIKey:  cfg.Services.Jellyseerr.APIKey,
				Timeout: 30 * time.Second,
			},
			logger,
		)
	}

	return svc
}

// SyncAll synchronizes media from all configured services.
func (s *SyncService) SyncAll(ctx context.Context) error {
	s.logger.Info("Starting full sync from all services")
	
	var allMedia []*models.Media
	
	// Sync from Radarr
	if s.radarrClient != nil {
		media, err := s.syncRadarr(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Radarr", zap.Error(err))
		} else {
			allMedia = append(allMedia, media...)
		}
	}
	
	// Sync from Sonarr
	if s.sonarrClient != nil {
		media, err := s.syncSonarr(ctx)
		if err != nil {
			s.logger.Error("Failed to sync Sonarr", zap.Error(err))
		} else {
			allMedia = append(allMedia, media...)
		}
	}
	
	// Update playback info from Jellyfin
	if s.jellyfinClient != nil {
		if err := s.updateJellyfinPlayback(ctx, allMedia); err != nil {
			s.logger.Error("Failed to update Jellyfin playback", zap.Error(err))
		}
	}
	
	// Save all media to database
	for _, media := range allMedia {
		if err := s.mediaRepo.CreateOrUpdate(media); err != nil {
			s.logger.Error("Failed to save media",
				zap.String("title", media.Title),
				zap.Error(err),
			)
		}
	}
	
	s.logger.Info("Sync completed",
		zap.Int("total_synced", len(allMedia)),
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

// updateJellyfinPlayback updates playback information from Jellyfin.
func (s *SyncService) updateJellyfinPlayback(ctx context.Context, mediaList []*models.Media) error {
	s.logger.Info("Updating Jellyfin playback info")
	
	for _, media := range mediaList {
		if media.JellyfinID == nil {
			continue
		}
		
		playbackInfo, err := s.jellyfinClient.GetPlaybackInfo(ctx, *media.JellyfinID)
		if err != nil {
			s.logger.Warn("Failed to get playback info",
				zap.String("jellyfin_id", *media.JellyfinID),
				zap.Error(err),
			)
			continue
		}
		
		if !playbackInfo.LastPlayed.IsZero() {
			media.LastWatched = &playbackInfo.LastPlayed
		}
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
		
	default:
		return fmt.Errorf("unknown service: %s", service)
	}
}
