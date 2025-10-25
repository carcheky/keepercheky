package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	repos       *repository.Repositories
	logger      *logger.Logger
	config      *config.Config
	syncService *service.SyncService
}

func NewSettingsHandler(repos *repository.Repositories, logger *logger.Logger, cfg *config.Config, syncService *service.SyncService) *SettingsHandler {
	return &SettingsHandler{
		repos:       repos,
		logger:      logger,
		config:      cfg,
		syncService: syncService,
	}
}

func (h *SettingsHandler) Index(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{
		"Title": "Settings - KeeperCheky",
	})
}

func (h *SettingsHandler) Get(c *fiber.Ctx) error {
	return c.JSON(h.config)
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	// TODO: Implement settings update
	return c.JSON(fiber.Map{
		"message": "Settings updated successfully",
	})
}

func (h *SettingsHandler) TestConnection(c *fiber.Ctx) error {
	service := c.Params("service")
	
	h.logger.Info("Testing connection", "service", service)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	err := h.syncService.TestConnection(ctx, service)
	if err != nil {
		h.logger.Error("Connection test failed",
			"service", service,
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"service": service,
			"message": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"status":  "success",
		"service": service,
		"message": "Connection test successful",
	})
}
