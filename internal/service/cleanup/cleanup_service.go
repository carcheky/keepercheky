package cleanup

import (
	"context"
	"fmt"
	"os"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"go.uber.org/zap"
)

// DeleteOptions represents the options for deleting media.
type DeleteOptions struct {
	Radarr      bool `json:"radarr"`
	Sonarr      bool `json:"sonarr"`
	Jellyfin    bool `json:"jellyfin"`
	DeleteFiles bool `json:"deleteFiles"`
	QBittorrent bool `json:"qbittorrent"`
}

// CleanupService handles media cleanup operations.
type CleanupService struct {
	mediaRepo         *repository.MediaRepository
	historyRepo       *repository.HistoryRepository
	radarrClient      clients.MediaClient
	sonarrClient      clients.MediaClient
	jellyfinClient    clients.StreamingClient
	qbittorrentClient *clients.QBittorrentClient
	logger            *zap.Logger
}

// NewCleanupService creates a new cleanup service.
func NewCleanupService(
	mediaRepo *repository.MediaRepository,
	historyRepo *repository.HistoryRepository,
	radarrClient clients.MediaClient,
	sonarrClient clients.MediaClient,
	jellyfinClient clients.StreamingClient,
	qbittorrentClient *clients.QBittorrentClient,
	logger *zap.Logger,
) *CleanupService {
	return &CleanupService{
		mediaRepo:         mediaRepo,
		historyRepo:       historyRepo,
		radarrClient:      radarrClient,
		sonarrClient:      sonarrClient,
		jellyfinClient:    jellyfinClient,
		qbittorrentClient: qbittorrentClient,
		logger:            logger,
	}
}

// DeleteMedia deletes media from selected services.
func (s *CleanupService) DeleteMedia(ctx context.Context, media *models.Media, options DeleteOptions) error {
	s.logger.Info("Starting media deletion",
		zap.Uint("id", media.ID),
		zap.String("title", media.Title),
		zap.String("type", media.Type),
		zap.Bool("radarr", options.Radarr),
		zap.Bool("sonarr", options.Sonarr),
		zap.Bool("jellyfin", options.Jellyfin),
		zap.Bool("delete_files", options.DeleteFiles),
		zap.Bool("qbittorrent", options.QBittorrent),
	)

	var errors []error
	deletedFrom := []string{}

	// Delete from Radarr
	if options.Radarr && media.RadarrID != nil && s.radarrClient != nil {
		s.logger.Info("Deleting from Radarr",
			zap.Int("radarr_id", *media.RadarrID),
		)

		err := s.radarrClient.DeleteItem(ctx, *media.RadarrID, options.DeleteFiles)
		if err != nil {
			s.logger.Error("Failed to delete from Radarr",
				zap.Int("radarr_id", *media.RadarrID),
				zap.Error(err),
			)
			errors = append(errors, fmt.Errorf("Radarr: %w", err))
		} else {
			deletedFrom = append(deletedFrom, "Radarr")
			s.logger.Info("Successfully deleted from Radarr")
		}
	}

	// Delete from Sonarr
	if options.Sonarr && media.SonarrID != nil && s.sonarrClient != nil {
		s.logger.Info("Deleting from Sonarr",
			zap.Int("sonarr_id", *media.SonarrID),
		)

		err := s.sonarrClient.DeleteItem(ctx, *media.SonarrID, options.DeleteFiles)
		if err != nil {
			s.logger.Error("Failed to delete from Sonarr",
				zap.Int("sonarr_id", *media.SonarrID),
				zap.Error(err),
			)
			errors = append(errors, fmt.Errorf("Sonarr: %w", err))
		} else {
			deletedFrom = append(deletedFrom, "Sonarr")
			s.logger.Info("Successfully deleted from Sonarr")
		}
	}

	// Delete from Jellyfin
	if options.Jellyfin && media.JellyfinID != nil && s.jellyfinClient != nil {
		s.logger.Info("Deleting from Jellyfin",
			zap.String("jellyfin_id", *media.JellyfinID),
		)

		err := s.jellyfinClient.DeleteItem(ctx, *media.JellyfinID)
		if err != nil {
			s.logger.Error("Failed to delete from Jellyfin",
				zap.String("jellyfin_id", *media.JellyfinID),
				zap.Error(err),
			)
			errors = append(errors, fmt.Errorf("Jellyfin: %w", err))
		} else {
			deletedFrom = append(deletedFrom, "Jellyfin")
			s.logger.Info("Successfully deleted from Jellyfin")
		}
	}

	// Delete files directly if requested and not already deleted by Radarr/Sonarr
	if options.DeleteFiles && !options.Radarr && !options.Sonarr && media.FilePath != "" {
		s.logger.Info("Deleting files from disk",
			zap.String("path", media.FilePath),
		)

		err := s.deleteFilesFromDisk(media.FilePath)
		if err != nil {
			s.logger.Error("Failed to delete files from disk",
				zap.String("path", media.FilePath),
				zap.Error(err),
			)
			errors = append(errors, fmt.Errorf("Disk: %w", err))
		} else {
			deletedFrom = append(deletedFrom, "Disk")
			s.logger.Info("Successfully deleted files from disk")
		}
	}

	// Create history entry
	status := "success"
	message := fmt.Sprintf("Deleted from: %v", deletedFrom)

	if len(errors) > 0 {
		status = "partial"
		message = fmt.Sprintf("Deleted from: %v. Errors: %d", deletedFrom, len(errors))
	}

	historyEntry := &models.History{
		MediaID:    media.ID,
		MediaTitle: media.Title,
		Action:     "deleted",
		Status:     status,
		Message:    message,
		DryRun:     false,
	}

	if err := s.historyRepo.Create(historyEntry); err != nil {
		s.logger.Error("Failed to create history entry",
			zap.Error(err),
		)
	}

	// Return combined error if any deletion failed
	if len(errors) > 0 {
		return fmt.Errorf("deletion completed with errors: %v", errors)
	}

	s.logger.Info("Media deletion completed successfully",
		zap.Uint("id", media.ID),
		zap.String("title", media.Title),
		zap.Strings("deleted_from", deletedFrom),
	)

	return nil
}

// deleteFilesFromDisk deletes files from the filesystem.
func (s *CleanupService) deleteFilesFromDisk(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	// If it's a directory, remove all contents
	if info.IsDir() {
		s.logger.Info("Removing directory recursively",
			zap.String("path", path),
		)
		return os.RemoveAll(path)
	}

	// If it's a file, remove it
	s.logger.Info("Removing file",
		zap.String("path", path),
	)
	return os.Remove(path)
}
