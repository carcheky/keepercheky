package handler

import (
	"context"
	"time"

	"github.com/carcheky/keepercheky/internal/service"
	"github.com/carcheky/keepercheky/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type SyncHandler struct {
	syncService *service.SyncService
	logger      *logger.Logger
}

func NewSyncHandler(syncService *service.SyncService, logger *logger.Logger) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
		logger:      logger,
	}
}

// Sync triggers a full synchronization from all services.
func (h *SyncHandler) Sync(c *fiber.Ctx) error {
	h.logger.Info("Manual sync triggered")
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	err := h.syncService.SyncAll(ctx)
	if err != nil {
		h.logger.Error("Sync failed", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Synchronization completed successfully",
	})
}
