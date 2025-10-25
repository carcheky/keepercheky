package handler

import (
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	repos  *repository.Repositories
	logger *logger.Logger
}

func NewDashboardHandler(repos *repository.Repositories, logger *logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		repos:  repos,
		logger: logger,
	}
}

func (h *DashboardHandler) Index(c *fiber.Ctx) error {
	return c.Render("pages/dashboard", fiber.Map{
		"Title": "Dashboard - KeeperCheky",
	}, "layouts/main")
}

func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	stats, err := h.repos.Media.GetStats()
	if err != nil {
		h.logger.Error("Failed to get stats", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get statistics",
		})
	}

	// Add placeholder values for features not yet implemented
	stats["to_delete"] = 0
	stats["leaving_soon"] = 0

	return c.JSON(stats)
}
