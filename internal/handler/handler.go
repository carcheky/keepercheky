package handler

import (
	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
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
}

func NewHandlers(db *gorm.DB, repos *repository.Repositories, logger *logger.Logger, cfg *config.Config) *Handlers {
	return &Handlers{
		Health:    NewHealthHandler(db, logger),
		Dashboard: NewDashboardHandler(repos, logger),
		Media:     NewMediaHandler(repos, logger),
		Schedule:  NewScheduleHandler(repos, logger),
		Settings:  NewSettingsHandler(repos, logger, cfg),
		Logs:      NewLogsHandler(repos, logger),
	}
}
