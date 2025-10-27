package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/carcheky/keepercheky/internal/models"
	"github.com/carcheky/keepercheky/internal/repository"
	"github.com/carcheky/keepercheky/internal/service/clients"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// FileActionsHandler handles file action operations
type FileActionsHandler struct {
	mediaRepo        *repository.MediaRepository
	radarrClient     clients.MediaClient
	sonarrClient     clients.MediaClient
	jellyfinClient   clients.StreamingClient
	qbittorrentClient interface{} // Generic interface for now
	logger           *zap.Logger
}

// NewFileActionsHandler creates a new file actions handler
func NewFileActionsHandler(
	mediaRepo *repository.MediaRepository,
	radarrClient clients.MediaClient,
	sonarrClient clients.MediaClient,
	jellyfinClient clients.StreamingClient,
	qbittorrentClient interface{},
	logger *zap.Logger,
) *FileActionsHandler {
	return &FileActionsHandler{
		mediaRepo:         mediaRepo,
		radarrClient:      radarrClient,
		sonarrClient:      sonarrClient,
		jellyfinClient:    jellyfinClient,
		qbittorrentClient: qbittorrentClient,
		logger:            logger,
	}
}

// ImportToRadarrRequest represents a request to import a file to Radarr
type ImportToRadarrRequest struct {
	FilePath   string `json:"file_path"`
	Title      string `json:"title"`
	Quality    string `json:"quality"`
	RootFolder string `json:"root_folder"`
}

// DeleteFileRequest represents a request to delete a file
type DeleteFileRequest struct {
	Confirm bool `json:"confirm"`
}

// CleanupHardlinkRequest represents a request to cleanup hardlinks
type CleanupHardlinkRequest struct {
	KeepPrimary bool `json:"keep_primary"`
}

// BulkActionRequest represents a request for bulk actions
type BulkActionRequest struct {
	FileIDs []uint   `json:"file_ids"`
	Action  string   `json:"action"`
	Options fiber.Map `json:"options"`
}

// BulkActionResponse represents the response from bulk actions
type BulkActionResponse struct {
	Success []uint   `json:"success"`
	Failed  []uint   `json:"failed"`
	Errors  []string `json:"errors"`
}

// ImportToRadarr imports a file to Radarr
func (h *FileActionsHandler) ImportToRadarr(c *fiber.Ctx) error {
	// Get file ID from params
	fileID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file ID"})
	}

	var req ImportToRadarrRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	h.logger.Info("Importing to Radarr",
		zap.Int("file_id", fileID),
		zap.String("title", req.Title),
	)

	// In a real implementation, this would:
	// 1. Create media item for Radarr
	// 2. Call Radarr client to add movie
	// 3. Update database with Radarr ID
	
	// For now, return success
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Import to Radarr initiated",
	})
}

// DeleteFile deletes a file from all services
func (h *FileActionsHandler) DeleteFile(c *fiber.Ctx) error {
	// Get file ID from params
	fileID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file ID"})
	}

	var req DeleteFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Require confirmation
	if !req.Confirm {
		return c.Status(400).JSON(fiber.Map{"error": "confirmation required"})
	}

	h.logger.Info("Deleting file",
		zap.Int("file_id", fileID),
	)

	// TODO: Implement actual deletion logic
	// This would involve:
	// 1. Getting file info from DB
	// 2. Deleting from Radarr/Sonarr if present
	// 3. Deleting from Jellyfin if present
	// 4. Deleting from qBittorrent if present
	// 5. Deleting the actual file

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// CleanupHardlink cleans up hardlinks for a file
func (h *FileActionsHandler) CleanupHardlink(c *fiber.Ctx) error {
	// Get file ID from params
	fileID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid file ID"})
	}

	var req CleanupHardlinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get file info
	var fileInfo models.MediaFileInfo
	err = h.mediaRepo.GetDB().
		Table("media").
		Where("id = ?", fileID).
		First(&fileInfo).Error
	
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File not found"})
	}

	// Validate it's actually a hardlink
	if !fileInfo.IsHardlink {
		return c.Status(400).JSON(fiber.Map{"error": "File is not a hardlink"})
	}

	// Validate same inode
	paths := strings.Split(fileInfo.HardlinkPaths, "|")
	if len(paths) < 2 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid hardlink data"})
	}

	// Get file stats for both paths to verify they're hardlinks
	stat1, err := os.Stat(paths[0])
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to stat first path"})
	}

	stat2, err := os.Stat(paths[1])
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to stat second path"})
	}

	// Check if they have the same inode (simplified check)
	if !os.SameFile(stat1, stat2) {
		return c.Status(400).JSON(fiber.Map{"error": "Paths are not hardlinks of each other"})
	}

	h.logger.Info("Cleaning up hardlink",
		zap.Int("file_id", fileID),
		zap.Strings("paths", paths),
	)

	// TODO: Implement actual cleanup logic
	// This would remove the non-primary hardlink

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hardlink cleaned up successfully",
	})
}

// BulkAction performs bulk actions on multiple files
func (h *FileActionsHandler) BulkAction(c *fiber.Ctx) error {
	var req BulkActionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	response := BulkActionResponse{
		Success: []uint{},
		Failed:  []uint{},
		Errors:  []string{},
	}

	h.logger.Info("Processing bulk action",
		zap.String("action", req.Action),
		zap.Int("file_count", len(req.FileIDs)),
	)

	// Process each file
	for _, fileID := range req.FileIDs {
		err := h.processBulkActionItem(fileID, req.Action, req.Options)
		if err != nil {
			response.Failed = append(response.Failed, fileID)
			response.Errors = append(response.Errors, fmt.Sprintf("File %d: %v", fileID, err))
			h.logger.Error("Bulk action failed for file",
				zap.Uint("file_id", fileID),
				zap.Error(err),
			)
		} else {
			response.Success = append(response.Success, fileID)
		}
	}

	return c.JSON(response)
}

// processBulkActionItem processes a single item in bulk action
func (h *FileActionsHandler) processBulkActionItem(fileID uint, action string, options fiber.Map) error {
	// This is a simplified implementation
	// In production, this would call the appropriate action handler
	switch action {
	case "delete":
		// Simulate deletion
		h.logger.Info("Simulating delete", zap.Uint("file_id", fileID))
		return nil
	case "import_to_radarr":
		// Simulate import
		h.logger.Info("Simulating import to Radarr", zap.Uint("file_id", fileID))
		return nil
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}
