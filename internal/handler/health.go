package handler

import (
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	repos  *repository.Repositories
	logger *logger.Logger
}

func NewHealthHandler(repos *repository.Repositories, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		repos:  repos,
		logger: logger,
	}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	// Get database connection
	db := h.repos.Media.(*repository.MediaRepository)
	sqlDB, err := db.GetDB().DB()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  "database ping failed",
		})
	}

	return c.JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": fiber.Now(),
	})
}
