package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	repos       *repository.Repositories
	logger      *logger.Logger
	syncService *service.SyncService
}

func NewDashboardHandler(repos *repository.Repositories, logger *logger.Logger, syncService *service.SyncService) *DashboardHandler {
	return &DashboardHandler{
		repos:       repos,
		logger:      logger,
		syncService: syncService,
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

// GetJellyseerrStats returns detailed Jellyseerr statistics.
func (h *DashboardHandler) GetJellyseerrStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestStats, err := h.syncService.GetJellyseerrRequestStats(ctx)
	if err != nil {
		h.logger.Error("Failed to get Jellyseerr stats",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(requestStats)
}

// GetJellyseerrRequests returns recent Jellyseerr requests.
func (h *DashboardHandler) GetJellyseerrRequests(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	requests, err := h.syncService.GetJellyseerrRequests(ctx)
	if err != nil {
		h.logger.Error("Failed to get Jellyseerr requests",
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"requests": requests,
		"count":    len(requests),
	})
}
