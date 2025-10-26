package handler

import (
	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/internal/service/cleanup"
	"github.com/carcheky/keepercheky/pkg/logger"
	"gorm.io/gorm"
)

type Handlers struct {
	Health    *HealthHandler
	Dashboard *DashboardHandler
	Media     *MediaHandler
	Schedule  *ScheduleHandler
	Settings  *SettingsHandler
	Logs      *LogsHandler
	Sync      *SyncHandler
	Files     *FilesHandler
}

func NewHandlers(db *gorm.DB, repos *repository.Repositories, logger *logger.Logger, cfg *config.Config) *Handlers {
	// Initialize OLD SyncService (for client access and backward compatibility)
	oldSyncService := service.NewSyncService(repos.Media, logger, cfg)

	// Initialize NEW FilesystemSyncService (filesystem-first approach)
	filesystemSyncService := service.NewFilesystemSyncService(
		repos.Media,
		repos.Settings,
		oldSyncService.GetRadarrClient(),
		oldSyncService.GetSonarrClient(),
		oldSyncService.GetJellyfinClient(),
		nil, // Jellyseerr not needed for filesystem sync
		oldSyncService.GetQBittorrentClient(),
		logger.Desugar(),
		cfg,
	)

	// Initialize CleanupService
	cleanupService := cleanup.NewCleanupService(
		repos.Media,
		repos.History,
		oldSyncService.GetRadarrClient(),
		oldSyncService.GetSonarrClient(),
		oldSyncService.GetJellyfinClient(),
		oldSyncService.GetQBittorrentClient(),
		logger.Desugar(),
	)

	return &Handlers{
		Health:    NewHealthHandler(db, logger),
		Dashboard: NewDashboardHandler(repos, logger),
		Media:     NewMediaHandler(repos, cleanupService, logger),
		Schedule:  NewScheduleHandler(repos, logger),
		Settings:  NewSettingsHandler(repos, logger, cfg, oldSyncService),
		Logs:      NewLogsHandler(repos, logger),
		Sync:      NewSyncHandler(filesystemSyncService, logger), // Use NEW filesystem-first sync
		Files:     NewFilesHandler(repos.Media, cfg, oldSyncService, logger.Desugar()),
	}
}
