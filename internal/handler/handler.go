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
}

func NewHandlers(db *gorm.DB, repos *repository.Repositories, logger *logger.Logger, cfg *config.Config) *Handlers {
	// Initialize SyncService
	syncService := service.NewSyncService(repos.Media, logger, cfg)

	// Initialize CleanupService
	cleanupService := cleanup.NewCleanupService(
		repos.Media,
		repos.History,
		syncService.GetRadarrClient(),
		syncService.GetSonarrClient(),
		syncService.GetJellyfinClient(),
		nil, // qBittorrent client not yet implemented
		logger.Desugar(),
	)

	return &Handlers{
		Health:    NewHealthHandler(db, logger),
		Dashboard: NewDashboardHandler(repos, logger),
		Media:     NewMediaHandler(repos, cleanupService, logger),
		Schedule:  NewScheduleHandler(repos, logger),
		Settings:  NewSettingsHandler(repos, logger, cfg, syncService),
		Logs:      NewLogsHandler(repos, logger),
		Sync:      NewSyncHandler(syncService, logger),
	}
}
