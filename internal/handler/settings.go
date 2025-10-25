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
	// Build exclusion tags string
	exclusionTags := ""
	if len(h.config.Cleanup.ExclusionTags) > 0 {
		exclusionTags = string(h.config.Cleanup.ExclusionTags[0])
		for i := 1; i < len(h.config.Cleanup.ExclusionTags); i++ {
			exclusionTags += ", " + h.config.Cleanup.ExclusionTags[i]
		}
	}

	// Get environment source map to inform UI which fields are from .env
	envSources := config.GetEnvSourceMap()

	// Return config in the format expected by frontend
	// Note: Config already contains values from environment variables (they have precedence in Viper)
	return c.JSON(fiber.Map{
		"services": fiber.Map{
			"radarr": fiber.Map{
				"enabled": h.config.Clients.Radarr.Enabled,
				"url":     h.config.Clients.Radarr.URL,
				"api_key": h.config.Clients.Radarr.APIKey,
			},
			"sonarr": fiber.Map{
				"enabled": h.config.Clients.Sonarr.Enabled,
				"url":     h.config.Clients.Sonarr.URL,
				"api_key": h.config.Clients.Sonarr.APIKey,
			},
			"jellyfin": fiber.Map{
				"enabled": h.config.Clients.Jellyfin.Enabled,
				"url":     h.config.Clients.Jellyfin.URL,
				"api_key": h.config.Clients.Jellyfin.APIKey,
			},
			"jellyseerr": fiber.Map{
				"enabled": h.config.Clients.Jellyseerr.Enabled,
				"url":     h.config.Clients.Jellyseerr.URL,
				"api_key": h.config.Clients.Jellyseerr.APIKey,
			},
			"qbittorrent": fiber.Map{
				"enabled":  h.config.Clients.QBittorrent.Enabled,
				"url":      h.config.Clients.QBittorrent.URL,
				"username": h.config.Clients.QBittorrent.Username,
				"password": h.config.Clients.QBittorrent.Password,
			},
		},
		"cleanup": fiber.Map{
			"dry_run":            h.config.Cleanup.DryRun,
			"days_to_keep":       h.config.Cleanup.DaysToKeep,
			"leaving_soon_days":  h.config.Cleanup.LeavingSoonDays,
			"exclusion_tags":     exclusionTags,
			"delete_unmonitored": h.config.Cleanup.DeleteUnmonitored,
		},
		"env_sources": envSources,
	})
}

func (h *SettingsHandler) Update(c *fiber.Ctx) error {
	type ConfigUpdate struct {
		Services struct {
			Radarr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"radarr"`
			Sonarr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"sonarr"`
			Jellyfin struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"jellyfin"`
			Jellyseerr struct {
				Enabled bool   `json:"enabled"`
				URL     string `json:"url"`
				APIKey  string `json:"api_key"`
			} `json:"jellyseerr"`
			QBittorrent struct {
				Enabled  bool   `json:"enabled"`
				URL      string `json:"url"`
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"qbittorrent"`
		} `json:"services"`
		Cleanup struct {
			DryRun            bool `json:"dry_run"`
			DaysToKeep        int  `json:"days_to_keep"`
			LeavingSoonDays   int  `json:"leaving_soon_days"`
			DeleteUnmonitored bool `json:"delete_unmonitored"`
		} `json:"cleanup"`
	}

	var update ConfigUpdate
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Update config in memory (environment variables will override these values on next load)
	h.config.Clients.Radarr.Enabled = update.Services.Radarr.Enabled
	h.config.Clients.Radarr.URL = update.Services.Radarr.URL
	h.config.Clients.Radarr.APIKey = update.Services.Radarr.APIKey

	h.config.Clients.Sonarr.Enabled = update.Services.Sonarr.Enabled
	h.config.Clients.Sonarr.URL = update.Services.Sonarr.URL
	h.config.Clients.Sonarr.APIKey = update.Services.Sonarr.APIKey

	h.config.Clients.Jellyfin.Enabled = update.Services.Jellyfin.Enabled
	h.config.Clients.Jellyfin.URL = update.Services.Jellyfin.URL
	h.config.Clients.Jellyfin.APIKey = update.Services.Jellyfin.APIKey

	h.config.Clients.Jellyseerr.Enabled = update.Services.Jellyseerr.Enabled
	h.config.Clients.Jellyseerr.URL = update.Services.Jellyseerr.URL
	h.config.Clients.Jellyseerr.APIKey = update.Services.Jellyseerr.APIKey

	h.config.Clients.QBittorrent.Enabled = update.Services.QBittorrent.Enabled
	h.config.Clients.QBittorrent.URL = update.Services.QBittorrent.URL
	h.config.Clients.QBittorrent.Username = update.Services.QBittorrent.Username
	h.config.Clients.QBittorrent.Password = update.Services.QBittorrent.Password

	h.config.Cleanup.DryRun = update.Cleanup.DryRun
	h.config.Cleanup.DaysToKeep = update.Cleanup.DaysToKeep
	h.config.Cleanup.LeavingSoonDays = update.Cleanup.LeavingSoonDays
	h.config.Cleanup.DeleteUnmonitored = update.Cleanup.DeleteUnmonitored

	// Parse exclusion tags
	// TODO: Parse exclusion_tags from string to []string

	h.logger.Info("Configuration updated",
		"radarr_enabled", h.config.Clients.Radarr.Enabled,
		"sonarr_enabled", h.config.Clients.Sonarr.Enabled,
		"jellyfin_enabled", h.config.Clients.Jellyfin.Enabled,
		"dry_run", h.config.Cleanup.DryRun,
	)

	// Save configuration to file
	if err := config.Save(h.config); err != nil {
		h.logger.Error("Failed to save configuration to file", "error", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Configuration updated in memory but failed to save to file: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Configuration saved successfully",
		"success": true,
	})
}

func (h *SettingsHandler) TestConnection(c *fiber.Ctx) error {
	service := c.Params("service")

	h.logger.Info("Testing connection", "service", service)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test basic connection first
	err := h.syncService.TestConnection(ctx, service)
	if err != nil {
		h.logger.Error("Connection test failed",
			"service", service,
			"error", err,
		)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// For Radarr, get complete system information
	if service == "radarr" {
		systemInfo, err := h.syncService.GetRadarrSystemInfo(ctx)
		if err != nil {
			h.logger.Error("Failed to get Radarr system info",
				"error", err,
			)
			// Still return success for connection, but without detailed info
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Connection test successful",
			})
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Connection successful",
			"system_info": systemInfo,
		})
	}

	// For Sonarr, get complete system information
	if service == "sonarr" {
		systemInfo, err := h.syncService.GetSonarrSystemInfo(ctx)
		if err != nil {
			h.logger.Error("Failed to get Sonarr system info",
				"error", err,
			)
			// Still return success for connection, but without detailed info
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Connection test successful",
			})
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Connection successful",
			"system_info": systemInfo,
		})
	}

	// For Jellyfin, get complete system information
	if service == "jellyfin" {
		systemInfo, err := h.syncService.GetJellyfinSystemInfo(ctx)
		if err != nil {
			h.logger.Error("Failed to get Jellyfin system info",
				"error", err,
			)
			// Still return success for connection, but without detailed info
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Connection test successful",
			})
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Connection successful",
			"system_info": systemInfo,
		})
	}

	// For Jellyseerr, get complete system information
	if service == "jellyseerr" {
		systemInfo, err := h.syncService.GetJellyseerrSystemInfo(ctx)
		if err != nil {
			h.logger.Error("Failed to get Jellyseerr system info",
				"error", err,
			)
			// Still return success for connection, but without detailed info
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Connection test successful",
			})
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Connection successful",
			"system_info": systemInfo,
		})
	}

	// For qBittorrent, get complete system information
	if service == "qbittorrent" {
		systemInfo, err := h.syncService.GetQBittorrentSystemInfo(ctx)
		if err != nil {
			h.logger.Error("Failed to get qBittorrent system info",
				"error", err,
			)
			// Still return success for connection, but without detailed info
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Connection test successful",
			})
		}

		return c.JSON(fiber.Map{
			"success":     true,
			"message":     "Connection successful",
			"system_info": systemInfo,
		})
	}

	// For other services, return basic success
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection test successful",
	})
}
