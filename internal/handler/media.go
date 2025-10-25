package handler

import (
	"strconv"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	repos  *repository.Repositories
	logger *logger.Logger
}

func NewMediaHandler(repos *repository.Repositories, logger *logger.Logger) *MediaHandler {
	return &MediaHandler{
		repos:  repos,
		logger: logger,
	}
}

func (h *MediaHandler) List(c *fiber.Ctx) error {
	return c.Render("pages/media", fiber.Map{
		"Title": "Media Management - KeeperCheky",
	}, "layouts/main")
}

func (h *MediaHandler) GetAll(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))
	mediaType := c.Query("type", "all")
	search := c.Query("search", "")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var media []models.Media
	var total int64
	var err error

	// If search query provided, use search
	if search != "" {
		media, err = h.repos.Media.Search(search)
		total = int64(len(media))
	} else {
		media, total, err = h.repos.Media.GetPaginated(page, pageSize, mediaType)
	}

	if err != nil {
		h.logger.Error("Failed to get media", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve media",
		})
	}

	return c.JSON(fiber.Map{
		"data":     media,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *MediaHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid media ID",
		})
	}

	media, err := h.repos.Media.GetByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get media", "id", id, "error", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	return c.JSON(media)
}

func (h *MediaHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid media ID",
		})
	}

	// TODO: Implement actual deletion logic with cleanup service
	if err := h.repos.Media.Delete(uint(id)); err != nil {
		h.logger.Error("Failed to delete media", "id", id, "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete media",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Media deleted successfully",
	})
}

func (h *MediaHandler) Exclude(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid media ID",
		})
	}

	// Check if media exists
	media, err := h.repos.Media.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	// Toggle excluded status
	newExcludedStatus := !media.Excluded
	if err := h.repos.Media.SetExcluded(uint(id), newExcludedStatus); err != nil {
		h.logger.Error("Failed to toggle exclude status", "id", id, "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update media",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Media status updated successfully",
		"excluded": newExcludedStatus,
	})
}

func (h *MediaHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.repos.Media.GetStats()
	if err != nil {
		h.logger.Error("Failed to get stats", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve statistics",
		})
	}

	return c.JSON(stats)
}
