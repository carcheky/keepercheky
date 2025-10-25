package handler

import (
	"github.com/carcheky/keepercheky/internal/config"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	repos  *repository.Repositories
	logger *logger.Logger
	config *config.Config
}

func NewSettingsHandler(repos *repository.Repositories, logger *logger.Logger, cfg *config.Config) *SettingsHandler {
	return &SettingsHandler{
		repos:  repos,
		logger: logger,
		config: cfg,
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
	
	// TODO: Implement connection testing
	h.logger.Info("Testing connection", "service", service)
	
	return c.JSON(fiber.Map{
		"status":  "success",
		"service": service,
		"message": "Connection test successful",
	})
}
