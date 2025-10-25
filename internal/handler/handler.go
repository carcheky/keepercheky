package handler

import (
	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
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
	syncService := service.NewSyncService(repos.Media, logger.Logger, cfg)
	
	return &Handlers{
		Health:    NewHealthHandler(db, logger),
		Dashboard: NewDashboardHandler(repos, logger),
		Media:     NewMediaHandler(repos, logger),
		Schedule:  NewScheduleHandler(repos, logger),
		Settings:  NewSettingsHandler(repos, logger, cfg, syncService),
		Logs:      NewLogsHandler(repos, logger),
		Sync:      NewSyncHandler(syncService, logger),
	}
}
