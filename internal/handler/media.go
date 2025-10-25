package handler

import (
	"strconv"

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
	})
}

func (h *MediaHandler) GetAll(c *fiber.Ctx) error {
	media, err := h.repos.Media.GetAll()
	if err != nil {
		h.logger.Error("Failed to get media", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve media",
		})
	}

	return c.JSON(media)
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

	media, err := h.repos.Media.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	media.Excluded = true
	if err := h.repos.Media.Update(media); err != nil {
		h.logger.Error("Failed to exclude media", "id", id, "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to exclude media",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Media excluded successfully",
	})
}
