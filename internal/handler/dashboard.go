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
	})
}

func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	// TODO: Implement stats calculation
	return c.JSON(fiber.Map{
		"total_media":  0,
		"total_size":   0,
		"to_delete":    0,
		"leaving_soon": 0,
	})
}
